#!/usr/bin/env bash
#set -x
# trap 'trap - SIGTERM && kill -- -$$' SIGINT SIGTERM EXIT

dir=$(pwd)
echo $dir

bin="$dir/benchmark/ycsb/ycsbtest"
defaultAddrs="127.0.0.1:40071,127.0.0.1:40072,127.0.0.1:40073,127.0.0.1:40074"
loadPath="$dir/temp/ycsb_data/workloada.dat"
runPath="$dir/temp/ycsb_data/run_workloada.dat"

size=${1:-4}
ndrivers=${2:-4}
nthreads=10


for (( c=5; c<=${size}; c++ ))
do 
defaultAddrs="${defaultAddrs},127.0.0.1:$((40070 + ${c}))"
done
echo "start test with veritas addrs: ${defaultAddrs}"

# veritasAddrs=${4:-"$defaultAddrs"}

$bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --veritas-addrs=${defaultAddrs} &

# wait; wait; wait; wait
