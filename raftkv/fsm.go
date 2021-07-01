package raftkv

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"

	pb "hybrid/proto/raftkv"
)

type fsm struct {
	kv KV

	logger *log.Logger
}

func NewFSM(kv KV) *fsm {
	return &fsm{
		logger: log.New(os.Stderr, "[fsm] ", log.LstdFlags),
		kv:     kv,
	}
}

func (f *fsm) Get(key string) (string, error) {
	v, err := f.kv.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (f *fsm) Apply(l *raft.Log) interface{} {
	var blk pb.Block
	if err := proto.Unmarshal(l.Data, &blk); err != nil {
		f.logger.Fatalf("failed to unmarshal raft log: %v", err)
	}
	keys, values := make([]string, 0), make([]string, 0)
	for _, req := range blk.Reqs {
		keys = append(keys, req.Key)
		values = append(values, req.Value)
	}

	return f.applySet(keys, values)
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.logger.Printf("Generate FSMSnapshot")
	return &fsmSnapshot{
		kv:     f.kv,
		logger: log.New(os.Stderr, "[fsmSnapshot] ", log.LstdFlags),
	}, nil
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	f.logger.Printf("Restore snapshot from FSMSnapshot")

	var (
		readBuf  []byte
		protoBuf *proto.Buffer
		err      error
		keyCount = 0
	)
	// decode message from protobuf
	f.logger.Printf("Read all data")
	if readBuf, err = ioutil.ReadAll(rc); err != nil {
		// read done completely
		f.logger.Printf("Snapshot restore failed: %v", err)
		return err
	}

	protoBuf = proto.NewBuffer(readBuf)

	f.logger.Printf("new protoBuf length %d bytes", len(protoBuf.Bytes()))

	// decode messages from 1M block file
	// the last message could decode failed with io.ErrUnexpectedEOF
	for {
		item := &pb.KVItem{}
		if err = protoBuf.DecodeMessage(item); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			f.logger.Printf("DecodeMessage failed: %v", err)
			return err
		}
		// apply item to store
		f.logger.Printf("Set key %v to %v count: %d", string(item.Key), string(item.Value), keyCount)
		err = f.kv.Set(item.Key, item.Value)
		if err != nil {
			f.logger.Printf("Snapshot load failed: %v", err)
			return err
		}
		keyCount++
	}

	f.logger.Printf("Restore total %d keys", keyCount)

	return nil
}

func (f *fsm) applySet(keys, values []string) interface{} {
	f.logger.Printf("Apply set %s to %s", keys[0], values[0])
	for i := range keys {
		if err := f.kv.Set([]byte(keys[i]), []byte(values[i])); err != nil {
			return err
		}
	}
	return nil
}

func (f *fsm) Close() error {
	f.kv.Close()
	return nil
}
