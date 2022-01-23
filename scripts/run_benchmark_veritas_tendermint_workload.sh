#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-workload-veritas-tendermint-$TSTAMP"
mkdir $LOGS

N=DEFAULT_NODES
DRIVERS=$DEFAULT_DRIVERS_VERITAS_TM
THREADS=$DEFAULT_THREADS_VERITAS_TM
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint.sh
    sleep 10
    ../bin/veritas-tendermint-bench --load-path=$DEFAULT_WORKLOAD_PATH/$W.dat --run-path=$DEFAULT_WORKLOAD_PATH/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS | tee $LOGS/veritas-workload-$W.txt
done
./stop_veritas_tendermint.sh
