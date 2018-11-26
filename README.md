# tfchain [![Build Status](https://travis-ci.org/threefoldfoundation/tfchain.svg?branch=master)](https://travis-ci.org/threefoldfoundation/tfchain)

tfchain is the official Go implementation of the ThreeFold blockchain client. It uses and is build on top of the [Rivine][rivine] protocol.

The ThreeFold blockchain is used as a ledger for the ThreeFold Token ("TFT"), a digital currency backed by neutral and sustainable internet capacity. You can learn more about the ThreeFold token on [threefoldtoken.com](https://threefoldtoken.com).

The ThreeFold blockchain is the technology behind the ThreeFold Grid ("TF Grid" or "Grid"), a grid where your data is local and controlled by you, shaping a new neutral and sustainable Internet. To learn more about this grid you can find and read the white paper at [http://whitepaper.threefold.me](http://whitepaper.threefold.me).

> ⚠ Attention ⚠
>
> `v1.0.8` is deprecated as of the 6th of September, 2018.
> Please upgrade to `v1.1.0`: [threefoldfoundation/tfchain@v1.1.0](https://github.com/threefoldfoundation/tfchain/releases/tag/v1.1.0).
>
> If you are still running `v1.0.8` or older, you are advised to upgrade to `v1.1.0` as soon as you can or risk getting out-of-sync with the rest of the network due to your ConsensusSet not accepting two new transaction types.

You can reach out to us for any development-related questions on [our public Telegram group](https://t.me/rivine). If you're a user of the tools distributed via this repository, you are welcome to join that chat as well. You're especially welcome if you want to contribute and wonder how.

[rivine]: http://github.com/threefoldtech/rivine

## install and use tfchain

You have 2 easy options to install tfchain in your work environment, as a prequisite to using the tfchain binaries and thus become a full tfchain node. Currently there is only a test network.

### tfchain docker container

The easiest is pulling and using the latest prebuilt docker container:

```bash
docker pull tfchain/tfchain
docker run -d --name tfchain tfchain/tfchain
```

This will pull and configure the latest tfchain container, and it will start a container using that image in the program, named `tfchaind`.

Should you want to use to use the CLI client, you can do so doing the running container:

```bash
docker exec -ti tfchain /tfchainc 
```

Note that this minimal tfchain docker container is not really meant to be used for much tfchain CLI client interaction. If you wish to do a lot of that, it is probably more easy/useful to run a the tfchain binary CLI from your host machine. If you are interested in that, you can check out [the "tfchain from source" section](#tfchain-from-source).

Even though the image for these containers is prebuilt available for you, should you wish to use a tfchain docker container from the hacked source code you can rebuilt those imagines using `docker-minimal`.

### tfchain from source

tfchain is developed and implemented using [Golang](http://golang.org). Using the golang toolchain it is very easy to download, update and install the tfchain binaries used to run a full node and interact with it:

```bash
go get -u github.com/threefoldfoundation/tfchain/cmd/... && \
    tfchaind &
```

> tfchain supports Go 1.9 and above. Older versions of Golang may work but aren't supported.

At this point (if all went right) you should have a tfchain daemon running in the background which is syncing with the test net. You can follow this syncing process using the CLI client: `tfchainc`.

Should you want to learn more, you can find additional daemon documentation of the daemon at [/doc/tfchaind.md](/doc/tfchaind.md) and the (CLI) client on [doc/tfchainc.md](doc/tfchainc.md).

### tfchain development and deployment

If you are a tfchain developer you probably also want to pull all submodules locally,
to do so, you can get this entire repo and all its submodules locally using the following command:

```bash
go get -d -u github.com/threefoldfoundation/tfchain && \
    cd $GOPATH/src/github.com/threefoldfoundation/tfchain && \
    submodule update --init --recursive
```

Once you have that repo, you can pull from than on the latest commit
(including submodule updates) using the following command:

```bash
cd $GOPATH/src/github.com/threefoldfoundation/tfchain && make update
```

> NOTE that if you didn't execute `submodule update --init --recursive` when cloning this repo
> `make update` will fail, as submodules have to initialized first.

## standard (net)

By default a tfchain daemon (`tfchaind`) will connect to the standard net(work).

This is the official network for ThreeFoldTokens.
You can learn more about the ThreeFold token on [threefoldtoken.com](https://threefoldtoken.com).

The standard net has following properties:

+ A new block every 2 minutes (on average);
+ Block stakes can be used 1 day after receiving;
+ One TFT equals to 10<sup>9</sup> of the smallest currency unit;
+ Payouts take roughly one day to mature;

All properties can be found in: [/cmd/tfchaind/main.go](/cmd/tfchaind/main.go) (in the `getStandardnetGenesis` function body).

The standard (net) uses the following bootstrap nodes:

+ bootstrap1.threefoldtoken.com:23112
+ bootstrap2.threefoldtoken.com:23112
+ bootstrap3.threefoldtoken.com:23112
+ bootstrap4.threefoldtoken.com:23112

Web-ui explorers for the standard (net) are available at https://explorer.threefoldtoken.com and https://explorer2.threefoldtoken.com (if one isn't available try the other one).

## testnet

A public testnet is also provided in the tfchaind binary. It can be connected to using the `--network testnet` flags.
This way you can connect to a real network, which behaves almost like the real network. It should be noted that this network
is used to test new features and bug fixes before moving them to the production network. It can also be reset by the developers at any time.

You can request `300 TFT` to play with on tfchain wallet, by pasting your wallet address at <http://faucet.testnet.threefoldtoken.com> and requesting the 300 TFT. You'll receive your
_testnet_ tokens as soon as the next block is created and added to the _testnet_ chain.

The testnet has the following properties:

+ 3000 blockstakes;
+ An average block time of 2 minutes;
+ One TFT equals to 10<sup>9</sup> of the smallest currency unit;
+ Block rewards and miner payouts take rougly one day to mature;

The blockstakes have been distributed equally to 3 block creator nodes. This means that it is not possible to create blocks yourself
at this time.

The following bootstrap nodes are being used by the testnet:

+ bootstrap1.testnet.threefoldtoken.com:23112
+ bootstrap2.testnet.threefoldtoken.com:23112
+ bootstrap3.testnet.threefoldtoken.com:23112
+ bootstrap4.testnet.threefoldtoken.com:23112

A web-ui explorer for the testnet is available at: https://explorer.testnet.threefoldtoken.com and https://explorer2.testnet.threefoldtoken.com (if one isn't available try the other one).

## devnet

Should you require a local devnet you can do so by starting the daemon
using the `--network devnet` and `--no-bootstrap` flags.
This will allow you to use the `devnet` network` and mine as soon as you have the block stakes requried to do so.
For obvious reasons no bootstrap nodes are required or even available for this network.

Once your daemon is up and running you can give your own wallet all genesis coins,
by loading the following mnemonic as your seed:

```
carbon boss inject cover mountain fetch fiber fit tornado cloth wing dinosaur proof joy intact fabric thumb rebel borrow poet chair network expire else
```

Should you want to do this with the provided `tfchainc` wallet you would have to do following steps:

```
$ tfchainc wallet unlock    # give passphrase
$ tfchainc wallet load seed # give passphrase and above mnemonic
$ fchainc wallet addresses # reload all our default addresses again
$ tfchainc stop # stops the tfchaind daemon and than manually restart it again
$ tfchainc wallet unlock    # give passphrase
$ tfchainc wallet  # should show that you have 100M coins and 3K block stakes
Wallet status:
Encrypted, Unlocked
Confirmed Balance:   100000000 TFT
Unconfirmed Delta:   + 0 TFT
BlockStakes:         3000 BS
```

Should you use another wallet/client, the steps might be different,
which is fine as long as you use the menmonic given above as seed,
as the genesis block stakes and coins are attached to that one.

## technical information

tfchain is using and build on top of [Rivine][rivine], a generic blockchain protocol, using Proof of Blockstake (PoB), rather than the also popular Proof of Work (PoW). It allows for custom blockchain implementations, hence tfchain is a custom implementation of the [Rivine][rivine] protocol.

This official (Golang) implementation is build using a vendored version of [the reference Golang implementation of Rivine][rivine].

For in-depth technical information you can check the [Rivine][rivine] docs at [github.com/threefoldtech/rivine/tree/master/doc](https://github.com/threefoldtech/rivine/tree/master/doc). There are no technical docs in this repository, as all the technology lives and is developed within the [Rivine repository][rivine].

## troubleshooting

If you run into troubles while using the official `tfchaind` binary, try to find a solution for your problem in [the Rivine troubleshooting FAQ](https://github.com/threefoldtech/rivine#troubleshooting) first. The tfchain daemon is built using the Rivine daemon package, so chances are your problem has already been solved when another user ran into the same or similar-enough issue.

If your problem is not related to the official `tfchaind` binary, or if you do not find a solution in [the Rivine troubleshooting FAQ](https://github.com/threefoldtech/rivine#troubleshooting), you should start looking into the `tfchain` issues, found at <https://github.com/threefoldfoundation/tfchain/issues?utf8=✓&q=is%3Aissue+is%3Aopen+is%3Aclosed>. You can search faster by guessing some keyword(s) that might be included in such an issue. Did you find an issue that is already talking about your problem? Post a comment in it describing your environment, version, situation, steps on how to reproduce and any info you think might be helpful. It does not matter if the issue is already closed, if you run into a problem, it is proof that the problem (still) exists.

Should you not find an issue (open or closed) for your problem, you can create a new issue. Make sure to include your environment, version, situation, steps on how to reproduce and any info you think might be helpful.

Please use GitHub issues as much as possible, it keeps everything tracked. Only reach out to the `tfchain` maintainers if you really have to or in case it is not issue-related. The team is small, and direct support is pretty time-expensive. 
