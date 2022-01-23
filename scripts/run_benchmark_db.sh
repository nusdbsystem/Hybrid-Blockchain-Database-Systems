#!/bin/bash
#
# Run benchmark (YCSB Workload A) locally on Redis, RediSQL, and MongoDB
#

if ! [ -x redis-server ]; then
	echo "Make sure you have redis-server binary in this folder"
	exit 1
fi
if ! [ -f redisql.so ]; then
        echo "Make sure you have redis-server library in this folder"
        exit 1
fi
if ! [ -x mongod ]; then
        echo "Make sure you have mongod binary in this folder"
        exit 1
fi

# Start servers
./redis-server --port 7777 --loadmodule ./redisql.so > redis.log 2>&1 &
mkdir -p data
rm -r data/*
./mongod --dbpath data > mongodb.log 2>&1 &
sleep 3

# Redis
echo "*** Redis Workload A"
../bin/db-redis --load-path temp/ycsb_data/workloada.dat --run-path temp/ycsb_data/run_workloada.dat --nthreads 6 --redis-addr 127.0.0.1:7777 --redis-db 0
echo "*** Redis Workload B"
../bin/db-redis --load-path temp/ycsb_data/workloadb.dat --run-path temp/ycsb_data/run_workloadb.dat --nthreads 6 --redis-addr 127.0.0.1:7777 --redis-db 0
echo "*** Redis Workload C"
../bin/db-redis --load-path temp/ycsb_data/workloadc.dat --run-path temp/ycsb_data/run_workloadc.dat --nthreads 6 --redis-addr 127.0.0.1:7777 --redis-db 0

# RedisQL
echo ""
echo "*** RediSQL Workload A"
../bin/db-redisql --load-path temp/ycsb_data/workloada.dat --run-path temp/ycsb_data/run_workloada.dat --nthreads 6 --redis-addr 127.0.0.1:7777 --redis-db 0
echo "*** RediSQL Workload B"
../bin/db-redisql --load-path temp/ycsb_data/workloadb.dat --run-path temp/ycsb_data/run_workloadb.dat --nthreads 6 --redis-addr 127.0.0.1:7777 --redis-db 0
echo "*** RediSQL Workload C"
../bin/db-redisql --load-path temp/ycsb_data/workloadc.dat --run-path temp/ycsb_data/run_workloadc.dat --nthreads 6 --redis-addr 127.0.0.1:7777 --redis-db 0

# MongoDB
echo ""
echo "*** MongoDB Workload A"
../bin/db-mongodb --load-path temp/ycsb_data/workloada.dat --run-path temp/ycsb_data/run_workloada.dat --nthreads 6 --mongo-addr 127.0.0.1 --mongo-port 27017
echo "*** MongoDB Workload B"
../bin/db-mongodb --load-path temp/ycsb_data/workloadb.dat --run-path temp/ycsb_data/run_workloadb.dat --nthreads 6 --mongo-addr 127.0.0.1 --mongo-port 27017
echo "*** MongoDB Workload C"
../bin/db-mongodb --load-path temp/ycsb_data/workloadc.dat --run-path temp/ycsb_data/run_workloadc.dat --nthreads 6 --mongo-addr 127.0.0.1 --mongo-port 27017

killall -9 redis-server
killall -9 mongod
