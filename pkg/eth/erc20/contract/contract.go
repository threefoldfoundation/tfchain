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
const OwnedBin = `0x608060405234801561001057600080fd5b5061002333640100000000610029810204565b506100c7565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906100a29060016401000000006100a7810204565b919050565b600091825260046020526040909120805460ff1916911515919091179055565b61037f806100d66000396000f3fe60806040526004361061004b5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663173825d981146100505780637065cb4814610085575b600080fd5b34801561005c57600080fd5b506100836004803603602081101561007357600080fd5b5035600160a060020a03166100b8565b005b34801561009157600080fd5b50610083600480360360208110156100a857600080fd5b5035600160a060020a0316610140565b6100c1336101b2565b15156100cc57600080fd5b600160a060020a03811615156100e157600080fd5b600160a060020a0381163314156100f757600080fd5b61010081610228565b5060408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b610149336101b2565b151561015457600080fd5b600160a060020a038116151561016957600080fd5b61017281610296565b5060408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a0909201909252805191012060009061022090610306565b90505b919050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906102239061031b565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120600090610223906001610333565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b600091825260046020526040909120805460ff191691151591909117905556fea165627a7a7230582099f82a33d75e3beaae0fdabad6c668689742ba4b9f982fbeb421865e32c764e80029`

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

// OwnedTokenStorageABI is the input ABI used to generate the binding from.
const OwnedTokenStorageABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_toRemove\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"removedOwner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"}]"

// OwnedTokenStorageBin is the compiled bytecode used for deploying new contracts.
const OwnedTokenStorageBin = `0x60c0604052600660809081527f545446543230000000000000000000000000000000000000000000000000000060a052610041906401000000006100d9810204565b60408051808201909152601981527f5454465420455243323020726570726573656e746174696f6e0000000000000060208201526100879064010000000061014c810204565b601261009b816401000000006101bc810204565b620f424060ff8216600a0a026100b98164010000000061022f810204565b50506100d33361029f640100000000026401000000009004565b5061040e565b6101496040516020018080602001828103825260068152602001807f73796d626f6c0000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208261031d640100000000026401000000009004565b50565b6101496040516020018080602001828103825260048152602001807f6e616d6500000000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208261031d640100000000026401000000009004565b6101496040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208260ff16610341640100000000026401000000009004565b61014960405160200180806020018281038252600b8152602001807f746f74616c537570706c790000000000000000000000000000000000000000008152506020019150506040516020818303038152906040528051906020012082610341640100000000026401000000009004565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120600090610318906001640100000000610353810204565b919050565b6000828152600160209081526040909120825161033c92840190610373565b505050565b60009182526020829052604090912055565b600091825260046020526040909120805460ff1916911515919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106103b457805160ff19168380011785556103e1565b828001600101855582156103e1579182015b828111156103e15782518255916020019190600101906103c6565b506103ed9291506103f1565b5090565b61040b91905b808211156103ed57600081556001016103f7565b90565b61037f8061041d6000396000f3fe60806040526004361061004b5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663173825d981146100505780637065cb4814610085575b600080fd5b34801561005c57600080fd5b506100836004803603602081101561007357600080fd5b5035600160a060020a03166100b8565b005b34801561009157600080fd5b50610083600480360360208110156100a857600080fd5b5035600160a060020a0316610140565b6100c1336101b2565b15156100cc57600080fd5b600160a060020a03811615156100e157600080fd5b600160a060020a0381163314156100f757600080fd5b61010081610228565b5060408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b610149336101b2565b151561015457600080fd5b600160a060020a038116151561016957600080fd5b61017281610296565b5060408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a0909201909252805191012060009061022090610306565b90505b919050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906102239061031b565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120600090610223906001610333565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b600091825260046020526040909120805460ff191691151591909117905556fea165627a7a723058208cc3e2563da34ebd3c444748af71bc9b8c81fa7a68784ee6dec8333ec6f2ad5d0029`

