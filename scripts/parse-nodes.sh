#!/bin/bash

if [ $# -lt 1 ]; then
	echo "Usage: $0 <log dir>"
	exit 1
fi
LOGS=$1

. ./env.sh

for N in $NODES; do
	FILE="$LOGS/veritas-nodes-$N.txt"
	TPS=`cat $FILE | grep Throughput | cut -d ' ' -f 12`
	LAT=`cat $FILE | grep latency | cut -d ' ' -f 3`
	echo "$TPS;$LAT"
done

