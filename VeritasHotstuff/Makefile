nodes := 4
clients := 4

.PHONY: all build clean download build init generate install test

all: build init generate install

build:
	@go build -o ./.bin/veritasnode $(GCFLAGS) ./cmd/veritas 
	@go build -o ./.bin/hotstuffserver $(GCFLAGS) ./hotstuffserver 
	@go build -o ./.bin/hotstuffclient $(GCFLAGS) ./cmd/client/hotstuff
	@go build -o ./.bin/benchmark_veritashf $(GCFLAGS) ./benchmark


download:
	@docker pull redis:latest
	@/bin/bash scripts/gen_ycsb_data.sh

clean:
	@rm -rfv .bin

init:
	@/bin/bash scripts/start_redis_db.sh $(nodes)

generate:
	@/bin/bash scripts/gen_keys.sh $(nodes)
	@/bin/bash scripts/gen_hotstuff_config.sh $(nodes)
	
install:
	@/bin/bash scripts/stop_nodes.sh
	@/bin/bash scripts/start_veritas_nodes.sh $(nodes) > veritas.log 2>&1 && cat veritas.log
	@/bin/bash scripts/start_hotstuff_cluster.sh $(nodes) > hotstuff.log 2>&1 && cat hotstuff.log

test:
	@echo "Test start with node size: $(nodes), client size: $(clients)"
	@/bin/bash scripts/start_ycsb_test.sh $(nodes) $(clients)
