rivine-extension-erc20
======

Go ERC20 consensus plugin extension for Rivine.

See <https://github.com/threefoldtech/rivine/blob/master/extensions/readme.md> for more information about Rivine (Go) extensions.

See <https://github.com/threefoldfoundation/tfchain/blob/master/doc/erc20.md#motivation> for a possible motivation for your chain's need of this extension.
This document also contains an implementation concept implemented on a higher level. While it might be at times specific to the needs of the Threefold Chain,
it might help you decide on the usefulness of this extension for your chain.

Product owners
--------------

* Rob Van Mieghem ([@robvanmieghem](https://github.com/robvanmieghem))
* Glen De Cauwsemaecker ([@glendc](https://github.com/glendc))


ERC20 Transaction Info
----------------------

The composition, encoding and signing of the three different ERC20 transactions are fully explained in the following subchapters.

Please note that you might want to make sure that you're familiar with the Rivine binary encoding, used to encode ERC20 transactions.
You can find more information about the Rivine binary encoding at <https://github.com/threefoldtech/rivine/blob/master/doc/encoding/RivineEncoding.md>.

#### ERC20 Convert Transaction

The ERC20 Convert Transaction is used to convert TFT to ERC20-funds. Converting meaning that the used TFT will
be burned and their value will be exchanged into the matching value of ERC20 funds, paid into the account as
defined by this transaction as well.

##### JSON Encoding an ERC20 Convert Transaction

```javascript
{
	// 0xD0, an example version number of an ERC20 Convert Transaction
	"version": 208, // the decimal representation of the above example version number
	"data": {
		// Required ERC20-valid address, fixed length of 20 bytes
		"address": "0x0123456789012345678901234567890123456789",
		// Required value of TFT to be burned towards funding the ERC20 funds,
		// note that at least 1000 TFT is required, but more can be burned as well,
		// the more TFT the more ERC20 funds you'll get, with the exact value as defined by market
		// at the time of the transaction.
		//
		// Note that the ERC20 bridge wil take a small cut from the money in order to pay the
		// Gas Limit on the ERC20 side of things.
		"value": "1000000000000",
		// Required Transaction Fee
		"txfee": "1000000000",
		// Coin Inputs that fund the burned value as well as the required Transaction Fee.
		"coininputs": [{
			"parentid": "9c61ec964105ec48bc95ffc0ac820ada600a2914a8dd4ef511ed7f218a3bf469",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:7469d51063cdb690cc8025db7d28faadc71ff69f7c372779bf3a1e801a923e02",
					"signature": "a0c683e8728710b4d3cd7eed4e1bd38a4be8145a2cf91b875986870aa98c6265d76cbb637d78500010e3ab1b651e31ab26b05de79938d7d0aee01f8566d08b09"
				}
			}
		}],
		// Optional Coin Output, to be used in case the sum of the coin inputs is
		// higher than the burned value and transaction fee combined.
		"refundcoinoutput": {
			"value": "99999476000000000",
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

###### Binary Encoding an ERC20 Convert Transaction

The binary encoding of an ERC20 Convert Transaction uses the Rivine encoding package. In order to understand the binary encoding of such a transaction, please see [the Rivine encoding documentation][rivine-encoding] in order to understand how an ERC20 Convert Transaction is binary encoded.

The same transaction that was shown as an example of a JSON-encoded ERC20 Convert Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
d001234567890123456789012345678901234567890a2e90edd000083b9aca00029c61ec964105ec48bc95ffc0ac820ada600a2914a8dd4ef511ed7f218a3bf46901c4017469d51063cdb690cc8025db7d28faadc71ff69f7c372779bf3a1e801a923e0280a0c683e8728710b4d3cd7eed4e1bd38a4be8145a2cf91b875986870aa98c6265d76cbb637d78500010e3ab1b651e31ab26b05de79938d7d0aee01f8566d08b090110016344fe5cb488000142011c17aaf2d54f63644f9ce91c06ff984182483d1b943e96b5e77cc36fdb887c84
```

###### Signing an ERC20 Convert Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

> Note though that for the signing of ERC20 Transactions the [Rivine encoding library][rivine-encoding] is used.

In order to sign an ERC20 Convert transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(RivineBinaryEncoding(
  - transactionVersion: 1 byte
  - specifier: 16 bytes, hardcoded to "erc20 convert tx"
  - address: 20 bytes,
  - value: ? bytes,
  - all extra objects (not the length)
  - length(coinInputs): int (8 bytes, little endian)
  - for each coin input:
    - parentID
  - txFee
  - ptr(refundCoinOutput))
)) : 32 bytes fixed-size crypto hash
```

#### ERC20 Coin Creation Transaction

The ERC20 Coin Creation Transaction is used to convert ERC20-funds to TFT. Converting meaning that the used
ERC20-funds will be burned and their value will be exchanged into the matching value of TFT, paid into
the account as defined by the ERC20 Withdrawal address.

##### JSON Encoding an ERC20 Coin Creation Transaction

```javascript
{
	// 0xD1, an example version number of an ERC20 Coin Creation Transaction
	"version": 2019, // the decimal representation of the above example version number
	"data": {
		// TFT Address to be paid into
		"address": "01f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609adbbcd3977363c",
		// Value, funded by burning ERC20-funds, to be paid into the TFT Wallet identified by the attached TFT address
		"value": "100000000000",
		// Regular Transaction Fee
		"txfee": "1000000000",
		// ERC20 BlockID of the parent block of the paired ERC20 Transaction.
		"blockid": "0x0000000000000000000000000000000000000000000000000000000000000000"
		// ERC20 TransationID in which the matching ERC20-funds got burned,
		// each transactionID can only be used once to fund a TFT coin exchange.
		"txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
	}
}
```

###### Binary Encoding an ERC20 Coin Creation Transaction

The binary encoding of an ERC20 Coin Creation Transaction uses the Rivine encoding package.
In order to understand the binary encoding of such a transaction, please see [the Rivine encoding documentation][rivine-encoding]
in order to understand how an ERC20 Coin Creation Transaction is binary encoded.

The same transaction that was shown as an example of a JSON-encoded ERC20 Coin Creation Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
d101f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609ad0a174876e800083b9aca000000000000000000000000000000000000000000000000000000000000000000
```

###### Signing an ERC20 Coin Creation Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

> Note though that for the signing of ERC20 Transactions the [Rivine encoding library][rivine-encoding] is used.

In order to sign an ERC20 Coin Creation transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(RivineBinaryEncoding(
  - transactionVersion: 1 byte
  - specifier: 16 bytes, hardcoded to "erc20 coingen tx"
  - all extra objects (not the length)
  - address: binary encoded unlock hash
  - value
  - txFee
  - ERC20 BlockID: 32 bytes
  - ERC20 TransactionID: 32 bytes
)) : 32 bytes fixed-size crypto hash
```

#### ERC20 Address Registration Transaction

The ERC20 Address Registration Transaction is used to register an ERC20 Address as the withdrawal address,
linked to the TFT address generated with the attached public key.

##### JSON Encoding an ERC20 Address Registration Transaction

```javascript
{
	// 0xD2, an example version number of an ERC20 Address Registration Transaction
	"version": 210, // the decimal representation of the above example version number
	"data": {
		// public key from which the TFT address is generated, and as a consequence also the ERC20 Address
		"pubkey": "ed25519:a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d",
		// the TFT address (optionally attached in the JSON format only) generated from the attached public key
		"tftaddress": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f",
		// the ERC20 address (optionally attached in the JSON format only) generated from the attached public key
		"erc20address": "0x828de486adc50aa52dab52a2ec284bcac75be211",
		// signature to proof the ownership of the attached public key
		"signature": "fe13823a96928a573f20a63f3b8d3cde08c506fa535d458120fdaa5f1c78f6939c81bf91e53393130fbfee32ff4e9cb6022f14ae7750d126a7b6c0202c674b02",
		// Registration Fee (hardcoded and required at 10 TFT)
		"regfee": "10000000000",
		// Regular Transaction Fee
		"txfee": "1000000000",
		// Coin Inputs to fund the fees
		"coininputs": [{
			"parentid": "a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee563",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
					"signature": "4fe14adcbded85476680bfd4fa8ff35d51ac34bb8a9b3f4904eac6eee4f53e19b6a39c698463499b9961524f026db2fb5c8173307f483c6458d401ecec2e7a0c"
				}
			}
		}],
		// Optional Refund CoinOutput
		// This the same as when sending 5tft  in a regular transaction but your inputssum up to say 100,
		// you also add an output of 95 to your own address then
		"refundcoinoutput": {
			"value": "99999999000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01370af706b547dd4e562a047e6265d7e7750771f9bff633b1a12dbd59b11712c6ef65edb1690d"
				}
			}
		}
	}
}
```

###### Binary Encoding an ERC20 Address Registration Transaction

The binary encoding of an ERC20 Address Registration Transaction uses the Rivine encoding package.
In order to understand the binary encoding of such a transaction, please see [the Rivine encoding documentation][rivine-encoding]
in order to understand how an ERC20 Address Registration Transaction is binary encoded.

The same transaction that was shown as an example of a JSON-encoded ERC20 Address Registration Transaction, can be represented in a hexadecimal string —when binary encoded— as:

```raw
d201a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d80fe13823a96928a573f20a63f3b8d3cde08c506fa535d458120fdaa5f1c78f6939c81bf91e53393130fbfee32ff4e9cb6022f14ae7750d126a7b6c0202c674b020a02540be400083b9aca0002a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee56301c401d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780804fe14adcbded85476680bfd4fa8ff35d51ac34bb8a9b3f4904eac6eee4f53e19b6a39c698463499b9961524f026db2fb5c8173307f483c6458d401ecec2e7a0c01100163457821ef3600014201370af706b547dd4e562a047e6265d7e7750771f9bff633b1a12dbd59b11712c6
```

###### Signing an ERC20 Address Registration Transaction

It is assumed that the reader of this chapter has already
read [Rivine's Introduction to Signing Transactions][rivine-signing-into] and all its referenced content.

> Note though that for the signing of ERC20 Transactions the [Rivine encoding library][rivine-encoding] is used.

In order to sign an ERC20 Address Registration transaction, you first need to compute the hash,
which is used as message, which we'll than to create a signature using the Ed25519 algorithm.

Computing that hash can be represented by following pseudo code:

```plain
blake2b_256_hash(RivineBinaryEncoding(
  - transactionVersion: 1 byte
  - specifier: 16 bytes, hardcoded to "erc20 addrreg tx"
  - public key
  - all extra objects (not the length)
  - for each coin input:
    - parentID
  - registration fee
  - transaction fee
  - ptr(refundCoinOutput))
)) : 32 bytes fixed-size crypto hash
```

[rivine]: https://github.com/threefoldtech/rivine
[sia-encoding]: https://github.com/threefoldtech/rivine/blob/master/doc/encoding/SiaEncoding.md
[rivine-encoding]: https://github.com/threefoldtech/rivine/blob/master/doc/encoding/RivineEncoding.md
[rivine-txs]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md
[rivine-tx-v1]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-v1-transactions
[rivine-condition-uh]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-an-unlockhashcondition
[rivine-condition-tl]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-timelockcondition
[rivine-condition-multisig]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#json-encoding-of-a-multisignaturecondition
[rivine-double-spending]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#double-spend-rules
[rivine-sign-tx]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#signing-a-v1-transaction
[rivine-signing-into]: https://github.com/threefoldtech/rivine/blob/master/doc/transactions/transaction.md#introduction-to-signing-transactions
