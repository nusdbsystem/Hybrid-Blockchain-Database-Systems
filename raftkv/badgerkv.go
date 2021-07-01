package raftkv

import (
	"log"
	"os"

	"github.com/dgraph-io/badger/v3"
)

var _ KV = (*BadgerKV)(nil)

type BadgerKV struct {
	db     *badger.DB
	logger *log.Logger
}

type KVItem struct {
	key   []byte
	value []byte
	err   error
}

func (i *KVItem) IsFinished() bool {
	return i.err == ErrSnapshotFinished
}

func NewBadgerKV(dir string) (KV, error) {
	opt := badger.DefaultOptions(dir)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &BadgerKV{
		db:     db,
		logger: log.New(os.Stderr, "[kv_badger] ", log.LstdFlags),
	}, nil
}

func (b *BadgerKV) Get(key []byte) ([]byte, error) {
	var value []byte

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

func (b *BadgerKV) Set(key, value []byte) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		if err1 := txn.Set(key, value); err1 != nil {
			return err1
		}
		return nil
	})
	return err
}

func (b *BadgerKV) SnapshotItems() <-chan DataItem {
	// create a new channel
	ch := make(chan DataItem, 1024)

	// generate items from snapshot to channel
	go b.db.View(func(txn *badger.Txn) error {
		defer close(ch)
		opt := badger.DefaultIteratorOptions
		opt.PrefetchSize = 10
		it := txn.NewIterator(opt)
		defer it.Close()

		keyCount := 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)

			kvi := &KVItem{
				key:   append([]byte{}, k...),
				value: append([]byte{}, v...),
				err:   err,
			}

			// write kvItem to channel with last error
			ch <- kvi
			keyCount++

			if err != nil {
				return err
			}
		}

		// just use nil kvItem to mark the end
		kvi := &KVItem{
			key:   nil,
			value: nil,
			err:   ErrSnapshotFinished,
		}
		ch <- kvi

		b.logger.Printf("Snapshot total %d keys", keyCount)

		return nil
	})
	// return channel to persist
	return ch
}

func (b *BadgerKV) Close() {
	b.db.Close()
}
