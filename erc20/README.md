# erc20

## Basic version

Code is in the [basic](./basic) subdirectory.

[a basic solidity contract](basic/basic_contract.sol), and [basic proxy contract](basic/proxy_contract.sol) +
compiled files in the contract subdirectory.


## Extended version

Code is in the [extended](./extended) subdirectory

This is a more advanced setup which separates the storage for the ERC20 token from the conctract logic. As such,
the storage model can be reused by both the proxy and the implementations. Although there are multiple different
contracts here, the only 2 which actually matter are the [basic token contract](./extended/basic_contract.sol) and
the [token proxy](./extended/token_proxy.sol) contracts.

## Proxy Contract setup
 
### motivation

Ethereum contracts are autonomous immutable code. Once deployed to the Ethereum blockchain though, they are essentially set in stone. This means that if a serious bug or issue appears and your contracts aren’t designed in a way that will allow them to be upgraded in your Dapp seamlessly, we're screwed.

To  solve this we make two main design choices:
- All main contracts must be upgradable
- Have a flexible, yet simple way to store data permanently

### Setup

The problem with a contract proxy is that the proxy does not actually call the implementing contract. Instead,
it loads the function code and executes it in it's own storage space. This means that the storage needs to be
defined in the proxy as well, and the storage layout needs to be the same.

![contract hierarchy diagram](erc20_setup.svg)

The two main contracts here are TokenProxy  and tokenV0..Vx

**tokenV0..Vx** are the actual upgradeable implementations

**TokenProxy** is the contract that will be called by the users. It delegates all calls to the current tokenV0..Vx implementation. The method are not defined here , every call it receives is delegated so if the implementation adds functionality, this contract does not have to be upgraded.
This way , the address of the deployed TokenProxy never changes.

helper contracts:
+ Proxy provides the functionality to delegate any unknown method to an embedded contract address
+ UpgradeabilityStorage is the contract storage holding the version and address of the actual token implementation
+ UpgradeabilityProxy adds the upgrade logic
+ TokenStorage is the storage for the erc20 token (balances, allowances, ...), which is immutable.
+ TokenProxy extends from UpgradeabilityProxy and TokenStorage, it is the contract which will delegate calls to the
  specific ERC20 contracts. (No actual implementation, only inheritence)
+ UpgradeableTokenStorage extends from UpgradeabilityStorage and TokenStorage. This includes all the required components
  for the upgraded token. (No actual implementation, only inheritence)
  The specific token contracts extend from UpgradeableTokenStorage, so they all have the same memory structure.

### Important

Notice that inheritance order in TokenProxy needs to be the same as the one in UpgradeableTokenStorage, otherwise the memory layout will
be different, causing the proxy to fail.

## Building and deployment

In both the basic and extended folders a `compile.sh` script is present to compile the contracts.

The solidity compiler (`solc`) is required. :
- Mac osx: `brew install solidity`  