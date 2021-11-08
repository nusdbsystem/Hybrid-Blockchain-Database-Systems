#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-clients-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

THREADS="4 8 16 32 64 128 192 256"

cd ..
RDIR=`pwd`
cd scripts

for TH in $THREADS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh        
    timeout 600 python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat http://192.168.20.2:9984,http://192.168.20.3:9984,http://192.168.20.4:9984,http://192.168.20.5:9984 $TH 2>&1 | tee $LOGSD/bigchaindb-clients-$TH.txt
done
./restart_cluster.sh
