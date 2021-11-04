package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgraph-io/badger/v3"
	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/veritastendermint"
)

var (
	dir        = kingpin.Flag("dir", "Directory of storing data").Required().String()
	configFile = kingpin.Flag("config", "Path to config.toml").Default("$HOME/.tendermint/config/config.toml").String()
)

func main() {
	kingpin.Parse()

	db, err := badger.Open(badger.DefaultOptions(*dir))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	svr, err := veritastendermint.NewServer(db, nil, *configFile)
	if err != nil {
		log.Fatalf("New tendermint server failed: %v", err)
	}
	svr.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	svr.StopAndWait()
}
