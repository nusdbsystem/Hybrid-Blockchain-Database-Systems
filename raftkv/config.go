package raftkv

type Config struct {
	Id        string
	RaftDir   string
	RaftBind  string
	RaftJoin  string
	BlockSize int
}

func NewConfig(id, raftDir, raftBind, join string, blkSize int) *Config {
	return &Config{
		Id:        id,
		RaftDir:   raftDir,
		RaftBind:  raftBind,
		RaftJoin:  join,
		BlockSize: blkSize,
	}
}
