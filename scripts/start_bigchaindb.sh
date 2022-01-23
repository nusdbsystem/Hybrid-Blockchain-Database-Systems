#!/bin/bash

. ./env.sh

N=$DEFAULT_NODES

set -x

if [ $# -ge 1 ]; then
	N=$1
fi
END_IDX=$(($N+1))

for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "export LC_CTYPE=C.UTF-8 && export BIGCHAINDB_WSSERVER_HOST=0.0.0.0 && export BIGCHAINDB_SERVER_BIND=0.0.0.0:9984 && export BIGCHAINDB_SERVER_WORKERS=4 && /usr/src/app/scripts/start-all.sh"
	sleep 5
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "killall -9 tendermint; rm -r .tendermint; /usr/local/bin/tendermint init"
	for jdx in `seq 2 $END_IDX`; do
		if [ $idx -ne $jdx ]; then
			echo "," >> ids_$jdx.txt
			ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "/usr/local/bin/tendermint show_node_id" >> ids_$jdx.txt
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
VALIDATORS=`tail +2 validators.txt | tr -d '\n' | base64 | tr -d '\n'`
POWERS=`tail +2 power.txt | tr -d '\n'`

for idx in `seq 2 $END_IDX`; do
	IDS=`tail +2 ids_$idx.txt | tr -d '\n'`
	IPS=`tail +2 ips_$idx.txt | tr -d '\n'`
	scp -o StrictHostKeyChecking=no tendermint_config_bigchaindb.py root@$IPPREFIX.$idx:/usr/src/app/scripts/tendermint_config.py
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "cd /usr/src/app/scripts && ./tendermint_config.py root $GENESIS generate $VALIDATORS $POWERS $IDS $IPS"
done

rm validators.txt power.txt ids*.txt ips*.txt

for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$IPPREFIX.$idx "killall -9 tendermint; sleep 1; /usr/local/bin/tendermint node --p2p.laddr 'tcp://0.0.0.0:26656' --proxy_app='tcp://0.0.0.0:26658' --consensus.create_empty_blocks=false --p2p.pex=false > tendermint.log 2>&1 &"
done
