package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	pb "hybrid/proto/raftkv"
	"hybrid/raftkv"
)

var (
	svrAddr    = kingpin.Flag("svr-addr", "Address of server").Default(":19001").String()
	raftAddr   = kingpin.Flag("raft-addr", "Address of raft module").Default("127.0.0.1:18001").String()
	dir        = kingpin.Flag("dir", "Dir for data and log").Required().String()
	raftLeader = kingpin.Flag("raft-leader", "Address of the existing raft cluster leader").String()
	redisAddr  = kingpin.Flag("redis-addr", "redis server address").String()
	redisDb    = kingpin.Flag("redis-db", "redis db number").Int()
	redisPwd   = kingpin.Flag("redis-pwd", "redis password").String()
	storage    = kingpin.Flag("store", "Underlying storage [redis/badger]").Default("redis").Enum("redis", "badger")
	blkSize    = kingpin.Flag("blk-size", "Block size in raft").Default("100").Int()
)

func main() {
	kingpin.Parse()

	dataDir, logDir, ledgerDir := filepath.Join(*dir, "data"), filepath.Join(*dir, "log"), filepath.Join(*dir, "ledger")

	var (
		kv  raftkv.KV
		err error
	)
	switch *storage {
	case "redis":
		kv, err = raftkv.NewRedisKV(*redisAddr, *redisPwd, *redisDb)
		if err != nil {
			panic(err)
		}
	case "badger":
		kv, err = raftkv.NewBadgerKV(dataDir)
		if err != nil {
			panic(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	node, err := raftkv.NewPeer(ctx, ledgerDir, kv, &raftkv.Config{
		Id:        fmt.Sprintf("%v", time.Now().UnixNano()),
		RaftDir:   logDir,
		RaftBind:  *raftAddr,
		RaftJoin:  *raftLeader,
		BlockSize: *blkSize,
	})
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterDBServer(s, node)
	lis, err := net.Listen("tcp", *svrAddr)
	if err != nil {
		panic(err)
	}
	go func() {
		log.Printf("Serving gRPC: %s", *svrAddr)
		s.Serve(lis)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	sig := <-ch
	cancel()
	log.Printf("Received signal %v, quiting gracefully", sig)
}
