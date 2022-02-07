#!/bin/bash
#
#

set -ev

replicas=${1:-4}

KEY_GEN_PATH=$PWD/.bin

$KEY_GEN_PATH/hotstuffkeygen -p 'r*' -n $replicas --hosts 127.0.0.1 --tls keys.$replicas
