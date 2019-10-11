# Cross chain atomic swap walkthrough using the Electrum thin client.

## required tools  
In order to execute atomic swaps as described in this document, you need to have the [Electrum wallet](https://electrum.org/) and the atomic swap tools that are available at <https://github.com/threefoldtech/atomicswap/releases>. 

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
 

Start Electrum on testnet and create a default wallet but do not set a password on it: 
`./Electrum --testnet`
On osx the default location is `/Applications/Electrum.app/Contents/MacOS`)

Configure and start Electrum as a daemon :
```
./Electrum --testnet  setconfig rpcuser user
./Electrum --testnet  setconfig rpcpassword pass
./Electrum --testnet  setconfig rpcport 7777
./Electrum --testnet daemon
```
While the daemon is running, make it load the wallet in a different shell:
```
./Electrum --testnet daemon load_wallet
```
Start the tfchain daemon on testnet:
`tfchaind --network testnet`

Alice creates a new bitcoin  address and provides this to Bob: 
```￼
./Electrum --testnet getunusedaddress
mw7GjaHMy8D4rcK1ycYeFEsoCTJkSo54cz
```

### initiate step
Bob initiates the process by using btcatomicswap to pay 0.1234BTC into a Bitcoin contract using Alice's Bitcoin address, sending the contract transaction, and sharing the secret hash (not the secret), and contract's transaction with Alice. The refund transaction can not be sent until the locktime expires, but should be saved in case a refund is necessary.

command:`btcatomicswap initiate <participant address> <amount>`

