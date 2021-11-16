#!/usr/bin/env bash
set -ex


shardIDs=${1:-1}
replicaIDs=${2:-4}

echo "Usage: ./scripts/gen_config.sh 4 1"
echo "Generate config files, shards: ${shardIDs}, replicas: ${replicaIDs}"
cd `dirname ${BASH_SOURCE-$0}`
. eth/env.sh
cd -
tomlDir="${ETH_CONFIG}/config.nodes.${shardIDs}.${replicaIDs}"
shardDir="${ETH_CONFIG}/config.eth.${shardIDs}.${replicaIDs}"
rm -rf ${tomlDir}/*
mkdir -p ${tomlDir}

for (( i=1; i<=${shardIDs}; i++ ))
do
	for (( c=1; c<=${replicaIDs}; c++ ))
	do
	tomlFile="${tomlDir}/config_${i}_${c}.toml"
	rm -f ${tomlFile}
	touch ${tomlFile}
	echo "self-id = ${i}_${c}" > ${tomlFile}
	echo "server-node-addr = \"127.0.0.1:$((50000 + ${c}))\"" >> ${tomlFile}
	echo "shard-type = \"ethereum\"" >> ${tomlFile}
	echo "shard-number = \"${shardIDs}\"" >> ${tomlFile}
	(cat "$shardDir/node_${i}_${c}.toml"; echo) >> ${tomlFile}
	# echo "fab-node = \"127.0.0.1:$((40000 + ${c}))\"" >> ${tomlFile}
	# echo "fab-config = \"connection${c}.yaml\"" >> ${tomlFile}
	echo '' >> ${tomlFile}

	echo '# This is the information that each replica is given about the other shards' >> ${tomlFile}
		for (( j=1; j<=${shardIDs}; j++ ))
		do
		echo '[[shards]]' >> ${tomlFile}
		echo "shard-id = ${j}" >> ${tomlFile}
		echo "shard-partition-key = \"eth${j}-\"" >> ${tomlFile}
		echo "shard-type = \"ethereum\"" >> ${tomlFile}
		echo "redis-address = \"127.0.0.1:$((60000 + ${j}))\"" >> ${tomlFile}
		(cat "$shardDir/shard_${j}.toml"; echo) >> ${tomlFile}
		#echo "eth-node = \"http://localhost:$((9000 + ${c} + 1000*${j}))\"" >> ${tomlFile}
		# echo "eth-node = \"$HOME/Data/eth_${shardIDs}_${c}/geth.ipc\"" >> ${tomlFile}
		# echo "eth-hexaddr = \"0x70fa2c27a4e365cdf64b2d8a6c36121eb80bb442\"" >> ${tomlFile}
		# echo "eth-hexkey = \"35fc8e4f2065b6813078a08069e3a946f203029ce2bc6a62339d30c37f978403\"" >> ${tomlFile}
		# echo "fab-node = \"127.0.0.1:$((40000 + ${j}))\"" >> ${tomlFile}
		# echo "fab-config = \"connection${j}.yaml\"" >> ${tomlFile}
		echo '' >> ${tomlFile}
		done
	echo "Generate config file ${tomlFile}"
	done
done

echo "Done!"