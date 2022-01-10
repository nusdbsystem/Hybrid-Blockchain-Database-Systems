#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-workload-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

THREADS=4
WORKLOADS="workloada workloadb workloadc"
IPPREFIX="192.168.30"

cd ..
RDIR=`pwd`
cd scripts

for W in $WORKLOADS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh
    sleep 5
    python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/$W.dat temp/ycsb_data/run_$W.dat http://$IPPREFIX.2:9984,http://$IPPREFIX.3:9984,http://$IPPREFIX.4:9984,http://$IPPREFIX.5:9984 $THREADS 2>&1 | tee $LOGSD/bigchaindb-workload-$W.txt
done
./stop_bigchaindb.sh
