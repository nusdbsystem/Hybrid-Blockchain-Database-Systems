#!/bin/bash
#
# Set latency
#

if [ $# -ge 1 ]; then
	NETWORK=$1
else
	echo "Usage: $0 <NETWORK_DELAY>"
	exit 1
fi

./multi_node.sh "tc qdisc del dev eth1 root"
echo "Network: $NETWORK"
case $NETWORK in
	"1ms")
		./multi_node.sh "tc qdisc add dev eth1 root netem delay 0.98ms"	
	;;
	"5ms")
        ./multi_node.sh "tc qdisc add dev eth1 root netem delay 2.48ms"
    ;;
	"10ms")
		./multi_node.sh "tc qdisc add dev eth1 root netem delay 4.98ms"
	;;
	"20ms")
		./multi_node.sh "tc qdisc add dev eth1 root netem delay 9.98ms"
	;;
	"30ms")
        ./multi_node.sh "tc qdisc add dev eth1 root netem delay 14.98ms"
	;;
	"40ms")
        ./multi_node.sh "tc qdisc add dev eth1 root netem delay 19.99ms"
    ;;
	"50ms")
        ./multi_node.sh "tc qdisc add dev eth1 root netem delay 24.98ms"
    ;;
	"60ms")
        ./multi_node.sh "tc qdisc add dev eth1 root netem delay 29.98ms"
    ;;
	"default")
		exit 0
	;;
	*)
		./multi_node.sh "sudo tc qdisc del dev eth0 root"
    	echo "$NETWORK -> reset latency"
esac