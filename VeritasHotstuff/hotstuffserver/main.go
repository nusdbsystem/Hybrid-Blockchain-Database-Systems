package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	hsconf "hybrid/VeritasHotstuff/cmd/config"

	flag "github.com/spf13/pflag"
)

func main() {

	configFile := flag.String("config", "", "The path to the config file")

	flag.Uint32("self-id", 0, "The id for this replica.")
	flag.Bool("tls", false, "Enable TLS")
	flag.Int("exit-after", 0, "Number of seconds after which the program should exit.")

	flag.String("privkey", "", "The path to the private key file (server)")
	flag.String("cert", "", "Path to the certificate (server)")
	flag.String("cert-key", "", "Path to the private key for the certificate (server)")
	flag.String("crypto", "ecdsa", "The name of the crypto implementation to use (server)")
	flag.String("consensus", "chainedhotstuff", "The name of the consensus implementation to use (server)")
	flag.Int("view-timeout", 10000, "How many milliseconds before a view is timed out (server)")
	flag.String("output", "", "Commands will be written here. (server)")
	flag.Int("batch-size", 100, "How many commands are batched together for each proposal (server)")
	flag.Bool("print-throughput", false, "Throughput measurements will be printed stdout (server)")
	flag.String("client-listen", "", "Override the listen address for the client server (server)")
	flag.String("peer-listen", "", "Override the listen address for the replica (peer) server (server)")

	flag.Parse()

	var conf hsconf.Options
	err := hsconf.ReadConfig(&conf, *configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Printf("Debug self-id: %d, privkey: %s\n", conf.SelfID, conf.Privkey)

	newServer(ctx, &conf)
}
