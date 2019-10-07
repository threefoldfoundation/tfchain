// +build !noeth

package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	tfeth "github.com/threefoldtech/rivine-extension-erc20/api"
	erc20bridge "github.com/threefoldtech/rivine-extension-erc20/api/bridge"
	"github.com/threefoldtech/rivine-extension-erc20/api/bridge/contract"
	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
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
// `github.com/threefoldtech/rivine-extension-erc20/api/bridge` for more info.
type ERC20NodeValidator struct {
	contract *erc20bridge.BridgeContract
	lc       *erc20bridge.LightClient
	abi      abi.ABI
}

// NewERC20NodeValidator creates a new INFURA-based ERC20NodeValidator.
// See the `ERC20NodeValidator` struct description for more information.
//
// If the cfg.Enabled property is False the tfchain `NopERC20TransactionValidator` implementation
// will be used and returned instead.
func NewERC20NodeValidator(cfg ERC20NodeValidatorConfig, cancel <-chan struct{}) (erc20types.ERC20TransactionValidator, error) {
	if !cfg.Enabled {
		return erc20types.NopERC20TransactionValidator{}, nil
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

	contract, err := erc20bridge.NewBridgeContract(netcfg.NetworkName, cfg.BootNodes, netcfg.ContractAddress.Hex(), cfg.Port, "", "", path.Join(cfg.DataDir, "lightnode"), cancel)
	if err != nil {
		return nil, err
	}
	return &ERC20NodeValidator{
		lc:       contract.LightClient(),
		abi:      abi,
		contract: contract,
	}, nil
}

func NewERC20NodeValidatorFromBridgeContract(contract *erc20bridge.BridgeContract) (erc20types.ERC20TransactionValidator, error) {
	if contract == nil {
		return nil, errors.New("no bridge contract is given, while a non-nil one is required for this validator constructor")
	}
	return &ERC20NodeValidator{
		contract: contract,
		lc:       contract.LightClient(),
		abi:      contract.ABI(),
	}, nil
}

// ValidateWithdrawTx implements ERC20TransactionValidator.ValidateWithdrawTx
func (ev *ERC20NodeValidator) ValidateWithdrawTx(_blockID, txID erc20types.ERC20Hash, expectedAddress erc20types.ERC20Address, expectedAmount types.Currency) error {
	withdraws, err := ev.contract.GetPastWithdraws(0, nil)
	for erc20bridge.IsNoPeerErr(err) {
		time.Sleep(time.Second * 5)
		log.Debug("Retrying to get past withdraws from peers")
		withdraws, err = ev.contract.GetPastWithdraws(0, nil)
	}
	if err != nil {
		return err
	}
	found := false
	var blockHash common.Hash
	for _, w := range withdraws {
		// looks like we found our transaction
		if w.TxHash() == common.Hash(txID) {
			found = true
			if (_blockID != erc20types.ERC20Hash{}) && common.Hash(_blockID) != w.BlockHash() {
				// IF a blockID is given, check if its the same. It might be different in case of a fork,
				// if so just add a statement in the logs.
				log.Info("Withdraw tx found in different block then specified", "expected", _blockID, "got", w.BlockHash().Hex())
			}
			if common.Address(expectedAddress) != w.Receiver() {
				return fmt.Errorf("Withdraw tx validation failed: invalid receiving address. Want address %s, got address %s", w.Receiver().Hex(), common.Address(expectedAddress).Hex())
			}
			if expectedAmount.Cmp(types.NewCurrency(w.Amount())) != 0 {
				return fmt.Errorf("Withdraw tx validation failed: invalid amount. Want %s, got %s", w.Amount().String(), expectedAmount.String())
			}
			// all event validations succeeded
			// remember block hash from the withdraw event so we can look up the
			// tx to check if it is old enough
			blockHash = w.BlockHash()
			break
		}
	}
	if !found {
		return fmt.Errorf("Withdraw tx validation failed: no matching withdraw event found - invalid tx ID %s", common.Hash(txID).Hex())
	}

	// Get the transaction
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, confirmations, err := ev.lc.FetchTransaction(ctx, blockHash, common.Hash(txID))
	// If we have no peers we can't verify and thus not continue syncing, so keep retrying
	for erc20bridge.IsNoPeerErr(err) {
		// wait 5 seconds before retrying
		time.Sleep(time.Second * 5)
		log.Debug("Retrying transaction fetch", "blockID", blockHash.Hex(), "txID", txID.String())
		_, confirmations, err = ev.lc.FetchTransaction(ctx, blockHash, common.Hash(txID))
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
func (ev *ERC20NodeValidator) GetStatus() (*erc20types.ERC20SyncStatus, error) {
	return ev.lc.GetStatus()
}

// GetBalanceInfo implements ERC20TransactionValidator.GetBalanceInfo
func (ev *ERC20NodeValidator) GetBalanceInfo() (*erc20types.ERC20BalanceInfo, error) {
	return ev.lc.GetBalanceInfo()
}

// Wait implements ERC20TransactionValidator.Wait
func (ev *ERC20NodeValidator) Wait(ctx context.Context) error {
	return ev.lc.Wait(ctx)
}
