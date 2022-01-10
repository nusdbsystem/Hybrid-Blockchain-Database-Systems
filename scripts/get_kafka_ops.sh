#!/bin/bash

NODES="4 8 16 32 64"

if [ $# -lt 1 ]; then
	echo "Usage: $0 <logs-nodes-veritas-kafka...>"
	exit 1
fi

# set -x

LOGSD=$1
echo "# Server Nodes	Reads	Writes"
for N in $NODES; do
	# LOG-END-OFFSET
	WRITES=`cat $LOGSD/veritas-nodes-$N-logs/kafka-counters.log | grep shared-log | tr -s ' ' | cut -d ' ' -f 5`
	# CURRENT-OFFSET
	READS=`cat $LOGSD/veritas-nodes-$N-logs/kafka-counters.log | grep shared-log | tr -s ' ' | cut -d ' ' -f 4`
	NW=`cat $LOGSD/veritas-nodes-$N-logs/kafka-counters.log | grep shared-log | tr -s ' ' | cut -d ' ' -f 5 | wc -l`
	NR=`cat $LOGSD/veritas-nodes-$N-logs/kafka-counters.log | grep shared-log | tr -s ' ' | cut -d ' ' -f 4 | wc -l`
	if [[ $N -ne $NW ]] || [[ $N -ne $NR ]]; then
		echo "Invalid number of counter records"
	fi
	SUMR=`echo $READS | tr ' ' '+' | bc -l`
	SUMW=`echo $WRITES | tr ' ' '+'`
	W=`echo "scale=2;($SUMW+0)/$NW" | bc -l`
	echo "$N		$SUMR	$W"
done



