#!/usr/bin/env bash
set -ex

bestnodes=${1:-4}
bestclients=${2:-256}

# Experiment 3
# DISTROS="uniform latest zipfian"
echo "========================================================"
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo " Experiment 3 start"
#make fast nodes=${bestnodes}
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_latest
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_zipfian
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo "========================================================"
