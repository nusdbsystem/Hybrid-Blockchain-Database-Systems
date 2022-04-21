#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-distribution-veritas-raft-$TSTAMP"
mkdir $LOGS

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

function copy_logs {
    LOGSDIR=$1        
    mkdir -p $LOGSDIR
    for I in `seq 2 $(($N+1))`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-raft-$IDX.log $LOGSDIR/
    done
}

# Uniform
./restart_cluster_veritas.sh
./start_veritas_raft.sh
sleep 5
../bin/veritas-raft-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-uniform.txt
copy_logs "$LOGS/veritas-uniform-logs"

# Latest
./restart_cluster_veritas.sh
./start_veritas_raft.sh
sleep 5
../bin/veritas-raft-bench --load-path=$DEFAULT_WORKLOAD_PATH"_latest/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_latest/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-latest.txt
copy_logs "$LOGS/veritas-latest-logs"

# Zipfian
./restart_cluster_veritas.sh
./start_veritas_raft.sh
sleep 5
../bin/veritas-raft-bench --load-path=$DEFAULT_WORKLOAD_PATH"_zipfian/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_zipfian/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-zipfian.txt
copy_logs "$LOGS/veritas-zipfian-logs"

