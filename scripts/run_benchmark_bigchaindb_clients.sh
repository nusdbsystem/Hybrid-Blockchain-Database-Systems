#!/bin/bash

TSTAMP=`date +%F-%H-%M-%S`
LOGSD="logs-clients-bigchaindb-$TSTAMP"
mkdir $LOGSD

set -x

IPPREFIX="192.168.30"
THREADS="4 8 16 32 64 128 192 256"

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

for TH in $THREADS; do
    ./restart_cluster_bigchaindb.sh
    ./start_bigchaindb.sh        
    sleep 10
    python3 $RDIR/BigchainDB/bench.py temp/ycsb_data/workloada.dat temp/ycsb_data/run_workloada.dat http://$IPPREFIX.2:9984,http://$IPPREFIX.3:9984,http://$IPPREFIX.4:9984,http://$IPPREFIX.5:9984 $TH 2>&1 | tee $LOGSD/bigchaindb-clients-$TH.txt
    ./stop_bigchaindb.sh
    sleep 3
    copy_logs $LOGSD/logs-bigchaindb-clients-$TH
done
# ./restart_cluster_bigchaindb.sh
