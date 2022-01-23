#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-distribution-bigchain-$TSTAMP"
mkdir $LOGSD

N=$DEFAULT_NODES
THREADS=$DEFAULT_THREADS_BIGCHAINDB

# Generate server addresses. BigchainDB port is 9984
ADDRS="http://$IPPREFIX.2:9984"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,http://$IPPREFIX.$IDX:9984"
done

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

# Uniform
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-uniform.txt
copy_logs $LOGSD/logs-bigchaindb-uniform

# Latest
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH""_latest/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH""_latest/run_$DEFAULT_WORKLOAD".dat

./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-latest.txt
copy_logs $LOGSD/logs-bigchaindb-latest

# Zipfian
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH""_zipfian/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH""_zipfian/run_$DEFAULT_WORKLOAD".dat

./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-zipfian.txt
copy_logs $LOGSD/logs-bigchaindb-zipfian
./stop_bigchaindb.sh
