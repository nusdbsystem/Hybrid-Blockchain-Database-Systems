package db

import (
	"context"

	pbv "hybrid/proto/veritas"
)

type TransactionDB struct {
	ts        int64
	signature string
	cli       pbv.NodeClient
	setBuffer map[string]string
}

func NewTransaction(ts int64, cli pbv.NodeClient, signature string) *TransactionDB {
	return &TransactionDB{
		ts:        ts,
		signature: signature,
		cli:       cli,
		setBuffer: make(map[string]string),
	}
}

func (db *TransactionDB) Get(key string) (string, error) {
	res, err := db.cli.Get(context.Background(), &pbv.GetRequest{
		Key:       key,
		Signature: db.signature,
	})
	if err != nil {
		return "", err
	}
	return res.GetValue(), nil
}

func (db *TransactionDB) Set(key, value string) error {
	db.setBuffer[key] = value
	return nil
}

func (db *TransactionDB) Commit() error {
	for k, v := range db.setBuffer {
		if _, err := db.cli.Set(context.Background(), &pbv.SetRequest{
			Signature: db.signature,
			Key:       k,
			Value:     v,
			Version:   db.ts,
		}); err != nil {
			return err
		}
	}
	return nil
}
