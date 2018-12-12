# erc20

## Basic version

Code is in the [basic](./basic) subdirectory.

[a basic solidity contract](basic/basic_contract.sol), and [basic proxy contract](basic/proxy_contract.sol) +
compiled files in the contract subdirectory.

To recompile the contracts after changes, run

```bash
solc --bin -o ./basic/basic_contract basic/basic_contract.sol
solc --bin -o ./basic/proxy_contract basic/proxy_contract.sol
```

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

![](erc20_setup.svg)

+ Proxy, UpgradeabilityProxy and UpgradeabilityStorage are generic contracts
  + Proxy delegates any unknown method to an embedded contract address
  + UpgradeabilityStorage is the contract storage holding the version and address of the actual token implementation
  + UpgraeabilityProxy extends from both Proxy and UpgradeabilityStorage
+ TokenStorage is the stoage for the erc20 token (balances, allowances, ...), which is immutable.
+ TokenProxy extends from UpgradeabilityProxy and TokenStorage, it is the contract which will delegate calls to the
  specific ERC20 contracts. (No actual implementation, only inheritence)
+ UpgradeableTokenStorage extends from UpgradeabilityStorage and TokenStorage. This includes all the required components
  for the upgraded token. (No actual implementation, only inheritence)
+ The specific token contracts extend from UpgradeableTokenStorage, so they all have the same memory structure.

### Important

Notice that inheritance order in TokenProxy needs to be the same as the one in UpgradeableTokenStorage, otherwise the memory layout will
be different, causing the proxy to fail.