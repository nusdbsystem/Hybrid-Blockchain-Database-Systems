#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-distribution-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8
THREADS=256  # best performance
DISTROS="uniform latest zipfian"

# Uniform
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-uniform.txt

# Latest
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=temp/ycsb_data_latest/workloada.dat --run-path=temp/ycsb_data_latest/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-latest.txt

# Zipfian
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
../bin/veritas-kafka-bench --load-path=temp/ycsb_data_zipfian/workloada.dat --run-path=temp/ycsb_data_zipfian/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas-zipfian.txt

./stop_veritas_kafka.sh
