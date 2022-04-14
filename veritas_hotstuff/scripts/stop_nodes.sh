#!/usr/bin/env bash
# set -x

echo "Stop all hotstuffservers"
pgrep -f "hotstuffserver"
pkill -f "hotstuffserver"

echo "Stop all veritasnodes"
pgrep -f "veritasnode"
pkill -f "veritasnode"

echo "##################### Stop nodes successfully! ##########################"
