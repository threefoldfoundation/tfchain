package consensus

import (
	"fmt"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/modules/consensus"
	"github.com/threefoldtech/rivine/types"
)

func GetStandardTransactionValidators() []modules.TransactionValidationFunction {
	return []modules.TransactionValidationFunction{
		consensus.ValidateTransactionFitsInABlock,
		consensus.ValidateTransactionArbitraryData,
		consensus.ValidateCoinInputsAreValid,
		consensus.ValidateCoinOutputsAreValid,
		consensus.ValidateBlockStakeInputsAreValid,
		consensus.ValidateBlockStakeOutputsAreValid,
		consensus.ValidateDoubleCoinSpends,
		consensus.ValidateDoubleBlockStakeSpends,
		consensus.ValidateCoinInputsAreFulfilled,
		consensus.ValidateBlockStakeInputsAreFulfilled,
	}
}

func GetStandardTransactionVersionMappedValidators() map[types.TransactionVersion][]modules.TransactionValidationFunction {
	const (
		secondsInOneDay                         = 86400 + config.StandardNetworkBlockFrequency // round up
		daysFromStartOfBlockchainUntil2ndOfJuly = 74
		txnFeeCheckBlockHeight                  = daysFromStartOfBlockchainUntil2ndOfJuly *
			(secondsInOneDay / config.StandardNetworkBlockFrequency)
		blockHeightSinceLegacyTransactionsAreDisabled = 385000
	)
	validator := &MinimumMinerFeeValidator{MinimumBlockHeight: txnFeeCheckBlockHeight}
	legacyValidator := &DisableTransactionSinceValidator{MinimumBlockHeight: blockHeightSinceLegacyTransactionsAreDisabled}
	return map[types.TransactionVersion][]modules.TransactionValidationFunction{
		types.TransactionVersionZero: {
			validator.Validate,
			legacyValidator.Validate,
			consensus.ValidateCoinOutputsAreBalanced,
			consensus.ValidateBlockStakeOutputsAreBalanced,
		},
		types.TransactionVersionOne: {
			validator.Validate,
			consensus.ValidateCoinOutputsAreBalanced,
			consensus.ValidateBlockStakeOutputsAreBalanced,
		},
	}
}

func GetTestnetTransactionValidators() []modules.TransactionValidationFunction {
	return []modules.TransactionValidationFunction{
		consensus.ValidateTransactionFitsInABlock,
		consensus.ValidateTransactionArbitraryData,
		consensus.ValidateCoinInputsAreValid,
		consensus.ValidateCoinOutputsAreValid,
		consensus.ValidateBlockStakeInputsAreValid,
		consensus.ValidateBlockStakeOutputsAreValid,
		consensus.ValidateDoubleCoinSpends,
		consensus.ValidateDoubleBlockStakeSpends,
		consensus.ValidateCoinInputsAreFulfilled,
		consensus.ValidateBlockStakeInputsAreFulfilled,
	}
}

func GetTestnetTransactionVersionMappedValidators() map[types.TransactionVersion][]modules.TransactionValidationFunction {
	const (
		secondsInOneDay                         = 86400 + config.TestNetworkBlockFrequency // round up
		daysFromStartOfBlockchainUntil2ndOfJuly = 90
		txnFeeCheckBlockHeight                  = daysFromStartOfBlockchainUntil2ndOfJuly *
			(secondsInOneDay / config.TestNetworkBlockFrequency)
		blockHeightSinceLegacyTransactionsAreDisabled = 385000
	)
	validator := &MinimumMinerFeeValidator{MinimumBlockHeight: txnFeeCheckBlockHeight}
	legacyValidator := &DisableTransactionSinceValidator{MinimumBlockHeight: blockHeightSinceLegacyTransactionsAreDisabled}
	return map[types.TransactionVersion][]modules.TransactionValidationFunction{
		types.TransactionVersionZero: {
			validator.Validate,
			legacyValidator.Validate,
			consensus.ValidateCoinOutputsAreBalanced,
			consensus.ValidateBlockStakeOutputsAreBalanced,
		},
		types.TransactionVersionOne: {
			validator.Validate,
			consensus.ValidateCoinOutputsAreBalanced,
			consensus.ValidateBlockStakeOutputsAreBalanced,
		},
	}
}

func GetDevnetTransactionValidators() []modules.TransactionValidationFunction {
	return []modules.TransactionValidationFunction{
		consensus.ValidateTransactionFitsInABlock,
		consensus.ValidateTransactionArbitraryData,
		consensus.ValidateCoinInputsAreValid,
		consensus.ValidateCoinOutputsAreValid,
		consensus.ValidateBlockStakeInputsAreValid,
		consensus.ValidateBlockStakeOutputsAreValid,
		consensus.ValidateMinerFeeIsPresent,
		consensus.ValidateMinerFeesAreValid,
		consensus.ValidateDoubleCoinSpends,
		consensus.ValidateDoubleBlockStakeSpends,
		consensus.ValidateCoinInputsAreFulfilled,
		consensus.ValidateBlockStakeInputsAreFulfilled,
	}
}

func GetDevnetTransactionVersionMappedValidators() map[types.TransactionVersion][]modules.TransactionValidationFunction {
	return map[types.TransactionVersion][]modules.TransactionValidationFunction{
		types.TransactionVersionZero: {
			consensus.ValidateInvalidByDefault, // no longer required to exist on devnet
		},
		types.TransactionVersionOne: {
			consensus.ValidateCoinOutputsAreBalanced,
			consensus.ValidateBlockStakeOutputsAreBalanced,
		},
	}
}

// MinimumMinerFeeValidator is a validator which allows to check
// the minimum miner fees (and whether they are available) only since a specific (block) height.
type MinimumMinerFeeValidator struct {
	MinimumBlockHeight types.BlockHeight
}

// Validate is a validator function that checks if all miner fees are valid.
// Until the minimum block height 0 fees are allowed, afterwards the minimum fee is checked
func (validator *MinimumMinerFeeValidator) Validate(tx modules.ConsensusTransaction, ctx types.TransactionValidationContext) error {
	if tx.BlockHeight < validator.MinimumBlockHeight {
		// no need to check
		return nil
	}
	if ctx.IsBlockCreatingTx {
		return nil // validation does not apply to to block creation tx
	}
	if len(tx.MinerFees) == 0 {
		return fmt.Errorf("tx %s does not contain any miner fees while at least one was expected", tx.ID().String())
	}
	for _, fee := range tx.MinerFees {
		if fee.Cmp(ctx.MinimumMinerFee) == -1 {
			return types.ErrTooSmallMinerFee
		}
	}
	return nil
}

// DisableTransactionSinceValidator is a validator which allows to
// no longer allow a transaction since a specific block height
type DisableTransactionSinceValidator struct {
	MinimumBlockHeight types.BlockHeight
}

// Validate is a validator function that checks if the block is still allowed
// in the current chain.
func (validator *DisableTransactionSinceValidator) Validate(tx modules.ConsensusTransaction, ctx types.TransactionValidationContext) error {
	if tx.BlockHeight < validator.MinimumBlockHeight {
		// no need to check
		return nil
	}
	return fmt.Errorf("transaction is no longer allowed since block height %d", tx.BlockHeight)
}
