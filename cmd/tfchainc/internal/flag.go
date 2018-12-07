package internal

import (
	"errors"
	"strings"

	"github.com/spf13/pflag"
	"github.com/threefoldfoundation/tfchain/pkg/types"
	rivinetypes "github.com/threefoldtech/rivine/types"
)

// BotNameArrayFlagVar defines a BotName Array flag with specified name and usage string.
// The arguments s points to a BotName slice variable in which to store the interpreted values of the flags.
// The value of each argument will not try to be separated by comma, each value has to be defined as a separate flag.
func BotNameArrayFlagVar(f *pflag.FlagSet, s *[]types.BotName, name string, usage string) {
	f.Var(&botNameArrayFlag{names: s}, name, usage)
}

// BotNameArrayFlagVarP defines a BotName Array flag with specified name, shorthand and usage string.
// The argument s points to a BotName slice variable in which to store the compiled values of the multiple flags.
// The value of each argument will not try to be separated by comma, each value has to be defined as a separate flag (using the same name or shorthand).
func BotNameArrayFlagVarP(f *pflag.FlagSet, s *[]types.BotName, name, shorthand string, usage string) {
	f.VarP(&botNameArrayFlag{names: s}, name, shorthand, usage)
}

type botNameArrayFlag struct {
	names   *[]types.BotName
	changed bool
}

// Set implements pflag.Value.Set
func (flag *botNameArrayFlag) Set(val string) error {
	if !flag.changed {
		*flag.names = make([]types.BotName, 0, 1)
		flag.changed = true
	}
	var newName types.BotName
	err := newName.LoadString(val)
	if err != nil {
		return err
	}
	for _, name := range *flag.names {
		if name.Equals(newName) {
			return errors.New(val + " is already set")
		}
	}
	*flag.names = append(*flag.names, newName)
	return nil
}

// Type implements pflag.Value.Type
func (flag *botNameArrayFlag) Type() string {
	return "BotNameArrayFlag"
}

// String implements pflag.Value.String
func (flag *botNameArrayFlag) String() string {
	vals := make([]string, 0, len(*flag.names))
	for _, name := range *flag.names {
		vals = append(vals, name.String())
	}
	return strings.Join(vals, ",")
}

// NetworkAddressArrayFlagVar defines a NetworkAddress Array flag with specified name and usage string.
// The arguments s points to a NetworkAddress slice variable in which to store the interpreted values of the flags.
// The value of each argument will not try to be separated by comma, each value has to be defined as a separate flag.
func NetworkAddressArrayFlagVar(f *pflag.FlagSet, s *[]types.NetworkAddress, name string, usage string) {
	f.Var(&networkAddressArrayFlag{addresses: s}, name, usage)
}

// NetworkAddressArrayFlagVarP defines a NetworkAddress Array flag with specified name, shorthand and usage string.
// The argument s points to a NetworkAddress slice variable in which to store the compiled values of the flags.
// The value of each argument will not try to be separated by comma, each value has to be defined as a separate flag (using the same name or shorthand).
func NetworkAddressArrayFlagVarP(f *pflag.FlagSet, s *[]types.NetworkAddress, name, shorthand string, usage string) {
	f.VarP(&networkAddressArrayFlag{addresses: s}, name, shorthand, usage)
}

type networkAddressArrayFlag struct {
	addresses *[]types.NetworkAddress
	changed   bool
}

// Set implements pflag.Value.Set
func (flag *networkAddressArrayFlag) Set(val string) error {
	if !flag.changed {
		*flag.addresses = make([]types.NetworkAddress, 0, 1)
		flag.changed = true
	}
	var newAddress types.NetworkAddress
	err := newAddress.LoadString(val)
	if err != nil {
		return err
	}
	for _, address := range *flag.addresses {
		if address.Equals(newAddress) {
			return errors.New(val + " is already set")
		}
	}
	*flag.addresses = append(*flag.addresses, newAddress)
	return nil
}

// Type implements pflag.Value.Type
func (flag *networkAddressArrayFlag) Type() string {
	return "NetworkAddressArrayFlag"
}

// String implements pflag.Value.String
func (flag *networkAddressArrayFlag) String() string {
	vals := make([]string, 0, len(*flag.addresses))
	for _, name := range *flag.addresses {
		vals = append(vals, name.String())
	}
	return strings.Join(vals, ",")
}

// PublicKeyFlagVar defines a PublicKey flag with specified name and usage string.
// The arguments pk points to a PublicKey variable in which to store the interpreted values of the flag.
func PublicKeyFlagVar(f *pflag.FlagSet, pk *rivinetypes.PublicKey, name string, usage string) {
	f.Var(&publicKeyFlag{publicKey: pk}, name, usage)
}

// PublicKeyFlagVarP defines a PublicKey flag with specified name, shorthand and usage string.
// The arguments pk points to a PublicKey variable in which to store the interpreted values of the flag.
func PublicKeyFlagVarP(f *pflag.FlagSet, pk *rivinetypes.PublicKey, name, shorthand string, usage string) {
	f.VarP(&publicKeyFlag{publicKey: pk}, name, shorthand, usage)
}

type publicKeyFlag struct {
	publicKey *rivinetypes.PublicKey
	changed   bool
}

// Set implements pflag.Value.Set
func (flag *publicKeyFlag) Set(val string) error {
	if !flag.changed {
		flag.publicKey = new(rivinetypes.PublicKey)
		flag.changed = true
	}
	return flag.publicKey.LoadString(val)
}

// Type implements pflag.Value.Type
func (flag *publicKeyFlag) Type() string {
	return "PublicKeyFlag"
}

// String implements pflag.Value.String
func (flag *publicKeyFlag) String() string {
	return flag.publicKey.String()
}
