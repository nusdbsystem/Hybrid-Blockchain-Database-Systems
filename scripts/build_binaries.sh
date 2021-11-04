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
