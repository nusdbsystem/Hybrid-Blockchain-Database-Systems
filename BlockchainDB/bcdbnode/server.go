package service

import (
	"context"
	"log"

	"hybrid/BlockchainDB/bcdbnode/config"
	pbv "hybrid/BlockchainDB/proto/blockchaindb"
	sharding "hybrid/BlockchainDB/shardingMgr"
)

var _ pbv.BCdbNodeServer = (*ServerNode)(nil)

type ServerNode struct {
	shardingMgr *sharding.ShardingMgr
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
	return &ServerNode{shardingMgr: shamgr}, nil
}

func (sv *ServerNode) Get(ctx context.Context, req *pbv.GetRequest) (*pbv.GetResponse, error) {
	val, err := sv.shardingMgr.Read(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}
	return &pbv.GetResponse{Value: []byte(val)}, nil
}

func (sv *ServerNode) Set(ctx context.Context, req *pbv.SetRequest) (*pbv.SetResponse, error) {
	// Use serverclient instance to set
	tx, err := sv.shardingMgr.Write(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		return nil, err
	}
	return &pbv.SetResponse{Tx: tx}, nil
}

func (sv *ServerNode) Verify(ctx context.Context, req *pbv.VerifyRequest) (*pbv.VerifyResponse, error) {
	result, err := sv.shardingMgr.Verify(ctx, req.GetOpt(), req.GetKey())
	if err != nil {
		return nil, err
	}
	return &pbv.VerifyResponse{Success: result}, nil
}
