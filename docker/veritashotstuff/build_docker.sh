#!/bin/bash
set -e

dir=$(dirname "$0")
echo ${dir}

if ! [ -d "${dir}/../../VeritasHotstuff/.bin" ]; then
    echo "Please build the binaries first! (cd VeritasHotstuff && make build)"
    exit 1
fi

rm -rf ${dir}/.bin
cd ${dir}/../../VeritasHotstuff/
make build
cd -
cp -r ${dir}/../../VeritasHotstuff/.bin ${dir}/	


if ! [ -f "${dir}/id_rsa.pub" ]; then
    if ! [ -f "$HOME/.ssh/id_rsa.pub" ]; then
        echo "You do not have a public SSH key. Please generate one! (ssh-keygen)"
        exit 1
    fi
    cp $HOME/.ssh/id_rsa.pub ${dir}/
fi

docker build -f ${dir}/Dockerfile -t veritas_hotstuff ${dir}/

rm -rf ${dir}/.bin ${dir}/id_rsa.pub