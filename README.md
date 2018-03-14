# tfchain

tfchain is the official Go implementation of the ThreeFold blockchain client. It uses and is build on top of the [Rivine][rivine] protocol.

The ThreeFold blockchain is used as a ledger for the ThreeFold Token ("TFT"), a digital currency backed by neutral and sustainable internet capacity. You can learn more about the ThreeFold token on [threefoldtoken.com](https://threefoldtoken.com).

The ThreeFold blockchain is the technology behind the ThreeFold Grid ("TF Grid" or "Grid"), a grid where your data is local and controlled by you, shaping a new neutral and sustainable Internet. To learn more about this grid you can find and read the whitepaper at [threefoldtoken.com/pdf/tf_whitepaper.pdf](https://threefoldtoken.com/pdf/tf_whitepaper.pdf).

[rivine]: http://github.com/rivine/rivine

## install and use tfchain

You have 2 easy options to install tfchain in your work environment, as a prequisite to using the tfchain binaries and thus become a full tfchain node. Currently there is only a test network.

### tfchain docker container

The easiest is pulling and using the latest prebuilt docker container:

```
$ docker pull tfchain/tfchain:v0.1.0
$ docker run -d --name tfchain tfchain/tfchain:v0.1.0
```

This will pull and configure the latest tfchain container, and it will start a container using that image in the program, named `tfchaind`.

Should you want to use to use the CLI client, you can do so doing the running container:

```
$ docker exec -ti tfchain /tfchainc 
```

Note that this minimal tfchain docker container is not really meant to be used for much tfchain CLI client interaction. If you wish to do a lot of that, it is probably more easy/useful to run a the tfchain binary CLI from your host machine. If you are interested in that, you can check out [the "tfchain from source" section](#tfchain-from-source).

Even though the image for these containers is prebuilt available for you, should you wish to use a tfchain docker container from the hacked source code you can rebuilt those imagines using `docker-minimal`.

### tfchain from source

tfchain is developed and implemented using [Golang](http://golang.org). Using the golang toolchain it is very easy to download, update and install the tfchain binaries used to run a full node and interact with it:

```
$ go get -u github.com/threefoldfoundation/tfchain/cmd/... && \
    tfchaind &
```

At this point (if all went right) you should have a tfchain daemon running in the background which is syncing with the test net. You can follow this syncing process using the CLI client: `tfchainc`.

Should you want to learn more, you can find additional daemon documentation of the daemon at [/doc/tfchaind.md](/doc/tfchaind.md) and the (CLI) client on [doc/tfchainc.md](doc/tfchainc.md).

## testnet

A testnet for the threefold chain is currently deployed, and can be connected to using the binaries provided on the release page.
The testnet parameters are:

    - 10 minute block time
    - 1 million blockstakes created in the genesis block
    - 100 million coins created in the genesis block
    - 10 coins block reward
    - 144 block maturity delay for miner payouts (~ 1 day)

### receiving tokens for testnet

Currently there is no automated way for receiving tokens on the testnet. [An issue for it has been created](https://github.com/threefoldfoundation/tfchain/issues/12), but there is no ETA for an implemented solution.

Should you require X amount of tokens for the testnet, prior to this automated solution being available, you can request them by creating an issue on GitHub where you give your address as well as the amount of tokens you require. The title of your github issue should be in the form of "`Testnet token request for <address>`".

Once we completed that transfer we'll close your issue. Should you require more tokens in the future for the same address, you can simply re-open the issue and request more tokens tat way. Please do not recreate multiple token requests for the same address. Do also not link multiple addresses in a single issue, as that might get confusing.

## technical information

tfchain is using and build on top of [Rivine][rivine], a generic blockchain protocol, using Proof of Blockstake (PoB), rather than the also popular Proof of Work (PoW). It allows for custom blockchain implementations, hence tfchain is a custom implementation of the [Rivine][rivine] protocol.

This official (Golang) implementation is build using a vendored version of [the reference Golang implementation of Rivine][rivine].

For in-depth technical information you can check the [Rivine][rivine] docs at [github.com/rivine/rivine/tree/master/doc](https://github.com/rivine/rivine/tree/master/doc). There are no technical docs in this repository, as all the technology lives and is developed within the [Rivine repository][rivine].
