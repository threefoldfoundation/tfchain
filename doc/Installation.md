# INSTALL #

This document explains how easy is to get tfchain up and running.

The file is [dockerfile](http://github.com/threefoldfoundation/tfchain/DockerBuilder)

You just need to use it

```bash
docker build -t tfchainbuilder -f DockerBuilder
docker run --rm -v $(pwd):/go/src/github.com/threefoldfoundation/tfchain tfchainbuilder
```

From this moment you have a docker container running with the software needed, just follow the instruction below on how to enter the container and use tfchainc, tfchaind will already be running.

## Using dockerfile ##

Using a dockerfile is the fastet way to get a tfchain tool set up and running

## Using docker ##

In this instalation we assume that you are working on computer with docker installed and running.

You need to have an image for a linux distribution, from now on we will use Ubuntu 16.04 as the one, but you may chose one that best suits you as long as a go development environment can run on it (all of the modern ones).

So we start downloading our own ubuntu image

```bash
docker pull ubuntu:16.04
```

Now we create a network, we will name it **tfnet** (you may change it if you like)

```bash
docker network create tfnet
```

We will build a container, this command creates and runs an instance of the image previously downloaded, the **/bin/bash** is just an app that is used to keep the container alive

```bash
docker run -d --name tfchain --net tfnet -i -t ubuntu:16.04 /bin/bash
```

From now on we can enter and interact with the container

```bash
docker exec -i -t tfchain /bin/bash
```

At this moment we should be inside the container and have root powers.

We need to add some software to the container. Among other we need to install the golang development environment, as we enter as root it just a series of commands that will do the needed steps

```bash
apt update      # We update the repositories
echo "deb http://archive.ubuntu.com/ubuntu xenial-backports main restricted universe multiverse" >> /etc/apt/sources.list   # we add the repository for golang
apt update                # Update again as a new repository has been added
apt install golang-1.10   # Install golang version 1.10
apt install git           # Golang uses git to load external code
mkdir -p /home/go         # We create a directory to keep golang files
export GOROOT=/usr/lib/go-1.10/
export GOPATH=/home/go
export PATH=$GOROOT/bin:$PATH
go get github.com/threefoldfoundation/tfchain/cmd/tfchaind
go get github.com/threefoldfoundation/tfchain/cmd/tfchainc
```

The **xenial-backports** line, loads the repository corresponding to the ubuntu version we are using (16.04 xenial) in this case. In other version of ubuntu this name will be different.

And that's all, the two last sentences have downloaded the needed code and compiled it, no errors expected. Two binaries must appear at **/home/go/bin/** tfchaind and tfchainc

From this moment we can start the daemon **tfchaind**

```bash
/home/go/bin/tfchaind &
```

We use the run in background ampersand **&** so we can keep the terminal available to continue working on it.

After a while we can check, that everything is running

```bash
/home/go/bin/tfchainc                    # Should present the consensus of the tfchaind
Synced: Yes
Block:  a7fbe936e32b5c66a467a7b7832c3dc71092d4819094f08fa360333508e209a5
Height: 2504
Target: [0 0 0 4 241 155 209 161 197 1 141 159 123 236 168 11 157 87 136 145 186 14 79 138 182 211 222 112 79 249 136 119]

/home/go/bin/tfchainc gateway list      # Will list all connected peers
4 active peers:
Version                    Outbound  Address
{65792 [0 0 0 0 0 0 0 0]}  Yes       185.69.166.13:23112
{65792 [0 0 0 0 0 0 0 0]}  Yes       185.69.166.12:23112
{65792 [0 0 0 0 0 0 0 0]}  Yes       185.69.166.11:23112
{65792 [0 0 0 0 0 0 0 0]}  Yes       185.69.166.14:23112
```

The consensus of the chain can differ from the sample provided, as when started (and more the first time) it needs to syncronice all the data the chain has, and while doing it, the display of consensus will show that Synced is No.

tfchainc offers may other commands to play with, you may explorer then using the --help flag.

### Wallet ###

Now you have a node on tfchain, you just need a wallet to be able to interact with the chain.

First you need to create one, when asked type a password, you will need this password always when you want to use the wallet, we don't recoment an empty password.

```bash
/home/go/bin/tfchain wallet init
You should provide a password, it may be empty if you wish.
Wallet password:
Reenter password:
Recovery seed:
vogue tossed threaten ditch toyed lucky pitched piano soccer lottery deepest asleep sadness rogue hiding eight goes energy yodel niece saucepan organs daft rarest sonic turnip maps dizzy acidic
```

After validating the password, you will be shown with a set of words, this words are the seed of your wallet, and you should keep this set of words safe and protected, as anyone with this set of words will be able to access your wallet and all the wealth it may contain.

You have just created a wallet, now you need to unlock it (with the password we typed before) and from that moment you will be able to send and receive money from other wallets on the network.

```bash
/home/go/bin/tfchain wallet unlock
Wallet password:
Unlocking the wallet. This may take several minutes...
Wallet unlocked
```

Now you are free to play and explore the wallet subcommand of tfchainc, as before, the flag --help will guide you on all the options available.

When I want to see all the transactions of our wallet

```bash
/home/go/bin/tfchain transactions
```

There is a web based explorer that you can consult at [tfexplorer](http://185.69.166.12:2015/)

### Multiple Nodes ###

The advantage of using docker is that we can create multiple nodes, just to play, or for fun or to explore the inners and working of tfchain.

If you want to create a new node, you just ned to follow the very same steps outlined in the "Using docker" section but, you will need to change the name of the container, otherwise docker will complain and it will not work.

```bash
docker run -d --name tfchain2 --net tfnet -i -t ubuntu:16.04 /bin/bash
```

Once you have started the tfchaind in the second instance, it will recognice the first daemon running as the peers will comunicate the addres.

If you don't see the address of your first tfchaind instance you can ask for a connection with it, where W.X.Y.Z is the addres of your first instance (btw: you may know the address of the node using the command "tfchain gateway address")

```bash
/home/go/bin/tfchain gateway connect W.X.Y.Z:23112
```

With docker you can create as much nodes as you wish or you workstation may stand.

## In your workstation ##

Most of the things donde to get, compile and run the tfchain on a docker container is applicable when you want to run the tfchain on your workstation.

You need a golang development environment (version 1.9 or higher), you can install it using already set repositories or manually (as you can see it in [golang](http://golang.org/dl)).

Once golang is up and running, you just need to set the GOROOT, GOPATH and PATH environment variables and from that moment, you just need to call go get

```bash
go get github.com/threefoldfoundation/tfchain/cmd/tfchaind
go get github.com/threefoldfoundation/tfchain/cmd/tfchainc
```

That will create the needed binaries and they will be placed in $GOPATH/bin. Now you can start the daemon **tfchaind** and use the client **tfchainc** to play with tfchain.

## Additional Documentation ##

You can find additional daemon documentation on [tfchaind](http://github.com/threefoldfoundation/tfchain/doc/tfchaind.md) and the client on [tfchainc](http://github.com/threefoldfoundation/tfchain/doc/tfchainc.md)
