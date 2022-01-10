#!/bin/bash

BWS="NoLimit 10000 1000 100"
RTTS="5ms 10ms 20ms 30ms 40ms 50ms 60ms"

TSTAMP=`date +%F-%H-%M-%S`
LOGS="logs-networking-bigchaindb-$TSTAMP"
mkdir $LOGS

THREADS=4
ADDRS="http://192.168.30.2:9984,http://192.168.30.3:9984,http://192.168.30.4:9984,http://192.168.30.5:9984"

cd ..
RDIR=`pwd`
cd scripts

set -x

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
        python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGS/bigchaindb-$BW-$RTT.txt
    done
done
