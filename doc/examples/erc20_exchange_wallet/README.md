# erc20 exchange wallet

This is a demo application showing the erc20 (T)TFT to regular (T)TFT conversion by simulating the wallet page of an exchange. 

The demo application itself is only connected to the ethereum network


## Running

Build using `go build`, run the produced binary (`./erc20_exchange_wallet`) and point your browser to http://localhost:8080 . 

### Withdrawing

In order to withdraw from the demo, it's account will need to be funded with both ether and tokens. The demo's
address to send to is displayed on the web page.

## ethereum test networks
You can pass will the `--ethnetwork` flag to specify the ethereum network to use.
The only 2 allowed options are `rinkeby` and `ropsten`, ropsten being the default as it is the one used by the tfchain testnet. By default the demo will use the contract addresses defined by tfchain.

## ccontract addresses  
If you are running a development setup, you will likely deploy your own contract. The address of this contract can then be passed with the `--contract-address` flag. 
