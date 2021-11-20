#!/bin/bash
#
# kill and remove docker containes
#

IMGNAME="blockchaindb"
PREFIX="blockchaindb"

idx=1
for id in `docker ps -a| grep $PREFIX | cut -d ' ' -f 1`; do
	echo "$PREFIX$idx"
	idx=$(($idx+1))
	docker kill $id
	docker rm $id
	echo ''
done

# idx=1
# for id in `docker ps | grep "redis-shard" | cut -d ' ' -f 1`; do
# 	echo "redis-shard$idx"
# 	idx=$(($idx+1))
# 	docker kill $id
# 	docker rm $id
# 	echo ''
# done
