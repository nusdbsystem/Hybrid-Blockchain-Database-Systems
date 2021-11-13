#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-workload-veritas-tendermint-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8 
THREADS=256
WORKLOADS="workloada workloadb workloadc"

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint.sh
    sleep 10
    ../bin/veritas-tendermint-bench --load-path=temp/ycsb_data/$W.dat --run-path=temp/ycsb_data/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 | tee $LOGS/veritas-workload-$W.txt
done
./stop_veritas_tendermint.sh
