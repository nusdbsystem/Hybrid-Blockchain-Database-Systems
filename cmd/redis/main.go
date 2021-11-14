package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/dbconn"
	"hybrid/veritas/benchmark"
)

var (
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
	redisAddr         = kingpin.Flag("redis-addr", "redis server address").Required().String()
	redisDb           = kingpin.Flag("redis-db", "redis db number").Required().Int()
	redisPwd          = kingpin.Flag("redis-pwd", "redis password").String()
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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

	cli, err := dbconn.NewRedisConn(*redisAddr, *redisPwd, *redisDb)
	check(err)

	var reqNum int64
	reqNum = 0

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
		avaLatency = float64(all) / (1000 * float64(atomic.LoadInt64(&reqNum)))
	}()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range loadBuf {
				cli.Set(context.Background(), kv[0], kv[1], 0).Result()
			}
		}()
	}
	wg.Wait()

	runFile, err := os.Open(*dataRun)
	if err != nil {
		panic(err)
	}
	defer runFile.Close()
	runBuf := make(chan *benchmark.Request, 20*(*driverConcurrency))
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(runBuf)
		if err := benchmark.LineByLine(runFile, func(line string) error {
			atomic.AddInt64(&reqNum, 1)
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

	var nokey int64
	nokey = 0

	for j := 0; j < *driverConcurrency; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for op := range runBuf {
				switch op.ReqType {
				case benchmark.GetOp:
					start := time.Now()
					_, err := cli.Get(context.Background(), op.Key).Result()
					latencyCh <- time.Since(start)
					if err == redis.Nil {
						atomic.AddInt64(&nokey, 1)
					}
				case benchmark.SetOp:
					start := time.Now()
					cli.Set(context.Background(), op.Key, op.Val, 0).Result()
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
	fmt.Printf("Throughput Redis of with %v concurrency to handle %v requests: %v req/s\n",
		*driverConcurrency, reqNum,
		int64(float64(atomic.LoadInt64(&reqNum))/time.Since(start).Seconds()),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
	fmt.Printf("Keys not found %d\n", nokey)
}
