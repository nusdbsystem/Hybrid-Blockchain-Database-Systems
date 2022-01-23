#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-nodes-veritas-tendermint-$TSTAMP"
mkdir $LOGS

DRIVERS=$DEFAULT_DRIVERS_VERITAS_TM
THREADS=$DEFAULT_THREADS_VERITAS_TM
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

for N in $NODES; do
    ./restart_cluster_veritas.sh $(($N+1))
    ./start_veritas_tendermint.sh $N
    
    ADDRS="$IPPREFIX.2:1990"
    for I in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,$IPPREFIX.$I:1990"
    done
   
    ../bin/veritas-tendermint-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS | tee $LOGS/veritas-nodes-$N.txt 
done
sudo ./unset_ovs_veritas.sh
./kill_containers_veritas.sh
