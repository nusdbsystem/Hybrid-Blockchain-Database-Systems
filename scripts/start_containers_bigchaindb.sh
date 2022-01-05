#!/bin/bash

N=4
if [ $# -gt 0 ]; then
	N=$1
else
	echo -e "Usage: $0 <# containers>"
	echo -e "\tDefault: $N containers"
fi

IMGNAME="bigchaindb"
PREFIX="bigchaindb"

CPUS_PER_CONTAINER=1

DFILE=dockers.txt
rm -rf $DFILE

cat /dev/null > $HOME/.ssh/known_hosts

set -x

for idx in `seq 1 $N`; do
	CPUID=$(($idx*$CPUS_PER_CONTAINER+1))
	#CPUID=$(($idx*$CPUS_PER_CONTAINER))
	CPUIDS=$CPUID
	for jdx in `seq 1 $(($CPUS_PER_CONTAINER-1))`; do
		CPUIDS="$CPUIDS,$(($CPUID+$jdx))"
	done
	numactl --localalloc docker run -d --publish-all=true --cap-add=SYS_ADMIN --cap-add=NET_ADMIN --security-opt seccomp:unconfined --cpuset-cpus=$CPUIDS --name=$PREFIX$idx $IMGNAME tail -f /dev/null 2>&1 >> $DFILE
done
while read ID; do
	docker exec $ID "/usr/sbin/sshd"
done < $DFILE
