pragma solidity ^0.5.0;

import "./proxy.sol";
import "./upgradeability_storage.sol";

contract UpgradeabilityProxy is Proxy, UpgradeabilityStorage {
  event Upgraded(string version, address indexed implementation);

  function upgradeTo(string memory version, address implementation) public {
    require(_implementation != implementation);
    _version = version;
    _implementation = implementation;
    emit Upgraded(version, implementation);
  }
}