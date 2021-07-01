package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"hybrid/tso"
)

var (
	addr = kingpin.Flag("addr", "tso server address").Default(":7070").String()
)

func main() {
	kingpin.Parse()

	client, err := tso.NewClient(*addr)
	if err != nil {
		panic(err)
	}
	if ts, err := client.TS(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ts)
	}
}
