#!/bin/bash

. ./env.sh

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-blksize-veritas_hotstuff-$TSTAMP"
mkdir $LOGSD

set -x

nodes=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}
delay=${5:-0}
ndrivers=${nodes}
nthreads=$(( ${clients} / ${ndrivers} ))
    
dir=$(pwd)
echo $dir
bin="$dir/../veritas_hotstuff/.bin/benchmark_veritashf"
defaultAddrs="192.168.20.2:50001"
loadPath="$DEFAULT_WORKLOAD_PATH/$DEFAULT_WORKLOAD".dat
runPath="$DEFAULT_WORKLOAD_PATH/run_$DEFAULT_WORKLOAD".dat

if [ ! -f ${bin} ]; then
    echo "Binary file ${bin} not found!"
    echo "Hint: "
    echo " Please build binaries by run command: "
    echo " cd ../veritas_hotstuff"
    echo " make build "
    echo " make docker (if never build veritas_hotstuff image before)"
    echo " cd -"
    echo "exit 1 "
    exit 1
fi

for (( c=2; c<=${nodes}; c++ ))
do 
defaultAddrs="${defaultAddrs},192.168.20.$((1+ ${c})):50001"
done
echo "start test with nodes addrs: ${defaultAddrs}"

# Block sizes
# BLKSIZES="10 100 1000 10000"
for TH in ${BLKSIZES}; do
    echo "Test start with node size: ${nodes}, client size: ${clients}, workload${workload}"
    ./restart_cluster_veritas_hotstuff.sh
    ./start_veritas_hotstuff.sh ${nodes} ${delay} ${TH}
   
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/veritas_hotstuff-blk-duration-${GAS}-${TH}.txt 2>&1
done

