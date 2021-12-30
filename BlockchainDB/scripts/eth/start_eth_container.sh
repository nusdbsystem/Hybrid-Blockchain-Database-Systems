#!/usr/bin/env bash
#set -x



replicaIDs=${1:-1}
shardID=${2:-1}
# docker pull ethereum/client-go:v1.8.23

docker rm -f $(sudo -S docker ps -aq  --filter ancestor=ethereum/client-go )

for (( c=1; c<=${replicaIDs}; c++ ))
do
docker run -itd --name geth${c}-shard${shardID} -p $((20070 + ${c})):8545 ethereum/client-go 
echo "geth start with port $((20070 + ${c}))"
done

echo "#########################################################################"
