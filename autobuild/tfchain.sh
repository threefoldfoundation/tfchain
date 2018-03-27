#!/bin/bash
set -ex

apt-get update
apt-get install git gcc wget -y

# make output directory
ARCHIVE=/tmp/archives
FLIST=/tmp/flist
mkdir -p $ARCHIVE

# install go
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
mkdir -p /root/go
export GOPATH=/root/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/go/bin

mkdir -p /root/go/src/github.com/threefoldfoundation
cp -ar /tfchain /root/go/src/github.com/threefoldfoundation/tfchain

TFCHAIN=$GOPATH/src/github.com/threefoldfoundation/tfchain
TFCHAIND=$TFCHAIN/cmd/tfchaind
TFCHAINC=$TFCHAIN/cmd/tfchainc

pushd $TFCHAIND
go build -ldflags "-linkmode external -extldflags -static" -o $FLIST/bin
popd

pushd $TFCHAINC
go build -ldflags "-linkmode external -extldflags -static" -o $FLIST/bin
popd

# make sure binary is executable
chmod +x $FLIST/bin/*

tar -czf "/tmp/archives/tfchain.tar.gz" -C $FLIST .