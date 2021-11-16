#!/bin/bash
set -ex

shardID=${1:-1}
nodes=${2:-4}
#bootnode
nodeID=1

cd `dirname ${BASH_SOURCE-$0}`
. env.sh

#${ETH_BIN}/geth --datadir=${ETH_DATA}_${nodeID} --rpc --rpcport "8000" --syncmode "full" --cache 4096 --gasprice 0 --networkid 10001 --mine --minerthreads 1 --unlock 0 console 2> ${ETH_DATA}_${nodeID}/geth.log
#--password <(echo -n "") js <(echo 'console.log(admin.nodeInfo.enode);') 
#--nodiscover 
#--targetgaslimit '67219750000000'

# ${ETH_BIN}/geth --datadir=${ETH_DATA}_${shardID}_${nodeID}  \
# --rpc --rpcport "$((9000 + ${nodeID} + 1000*${shardID}))" \
# --port "$((30303 + ${nodeID} + 1000*(${shardID}-1)))" \
# -networkid $((1000 + ${shardID})) \
# --syncmode "full" --cache 4096 --gasprice 0 -\
# --mine --minerthreads 1 \
# --unlock 0 --password <(echo -n "") 2> ${ETH_DATA}_${shardID}_${nodeID}/eth.log &

# pkill -f "geth" || true
# kill $(ps -ef|grep "geth"|grep -v "grep"|awk '{print $2}') || true
pgeth=`ps -ef|grep "geth"|grep -v "grep"|wc -l`
echo ${pgeth}
if (( ${pgeth} > 0 )); then 
   kill $(ps -ef|grep "geth"|grep -v "grep"|awk '{print $2}')
fi
sleep 2
# start bootnode
# --miner.gaslimit 67219750000000
# --netrestrict --gcmode 'archive'
${ETH_BIN}/geth --datadir=${ETH_DATA}_${shardID}_${nodeID}  \
--rpc --rpcaddr 'localhost' --rpcport "$((9000 + ${nodeID} + 1000*${shardID}))" \
--port "$((30303 + ${nodeID} + 1000*(${shardID}-1)))" \
--gasprice 0 --targetgaslimit 10000000 --mine --minerthreads 1 --unlock 0 --password <(echo -n "") \
--syncmode 'full' \
--nat extip:127.0.0.1 \
-networkid $((1000 + ${shardID})) 2> ${ETH_DATA}_${shardID}_${nodeID}/eth.log &

echo "Sleep 4s to wait for bootnode start..."
sleep 4 

bootenode=`geth attach ${ETH_DATA}_${shardID}_${nodeID}/geth.ipc --exec admin.nodeInfo.enode | tr -d '"'`

for (( j=2; j<=${nodes}; j++ ))
do
${ETH_BIN}/geth --datadir=${ETH_DATA}_${shardID}_${j}  \
--rpc --rpcaddr 'localhost' --rpcport "$((9000 + ${j} + 1000*${shardID}))" \
--port "$((30303 + ${j} + 1000*(${shardID}-1)))" \
--gasprice 0 --targetgaslimit 10000000 --mine --minerthreads 1 --unlock 0 --password <(echo -n "") \
--syncmode 'full' \
-networkid $((1000 + ${shardID})) \
--bootnodes ${bootenode} 2> ${ETH_DATA}_${shardID}_${j}/eth.log &
echo "member node: ${ETH_DATA}_${shardID}_${j}"
done

echo "Sleep 2s to add peers to network..."
sleep 2
# check bootnode admin peers
geth attach ${ETH_DATA}_${shardID}_${nodeID}/geth.ipc --exec admin.peers

#geth --unlock ${BootSignerAddress} --gasprice 0 --password <(echo -n "")
