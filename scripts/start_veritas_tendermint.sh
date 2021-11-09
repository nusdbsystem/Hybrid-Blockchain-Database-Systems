#!/bin/bash

set -x

N=4
if [ $# -gt 0 ]; then
        N=$1
else
        echo -e "Usage: $0 <# containers>"
        echo -e "\tDefault: $N containers"
fi

IMGNAME="veritas"
PREFIX="veritas"

END_IDX=$(($N+1))
PREFIX="192.168.20."

set -x

for idx in `seq 2 $END_IDX`; do
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "killall -9 tendermint; rm -r .tendermint; /usr/local/bin/tendermint init"
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
VALIDATORS=`tail +2 validators.txt | tr -d '\n' | base64 | tr -d '\n'`
POWERS=`tail +2 power.txt | tr -d '\n'`

for idx in `seq 2 $END_IDX`; do
	IDS=`tail +2 ids_$idx.txt | tr -d '\n'`
	IPS=`tail +2 ips_$idx.txt | tr -d '\n'`
	scp -o StrictHostKeyChecking=no tendermint_config.py root@$PREFIX$idx:
	ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "./tendermint_config.py root $GENESIS generate $VALIDATORS $POWERS $IDS $IPS"
done

rm validators.txt power.txt ids*.txt ips*.txt

#for idx in `seq 2 $END_IDX`; do
	#ssh -o StrictHostKeyChecking=no root@$PREFIX$idx "killall -9 tendermint; sleep 1; /usr/local/bin/tendermint node --consensus.create_empty_blocks=false --p2p.pex=false > tendermint.log 2>&1 &"
#done

for I in `seq 1 $N`; do
	ADDR="192.168.20.$(($I+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "cd /; mkdir -p /veritas/data; nohup /bin/veritas-tendermint --dir=/veritas/data --config=/root/.tendermint/config/config.toml > veritas-tm-$I.log 2>&1 &"
done
