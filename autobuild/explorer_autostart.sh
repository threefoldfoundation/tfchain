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
GOFILE=go1.12.linux-amd64.tar.gz
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
EXPLORER_AUTOSTART=$TFCHAIN/autobuild/startup_explorer.toml 

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

CADDY_PARENT_DIR=$GOPATH/src/github.com/mholt
CADDY_SRC=$CADDY_PARENT_DIR/mholt
pushd CADDY_PARENT_DIR
	git clone https://github.com/mholt/caddy
popd

pushd $CADDY_SRC
	cp $CADDYMAN/caddyman.sh .
	GO111MODULE=on ./caddyman.sh install iyo
popd

cp $GOPATH/bin/caddy $TFCHAIN_FLIST/bin

# make sure binary is executable
chmod +x $TFCHAIN_FLIST/bin/*
cp -R $EXPLORER $TFCHAIN_FLIST/var/www/
cp $EXPLORER_AUTOSTART $TFCHAIN_FLIST/.startup.toml

tar -czf "/tmp/archives/explorer-autostart.tar.gz" -C $TFCHAIN_FLIST .
