# ERC-20 threefold tokens

[ERC-20](https://theethereum.wiki/w/index.php/ERC20_Token_Standard) is a technical standard used for smart contracts on the Ethereum blockchain for implementing tokens.

## Motivation

A lot of the bigger exchanges support ERC20 tokens but it would require some effort to support TFT directly. By creating an ERC20 token (TFT20) to represent TFT, these exchanges can support and allow internal TFT(20) transfers without any effort. 

## Concept

A tft erc-20 token is deployed on the ethereum blockchain, **tft-erc20** called from now on.
A second component is a bridge, doing the conversions from tft to tft-erc20 and back. 

In order to convert normal tft to tft-erc20, the tft wallet creates a convert transaction with an Ethereum target adress and publishes it on the threefold chain.

The bridge picks up this transactions and creates the tft-erc20 tokens (on face value of the TFT tokens burned) and adds them to the target Ethereum address. 
From this moment on, they are normal erc-20 tokens that can be traded, send to other exchanges or Ethereum wallets.

For the reverse, withdrawing from an exchange and converting them back to real tft, a tft wallet publishes a receiving address on the TFT chain first.

The bridge registers the  withdrawal address in the tft-erc20 smart contract and all tft-erc20 received on that address willl be converted back to tft. 
A user can withdraw from an exchange by withdrawing to the registered address. The tft-erc20 contracts destroys the tft-20 tokens and submits a withdrawal event.
The bridge picks up this transaction and publishes a coin creation transaction including the Ethereum transaction ID. 
Every node can enable a validator to only accept valid tft-erc20 to TFT transactions by checking the withdrawal transactions on the Ethereum chain but as long as more than half of the block creator nodes run it, only valid conversions will be included in the threefold chain.

This allows a user to easily convert TFT to tft-erc20 and back, making the process to transfer TFT to an exchange that supports TFT20 and back almost seamless.

The reason for the registration of withdrawal adresses is that an ethereum adress is just 20 bytes, the contract would not be able to see the difference between a transfer and a withdrawal if withdrawing from an exchange without this.
Also, there is no free choice of ethereum address since it would be possible to steal someone else's token otherwise and requiring an ethereum private key and signature would only complicate things a lot.

### demo/test exchange wallet

[A small web application is available that mimics the balance page of an exchange for demo, test and development purposes](examples/erc20_exchange_wallet).

### Ethereum test networks

We use the ethereum ropsten testnetwork for our testnet, the contract is deployed at `0xb821227dBa4Ef9585D31aa494406FD5E47a3db37`.

Rinkeby testnetwork is used during development.

Faucets:
-  rinkeby: https://faucet.rinkeby.io
  Requires a Twitter, Google+ or Facebook post.
- ropsten: https://faucet.ropsten.be 
  1 test ETH/ day can be requested here

## Technical

- [Explanation of the Ethereum contract](../erc20/README.md)
- [Detailed  description of all ERC20 related transactions](transactions.md#erc20-transactions)
    - Transactions relevant for wallets supporting this functionality:
        - [Convert to erc20 transaction](transactions.md#erc20-convert-transaction)
        - [Withdrawal address registration transaction](transactions.md#erc20-address-registration-transaction)
        - [The coin creation transaction](https://github.com/threefoldfoundation/tfchain/blob/bridge_tft_erc20/doc/transactions.md#erc20-coin-creation-transaction) is created by the bridge, a wallet only needs to be able to understand it to take the coin outputs in to account to be able to spend them.

The Go code defining these transactions resides in [transactions_erc20.go](../pkg/types/transactions_erc20.go)

It is possible to query the status of a withdrawal address registration through the explorer by looking up the ERC20 address using the regular `/explorer/hashes/:hash` endpoint:

* If a [ErrStatusNotFound](https://godoc.org/github.com/threefoldtech/rivine/pkg/api#ErrStatusNotFound) error (HTTP Status 204) is returned, than the ERC20 address can be seen as `unregistered`;
* If a reply is returned than the `confirmations` child property from the `erc20info` root property can be checked if a refined (registration) status check is desired:
   * if it is `0`, than the (registration) status can be seen as `unconfirmed` (still in Tx pool)
   * Otherwise the value is `>=1`:
       * `1` being just created, and thus the highest block is the block in which the address is created
       * Should you want want, you can only see the address as confirmed if it has x or more confirmations (e.g. 6),
         and as such differentiate for example between `confirming` and `confirmed`;
