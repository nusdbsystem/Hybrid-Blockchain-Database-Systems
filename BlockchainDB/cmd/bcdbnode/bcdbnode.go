package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	service "hybrid/BlockchainDB/bcdbnode"
	"hybrid/BlockchainDB/bcdbnode/config"
	pbv "hybrid/BlockchainDB/proto/blockchaindb"
)

func main() {

	configFile := flag.String("config", "config.toml", "The path to the config file")
	flag.Parse()
	var conf config.Options
	err := config.ReadConfig(&conf, *configFile) //default config file "config.toml"
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()

	svr, err := service.NewServerNode(&conf)
	if err != nil {
		log.Fatalf("New ServerNode err %v", err)
	}

	pbv.RegisterBCdbNodeServer(s, svr)
	lis, err := net.Listen("tcp", conf.ServerNodeAddr)
	if err != nil {
		log.Fatalf("Node listen err %v", err)
	} else {
		log.Println("Node listen address: " + conf.ServerNodeAddr)
	}

	go func() {
		log.Println("Node Serving gRPC: ", conf.ServerNodeAddr)
		s.Serve(lis)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	sig := <-ch
	log.Printf("Received signal %v, quiting gracefully", sig)
}
