#!/bin/bash

. ../env.sh

set -x

THREADS=$DEFAULT_THREADS_BIGCHAINDB
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-nodes-bigchaindb-$TSTAMP"
mkdir $LOGSD

cd ..
RDIR=`pwd`
cd scripts

for N in $NODES; do
    ./restart_cluster_bigchaindb.sh $N
    ./start_bigchaindb.sh $N
    
    ADDRS="http://$IPPREFIX:9984"
    for I in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,http://$IPPREFIX.$I:9984"
    done
    
    python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-nodes-$N.txt    
done
