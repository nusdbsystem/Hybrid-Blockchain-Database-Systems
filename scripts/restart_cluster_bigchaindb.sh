#!/bin/bash

. ./env.sh

N=$DEFAULT_NODES

if [ $# -gt 0 ]; then
	N=$1
else
	echo -e "Usage: $0 <# containers>"
	echo -e "\tDefault: $N containers"
fi

sudo ./unset_ovs_bigchaindb.sh $N
./kill_containers_bigchaindb.sh $N
./start_containers_bigchaindb.sh $N
sudo ./set_ovs_bigchaindb.sh $N