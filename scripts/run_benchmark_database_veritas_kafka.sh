#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-database-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=4   # 4 nodes, 4 drivers
THREADS=16  # best performance
WORKLOADS="workloada workloadb workloadc"

# Redis KV
for W in $WORKLOADS; do
    ./stop_veritas_kafka.sh
    ./start_veritas_kafka.sh
    sleep 30
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/$W.dat --run-path=temp/ycsb_data/run_$W.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-rediskv-$W.txt
done
./stop_veritas_kafka.sh

# Redis SQL