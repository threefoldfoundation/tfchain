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

### 3Bot Transactions

The composition, encoding and signing of the three different 3Bot transactions are fully explained in the following subchapters.

Please note that you might want to read a high level technical overview, found at [3bot.md](3bot.md), prior to reading this chapter. Further you might also want to make sure that you're familiar with the tfchain binary encoding, as the 3Bot transactions are the first transaction versions where this encoding library is used. You can find more information about the tfchain binary encoding at [binary_encoding.md](binary_encoding.md).

#### 3Bot Registration Transaction

The 3Bot Registration Transaction is used to register a new 3Bot. It has to be new, meaning that the public key cannot be linked yet to an existing 3Bot. Within the consensus engine and explorers, 3Bots are represented by 3Bot records. One such record will be created as the result of this transaction. Other 3Bot transaction types will modify this created record, rather than creating a new one. As part of the registration process, the 3Bot is also assigned a unique 32-bit identifier, which is to be used to identify the 3Bot.

##### JSON Encoding a 3Bot Registration Transaction

```javascript
{
	// 0x90,
	// the version of a 3Bot Registration Transaction
	"version": 144,
	// the 3Bot Registration Transaction data
	"data": {
		// optional network addresses that can be used to reach the 3Bot on its Public API,
		// maximum 10 addresses are allowed, and no fees are paid for it
		// during registration.
		"addresses": ["91.198.174.192", "example.org"],
		// optional names that can be used to reach the 3Bot on its Public API, using these
		// aliases rather than directly using the 3Bot addresses, up to 5 names are
		// allowed, and the first one is free-of-charge during registration.
		// Note that even though both addresses and names are optional,
		// at least one name or one address is required, so one of twoproperty is optional,
		// not both.
		"names": ["chatbot.example"],
		// The number of months that are prepaid in advance,
		// at least one is required, 12 gives a 30% discount,
		// 24 is the limit and gives a 50% discount.
		"nrofmonths": 1,
		// Required transaction fee, has to be equal to or greater than 0.1 TFT
		"txfee": "1000000000",
		// Coin Inputs used to fund the Tx and 3Bot fees
		"coininputs": [{
			"parentid": "6baaa92439370a5110fdc244286a49b40b282b2af5af81e7a8e31c3658f16c04",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d",
					"signature": "b7da6d67e98c15ff83709419269a6b2f7b041c7f3605e927113d53fd1f6fafec4db2a991df7e7b3824d8fa806d2809d59d9d6560e5121b048910c40a5ed40f0d"
				}
			}
		}],
		// Optional (single) Refund Coin Output, can be used in case the coin input,
		// defines more input coins than required for the 3Bot and Tx fees.
		"refundcoinoutput": {
			"value": "99999798000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "0173f82c3ee74286c33fee8d883a7e9e759c6230b9e4e956ef233d7202bde69da45054270eef99"
				}
			}
		},
		// Identification of the 3Bot, containing its public key and
		// signature (of this Tx). The Public key is only given during registration,
		// afterwards it has to be looked up using its unique identifier.
		"identification": {
			"publickey": "ed25519:4e42a2fcfc0963d6fa7bb718fd088d9b6544331e8562d2743e730cdfbedeb55a",
			"signature": "a5ec12a859e56e8ddad951007591ad989dafc90d9aaabe8c879de42d4ff6edcd40213a002da251444b6fb6a29d78e4bc6bcfc844052969da7d0a67d91fa9c001"
		}
	}
}
```

###### Binary Encoding a 3Bot Registration Transaction

The binary encoding of a 3Bot Registration Transaction uses the tfchain encoding package. In order to understand the binary encoding of such a transaction, please see [the tfchain encoding documentation][tfchain-encoding] in order to understand how a 3Bot Registration Transaction is binary encoded.

