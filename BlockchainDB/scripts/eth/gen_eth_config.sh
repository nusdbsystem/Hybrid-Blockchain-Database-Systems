#!/bin/bash
#args: number_of_nodes, number of networks
#replicaIDs=${1:-1}
set -ex

shardIDs=${1:-1}
nodeIDs=${2:-4}

cd `dirname ${BASH_SOURCE-$0}`
. env.sh

rm -rf ${ETH_DATA}*
genesisDir=${ETH_CONFIG}.${shardIDs}.${nodeIDs}
genesisTemplate=${ETH_HOME}/networks/CustomGenesis.template
mkdir -p $genesisDir

echo '# This is custom genesis config template given about each shard'
template=$(<${genesisTemplate})
echo "${template}"

for (( j=1; j<=${shardIDs}; j++ ))
do
genesisFile="${genesisDir}/CustomGenesis_${j}.json"
rm -f ${genesisFile}
touch ${genesisFile}
chainIdByShard=$((1000 + ${j}))
cp $genesisTemplate ${genesisFile}
    for (( i=1; i<=${nodeIDs}; i++ ))
    do
    signer1=`geth --datadir=${ETH_DATA}_${j}_${i} --password <(echo -n "") account new | cut -d '{' -f2 | cut -d '}' -f1`
    # sed -i "s/Signer${i}/$signer1/" ${genesisFile}
    if (( ${i} < 2 )); then
        shardsigner=${signer1}
        allocSigners=\"${signer1}\"': { "balance": "90000000" }'
    else
        allocSigners=${allocSigners}', '\"${signer1}\"': { "balance": "90000000" }'
    fi
    # set 4 signers
    if (( ${i} <= ${nodeIDs} )); then
        signers=${signers}${signer1}
    fi
    echo "eth-node = \"${HOME}/Data/eth_${shardIDs}_${i}/geth.ipc\"" > ${genesisDir}/node_${j}_${i}.toml
    echo "eth-account-address = \"${signer1}\"" >> ${genesisDir}/node_${j}_${i}.toml
    hexkey=$(jq -r '.crypto.ciphertext' <<< cat ${HOME}/Data/eth_${j}_${i}/keystore/UTC*)
    echo "eth-hexkey = \"${hexkey}\"" >> ${genesisDir}/node_${j}_${i}.toml
    echo "Generate node account file  ${genesisDir}/node_${j}_${i}.toml"
    done
extraData="0x0000000000000000000000000000000000000000000000000000000000000000${signers}0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
sed -i "s/ChainIdByShard/${chainIdByShard}/" ${genesisFile}
sed -i "s/ExtraData/${extraData}/" ${genesisFile}
sed -i "s/AllocSigners/${allocSigners}/" ${genesisFile}

echo "Generate genesis file  ${genesisFile}"

echo "eth-node = \"${HOME}/Data/eth_${shardIDs}_1/geth.ipc\"" > ${genesisDir}/shard_${j}.toml
echo "eth-boot-signer-address = \"${shardsigner}\"" >> ${genesisDir}/shard_${j}.toml
hexkey=$(jq -r '.crypto.ciphertext' <<< cat ${HOME}/Data/eth_${j}_1/keystore/UTC*)
echo "eth-hexkey = \"${hexkey}\"" >> ${genesisDir}/shard_${j}.toml
echo "Generate shard file  ${genesisDir}/shard_${j}.toml"

echo "chainId: $chainIdByShard"
done
