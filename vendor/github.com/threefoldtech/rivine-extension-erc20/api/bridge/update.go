package bridge

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/threefoldtech/rivine/modules"

	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
)

// ProcessConsensusChange implements modules.ConsensusSetSubscriber,
// used to apply/revert blocks.
func (bridge *Bridge) ProcessConsensusChange(css modules.ConsensusChange) {
	bridge.mut.Lock()
	defer bridge.mut.Unlock()

	// In case there is a reverted block remove it from the buffer
	for range css.RevertedBlocks {
		bridge.buffer.rewindBlock()
	}
	for _, block := range css.AppliedBlocks {
		height, _ := bridge.cs.BlockHeightOfBlock(block)
		log.Debug("Processing TfChain block", "block", height)

		// Add block in buffer, and get the block which we held long enough to be confident it won't be rewinded
		oldBlock := bridge.buffer.pushBlock(block, css.ID)

		// If there was no block in the buffer yet don't do anything else
		if oldBlock == nil {
			continue
		}

		for _, tx := range oldBlock.Transactions {
			if tx.Version == bridge.txVersions.ERC20Conversion {
				log.Warn("Found convert transacton")
				txConvert, err := erc20types.ERC20ConvertTransactionFromTransaction(tx, bridge.txVersions.ERC20Conversion)
				if err != nil {
					log.Error("Found a TFT convert transaction version, but can't create a conversion transaction from it")
					return
				}
				// Send the mint transaction, this requires gas
				if err = bridge.mint(txConvert.Address, txConvert.Value, tx.ID()); err != nil {
					log.Error("Failed to push mint transaction", "error", err)
					return
				}
				log.Info("Created mint transaction on eth network")
			} else if tx.Version == bridge.txVersions.ERC20AddressRegistration {
				log.Warn("Found erc20 address registration")
				txRegistration, err := erc20types.ERC20AddressRegistrationTransactionFromTransaction(tx, bridge.txVersions.ERC20AddressRegistration)
				if err != nil {
					log.Error("Found a TFT ERC20 Address registration transaction version, but can't create the right transaction for it")
					return
				}
				// send the address registration transaction
				if err = bridge.registerWithdrawalAddress(txRegistration.PublicKey); err != nil {
					log.Error("Failed to push withdrawal address registration transaction", "err", err)
					return
				}
				log.Info("Registered withdrawal address on eth network")
			}
		}

		// update stats
		bridge.persist.Height++
		bridge.persist.RecentChange = oldBlock.ConsensusChangeID
		bridge.save()
	}
}
