package tso

const (
	GET uint8 = iota
	REPLY
)

type GetTS struct {
	Num int32
}

type ReplyTS struct {
	Timestamp int64
}

type LogTS struct {
	crc uint32
	ts  int64
}
