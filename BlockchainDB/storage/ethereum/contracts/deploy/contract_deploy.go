package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"hybrid/BlockchainDB/bcdbnode/config"
	ClientSDK "hybrid/BlockchainDB/storage/ethereum/clientSDK"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	//local
	// ethnode := "/home/tianwen/Data/eth_1_1/geth.ipc"

	//ganache
	// ethnode := "http://localhost:7545"

	//config from file
	configFile := flag.String("config", "config/config.eth.1.4/shard_1", "The path to the config file")
	flag.Parse()
	var conf config.Options
	err := config.ReadConfig(&conf, *configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %v\n", err)
		os.Exit(1)
	}
	ethnode := conf.EthNode
	hexkey := conf.EthHexKey

	hexaddress, tx, _, err := ClientSDK.DeployEthereumKVStoreContract(ethnode, hexkey)
	if err != nil {
		log.Fatal("DeployEthereumKVStoreContract", err)
	}

	fmt.Printf("eth-hexaddr = \"%v\"\n", hexaddress)
	fmt.Printf("contract-tx = \"%v\"\n", tx)
	time.Sleep(10 * time.Second)

	//Debug
	client, err := ethclient.Dial(ethnode)
	if err != nil {
		fmt.Println("error ethclient Dail "+ethnode, err)
	}
	address := common.HexToAddress(hexaddress)
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	isContract := len(bytecode) > 0
	fmt.Printf("contract = %v\n", isContract) // is contract: true
}
