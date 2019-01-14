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

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_toRemove\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"removedOwner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
const OwnedBin = `0x608060405234801561001057600080fd5b5061002333640100000000610028810204565b6100c1565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a0909201909252805191012061009e9060016401000000006100a1810204565b50565b600091825260046020526040909120805460ff1916911515919091179055565b610365806100d06000396000f3fe608060405234801561001057600080fd5b5060043610610052577c01000000000000000000000000000000000000000000000000000000006000350463173825d981146100575780637065cb481461007f575b600080fd5b61007d6004803603602081101561006d57600080fd5b5035600160a060020a03166100a5565b005b61007d6004803603602081101561009557600080fd5b5035600160a060020a031661012c565b6100ae3361019d565b15156100b957600080fd5b600160a060020a03811615156100ce57600080fd5b600160a060020a0381163314156100e457600080fd5b6100ed81610211565b60408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b6101353361019d565b151561014057600080fd5b600160a060020a038116151561015557600080fd5b61015e8161027f565b60408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a0909201909252805191012060009061020b906102ec565b92915050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a0909201909252805191012061027c90610301565b50565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a0909201909252805191012061027c906001610319565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b600091825260046020526040909120805460ff191691151591909117905556fea165627a7a72305820ef0d2835b5783e6eb4a2578fb023ac95e65272b9dabe39f50813bb073d39e7590029`

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

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_Owned *OwnedTransactor) AddOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "addOwner", _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_Owned *OwnedSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.AddOwner(&_Owned.TransactOpts, _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_Owned *OwnedTransactorSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.AddOwner(&_Owned.TransactOpts, _newOwner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_Owned *OwnedTransactor) RemoveOwner(opts *bind.TransactOpts, _toRemove common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "removeOwner", _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_Owned *OwnedSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _Owned.Contract.RemoveOwner(&_Owned.TransactOpts, _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_Owned *OwnedTransactorSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _Owned.Contract.RemoveOwner(&_Owned.TransactOpts, _toRemove)
}

// OwnedAddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the Owned contract.
type OwnedAddedOwnerIterator struct {
	Event *OwnedAddedOwner // Event containing the contract specifics and raw log

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
func (it *OwnedAddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedAddedOwner)
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
		it.Event = new(OwnedAddedOwner)
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
func (it *OwnedAddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedAddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedAddedOwner represents a AddedOwner event raised by the Owned contract.
type OwnedAddedOwner struct {
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_Owned *OwnedFilterer) FilterAddedOwner(opts *bind.FilterOpts) (*OwnedAddedOwnerIterator, error) {

	logs, sub, err := _Owned.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &OwnedAddedOwnerIterator{contract: _Owned.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_Owned *OwnedFilterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *OwnedAddedOwner) (event.Subscription, error) {

	logs, sub, err := _Owned.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedAddedOwner)
				if err := _Owned.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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

// OwnedRemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the Owned contract.
type OwnedRemovedOwnerIterator struct {
	Event *OwnedRemovedOwner // Event containing the contract specifics and raw log

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
func (it *OwnedRemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedRemovedOwner)
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
		it.Event = new(OwnedRemovedOwner)
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
func (it *OwnedRemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedRemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedRemovedOwner represents a RemovedOwner event raised by the Owned contract.
type OwnedRemovedOwner struct {
	RemovedOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_Owned *OwnedFilterer) FilterRemovedOwner(opts *bind.FilterOpts) (*OwnedRemovedOwnerIterator, error) {

	logs, sub, err := _Owned.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &OwnedRemovedOwnerIterator{contract: _Owned.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_Owned *OwnedFilterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *OwnedRemovedOwner) (event.Subscription, error) {

	logs, sub, err := _Owned.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedRemovedOwner)
				if err := _Owned.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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

// OwnedUpgradeableTokenStorageABI is the input ABI used to generate the binding from.
const OwnedUpgradeableTokenStorageABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_toRemove\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_version\",\"type\":\"string\"},{\"name\":\"_implementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"version\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"removedOwner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"}]"

// OwnedUpgradeableTokenStorageBin is the compiled bytecode used for deploying new contracts.
const OwnedUpgradeableTokenStorageBin = `0x60c0604052600560809081527f544654323000000000000000000000000000000000000000000000000000000060a0526200004390640100000000620000e5810204565b60408051808201909152601881527f54465420455243323020726570726573656e746174696f6e000000000000000060208201526200008b906401000000006200015a810204565b6009620000a181640100000000620001cc810204565b64174876e80060ff8216600a0a02620000c38164010000000062000241810204565b5050620000df33620002b3640100000000026401000000009004565b62000428565b620001576040516020018080602001828103825260068152602001807f73796d626f6c000000000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826200032b640100000000026401000000009004565b50565b620001576040516020018080602001828103825260048152602001807f6e616d650000000000000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826200032b640100000000026401000000009004565b620001576040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208260ff1662000351640100000000026401000000009004565b6200015760405160200180806020018281038252600b8152602001807f746f74616c537570706c79000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208262000351640100000000026401000000009004565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206200015790600164010000000062000363810204565b600082815260016020908152604090912082516200034c9284019062000383565b505050565b60009182526020829052604090912055565b600091825260046020526040909120805460ff1916911515919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620003c657805160ff1916838001178555620003f6565b82800160010185558215620003f6579182015b82811115620003f6578251825591602001919060010190620003d9565b506200040492915062000408565b5090565b6200042591905b808211156200040457600081556001016200040f565b90565b61090080620004386000396000f3fe608060405234801561001057600080fd5b5060043610610073577c01000000000000000000000000000000000000000000000000000000006000350463173825d9811461007857806354fd4d50146100a05780635a8b1a9f1461011d5780635c60da1b146101ce5780637065cb48146101f2575b600080fd5b61009e6004803603602081101561008e57600080fd5b5035600160a060020a0316610218565b005b6100a861029f565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e25781810151838201526020016100ca565b50505050905090810190601f16801561010f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61009e6004803603604081101561013357600080fd5b81019060208101813564010000000081111561014e57600080fd5b82018360208201111561016057600080fd5b8035906020019184600183028401116401000000008311171561018257600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050509035600160a060020a031691506102af9050565b6101d661038c565b60408051600160a060020a039092168252519081900360200190f35b61009e6004803603602081101561020857600080fd5b5035600160a060020a0316610396565b61022133610407565b151561022c57600080fd5b600160a060020a038116151561024157600080fd5b600160a060020a03811633141561025757600080fd5b6102608161047b565b60408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b60606102a96104e9565b90505b90565b6102b833610407565b15156102c357600080fd5b80600160a060020a03166102d5610544565b600160a060020a031614156102e957600080fd5b6102f2826105a6565b6102fb81610607565b80600160a060020a0316826040518082805190602001908083835b602083106103355780518252601f199092019160209182019101610316565b5181516020939093036101000a60001901801990911692169190911790526040519201829003822093507f8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e92506000919050a35050565b60006102a9610544565b61039f33610407565b15156103aa57600080fd5b600160a060020a03811615156103bf57600080fd5b6103c881610668565b60408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120600090610475906106d5565b92915050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206104e6906106ea565b50565b6040805160208082018190526007828401527f76657273696f6e00000000000000000000000000000000000000000000000000606083810191909152835180840382018152608090930190935281519101206102a990610702565b60006102a960405160200180806020018281038252600e8152602001807f696d706c656d656e746174696f6e000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001206107a2565b6104e66040516020018080602001828103825260078152602001807f76657273696f6e0000000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826107bd565b6104e660405160200180806020018281038252600e8152602001807f696d706c656d656e746174696f6e00000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826107e1565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206104e690600161081c565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b60008181526001602081815260409283902080548451600294821615610100026000190190911693909304601f810183900483028401830190945283835260609390918301828280156107965780601f1061076b57610100808354040283529160200191610796565b820191906000526020600020905b81548152906001019060200180831161077957829003601f168201915b50505050509050919050565b600090815260026020526040902054600160a060020a031690565b600082815260016020908152604090912082516107dc9284019061083c565b505050565b600091825260026020526040909120805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909216919091179055565b600091825260046020526040909120805460ff1916911515919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061087d57805160ff19168380011785556108aa565b828001600101855582156108aa579182015b828111156108aa57825182559160200191906001019061088f565b506108b69291506108ba565b5090565b6102ac91905b808211156108b657600081556001016108c056fea165627a7a72305820d49f49f1454a626ac321a5346d91e30e6f73a2477df32bbcec6c04a89f16cbef0029`

// DeployOwnedUpgradeableTokenStorage deploys a new Ethereum contract, binding an instance of OwnedUpgradeableTokenStorage to it.
func DeployOwnedUpgradeableTokenStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OwnedUpgradeableTokenStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedUpgradeableTokenStorageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedUpgradeableTokenStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnedUpgradeableTokenStorage{OwnedUpgradeableTokenStorageCaller: OwnedUpgradeableTokenStorageCaller{contract: contract}, OwnedUpgradeableTokenStorageTransactor: OwnedUpgradeableTokenStorageTransactor{contract: contract}, OwnedUpgradeableTokenStorageFilterer: OwnedUpgradeableTokenStorageFilterer{contract: contract}}, nil
}

// OwnedUpgradeableTokenStorage is an auto generated Go binding around an Ethereum contract.
type OwnedUpgradeableTokenStorage struct {
	OwnedUpgradeableTokenStorageCaller     // Read-only binding to the contract
	OwnedUpgradeableTokenStorageTransactor // Write-only binding to the contract
	OwnedUpgradeableTokenStorageFilterer   // Log filterer for contract events
}

// OwnedUpgradeableTokenStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedUpgradeableTokenStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedUpgradeableTokenStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedUpgradeableTokenStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedUpgradeableTokenStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedUpgradeableTokenStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedUpgradeableTokenStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedUpgradeableTokenStorageSession struct {
	Contract     *OwnedUpgradeableTokenStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OwnedUpgradeableTokenStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedUpgradeableTokenStorageCallerSession struct {
	Contract *OwnedUpgradeableTokenStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// OwnedUpgradeableTokenStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedUpgradeableTokenStorageTransactorSession struct {
	Contract     *OwnedUpgradeableTokenStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// OwnedUpgradeableTokenStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedUpgradeableTokenStorageRaw struct {
	Contract *OwnedUpgradeableTokenStorage // Generic contract binding to access the raw methods on
}

// OwnedUpgradeableTokenStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedUpgradeableTokenStorageCallerRaw struct {
	Contract *OwnedUpgradeableTokenStorageCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedUpgradeableTokenStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedUpgradeableTokenStorageTransactorRaw struct {
	Contract *OwnedUpgradeableTokenStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnedUpgradeableTokenStorage creates a new instance of OwnedUpgradeableTokenStorage, bound to a specific deployed contract.
func NewOwnedUpgradeableTokenStorage(address common.Address, backend bind.ContractBackend) (*OwnedUpgradeableTokenStorage, error) {
	contract, err := bindOwnedUpgradeableTokenStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorage{OwnedUpgradeableTokenStorageCaller: OwnedUpgradeableTokenStorageCaller{contract: contract}, OwnedUpgradeableTokenStorageTransactor: OwnedUpgradeableTokenStorageTransactor{contract: contract}, OwnedUpgradeableTokenStorageFilterer: OwnedUpgradeableTokenStorageFilterer{contract: contract}}, nil
}

// NewOwnedUpgradeableTokenStorageCaller creates a new read-only instance of OwnedUpgradeableTokenStorage, bound to a specific deployed contract.
func NewOwnedUpgradeableTokenStorageCaller(address common.Address, caller bind.ContractCaller) (*OwnedUpgradeableTokenStorageCaller, error) {
	contract, err := bindOwnedUpgradeableTokenStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorageCaller{contract: contract}, nil
}

// NewOwnedUpgradeableTokenStorageTransactor creates a new write-only instance of OwnedUpgradeableTokenStorage, bound to a specific deployed contract.
func NewOwnedUpgradeableTokenStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedUpgradeableTokenStorageTransactor, error) {
	contract, err := bindOwnedUpgradeableTokenStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorageTransactor{contract: contract}, nil
}

// NewOwnedUpgradeableTokenStorageFilterer creates a new log filterer instance of OwnedUpgradeableTokenStorage, bound to a specific deployed contract.
func NewOwnedUpgradeableTokenStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedUpgradeableTokenStorageFilterer, error) {
	contract, err := bindOwnedUpgradeableTokenStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorageFilterer{contract: contract}, nil
}

// bindOwnedUpgradeableTokenStorage binds a generic wrapper to an already deployed contract.
func bindOwnedUpgradeableTokenStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedUpgradeableTokenStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OwnedUpgradeableTokenStorage.Contract.OwnedUpgradeableTokenStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.OwnedUpgradeableTokenStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.OwnedUpgradeableTokenStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OwnedUpgradeableTokenStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OwnedUpgradeableTokenStorage.contract.Call(opts, out, "implementation")
	return *ret0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageSession) Implementation() (common.Address, error) {
	return _OwnedUpgradeableTokenStorage.Contract.Implementation(&_OwnedUpgradeableTokenStorage.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageCallerSession) Implementation() (common.Address, error) {
	return _OwnedUpgradeableTokenStorage.Contract.Implementation(&_OwnedUpgradeableTokenStorage.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _OwnedUpgradeableTokenStorage.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageSession) Version() (string, error) {
	return _OwnedUpgradeableTokenStorage.Contract.Version(&_OwnedUpgradeableTokenStorage.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageCallerSession) Version() (string, error) {
	return _OwnedUpgradeableTokenStorage.Contract.Version(&_OwnedUpgradeableTokenStorage.CallOpts)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactor) AddOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.contract.Transact(opts, "addOwner", _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.AddOwner(&_OwnedUpgradeableTokenStorage.TransactOpts, _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactorSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.AddOwner(&_OwnedUpgradeableTokenStorage.TransactOpts, _newOwner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactor) RemoveOwner(opts *bind.TransactOpts, _toRemove common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.contract.Transact(opts, "removeOwner", _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.RemoveOwner(&_OwnedUpgradeableTokenStorage.TransactOpts, _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactorSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.RemoveOwner(&_OwnedUpgradeableTokenStorage.TransactOpts, _toRemove)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactor) UpgradeTo(opts *bind.TransactOpts, _version string, _implementation common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.contract.Transact(opts, "upgradeTo", _version, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageSession) UpgradeTo(_version string, _implementation common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.UpgradeTo(&_OwnedUpgradeableTokenStorage.TransactOpts, _version, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageTransactorSession) UpgradeTo(_version string, _implementation common.Address) (*types.Transaction, error) {
	return _OwnedUpgradeableTokenStorage.Contract.UpgradeTo(&_OwnedUpgradeableTokenStorage.TransactOpts, _version, _implementation)
}

// OwnedUpgradeableTokenStorageAddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the OwnedUpgradeableTokenStorage contract.
type OwnedUpgradeableTokenStorageAddedOwnerIterator struct {
	Event *OwnedUpgradeableTokenStorageAddedOwner // Event containing the contract specifics and raw log

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
func (it *OwnedUpgradeableTokenStorageAddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedUpgradeableTokenStorageAddedOwner)
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
		it.Event = new(OwnedUpgradeableTokenStorageAddedOwner)
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
func (it *OwnedUpgradeableTokenStorageAddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedUpgradeableTokenStorageAddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedUpgradeableTokenStorageAddedOwner represents a AddedOwner event raised by the OwnedUpgradeableTokenStorage contract.
type OwnedUpgradeableTokenStorageAddedOwner struct {
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageFilterer) FilterAddedOwner(opts *bind.FilterOpts) (*OwnedUpgradeableTokenStorageAddedOwnerIterator, error) {

	logs, sub, err := _OwnedUpgradeableTokenStorage.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorageAddedOwnerIterator{contract: _OwnedUpgradeableTokenStorage.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageFilterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *OwnedUpgradeableTokenStorageAddedOwner) (event.Subscription, error) {

	logs, sub, err := _OwnedUpgradeableTokenStorage.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedUpgradeableTokenStorageAddedOwner)
				if err := _OwnedUpgradeableTokenStorage.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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

// OwnedUpgradeableTokenStorageRemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the OwnedUpgradeableTokenStorage contract.
type OwnedUpgradeableTokenStorageRemovedOwnerIterator struct {
	Event *OwnedUpgradeableTokenStorageRemovedOwner // Event containing the contract specifics and raw log

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
func (it *OwnedUpgradeableTokenStorageRemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedUpgradeableTokenStorageRemovedOwner)
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
		it.Event = new(OwnedUpgradeableTokenStorageRemovedOwner)
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
func (it *OwnedUpgradeableTokenStorageRemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedUpgradeableTokenStorageRemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedUpgradeableTokenStorageRemovedOwner represents a RemovedOwner event raised by the OwnedUpgradeableTokenStorage contract.
type OwnedUpgradeableTokenStorageRemovedOwner struct {
	RemovedOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageFilterer) FilterRemovedOwner(opts *bind.FilterOpts) (*OwnedUpgradeableTokenStorageRemovedOwnerIterator, error) {

	logs, sub, err := _OwnedUpgradeableTokenStorage.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorageRemovedOwnerIterator{contract: _OwnedUpgradeableTokenStorage.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageFilterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *OwnedUpgradeableTokenStorageRemovedOwner) (event.Subscription, error) {

	logs, sub, err := _OwnedUpgradeableTokenStorage.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedUpgradeableTokenStorageRemovedOwner)
				if err := _OwnedUpgradeableTokenStorage.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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

// OwnedUpgradeableTokenStorageUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the OwnedUpgradeableTokenStorage contract.
type OwnedUpgradeableTokenStorageUpgradedIterator struct {
	Event *OwnedUpgradeableTokenStorageUpgraded // Event containing the contract specifics and raw log

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
func (it *OwnedUpgradeableTokenStorageUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedUpgradeableTokenStorageUpgraded)
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
		it.Event = new(OwnedUpgradeableTokenStorageUpgraded)
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
func (it *OwnedUpgradeableTokenStorageUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedUpgradeableTokenStorageUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedUpgradeableTokenStorageUpgraded represents a Upgraded event raised by the OwnedUpgradeableTokenStorage contract.
type OwnedUpgradeableTokenStorageUpgraded struct {
	Version        common.Hash
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0x8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e.
//
// Solidity: e Upgraded(version indexed string, implementation indexed address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageFilterer) FilterUpgraded(opts *bind.FilterOpts, version []string, implementation []common.Address) (*OwnedUpgradeableTokenStorageUpgradedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _OwnedUpgradeableTokenStorage.contract.FilterLogs(opts, "Upgraded", versionRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &OwnedUpgradeableTokenStorageUpgradedIterator{contract: _OwnedUpgradeableTokenStorage.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0x8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e.
//
// Solidity: e Upgraded(version indexed string, implementation indexed address)
func (_OwnedUpgradeableTokenStorage *OwnedUpgradeableTokenStorageFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *OwnedUpgradeableTokenStorageUpgraded, version []string, implementation []common.Address) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _OwnedUpgradeableTokenStorage.contract.WatchLogs(opts, "Upgraded", versionRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedUpgradeableTokenStorageUpgraded)
				if err := _OwnedUpgradeableTokenStorage.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a72305820b16291a69ac916223f67fbb2dd912316c35f39ac46fe66d50c2ff58132a9df280029`

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

// StorageABI is the input ABI used to generate the binding from.
const StorageABI = "[]"

// StorageBin is the compiled bytecode used for deploying new contracts.
const StorageBin = `0x6080604052348015600f57600080fd5b50603580601d6000396000f3fe6080604052600080fdfea165627a7a7230582004cf9335d2bd9f90093462cdfce184c24cfa24c88a19ffa36c93d992239f2e260029`

// DeployStorage deploys a new Ethereum contract, binding an instance of Storage to it.
func DeployStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Storage, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Storage{StorageCaller: StorageCaller{contract: contract}, StorageTransactor: StorageTransactor{contract: contract}, StorageFilterer: StorageFilterer{contract: contract}}, nil
}

// Storage is an auto generated Go binding around an Ethereum contract.
type Storage struct {
	StorageCaller     // Read-only binding to the contract
	StorageTransactor // Write-only binding to the contract
	StorageFilterer   // Log filterer for contract events
}

// StorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageSession struct {
	Contract     *Storage          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageCallerSession struct {
	Contract *StorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageTransactorSession struct {
	Contract     *StorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageRaw struct {
	Contract *Storage // Generic contract binding to access the raw methods on
}

// StorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageCallerRaw struct {
	Contract *StorageCaller // Generic read-only contract binding to access the raw methods on
}

// StorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageTransactorRaw struct {
	Contract *StorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorage creates a new instance of Storage, bound to a specific deployed contract.
func NewStorage(address common.Address, backend bind.ContractBackend) (*Storage, error) {
	contract, err := bindStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Storage{StorageCaller: StorageCaller{contract: contract}, StorageTransactor: StorageTransactor{contract: contract}, StorageFilterer: StorageFilterer{contract: contract}}, nil
}

// NewStorageCaller creates a new read-only instance of Storage, bound to a specific deployed contract.
func NewStorageCaller(address common.Address, caller bind.ContractCaller) (*StorageCaller, error) {
	contract, err := bindStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageCaller{contract: contract}, nil
}

// NewStorageTransactor creates a new write-only instance of Storage, bound to a specific deployed contract.
func NewStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageTransactor, error) {
	contract, err := bindStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageTransactor{contract: contract}, nil
}

// NewStorageFilterer creates a new log filterer instance of Storage, bound to a specific deployed contract.
func NewStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageFilterer, error) {
	contract, err := bindStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageFilterer{contract: contract}, nil
}

// bindStorage binds a generic wrapper to an already deployed contract.
func bindStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Storage *StorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Storage.Contract.StorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Storage *StorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Storage.Contract.StorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Storage *StorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Storage.Contract.StorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Storage *StorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Storage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Storage *StorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Storage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Storage *StorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Storage.Contract.contract.Transact(opts, method, params...)
}

// TTFT20ABI is the input ABI used to generate the binding from.
const TTFT20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_toRemove\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isWithdrawalAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"registerWithdrawalAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_version\",\"type\":\"string\"},{\"name\":\"_implementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_txid\",\"type\":\"string\"}],\"name\":\"isMintID\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"},{\"name\":\"txid\",\"type\":\"string\"}],\"name\":\"mintTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"RegisterWithdrawalAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"txid\",\"type\":\"string\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"version\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"removedOwner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"}]"

// TTFT20Bin is the compiled bytecode used for deploying new contracts.
const TTFT20Bin = `0x60c0604052600560809081527f544654323000000000000000000000000000000000000000000000000000000060a0526200004390640100000000620000e5810204565b60408051808201909152601881527f54465420455243323020726570726573656e746174696f6e000000000000000060208201526200008b906401000000006200015a810204565b6009620000a181640100000000620001cc810204565b64174876e80060ff8216600a0a02620000c38164010000000062000241810204565b5050620000df33620002b3640100000000026401000000009004565b62000428565b620001576040516020018080602001828103825260068152602001807f73796d626f6c000000000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826200032b640100000000026401000000009004565b50565b620001576040516020018080602001828103825260048152602001807f6e616d650000000000000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826200032b640100000000026401000000009004565b620001576040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208260ff1662000351640100000000026401000000009004565b6200015760405160200180806020018281038252600b8152602001807f746f74616c537570706c79000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208262000351640100000000026401000000009004565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206200015790600164010000000062000363810204565b600082815260016020908152604090912082516200034c9284019062000383565b505050565b60009182526020829052604090912055565b600091825260046020526040909120805460ff1916911515919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620003c657805160ff1916838001178555620003f6565b82800160010185558215620003f6579182015b82811115620003f6578251825591602001919060010190620003d9565b506200040492915062000408565b5090565b6200042591905b808211156200040457600081556001016200040f565b90565b611a4880620004386000396000f3fe608060405260043610610126576000357c0100000000000000000000000000000000000000000000000000000000900480635a8b1a9f116100b257806395d89b411161008157806395d89b411461049c578063a9059cbb146104b1578063dd62ed3e146104ea578063dd6ad77e14610525578063e67524a3146105d857610126565b80635a8b1a9f146103475780635c60da1b146104055780637065cb481461043657806370a082311461046957610126565b8063206a943c116100f9578063206a943c1461025e57806323b872dd14610291578063313ce567146102d457806334ca6a71146102ff57806354fd4d501461033257610126565b806306fdde031461012b578063095ea7b3146101b5578063173825d91461020257806318160ddd14610237575b600080fd5b34801561013757600080fd5b506101406106a0565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561017a578181015183820152602001610162565b50505050905090810190601f1680156101a75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101c157600080fd5b506101ee600480360360408110156101d857600080fd5b50600160a060020a0381351690602001356106b0565b604080519115158252519081900360200190f35b34801561020e57600080fd5b506102356004803603602081101561022557600080fd5b5035600160a060020a0316610707565b005b34801561024357600080fd5b5061024c61078e565b60408051918252519081900360200190f35b34801561026a57600080fd5b506101ee6004803603602081101561028157600080fd5b5035600160a060020a03166107b1565b34801561029d57600080fd5b506101ee600480360360608110156102b457600080fd5b50600160a060020a038135811691602081013590911690604001356107bc565b3480156102e057600080fd5b506102e96108ce565b6040805160ff9092168252519081900360200190f35b34801561030b57600080fd5b506102356004803603602081101561032257600080fd5b5035600160a060020a03166108d8565b34801561033e57600080fd5b506101406109f5565b34801561035357600080fd5b506102356004803603604081101561036a57600080fd5b81019060208101813564010000000081111561038557600080fd5b82018360208201111561039757600080fd5b803590602001918460018302840111640100000000831117156103b957600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050509035600160a060020a031691506109ff9050565b34801561041157600080fd5b5061041a610adc565b60408051600160a060020a039092168252519081900360200190f35b34801561044257600080fd5b506102356004803603602081101561045957600080fd5b5035600160a060020a0316610ae6565b34801561047557600080fd5b5061024c6004803603602081101561048c57600080fd5b5035600160a060020a0316610b57565b3480156104a857600080fd5b50610140610b62565b3480156104bd57600080fd5b506101ee600480360360408110156104d457600080fd5b50600160a060020a038135169060200135610b6c565b3480156104f657600080fd5b5061024c6004803603604081101561050d57600080fd5b50600160a060020a0381358116916020013516610c3e565b34801561053157600080fd5b506101ee6004803603602081101561054857600080fd5b81019060208101813564010000000081111561056357600080fd5b82018360208201111561057557600080fd5b8035906020019184600183028401116401000000008311171561059757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610c51945050505050565b3480156105e457600080fd5b50610235600480360360608110156105fb57600080fd5b600160a060020a038235169160208101359181019060608101604082013564010000000081111561062b57600080fd5b82018360208201111561063d57600080fd5b8035906020019184600183028401116401000000008311171561065f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610c5c945050505050565b60606106aa610d99565b90505b90565b60006106bd338484610df4565b604080518381529051600160a060020a0385169133917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259181900360200190a35060015b92915050565b61071033610e72565b151561071b57600080fd5b600160a060020a038116151561073057600080fd5b600160a060020a03811633141561074657600080fd5b61074f81610ee0565b60408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b60006106aa61079d6000610f4e565b6107a5610fbc565b9063ffffffff61101e16565b600061070182611033565b60006107d684336107d1856107a589336110e8565b610df4565b6107ec846107e7846107a588610f4e565b611163565b6107f583611033565b8015610806575061080683836111d3565b1561085b5782600160a060020a031684600160a060020a03167f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb846040518082815260200191505060405180910390a36108c4565b610878836107e78461086c87610f4e565b9063ffffffff6111fb16565b82600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a35b5060019392505050565b60006106aa61120b565b6108e133610e72565b15156108ec57600080fd5b6108f581611033565b1561094b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806119f86025913960400191505060405180910390fd5b6109548161126d565b600061095f82610f4e565b905061096c8260006111d3565b156109bd5761097c826000611163565b604080518281529051600160a060020a0384169182917f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9181900360200190a35b604051600160a060020a038316907f77bc19082a31daad021d73c26bb4f6e74100a41c98099405e92a9323d133e60290600090a25050565b60606106aa611316565b610a0833610e72565b1515610a1357600080fd5b80600160a060020a0316610a25611371565b600160a060020a03161415610a3957600080fd5b610a42826113d3565b610a4b81611434565b80600160a060020a0316826040518082805190602001908083835b60208310610a855780518252601f199092019160209182019101610a66565b5181516020939093036101000a60001901801990911692169190911790526040519201829003822093507f8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e92506000919050a35050565b60006106aa611371565b610aef33610e72565b1515610afa57600080fd5b600160a060020a0381161515610b0f57600080fd5b610b1881611495565b60408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b600061070182610f4e565b60606106aa611502565b6000610b7f336107e7846107a533610f4e565b610b8883611033565b8015610b995750610b9983836111d3565b15610be357604080518381529051600160a060020a0385169133917f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9181900360200190a3610c35565b610bf4836107e78461086c87610f4e565b604080518381529051600160a060020a0385169133917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a35b50600192915050565b6000610c4a83836110e8565b9392505050565b60006107018261155d565b610c6533610e72565b1515610c7057600080fd5b610c798161155d565b15610ce557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f544654207472616e736163746f6e20494420616c7265616479206b6e6f776e00604482015290519081900360640190fd5b610cee8161169b565b610cff836107e78461086c87610f4e565b806040518082805190602001908083835b60208310610d2f5780518252601f199092019160209182019101610d10565b51815160209384036101000a6000190180199092169116179052604080519290940182900382208883529351939550600160a060020a03891694507f85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a9391829003019150a3505050565b6040805160208082018190526004828401527f6e616d6500000000000000000000000000000000000000000000000000000000606083810191909152835180840382018152608090930190935281519101206106aa906117d9565b60408051600160a060020a03808616828401528416606080830191909152602080830191909152600760808301527f616c6c6f7765640000000000000000000000000000000000000000000000000060a0808401919091528351808403909101815260c09092019092528051910120610e6d9082611879565b505050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906107019061188b565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120610f4b906118a0565b50565b60408051600160a060020a038316818301526020808201839052600760608301527f62616c616e6365000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120600090610701906118b8565b60006106aa60405160200180806020018281038252600b8152602001807f746f74616c537570706c79000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001206118b8565b60008282111561102d57600080fd5b50900390565b60006107018260405160200180806020018060200184600160a060020a0316600160a060020a03168152602001838103835260078152602001807f61646472657373000000000000000000000000000000000000000000000000008152506020018381038252600a8152602001807f7769746864726177616c0000000000000000000000000000000000000000000081525060200193505050506040516020818303038152906040528051906020012061188b565b60408051600160a060020a03808516828401528316606080830191909152602080830191909152600760808301527f616c6c6f7765640000000000000000000000000000000000000000000000000060a0808401919091528351808403909101815260c09092019092528051910120600090610c4a906118b8565b60408051600160a060020a038416818301526020808201839052600760608301527f62616c616e6365000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206111cf9082611879565b5050565b600060016111df61120b565b0360ff16600a0a6111f38361086c86610f4e565b119392505050565b8181018281101561070157600080fd5b60006106aa6040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001206118b8565b60408051600160a060020a038316606080830191909152602080830191909152600760808301527f616464726573730000000000000000000000000000000000000000000000000060a08084019190915282840152600a60c08301527f7769746864726177616c0000000000000000000000000000000000000000000060e080840191909152835180840390910181526101009092019092528051910120610f4b9060016118ca565b6040805160208082018190526007828401527f76657273696f6e00000000000000000000000000000000000000000000000000606083810191909152835180840382018152608090930190935281519101206106aa906117d9565b60006106aa60405160200180806020018281038252600e8152602001807f696d706c656d656e746174696f6e000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001206118ea565b610f4b6040516020018080602001828103825260078152602001807f76657273696f6e000000000000000000000000000000000000000000000000008152506020019150506040516020818303038152906040528051906020012082611905565b610f4b60405160200180806020018281038252600e8152602001807f696d706c656d656e746174696f6e0000000000000000000000000000000000008152506020019150506040516020818303038152906040528051906020012082611924565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120610f4b9060016118ca565b6040805160208082018190526006828401527f73796d626f6c0000000000000000000000000000000000000000000000000000606083810191909152835180840382018152608090930190935281519101206106aa906117d9565b6000610701826040516020018080602001806020018060200180602001858103855260048152602001807f6d696e74000000000000000000000000000000000000000000000000000000008152506020018581038452600b8152602001807f7472616e73616374696f6e000000000000000000000000000000000000000000815250602001858103835260028152602001807f6964000000000000000000000000000000000000000000000000000000000000815250602001858103825286818151815260200191508051906020019080838360005b8381101561164b578181015183820152602001611633565b50505050905090810190601f1680156116785780820380516001836020036101000a031916815260200191505b50955050505050506040516020818303038152906040528051906020012061188b565b610f4b816040516020018080602001806020018060200180602001858103855260048152602001807f6d696e74000000000000000000000000000000000000000000000000000000008152506020018581038452600b8152602001807f7472616e73616374696f6e000000000000000000000000000000000000000000815250602001858103835260028152602001807f6964000000000000000000000000000000000000000000000000000000000000815250602001858103825286818151815260200191508051906020019080838360005b8381101561178757818101518382015260200161176f565b50505050905090810190601f1680156117b45780820380516001836020036101000a031916815260200191505b50955050505050506040516020818303038152906040528051906020012060016118ca565b60008181526001602081815260409283902080548451600294821615610100026000190190911693909304601f8101839004830284018301909452838352606093909183018282801561186d5780601f106118425761010080835404028352916020019161186d565b820191906000526020600020905b81548152906001019060200180831161185057829003601f168201915b50505050509050919050565b60009182526020829052604090912055565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b60009081526020819052604090205490565b600091825260046020526040909120805460ff1916911515919091179055565b600090815260026020526040902054600160a060020a031690565b60008281526001602090815260409091208251610e6d9284019061195f565b600091825260026020526040909120805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909216919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106119a057805160ff19168380011785556119cd565b828001600101855582156119cd579182015b828111156119cd5782518255916020019190600101906119b2565b506119d99291506119dd565b5090565b6106ad91905b808211156119d957600081556001016119e356fe5769746864726177616c206164647265737320616c72656164792072656769737465726564a165627a7a72305820bdcd7b1510c1e88ee05f2e414243652d109e8b7ad7c747660918999c907d67c50029`

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
// Solidity: function allowance(tokenOwner address, spender address) constant returns(remaining uint256)
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
// Solidity: function allowance(tokenOwner address, spender address) constant returns(remaining uint256)
func (_TTFT20 *TTFT20Session) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _TTFT20.Contract.Allowance(&_TTFT20.CallOpts, tokenOwner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(tokenOwner address, spender address) constant returns(remaining uint256)
func (_TTFT20 *TTFT20CallerSession) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _TTFT20.Contract.Allowance(&_TTFT20.CallOpts, tokenOwner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(tokenOwner address) constant returns(balance uint256)
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
// Solidity: function balanceOf(tokenOwner address) constant returns(balance uint256)
func (_TTFT20 *TTFT20Session) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _TTFT20.Contract.BalanceOf(&_TTFT20.CallOpts, tokenOwner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(tokenOwner address) constant returns(balance uint256)
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

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_TTFT20 *TTFT20Caller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "implementation")
	return *ret0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_TTFT20 *TTFT20Session) Implementation() (common.Address, error) {
	return _TTFT20.Contract.Implementation(&_TTFT20.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_TTFT20 *TTFT20CallerSession) Implementation() (common.Address, error) {
	return _TTFT20.Contract.Implementation(&_TTFT20.CallOpts)
}

// IsMintID is a free data retrieval call binding the contract method 0xdd6ad77e.
//
// Solidity: function isMintID(_txid string) constant returns(bool)
func (_TTFT20 *TTFT20Caller) IsMintID(opts *bind.CallOpts, _txid string) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "isMintID", _txid)
	return *ret0, err
}

// IsMintID is a free data retrieval call binding the contract method 0xdd6ad77e.
//
// Solidity: function isMintID(_txid string) constant returns(bool)
func (_TTFT20 *TTFT20Session) IsMintID(_txid string) (bool, error) {
	return _TTFT20.Contract.IsMintID(&_TTFT20.CallOpts, _txid)
}

// IsMintID is a free data retrieval call binding the contract method 0xdd6ad77e.
//
// Solidity: function isMintID(_txid string) constant returns(bool)
func (_TTFT20 *TTFT20CallerSession) IsMintID(_txid string) (bool, error) {
	return _TTFT20.Contract.IsMintID(&_TTFT20.CallOpts, _txid)
}

// IsWithdrawalAddress is a free data retrieval call binding the contract method 0x206a943c.
//
// Solidity: function isWithdrawalAddress(_addr address) constant returns(bool)
func (_TTFT20 *TTFT20Caller) IsWithdrawalAddress(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "isWithdrawalAddress", _addr)
	return *ret0, err
}

// IsWithdrawalAddress is a free data retrieval call binding the contract method 0x206a943c.
//
// Solidity: function isWithdrawalAddress(_addr address) constant returns(bool)
func (_TTFT20 *TTFT20Session) IsWithdrawalAddress(_addr common.Address) (bool, error) {
	return _TTFT20.Contract.IsWithdrawalAddress(&_TTFT20.CallOpts, _addr)
}

// IsWithdrawalAddress is a free data retrieval call binding the contract method 0x206a943c.
//
// Solidity: function isWithdrawalAddress(_addr address) constant returns(bool)
func (_TTFT20 *TTFT20CallerSession) IsWithdrawalAddress(_addr common.Address) (bool, error) {
	return _TTFT20.Contract.IsWithdrawalAddress(&_TTFT20.CallOpts, _addr)
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

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_TTFT20 *TTFT20Caller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TTFT20.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_TTFT20 *TTFT20Session) Version() (string, error) {
	return _TTFT20.Contract.Version(&_TTFT20.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_TTFT20 *TTFT20CallerSession) Version() (string, error) {
	return _TTFT20.Contract.Version(&_TTFT20.CallOpts)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_TTFT20 *TTFT20Transactor) AddOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "addOwner", _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_TTFT20 *TTFT20Session) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.AddOwner(&_TTFT20.TransactOpts, _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_TTFT20 *TTFT20TransactorSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.AddOwner(&_TTFT20.TransactOpts, _newOwner)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "approve", spender, tokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20Session) Approve(spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Approve(&_TTFT20.TransactOpts, spender, tokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20TransactorSession) Approve(spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Approve(&_TTFT20.TransactOpts, spender, tokens)
}

// MintTokens is a paid mutator transaction binding the contract method 0xe67524a3.
//
// Solidity: function mintTokens(receiver address, tokens uint256, txid string) returns()
func (_TTFT20 *TTFT20Transactor) MintTokens(opts *bind.TransactOpts, receiver common.Address, tokens *big.Int, txid string) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "mintTokens", receiver, tokens, txid)
}

// MintTokens is a paid mutator transaction binding the contract method 0xe67524a3.
//
// Solidity: function mintTokens(receiver address, tokens uint256, txid string) returns()
func (_TTFT20 *TTFT20Session) MintTokens(receiver common.Address, tokens *big.Int, txid string) (*types.Transaction, error) {
	return _TTFT20.Contract.MintTokens(&_TTFT20.TransactOpts, receiver, tokens, txid)
}

// MintTokens is a paid mutator transaction binding the contract method 0xe67524a3.
//
// Solidity: function mintTokens(receiver address, tokens uint256, txid string) returns()
func (_TTFT20 *TTFT20TransactorSession) MintTokens(receiver common.Address, tokens *big.Int, txid string) (*types.Transaction, error) {
	return _TTFT20.Contract.MintTokens(&_TTFT20.TransactOpts, receiver, tokens, txid)
}

// RegisterWithdrawalAddress is a paid mutator transaction binding the contract method 0x34ca6a71.
//
// Solidity: function registerWithdrawalAddress(addr address) returns()
func (_TTFT20 *TTFT20Transactor) RegisterWithdrawalAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "registerWithdrawalAddress", addr)
}

// RegisterWithdrawalAddress is a paid mutator transaction binding the contract method 0x34ca6a71.
//
// Solidity: function registerWithdrawalAddress(addr address) returns()
func (_TTFT20 *TTFT20Session) RegisterWithdrawalAddress(addr common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.RegisterWithdrawalAddress(&_TTFT20.TransactOpts, addr)
}

// RegisterWithdrawalAddress is a paid mutator transaction binding the contract method 0x34ca6a71.
//
// Solidity: function registerWithdrawalAddress(addr address) returns()
func (_TTFT20 *TTFT20TransactorSession) RegisterWithdrawalAddress(addr common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.RegisterWithdrawalAddress(&_TTFT20.TransactOpts, addr)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_TTFT20 *TTFT20Transactor) RemoveOwner(opts *bind.TransactOpts, _toRemove common.Address) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "removeOwner", _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_TTFT20 *TTFT20Session) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.RemoveOwner(&_TTFT20.TransactOpts, _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_TTFT20 *TTFT20TransactorSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.RemoveOwner(&_TTFT20.TransactOpts, _toRemove)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "transfer", to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20Session) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Transfer(&_TTFT20.TransactOpts, to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20TransactorSession) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.Transfer(&_TTFT20.TransactOpts, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "transferFrom", from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20Session) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferFrom(&_TTFT20.TransactOpts, from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, tokens uint256) returns(success bool)
func (_TTFT20 *TTFT20TransactorSession) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _TTFT20.Contract.TransferFrom(&_TTFT20.TransactOpts, from, to, tokens)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_TTFT20 *TTFT20Transactor) UpgradeTo(opts *bind.TransactOpts, _version string, _implementation common.Address) (*types.Transaction, error) {
	return _TTFT20.contract.Transact(opts, "upgradeTo", _version, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_TTFT20 *TTFT20Session) UpgradeTo(_version string, _implementation common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.UpgradeTo(&_TTFT20.TransactOpts, _version, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_TTFT20 *TTFT20TransactorSession) UpgradeTo(_version string, _implementation common.Address) (*types.Transaction, error) {
	return _TTFT20.Contract.UpgradeTo(&_TTFT20.TransactOpts, _version, _implementation)
}

// TTFT20AddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the TTFT20 contract.
type TTFT20AddedOwnerIterator struct {
	Event *TTFT20AddedOwner // Event containing the contract specifics and raw log

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
func (it *TTFT20AddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20AddedOwner)
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
		it.Event = new(TTFT20AddedOwner)
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
func (it *TTFT20AddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20AddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20AddedOwner represents a AddedOwner event raised by the TTFT20 contract.
type TTFT20AddedOwner struct {
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_TTFT20 *TTFT20Filterer) FilterAddedOwner(opts *bind.FilterOpts) (*TTFT20AddedOwnerIterator, error) {

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &TTFT20AddedOwnerIterator{contract: _TTFT20.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_TTFT20 *TTFT20Filterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *TTFT20AddedOwner) (event.Subscription, error) {

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20AddedOwner)
				if err := _TTFT20.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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
// Solidity: e Approval(tokenOwner indexed address, spender indexed address, tokens uint256)
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
// Solidity: e Approval(tokenOwner indexed address, spender indexed address, tokens uint256)
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
	Txid     common.Hash
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a.
//
// Solidity: e Mint(receiver indexed address, tokens uint256, txid indexed string)
func (_TTFT20 *TTFT20Filterer) FilterMint(opts *bind.FilterOpts, receiver []common.Address, txid []string) (*TTFT20MintIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Mint", receiverRule, txidRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20MintIterator{contract: _TTFT20.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a.
//
// Solidity: e Mint(receiver indexed address, tokens uint256, txid indexed string)
func (_TTFT20 *TTFT20Filterer) WatchMint(opts *bind.WatchOpts, sink chan<- *TTFT20Mint, receiver []common.Address, txid []string) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	var txidRule []interface{}
	for _, txidItem := range txid {
		txidRule = append(txidRule, txidItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Mint", receiverRule, txidRule)
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
// Solidity: e RegisterWithdrawalAddress(addr indexed address)
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
// Solidity: e RegisterWithdrawalAddress(addr indexed address)
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

// TTFT20RemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the TTFT20 contract.
type TTFT20RemovedOwnerIterator struct {
	Event *TTFT20RemovedOwner // Event containing the contract specifics and raw log

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
func (it *TTFT20RemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20RemovedOwner)
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
		it.Event = new(TTFT20RemovedOwner)
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
func (it *TTFT20RemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20RemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20RemovedOwner represents a RemovedOwner event raised by the TTFT20 contract.
type TTFT20RemovedOwner struct {
	RemovedOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_TTFT20 *TTFT20Filterer) FilterRemovedOwner(opts *bind.FilterOpts) (*TTFT20RemovedOwnerIterator, error) {

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &TTFT20RemovedOwnerIterator{contract: _TTFT20.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_TTFT20 *TTFT20Filterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *TTFT20RemovedOwner) (event.Subscription, error) {

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20RemovedOwner)
				if err := _TTFT20.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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
// Solidity: e Transfer(from indexed address, to indexed address, tokens uint256)
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
// Solidity: e Transfer(from indexed address, to indexed address, tokens uint256)
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

// TTFT20UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TTFT20 contract.
type TTFT20UpgradedIterator struct {
	Event *TTFT20Upgraded // Event containing the contract specifics and raw log

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
func (it *TTFT20UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TTFT20Upgraded)
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
		it.Event = new(TTFT20Upgraded)
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
func (it *TTFT20UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TTFT20UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TTFT20Upgraded represents a Upgraded event raised by the TTFT20 contract.
type TTFT20Upgraded struct {
	Version        common.Hash
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0x8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e.
//
// Solidity: e Upgraded(version indexed string, implementation indexed address)
func (_TTFT20 *TTFT20Filterer) FilterUpgraded(opts *bind.FilterOpts, version []string, implementation []common.Address) (*TTFT20UpgradedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Upgraded", versionRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20UpgradedIterator{contract: _TTFT20.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0x8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e.
//
// Solidity: e Upgraded(version indexed string, implementation indexed address)
func (_TTFT20 *TTFT20Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TTFT20Upgraded, version []string, implementation []common.Address) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Upgraded", versionRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TTFT20Upgraded)
				if err := _TTFT20.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
	From     common.Address
	Receiver common.Address
	Tokens   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: e Withdraw(from indexed address, receiver indexed address, tokens uint256)
func (_TTFT20 *TTFT20Filterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address, receiver []common.Address) (*TTFT20WithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _TTFT20.contract.FilterLogs(opts, "Withdraw", fromRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &TTFT20WithdrawIterator{contract: _TTFT20.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: e Withdraw(from indexed address, receiver indexed address, tokens uint256)
func (_TTFT20 *TTFT20Filterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *TTFT20Withdraw, from []common.Address, receiver []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _TTFT20.contract.WatchLogs(opts, "Withdraw", fromRule, receiverRule)
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

// TokenStorageABI is the input ABI used to generate the binding from.
const TokenStorageABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TokenStorageBin is the compiled bytecode used for deploying new contracts.
const TokenStorageBin = `0x608060405234801561001057600080fd5b5060408051808201909152600581527f54465432300000000000000000000000000000000000000000000000000000006020820152610057906401000000006100d8810204565b60408051808201909152601881527f54465420455243323020726570726573656e746174696f6e0000000000000000602082015261009d9064010000000061014b810204565b60096100b1816401000000006101bb810204565b64174876e80060ff8216600a0a026100d18164010000000061022e810204565b505061036f565b6101486040516020018080602001828103825260068152602001807f73796d626f6c0000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208261029e640100000000026401000000009004565b50565b6101486040516020018080602001828103825260048152602001807f6e616d6500000000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208261029e640100000000026401000000009004565b6101486040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208260ff166102c2640100000000026401000000009004565b61014860405160200180806020018281038252600b8152602001807f746f74616c537570706c7900000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826102c2640100000000026401000000009004565b600082815260016020908152604090912082516102bd928401906102d4565b505050565b60009182526020829052604090912055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061031557805160ff1916838001178555610342565b82800160010185558215610342579182015b82811115610342578251825591602001919060010190610327565b5061034e929150610352565b5090565b61036c91905b8082111561034e5760008155600101610358565b90565b60358061037d6000396000f3fe6080604052600080fdfea165627a7a72305820f0c9bbd1ab537a932cb2bb6985808ac8257813ad7ba3c795e1ab2e453c04189a0029`

// DeployTokenStorage deploys a new Ethereum contract, binding an instance of TokenStorage to it.
func DeployTokenStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TokenStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenStorageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenStorage{TokenStorageCaller: TokenStorageCaller{contract: contract}, TokenStorageTransactor: TokenStorageTransactor{contract: contract}, TokenStorageFilterer: TokenStorageFilterer{contract: contract}}, nil
}

// TokenStorage is an auto generated Go binding around an Ethereum contract.
type TokenStorage struct {
	TokenStorageCaller     // Read-only binding to the contract
	TokenStorageTransactor // Write-only binding to the contract
	TokenStorageFilterer   // Log filterer for contract events
}

// TokenStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenStorageSession struct {
	Contract     *TokenStorage     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenStorageCallerSession struct {
	Contract *TokenStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenStorageTransactorSession struct {
	Contract     *TokenStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenStorageRaw struct {
	Contract *TokenStorage // Generic contract binding to access the raw methods on
}

// TokenStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenStorageCallerRaw struct {
	Contract *TokenStorageCaller // Generic read-only contract binding to access the raw methods on
}

// TokenStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenStorageTransactorRaw struct {
	Contract *TokenStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenStorage creates a new instance of TokenStorage, bound to a specific deployed contract.
func NewTokenStorage(address common.Address, backend bind.ContractBackend) (*TokenStorage, error) {
	contract, err := bindTokenStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenStorage{TokenStorageCaller: TokenStorageCaller{contract: contract}, TokenStorageTransactor: TokenStorageTransactor{contract: contract}, TokenStorageFilterer: TokenStorageFilterer{contract: contract}}, nil
}

// NewTokenStorageCaller creates a new read-only instance of TokenStorage, bound to a specific deployed contract.
func NewTokenStorageCaller(address common.Address, caller bind.ContractCaller) (*TokenStorageCaller, error) {
	contract, err := bindTokenStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenStorageCaller{contract: contract}, nil
}

// NewTokenStorageTransactor creates a new write-only instance of TokenStorage, bound to a specific deployed contract.
func NewTokenStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenStorageTransactor, error) {
	contract, err := bindTokenStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenStorageTransactor{contract: contract}, nil
}

// NewTokenStorageFilterer creates a new log filterer instance of TokenStorage, bound to a specific deployed contract.
func NewTokenStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenStorageFilterer, error) {
	contract, err := bindTokenStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenStorageFilterer{contract: contract}, nil
}

// bindTokenStorage binds a generic wrapper to an already deployed contract.
func bindTokenStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenStorage *TokenStorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenStorage.Contract.TokenStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenStorage *TokenStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStorage.Contract.TokenStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenStorage *TokenStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenStorage.Contract.TokenStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenStorage *TokenStorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenStorage *TokenStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenStorage *TokenStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenStorage.Contract.contract.Transact(opts, method, params...)
}

// UpgradeableABI is the input ABI used to generate the binding from.
const UpgradeableABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_toRemove\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_version\",\"type\":\"string\"},{\"name\":\"_implementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"version\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"removedOwner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"}]"

// UpgradeableBin is the compiled bytecode used for deploying new contracts.
const UpgradeableBin = `0x60806040526100163364010000000061001b810204565b6100b4565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120610091906001640100000000610094810204565b50565b600091825260046020526040909120805460ff1916911515919091179055565b610900806100c36000396000f3fe608060405234801561001057600080fd5b5060043610610073577c01000000000000000000000000000000000000000000000000000000006000350463173825d9811461007857806354fd4d50146100a05780635a8b1a9f1461011d5780635c60da1b146101ce5780637065cb48146101f2575b600080fd5b61009e6004803603602081101561008e57600080fd5b5035600160a060020a0316610218565b005b6100a861029f565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e25781810151838201526020016100ca565b50505050905090810190601f16801561010f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61009e6004803603604081101561013357600080fd5b81019060208101813564010000000081111561014e57600080fd5b82018360208201111561016057600080fd5b8035906020019184600183028401116401000000008311171561018257600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050509035600160a060020a031691506102af9050565b6101d661038c565b60408051600160a060020a039092168252519081900360200190f35b61009e6004803603602081101561020857600080fd5b5035600160a060020a0316610396565b61022133610407565b151561022c57600080fd5b600160a060020a038116151561024157600080fd5b600160a060020a03811633141561025757600080fd5b6102608161047b565b60408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b60606102a96104e9565b90505b90565b6102b833610407565b15156102c357600080fd5b80600160a060020a03166102d5610544565b600160a060020a031614156102e957600080fd5b6102f2826105a6565b6102fb81610607565b80600160a060020a0316826040518082805190602001908083835b602083106103355780518252601f199092019160209182019101610316565b5181516020939093036101000a60001901801990911692169190911790526040519201829003822093507f8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e92506000919050a35050565b60006102a9610544565b61039f33610407565b15156103aa57600080fd5b600160a060020a03811615156103bf57600080fd5b6103c881610668565b60408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120600090610475906106d5565b92915050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206104e6906106ea565b50565b6040805160208082018190526007828401527f76657273696f6e00000000000000000000000000000000000000000000000000606083810191909152835180840382018152608090930190935281519101206102a990610702565b60006102a960405160200180806020018281038252600e8152602001807f696d706c656d656e746174696f6e000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001206107a2565b6104e66040516020018080602001828103825260078152602001807f76657273696f6e0000000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826107bd565b6104e660405160200180806020018281038252600e8152602001807f696d706c656d656e746174696f6e00000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826107e1565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206104e690600161081c565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b60008181526001602081815260409283902080548451600294821615610100026000190190911693909304601f810183900483028401830190945283835260609390918301828280156107965780601f1061076b57610100808354040283529160200191610796565b820191906000526020600020905b81548152906001019060200180831161077957829003601f168201915b50505050509050919050565b600090815260026020526040902054600160a060020a031690565b600082815260016020908152604090912082516107dc9284019061083c565b505050565b600091825260026020526040909120805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909216919091179055565b600091825260046020526040909120805460ff1916911515919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061087d57805160ff19168380011785556108aa565b828001600101855582156108aa579182015b828111156108aa57825182559160200191906001019061088f565b506108b69291506108ba565b5090565b6102ac91905b808211156108b657600081556001016108c056fea165627a7a72305820a64c7b2a4f8f8d73d7d703a5c0f5d0bbec3e5f4ac2e936c6e1a6b62b73c207bf0029`

// DeployUpgradeable deploys a new Ethereum contract, binding an instance of Upgradeable to it.
func DeployUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Upgradeable, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeableABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Upgradeable{UpgradeableCaller: UpgradeableCaller{contract: contract}, UpgradeableTransactor: UpgradeableTransactor{contract: contract}, UpgradeableFilterer: UpgradeableFilterer{contract: contract}}, nil
}

// Upgradeable is an auto generated Go binding around an Ethereum contract.
type Upgradeable struct {
	UpgradeableCaller     // Read-only binding to the contract
	UpgradeableTransactor // Write-only binding to the contract
	UpgradeableFilterer   // Log filterer for contract events
}

// UpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeableSession struct {
	Contract     *Upgradeable      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeableCallerSession struct {
	Contract *UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// UpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeableTransactorSession struct {
	Contract     *UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeableRaw struct {
	Contract *Upgradeable // Generic contract binding to access the raw methods on
}

// UpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeableCallerRaw struct {
	Contract *UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeableTransactorRaw struct {
	Contract *UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeable creates a new instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeable(address common.Address, backend bind.ContractBackend) (*Upgradeable, error) {
	contract, err := bindUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Upgradeable{UpgradeableCaller: UpgradeableCaller{contract: contract}, UpgradeableTransactor: UpgradeableTransactor{contract: contract}, UpgradeableFilterer: UpgradeableFilterer{contract: contract}}, nil
}

// NewUpgradeableCaller creates a new read-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*UpgradeableCaller, error) {
	contract, err := bindUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableCaller{contract: contract}, nil
}

// NewUpgradeableTransactor creates a new write-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeableTransactor, error) {
	contract, err := bindUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableTransactor{contract: contract}, nil
}

// NewUpgradeableFilterer creates a new log filterer instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeableFilterer, error) {
	contract, err := bindUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeableFilterer{contract: contract}, nil
}

// bindUpgradeable binds a generic wrapper to an already deployed contract.
func bindUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_Upgradeable *UpgradeableCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Upgradeable.contract.Call(opts, out, "implementation")
	return *ret0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_Upgradeable *UpgradeableSession) Implementation() (common.Address, error) {
	return _Upgradeable.Contract.Implementation(&_Upgradeable.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() constant returns(address)
func (_Upgradeable *UpgradeableCallerSession) Implementation() (common.Address, error) {
	return _Upgradeable.Contract.Implementation(&_Upgradeable.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_Upgradeable *UpgradeableCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Upgradeable.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_Upgradeable *UpgradeableSession) Version() (string, error) {
	return _Upgradeable.Contract.Version(&_Upgradeable.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_Upgradeable *UpgradeableCallerSession) Version() (string, error) {
	return _Upgradeable.Contract.Version(&_Upgradeable.CallOpts)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_Upgradeable *UpgradeableTransactor) AddOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "addOwner", _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_Upgradeable *UpgradeableSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Upgradeable.Contract.AddOwner(&_Upgradeable.TransactOpts, _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_Upgradeable *UpgradeableTransactorSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Upgradeable.Contract.AddOwner(&_Upgradeable.TransactOpts, _newOwner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_Upgradeable *UpgradeableTransactor) RemoveOwner(opts *bind.TransactOpts, _toRemove common.Address) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "removeOwner", _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_Upgradeable *UpgradeableSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _Upgradeable.Contract.RemoveOwner(&_Upgradeable.TransactOpts, _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_Upgradeable *UpgradeableTransactorSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _Upgradeable.Contract.RemoveOwner(&_Upgradeable.TransactOpts, _toRemove)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_Upgradeable *UpgradeableTransactor) UpgradeTo(opts *bind.TransactOpts, _version string, _implementation common.Address) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "upgradeTo", _version, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_Upgradeable *UpgradeableSession) UpgradeTo(_version string, _implementation common.Address) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeTo(&_Upgradeable.TransactOpts, _version, _implementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x5a8b1a9f.
//
// Solidity: function upgradeTo(_version string, _implementation address) returns()
func (_Upgradeable *UpgradeableTransactorSession) UpgradeTo(_version string, _implementation common.Address) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeTo(&_Upgradeable.TransactOpts, _version, _implementation)
}

// UpgradeableAddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the Upgradeable contract.
type UpgradeableAddedOwnerIterator struct {
	Event *UpgradeableAddedOwner // Event containing the contract specifics and raw log

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
func (it *UpgradeableAddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeableAddedOwner)
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
		it.Event = new(UpgradeableAddedOwner)
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
func (it *UpgradeableAddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeableAddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeableAddedOwner represents a AddedOwner event raised by the Upgradeable contract.
type UpgradeableAddedOwner struct {
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_Upgradeable *UpgradeableFilterer) FilterAddedOwner(opts *bind.FilterOpts) (*UpgradeableAddedOwnerIterator, error) {

	logs, sub, err := _Upgradeable.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &UpgradeableAddedOwnerIterator{contract: _Upgradeable.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_Upgradeable *UpgradeableFilterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *UpgradeableAddedOwner) (event.Subscription, error) {

	logs, sub, err := _Upgradeable.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeableAddedOwner)
				if err := _Upgradeable.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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

// UpgradeableRemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the Upgradeable contract.
type UpgradeableRemovedOwnerIterator struct {
	Event *UpgradeableRemovedOwner // Event containing the contract specifics and raw log

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
func (it *UpgradeableRemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeableRemovedOwner)
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
		it.Event = new(UpgradeableRemovedOwner)
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
func (it *UpgradeableRemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeableRemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeableRemovedOwner represents a RemovedOwner event raised by the Upgradeable contract.
type UpgradeableRemovedOwner struct {
	RemovedOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_Upgradeable *UpgradeableFilterer) FilterRemovedOwner(opts *bind.FilterOpts) (*UpgradeableRemovedOwnerIterator, error) {

	logs, sub, err := _Upgradeable.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &UpgradeableRemovedOwnerIterator{contract: _Upgradeable.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_Upgradeable *UpgradeableFilterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *UpgradeableRemovedOwner) (event.Subscription, error) {

	logs, sub, err := _Upgradeable.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeableRemovedOwner)
				if err := _Upgradeable.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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

// UpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Upgradeable contract.
type UpgradeableUpgradedIterator struct {
	Event *UpgradeableUpgraded // Event containing the contract specifics and raw log

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
func (it *UpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeableUpgraded)
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
		it.Event = new(UpgradeableUpgraded)
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
func (it *UpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeableUpgraded represents a Upgraded event raised by the Upgradeable contract.
type UpgradeableUpgraded struct {
	Version        common.Hash
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0x8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e.
//
// Solidity: e Upgraded(version indexed string, implementation indexed address)
func (_Upgradeable *UpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, version []string, implementation []common.Address) (*UpgradeableUpgradedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Upgradeable.contract.FilterLogs(opts, "Upgraded", versionRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeableUpgradedIterator{contract: _Upgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0x8e05e0e35ff592971ca8b477d4285a33a61ded208d644042667b78693a472f5e.
//
// Solidity: e Upgraded(version indexed string, implementation indexed address)
func (_Upgradeable *UpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UpgradeableUpgraded, version []string, implementation []common.Address) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Upgradeable.contract.WatchLogs(opts, "Upgraded", versionRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeableUpgraded)
				if err := _Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
