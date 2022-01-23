#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-workload-veritas-kafka-$TSTAMP"
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

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka.sh
    sleep 30
    ../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH/$W.dat --run-path=$DEFAULT_WORKLOAD_PATH/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-kafka-$W.txt
    SLOGS="$LOGS/veritas-kafka-$W-logs"
    mkdir -p $SLOGS
    for I in `seq 2 5`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-$IDX.log $SLOGS/
    done
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.6:/kafka_2.12-2.7.0/zookeeper.log $SLOGS/
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.6:/kafka_2.12-2.7.0/kafka.log $SLOGS/
done
./stop_veritas_kafka.sh
