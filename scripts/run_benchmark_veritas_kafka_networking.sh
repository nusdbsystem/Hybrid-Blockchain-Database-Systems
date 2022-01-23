#!/bin/bash

. ./env.sh

set -x

BWS="NoLimit 10000 1000 100"
RTTS="5ms 10ms 20ms 30ms 40ms 50ms 60ms"

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-networking-veritas-kafka-$TSTAMP"
mkdir $LOGS

N=$(($DEFAULT_NODES + 1))
DRIVERS=$DEFAULT_DRIVERS_VERITAS_KAFKA
THREADS=$DEFAULT_THREADS_VERITAS_KAFKA
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. Veritas port is 1990
ADDRS="$IPPREFIX.2:1990"
for IDX in `seq 3 $N`; do
	ADDRS="$ADDRS,$IPPREFIX.$IDX:1990"
done

for BW in $BWS; do    
    for RTT in $RTTS; do
	    LOGSD="$LOGS/logs-$BW-$RTT"
	    mkdir $LOGSD
        ./restart_cluster_veritas.sh
        if [[ "$BW" != "NoLimit" ]]; then
            sudo ./set_ovs_bw_limit.sh $BW 1
        fi
	    ./set_tc.sh $RTT
	    sleep 3
        ./start_veritas_kafka.sh
	    ./run_iperf_ping.sh 2>&1 | tee $LOGSD/net.txt
	    sleep 3        
	    ../bin/veritas-kafka-bench --load-path=$WORKLOAD_FILE --run-path=$WORKLOAD_RUN_FILE --ndrivers=$DRIVERS --nthreads=$THREADS --veritas-addrs=$ADDRS --tso-addr=:7070 2>&1 | tee $LOGS/veritas-$BW-$RTT.txt
    done
done
