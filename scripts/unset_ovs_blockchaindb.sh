#!/bin/bash

if [ $EUID -ne 0 ]; then
	echo "This script must be run as root!"
	exit 1
fi

N=${1:-4}
shard=${2:-1}
PREFIX="blockchaindb"

for idx in `seq 1 $N`; do
	idx2=$(($idx+1))
	ovs-docker del-port ovs-br1 eth1 $PREFIX$idx
done
# ovs-docker del-port ovs-br1 eth1 redis-shard${shard}
ovs-vsctl del-br ovs-br1