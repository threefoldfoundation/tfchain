# Explorer
A block explorer for Threefold Chain

A public instance of this explorer for the demo/test net is availabe at https://explorer.testnet.threefoldtoken.com

## Run it yourself

### Prerequisites
* Caddyserver
* Tfchaind daemon


Make sure you have a tfchaind running with the explorer module enabled:
`tfchaind -M cgte`

Now start caddy from the `caddy` folder of this repository:
`caddy -conf Caddyfile.local`
and browse to http://localhost:2015
