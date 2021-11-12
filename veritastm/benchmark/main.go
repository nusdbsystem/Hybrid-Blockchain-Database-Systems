package main

import (
	"context"
	"fmt"

	"hybrid/veritastm"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverNum         = kingpin.Flag("ndrivers", "Number of drivers for sending requests").Default("4").Int()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
	veritasAddrs      = kingpin.Flag("veritas-addrs", "Address of Veritas nodes").Required().String()
)

func main() {
	kingpin.Parse()

	addrs := strings.Split(*veritasAddrs, ",")
	clis := make([]*veritastm.Driver, 0)
	defer func() {
		for _, cli := range clis {
			cli.Close()
		}
	}()
	for i := 0; i < *driverNum; i++ {
		cli, err := veritastm.Open(addrs[i%len(addrs)], veritastm.GenRandString(16))
		if err != nil {
			panic(err)
		}
		clis = append(clis, cli)
	}

	fmt.Println("Start loading ...")
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
		if err := veritastm.LineByLine(loadFile, func(line string) error {
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
			all += ts.Microseconds()
		}
		avaLatency = float64(all) / (1000 * float64(reqNum.Load()))
	}()
	for i := 0; i < *driverConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range loadBuf {
				clis[0].Set(context.Background(), kv[0], kv[1])
			}
		}()
	}
	wg.Wait()
	fmt.Println("End loading ...")

	fmt.Println("Start running ...")
	runFile, err := os.Open(*dataRun)
	if err != nil {
		panic(err)
	}
	defer runFile.Close()
	runBuf := make(chan *veritastm.Request, 20*(*driverNum)*(*driverConcurrency))
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(runBuf)
		if err := veritastm.LineByLine(runFile, func(line string) error {
			reqNum.Add(1)
			operands := strings.SplitN(line, " ", 4)
			r := &veritastm.Request{
				Key: operands[2],
			}
			if operands[0] == "READ" {
				r.ReqType = veritastm.GetOp
			} else {
				r.ReqType = veritastm.SetOp
				r.Val = operands[3]
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
					case veritastm.GetOp:
						start := time.Now()
						clis[seq].Get(context.Background(), op.Key)
						latencyCh <- time.Since(start)
					case veritastm.SetOp:
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
