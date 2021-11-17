#!/bin/bash

N=4
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers>"
        echo -e "\tDefault: $N containers"        
fi

END_IDX=$(($N+1))
PREFIX="192.168.20."

TMCONFIGFILE="../BigchainDB/config.toml"
if ! [ -f "$TMCONFIGFILE" ]; then
	echo "BigchainDB Tendermint config file not found!"
	exit 1
fi

set -x

for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "cd /usr/src/app/scripts && ./start-all.sh"	
	scp -o StrictHostKeyChecking=no $TMCONFIGFILE root@$PREFIX$idx:/root/.tendermint/config/
	sleep 5
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "killall -9 tendermint; sleep 1; /usr/local/bin/tendermint node --p2p.laddr 'tcp://192.168.20.$idx:26656' --proxy_app='tcp://0.0.0.0:26658' --p2p.pex=false > tendermint.log 2>&1 &"
done

exit 0

for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "cd /usr/src/app/scripts && ./start-all.sh"
	sleep 5
	# ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "killall -9 tendermint; rm -r .tendermint; /usr/local/bin/tendermint init"
	for jdx in `seq 2 $END_IDX`; do
		if [ $idx -ne $jdx ]; then
			echo "," >> ids_$jdx.txt
			ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "/usr/local/bin/tendermint show_node_id" >> ids_$jdx.txt
			echo "," >> ips_$jdx.txt
		        echo $PREFIX$idx >> ips_$jdx.txt
		fi
	done
	echo "," >> validators.txt
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "cat .tendermint/config/genesis.json" | jq .validators[0] >> validators.txt
	GENESIS=`ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "cat .tendermint/config/genesis.json" | jq .genesis_time`
	echo "," >> power.txt
	echo "default" >> power.txt
done
for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "killall -9 tendermint"
done
VALIDATORS=`tail +2 validators.txt | tr -d '\n' | base64 | tr -d '\n'`
POWERS=`tail +2 power.txt | tr -d '\n'`

for idx in `seq 2 $END_IDX`; do
	IDS=`tail +2 ids_$idx.txt | tr -d '\n'`
	IPS=`tail +2 ips_$idx.txt | tr -d '\n'`
	scp -o StrictHostKeyChecking=no tendermint_config_v1.py root@$PREFIX$idx:/usr/src/app/scripts/
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "cd /usr/src/app/scripts && ./tendermint_config_v1.py root $GENESIS generate $VALIDATORS $POWERS $IDS $IPS"
done

rm validators.txt power.txt ids*.txt ips*.txt

for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "killall -9 tendermint; sleep 1; /usr/local/bin/tendermint node --p2p.laddr 'tcp://192.168.20.$idx:26656' --proxy_app='tcp://0.0.0.0:26658' --p2p.pex=false > tendermint.log 2>&1 &"
done
# --consensus.create_empty_blocks=false --p2p.pex=false