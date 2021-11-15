#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-distribution-bigchain-$TSTAMP"
mkdir $LOGS

set -x

THREADS=32
DISTROS="uniform latest zipfian"

# Uniform
./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat http://192.168.20.2:9984,http://192.168.20.3:9984,http://192.168.20.4:9984,http://192.168.20.5:9984 $THREADS 2>&1 | tee $LOGSD/bigchaindb-uniform.txt

# Latest
./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py temp/ycsb_data_latest/workloada.dat temp/ycsb_data_latest/run_workloada.dat http://192.168.20.2:9984,http://192.168.20.3:9984,http://192.168.20.4:9984,http://192.168.20.5:9984 $THREADS 2>&1 | tee $LOGSD/bigchaindb-latest.txt

# Zipfian
./restart_cluster_veritas.sh
./start_veritas_kafka.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py temp/ycsb_data_zipfian/workloada.dat temp/ycsb_data_zipfian/run_workloada.dat http://192.168.20.2:9984,http://192.168.20.3:9984,http://192.168.20.4:9984,http://192.168.20.5:9984 $THREADS 2>&1 | tee $LOGSD/bigchaindb-zipfian.txt

./stop_bigchaindb.sh
