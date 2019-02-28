package erc20

import (
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestDenominate(t *testing.T) {
	f := func(str string) *big.Int {
		i := new(big.Int)
		err := i.UnmarshalText([]byte(str))
		if err != nil {
			t.Errorf("failed to unmarshal text %q: %v", str, err)
		}
		return i
	}
	testCases := []struct {
		Input  *big.Int
		Output string
	}{
		{nil, "0 ETH"},
		{new(big.Int), "0 ETH"},
		{big.NewInt(0), "0 ETH"},
		{big.NewInt(-0), "0 ETH"},
		{big.NewInt(1), "0.000000000000000001 ETH"},
		{big.NewInt(-1), "-0.000000000000000001 ETH"},
		{big.NewInt(-120000000000), "-0.00000012 ETH"},
		{big.NewInt(math.MaxUint64 >> 1), "9.223372036854775807 ETH"},
		{f("87654321012345678942"), "87.654321012345678942 ETH"},
		{f(strings.Repeat("123456789123456789", 4)), strings.Repeat("123456789123456789", 3) + ".123456789123456789 ETH"},
	}
	for caseIndex, testCase := range testCases {
		output := Denominate(testCase.Input)
		if output != testCase.Output {
			t.Errorf("testcase #%d failed: %s != %s", caseIndex, output, testCase.Output)
		}
	}
}
