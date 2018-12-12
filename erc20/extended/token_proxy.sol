pragma solidity ^0.5.0;

import "./upgradeability_proxy.sol";
import "./token_storage.sol";

contract TokenProxy is UpgradeabilityProxy, TokenStorage {}