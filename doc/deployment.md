# Deployment

## Test net deployment

In order to launch and support the testnet, 4 zero-os nodes have been deployed, each running a `tfchaind` instance
Each node has a public IP address, used by the daemons' rpc calls. The http interface is exposed through the zerotier network
(except for the explorer).

The containers are launched from an flist on the gig hub. First a minimal flist was created using the `Makefile` in the root
of this repository (using the `release-images` target). The resulting flist was then merged with the ubuntu-16.04 flist, which
was already present on the hub. These flists were then used to create a container on each node. For the explorer node, an ssh 
server was also installed, since it required more manual setup (caddy was not yet installed, neither was the explorer repo with
the caddyfile).

The nodes (and thus the container running on the node) can be accessed directly with the python client (provided you are authorized
on the ZeroTier network).

Node overview:
| ZeroTier IP | Public IP     | Exposed port - subnet | Notes                              |
| ----------- | ------------- | --------------------- | ---------------------------------- |
| 10.250.1.11 | 185.69.166.14 | 23112 - 0.0.0.0       | Blockcreator node                  |
|             |               | 23110 - 10.250.1.0/24 |                                    |
| 10.250.1.12 | 185.69.166.11 | 23112 - 0.0.0.0       | Blockcreator node                  |
|             |               | 23110 - 10.250.1.0/24 |                                    |
| 10.250.1.13 | 185.69.166.13 | 23112 - 0.0.0.0       | Blockcreator node                  |
|             |               | 23110 - 10.250.1.0/24 |                                    |
| 10.250.1.14 | 185.69.166.12 | 23112 - 0.0.0.0       | Explorer node                      |
|             |               | 2222 - 10.250.1.0/24  | Ssh server (authorized key access) |
|             |               | 2015 - 0.0.0.0        | Explorer UI                        |