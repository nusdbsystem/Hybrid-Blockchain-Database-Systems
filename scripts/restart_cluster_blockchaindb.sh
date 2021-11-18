#!/bin/bash

N=${1:-4}


sudo ./unset_ovs_blockchaindb.sh $N
./kill_containers_blockchaindb.sh $N
./start_containers_blockchaindb.sh $N
sudo ./set_ovs_blockchaindb.sh $N