The same transaction that was shown as an example of a JSON-encoded 3Bot Registration Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
90e112115bc6aec02c6578616d706c652e6f72671e63686174626f742e6578616d706c6504000000000000003b9aca00026baaa92439370a5110fdc244286a49b40b282b2af5af81e7a8e31c3658f16c04018000000000000000656432353531390000000000000000002000000000000000a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d4000000000000000b7da6d67e98c15ff83709419269a6b2f7b041c7f3605e927113d53fd1f6fafec4db2a991df7e7b3824d8fa806d2809d59d9d6560e5121b048910c40a5ed40f0d08000000000000000163454955669c000121000000000000000173f82c3ee74286c33fee8d883a7e9e759c6230b9e4e956ef233d7202bde69da4004e42a2fcfc0963d6fa7bb718fd088d9b6544331e8562d2743e730cdfbedeb55aa5ec12a859e56e8ddad951007591ad989dafc90d9aaabe8c879de42d4ff6edcd40213a002da251444b6fb6a29d78e4bc6bcfc844052969da7d0a67d91fa9c001
```

###### Signing a 3Bot Registration Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

> Note though that for the signing of 3Bot transactions the [tfchain encoding library][tfchain-encoding] is used.

In order to sign a 3Bot transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(TFChainBinaryEncoding(
  - transactionVersion: 1 byte, hardcoded to `0x90` (144 in decimal)
  - specifier: 16 bytes, hardcoded to "bot register tx\0"
  - TFChainBinaryEncoding(addresses, names, nrOfMonths)
  - length(coinInputs): int (8 bytes, little endian)
  - for each coin input:
    - parentID
  - TFChainBinaryEncoding(txFee, ptr(refundCoinOutput), publicKey)
)) : 32 bytes fixed-size crypto hash
```

#### 3Bot Record Update Transaction

The 3Bot Record Update Transaction is used to update an existing 3Bot,
meaning it has to be created in an existing block on the chain, using a 3Bot Registration Transaction.
The 3Bot can be inactive at the point of update, as long as it is made active again at that point by
paying for at least one month.

##### JSON Encoding a 3Bot Record Update Transaction

```javascript
{
	// 0x91,
	// the 3Bot Record Update Transaction Version
	"version": 145,
	// the 3Bot Record Update Transaction data
	"data": {
		// unique identifier of the 3Bot, assdigned during registration
		"id": 2,
		// optional, network addresses to add/remove:
		// - only network addresses currently existing in the 3Bot record can be removed;
		// - after applying the removed and added network addresses, the record
		//   is allowed a maximum of 10 addresses;
		// - 20 TFT is to be paid to modify the network addresses of an existing 3Bot.
		"addresses": {
			"add": ["example.com"],
			"remove": ["example.org"]
		},
		// optional, 3Bot names to add/remove:
		// - only 3Bot names currently linked to the 3Bot record can be removed;
		// - after applying the removed and added 3Bot names, the record
		//   is allowed a maximum of 5 names;
		// - 50 TFT is to be paid per added bot name,
		//   removing bot names requires no additional fee.
		"names": {
			"add": ["voicebot.example", "voicebot.example.myorg"],
			"remove": ["chatbot.example"]
		},
		// number of months to be added:
		//  - if the 3Bot was inactive this field is required, the 3Bot expiration date will
		//    be reset starting from this transaction's block time and adding the number of months to it;
		//  - if the 3Bot was (still) active this field is optional, and the number of months
		//    will be added to the current expiration time;
		"nrofmonths": 0,
		// Required transaction fee, has to be equal to or greater than 0.1 TFT
		"txfee": "1000000000",
		// Coin Inputs used to fund the Tx and 3Bot fees
		"coininputs": [{
			"parentid": "716f00dcaa5604f665aafad40d7704dba416de060174b0d3dc847bf61936f14f",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:7469d51063cdb690cc8025db7d28faadc71ff69f7c372779bf3a1e801a923e02",
					"signature": "33d02ebecc5e54de79d474473228511522219b922006ab68fc4732881dbd219938e402558e1eb4c981f27df17c28728c21c8bf3027341b7602421e715968bc00"
				}
			}
		}],
		// Optional (single) Refund Coin Output, can be used in case the coin input,
		// defines more input coins than required for the 3Bot and Tx fees.
		"refundcoinoutput": {
			"value": "99999677000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01af49ca1223d84089b60b40d6ae171dc951e331938fd75fd39e0167a989f3a83b0781f2372113"
				}
			}
		},
		// Signature of this 3Bot Record Update Transaction, signed by the 3Bot,
		// can be validated by looking up the 3Bot's public key (using its ID).
		"signature": "4239cbe196f188f051d01cbe490bb52c8667cda56b1882f4de5aed364b28fb4e375d1c237108d80f215ed96276d15ef38a992d1d1c0c019da951e680bfc79c05"
	}
}
```

