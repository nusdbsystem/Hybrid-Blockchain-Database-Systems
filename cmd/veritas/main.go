package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/dbconn"
	"hybrid/kafkarole"
	pbv "hybrid/proto/veritas"
	"hybrid/veritas"
)

var (
	signature  = kingpin.Flag("signature", "server signature").Required().String()
	blockSize  = kingpin.Flag("blk-size", "block size").Default("100").Int()
	parties    = kingpin.Flag("parties", "party1,party2,...").Required().String()
	addr       = kingpin.Flag("addr", "server address").Required().String()
	kafkaAddr  = kingpin.Flag("kafka-addr", "kafka server address").Required().String()
	kafkaGroup = kingpin.Flag("kafka-group", "kafka group id").Required().String()
	kafkaTopic = kingpin.Flag("kafka-topic", "kafka topic").Required().String()
	redisAddr  = kingpin.Flag("redis-addr", "redis server address").Required().String()
	redisDb    = kingpin.Flag("redis-db", "redis db number").Required().Int()
	redisPwd   = kingpin.Flag("redis-pwd", "redis password").String()
	ledgerPath = kingpin.Flag("ledger-path", "ledger path").Required().String()
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	kingpin.Parse()

	r, err := dbconn.NewRedisConn(*redisAddr, *redisPwd, *redisDb)
	check(err)

	c, err := kafkarole.NewConsumer(*kafkaAddr, *kafkaGroup, []string{*kafkaTopic})
	check(err)
	p, err := kafkarole.NewProducer(*kafkaAddr, *kafkaTopic)
	check(err)

	pm := make(map[string]struct{})
	pm[*signature] = struct{}{}

	ps := strings.Split(*parties, ",")
	for _, s := range ps {
		pm[s] = struct{}{}
	}

	s := grpc.NewServer()
	svr := veritas.NewServer(r, c, p, *ledgerPath, &veritas.Config{
		Signature: *signature,
		Topic:     *kafkaTopic,
		Parties:   pm,
		BlockSize: *blockSize,
	})
	pbv.RegisterNodeServer(s, svr)
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Printf("Serving gRPC: %s", *addr)
		s.Serve(lis)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	sig := <-ch
	log.Printf("Received signal %v, quiting gracefully", sig)
}
