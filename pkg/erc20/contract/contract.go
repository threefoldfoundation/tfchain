// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ApproveAndCallFallBackABI is the input ABI used to generate the binding from.
const ApproveAndCallFallBackABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ApproveAndCallFallBackBin is the compiled bytecode used for deploying new contracts.
const ApproveAndCallFallBackBin = `0x`

// DeployApproveAndCallFallBack deploys a new Ethereum contract, binding an instance of ApproveAndCallFallBack to it.
func DeployApproveAndCallFallBack(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ApproveAndCallFallBack, error) {
	parsed, err := abi.JSON(strings.NewReader(ApproveAndCallFallBackABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ApproveAndCallFallBackBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ApproveAndCallFallBack{ApproveAndCallFallBackCaller: ApproveAndCallFallBackCaller{contract: contract}, ApproveAndCallFallBackTransactor: ApproveAndCallFallBackTransactor{contract: contract}, ApproveAndCallFallBackFilterer: ApproveAndCallFallBackFilterer{contract: contract}}, nil
}

// ApproveAndCallFallBack is an auto generated Go binding around an Ethereum contract.
type ApproveAndCallFallBack struct {
	ApproveAndCallFallBackCaller     // Read-only binding to the contract
	ApproveAndCallFallBackTransactor // Write-only binding to the contract
	ApproveAndCallFallBackFilterer   // Log filterer for contract events
}

// ApproveAndCallFallBackCaller is an auto generated read-only Go binding around an Ethereum contract.
type ApproveAndCallFallBackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApproveAndCallFallBackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ApproveAndCallFallBackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApproveAndCallFallBackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ApproveAndCallFallBackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApproveAndCallFallBackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ApproveAndCallFallBackSession struct {
	Contract     *ApproveAndCallFallBack // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ApproveAndCallFallBackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ApproveAndCallFallBackCallerSession struct {
	Contract *ApproveAndCallFallBackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ApproveAndCallFallBackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ApproveAndCallFallBackTransactorSession struct {
	Contract     *ApproveAndCallFallBackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ApproveAndCallFallBackRaw is an auto generated low-level Go binding around an Ethereum contract.
type ApproveAndCallFallBackRaw struct {
	Contract *ApproveAndCallFallBack // Generic contract binding to access the raw methods on
}

// ApproveAndCallFallBackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ApproveAndCallFallBackCallerRaw struct {
	Contract *ApproveAndCallFallBackCaller // Generic read-only contract binding to access the raw methods on
}

// ApproveAndCallFallBackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ApproveAndCallFallBackTransactorRaw struct {
	Contract *ApproveAndCallFallBackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewApproveAndCallFallBack creates a new instance of ApproveAndCallFallBack, bound to a specific deployed contract.
func NewApproveAndCallFallBack(address common.Address, backend bind.ContractBackend) (*ApproveAndCallFallBack, error) {
	contract, err := bindApproveAndCallFallBack(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ApproveAndCallFallBack{ApproveAndCallFallBackCaller: ApproveAndCallFallBackCaller{contract: contract}, ApproveAndCallFallBackTransactor: ApproveAndCallFallBackTransactor{contract: contract}, ApproveAndCallFallBackFilterer: ApproveAndCallFallBackFilterer{contract: contract}}, nil
}

// NewApproveAndCallFallBackCaller creates a new read-only instance of ApproveAndCallFallBack, bound to a specific deployed contract.
func NewApproveAndCallFallBackCaller(address common.Address, caller bind.ContractCaller) (*ApproveAndCallFallBackCaller, error) {
	contract, err := bindApproveAndCallFallBack(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ApproveAndCallFallBackCaller{contract: contract}, nil
}

// NewApproveAndCallFallBackTransactor creates a new write-only instance of ApproveAndCallFallBack, bound to a specific deployed contract.
func NewApproveAndCallFallBackTransactor(address common.Address, transactor bind.ContractTransactor) (*ApproveAndCallFallBackTransactor, error) {
	contract, err := bindApproveAndCallFallBack(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ApproveAndCallFallBackTransactor{contract: contract}, nil
}

// NewApproveAndCallFallBackFilterer creates a new log filterer instance of ApproveAndCallFallBack, bound to a specific deployed contract.
func NewApproveAndCallFallBackFilterer(address common.Address, filterer bind.ContractFilterer) (*ApproveAndCallFallBackFilterer, error) {
	contract, err := bindApproveAndCallFallBack(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ApproveAndCallFallBackFilterer{contract: contract}, nil
}

// bindApproveAndCallFallBack binds a generic wrapper to an already deployed contract.
func bindApproveAndCallFallBack(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ApproveAndCallFallBackABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ApproveAndCallFallBack *ApproveAndCallFallBackRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ApproveAndCallFallBack.Contract.ApproveAndCallFallBackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ApproveAndCallFallBack *ApproveAndCallFallBackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.Contract.ApproveAndCallFallBackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ApproveAndCallFallBack *ApproveAndCallFallBackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.Contract.ApproveAndCallFallBackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ApproveAndCallFallBack *ApproveAndCallFallBackCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ApproveAndCallFallBack.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ApproveAndCallFallBack *ApproveAndCallFallBackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ApproveAndCallFallBack *ApproveAndCallFallBackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.Contract.contract.Transact(opts, method, params...)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 tokens, address token, bytes data) returns()
func (_ApproveAndCallFallBack *ApproveAndCallFallBackTransactor) ReceiveApproval(opts *bind.TransactOpts, from common.Address, tokens *big.Int, token common.Address, data []byte) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.contract.Transact(opts, "receiveApproval", from, tokens, token, data)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 tokens, address token, bytes data) returns()
func (_ApproveAndCallFallBack *ApproveAndCallFallBackSession) ReceiveApproval(from common.Address, tokens *big.Int, token common.Address, data []byte) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.Contract.ReceiveApproval(&_ApproveAndCallFallBack.TransactOpts, from, tokens, token, data)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 tokens, address token, bytes data) returns()
func (_ApproveAndCallFallBack *ApproveAndCallFallBackTransactorSession) ReceiveApproval(from common.Address, tokens *big.Int, token common.Address, data []byte) (*types.Transaction, error) {
	return _ApproveAndCallFallBack.Contract.ReceiveApproval(&_ApproveAndCallFallBack.TransactOpts, from, tokens, token, data)
}

// ERC20InterfaceABI is the input ABI used to generate the binding from.
const ERC20InterfaceABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20InterfaceBin is the compiled bytecode used for deploying new contracts.
const ERC20InterfaceBin = `0x`

// DeployERC20Interface deploys a new Ethereum contract, binding an instance of ERC20Interface to it.
func DeployERC20Interface(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Interface, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20InterfaceBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Interface{ERC20InterfaceCaller: ERC20InterfaceCaller{contract: contract}, ERC20InterfaceTransactor: ERC20InterfaceTransactor{contract: contract}, ERC20InterfaceFilterer: ERC20InterfaceFilterer{contract: contract}}, nil
}

// ERC20Interface is an auto generated Go binding around an Ethereum contract.
type ERC20Interface struct {
	ERC20InterfaceCaller     // Read-only binding to the contract
	ERC20InterfaceTransactor // Write-only binding to the contract
	ERC20InterfaceFilterer   // Log filterer for contract events
}

// ERC20InterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20InterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20InterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20InterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20InterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20InterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20InterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20InterfaceSession struct {
	Contract     *ERC20Interface   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20InterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20InterfaceCallerSession struct {
	Contract *ERC20InterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ERC20InterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20InterfaceTransactorSession struct {
	Contract     *ERC20InterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ERC20InterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20InterfaceRaw struct {
	Contract *ERC20Interface // Generic contract binding to access the raw methods on
}

// ERC20InterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20InterfaceCallerRaw struct {
	Contract *ERC20InterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20InterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20InterfaceTransactorRaw struct {
	Contract *ERC20InterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Interface creates a new instance of ERC20Interface, bound to a specific deployed contract.
func NewERC20Interface(address common.Address, backend bind.ContractBackend) (*ERC20Interface, error) {
	contract, err := bindERC20Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Interface{ERC20InterfaceCaller: ERC20InterfaceCaller{contract: contract}, ERC20InterfaceTransactor: ERC20InterfaceTransactor{contract: contract}, ERC20InterfaceFilterer: ERC20InterfaceFilterer{contract: contract}}, nil
}

// NewERC20InterfaceCaller creates a new read-only instance of ERC20Interface, bound to a specific deployed contract.
func NewERC20InterfaceCaller(address common.Address, caller bind.ContractCaller) (*ERC20InterfaceCaller, error) {
	contract, err := bindERC20Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20InterfaceCaller{contract: contract}, nil
}

// NewERC20InterfaceTransactor creates a new write-only instance of ERC20Interface, bound to a specific deployed contract.
func NewERC20InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20InterfaceTransactor, error) {
	contract, err := bindERC20Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20InterfaceTransactor{contract: contract}, nil
}

// NewERC20InterfaceFilterer creates a new log filterer instance of ERC20Interface, bound to a specific deployed contract.
func NewERC20InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20InterfaceFilterer, error) {
	contract, err := bindERC20Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20InterfaceFilterer{contract: contract}, nil
}

// bindERC20Interface binds a generic wrapper to an already deployed contract.
func bindERC20Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Interface *ERC20InterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Interface.Contract.ERC20InterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Interface *ERC20InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Interface.Contract.ERC20InterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Interface *ERC20InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Interface.Contract.ERC20InterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Interface *ERC20InterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Interface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Interface *ERC20InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Interface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Interface *ERC20InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Interface.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceCaller) Allowance(opts *bind.CallOpts, tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.contract.Call(opts, out, "allowance", tokenOwner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceSession) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Interface.Contract.Allowance(&_ERC20Interface.CallOpts, tokenOwner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceCallerSession) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Interface.Contract.Allowance(&_ERC20Interface.CallOpts, tokenOwner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceCaller) BalanceOf(opts *bind.CallOpts, tokenOwner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.contract.Call(opts, out, "balanceOf", tokenOwner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _ERC20Interface.Contract.BalanceOf(&_ERC20Interface.CallOpts, tokenOwner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceCallerSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _ERC20Interface.Contract.BalanceOf(&_ERC20Interface.CallOpts, tokenOwner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceSession) TotalSupply() (*big.Int, error) {
	return _ERC20Interface.Contract.TotalSupply(&_ERC20Interface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Interface.Contract.TotalSupply(&_ERC20Interface.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceTransactor) Approve(opts *bind.TransactOpts, spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.contract.Transact(opts, "approve", spender, tokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceSession) Approve(spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.Contract.Approve(&_ERC20Interface.TransactOpts, spender, tokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceTransactorSession) Approve(spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.Contract.Approve(&_ERC20Interface.TransactOpts, spender, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceTransactor) Transfer(opts *bind.TransactOpts, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.contract.Transact(opts, "transfer", to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceSession) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.Contract.Transfer(&_ERC20Interface.TransactOpts, to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceTransactorSession) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.Contract.Transfer(&_ERC20Interface.TransactOpts, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.contract.Transact(opts, "transferFrom", from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceSession) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.Contract.TransferFrom(&_ERC20Interface.TransactOpts, from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceTransactorSession) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _ERC20Interface.Contract.TransferFrom(&_ERC20Interface.TransactOpts, from, to, tokens)
}

// ERC20InterfaceApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20Interface contract.
type ERC20InterfaceApprovalIterator struct {
	Event *ERC20InterfaceApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC20InterfaceApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20InterfaceApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC20InterfaceApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC20InterfaceApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20InterfaceApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20InterfaceApproval represents a Approval event raised by the ERC20Interface contract.
type ERC20InterfaceApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_ERC20Interface *ERC20InterfaceFilterer) FilterApproval(opts *bind.FilterOpts, tokenOwner []common.Address, spender []common.Address) (*ERC20InterfaceApprovalIterator, error) {

	var tokenOwnerRule []interface{}
	for _, tokenOwnerItem := range tokenOwner {
		tokenOwnerRule = append(tokenOwnerRule, tokenOwnerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Interface.contract.FilterLogs(opts, "Approval", tokenOwnerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20InterfaceApprovalIterator{contract: _ERC20Interface.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_ERC20Interface *ERC20InterfaceFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20InterfaceApproval, tokenOwner []common.Address, spender []common.Address) (event.Subscription, error) {

	var tokenOwnerRule []interface{}
	for _, tokenOwnerItem := range tokenOwner {
		tokenOwnerRule = append(tokenOwnerRule, tokenOwnerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Interface.contract.WatchLogs(opts, "Approval", tokenOwnerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20InterfaceApproval)
				if err := _ERC20Interface.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ERC20InterfaceTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Interface contract.
type ERC20InterfaceTransferIterator struct {
	Event *ERC20InterfaceTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC20InterfaceTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20InterfaceTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC20InterfaceTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC20InterfaceTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20InterfaceTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20InterfaceTransfer represents a Transfer event raised by the ERC20Interface contract.
type ERC20InterfaceTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_ERC20Interface *ERC20InterfaceFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20InterfaceTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Interface.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20InterfaceTransferIterator{contract: _ERC20Interface.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_ERC20Interface *ERC20InterfaceFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20InterfaceTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Interface.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20InterfaceTransfer)
				if err := _ERC20Interface.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
const OwnedBin = `0x608060405234801561001057600080fd5b5060008054600160a060020a03191633179055610208806100326000396000f3fe60806040526004361061005b577c0100000000000000000000000000000000000000000000000000000000600035046379ba509781146100605780638da5cb5b14610077578063d4ee1d90146100a8578063f2fde38b146100bd575b600080fd5b34801561006c57600080fd5b506100756100f0565b005b34801561008357600080fd5b5061008c610178565b60408051600160a060020a039092168252519081900360200190f35b3480156100b457600080fd5b5061008c610187565b3480156100c957600080fd5b50610075600480360360208110156100e057600080fd5b5035600160a060020a0316610196565b600154600160a060020a0316331461010757600080fd5b60015460008054604051600160a060020a0393841693909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600180546000805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600054600160a060020a031681565b600154600160a060020a031681565b600054600160a060020a031633146101ad57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039290921691909117905556fea165627a7a723058207c17f422cb9cad1f59cd9eef1c11f078517992807c981f8297b8256843f25ee60029`

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// NewOwner is a free data retrieval call binding the contract method 0xd4ee1d90.
//
// Solidity: function newOwner() constant returns(address)
func (_Owned *OwnedCaller) NewOwner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "newOwner")
	return *ret0, err
}

// NewOwner is a free data retrieval call binding the contract method 0xd4ee1d90.
//
// Solidity: function newOwner() constant returns(address)
func (_Owned *OwnedSession) NewOwner() (common.Address, error) {
	return _Owned.Contract.NewOwner(&_Owned.CallOpts)
}

// NewOwner is a free data retrieval call binding the contract method 0xd4ee1d90.
//
// Solidity: function newOwner() constant returns(address)
func (_Owned *OwnedCallerSession) NewOwner() (common.Address, error) {
	return _Owned.Contract.NewOwner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_Owned *OwnedTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_Owned *OwnedSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_Owned *OwnedTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _newOwner)
}

// OwnedOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Owned contract.
type OwnedOwnershipTransferredIterator struct {
	Event *OwnedOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnedOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnedOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnedOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnershipTransferred represents a OwnershipTransferred event raised by the Owned contract.
type OwnedOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed _from, address indexed _to)
func (_Owned *OwnedFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*OwnedOwnershipTransferredIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnershipTransferred", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &OwnedOwnershipTransferredIterator{contract: _Owned.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed _from, address indexed _to)
func (_Owned *OwnedFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnedOwnershipTransferred, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnershipTransferred", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedOwnershipTransferred)
				if err := _Owned.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a7230582066e53a8912e3079393bf04474cb2de22946a547972b11873e061425db04ac5390029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// TTFT20ABI is the input ABI used to generate the binding from.
const TTFT20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"registerWithdrawalAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"approveAndCall\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferAnyERC20Token\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"},{\"name\":\"txid\",\"type\":\"string\"}],\"name\":\"mintTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"RegisterWithdrawalAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"txid\",\"type\":\"string\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// TTFT20Bin is the compiled bytecode used for deploying new contracts.
const TTFT20Bin = `0x608060405234801561001057600080fd5b5060008054600160a060020a031916331790556040805180820190915260068082527f545446543230000000000000000000000000000000000000000000000000000060209092019182526100679160029161012a565b506040805180820190915260198082527f5454465420455243323020726570726573656e746174696f6e0000000000000060209092019182526100ac9160039161012a565b5060048054601260ff19909116179081905560ff16600a0a620f424002600581905560008054600160a060020a0390811682526006602090815260408084208590558354815195865290519216937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a36101c5565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061016b57805160ff1916838001178555610198565b82800160010185558215610198579182015b8281111561019857825182559160200191906001019061017d565b506101a49291506101a8565b5090565b6101c291905b808211156101a457600081556001016101ae565b90565b610eff80620001d56000396000f3fe6080604052600436106100ea577c0100000000000000000000000000000000000000000000000000000000600035046306fdde0381146100ef578063095ea7b31461017957806318160ddd146101c657806323b872dd146101ed578063313ce5671461023057806334ca6a711461025b57806370a082311461029057806379ba5097146102c35780638da5cb5b146102d857806395d89b4114610309578063a9059cbb1461031e578063cae9ca5114610357578063d4ee1d901461041f578063dc39d06d14610434578063dd62ed3e1461046d578063e67524a3146104a8578063f2fde38b14610570575b600080fd5b3480156100fb57600080fd5b506101046105a3565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561013e578181015183820152602001610126565b50505050905090810190601f16801561016b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561018557600080fd5b506101b26004803603604081101561019c57600080fd5b50600160a060020a038135169060200135610631565b604080519115158252519081900360200190f35b3480156101d257600080fd5b506101db610698565b60408051918252519081900360200190f35b3480156101f957600080fd5b506101b26004803603606081101561021057600080fd5b50600160a060020a038135811691602081013590911690604001356106db565b34801561023c57600080fd5b506102456107e6565b6040805160ff9092168252519081900360200190f35b34801561026757600080fd5b5061028e6004803603602081101561027e57600080fd5b5035600160a060020a03166107ef565b005b34801561029c57600080fd5b506101db600480360360208110156102b357600080fd5b5035600160a060020a03166108d9565b3480156102cf57600080fd5b5061028e6108f4565b3480156102e457600080fd5b506102ed61097c565b60408051600160a060020a039092168252519081900360200190f35b34801561031557600080fd5b5061010461098b565b34801561032a57600080fd5b506101b26004803603604081101561034157600080fd5b50600160a060020a0381351690602001356109e3565b34801561036357600080fd5b506101b26004803603606081101561037a57600080fd5b600160a060020a03823516916020810135918101906060810160408201356401000000008111156103aa57600080fd5b8201836020820111156103bc57600080fd5b803590602001918460018302840111640100000000831117156103de57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610af9945050505050565b34801561042b57600080fd5b506102ed610c5a565b34801561044057600080fd5b506101b26004803603604081101561045757600080fd5b50600160a060020a038135169060200135610c69565b34801561047957600080fd5b506101db6004803603604081101561049057600080fd5b50600160a060020a0381358116916020013516610d24565b3480156104b457600080fd5b5061028e600480360360608110156104cb57600080fd5b600160a060020a03823516916020810135918101906060810160408201356401000000008111156104fb57600080fd5b82018360208201111561050d57600080fd5b8035906020019184600183028401116401000000008311171561052f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610d4f945050505050565b34801561057c57600080fd5b5061028e6004803603602081101561059357600080fd5b5035600160a060020a0316610e68565b6003805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156106295780601f106105fe57610100808354040283529160200191610629565b820191906000526020600020905b81548152906001019060200180831161060c57829003601f168201915b505050505081565b336000818152600760209081526040808320600160a060020a038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a35060015b92915050565b600080805260066020527f54cdd369e4e8a8515e52ca72ec816c2101831ad1f18bf44102ed171459c9b4f8546005546106d69163ffffffff610eae16565b905090565b600160a060020a038316600090815260066020526040812054610704908363ffffffff610eae16565b600160a060020a0385166000908152600660209081526040808320939093556007815282822033835290522054610741908363ffffffff610eae16565b600160a060020a038086166000908152600760209081526040808320338452825280832094909455918616815260069091522054610785908363ffffffff610ec316565b600160a060020a0380851660008181526006602090815260409182902094909455805186815290519193928816927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a35060019392505050565b60045460ff1681565b600054600160a060020a0316331461080657600080fd5b600160a060020a0381166000908152600860209081526040808320805460ff19166001179055600690915281205411156108a257600160a060020a03811660008181526006602090815260409182902054825190815291517f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649281900390910190a2600160a060020a0381166000908152600660205260408120555b604051600160a060020a038216907f77bc19082a31daad021d73c26bb4f6e74100a41c98099405e92a9323d133e60290600090a250565b600160a060020a031660009081526006602052604090205490565b600154600160a060020a0316331461090b57600080fd5b60015460008054604051600160a060020a0393841693909116917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600180546000805473ffffffffffffffffffffffffffffffffffffffff19908116600160a060020a03841617909155169055565b600054600160a060020a031681565b6002805460408051602060018416156101000260001901909316849004601f810184900484028201840190925281815292918301828280156106295780601f106105fe57610100808354040283529160200191610629565b600160a060020a03821660009081526008602052604081205460ff1615610a4857604080518381529051600160a060020a038516917f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364919081900360200190a2610af0565b33600090815260066020526040902054610a68908363ffffffff610eae16565b3360009081526006602052604080822092909255600160a060020a03851681522054610a9a908363ffffffff610ec316565b600160a060020a0384166000818152600660209081526040918290209390935580518581529051919233927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35b50600192915050565b336000818152600760209081526040808320600160a060020a038816808552908352818420879055815187815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a36040517f8f4ffcb10000000000000000000000000000000000000000000000000000000081523360048201818152602483018690523060448401819052608060648501908152865160848601528651600160a060020a038a1695638f4ffcb195948a94938a939192909160a490910190602085019080838360005b83811015610be9578181015183820152602001610bd1565b50505050905090810190601f168015610c165780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b158015610c3857600080fd5b505af1158015610c4c573d6000803e3d6000fd5b506001979650505050505050565b600154600160a060020a031681565b60008054600160a060020a03163314610c8157600080fd5b60008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831660048201526024810186905290519186169263a9059cbb926044808401936020939083900390910190829087803b158015610cf157600080fd5b505af1158015610d05573d6000803e3d6000fd5b505050506040513d6020811015610d1b57600080fd5b50519392505050565b600160a060020a03918216600090815260076020908152604080832093909416825291909152205490565b600054600160a060020a03163314610d6657600080fd5b600160a060020a038316600090815260066020526040902054610d8f908363ffffffff610ec316565b6006600085600160a060020a0316600160a060020a031681526020019081526020016000208190555082600160a060020a03167f85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a83836040518083815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610e28578181015183820152602001610e10565b50505050905090810190601f168015610e555780820380516001836020036101000a031916815260200191505b50935050505060405180910390a2505050565b600054600160a060020a03163314610e7f57600080fd5b6001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600082821115610ebd57600080fd5b50900390565b8181018281101561069257600080fdfea165627a7a7230582095b9402f247fa93b55afe31f5cee2afa92b5445ebcb2a34e0a1437092b193cb30029`

// DeployTTFT20 deploys a new Ethereum contract, binding an instance of TTFT20 to it.
func DeployTTFT20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TTFT20, error) {
	parsed, err := abi.JSON(strings.NewReader(TTFT20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TTFT20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TTFT20{TTFT20Caller: TTFT20Caller{contract: contract}, TTFT20Transactor: TTFT20Transactor{contract: contract}, TTFT20Filterer: TTFT20Filterer{contract: contract}}, nil
}

// TTFT20 is an auto generated Go binding around an Ethereum contract.
type TTFT20 struct {
	TTFT20Caller     // Read-only binding to the contract
	TTFT20Transactor // Write-only binding to the contract
	TTFT20Filterer   // Log filterer for contract events
}

// TTFT20Caller is an auto generated read-only Go binding around an Ethereum contract.
type TTFT20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TTFT20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TTFT20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TTFT20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TTFT20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TTFT20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TTFT20Session struct {
	Contract     *TTFT20           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TTFT20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TTFT20CallerSession struct {
	Contract *TTFT20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TTFT20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TTFT20TransactorSession struct {
	Contract     *TTFT20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TTFT20Raw is an auto generated low-level Go binding around an Ethereum contract.
type TTFT20Raw struct {
	Contract *TTFT20 // Generic contract binding to access the raw methods on
}

// TTFT20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TTFT20CallerRaw struct {
	Contract *TTFT20Caller // Generic read-only contract binding to access the raw methods on
}

// TTFT20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TTFT20TransactorRaw struct {
	Contract *TTFT20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTTFT20 creates a new instance of TTFT20, bound to a specific deployed contract.
func NewTTFT20(address common.Address, backend bind.ContractBackend) (*TTFT20, error) {
	contract, err := bindTTFT20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TTFT20{TTFT20Caller: TTFT20Caller{contract: contract}, TTFT20Transactor: TTFT20Transactor{contract: contract}, TTFT20Filterer: TTFT20Filterer{contract: contract}}, nil
}

// NewTTFT20Caller creates a new read-only instance of TTFT20, bound to a specific deployed contract.
func NewTTFT20Caller(address common.Address, caller bind.ContractCaller) (*TTFT20Caller, error) {
	contract, err := bindTTFT20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TTFT20Caller{contract: contract}, nil
}

// NewTTFT20Transactor creates a new write-only instance of TTFT20, bound to a specific deployed contract.
func NewTTFT20Transactor(address common.Address, transactor bind.ContractTransactor) (*TTFT20Transactor, error) {
	contract, err := bindTTFT20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TTFT20Transactor{contract: contract}, nil
}

// NewTTFT20Filterer creates a new log filterer instance of TTFT20, bound to a specific deployed contract.
func NewTTFT20Filterer(address common.Address, filterer bind.ContractFilterer) (*TTFT20Filterer, error) {
	contract, err := bindTTFT20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TTFT20Filterer{contract: contract}, nil
}

// bindTTFT20 binds a generic wrapper to an already deployed contract.
func bindTTFT20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TTFT20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TTFT20 *TTFT20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TTFT20.Contract.TTFT20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TTFT20 *TTFT20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TTFT20.Contract.TTFT20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TTFT20 *TTFT20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TTFT20.Contract.TTFT20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TTFT20 *TTFT20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TTFT20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TTFT20 *TTFT20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TTFT20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TTFT20 *TTFT20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TTFT20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_TTFT20 *TTFT20Caller) Allowance(opts *bind.CallOpts, tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "allowance", tokenOwner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_TTFT20 *TTFT20Session) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _TTFT20.Contract.Allowance(&_TTFT20.CallOpts, tokenOwner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_TTFT20 *TTFT20CallerSession) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _TTFT20.Contract.Allowance(&_TTFT20.CallOpts, tokenOwner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_TTFT20 *TTFT20Caller) BalanceOf(opts *bind.CallOpts, tokenOwner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "balanceOf", tokenOwner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_TTFT20 *TTFT20Session) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _TTFT20.Contract.BalanceOf(&_TTFT20.CallOpts, tokenOwner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_TTFT20 *TTFT20CallerSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _TTFT20.Contract.BalanceOf(&_TTFT20.CallOpts, tokenOwner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TTFT20 *TTFT20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TTFT20 *TTFT20Session) Decimals() (uint8, error) {
	return _TTFT20.Contract.Decimals(&_TTFT20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TTFT20 *TTFT20CallerSession) Decimals() (uint8, error) {
	return _TTFT20.Contract.Decimals(&_TTFT20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TTFT20 *TTFT20Caller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TTFT20 *TTFT20Session) Name() (string, error) {
	return _TTFT20.Contract.Name(&_TTFT20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TTFT20 *TTFT20CallerSession) Name() (string, error) {
	return _TTFT20.Contract.Name(&_TTFT20.CallOpts)
}

// NewOwner is a free data retrieval call binding the contract method 0xd4ee1d90.
//
// Solidity: function newOwner() constant returns(address)
func (_TTFT20 *TTFT20Caller) NewOwner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "newOwner")
	return *ret0, err
}

// NewOwner is a free data retrieval call binding the contract method 0xd4ee1d90.
//
// Solidity: function newOwner() constant returns(address)
func (_TTFT20 *TTFT20Session) NewOwner() (common.Address, error) {
	return _TTFT20.Contract.NewOwner(&_TTFT20.CallOpts)
}

// NewOwner is a free data retrieval call binding the contract method 0xd4ee1d90.
//
// Solidity: function newOwner() constant returns(address)
func (_TTFT20 *TTFT20CallerSession) NewOwner() (common.Address, error) {
	return _TTFT20.Contract.NewOwner(&_TTFT20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TTFT20 *TTFT20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TTFT20 *TTFT20Session) Owner() (common.Address, error) {
	return _TTFT20.Contract.Owner(&_TTFT20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TTFT20 *TTFT20CallerSession) Owner() (common.Address, error) {
	return _TTFT20.Contract.Owner(&_TTFT20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TTFT20 *TTFT20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TTFT20 *TTFT20Session) Symbol() (string, error) {
	return _TTFT20.Contract.Symbol(&_TTFT20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TTFT20 *TTFT20CallerSession) Symbol() (string, error) {
	return _TTFT20.Contract.Symbol(&_TTFT20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TTFT20 *TTFT20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TTFT20 *TTFT20Session) TotalSupply() (*big.Int, error) {
	return _TTFT20.Contract.TotalSupply(&_TTFT20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TTFT20 *TTFT20CallerSession) TotalSupply() (*big.Int, error) {
	return _TTFT20.Contract.TotalSupply(&_TTFT20.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TTFT20 *TTFT20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TTFT20 *TTFT20Session) AcceptOwnership() (*types.Transaction, error) {
	return _TTFT20.Contract.AcceptOwnership(&_TTFT20.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_TTFT20 *TTFT20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TTFT20.Contract.AcceptOwnership(&_TTFT20.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "approve", spender, tokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Session) Approve(spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Approve(&_TTFT20.TransactOpts, spender, tokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20TransactorSession) Approve(spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Approve(&_TTFT20.TransactOpts, spender, tokens)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 tokens, bytes data) returns(bool success)
func (_TTFT20 *TTFT20Transactor) ApproveAndCall(opts *bind.TransactOpts, spender common.Address, tokens *big.Int, data []byte) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "approveAndCall", spender, tokens, data)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 tokens, bytes data) returns(bool success)
func (_TTFT20 *TTFT20Session) ApproveAndCall(spender common.Address, tokens *big.Int, data []byte) (*types.Transaction, error) {
	return _TTFT20.Contract.ApproveAndCall(&_TTFT20.TransactOpts, spender, tokens, data)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 tokens, bytes data) returns(bool success)
func (_TTFT20 *TTFT20TransactorSession) ApproveAndCall(spender common.Address, tokens *big.Int, data []byte) (*types.Transaction, error) {
	return _TTFT20.Contract.ApproveAndCall(&_TTFT20.TransactOpts, spender, tokens, data)
}

// MintTokens is a paid mutator transaction binding the contract method 0xe67524a3.
//
// Solidity: function mintTokens(address receiver, uint256 tokens, string txid) returns()
func (_TTFT20 *TTFT20Transactor) MintTokens(opts *bind.TransactOpts, receiver common.Address, tokens *big.Int, txid string) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "mintTokens", receiver, tokens, txid)
}

// MintTokens is a paid mutator transaction binding the contract method 0xe67524a3.
//
// Solidity: function mintTokens(address receiver, uint256 tokens, string txid) returns()
func (_TTFT20 *TTFT20Session) MintTokens(receiver common.Address, tokens *big.Int, txid string) (*types.Transaction, error) {
	return _TTFT20.Contract.MintTokens(&_TTFT20.TransactOpts, receiver, tokens, txid)
}

// MintTokens is a paid mutator transaction binding the contract method 0xe67524a3.
//
// Solidity: function mintTokens(address receiver, uint256 tokens, string txid) returns()
func (_TTFT20 *TTFT20TransactorSession) MintTokens(receiver common.Address, tokens *big.Int, txid string) (*types.Transaction, error) {
	return _TTFT20.Contract.MintTokens(&_TTFT20.TransactOpts, receiver, tokens, txid)
}

// RegisterWithdrawalAddress is a paid mutator transaction binding the contract method 0x34ca6a71.
//
// Solidity: function registerWithdrawalAddress(address addr) returns()
func (_TTFT20 *TTFT20Transactor) RegisterWithdrawalAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "registerWithdrawalAddress", addr)
}

// RegisterWithdrawalAddress is a paid mutator transaction binding the contract method 0x34ca6a71.
//
// Solidity: function registerWithdrawalAddress(address addr) returns()
func (_TTFT20 *TTFT20Session) RegisterWithdrawalAddress(addr common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.RegisterWithdrawalAddress(&_TTFT20.TransactOpts, addr)
}

// RegisterWithdrawalAddress is a paid mutator transaction binding the contract method 0x34ca6a71.
//
// Solidity: function registerWithdrawalAddress(address addr) returns()
func (_TTFT20 *TTFT20TransactorSession) RegisterWithdrawalAddress(addr common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.RegisterWithdrawalAddress(&_TTFT20.TransactOpts, addr)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "transfer", to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Session) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Transfer(&_TTFT20.TransactOpts, to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20TransactorSession) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Transfer(&_TTFT20.TransactOpts, to, tokens)
}

// TransferAnyERC20Token is a paid mutator transaction binding the contract method 0xdc39d06d.
//
// Solidity: function transferAnyERC20Token(address tokenAddress, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Transactor) TransferAnyERC20Token(opts *bind.TransactOpts, tokenAddress common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "transferAnyERC20Token", tokenAddress, tokens)
}

// TransferAnyERC20Token is a paid mutator transaction binding the contract method 0xdc39d06d.
//
// Solidity: function transferAnyERC20Token(address tokenAddress, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Session) TransferAnyERC20Token(tokenAddress common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferAnyERC20Token(&_TTFT20.TransactOpts, tokenAddress, tokens)
}

// TransferAnyERC20Token is a paid mutator transaction binding the contract method 0xdc39d06d.
//
// Solidity: function transferAnyERC20Token(address tokenAddress, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20TransactorSession) TransferAnyERC20Token(tokenAddress common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferAnyERC20Token(&_TTFT20.TransactOpts, tokenAddress, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "transferFrom", from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20Session) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferFrom(&_TTFT20.TransactOpts, from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_TTFT20 *TTFT20TransactorSession) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferFrom(&_TTFT20.TransactOpts, from, to, tokens)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_TTFT20 *TTFT20Transactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_TTFT20 *TTFT20Session) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferOwnership(&_TTFT20.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_TTFT20 *TTFT20TransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferOwnership(&_TTFT20.TransactOpts, _newOwner)
}

// TTFT20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TTFT20 contract.
type TTFT20ApprovalIterator struct {
	Event *TTFT20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TTFT20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TTFT20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TTFT20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20Approval represents a Approval event raised by the TTFT20 contract.
type TTFT20Approval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_TTFT20 *TTFT20Filterer) FilterApproval(opts *bind.FilterOpts, tokenOwner []common.Address, spender []common.Address) (*TTFT20ApprovalIterator, error) {

	var tokenOwnerRule []interface{}
	for _, tokenOwnerItem := range tokenOwner {
		tokenOwnerRule = append(tokenOwnerRule, tokenOwnerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Approval", tokenOwnerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20ApprovalIterator{contract: _TTFT20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_TTFT20 *TTFT20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TTFT20Approval, tokenOwner []common.Address, spender []common.Address) (event.Subscription, error) {

	var tokenOwnerRule []interface{}
	for _, tokenOwnerItem := range tokenOwner {
		tokenOwnerRule = append(tokenOwnerRule, tokenOwnerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Approval", tokenOwnerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20Approval)
				if err := _TTFT20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TTFT20MintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the TTFT20 contract.
type TTFT20MintIterator struct {
	Event *TTFT20Mint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TTFT20MintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20Mint)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TTFT20Mint)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TTFT20MintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20MintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20Mint represents a Mint event raised by the TTFT20 contract.
type TTFT20Mint struct {
	Receiver common.Address
	Tokens   *big.Int
	Txid     string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a.
//
// Solidity: event Mint(address indexed receiver, uint256 tokens, string txid)
func (_TTFT20 *TTFT20Filterer) FilterMint(opts *bind.FilterOpts, receiver []common.Address) (*TTFT20MintIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Mint", receiverRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20MintIterator{contract: _TTFT20.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a.
//
// Solidity: event Mint(address indexed receiver, uint256 tokens, string txid)
func (_TTFT20 *TTFT20Filterer) WatchMint(opts *bind.WatchOpts, sink chan<- *TTFT20Mint, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Mint", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20Mint)
				if err := _TTFT20.contract.UnpackLog(event, "Mint", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TTFT20OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TTFT20 contract.
type TTFT20OwnershipTransferredIterator struct {
	Event *TTFT20OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TTFT20OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TTFT20OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TTFT20OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20OwnershipTransferred represents a OwnershipTransferred event raised by the TTFT20 contract.
type TTFT20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed _from, address indexed _to)
func (_TTFT20 *TTFT20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*TTFT20OwnershipTransferredIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "OwnershipTransferred", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20OwnershipTransferredIterator{contract: _TTFT20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed _from, address indexed _to)
func (_TTFT20 *TTFT20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TTFT20OwnershipTransferred, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "OwnershipTransferred", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20OwnershipTransferred)
				if err := _TTFT20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TTFT20RegisterWithdrawalAddressIterator is returned from FilterRegisterWithdrawalAddress and is used to iterate over the raw logs and unpacked data for RegisterWithdrawalAddress events raised by the TTFT20 contract.
type TTFT20RegisterWithdrawalAddressIterator struct {
	Event *TTFT20RegisterWithdrawalAddress // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TTFT20RegisterWithdrawalAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20RegisterWithdrawalAddress)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TTFT20RegisterWithdrawalAddress)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TTFT20RegisterWithdrawalAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20RegisterWithdrawalAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20RegisterWithdrawalAddress represents a RegisterWithdrawalAddress event raised by the TTFT20 contract.
type TTFT20RegisterWithdrawalAddress struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRegisterWithdrawalAddress is a free log retrieval operation binding the contract event 0x77bc19082a31daad021d73c26bb4f6e74100a41c98099405e92a9323d133e602.
//
// Solidity: event RegisterWithdrawalAddress(address indexed addr)
func (_TTFT20 *TTFT20Filterer) FilterRegisterWithdrawalAddress(opts *bind.FilterOpts, addr []common.Address) (*TTFT20RegisterWithdrawalAddressIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "RegisterWithdrawalAddress", addrRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20RegisterWithdrawalAddressIterator{contract: _TTFT20.contract, event: "RegisterWithdrawalAddress", logs: logs, sub: sub}, nil
}

// WatchRegisterWithdrawalAddress is a free log subscription operation binding the contract event 0x77bc19082a31daad021d73c26bb4f6e74100a41c98099405e92a9323d133e602.
//
// Solidity: event RegisterWithdrawalAddress(address indexed addr)
func (_TTFT20 *TTFT20Filterer) WatchRegisterWithdrawalAddress(opts *bind.WatchOpts, sink chan<- *TTFT20RegisterWithdrawalAddress, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "RegisterWithdrawalAddress", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20RegisterWithdrawalAddress)
				if err := _TTFT20.contract.UnpackLog(event, "RegisterWithdrawalAddress", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TTFT20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TTFT20 contract.
type TTFT20TransferIterator struct {
	Event *TTFT20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TTFT20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TTFT20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TTFT20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20Transfer represents a Transfer event raised by the TTFT20 contract.
type TTFT20Transfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_TTFT20 *TTFT20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TTFT20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20TransferIterator{contract: _TTFT20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_TTFT20 *TTFT20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TTFT20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20Transfer)
				if err := _TTFT20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TTFT20WithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the TTFT20 contract.
type TTFT20WithdrawIterator struct {
	Event *TTFT20Withdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TTFT20WithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20Withdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TTFT20Withdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TTFT20WithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20WithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20Withdraw represents a Withdraw event raised by the TTFT20 contract.
type TTFT20Withdraw struct {
	Receiver common.Address
	Tokens   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed receiver, uint256 tokens)
func (_TTFT20 *TTFT20Filterer) FilterWithdraw(opts *bind.FilterOpts, receiver []common.Address) (*TTFT20WithdrawIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Withdraw", receiverRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20WithdrawIterator{contract: _TTFT20.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed receiver, uint256 tokens)
func (_TTFT20 *TTFT20Filterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *TTFT20Withdraw, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Withdraw", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20Withdraw)
				if err := _TTFT20.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
