#!/bin/bash

if [ $# -lt 1 ]; then
	echo "Usage: $0 <logs dir>"
	exit 1
fi
LOGS=$1

CLIENTS="4 8 16 32 64 128 192 256"
NODES=4

for CLI in $CLIENTS; do 
	AVG=0
	for IDX in `seq 1 $NODES`; do 
		SUM=`cat $LOGS/logs-bigchaindb-clients-$CLI/node-$IDX/tendermint.log | grep validTxs | tr -s ' ' | cut -d ' ' -f 6 | cut -d '=' -f 2 | tr '\n' '+'`
	       	N=`cat $LOGS/logs-bigchaindb-clients-$CLI/node-$IDX/tendermint.log | grep validTxs | wc -l`
	       	AVG=`echo "$AVG+("$SUM"0)/$N" | bc -l`
	done
	THR=`cat $LOGS/bigchaindb-clients-$CLI.txt | grep Throughput | cut -d ' ' -f 5`
	LAT=`cat $LOGS/bigchaindb-clients-$CLI.txt | grep Latency | cut -d ' ' -f 2`
	AVGBLKSIZE=`echo "$AVG/4.0" | bc -l`
	echo "$CLI;$THR;$LAT;$AVGBLKSIZE"
done
