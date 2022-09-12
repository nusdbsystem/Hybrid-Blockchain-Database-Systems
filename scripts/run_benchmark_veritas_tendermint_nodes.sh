#!/bin/bash

. ./env.sh

set -x

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-nodes-veritas-tendermint-$TSTAMP"
mkdir $LOGS

DRIVERS=$DEFAULT_DRIVERS_VERITAS_TM
THREADS=$DEFAULT_THREADS_VERITAS_TM
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

for N in $NODES; do
    ./restart_cluster_veritas.sh $(($N+1))
    ./start_veritas_tendermint.sh $N
    
    ADDRS="$IPPREFIX.2:1990"
    for I in `seq 3 $(($N+1))`; do
        ADDRS="$ADDRS,$IPPREFIX.$I:1990"
    done

    for I in `seq 2 $(($N+1))`; do
        IDX=$(($I-1))
        ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$I "rm -f /dstat-$IDX.csv; dstat --noheaders --nocolor --output /dstat-$IDX.csv > /dev/null 2>&1 &"
    done
   
    ../bin/veritas-tendermint-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS | tee $LOGS/veritas-nodes-$N.txt

    for I in `seq 2 $(($N+1))`; do
        IDX=$(($I-1))
        ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$I "killall -SIGINT dstat"
        scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/veritas-$IDX.log $SLOGS/
        scp -o StrictHostKeyChecking=no root@$IPPREFIX.$I:/dstat-$IDX.csv $SLOGS/
    done

done
sudo ./unset_ovs_veritas.sh
./kill_containers_veritas.sh
