# Hybrid-Blockchain-Database
## Veritas (Kafka)

### Setup

- Kafka

```bash
# run
$ docker pull zookeeper:3.4
$ docker pull wurstmeister/kafka:2.12-2.4.1
$ docker-compose -f docker_compose.yml up -d
# test
$ docker exec -it kafka1 /bin/bash
$ cd /opt/kafka_2.12-2.4.1/bin/
$ kafka-topics.sh --create --topic ${TOPIC} --replication-factor 3 --partitions 3 --zookeeper 10.0.0.4:2184
$ kafka-topics.sh --list --zookeeper 10.0.0.4:2184
$ sh /opt/kafka_2.12-2.4.1/bin/kafka-console-producer.sh --broker-list 10.0.0.4:9092 --topic ${TOPIC}
$ sh /opt/kafka_2.12-2.4.1/bin/kafka-console-consumer.sh --bootstrap-server 10.0.0.4:9092 --topic ${TOPIC} --from-beginning
```

```bash
# NOTE: Your local environment must have Java 8+ installed
# sudo apt install default-jdk
# Get Kafka
$ wget https://apachemirror.sg.wuchna.com/kafka/2.7.0/kafka_2.12-2.7.0.tgz
$ tar -xzf kafka_2.12-2.7.0.tgz
$ cd kafka_2.12-2.7.0
# Start Kafka
# Start the ZooKeeper service
$ bin/zookeeper-server-start.sh config/zookeeper.properties
# Start the Kafka broker service
$ bin/kafka-server-start.sh config/server.properties
```

```bash
# Create a topic
$ bin/kafka-topics.sh --create --topic ${topic_name} --bootstrap-server localhost:9092
# Get topic details
$ bin/kafka-topics.sh --describe --topic ${topic_name} --bootstrap-server localhost:9092
```
- Redis
```bash
$ docker pull redis:latest
$ docker run -itd --name ${NAME} -p 6379:6379 redis
```

- redisql
```bash
$ docker pull dalongrong/redisql:latest
$ docker run -itd --name ${NAME} -p 6379:6379 dalongrong/redisql
```

- Timestamp Oracle

```bash
# Install go
# sudo add-apt-repository ppa:longsleep/golang-backports
# sudo apt update
# sudo apt install golang-go

$ go run cmd/tso/main.go --addr=":7070"
```

### Build & Run

```bash
$ go build -o veritas-server cmd/veritas/main.go
$ ./veritas-server --help
usage: veritas-server --signature=SIGNATURE --parties=PARTIES --addr=ADDR --kafka-addr=KAFKA-ADDR --kafka-group=KAFKA-GROUP --kafka-topic=KAFKA-TOPIC --redis-addr=REDIS-ADDR --redis-db=REDIS-DB --ledger-path=LEDGER-PATH [<flags>]

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --signature=SIGNATURE      server signature
  --blk-size=100             block size
  --parties=PARTIES          party1,party2,...
  --addr=ADDR                server address
  --kafka-addr=KAFKA-ADDR    kafka server address
  --kafka-group=KAFKA-GROUP  kafka group id
  --kafka-topic=KAFKA-TOPIC  kafka topic
  --redis-addr=REDIS-ADDR    redis server address
  --redis-db=REDIS-DB        redis db number
  --redis-pwd=REDIS-PWD      redis password
  --ledger-path=LEDGER-PATH  ledger path
```

## Veritas (Raft)

### Build & Run

```bash
$ go build -o raftkv-server cmd/raftkv/main.go
$ ./raftkv-server --help
usage: raftkv-server --dir=DIR [<flags>]

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --svr-addr=":19001"            Address of server
  --raft-addr="127.0.0.1:18001"  Address of raft module
  --dir=DIR                      Dir for data and log
  --raft-leader=RAFT-LEADER      Address of the existing raft cluster leader
  --redis-addr=REDIS-ADDR        redis server address
  --redis-db=REDIS-DB            redis db number
  --redis-pwd=REDIS-PWD          redis password
  --store=redis                  Underlying storage [redis/badger]
  --blk-size=100                 Block size in raft
```

## Veritas (TM)

### Build & Run

```bash
$ cd veritastendermint
$ go build -o veritas-tm
```

## BigChainDB

### Requirements
- Python 3.5+
- A recent Python 3 version of pip
- A recent Python 3 version of setuptools
- cryptography and cryptoconditions
```bash
$ pip3 install --upgrade setuptools
$ pip3 install bigchaindb_driver
$ sudo apt-get update
$ sudo apt-get install python3-dev libssl-dev libffi-dev
$ pip3 install bigchaindb_driver
```
### All-in-one node [For Dev]
```bash
# NOTE: Your local environment must have Docker installed
$ docker pull bigchaindb/bigchaindb:all-in-one
$ docker run \
  --detach \
  --name bigchaindb \
  --publish 9984:9984 \
  --publish 9985:9985 \
  --publish 27017:27017 \
  --publish 26657:26657 \
  --volume $HOME/bigchaindb_docker/mongodb/data/db:/data/db \
  --volume $HOME/bigchaindb_docker/mongodb/data/configdb:/data/configdb \
  --volume $HOME/bigchaindb_docker/tendermint:/tendermint \
  bigchaindb/bigchaindb:all-in-one
```
### Network
[bigchaindb ansible playbook](https://github.com/bigchaindb/bigchaindb-node-ansible)

#### Requirements

```bash
$ chmod +x BigchainDB/setup.sh
$ sudo ./setup.sh
```

#### Tendermint Config

[Tendermint Doc](https://docs.tendermint.com/master/nodes/configuration.html)

```bash
$ tendermint node --consensus.create_empty_blocks=false
$ tendermint node --p2p.seeds "${node1_publickey}@${host1}:${port1},${node2_publickey}@${host2}:$port2"
# --p2p.persistent_peers is used for persistent connection
```

## Kafka Benchmark
```bash
$ cd ${path_to_kafka}
# Consumer
$ bin/kafka-consumer-perf-test.sh --broker-list 0.0.0.0:9092  --topic test-tps --messages 1000000 --fetch-size 1048576 --threads 10
# Producer
$ bin/kafka-producer-perf-test.sh --num-records 100000000 --record-size 1000 --topic test-tps --throughput 10000000 --producer-props bootstrap.servers=0.0.0.0:9092
```

