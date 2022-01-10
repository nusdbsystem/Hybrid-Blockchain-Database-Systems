#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-distribution-bigchain-$TSTAMP"
mkdir $LOGSD

set -x

THREADS=4
DISTROS="uniform latest zipfian"
IPPREFIX="192.168.30"
ADDRS="http://$IPPREFIX.2:9984,http://$IPPREFIX.3:9984,http://$IPPREFIX.4:9984,http://$IPPREFIX.5:9984"

cd ..
RDIR=`pwd`
cd scripts

function copy_logs {
	DEST=$1
	mkdir -p $DEST
	for IDX in `seq 2 5`; do
		DEST_NODE=$DEST/node-$(($IDX-1))
		mkdir -p $DEST_NODE
		scp root@$IPPREFIX.$IDX:bigchaindb* $DEST_NODE/
		scp root@$IPPREFIX.$IDX:mongodb.log $DEST_NODE/
		scp root@$IPPREFIX.$IDX:tendermint.log $DEST_NODE/
	done
}

# Uniform
./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-uniform.txt
copy_logs $LOGSD/logs-bigchaindb-uniform

# Latest
./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py temp/ycsb_data_latest/workloada.dat temp/ycsb_data_latest/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-latest.txt
copy_logs $LOGSD/logs-bigchaindb-latest

# Zipfian
./restart_cluster_bigchaindb.sh
./start_bigchaindb.sh
sleep 5
python3 $RDIR/BigchainDB/bench.py temp/ycsb_data_zipfian/workloada.dat temp/ycsb_data_zipfian/run_workloada.dat $ADDRS $THREADS 2>&1 | tee $LOGSD/bigchaindb-zipfian.txt
copy_logs $LOGSD/logs-bigchaindb-zipfian
./stop_bigchaindb.sh
