package types

import (
	"github.com/threefoldtech/rivine/types"
)

const (
	// TransactionVersionMinterDefinition defines the Transaction version
	// for a MinterDefinition Transaction.
	//
	// See the `MinterDefinitionTransactionController` and `MinterDefinitionTransaction`
	// types for more information.
	TransactionVersionMinterDefinition types.TransactionVersion = iota + 128
	// TransactionVersionCoinCreation defines the Transaction version
	// for a CoinCreation Transaction.
	//
	// See the `CoinCreationTransactionController` and `CoinCreationTransaction`
	// types for more information.
	TransactionVersionCoinCreation
)

const (
	// TransactionVersionERC20Conversion defines the Transaction version
	// for an ERC20ConvertTransaction, used to convert TFT into ERC20 funds.
	TransactionVersionERC20Conversion types.TransactionVersion = iota + 208
	// TransactionVersionERC20CoinCreation defines the Transaction version
	// for an ERC20CoinCreationTransaction, used to convert ERC20 funds into TFT.
	TransactionVersionERC20CoinCreation
	// TransactionVersionERC20AddressRegistration defines the Transaction version
	// for an TransactionVersionERC20AddressRegistration, used to register an ERC20 address,
	// linked to an TFT address.
	TransactionVersionERC20AddressRegistration
)
