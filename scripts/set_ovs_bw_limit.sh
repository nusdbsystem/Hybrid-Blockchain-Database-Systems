#!/bin/bash
if [ $EUID -ne 0 ]; then
	echo "This script must be run as root!"
	exit 1
fi
if [ $# -lt 2 ]; then
	echo "Usage: $0 <bandwidth [Mbps]> <latency [ms]>"
	exit 1
fi
BW=$(($1*1000000))
LAT=$2
PORTS=`ovs-appctl dpif/show | tail +3 | head -n 20 | cut -d ' ' -f 1 | tr -d '\t'`
for PORT in $PORTS; do
	echo $PORT
	ovs-vsctl -- set port $PORT qos=@newqos -- --id=@newqos create qos type=linux-htb other-config:max-rate=$BW queues=0=@q0 -- --id=@q0 create queue other-config:min-rate=$BW other-config:max-rate=$BW
done