package test

import (
	Connectors "hybrid/BlockchainDB/blockchainconnectors"
)

func NewShardingMgr() (Connectors.BlockchainConnector, error) {
	var bcConn Connectors.BlockchainConnector
	bcConn = &Connectors.Testconnector{}
	return bcConn, nil
	// bcConn, err := EthClientSDK.NewEthereumKVStoreInstance(conf.EthNode, conf.EthHexAddr, conf.EthHexKey)
	// if err != nil {
	// 	log.Println("Failed to NewEthereumKVStoreInstance", err)
	// 	return nil, err
	// }
	//return &ShardingMgr{BCConn: bcConn}, nil
}
