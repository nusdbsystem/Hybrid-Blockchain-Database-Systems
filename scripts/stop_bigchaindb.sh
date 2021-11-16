#!/bin/bash

./multi_node.sh "killall -9 tendermint; killall -9 bigchaindb; killall -9 bigchaindb_ws; killall -9 bigchaindb_exchange; killall -9 mongod; killall -9 gunicorn"
./multi_node.sh "rm -rf /data/db/*; rm -rf /root/.tendermint"