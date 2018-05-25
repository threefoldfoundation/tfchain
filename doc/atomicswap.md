# Cross chain atomic swaps

## Theory

A cross-chain swap is a trade between two users of different cryptocurrencies. For example, one party may send Threefold tokens to a second party's Threefold address, while the second party would send Bitcoin to the first party's Bitcoin address. However, as the blockchains are unrelated and transactions cannot be reversed, this provides no protection against one of the parties not honoring their end of the deal. One common solution to this problem is to introduce a mutually-trusted third party for escrow. An atomic cross-chain swap solves this problem without the need for a third party. On top of that it achieves waterproof validation without introducing the problems and complexities introduced by a escrow-based validation system.

Atomic swaps involve each party paying into a contract transaction, one contract for each blockchain. The contracts contain an output that is spendable by either party, but the rules required for redemption are different for each party involved. 

## required tools  
In order to execute atomic swaps as described in this document, you need to run the core tfchain daemon and client, available from https://github.com/threefoldfoundation/tfchain/releases and the decred atomic swap tools are available at https://github.com/rivine/atomicswap/releases. 

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

## Technical details of the example
This example is a walkthrough of an actual atomic swap  on the threefold and bitcoin testnets.
 

Start bitcoin core qt in server mode on testnet: 
`./Bitcoin-Qt  -testnet -server -rpcuser=user -rpcpassword=pass -rpcport=18332`
Start the tfchain daemon on testnet:
`tfchaind --network testnet`

Alice creates a new bitcoin  address (as of bitcoin core 0.16, make sure to specify the 'legacy' address type since we need a p2pkh address): 
```￼
getnewaddress "" legacy
muQ1J2UMfekrRJqEgXM59AuFm2az7y94V3
```

### initiate step
Bob initiates the process by using btcatomicswap to pay 0.1234BTC into the Bitcoin contract using Alice's Bit coin address, sending the contract transaction, and sharing the secret hash (not the secret), and contract's transaction with Alice. The refund transaction can not be sent until the locktime expires, but should be saved in case a refund is necessary.

command:`btcatomicswap initiate <participant address> <amount>`

```
`Secret:      83639420c8683e288152f55b73598bc33f0ce6712316a8edc6e9772787845d2f
Secret hash: ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5

Contract fee: 0.00000316 BTC (0.00000798 BTC/kB)
Refund fee:   0.00000299 BTC (0.00001024 BTC/kB)

Contract (2N2b1Lcs2ic3kt1bj8aoHUDhy9YawY47vhd):
6382012088a820ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa58876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba67041e990a5bb17576a9143a1e96e746e5d673e53f8cbca956f38e230f081e6888ac

Contract transaction (1dbfd5964895efd6ab7bbc6772b83a8c79bbeed343d28a669cbc5766be449395):
02000000000102215d7d1855d16bf9d3c9c349ef8d0b7b58cb5c89922148bf4912d462ff4caef6010000001716001437c1c3347cb976780fe386f0631582fb751e21d8feffffff9400ee75b3ea0600ce429ecb9a0d5d6565e027d9d8dfd051285388169a676ed8000000006b483045022100dc32890b45892dbfa313fddc67e8c176228f82e10635ef56905bdd307f7a2bfa022069a2b8e96f550a005e4650a8ac0704e0436a902be97bd7ce78aa0ce485338ddb0121039eb1a57645dd67e5c90fcc3fc0b5b5c4e0590156cfee2baa3f8662c10e193b9bfeffffff02160772000000000017a9147804ac0238b575b11af4da3f8c4f192f37b534f587204bbc000000000017a9146676e4612fe637077c2b8123156cdd8e93e27e88870247304402200b4ce38eab830a4ef40dc8777054a09667734057b05b449952f49faa1a72e65602206dee994d47d5155981475cbb8a4bd296a6a4e78d4d5b9de913ae50ae8cb3ff01012103286e9acec10501ff78da4e4e5b956a060bab9a6898f193829a6b20ef98b1ebdb0000000000

