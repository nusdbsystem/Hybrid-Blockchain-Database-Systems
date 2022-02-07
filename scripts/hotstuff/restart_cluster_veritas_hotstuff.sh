#!/bin/bash

N=${1:-4}


sudo ./unset_ovs_veritas_hotstuff.sh $N
./kill_containers_veritas_hotstuff.sh
sleep 2
./start_containers_veritas_hotstuff.sh $N
sleep 2
sudo ./set_ovs_veritas_hotstuff.sh $N