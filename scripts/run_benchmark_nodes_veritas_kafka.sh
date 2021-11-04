#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-nodes-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

NODES="4 8 16 32 64"
# Veritas (Kafka) uses 64 clients
DRIVERS=4
THREADS=16

for N in $NODES; do
    sudo ./unset_ovs_veritas.sh
    ./kill_containers_veritas.sh
    ./start_containers_veritas.sh $(($N+1))
    sudo ./set_ovs_veritas.sh
    ./start_veritas_kafka.sh $(($N+1))
    
    ADDRS="192.168.20.2:1990"
    for I in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,192.168.20.$I:1990"
    done
    
    ../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-nodes-$N.txt
done
sudo ./unset_ovs_veritas.sh
./kill_containers_veritas.sh