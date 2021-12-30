#!/bin/bash
set -ex

shards=${1:-1}
nodes=${2:-4}
period=${3:-5}
gaslimit=${4:-10000000}
txdelay=${5:-0}
PREFIX="192.168.20."
dir=$(dirname "$0")
genesisTemplate=${dir}/../BlockchainDB/storage/ethereum/networks/CustomGenesis.template
genesisDir="${dir}/config/config.eth.${shards}.${nodes}"
mkdir -p ${genesisDir}
tomlDir="${dir}/config/config.nodes.${shards}.${nodes}"
rm -rf ${tomlDir}/*
mkdir -p ${tomlDir}

# multi-shards

# Part 1: Start ethereum nodes 
for (( shardId=1; shardId<=${shards}; shardId++ )); 
do
chainIdByShard=$((1000 + ${shardId}))
bootnode=$((1 +(${shardId} - 1) * ${nodes} ))

genesisFile="${genesisDir}/CustomGenesis_${shardId}.json"
rm -f ${genesisFile}
touch ${genesisFile}
cp ${genesisTemplate} ${genesisFile}

echo "Start blockchaindb server containers, network size(${shardId} shard, ${nodes} nodes)"

# 1
echo "##################### 1.generate ethereum genesis config ##########"
echo "chainId: $chainIdByShard"
for (( i=1; i<=${nodes}; i++ )); do
    IPX=$((${i} + ${bootnode}))
	#killall -9 geth; 
	signer1=`ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "killall -9 geth; rm -rf /Data/* && /usr/local/go/bin/geth --datadir=/Data/eth_${shardId}_${i} --password <(echo -n "") account new | cut -d '{' -f2 | cut -d '}' -f1"`

	if (( ${i} < 2 )); then
        shardsigner=${signer1}
        allocSigners=\"${signer1}\"': { "balance": "90000000" }'
    else
        allocSigners=${allocSigners}', '\"${signer1}\"': { "balance": "90000000" }'
    fi
    # set signers
    if (( ${i} <= ${nodes} )); then
        signers=${signers}${signer1}
    fi
    echo "eth-node = \"/Data/eth_${shardId}_${i}/geth.ipc\"" > ${genesisDir}/node_${i}_${shardId}.toml
    echo "eth-account-address = \"${signer1}\"" >> ${genesisDir}/node_${i}_${shardId}.toml
    hexkey=`ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "jq -r '.crypto.ciphertext' <<< cat /Data/eth_${shardId}_${i}/keystore/UTC*"`
    echo "eth-hexkey = \"${hexkey}\"" >> ${genesisDir}/node_${i}_${shardId}.toml
    echo "Generate node account file  ${genesisDir}/node_${i}_${shardId}.toml"
done
extraData="0x0000000000000000000000000000000000000000000000000000000000000000${signers}0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
sed -i "s/ChainIdByShard/${chainIdByShard}/" ${genesisFile}
sed -i "s/PeriodX/${period}/" ${genesisFile}
sed -i "s/GasLimitX/${gaslimit}/" ${genesisFile}
sed -i "s/ExtraData/${extraData}/" ${genesisFile}
sed -i "s/AllocSigners/${allocSigners}/" ${genesisFile}
echo "Generate genesis file  ${genesisFile}"
echo "eth-node = \"/Data/eth_${shardId}_1/geth.ipc\"" > ${genesisDir}/shard_${shardId}.toml
echo "eth-boot-signer-address = \"${shardsigner}\"" >> ${genesisDir}/shard_${shardId}.toml
hexkey=`ssh -o StrictHostKeyChecking=no root@${PREFIX}$((${bootnode}+1)) "jq -r '.crypto.ciphertext' <<< cat /Data/eth_${shardId}_1/keystore/UTC*"`
echo "eth-hexkey = \"${hexkey}\"" >> ${genesisDir}/shard_${shardId}.toml
echo "Generate shard file  ${genesisDir}/shard_${shardId}.toml"


# 2
echo "##################### 2.init geth nodes using genesis file ##########"
for (( i=1; i<=${nodes}; i++ )); do
    IPX=$((${i} + ${bootnode}))
    echo "Using custom genesis file: ${genesisFile}, datadir: /Data/eth_${shardId}_${i}"
    scp -o StrictHostKeyChecking=no ${genesisFile} root@${PREFIX}${IPX}:/root/BlockchainDB/config/
    ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/usr/local/go/bin/geth --datadir=/Data/eth_${shardId}_${i} init /root/BlockchainDB/config/CustomGenesis_${shardId}.json"
done


# 3
echo "##################### 3.start geth bootnode and add peers ##########"
# start bootnode
IPX=$((${bootnode} + 1))
ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/usr/local/go/bin/geth --datadir=/Data/eth_${shardId}_1 \
--rpc --rpcaddr '${PREFIX}${IPX}' --rpcport "$((9000 + ${shardId}))" \
--port "$((30303 + ${bootnode} + 1000*(${shardId}-1)))" \
--gasprice 0 --targetgaslimit 10000000 --mine --minerthreads 1 --unlock 0 --password <(echo -n "") \
--syncmode 'full' \
--nat extip:${PREFIX}${IPX} \
-networkid $((1000 + ${shardId})) > /Data/eth_${shardId}_1/eth.log 2>&1 &"

echo "Sleep 2s to wait for bootnode start..."
sleep 2
# start peernode
enodeAddr=`ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/usr/local/go/bin/geth attach /Data/eth_${shardId}_1/geth.ipc --exec admin.nodeInfo.enode  " | tr -d '"' `
for (( j=2; j<=${nodes}; j++ )); do
	IPX=$((${j} + ${bootnode}))
	ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/usr/local/go/bin/geth --datadir=/Data/eth_${shardId}_${j} \
	--rpc --rpcaddr '${PREFIX}${IPX}' --rpcport "$((9000 + ${shardId}))" \
	--port "$((30303 + ${j} + 1000*(${shardId}-1)))" \
	--gasprice 0 --targetgaslimit 10000000 --mine --minerthreads 1 --unlock 0 --password <(echo -n "") \
	--syncmode 'full' \
	-networkid $((1000 + ${shardId})) \
	--bootnodes ${enodeAddr} > /Data/eth_${shardId}_${j}/eth.log 2>&1 &"
	echo "member node: /Data/eth_${shardId}_${j}"
done
echo "Sleep 2s to add peers to network..."
sleep 2
# check bootnode admin peers
IPX=$((${bootnode}+1))
ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/usr/local/go/bin/geth attach /Data/eth_${shardId}_1/geth.ipc --exec net.peerCount"


#4
echo "##################### 4.deploy KVContract to eth network ##########"
# Deploy smart contract to each shard
IPX=$((${bootnode}+1))
scp -o StrictHostKeyChecking=no ${genesisDir}/shard_${shardId}.toml root@${PREFIX}${IPX}:/root/BlockchainDB/config/
ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/root/BlockchainDB/bin/deploy_contract --config=/BlockchainDB/config/shard_${shardId}"  | tee -a ${genesisDir}/*_${shardId}.toml
echo "Deploy contract to bcdbnode$c wtih ${genesisDir}/shard_${shardId}.toml"
done


# Part 2: Start bcdb nodes 
for (( shardId=1; shardId<=${shards}; shardId++ )); 
do
bootnode=$((1 +(${shardId} - 1) * ${nodes} ))
#5
echo "##################### 5.generate bcdbnode config ##########"

for (( c=1; c<=${nodes}; c++ )); do
    IPX=$((${c} + ${bootnode}))
	tomlFile="${tomlDir}/config_${shardId}_${c}.toml"
	rm -f ${tomlFile}
	touch ${tomlFile}
	echo "self-id = ${shardId}_${c}" > ${tomlFile}
	echo "server-node-addr = \"${PREFIX}${IPX}:50001\"" >> ${tomlFile}
	echo "delay = ${txdelay}" >> ${tomlFile}
	echo "shard-type = \"ethereum\"" >> ${tomlFile}
	echo "shard-number = \"${shards}\"" >> ${tomlFile}
	(cat "${genesisDir}/node_${c}_${shardId}.toml"; echo) >> ${tomlFile}
	echo '' >> ${tomlFile}

	echo '# This is the information that each replica is given about the other shards' >> ${tomlFile}
	for (( j=1; j<=${shards}; j++ )); do
		BOOTID=$((2 +(${j} - 1) * ${nodes} ))
		echo '[[shards]]' >> ${tomlFile}
		echo "shard-id = ${j}" >> ${tomlFile}
		echo "shard-partition-key = \"eth${j}-\"" >> ${tomlFile}
		echo "shard-type = \"ethereum\"" >> ${tomlFile}
		(cat "${genesisDir}/shard_${j}.toml"; echo) >> ${tomlFile}
		# echo '' >> ${tomlFile}
		done
	echo "Generate config file ${tomlFile}"
done


#6
echo "##################### 6.start bcdbnodes  ##########"
for (( c=1; c<=$nodes; c++ )); do 
    IPX=$((${c} + ${bootnode}))
	tomlFile="${tomlDir}/config_${shardId}_${c}.toml"
	scp -o StrictHostKeyChecking=no ${tomlFile} root@${PREFIX}${IPX}:/root/BlockchainDB/config/
    ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "killall -9 bcdbnode; sleep 2; /root/BlockchainDB/bin/bcdbnode --config=/BlockchainDB/config/config_${shardId}_${c} > /root/BlockchainDB/logs/node.${shardId}.${c}.log 2>&1 &"
    echo "bcdbnode$c ${PREFIX}${IPX} start with config file config.nodes.${shards}.${replicas}/config_${shardId}_${c}.toml"
done
echo "##################### Start blockchaindb done! ##########"


#7
echo "##################### 7.verify test  ##########"
go run ${dir}/../BlockchainDB/cmd/tests/main.go --addr=${PREFIX}$((1 + ${bootnode})):50001
done