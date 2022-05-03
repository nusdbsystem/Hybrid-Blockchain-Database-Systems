package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	benchmark "github.com/nusdbsystem/hybrid/veritas_hotstuff/benchmark/ycsb"
	pbv "github.com/nusdbsystem/hybrid/veritas_hotstuff/proto/veritashs"
)

var (
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverNum         = kingpin.Flag("ndrivers", "Number of drivers for sending requests").Default("4").Int()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
	serverAddrs       = kingpin.Flag("server-addrs", "Address of veritas nodes").Required().String()
)

func main() {
	kingpin.Parse()
	// defer profile.Start(profile.CPUProfile, profile.NoShutdownHook, profile.ProfilePath("./tmp")).Stop()

	fmt.Println("Time start: ", time.Now())
	lastopt := ""
	lastkey := ""
	addrs := strings.Split(*serverAddrs, ",")
	clis := make([]pbv.VeritasNodeClient, 0)
	conns := make([]*grpc.ClientConn, 0)
	defer func() {
		for _, conn := range conns {
			conn.Close()
		}
	}()

	for i := 0; i < *driverNum; i++ {
		conn, err := grpc.Dial(addrs[i%len(addrs)], grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		cli := pbv.NewVeritasNodeClient(conn)

		conns = append(conns, conn)
		clis = append(clis, cli)
	}

	reqNum := atomic.NewInt64(0)

	loadFile, err := os.Open(*dataLoad)
	if err != nil {
		panic(err)
	}
	defer loadFile.Close()

	loadBuf := make(chan [3]string, 1024)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(loadBuf)
		if err := benchmark.LineByLine(loadFile, func(line string) error {
			operands := strings.SplitN(line, " ", 5)
			loadBuf <- [3]string{operands[2], operands[3], operands[4]}
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	latencyCh := make(chan time.Duration, 1024)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	var avaLatency float64
	go func() {
		defer wg2.Done()
		all := int64(0)
		for ts := range latencyCh {
			all += ts.Milliseconds()
		}
		avaLatency = float64(all) / float64(reqNum.Load())
	}()
	// Init Data
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range loadBuf {
				lastkey = kv[0]
				ver, err := strconv.ParseInt(kv[2], 10, 64)
				if err != nil {
					panic(err)
				}
				if _, err := clis[0].Set(context.Background(), &pbv.SetRequest{
					Key:     kv[0],
					Value:   kv[1],
					Version: ver,
				}); err != nil {
					panic(err)
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println("Init Data done " + lastkey)

	runFile, err := os.Open(*dataRun)
	if err != nil {
		panic(err)
	}
	defer runFile.Close()

	runBuf := make(chan *benchmark.Request, 20*(*driverNum)*(*driverConcurrency))
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(runBuf)
		if err := benchmark.LineByLine(runFile, func(line string) error {
			reqNum.Add(1)
			operands := strings.SplitN(line, " ", 5)
			ver, err := strconv.ParseInt(operands[4], 10, 64)
			if err != nil {
				panic(err)
			}
			r := &benchmark.Request{
				Key: operands[2],
			}
			if operands[0] == "READ" {
				r.ReqType = benchmark.GetOp
			} else {
				r.ReqType = benchmark.SetOp
				r.Val = operands[3]
				r.Version = ver
			}
			runBuf <- r
			return nil
		}); err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)

	start := time.Now()
	for i := 0; i < *driverNum; i++ {
		for j := 0; j < *driverConcurrency; j++ {
			wg.Add(1)
			go func(seq int) {
				defer wg.Done()
				for op := range runBuf {
					switch op.ReqType {
					case benchmark.GetOp:
						beginOp := time.Now()
						if _, err := clis[seq].Get(context.Background(), &pbv.GetRequest{Key: op.Key}); err != nil {
							// panic(err)
							fmt.Println("GetOp ", err)
						}
						latencyCh <- time.Since(beginOp)
					case benchmark.SetOp:
						beginOp := time.Now()
						_, err := clis[seq].Set(context.Background(), &pbv.SetRequest{
							Key:     op.Key,
							Value:   op.Val,
							Version: op.Version,
						})
						if err != nil {
							fmt.Println("SetOp ", err)
						}
						latencyCh <- time.Since(beginOp)
						lastopt = "set"
						lastkey = op.Key
					default:
						panic(fmt.Sprintf("invalid operation: %v", op.ReqType))
					}
				}
			}(i)
		}
	}
	fmt.Println("get/set opt is ongoing ... ", lastkey)
	wg.Wait()
	fmt.Println("wg.Wait() ... ", lastkey)
	close(latencyCh)
	fmt.Println("close(latencyCh) ... ", lastkey)
	wg2.Wait()
	fmt.Println("wg2.Wait() ... ", lastkey)

	fmt.Println("Last opt verify is ongoing ... ", lastopt)
	fmt.Println("Last key verify is ongoing ... ", lastkey)

	fmt.Println("#########################################################################")
	fmt.Printf("Experiment: %v servers %v drivers with %v concurrency to handle %v requests(loadpath %v workload %v ) --> Throughput: %v req/s, ",
		len(addrs), *driverNum, *driverConcurrency, reqNum.Load(), *dataRun, *dataLoad,
		int64(float64(reqNum.Load())/time.Since(start).Seconds()),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
	fmt.Println("#########################################################################")
	fmt.Println("Time stop: ", time.Now())
}
