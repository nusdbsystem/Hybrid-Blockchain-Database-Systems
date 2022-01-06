package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/veritas/benchmark"
	"hybrid/veritas/driver"
)

var (
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverNum         = kingpin.Flag("ndrivers", "Number of drivers for sending requests").Default("4").Int()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
	veritasAddrs      = kingpin.Flag("veritas-addrs", "Address of veritas nodes").Required().String()
	tsoAddr           = kingpin.Flag("tso-addr", "TSO server address").Required().String()
)

func encodeVal(val string) string {
        runes := []rune(val)
        for i := 0; i < len(runes); i++ {
                if (runes[i] >= 'a' && runes[i] <= 'z') ||
                        (runes[i] >= 'A' && runes[i] <= 'Z') ||
                        (runes[i] >= '0' && runes[i] <= '9') {
                        continue
                }
                runes[i] = '0'
        }
        return string(runes)
}

func main() {
	kingpin.Parse()

	addrs := strings.Split(*veritasAddrs, ",")
	clis := make([]*driver.Driver, 0)
	defer func() {
		for _, cli := range clis {
			cli.Close()
		}
	}()
	for i := 0; i < *driverNum; i++ {
		cli, err := driver.Open(addrs[i%len(addrs)], *tsoAddr, benchmark.GenRandString(16))
		if err != nil {
			panic(err)
		}
		clis = append(clis, cli)
	}

	reqNum := atomic.NewInt64(0)

	loadFile, err := os.Open(*dataLoad)
	if err != nil {
		panic(err)
	}
	defer loadFile.Close()
	loadBuf := make(chan [2]string, 10)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(loadBuf)
		if err := benchmark.LineByLine(loadFile, func(line string) error {
			operands := strings.SplitN(line, " ", 4)
			loadBuf <- [2]string{operands[2], encodeVal(operands[3])}
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
			all += ts.Microseconds()
		}
		avaLatency = float64(all) / (1000 * float64(reqNum.Load()))
	}()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range loadBuf {
				clis[0].Set(context.Background(), kv[0], kv[1])
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
				r.Val = encodeVal(operands[3])
			}
			runBuf <- r
			return nil
		}); err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)

	for i := 0; i < *driverNum; i++ {
		for j := 0; j < *driverConcurrency; j++ {
			wg.Add(1)
			go func(seq int) {
				defer wg.Done()
				for op := range runBuf {
					switch op.ReqType {
					case benchmark.GetOp:
						start := time.Now()
						clis[seq].Get(context.Background(), op.Key)
						latencyCh <- time.Since(start)
					case benchmark.SetOp:
						start := time.Now()
						clis[seq].Set(context.Background(), op.Key, op.Val)
						latencyCh <- time.Since(start)
					default:
						panic(fmt.Sprintf("invalid operation: %v", op.ReqType))
					}
				}
			}(i)
		}
	}
	start := time.Now()
	wg.Wait()
	close(latencyCh)
	wg2.Wait()
	fmt.Printf("Throughput of %v drivers with %v concurrency to handle %v requests: %v req/s\n",
		*driverNum, *driverConcurrency, reqNum,
		int64(float64(reqNum.Load())/time.Since(start).Seconds()),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
}