// DeployOwnedTokenStorage deploys a new Ethereum contract, binding an instance of OwnedTokenStorage to it.
func DeployOwnedTokenStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OwnedTokenStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedTokenStorageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedTokenStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnedTokenStorage{OwnedTokenStorageCaller: OwnedTokenStorageCaller{contract: contract}, OwnedTokenStorageTransactor: OwnedTokenStorageTransactor{contract: contract}, OwnedTokenStorageFilterer: OwnedTokenStorageFilterer{contract: contract}}, nil
}

// OwnedTokenStorage is an auto generated Go binding around an Ethereum contract.
type OwnedTokenStorage struct {
	OwnedTokenStorageCaller     // Read-only binding to the contract
	OwnedTokenStorageTransactor // Write-only binding to the contract
	OwnedTokenStorageFilterer   // Log filterer for contract events
}

// OwnedTokenStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedTokenStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTokenStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTokenStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTokenStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedTokenStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTokenStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedTokenStorageSession struct {
	Contract     *OwnedTokenStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnedTokenStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedTokenStorageCallerSession struct {
	Contract *OwnedTokenStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// OwnedTokenStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTokenStorageTransactorSession struct {
	Contract     *OwnedTokenStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// OwnedTokenStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedTokenStorageRaw struct {
	Contract *OwnedTokenStorage // Generic contract binding to access the raw methods on
}

// OwnedTokenStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedTokenStorageCallerRaw struct {
	Contract *OwnedTokenStorageCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTokenStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTokenStorageTransactorRaw struct {
	Contract *OwnedTokenStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnedTokenStorage creates a new instance of OwnedTokenStorage, bound to a specific deployed contract.
func NewOwnedTokenStorage(address common.Address, backend bind.ContractBackend) (*OwnedTokenStorage, error) {
	contract, err := bindOwnedTokenStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnedTokenStorage{OwnedTokenStorageCaller: OwnedTokenStorageCaller{contract: contract}, OwnedTokenStorageTransactor: OwnedTokenStorageTransactor{contract: contract}, OwnedTokenStorageFilterer: OwnedTokenStorageFilterer{contract: contract}}, nil
}

// NewOwnedTokenStorageCaller creates a new read-only instance of OwnedTokenStorage, bound to a specific deployed contract.
func NewOwnedTokenStorageCaller(address common.Address, caller bind.ContractCaller) (*OwnedTokenStorageCaller, error) {
	contract, err := bindOwnedTokenStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTokenStorageCaller{contract: contract}, nil
}

// NewOwnedTokenStorageTransactor creates a new write-only instance of OwnedTokenStorage, bound to a specific deployed contract.
func NewOwnedTokenStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTokenStorageTransactor, error) {
	contract, err := bindOwnedTokenStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTokenStorageTransactor{contract: contract}, nil
}

// NewOwnedTokenStorageFilterer creates a new log filterer instance of OwnedTokenStorage, bound to a specific deployed contract.
func NewOwnedTokenStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedTokenStorageFilterer, error) {
	contract, err := bindOwnedTokenStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedTokenStorageFilterer{contract: contract}, nil
}

