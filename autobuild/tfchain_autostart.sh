#!/bin/bash
set -ex

apt-get update
apt-get install git gcc wget -y

# make output directory
ARCHIVE=/tmp/archives
TFCHAIN_FLIST=/tmp/tfchain
TFCHAIN_AUTOSTART_FLIST=/tmp/tfchain_autostart_flist


mkdir -p $ARCHIVE
mkdir -p $TFCHAIN_FLIST/bin
mkdir -p $TFCHAIN_AUTOSTART_FLIST/bin


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
TFCHAIN_AUTOSTART_FILE="$TFCHAIN/autobuild/startup_blockcreator.toml"

pushd $TFCHAIND
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $TFCHAIN_FLIST/bin/tfchaind
popd

pushd $TFCHAINC
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $TFCHAIN_FLIST/bin/tfchainc
popd

# make sure binary is executable
chmod +x $TFCHAIN_FLIST/bin/*

cp $TFCHAIN_FLIST -R $TFCHAIN_AUTOSTART_FLIST
cp $TFCHAIN_AUTOSTART_FILE $TFCHAIN_AUTOSTART_FLIST/.startup.toml


tar -czf "/tmp/archives/tfchain_autostart.tar.gz" -C $TFCHAIN_AUTOSTART_FLIST .
