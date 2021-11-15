#!/bin/bash

./multi_node.sh "killall -9 redis-server; killall -9 mongod; killall -9 tendermint; killall -9 veritas-tendermint"
./multi_node.sh "rm -r /veritas/data/*"