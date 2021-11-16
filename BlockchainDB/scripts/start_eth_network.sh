#!/usr/bin/env bash
set -ex

echo "restart: kill all ethnodes"
# pkill -f "geth" || true
pgeth=`ps -ef|grep "geth"|grep -v "grep"|wc -l`
echo ${pgeth}
if (( ${pgeth} > 0 )); then 
   kill $(ps -ef|grep "geth"|grep -v "grep"|awk '{print $2}')
fi

sleep 2

echo "Start ethereum nodes, Please input shard size(default 1), node size(default 4)"
shards=${1:-1}
nodes=${2:-4}


dir=$(dirname "$0")
echo "##################### generate ethereum genesis config ##########"
${dir}/eth/gen_eth_config.sh ${shards} ${nodes}

echo "##################### init geth nodes using genesis file ##########"
${dir}/eth/init_eth_account.sh ${shards} ${nodes}

echo "##################### start geth bootnode and add peers ##########"
${dir}/eth/start_eth_node.sh ${shards} ${nodes}

echo "##################### deploy KVContract to eth network ##########"
${dir}/eth/deploy_contract.sh ${shards} ${nodes}


echo "##################### Setup ethereum network successfully! ##########"
