package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mediocregopher/radix/v3"

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

func rediSQLGet(r *radix.Pool, key string) (string, error) {
	var items [][]string
	err := r.Do(radix.Cmd(
		&items,
		"REDISQL.EXEC",
		"VERITAS",
		fmt.Sprintf("SELECT kv.value FROM kv WHERE kv.key=\"%s\";", key),
	))
	if err != nil {
		return "", err
	}
	/*
		if err != nil {
			fmt.Printf("Get %v |%s| |%v|\n", err, key, item)
			return "", err
		}
	*/
	return items[0][0], nil
}

func rediSQLSet(r *radix.Pool, key string, val string) error {
	err := r.Do(radix.Cmd(
		nil,
		"REDISQL.EXEC",
		"VERITAS",
		fmt.Sprintf("INSERT INTO kv VALUES(\"%s\", \"%s\");", key, val),
	))
	if err != nil {
		fmt.Printf("Set %v\n", err)
	}
	return err
}

func main() {
	kingpin.Parse()

	cli, err := dbconn.NewRedisqlConn(*redisAddr, *driverConcurrency)
	check(err)

	cli.Do(radix.Cmd(nil, "DEL", "VERITAS"))
	// create RediSQL DB
	if err = cli.Do(radix.Cmd(nil, "REDISQL.CREATE_DB", "VERITAS")); err != nil {
		panic(err)
	}
	if err := cli.Do(radix.Cmd(
		nil,
		"REDISQL.EXEC",
		"VERITAS",
		"CREATE TABLE IF NOT EXISTS kv(key TEXT, value TEXT);",
	)); err != nil {
		panic(err)
	}
	if err := cli.Do(radix.Cmd(
		nil,
		"REDISQL.EXEC",
		"VERITAS",
		"CREATE INDEX cust_key_ind ON kv(key);",
	)); err != nil {
		panic(err)
	}

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
				rediSQLSet(cli, kv[0], kv[1])
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
					_, err := rediSQLGet(cli, op.Key)
					latencyCh <- time.Since(start)
					if err != nil {
						atomic.AddInt64(&nokey, 1)
					}
				case benchmark.SetOp:
					start := time.Now()
					rediSQLSet(cli, op.Key, op.Val)
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
	fmt.Printf("Throughput of RediSQL with concurrency %v to handle %v requests: %v req/s\n",
		*driverConcurrency, reqNum,
		int64(float64(reqNum)/time.Since(start).Seconds()),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
	fmt.Printf("Keys not found %d\n", nokey)
}
