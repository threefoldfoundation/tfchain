pragma solidity ^0.5.0;

import "./token_storage.sol";
import "./upgradeable.sol";

contract OwnedUpgradeableTokenStorage is TokenStorage, Upgradeable {}