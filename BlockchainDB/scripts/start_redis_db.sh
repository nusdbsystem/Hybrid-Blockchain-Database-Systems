#!/usr/bin/env bash
#set -x


shardIDs=${1:-1}

docker ps |grep shard
# docker rm -f $(sudo -S docker ps -aq  --filter ancestor=redis)
idx=1
for id in `docker ps | grep "redis-shard" | cut -d ' ' -f 1`; do
	echo "redis-shard$idx"
	idx=$(($idx+1))
	docker kill $id
	docker rm $id
done


for (( c=1; c<=${shardIDs}; c++ ))
do
docker run -itd --name shard${c}-redis -p $((60000 + ${c})):6379 redis
echo "redis db start with port $((60000 + ${c}))"
done
docker ps |grep shard
echo "##################### Start redis dbs successfully! #####################"