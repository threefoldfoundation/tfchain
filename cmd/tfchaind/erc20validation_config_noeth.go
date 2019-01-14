// +build noeth

package main

import (
	flag "github.com/spf13/pflag"
)

// SetFlags defines no flags at all,
// when compiling using the `noeth` buildig flag, disabling any eth dependency.
func (cfg *ERC20NodeValidatorConfig) SetFlags(*flag.FlagSet) {}
