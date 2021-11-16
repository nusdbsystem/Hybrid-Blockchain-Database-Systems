package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"hybrid/BlockchainDB/storage/redis/redisconn"
)

var _ KV = (*RedisKV)(nil)

type RedisKV struct {
	cli *redis.Client
}

type KV interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Close()
}

func NewRedisKV(addr, pwd string, db int) (*RedisKV, error) {
	cli, err := redisconn.NewRedisConn(addr, pwd, db)
	if err != nil {
		return nil, err
	}
	return &RedisKV{cli: cli}, nil
}

func (r *RedisKV) Get(key []byte) ([]byte, error) {
	val, err := r.cli.Get(context.Background(), string(key)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r *RedisKV) Set(key, value []byte) error {
	return r.cli.Set(context.Background(), string(key), string(value), 0).Err()
}

func (r *RedisKV) Close() {
	r.cli.Close()
}
