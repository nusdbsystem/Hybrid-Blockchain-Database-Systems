#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-distribution-veritas-tendermint-$TSTAMP"
mkdir $LOGS

N=$DEFAULT_NODES
DRIVERS=$DEFAULT_DRIVERS_VERITAS_TM
THREADS=$DEFAULT_THREADS_VERITAS_TM
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

# Uniform
./restart_cluster_veritas.sh
./start_veritas_tendermint.sh
sleep 5
../bin/veritas-tendermint-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS | tee $LOGS/veritas-uniform.txt

# Latest
./restart_cluster_veritas.sh
./start_veritas_tendermint.sh
sleep 5
../bin/veritas-tendermint-bench --load-path=$DEFAULT_WORKLOAD_PATH"_latest/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_latest/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS | tee $LOGS/veritas-latest.txt

# Zipfian
./restart_cluster_veritas.sh
./start_veritas_tendermint.sh
sleep 5
../bin/veritas-tendermint-bench --load-path=$DEFAULT_WORKLOAD_PATH"_zipfian/"$DEFAULT_WORKLOAD".dat" --run-path=$DEFAULT_WORKLOAD_PATH"_zipfian/run_"$DEFAULT_WORKLOAD".dat" --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS | tee $LOGS/veritas-zipfian.txt

./stop_veritas_tendermint.sh
