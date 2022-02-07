package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	benchmark "hybrid/VeritasHotstuff/benchmark/ycsb"
	pbv "hybrid/VeritasHotstuff/proto/veritas"
)

var (
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverNum         = kingpin.Flag("ndrivers", "Number of drivers for sending requests").Default("4").Int()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
	serverAddrs       = kingpin.Flag("veritas-addrs", "Address of veritas nodes").Required().String()
)

func main() {
	kingpin.Parse()

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

	loadBuf := make(chan [2]string, 1024)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(loadBuf)
		if err := benchmark.LineByLine(loadFile, func(line string) error {
			operands := strings.SplitN(line, " ", 4)
			loadBuf <- [2]string{operands[2], operands[3]}
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
				if _, err := clis[0].Set(context.Background(), &pbv.SetRequest{
					Key:   kv[0],
					Value: kv[1],
				}); err != nil {
					panic(err)
				}
			}
		}()
	}
	wg.Wait()

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
			operands := strings.SplitN(line, " ", 4)
			r := &benchmark.Request{
				Key: operands[2],
			}
			if operands[0] == "READ" {
				r.ReqType = benchmark.GetOp
			} else {
				r.ReqType = benchmark.SetOp
				r.Val = operands[3]
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
							fmt.Println(err)
						}
						latencyCh <- time.Since(beginOp)
					case benchmark.SetOp:
						beginOp := time.Now()
						if _, err := clis[seq].Set(context.Background(), &pbv.SetRequest{
							Key:   op.Key,
							Value: op.Val,
						}); err != nil {
							panic(err)
						}
						latencyCh <- time.Since(beginOp)
					default:
						panic(fmt.Sprintf("invalid operation: %v", op.ReqType))
					}
				}
			}(i)
		}
	}
	wg.Wait()
	close(latencyCh)
	wg2.Wait()

	fmt.Println("#########################################################################")
	fmt.Printf("Throughput of %v drivers with %v concurrency to handle %v requests: %v req/s\n",
		*driverNum, *driverConcurrency, reqNum.Load(),
		int64(float64(reqNum.Load())/time.Since(start).Seconds()),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
	fmt.Println("#########################################################################")
}
