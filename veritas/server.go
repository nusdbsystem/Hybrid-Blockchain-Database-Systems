package veritas

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	pbv "hybrid/proto/veritas"
	"hybrid/veritas/keylocker"
	"hybrid/veritas/ledger"
)

type server struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *Config

	l *ledger.LogLedger

	cli *redis.Client

	puller *kafka.Consumer
	pusher *kafka.Producer
	locker keylocker.KLocker

	mu     *sync.RWMutex
	buffer map[int64]*BlockPurpose

	msgCh chan *pbv.SharedLog

	getCache *cache.Cache
}

type BlockPurpose struct {
	blk      *pbv.Block
	approved map[string]struct{}
}

func NewServer(redisCli *redis.Client, consumer *kafka.Consumer, producer *kafka.Producer, ledgerPath string, config *Config) *server {
	ctx, cancel := context.WithCancel(context.Background())
	l, err := ledger.NewLedger(ledgerPath, true)
	if err != nil {
		log.Fatalf("Create ledger failed: %v", err)
	}
	s := &server{
		ctx:      ctx,
		cancel:   cancel,
		l:        l,
		config:   config,
		cli:      redisCli,
		puller:   consumer,
		pusher:   producer,
		locker:   &keylocker.KMutex{},
		mu:       &sync.RWMutex{},
		buffer:   make(map[int64]*BlockPurpose),
		getCache: cache.New(5*time.Minute, 10*time.Minute),
		msgCh:    make(chan *pbv.SharedLog, 10000),
	}
	if err := s.puller.Subscribe(config.Topic, nil); err != nil {
		log.Fatalf("Subscribe topic %v failed: %v", config.Topic, err)
	}
	go s.applyLoop()
	go s.batchLoop()
	return s
}

func (s *server) applyLoop() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		msg, err := s.puller.ReadMessage(-1)
		if err != nil {
			log.Fatalf("Consumer read msg failed: %v", err)
		}
		var blk pbv.Block
		if err := proto.Unmarshal(msg.Value, &blk); err != nil {
			log.Fatalf("Parse msg failed: %v", err)
		}
		switch blk.GetType() {
		case pbv.MessageType_Approve:
			if blk.Signature == s.config.Signature {
				break
			}
			s.mu.RLock()
			_, ok := s.buffer[blk.Txs[0].Seq]
			s.mu.RUnlock()
			if !ok {
				if _, o := s.config.Parties[blk.Signature]; !o {
					break
				}
				s.mu.Lock()
				s.buffer[blk.Txs[0].Seq] = &BlockPurpose{
					blk:      &blk,
					approved: make(map[string]struct{}),
				}
				s.mu.Unlock()
			}
			s.mu.Lock()
			s.buffer[blk.Txs[0].Seq].approved[blk.Signature] = struct{}{}
			s.mu.Unlock()
			s.mu.RLock()
			_, ok = s.buffer[blk.Txs[0].Seq].approved[s.config.Signature]
			s.mu.RUnlock()
			if !ok {
				verifyRes := pbv.MessageType_Approve
			LOOP:
				for _, sl := range blk.Txs {
					for _, t := range sl.Sets {
						res, err := s.cli.Get(s.ctx, t.GetKey()).Result()
						if err == redis.Nil {
							continue
						} else if err != nil {
							log.Fatalf("Commit log %v get failed: %v", blk.Txs[0].GetSeq(), err)
						}
						v, err := Decode(res)
						if err != nil {
							log.Fatalf("Commit log %v decode failed: %v", blk.Txs[0].GetSeq(), err)
						}
						if v.Version >= t.Version {
							verifyRes = pbv.MessageType_Abort
							break LOOP
						}
					}
				}
				go func() {
					approveLog, err := proto.Marshal(&pbv.Block{
						Txs:       blk.Txs,
						Type:      verifyRes,
						Signature: s.config.Signature,
					})
					if err != nil {
						log.Fatalf("Approve log %v failed: %v", blk.Txs[0].GetSeq(), err)
					}
					if err := s.pusher.Produce(&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &s.config.Topic, Partition: kafka.PartitionAny},
						Value:          approveLog,
					}, nil); err != nil {
						log.Fatalf("%s produce approve %v failded: %v", s.config.Signature, blk.Txs[0].Seq, err)
					}
				}()
				if verifyRes == pbv.MessageType_Approve {
					s.mu.Lock()
					s.buffer[blk.Txs[0].Seq].approved[s.config.Signature] = struct{}{}
					s.mu.Unlock()
				} else {
					s.mu.Lock()
					delete(s.buffer, blk.Txs[0].Seq)
					s.mu.Unlock()
					break
				}
			}
			s.mu.RLock()
			approveLen := len(s.buffer[blk.Txs[0].Seq].approved)
			s.mu.RUnlock()
			if approveLen != len(s.config.Parties) {
				break
			}
			var blkBuf *BlockPurpose
			s.mu.RLock()
			blkBuf = s.buffer[blk.Txs[0].Seq]
			s.mu.RUnlock()
			for _, sl := range blkBuf.blk.Txs {
				for _, t := range sl.Sets {
					entry, err := Encode(t.GetValue(), t.GetVersion())
					if err != nil {
						log.Fatalf("Commit log %v encode failed: %v", blkBuf.blk.Txs[0].GetSeq(), err)
					}
					if err := s.cli.Set(s.ctx, t.GetKey(), entry, 0).Err(); err != nil {
						log.Fatalf("Commit log %v redis set failed: %v", blkBuf.blk.Txs[0].GetSeq(), err)
					}
					if _, ok := s.getCache.Get(t.Key); ok {
						s.getCache.Set(t.Key, t.Value, cache.DefaultExpiration)
					}
					if err := s.l.Append([]byte(t.GetKey()), []byte(t.GetValue()+"-"+fmt.Sprintf("%v", t.GetVersion()))); err != nil {
						log.Fatalf("Append to ledger failed: %v", err)
					}
				}
			}
			s.mu.Lock()
			delete(s.buffer, blkBuf.blk.Txs[0].Seq)
			s.mu.Unlock()
		case pbv.MessageType_Abort:
			delete(s.buffer, blk.Txs[0].Seq)
		default:
			log.Fatalf("Invalid shared log type: %v", blk.GetType())
		}
	}
}

