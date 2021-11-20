#!/bin/bash
set -ex

nodes=${1:-4}
shard=${2:-1}
IMGNAME="blockchaindb"
PREFIX="blockchaindb"

CPUS_PER_CONTAINER=1

DFILE=dockers.txt
rm -rf $DFILE

for idx in `seq 1 ${nodes}`; do
	CPUID=$(($idx*$CPUS_PER_CONTAINER))
	CPUIDS=$CPUID
	for jdx in `seq 1 $(($CPUS_PER_CONTAINER-1))`; do
		CPUIDS="$CPUIDS,$(($CPUID+$jdx))"
	done
	docker run -d --publish-all=true --cap-add=SYS_ADMIN --cap-add=NET_ADMIN --security-opt seccomp:unconfined --cpuset-cpus=$CPUIDS --name=$PREFIX$idx $IMGNAME tail -f /dev/null 2>&1 >> $DFILE
done
# docker run -d --publish-all=true --cap-add=SYS_ADMIN --cap-add=NET_ADMIN --security-opt seccomp:unconfined --cpuset-cpus=$CPUIDS --name=redis-shard${shard} redis tail -f /dev/null 2>&1 >> $DFILE
while read ID; do
	docker exec $ID "/usr/sbin/sshd"
done < $DFILE