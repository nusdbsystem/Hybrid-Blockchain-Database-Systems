#!/usr/bin/env bash
set -ex

# echo "restart: kill all previous bcdbnode"
# pgrep -f "bcdbnode"
# pkill -f "bcdbnode"
# kill -9 $(ps -ef|grep "geth"|grep -v "grep"|awk '{print $2}')
# sleep 5
#dir=$(dirname "$0")

echo "Start blockchaindb server nodes, Please input server node size(default 4)"
shardIDs=${1:-1}
replicaIDs=${2:-4}
dir=$(dirname "$0")
echo $dir

cd `dirname ${BASH_SOURCE-$0}`
. eth/env.sh
cd -

bin="${ETH_BIN}/bcdbnode"
tomlDir="config/config.nodes.${shardIDs}.${replicaIDs}"
mkdir -p nodelog
if [ ! -f ${bin} ]; then
    echo "Binary file ${bin} not found!"
    echo "Hint: "
    echo " Please build binaries by run command: make build "
    echo "exit 1 "
    exit 1
fi
for (( i=1; i<=${shardIDs}; i++ ))
do
    for (( c=1; c<=$replicaIDs; c++ ))
    do 
    $bin --config="${tomlDir}/config_${i}_${c}" > nodelog/node.${i}.${c}.log 2>&1 &
    echo "bcdbnode$c start with config file config.nodes.${shardIDs}.${replicaIDs}/config_${i}_${c}.toml"
    done
done
echo "#########################################################################"
echo "##################### Start blockchaindb server nodes successfully! ##########"
echo "#########################################################################"

