# tfchainc

Tfchainc is a tfchaind client, it uses the REST API to communicate with the daemon.

The usage is pretty simple, as you only need to execute it.

```bash
tfchainc --help
Tfchain Client v0.1.1

Usage:
  tfchainc [flags]
  tfchainc [command]

Available Commands:
  consensus   Print the current state of consensus
  gateway     Perform gateway actions
  help        Help about any command
  stop        Stop the rivine daemon
  update      Update rivine
  version     Print version information
  wallet      Perform wallet actions

Flags:
  -a, --addr string   which host/port to communicate with (i.e. the host/port tfchaind is listening on) (default "localhost:23110")
  -h, --help          help for ./tfchainc

Use "./tfchainc [command] --help" for more information about a command.
```

The commands let you interact with the daemon

* consensus, will inform you about if the node has consensus, ie: it has the same information other nodes on the network has, or if it is still syncing such information

* gateway, shows you information related to the communications, such as own address and peers connected to your node, also let you create/remove new/existing connections.

* stop, allows you stop the tfchaind in a controlled manner

* update, let you check for newer versions of the software

* wallet, prints information on your wallet, such as addresses, transactions and balances,it lets you send coins, and enables you to initialize, lock/unlock your wallet, or create new addresses.
