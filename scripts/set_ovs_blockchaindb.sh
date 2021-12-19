#!/bin/bash
set -ex

if [ $EUID -ne 0 ]; then
	echo "This script must be run as root!"
	exit 1
fi

N=${1:-4}
shard=${2:-1}
PREFIX="blockchaindb"
NET_PREFIX="192.168.20"

ovs-vsctl add-br ovs-br1
ifconfig ovs-br1 $NET_PREFIX.1  netmask 255.255.255.0 up
for idx in `seq 1 $N`; do
	idx2=$(($idx+1))
	ovs-docker add-port ovs-br1 eth1 $PREFIX$idx --ipaddress=$NET_PREFIX.${idx2}/24
done
# ovs-docker add-port ovs-br1 eth1 redis-shard${shard} --ipaddress=$NET_PREFIX.$((${idx2}+${shard}))/24
