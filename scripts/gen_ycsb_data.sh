#!/bin/bash

# shellcheck disable=SC2236
if [ ! -n "$1" ] ;then
    RECORD_COUNT=100000
else
    # shellcheck disable=SC2034
    RECORD_COUNT=$1
fi

# shellcheck disable=SC2236
if [ ! -n "$2" ] ;then
    OPERATION_COUNT=100000
else
    # shellcheck disable=SC2034
    OPERATION_COUNT=$2
fi

echo "Record count: ${RECORD_COUNT}"
echo "Operation count: ${OPERATION_COUNT}"

# Check Java is installed
if ! [ -x "$(command -v java)" ]; then
  echo 'Error: Java is not installed!' >&2
  exit 1
fi
# shellcheck disable=SC2046
# shellcheck disable=SC2005
echo $(java -version 2>&1 |awk 'NR==1')

# shellcheck disable=SC2034
WORK_DIR="${PWD}/temp"
if [ ! -d "${WORK_DIR}" ]; then
  mkdir -p "${WORK_DIR}"
fi

WLOAD_DIR="${PWD}/../workloads"

# shellcheck disable=SC2164
pushd "${WORK_DIR}"

# Get ycsb generator tool
if [ ! -d "${WORK_DIR}/ycsb-0.17.0" ]; then
  echo "Generator not found, downloading from github..."
  curl -O --location https://github.com/brianfrankcooper/YCSB/releases/download/0.17.0/ycsb-0.17.0.tar.gz
  tar xfvz ycsb-0.17.0.tar.gz
fi

# Dir for ycsb data
if [ ! -d "${WORK_DIR}/ycsb_data" ]; then
  # shellcheck disable=SC2086
  # uniform (default)
  mkdir ${WORK_DIR}/ycsb_data
  # latest
  mkdir ${WORK_DIR}/ycsb_data_latest
  # zipfian
  mkdir ${WORK_DIR}/ycsb_data_zipfian
fi

echo "Generating YCSB data..."
# shellcheck disable=SC2164
pushd ycsb-0.17.0

# Uniform
echo "Workload A..."
bin/ycsb.sh load basic -P $WLOAD_DIR/uniform/workloada -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloada.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/uniform/workloada -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloada.dat
echo "Workload B..."
bin/ycsb.sh load basic -P $WLOAD_DIR/uniform/workloadb -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloadb.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/uniform/workloadb -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloadb.dat
echo "Workload C..."
bin/ycsb.sh load basic -P $WLOAD_DIR/uniform/workloadc -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloadc.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/uniform/workloadc -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloadc.dat

# Latest
echo "Workload A..."
bin/ycsb.sh load basic -P $WLOAD_DIR/latest/workloada -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_latest/workloada.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/latest/workloada -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_latest/run_workloada.dat
echo "Workload B..."
bin/ycsb.sh load basic -P $WLOAD_DIR/latest/workloadb -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_latest/workloadb.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/latest/workloadb -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_latest/run_workloadb.dat
echo "Workload C..."
bin/ycsb.sh load basic -P $WLOAD_DIR/latest/workloadc -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_latest/workloadc.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/latest/workloadc -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_latest/run_workloadc.dat

# Zipfian
echo "Workload A..."
bin/ycsb.sh load basic -P $WLOAD_DIR/zipfian/workloada -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_zipfian/workloada.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/zipfian/workloada -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_zipfian/run_workloada.dat
echo "Workload B..."
bin/ycsb.sh load basic -P $WLOAD_DIR/zipfian/workloadb -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_zipfian/workloadb.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/zipfian/workloadb -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_zipfian/run_workloadb.dat
echo "Workload C..."
bin/ycsb.sh load basic -P $WLOAD_DIR/zipfian/workloadc -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_zipfian/workloadc.dat
bin/ycsb.sh run basic -P $WLOAD_DIR/zipfian/workloadc -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_zipfian/run_workloadc.dat

# Different Record Sizes
SIZES="512B 2kB 8kB 32kB 128kB"
for SIZE in $SIZES; do
  mkdir -p "${WORK_DIR}"/ycsb_data_$SIZE
  bin/ycsb.sh load basic -P $WLOAD_DIR/sizes/workloada_$SIZE -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data_$SIZE/workloada.dat
  bin/ycsb.sh run basic -P $WLOAD_DIR/uniform/workloada_$SIZE -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data_$SIZE/run_workloada.dat
done

# shellcheck disable=SC2164
popd
# shellcheck disable=SC2164
popd
