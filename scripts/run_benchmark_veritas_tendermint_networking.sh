#!/bin/bash

. ./env.sh

set -x

BWS="NoLimit 10000 1000 100"
RTTS="5ms 10ms 20ms 30ms 40ms 50ms 60ms"

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-networking-veritas-tendermint-$TSTAMP"
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

for BW in $BWS; do    
    for RTT in $RTTS; do
	LOGSD="$LOGS/logs-$BW-$RTT"
	mkdir $LOGSD
	./restart_cluster_veritas.sh
	if [[ "$BW" != "NoLimit" ]]; then
    	sudo ./set_ovs_bs_limit.sh $BW 1
	fi
	./set_tc.sh $RTT
	sleep 3
    ./start_veritas_tendermint.sh
	./run_iperf_ping.sh 2>&1 | tee $LOGSD/net.txt
	sleep 3
	../bin/veritas-tendermint-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS 2>&1 | tee $LOGS/veritas-net-$BW-$RTT.txt
    done
done
