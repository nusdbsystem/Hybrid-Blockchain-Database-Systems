#!/bin/bash
BIGCHAINDB="bigchaindb-2.2.2"
if ! [ -d "$BIGCHAINDB" ]; then
    wget https://github.com/bigchaindb/bigchaindb/archive/refs/tags/v2.2.2.tar.gz
    tar xf v2.2.2.tar.gz
fi
if ! [ -d "$BIGCHAINDB/scripts" ]; then
    cp -r ../../BigchainDB/scripts $BIGCHAINDB/
fi
if ! [ -f "id_rsa.pub" ]; then
    if ! [ -f "$HOME/.ssh/id_rsa.pub" ]; then
        echo "You do not have a public SSH key. Please generate one! (ssh-keygen)"
        exit 1
    fi
    cp $HOME/.ssh/id_rsa.pub .
fi
docker build -f Dockerfile -t bigchaindb .