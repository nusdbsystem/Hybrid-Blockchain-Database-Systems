package main

import (
	"fmt"
	"hybrid/tso"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	tsoAddr     = kingpin.Flag("addr", "tso server address").Default(":7070").String()
	concurrency = kingpin.Flag("concurrency", "client num").Default("20").Int()
	reqNum      = kingpin.Flag("req-num", "request num").Default("1000000").Int()
)

func main() {
	kingpin.Parse()

	avaReqNum := *reqNum / (*concurrency)
	wg := &sync.WaitGroup{}
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			cli, err := tso.NewClient(*tsoAddr)
			if err != nil {
				panic(err)
			}
			defer cli.Close()

			for j := 0; j < avaReqNum; j++ {
				if _, err := cli.TS(); err != nil {
					panic(err)
				}
			}
		}()
	}

	start := time.Now()
	wg.Wait()
	fmt.Printf("%v clients %v requests: %v req/s\n", *concurrency, *reqNum, int64(float64(*reqNum)/time.Since(start).Seconds()))
}