Refund transaction (afe1dbee0fd4e58b27ab22a17dc27b73963dd933fe3b4d1366c10edc57e4ed72):
0200000001959344be6657bc9c668ad243d3eebb798c3ab87267bc7babd6ef954896d5bf1d01000000cf483045022100949782c8cbf8010739b50ec45bfa4ace7d6800e6e2e77bbab36ed0f17192184e0220782fd99f9c927a06a0d835ad79d06110391eb91d446af8e06217613755269919012102b3516f78fe5f4712c84fc2373cadd3e458e1d6c09fc5a74c8df47c8b024dc947004c616382012088a820ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa58876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba67041e990a5bb17576a9143a1e96e746e5d673e53f8cbca956f38e230f081e6888ac0000000001f549bc00000000001976a91455f6301a86a0728ca73cfad8fc088e85bd730c9f88ac1e990a5b

Publish contract transaction? [y/N] y
Published contract transaction (1dbfd5964895efd6ab7bbc6772b83a8c79bbeed343d28a669cbc5766be449395)
```
You can check the transaction [on a bitcoin testnet blockexplorer](https://testnet.blockexplorer.com/tx/1dbfd5964895efd6ab7bbc6772b83a8c79bbeed343d28a669cbc5766be449395) where you can see that 0.1234 BTC is sent to 2N2b1Lcs2ic3kt1bj8aoHUDhy9YawY47vhd (= the contract script hash) being a [p2sh](https://en.bitcoin.it/wiki/Pay_to_script_hash) address in the bitcoin testnet. 


 ### audit contract

Bob sends Alice the contract and the contract transaction. Alice should now verify if
- the script is correct 
- the locktime is far enough in the future
- the amount is correct
- she is the recipient 

command:`btcatomicswap auditcontract <contract> <contract transaction>`

 ```
./btcatomicswap auditcontract  6382012088a820ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa58876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba67041e990a5bb17576a9143a1e96e746e5d673e53f8cbca956f38e230f081e6888ac 02000000000102215d7d1855d16bf9d3c9c349ef8d0b7b58cb5c89922148bf4912d462ff4caef6010000001716001437c1c3347cb976780fe386f0631582fb751e21d8feffffff9400ee75b3ea0600ce429ecb9a0d5d6565e027d9d8dfd051285388169a676ed8000000006b483045022100dc32890b45892dbfa313fddc67e8c176228f82e10635ef56905bdd307f7a2bfa022069a2b8e96f550a005e4650a8ac0704e0436a902be97bd7ce78aa0ce485338ddb0121039eb1a57645dd67e5c90fcc3fc0b5b5c4e0590156cfee2baa3f8662c10e193b9bfeffffff02160772000000000017a9147804ac0238b575b11af4da3f8c4f192f37b534f587204bbc000000000017a9146676e4612fe637077c2b8123156cdd8e93e27e88870247304402200b4ce38eab830a4ef40dc8777054a09667734057b05b449952f49faa1a72e65602206dee994d47d5155981475cbb8a4bd296a6a4e78d4d5b9de913ae50ae8cb3ff01012103286e9acec10501ff78da4e4e5b956a060bab9a6898f193829a6b20ef98b1ebdb0000000000
Contract address:        3B2oGsw179YQgDyBTTBQrGihwCNmnpEVy7
Contract value:          0.1234 BTC
Recipient address:       1Et3zyPNrdKbeCMcxxNhKFgvu2zHE1GTZy
Author's refund address: 16JJqX8ny8q1wUvwKEU6FfjRk3yHjtdKsd

Secret hash: ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5

Locktime: 2018-05-27 11:40:14 +0000 UTC
Locktime reached in 47h54m59s
```

WARNING:
A check on the blockchain should be done as the auditcontract does not do that so an already spent output could have been used as an input. Checking if the contract has been mined in a block should suffice

### Participate

Alice trusts the contract so she participates in the atomic swap by paying the tokens into a threefold token  contract using the same secret hash. 

Bob creates a new threefold address ( or uses an existing one): 
```
tfchainc wallet address
Created new address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
```

Bob sends this address to Alice who uses it to participate in the swap.
command:`tfchainc atomicswap participate <initiator address> <amount> <secret hash>`
```
$./tfchainc atomicswap  participate 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb 567 ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5
Contract address: 02699891175f7a84f000deeb7d43d0c2c339b2de6ffe3561212f775d244b0f0ead2ce95cee274d
Contract value: 567 TFT
Recipient address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
Refund address: 016e3a2832cdbde6b714135ad1b38d150fef6053058a4f3361631bc1bc918573ddd7182cafd130

