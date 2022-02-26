#!/bin/bash

# . ./env.sh

nodes=${1:-4}
# Nodes
# for (( c=1; c<=$nodes; c++ )); do 
# 	ADDR=$IPPREFIX".$(($c+1))"
# 	ssh -o StrictHostKeyChecking=no root@$ADDR "killall -9 veritasnode; killall -9 hotstuffserver"
# done

ps -ef | grep veritasnode| wc -l
pkill -f "veritasnode"
ps -ef | grep hotstufferver| wc -l
pkill -f "hotstuffserver"