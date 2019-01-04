package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	tfeth "github.com/threefoldfoundation/tfchain/pkg/eth"
	"github.com/threefoldfoundation/tfchain/pkg/eth/erc20"
	"github.com/threefoldfoundation/tfchain/pkg/eth/erc20/contract"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/build"

	"github.com/threefoldtech/rivine/types"
)

// ERC20NodeValidator implements the ERC20TransactionValidator,
// getting the transactions using the INFURA service.
//
// Goal is to move away from this ASAP and be able to fetch transactions using a Light client,
// problem is however that the LES/v2 protocol is not implemented yet by default in go-ethereum,
// it is however implemented server-side. So should we want, we can fork Ethereum,
// and contribute the client-side calls we care about as to be able to do it all from
// a light client.
type ERC20NodeValidator struct {
	lc  *erc20.LightClient
	abi abi.ABI
}

// ERC20NodeValidatorConfig is all info required to create a ERC20NodeValidator.
// See the `ERC20NodeValidator` struct for more information.
type ERC20NodeValidatorConfig struct {
	Enabled     bool
	NetworkName string
	DataDir     string
	Port        int
}

// NewERC20NodeValidator creates a new INFURA-based ERC20NodeValidator.
// See the `ERC20NodeValidator` struct description for more information.
//
// If the cfg.Enabled property is False the tfchain `NopERC20TransactionValidator` implementation
// will be used and returned instead.
func NewERC20NodeValidator(cfg ERC20NodeValidatorConfig, cancel <-chan struct{}) (tftypes.ERC20TransactionValidator, error) {
	if !cfg.Enabled {
		return tftypes.NopERC20TransactionValidator{}, nil
	}

	// Create the persistent dir if it doesn't exist already
	err := os.MkdirAll(cfg.DataDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while creating the persistent (data) dir: %v", err)
	}

	// Define the Ethereum Logger,
	// logging both to a file and the STDERR, with a lower verbosity for the latter.
	ethLogFmtr := log.TerminalFormat(true)
	ethLogFileHandler, err := log.FileHandler(path.Join(cfg.DataDir, "node.log"), ethLogFmtr)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while ETH file-logger: %v", err)
	}
	ethFileLogLvl, ethStreamLogLvl := log.LvlInfo, log.LvlWarn
	if build.DEBUG {
		ethFileLogLvl, ethStreamLogLvl = log.LvlDebug, log.LvlInfo
	}
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.Lvl(ethFileLogLvl), ethLogFileHandler),
		log.LvlFilterHandler(log.Lvl(ethStreamLogLvl), log.StreamHandler(os.Stderr, ethLogFmtr))))

	// parse the ERC20 smart contract
	abi, err := abi.JSON(strings.NewReader(contract.TTFT20ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while parsing contract ABI: %v", err)
	}

	// get the ETH network config
	netcfg, err := tfeth.GetEthNetworkConfiguration(cfg.NetworkName)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while fetching the ETH network config: %v", err)
	}

	// get the ETH bootstrap nodes
	bootstrapNodes, err := netcfg.GetBootnodes()
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while fetching the ETH bootstrap node info: %v", err)
	}

	// create the ethereum light client
	lc, err := erc20.NewLightClient(erc20.LightClientConfig{
		Port:           cfg.Port,
		DataDir:        path.Join(cfg.DataDir, "lightnode"),
		BootstrapNodes: bootstrapNodes,
		NetworkID:      netcfg.NetworkID,
		NetworkName:    netcfg.NetworkName,
		GenesisBlock:   netcfg.GenesisBlock,
	}, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while creating the light client: %v", err)
	}

	return &ERC20NodeValidator{
		lc:  lc,
		abi: abi,
	}, nil
}

// ValidateWithdrawTx implements ERC20TransactionValidator.ValidateWithdrawTx
func (ev *ERC20NodeValidator) ValidateWithdrawTx(blockID, txID tftypes.ERC20Hash, expectedAddress tftypes.ERC20Address, expecedAmount types.Currency) error {
	// Get the transaction
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	tx, err := ev.lc.FetchTransaction(ctx, common.Hash(blockID), common.Hash(txID))
	if err != nil {
		return fmt.Errorf("failed to fetch ERC20 Tx: %v", err)
	}

	// Extract the data
	// (and as a necessary step also validate if the Tx and input are of the correct type)
	txData := tx.Data()
	if len(txData) <= 4 {
		return fmt.Errorf("invalid ERC20 Tx: unexpected Tx data length: %v", len(txData))
	}

	// first 4 bytes contain the id, so let's get method using that ID
	method, err := ev.abi.MethodById(txData[:4])
	if err != nil {
		return fmt.Errorf("invalid ERC20 Tx: failed to get method using its parsed id: %v", err)
	}
	if method.Name != "transfer" {
		return fmt.Errorf("invalid ERC20 Tx: unexpected name for unpacked method ID: %s", method.Name)
	}
	// unpack the input into a struct we can work with
	params := struct {
		To     common.Address
		Tokens *big.Int
	}{}
	err = method.Inputs.Unpack(&params, txData[4:])
	if err != nil {
		return fmt.Errorf("error while unpacking transfer ERC20 Tx %v input: %v", txID, err)
	}

	// validate if the address is correct
	toAddress := params.To.String()
	expectedToAddress := common.Address(expectedAddress).String()
	if toAddress != expectedToAddress {
		return fmt.Errorf("unexpected to address %v: expected address %v", toAddress, expectedToAddress)
	}

	// validate if the amount of tokens withdrawn is correct
	amount := types.NewCurrency(params.Tokens)
	if amount.Equals(expecedAmount) {
		return fmt.Errorf("unexpected transferred TFT value %v: expected value %v",
			amount.String(), expecedAmount.String())
	}

	// all is good, return nil to indicate this
	return nil
}
