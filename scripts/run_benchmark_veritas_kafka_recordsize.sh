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
    ../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH"_$TXSIZE/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_$TXSIZE/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-kafka-txsize-$TXSIZE.txt
    # copy logs
    SLOGS="$LOGS/veritas-txsize-$TXSIZE-logs"
    mkdir -p $SLOGS
    for I in `seq 2 $N`; do
            IDX=$(($I-1))
            scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-$IDX.log $SLOGS/
    done
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$(($N+1)):/kafka_2.12-2.7.0/zookeeper.log $SLOGS/
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$(($N+1)):/kafka_2.12-2.7.0/kafka.log $SLOGS/
done
./stop_veritas_kafka.sh
