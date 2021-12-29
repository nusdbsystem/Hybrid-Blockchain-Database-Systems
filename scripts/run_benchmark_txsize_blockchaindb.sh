#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-txsize-blockchaindb-$TSTAMP"
mkdir $LOGSD

set -x

size=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}
shards=${5:-1}
ndrivers=${size}

dir=$(pwd)
echo $dir
bin="$dir/../BlockchainDB/.bin/benchmark_bcdb"
defaultAddrs="192.168.20.2:50001"
nthreads=$(( ${clients} / ${ndrivers} ))

if [ ! -f ${bin} ]; then
    echo "Binary file ${bin} not found!"
    echo "Hint: "
    echo " Please build binaries by run command: "
    echo " cd ../BlockchainDB"
    echo " make build "
    echo " make docker (if never build blockchaindb image before)"
    echo " cd -"
    echo "exit 1 "
    exit 1
fi

for (( c=2; c<=${size}; c++ ))
do 
defaultAddrs="${defaultAddrs},192.168.20.$((1+ ${c})):50001"
done
echo "start test with bcdbnode addrs: ${defaultAddrs}"


nTXSIZES="ycsb_data_512B ycsb_data_2kB ycsb_data_8kB ycsb_data_32kB ycsb_data_128kB"
DURATION=5
# default gas_limit
#GAS=10000000
# increase gas_limit for tx_size tests
GAS=100000000
for TH in $nTXSIZES; do
    echo "Test start with node size: ${size}, client size: ${clients}, workload${workload}, TxSize: ${TH}"
    loadPath="$dir/temp/${TH}/workload${workload}.dat"
    runPath="$dir/temp/${TH}/run_workload${workload}.dat"
    ./restart_cluster_blockchaindb.sh 
    ./start_blockchaindb.sh ${shards} ${size} ${DURATION} ${GAS}       
    sleep 10
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/blockchaindb-txsize-$TH.txt 2>&1
done

