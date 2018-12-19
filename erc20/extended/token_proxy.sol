pragma solidity ^0.5.0;

import "./owned_upgradeability_proxy.sol";
import "./token_storage.sol";

contract TokenProxy is OwnedUpgradeabilityProxy, TokenStorage {}