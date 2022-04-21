#!/bin/bash

if [ $# -lt 1 ]; then
	echo "Usage: $0 <log dir>"
	exit 1
fi
LOGS=$1

. ./env.sh

for DIST in $DISTRIBUTIONS; do
	FILE="$LOGS/veritas-$DIST.txt"
	TPS=`cat $FILE | grep Throughput | cut -d ' ' -f 12`
	LAT=`cat $FILE | grep latency | cut -d ' ' -f 3`
	ABT=`cat $LOGS/veritas-$DIST-logs/veritas-* | grep Abort | wc -l`
	echo "$TPS;$LAT;$ABT"
done

