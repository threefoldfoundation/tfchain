pragma solidity ^0.5.0;

import "./storage.sol";

contract TokenStorage is Storage {

    // -----------------------------------------------------
    // Usual storage
    // -----------------------------------------------------

    // string public symbol;
    // string public  name;
    // uint8 public decimals;
    // uint _totalSupply;

    // mapping(address => uint) balances;
    // mapping(address => mapping(address => uint)) allowed;

    // -----------------------------------------------------
    // getter utilities
    // -----------------------------------------------------
    function getSymbol() public view returns (string memory) {
        return getString(keccak256(abi.encode("symbol")));
    }

    function getName() public view returns (string memory) {
        return getString(keccak256(abi.encode("name")));
    }

    function getDecimals() public view returns (uint8) {
        return uint8(getUint(keccak256(abi.encode("decimals"))))    ;
    }

    function getTotalSupply() public view returns (uint) {
        return getUint(keccak256(abi.encode("totalSupply")));
    }

    function getBalance(address _account) public view returns (uint) {
        return getUint(keccak256(abi.encode("balance", _account)));
    }

    function getAllowed(address _account, address _spender) public view returns (uint) {
        return getUint(keccak256(abi.encode("allowed", _account, _spender)));
    }

    // -----------------------------------------------------
    // setter utilities
    // -----------------------------------------------------
    function setSymbol(string memory _symbol) internal {
        setString(keccak256(abi.encode("symbol")), _symbol);
    }

    function setName(string memory _name) internal {
        setString(keccak256(abi.encode("name")), _name);
    }

    function setDecimals(uint8 _decimals) internal {
        setUint(keccak256(abi.encode("decimals")), _decimals);
    }

    function setTotalSupply(uint _totalSupply) internal {
        setUint(keccak256(abi.encode("totalSupply")), _totalSupply);
    }

    function setBalance(address _account, uint _balance) internal {
        setUint(keccak256(abi.encode("balance", _account)), _balance);
    }

    function setAllowed(address _account, address _spender, uint _allowance) internal {
        setUint(keccak256(abi.encode("allowed", _account, _spender)), _allowance);
    }
}