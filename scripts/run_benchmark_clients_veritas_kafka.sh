#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-clients-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8
THREADS="4 8 16 32 64 128 192 256"

for TH in $THREADS; do
    ./restart_cluster.sh
    ./start_veritas_kafka.sh    
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$TH --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-clients-$TH.txt
    # copy logs
    SLOGS="$LOGS/veritas-clients-$TH-logs"
    mkdir -p $SLOGS
    for I in `seq 2 5`; do
	IDX=$(($I-1))
	scp -o StrictHostKeyChecking=no root@192.168.20.$I:/veritas-$IDX.log $SLOGS/
    done
    scp -o StrictHostKeyChecking=no root@192.168.20.$I:/kafka_2.12-2.7.0/zookeeper.log $SLOGS/
    scp -o StrictHostKeyChecking=no root@192.168.20.$I:/kafka_2.12-2.7.0/kafka.log $SLOGS/
done
./restart_cluster.sh
