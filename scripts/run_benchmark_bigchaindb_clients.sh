#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-clients-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

IPPREFIX="192.168.30"
# THREADS="4 8 16 32 64 128 192 256"
THREADS="32"

cd ..
RDIR=`pwd`
cd scripts

for TH in $THREADS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh        
    sleep 10
    # timeout 600 python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat http://$IPPREFIX.2:9984,http://$IPPREFIX.3:9984,http://$IPPREFIX.4:9984,http://$IPPREFIX.5:9984 $TH 2>&1 | tee $LOGSD/bigchaindb-clients-$TH.txt
    python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloadc.dat temp/ycsb_data/run_workloadc.dat http://$IPPREFIX.2:9984,http://$IPPREFIX.3:9984,http://$IPPREFIX.4:9984,http://$IPPREFIX.5:9984 $TH 2>&1 
    # | tee $LOGSD/bigchaindb-clients-$TH.txt
done
# ./restart_cluster_bigchaindb.sh