func (s *server) batchLoop() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	block := &pbv.Block{
		Txs:       make([]*pbv.SharedLog, 0),
		Type:      pbv.MessageType_Approve,
		Signature: s.config.Signature,
	}
	defer close(s.msgCh)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-t.C:
			if len(block.Txs) > 0 {
				blkLog, err := proto.Marshal(block)
				if err != nil {
					log.Fatalf("Block log failed: %v", err)
				}
				if err := s.pusher.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &s.config.Topic, Partition: kafka.PartitionAny},
					Value:          blkLog,
				}, nil); err != nil {
					log.Fatalf("%s produce block log failded: %v", s.config.Signature, err)
				}
				blockTmp := &pbv.Block{
					Txs:       make([]*pbv.SharedLog, len(block.Txs)),
					Type:      block.Type,
					Signature: block.Signature,
				}
				copy(blockTmp.Txs, block.Txs)
				s.mu.Lock()
				s.buffer[blockTmp.Txs[0].Seq] = &BlockPurpose{
					blk:      blockTmp,
					approved: map[string]struct{}{s.config.Signature: {}},
				}
				s.mu.Unlock()
				block.Txs = make([]*pbv.SharedLog, 0)
			}
		case txn := <-s.msgCh:
			block.Txs = append(block.Txs, txn)
			if len(block.Txs) >= s.config.BlockSize {
				blkLog, err := proto.Marshal(block)
				if err != nil {
					log.Fatalf("Block log failed: %v", err)
				}
				if err := s.pusher.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &s.config.Topic, Partition: kafka.PartitionAny},
					Value:          blkLog,
				}, nil); err != nil {
					log.Fatalf("%s produce block log failded: %v", s.config.Signature, err)
				}
				blockTmp := &pbv.Block{
					Txs:       make([]*pbv.SharedLog, len(block.Txs)),
					Type:      block.Type,
					Signature: block.Signature,
				}
				copy(blockTmp.Txs, block.Txs)
				s.mu.Lock()
				s.buffer[blockTmp.Txs[0].Seq] = &BlockPurpose{
					blk:      blockTmp,
					approved: map[string]struct{}{s.config.Signature: {}},
				}
				s.mu.Unlock()
				block.Txs = make([]*pbv.SharedLog, 0)
			}
		}
	}
}

func (s *server) Get(ctx context.Context, req *pbv.GetRequest) (*pbv.GetResponse, error) {
	if v, ok := s.getCache.Get(req.Key); ok {
		return &pbv.GetResponse{Value: v.(string)}, nil
	}
	res, err := s.cli.Get(ctx, req.GetKey()).Result()
	if err != nil {
		return nil, err
	}
	v, err := Decode(res)
	if err != nil {
		return nil, err
	}

	s.getCache.Set(req.Key, v.Val, cache.DefaultExpiration)

	return &pbv.GetResponse{Value: v.Val}, nil
}

func (s *server) Set(ctx context.Context, req *pbv.SetRequest) (*pbv.SetResponse, error) {
	s.locker.Lock(req.GetKey())
	defer s.locker.Unlock(req.GetKey())

	sets := []*pbv.SetRequest{{
		Signature: req.GetSignature(),
		Key:       req.GetKey(),
		Value:     req.GetValue(),
		Version:   req.GetVersion(),
	}}
	s.msgCh <- &pbv.SharedLog{
		Seq:  req.GetVersion(),
		Sets: sets,
	}
	return &pbv.SetResponse{}, nil
}

func (s *server) Verify(ctx context.Context, req *pbv.VerifyRequest) (*pbv.VerifyResponse, error) {
	proof, err := s.l.ProveKey([]byte(req.GetKey()))
	if err != nil {
		return nil, err
	}
	return &pbv.VerifyResponse{
		RootDigest:            s.l.GetRootDigest(),
		SideNodes:             proof.SideNodes,
		NonMembershipLeafData: proof.NonMembershipLeafData,
	}, nil
}

func (s *server) BatchSet(ctx context.Context, req *pbv.BatchSetRequest) (*pbv.BatchSetResponse, error) {
	for _, r := range req.GetSets() {
		s.locker.Lock(r.GetKey())
	}
	defer func() {
		for _, r := range req.GetSets() {
			s.locker.Unlock(r.GetKey())
		}
	}()

	s.msgCh <- &pbv.SharedLog{
		Seq:  req.GetSets()[0].GetVersion(),
		Sets: req.GetSets(),
	}

	return &pbv.BatchSetResponse{}, nil
}
