#!/bin/bash

# OpenVSwitch
sudo apt install openvswitch-switch

# kafka
wget https://archive.apache.org/dist/kafka/2.7.0/kafka_2.12-2.7.0.tgz
tar xf kafka_2.12-2.7.0.tgz

# Redis
docker pull redis:latest