# Hybrid Blockchain Database Systems

This repository contains the code and instructions to reproduce the experiments presented in the paper "Hybrid Blockchain Database Systems: Design and Performance".

## Reproduce

### Pre-requisites

- Install dependecies
- Download and prepare YCSB data
- Build binaries
- Build docker images

```bash
$ cd scripts
$ ./install_dependencies.sh
$ ./gen_ycsb_data.sh 
$ ./build_binaries.sh
$ cd ../docker/veritas
$ ./build_docker.sh
$ cd ../bigchaindb
$ ./build_docker.sh
$ cd ../blockchaindb
$ ./build_docker.sh
```

For each system, we run the following set of experiments:

- Effect of number of clients
- Effect of number of nodes (peers)
- Effect of key access distribution (uniform, latest, zipfian)
- Effect of different YCSB workloads (WorkloadA, WorkloadB, WorkloadC)
- Effect of block size
- Effect of transaction size (record size)
- Effect of transaction processing time
- Effect of Networking

### Veritas + Kafka

Run the following scripts to reproduce the above-mentioned experiments:

```
./run_benchmark_veritas_kafka_clients.sh
./run_benchmark_veritas_kafka_nodes.sh
./run_benchmark_veritas_kafka_distribution.sh
./run_benchmark_veritas_kafka_workload.sh
./run_benchmark_veritas_kafka_blocksize.sh
./run_benchmark_veritas_kafka_recordsize.sh
./run_benchmark_veritas_kafka_proctime.sh
./run_benchmark_veritas_kafka_networking.sh
```

Each script will generate a folder of the form ``logs-...-<timestamp>``. Check the results recorded in the files of these folders.

In addition, we run the following experiments for Veritas + Kafka:

- Effect of Veritas Peer Count on Kafka Operations (Figure 8)
- Effect of the Underlying Database (Figure 11)
- Effect of Zookeeper TSO (not included in the paper)

```
./run_benchmark_veritas_kafka_database.sh
./run_benchmark_veritas_kafka_clients_tso_zk.sh
```

#### Effect of Veritas Peer Count on Kafka Operations

This reuses the logs from running ``./run_benchmark_veritas_kafka_nodes.sh``, which are in the form ``logs-nodes-veritas-<timestamp>``. In particular, we are interested in the following Kafka counters: 

- CURRENT-OFFSET - the current position of a consumer. There are N consumers in total, so we need the sum of these counters to get the number of read operations.
- LOG-END-OFFSET - the offset of the last message written to the topic. This represents the number of write operations.

Please run ``./get_kafka_ops.sh logs-nodes-veritas-...`` to get the read and write operations for a given log folder. Note that Figure 8 represents the average of 3 such logs (runs).


### Veritas + Tendermint

We use the same ``veritas`` docker images. Then, we run the following scripts:

```
./run_benchmark_veritas_tendermint_clients.sh
./run_benchmark_veritas_tendermint_nodes.sh
./run_benchmark_veritas_tendermint_distribution.sh
./run_benchmark_veritas_tendermint_workload.sh
./run_benchmark_veritas_tendermint_recordsize.sh
./run_benchmark_veritas_tendermint_proctime.sh
./run_benchmark_veritas_tendermint_networking.sh
```

### BlockchainDB

We use ``blockchaindb`` docker images. Then, we run the following scripts:

```
./run_benchmark_blockchaindb_clients.sh
./run_benchmark_blockchaindb_nodes.sh
./run_benchmark_blockchaindb_distribution.sh
./run_benchmark_blockchaindb_workload.sh
./run_benchmark_blockchaindb_blocksize.sh
./run_benchmark_blockchaindb_recordsize.sh
./run_benchmark_blockchaindb_proctime.sh
./run_benchmark_blockchaindb_networking.sh
./run_benchmark_blockchaindb_sharding.sh
```

### 

### BigchainDB

We use ``bigchaindb`` docker images. Then, we run the following scripts:

```
./run_benchmark_bigchaindb_clients.sh
./run_benchmark_bigchaindb_nodes.sh
./run_benchmark_bigchaindb_distribution.sh
./run_benchmark_bigchaindb_workload.sh
./run_benchmark_bigchaindb_recordsize.sh
./run_benchmark_bigchaindb_proctime.sh
./run_benchmark_bigchaindb_networking.sh
```

### BigchainDB Parallel Validation

To run BigchainDB with Parallel Validation (PV), modify lines 16 and 17 of [BigchainDB/scripts/start-all.sh](BigchainDB/scripts/start-all.sh), such as:

```
# bigchaindb start > /dev/null 2>&1 &
bigchaindb start --experimental-parallel-validation > /dev/null 2>&1 &
```

Next, re-build the Docker image:

```
cd docker/bigchaindb
rm -r bigchaindb-2.2.2/scripts
./build_docker.sh
```

Then repeat all the steps of BigchainDB:

```
./run_benchmark_bigchaindb_clients.sh
./run_benchmark_bigchaindb_nodes.sh
./run_benchmark_bigchaindb_distribution.sh
./run_benchmark_bigchaindb_workload.sh
./run_benchmark_bigchaindb_recordsize.sh
./run_benchmark_bigchaindb_proctime.sh
./run_benchmark_bigchaindb_networking.sh
```


### Aborted Transaction

The number of aborted transactions is reported by Veritas (Kafka) in its server log, such as:

```
2021/11/13 13:57:20 Abort transaction 1
```

By counting the number of such messages, one can get the total number of aborted transactions:

```
cat logs-distribution-veritas-kafka-.../veritas-latest-logs/veritas-* | grep Abort | wc -l
```

There are no aborted transactions in Veritas (TM) and BlockchainDB.

## License

MIT License

## Authors

Zerui Ge, Dumitrel Loghin, Tianwen Wang, Pingcheng Ruan, Beng Chin Ooi 
