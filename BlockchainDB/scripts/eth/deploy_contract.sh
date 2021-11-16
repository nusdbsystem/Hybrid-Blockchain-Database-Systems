#!/bin/bash
# set -ex

shardIDs=${1:-1}
nodeIDs=${2:-4}

cd `dirname ${BASH_SOURCE-$0}`
. env.sh
cd -

dir=$(dirname "$0")
echo $dir
bin="${dir}/../../storage/ethereum/contracts/deploy/deyploycontract"
configDir="${dir}/../../config.eth.${shardIDs}.${nodeIDs}"
ls ${configDir}

if [ ! -f ${bin} ]; then
    echo "Binary file ${bin} not found!"
    echo "Hint: "
    echo " Please build binaries by run command: make build "
    echo "exit 1 "
    exit 1
fi

for (( c=1; c<=${shardIDs}; c++ ))
do 
#$bin --config="${configDir}/shard_${c}" >> ${configDir}/*.toml
$bin --config="${configDir}/shard_${c}" | tee -a ${configDir}/*.toml
echo "Deploy contract to bcdbnode$c wtih ${configDir}/shard_${c}.conf"
done
