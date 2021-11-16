#!/bin/bash
#
# fabric repository: run script to deploy kvstore chaincode
#./network.sh -h
./network.sh down
./network.sh up createChannel
docker ps -a

./network.sh createChannel -c kvchannel
./network.sh deployCC -c kvchannel -ccn kvstore -ccp ../kvstore/chaincode-go -ccl go

echo "=================== Success ==================="
