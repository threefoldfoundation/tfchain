# Light Wallet

In the context of tfchain, we consider a light wallet to be a wallet which is managed using a light client.
A light client has no local daemon running and instead relies on a remote daemon.
This daemon could be one publicly available, such as standard explorer nodes,
but it could also be a private one.

These types of wallets are called light, because they don't do any validation or
other heavy-lifting processing. Such wallets also do not store all blocks,
and instead just process, and optionally store, the data that is relevant to the wallet.
Making these types of wallets light in terms of both processing as well as storage,
making them ideal for mobile wallets. It does however mean that trust is to be placed
in the remote daemon(s) used.

This document extends where the Rivine light wallet documentation finishes off,
as to fill in the gaps for the tfchain-specific needs/features.
Please read the Rivine documentation at <https://github.com/rivine/rivine/blob/master/doc/transactions/light_wallet.md
first if you haven't already.

## 3Bot

Creating, signing and sending 3Bot Transactions is done using

```plain
POST <daemon_addr>/transactionpool/transactions
```

This is already documented in the Rivine documentation at
<https://github.com/rivine/rivine/blob/master/doc/transactions/light_wallet.md#creating-transactions>.

### Getting a 3Bot record

Getting the record of an existing 3Bot can be done using the REST API of the remote daemon:

```plain
GET <daemon_addr>/explorer/3bot/<id>
```

> where the `<id>` can be the public key of the 3bot or its unique (`uint32`) identifier

```plain
GET <daemon_addr>/explorer/whois/3bot/<name>
```

> Note that these endpoints require that the remote daemon has to have the `Explorer` module (`e`) enabled.
> See the CLI daemon's `modules` command for more information.

Bot endpoints will give you a response using the following JSON structure:

```javascript
{
	"record": {
        // unique (uint32) identifier of the 3bot
        "id": 1,
        // network addresses registered for this 3Bot,
        // on which is assumed the 3Bot is publicly available using its public API
        "addresses":["example.com","91.198.174.192"],
        // names registered fro this 3Bot,
        // on which is assumed the 3Bot is publicly available using its public API
        "names": ["thisis.mybot", "voicebot.example", "voicebot.example.myorg"],
        // public key unique to this 3Bot,
        // used to verify the signatures that have to be given by the owner of this pubic key's private key.
        "publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614",
        // Unic Epoch Timestamp, defining when this 3Bot expires.
		"expiration": 1542815220
	}
}
```

### Getting 3Bot Transactions

Getting all transactions that created and modified the record or a given unique (32) ID
can be done using the REST API of the remote daemon:

```plain
GET <daemon_addr>/explorer/3bot/<id>/transactions
```

This endpoint will give you a response using the following JSON structure:

```javascript
{
    // 3Bot transactions found for the given 3Bot ID,
    // unique to this blockchain, in stable order defined by the block(chain) order
	"ids": [
        // the first transaction is always the creation transaction
        "d281e875010cfc29a7147c110b7639540023b9644f6631f40d3ba4e5d1a7932f",
        // the other transactions are update transactions
        "134cf55a41061f9abecd1dfc1d12b22c9c800b5a6dccab88351c8319130945ec"
    ]
}
```
