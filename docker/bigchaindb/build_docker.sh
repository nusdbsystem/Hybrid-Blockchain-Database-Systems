#!/bin/bash
BIGCHAINDB="bigchaindb-2.2.2"
if ! [ -d "$BIGCHAINDB" ]; then
    wget https://github.com/bigchaindb/bigchaindb/archive/refs/tags/v2.2.2.tar.gz
    tar xf v2.2.2.tar.gz
fi
docker build -f Dockerfile -t bigchaindb .