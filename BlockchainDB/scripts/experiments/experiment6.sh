#!/usr/bin/env bash
set -ex

bestnodes=${1:-4}
bestclients=${2:-256}
TXSIZES="512B 2kB 8kB 32kB 128kB"

# Experiment 6
# DISTROS="uniform latest zipfian"
echo "========================================================"
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo " Experiment 3 start"
#make fast nodes=${bestnodes}
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_512B
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_2kB
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_8kB
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_32kB
make test nodes=${bestnodes} clients=${bestclients} distribution=ycsb_data_128kB
printf -v date '%(%Y-%m-%d %H:%M:%S)T\n' -1 
echo "========================================================"
