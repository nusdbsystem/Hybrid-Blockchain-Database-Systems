#!/bin/bash

. ./env.sh

set -x

N=$DEFAULT_NODES
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# servers>"
        echo -e "\tDefault: $N servers" 
fi

# Nodes
for I in `seq 1 $N`; do
	ADDR=$IPPREFIX".$(($I+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "redis-cli flushdb; killall -9 redis-server; killall -9 veritas-raft"
done
