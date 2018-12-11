pragma solidity ^0.5.0;

import "./upgradeability_storage.sol";
import "./token_storage.sol";

contract UpgradeableTokenStorage is UpgradeabilityStorage, TokenStorage {}