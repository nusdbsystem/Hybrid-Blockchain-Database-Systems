#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-clients-blockchaindb-$TSTAMP"
mkdir $LOGSD

set -x

size=${1:-4}
clients=${2:-4} 
workload=${3:-a}
distribution=${4:-ycsb_data}
ndrivers=${size}

dir=$(pwd)
echo $dir
bin="$dir/../BlockchainDB/.bin/benchmark_bcdb"
defaultAddrs="192.168.20.2:50001"
loadPath="$dir/../temp/${distribution}/workload${workload}.dat"
runPath="$dir/../temp/${distribution}/run_workload${workload}.dat"

if [ ! -f ${bin} ]; then
    echo "Binary file ${bin} not found!"
    echo "Hint: "
    echo " Please build binaries by run command: make build "
    echo "exit 1 "
    exit 1
fi

for (( c=2; c<=${size}; c++ ))
do 
defaultAddrs="${defaultAddrs},192.168.20.$((1+ ${c})):50001"
done
echo "start test with bcdbnode addrs: ${defaultAddrs}"


nCLIENTS="4 8 16 32 64 128 192 256"

for TH in $nCLIENTS; do
    echo "Test start with node size: ${size}, client size: ${TH}, workload${workload}"
    nthreads=$(( ${TH} / ${ndrivers} ))
    ./restart_cluster_blockchaindb.sh
    ./start_blockchaindb.sh        
    sleep 6
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} 2>&1 | tee $LOGSD/blockchaindb-clients-$TH.txt
done
./restart_cluster_blockchaindb.sh
