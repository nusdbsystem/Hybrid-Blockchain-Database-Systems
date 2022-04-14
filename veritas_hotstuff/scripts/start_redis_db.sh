#!/usr/bin/env bash
#set -x



replicaIDs=${1:-4}

#docker pull redis:latest

docker rm -f $(sudo -S docker ps -aq  --filter ancestor=redis)

for (( c=1; c<=${replicaIDs}; c++ ))
do 
docker run -itd --name hs${c}-redis -p $((30070 + ${c})):6379 redis
echo "redis db start with port $((30070 + ${c}))"
done

echo "#########################################################################"
echo "##################### Start redis dbs successfully! #####################"
echo "#########################################################################"
