#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-blksize-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8 
THREADS=256
BLKSIZES="10 100 1000 10000"

for BLK in $BLKSIZES; do
    ./restart_cluster_veritas.sh
    ./start_veritas_kafka.sh 5 $BLK
    sleep 30
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-kafka-blksize-$BLK.txt
done
./stop_veritas_kafka.sh
