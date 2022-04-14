package raftkv

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"

	pb "hybrid/proto/raftkv"
)

type fsmSnapshot struct {
	kv     KV
	logger *log.Logger
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	f.logger.Printf("Persist action in fsmSnapshot")

	ch := f.kv.SnapshotItems()

	keyCount := 0

	for {
		buff := proto.NewBuffer([]byte{})

		dataItem := <-ch
		item := dataItem.(*KVItem)

		if item.IsFinished() {
			break
		}

		protoKVItem := &pb.KVItem{
			Key:   item.key,
			Value: item.value,
		}

		keyCount++

		buff.EncodeMessage(protoKVItem)

		if _, err := sink.Write(buff.Bytes()); err != nil {
			return err
		}
	}
	f.logger.Printf("Persist total %d keys", keyCount)

	return nil
}

func (f *fsmSnapshot) Release() {
	f.logger.Printf("Release action in fsmSnapshot")
}
