#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-txsizes-bigchaindb-$TSTAMP"
mkdir $LOGS

set -x

THREADS=32
TXSIZES="512B 2kB 8kB 32kB 128kB"

cd ..
RDIR=`pwd`
cd scripts

for TXSIZE in $TXSIZES; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh
    sleep 5
    timeout 600 python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat http://192.168.20.2:9984,http://192.168.20.3:9984,http://192.168.20.4:9984,http://192.168.20.5:9984 $THREADS 2>&1 | tee $LOGSD/bigchaindb-txsize-$TXSIZE.txt
done
./stop_bigchaindb.sh
