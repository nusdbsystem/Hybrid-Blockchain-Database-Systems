#!/bin/bash

./multi_node.sh "killall -9 tendermint; killall -9 veritas-tendermint"
./multi_node.sh "rm -r /veritas/data/*"