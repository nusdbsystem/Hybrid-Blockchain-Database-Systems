#!/bin/bash

. ./env.sh

set -x

N=$DEFAULT_NODES
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. BigchainDB port is 9984
ADDRS="http://$IPPREFIX.2:9984"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,http://$IPPREFIX.$IDX:9984"
done

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-clients-bigchaindb-$TSTAMP"
mkdir $LOGSD

cd ..
RDIR=`pwd`
cd scripts

function copy_logs {
	DEST=$1
	mkdir -p $DEST
	for IDX in `seq 2 $(($N+1))`; do
		DEST_NODE=$DEST/node-$(($IDX-1))
		mkdir -p $DEST_NODE
		scp root@$IPPREFIX.$IDX:bigchaindb* $DEST_NODE/
		scp root@$IPPREFIX.$IDX:mongodb.log $DEST_NODE/
		scp root@$IPPREFIX.$IDX:tendermint.log $DEST_NODE/
	done
}

# Threads list is defined in env.sh
for TH in $THREADS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh        
    sleep 10
    python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $TH 2>&1 | tee $LOGSD/bigchaindb-clients-$TH.txt
	copy_logs $LOGSD/logs-bigchaindb-clients-$TH
    ./stop_bigchaindb.sh
done
# ./restart_cluster_bigchaindb.sh
