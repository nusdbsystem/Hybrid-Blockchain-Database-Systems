#!/bin/bash

if [ $EUID -ne 0 ]; then
	echo "This script must be run as root!"
	exit 1
fi

. ./env.sh

N=$DEFAULT_NODES
if [ $# -ge 1 ]; then
	N=$1
fi
PREFIX="bigchaindb"

ovs-vsctl add-br ovs-br1
ifconfig ovs-br1 $IPPREFIX.1 netmask 255.255.255.0 up
for idx in `seq 1 $N`; do
	idx2=$(($idx+1))
	ovs-docker add-port ovs-br1 eth1 $PREFIX$idx --ipaddress=$IPPREFIX.$idx2/24
done
