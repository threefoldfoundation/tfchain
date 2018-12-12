package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/threefoldtech/rivine/modules"
)

// NetAddressArrayFlagVar defines a []modules.NetAddress flag with specified name and usage string.
// The argument s points to a []modules.NetAddress variable in which to store the validated values of the flags.
// The value of each argument will not try to be separated by comma, each value has to be defined as a separate flag (using the same name).
func NetAddressArrayFlagVar(f *pflag.FlagSet, s *[]modules.NetAddress, name string, usage string) {
	f.Var(&netAddressArray{array: s}, name, usage)
}

// NetAddressArrayFlagVarP defines a []modules.NetAddress flag with specified name, shorthand and usage string.
// The argument s points to a []modules.NetAddress variable in which to store the validated values of the flags.
// The value of each argument will not try to be separated by comma, each value has to be defined as a separate flag (using the same name or shorthand).
func NetAddressArrayFlagVarP(f *pflag.FlagSet, s *[]modules.NetAddress, name, shorthand string, usage string) {
	f.VarP(&netAddressArray{array: s}, name, shorthand, usage)
}

type netAddressArray struct {
	array   *[]modules.NetAddress
	changed bool
}

// Set implements pflag.Value.Set
func (flag *netAddressArray) Set(val string) error {
	if !flag.changed {
		*flag.array = make([]modules.NetAddress, 0)
		flag.changed = true
	}
	na := modules.NetAddress(val)
	err := na.IsStdValid()
	if err != nil {
		return fmt.Errorf("invalid network address %v: %v", val, err)
	}
	*flag.array = append(*flag.array, na)
	return nil
}

// Type implements pflag.Value.Type
func (flag *netAddressArray) Type() string {
	return "NetAddressArray"
}

// String implements pflag.Value.String
func (flag *netAddressArray) String() string {
	if flag.array == nil || len(*flag.array) == 0 {
		return ""
	}
	var str string
	for _, na := range *flag.array {
		str += string(na) + ","
	}
	return str[:len(str)-1]
}
