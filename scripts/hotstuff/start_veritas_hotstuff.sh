#!/usr/bin/env bash
set -x

nodes=${1:-4}
txdelay=${5:-0}
PREFIX="192.168.20."
dir=$(dirname "$0")
tomlDir="$dir/toml.${nodes}"
rm -rf ${tomlDir}/*
mkdir -p ${tomlDir}
mkdir -p $dir/keys.${nodes}
KEY_GEN_PATH=$dir/.bin



# 1
echo "##################### 1.generate config ######################"
for (( c=1; c<=${nodes}; c++ ))
do
IPX=$((${i} + 1))
$KEY_GEN_PATH/hotstuffkeygen -p 'r*' -n $replicas --hosts ${PREFIX}${IPX} --tls keys.$nodes

tomlFile="${tomlDir}/hotstuff${c}.toml"
rm -f ${tomlFile}
touch ${tomlFile}
echo "self-id = ${c}" > ${tomlFile}
echo 'pacemaker = "fixed"' >> ${tomlFile}
# round-robin
echo 'leader-id = 1' >> ${tomlFile}
echo "root-cas = [\"keys.${nodes}/ca.crt\"]" >> ${tomlFile}
#echo 'consensus = "fasthotstuff"' >> ${tomlFile}
#echo 'benchmark = true' >> ${tomlFile}
# echo 'tls = true' >> ${tomlFile}
echo 'batch-size = 100' >> ${tomlFile}
#echo "input = \"input.txt\"" >> ${tomlFile}
#echo "print-commands = \"output/result1.txt\"" >> ${tomlFile}
echo "client-listen = \"${PREFIX}${IPX}:20070\"" >> ${tomlFile}
echo "self-veritas-node = \"${PREFIX}${IPX}:40070\"" >> ${tomlFile}
echo "self-redis-address = \"${PREFIX}${IPX}:6379\"" >> ${tomlFile}
echo '' >> ${tomlFile}

echo '# This is the information that each replica is given about the other replicas' >> ${tomlFile}
	for (( j=1; j<=${nodes}; j++ ))
	do
	IPY=$((${j} + 1))
	echo '[[replicas]]' >> ${tomlFile}
	echo "id = ${j}" >> ${tomlFile}
	echo "peer-address = \"${PREFIX}${IPY}:10070\"" >> ${tomlFile}
	echo "client-address = \"${PREFIX}${IPY}:20070\"" >> ${tomlFile}
	echo "redis-address = \"${PREFIX}${IPY}:6379\"" >> ${tomlFile}
	echo "pubkey = \"keys.${nodes}/r${j}.key.pub\"" >> ${tomlFile}
	echo "cert = \"keys.${nodes}/r${j}.crt\"" >> ${tomlFile}
	echo '' >> ${tomlFile}
	done
echo "Generate config file toml.${nodes}/hotstuff${c}.toml"
done

# 2
echo "##################### 3.start hotstuff nodes ####################"
for (( j=1; j<=${nodes}; j++ )); do
	IPX=$((${j} + 1))
	scp -o StrictHostKeyChecking=no ${tomlDir}/hotstuff${c}.toml root@${PREFIX}${IPX}:/root/BlockchainDB/config/
	ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/root/VeritasHotstuff/bin/hotstuffserver --config=/VeritasHotstuff/config/hotstuff${j} --privkey=keys/r${j}.key --output=output/${j}.out &"
	
done
sleep 10

# 3
echo "##################### 3.start veritasnode nodes ####################"
for (( j=1; j<=${nodes}; j++ )); do
	IPX=$((${j} + 1))
	ssh -o StrictHostKeyChecking=no root@${PREFIX}${IPX} "/root/VeritasHotstuff/bin/veritasnode --config=/VeritasHotstuff/config/hotstuff${j} &"
	
done