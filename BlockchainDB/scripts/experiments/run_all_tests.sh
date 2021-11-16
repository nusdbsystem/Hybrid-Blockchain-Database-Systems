#!/usr/bin/env bash
set -ex


bestnodes=${1:-4}
bestclients=${2:-256}
dir=$(pwd)
mkdir -p experiments.log

# Experiment 4
#${dir}/scripts/experiments/experiment4.sh ${bestnodes} ${bestclients} >> experiments.log/experiment4.log 2>&1


# Experiment 3
${dir}/scripts/experiments/experiment3.sh ${bestnodes} ${bestclients} >> experiments.log/experiment3.log 2>&1

# Experiment 6
${dir}/scripts/experiments/experiment6.sh ${bestnodes} ${bestclients} >> experiments.log/experiment6.log 2>&1

# Experiment 1
${dir}/scripts/experiments/experiment1.sh 4 >> experiments.log/experiment1.log 2>&1


# Experiment 2
#${dir}/scripts/experiments/experiment2.sh ${bestclients} >> experiments.log/experiment2.log 2>&1

