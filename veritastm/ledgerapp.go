package veritastm

import (
	"bytes"
	"context"
	"fmt"
	"hybrid/veritas/ledger"
	"log"

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
	parts := bytes.Split(req.Tx, []byte("="))
	key, val := parts[0], parts[1]
	l.ledger.Append(key, val)
	err := l.db.Set(context.Background(), string(key), string(val), 0).Err()
	if err != nil {
		fmt.Printf("Error in Set DeliverTx: %v\n", err)
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
