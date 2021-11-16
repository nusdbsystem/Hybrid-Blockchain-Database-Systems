#!/bin/bash
dir=$(dirname "${BASH_SOURCE[0]}")
#go get -d google.golang.org/grpc
#go get -d github.com/golang/protobuf/{proto,protoc-gen-go}

protoc --go_out=plugins=grpc:. ./proto/blockchaindb/blockchaindb.proto

