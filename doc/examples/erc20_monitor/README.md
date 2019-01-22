# erc20 monitor

This is a demo application showing the erc20 (T)TFT to regular (T)TFT conversion. The application leverages the light
client code used by the bridge and tfchain daemon eth validator modules. The demo application itself is only connected
to the ethereum network

## Running

If you want to run this example yourself, simply build the example (`go build`), and run the produced binary
in this directory. You will need to pass at least the `--ethnetwork` flag to specify the ethereum network to use.
The only 2 allowed options are `rinkeby` and `ropsten`. By default the demo will use the contract addresses defined
by tfchain. If you are running a development setup, you will likely need to deploy your own contract. The address
of this contract can then be passed with the `--contract-address` flag. For a list of all supported flags, pass
the `--help` flag.

## Withdrawing

In order to withdraw from the demo, it's account will need to be funded with both ethers and tokens. The demo's
address to send to is displayed on the web page.