package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	pbv "github.com/nusdbsystem/hybrid/veritas_hotstuff/proto/veritashs"

	"github.com/nusdbsystem/hybrid/veritas_hotstuff/cmd/config"
	"github.com/nusdbsystem/hybrid/veritas_hotstuff/svrnode"

	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

func main() {

	configFile := flag.String("config", "", "The path to the config file")
	//priv := flag.String("privkey", "", "The path to the private key file")
	flag.Parse()

	fmt.Println(*configFile)
	var conf config.Options
	err := config.ReadConfig(&conf, *configFile) //default config file "hotstuff.toml"
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	svr, err := svrnode.NewServerNode(&conf, *configFile)
	if err != nil {
		//panic(err)
		fmt.Printf("%v", err)

	}
	pbv.RegisterVeritasNodeServer(s, svr)
	lis, err := net.Listen("tcp", conf.ServerNodeAddr)

	if err != nil {
		//panic(err)
		fmt.Printf("%v", err)
	} else {
		log.Println("Veritas Node listen address: " + conf.ServerNodeAddr)
	}

	go func() {
		log.Printf("Veritas Node Serving gRPC: %s", conf.ServerNodeAddr)
		s.Serve(lis)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	sig := <-ch
	log.Printf("Received signal %v, quiting gracefully", sig)
}
