package erc20

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/threefoldfoundation/tfchain/pkg/persist"
	tfchaintypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/types"
)

const (
	// blockDelay is the amount of blocks to wait before
	// pushing tft transactions to the ethereum contract,
	// and to push ethereum transactions to the TF chain
	blockDelay = 6
)

// Bridge is a high lvl structure which listens on contract events and bridge-related
// tfchain transactions, and handles them
type Bridge struct {
	cs   modules.ConsensusSet
	txdb *persist.TransactionDB
	tp   modules.TransactionPool

	persistDir string
	persist    persistence
	buffer     *blockBuffer

	bcInfo   types.BlockchainInfo
	chainCts types.ChainConstants

	bridgeContract *bridgeContract

	mut sync.Mutex
}

// NewBridge creates a new Bridge.
func NewBridge(cs modules.ConsensusSet, txdb *persist.TransactionDB, tp modules.TransactionPool, ethPort uint16, accountJSON, accountPass string, ethNetworkName string, datadir string, bcInfo types.BlockchainInfo, chainCts types.ChainConstants, cancel <-chan struct{}) (*Bridge, error) {
	contract, err := newBridgeContract(ethNetworkName, int(ethPort), accountJSON, accountPass, filepath.Join(datadir, "eth"), cancel)
	if err != nil {
		return nil, err
	}

	bridge := &Bridge{
		cs:             cs,
		txdb:           txdb,
		tp:             tp,
		persistDir:     datadir,
		bcInfo:         bcInfo,
		chainCts:       chainCts,
		bridgeContract: contract,
	}

	err = bridge.initPersist()
	if err != nil {
		return nil, errors.New("bridge persistence startup failed: " + err.Error())
	}

	bridge.buffer = newBlockBuffer(blockDelay)

	err = cs.ConsensusSetSubscribe(bridge, bridge.persist.RecentChange, cancel)
	if err != nil {
		return nil, fmt.Errorf("bridged: failed to subscribe to consensus set: %v", err)
	}

	heads := make(chan *ethtypes.Header)

	go bridge.bridgeContract.loop(heads)

	// subscribing to these events is not needed for operational purposes, but might be nice to get some info
	go bridge.bridgeContract.subscribeTransfers()
	go bridge.bridgeContract.subscribeMint()
	go bridge.bridgeContract.subscribeRegisterWithdrawAddress()

	withdrawChan := make(chan withdrawEvent)
	go bridge.bridgeContract.subscribeWithdraw(withdrawChan, bridge.persist.EthHeight)
	go func() {
		txMap := make(map[tfchaintypes.ERC20Hash]withdrawEvent)
		for {
			select {
			// Remember new withdraws
			case we := <-withdrawChan:
				// Check if the withdraw is valid
				_, found, err := txdb.GetTFTAddressForERC20Address(tfchaintypes.ERC20Address(we.receiver))
				if err != nil {
					log.Error(fmt.Sprintf("Retrieving TFT address for registered ERC20 address %v errored: %v", we.receiver, err))
					continue
				}
				if !found {
					log.Error(fmt.Sprintf("Failed to retrieve TFT address for registered ERC20 Withdrawal address %v", we.receiver))
					continue
				}
				// remember the withdraw
				txMap[tfchaintypes.ERC20Hash(we.txHash)] = we

			// If we get a new head, check every withdraw we have to see if it has matured
			case head := <-heads:
				for id := range txMap {
					we := txMap[id]
					if head.Number.Uint64() >= we.blockHeight+blockDelay {
						// we waited long enough, create transaction and push it
						uh, found, err := txdb.GetTFTAddressForERC20Address(tfchaintypes.ERC20Address(we.receiver))
						if err != nil {
							log.Error(fmt.Sprintf("Retrieving TFT address for registered ERC20 address %v errored: %v", we.receiver, err))
							continue
						}
						if !found {
							log.Error(fmt.Sprintf("Failed to retrieve TFT address for registered ERC20 Withdrawal address %v", we.receiver))
							continue
						}

						tx := tfchaintypes.ERC20CoinCreationTransaction{}
						tx.Address = uh

						// define the bridgeFee, txFee
						tx.TransactionFee = chainCts.MinimumTransactionFee
						tx.BridgeFee = chainCts.CurrencyUnits.OneCoin.Mul64(tfchaintypes.HardcodedERC20BridgeFeeOneCoinMultiplier)

						// define the value, which is the value withdrawn minus the fees
						tx.Value = types.NewCurrency(we.amount).Sub(tx.TransactionFee).Sub(tx.BridgeFee)

						// fill in the other info
						tx.TransactionID = tfchaintypes.ERC20Hash(we.txHash)
						tx.BlockID = tfchaintypes.ERC20Hash(we.blockHash)

						if err := bridge.commitWithdrawTransaction(tx); err != nil {
							log.Error("Failed ot create ERC20 Withdraw transaction", "err", err)
							continue
						}

						log.Info("Created ERC20 -> TFT transaction", "txid", tx.Transaction().ID())

						// forget about our tx
						delete(txMap, id)
					}
				}

				bridge.persist.EthHeight = head.Number.Uint64() - blockDelay
				// Check for underflow
				if bridge.persist.EthHeight > head.Number.Uint64() {
					bridge.persist.EthHeight = 0
				}
				if err := bridge.save(); err != nil {
					log.Error("Failed to save bridge persistency", "err", err)
				}
			}
		}
	}()

	return bridge, nil
}

// commitWithdrawTransaction verifies and (if successfull) commits a withdraw transaction on
// the tfchain (thus creating new tokens)
func (bridge *Bridge) commitWithdrawTransaction(tx tfchaintypes.ERC20CoinCreationTransaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// verify that the transaction is still present in the block (no forks occurred)
	_, err := bridge.bridgeContract.lc.FetchTransaction(ctx, common.Hash(tx.BlockID), common.Hash(tx.TransactionID))
	if err != nil {
		return err
	}
	// accept transaction, ignore duplicate errors which might occur if we are syncing a new bridge
	if err := bridge.tp.AcceptTransactionSet([]types.Transaction{tx.Transaction()}); err != nil && err != modules.ErrDuplicateTransactionSet {
		return err
	}
	return nil
}

// Close bridge
func (bridge *Bridge) Close() {
	bridge.mut.Lock()
	defer bridge.mut.Unlock()
	bridge.bridgeContract.close()
	bridge.cs.Unsubscribe(bridge)
}

func (bridge *Bridge) mint(receiver tfchaintypes.ERC20Address, amount types.Currency, txID types.TransactionID) error {
	return bridge.bridgeContract.mint(common.Address(receiver), amount.Big(), txID.String())
}

func (bridge *Bridge) registerWithdrawalAddress(key types.PublicKey) error {
	// convert public key to unlockhash to eth address
	erc20addr := tfchaintypes.ERC20AddressFromUnlockHash(types.NewPubKeyUnlockHash(key))
	return bridge.bridgeContract.registerWithdrawalAddress(common.Address(erc20addr))
}
