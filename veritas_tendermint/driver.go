package veritastm

import (
	"context"

	"google.golang.org/grpc"

	pbv "github.com/nusdbsystem/hybridproto/veritas"
)

type Driver struct {
	signature string
	cc        *grpc.ClientConn
	dbCli     pbv.NodeClient
}

func Open(serverAddr, signature string) (*Driver, error) {
	cc, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	dbCli := pbv.NewNodeClient(cc)

	return &Driver{
		signature: signature,
		cc:        cc,
		dbCli:     dbCli,
	}, nil
}

func (d *Driver) Get(ctx context.Context, key string) (string, error) {
	res, err := d.dbCli.Get(ctx, &pbv.GetRequest{
		Signature: d.signature,
		Key:       key,
	})
	if err != nil {
		return "", err
	}
	return res.GetValue(), nil
}

func (d *Driver) Set(ctx context.Context, key, value string, version int64) (string, error) {
	res, err := d.dbCli.Set(ctx, &pbv.SetRequest{
		Signature: d.signature,
		Key:       key,
		Value:     value,
		Version:   version,
	})

	if err != nil {
		return "", err
	}

	return res.Txid, nil
}

func (d *Driver) Close() error {
	return d.cc.Close()
}
