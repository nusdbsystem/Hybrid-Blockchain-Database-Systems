#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-nodes-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

cd ..
RDIR=`pwd`
cd scripts

NODES="4 8 16 32 64"
THREADS=4

for N in $NODES; do
    ./restart_cluster_bigchaindb.sh $N
    ./start_bigchaindb.sh $N
    
    ADDRS="http://192.168.30.2:9984"
    for I in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,http://192.168.30.$I:9984"
    done
    
    python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-nodes-$N.txt    
done
