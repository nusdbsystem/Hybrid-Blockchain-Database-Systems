#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-distribution-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8
THREADS=256  # best performance
DISTROS="uniform latest zipfian"

function copy_logs {
    LOGSDIR=$1
        
    mkdir -p $LOGSDIR
    for I in `seq 2 5`; do
	    IDX=$(($I-1))
	    scp -o StrictHostKeyChecking=no root@192.168.20.$I:/veritas-$IDX.log $LOGSDIR/
    done
    scp -o StrictHostKeyChecking=no root@192.168.20.$I:/kafka_2.12-2.7.0/zookeeper.log $LOGSDIR/
    scp -o StrictHostKeyChecking=no root@192.168.20.$I:/kafka_2.12-2.7.0/kafka.log $LOGSDIR/
}

# Uniform
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-uniform.txt
copy_logs "$LOGS/veritas-uniform-logs"

# Latest
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=temp/ycsb_data_latest/workloada.dat --run-path=temp/ycsb_data_latest/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-latest.txt
copy_logs "$LOGS/veritas-latest-logs"

# Zipfian
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=temp/ycsb_data_zipfian/workloada.dat --run-path=temp/ycsb_data_zipfian/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-zipfian.txt
copy_logs "$LOGS/veritas-zipfian-logs"

./stop_veritas_kafka.sh
