#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-workload-veritas-kafka-redisql-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8 
THREADS=256
WORKLOADS="workloada workloadb workloadc"

# Redis KV
for W in $WORKLOADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka_redisql.sh
    sleep 30
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/$W.dat --run-path=temp/ycsb_data/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-redisql-$W.txt
    SLOGS="$LOGS/veritas-redisql-$W-logs"
    mkdir -p $SLOGS
    for I in `seq 2 5`; do
        IDX=$(($I-1))
        scp -o StrictHostKeyChecking=no root@192.168.20.$I:/veritas-$IDX.log $SLOGS/
    done
    scp -o StrictHostKeyChecking=no root@192.168.20.6:/kafka_2.12-2.7.0/zookeeper.log $SLOGS/
    scp -o StrictHostKeyChecking=no root@192.168.20.6:/kafka_2.12-2.7.0/kafka.log $SLOGS/
done
./stop_veritas_kafka.sh
