#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-clients-veritas-raft-$TSTAMP"
mkdir $LOGS

N=$DEFAULT_NODES
DRIVERS=$DEFAULT_DRIVERS_VERITAS_RAFT
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas Raft port is 1900
ADDRS="$IPPREFIX.2:1900"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1900"
done

for TH in $THREADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_raft.sh
    sleep 10
    ../bin/veritas-raft-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$TH --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-clients-$TH.txt
    # copy logs
    SLOGS="$LOGS/veritas-clients-$TH-logs"
    mkdir -p $SLOGS
    for I in `seq 1 $N`; do
	    IDX=$(($I+1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$IDX:/veritas-raft-$I.log $SLOGS/
    done    
done
./restart_cluster_veritas.sh
