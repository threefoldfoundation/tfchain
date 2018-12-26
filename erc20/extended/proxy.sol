pragma solidity ^0.5.0;

import "./token_storage.sol";

// inherit from TokenStorage so we have the constructor, since the token variables need to be stored in the
// proxy's storage
contract Proxy is TokenStorage {
    function () external payable {
        // directly get the implementation contract address from the storage. This way we don't need to depend
        // on the upgradeable contract
        address _impl = getAddress(keccak256(abi.encode("implementation")));
        require(_impl != address(0));
        bytes memory data = msg.data;

        assembly {
            let result := delegatecall(gas, _impl, add(data, 0x20), mload(data), 0, 0)
            let size := returndatasize
            let ptr := mload(0x40)
            returndatacopy(ptr, 0, size)
            switch result
            case 0 { revert(ptr, size) }
            default { return(ptr, size) }
        }
    }
}