###### Binary Encoding a 3Bot Record Update Transaction

The binary encoding of a 3Bot Record Update Transaction uses the tfchain encoding package. In order to understand the binary encoding of such a transaction, please see [the tfchain encoding documentation][tfchain-encoding] in order to understand how a 3Bot Record Update Transaction is binary encoded.

The same transaction that was shown as an example of a JSON-encoded 3Bot Record Update Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
9102000000e0112c6578616d706c652e636f6d2c6578616d706c652e6f72671220766f696365626f742e6578616d706c652c766f696365626f742e6578616d706c652e6d796f72671e63686174626f742e6578616d706c6504000000000000003b9aca0002716f00dcaa5604f665aafad40d7704dba416de060174b0d3dc847bf61936f14f0180000000000000006564323535313900000000000000000020000000000000007469d51063cdb690cc8025db7d28faadc71ff69f7c372779bf3a1e801a923e02400000000000000033d02ebecc5e54de79d474473228511522219b922006ab68fc4732881dbd219938e402558e1eb4c981f27df17c28728c21c8bf3027341b7602421e715968bc0008000000000000000163452d293d220001210000000000000001af49ca1223d84089b60b40d6ae171dc951e331938fd75fd39e0167a989f3a83b804239cbe196f188f051d01cbe490bb52c8667cda56b1882f4de5aed364b28fb4e375d1c237108d80f215ed96276d15ef38a992d1d1c0c019da951e680bfc79c05
```

###### Signing a 3Bot Record Update Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

> Note though that for the signing of 3Bot transactions the [tfchain encoding library][tfchain-encoding] is used.

In order to sign a 3Bot transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(TFChainBinaryEncoding(
  - transactionVersion: 1 byte, hardcoded to `0x90` (144 in decimal)
  - specifier: 16 bytes, hardcoded to "bot recupdate tx"
  - identifier of the 3Bot (uint32)
  - TFChainBinaryEncoding(addresses_add, addresses_remove, names_add, names_remove, nrOfMonths)
  - length(coinInputs): int (8 bytes, little endian)
  - for each coin input:
    - parentID
  - TFChainBinaryEncoding(txFee, ptr(refundCoinOutput))
)) : 32 bytes fixed-size crypto hash
```

#### 3Bot Name Transfer Transaction

The 3Bot Name Transfer Transaction is used to transfer one or multiple 3Bot names from
one existing active 3Bot to another existing active 3Bot.

##### JSON Encoding a 3Bot Name Transfer Transaction

```javascript
{
	// 0x92,
	// the version of a 3Bot Name Transfer Transaction
	"version": 146,
	// the Name Transfer Transaction Data
	"data": {
		// unique identifier and signature of the sending 3Bot,
		// meaning the 3Bot transferring names it owned to the receiver 3Bot.
		"sender": {
			"id": 2,
			"signature": "86755218edff668c46f5c7a0bb3788a35e6d0b7de317aa3cb2a02d29762133302bcec739c9a20c8e39ba2dc950f7e604ccfb801124f4954785a2444c9ba78109"
		},
		// unique identifier and signature of the receiver 3Bot,
		// meaning the 3Bot receiving names from the sender 3Bot.
		"receiver": {
			"id": 1,
			"signature": "15feaecd462b9223db1703f0e650146981a96b4972901bf3a8ca4480224ceac173418b7724155f0f91e0a30af28ef75afa829ab028c1c0b8360ba42bfffa9404"
		},
		// names to be transferred from the sender to the receiver 3Bot.
		"names": ["voicebot.example", "voicebot.example.myorg"],
		// Required transaction fee, has to be equal to or greater than 0.1 TFT
		"txfee": "1000000000",
		// Coin Inputs used to fund the Tx and 3Bot fees
		"coininputs": [{
			"parentid": "07d4c70711d634c6922c16b994168bac11b359e0c61c231209132ad4dfa8c1b2",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:300d034c02cfcc58ddf2b3059547ef91184f49f4a84bc3ec0123051bacfb987e",
					"signature": "5556ad839ebde45fe09b2bfabdff661b5e4841b7d3668def2bd7a6eca0a621519dcaf456368a96e0bb89de179d57c90745c7d5e4dac86a52766f49deb2e65208"
				}
			}
		}],
		// Optional (single) Refund Coin Output, can be used in case the coin input,
		// defines more input coins than required for the 3Bot and Tx fees.
		"refundcoinoutput": {
			"value": "99999576000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "011c17aaf2d54f63644f9ce91c06ff984182483d1b943e96b5e77cc36fdb887c846b60460bceb0"
				}
			}
		}
	}
}
```

