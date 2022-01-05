#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-txsizes-veritas-tendermint-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8 
THREADS=256
TXSIZES="512B 2kB 8kB 32kB 128kB"

for TXSIZE in $TXSIZES; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint.sh
    sleep 10
    ../bin/veritas-tendermint-bench --load-path=temp/ycsb_data_$TXSIZE/workloada.dat --run-path=temp/ycsb_data_$TXSIZE/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 2>&1 | tee $LOGS/veritas-tendermint-txsize-$TXSIZE.txt
done
./stop_veritas_tendermint.sh
