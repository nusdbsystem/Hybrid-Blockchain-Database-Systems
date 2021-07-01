package main

import (
	"bytes"

	"github.com/dgraph-io/badger/v3"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type KVStore struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

var _ abcitypes.Application = (*KVStore)(nil)

func NewKVStore(db *badger.DB) *KVStore {
	return &KVStore{
		db: db,
	}
}

func (KVStore) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (KVStore) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (kv *KVStore) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := kv.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	parts := bytes.Split(req.Tx, []byte("="))
	key, value := parts[0], parts[1]

	if err := kv.currentBatch.Set(key, value); err != nil {
		panic(err)
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (kv *KVStore) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := kv.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

func (kv *KVStore) Commit() abcitypes.ResponseCommit {
	kv.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (kv *KVStore) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data
	if err := kv.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err == badger.ErrKeyNotFound {
			resQuery.Log = "key does not exist"
		} else if err != nil {
			return err
		} else {
			return item.Value(func(val []byte) error {
				resQuery.Log = "key exists"
				resQuery.Value = val
				return nil
			})
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return
}

func (KVStore) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (kv *KVStore) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	kv.currentBatch = kv.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (KVStore) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (KVStore) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (KVStore) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (KVStore) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (KVStore) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}

func (kv *KVStore) isValid(tx []byte) (code uint32) {
	// Check format.
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}

	key, value := parts[0], parts[1]

	// Check if the same key=value already exists.
	if err := kv.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return item.Value(func(val []byte) error {
				if bytes.Equal(val, value) {
					code = 2
				}
				return nil
			})
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return code
}
