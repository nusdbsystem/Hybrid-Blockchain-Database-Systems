package dbconn

import (
	"github.com/mediocregopher/radix/v3"
)

func NewRedisqlConn(addr string, connNum int) (*radix.Pool, error) {
	return radix.NewPool("tcp", addr, connNum)
}
