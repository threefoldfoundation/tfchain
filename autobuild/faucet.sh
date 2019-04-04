#!/bin/bash
set -ex

apt-get update
apt-get install git gcc wget -y

# make output directory
ARCHIVE=/tmp/archives
FAUCET_FLIST=/tmp/faucet


mkdir -p $ARCHIVE
mkdir -p $FAUCET_FLIST/bin

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
FAUCET=$TFCHAIN/frontend/tftfaucet


pushd $FAUCET
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $FAUCET_FLIST/bin/tftfaucet
popd

# make sure binary is executable
chmod +x $FAUCET_FLIST/bin/*


# CADDY

CADDYMAN=/tmp/caddyman

pushd /tmp
    git clone https://github.com/Incubaid/caddyman.git
popd 

pushd $CADDYMAN
./caddyman.sh install iyo
popd

cp $GOPATH/bin/caddy $FAUCET_FLIST/bin


tar -czf "/tmp/archives/faucet.tar.gz" -C $FAUCET_FLIST .
