#!/bin/bash
# installing ethereum and docker
sudo apt-get install software-properties-common
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo add-apt-repository -y ppa:ethereum/ethereum-dev
sudo apt-get install apt-transport-https ca-certificates
sudo apt-get update
sudo apt-get install -y ethereum
sudo apt-get install solc



# Tools
docker pull ethereum/client-go:v1.8.23
sudo apt-get install jq
