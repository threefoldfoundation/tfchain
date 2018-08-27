# Cross chain atomic swap walkthrough using the official Ethereum Go client.

## required tools

In order to execute atomic swaps as described in this document, you need to run the core tfchain daemon and client, available from <https://github.com/threefoldfoundation/tfchain/releases>, the [Official Go Ethereum full node client](https://geth.ethereum.org/downloads/) and the `ethatomicswap` tool that is available at <https://github.com/rivine/atomicswap/releases>.

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