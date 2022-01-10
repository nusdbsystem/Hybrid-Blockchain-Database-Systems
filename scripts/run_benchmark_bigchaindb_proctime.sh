#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-txdelay-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

THREADS=4
TXDELAYS="0 10 100 1000"
ADDRS="http://192.168.30.2:9984,http://192.168.30.3:9984,http://192.168.30.4:9984,http://192.168.30.5:9984"

cd ..
RDIR=`pwd`
cd scripts

for TXD in $TXDELAYS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh 4 $TXD
    sleep 5
    python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-txdelay-$TXD.txt
done
./stop_bigchaindb.sh
