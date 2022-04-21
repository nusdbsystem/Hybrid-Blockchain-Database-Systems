#!/bin/bash

. ./env.sh

set -x

N=$DEFAULT_NODES
DRIVERS=$DEFAULT_DRIVERS_VERITAS_RAFT
THREADS=$DEFAULT_THREADS_VERITAS_RAFT
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1900
ADDRS="$IPPREFIX.2:1900"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1900"
done

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-blksize-veritas-raft-$TSTAMP"
mkdir $LOGS

for BLK in $BLKSIZES; do
    ./restart_cluster_veritas.sh
    ./start_veritas_raft.sh $N $BLK
    sleep 10
    ../bin/veritas-raft-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-raft-blksize-$BLK.txt
done
