// +build !noeth

package main

import (
	flag "github.com/spf13/pflag"
)

// SetFlags defines the ERC20NodeValidatorConfig as flags.
func (cfg *ERC20NodeValidatorConfig) SetFlags(flags *flag.FlagSet) {
	flags.BoolVar(
		&cfg.Enabled,
		"ethvalidation", false,
		"enable full validation of ERC20 validation, attaching this node in light-mode to the ETH network",
	)
	flags.IntVar(
		&cfg.Port,
		"ethport", 3003,
		"network port used by peers on the ETH network to connect to this node should ethvalidation be enabled",
	)
	flags.StringVar(
		&cfg.NetworkName,
		"ethnetwork", "",
		"The ethereum network, {main, rinkeby, ropsten}, defaults to the TFT-linked network",
	)
}
