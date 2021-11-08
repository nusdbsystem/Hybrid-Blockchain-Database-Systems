#!/bin/bash
#
# kill and remove docker containes
#

IMGNAME="bigchaindb"
PREFIX="bigchaindb"

idx=1
for id in `docker ps | grep $PREFIX | cut -d ' ' -f 1`; do
	echo "$PREFIX$idx"
	idx=$(($idx+1))
	docker kill $id
	docker rm $id
done
