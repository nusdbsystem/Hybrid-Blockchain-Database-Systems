#!/bin/bash

set -x

if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers>"
        echo -e "\tDefault: 5 containers"
        N=5
fi

../bin/veritas-tso --addr=":7070" > tso.log 2>&1 &

KAFKA_ADDR="192.168.20.$(($N+1))"
for I in `seq 1 $(($N-1))`; do
	ADDR="192.168.20.$(($I+1))"        
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; redis-server > redis.log 2>&1 &"
	sleep 3
	ssh -o StrictHostKeyChecking=no root@$ADDR "mkdir -p /veritas/data; cd /bin; nohup ./veritas-kafka --addr=:1993 --dir=/veritas/data --config=/veritas/config.toml --redis-addr=0.0.0.0:6379 --redis-db=0 > veritas-$I.log 2>&1 &"
done
