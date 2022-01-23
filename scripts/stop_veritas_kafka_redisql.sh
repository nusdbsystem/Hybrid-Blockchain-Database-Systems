#!/bin/bash

. ./env.sh

set -x

N=$(($DEFAULT_NODES + 1))
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers>"
        echo -e "\tDefault: 5 containers"        
fi

#TSO
killall -9 veritas-tso

# Kafka
KAFKA_ADDR=$IPPREFIX".$(($N+1))"
RES=`kafkacat -b $KAFKA_ADDR:9092 -L 2>&1`
if [[ "$RES" =~ "% ERROR" ]]; then
	echo "Kafka is down."
else
	ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0; bin/kafka-topics.sh --delete --topic shared-log --bootstrap-server 0.0.0.0:9092"
	sleep 30
	ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0; bin/kafka-server-stop.sh"
	sleep 5
fi
ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0; bin/zookeeper-server-stop.sh"

# Nodes
for I in `seq 1 $(($N-1))`; do
	ADDR=$IPPREFIX".$(($I+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "redis-cli flushdb; killall -9 redis-server; killall -9 veritas-kafka-redisql"
done
