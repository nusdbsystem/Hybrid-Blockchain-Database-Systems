#!/bin/bash

./multi_node.sh "killall -9 tendermint; killall -9 bigchaindb; killall -9 bigchaindb_ws; killall -9 bigchaindb_exchange; killall -9 mongod"
./multi_node.sh "rm -r /data/db/*"