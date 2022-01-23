#!/bin/bash

. ./env.sh

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-database-veritas-kafka-$TSTAMP"
mkdir $LOGS

N=$(($DEFAULT_NODES + 1))
DRIVERS=$DEFAULT_DRIVERS_VERITAS_KAFKA
THREADS=$DEFAULT_THREADS_VERITAS_KAFKA

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $N`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka.sh
    sleep 10
    ../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH/$W.dat --run-path=$DEFAULT_WORKLOAD_PATH/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-rediskv-$W.txt
    mkdir -p $LOGS/logs-rediskv-$W
    for IDX in `seq 2 5`; do
	    IDXX=$(($IDX-1))
	    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$IDX:/veritas-$IDXX.log $LOGS/logs-rediskv-$W/
    done
    IDX=6
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$IDX:/kafka_2.12-2.7.0/kafka.log $LOGS/logs-rediskv-$W/
done
./stop_veritas_kafka.sh

# Redis SQL
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka_redisql.sh
    sleep 10
    ../bin/veritas-kafka-bench --load-path=$DEFAULT_WORKLOAD_PATH/$W.dat --run-path=$DEFAULT_WORKLOAD_PATH/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-redisql-$W.txt
    mkdir -p $LOGS/logs-redisql-$W
    for IDX in `seq 2 5`; do
            IDXX=$(($IDX-1))
            scp -o StrictHostKeyChecking=no root@$IPPREFIX.$IDX:/veritas-$IDXX.log $LOGS/logs-redisql-$W/
    done
    IDX=6
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$IDX:/kafka_2.12-2.7.0/kafka.log $LOGS/logs-redisql-$W/
done
./stop_veritas_kafka_redisql.sh
