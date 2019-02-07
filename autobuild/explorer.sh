#!/bin/bash
set -ex

apt-get update
apt-get install git gcc wget -y

# make output directory
ARCHIVE=/tmp/archives
TFCHAIN_FLIST=/tmp/tfchain

mkdir -p $ARCHIVE
mkdir -p $TFCHAIN_FLIST/bin
mkdir -p $TFCHAIN_FLIST/var/www

# install go
GOFILE=go1.11.linux-amd64.tar.gz
wget https://dl.google.com/go/$GOFILE
tar -C /usr/local -xzf $GOFILE
mkdir -p /root/go
export GOPATH=/root/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/go/bin

mkdir -p /root/go/src/github.com/threefoldfoundation
cp -ar /tfchain /root/go/src/github.com/threefoldfoundation/tfchain

TFCHAIN=$GOPATH/src/github.com/threefoldfoundation/tfchain
TFCHAIND=$TFCHAIN/cmd/tfchaind
TFCHAINC=$TFCHAIN/cmd/tfchainc
EXPLORER=$TFCHAIN/frontend/explorer


pushd $TFCHAIND
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $TFCHAIN_FLIST/bin/tfchaind
popd

pushd $TFCHAINC
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $TFCHAIN_FLIST/bin/tfchainc
popd



# CADDY
CADDYMAN=/tmp/caddyman

pushd /tmp
    git clone https://github.com/Incubaid/caddyman.git
popd 

pushd $CADDYMAN
./caddyman.sh install iyo
popd

cp $GOPATH/bin/caddy $TFCHAIN_FLIST/bin

# make sure binary is executable
chmod +x $TFCHAIN_FLIST/bin/*
cp -R $EXPLORER $TFCHAIN_FLIST/var/www/

tar -czf "/tmp/archives/explorer.tar.gz" -C $TFCHAIN_FLIST .