#!/bin/bash

cd ..
mkdir -p bin
cd bin

set -x

go version

go build -o preprocess ../cmd/preprocess/main.go
go build -o veritas-kafka ../cmd/veritas/main.go
go build -o veritas-kafka-bench ../veritas/benchmark/ycsbbench/main.go
go build -o veritas-tendermint ../cmd/veritastm/main.go
go build -o veritas-tendermint-bench ../veritastm/benchmark/main.go

exit 1

# DB benchmarks
go build -o db-redis ../cmd/redis/main.go
go build -o db-redisql ../cmd/redisql/main.go
go build -o db-mongodb ../cmd/mongodb/main.go
