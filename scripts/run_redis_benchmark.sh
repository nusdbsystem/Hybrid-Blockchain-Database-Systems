#!/bin/bash

go build -o redis-bench ../cmd/redis/main.go
go build -o redisql-bench ../cmd/redisql/main.go

if ! [ -f "RediSQL.so" ]; then
    wget https://github.com/RedBeardLab/rediSQL/releases/download/v1.1.1/RediSQL_v1.1.1_9b110f__release.so
    mv RediSQL_v1.1.1_9b110f__release.so RediSQL.so 
    chmod u+x RediSQL.so
fi
killall -9 redis-server
sleep 5
redis-server --loadmodule `pwd`/RediSQL.so --port 6500 > redis.log 2>&1 &
sleep 5

echo "Redis ..."
echo "FLUSHALL" | redis-cli -p 6500
./redis-bench --load-path temp/ycsb_data/workloada.dat --run-path temp/ycsb_data/run_workloada.dat --nthreads 6 --redis-addr 127.0.0.1:6500 --redis-db 0

echo "RediSQL ..."
echo "FLUSHALL" | redis-cli -p 6500
./redisql-bench --load-path temp/ycsb_data/workloada.dat --run-path temp/ycsb_data/run_workloada.dat --nthreads 6 --redis-addr 127.0.0.1:6500 --redis-db 0
