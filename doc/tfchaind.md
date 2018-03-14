# tfchaind

tfchaind is the daemon running the tfchain nnode, also provides a REST API to allow clients to connect and interact with it.

The usage is pretty simple, as you only need to execute it.
You can access the embedded help using the --help flag.

```bash
tfchaind --help
Tfchain Daemon v0.1.1

Usage:
  tfchaind [flags]
  tfchaind [command]

Available Commands:
  help        Help about any command
  modules     List available modules for use with -M, --modules flag
  version     Print version information

Flags:
      --agent string               required substring for the user agent (default "Rivine-Agent")
      --api-addr string            which host:port the API server listens on (default "localhost:23110")
      --authenticate-api           enable API password protection
      --disable-api-security       allow tfchaind to listen on a non-localhost address (DANGEROUS)
  -h, --help                       help for ./tfchaind
  -M, --modules string             enabled modules, see 'tfchaind modules' for more info (default "cgtwb")
      --no-bootstrap               disable bootstrapping on this run
      --profile                    enable profiling
      --profile-directory string   location of the profiling directory (default "profiles")
      --rpc-addr string            which port the gateway listens on (default ":23112")
  -d, --tfchain-directory string   location of the tfchain directory

```

Tfchaind has some modules that are optional, so you may run a node without some of the functionality. Additional info on the modules as its dependencies can be shown using the **./tchaind modules** command.

The modules available now are:

* Gateway (abreviated as "g"): is the network connection manager, let your node identify itself and connect with other peers.

* Consensus Set (abreviated as "c"): keeps the chain in sync with the rest of the network.

* Transaction Pool (aka "t"): keeps a pool of unconfirmed transactions.

* Wallet (aka "w"): stores and manages coins and blockstakes.

* BlockCreator (aka "b"): creates new blocks for the chain.

* Explorer (aka "e"): provides statistics, transactions and objects info on the chain.

Some modules have dependencies on other modules.