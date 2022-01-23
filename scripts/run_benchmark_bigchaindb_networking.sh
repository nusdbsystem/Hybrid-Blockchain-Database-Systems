#!/bin/bash

. ./env.sh

set -x

BWS="NoLimit 10000 1000 100"
RTTS="5ms 10ms 20ms 30ms 40ms 50ms 60ms"

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-networking-bigchaindb-$TSTAMP"
mkdir $LOGS

N=$DEFAULT_NODES
THREADS=$DEFAULT_THREADS_BIGCHAINDB
WORKLOAD_FILE="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
WORKLOAD_RUN_FILE="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

# Generate server addresses. BigchainDB port is 9984
ADDRS="http://$IPPREFIX.2:9984"
for IDX in `seq 3 $(($N+1))`; do
	ADDRS="$ADDRS,http://$IPPREFIX.$IDX:9984"
done

cd ..
RDIR=`pwd`
cd scripts

for BW in $BWS; do    
    for RTT in $RTTS; do
    	LOGSD="$LOGS/logs-$BW-$RTT"
	    mkdir -p $LOGSD
        ./restart_cluster_bigchaindb.sh
        if [[ "$BW" != "NoLimit" ]]; then
            sudo ./set_ovs_bw_limit.sh $BW 1
        fi
	    ./set_tc.sh $RTT
	    sleep 3
        ./start_bigchaindb.sh
	    ./run_iperf_ping.sh 2>&1 | tee $LOGSD/net.txt
	    sleep 3
        python3 $RDIR/BigchainDB/bench.py $WORKLOAD_FILE $WORKLOAD_RUN_FILE $ADDRS $THREADS 2>&1 | tee $LOGS/bigchaindb-$BW-$RTT.txt
    done
done
