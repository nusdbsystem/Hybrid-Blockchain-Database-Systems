#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-txsize-veritas_hotstuff-$TSTAMP"
mkdir $LOGSD

set -x

nodes=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}

ndrivers=${nodes}

dir=$(pwd)
echo $dir
bin="$dir/../VeritasHotstuff/.bin/benchmark_veritashf"
defaultAddrs="192.168.20.2:50001"
nthreads=$(( ${clients} / ${ndrivers} ))

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


nTXSIZES="ycsb_data_512B ycsb_data_2kB ycsb_data_8kB ycsb_data_32kB ycsb_data_128kB"

for TH in $nTXSIZES; do
    echo "Test start with node size: ${nodes}, client size: ${clients}, workload${workload}, TxSize: ${TH}"
    loadPath="$dir/temp/${TH}/workload${workload}.dat"
    runPath="$dir/temp/${TH}/run_workload${workload}.dat"
    ./restart_cluster_veritas_hotstuff.sh 
    ./start_veritas_hotstuff.sh       
    sleep 10
    $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/veritas_hotstuff-txsize-$TH.txt 2>&1
done

