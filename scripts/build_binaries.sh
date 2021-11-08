#!/bin/bash

cd ..
mkdir -p bin
cd bin

set -x

go build -o veritas-kafka ../cmd/veritas/main.go
go build -o veritas-tso ../cmd/tso/main.go
go build -o veritas-kafka-bench ../veritas/benchmark/ycsbbench/main.go
go build -o veritas-tendermint ../cmd/veritastendermint/main.go
go build -o veritas-tendermint-bench ../veritastendermint/benchmark/ycsb.go

# build Veritas Kafka + ZooKeeper TSO
cd ..
git apply scripts/veritas-tso-zk.patch
cd bin
go build -o veritas-kafka-zk ../cmd/veritas/main.go
go build -o veritas-kafka-zk-bench ../veritas/benchmark/ycsbbench/main.go
git checkout ../veritas/driver/driver.go

# build Veritas Kafka + RediSQL
cd ..
git apply scripts/veritas-redisql.patch
cd bin
go build -o veritas-kafka-redisql ../cmd/veritas-redisql/main.go
git checkout ../veritas/server.go

# build Veritas Kafka + Transaction Delay
cd ..
git apply scripts/veritas-txdelay.patch
cd bin
go build -o veritas-kafka-txdelay ../cmd/veritas/main.go
git checkout ../cmd/veritas/main.go ../veritas/config.go ../veritas/server.go