// bindOwnedTokenStorage binds a generic wrapper to an already deployed contract.
func bindOwnedTokenStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedTokenStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnedTokenStorage *OwnedTokenStorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OwnedTokenStorage.Contract.OwnedTokenStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnedTokenStorage *OwnedTokenStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.OwnedTokenStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnedTokenStorage *OwnedTokenStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.OwnedTokenStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnedTokenStorage *OwnedTokenStorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OwnedTokenStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnedTokenStorage *OwnedTokenStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnedTokenStorage *OwnedTokenStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.contract.Transact(opts, method, params...)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_OwnedTokenStorage *OwnedTokenStorageTransactor) AddOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _OwnedTokenStorage.contract.Transact(opts, "addOwner", _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_OwnedTokenStorage *OwnedTokenStorageSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.AddOwner(&_OwnedTokenStorage.TransactOpts, _newOwner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(_newOwner address) returns()
func (_OwnedTokenStorage *OwnedTokenStorageTransactorSession) AddOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.AddOwner(&_OwnedTokenStorage.TransactOpts, _newOwner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_OwnedTokenStorage *OwnedTokenStorageTransactor) RemoveOwner(opts *bind.TransactOpts, _toRemove common.Address) (*types.Transaction, error) {
	return _OwnedTokenStorage.contract.Transact(opts, "removeOwner", _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_OwnedTokenStorage *OwnedTokenStorageSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.RemoveOwner(&_OwnedTokenStorage.TransactOpts, _toRemove)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(_toRemove address) returns()
func (_OwnedTokenStorage *OwnedTokenStorageTransactorSession) RemoveOwner(_toRemove common.Address) (*types.Transaction, error) {
	return _OwnedTokenStorage.Contract.RemoveOwner(&_OwnedTokenStorage.TransactOpts, _toRemove)
}

// OwnedTokenStorageAddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the OwnedTokenStorage contract.
type OwnedTokenStorageAddedOwnerIterator struct {
	Event *OwnedTokenStorageAddedOwner // Event containing the contract specifics and raw log

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
func (it *OwnedTokenStorageAddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedTokenStorageAddedOwner)
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
		it.Event = new(OwnedTokenStorageAddedOwner)
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
func (it *OwnedTokenStorageAddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedTokenStorageAddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedTokenStorageAddedOwner represents a AddedOwner event raised by the OwnedTokenStorage contract.
type OwnedTokenStorageAddedOwner struct {
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_OwnedTokenStorage *OwnedTokenStorageFilterer) FilterAddedOwner(opts *bind.FilterOpts) (*OwnedTokenStorageAddedOwnerIterator, error) {

	logs, sub, err := _OwnedTokenStorage.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &OwnedTokenStorageAddedOwnerIterator{contract: _OwnedTokenStorage.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: e AddedOwner(newOwner address)
func (_OwnedTokenStorage *OwnedTokenStorageFilterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *OwnedTokenStorageAddedOwner) (event.Subscription, error) {

	logs, sub, err := _OwnedTokenStorage.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedTokenStorageAddedOwner)
				if err := _OwnedTokenStorage.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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

// OwnedTokenStorageRemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the OwnedTokenStorage contract.
type OwnedTokenStorageRemovedOwnerIterator struct {
	Event *OwnedTokenStorageRemovedOwner // Event containing the contract specifics and raw log

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
func (it *OwnedTokenStorageRemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedTokenStorageRemovedOwner)
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
		it.Event = new(OwnedTokenStorageRemovedOwner)
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
func (it *OwnedTokenStorageRemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedTokenStorageRemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedTokenStorageRemovedOwner represents a RemovedOwner event raised by the OwnedTokenStorage contract.
type OwnedTokenStorageRemovedOwner struct {
	RemovedOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_OwnedTokenStorage *OwnedTokenStorageFilterer) FilterRemovedOwner(opts *bind.FilterOpts) (*OwnedTokenStorageRemovedOwnerIterator, error) {

	logs, sub, err := _OwnedTokenStorage.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &OwnedTokenStorageRemovedOwnerIterator{contract: _OwnedTokenStorage.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: e RemovedOwner(removedOwner address)
func (_OwnedTokenStorage *OwnedTokenStorageFilterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *OwnedTokenStorageRemovedOwner) (event.Subscription, error) {

	logs, sub, err := _OwnedTokenStorage.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedTokenStorageRemovedOwner)
				if err := _OwnedTokenStorage.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a72305820e9c30c701aa55547539049ce5993fd232f4838f74349d0980276d05fdf6414830029`

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
const StorageBin = `0x6080604052348015600f57600080fd5b50603580601d6000396000f3fe6080604052600080fdfea165627a7a72305820de3af130e6b28988bc95f8627c87721d624affc25874186b432ab8d21f4d71490029`

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
const TTFT20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_toRemove\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"registerWithdrawalAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"},{\"name\":\"txid\",\"type\":\"string\"}],\"name\":\"mintTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"RegisterWithdrawalAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"txid\",\"type\":\"string\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"removedOwner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"}]"

// TTFT20Bin is the compiled bytecode used for deploying new contracts.
const TTFT20Bin = `0x60c0604052600660809081527f545446543230000000000000000000000000000000000000000000000000000060a0526200004390640100000000620000e4810204565b60408051808201909152601981527f5454465420455243323020726570726573656e746174696f6e0000000000000060208201526200008b9064010000000062000159810204565b6012620000a181640100000000620001cb810204565b620f424060ff8216600a0a02620000c18164010000000062000240810204565b5050620000dd33620002b2640100000000026401000000009004565b506200042f565b620001566040516020018080602001828103825260068152602001807f73796d626f6c0000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208262000332640100000000026401000000009004565b50565b620001566040516020018080602001828103825260048152602001807f6e616d6500000000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208262000332640100000000026401000000009004565b620001566040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208260ff1662000358640100000000026401000000009004565b6200015660405160200180806020018281038252600b8152602001807f746f74616c537570706c79000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208262000358640100000000026401000000009004565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906200032d9060016401000000006200036a810204565b919050565b6000828152600160209081526040909120825162000353928401906200038a565b505050565b60009182526020829052604090912055565b600091825260046020526040909120805460ff1916911515919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620003cd57805160ff1916838001178555620003fd565b82800160010185558215620003fd579182015b82811115620003fd578251825591602001919060010190620003e0565b506200040b9291506200040f565b5090565b6200042c91905b808211156200040b576000815560010162000416565b90565b6113a3806200043f6000396000f3fe6080604052600436106100c45763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde0381146100c9578063095ea7b314610153578063173825d9146101a057806318160ddd146101d557806323b872dd146101fc578063313ce5671461023f57806334ca6a711461026a5780637065cb481461029d57806370a08231146102d057806395d89b4114610303578063a9059cbb14610318578063dd62ed3e14610351578063e67524a31461038c575b600080fd5b3480156100d557600080fd5b506100de610454565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610118578181015183820152602001610100565b50505050905090810190601f1680156101455780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015f57600080fd5b5061018c6004803603604081101561017657600080fd5b50600160a060020a038135169060200135610463565b604080519115158252519081900360200190f35b3480156101ac57600080fd5b506101d3600480360360208110156101c357600080fd5b5035600160a060020a03166104ba565b005b3480156101e157600080fd5b506101ea610542565b60408051918252519081900360200190f35b34801561020857600080fd5b5061018c6004803603606081101561021f57600080fd5b50600160a060020a03813581169160208101359091169060400135610565565b34801561024b57600080fd5b50610254610666565b6040805160ff9092168252519081900360200190f35b34801561027657600080fd5b506101d36004803603602081101561028d57600080fd5b5035600160a060020a0316610670565b3480156102a957600080fd5b506101d3600480360360208110156102c057600080fd5b5035600160a060020a0316610727565b3480156102dc57600080fd5b506101ea600480360360208110156102f357600080fd5b5035600160a060020a0316610799565b34801561030f57600080fd5b506100de6107ac565b34801561032457600080fd5b5061018c6004803603604081101561033b57600080fd5b50600160a060020a0381351690602001356107b6565b34801561035d57600080fd5b506101ea6004803603604081101561037457600080fd5b50600160a060020a0381358116916020013516610877565b34801561039857600080fd5b506101d3600480360360608110156103af57600080fd5b600160a060020a03823516916020810135918101906060810160408201356401000000008111156103df57600080fd5b8201836020820111156103f157600080fd5b8035906020019184600183028401116401000000008311171561041357600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061088a945050505050565b606061045e6109c7565b905090565b6000610470338484610a22565b604080518381529051600160a060020a0385169133917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259181900360200190a35060015b92915050565b6104c333610aa0565b15156104ce57600080fd5b600160a060020a03811615156104e357600080fd5b600160a060020a0381163314156104f957600080fd5b61050281610b0e565b5060408051600160a060020a038316815290517ff8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf9181900360200190a150565b600061045e6105516000610b7c565b610559610bea565b9063ffffffff610c4c16565b600061057f843361057a856105598933610c61565b610a22565b610595846105908461055988610b7c565b610cdc565b61059e83610d4c565b156105f35782600160a060020a031684600160a060020a03167f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb846040518082815260200191505060405180910390a361065c565b610610836105908461060487610b7c565b9063ffffffff610e0116565b82600160a060020a031684600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a35b5060019392505050565b600061045e610e11565b61067933610aa0565b151561068457600080fd5b61068d81610e73565b600061069882610b7c565b905060008111156106ef576106ae826000610cdc565b604080518281529051600160a060020a0384169182917f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9181900360200190a35b604051600160a060020a038316907f77bc19082a31daad021d73c26bb4f6e74100a41c98099405e92a9323d133e60290600090a25050565b61073033610aa0565b151561073b57600080fd5b600160a060020a038116151561075057600080fd5b61075981610f1f565b5060408051600160a060020a038316815290517f9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea269181900360200190a150565b60006107a482610b7c565b90505b919050565b606061045e610f8f565b60006107c9336105908461055933610b7c565b6107d283610d4c565b1561081c57604080518381529051600160a060020a0385169133917f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb9181900360200190a361086e565b61082d836105908461060487610b7c565b604080518381529051600160a060020a0385169133917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a35b50600192915050565b60006108838383610c61565b9392505050565b61089333610aa0565b151561089e57600080fd5b6108a781610fea565b1561091357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f544654207472616e736163746f6e20494420616c7265616479206b6e6f776e00604482015290519081900360640190fd5b61091c81611128565b61092d836105908461060487610b7c565b806040518082805190602001908083835b6020831061095d5780518252601f19909201916020918201910161093e565b51815160209384036101000a6000190180199092169116179052604080519290940182900382208883529351939550600160a060020a03891694507f85a66b9141978db9980f7e0ce3b468cebf4f7999f32b23091c5c03e798b1ba7a9391829003019150a3505050565b6040805160208082018190526004828401527f6e616d65000000000000000000000000000000000000000000000000000000006060838101919091528351808403820181526080909301909352815191012061045e90611266565b60408051600160a060020a03808616828401528416606080830191909152602080830191909152600760808301527f616c6c6f7765640000000000000000000000000000000000000000000000000060a0808401919091528351808403909101815260c09092019092528051910120610a9b9082611306565b505050565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906107a490611318565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906107a79061132d565b60408051600160a060020a038316818301526020808201839052600760608301527f62616c616e6365000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906107a490611345565b600061045e60405160200180806020018281038252600b8152602001807f746f74616c537570706c7900000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120611345565b600082821115610c5b57600080fd5b50900390565b60408051600160a060020a03808516828401528316606080830191909152602080830191909152600760808301527f616c6c6f7765640000000000000000000000000000000000000000000000000060a0808401919091528351808403909101815260c0909201909252805191012060009061088390611345565b60408051600160a060020a038416818301526020808201839052600760608301527f62616c616e6365000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a09092019092528051910120610d489082611306565b5050565b60006107a48260405160200180806020018060200184600160a060020a0316600160a060020a03168152602001838103835260078152602001807f61646472657373000000000000000000000000000000000000000000000000008152506020018381038252600a8152602001807f7769746864726177616c00000000000000000000000000000000000000000000815250602001935050505060405160208183030381529060405280519060200120611318565b818101828110156104b457600080fd5b600061045e6040516020018080602001828103825260088152602001807f646563696d616c7300000000000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120611345565b60408051600160a060020a038316606080830191909152602080830191909152600760808301527f616464726573730000000000000000000000000000000000000000000000000060a08084019190915282840152600a60c08301527f7769746864726177616c0000000000000000000000000000000000000000000060e080840191909152835180840390910181526101009092019092528051910120610f1c906001611357565b50565b60408051600160a060020a038316818301526020808201839052600560608301527f6f776e65720000000000000000000000000000000000000000000000000000006080808401919091528351808403909101815260a090920190925280519101206000906107a7906001611357565b6040805160208082018190526006828401527f73796d626f6c00000000000000000000000000000000000000000000000000006060838101919091528351808403820181526080909301909352815191012061045e90611266565b60006107a4826040516020018080602001806020018060200180602001858103855260048152602001807f6d696e74000000000000000000000000000000000000000000000000000000008152506020018581038452600b8152602001807f7472616e73616374696f6e000000000000000000000000000000000000000000815250602001858103835260028152602001807f6964000000000000000000000000000000000000000000000000000000000000815250602001858103825286818151815260200191508051906020019080838360005b838110156110d85781810151838201526020016110c0565b50505050905090810190601f1680156111055780820380516001836020036101000a031916815260200191505b509550505050505060405160208183030381529060405280519060200120611318565b610f1c816040516020018080602001806020018060200180602001858103855260048152602001807f6d696e74000000000000000000000000000000000000000000000000000000008152506020018581038452600b8152602001807f7472616e73616374696f6e000000000000000000000000000000000000000000815250602001858103835260028152602001807f6964000000000000000000000000000000000000000000000000000000000000815250602001858103825286818151815260200191508051906020019080838360005b838110156112145781810151838201526020016111fc565b50505050905090810190601f1680156112415780820380516001836020036101000a031916815260200191505b5095505050505050604051602081830303815290604052805190602001206001611357565b60008181526001602081815260409283902080548451600294821615610100026000190190911693909304601f810183900483028401830190945283835260609390918301828280156112fa5780601f106112cf576101008083540402835291602001916112fa565b820191906000526020600020905b8154815290600101906020018083116112dd57829003601f168201915b50505050509050919050565b60009182526020829052604090912055565b60009081526004602052604090205460ff1690565b6000908152600460205260409020805460ff19169055565b60009081526020819052604090205490565b600091825260046020526040909120805460ff191691151591909117905556fea165627a7a723058208e471f913aff6639e16ee4f26ca52888fd6ecc35bc30a38570f731c46e7ae78e0029`

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
const TokenStorageBin = `0x608060405234801561001057600080fd5b5060408051808201909152600681527f54544654323000000000000000000000000000000000000000000000000000006020820152610057906401000000006100d6810204565b60408051808201909152601981527f5454465420455243323020726570726573656e746174696f6e00000000000000602082015261009d90640100000000610149810204565b60126100b1816401000000006101b9810204565b620f424060ff8216600a0a026100cf8164010000000061022c810204565b505061036d565b6101466040516020018080602001828103825260068152602001807f73796d626f6c0000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208261029c640100000000026401000000009004565b50565b6101466040516020018080602001828103825260048152602001807f6e616d6500000000000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208261029c640100000000026401000000009004565b6101466040516020018080602001828103825260088152602001807f646563696d616c73000000000000000000000000000000000000000000000000815250602001915050604051602081830303815290604052805190602001208260ff166102c0640100000000026401000000009004565b61014660405160200180806020018281038252600b8152602001807f746f74616c537570706c7900000000000000000000000000000000000000000081525060200191505060405160208183030381529060405280519060200120826102c0640100000000026401000000009004565b600082815260016020908152604090912082516102bb928401906102d2565b505050565b60009182526020829052604090912055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061031357805160ff1916838001178555610340565b82800160010185558215610340579182015b82811115610340578251825591602001919060010190610325565b5061034c929150610350565b5090565b61036a91905b8082111561034c5760008155600101610356565b90565b60358061037b6000396000f3fe6080604052600080fdfea165627a7a72305820f3a729435cb97d47612c35ac105071011e6b3671be7a6fdf12d61e5b60f574b50029`

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
