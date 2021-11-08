#!/bin/bash

set -x

if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers>"
        echo -e "\tDefault: 5 containers"
        N=5
fi

IMGNAME="veritas"
PREFIX="veritas"

KAFKA_ADDR="192.168.20.$(($N+1))"

for I in `seq 1 $(($N-1))`; do
	ADDR="192.168.20.$(($I+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; mkdir -p /veritas/data; nohup /bin/veritas-tendermint --dir=/veritas/data --config=/config.toml > veritas-tm-$I.log 2>&1 &"
done
