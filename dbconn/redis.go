package dbconn

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisConn(addr, pwd string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	return rdb, nil
}
