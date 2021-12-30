package service

import (
	"context"
	"log"
	"time"

	"hybrid/BlockchainDB/bcdbnode/config"
	pbv "hybrid/BlockchainDB/proto/blockchaindb"
	sharding "hybrid/BlockchainDB/shardingMgr"
)

var _ pbv.BCdbNodeServer = (*ServerNode)(nil)

type ServerNode struct {
	shardingMgr *sharding.ShardingMgr
	txDelay     int
}

func NewServerNode(conf *config.Options) (*ServerNode, error) {
	// version 1.0
	// ethereumconn, err := EthClientSDK.NewEthereumKVStoreInstance(conf.EthNode, conf.EthHexAddr, conf.EthHexKey)
	// if err != nil {
	// 	log.Println("Failed to NewEthereumKVStoreInstance", err)
	// 	return nil, err
	// }
	// fabricconn, err := FabClientSDK.NewFabricKVStoreInstance()
	// if err != nil {
	// 	fmt.Println("Failed to NewFabricKVStoreInstance", err)
	// 	return nil, err
	// }
	// shamgr := &sharding.ShardingMgr{EthConn: ethereumconn, FabConn: fabricconn}

	// version 2.0
	shamgr, err := sharding.NewShardingMgr(conf)
	if err != nil {
		log.Println("Failed to NewShardingMgr", err)
		return nil, err
	}
	if conf.Delay > 0 {
		log.Println("Enable txDelay Experiment(ms): ", conf.Delay)
	}
	return &ServerNode{shardingMgr: shamgr, txDelay: conf.Delay}, nil
}

func (sv *ServerNode) Get(ctx context.Context, req *pbv.GetRequest) (*pbv.GetResponse, error) {
	if sv.txDelay > 0 {
		time.Sleep(time.Duration(sv.txDelay) * time.Millisecond)
	}
	val, err := sv.shardingMgr.Read(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}
	return &pbv.GetResponse{Value: []byte(val)}, nil
}

func (sv *ServerNode) Set(ctx context.Context, req *pbv.SetRequest) (*pbv.SetResponse, error) {
	if sv.txDelay > 0 {
		time.Sleep(time.Duration(sv.txDelay) * time.Millisecond)
	}
	tx, err := sv.shardingMgr.Write(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		return nil, err
	}
	return &pbv.SetResponse{Tx: tx}, nil
}

func (sv *ServerNode) Verify(ctx context.Context, req *pbv.VerifyRequest) (*pbv.VerifyResponse, error) {
	result, err := sv.shardingMgr.Verify(ctx, req.GetOpt(), req.GetKey(), req.GetTx())
	if err != nil {
		return nil, err
	}
	return &pbv.VerifyResponse{Success: result}, nil
}
