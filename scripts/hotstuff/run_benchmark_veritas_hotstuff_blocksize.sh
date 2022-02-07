#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-blksize-veritas_hotstuff-$TSTAMP"
mkdir $LOGSD

set -x

nodes=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}
ndrivers=${nodes}
nthreads=$(( ${clients} / ${ndrivers} ))
    
dir=$(pwd)
echo $dir
bin="$dir/../VeritasHotstuff/.bin/benchmark_veritashf"
defaultAddrs="192.168.20.2:50001"
loadPath="$dir/temp/${distribution}/workload${workload}.dat"
runPath="$dir/temp/${distribution}/run_workload${workload}.dat"

if [ ! -f ${bin} ]; then
    echo "Binary file ${bin} not found!"
    echo "Hint: "
    echo " Please build binaries by run command: "
    echo " cd ../VeritasHotstuff"
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
BLKSIZES="10 100 1000 10000"
for TH in $BLKSIZES; do
    echo "Test start with node size: ${nodes}, client size: ${clients}, workload${workload}"
    ./restart_cluster_veritas_hotstuff.sh
    ./start_veritas_hotstuff.sh ${nodes} ${TH}
    sleep 6
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/veritas_hotstuff-blk-duration-${GAS}-${TH}.txt 2>&1
done

