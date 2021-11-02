#!/bin/bash

MAIN_DIR=`pwd`
KAFKA_DIR="temp/kafka"
TOPIC="veritas"

cd $KAFKA_DIR
# Start Zookeeper
./bin/zookeeper-server-start.sh config/zookeeper.properties > zookeeper.log 2>&1 &
# Start Kafka
./bin/kafka-server-start.sh config/server.properties > kafka.log 2>&1 &
# Create topic
./bin/kafka-topics.sh --create --topic $TOPIC --bootstrap-server localhost:9092
# Check topic
./bin/kafka-topics.sh --describe --topic $TOPIC --bootstrap-server localhost:9092
# Redis
docker run -itd --name redis_container -p 6379:6379 redis
# TSO
go run cmd/tso/main.go --addr=":7070" > tso.log 2>&1 &



