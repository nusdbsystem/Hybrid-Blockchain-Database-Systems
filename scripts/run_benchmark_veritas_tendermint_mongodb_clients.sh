#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-clients-veritas-tendermint-mongodb-$TSTAMP"
mkdir $LOGS

N=$DEFAULT_NODES
DRIVERS=$DEFAULT_DRIVERS_VERITAS_TM
THREADS=$DEFAULT_THREADS_VERITAS_TM
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

for TH in $THREADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint_mongodb.sh
    ../bin/veritas-tendermint-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$TH --veritas-addrs=$ADDRS | tee $LOGS/veritas-clients-$TH.txt
    # copy logs
    SLOGS="$LOGS/veritas-clients-$TH-logs"
    mkdir -p $SLOGS
    for I in `seq 2 $(($N+1))`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-tm-$IDX.log $SLOGS/
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/root/tendermint.log $SLOGS/tendermint-$IDX.log
    done    
done
#./restart_cluster_veritas.sh
