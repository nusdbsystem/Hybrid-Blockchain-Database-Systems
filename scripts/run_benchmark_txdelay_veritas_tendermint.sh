#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-txdelay-veritas-tendermint-$TSTAMP"
mkdir $LOGS

set -x

DRIVERS=8 
THREADS=256
TXDELAYS="0 10 100 1000"

for TXD in $TXDELAYS; do
    ./restart_cluster_veritas.sh
    ./start_veritas_tendermint_delay.sh 4 $TXD
    sleep 30
    ../bin/veritas-tendermint-bench-txdelay --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 2>&1 | tee $LOGS/veritas-tendermint-txdelay-$TXD.txt
done
./stop_veritas_tendermint.sh
