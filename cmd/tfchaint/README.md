# Light-client

A light client which can be used to manage tfchain wallets.

## Building

`go build` in this directory

## Creating a wallet

To use the light client, you will need a wallet. You can either create a new wallet, or load an existing one. If an existing one is loaded, you will be prompted
for the seed. Every wallet needs to have a name, this way they can be identified, and multiple wallets can be managed.

Create a new wallet:

```bash
# create a new wallet
./light-client init $walletname

# or load one from an existing seed
./light-client recover $walletname
# you will then be prompted to enter the seed
```

By default, the light client only generates a single address. You can generate more when loading the wallet by passing the `--key-amount` flag, followed by the amount
of addresses to load.

At this time, only testnet is supported.

## Using a wallet

Once a wallet is created, it can be accessed via `./light-client $walletname`. If no subcommand is given, the balance of the generated addresses will be requested. Additional subcommands are also available. 

Listing addresses managed by the wallet is done via the `addresses` subcommand, creating a transaction is done via the `send` command

```bash
# list addresses, one can then be copied and pasted in the faucet webpage to get some testing funds
./light-client $walletname addresses

# once some tokens have been received, you can check the balance
./light-client $walletname

# and a transaction can be created
./light-client $walletname send $amount $address
```

There are some additional options for sending money, such as sending to a multisig address, or time locking the output. For a detailed description of the arguments, and the available flags, you can pass the `-h` or `--help` flag to the command (as well as all other commands). This will print more detailed information about the options.