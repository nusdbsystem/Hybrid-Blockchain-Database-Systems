#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-txsizes-veritas-kafka-$TSTAMP"
mkdir $LOGS

N=$(($DEFAULT_NODES + 1))
DRIVERS=$DEFAULT_DRIVERS_VERITAS_KAFKA
THREADS=$DEFAULT_THREADS_VERITAS_KAFKA
BLKSIZE=$DEFAULT_BLOCK_SIZE
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $N`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

for TXSIZE in $TXSIZES; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka.sh
    sleep 30
    ../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH"_$TXSIZE/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_$TXSIZE/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-kafka-txsize-$TXSIZE.txt
done
./stop_veritas_kafka.sh
