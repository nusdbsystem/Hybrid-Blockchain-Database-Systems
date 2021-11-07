#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-database-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8    # 8 drivers
THREADS=256  # best performance
# WORKLOADS="workloada workloadb workloadc"
WORKLOADS="workloada"
PREFIX="192.168.20."

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster.sh
    ./start_veritas_kafka.sh
    sleep 10
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/$W.dat --run-path=temp/ycsb_data/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-rediskv-$W.txt
    mkdir -p $LOGS/logs-rediskv-$W
    for IDX in `seq 2 5`; do
	    IDXX=$(($IDX-1))
	    scp -o StrictHostKeyChecking=no root@$PREFIX$IDX:/veritas-$IDXX.log $LOGS/logs-rediskv-$W/
    done
    IDX=6
    scp -o StrictHostKeyChecking=no root@$PREFIX$IDX:/kafka_2.12-2.7.0/kafka.log $LOGS/logs-rediskv-$W/
done
./stop_veritas_kafka.sh

# Redis SQL
for W in $WORKLOADS; do
    ./restart_cluster.sh
    ./start_veritas_kafka_redisql.sh
    sleep 10
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/$W.dat --run-path=temp/ycsb_data/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-redisql-$W.txt
    mkdir -p $LOGS/logs-redisql-$W
    for IDX in `seq 2 5`; do
            IDXX=$(($IDX-1))
            scp -o StrictHostKeyChecking=no root@$PREFIX$IDX:/veritas-$IDXX.log $LOGS/logs-redisql-$W/
    done
    IDX=6
    scp -o StrictHostKeyChecking=no root@$PREFIX$IDX:/kafka_2.12-2.7.0/kafka.log $LOGS/logs-redisql-$W/
done
./stop_veritas_kafka_redisql.sh
