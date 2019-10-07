package bridge

import (
	"math"
	"math/big"
)

// Denominate converts gwei units into ether units
func Denominate(gwei *big.Int) string {
	const precision = 256
	fbalance := new(big.Float).SetPrec(precision)
	fbalance.SetString(gwei.String())
	ethValue := new(big.Float).SetPrec(precision)
	divisor := big.NewFloat(math.Pow10(18)).SetPrec(precision)
	return ethValue.Quo(fbalance, divisor).SetPrec(precision).Text('f', -1) + " ETH"
}
