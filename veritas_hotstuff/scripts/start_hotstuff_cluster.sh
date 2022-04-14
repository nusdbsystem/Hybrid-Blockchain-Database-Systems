#!/usr/bin/env bash
# set -x
# trap 'trap - SIGTERM && kill -- -$$' SIGINT SIGTERM EXIT
#export HOTSTUFF_LOG=debug

# echo "restart: kill all previous hotstuffservers"
# pgrep -f "hotstuffserver"
# pkill -f "hotstuffserver"
# sleep 5

dir=$(pwd)
echo $dir
bin="$dir/hotstuffserver/hotstuffserver"

mkdir -p output
echo "Start hotstuff nodes, Please input server node size(default 4)"
replicaIDs=${1:-4}

for (( c=1; c<=${replicaIDs}; c++ ))
do 
$bin --config="toml.${replicaIDs}/hotstuff"${c} --privkey=keys.${replicaIDs}/r${c}.key --output=output/${c}.out &
echo "hotstuffserver$c start with config file toml.${replicaIDs}/hotstuff$c.toml"
done


echo "#########################################################################"
echo "##################### Start hotstuff nodes successfully! ################"
echo "#########################################################################"

# wait; wait; wait; wait
