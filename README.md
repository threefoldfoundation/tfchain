# tfchain

tfchain is the official Go implementation of the ThreeFold blockchain client. It uses and is build on top of the [Rivine][rivine] protocol.

The ThreeFold blockchain is used as a ledger for the ThreeFold Token ("TFT"), a digital currency backed by neutral and sustainable internet capacity. You can learn more about the ThreeFold token on [threefoldtoken.com](https://threefoldtoken.com).

The ThreeFold blockchain is the technology behind the ThreeFold Grid ("TF Grid" or "Grid"), a grid where your data is local and controlled by you, shaping a new neutral and sustainable Internet. To learn more about this grid you can find and read the whitepaper at [threefoldtoken.com/pdf/tf_whitepaper.pdf](https://threefoldtoken.com/pdf/tf_whitepaper.pdf).

[rivine]: http://github.com/rivine/rivine

## testnet

A testnet for the threefold chain is currently deployed, and can be connected to using the binaries provided on the release page.
The testnet parameters are:

    - 10 minute block time
    - 1 million blockstakes created in the genesis block
    - 100 million coins created in the genesis block
    - 10 coins block reward
    - 144 block maturity delay for miner payouts (~ 1 day)

### Receiving tokens for testnet

Currently there is no automated way for receiving tokens on the testnet. [An issue for it has been created](https://github.com/threefoldfoundation/tfchain/issues/12), but there is no ETA for an implemented solution.

Should you require X amount of tokens for the testnet, prior to this automated solution being available, you can request them by creating an issue on GitHub where you give your address as well as the amount of tokens you require. The title of your github issue should be in the form of "`Testnet token request for <address>`".

Once we completed that transfer we'll close your issue. Should you require more tokens in the future for the same address, you can simply re-open the issue and request more tokens tat way. Please do not recreate multiple token requests for the same address. Do also not link multiple addresses in a single issue, as that might get confusing.

## technical information

tfchain is using and build on top of [Rivine][rivine], a generic blockchain protocol, using Proof of Blockstake (PoB), rather than the also popular Proof of Work (PoW). It allows for custom blockchain implementations, hence tfchain is a custom implementation of the [Rivine][rivine] protocol.

This official (Golang) implementation is build using a vendored version of [the reference Golang implementation of Rivine][rivine].

For in-depth technical information you can check the [Rivine][rivine] docs at [github.com/rivine/rivine/tree/master/doc](https://github.com/rivine/rivine/tree/master/doc). There are no technical docs in this repository, as all the technology lives and is developed within the [Rivine repository][rivine].

## Install ##

Start by cloning the tfchain repository (git clone).

In the Makefile there are various subcomands to make things easier.
* xc , creates binaries for Windows Mac and Linux.
* release-images, creates a docker image with a copy of the binaries. Needs credentials from [Itsyouonline](https://itsyou.online/)

In any machine with a golang 1.9 (or higher) development environment you can just do

```bash
go get github.com/threefoldfoundation/tfchain/cmd/tfchaind
go get github.com/threefoldfoundation/tfchain/cmd/tfchainc
```

and that will download the needed code, compile it and install it in $GOROOT/bin or $GOPATH/bin

You can find additional daemon documentation of the daemon on [tfchaind](http://github.com/threefoldfoundation/tfchain/doc/tfchaind.md) and the client on [tfchainc](http://github.com/threefoldfoundation/tfchain/doc/tfchainc.md)