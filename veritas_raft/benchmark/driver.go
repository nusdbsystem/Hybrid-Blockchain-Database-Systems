package benchmark

import (
	"context"

	"google.golang.org/grpc"

	pbv "github.com/nusdbsystem/hybridproto/raftkv"
)

type Driver struct {
	signature string
	cc        *grpc.ClientConn
	dbCli     pbv.DBClient
}

func Open(serverAddr, signature string) (*Driver, error) {
	cc, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	dbCli := pbv.NewDBClient(cc)

	return &Driver{
		signature: signature,
		cc:        cc,
		dbCli:     dbCli,
	}, nil
}

func (d *Driver) Get(ctx context.Context, key string) (string, int64, error) {
	res, err := d.dbCli.Get(ctx, &pbv.GetRequest{
		Key: key,
	})
	if err != nil {
		return "", -1, err
	}
	return res.GetValue(), res.GetVersion(), nil
}

func (d *Driver) Set(ctx context.Context, key, value string, version int64) error {
	if _, err := d.dbCli.Set(ctx, &pbv.SetRequest{
		Key:     key,
		Value:   value,
		Version: version,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Close() error {
	return d.cc.Close()
}
