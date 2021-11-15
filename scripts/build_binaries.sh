#!/bin/bash

cd ..
mkdir -p bin
cd bin

set -x

go build -o veritas-kafka ../cmd/veritas/main.go
go build -o veritas-tso ../cmd/tso/main.go
go build -o veritas-kafka-bench ../veritas/benchmark/ycsbbench/main.go
go build -o veritas-tendermint ../cmd/veritastm/main.go
go build -o veritas-tendermint-bench ../veritastm/benchmark/main.go

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
go build -o veritas-kafka-bench-txdelay ../veritas/benchmark/ycsbbench/main.go
go build -o veritas-tendermint-txdelay ../cmd/veritastm/main.go
go build -o veritas-tendermint-bench-txdelay ../veritastm/benchmark/main.go
git checkout ../cmd/veritas/main.go ../cmd/veritastm/main.go ../veritas/benchmark/ycsbbench/main.go ../veritas/config.go ../veritas/driver/driver.go ../veritas/server.go ../veritastm/benchmark/main.go ../veritastm/config.go ../veritastm/driver.go ../veritastm/server.go

# build Veritas Tendermint + MongoDB
cd ..
git apply scripts/veritas-tendermint-mongodb.patch
cd bin
go build -o veritas-tendermint-mongodb ../cmd/veritastm-mongodb/main.go
git checkout ../veritastm/ledgerapp.go ../veritastm/server.go
