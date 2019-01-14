package main

// ERC20NodeValidatorConfig is all info required to create a ERC20NodeValidator.
// See the `ERC20NodeValidator` struct for more information.
type ERC20NodeValidatorConfig struct {
	Enabled     bool
	NetworkName string
	DataDir     string
	Port        int
}
