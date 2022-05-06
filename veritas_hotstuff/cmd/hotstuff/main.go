package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
	hs "github.com/wtwinlab/hotstuff/cmd/hotstuffserver/clientsrv"
)

func main() {

	configFile := flag.String("config", "", "The path to the config file")
	flag.Parse()

	fmt.Printf("Debug configFile: %s\n", *configFile)
	//version 0.2.2
	hs.InitHotstuffServer(*configFile)

}
