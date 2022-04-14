package merkletree

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

type KVStore interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	Close() error
}

type InvalidKeyError struct {
	Key []byte
}

func (e *InvalidKeyError) Error() string {
	return fmt.Sprintf("invalid key: %s", e.Key)
}

type SimpleMap struct {
	m map[string][]byte
}

func NewSimpleMap() *SimpleMap {
	return &SimpleMap{
		m: make(map[string][]byte),
	}
}

func (sm *SimpleMap) Get(key []byte) ([]byte, error) {
	if value, ok := sm.m[string(key)]; ok {
		return value, nil
	}
	return nil, &InvalidKeyError{Key: key}
}

func (sm *SimpleMap) Set(key []byte, value []byte) error {
	sm.m[string(key)] = value
	return nil
}

func (sm *SimpleMap) Close() error {
	sm.m = nil
	runtime.GC()
	return nil
}

type BadgerStore struct {
	db *badger.DB
}

// NewBadgerStore creates a new empty BadgerStore.
func NewBadgerStore(path string) (*BadgerStore, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &BadgerStore{
		db: db,
	}, nil
}

// Get gets the value for a key.
func (bs *BadgerStore) Get(key []byte) ([]byte, error) {
	var value []byte

	if err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if errors.Is(err, badger.ErrKeyNotFound) {
			return &InvalidKeyError{Key: key}
		} else if err != nil {
			return err
		}
		if err := item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return value, nil
}

// Set updates the value for a key.
func (bs *BadgerStore) Set(key []byte, value []byte) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, value); err != nil {
			return err
		}
		return nil
	})
}

// Close closes BadgerStore
func (bs *BadgerStore) Close() error {
	return bs.db.Close()
}
