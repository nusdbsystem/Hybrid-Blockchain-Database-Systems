#!/bin/bash

N=5
CDIR=`pwd`

if [ $# -gt 0 ]; then
	N=$1
else
	echo -e "Usage: $0 <# containers>"
	echo -e "\tDefault: $N containers"
fi

IMGNAME="veritas"
PREFIX="veritas"

DFILE=dockers.txt
rm -rf $DFILE

for idx in `seq 1 $N`; do
	CPUID=$(($idx+0))
	docker run -v "$CDIR/temp/ycsb_data":/data -d --publish-all=true --cap-add=SYS_ADMIN --cap-add=NET_ADMIN --security-opt seccomp:unconfined --cpuset-cpus=$CPUID --name=$PREFIX$idx $IMGNAME tail -f /dev/null 2>&1 >> $DFILE
done
while read ID; do
	docker exec $ID service ssh start
done < $DFILE