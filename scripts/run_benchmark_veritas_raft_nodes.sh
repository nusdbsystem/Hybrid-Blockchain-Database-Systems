#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-nodes-veritas-raft-$TSTAMP"
mkdir $LOGS

DRIVERS=$DEFAULT_DRIVERS_VERITAS_RAFT
THREADS=$DEFAULT_THREADS_VERITAS_RAFT
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

NODES="32"
for N in $NODES; do
    ./restart_cluster_veritas.sh $(($N+1))
    ./start_veritas_raft.sh $N

    # Generate server addresses. Veritas port is 1900
    ADDRS="$IPPREFIX.2:1900"
    for IDX in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,$IPPREFIX.$IDX:1900"
    done
        
    ../bin/veritas-raft-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-nodes-$N.txt

    sleep 10
    SLOGS=$LOGS/veritas-nodes-$N-logs
    mkdir -p $SLOGS
    for I in `seq 2 $(($N+1))`; do
        IDX=$(($I-1))
        scp -o StrictHostKeyChecking=no root@1$IPPREFIX.$I:/veritas-raft-$IDX.log $SLOGS/
    done
done
sudo ./unset_ovs_veritas.sh
./kill_containers_veritas.sh
