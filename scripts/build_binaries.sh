#!/bin/bash

cd ..
mkdir -p bin
cd bin

set -x

go version

go build -o preprocess ../cmd/preprocess/main.go
go build -o veritas-kafka ../cmd/veritas/main.go
go build -o veritas-kafka-bench ../veritas_kafka/benchmark/ycsbbench/main.go
go build -o veritas-tendermint ../cmd/veritas-tendermint/main.go
go build -o veritas-tendermint-bench ../veritas_tendermint/benchmark/main.go
go build -o veritas-raft ../cmd/veritas-raft/main.go
go build -o veritas-raft-bench ../veritas_raft/benchmark/ycsbbench/main.go

