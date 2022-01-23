#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-txdelay-veritas-tendermint-$TSTAMP"
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

for TXD in $TXDELAYS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint_delay.sh 4 $TXD
    sleep 30
    ../bin/veritas-tendermint-bench-txdelay --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-tendermint-txdelay-$TXD.txt
done
./stop_veritas_tendermint.sh
