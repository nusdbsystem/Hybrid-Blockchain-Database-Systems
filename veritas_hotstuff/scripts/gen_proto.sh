#!/bin/bash

# go get -u google.golang.org/grpc
# go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
# protoc --go_out=plugins=grpc:. ./proto/veritas/veritas.proto

go get -u github.com/golang/protobuf/protoc-gen-go 
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
protoc --go_out=paths=source_relative:. --go-grpc_out=. ./proto/veritashs.proto


