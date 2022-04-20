# Hybrid Blockchain Database Systems

This branch contains the code and instructions to reproduce the experiments presented 
in the paper "Hybrid Blockchain Database Systems: Design and Performance".

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
$ ./preprocess_ycsb_data.sh
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

- Effect of the Underlying Database
- Effect of Zookeeper TSO

```
./run_benchmark_veritas_kafka_database.sh
./run_benchmark_veritas_kafka_clients_tso_zk.sh
```

Note that Veritas Kafka also needs a node for Kafka.

To get the Kafka ops plotted in Figure 8, run ``./get_kafka_ops.sh <logs>`` on the logs obtained after running ``./run_benchmark_veritas_kafka_nodes.sh`` (effect of numbe rof nodes).

### Veritas + Raft

Run the following scripts to reproduce the above-mentioned experiments:

```
./run_benchmark_veritas_raft_clients.sh
./run_benchmark_veritas_raft_nodes.sh
./run_benchmark_veritas_raft_distribution.sh
./run_benchmark_veritas_raft_workload.sh
./run_benchmark_veritas_raft_blocksize.sh
./run_benchmark_veritas_raft_recordsize.sh
./run_benchmark_veritas_raft_proctime.sh
./run_benchmark_veritas_raft_networking.sh
```


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


### Veritas + HotStuff

Run the following scripts to reproduce the above-mentioned experiments:

```
./run_benchmark_veritas_hotstuff_clients.sh
./run_benchmark_veritas_hotstuff_nodes.sh
./run_benchmark_veritas_hotstuff_distribution.sh
./run_benchmark_veritas_hotstuff_workload.sh
./run_benchmark_veritas_hotstuff_blocksize.sh
./run_benchmark_veritas_hotstuff_recordsize.sh
./run_benchmark_veritas_hotstuff_proctime.sh
./run_benchmark_veritas_hotstuff_networking.sh
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