```

$ ./btcatomicswap -testnet --rpcuser=user --rpcpass=pass -s  "localhost:7777" initiate mw7GjaHMy8D4rcK1ycYeFEsoCTJkSo54cz  0.1234 
Secret:      0720ee3138a7646e12f1c5b9adb5be59502f738cdbfed198cb676b4842865915
Secret hash: 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c

Contract fee: 0.00262528 BTC (0.01177256 BTC/kB)
Refund fee:   0.0034829 BTC (0.01196873 BTC/kB)

Contract (2N9rDE3K93eKm6gBfMFpDK3c6jVBJkdCQYW):
6382012088a82017b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c8876a914ab076c229f91119fb0b732fc601682875e8a124467043d9b7e5bb17576a914067c329d318fb219b5fa7aab6544355ada62d9d86888ac

Contract transaction (3c4da5200ec232140cb80f10eafd20ae0e0ff9acf461aa52f931b1f7117be858):
010000000121cd88b06bf680267cfed188ebd589f4bc19264cf20c919b3a845237bc3d998a010000006a47304402201b6f8ab2ac153f2eac87b9a4869a35acf66a300cedaa6346f4ebd72bccd89b2302201f2b7c9e3fccf2466362752b3257a4386ceeab353e93ad5b5476fb5aafbe1cf4012103da30f32bcb4e25f54dc8a507dc8c721176834eaf96bcc9e12afd2a2d17b18003fdffffff021f4bbc000000000017a914b61fef374878697ecd2c6628b34b2d9ccc22521d87c144e904000000001976a91418bb87066a98c37dfceb5ffa00fddbe2d1cfe34b88ac62211500

Refund transaction (d45c42783643cbf215d1f765414dfbb2d52e1ef32fc3bddbd5db2c17006cfaf5):
020000000158e87b11f7b131f952aa61f4acf90f0eae20fdea100fb80c1432c20e20a54d3c00000000ce47304402204562d490f49ba1e612b7ba038aedd1c41bbd39d3325795d39fc76697d17ab8a702202a0225920b3f9ff53e7ba272f12f5ccce30f9482aa86d6ae281084ec75c812110121026fb3b89c6c7c0fcac0cbb2f3be5c6e3cc999819f22a3e1235a2202188622d0b1004c616382012088a82017b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c8876a914ab076c229f91119fb0b732fc601682875e8a124467043d9b7e5bb17576a914067c329d318fb219b5fa7aab6544355ada62d9d86888ac00000000019dfab600000000001976a914067c329d318fb219b5fa7aab6544355ada62d9d888ac3d9b7e5b
```
You can check the transaction [on a bitcoin testnet blockexplorer](https://testnet.blockexplorer.com/tx/3c4da5200ec232140cb80f10eafd20ae0e0ff9acf461aa52f931b1f7117be858) where you can see that 0.1234 BTC is sent to 2N9rDE3K93eKm6gBfMFpDK3c6jVBJkdCQYW (= the contract script hash) being a [p2sh](https://en.bitcoin.it/wiki/Pay_to_script_hash) address in the bitcoin testnet. 


 ### audit contract

Bob sends Alice the contract and the contract transaction. Alice should now verify if
- the script is correct 
- the locktime is far enough in the future
- the amount is correct
- she is the recipient 

command:`btcatomicswap auditcontract <contract> <contract transaction>`

 ```
$ ./btcatomicswap --testnet --rpcuser=user --rpcpass=pass  -s  "localhost:7777" auditcontract 6382012088a82017b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c8876a914ab076c229f91119fb0b732fc601682875e8a124467043d9b7e5bb17576a914067c329d318fb219b5fa7aab6544355ada62d9d86888ac 010000000121cd88b06bf680267cfed188ebd589f4bc19264cf20c919b3a845237bc3d998a010000006a47304402201b6f8ab2ac153f2eac87b9a4869a35acf66a300cedaa6346f4ebd72bccd89b2302201f2b7c9e3fccf2466362752b3257a4386ceeab353e93ad5b5476fb5aafbe1cf4012103da30f32bcb4e25f54dc8a507dc8c721176834eaf96bcc9e12afd2a2d17b18003fdffffff021f4bbc000000000017a914b61fef374878697ecd2c6628b34b2d9ccc22521d87c144e904000000001976a91418bb87066a98c37dfceb5ffa00fddbe2d1cfe34b88ac62211500 
Contract address:        2N9rDE3K93eKm6gBfMFpDK3c6jVBJkdCQYW
Contract value:          0.12339999 BTC
Recipient address:       mw7GjaHMy8D4rcK1ycYeFEsoCTJkSo54cz
Author's refund address: mg7F8cXEkM4gGdkXZHvjogqrTWyTUV7gRe

Secret hash: 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c

Locktime: 2018-08-23 11:32:13 +0000 UTC
Locktime reached in 47h55m26s
```

WARNING:
A check on the blockchain should be done as the auditcontract does not do that so an already spent output could have been used as an input. Checking if the contract has been mined in a block should suffice

### Participate

Alice trusts the contract so she participates in the atomic swap by paying the tokens into a threefold token  contract using the same secret hash. 

Bob creates a new threefold address ( or uses an existing one): 
```
tfchainc wallet address
Created new address: 01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431
```

Bob sends this address to Alice who uses it to participate in the swap.
command:`tfchainc atomicswap participate <initiator address> <amount> <secret hash>`
```
$ tfchainc atomicswap participate 01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431 567 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c
Contract address: 02504acff12707d6df959f68d94aeb863f985320a6b604a8ccdb3a0fc5212b3e7661d942779726
Contract value: 567 TFT
Receiver's address: 01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431
Sender's (contract creator) address: 0145810e94de7b60d3082d0249789853f04082de7377e094a9f726f17121fb537f2307224b0fb4

SecretHash: 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c

TimeLock: 1534938493 (2018-08-22 13:48:13 +0200 CEST)
TimeLock reached in: 23h59m59.766589s
Publish atomic swap transaction? [Y/N] Y

published contract transaction

OutputID: f3ec88e3238f922882f5dc0b4d7c80256187d672a4effeb6563849c0245ee53b
TransactionID: b50f9ba4cc2b8a3242f04e518f16485cf3fd027cde8f54d00443be2d20b5ba85

Contract Info:

Contract address: 02504acff12707d6df959f68d94aeb863f985320a6b604a8ccdb3a0fc5212b3e7661d942779726
Contract value: 567 TFT
Receiver's address: 01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431
Sender's (contract creator) address: 0145810e94de7b60d3082d0249789853f04082de7377e094a9f726f17121fb537f2307224b0fb4

SecretHash: 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c

TimeLock: 1534938493 (2018-08-22 13:48:13 +0200 CEST)
TimeLock reached in: 23h59m49.914711s 
``` 

The above command will create a transaction with `567` TFT as the Output  value of the output (`f3ec88e3238f922882f5dc0b4d7c80256187d672a4effeb6563849c0245ee53b`). The output can be claimed by Bobs address (`01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431`)  and Bob will  to also have to provide the secret that hashes to the hashed secret `17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c`.

Alice now informs Bob that the Threefold contract transaction has been created and provides him with the contract details.

### audit Threefold contract

Just as Alice had to audit Bob's contract, Bob now has to do the same with Alice's contract before withdrawing. 
Bob verifies if:
- the amount of threefold tokens () defined in the output is correct
- the attached script is correct
- the locktime, hashed secret (`ed7c9cb48bf06db077641a09a0b7f7c3cc688760b771811fc0a0d07bdd3c6fa5`) and wallet address, defined in the attached script, are correct

command:`tfchainc atomicswap auditcontract outputid`
flags are available to automatically check the information in the contract.
```
$ ./tfchainc atomicswap auditcontract f3ec88e3238f922882f5dc0b4d7c80256187d672a4effeb6563849c0245ee53b
Atomic Swap Contract (condition) found:

Contract value: 567 TFT

Receiver's address: 01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431
Sender's (contract creator) address: 0145810e94de7b60d3082d0249789853f04082de7377e094a9f726f17121fb537f2307224b0fb4
Secret Hash: 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c
TimeLock: 1534938493 (2018-08-22 13:48:13 +0200 CEST)
TimeLock reached in: 23h57m6.531552s

found Atomic Swap Contract is valid
```

The audit also checks if that the given contract's output  has not already been spend.

### redeem tokens

Now that both Bob and Alice have paid into their respective contracts, Bob may withdraw from the Threefold contract. This step involves publishing a transaction which reveals the secret to Alice, allowing her to withdraw from the Bitcoin contract.

command:`tfchainc atomicswap redeem outputid secret`

```
$ tfchainc atomicswap redeem f3ec88e3238f922882f5dc0b4d7c80256187d672a4effeb6563849c0245ee53b 0720ee3138a7646e12f1c5b9adb5be59502f738cdbfed198cb676b4842865915
Contract address: 02504acff12707d6df959f68d94aeb863f985320a6b604a8ccdb3a0fc5212b3e7661d942779726
Contract value: 567 TFT
Receiver's address: 01f90b58b929ba434c256eea2b7895eb159ea72f76f8b7084f985163bc84bdf06bd853161fc431
Sender's (contract creator) address: 0145810e94de7b60d3082d0249789853f04082de7377e094a9f726f17121fb537f2307224b0fb4

SecretHash: 17b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c
Secret: 0720ee3138a7646e12f1c5b9adb5be59502f738cdbfed198cb676b4842865915

TimeLock: 1534938493 (2018-08-22 13:48:13 +0200 CEST)
TimeLock reached in: 23h52m8.057129s
Publish atomic swap redeem transaction? [Y/N] Y

published atomic swap redeem transaction
transaction ID: 21dd60aa8c0720840c701721a2dc76c1f9e68cd3a9d2b64eb5280f79e03396b1
>   Note that this does not mean for 100% you'll have the money.
> Due to potential forks, double spending, and any other possible issues your
> redeem might be declined by the network. Please check the network
> (e.g. using a public explorer node or your own full node) to ensure
> your payment went through. If not, try to audit the contract (again).
```

### redeem bitcoins

Now that Bob has withdrawn from the rivine contract and revealed the secret. If bob is really nice he could simply give the secret to Alice. However,even if he doesn't do this Alice can extract the secret from this redemption transaction. Alice may watch a block explorer to see when the rivine contract output was spent and look up the redeeming transaction.

Alice can automatically extract the secret from the input where it is used by Bob, by simply giving the outputID of the contract. Either you do this using a public web-based explorer, by looking up the outputID as hash. Or you let the command line client do it automatically for you:

command:`tfchainc atomicswap extractsecret transactionid`
```
$/tfchainc atomicswap extractsecret 21dd60aa8c0720840c701721a2dc76c1f9e68cd3a9d2b64eb5280f79e03396b1
atomic swap contract was redeemed
extracted secret: 0720ee3138a7646e12f1c5b9adb5be59502f738cdbfed198cb676b4842865915
```

With the secret known, Alice may redeem from Bob's Bitcoin contract:
command: `btcatomicswap redeem <contract> <contract transaction> <secret>`
```
./btcatomicswap --testnet --rpcuser=user --rpcpass=pass -s  "localhost:7777" redeem 6382012088a82017b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c8876a914ab076c229f91119fb0b732fc601682875e8a124467043d9b7e5bb17576a914067c329d318fb219b5fa7aab6544355ada62d9d86888ac 010000000121cd88b06bf68026

7cfed188ebd589f4bc19264cf20c919b3a845237bc3d998a010000006a47304402201b6f8ab2ac153f2eac87b9a4869a35acf66a300cedaa6346f4ebd72bccd89b2302201f2b7c9e3fccf2466362752b3257a4386ceeab353e93ad5b5476fb5aafbe1cf4012103da30f32bcb4e25f54dc8a507dc8c721176834eaf96bcc9e12afd2a2d17b18003fdffffff021f4bbc000000000017a914b61fef374878697ecd2c6628b34b2d9ccc22521d87c144e904000000001976a91418bb87066a98c37dfceb5ffa00fddbe2d1cfe34b88ac62211500 0720ee3138a7646e12f1c5b9adb5be59502f738cdbfed198cb676b4842865915
Redeem fee: 0.00495 BTC (0.01527778 BTC/kB)

Redeem transaction (f46d0198cdeeedbf5ad81c0a79f7b8296fde281deff1a70d418cd012eb853763):
020000000158e87b11f7b131f952aa61f4acf90f0eae20fdea100fb80c1432c20e20a54d3c00000000ef47304402200771ee5e74bc216869681f374d35e58654be302485bd4f5b42f999d789e60f6b02201586305db69fcc44a8f75817adf212f37af8ee3a4750bdb13e7c1be528b019be012103cf41fe2eca9e9c210574e0d92774f6d9da86fbf344cdab3712752788bb948386200720ee3138a7646e12f1c5b9adb5be59502f738cdbfed198cb676b4842865915514c616382012088a82017b723d2a2f23480fc464194038419e178306525c3279881843c4e2618293d6c8876a914ab076c229f91119fb0b732fc601682875e8a124467043d9b7e5bb17576a914067c329d318fb219b5fa7aab6544355ada62d9d86888acffffffff0187bdb400000000001976a914067c329d318fb219b5fa7aab6544355ada62d9d888ac3d9b7e5b

Publish redeem transaction? [y/N] Y
Published redeem transaction (f46d0198cdeeedbf5ad81c0a79f7b8296fde281deff1a70d418cd012eb853763)
```
This transaction can be verified [on a bitcoin testnet blockexplorer](https://testnet.blockexplorer.com/tx/8c77722af341c56f968175c09c57f70841e42d2651b456dfa88a04d63046df76) .
The cross-chain atomic swap is now completed and successful.

## References

- [Electrum](https://electrum.org)

