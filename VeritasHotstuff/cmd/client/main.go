package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	hs "hybrid/VeritasHotstuff/cmd/client/hotstuff"
	hsconf "hybrid/VeritasHotstuff/cmd/config"

	flag "github.com/spf13/pflag"
)

func main() {

	configFile := flag.String("config", "", "The path to the config file")

	flag.Uint32("self-id", 0, "The id for this replica.")
	flag.Bool("tls", false, "Enable TLS")
	flag.Int("exit-after", 0, "Number of seconds after which the program should exit.")

	flag.String("input", "", "Optional file to use for payload data (client)")
	flag.Bool("benchmark", false, "If enabled, a BenchmarkData protobuf will be written to stdout. (client)")
	flag.Int("rate-limit", 0, "Limit the request-rate to approximately (in requests per second). (client)")
	flag.Int("payload-size", 0, "The size of the payload in bytes (client)")
	flag.Uint64("max-inflight", 10000, "The maximum number of messages that the client can wait for at once (client)")

	flag.Parse()

	var conf hsconf.Options
	err := hsconf.ReadConfig(&conf, *configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	hs.NewClient(ctx, &conf)
}
