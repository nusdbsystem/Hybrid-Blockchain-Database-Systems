#!/bin/bash
#
# Run command on multiple nodes
#

. ./env.sh

N=$DEFAULT_NODES
START=2
END=$(($START+$N-1))

function reset_term_color {
	echo -e "\e[0m" 
}

if [[ $# -eq 0 ]] ; then
	echo "usage: $0 <cmd> [<cmd_arg1> ... <cmd_argn> ]"	
	exit
fi

rand=$RANDOM
hosts=""
domain=""
for i in `seq $START $END`; do
	hosts="$hosts $IPPREFIX.$i"
done

#echo "Running $@ on $hosts"

for host in $hosts; do
	echo -e "\033[31m$host\e[0m"> qres-$host-$rand.log ; 
	ssh -o StrictHostKeyChecking=no root@$host$domain $@ >> qres-$host-$rand.log &
done

wait

for host in $hosts ; do
	cat qres-$host-$rand.log
	rm qres-$host-$rand.log &
done

wait

reset_term_color