package veritastendermint

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-redis/redis/v8"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"

	pbv "hybrid/proto/veritas"

	"github.com/dgraph-io/badger/v3"
	nm "github.com/tendermint/tendermint/node"
)

type Server struct {
	node *nm.Node
	kv   *redis.Client
	l    *badger.DB
	mu   *sync.Mutex
}

func NewServer(db *badger.DB, r *redis.Client, configFile string) (*Server, error) {
	s := &Server{}
	var err error
	if err = s.newTendermint(NewKVStore(db), configFile); err != nil {
		return nil, err
	}
	s.kv = r
	s.l = db
	s.mu = &sync.Mutex{}
	return s, nil
}

func (s *Server) Run() error {
	return s.node.Start()
}

func (s *Server) StopAndWait() error {
	err := s.node.Stop()
	s.node.Wait()
	return err
}

func (s *Server) newTendermint(app abci.Application, configFile string) error {
	// Read config.
	if info, err := os.Stat(configFile); err != nil || !info.IsDir() {
		return nil
	}
	config := cfg.DefaultConfig()
	config.RootDir = filepath.Dir(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "viper failed to read config file")
	}
	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "viper failed to unmarshal config")
	}
	if err := config.ValidateBasic(); err != nil {
		return errors.Wrap(err, "config is invalid")
	}

	// Create logger.
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var err error
	logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)
	if err != nil {
		return errors.Wrap(err, "failed to parse log level")
	}

	// Read private validator.
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	// Read node key.
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return errors.Wrap(err, "failed to load node's key")
	}

	// Create node
	s.node, err = nm.NewNode(
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create new Tendermint node")
	}

	return nil
}

var _ pbv.NodeServer = (*Server)(nil)

func (s *Server) Get(ctx context.Context, req *pbv.GetRequest) (*pbv.GetResponse, error) {
	res := s.kv.Get(ctx, req.GetKey())
	if err := res.Err(); err != nil {
		return nil, err
	}
	return &pbv.GetResponse{Value: res.Val()}, nil
}

func (s *Server) Set(ctx context.Context, req *pbv.SetRequest) (*pbv.SetResponse, error) {
	if err := s.kv.Set(ctx, req.GetKey(), req.GetValue(), 0).Err(); err != nil {
		return nil, err
	}
	s.l.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(req.Key), []byte(req.Value))
		return nil
	})
	return &pbv.SetResponse{}, nil
}

func (s *Server) BatchSet(ctx context.Context, req *pbv.BatchSetRequest) (*pbv.BatchSetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pbv.BatchSetResponse{}, nil
}

func (s *Server) Verify(ctx context.Context, req *pbv.VerifyRequest) (*pbv.VerifyResponse, error) {
	return nil, nil
}
