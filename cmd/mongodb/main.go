package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/dbconn"
	"hybrid/veritas/benchmark"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dataLoad          = kingpin.Flag("load-path", "Path of YCSB initial data").Required().String()
	dataRun           = kingpin.Flag("run-path", "Path of YCSB operation data").Required().String()
	driverConcurrency = kingpin.Flag("nthreads", "Number of threads for each driver").Default("10").Int()
	mongoAddr         = kingpin.Flag("mongo-addr", "monogodb server address").Required().String()
	mongoPort         = kingpin.Flag("mongo-port", "mongodb port").Required().String()
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

func MongoGet(cli *mongo.Client, key string) (string, error) {
	collection := cli.Database("test").Collection("kv")
	var result map[string]string
	filter := bson.M{"key": key}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		// fmt.Println(err)
		return "", err
	}
	return result["val"], nil
}

func MongoSet(cli *mongo.Client, key string, val string) error {
	collection := cli.Database("test").Collection("kv")
	_, err := collection.InsertOne(context.Background(), bson.M{"key": key, "val": val})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	kingpin.Parse()

	cli, err := dbconn.NewMongoConn(context.Background(), *mongoAddr, *mongoPort)
	check(err)
	/*
	mod := mongo.IndexModel{
		Keys: bson.M{
			"key": 1, // index in ascending order
		}, Options: nil,
	}
	collection := cli.Database("test").Collection("kv")
	_, err = collection.Indexes().CreateOne(context.Background(), mod)
	check(err)
	*/
	var reqNum int64
	reqNum = 0
	var reqNumSet int64
	reqNumSet = 0
	var reqNumGet int64
	reqNumGet = 0

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
	latencyCh := make(chan time.Duration, 100000)
	latencySetCh := make(chan time.Duration, 100000)
	latencyGetCh := make(chan time.Duration, 100000)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	var avaLatency float64
	var avgSetLatency float64
	var avgGetLatency float64
	go func() {
		defer wg2.Done()
		all := int64(0)
		for ts := range latencyCh {
			all += ts.Microseconds()
		}
		avaLatency = float64(all) / (1000 * float64(atomic.LoadInt64(&reqNum)))
		setlat := int64(0)
		for ts := range latencySetCh {
                        setlat += ts.Microseconds()
                }
                avgSetLatency = float64(setlat) / (1000 * float64(atomic.LoadInt64(&reqNumSet)))
		getlat := int64(0)
                for ts := range latencyGetCh {
                        getlat += ts.Microseconds()
                }
                avgGetLatency = float64(getlat) / (1000 * float64(atomic.LoadInt64(&reqNumGet)))
	}()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range loadBuf {
				MongoSet(cli, kv[0], kv[1])
			}
		}()
	}
	wg.Wait()
	// fmt.Println("Init done.")

	runFile, err := os.Open(*dataRun)
	if err != nil {
		panic(err)
	}
	defer runFile.Close()
	runBuf := make(chan *benchmark.Request, 20*(*driverConcurrency))
	// wg.Add(1)
	go func() {
		// defer wg.Done()
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
	// wg.Wait()
	// time.Sleep(5 * time.Second)
	// fmt.Println("Read file done.")

	var nokey int64
	nokey = 0

	start := time.Now()
	for j := 0; j < *driverConcurrency; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for op := range runBuf {
				switch op.ReqType {
				case benchmark.GetOp:
					start := time.Now()
					_, err := MongoGet(cli, op.Key)
					// var err error
					// err = nil
					dt := time.Since(start)
					latencyCh <- dt
					latencyGetCh <- dt
					atomic.AddInt64(&reqNumGet, 1)
					if err != nil {
						atomic.AddInt64(&nokey, 1)
					}
				case benchmark.SetOp:
					start := time.Now()
					MongoSet(cli, op.Key, op.Val)
					dt := time.Since(start)
					latencyCh <- dt
					latencySetCh <- dt
                                        atomic.AddInt64(&reqNumSet, 1)
				default:
					panic(fmt.Sprintf("invalid operation: %v", op.ReqType))
				}
			}
		}()
	}
	wg.Wait()
	dtime := time.Since(start).Seconds()
	close(latencyCh)
	close(latencyGetCh)
	close(latencySetCh)
	wg2.Wait()
	fmt.Printf("Throughput of MongoDB with client concurrency %v handling %v requests: %v req/s\n",
		*driverConcurrency, reqNum,
		int64(float64(atomic.LoadInt64(&reqNum))/dtime),
	)
	fmt.Printf("Average latency: %v ms\n", avaLatency)
	fmt.Printf("Average Get latency: %v ms\n", avgGetLatency)
	fmt.Printf("Average Set latency: %v ms\n", avgSetLatency)
	fmt.Printf("Keys not found %d\n", nokey)
}
