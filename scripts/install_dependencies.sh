#!/bin/bash

MYUSER=`whoami`

# Docker - https://docs.docker.com/engine/install/ubuntu/
sudo apt-get update
sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get -y install docker-ce docker-ce-cli containerd.io
sudo adduser $MYUSER docker

# OpenVSwitch, KafkaCat, jq, make, gcc, java, pip
sudo apt -y install openvswitch-switch kafkacat jq make gcc default-jdk python3-pip

# BigchainDB driver
pip3 install bigchaindb_driver==0.6.2
# update timeout from 20s to 120s
cp ../BigchainDB/driver.py .local/lib/python3.6/site-packages/bigchaindb_driver/driver.py

# Go 1.15.6
mkdir -p temp
cd temp
wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
tar xf go1.15.6.linux-amd64.tar.gz
mkdir gopath
GOROOT=`pwd`/go
GOPATH=`pwd`/gopath
echo "" >> /home/$MYUSER/.bashrc
echo "export GOROOT=$GOROOT" >> /home/$MYUSER/.bashrc
echo "export GOPATH=$GOPATH" >> /home/$MYUSER/.bashrc
echo "export PATH=$PATH:$GOROOT/bin" >> /home/$MYUSER/.bashrc

echo "*** Please log out or reboot your system!"