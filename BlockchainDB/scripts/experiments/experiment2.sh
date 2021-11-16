#!/usr/bin/env bash
#set -x

bestclients=${1:-16}

# Experiment 2
echo "========================================================"
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo " Experiment 2 start"
make fast nodes=8
make test nodes=8 clients=${bestclients}

make fast nodes=16
make test nodes=16 clients=${bestclients}

make fast nodes=32
make test nodes=32 clients=${bestclients}

make fast nodes=64
make test nodes=64 clients=${bestclients}

printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo "========================================================"
