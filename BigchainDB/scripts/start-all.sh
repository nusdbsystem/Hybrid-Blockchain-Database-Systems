#!/bin/bash

# preconf
# ./ntpset
cp .bigchaindb /root/.bigchaindb
cd /root

# MongoDB
[ "$(stat -c %U /data/db)" = mongodb ] || chown -R mongodb /data/db
nohup mongod --bind_ip_all > mongodb.log 2>&1 &

# Tendermint Init
/usr/local/bin/tendermint init

# BigchainDB
bigchaindb start > /dev/null 2>&1 &
# bigchaindb start --experimental-parallel-validation > /dev/null 2>&1 &

# Tendermint Start
/usr/local/bin/tendermint node --p2p.laddr "tcp://0.0.0.0:26656" --proxy_app="tcp://0.0.0.0:26658" --consensus.create_empty_blocks=false --p2p.pex=false > tendermint.log 2>&1 &
