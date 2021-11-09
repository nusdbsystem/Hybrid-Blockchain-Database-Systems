#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-clients-veritas-tendermint-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=4
THREADS="4 8 16 32 64 128 192 256"

for TH in $THREADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint.sh    
    ../bin/veritas-tendermint-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --nthreads=$TH --urls=http://192.168.20.2:26656,http://192.168.20.3:26656,http://192.168.20.4:26656,http://192.168.20.5:26656 2>&1 | tee $LOGS/veritas-clients-$TH.txt
    # copy logs
    SLOGS="$LOGS/veritas-clients-$TH-logs"
    mkdir -p $SLOGS
    for I in `seq 2 5`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@192.168.20.$I:/veritas-tm-$IDX.log $SLOGS/
    done    
done
./restart_cluster_veritas.sh
