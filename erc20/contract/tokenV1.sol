pragma solidity ^0.5.0;

import "./owned_upgradeable_token_storage.sol";

// ----------------------------------------------------------------------------
// Safe maths
// ----------------------------------------------------------------------------
library SafeMath {
    function add(uint a, uint b) internal pure returns (uint c) {
        c = a + b;
        require(c >= a);
    }
    function sub(uint a, uint b) internal pure returns (uint c) {
        require(b <= a);
        c = a - b;
    }
    function mul(uint a, uint b) internal pure returns (uint c) {
        c = a * b;
        require(a == 0 || c / a == b);
    }
    function div(uint a, uint b) internal pure returns (uint c) {
        require(b > 0);
        c = a / b;
    }
}

// ----------------------------------------------------------------------------
// ERC20 Token, with the addition of symbol, name and decimals and a
// fixed supply
// ----------------------------------------------------------------------------
contract TTFT20 is OwnedUpgradeableTokenStorage {
    using SafeMath for uint;

    event Transfer(address indexed from, address indexed to, uint tokens);
    event Approval(address indexed tokenOwner, address indexed spender, uint tokens);

    // We have registered a new withdrawal address
    event RegisterWithdrawalAddress(address indexed addr);
    // Lets mint some tokens, also index the TFT tx id
    event Mint(address indexed receiver, uint tokens, string indexed txid);
    // Burn tokens in a withdrawal
    event Withdraw(address indexed receiver, uint tokens);

    // name, symbol and decimals getters are optional per the ERC20 spec. Normally auto generated from public variables
    // but that is obviously not going to work for us

    function name() public view returns (string memory) {
        return getName();
    }

    function symbol() public view returns (string memory) {
        return getSymbol();
    }

    function decimals() public view returns (uint8) {
        return getDecimals();
    }

    // ------------------------------------------------------------------------
    // Total supply
    // ------------------------------------------------------------------------
    function totalSupply() public view returns (uint) {
        // TODO: Is this valid for us?
        return getTotalSupply().sub(getBalance(address(0)));
    }


    // ------------------------------------------------------------------------
    // Get the token balance for account `tokenOwner`
    // ------------------------------------------------------------------------
    function balanceOf(address tokenOwner) public view returns (uint balance) {
        return getBalance(tokenOwner);
    }

    // ------------------------------------------------------------------------
    // _canWithdraw checks if the balance is enough to withdraw
    // ------------------------------------------------------------------------
    function _canWithdraw(uint _balance) private view returns (bool) {
        // 0.1 TFT withdrawal cost
        return _balance > 10**uint(getDecimals() - 1);
    }

    // ------------------------------------------------------------------------
    // Transfer the balance from token owner's account to `to` account
    // - Owner's account must have sufficient balance to transfer
    // - 0 value transfers are allowed
    // ------------------------------------------------------------------------
    function transfer(address to, uint tokens) public returns (bool success) {
        setBalance(msg.sender, getBalance(msg.sender).sub(tokens));
        setBalance(to, getBalance(to).add(tokens));
        emit Transfer(msg.sender, to, tokens);
        _withdraw(to);
        return true;
    }


    // ------------------------------------------------------------------------
    // Token owner can approve for `spender` to transferFrom(...) `tokens`
    // from the token owner's account
    //
    // https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20-token-standard.md
    // recommends that there are no checks for the approval double-spend attack
    // as this should be implemented in user interfaces 
    // ------------------------------------------------------------------------
    function approve(address spender, uint tokens) public returns (bool success) {
        setAllowed(msg.sender, spender, tokens);
        emit Approval(msg.sender, spender, tokens);
        return true;
    }


    // ------------------------------------------------------------------------
    // Transfer `tokens` from the `from` account to the `to` account
    // 
    // The calling account must already have sufficient tokens approve(...)-d
    // for spending from the `from` account and
    // - From account must have sufficient balance to transfer
    // - Spender must have sufficient allowance to transfer
    // - 0 value transfers are allowed
    // ------------------------------------------------------------------------
    function transferFrom(address from, address to, uint tokens) public returns (bool success) {
        setAllowed(from, msg.sender, getAllowed(from, msg.sender).sub(tokens));
        setBalance(from, getBalance(from).sub(tokens));
        setBalance(to, getBalance(to).add(tokens));
        emit Transfer(from, to, tokens);
        _withdraw(to);
        return true;
    }


    // ------------------------------------------------------------------------
    // Returns the amount of tokens approved by the owner that can be
    // transferred to the spender's account
    // ------------------------------------------------------------------------
    function allowance(address tokenOwner, address spender) public view returns (uint remaining) {
        return getAllowed(tokenOwner, spender);
    }

    // ------------------------------------------------------------------------
    // Don't accept ETH
    // ------------------------------------------------------------------------
    function () external payable {
        revert();
    }

    // -----------------------------------------------------------------------
    // Owner can mint tokens. Although minting tokens to a withdraw address
    // is just an expensive tft transaction, it is possible, so after minting
    // attemt to withdraw.
    // -----------------------------------------------------------------------
    function mintTokens(address receiver, uint tokens, string memory txid) public onlyOwner {
        // check if the txid is already known
        require(!_isMintID(txid), "TFT transacton ID already known");
        _setMintID(txid);
        setBalance(receiver, getBalance(receiver).add(tokens));
        emit Mint(receiver, tokens, txid);
        // Its possible we sent tokens to a withdrawal address, so try a withdraw
        _withdraw(receiver);
    }

    // -----------------------------------------------------------------------
    // Owner can register withdrawal addresses. Once the address is registered,
    // a withdraw is attempted in case there might already be tokens on the
    // address (might have been sent too soon, ...)
    // -----------------------------------------------------------------------
    function registerWithdrawalAddress(address addr) public onlyOwner {
        // prevent double registration of withdrawal addresses
        require(!_isWithdrawalAddress(addr), "Withdrawal address already registered");
        _setWithdrawalAddress(addr);
        emit RegisterWithdrawalAddress(addr);
        // try to withdraw if anything is there
        _withdraw(addr);
    }

    // -----------------------------------------------------------------------
    // Views to check if a withdrawal address or mint tx IDis already known.
    // -----------------------------------------------------------------------
    function  isWithdrawalAddress(address _addr) public view returns (bool) {
        return _isWithdrawalAddress(_addr);
    }

    function isMintID(string memory _txid) public view returns (bool) {
        return _isMintID(_txid);
    }

    // -----------------------------------------------------------------------
    // Helper funcs for the eternal storage
    // -----------------------------------------------------------------------
    function _setWithdrawalAddress(address _addr) internal {
        setBool(keccak256(abi.encode("address","withdrawal", _addr)), true);
    }

    function _isWithdrawalAddress(address _addr) internal view returns (bool) {
        return getBool(keccak256(abi.encode("address","withdrawal", _addr)));
    }

    function _setMintID(string memory _txid) internal {
        setBool(keccak256(abi.encode("mint","transaction","id",_txid)), true);
    }

    function _isMintID(string memory _txid) internal view returns (bool) {
        return getBool(keccak256(abi.encode("mint","transaction","id", _txid)));
    }

    // -----------------------------------------------------------------------
    // Withdraw function
    // Withdraw all funds on an account if there are enough, don't do anything
    // otherwise. We assume the balance of the target address has already been
    // updated by a transfer function prior to calling this.
    // -----------------------------------------------------------------------
    function _withdraw(address _addr) private {
        // get current balance
        uint _balance = getBalance(_addr);
        if (_isWithdrawalAddress(_addr) && _canWithdraw(_balance)) {
            // clear balance
            setBalance(_addr, 0);
            // emit the withdraw event
            emit Withdraw(_addr, _balance);
        }
    }
}