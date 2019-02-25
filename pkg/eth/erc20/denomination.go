package erc20

import (
	"math"
	"math/big"
)

// Denominate converts gwei units into ether units
func Denominate(gwei *big.Int) string {
	fbalance := new(big.Float)
	fbalance.SetString(gwei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18))).String()

	return ethValue + " ETH"
}
