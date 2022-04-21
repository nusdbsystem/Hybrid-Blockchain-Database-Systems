#!/bin/bash

if [ $# -lt 1 ]; then
	echo "Usage: $0 <log dir>"
	exit 1
fi
LOGS=$1

. ./env.sh
BWS="NoLimit 10000 1000 100"
RTTS="5ms 10ms 20ms 30ms 40ms 50ms 60ms"

for BW in $BWS; do
	for RTT in $RTTS; do
		FILE="$LOGS/veritas-$BW-$RTT.txt"
		TPS=`cat $FILE | grep Throughput | cut -d ' ' -f 12`
		LAT=`cat $FILE | grep latency | cut -d ' ' -f 3`
		echo "$TPS;$LAT"
	done
	echo ";"
done

