package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	SelfID         string `mapstructure:"self-id"`
	ServerNodeAddr string `mapstructure:"server-node-addr"`
	Type           string `mapstructure:"shard-type"`
	Delay          int    `mapstructure:"delay"`
	EthNode        string `mapstructure:"eth-node"`
	EthHexAddr     string `mapstructure:"eth-hexaddr"`
	EthHexKey      string `mapstructure:"eth-hexkey"`
	EthBootSigner  string `mapstructure:"eth-boot-signer-address"`
	FabNode        string `mapstructure:"fab-node"`
	FabConfig      string `mapstructure:"fab-config"`
	ShardNumber    int    `mapstructure:"shard-number"`
	Shards         []Shard
}

type Shard struct {
	ID           string `mapstructure:"shard-id"`
	Type         string `mapstructure:"shard-type"`
	PartitionKey string `mapstructure:"shard-patition-key"`
	EthNode      string `mapstructure:"eth-node"`
	EthHexAddr   string `mapstructure:"eth-hexaddr"`
	EthHexKey    string `mapstructure:"eth-hexkey"`
	FabNode      string `mapstructure:"fab-node"`
	FabConfig    string `mapstructure:"fab-config"`
}

func ReadConfig(opts interface{}, configfile string) (err error) {
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}
	if configfile != "" {
		//viper.SetConfigFile(configfile)
		viper.SetConfigName(configfile)
	} else {
		viper.SetConfigName("blockchaindb")
	}

	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(opts)
	if err != nil {
		return err
	}

	return nil
}
