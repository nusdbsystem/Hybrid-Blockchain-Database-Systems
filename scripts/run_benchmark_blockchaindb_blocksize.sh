#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-blksize-blockchaindb-$TSTAMP"
mkdir $LOGSD

set -x

size=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}
shards=${5:-1}
ndrivers=${size}
nthreads=$(( ${clients} / ${ndrivers} ))
    
dir=$(pwd)
echo $dir
bin="$dir/../BlockchainDB/.bin/benchmark_bcdb"
defaultAddrs="192.168.20.2:50001"
loadPath="$dir/temp/${distribution}/workload${workload}.dat"
runPath="$dir/temp/${distribution}/run_workload${workload}.dat"

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

# group 1
#nDURATIONS="1 5 10 15"
#nGASLIMITS="10000000 100000000"
# group 2
nDURATIONS="5"
nGASLIMITS="100000000 80000000 60000000 40000000 20000000 10000000"

for GAS in $nGASLIMITS; do
    for TH in $nDURATIONS; do
        echo "Test start with node size: ${size}, client size: ${clients}, workload${workload}"
        ./restart_cluster_blockchaindb.sh
        ./start_blockchaindb.sh ${shards} ${size} ${TH} ${GAS}
        sleep 6
        $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/blockchaindb-blk-duration-${GAS}-${TH}.txt 2>&1
    done
done

