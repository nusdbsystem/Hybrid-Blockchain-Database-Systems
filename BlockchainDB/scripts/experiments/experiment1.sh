#!/usr/bin/env bash
#set -x

size=${1:-4}

# Experiment 1
echo "========================================================"
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo " Experiment 1 start"
#make fast shards=1 nodes=${size}
make test nodes=${size} clients=4
make test nodes=${size} clients=8
make test nodes=${size} clients=16
make test nodes=${size} clients=32
make test nodes=${size} clients=64
make test nodes=${size} clients=128
make test nodes=${size} clients=192
make test nodes=${size} clients=256
echo " Experiment 1 stop"
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo "========================================================"
