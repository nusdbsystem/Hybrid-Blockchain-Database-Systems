package raftkv

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pb "hybrid/proto/raftkv"
	"hybrid/veritas/ledger"
)

var _ pb.DBServer = (*Peer)(nil)

type Peer struct {
	ctx context.Context

	raft *raft.Raft // The consensus mechanism
	fsm  *fsm

	msgCh chan *SetMessage

	config *Config
	logger *log.Logger

	l *ledger.LogLedger
}

type SetMessage struct {
	req *pb.SetRequest
	ch  chan error
}

func NewPeer(ctx context.Context, ledgerPath string, kv KV, config *Config) (*Peer, error) {
	l, err := ledger.NewLedger(ledgerPath, true)
	if err != nil {
		return nil, err
	}
	p := &Peer{
		ctx:    ctx,
		logger: log.New(os.Stderr, "[store] ", log.LstdFlags),
		fsm:    NewFSM(kv),
		config: config,
		l:      l,
		msgCh:  make(chan *SetMessage, 10*config.BlockSize),
	}
	if err := p.open(); err != nil {
		return nil, err
	}

	go p.applyLoop()

	return p, nil
}

func (p *Peer) open() error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(p.config.Id)
	config.SnapshotThreshold = 1024

	addr, err := net.ResolveTCPAddr("tcp", p.config.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(p.config.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	ss, err := raft.NewFileSnapshotStore(p.config.RaftDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	// boltDB implement log store and stable store interface
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(p.config.RaftDir, fmt.Sprintf("raft.db-%s", p.config.Id)))
	if err != nil {
		return err
	}

	// raft system
	r, err := raft.NewRaft(config, p.fsm, boltDB, boltDB, ss, transport)
	if err != nil {
		return err
	}
	p.raft = r

	if p.config.RaftJoin == "" {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		p.raft.BootstrapCluster(configuration)
	} else {
		if err := func() error {
			cc, err := grpc.Dial(p.config.RaftJoin, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer cc.Close()
			cli := pb.NewDBClient(cc)
			if _, err := cli.Join(context.Background(), &pb.JoinRequest{
				PeerAddr: p.config.RaftBind,
				PeerId:   p.config.Id}); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Peer) applyLoop() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	setMsgBuffer := make([]*SetMessage, 0)
	defer close(p.msgCh)
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-t.C:
			if len(setMsgBuffer) > 0 {
				go func(buf []*SetMessage) {
					block := &pb.Block{Reqs: make([]*pb.SetRequest, 0)}
					for _, msg := range buf {
						block.Reqs = append(block.Reqs, msg.req)
					}
					blkLog, err := proto.Marshal(block)
					if err != nil {
						log.Fatalf("Block log failed: %v", err)
					}
					f := p.raft.Apply(blkLog, 10*time.Second)
					if err := f.Error(); err != nil {
						log.Fatalf("Apply log failed: %v", err)
					}
					for _, msg := range buf {
						msg.ch <- nil
						p.l.Append([]byte(msg.req.Key), []byte(msg.req.Value))
					}
				}(setMsgBuffer)
				setMsgBuffer = make([]*SetMessage, 0)
			}
		case setMsg := <-p.msgCh:
			setMsgBuffer = append(setMsgBuffer, setMsg)
			if len(setMsgBuffer) >= p.config.BlockSize {
				go func(buf []*SetMessage) {
					block := &pb.Block{Reqs: make([]*pb.SetRequest, 0)}
					for _, msg := range buf {
						block.Reqs = append(block.Reqs, msg.req)
					}
					blkLog, err := proto.Marshal(block)
					if err != nil {
						log.Fatalf("Block log failed: %v", err)
					}
					f := p.raft.Apply(blkLog, 10*time.Second)
					if err := f.Error(); err != nil {
						log.Fatalf("Apply log failed: %v", err)
					}
					for _, msg := range buf {
						msg.ch <- nil
					}
				}(setMsgBuffer)
				setMsgBuffer = make([]*SetMessage, 0)
			}
		}
	}
}

func (p *Peer) isLeader() bool {
	return p.raft.State() == raft.Leader
}

func (p *Peer) IsLeader(ctx context.Context, req *pb.IsLeaderRequest) (*pb.IsLeaderResponse, error) {
	return &pb.IsLeaderResponse{IsLeader: p.isLeader()}, nil
}

func (p *Peer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	if !p.isLeader() {
		return nil, raft.ErrNotLeader
	}

	val, err := p.fsm.Get(req.GetKey())
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{Value: val}, nil
}

func (p *Peer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	if !p.isLeader() {
		return nil, raft.ErrNotLeader
	}

	ch := make(chan error, 1)
	p.msgCh <- &SetMessage{
		req: &pb.SetRequest{
			Key:   req.GetKey(),
			Value: req.GetValue(),
		},
		ch: ch,
	}

	if err := <-ch; err != nil {
		return nil, err
	}

	return &pb.SetResponse{}, nil
}

func (p *Peer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	proof, err := p.l.ProveKey([]byte(req.GetKey()))
	if err != nil {
		return nil, err
	}
	return &pb.VerifyResponse{
		RootDigest:            p.l.GetRootDigest(),
		SideNodes:             proof.SideNodes,
		NonMembershipLeafData: proof.NonMembershipLeafData,
	}, nil
}

func (p *Peer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	if !p.isLeader() {
		return nil, raft.ErrNotLeader
	}

	p.logger.Printf("received join request for remote node %s, addr %s", req.GetPeerId(), req.GetPeerAddr())

	cf := p.raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		p.logger.Printf("get raft configuration failed: %v", err)
		return nil, err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(req.GetPeerId()) {
			p.logger.Printf("node %s already joined raft cluster", req.GetPeerId())
			return &pb.JoinResponse{}, nil
		}
	}

	f := p.raft.AddVoter(raft.ServerID(req.GetPeerId()), raft.ServerAddress(req.GetPeerAddr()), 0, 0)
	if err := f.Error(); err != nil {
		return nil, err
	}

	p.logger.Printf("node %s at %s joined successfully", req.GetPeerId(), req.GetPeerAddr())

	return &pb.JoinResponse{}, nil
}

func (p *Peer) Leave(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	if !p.isLeader() {
		return nil, raft.ErrNotLeader
	}

	p.logger.Printf("received leave request for remote node %s", req.GetPeerId())

	cf := p.raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		p.logger.Printf("get raft configuration failed: %v", err)
		return nil, err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(req.GetPeerId()) {
			f := p.raft.RemoveServer(server.ID, 0, 0)
			if err := f.Error(); err != nil {
				p.logger.Printf("remove server %s failed", req.GetPeerId())
				return nil, err
			}
			p.logger.Printf("node %s leaved successfully", req.GetPeerId())
			return &pb.LeaveResponse{}, nil
		}
	}

	p.logger.Printf("node %s not exists in raft group", req.GetPeerId())

	return &pb.LeaveResponse{}, nil
}

func (p *Peer) Snapshot() error {
	p.logger.Printf("doing snapshot manually")
	f := p.raft.Snapshot()
	return f.Error()
}
