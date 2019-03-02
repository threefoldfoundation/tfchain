// +build !noeth

package main

import (
	"context"
	"fmt"
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

const (
	// MinimumERC20CoinCreationConfirmationsRequired defines the amount of minimum confirmations required,
	// in order for the ERC20 Node validator to accept a CoinCreation Tx, backed by an ERC20 Tx.
	MinimumERC20CoinCreationConfirmationsRequired = 25
)

// ERC20NodeValidator implements the ERC20TransactionValidator,
// getting the transactions using the LES/v2 protocol, see the
// `github.com/threefoldfoundation/tfchain/pkg/eth` for more info.
type ERC20NodeValidator struct {
	contract *erc20.BridgeContract
	lc       *erc20.LightClient
	abi      abi.ABI
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
	ethStreamLogLvl := log.LvlWarn
	ethFileLogLvl := log.Lvl(cfg.EthLogLevel)

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
	bootstrapNodes, err := netcfg.GetBootnodes(cfg.BootNodes)
	log.Info("bootnodes", "nodes", bootstrapNodes)

	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while fetching the ETH bootstrap node info: %v", err)
	}

	contract, err := erc20.NewBridgeContract(netcfg.NetworkName, cfg.BootNodes, netcfg.ContractAddress.Hex(), cfg.Port, "", "", path.Join(cfg.DataDir, "lightnode"), cancel)
	if err != nil {
		return nil, err
	}
	return &ERC20NodeValidator{
		lc:       contract.LightClient(),
		abi:      abi,
		contract: contract,
	}, nil
}

// ValidateWithdrawTx implements ERC20TransactionValidator.ValidateWithdrawTx
func (ev *ERC20NodeValidator) ValidateWithdrawTx(blockID, txID tftypes.ERC20Hash, expectedAddress tftypes.ERC20Address, expectedAmount types.Currency) error {
	withdraws, err := ev.contract.GetPastWithdraws(0, nil)
	for erc20.IsNoPeerErr(err) {
		time.Sleep(time.Second * 5)
		log.Debug("Retrying to get past withdraws from peers")
		withdraws, err = ev.contract.GetPastWithdraws(0, nil)
	}
	if err != nil {
		return err
	}
	found := false
	for _, w := range withdraws {
		// looks like we found our transaction
		if w.TxHash() == common.Hash(txID) {
			found = true
			if common.Hash(blockID) != w.BlockHash() {
				return fmt.Errorf("withdraw tx validation failed: invalid block ID. Want ID %s, got ID %s", w.BlockHash().Hex(), common.Hash(blockID).Hex())
			}
			if common.Address(expectedAddress) != w.Receiver() {
				return fmt.Errorf("Withdraw tx validation failed: invalid receiving address. Want address %s, got address %s", w.Receiver().Hex(), common.Address(expectedAddress).Hex())
			}
			if expectedAmount.Cmp(types.NewCurrency(w.Amount())) != 0 {
				return fmt.Errorf("Withdraw tx validation failed: invalid amount. Want %s, got %s", w.Amount().String(), expectedAmount.String())
			}
			// all event validations succeeded
			break
		}
	}
	if !found {
		return fmt.Errorf("Withdraw tx validation failed: no matching withdraw event found - invalid tx ID %s", common.Hash(txID).Hex())
	}

	// Get the transaction
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, confirmations, err := ev.lc.FetchTransaction(ctx, common.Hash(blockID), common.Hash(txID))
	// If we have no peers we can't verify and thus not continue syncing, so keep retrying
	for erc20.IsNoPeerErr(err) {
		// wait 5 seconds before retrying
		time.Sleep(time.Second * 5)
		log.Debug("Retrying transaction fetch", "blockID", blockID, "txID", txID)
		_, confirmations, err = ev.lc.FetchTransaction(ctx, common.Hash(blockID), common.Hash(txID))
	}
	if err != nil {
		return fmt.Errorf("failed to fetch ERC20 Tx: %v", err)
	}

	// Validate we have sufficient amount of confirmations available
	if confirmations < MinimumERC20CoinCreationConfirmationsRequired {
		return fmt.Errorf("invalid ERC20 Tx: insufficient block confirmations: %d", confirmations)
	}

	// all is good, return nil to indicate this
	return nil
}

// GetStatus implements ERC20TransactionValidator.GetStatus
func (ev *ERC20NodeValidator) GetStatus() (*tftypes.ERC20SyncStatus, error) {
	return ev.lc.GetStatus()
}

// GetBalanceInfo implements ERC20TransactionValidator.GetBalanceInfo
func (ev *ERC20NodeValidator) GetBalanceInfo() (*tftypes.ERC20BalanceInfo, error) {
	return ev.lc.GetBalanceInfo()
}

// Wait implements ERC20TransactionValidator.Wait
func (ev *ERC20NodeValidator) Wait(ctx context.Context) error {
	return ev.lc.Wait(ctx)
}
