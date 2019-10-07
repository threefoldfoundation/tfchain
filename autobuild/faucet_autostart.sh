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
FAUCET=$TFCHAIN/frontend/faucet
FAUCET_AUTOSTART_FILE="$TFCHAIN/autobuild/startup_faucet.toml"

TFCHAIND=$TFCHAIN/cmd/tfchaind
TFCHAINC=$TFCHAIN/cmd/tfchainc


pushd $TFCHAIND
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $FAUCET_FLIST/bin/tfchaind
popd

pushd $FAUCET
go build -ldflags "-linkmode external -s -w -extldflags -static" -o $FAUCET_FLIST/bin/faucet
popd

# make sure binary is executable
chmod +x $FAUCET_FLIST/bin/*


# CADDY
CADDYMAN=/tmp/caddyman

pushd /tmp
    git clone https://github.com/Incubaid/caddyman.git
popd 

CADDY_PARENT_DIR=$GOPATH/src/github.com/mholt
CADDY_SRC=$CADDY_PARENT_DIR/caddy
pushd $CADDY_PARENT_DIR
	git clone https://github.com/mholt/caddy
popd

pushd $CADDY_SRC
	cp $CADDYMAN/caddyman.sh .
	GO111MODULE=on ./caddyman.sh install iyo
popd

ls $GOPATH/bin/
cp $GOPATH/bin/caddy $FAUCET_FLIST/bin

cp $FAUCET_AUTOSTART_FILE $FAUCET_FLIST/.startup.toml

tar -czf "/tmp/archives/faucet.tar.gz" -C $FAUCET_FLIST .
