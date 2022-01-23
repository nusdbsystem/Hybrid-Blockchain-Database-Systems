#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-workload-bigchaindb-$TSTAMP"
mkdir $LOGSD

N=$DEFAULT_NODES
THREADS=$DEFAULT_THREADS_BIGCHAINDB
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. BigchainDB port is 9984
ADDRS="http://$IPPREFIX.2:9984"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,http://$IPPREFIX.$IDX:9984"
done

cd ..
RDIR=`pwd`
cd scripts

for W in $WORKLOADS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh
    sleep 5
    WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$W".dat
    WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$W".dat
    python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-workload-$W.txt
done
./stop_bigchaindb.sh
