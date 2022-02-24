#!/usr/bin/env bash
set -ex

nodes=${1:-4}
txdelay=${2:-0}
PREFIX="192.168.20."
dir=$(dirname "$0")
tomlDir="$dir/toml.${nodes}"
rm -rf ${tomlDir}/*
mkdir -p ${tomlDir}
mkdir -p $dir/keys
KEY_GEN_PATH=$dir/../../VeritasHotstuff/.bin



# 1
echo "##################### 1.generate config ######################"
for (( c=1; c<=${nodes}; c++ ))
do
IPX=$((${c} + 1))
$KEY_GEN_PATH/hotstuffkeygen -p 'r*' -n ${nodes} --hosts ${PREFIX}${IPX} --tls keys

tomlFile="${tomlDir}/hotstuff${c}.toml"
rm -f ${tomlFile}
touch ${tomlFile}
echo "self-id = ${c}" > ${tomlFile}
echo 'pacemaker = "fixed"' >> ${tomlFile}
# echo 'pacemaker = "round-robin"' >> ${tomlFile}
echo 'leader-id = 1' >> ${tomlFile}
echo "root-cas = [\"keys/ca.crt\"]" >> ${tomlFile}
echo "privkey = \"keys/r${c}.key\"" >> ${tomlFile}
echo 'rate-limit = 0' >> ${tomlFile}
#echo 'benchmark = false' >> ${tomlFile}
echo 'tls = false' >> ${tomlFile}
echo 'batch-size = 100' >> ${tomlFile}
echo "max-inflight = 10000" >> ${tomlFile}
echo "view-timeout = 1000" >> ${tomlFile}
echo "print-commands = true" >> ${tomlFile}
echo "client-listen = \"${PREFIX}${IPX}:20070\"" >> ${tomlFile}
echo "self-veritas-node = \"${PREFIX}${IPX}:50001\"" >> ${tomlFile}
echo "self-redis-address = \"${PREFIX}${IPX}:6379\"" >> ${tomlFile}
echo  >> ${tomlFile}

echo '# This is the information that each replica is given about the other replicas' >> ${tomlFile}
	for (( j=1; j<=${nodes}; j++ ))
	do
	IPY=$((${j} + 1))
	echo '[[replicas]]' >> ${tomlFile}
	echo "id = ${j}" >> ${tomlFile}
	echo "peer-address = \"${PREFIX}${IPY}:10070\"" >> ${tomlFile}
	echo "client-address = \"${PREFIX}${IPY}:20070\"" >> ${tomlFile}
	echo "redis-address = \"${PREFIX}${IPY}:6379\"" >> ${tomlFile}
	echo "pubkey = \"keys/r${j}.key.pub\"" >> ${tomlFile}
	echo "cert = \"keys/r${j}.crt\"" >> ${tomlFile}
	echo  >> ${tomlFile}
	done
echo "Generate config file toml.${nodes}/hotstuff${c}.toml"
done

# 2
echo "##################### 3.start hotstuff nodes ####################"

for (( j=1; j<=${nodes}; j++ )); do
	IPX=$((${j} + 1))
	scp -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" ${tomlDir}/hotstuff${j}.toml root@${PREFIX}${IPX}:/root/VeritasHotstuff/
	scp -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" keys/r${j}.key root@${PREFIX}${IPX}:/root/VeritasHotstuff/keys/
	scp -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" keys/r*.key.pub root@${PREFIX}${IPX}:/root/VeritasHotstuff/keys/
	ssh -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" root@${PREFIX}${IPX} "service redis-server start"
done
for (( j=1; j<=${nodes}; j++ )); do
	IPX=$((${j} + 1))
	ssh -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" root@${PREFIX}${IPX} "cd /root/VeritasHotstuff/; export HOTSTUFF_LOG=debug && bin/hotstuffserver --config=hotstuff${j} > /root/VeritasHotstuff/logs/hotstuff.log 2>&1 &" &
done
#
echo "Waiting for connections to the replicas..."
# wait
sleep 10

# 3
echo "##################### 3.start veritasnode nodes ####################"
for (( j=1; j<=${nodes}; j++ ))
do
IPX=$((${j} + 1))
ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" root@${PREFIX}${IPX} "cd /root/VeritasHotstuff/; export HOTSTUFF_LOG=debug && bin/veritasnode --config=hotstuff${j} > /root/VeritasHotstuff/logs/veritas.log 2>&1 &" &
echo "${IPX}"
done
sleep 10
exit 0