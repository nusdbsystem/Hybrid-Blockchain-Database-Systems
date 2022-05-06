package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wtwinlab/hotstuff"
)

type Options struct {
	ClientAddr     string      `mapstructure:"client-listen"`
	PeerAddr       string      `mapstructure:"peer-listen"`
	ServerNodeAddr string      `mapstructure:"self-veritas-node"`
	RedisAddr      string      `mapstructure:"self-redis-address"`
	LedgerPath     string      `mapstructure:"self-ledger-path"`
	SelfID         hotstuff.ID `mapstructure:"self-id"`
	Replicas       []Replica
}

type Replica struct {
	ID         hotstuff.ID
	PeerAddr   string `mapstructure:"peer-address"`
	ClientAddr string `mapstructure:"client-address"`
	RedisAddr  string `mapstructure:"redis-address"`
	LedgerPath string `mapstructure:"ledger-path"`
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
		viper.SetConfigName("hotstuff")
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
