#!/bin/bash
set -e

dir=$(dirname "$0")
echo ${dir}



rm -rf ${dir}/.bin
cd ${dir}/../../veritas_hotstuff/
make build
cd -
git submodule update --init
cd ${dir}/../../veritas_hotstuff/hotstuff
git checkout veritas
go build -o ../.bin/hotstuffkeygen ./cmd/hotstuffkeygen
cd -

if ! [ -d "${dir}/../../veritas_hotstuff/.bin" ]; then
    echo "Please build the binaries first! (cd veritas_hotstuff && make build)"
    exit 1
fi
cp -r ${dir}/../../veritas_hotstuff/.bin ${dir}/	


if ! [ -f "${dir}/id_rsa.pub" ]; then
    if ! [ -f "$HOME/.ssh/id_rsa.pub" ]; then
        echo "You do not have a public SSH key. Please generate one! (ssh-keygen)"
        exit 1
    fi
    cp $HOME/.ssh/id_rsa.pub ${dir}/
fi

docker build -f ${dir}/Dockerfile -t veritas_hotstuff ${dir}/

rm -rf ${dir}/.bin ${dir}/id_rsa.pub