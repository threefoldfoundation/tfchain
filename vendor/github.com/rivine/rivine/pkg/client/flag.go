package client

import (
	"fmt"
	"os"

	"github.com/rivine/rivine/types"
	"github.com/spf13/pflag"
)

type coinFlag struct {
	str string
	cli *CommandLineClient
}

// String implements pflag.Value.String
func (c coinFlag) String() string {
	return c.str
}

// Set implements pflag.Value.Set
func (c *coinFlag) Set(s string) error {
	c.str = s
	return nil
}

// Type implements pflag.Value.Type
func (c coinFlag) Type() string {
	return "Coin"
}

func parseCoinArg(cc CurrencyConvertor, str string) types.Currency {
	amount, err := cc.ParseCoinString(str)
	if err != nil {
		fmt.Fprintln(os.Stderr, cc.CoinArgDescription("amount"))
		DieWithExitCode(ExitCodeUsage, "failed to parse coin-typed argument: ", err)
		return types.Currency{}
	}
	return amount
}

func (c coinFlag) Amount() types.Currency {
	if c.str == "" {
		return types.Currency{}
	}
	return parseCoinArg(c.cli.CreateCurrencyConvertor(), c.str)
}

var (
	_ pflag.Value = (*coinFlag)(nil)
)
