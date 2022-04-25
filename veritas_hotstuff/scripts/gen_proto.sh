#!/bin/bash

go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=grpc:. ./proto/veritashs/veritashs.proto



