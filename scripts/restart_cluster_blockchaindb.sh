#!/bin/bash

N=${1:-4}


sudo ./unset_ovs_blockchaindb.sh $N
./kill_containers_blockchaindb.sh
sleep 2
./start_containers_blockchaindb.sh $N
sleep 2
sudo ./set_ovs_blockchaindb.sh $N