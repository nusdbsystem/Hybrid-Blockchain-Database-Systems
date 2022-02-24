package config

import (
	"github.com/EinWTW/hotstuff"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	BatchSize       int         `mapstructure:"batch-size"`
	Benchmark       bool        `mapstructure:"benchmark"`
	Cert            string      `mapstructure:"cert"`
	CertKey         string      `mapstructure:"cert-key"`
	Crypto          string      `mapstructure:"crypto"`
	Consensus       string      `mapstructure:"consensus"`
	ClientAddr      string      `mapstructure:"client-listen"`
	ExitAfter       int         `mapstructure:"exit-after"`
	Input           string      `mapstructure:"input"`
	LeaderID        hotstuff.ID `mapstructure:"leader-id"`
	MaxInflight     uint64      `mapstructure:"max-inflight"`
	Output          string      `mapstructure:"print-commands"`
	PayloadSize     int         `mapstructure:"payload-size"`
	PeerAddr        string      `mapstructure:"peer-listen"`
	PmType          string      `mapstructure:"pacemaker"`
	PrintThroughput bool        `mapstructure:"print-throughput"`
	Privkey         string      `mapstructure:"privkey"`
	RateLimit       int         `mapstructure:"rate-limit"`
	RootCAs         []string    `mapstructure:"root-cas"`
	SelfID          hotstuff.ID `mapstructure:"self-id"`
	ServerNodeAddr  string      `mapstructure:"self-veritas-node"`
	RedisAddr       string      `mapstructure:"self-redis-address"`
	TLS             bool        `mapstructure:"tls"`
	ViewTimeout     float64     `mapstructure:"view-timeout"`
	ConnectTimeout  float64     `mapstructure:"view-timeout"`
	Replicas        []Replica
}

type Replica struct {
	ID         hotstuff.ID
	PeerAddr   string `mapstructure:"peer-address"`
	ClientAddr string `mapstructure:"client-address"`
	RedisAddr  string `mapstructure:"redis-address"`
	Pubkey     string `mapstructure:"pubkey"`
	Cert       string `mapstructure:"cert"`
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
