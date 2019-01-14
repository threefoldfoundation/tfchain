// +build noeth

package main

import (
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
)

// NewERC20NodeValidator returns always the NopERC20TransactionValidator,
// as erc20 validation is disabled in tfchain daemons running the `noeth` build flag.
func NewERC20NodeValidator(ERC20NodeValidatorConfig, <-chan struct{}) (tftypes.ERC20TransactionValidator, error) {
	return tftypes.NopERC20TransactionValidator{}, nil
}
