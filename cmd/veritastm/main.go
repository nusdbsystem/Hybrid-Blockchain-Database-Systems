package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/dbconn"
	pbv "hybrid/proto/veritas"
	"hybrid/veritastm"

	abciserver "github.com/tendermint/tendermint/abci/server"
)

var (
	signature  = kingpin.Flag("signature", "server signature").Required().String()
	blockSize  = kingpin.Flag("blk-size", "block size").Default("100").Int()
	parties    = kingpin.Flag("parties", "party1,party2,...").Required().String()
	addr       = kingpin.Flag("addr", "server address").Required().String()
	redisAddr  = kingpin.Flag("redis-addr", "redis server address").Required().String()
	redisDb    = kingpin.Flag("redis-db", "redis db number").Required().Int()
	redisPwd   = kingpin.Flag("redis-pwd", "redis password").String()
	ledgerPath = kingpin.Flag("ledger-path", "ledger path").Required().String()
	tmSocket   = kingpin.Flag("tendermint-socket", "tendermint socket").Required().String()
	abciSocket = kingpin.Flag("abci-socket", "abci socket").Required().String()
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

	pm := make(map[string]struct{})
	pm[*signature] = struct{}{}

	ps := strings.Split(*parties, ",")
	for _, s := range ps {
		pm[s] = struct{}{}
	}

	s := grpc.NewServer()
	svr := veritastm.NewServer(r, &veritastm.Config{
		Signature:   *signature,
		Parties:     pm,
		BlockSize:   *blockSize,
		LedgerPath:  *ledgerPath,
		ABCIRPCAddr: *abciSocket,
	})
	pbv.RegisterNodeServer(s, svr)
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	server := abciserver.NewSocketServer(*tmSocket, svr.Ledger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()

	go func() {
		fmt.Printf("Serving gRPC on port: %s\n", *addr)
		s.Serve(lis)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	sig := <-ch
	fmt.Printf("Received signal %v, quiting gracefully", sig)
}
