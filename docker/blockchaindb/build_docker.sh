#!/bin/bash

if ! [ -d "../../BlockchainDB/.bin" ]; then
    echo "Please build the binaries first! (cd BlockchainDB && make build)"
    exit 1
fi

rm -rf .bin
cp -r ../../BlockchainDB/.bin .	

if ! [ -f "id_rsa.pub" ]; then
    if ! [ -f "$HOME/.ssh/id_rsa.pub" ]; then
        echo "You do not have a public SSH key. Please generate one! (ssh-keygen)"
        exit 1
    fi
    cp $HOME/.ssh/id_rsa.pub .
fi

docker build -f Dockerfile -t blockchaindb .