package raftkv

import "errors"

var ErrSnapshotFinished = errors.New("snapshot finished successfully")

type KV interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error

	SnapshotItems() <-chan DataItem

	Close()
}

type DataItem interface{}
