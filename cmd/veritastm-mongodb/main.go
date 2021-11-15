package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	pbv "hybrid/proto/veritas"
	"hybrid/veritastm"

	abciserver "github.com/tendermint/tendermint/abci/server"
)

var (
	signature  = kingpin.Flag("signature", "server signature").Required().String()
	blockSize  = kingpin.Flag("blk-size", "block size").Default("100").Int()
	parties    = kingpin.Flag("parties", "party1,party2,...").Required().String()
	addr       = kingpin.Flag("addr", "server address").Required().String()
	mongoAddr  = kingpin.Flag("mongodb-addr", "MongoDB address").Required().String()
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

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(*mongoAddr))
	check(err)

	pm := make(map[string]struct{})
	pm[*signature] = struct{}{}

	ps := strings.Split(*parties, ",")
	for _, s := range ps {
		pm[s] = struct{}{}
	}

	s := grpc.NewServer()
	svr := veritastm.NewServer(client, &veritastm.Config{
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
		fmt.Printf("Veritas + Tendermint+ MongoDB: Serving gRPC on port: %s\n", *addr)
		s.Serve(lis)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	sig := <-ch
	fmt.Printf("Received signal %v, quiting gracefully", sig)
}
