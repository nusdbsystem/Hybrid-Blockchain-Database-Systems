#!/bin/bash

cd ..
mkdir -p bin
cd bin

go build -o veritas-kafka ../cmd/veritas/main.go
go build -o veritas-tso ../cmd/tso/main.go
go build -o veritas-kafka-bench ../veritas/benchmark/ycsbbench/main.go