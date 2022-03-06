#!/bin/bash

SRC_DIR="temp1"
DST_DIR="temp"
BENCH_LIST="ycsb_data"
WORKLOAD_LIST="workloada.dat workloadb.dat workloadc.dat"
BIN="../bin/preprocess"

mkdir -p $DST_DIR
for BENCH in $BENCH_LIST; do
    for WORKLOAD in $WORKLOAD_LIST; do
        mkdir -p $DST_DIR/$BENCH
        $BIN --load-read-path=$SRC_DIR/$BENCH/$WORKLOAD --load-write-path=$DST_DIR/$BENCH/$WORKLOAD --run-read-path=$SRC_DIR/$BENCH/run_$WORKLOAD --run-write-path=$DST_DIR/$BENCH/run_$WORKLOAD        
    done
done