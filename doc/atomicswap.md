# Cross chain atomic swaps

## Theory

A cross-chain swap is a trade between two users of different cryptocurrencies. For example, one party may send Threefold tokens to a second party's Threefold address, while the second party would send Bitcoin to the first party's Bitcoin address. However, as the blockchains are unrelated and transactions cannot be reversed, this provides no protection against one of the parties not honoring their end of the deal. One common solution to this problem is to introduce a mutually-trusted third party for escrow. An atomic cross-chain swap solves this problem without the need for a third party. On top of that it achieves waterproof validation without introducing the problems and complexities introduced by a escrow-based validation system.

Atomic swaps involve each party paying into a contract transaction, one contract for each blockchain. The contracts contain an output that is spendable by either party, but the rules required for redemption are different for each party involved. 

## required tools  
In order to execute atomic swaps as described in this document, you need to run the core tfchain daemon and client, available from https://github.com/threefoldfoundation/tfchain/releases and the decred atomic swap tools are available at https://github.com/rivine/decredatomicswap/releases. 

The original decred atomic swap project is [Decred atomic swaps](https://github.com/decred/atomicswap), the rivine fork just supplies the binaries.
 
## Example

Let's assume Bob wants to buy 567 TFT from Alice for 0.1234BTC

Bob creates a bitcoin address and Alice creates a threefold address.

Bob initiates the swap, he generates a 32-byte secret and hashes it
using the SHA256 algorithm, resulting in a 32-byte hashed secret.

Bob now creates a swap transaction, as a smart contract, and publishes it on the Bitcoin chain, it has 0.1234BTC as an output and the output can be redeemed (used as input) using either 1 of the following conditions:
- timeout has passed (48hours) and claimed by Bob's refund address;
- the money is claimed by Alice's registered address and the secret is given that hashes to the hashed secret created by Bob 

This means Alice can claim the bitcoin if she has the secret and if the atomic swap process fails, Bob can always reclaim it's btc after the timeout.

 Bob sends this contract and the transaction id of this transaction on the bitcoin chain to Alice, making sure he does not share the secret of course.

 Now Alice validates if everything is as agreed (=audit)after which She creates a similar transaction on the Rivine chain but with a timeout for refund of only 24 hours and she uses the same hashsecret as the first contract for Bob to claim the tokens.
 This transaction has 9876 tokens as an output and the output can be redeemed( used as input) using either 1 of the following conditions:
- timeout has passed ( 24hours) and claimed by the sellers refund address
- the secret is given that hashes to the hashsecret Bob created (= same one as used in the bitcoin swap transaction) and claimed by the buyers's address

In order for Bob to claim the threefold tokens, he has to use and as such disclose the secret.

The magic of the atomic swap lies in the fact that the same secret is used to claim the tokens in both swap transactions but it is not disclosed in the contracts because only the hash of the secret is used there. The moment Bob claims the threefold tokens, he discloses the secret and Alice has enough time lef to claim the bitcoin because the timeout of the first contract is longer than the one of the second contract.
Of course, either Bob or Alice can be the initiator or the participant.

### Technical details of the example
This example is a walkthrough of an actual atomic swap  on the threefold and bitcoin testnets.
 

Start bitcoin core qt in server mode on testnet: 
`./Bitcoin-Qt  -testnet -server -rpcuser=user -rpcpassword=pass -rpcport=18332`
Start the tfchain daemon on testnet:
`tfchaind --network testnet`

Alice creates a new bitcoin  address (as of bitcoin core 0.16, make sure to specify the 'legacy' address type since we need a p2pkh address)and provides this to Bob: 
```￼
getnewaddress "" legacy
muQ1J2UMfekrRJqEgXM59AuFm2az7y94V3
```

#### initiate step
Bob initiates the process by using btcatomicswap to pay 0.1234BTC into the Bitcoin contract using Alice's Bit coin address, sending the contract transaction, and sharing the secret hash (not the secret), and contract's transaction with Alice. The refund transaction can not be sent until the locktime expires, but should be saved in case a refund is necessary.

command:`btcatomicswap initiate <participant address> <amount>`

```
$ ./btcatomicswap --testnet --rpcuser=user --rpcpass=pass initiate muQ1J2UMfekrRJqEgXM59AuFm2az7y94V3  0.1234
Secret:      9cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257
Secret hash: 8b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b

Contract fee: 0.00001545 BTC (0.00003902 BTC/kB)
Refund fee:   0.00001462 BTC (0.00005024 BTC/kB)

Contract (2N22npr3a5JiSMpZDeeQeY3FVfjBwzvkwuW):
6382012088a8208b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b8876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba6704d8e3105bb17576a914e0bd0eb7e5c382b701f4783d07d9c16b43e4a84c6888ac

Contract transaction (b0a4d6a6423ae62d77e0ae79d5b6b39a8946cbcd9d3f127999f96c69bf67421f):
020000000001027cb7f7a8861524a5dfe6057b0142cf61abbb39319e573d0774a0e0bfad5c5635000000006b483045022100dcefcd54f9ef5bd7dfb670a1312a0890b719aa007f4787760781383e12602ea70220198de03654c7b2929973765785998a3a8f80e3bb8afd7977c684593bf24d9f24012102f2d0748df6ebaed3dbe2ae119abaf5a4f679799d1f84822e0c930244a844b8c6feffffff959344be6657bc9c668ad243d3eebb798c3ab87267bc7babd6ef954896d5bf1d00000000171600145bae19216eab65fd4c08c255d35d0c96731cfacdfeffffff02204bbc000000000017a914605f22bf839c23e7d511872a8e141d09d9f15009874d846c000000000017a9143629e8ffdfcaf9821d9f9ff29daf4ba301bed5ee87000247304402202e286013676c6c4223d2d3d65987c4143a31d9b4c953ef4a5ff560e80a41d7aa022060a01ce679ad66bcfe60d88d026c308797737557c41a845a815d5a5e74bf0832012103b19836fd942b5fca3808b83923557add6022a1e12e0d0c79dfb9adbbe6e9bb8b00000000

Refund transaction (324c8ee26ffb91e07f56844863baeff8a47c4cefe787205e2256c5e074b544ac):
02000000011f4267bf696cf99979123f9dcdcb46899ab3b6d579aee0772de63a42a6d6a4b000000000ce47304402204523c67fa03bbf4312b987dea262b65e69ffd9cdf7cfe34a18a8cc8af832975b02200966b0bc0a90f361320af291ff281c862a821e37d35895842d866da8af0321990121032e4cca08710c81e2c84ad8124b335c77a0920f5d0fcae09283a9095c640a23df004c616382012088a8208b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b8876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba6704d8e3105bb17576a914e0bd0eb7e5c382b701f4783d07d9c16b43e4a84c6888ac00000000016a45bc00000000001976a9141014ecc33fb66726ed9508b6caeb78fa3520ca3d88acd8e3105b

Publish contract transaction? [y/N] y
Published contract transaction (b0a4d6a6423ae62d77e0ae79d5b6b39a8946cbcd9d3f127999f96c69bf67421f)
```
You can check the transaction [on a bitcoin testnet blockexplorer](https://testnet.blockexplorer.com/tx/b0a4d6a6423ae62d77e0ae79d5b6b39a8946cbcd9d3f127999f96c69bf67421f) where you can see that 0.1234 BTC is sent to 2N22npr3a5JiSMpZDeeQeY3FVfjBwzvkwuW (= the contract script hash) being a [p2sh](https://en.bitcoin.it/wiki/Pay_to_script_hash) address in the bitcoin testnet. 


 #### audit contract

Bob sends Alice the contract and the contract transaction. Alice should now verify if
- the script is correct 
- the locktime is far enough in the future
- the amount is correct
- she is the recipient 

command:`btcatomicswap auditcontract <contract> <contract transaction>`

 ```
$ ./btcatomicswap --testnet --rpcuser=user --rpcpass=pass auditcontract 6382012088a8208b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b8876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba6704d8e3105bb17576a914e0bd0eb7e5c382b701f4783d07d9c16b43e4a84c6888ac 020000000001027cb7f7a8861524a5dfe6057b0142cf61abbb39319e573d0774a0e0bfad5c5635000000006b483045022100dcefcd54f9ef5bd7dfb670a1312a0890b719aa007f4787760781383e12602ea70220198de03654c7b2929973765785998a3a8f80e3bb8afd7977c684593bf24d9f24012102f2d0748df6ebaed3dbe2ae119abaf5a4f679799d1f84822e0c930244a844b8c6feffffff959344be6657bc9c668ad243d3eebb798c3ab87267bc7babd6ef954896d5bf1d00000000171600145bae19216eab65fd4c08c255d35d0c96731cfacdfeffffff02204bbc000000000017a914605f22bf839c23e7d511872a8e141d09d9f15009874d846c000000000017a9143629e8ffdfcaf9821d9f9ff29daf4ba301bed5ee87000247304402202e286013676c6c4223d2d3d65987c4143a31d9b4c953ef4a5ff560e80a41d7aa022060a01ce679ad66bcfe60d88d026c308797737557c41a845a815d5a5e74bf0832012103b19836fd942b5fca3808b83923557add6022a1e12e0d0c79dfb9adbbe6e9bb8b00000000
Contract address:        2N22npr3a5JiSMpZDeeQeY3FVfjBwzvkwuW
Contract value:          0.1234 BTC
Recipient address:       muQ1J2UMfekrRJqEgXM59AuFm2az7y94V3
Author's refund address: n21G8y9tSuaBRDk2yy8zzEgGzWPDqhwYzn

Secret hash: 8b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b

Locktime: 2018-06-01 06:12:40 +0000 UTC
Locktime reached in 47h57m25s
```

WARNING:
A check on the blockchain should be done as the auditcontract does not do that so an already spent output could have been used as an input. Checking if the contract has been mined in a block should suffice

#### Participate

Alice trusts the contract so she participates in the atomic swap by paying the tokens into a threefold token  contract using the same secret hash. 

Bob creates a new threefold address ( or uses an existing one): 
```
tfchainc wallet address
Created new address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
```

Bob sends this address to Alice who uses it to participate in the swap.
command:`tfchainc atomicswap participate <initiator address> <amount> <secret hash>`
```
$ ./tfchainc atomicswap participate 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb 567 8b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b
Contract address: 0202d7f00771cbaa0fa34481709004e953fa979881491f35b7de566e69579ba19114ec9b68d8bc
Contract value: 567 TFT
Receiver's address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
Sender's (contract creator) address: 0179fb6a617f52d60799fe610665b83b7372683201c06da24db54ad1878e5f1d8ff8c1b41ba3a2

SecretHash: 8b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b

TimeLock: 1527747605 (2018-05-31 08:20:05 +0200 CEST)
TimeLock reached in: 23h59m59.939846619s

Publish atomic swap (participation) transaction? [Y/N] Y
published contract transaction
OutputID: 15504a2eee64101ca04f4246bf9358db501397947a37b231454be66d1a3e5e7e
TransactionID: 7c5d9694b153945cea505a850b4a694ad6c2cb33d0f8d49c0416988783345436
``` 

The above command will create a transaction with `567` TFT as the Output  value of the output (`15504a2eee64101ca04f4246bf9358db501397947a37b231454be66d1a3e5e7e`). The output can be claimed by Bobs address (`01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb`)  and Bob will  to also have to provide the secret that hashes to the hashed secret `ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5`.

Alice now informs Bob that the Threefold contract transaction has been created and provides him with the contract details.

#### audit Threefold contract

Just as Alice had to audit Bob's contract, Bob now has to do the same with Alice's contract before withdrawing. 
Bob verifies if:
- the amount of threefold tokens () defined in the output is correct
- the attached script is correct
- the locktime, hashed secret (`ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5`) and wallet address, defined in the attached script, are correct

command:`./tfchainc atomicswap auditcontract outputid`
flags are available to automatically check the information in the contract.
```
$ ./tfchainc atomicswap auditcontract 15504a2eee64101ca04f4246bf9358db501397947a37b231454be66d1a3e5e7e
Atomic Swap Contract (condition) found:

Contract value: 567 TFT

Receiver's address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
Sender's (contract creator) address: 0179fb6a617f52d60799fe610665b83b7372683201c06da24db54ad1878e5f1d8ff8c1b41ba3a2
Secret Hash: 8b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b
TimeLock: 1527747605 (2018-05-31 08:20:05 +0200 CEST)
TimeLock reached in: 23h52m36.117162779s

Found Atomic Swap Contract is valid :)
```

The audit also checks if that the given contract's output   has not already been spend.

#### redeem tokens

Now that both Bob and Alice have paid into their respective contracts, Bob may withdraw from the Threefold contract. This step involves publishing a transaction which reveals the secret to Alice, allowing her to withdraw from the Bitcoin contract.

command:`/tfchainc atomicswap redeem outputid secret`

```
$ ./tfchainc atomicswap redeem 15504a2eee64101ca04f4246bf9358db501397947a37b231454be66d1a3e5e7e 9cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257
Contract address: 0202d7f00771cbaa0fa34481709004e953fa979881491f35b7de566e69579ba19114ec9b68d8bc
Contract value: 567 TFT
Receiver's address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
Sender's (contract creator) address: 0179fb6a617f52d60799fe610665b83b7372683201c06da24db54ad1878e5f1d8ff8c1b41ba3a2

SecretHash: 8b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b
Secret: 9cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257

TimeLock: 1527747605 (2018-05-31 08:20:05 +0200 CEST)
TimeLock reached in: 23h48m8.199200656s

Publish atomic swap redeem transaction? [Y/N] Y

Published atomic swap redeem transaction!
Transaction ID: 1d4428d7651710c9630a3c150277fd24c504f8a30e0a9d338f04e819aeed48db
>   NOTE that this does NOT mean for 100% you'll have the money!
> Due to potential forks, double spending, and any other possible issues your
> redeem might be declined by the network. Please check the network
> (e.g. using a public explorer node or your own full node) to ensure
> your payment went through. If not, try to audit the contract (again).
```

#### redeem bitcoins

Now that Bob has withdrawn from the threefold contract and revealed the secret. If bob is really nice he could simply give the secret to Alice. However,even if he doesn't do this Alice can extract the secret from this redemption transaction. Alice may watch a block explorer to see when the threefold contract output was spent and look up the redeeming transaction.

Alice can automatically extract the secret from the input where it is used by Bob, by simply giving the outputID of the contract. Either you do this using a public web-based explorer, by looking up the outputID as hash. Or you let the command line client do it automatically for you:

command:`tfchainc atomicswap extractsecret outputid`
```
$./tfchainc atomicswap extractsecret 1d4428d7651710c9630a3c150277fd24c504f8a30e0a9d338f04e819aeed48db
atomic swap contract was redeemed by participator
extracted secret: 9cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257
```

NOTE: in this call I gave a public explorer address as I have no explorer node running myself.
Therefore I can use a public explorer to look it up for me instead.
Should you have a local explorer node running on the default address, you can simply omit the flag and use 
`$tfchainc extractsecret abcdef01234567890abcdef01234567890abcdef01234567890abcdef0123452` .

With the secret known (extracted from the coinInput with parent OutputID `9cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257`), Alice may redeem from Bob's Bitcoin contract:
command: `btcatomicswap redeem <contract> <contract transaction> <secret>`
```
./btcatomicswap --testnet --rpcuser=user --rpcpass=pass redeem  6382012088a8208b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b8876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba6704d8e3105bb17576a914e0bd0eb7e5c382b701f4783d07d9c16b43e4a84c6888ac  020000000001027cb7f7a8861524a5dfe6057b0142cf61abbb39319e573d0774a0e0bfad5c5635000000006b483045022100dcefcd54f9ef5bd7dfb670a1312a0890b719aa007f4787760781383e12602ea70220198de03654c7b2929973765785998a3a8f80e3bb8afd7977c684593bf24d9f24012102f2d0748df6ebaed3dbe2ae119abaf5a4f679799d1f84822e0c930244a844b8c6feffffff959344be6657bc9c668ad243d3eebb798c3ab87267bc7babd6ef954896d5bf1d00000000171600145bae19216eab65fd4c08c255d35d0c96731cfacdfeffffff02204bbc000000000017a914605f22bf839c23e7d511872a8e141d09d9f15009874d846c000000000017a9143629e8ffdfcaf9821d9f9ff29daf4ba301bed5ee87000247304402202e286013676c6c4223d2d3d65987c4143a31d9b4c953ef4a5ff560e80a41d7aa022060a01ce679ad66bcfe60d88d026c308797737557c41a845a815d5a5e74bf0832012103b19836fd942b5fca3808b83923557add6022a1e12e0d0c79dfb9adbbe6e9bb8b00000000 9cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257
Redeem fee: 0.00022075 BTC (0.00067923 BTC/kB)

Redeem transaction (8c77722af341c56f968175c09c57f70841e42d2651b456dfa88a04d63046df76):
02000000011f4267bf696cf99979123f9dcdcb46899ab3b6d579aee0772de63a42a6d6a4b000000000f0483045022100e4f416918bc7004326d438b5d59f10289f32c6524e15e1f40eb6eefbc0c5550202204c954ae1abe1db45f7e91a1bfb0c659ac6f2731ff07c3fe485ecccaeb1ee16ce0121023d3bdd190a65b033905a5a521596598518ac884b6b974f4395fe3e87525b253b209cddc24ba8e77d868c97e98374f4a2447aab114fa6f62a35d53f636c092f5257514c616382012088a8208b445001958277e6372424625d31e649e32812eeb62eece03ff616a31ebd0f6b8876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba6704d8e3105bb17576a914e0bd0eb7e5c382b701f4783d07d9c16b43e4a84c6888acffffffff01e5f4bb00000000001976a91404458d8235eaaa929bf5af362f06196fee01009488acd8e3105b

Publish redeem transaction? [y/N] y
Published redeem transaction (8c77722af341c56f968175c09c57f70841e42d2651b456dfa88a04d63046df76)
```
This transaction can be verified [on a bitcoin testnet blockexplorer](https://testnet.blockexplorer.com/tx/8c77722af341c56f968175c09c57f70841e42d2651b456dfa88a04d63046df76) .
The cross-chain atomic swap is now completed and successful.

## Ethereum Example

This example is a walkthrough of an actual atomic swap on the threefold and ethereum (Rinkeby) testnets.
In our example Bob wants to buy 25000 TFT from Alice for 5.41 ETH.

First testnet TFT-ETH atomic swap (as created while documenting this example):

| Description | Link to (raw) transaction |
| - | - |
| Ethereum contract created by A | [e5cb8162f6d4e5948ee9c7690aba4c5a641c9547aec3f90cb5ffbc30f513c300](https://rinkeby.etherscan.io/tx/0xe5cb8162f6d4e5948ee9c7690aba4c5a641c9547aec3f90cb5ffbc30f513c300) |
| Threefold contract created by B | [37e2ffea404d5dd053e952a06182552c7b3a9d5ee672dd2501ebd0c868c9f215](https://explorer.testnet.threefoldtoken.com/hash.html?hash=37e2ffea404d5dd053e952a06182552c7b3a9d5ee672dd2501ebd0c868c9f215) |
| A's Threefold redemption | [e3d4558124736231f79df3c1c1df36c737002d6c10b1097cc098ce19273d3e12](https://explorer.testnet.threefoldtoken.com/hash.html?hash=e3d4558124736231f79df3c1c1df36c737002d6c10b1097cc098ce19273d3e12) |
| B's Ethereum redemption | [9a146166bb61573fdeabc9a02da5a6a739e8e93c560b63e8b1faabe3d9a59767](https://rinkeby.etherscan.io/tx/0x9a146166bb61573fdeabc9a02da5a6a739e8e93c560b63e8b1faabe3d9a59767) |

Not everything is explained in detail, as by now it is assumed that you know already how an atomic swap works,
on a conceptual level as well as how it is to be done for TFT using the `tfchainc` tool.

### Setup

Start `geth` in server mode on the Rinkeby testnet with a console attached to it:

```
geth --networkid=4 --datadir=./data --syncmode=light --ethstats='yournode:Respect my authoritah!@stats.rinkeby.io' --bootnodes=enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303 --rpc console
```

> Here we assumes you have a configured daemon as well as an account configured already.
>
> + You can learn how to connect yourself to the Rinkeby testnet at <https://www.rinkeby.io/#geth>:
>   + A light client is sufficient in order to reproduce this example
> + You can learn how to create an account at <https://github.com/ethereum/go-ethereum/wiki/Managing-your-accounts#creating-an-account>:
>   + Make sure to use the same flags where needed (`--datadir`) as you do when starting `geth` in server-mode;

Once your Ethereum daemon is running make sure to unlock the accounts (e.g. `0x3cc65c21435d484d073b8abd361ab852f26cc504`) to be used.
You can do so from the console which you attached while starting the `geth` in server-mode:

```
> personal.unlockAccount("0x3cc65c21435d484d073b8abd361ab852f26cc504")
Unlock account 0x3cc65c21435d484d073b8abd361ab852f26cc504
Passphrase: 
true
```

> You can find more information about this command and more at: <https://github.com/ethereum/go-ethereum/wiki/Management-APIs#examples-1>

It is assumed that you already know how to start the ThreeFold Chain Daemon (`tfchaind`) and how to create and unlock a wallet.

### Walkthrough

What follows is a full Walkthrough of a TFT-ETH atomic swap on testnet (Rinkeby).
The process is however exactly the same for mainnet. The only difference would be
that you have to omit the `-testnet` flag while using the `ethatomicswap` tool,
as well that you start your `tfchaind` proces without the `--network` flag.

#### Initiate (ETH)

Bob initiates the process by using `ethatomicswap` to pay 5.41 ETH into the Bitcoin contract using Alice's Ethereumaddress, sending the contract transaction, and sharing the secret hash (not the secret), and contract's transaction with Alice. The refund transaction can not be sent until the locktime expires, but should be saved in case a refund is necessary.

command: `ethatomicswap initiate <participant address> <amount>`

```
$ ethatomicswap -testnet initiate 3701285fd20e0556c3a51bb6fabda5be25d0c40f 5.41
Amount: 5410000000000000000 Wei (5.41 ETH)

Secret:      5486225923792fd14b2e31317381154fba96e6c88f7cc6e7dd758b3c8d26c642
Secret hash: 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd

Author's refund address: 3cc65c21435d484d073b8abd361ab852f26cc504

Contract fee: 0.000175469 ETH
Refund fee:   0.00021 ETH (max)

Chain ID:         4
Contract Address: 2661cbaa149721f7c5fab3fa88c1ea564a683631
Contract transaction (e5cb8162f6d4e5948ee9c7690aba4c5a641c9547aec3f90cb5ffbc30f513c300):
f8d152843b9aca008302ad6d942661cbaa149721f7c5fab3fa88c1ea564a683631884b142e562add0000b864ae052147000000000000000000000000000000000000000000000000000000000002a300322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd0000000000000000000000003701285fd20e0556c3a51bb6fabda5be25d0c40f2ca057e2bfe7767d49046642c33c4c899208f34f755bbf00d2601757dcc1082c87a6a011a9c6eac94d2c3bf25f318885744fc3ea3856b4bd58e6bc4bb8abf0c70e391a

Publish contract transaction? [y/N] y
Published contract transaction (e5cb8162f6d4e5948ee9c7690aba4c5a641c9547aec3f90cb5ffbc30f513c300)
```

You can check the transaction [on the Rinkeby (etherscan) testnet block explorer](https://rinkeby.etherscan.io/tx/0xe5cb8162f6d4e5948ee9c7690aba4c5a641c9547aec3f90cb5ffbc30f513c300) where you can see that 5.41 ETH is sent to 0x2661cbaa149721f7c5fab3fa88c1ea564a683631 (= the smart contract's address, deployed on the Rinkeby testnet by us).


#### Audit Contract (ETH)

Bob sends Alice the contract and the contract transaction. Alice should now verify if
- the correct smart contract is used
- the locktime is far enough in the future
- the amount is correct
- she is the recipient

command:`ethatomicswap auditcontract <contract transaction>`

 ```
$ ethatomicswap -testnet auditcontract f8d152843b9aca008302ad6d942661cbaa149721f7c5fab3fa88c1ea564a683631884b142e562add0000b864ae052147000000000000000000000000000000000000000000000000000000000002a300322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd0000000000000000000000003701285fd20e0556c3a51bb6fabda5be25d0c40f2ca057e2bfe7767d49046642c33c4c899208f34f755bbf00d2601757dcc1082c87a6a011a9c6eac94d2c3bf25f318885744fc3ea3856b4bd58e6bc4bb8abf0c70e391a
Contract address:        2661cbaa149721f7c5fab3fa88c1ea564a683631
Contract value:          5.41 ETH
Recipient address:       3701285fd20e0556c3a51bb6fabda5be25d0c40f
Author's refund address: 3cc65c21435d484d073b8abd361ab852f26cc504

Secret hash: 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd

Locktime: 2018-07-22 12:41:34 +0000 UTC
Locktime reached in 47h53m28s
```

It also already checks for you if this transaction was succesfully registered and exists on the blockchain.

WARNING: it does not check for you if the contract has not yet been redeemed/refunded.

#### Participate (TFT)

Alice trusts the contract so she participates in the atomic swap by paying the tokens (25000 TFT as agreed) into a threefold token contract using the same secret hash. 

command:`tfchainc atomicswap participate <initiator address> <amount> <secret hash>`

```
$ tfchainc atomicswap participate 01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165 25000 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd
Contract address: 02ed044fa056e6fda24261628576d61f7ed580a649496a582ac0691f90f4ef89002a1a5f368c57
Contract value: 25000 TFT
Receiver's address: 01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165
Sender's (contract creator) address: 01d8a5b1f0e92d5c333b8368189c8e09c4e7aea7e7186967035242c6ee94bfa8f9cf19725d5e96

SecretHash: 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd

TimeLock: 1532177577 (2018-07-21 14:52:57 +0200 CEST)
TimeLock reached in: 23h59m59.080485111s
Publish atomic swap transaction? [Y/N] y

published contract transaction

OutputID: 19c162756112c6d950593c8d740999d976a9743390415ee870104f336551d73c
TransactionID: 37e2ffea404d5dd053e952a06182552c7b3a9d5ee672dd2501ebd0c868c9f215

Contract Info:

Contract address: 02ed044fa056e6fda24261628576d61f7ed580a649496a582ac0691f90f4ef89002a1a5f368c57
Contract value: 25000 TFT
Receiver's address: 01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165
Sender's (contract creator) address: 01d8a5b1f0e92d5c333b8368189c8e09c4e7aea7e7186967035242c6ee94bfa8f9cf19725d5e96

SecretHash: 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd

TimeLock: 1532177577 (2018-07-21 14:52:57 +0200 CEST)
TimeLock reached in: 23h59m56.822133808s
```

You can check the transaction [on the official tfchain block explorer](https://explorer.testnet.threefoldtoken.com/hash.html?hash=37e2ffea404d5dd053e952a06182552c7b3a9d5ee672dd2501ebd0c868c9f215) where you can see that 25000 TFT is sent to 01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165 (= the TFT wallet address of Bob).

Alice now informs Bob that the Threefold contract transaction has been created and provides him with the contract details.

#### Audit Contract (TFT)

Just as Alice had to audit Bob's contract, Bob now has to do the same with Alice's contract before withdrawing. 
Bob verifies if:
- the amount of threefold tokens defined in the output is correct
- the locktime, secret hash (`322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd`) and wallet address, defined in the attached condition, are correct

command: `tfchainc atomicswap auditcontract outputid`

> Tip: flags are available to automatically check the information in the contract.
> You can see what those flags are and how to use them with
> `tfchainc atomicswap auditcontract --help`.

```
$ tfchainc atomicswap auditcontract 19c162756112c6d950593c8d740999d976a9743390415ee870104f336551d73c
Atomic Swap Contract (condition) found:

Contract value: 25000 TFT

Receiver's address: 01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165
Sender's (contract creator) address: 01d8a5b1f0e92d5c333b8368189c8e09c4e7aea7e7186967035242c6ee94bfa8f9cf19725d5e96
Secret Hash: 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd
TimeLock: 1532177577 (2018-07-21 14:52:57 +0200 CEST)
TimeLock reached in: 23h55m36.579063567s

found Atomic Swap Contract is valid
```

The audit also checks if that the given contract exists in the consensus set
(meaning it is registered on the blockchain) and that its output has not already been spend.

#### Redeem (TFT)

Now that both Bob and Alice have paid into their respective contracts, Bob may withdraw from the Threefold contract. This step involves publishing a transaction which reveals the secret to Alice, allowing her to withdraw from the Ethereum contract.

command: `tfchainc atomicswap redeem outputid secret`

```
$ tfchainc atomicswap redeem 19c162756112c6d950593c8d740999d976a9743390415ee870104f336551d73c 5486225923792fd14b2e31317381154fba96e6c88f7cc6e7dd758b3c8d26c642
Contract address: 02ed044fa056e6fda24261628576d61f7ed580a649496a582ac0691f90f4ef89002a1a5f368c57
Contract value: 25000 TFT
Receiver's address: 01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165
Sender's (contract creator) address: 01d8a5b1f0e92d5c333b8368189c8e09c4e7aea7e7186967035242c6ee94bfa8f9cf19725d5e96

SecretHash: 322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd
Secret: 5486225923792fd14b2e31317381154fba96e6c88f7cc6e7dd758b3c8d26c642

TimeLock: 1532177577 (2018-07-21 14:52:57 +0200 CEST)
TimeLock reached in: 23h52m8.284260937s
Publish atomic swap redeem transaction? [Y/N] y

published atomic swap redeem transaction
transaction ID: e3d4558124736231f79df3c1c1df36c737002d6c10b1097cc098ce19273d3e12
>   Note that this does not mean for 100% you'll have the money.
> Due to potential forks, double spending, and any other possible issues your
> redeem might be declined by the network. Please check the network
> (e.g. using a public explorer node or your own full node) to ensure
> your payment went through. If not, try to audit the contract (again).
```

You can check the transaction [on the official tfchain block explorer](https://explorer.testnet.threefoldtoken.com/hash.html?hash=e3d4558124736231f79df3c1c1df36c737002d6c10b1097cc098ce19273d3e12) where you can see that 25000 TFT (minus the transaction fee) is claimed by Bob (01a75b03da048b933d2d04cc22283c170eb8300f1806b02bf7138f4092a3a385703864e427d165).

#### Redeem (ETH)

Now that Bob has withdrawn from the tfchain contract and revealed the secret. If bob is really nice he could simply give the secret to Alice. However, even if he doesn't do this Alice can extract the secret from this redemption transaction. Alice may watch a block explorer to see when the threefold contract output was spent and look up the redeeming transaction.

Alice can automatically extract the secret from the input where it is used by Bob, by simply giving the transactionID of the contract. Either you do this using a public web-based explorer, by looking up the transactionID as hash. Or you let the command line client do it automatically for you:

command: `tfchainc atomicswap extractsecret transactionID`

> Tip: the `--secrethash` flag is available which allows you to automatically
> validate the extracted secret against the expected secrethash.

```
$ tfchainc atomicswap extractsecret 
atomic swap contract was redeemed
extracted secret: 5486225923792fd14b2e31317381154fba96e6c88f7cc6e7dd758b3c8d26c64
```

With the secret known, Alice has everything she needs
in order to redeem from Bob's Ethereum contract:

command: `ethatomicswap redeem <contract transaction> <secret>`

```
ethatomicswap -testnet redeem f8d152843b9aca008302ad6d942661cbaa149721f7c5fab3fa88c1ea564a683631884b142e562add0000b864ae052147000000000000000000000000000000000000000000000000000000000002a300322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd0000000000000000000000003701285fd20e0556c3a51bb6fabda5be25d0c40f2ca057e2bfe7767d49046642c33c4c899208f34f755bbf00d2601757dcc1082c87a6a011a9c6eac94d2c3bf25f318885744fc3ea3856b4bd58e6bc4bb8abf0c70e391a 5486225923792fd14b2e31317381154fba96e6c88f7cc6e7dd758b3c8d26c64
Redeem fee: 0.000065077 ETH

Chain ID:         4
Contract Address: 2661cbaa149721f7c5fab3fa88c1ea564a683631
Redeem transaction (9a146166bb61573fdeabc9a02da5a6a739e8e93c560b63e8b1faabe3d9a59767):
f8a81c843b9aca0082fe35942661cbaa149721f7c5fab3fa88c1ea564a68363180b844b31597ad5486225923792fd14b2e31317381154fba96e6c88f7cc6e7dd758b3c8d26c642322382392c5c402321a0507c9a8fce1c2599e72b631eff4f207104e4f225acfd2ca0b7c5960e73ea953bb45cbac5848533c0c5e7084295b852ff562ff540598a2a32a023bf4be9e26d5f45ef9734b07dc9e6f3387ad59d000fb744bdae0e22d75d4746

Publish redeem transaction? [y/N] y
Published redeem transaction (9a146166bb61573fdeabc9a02da5a6a739e8e93c560b63e8b1faabe3d9a59767)
```

This transaction can be verified [on the Rinkeby (etherscan) testnet block explorer](https://rinkeby.etherscan.io/tx/0x9a146166bb61573fdeabc9a02da5a6a739e8e93c560b63e8b1faabe3d9a59767).

The cross-chain atomic swap is now completed and successful.
Bob has bought 25000 TFT from ALice for 5.41 ETH, minus the redeem transaction fees.

## References

Rivine atomic swaps are an implementation of [Decred atomic swaps](https://github.com/decred/atomicswap).

[Bitcoin scripts and opcodes](https://en.bitcoin.it/wiki/Script)
