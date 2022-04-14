package veritastm

import (
	"bytes"
	"context"
	"fmt"
	"hybrid/veritas_kafka/ledger"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type LedgerApp struct {
	ledger *ledger.LogLedger
	db     *redis.Client
}

var _ abcitypes.Application = (*LedgerApp)(nil)

func NewLedgerApp(config *Config, rcli *redis.Client) *LedgerApp {
	// create ledger
	l, err := ledger.NewLedger(config.LedgerPath, true)
	if err != nil {
		log.Fatalf("Create ledger failed: %v", err)
	}
	return &LedgerApp{
		ledger: l,
		db:     rcli,
	}
}

func (LedgerApp) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{
		Data: "Hello from Veritas + Tendermint",
	}
}

func (l *LedgerApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	// first, check MVCC
	parts := bytes.Split(req.Tx, []byte("="))
	vparts := bytes.Split(parts[1], []byte("#"))
	key := parts[0]
	val := vparts[0]
	ver, err := strconv.ParseInt(string(vparts[1]), 10, 64)
	if err != nil {
		fmt.Printf("Error in Set DeliverTx version parseInt: %v\n", err)
		return abcitypes.ResponseDeliverTx{Code: 1}
	}
	res, err := l.db.Get(context.Background(), string(key)).Result()
	var localVer int64
	localVer = 0
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("Error in Set DeliverTx DB Get: %v\n", err)
			return abcitypes.ResponseDeliverTx{Code: 2}
		}
	} else {
		vres, err := Decode(res)
		if err != nil {
			fmt.Printf("Error in Set DeliverTx Decode: %v\n", err)
			return abcitypes.ResponseDeliverTx{Code: 3}
		}
		localVer = vres.Version
	}
	if localVer > ver {
		fmt.Printf("Abort tx due to MVCC\n")
		return abcitypes.ResponseDeliverTx{Code: 4}
	}
	// second, append to ledger and set to DB
	entry, err := Encode(string(val), ver+1)
	if err != nil {
		fmt.Printf("Error in Set DeliverTx Encode: %v\n", err)
		return abcitypes.ResponseDeliverTx{Code: 5}
	}
	l.ledger.Append(key, []byte(entry))
	err = l.db.Set(context.Background(), string(key), string(entry), 0).Err()
	if err != nil {
		fmt.Printf("Error in Set DeliverTx DB Set: %v\n", err)
		return abcitypes.ResponseDeliverTx{Code: 6}
	}
	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (LedgerApp) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	return abcitypes.ResponseCheckTx{Code: 0}
}

func (LedgerApp) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{}
}

func (LedgerApp) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {
	return abcitypes.ResponseQuery{Code: 0}
}

func (LedgerApp) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (LedgerApp) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

func (LedgerApp) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (LedgerApp) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (LedgerApp) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (LedgerApp) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (LedgerApp) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
