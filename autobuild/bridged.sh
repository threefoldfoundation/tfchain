#!/bin/bash
set -ex

apt-get update
apt-get install git gcc wget -y

# make output directory
ARCHIVE=/tmp/archives
TFCHAIN_FLIST=/tmp/tfchain
BRIDGED_FLIST=/tmp/bridged_flist


mkdir -p $ARCHIVE
mkdir -p $TFCHAIN_FLIST/bin
mkdir -p $BRIDGED_FLIST/bin

# install go
GOFILE=go1.12.linux-amd64.tar.gz
wget https://dl.google.com/go/$GOFILE
tar -C /usr/local -xzf $GOFILE
mkdir -p /root/go
export GOPATH=/root/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/go/bin

mkdir -p /root/go/src/github.com/threefoldfoundation
cp -ar /tfchain /root/go/src/github.com/threefoldfoundation/tfchain

TFCHAIN=$GOPATH/src/github.com/threefoldfoundation/tfchain
BRIDGED=$TFCHAIN/cmd/bridged


pushd $BRIDGED
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $TFCHAIN_FLIST/bin/bridged
popd

# make sure binary is executable
chmod +x $TFCHAIN_FLIST/bin/*

cp $TFCHAIN_FLIST/bin/bridged  $BRIDGED_FLIST/bin/



tar -czf "/tmp/archives/bridged.tar.gz" -C $BRIDGED_FLIST .
