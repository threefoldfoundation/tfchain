pragma solidity ^0.5.0;

contract Owned {
    mapping(address => bool) public owner;
    event AddedOwner(address newOwner);
    event RemovedOwner(address removedOwner);

    function Ownable() public {
        owner[msg.sender] = true;
    }

    modifier onlyOwner() {
        require(owner[msg.sender]);
        _;
    }

    function addOwner(address _newOwner) onlyOwner public {
        require(_newOwner != address(0));
        owner[_newOwner] = true;
        emit AddedOwner(_newOwner);
    }

    function removeOwner(address _toRemove) onlyOwner public {
        require(_toRemove != address(0));
        require(_toRemove != msg.sender);
        owner[_toRemove] = false;
        emit RemovedOwner(_toRemove);
    }
}