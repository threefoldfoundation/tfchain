package erc20

import (
	"math/big"
)

// Denominate converts gwei units into ether units
func Denominate(gwei *big.Int) *big.Int {
	ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	balance := new(big.Int).Div(gwei, ether)
	return balance
}
