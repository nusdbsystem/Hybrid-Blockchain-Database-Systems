#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-nodes-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

NODES="4 8 16 32 64"
DRIVERS=8
THREADS=256

for N in $NODES; do
    ./restart_cluster_veritas.sh $(($N+1))
    ./start_veritas_kafka.sh $(($N+1))
    
    ADDRS="192.168.20.2:1990"
    for I in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,192.168.20.$I:1990"
    done
    
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-nodes-$N.txt

    sleep 10
    SLOGS=$LOGS/veritas-nodes-$N-logs
    mkdir -p $SLOGS
    for I in `seq 2 $(($N+1))`; do
        IDX=$(($I-1))
        scp -o StrictHostKeyChecking=no root@192.168.20.$I:/veritas-$IDX.log $SLOGS/
    done
    KAFKA_HOST="192.168.20.$(($N+2))"
    scp -o StrictHostKeyChecking=no root@$KAFKA_HOST:/kafka_2.12-2.7.0/zookeeper.log $SLOGS/
    scp -o StrictHostKeyChecking=no root@$KAFKA_HOST:/kafka_2.12-2.7.0/kafka.log $SLOGS/
    ssh -o StrictHostKeyChecking=no root@$KAFKA_HOST "cd /kafka_2.12-2.7.0 && ./bin/kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list" | tee -a $SLOGS/kafka-counters.log
    for I in `seq 1 $N`; do
	    ssh -o StrictHostKeyChecking=no root@$KAFKA_HOST "cd /kafka_2.12-2.7.0 && ./bin/kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group $I" | tee -a $SLOGS/kafka-counters.log
    done
done
sudo ./unset_ovs_veritas.sh
./kill_containers_veritas.sh
