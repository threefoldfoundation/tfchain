pragma solidity ^0.5.0;

import "./storage.sol";

contract Owned is Storage {
    
    // -----------------------------------------------------
    // Usual storage
    // -----------------------------------------------------

    // mapping(address => bool) public owner;

    // -----------------------------------------------------
    // Events
    // -----------------------------------------------------

    event AddedOwner(address newOwner);
    event RemovedOwner(address removedOwner);

    // -----------------------------------------------------
    // storage utilities
    // -----------------------------------------------------

    function _isOwner(address _caller) internal view returns (bool) {
        return getBool(keccak256(abi.encode("owner",_caller)));
    }

    function _addOwner(address _newOwner) internal {
        setBool(keccak256(abi.encode("owner", _newOwner)), true);
    }

    function _deleteOwner(address _owner) internal {
        deleteBool(keccak256(abi.encode("owner", _owner)));
    }

    // -----------------------------------------------------
    // Main contract
    // -----------------------------------------------------

    constructor() public {
        _addOwner(msg.sender);
    }

    modifier onlyOwner() {
        require(_isOwner(msg.sender));
        _;
    }

    function addOwner(address _newOwner) onlyOwner public {
        require(_newOwner != address(0));
        _addOwner(_newOwner);
        emit AddedOwner(_newOwner);
    }

    function removeOwner(address _toRemove) onlyOwner public {
        require(_toRemove != address(0));
        require(_toRemove != msg.sender);
        _deleteOwner(_toRemove);
        emit RemovedOwner(_toRemove);
    }

}