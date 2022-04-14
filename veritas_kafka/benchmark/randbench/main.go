package main

import (
	"context"
	"fmt"
	"hybrid/veritas/benchmark"
	"math/rand"
	"strings"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/veritas/driver"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	concurrency  = kingpin.Flag("concurrency", "Concurrency of sending requests").Default("2").Int()
	reqNum       = kingpin.Flag("req-num", "Request number").Default("10000").Int()
	veritasAddrs = kingpin.Flag("veritas-addrs", "Addresses of veritas nodes, split by ,").Required().String()
	signatures   = kingpin.Flag("signatures", "Signatures of clients, split by ,").Required().String()
	tsoAddr      = kingpin.Flag("tso-addr", "TSO server address").Required().String()
	method       = kingpin.Flag("method", "set/get").Required().Enum("set", "get")
	keySize      = kingpin.Flag("key-size", "key size").Default("128").Int()
	valSize      = kingpin.Flag("val-size", "val size").Default("128").Int()
)

func main() {
	kingpin.Parse()

	addrs := strings.Split(*veritasAddrs, ",")
	sigs := strings.Split(*signatures, ",")

	clis := make([]*driver.Driver, 0)
	defer func() {
		for _, cli := range clis {
			cli.Close()
		}
	}()
	for i, addr := range addrs {
		cli, err := driver.Open(addr, *tsoAddr, sigs[i])
		if err != nil {
			panic(err)
		}
		clis = append(clis, cli)
	}

	// Pre-Condition
	runBuf := make(chan *benchmark.Request, 20*(*concurrency))
	latencyCh := make(chan time.Duration, 1024)
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	var avaLatency float64
	go func() {
		defer wg2.Done()
		all := int64(0)
		for ts := range latencyCh {
			all += ts.Microseconds()
		}
		avaLatency = float64(all) / float64(*reqNum)
	}()
	switch *method {
	case "set":
		wg.Add(1)
		go func() {
			defer close(runBuf)
			defer wg.Done()
			for i := 0; i < *reqNum; i++ {
				runBuf <- &benchmark.Request{
					ReqType: benchmark.SetOp,
					Key:     benchmark.GenRandString(*keySize),
					Val:     benchmark.GenRandString(*valSize),
				}
			}
		}()
	case "get":
		k, v := benchmark.GenRandString(*keySize), benchmark.GenRandString(*valSize)
		if err := clis[0].Set(context.Background(), k, v); err != nil {
			panic(err)
		}
		wg.Add(1)
		go func() {
			defer close(runBuf)
			defer wg.Done()
			for i := 0; i < *reqNum; i++ {
				runBuf <- &benchmark.Request{
					ReqType: benchmark.GetOp,
					Key:     k,
				}
			}
		}()
	default:
		panic(fmt.Errorf("invalid method: %s", *method))
	}
	time.Sleep(5 * time.Second)

	// Benchmark
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for op := range runBuf {
				switch op.ReqType {
				case benchmark.GetOp:
					start := time.Now()
					clis[0].Get(context.Background(), op.Key)
					latencyCh <- time.Since(start)
				case benchmark.SetOp:
					start := time.Now()
					clis[0].Set(context.Background(), op.Key, op.Val)
					latencyCh <- time.Since(start)
				default:
					panic(fmt.Sprintf("invalid operation: %v", op.ReqType))
				}
			}
		}()
	}
	start := time.Now()
	wg.Wait()
	close(latencyCh)
	wg2.Wait()
	fmt.Printf("Throughput: %v clients %v requests: %v req/s\n",
		*concurrency, *reqNum,
		int64(float64(*reqNum)/time.Since(start).Seconds()),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
}
