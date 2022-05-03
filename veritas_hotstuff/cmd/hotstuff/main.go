package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
	hs "github.com/wtwinlab/hotstuff/cmd/hotstuffserver/clientsrv"
)

func main() {

	configFile := flag.String("config", "", "The path to the config file")
	//priv := flag.String("privkey", "", "The path to the private key file")
	flag.Parse()

	fmt.Printf("Debug configFile: %s\n", *configFile)
	//versoin 10.2.2
	hs.InitHotstuffServer(*configFile) //(ctx, &conf)

}
