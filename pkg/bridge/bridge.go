package bridge

import (
	"fmt"
	"math/big"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lunny/log"
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

	bcInfo   types.BlockchainInfo
	chainCts types.ChainConstants

	bridgeContract *bridgeContract

	mut sync.Mutex
}

// New creates a new Bridge.
func New(cs modules.ConsensusSet, txdb *persist.TransactionDB, tp modules.TransactionPool, ethPort uint16, accountJSON, accountPass string, ethNetworkName string, datadir string, bcInfo types.BlockchainInfo, chainCts types.ChainConstants, cancel <-chan struct{}) (*Bridge, error) {

	contract, err := newBridgeContract(ethNetworkName, int(ethPort), accountJSON, accountPass, filepath.Join(datadir, "eth"))
	if err != nil {
		return nil, err
	}

	bridge := &Bridge{
		cs:             cs,
		txdb:           txdb,
		tp:             tp,
		bcInfo:         bcInfo,
		chainCts:       chainCts,
		bridgeContract: contract,
	}

	err = cs.ConsensusSetSubscribe(bridge, txdb.GetLastConsensusChangeID(), cancel)
	if err != nil {
		return nil, fmt.Errorf("bridged: failed to subscribe to consensus set: %v", err)
	}

	go bridge.bridgeContract.loop()
	go bridge.bridgeContract.subscribeTransfers()
	go bridge.bridgeContract.subscribeMint()
	go bridge.bridgeContract.subscribeRegisterWithdrawAddress()

	withdrawChan := make(chan withdrawEvent)
	go bridge.bridgeContract.subscribeWithdraw(withdrawChan)
	go func() {
		for {
			we := <-withdrawChan
			uh, found, err := txdb.GetTFTAddressForERC20Address(tfchaintypes.ERC20Address(we.receiver))
			if err != nil {
				log.Error("Retireving TFT address for registered ERC20 address errored: ", err)
				return
			}
			if !found {
				log.Error("Failed to retrieve TFT address for registered ERC20 Withdrawal address")
				return
			}

			tx := tfchaintypes.ERC20CoinCreationTransaction{}
			tx.Address = uh
			tx.Value = types.NewCurrency(we.amount)
			tx.TransactionID = tfchaintypes.ERC20TransactionID(we.txHash)
			tx.TransactionFee = types.NewCurrency(OneToken)
			if err := tp.AcceptTransactionSet([]types.Transaction{tx.Transaction()}); err != nil {
				log.Error("Failed to push ERC20 -> TFT transaction", "err", err)
				return
			}
			log.Info("Created ERC20 -> TFT transaction", "txid", tx.Transaction().ID())
		}
	}()

	return bridge, nil
}

// Close bridged.
func (bridged *Bridge) Close() {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()
	bridged.bridgeContract.close()
	bridged.cs.Unsubscribe(bridged)
}

// ProcessConsensusChange implements modules.ConsensusSetSubscriber,
// used to apply/revert blocks.
func (bridged *Bridge) ProcessConsensusChange(css modules.ConsensusChange) {
	bridged.mut.Lock()
	defer bridged.mut.Unlock()

	// TODO: add delay

	for _, block := range css.AppliedBlocks {
		for _, tx := range block.Transactions {
			if tx.Version == tfchaintypes.TransactionVersionERC20Conversion {
				log.Warn("Found convert transacton")
				txConvert, err := tfchaintypes.ERC20ConvertTransactionFromTransaction(tx)
				if err != nil {
					log.Error("Found a TFT convert transaction version, but can't create a conversion transaction from it")
					return
				}
				// Send the mint transaction, this requires gas
				if err = bridged.Mint(txConvert.Address, txConvert.Value, tx.ID()); err != nil {
					log.Error("Failed to push mint transaction", "error", err)
					return
				}
				log.Info("Created mint transaction on eth network")
			} else if tx.Version == tfchaintypes.TransactionVersionERC20AddressRegistration {
				log.Warn("Found erc20 address registration")
				txRegistration, err := tfchaintypes.ERC20AddressRegistrationTransactionFromTransaction(tx)
				if err != nil {
					log.Error("Found a TFT ERC20 Address registration transaction version, but can't create the right transaction for it")
					return
				}
				// send the address registration transaction
				if err = bridged.RegisterWithdrawalAddress(txRegistration.PublicKey); err != nil {
					log.Error("Failed to push withdrawal address registration transaction", "err", err)
					return
				}
				log.Info("Registered withdrawal address on eth network")
			}
		}
	}
}

var (
	// 18 digit precision
	precision = big.NewInt(1000000000000000000)
)

func (bridge *Bridge) Mint(receiver tfchaintypes.ERC20Address, amount types.Currency, txID types.TransactionID) error {
	return bridge.bridgeContract.mint(common.Address(receiver), big.NewInt(0).Mul(amount.Big(), precision), txID.String())
}

func (bridge *Bridge) RegisterWithdrawalAddress(key types.PublicKey) error {
	// convert public key to unlockhash to eth address
	erc20addr := tfchaintypes.ERC20AddressFromUnlockHash(types.NewPubKeyUnlockHash(key))
	return bridge.bridgeContract.registerWithdrawalAddress(common.Address(erc20addr))
}
