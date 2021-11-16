#!/bin/bash
#
# geth-v1.8.23

git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum
git checkout v1.8.23
export GO111MODULE=off
make