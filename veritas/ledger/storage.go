package ledger

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"

	"hybrid/veritas/ledger/merkletree"
)

type LogLedger struct {
	store *badger.DB
	smt   *merkletree.SparseMerkleTree
}

func NewLedger(ledgerPath string, withMerkleTree bool) (*LogLedger, error) {
	db, err := badger.Open(badger.DefaultOptions(ledgerPath))
	if err != nil {
		return nil, err
	}
	if !withMerkleTree {
		return &LogLedger{
			store: db,
			smt:   nil,
		}, nil
	}
	tree := merkletree.NewSparseMerkleTree(merkletree.NewSimpleMap(), sha256.New())
	l := &LogLedger{
		store: db,
		smt:   tree,
	}
	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if err := item.Value(func(v []byte) error {
				if _, err := tree.Update(item.Key(), v); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *LogLedger) Append(key, val []byte) error {
	if l.smt != nil {
		if _, err := l.smt.Update(key, val); err != nil {
			return fmt.Errorf("append to merkle tree: %w", err)
		}
	}
	if err := l.store.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, val); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if l.smt != nil {
			l.smt.Delete(key)
		}
		return err
	}
	return nil
}

func (l *LogLedger) GetRootDigest() []byte {
	if l.smt != nil {
		return l.smt.Root()
	}
	return []byte{}
}

func (l *LogLedger) ProveKey(key []byte) (merkletree.SparseMerkleProof, error) {
	if l.smt != nil {
		return l.smt.Prove(key)
	}
	return merkletree.SparseMerkleProof{}, errors.New("no merkle tree")
}
