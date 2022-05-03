package svrnode

import (
	"context"
	"fmt"

	"hybrid/veritas_hotstuff/cmd/config"
	pbv "hybrid/veritas_hotstuff/proto/veritashs"
	"hybrid/veritas_hotstuff/storage"

	hsc "github.com/wtwinlab/hotstuff/cmd/hotstuffclient/client"

	//hsc "github.com/wtwinlab/hotstuff/client"
	veritas "hybrid/veritas_kafka"

	"google.golang.org/protobuf/proto"
)

var _ pbv.VeritasNodeServer = (*ServerNode)(nil)

type ServerNode struct {
	sharedTable *storage.RedisKV
	// hotstuffclient instance
	hotstuffclient *hsc.HotstuffClient
}

func NewServerNode(conf *config.Options, configFile string) (*ServerNode, error) {
	rdb, err := storage.NewRedisKV(conf.RedisAddr, "", 1)
	//log.Println("New Server Node with redis address: " + conf.RedisAddr)
	if err != nil {
		return nil, err
	}

	// Init serverclient instance of ServerNode with default config file "hotstuff.toml"
	cli, err := hsc.InitHotstuffClient(configFile) //conf

	if err != nil {
		return nil, err
	}
	return &ServerNode{sharedTable: rdb, hotstuffclient: cli}, nil
}

func (sn *ServerNode) Get(ctx context.Context, req *pbv.GetRequest) (*pbv.GetResponse, error) {
	res, err := sn.sharedTable.Get([]byte(req.GetKey()))
	if err != nil {
		return nil, err
	}
	v, err := veritas.Decode(string(res))
	if err != nil {
		return nil, err
	}

	return &pbv.GetResponse{Value: v.Val, Version: v.Version}, nil
}

func (sn *ServerNode) Set(ctx context.Context, req *pbv.SetRequest) (*pbv.SetResponse, error) {
	// Use serverclient instance to set
	cmd, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = sn.hotstuffclient.SendCommands(ctx, cmd)
	if err != nil {
		return nil, err
	}
	fmt.Println("Debug20220222-srv " + string(req.Key))
	return &pbv.SetResponse{}, nil
}

// func (sn *ServerNode) BatchSet(ctx context.Context, reqs *pbv.BatchSetRequest) (*pbv.BatchSetResponse, error) {
// 	// Use serverclient instance to set
// 	err := sn.hotstuffclient.SendCommands(ctx, reqs.Sets)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pbv.BatchSetResponse{}, nil
// }
