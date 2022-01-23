#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-clients-veritas-kafka-tso-zk-$TSTAMP"
mkdir $LOGS

N=$(($DEFAULT_NODES + 1))
DRIVERS=$DEFAULT_DRIVERS_VERITAS_KAFKA
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $N`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done
TSO_ADDR="$IPPREFIX.6:2181"

for TH in $THREADS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka_tso_zk.sh    
    ../bin/veritas-kafka-zk-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$TH --veritas-addrs=$ADDRS --tso-addr=$TSO_ADDR 2>&1 | tee $LOGS/veritas-clients-$TH.txt
    # copy logs
    SLOGS="$LOGS/veritas-clients-$TH-logs"
    mkdir -p $SLOGS
    for I in `seq 2 5`; do
        IDX=$(($I-1))
        scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-$IDX.log $SLOGS/
    done
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/kafka_2.12-2.7.0/zookeeper.log $SLOGS/
    scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/kafka_2.12-2.7.0/kafka.log $SLOGS/
done
./restart_cluster_veritas.sh
