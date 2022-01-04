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
- Effect of transaction size (value size)
- Effect of transaction processing time
- Effect of Networking

### Veritas + Kafka

Run the following scripts to reproduce the above-mentioned experiments:

```
./run_benchmark_clients_veritas_kafka.sh
./run_benchmark_nodes_veritas_kafka.sh
./run_benchmark_distribution_veritas_kafka.sh
./run_benchmark_workload_veritas_kafka.sh
./run_benchmark_blksize_veritas_kafka.sh
./run_benchmark_txsize_veritas_kafka.sh
./run_benchmark_txdelay_veritas_kafka.sh
./run_benchmark_networking_veritas_kafka.sh
```

Each script will generate a folder of the form ``logs-...-<timestamp>``. Check the results recorded in the files of these folders.

In addition, we run the following experiments for Veritas + Kafka:

- Effect of the Underlying Database
- Effect of Zookeeper TSO

```
./run_benchmark_database_veritas_kafka.sh
./run_benchmark_clients_veritas_kafka_tso_zk.sh
```

Note that Veritas Kafka also needs a node for Kafka.

### Veritas + Tendermint

We use the same ``veritas`` docker images. Then, we run the following scripts:

```
./run_benchmark_clients_veritas_tendermint.sh
./run_benchmark_nodes_veritas_tendermint.sh
./run_benchmark_distribution_veritas_tendermint.sh
./run_benchmark_workload_veritas_tendermint.sh
./run_benchmark_txsize_veritas_tendermint.sh
./run_benchmark_txdelay_veritas_tendermint.sh
./run_benchmark_networking_veritas_tendermint.sh
```

### BlockchainDB

We use ``blockchaindb`` docker images. Then, we run the following scripts:

```
./run_benchmark_clients_blockchaindb.sh
./run_benchmark_nodes_blockchaindb.sh
./run_benchmark_distribution_blockchaindb.sh
./run_benchmark_workload_blockchaindb.sh
./run_benchmark_blksize_blockchaindb.sh
./run_benchmark_txsize_blockchaindb.sh
./run_benchmark_txdelay_blockchaindb.sh
./run_benchmark_sharding_blockchaindb.sh
```

### 

### BigchainDB

We use ``bigchaindb`` docker images. Then, we run the following scripts:

```
./run_benchmark_clients_bigchaindb.sh
./run_benchmark_nodes_bigchaindb.sh
./run_benchmark_distribution_bigchaindb.sh
./run_benchmark_workload_bigchaindb.sh
./run_benchmark_blksize_bigchaindb.sh
./run_benchmark_txsize_bigchaindb.sh
./run_benchmark_txdelay_bigchaindb.sh
./run_benchmark_networking_bigchaindb.sh
```

### BigchainDB Parallel Validation

To run BigchainDB with Parallel Validation (PV), modify lines 16 and 17 of [BigchainDB/scripts/start-all.sh](BigchainDB/scripts/start-all.sh), such as:

```
# bigchaindb start > /dev/null 2>&1 &
bigchaindb start --experimental-parallel-validation > /dev/null 2>&1 &
```

Then repeat all the steps of BigchainDB.


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