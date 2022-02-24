#!/bin/bash
#
# kill and remove docker containes
#

IMGNAME="veritas_hotstuff"
PREFIX="veritas_hotstuff"

idx=1
for id in `docker ps -a| grep $PREFIX | cut -d ' ' -f 1`; do
	echo "$PREFIX$idx"
	idx=$(($idx+1))
	docker kill $id
	docker rm $id
	echo ''
done
