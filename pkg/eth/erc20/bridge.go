package erc20

import (
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/threefoldfoundation/tfchain/pkg/persist"
	tfchaintypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/types"
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
	contract, err := newBridgeContract(ethNetworkName, int(ethPort), accountJSON, accountPass, filepath.Join(datadir, "eth"))
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

	go bridge.bridgeContract.loop()
	go bridge.bridgeContract.subscribeTransfers()
	go bridge.bridgeContract.subscribeMint()
	go bridge.bridgeContract.subscribeRegisterWithdrawAddress()

	withdrawChan := make(chan withdrawEvent)
	go bridge.bridgeContract.subscribeWithdraw(withdrawChan, 3544963)
	go func() {
		for {
			we := <-withdrawChan
			uh, found, err := txdb.GetTFTAddressForERC20Address(tfchaintypes.ERC20Address(we.receiver))
			if err != nil {
				log.Error(fmt.Sprintf("Retrieving TFT address for registered ERC20 address %v errored: %v", we.receiver, err))
				return
			}
			if !found {
				log.Error(fmt.Sprintf("Failed to retrieve TFT address for registered ERC20 Withdrawal address %v", we.receiver))
				return
			}

			tx := tfchaintypes.ERC20CoinCreationTransaction{}
			tx.Address = uh

			// calculate the amount of tokens we need to hand out
			erc20Tokens := big.NewInt(0).Div(we.amount, erc20Precision)
			tfTokens := big.NewInt(0).Mul(erc20Tokens, tftPrecision)

			// define the bridgeFee, txFee
			tx.TransactionFee = chainCts.MinimumTransactionFee
			tx.BridgeFee = chainCts.CurrencyUnits.OneCoin.Mul64(tfchaintypes.HardcodedERC20BridgeFeeOneCoinMultiplier)

			// define the value, which is the value withdrawn minus the fees
			tx.Value = types.NewCurrency(tfTokens).Sub(tx.TransactionFee).Sub(tx.BridgeFee)

			// fill in the other info
			tx.TransactionID = tfchaintypes.ERC20TransactionID(we.txHash)
			if err := tp.AcceptTransactionSet([]types.Transaction{tx.Transaction()}); err != nil {
				log.Error("Failed to push ERC20 -> TFT transaction", "err", err)
				return
			}
			log.Info("Created ERC20 -> TFT transaction", "txid", tx.Transaction().ID())
		}
	}()

	return bridge, nil
}

// Close bridge
func (bridge *Bridge) Close() {
	bridge.mut.Lock()
	defer bridge.mut.Unlock()
	bridge.bridgeContract.close()
	bridge.cs.Unsubscribe(bridge)
}

var (
	// 18 digit precision
	erc20Precision = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)
	// 9 digit precision
	tftPrecision = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(9), nil)
)

func (bridge *Bridge) mint(receiver tfchaintypes.ERC20Address, amount types.Currency, txID types.TransactionID) error {
	tfTokens := big.NewInt(0).Div(amount.Big(), tftPrecision)
	erc20Tokens := big.NewInt(0).Mul(tfTokens, erc20Precision)
	return bridge.bridgeContract.mint(common.Address(receiver), erc20Tokens, txID.String())
}

func (bridge *Bridge) registerWithdrawalAddress(key types.PublicKey) error {
	// convert public key to unlockhash to eth address
	erc20addr := tfchaintypes.ERC20AddressFromUnlockHash(types.NewPubKeyUnlockHash(key))
	return bridge.bridgeContract.registerWithdrawalAddress(common.Address(erc20addr))
}
