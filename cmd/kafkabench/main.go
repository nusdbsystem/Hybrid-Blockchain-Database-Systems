package main

import (
	"fmt"
	"hybrid/kafkarole"
	"log"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	op          = kingpin.Flag("op", "read/write").Required().Enum("read", "write")
	kafkaAddr   = kingpin.Flag("kafka-addr", "kafka server address").Default("0.0.0.0:9092").String()
	kafkaGroup  = kingpin.Flag("kafka-group", "kafka group id").Default("0").String()
	kafkaTopic  = kingpin.Flag("kafka-topic", "kafka topic").Default("shared-log").String()
	concurrency = kingpin.Flag("concurrency", "operation concurrency").Default("20").Int()
	msgSize     = kingpin.Flag("size", "size of message").Default("128").Int()
	reqNum      = kingpin.Flag("req-num", "number of request").Default("100000").Int()
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GenRandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func main() {
	kingpin.Parse()

	c, err := kafkarole.NewConsumer(*kafkaAddr, *kafkaGroup, []string{*kafkaTopic})
	check(err)
	p, err := kafkarole.NewProducer(*kafkaAddr, *kafkaTopic)
	check(err)

	msgs := make([]string, 0)
	for i := 0; i < *reqNum; i++ {
		msgs = append(msgs, GenRandString(*msgSize))
	}

	wg := sync.WaitGroup{}
	switch *op {
	case "read":
		log.Println("write msgs for read...")
		for _, msg := range msgs {
			check(p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: kafkaTopic, Partition: kafka.PartitionAny},
				Value:          []byte(msg),
			}, nil))
		}
		log.Println("start reading...")
		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					_, err := c.ReadMessage(100 * time.Millisecond)
					if err.(kafka.Error).Code() == kafka.ErrTimedOut {
						return
					} else if err != nil {
						panic(err)
					}
				}
			}()
		}
	case "write":
		log.Println("start writing...")
		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go func(seq int) {
				defer wg.Done()
				for j := seq; j < *reqNum; j += *concurrency {
					dc := make(chan kafka.Event, 100)
					check(p.Produce(&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: kafkaTopic, Partition: kafka.PartitionAny},
						Value:          []byte(msgs[j]),
					}, dc))
					<-dc
				}
			}(i)
		}
	default:
		panic(fmt.Errorf("invalid operation: %s", *op))
	}
	start := time.Now()
	wg.Wait()
	allTime := time.Since(start)
	log.Println("finish test in", allTime.String())

	fmt.Printf("Throughput: %v op/s %v MB/s\n",
		int64(float64(*reqNum)/allTime.Seconds()),
		int64(float64(*reqNum*(*msgSize))/allTime.Seconds()),
	)

}
