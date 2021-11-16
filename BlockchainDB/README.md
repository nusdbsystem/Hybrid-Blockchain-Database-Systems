# BlockchainDB

##### 1.Prepare

```
make download
make build
```



#### 2.(option one) Start all by one step

```
make fast shards=1 nodes=4

```

#### 2.(option two) Start step by step

##### 2.1 Start blockchain network 

(default: ethereum poa)

```
make ethnet shards=1 nodes=4

```

##### 2.2 Start bcdb nodes

```
make install shards=1 nodes=4

```



#### 3. Run ycsb tests

```
make test nodes=4 clients=4 workload=a
```

Check test result: test.4.4.log
