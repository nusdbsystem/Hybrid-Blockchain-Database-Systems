#!/bin/bash

. ./env.sh

nodes=${1:-4}
# Nodes
for (( c=1; c<=$nodes; c++ )); do 
	ADDR=$IPPREFIX".$(($c+1))"
	ssh -o StrictHostKeyChecking=no root@$ADDR "killall -9 bcdbnode; killall -9 geth"
done
