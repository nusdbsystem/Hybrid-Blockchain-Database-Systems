#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-veritas-kafka-$TSTAMP"
mkdir $LOGS

set -x

../bin/veritas-kafka-bench --load-path=temp/ycsb_data/workloada.dat --run-path=temp/ycsb_data/run_workloada.dat --ndrivers=8 --nthreads=128 --veritas-addrs=192.168.20.2:1990,192.168.20.3:1990,192.168.20.4:1990,192.168.20.5:1990 --tso-addr=:7070 2>&1 | tee $LOGS/veritas.txt
