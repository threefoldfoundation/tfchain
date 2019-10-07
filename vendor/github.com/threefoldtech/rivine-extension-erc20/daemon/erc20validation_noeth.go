// +build noeth

package daemon

import (
	erc20bridge "github.com/threefoldtech/rivine-extension-erc20/api/bridge"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
)

// NewERC20NodeValidator returns always the NopERC20TransactionValidator,
// as erc20 validation is disabled in tfchain daemons running the `noeth` build flag.
func NewERC20NodeValidator(ERC20NodeValidatorConfig, <-chan struct{}) (erc20types.ERC20TransactionValidator, error) {
	return erc20types.NopERC20TransactionValidator{}, nil
}

func NewERC20NodeValidatorFromBridgeContract(contract *erc20bridge.BridgeContract) (erc20types.ERC20TransactionValidator, error) {
	return erc20types.NopERC20TransactionValidator{}, nil
}
