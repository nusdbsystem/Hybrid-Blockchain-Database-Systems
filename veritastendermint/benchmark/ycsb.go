package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/veritas/benchmark"
)

var (
	urls              = kingpin.Flag("urls", "URLs of tendermint nodes, split by ,").Required().String()
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
)

func main() {
	kingpin.Parse()

	addrs := strings.Split(*urls, ",")
	reqNum := atomic.NewFloat64(0)

	loadFile, err := os.Open(*dataLoad)
	if err != nil {
		panic(err)
	}
	defer loadFile.Close()
	loadBuf := make(chan [2]string, 256)
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
		avaLatency = float64(all)/reqNum.Load()
	}()
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			for kv := range loadBuf {
				h1, h2 := md5.New(), md5.New()
				h1.Write([]byte(kv[0]))
				h2.Write([]byte(kv[1]))
				resp, err := http.Get(fmt.Sprintf("%s/broadcast_tx_commit?tx=\"%s\"", addrs[seq%len(addrs)], hex.EncodeToString(h1.Sum(nil))+"="+hex.EncodeToString(h2.Sum(nil))))
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
				if resp.StatusCode != 200 {
					panic(fmt.Sprintf("send request error: %v", resp.StatusCode))
				}
			}
		}(i)
	}
	wg.Wait()

	runFile, err := os.Open(*dataRun)
	if err != nil {
		panic(err)
	}
	defer runFile.Close()
	runBuf := make(chan *benchmark.Request, 100*(*driverConcurrency))
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

	for i := 0; i < *driverConcurrency; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			for op := range runBuf {
				switch op.ReqType {
				case benchmark.GetOp:
					start := time.Now()
					h1 := md5.New()
					h1.Write([]byte(op.Key))
					resp, err := http.Get(fmt.Sprintf("%s/abci_query?data=\"%s\"", addrs[seq%len(addrs)], hex.EncodeToString(h1.Sum(nil))))
					if err != nil {
						panic(err)
					}
					resp.Body.Close()
					if resp.StatusCode != 200 {
						panic(fmt.Sprintf("send request error: %v", resp.StatusCode))
					}
					latencyCh <- time.Since(start)
				case benchmark.SetOp:
					start := time.Now()
					h1, h2 := md5.New(), md5.New()
					h1.Write([]byte(op.Key))
					h2.Write([]byte(op.Val))
					resp, err := http.Get(fmt.Sprintf("%s/broadcast_tx_commit?tx=\"%s\"", addrs[seq%len(addrs)], hex.EncodeToString(h1.Sum(nil))+"="+hex.EncodeToString(h2.Sum(nil))))
					if err != nil {
						panic(err)
					}
					resp.Body.Close()
					if resp.StatusCode != 200 {
						panic(fmt.Sprintf("send request error: %v", resp.StatusCode))
					}
					latencyCh <- time.Since(start)
				default:
					panic(fmt.Sprintf("invalid operation: %v", op.ReqType))
				}
			}
		}(i)
	}
	start := time.Now()
	wg.Wait()
	close(latencyCh)
	wg2.Wait()

	fmt.Printf("Throughput: %v req/s\n", int64(reqNum.Load()/time.Since(start).Seconds()))
	fmt.Printf("Average latency: %v ms\n", avaLatency)
}
