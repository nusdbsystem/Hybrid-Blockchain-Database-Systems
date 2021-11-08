#!/bin/bash

if [ $EUID -ne 0 ]; then
	echo "This script must be run as root!"
	exit 1
fi

N=4
if [ $# -gt 0 ]; then
	N=$1
else
	echo -e "Usage: $0 <# containers>"
	echo -e "\tDefault: $N containers"
fi

PREFIX="bigchaindb"

for idx in `seq 1 $N`; do
	idx2=$(($idx+1))
	ovs-docker del-port ovs-br1 eth1 $PREFIX$idx
done
ovs-vsctl del-br ovs-br1