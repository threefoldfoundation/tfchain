package eth

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
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
// The default bootnodes can be overridden by passing a non nil or empty bootnodes parameter
func (config NetworkConfiguration) GetBootnodes(bootnodes []string) ([]*enode.Node, error) {
	if bootnodes == nil || len(bootnodes) == 0 {
		bootnodes = config.bootnodes
	}
	var nodes []*enode.Node
	for _, boot := range bootnodes {
		if url, err := enode.ParseV4(boot); err == nil {
			nodes = append(nodes, url)
		} else {
			err = errors.New("Failed to parse bootnode URL" + "url" + boot + "err" + err.Error())
			return nil, err
		}
	}
	return nodes, nil
}

var ethNetworkConfigurations = map[string]NetworkConfiguration{
	"main": {
		1,
		"main",
		core.DefaultGenesisBlock(),
		//Todo: replace with actual address
		common.HexToAddress("0x21826CC49B92029553af86F4e7A62C427E61e53a"),
		params.MainnetBootnodes,
	},
	"ropsten": {
		3,
		"ropsten",
		core.DefaultTestnetGenesisBlock(),
		common.HexToAddress("0xb821227dBa4Ef9585D31aa494406FD5E47a3db37"),
		params.TestnetBootnodes,
	},
	"rinkeby": {
		4,
		"rinkeby",
		core.DefaultRinkebyGenesisBlock(),
		common.HexToAddress("0x3bb58ffA340861b2Bac19c8b18262375F68c0AA5"),
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
