#!/bin/bash

. ./env.sh

set -x

N=$DEFAULT_NODES
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers>"
        echo -e "\tDefault: $N containers"
fi
TXDELAY=0
if [ $# -gt 1 ]; then
	TXDELAY=$2
else
        echo -e "Usage: $0 <# containers> <txdelay>"
        echo -e "\tDefault: $N containers"
        echo -e "\tDefault: $TXDELAY block size"
fi

IMGNAME="veritas:latest"
PREFIX="veritas"

END_IDX=$(($N+1))

# Configure Tendermint
for idx in `seq 2 $END_IDX`; do	
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "killall -9 tendermint; rm -r .tendermint; /usr/local/bin/tendermint init validator"
	for jdx in `seq 2 $END_IDX`; do
		if [ $idx -ne $jdx ]; then
			echo "," >> ids_$jdx.txt
			ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "/usr/local/bin/tendermint show-node-id" >> ids_$jdx.txt
			echo "," >> ips_$jdx.txt
		    echo $IPPREFIX.$idx >> ips_$jdx.txt
		fi
	done
	echo "," >> validators.txt
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "cat .tendermint/config/genesis.json" | jq .validators[0] >> validators.txt
	GENESIS=`ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "cat .tendermint/config/genesis.json" | jq .genesis_time`
	echo "," >> power.txt
	echo "default" >> power.txt
done
for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "killall -9 tendermint"
done
VALIDATORS=`tail +2 validators.txt | tr -d '\n' | base64 | tr -d '\n'`
POWERS=`tail +2 power.txt | tr -d '\n'`

for idx in `seq 2 $END_IDX`; do
	IDS=`tail +2 ids_$idx.txt | tr -d '\n'`
	IPS=`tail +2 ips_$idx.txt | tr -d '\n'`
	scp -o StrictHostKeyChecking=no tendermint_config.py root@$IPPREFIX.$idx:
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "./tendermint_config.py root $GENESIS generate $VALIDATORS $POWERS $IDS $IPS"
done

rm validators.txt power.txt ids*.txt ips*.txt

# Veritas Nodes
NODES=node1
for I in `seq 1 $N`; do
        NODES="$NODES,node$I"
done

# Tendermint socket
TMSOCK1="tcp://0.0.0.0:26658"
# ABCI socket
TMSOCK2="tcp://127.0.0.1:26657"

# Start Veritas
for I in `seq 1 $N`; do
	ADDR=$IPPREFIX".$(($I+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; redis-server > redis.log 2>&1 &"
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; rm -rf veritas; mkdir -p /veritas/data; nohup /bin/veritas-tendermint-txdelay --signature=node$I --parties=${NODES} --blk-size=100 --addr=:1990 --redis-addr=0.0.0.0:6379 --redis-db=0 --ledger-path=veritas$I --tendermint-socket=$TMSOCK1 --abci-socket=$TMSOCK2 --tx-delay=$TXDELAY > veritas-$I.log 2>&1 &"
done

# Start Tendermint
for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "killall -9 tendermint; sleep 1; /usr/local/bin/tendermint start --proxy-app=$TMSOCK1 > tendermint.log 2>&1 &"
done