Hashed Secret: ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5

Locktime: 1527335287 (2018-05-26 13:48:07 +0200 CEST)
Locktime reached in: 23h59m59.621948718s

Publish atomic swap (participation) transaction? [Y/N] Y
published contract transaction
OutputID: b1a7fde7425416c34d395cd89d839f1dfef55d3a68b9cc0cec85803b110cb33f
``` 

The above command will create a transaction with `567` TFT as the Output  value of the output (`b1a7fde7425416c34d395cd89d839f1dfef55d3a68b9cc0cec85803b110cb33f`). The output can be claimed by Bobs address (`01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb`)  and Bob will  to also have to provide the secret that hashes to the hashed secret `ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5`.

Alice now informs Bob that the Threefold contract transaction has been created and provides him with the contract details.

### audit Threefold contract

Just as Alice had to audit Bob's contract, Bob now has to do the same with Alice's contract before withdrawing. 
Bob verifies if:
- the amount of threefold tokens () defined in the output is correct
- the attached script is correct
- the locktime, hashed secret (`ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5`) and wallet address, defined in the attached script, are correct

command:`tfchainc atomicswap audit outputid|unlockhash dest src timelock hashedsecret [amount] [flags]`
```
$./tfchainc atomicswap  audit b1a7fde7425416c34d395cd89d839f1dfef55d3a68b9cc0cec85803b110cb33f 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb 016e3a2832cdbde6b714135ad1b38d150fef6053058a4f3361631bc1bc918573ddd7182cafd130 1527335287  ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5 567
An unspend atomic swap contract could be found for the given outputID,
and the given contract information matches the found contract's information, all good! :)

Contract address: 02699891175f7a84f000deeb7d43d0c2c339b2de6ffe3561212f775d244b0f0ead2ce95cee274d
Recipient address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
Refund address: 016e3a2832cdbde6b714135ad1b38d150fef6053058a4f3361631bc1bc918573ddd7182cafd130

Hashed Secret: ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5

Locktime: 1527335287 (2018-05-26 13:48:07 +0200 CEST)
Locktime reached in: 23h47m29.225500748s
```

The audit also checks if that the given contract's output   has not already been spend.

### redeem tokens

Now that both Bob and Alice have paid into their respective contracts, Bob may withdraw from the Threefold contract. This step involves publishing a transaction which reveals the secret to Alice, allowing her to withdraw from the Bitcoin contract.

command:`tfchainc atomicswap claim outputid secret  amount`

```
$./tfchainc atomicswap claim  b1a7fde7425416c34d395cd89d839f1dfef55d3a68b9cc0cec85803b110cb33f  83639420c8683e288152f55b73598bc33f0ce6712316a8edc6e9772787845d2f 567
An unspend atomic swap contract could be found for the given outputID,
and the given contract information matches the found contract's information, all good! :)

Contract address: 02699891175f7a84f000deeb7d43d0c2c339b2de6ffe3561212f775d244b0f0ead2ce95cee274d
Recipient address: 01ba82c45bc004a7a4a169c7daade3422c59158981044e4f341e7cea57a2852a36ea43e9fc25bb
Refund address: 016e3a2832cdbde6b714135ad1b38d150fef6053058a4f3361631bc1bc918573ddd7182cafd130

Hashed Secret: ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5

Locktime: 1527335287 (2018-05-26 13:48:07 +0200 CEST)
Locktime reached in: 23h40m17.370839819s

Publish atomic swap claim transaction? [Y/N] Y

