pragma solidity ^0.5.0;

import "./token_storage.sol";
import "./owned.sol";

contract OwnedTokenStorage is TokenStorage, Owned {}