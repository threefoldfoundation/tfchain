# Transactions

The main purpose of a transaction is to transfer coins and/or block stakes between addresses.
For each coin that is spend (registered as coin output), including a transaction fee,
there must be one or multiple (coin) inputs backing it up. Meaning the sum of coin inputs,
must equal the sum of miner fees plus coin outputs. Coin inputs must be previously-registered outputs,
which haven't been used as input yet. The same rule applies to block stakes, another kind of asset,
except that miner fees are to be still paid in coins not in block stakes.

> [Minter Definition Transactions](#minter-definition-transactions) and [Coin Creation Transactions](#coin-creation-transactions) are exceptions to the rule, a transaction of this type define coin outputs and/or minder fees with no coin inputs defined. Meaning these coins are created with no backing of previous outputs.

As each output is backed by one or multiple inputs, it is not uncommon to have a too big
amount of input registered. If so, it is is required to register the extra
amount as another output, be at addressed to your own wallet.
This works very similar to the change money you get in a supermarket by paying too much.

## Types of Transactions

Each transaction has a version, which is to be decoded as the very first step,
and identifies the type of transaction. Knowing the version (= type),
it can be deduced how to decode the rest of the data, if possible at all.

All regular transactions are in general `0x01` transactions, a format defined by [Rivine][rivine].
This version deprecates `0x00` transactions, which are still valid but not used any longer
for the creation of new (regular) transactions by the official tfchain reference implementation found in this repo.

> In this doc we call any transaction which only means it is to send TFT and or block stakes between addresses a regular transaction.

In this document we will only describe the Transaction Types specific to tfchain.

Please read [the Rivine Transaction docs][rivine-txs] for more information about [0x00](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-v0-transactions) and [0x01](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-v1-transactions) transactions. These are regular transaction versions which are supported and used within tfchain as well, but are implemented in [Rivine][rivine]. In the Rivine documentation you'll also be able to find more information about [the double spending rule][rivine-double-spending] and [signing of transactions][rivine-sign-tx].

### Minter Definition Transactions

Minter Definition Transactions are used to redefine the creators of coins (AKA minters). These transactions can only be created by the Coin Creators. The (previously-defined) mint condition —meaning the mint condition active at the height of the (to be) created Minter Definition Transaction— defines who the coin creators are and thus who can redefine who the coin creators are to become. A mint condition can be any of the following conditions:

* (1) An [UnlockHash Condition][rivine-condition-uh]: it is a single person (or multiple people owning the same private key for the same wallet);
* (2) A [MultiSignature Condition][rivine-condition-multisig]: it is a multi signature wallet, most likely meaning multiple people owning different private keys that all have a certain degree of control in the same multi signature wallet, allowing them to come to a consensus on the creation of coins and redefinition of who the coin creators are (to be);
* (3) An [TimeLocked Condition][rivine-condition-tl]: the minting powers (creation of coins and redefinition of the coin creators) are locked until a certain (unix epoch) timestamp or block height. Once this height or time is reached, the internal condition (an [UnlockHash Condition][rivine-condition-uh] (1) or a [MultiSignature Condition][rivine-condition-multisig]) (2) is the condition that defines who can create coins and redefine who the coin creators are to be. Prior to this timestamp or block height no-one can create coins or redefine who the coin creators are to be, not even the ones defined by the internal condition of the currently active [TimeLocked Condition][rivine-condition-tl].

The Coin Creation transactions defines 4 fields:

* `mintfulfillment`: the fulfillment which has to fulfill the consensus-defined MintCondition, just the same as that a Coin Input's fulfillment has to fulfill the condition of the Coin Output it is about to spend;
* `mintcondition`: the condition which will become the new mint condition (that has to be fulfilled in order to create coins and redefine the mint condition, in other words the condition that defines who the coin creators are) once the transaction is part of a created block and until there is a newer block with an accepted mint condition;
* `minerfees`: defines the transaction fee(s) (works the same as in regular transactions);
* `arbitrarydata`: describes the capacity that is created/added, creating these coins as a result;
* `nonce`: a crypto-random 8-byte array, used to ensure the uniqueness of this transaction's ID;

> The Genesis Mint Condition is hardcoded and the specific condition for each network
> (`devnet`, `testnet` and `standard`) can be found in the source code of tfchain.
> See for more information the Godoc (and linked source code) for each network:
> * `standard`: https://godoc.org/github.com/threefoldfoundation/tfchain/pkg/config#GetStandardnetGenesisMintCondition
> * `testnet`: https://godoc.org/github.com/threefoldfoundation/tfchain/pkg/config#GetTestnetGenesisMintCondition
> * `devnet`: https://godoc.org/github.com/threefoldfoundation/tfchain/pkg/config#GetDevnetGenesisMintCondition

In practice a MultiSignature Condition will always be used as MintCondition,
this is however not a consensus-defined requirement, as discussed earlier.

#### JSON Encoding a Minter Definition Transaction

```javascript
{
	// 0x80, the version number of a Minter Definition Transaction
	"version": 128,
	// Coin Creation Transaction Data
	"data": {
		// crypto-random 8-byte array (base64-encoded to a string) to ensure
		// the uniqueness of this transaction's ID
		"nonce": "FoAiO8vN2eU=",
		// fulfillment which fulfills the MintCondition,
		// can be any type of fulfillment as long as it is
		// valid AND fulfills the MintCondition
		"mintfulfillment": {
			"type": 1,
			"data": {
				"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
				"signature": "bdf023fbe7e0efec584d254b111655e1c2f81b9488943c3a712b91d9ad3a140cb0949a8868c5f72e08ccded337b79479114bdb4ed05f94dfddb359e1a6124602"
			}
		},
		// condition which will become the new MintCondition
		// once the transaction is part of a created block and
		// until there is a newer block with another accepted MintCondition
		"mintcondition": {
			"type": 1,
			"data": {
				"unlockhash": "01e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73e047fe6a0703"
			}
		},
		// the transaction fees to be paid, also paid in
		// newly created) coins, rather than inputs
		"minerfees": ["1000000000"],
		// arbitrary data, can contain anything as long as it
		// fits within 83 bytes, but is in practice used
		// to link the capacity added/created
		// with as a consequence the creation of
		// these transaction and its coin outputs
		"arbitrarydata": "dGVzdC4uLiAxLCAyLi4uIDM="
	}
}
```

See the [Rivine documentation about JSON-encoding of v1 Transactions][rivine-tx-v1] for more information about the primitive data types used, as well as the meaning and encoding of the different types of unlock conditions and fulfillments.

#### Binary Encoding a Minter Definition Transaction

The binary encoding of a Minter Definition Transaction uses the Rivine encoding package. In order to understand the binary encoding of such a transaction, please see [the Rivine encoding documentation][rivine-encoding] and [Rivine binary encoding of a v1 transaction](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#binary-encoding-of-v1-transactions) in order to understand how a Minter Definition Transaction is binary encoded. That documentation also contains [an in-detail documented example of a binary encoding v1 transaction](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#example-of-a-binary-encoded-v1-transaction).

The same transaction that was shown as an example of a JSON-encoded Minter Definition Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
801680223bcbcdd9e5018000000000000000656432353531390000000000000000002000000000000000d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d77804000000000000000bdf023fbe7e0efec584d254b111655e1c2f81b9488943c3a712b91d9ad3a140cb0949a8868c5f72e08ccded337b79479114bdb4ed05f94dfddb359e1a612460201210000000000000001e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73010000000000000004000000000000003b9aca00180000000000000061206d696e74657220646566696e6974696f6e2074657374
```

#### Signing a Minter Definition Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

In order to sign a v1 transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(BinaryEncoding(
  - transactionVersion: 1 byte, hardcoded to `0x81` (129 in decimal)
  - specifier: 16 bytes, hardcoded to "minter defin tx\0"
  - nonce: 8 bytes
  - binaryEncoding(mintCondition)
  - length(minerFees): int64 (8 bytes, little endian)
  for each minerFee:
    - fee: Currency (8 bytes length + n bytes, little endian encoded)
  - arbitraryData: 8 bytes length + n bytes
)) : 32 bytes fixed-size crypto hash
```

### Coin Creation Transactions

Coin Creation Transactions are used for the creation of new coins. These transactions can only be created by the Coin Creators (also called minters). The Mint Condition defines who the coin creators are. If it is an [UnlockHash Condition][rivine-condition-uh] it is a single person, while it will be a [MultiSignature Condition][rivine-condition-multisig] in case there are multiple coin creators that have to come to a consensus.

The Coin Creation transactions defines 4 fields:

* `mintfulfillment`: the fulfillment which has to fulfill the consensus-defined MintCondition, just the same as that a Coin Input's fulfillment has to fulfill the condition of the Coin Output it is about to spend;
* `coinoutputs`: defines coin outputs, the destination of coins (works the same as in regular transactions);
* `minerfees`: defines the transaction fee(s) (works the same as in regular transactions);
* `arbitrarydata`: describes the capacity that is created/added, creating these coins as a result;
* `nonce`: a crypto-random 8-byte array, used to ensure the uniqueness of this transaction's ID;

> The Genesis Mint Condition is hardcoded and the specific condition for each network
> (`devnet`, `testnet` and `standard`) can be found in the source code of tfchain.
> See for more information the Godoc (and linked source code) for each network:
> * `standard`: https://godoc.org/github.com/threefoldfoundation/tfchain/pkg/config#GetStandardnetGenesisMintCondition
> * `testnet`: https://godoc.org/github.com/threefoldfoundation/tfchain/pkg/config#GetTestnetGenesisMintCondition
> * `devnet`: https://godoc.org/github.com/threefoldfoundation/tfchain/pkg/config#GetDevnetGenesisMintCondition

In practice a MultiSignature Condition will always be used as MintCondition,
this is however not a consensus-defined requirement. You can ready more about this in the chapter on
[Minter Definition Transactions](#minter-definition-transactions).

#### JSON Encoding a Coin Creation Transaction

```javascript
{
	// 0x81, the version number of a Coin Creation Transaction
	"version": 129,
	// Coin Creation Transaction Data
	"data": {
		// crypto-random 8-byte array (base64-encoded to a string) to ensure
		// the uniqueness of this transaction's ID
		"nonce": "1oQFzIwsLs8=",
		// fulfillment which fulfills the MintCondition,
		// can be any type of fulfillment as long as it is
		// valid AND fulfills the MintCondition
		"mintfulfillment": {
			"type": 1,
			"data": {
				"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
				"signature": "ad59389329ed01c5ee14ce25ae38634c2b3ef694a2bdfa714f73b175f979ba6613025f9123d68c0f11e8f0a7114833c0aab4c8596d4c31671ec8a73923f02305"
			}
		},
		// defines the recipients (as conditions) who are to
		// receive the paired (newly created) coin values
		"coinoutputs": [{
			"value": "500000000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01e3cbc41bd3cdfec9e01a6be46a35099ba0e1e1b793904fce6aa5a444496c6d815f5e3e981ccf"
				}
			}
		}],
		// the transaction fees to be paid, also paid in
		// newly created) coins, rather than inputs
		"minerfees": ["1000000000"],
		// arbitrary data, can contain anything as long as it
		// fits within 83 bytes, but is in practice used
		// to link the capacity added/created
		// with as a consequence the creation of
		// these transaction and its coin outputs
		"arbitrarydata": "dGVzdC4uLiAxLCAyLi4uIDM="
	}
}
```

See the [Rivine documentation about JSON-encoding of v1 Transactions][rivine-tx-v1] for more information about the primitive data types used, as well as the meaning and encoding of the different types of unlock conditions and fulfillments.

#### Binary Encoding a Coin Creation Transaction

The binary encoding of a Coin Creation Transaction uses the Rivine encoding package. In order to understand the binary encoding of such a transaction, please see [the Rivine encoding documentation][rivine-encoding] and [Rivine binary encoding of a v1 transaction](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#binary-encoding-of-v1-transactions) in order to understand how a Coin Creation Transaction is binary encoded. That documentation also contains [an in-detail documented example of a binary encoding v1 transaction](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#example-of-a-binary-encoded-v1-transaction).

The same transaction that was shown as an example of a JSON-encoded Coin Creation Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
8133a6432220334946018000000000000000656432353531390000000000000000002000000000000000d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d77804000000000000000a074b976556d6ea2e4ae8d51fbbb5ec99099f11918201abfa31cf80d415c8d5bdfda5a32d9cc167067b6b798e80c6c1a45f6fd9e0f01ac09053e767b15d310050100000000000000070000000000000001c6bf5263400001210000000000000001e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73010000000000000004000000000000003b9aca0012000000000000006d6f6e65792066726f6d2074686520736b79
```

#### Signing a Coin Creation Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

In order to sign a v1 transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(BinaryEncoding(
  - transactionVersion: 1 byte, hardcoded to `0x81` (129 in decimal)
  - specifier: 16 bytes, hardcoded to "coin mint tx\0\0\0\0"
  - nonce: 8 bytes
  - length(coinOutputs): int64 (8 bytes, little endian)
  for each coinOutput:
    - value: Currency (8 bytes length + n bytes, little endian encoded)
    - binaryEncoding(condition)
  - length(minerFees): int64 (8 bytes, little endian)
  for each minerFee:
    - fee: Currency (8 bytes length + n bytes, little endian encoded)
  - arbitraryData: 8 bytes length + n bytes
)) : 32 bytes fixed-size crypto hash
```

[rivine]: https://github.com/rivine/rivine
[rivine-encoding]: https://github.com/rivine/rivine/blob/master/doc/Encoding.md
[rivine-txs]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md
[rivine-tx-v1]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-v1-transactions
[rivine-condition-uh]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-an-unlockhashcondition
[rivine-condition-tl]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-timelockcondition
[rivine-condition-multisig]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-multisignaturecondition
[rivine-double-spending]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#double-spend-rules
[rivine-sign-tx]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#signing-a-v1-transaction
[rivine-signing-into]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#introduction-to-signing-transactions