Published atomic swap claim transaction!
Transaction ID: e855c58903e9041105ebdd13c9beedc9a0943faee7fc4b999be96d6f22ac33ea
>   NOTE that this does NOT mean for 100% you'll have the money!
> Due to potential forks, double spending, and any other possible issues your
> claim might be declined by the network. Please check the network
> (e.g. using a public explorer node or your own full node) to ensure
> your payment went through. If not, try to audit the contract (again).
```
 This 
### redeem bitcoins

Now that Bob has withdrawn from the rivine contract and revealed the secret. If bob is really nice he could simply give the secret to Alice. However,even if he doesn't do this Alice can extract the secret from this redemption transaction. Alice may watch a block explorer to see when the rivine contract output was spent and look up the redeeming transaction.

Alice can automatically extract the secret from the input where it is used by Bob, by simply giving the outputID of the contract. Either you do this using a public web-based explorer, by looking up the outputID as hash. Or you let the command line client do it automatically for you:

command:` tfchainc atomicswap extractsecret outputid`
```
./tfchainc atomicswap --addr explorer.testnet.threefoldtoken.com extractsecret b1a7fde7425416c34d395cd89d839f1dfef55d3a68b9cc0cec85803b110cb33f
TODO:add output
```

NOTE: in this call I gave a public explorer address as I have no explorer node running myself.
Therefore I can use a public explorer to look it up for me instead.
Should you have a local explorer node running on the default address, you can simply omit the flag and use 
`$tfchainc extractsecret abcdef01234567890abcdef01234567890abcdef01234567890abcdef0123452` .

With the secret known (extracted from the coinInput with parent OutputID `abcdef01234567890abcdef01234567890abcdef01234567890abcdef0123452`), Alice may redeem from Bob's Bitcoin contract:
command: `btcatomicswap redeem <contract> <contract transaction> <secret>`
```
./btcatomicswap --testnet --rpcuser=user --rpcpass=pass redeem  6382012088a820ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa58876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba67041e990a5bb17576a9143a1e96e746e5d673e53f8cbca956f38e230f081e6888ac  02000000000102215d7d1855d16bf9d3c9c349ef8d0b7b58cb5c89922148bf4912d462ff4caef6010000001716001437c1c3347cb976780fe386f0631582fb751e21d8feffffff9400ee75b3ea0600ce429ecb9a0d5d6565e027d9d8dfd051285388169a676ed8000000006b483045022100dc32890b45892dbfa313fddc67e8c176228f82e10635ef56905bdd307f7a2bfa022069a2b8e96f550a005e4650a8ac0704e0436a902be97bd7ce78aa0ce485338ddb0121039eb1a57645dd67e5c90fcc3fc0b5b5c4e0590156cfee2baa3f8662c10e193b9bfeffffff02160772000000000017a9147804ac0238b575b11af4da3f8c4f192f37b534f587204bbc000000000017a9146676e4612fe637077c2b8123156cdd8e93e27e88870247304402200b4ce38eab830a4ef40dc8777054a09667734057b05b449952f49faa1a72e65602206dee994d47d5155981475cbb8a4bd296a6a4e78d4d5b9de913ae50ae8cb3ff01012103286e9acec10501ff78da4e4e5b956a060bab9a6898f193829a6b20ef98b1ebdb0000000000  83639420c8683e288152f55b73598bc33f0ce6712316a8edc6e9772787845d2f
Redeem fee: 0.00359616 BTC (0.01109926 BTC/kB)

Redeem transaction (35565cadbfe0a074073d579e3139bbab61cf42017b05e6dfa5241586a8f7b77c):
0200000001959344be6657bc9c668ad243d3eebb798c3ab87267bc7babd6ef954896d5bf1d01000000ef4730440220300e89b7fc9a0aa0c8e32ba31cca7a1bdcf6237adab21be72c99dacd82b7566b0220379f2c79db197fbf7b655028f6e46202c2f3ee4c956c17c3603a49c3d1faee110121023d3bdd190a65b033905a5a521596598518ac884b6b974f4395fe3e87525b253b2083639420c8683e288152f55b73598bc33f0ce6712316a8edc6e9772787845d2f514c616382012088a820ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa58876a91498415a65a8b96b72a3cc26b81e37af698d4ff4ba67041e990a5bb17576a9143a1e96e746e5d673e53f8cbca956f38e230f081e6888acffffffff0160ceb600000000001976a91467e8cb4c52cb34b50e3036d969ff00299d4ce84f88ac1e990a5b

Publish redeem transaction? [y/N] y
Published redeem transaction (35565cadbfe0a074073d579e3139bbab61cf42017b05e6dfa5241586a8f7b77c)
```
his transaction can be verified [on a bitcoin testnet blockexplorer](https://testnet.blockexplorer.com/tx/35565cadbfe0a074073d579e3139bbab61cf42017b05e6dfa5241586a8f7b77c) .
The cross-chain atomic swap is now completed and successful.

## References

Rivine atomic swaps are an implementation of [Decred atomic swaps](https://github.com/decred/atomicswap).

[Bitcoin scripts and opcodes](https://en.bitcoin.it/wiki/Script)
