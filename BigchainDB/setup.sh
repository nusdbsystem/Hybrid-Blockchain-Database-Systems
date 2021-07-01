#!/bin/bash

sudo apt update
sudo apt full-upgrade -y
sudo apt install -y python3-pip libssl-dev unzip mongodb git tmux jq

wget https://github.com/tendermint/tendermint/releases/download/v0.22.8/tendermint_0.22.8_linux_amd64.zip -O tendermint.zip
unzip tendermint.zip
rm tendermint.zip
sudo mv tendermint /usr/local/bin

sudo ufw allow 22/tcp
sudo ufw allow 9984/tcp
sudo ufw allow 9985/tcp
sudo ufw allow 26656/tcp
sudo ufw allow 26657/tcp
yes | sudo ufw enable

git clone https://github.com/bigchaindb/bigchaindb.git

sudo pip3 install -e bigchaindb/
sudo pip3 install -e benchmark/

pip3 uninstall -y requests
pip3 install requests==2.23.0
pip3 uninstall -y python-rapidjson
pip3 install python-rapidjson==0.9.1
pip3 uninstall -y cryptoconditions
pip3 install cryptoconditions==0.8.0
pip3 uninstall -y aiohttp
pip3 install aiohttp==3.6.2
pip3 install psutil==5.7.0
