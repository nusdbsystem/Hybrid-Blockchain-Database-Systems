#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-clients-blockchaindb-$TSTAMP"
mkdir $LOGSD

set -x

size=${1:-4}
clients=${2:-4} 
workload=${3:-a}
distribution=${4:-ycsb_data}

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


echo "start test with bcdbnode addrs: ${defaultAddrs}"


nNODES="4 8 16 32 64"

for TH in $nNODES; do
    size=${TH}
    for (( c=2; c<=${size}; c++ ))
    do 
    defaultAddrs="${defaultAddrs},192.168.20.$((1+ ${c})):50001"
    done

    echo "Test start with node size: ${size}, client size: ${clients}, workload${workload}"
    ndrivers=${TH}
    nthreads=$(( ${clients} / ${ndrivers} ))
    ./restart_cluster_blockchaindb.sh ${TH}
    ./start_blockchaindb.sh ${TH}      
    sleep 10
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} 2>&1 | tee $LOGSD/blockchaindb-nodes-$TH.txt
done
./restart_cluster_blockchaindb.sh
