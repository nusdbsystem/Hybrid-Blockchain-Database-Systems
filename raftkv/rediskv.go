package raftkv

import (
	"context"
	"hybrid/dbconn"

	"github.com/go-redis/redis/v8"
)

var _ KV = (*RedisKV)(nil)

type RedisKV struct {
	cli *redis.Client
}

func NewRedisKV(addr, pwd string, db int) (KV, error) {
	cli, err := dbconn.NewRedisConn(addr, pwd, db)
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

func (r *RedisKV) SnapshotItems() <-chan DataItem {
	ch := make(chan DataItem, 1024)

	go func() {
		defer close(ch)
		iter := r.cli.Scan(context.Background(), 0, "*", 0).Iterator()
		for iter.Next(context.Background()) {
			key := iter.Val()
			val, err := r.cli.Get(context.Background(), key).Result()
			if err != nil {
				panic(err)
			}
			kvi := &KVItem{
				key:   append([]byte{}, []byte(key)...),
				value: append([]byte{}, []byte(val)...),
				err:   nil,
			}
			ch <- kvi
		}
		if err := iter.Err(); err != nil {
			panic(err)
		}
		kvi := &KVItem{
			key:   nil,
			value: nil,
			err:   ErrSnapshotFinished,
		}
		ch <- kvi
	}()

	return ch
}

func (r *RedisKV) Close() {
	r.cli.Close()
}
