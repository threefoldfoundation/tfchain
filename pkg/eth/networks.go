package eth

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/params"
)

//NetworkConfiguration defines the Ethereum network specific configuration needed by the bridge
type NetworkConfiguration struct {
	NetworkID       uint64
	NetworkName     string
	GenesisBlock    *core.Genesis
	ContractAddress common.Address
	bootnodes       []string
}

//GetBootnodes returns the bootnodes for the specific network as  slice of *discv5.Node
func (config NetworkConfiguration) GetBootnodes() ([]*discv5.Node, error) {
	var nodes []*discv5.Node
	for _, boot := range config.bootnodes {
		if url, err := discv5.ParseNode(boot); err == nil {
			nodes = append(nodes, url)
		} else {
			err = errors.New("Failed to parse bootnode URL" + "url" + boot + "err" + err.Error())
			return nil, err
		}
	}
	return nodes, nil
}

var ethNetworkConfigurations = map[string]NetworkConfiguration{
	"main": NetworkConfiguration{
		1,
		"main",
		core.DefaultGenesisBlock(),
		//Todo: replace with actual address
		common.HexToAddress("0x21826CC49B92029553af86F4e7A62C427E61e53a"),
		params.MainnetBootnodes,
	},
	"ropsten": NetworkConfiguration{
		3,
		"ropsten",
		core.DefaultTestnetGenesisBlock(),
		//Todo: replace with actual address
		common.HexToAddress("0x21826CC49B92029553af86F4e7A62C427E61e53a"),
		params.TestnetBootnodes,
	},
	"rinkeby": NetworkConfiguration{
		4,
		"rinkeby",
		core.DefaultRinkebyGenesisBlock(),
		common.HexToAddress("0xa99000aD598cdF135A23A10BFAb8a239AbC412a3"),
		params.RinkebyBootnodes,
	},
}

//GetEthNetworkConfiguration returns the EthNetworkConAfiguration for a specific network
func GetEthNetworkConfiguration(networkname string) (networkconfig NetworkConfiguration, err error) {
	networkconfig, found := ethNetworkConfigurations[networkname]
	if !found {
		err = fmt.Errorf("Ethereum network %s not supported", networkname)
	}
	return
}
