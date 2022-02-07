#!/usr/bin/env bash
# set -x
# trap 'trap - SIGTERM && kill -- -$$' SIGINT SIGTERM EXIT
#export HOTSTUFF_LOG=debug

# echo "restart: kill all previous veritasnodes"
# pgrep -f "veritasnode"
# pkill -f "veritasnode"
# sleep 5

dir=$(pwd)
bin="$dir/cmd/veritas/veritasnode"

echo "Start veritas server nodes, Please input server node size(default 4)"
replicaIDs=${1:-4}

for (( c=1; c<=$replicaIDs; c++ ))
do 
$bin --config="toml.${replicaIDs}/hotstuff"${c} &
echo "veritasnode$c start with config file toml.${replicaIDs}/hotstuff$c.toml"
done

echo "#########################################################################"
echo "##################### Start veritas server nodes successfully! ##########"
echo "#########################################################################"
# wait; wait; wait; wait
