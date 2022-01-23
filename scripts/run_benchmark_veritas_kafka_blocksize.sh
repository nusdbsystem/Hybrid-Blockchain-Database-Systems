#!/bin/bash

. ./env.sh

set -x

N=$(($DEFAULT_NODES + 1))
DRIVERS=$DEFAULT_DRIVERS_VERITAS_KAFKA
THREADS=$DEFAULT_THREADS_VERITAS_KAFKA
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $N`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-blksize-veritas-kafka-$TSTAMP"
mkdir $LOGS

for BLK in $BLKSIZES; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka.sh 5 $BLK
    sleep 30
    ../bin/veritas-kafka-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-kafka-blksize-$BLK.txt
done
./stop_veritas_kafka.sh
