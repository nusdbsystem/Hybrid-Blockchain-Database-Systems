#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-distribution-veritas-kafka-$TSTAMP"
mkdir $LOGS

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

function copy_logs {
    LOGSDIR=$1        
    mkdir -p $LOGSDIR
    for I in `seq 2 5`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-$IDX.log $LOGSDIR/
    done
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/kafka_2.12-2.7.0/zookeeper.log $LOGSDIR/
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/kafka_2.12-2.7.0/kafka.log $LOGSDIR/
}

# Uniform
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-uniform.txt
copy_logs "$LOGS/veritas-uniform-logs"

# Latest
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH"_latest/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_latest/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-latest.txt
copy_logs "$LOGS/veritas-latest-logs"

# Zipfian
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH"_zipfian/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_zipfian/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-zipfian.txt
copy_logs "$LOGS/veritas-zipfian-logs"

./stop_veritas_kafka.sh
