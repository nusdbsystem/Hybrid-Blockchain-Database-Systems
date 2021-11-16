#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-txdelay-bigchaindb-$TSTAMP"
mkdir $LOGS

set -x

THREADS=32
TXDELAYS="0 10 100 1000"

cd ..
RDIR=`pwd`
cd scripts

for TXD in $TXDELAYS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh 4 $TXD
    sleep 5
    timeout 600 python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat http://192.168.20.2:9984,http://192.168.20.3:9984,http://192.168.20.4:9984,http://192.168.20.5:9984 $THREADS 2>&1 | tee $LOGSD/bigchaindb-txdelay-$TXD.txt
done
./stop_bigchaindb.sh