###### Binary Encoding a 3Bot Name Transfer Transaction

The binary encoding of a 3Bot Name Transfer Transaction uses the tfchain encoding package. In order to understand the binary encoding of such a transaction, please see [the tfchain encoding documentation][tfchain-encoding] in order to understand how a 3Bot Name Transfer Transaction is binary encoded.

The same transaction that was shown as an example of a JSON-encoded 3Bot Name Transfer Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
92020000008086755218edff668c46f5c7a0bb3788a35e6d0b7de317aa3cb2a02d29762133302bcec739c9a20c8e39ba2dc950f7e604ccfb801124f4954785a2444c9ba78109010000008015feaecd462b9223db1703f0e650146981a96b4972901bf3a8ca4480224ceac173418b7724155f0f91e0a30af28ef75afa829ab028c1c0b8360ba42bfffa94041220766f696365626f742e6578616d706c652c766f696365626f742e6578616d706c652e6d796f726704000000000000003b9aca000207d4c70711d634c6922c16b994168bac11b359e0c61c231209132ad4dfa8c1b2018000000000000000656432353531390000000000000000002000000000000000300d034c02cfcc58ddf2b3059547ef91184f49f4a84bc3ec0123051bacfb987e40000000000000005556ad839ebde45fe09b2bfabdff661b5e4841b7d3668def2bd7a6eca0a621519dcaf456368a96e0bb89de179d57c90745c7d5e4dac86a52766f49deb2e65208080000000000000001634515a52b7000012100000000000000011c17aaf2d54f63644f9ce91c06ff984182483d1b943e96b5e77cc36fdb887c84
```

###### Signing a 3Bot Record Update Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

> Note though that for the signing of 3Bot transactions the [tfchain encoding library][tfchain-encoding] is used.

In order to sign a 3Bot transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(TFChainBinaryEncoding(
  - transactionVersion: 1 byte, hardcoded to `0x90` (144 in decimal)
  - specifier: 16 bytes, hardcoded to "bot nametrans tx"
  - identifier of the sender 3Bot (uint32)
  - identifier of the receiver 3Bot (uint32)
  - TFChainBinaryEncoding(names)
  - length(coinInputs): int (8 bytes, little endian)
  - for each coin input:
    - parentID
  - TFChainBinaryEncoding(txFee, ptr(refundCoinOutput))
)) : 32 bytes fixed-size crypto hash
```

[rivine]: https://github.com/rivine/rivine
[rivine-encoding]: https://github.com/rivine/rivine/blob/master/doc/Encoding.md
[tfchain-encoding]: binary_encoding.md
[rivine-txs]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md
[rivine-tx-v1]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-v1-transactions
[rivine-condition-uh]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-an-unlockhashcondition
[rivine-condition-tl]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-timelockcondition
[rivine-condition-multisig]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-multisignaturecondition
[rivine-double-spending]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#double-spend-rules
[rivine-sign-tx]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#signing-a-v1-transaction
[rivine-signing-into]: https://github.com/rivine/rivine/blob/master/doc/transactions/transaction.md#introduction-to-signing-transactions
