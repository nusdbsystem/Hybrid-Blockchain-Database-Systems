package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	pbv "hybrid/VeritasHotstuff/proto/veritas"

	"hybrid/VeritasHotstuff/cmd/config"
	"hybrid/VeritasHotstuff/svrnode"

	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var (
// svrAddr   = kingpin.Flag("svr-addr", "server address").Default("0.0.0.0:1997").String()
// redisAddr = kingpin.Flag("redis-addr", "redis server address").Required().String()

)

func main() {
	//kingpin.Parse()
	configFile := flag.String("config", "", "The path to the config file")
	flag.Parse()
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	var conf config.Options
	err := config.ReadConfig(&conf, *configFile) //default config file "hotstuff.toml"
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	svr, err := svrnode.NewServerNode(&conf)
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
