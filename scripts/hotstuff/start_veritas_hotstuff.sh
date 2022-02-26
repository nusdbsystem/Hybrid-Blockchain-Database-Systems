#!/usr/bin/env bash
set -ex

nodes=${1:-4}
delay=${2:-0}
blksize=${3:-100}
PREFIX="192.168.20."
dir=$(dirname "$0")
tomlDir="$dir/toml.${nodes}"
rm -rf ${tomlDir}/*
mkdir -p ${tomlDir}
rm -rf $dir/keys/*
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
echo "root-cas = [\"VeritasHotstuff/keys/ca.crt\"]" >> ${tomlFile}
echo "privkey = \"VeritasHotstuff/keys/r${c}.key\"" >> ${tomlFile}
echo 'rate-limit = 0' >> ${tomlFile}
echo "tx-delay = ${delay}" >> ${tomlFile}
echo 'tls = false' >> ${tomlFile}
echo "batch-size = ${blksize}" >> ${tomlFile}
echo "max-inflight = 100000" >> ${tomlFile}
echo "view-timeout = 1000" >> ${tomlFile}
echo "print-commands = false" >> ${tomlFile}
echo "print-throughput = false" >> ${tomlFile}
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
	echo "pubkey = \"VeritasHotstuff/keys/r${j}.key.pub\"" >> ${tomlFile}
	echo "cert = \"VeritasHotstuff/keys/r${j}.crt\"" >> ${tomlFile}
	echo  >> ${tomlFile}
	done
echo "Generate config file toml.${nodes}/hotstuff${c}.toml"
done

# 2
echo "##################### 3.start hotstuff nodes ####################"
#export HOTSTUFF_LOG=error;
#killall -9 hotstuffserver; sleep 2; 
for (( j=1; j<=${nodes}; j++ )); do
	IPX=$((${j} + 1))
	scp -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" ${tomlDir}/hotstuff${j}.toml root@${PREFIX}${IPX}:/root/VeritasHotstuff/
	scp -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" keys/r${j}.key root@${PREFIX}${IPX}:/root/VeritasHotstuff/keys/
	scp -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" keys/r*.key.pub root@${PREFIX}${IPX}:/root/VeritasHotstuff/keys/
	ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" root@${PREFIX}${IPX} "service redis-server start"
done
sleep 5
for (( j=1; j<=${nodes}; j++ )); do
	IPX=$((${j} + 1))
	ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" root@${PREFIX}${IPX} "/root/VeritasHotstuff/bin/hotstuffserver --config=VeritasHotstuff/hotstuff${j} > /root/VeritasHotstuff/logs/hotstuff.log 2>&1 " &
done
#
echo "Waiting for connections to the replicas..."
# wait
sleep 30

# 3
echo "##################### 3.start veritasnode nodes ####################"
#killall -9 veritasnode; sleep 2; 
for (( j=1; j<=${nodes}; j++ ))
do
IPX=$((${j} + 1))
ssh -o LogLevel=ERROR -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" root@${PREFIX}${IPX} "/root/VeritasHotstuff/bin/veritasnode --config=VeritasHotstuff/hotstuff${j} > /root/VeritasHotstuff/logs/veritas.log 2>&1 " &
echo "${IPX}"
done
sleep 10
exit 0