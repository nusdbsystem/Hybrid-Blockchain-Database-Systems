#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-clients-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=4   # 4 nodes, 4 drivers
THREADS="1 2 4 8 16 32 48 64"

for TH in $THREADS; do
    ./stop_veritas_kafka.sh
    ./start_veritas_kafka.sh
    CLIENTS=$(($DRIVERS * $TH))
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$TH --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-clients-$CLIENTS.txt
done
./stop_veritas_kafka.sh
