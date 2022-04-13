#!/bin/bash

. ./env.sh

set -x

N=$DEFAULT_NODES
BLKSIZE=100
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers> <blk size>"
        echo -e "\tDefault: $N containers"
	echo -e "\tDefault: $BLKSIZE block size"
fi
if [ $# -gt 1 ]; then
	BLKSIZE=$2
else
        echo -e "Usage: $0 <# containers> <blk size>"
        echo -e "\tDefault: $N containers"
        echo -e "\tDefault: $BLKSIZE block size"
fi

# Start
LEADER=""
for I in `seq 1 $N`; do
        ADDR=$IPPREFIX".$(($I+1))"        
        ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; redis-server > redis.log 2>&1 &"
        ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; mkdir veritas-raft; nohup /bin/veritas-raft --svr-addr=$ADDR:1900 --raft-addr=$ADDR:1800 --raft-leader=$LEADER --dir=/veritas-raft --blk-size=$BLKSIZE --redis-addr=0.0.0.0:6379 --redis-db=0 > veritas-raft-$I.log 2>&1 &"        
        if [ $I -eq 1 ]; then
                LEADER="$ADDR:1900"
        fi
        sleep 5
done
