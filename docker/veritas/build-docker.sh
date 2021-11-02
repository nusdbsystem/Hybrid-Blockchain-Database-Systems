#!/bin/bash

#if ! [ -d "bin" ]; then
rm -rd bin
cp -r ../../bin .	
#fi
if ! [ -f "id_rsa.pub" ]; then
    if ! [ -f "~/.ssh/id_rsa.pub" ]; then
        echo "You do not have a public SSH key. Please generate one with 'ssh-keygen'!"
        exit 1
    fi
    cp ~/.ssh/id_rsa.pub .
fi
docker build -f Dockerfile -t veritas .