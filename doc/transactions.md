# Transactions

The main purpose of a transaction is to transfer coins and/or block stakes between addresses.
For each coin/ that is spend (registered as coin output), including a transaction fee,
there must be one or multiple (coin) inputs backing it up. Meaning the sum of coin inputs,
must equal the sum of miner fees plus coin outputs. Coin inputs must be previously-registered outputs,
which haven't been used as input yet. The same rule applies to block stakes, another kind of asset,
except that miner fees are to be paid in coins not in block stakes.

> [Coin Creation Transactions](#coin-creation-transactions) are exceptions to the rule, the transaction of this type define coin outputs and miner fees with no coin inputs defined. Hence why these transactions are called [Coin Creation transactions](#coin-creation-transactions). More about this type of transactions later.

As each output is backed by one or multiple inputs, it is not uncommon to have a too big
amount of input registered. If so, it is the convention to simply register the extra
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

### Coin Creation Transactions

Coin Creation Transactions are used for the creation of new coins. These transactions can only be created by the Coin Creators (also called minters). The Mint Condition defines who the coin creators are. If it is an [UnlockHash Condition][rivine-condition-uh] it is a single person, while it will be a [MultiSignature Condition][rivine-condition-multisig] in case there are multiple coin creators that have to come to a consensus.

The Coin Creation transactions defines 4 fields:

* `mintfulfillment`: the fulfillment which has to fulfill the consensus-defined MintCondition, just the same as that a Coin Input's fulfillment has to fulfill the condition of the Coin Output it is about to spend;
* `coinoutputs`: defines coin outputs, the destination of coins (works the same as in regular transactions);
* `minerfees`: defines the transaction fee(s) (works the same as in regular transactions);
* `arbitrarydata`: describes the capacity that is created/added, creating these coins as a result;
* `nonce`: a crypto-random 8-byte array, used to ensure the uniqueness of this transaction's ID;

> The Mint Condition is hardcoded and the specific condition for each network
> (`devnet`, `testnet` and `standard`) be found in the source code of tfchain.

In practice a MultiSignature Condition will always be used as MintCondition, but this is not a consensus-defined requirement.

#### JSON Encoding

```javascript
{
	// 0x81, the version number of a Coin Creation Transaction
	"version": 129,
	// Coin Creation Transaction Data
	"data": {
		// crypto-random 8-byte array to ensure
		// the uniqueness of this transaction's ID
		"nonce": [51, 166, 67, 34, 32, 51, 73, 70],
		// fulfillment which fulfills the MintCondition,
		// can be any type of fulfillment as long as it is
		// valid AND fulfills the MintCondition
		"mintfulfillment": {
			"type": 1,
			"data": {
				"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
				"signature": "a074b976556d6ea2e4ae8d51fbbb5ec99099f11918201abfa31cf80d415c8d5bdfda5a32d9cc167067b6b798e80c6c1a45f6fd9e0f01ac09053e767b15d31005"
			}
		},
		// defines the recipients (as conditions) who are to
		// receive the paired (newly created) coin values
		"coinoutputs": [{
			"value": "500000000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73e047fe6a0703"
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
		"arbitrarydata": "bW9uZXkgZnJvbSB0aGUgc2t5"
	}
}
```

See the [Rivine documentation about JSON-encoding of v1 Transactions][rivine-tx-v1] for more information about the primitive data types used, as well as the meaning and encoding of the different types of unlock conditions and fulfillments.

#### Binary Encoding

The binary encoding of a Coin Creation Transaction uses the Rivine encoding package. In order to understand the binary encoding of such a transaction, please see [the Rivine encoding documentation][rivine-encoding] and [Rivine binary encoding of a v1 transaction](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#binary-encoding-of-v1-transactions) in order to understand how a Coin Creation Transaction is binary encoded. That documentation also contains [an in-detail documented example of a binary encoding v1 transaction](https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#example-of-a-binary-encoded-v1-transaction).

The same transaction that was shown as an example of a JSON-encoded Coin Creation Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
8133a6432220334946018000000000000000656432353531390000000000000000002000000000000000d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d77804000000000000000a074b976556d6ea2e4ae8d51fbbb5ec99099f11918201abfa31cf80d415c8d5bdfda5a32d9cc167067b6b798e80c6c1a45f6fd9e0f01ac09053e767b15d310050100000000000000070000000000000001c6bf5263400001210000000000000001e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73010000000000000004000000000000003b9aca0012000000000000006d6f6e65792066726f6d2074686520736b79
```

[rivine]: https://github.com/rivine/rivine
[rivine-encoding]: https://github.com/rivine/rivine/blob/master/doc/Encoding.md
[rivine-txs]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md
[rivine-tx-v1]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-v1-transactions
[rivine-condition-uh]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-an-unlockhashcondition
[rivine-condition-multisig]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-multisignaturecondition
[rivine-double-spending]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#double-spend-rules
[rivine-sign-tx]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#signing-a-v1-transaction
