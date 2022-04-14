package main

import (
	"context"
	"fmt"
	"hybrid/veritas/benchmark"
	"hybrid/veritastm"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	veritasAddrs = kingpin.Flag("veritas-addrs", "Address of veritas node").Required().String()
)

func main() {
	ctx := context.Background()

	kingpin.Parse()

	cli, err := veritastm.Open(*veritasAddrs, benchmark.GenRandString(16))
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	defer cli.Close()

	res, err := cli.Set(ctx, "abc", "xyz")
	fmt.Printf("Set error %v\n", err)
	fmt.Printf("Set result %v\n", res)

	res, err = cli.Get(ctx, "abc")
	fmt.Printf("Get key %v value %v\n", "abc", res)
}
