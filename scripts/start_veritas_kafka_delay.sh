#!/bin/bash

. ./env.sh

set -x

N=$(($DEFAULT_NODES + 1))
BLKSIZE=100
TXDELAY=0
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers> <blk size>"
        echo -e "\tDefault: $N containers"
	echo -e "\tDefault: $BLKSIZE block size"
        echo -e "\tDefault: $TXDELAY block size"
fi
if [ $# -gt 1 ]; then
	BLKSIZE=$2
else
        echo -e "Usage: $0 <# containers> <blk size>"
        echo -e "\tDefault: $N containers"
        echo -e "\tDefault: $BLKSIZE block size"
        echo -e "\tDefault: $TXDELAY block size"
fi
if [ $# -gt 2 ]; then
	TXDELAY=$3
else
        echo -e "Usage: $0 <# containers> <blk size>"
        echo -e "\tDefault: $N containers"
        echo -e "\tDefault: $BLKSIZE block size"
        echo -e "\tDefault: $TXDELAY block size"
fi

#TSO
../bin/veritas-tso --addr=":7070" > tso.log 2>&1 &

# Kafka
KAFKA_ADDR=$IPPREFIX".$(($N+1))"
ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0/config && echo 'advertised.listeners=PLAINTEXT://$KAFKA_ADDR:9092' >> server.properties"
ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0; bin/zookeeper-server-start.sh config/zookeeper.properties > zookeeper.log 2>&1 &"
sleep 10s
ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0; bin/kafka-server-start.sh config/server.properties > kafka.log 2>&1 &"
sleep 10s
ssh -o StrictHostKeyChecking=no root@$KAFKA_ADDR "cd /kafka_2.12-2.7.0; bin/kafka-topics.sh --create --topic shared-log --bootstrap-server 0.0.0.0:9092"
sleep 10

# Nodes
NODES=node1
for I in `seq 1 $(($N-1))`; do
        NODES="$NODES,node$I"
done

# Start
for I in `seq 1 $(($N-1))`; do
	ADDR=$IPPREFIX".$(($I+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; redis-server > redis.log 2>&1 &"
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; nohup /bin/veritas-kafka-txdelay --signature=node$I --parties=${NODES} --tx-delay=$TXDELAY --blk-size=$BLKSIZE --addr=:1990 --kafka-addr=$KAFKA_ADDR:9092 --kafka-group=$I --kafka-topic=shared-log --redis-addr=0.0.0.0:6379 --redis-db=0 --ledger-path=veritas$I > veritas-$I.log 2>&1 &"
done
