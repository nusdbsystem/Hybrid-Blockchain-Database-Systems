#!/usr/bin/env bash
set -ex

bestnodes=${1:-4}
bestclients=${2:-256}

# Experiment 4
echo "========================================================"
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo " Experiment 4 start"
#make fast nodes=${bestnodes}
make test nodes=${bestnodes} clients=${bestclients} workload=a
make test nodes=${bestnodes} clients=${bestclients} workload=b
make test nodes=${bestnodes} clients=${bestclients} workload=c
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo "========================================================"
