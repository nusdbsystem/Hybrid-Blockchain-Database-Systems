#!/bin/bash

# shellcheck disable=SC2236
if [ ! -n "$1" ] ;then
    RECORD_COUNT=10000
else
    # shellcheck disable=SC2034
    RECORD_COUNT=$1
fi

# shellcheck disable=SC2236
if [ ! -n "$2" ] ;then
    OPERATION_COUNT=1000000
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
  mkdir ${WORK_DIR}/ycsb_data
fi

echo "Generating YCSB data..."
# shellcheck disable=SC2164
pushd ycsb-0.17.0
echo "Workload A..."
bin/ycsb.sh load basic -P workloads/workloada -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloada.dat
bin/ycsb.sh run basic -P workloads/workloada -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloada.dat
echo "Workload B..."
bin/ycsb.sh load basic -P workloads/workloadb -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloadb.dat
bin/ycsb.sh run basic -P workloads/workloadb -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloadb.dat
echo "Workload C..."
bin/ycsb.sh load basic -P workloads/workloadc -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloadc.dat
bin/ycsb.sh run basic -P workloads/workloadc -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloadc.dat
echo "Workload D..."
bin/ycsb.sh load basic -P workloads/workloadd -p recordcount="${RECORD_COUNT}"> "${WORK_DIR}"/ycsb_data/workloadd.dat
bin/ycsb.sh run basic -P workloads/workloadd -p recordcount="${RECORD_COUNT}" -p operationcount="${OPERATION_COUNT}"> "${WORK_DIR}"/ycsb_data/run_workloadd.dat

# shellcheck disable=SC2164
popd
# shellcheck disable=SC2164
popd
