#!/bin/bash

. ./env.sh

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-workload-veritas_hotstuff-$TSTAMP"
mkdir $LOGSD

set -x

nodes=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}

ndrivers=${DEFAULT_DRIVERS_VERITAS_HOTSTUFF}

dir=$(pwd)
echo $dir
bin="$dir/../veritas_hotstuff/.bin/benchmark_veritashf"
defaultAddrs="192.168.20.2:50001"
nthreads=$(( ${clients} / ${ndrivers} ))

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


for TH in $WORKLOADS; do
    echo "Test start with node size: ${nodes}, client size: ${clients}, workload${TH}"
    loadPath="$DEFAULT_WORKLOAD_PATH/${TH}.dat"
    runPath="$DEFAULT_WORKLOAD_PATH/run_${TH}.dat"
    ./restart_cluster_veritas_hotstuff.sh
    ./start_veritas_hotstuff.sh        
    
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/veritas_hotstuff-workload-$TH.txt 2>&1
done

