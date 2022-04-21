#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-workload-veritas-raft-$TSTAMP"
mkdir $LOGS

N=$DEFAULT_NODES
DRIVERS=$DEFAULT_DRIVERS_VERITAS_RAFT
THREADS=$DEFAULT_THREADS_VERITAS_RAFT
BLKSIZE=$DEFAULT_BLOCK_SIZE
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1900
ADDRS="$IPPREFIX.2:1900"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1900"
done

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_raft.sh
    sleep 10
    ../bin/veritas-raft-bench --load-path=$DEFAULT_WORKLOAD_PATH/$W.dat --run-path=$DEFAULT_WORKLOAD_PATH/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-raft-$W.txt
    SLOGS="$LOGS/veritas-raft-$W-logs"
    mkdir -p $SLOGS
    for I in `seq 2 $(($N+1))`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-raft-$IDX.log $SLOGS/
    done
done
