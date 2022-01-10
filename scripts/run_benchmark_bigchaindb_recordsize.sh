#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-txsizes-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

THREADS=4
TXSIZES="512B 2kB 8kB 32kB"
# TXSIZES="128kB"
IPPREFIX="192.168.30"
ADDRS="http://$IPPREFIX.2:9984,http://$IPPREFIX.3:9984,http://$IPPREFIX.4:9984,http://$IPPREFIX.5:9984"

cd ..
RDIR=`pwd`
cd scripts

for TXSIZE in $TXSIZES; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh
    sleep 5
    python3 $RDIR/BigchainDB/bench.py temp/ycsb_data_$TXSIZE/workloada.dat temp/ycsb_data_$TXSIZE/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-txsize-$TXSIZE.txt
done
./stop_bigchaindb.sh
