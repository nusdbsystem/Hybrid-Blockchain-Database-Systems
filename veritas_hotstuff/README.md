# Veritas + HotStuff



For benchmark experiments, you don't need to follow the steps. Just build the Dockerfile and run benchmark scripts, it will handle everything.



### Download

1. ###### Download go pkg

   <!--Run `go mod tidy` to install go packages.-->

2. ###### Download redis image 

   <!--docker pull redis:latest-->

3. ###### Download YCSB data

   <!--`$./scripts/gen_ycsb_data.sh`-->

   ###### *Run the following command to download all:

   ```
    $ make download
   ```

   

### Build, Generate & Install

1. ##### Build binaries

   `$ make build`

2. ##### Init redis db

   <!--Usage: ./scripts/start_redis_db.sh ${networkSize}-->

   <!--$ ./scripts/start_redis_db.sh 4-->

   `$ make init nodes=4`

3. ##### Generate keys & Hotstuff config files 

   <!--Usage:  ./scripts/gen_keys.sh ${networkSize}-->

   <!--e.g. `$./scripts/gen_keys.sh 4`-->

   <!--Usage:  ./scripts/gen_hotstuff_config.sh ${networkSize}-->

   <!--e.g.`$./scripts/gen_hotstuff_config.sh 4`-->

   `$ make generate nodes=4`

4. ##### Start Hotstuff Servers & Veritas Nodes

   <!--Usage: ./scripts/start_hotstuff_cluster.sh ${networkSize}-->

   <!--e.g.`$ export HOTSTUFF_LOG=info && scripts/start_hotstuff_cluster.sh 4`--> 

   <!--*Start hotstuff client(Debug only)-->

   <!--`$ export HOTSTUFF_LOG=debug && cmd/client/hotstuffclient`-->

   <!--The client read data from input.txt and send to hotstuff servers.-->

   <!--The server execute commands, io.write commands to "output/result.txt"--> 

   <!--Usage: ./scripts/start_veritas_nodes.sh ${networkSize}-->

   <!--e.g.`$ export HOTSTUFF_LOG=debug && scripts/start_veritas_nodes.sh 4`--> 

   `$ make install nodes=4`

   ##### *Or simply run the following command for all the above steps:

   ```
   $ make all nodes=4
   ```

   

### Run Benchmark Tests

Run ycsb tests with workload, e.g.`./scripts/start_ycsb_test.sh ${nodeSize} ${clientSize} ` 

(Default workloada.dat) 

##### *run test with the following command:

```
$ make test nodes=4 clients=4
```



### Record Test Results

1. Throughput

2. Latency

