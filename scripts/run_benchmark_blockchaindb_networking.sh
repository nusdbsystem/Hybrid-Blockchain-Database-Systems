#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-networking-blockchaindb-$TSTAMP"
mkdir $LOGSD

set -x

size=${1:-4}
clients=${2:-256} 
workload=${3:-a}
distribution=${4:-ycsb_data}
shards=${5:-1}
ndrivers=${size}
DURATION=5
GAS=10000000

dir=$(pwd)
echo $dir
bin="$dir/../BlockchainDB/.bin/benchmark_bcdb"
defaultAddrs="192.168.20.2:50001"
nthreads=$(( ${clients} / ${ndrivers} ))
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

BWS="NoLimit 10000 1000 100"
RTTS="5ms 10ms 20ms 30ms 40ms 50ms 60ms"

for BW in $BWS; do    
    for RTT in $RTTS; do
	    LOGSDD="$LOGS/logs-$BW-$RTT"
	    mkdir $LOGSDD
        echo "Test start with node size: ${size}, client size: ${clients}, workload${workload}, TxSize: ${TH}"
        ./restart_cluster_blockchaindb.sh 
        if [[ "$BW" != "NoLimit" ]]; then
            sudo ./set_ovs_bs_limit.sh $BW 1
        fi
	    ./set_tc.sh $RTT
	    sleep 3
        ./start_blockchaindb.sh    
        ./run_iperf_ping.sh 2>&1 | tee $LOGSDD/net.txt 
        sleep 3
        $bin --load-path=$loadPath --run-path=$runPath --ndrivers=$ndrivers --nthreads=$nthreads --server-addrs=${defaultAddrs} > $LOGSD/blockchaindb--$BW-$RTT.txt 2>&1
    done
done

