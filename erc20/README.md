# erc20

[a basic solidity contract](basic_contract.sol), and [basic proxy contract](proxy_contract.sol)
+ compiled files in the contract subdirectory.

To recompile the contracts after changes, run

```bash
solc --bin -o ./basic_contract basic_contract.sol
solc --bin -o ./proxy_contract proxy_contract.sol
```

## Proxy

TODO: Figure out how the storage works and make a better implementation. Investigate
if it is possible to move the storage to a fixed contract, so we can only change
the method implementation (which should be fine).
