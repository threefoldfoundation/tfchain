package types

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	tfencoding "github.com/threefoldfoundation/tfchain/pkg/encoding"

	"github.com/rivine/rivine/build"
	"github.com/rivine/rivine/crypto"
	"github.com/rivine/rivine/encoding"
	"github.com/rivine/rivine/types"
)

const (
	// TransactionVersionMinterDefinition defines the Transaction version
	// for a MinterDefinition Transaction.
	//
	// See the `MinterDefinitionTransactionController` and `MinterDefinitionTransaction`
	// types for more information.
	TransactionVersionMinterDefinition types.TransactionVersion = iota + 128
	// TransactionVersionCoinCreation defines the Transaction version
	// for a CoinCreation Transaction.
	//
	// See the `CoinCreationTransactionController` and `CoinCreationTransaction`
	// types for more information.
	TransactionVersionCoinCreation
)

const (
	// TransactionVersionBotRegistration defines the Transaction version
	// for a BotRegistration Transaction, used to register a new 3bot,
	// where new means that the used public key cannot yet exist.
	TransactionVersionBotRegistration types.TransactionVersion = iota + 144
	// TransactionVersionBotRecordUpdate defines the Transaction version
	// for a Tx used to update a 3bot Record by the owner. where owner
	// means the 3bot that created the record to be updated initially using the BotRegistration Tx.
	TransactionVersionBotRecordUpdate
	// TransactionVersionBotNameTransfer defines the Transaction version
	// for a Tx used to transfer one or multiple names from the active
	// 3bot that up to the point of that Tx to another 3bot.
	TransactionVersionBotNameTransfer
)

// These Specifiers are used internally when calculating a Transaction's ID.
// See Rivine's Specifier for more details.
var (
	SpecifierMintDefinitionTransaction  = types.Specifier{'m', 'i', 'n', 't', 'e', 'r', ' ', 'd', 'e', 'f', 'i', 'n', ' ', 't', 'x'}
	SpecifierCoinCreationTransaction    = types.Specifier{'c', 'o', 'i', 'n', ' ', 'm', 'i', 'n', 't', ' ', 't', 'x'}
	SpecifierBotRegistrationTransaction = types.Specifier{'b', 'o', 't', ' ', 'r', 'e', 'g', 'i', 's', 't', 'e', 'r', ' ', 't', 'x'}
	SpecifierBotRecordUpdateTransaction = types.Specifier{'b', 'o', 't', ' ', 'r', 'e', 'c', 'u', 'p', 'd', 'a', 't', 'e', ' ', 't', 'x'}
	SpecifierBotNameTransferTransaction = types.Specifier{'b', 'o', 't', ' ', 'n', 'a', 'm', 'e', 't', 'r', 'a', 'n', 's', ' ', 't', 'x'}
)

// Bot validation errors
var (
	ErrBotKeyAlreadyRegistered  = errors.New("bot key is already registered")
	ErrBotNameAlreadyRegistered = errors.New("bot name is already registered")
)

// TFChainReadDB is the Read-Only Database that is required in order to fetch the
// different transaction-related data from required by Tfchain transactions.
type TFChainReadDB interface {
	MintConditionGetter
	BotRecordReadRegistry
}

// RegisterTransactionTypesForStandardNetwork registers he transaction controllers
// for all transaction versions supported on the standard network.
func RegisterTransactionTypesForStandardNetwork(db TFChainReadDB, oneCoin types.Currency, cfg config.DaemonNetworkConfig) {
	const (
		secondsInOneDay                         = 86400 + config.StandardNetworkBlockFrequency // round up
		daysFromStartOfBlockchainUntil2ndOfJuly = 74
		txnFeeCheckBlockHeight                  = daysFromStartOfBlockchainUntil2ndOfJuly *
			(secondsInOneDay / config.StandardNetworkBlockFrequency)
	)
	// overwrite rivine-defined transaction versions
	types.RegisterTransactionVersion(types.TransactionVersionZero, LegacyTransactionController{
		LegacyTransactionController:    types.LegacyTransactionController{},
		TransactionFeeCheckBlockHeight: txnFeeCheckBlockHeight,
	})
	types.RegisterTransactionVersion(types.TransactionVersionOne, DefaultTransactionController{
		DefaultTransactionController:   types.DefaultTransactionController{},
		TransactionFeeCheckBlockHeight: txnFeeCheckBlockHeight,
	})

	// define tfchain-specific transaction versions

	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: db,
	})
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: db,
	})

	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry:              db,
		RegistryPoolCondition: cfg.FoundationPoolCondition,
		OneCoin:               oneCoin,
	})
}

// RegisterTransactionTypesForTestNetwork registers he transaction controllers
// for all transaction versions supported on the test network.
func RegisterTransactionTypesForTestNetwork(db TFChainReadDB, oneCoin types.Currency, cfg config.DaemonNetworkConfig) {
	const (
		secondsInOneDay                         = 86400 + config.TestNetworkBlockFrequency // round up
		daysFromStartOfBlockchainUntil2ndOfJuly = 90
		txnFeeCheckBlockHeight                  = daysFromStartOfBlockchainUntil2ndOfJuly *
			(secondsInOneDay / config.TestNetworkBlockFrequency)
	)
	// overwrite rivine-defined transaction versions
	types.RegisterTransactionVersion(types.TransactionVersionZero, LegacyTransactionController{
		LegacyTransactionController:    types.LegacyTransactionController{},
		TransactionFeeCheckBlockHeight: txnFeeCheckBlockHeight,
	})
	types.RegisterTransactionVersion(types.TransactionVersionOne, DefaultTransactionController{
		DefaultTransactionController:   types.DefaultTransactionController{},
		TransactionFeeCheckBlockHeight: txnFeeCheckBlockHeight,
	})

	// define tfchain-specific transaction versions

	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: db,
	})
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: db,
	})

	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry:              db,
		RegistryPoolCondition: cfg.FoundationPoolCondition,
		OneCoin:               oneCoin,
	})
}

// RegisterTransactionTypesForDevNetwork registers he transaction controllers
// for all transaction versions supported on the dev network.
func RegisterTransactionTypesForDevNetwork(db TFChainReadDB, oneCoin types.Currency, cfg config.DaemonNetworkConfig) {
	// overwrite rivine-defined transaction versions
	types.RegisterTransactionVersion(types.TransactionVersionZero, LegacyTransactionController{
		LegacyTransactionController:    types.LegacyTransactionController{},
		TransactionFeeCheckBlockHeight: 0,
	})
	types.RegisterTransactionVersion(types.TransactionVersionOne, DefaultTransactionController{
		DefaultTransactionController:   types.DefaultTransactionController{},
		TransactionFeeCheckBlockHeight: 0,
	})

	// define tfchain-specific transaction versions

	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: db,
	})
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: db,
	})

	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry:              db,
		RegistryPoolCondition: cfg.FoundationPoolCondition,
		OneCoin:               oneCoin,
	})
}

type (
	// MintConditionGetter allows you to get the mint condition at a given block height.
	//
	// For the daemon this interface could be implemented directly by the DB object
	// that keeps track of the mint condition state, while for a client this could
	// come via the REST API from a tfchain daemon in a more indirect way.
	MintConditionGetter interface {
		// GetActiveMintCondition returns the active active mint condition.
		GetActiveMintCondition() (types.UnlockConditionProxy, error)
		// GetMintConditionAt returns the mint condition at a given block height.
		GetMintConditionAt(height types.BlockHeight) (types.UnlockConditionProxy, error)
	}
)

type (
	// DefaultTransactionController wraps around Rivine's DefaultTransactionController,
	// as to ensure that we use check the MinimumTransactionFee,
	// only since a certain block height, and otherwise just ensure it is bigger than 0.
	//
	// In order to achieve this, the TransactionValidation interface is
	// implemented on top of the regular DefaultTransactionController.
	DefaultTransactionController struct {
		types.DefaultTransactionController
		TransactionFeeCheckBlockHeight types.BlockHeight
	}
	// LegacyTransactionController wraps around Rivine's LegacyTransactionController,
	// as to ensure that we use check the MinimumTransactionFee,
	// only since a certain block height, and otherwise just ensure it is bigger than 0.
	//
	// In order to achieve this, the TransactionValidation interface is
	// implemented on top of the regular LegacyTransactionController.
	LegacyTransactionController struct {
		types.LegacyTransactionController
		TransactionFeeCheckBlockHeight types.BlockHeight
	}

	// CoinCreationTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 129. It allows for the creation of Coin Outputs,
	// without requiring coin inputs, but can only be used by the defined Coin Minters.
	CoinCreationTransactionController struct {
		// MintConditionGetter is used to get a mint condition at the context-defined block height.
		//
		// The found MintCondition defines the condition that has to be fulfilled
		// in order to mint new coins into existence (in the form of non-backed coin outputs).
		MintConditionGetter MintConditionGetter
	}

	// MinterDefinitionTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 128. It allows the transfer of coin minting powers.
	MinterDefinitionTransactionController struct {
		// MintConditionGetter is used to get a mint condition at the context-defined block height.
		//
		// The found MintCondition defines the condition that has to be fulfilled
		// in order to mint new coins into existence (in the form of non-backed coin outputs).
		MintConditionGetter MintConditionGetter
	}
)

// ensure our controllers implement all desired interfaces
var (
	// ensure at compile time that DefaultTransactionController
	// implements the desired interfaces
	_ types.TransactionController = DefaultTransactionController{}
	_ types.TransactionValidator  = DefaultTransactionController{}

	// ensure at compile time that LegacyTransactionController
	// implements the desired interfaces
	_ types.TransactionController = LegacyTransactionController{}
	_ types.TransactionValidator  = LegacyTransactionController{}
	_ types.InputSigHasher        = LegacyTransactionController{}
	_ types.TransactionIDEncoder  = LegacyTransactionController{}

	// ensure at compile time that CoinCreationTransactionController
	// implements the desired interfaces
	_ types.TransactionController      = CoinCreationTransactionController{}
	_ types.TransactionExtensionSigner = CoinCreationTransactionController{}
	_ types.TransactionValidator       = CoinCreationTransactionController{}
	_ types.CoinOutputValidator        = CoinCreationTransactionController{}
	_ types.BlockStakeOutputValidator  = CoinCreationTransactionController{}
	_ types.InputSigHasher             = CoinCreationTransactionController{}
	_ types.TransactionIDEncoder       = CoinCreationTransactionController{}

	// ensure at compile time that MinterDefinitionTransactionController
	// implements the desired interfaces
	_ types.TransactionController      = MinterDefinitionTransactionController{}
	_ types.TransactionExtensionSigner = MinterDefinitionTransactionController{}
	_ types.TransactionValidator       = MinterDefinitionTransactionController{}
	_ types.CoinOutputValidator        = MinterDefinitionTransactionController{}
	_ types.BlockStakeOutputValidator  = MinterDefinitionTransactionController{}
	_ types.InputSigHasher             = MinterDefinitionTransactionController{}
	_ types.TransactionIDEncoder       = MinterDefinitionTransactionController{}
)

// DefaultTransactionController

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (dtc DefaultTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) error {
	if ctx.Confirmed && ctx.BlockHeight < dtc.TransactionFeeCheckBlockHeight {
		// as to ensure the miner fee is at least bigger than 0,
		// we however only want to put this restriction within the consensus set,
		// the stricter miner fee checks should apply immediately to the transaction pool logic
		constants.MinimumMinerFee = types.NewCurrency64(1)
	}
	return types.DefaultTransactionValidation(t, ctx, constants)
}

// LegacyTransactionController

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (ltc LegacyTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) error {
	if ctx.Confirmed && ctx.BlockHeight < ltc.TransactionFeeCheckBlockHeight {
		// as to ensure the miner fee is at least bigger than 0,
		// we however only want to put this restriction within the consensus set,
		// the stricter miner fee checks should apply immediately to the transaction pool logic
		constants.MinimumMinerFee = types.NewCurrency64(1)
	}
	return types.DefaultTransactionValidation(t, ctx, constants)
}

// CoinCreationTransactionController

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (cctc CoinCreationTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	cctx, err := CoinCreationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a CoinCreationTx: %v", err)
	}
	return encoding.NewEncoder(w).Encode(cctx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (cctc CoinCreationTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var cctx CoinCreationTransaction
	err := encoding.NewDecoder(r).Decode(&cctx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to binary-decode tx as a CoinCreationTx: %v", err)
	}
	// return coin creation tx as regular tfchain tx data
	return cctx.TransactionData(), nil
}

// JSONEncodeTransactionData implements TransactionController.JSONEncodeTransactionData
func (cctc CoinCreationTransactionController) JSONEncodeTransactionData(txData types.TransactionData) ([]byte, error) {
	cctx, err := CoinCreationTransactionFromTransactionData(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert txData to a CoinCreationTx: %v", err)
	}
	return json.Marshal(cctx)
}

// JSONDecodeTransactionData implements TransactionController.JSONDecodeTransactionData
func (cctc CoinCreationTransactionController) JSONDecodeTransactionData(data []byte) (types.TransactionData, error) {
	var cctx CoinCreationTransaction
	err := json.Unmarshal(data, &cctx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to json-decode tx as a CoinCreationTx: %v", err)
	}
	// return coin creation tx as regular tfchain tx data
	return cctx.TransactionData(), nil
}

// SignExtension implements TransactionExtensionSigner.SignExtension
func (cctc CoinCreationTransactionController) SignExtension(extension interface{}, sign func(*types.UnlockFulfillmentProxy, types.UnlockConditionProxy) error) (interface{}, error) {
	// (tx) extension (data) is expected to be a pointer to a valid CoinCreationTransactionExtension,
	// which contains the nonce and the mintFulfillment that can be used to fulfill the globally defined mint condition
	ccTxExtension, ok := extension.(*CoinCreationTransactionExtension)
	if !ok {
		return nil, errors.New("invalid extension data for a CoinCreationTransaction")
	}

	// get the active mint condition and use it to sign
	// NOTE: this does mean that if the mint condition suddenly this transaction will be invalid,
	// however given that only the minters (that create this coin transaction) can change the mint condition,
	// it is unlikely that this ever gives problems
	mintCondition, err := cctc.MintConditionGetter.GetActiveMintCondition()
	if err != nil {
		return nil, fmt.Errorf("failed to get the active mint condition: %v", err)
	}
	err = sign(&ccTxExtension.MintFulfillment, mintCondition)
	if err != nil {
		return nil, fmt.Errorf("failed to sign mint fulfillment of coin creation tx: %v", err)
	}
	return ccTxExtension, nil
}

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (cctc CoinCreationTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) (err error) {
	err = types.TransactionFitsInABlock(t, constants.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get CoinCreationTxn
	cctx, err := CoinCreationTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to use tx as a coin creation tx: %v", err)
	}

	// get MintCondition
	mintCondition, err := cctc.MintConditionGetter.GetMintConditionAt(ctx.BlockHeight)
	if err != nil {
		return fmt.Errorf("failed to get mint condition at block height %d: %v", ctx.BlockHeight, err)
	}

	// check if MintFulfillment fulfills the Globally defined MintCondition for the context-defined block height
	err = mintCondition.Fulfill(cctx.MintFulfillment, types.FulfillContext{
		InputIndex:  0, // InputIndex is ignored for coin creation signature
		BlockHeight: ctx.BlockHeight,
		BlockTime:   ctx.BlockTime,
		Transaction: t,
	})
	if err != nil {
		return fmt.Errorf("failed to fulfill mint condition: %v", err)
	}
	// ensure the Nonce is not Nil
	if cctx.Nonce == (TransactionNonce{}) {
		return errors.New("nil nonce is not allowed for a coin creation transaction")
	}

	// validate the rest of the content
	err = types.ArbitraryDataFits(cctx.ArbitraryData, constants.ArbitraryDataSizeLimit)
	if err != nil {
		return
	}
	for _, fee := range cctx.MinerFees {
		if fee.Cmp(constants.MinimumMinerFee) == -1 {
			return types.ErrTooSmallMinerFee
		}
	}
	// check if all condtions are standard and that the parent outputs have non-zero values
	for _, sco := range cctx.CoinOutputs {
		if sco.Value.IsZero() {
			return types.ErrZeroOutput
		}
		err = sco.Condition.IsStandardCondition(ctx)
		if err != nil {
			return err
		}
	}
	return
}

// ValidateCoinOutputs implements CoinOutputValidator.ValidateCoinOutputs
func (cctc CoinCreationTransactionController) ValidateCoinOutputs(t types.Transaction, ctx types.FundValidationContext, coinInputs map[types.CoinOutputID]types.CoinOutput) (err error) {
	return nil // always valid, coin outputs are created not backed
}

// ValidateBlockStakeOutputs implements BlockStakeOutputValidator.ValidateBlockStakeOutputs
func (cctc CoinCreationTransactionController) ValidateBlockStakeOutputs(t types.Transaction, ctx types.FundValidationContext, blockStakeInputs map[types.BlockStakeOutputID]types.BlockStakeOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within a coin creation transaction
}

// InputSigHash implements InputSigHasher.InputSigHash
func (cctc CoinCreationTransactionController) InputSigHash(t types.Transaction, _ uint64, extraObjects ...interface{}) (crypto.Hash, error) {
	cctx, err := CoinCreationTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a coin creation tx: %v", err)
	}

	h := crypto.NewHash()
	enc := encoding.NewEncoder(h)

	enc.EncodeAll(
		t.Version,
		SpecifierCoinCreationTransaction,
		cctx.Nonce,
	)

	if len(extraObjects) > 0 {
		enc.EncodeAll(extraObjects...)
	}

	enc.EncodeAll(
		cctx.CoinOutputs,
		cctx.MinerFees,
		cctx.ArbitraryData,
	)

	var hash crypto.Hash
	h.Sum(hash[:0])
	return hash, nil
}

// EncodeTransactionIDInput implements TransactionIDEncoder.EncodeTransactionIDInput
func (cctc CoinCreationTransactionController) EncodeTransactionIDInput(w io.Writer, txData types.TransactionData) error {
	cctx, err := CoinCreationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a CoinCreationTx: %v", err)
	}
	return encoding.NewEncoder(w).EncodeAll(SpecifierCoinCreationTransaction, cctx)
}

// MinterDefinitionTransactionController

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (mdtc MinterDefinitionTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	mdtx, err := MinterDefinitionTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a MinterDefinitionTx: %v", err)
	}
	return encoding.NewEncoder(w).Encode(mdtx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (mdtc MinterDefinitionTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var mdtx MinterDefinitionTransaction
	err := encoding.NewDecoder(r).Decode(&mdtx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to binary-decode tx as a MinterDefinitionTx: %v", err)
	}
	// return minter definition tx as regular tfchain tx data
	return mdtx.TransactionData(), nil
}

// JSONEncodeTransactionData implements TransactionController.JSONEncodeTransactionData
func (mdtc MinterDefinitionTransactionController) JSONEncodeTransactionData(txData types.TransactionData) ([]byte, error) {
	mdtx, err := MinterDefinitionTransactionFromTransactionData(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert txData to a MinterDefinitionTx: %v", err)
	}
	return json.Marshal(mdtx)
}

// JSONDecodeTransactionData implements TransactionController.JSONDecodeTransactionData
func (mdtc MinterDefinitionTransactionController) JSONDecodeTransactionData(data []byte) (types.TransactionData, error) {
	var mdtx MinterDefinitionTransaction
	err := json.Unmarshal(data, &mdtx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to json-decode tx as a MinterDefinitionTx: %v", err)
	}
	// return minter definition tx as regular tfchain tx data
	return mdtx.TransactionData(), nil
}

// SignExtension implements TransactionExtensionSigner.SignExtension
func (mdtc MinterDefinitionTransactionController) SignExtension(extension interface{}, sign func(*types.UnlockFulfillmentProxy, types.UnlockConditionProxy) error) (interface{}, error) {
	// (tx) extension (data) is expected to be a pointer to a valid MinterDefinitionTransactionExtension,
	// which contains the nonce and the mintFulfillment that can be used to fulfill the globally defined mint condition
	mdTxExtension, ok := extension.(*MinterDefinitionTransactionExtension)
	if !ok {
		return nil, errors.New("invalid extension data for a MinterDefinitionTx")
	}

	// get the active mint condition and use it to sign
	// NOTE: this does mean that if the mint condition suddenly this transaction will be invalid,
	// however given that only the minters (that create this coin transaction) can change the mint condition,
	// it is unlikely that this ever gives problems
	mintCondition, err := mdtc.MintConditionGetter.GetActiveMintCondition()
	if err != nil {
		return nil, fmt.Errorf("failed to get the active mint condition: %v", err)
	}
	err = sign(&mdTxExtension.MintFulfillment, mintCondition)
	if err != nil {
		return nil, fmt.Errorf("failed to sign mint fulfillment of MinterDefinitionTx: %v", err)
	}
	return mdTxExtension, nil
}

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (mdtc MinterDefinitionTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) (err error) {
	err = types.TransactionFitsInABlock(t, constants.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get MinterDefinitionTx
	mdtx, err := MinterDefinitionTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to use tx as a coin creation tx: %v", err)
	}

	// check if the MintCondition is valid
	err = mdtx.MintCondition.IsStandardCondition(ctx)
	if err != nil {
		return fmt.Errorf("defined mint condition is not standard within the given blockchain context: %v", err)
	}
	// check if the valid mint condition has a type we want to support, one of:
	//   * PubKey-UnlockHashCondtion
	//   * MultiSigConditions
	//   * TimeLockConditions (if the internal condition type is supported)
	err = validateMintCondition(mdtx.MintCondition)
	if err != nil {
		return err
	}

	// get MintCondition
	mintCondition, err := mdtc.MintConditionGetter.GetMintConditionAt(ctx.BlockHeight)
	if err != nil {
		return fmt.Errorf("failed to get mint condition at block height %d: %v", ctx.BlockHeight, err)
	}

	// check if MintFulfillment fulfills the Globally defined MintCondition for the context-defined block height
	err = mintCondition.Fulfill(mdtx.MintFulfillment, types.FulfillContext{
		InputIndex:  0, // InputIndex is ignored for coin creation signature
		BlockHeight: ctx.BlockHeight,
		BlockTime:   ctx.BlockTime,
		Transaction: t,
	})
	if err != nil {
		return fmt.Errorf("failed to fulfill mint condition: %v", err)
	}
	// ensure the Nonce is not Nil
	if mdtx.Nonce == (TransactionNonce{}) {
		return errors.New("nil nonce is not allowed for a mint condition transaction")
	}

	// validate the rest of the content
	err = types.ArbitraryDataFits(mdtx.ArbitraryData, constants.ArbitraryDataSizeLimit)
	if err != nil {
		return
	}
	for _, fee := range mdtx.MinerFees {
		if fee.Cmp(constants.MinimumMinerFee) == -1 {
			return types.ErrTooSmallMinerFee
		}
	}
	return
}

func validateMintCondition(condition types.UnlockCondition) error {
	switch ct := condition.ConditionType(); ct {
	case types.ConditionTypeMultiSignature:
		// always valid
		return nil

	case types.ConditionTypeUnlockHash:
		// only valid for unlock hash type 1 (PubKey)
		if condition.UnlockHash().Type == types.UnlockTypePubKey {
			return nil
		}
		return errors.New("unlockHash conditions can be used as mint conditions, if the unlock hash type is PubKey")

	case types.ConditionTypeTimeLock:
		// ensure to unpack a proxy condition first
		if cp, ok := condition.(types.UnlockConditionProxy); ok {
			condition = cp.Condition
		}
		// time lock conditions are allowed as long as the internal condition is allowed
		cg, ok := condition.(types.MarshalableUnlockConditionGetter)
		if !ok {
			err := fmt.Errorf("unexpected Go-type for TimeLockCondition: %T", condition)
			if build.DEBUG {
				panic(err)
			}
			return err
		}
		return validateMintCondition(cg.GetMarshalableUnlockCondition())

	default:
		// all other types aren't allowed
		return fmt.Errorf("condition type %d cannot be used as a mint condition", ct)
	}
}

// ValidateCoinOutputs implements CoinOutputValidator.ValidateCoinOutputs
func (mdtc MinterDefinitionTransactionController) ValidateCoinOutputs(t types.Transaction, ctx types.FundValidationContext, coinInputs map[types.CoinOutputID]types.CoinOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within a minter definition transaction
}

// ValidateBlockStakeOutputs implements BlockStakeOutputValidator.ValidateBlockStakeOutputs
func (mdtc MinterDefinitionTransactionController) ValidateBlockStakeOutputs(t types.Transaction, ctx types.FundValidationContext, blockStakeInputs map[types.BlockStakeOutputID]types.BlockStakeOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within a minter definition transaction
}

// InputSigHash implements InputSigHasher.InputSigHash
func (mdtc MinterDefinitionTransactionController) InputSigHash(t types.Transaction, _ uint64, extraObjects ...interface{}) (crypto.Hash, error) {
	mdtx, err := MinterDefinitionTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a MinterDefinitionTx: %v", err)
	}

	h := crypto.NewHash()
	enc := encoding.NewEncoder(h)

	enc.EncodeAll(
		t.Version,
		SpecifierMintDefinitionTransaction,
		mdtx.Nonce,
	)

	if len(extraObjects) > 0 {
		enc.EncodeAll(extraObjects...)
	}

	enc.EncodeAll(
		mdtx.MintCondition,
		mdtx.MinerFees,
		mdtx.ArbitraryData,
	)

	var hash crypto.Hash
	h.Sum(hash[:0])
	return hash, nil
}

// EncodeTransactionIDInput implements TransactionIDEncoder.EncodeTransactionIDInput
func (mdtc MinterDefinitionTransactionController) EncodeTransactionIDInput(w io.Writer, txData types.TransactionData) error {
	mdtx, err := MinterDefinitionTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a MinterDefinitionTx: %v", err)
	}
	return encoding.NewEncoder(w).EncodeAll(SpecifierMintDefinitionTransaction, mdtx)
}

type (
	// CoinCreationTransaction is to be created only by the defined Coin Minters,
	// as a medium in order to create coins (coin outputs), without backing them
	// (so without having to spend previously unspend coin outputs, see: coin inputs).
	CoinCreationTransaction struct {
		// Nonce used to ensure the uniqueness of a CoinCreationTransaction's ID and signature.
		Nonce TransactionNonce `json:"nonce"`
		// MintFulfillment defines the fulfillment which is used in order to
		// fulfill the globally defined MintCondition.
		MintFulfillment types.UnlockFulfillmentProxy `json:"mintfulfillment"`
		// CoinOutputs defines the coin outputs,
		// which contain the freshly created coins, adding to the total pool of coins
		// available in the tfchain network.
		CoinOutputs []types.CoinOutput `json:"coinoutputs"`
		// Minerfees, a fee paid for this coin creation transaction.
		MinerFees []types.Currency `json:"minerfees"`
		// ArbitraryData can be used for any purpose,
		// but is mostly to be used in order to define the reason/origins
		// of the coin creation.
		ArbitraryData []byte `json:"arbitrarydata,omitempty"`
	}
	// CoinCreationTransactionExtension defines the CoinCreationTx Extension Data
	CoinCreationTransactionExtension struct {
		Nonce           TransactionNonce
		MintFulfillment types.UnlockFulfillmentProxy
	}
)

// CoinCreationTransactionFromTransaction creates a CoinCreationTransaction,
// using a regular in-memory tfchain transaction.
//
// Past the (tx) Version validation it piggy-backs onto the
// `CoinCreationTransactionFromTransactionData` constructor.
func CoinCreationTransactionFromTransaction(tx types.Transaction) (CoinCreationTransaction, error) {
	if tx.Version != TransactionVersionCoinCreation {
		return CoinCreationTransaction{}, fmt.Errorf(
			"a coin creation transaction requires tx version %d",
			TransactionVersionCoinCreation)
	}
	return CoinCreationTransactionFromTransactionData(types.TransactionData{
		CoinInputs:        tx.CoinInputs,
		CoinOutputs:       tx.CoinOutputs,
		BlockStakeInputs:  tx.BlockStakeInputs,
		BlockStakeOutputs: tx.BlockStakeOutputs,
		MinerFees:         tx.MinerFees,
		ArbitraryData:     tx.ArbitraryData,
		Extension:         tx.Extension,
	})
}

// CoinCreationTransactionFromTransactionData creates a CoinCreationTransaction,
// using the TransactionData from a regular in-memory tfchain transaction.
func CoinCreationTransactionFromTransactionData(txData types.TransactionData) (CoinCreationTransaction, error) {
	// (tx) extension (data) is expected to be a pointer to a valid CoinCreationTransactionExtension,
	// which contains the nonce and the mintFulfillment that can be used to fulfill the globally defined mint condition
	extensionData, ok := txData.Extension.(*CoinCreationTransactionExtension)
	if !ok {
		return CoinCreationTransaction{}, errors.New("invalid extension data for a CoinCreationTransaction")
	}
	// at least one coin output as well as one miner fee is required
	if len(txData.CoinOutputs) == 0 || len(txData.MinerFees) == 0 {
		return CoinCreationTransaction{}, errors.New("at least one coin output and miner fee is required for a CoinCreationTransaction")
	}
	// no coin inputs, block stake inputs or block stake outputs are allowed
	if len(txData.CoinInputs) != 0 || len(txData.BlockStakeInputs) != 0 || len(txData.BlockStakeOutputs) != 0 {
		return CoinCreationTransaction{}, errors.New("no coin inputs and block stake inputs/outputs are allowed in a CoinCreationTransaction")
	}
	// return the CoinCreationTransaction, with the data extracted from the TransactionData
	return CoinCreationTransaction{
		Nonce:           extensionData.Nonce,
		MintFulfillment: extensionData.MintFulfillment,
		CoinOutputs:     txData.CoinOutputs,
		MinerFees:       txData.MinerFees,
		// ArbitraryData is optional
		ArbitraryData: txData.ArbitraryData,
	}, nil
}

// TransactionData returns this CoinCreationTransaction
// as regular tfchain transaction data.
func (cctx *CoinCreationTransaction) TransactionData() types.TransactionData {
	return types.TransactionData{
		CoinOutputs:   cctx.CoinOutputs,
		MinerFees:     cctx.MinerFees,
		ArbitraryData: cctx.ArbitraryData,
		Extension: &CoinCreationTransactionExtension{
			Nonce:           cctx.Nonce,
			MintFulfillment: cctx.MintFulfillment,
		},
	}
}

// Transaction returns this CoinCreationTransaction
// as regular tfchain transaction, using TransactionVersionCoinCreation as the type.
func (cctx *CoinCreationTransaction) Transaction() types.Transaction {
	return types.Transaction{
		Version:       TransactionVersionCoinCreation,
		CoinOutputs:   cctx.CoinOutputs,
		MinerFees:     cctx.MinerFees,
		ArbitraryData: cctx.ArbitraryData,
		Extension: &CoinCreationTransactionExtension{
			Nonce:           cctx.Nonce,
			MintFulfillment: cctx.MintFulfillment,
		},
	}
}

type (
	// MinterDefinitionTransaction is to be created only by the defined Coin Minters,
	// as a medium in order to transfer minting powers.
	MinterDefinitionTransaction struct {
		// Nonce used to ensure the uniqueness of a MinterDefinitionTransaction's ID and signature.
		Nonce TransactionNonce `json:"nonce"`
		// MintFulfillment defines the fulfillment which is used in order to
		// fulfill the globally defined MintCondition.
		MintFulfillment types.UnlockFulfillmentProxy `json:"mintfulfillment"`
		// MintCondition defines a new condition that defines who become(s) the new minter(s),
		// and thus defines who can create coins as well as update who is/are the current minter(s)
		//
		// UnlockHash (unlockhash type 1) and MultiSigConditions are allowed,
		// as well as TimeLocked conditions which have UnlockHash- and MultiSigConditions as
		// internal condition.
		MintCondition types.UnlockConditionProxy `json:"mintcondition"`
		// Minerfees, a fee paid for this minter definition transaction.
		MinerFees []types.Currency `json:"minerfees"`
		// ArbitraryData can be used for any purpose,
		// but is mostly to be used in order to define the reason/origins
		// of the transfer of minting power.
		ArbitraryData []byte `json:"arbitrarydata,omitempty"`
	}
	// MinterDefinitionTransactionExtension defines the MinterDefinitionTx Extension Data
	MinterDefinitionTransactionExtension struct {
		Nonce           TransactionNonce
		MintFulfillment types.UnlockFulfillmentProxy
		MintCondition   types.UnlockConditionProxy
	}
)

// MinterDefinitionTransactionFromTransaction creates a MinterDefinitionTransaction,
// using a regular in-memory tfchain transaction.
//
// Past the (tx) Version validation it piggy-backs onto the
// `MinterDefinitionTransactionFromTransactionData` constructor.
func MinterDefinitionTransactionFromTransaction(tx types.Transaction) (MinterDefinitionTransaction, error) {
	if tx.Version != TransactionVersionMinterDefinition {
		return MinterDefinitionTransaction{}, fmt.Errorf(
			"a minter definition transaction requires tx version %d",
			TransactionVersionCoinCreation)
	}
	return MinterDefinitionTransactionFromTransactionData(types.TransactionData{
		CoinInputs:        tx.CoinInputs,
		CoinOutputs:       tx.CoinOutputs,
		BlockStakeInputs:  tx.BlockStakeInputs,
		BlockStakeOutputs: tx.BlockStakeOutputs,
		MinerFees:         tx.MinerFees,
		ArbitraryData:     tx.ArbitraryData,
		Extension:         tx.Extension,
	})
}

// MinterDefinitionTransactionFromTransactionData creates a MinterDefinitionTransaction,
// using the TransactionData from a regular in-memory tfchain transaction.
func MinterDefinitionTransactionFromTransactionData(txData types.TransactionData) (MinterDefinitionTransaction, error) {
	// (tx) extension (data) is expected to be a pointer to a valid MinterDefinitionTransactionExtension,
	// which contains the nonce, the mintFulfillment that can be used to fulfill the currently globally defined mint condition,
	// as well as a mintCondition to replace the current in-place mintCondition.
	extensionData, ok := txData.Extension.(*MinterDefinitionTransactionExtension)
	if !ok {
		return MinterDefinitionTransaction{}, errors.New("invalid extension data for a MinterDefinitionTransaction")
	}
	// at least one miner fee is required
	if len(txData.MinerFees) == 0 {
		return MinterDefinitionTransaction{}, errors.New("at least one miner fee is required for a MinterDefinitionTransaction")
	}
	// no coin inputs, block stake inputs or block stake outputs are allowed
	if len(txData.CoinInputs) != 0 || len(txData.CoinOutputs) != 0 || len(txData.BlockStakeInputs) != 0 || len(txData.BlockStakeOutputs) != 0 {
		return MinterDefinitionTransaction{}, errors.New(
			"no coin inputs/outputs and block stake inputs/outputs are allowed in a MinterDefinitionTransaction")
	}
	// return the MinterDefinitionTransaction, with the data extracted from the TransactionData
	return MinterDefinitionTransaction{
		Nonce:           extensionData.Nonce,
		MintFulfillment: extensionData.MintFulfillment,
		MintCondition:   extensionData.MintCondition,
		MinerFees:       txData.MinerFees,
		// ArbitraryData is optional
		ArbitraryData: txData.ArbitraryData,
	}, nil
}

// TransactionData returns this CoinCreationTransaction
// as regular tfchain transaction data.
func (cctx *MinterDefinitionTransaction) TransactionData() types.TransactionData {
	return types.TransactionData{
		MinerFees:     cctx.MinerFees,
		ArbitraryData: cctx.ArbitraryData,
		Extension: &MinterDefinitionTransactionExtension{
			Nonce:           cctx.Nonce,
			MintFulfillment: cctx.MintFulfillment,
			MintCondition:   cctx.MintCondition,
		},
	}
}

// Transaction returns this CoinCreationTransaction
// as regular tfchain transaction, using TransactionVersionCoinCreation as the type.
func (cctx *MinterDefinitionTransaction) Transaction() types.Transaction {
	return types.Transaction{
		Version:       TransactionVersionMinterDefinition,
		MinerFees:     cctx.MinerFees,
		ArbitraryData: cctx.ArbitraryData,
		Extension: &MinterDefinitionTransactionExtension{
			Nonce:           cctx.Nonce,
			MintFulfillment: cctx.MintFulfillment,
			MintCondition:   cctx.MintCondition,
		},
	}
}

// 3bot Multiplier fees that have to be multiplied with the OneCoin definition,
// in order to know the amount in the used chain currency (TFT).
const (
	BotFeePerAdditionalNameMultiplier           = 50
	BotFeeForNetworkAddressInfoChangeMultiplier = 20
	BotRegistrationFeeMultiplier                = 90
	BotMonthlyFeeMultiplier                     = 10
)

// [DONE] define the binary marshalling for each of the 3bot Tx's
//   TODO: ^TEST THIS LOGIC^
// TODO: define the Tx controllers for each of the 3bot Tx's
//   TODO: ^TEST THIS LOGIC^

type (
	// BotRegistrationTransaction defines the Transaction (with version 0x90)
	// used to register a new 3bot, where new means that the used public key
	// (identification) cannot yet exist.
	BotRegistrationTransaction struct {
		// Addresses contains the optional network addresses used to reach the 3bot.
		// Normally at least one is given, none are required however.
		// All addresses (max 10) can be of any of the following types: IPv4, IPv6, hostname
		Addresses []NetworkAddress `json:"addresses"`
		// Names contains the optional names (max 5) that can be used to reach the bot,
		// using a name, instead of one of its network addresses, comparable to how DNS works.
		Names []BotName `json:"names"`

		// NrOfMonths defines the amount of months that
		// is desired to be paid upfront. Note that the amount of
		// months defined here indicates how much additional fees are to be paid.
		// The NrOfMonths has to be within this inclusive range [1,24].
		NrOfMonths uint8 `json:"nrofmonths"`

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency `json:"txfee"`

		// CoinInputs are only used for the required fees,
		// which contains the regular Tx fee as well as the additional fees,
		// to be paid for a 3bot registration. At least one CoinInput is required.
		CoinInputs []types.CoinInput `json:"coininputs"`
		// RefundCoinOutput is an optional coin output that can be used
		// to refund coins paid as inputs for the required fees.
		RefundCoinOutput *types.CoinOutput `json:"refundcoinoutput"`

		// Identification is used to identify the 3bot and verify its identity.
		// The identification is only given at registration, for all other
		// 3bot Tx types it is identified by a combination of its unique ID and signature.
		Identification PublicKeySignaturePair `json:"identification"`
	}
	// BotRegistrationTransactionExtension defines the BotRegistrationTransaction Extension Data
	BotRegistrationTransactionExtension struct {
		Addresses      []NetworkAddress
		Names          []BotName
		NrOfMonths     uint8
		Identification PublicKeySignaturePair
	}
)

// BotRegistrationTransactionFromTransaction creates a BotRegistrationTransaction,
// using a regular in-memory tfchain transaction.
//
// Past the (tx) Version validation it piggy-backs onto the
// `BotRegistrationTransactionFromTransactionData` constructor.
func BotRegistrationTransactionFromTransaction(tx types.Transaction) (BotRegistrationTransaction, error) {
	if tx.Version != TransactionVersionBotRegistration {
		return BotRegistrationTransaction{}, fmt.Errorf(
			"a bot registration transaction requires tx version %d",
			TransactionVersionBotRegistration)
	}
	return BotRegistrationTransactionFromTransactionData(types.TransactionData{
		CoinInputs:        tx.CoinInputs,
		CoinOutputs:       tx.CoinOutputs,
		BlockStakeInputs:  tx.BlockStakeInputs,
		BlockStakeOutputs: tx.BlockStakeOutputs,
		MinerFees:         tx.MinerFees,
		ArbitraryData:     tx.ArbitraryData,
		Extension:         tx.Extension,
	})
}

// BotRegistrationTransactionFromTransactionData creates a BotRegistrationTransaction,
// using the TransactionData from a regular in-memory tfchain transaction.
func BotRegistrationTransactionFromTransactionData(txData types.TransactionData) (BotRegistrationTransaction, error) {
	// (tx) extension (data) is expected to be a pointer to a valid BotRegistrationTransaction,
	// which contains all the properties unique to a 3bot (registration) Tx
	extensionData, ok := txData.Extension.(*BotRegistrationTransactionExtension)
	if !ok {
		return BotRegistrationTransaction{}, errors.New("invalid extension data for a BotRegistrationTransaction")
	}
	// at least one coin input as well as one miner fee is required
	if len(txData.CoinInputs) == 0 || len(txData.MinerFees) != 1 {
		return BotRegistrationTransaction{}, errors.New("at least one coin input and exactly one miner fee is required for a BotRegistrationTransaction")
	}
	// no block stake inputs or block stake outputs are allowed
	if len(txData.BlockStakeInputs) != 0 || len(txData.BlockStakeOutputs) != 0 {
		return BotRegistrationTransaction{}, errors.New("no block stake inputs/outputs are allowed in a BotRegistrationTransaction")
	}

	tx := BotRegistrationTransaction{
		Addresses:      extensionData.Addresses,
		Names:          extensionData.Names,
		NrOfMonths:     extensionData.NrOfMonths,
		TransactionFee: txData.MinerFees[0],
		CoinInputs:     txData.CoinInputs,
		Identification: extensionData.Identification,
	}
	switch len(txData.CoinOutputs) {
	case 2:
		// take refund coin output
		// convention always assumed to be the required BotFee
		tx.RefundCoinOutput = &txData.CoinOutputs[1]
		return tx, nil
	case 1:
		// nothing to do
		return tx, nil
	default:
		// return with an error, maximum one coin output can be set, not more
		return BotRegistrationTransaction{}, errors.New("only one coinoutput is allowed (and required)")
	}
}

// TransactionData returns this BotRegistrationTransaction
// as regular tfchain transaction data.
func (brtx *BotRegistrationTransaction) TransactionData(oneCoin types.Currency, registryPoolCondition types.UnlockConditionProxy) types.TransactionData {
	txData := types.TransactionData{
		CoinInputs: brtx.CoinInputs,
		CoinOutputs: []types.CoinOutput{
			{
				Value:     brtx.RequiredBotFee(oneCoin),
				Condition: registryPoolCondition,
			},
		},
		MinerFees: []types.Currency{brtx.TransactionFee},
		Extension: &BotRegistrationTransactionExtension{
			Addresses:      brtx.Addresses,
			Names:          brtx.Names,
			NrOfMonths:     brtx.NrOfMonths,
			Identification: brtx.Identification,
		},
	}
	if brtx.RefundCoinOutput != nil {
		txData.CoinOutputs = append(txData.CoinOutputs, *brtx.RefundCoinOutput)
	}
	return txData
}

// Transaction returns this BotRegistrationTransaction
// as regular tfchain transaction, using TransactionVersionBotRegistration as the type.
func (brtx *BotRegistrationTransaction) Transaction(oneCoin types.Currency, registryPoolCondition types.UnlockConditionProxy) types.Transaction {
	tx := types.Transaction{
		Version:    TransactionVersionBotRegistration,
		CoinInputs: brtx.CoinInputs,
		CoinOutputs: []types.CoinOutput{
			{
				Value:     brtx.RequiredBotFee(oneCoin),
				Condition: registryPoolCondition,
			},
		},
		MinerFees: []types.Currency{brtx.TransactionFee},
		Extension: &BotRegistrationTransactionExtension{
			Addresses:      brtx.Addresses,
			Names:          brtx.Names,
			NrOfMonths:     brtx.NrOfMonths,
			Identification: brtx.Identification,
		},
	}
	if brtx.RefundCoinOutput != nil {
		tx.CoinOutputs = append(tx.CoinOutputs, *brtx.RefundCoinOutput)
	}
	return tx
}

// RequiredBotFee computes the required Bot Fee, that is to be applied as a required
// additional fee on top of the regular required (minimum) Tx fee.
func (brtx *BotRegistrationTransaction) RequiredBotFee(oneCoin types.Currency) types.Currency {
	// a static registration fee has to be paid
	fee := oneCoin.Mul64(BotRegistrationFeeMultiplier)
	// the amount of desired months also has to be paid
	fee = fee.Add(ComputeMonthlyBotFees(brtx.NrOfMonths, oneCoin))
	// if more than one name is defined it also has to be paid
	if n := len(brtx.Names); n > 1 {
		fee = fee.Add(oneCoin.Mul64(uint64(n-1) * BotFeePerAdditionalNameMultiplier))
	}
	// no fee has to be paid for the used network addresses during registration
	// return the total fees
	return fee
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (brtx BotRegistrationTransaction) MarshalSia(w io.Writer) error {
	// the tfchain binary encoder used for this implementation
	enc := tfencoding.NewEncoder(w)

	// encode the nr of months, flags and paired lenghts
	addrLen := len(brtx.Addresses)
	nameLen := len(brtx.Names)
	maf := &BotMonthsAndFlagsData{
		NrOfMonths:   brtx.NrOfMonths,
		HasAddresses: addrLen != 0,
		HasNames:     nameLen != 0,
		HasRefund:    brtx.RefundCoinOutput != nil,
	}
	err := enc.EncodeAll(maf, (uint8(addrLen) | (uint8(nameLen) << 4)))
	if err != nil {
		return err
	}
	// encode all addresses
	for _, addr := range brtx.Addresses {
		err = enc.Encode(addr)
		if err != nil {
			return err
		}
	}
	// encode all names
	for _, name := range brtx.Names {
		err = enc.Encode(name)
		if err != nil {
			return err
		}
	}
	// encode TxFee and CoinInputs
	err = enc.EncodeAll(brtx.TransactionFee, brtx.CoinInputs)
	if err != nil {
		return err
	}
	// encode refund coin output, if given
	if maf.HasRefund {
		// deref to ensure we do not also encode one byte
		// for the pointer indication
		err = enc.Encode(*brtx.RefundCoinOutput)
		if err != nil {
			return err
		}
	}
	// encode the identification at the end
	return enc.Encode(brtx.Identification)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (brtx *BotRegistrationTransaction) UnmarshalSia(r io.Reader) error {
	dec := tfencoding.NewDecoder(r)

	var maf BotMonthsAndFlagsData
	err := dec.Decode(&maf)
	if err != nil {
		return err
	}

	// assign number of months
	brtx.NrOfMonths = maf.NrOfMonths

	// decode the pair length (length of both names and addresses in one byte)
	var pairLength uint8
	err = dec.Decode(&pairLength)
	if err != nil {
		return err
	}

	addrLen, nameLen := pairLength&15, pairLength>>4

	// decode all addresses and all names and store them in this Tx
	if addrLen > 0 {
		brtx.Addresses = make([]NetworkAddress, addrLen)
		for i := range brtx.Addresses {
			err = dec.Decode(&brtx.Addresses[i])
			if err != nil {
				return err
			}
		}
	} else {
		brtx.Addresses = nil
	}
	if nameLen > 0 {
		brtx.Names = make([]BotName, nameLen)
		for i := range brtx.Names {
			err = dec.Decode(&brtx.Names[i])
			if err != nil {
				return err
			}
		}
	} else {
		brtx.Names = nil
	}

	// decode tx fee and coin inputs
	err = dec.DecodeAll(&brtx.TransactionFee, &brtx.CoinInputs)
	if err != nil {
		return err
	}

	// decode the refund coin output, only if its flag is defined
	if maf.HasRefund {
		brtx.RefundCoinOutput = new(types.CoinOutput)
		err = dec.Decode(brtx.RefundCoinOutput)
		if err != nil {
			return err
		}
	} else {
		brtx.RefundCoinOutput = nil // explicitly set it nil
	}

	// decode identification as the last step
	return dec.Decode(&brtx.Identification)
}

type (
	// BotRecordUpdateTransaction defines the Transaction (with version 0x91)
	// used to update a 3bot Record by the owner. where owner
	// means the 3bot that created the record to be updated initially using the BotRegistration Tx.
	BotRecordUpdateTransaction struct {
		// Identifier of the 3bot, used to find the 3bot record to be updated,
		// and verify that the Tx is authorized to do so.
		Identifier BotID `json:"id"`

		// Addresses can be used to add and/or remove network addresses
		// to/from the existing 3bot record. Note that after each Tx,
		// no more than 10 addresses can be linked to a single 3bot record.
		Addresses struct {
			Add    []NetworkAddress `json:"add"`
			Remove []NetworkAddress `json:"remove"`
		} `json:"addresses"`

		// Names can be used to add and/or remove names
		// to/from the existing 3bot record. Note that after each Tx,
		// no more than 5 names can be linked to a single 3bot record.
		Names struct {
			Add    []BotName `json:"add"`
			Remove []BotName `json:"remove"`
		} `json:"names"`

		// NrOfMonths defines the optional amount of months that
		// is desired to be paid upfront in this update. Note that the amount of
		// months defined here defines how much additional fees are to be paid.
		// The NrOfMonths has to be within this inclusive range [0,24].
		NrOfMonths uint8 `json:"nrofmonths"`

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency `json:"txfee"`

		// CoinInputs are only used for the required fees,
		// which contains the regular Tx fee as well as the additional fees,
		// to be paid for a 3bot record update. At least one CoinInput is required.
		// If this 3bot record update is only to pay for extending the 3bot activity,
		// than no fees are required other than the monthly fees as defined by this bots usage.
		CoinInputs []types.CoinInput `json:"coininputs"`
		// RefundCoinOutput is an optional coin output that can be used
		// to refund coins paid as inputs for the required fees.
		RefundCoinOutput *types.CoinOutput `json:"refundcoinoutput"`

		// Signature is used to proof the ownership of the 3bot record to be updated,
		// and is verified using the public key defined in the 3bot linked
		// to the given (3bot) identifier.
		Signature types.ByteSlice `json:"signature"`
	}
	// BotRecordUpdateTransactionExtension defines the BotRecordUpdateTransaction Extension Data
	BotRecordUpdateTransactionExtension struct {
		Identifier       BotID
		Signature        types.ByteSlice
		AddressesAdded   []NetworkAddress
		AddressesRemoved []NetworkAddress
		NamesAdded       []BotName
		NamesRemoved     []BotName
		NrOfMonths       uint8
	}
)

// RequiredBotFee computes the required Bot Fee, that is to be applied as a required
// additional fee on top of the regular required (minimum) Tx fee.
func (brutx *BotRecordUpdateTransaction) RequiredBotFee(oneCoin types.Currency) (fee types.Currency) {
	// all additional months have to be paid
	if brutx.NrOfMonths > 0 {
		fee = fee.Add(ComputeMonthlyBotFees(brutx.NrOfMonths, oneCoin))
	}
	// a Tx that modifies the network address info of a 3bot record also has to be paid
	if len(brutx.Addresses.Add) > 0 || len(brutx.Addresses.Remove) > 0 {
		fee = fee.Add(oneCoin.Mul64(BotFeeForNetworkAddressInfoChangeMultiplier))
	}
	// each additional name has to be paid as well
	// (regardless of the fact that the 3bot has a name or not)
	if n := len(brutx.Names.Add); n > 0 {
		fee = fee.Add(oneCoin.Mul64(BotFeePerAdditionalNameMultiplier * uint64(n)))
	}
	// return the total fees
	return fee
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (brutx BotRecordUpdateTransaction) MarshalSia(w io.Writer) error {
	// collect length of all the name/addr slices
	addrAddLen, addrRemoveLen := len(brutx.Addresses.Add), len(brutx.Addresses.Remove)
	nameAddLen, nameRemoveLen := len(brutx.Names.Add), len(brutx.Names.Remove)

	// the tfchain binary encoder used for this implementation
	enc := tfencoding.NewEncoder(w)

	// encode the identifier, nr of months, flags and paired lenghts
	maf := BotMonthsAndFlagsData{
		NrOfMonths:   brutx.NrOfMonths,
		HasAddresses: addrAddLen > 0 || addrRemoveLen > 0,
		HasNames:     nameAddLen > 0 || nameRemoveLen > 0,
		HasRefund:    brutx.RefundCoinOutput != nil,
	}
	err := enc.EncodeAll(brutx.Identifier, maf)
	if err != nil {
		return err
	}

	// encode addressed added and removed, if defined
	if maf.HasAddresses {
		err = enc.Encode(uint8(addrAddLen) | (uint8(addrRemoveLen) << 4))
		if err != nil {
			return err
		}
		for _, addr := range brutx.Addresses.Add {
			err = enc.Encode(addr)
			if err != nil {
				return err
			}
		}
		for _, addr := range brutx.Addresses.Remove {
			err = enc.Encode(addr)
			if err != nil {
				return err
			}
		}
	}

	// encode names added and removed, if defined
	if maf.HasNames {
		err = enc.Encode(uint8(nameAddLen) | (uint8(nameRemoveLen) << 4))
		if err != nil {
			return err
		}
		for _, name := range brutx.Names.Add {
			err = enc.Encode(name)
			if err != nil {
				return err
			}
		}
		for _, name := range brutx.Names.Remove {
			err = enc.Encode(name)
			if err != nil {
				return err
			}
		}
	}

	// encode TxFee and CoinInputs
	err = enc.EncodeAll(brutx.TransactionFee, brutx.CoinInputs)
	if err != nil {
		return err
	}
	// encode refund coin output, if given
	if maf.HasRefund {
		// deref to ensure we do not also encode one byte
		// for the pointer indication
		err = enc.Encode(*brutx.RefundCoinOutput)
		if err != nil {
			return err
		}
	}
	// encode the signature at the end
	return enc.Encode(brutx.Signature)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (brutx *BotRecordUpdateTransaction) UnmarshalSia(r io.Reader) error {
	dec := tfencoding.NewDecoder(r)

	// unmarshal identifier, NrOfMonths and flags
	var maf BotMonthsAndFlagsData
	err := dec.DecodeAll(&brutx.Identifier, &maf)
	if err != nil {
		return err
	}

	// assign number of months
	brutx.NrOfMonths = maf.NrOfMonths

	// decode addressed added and removed, if defined
	if maf.HasAddresses {
		var pairLength uint8
		err = dec.Decode(&pairLength)
		if err != nil {
			return err
		}
		addrAddLen, addrRemoveLen := pairLength&15, pairLength>>4
		if addrAddLen > 0 {
			brutx.Addresses.Add = make([]NetworkAddress, addrAddLen)
			for i := range brutx.Addresses.Add {
				err = dec.Decode(&brutx.Addresses.Add[i])
				if err != nil {
					return err
				}
			}
		} else {
			brutx.Addresses.Add = nil
		}
		if addrRemoveLen > 0 {
			brutx.Addresses.Remove = make([]NetworkAddress, addrRemoveLen)
			for i := range brutx.Addresses.Remove {
				err = dec.Decode(&brutx.Addresses.Remove[i])
				if err != nil {
					return err
				}
			}
		} else {
			brutx.Addresses.Remove = nil
		}
	} else {
		// explicitly set added/removed address to nil
		brutx.Addresses.Add, brutx.Addresses.Remove = nil, nil
	}

	if maf.HasNames {
		var pairLength uint8
		err = dec.Decode(&pairLength)
		if err != nil {
			return err
		}
		nameAddLen, nameRemoveLen := pairLength&15, pairLength>>4
		if nameAddLen > 0 {
			brutx.Names.Add = make([]BotName, nameAddLen)
			for i := range brutx.Names.Add {
				err = dec.Decode(&brutx.Names.Add[i])
				if err != nil {
					return err
				}
			}
		} else {
			brutx.Names.Add = nil
		}
		if nameRemoveLen > 0 {
			brutx.Names.Remove = make([]BotName, nameRemoveLen)
			for i := range brutx.Names.Remove {
				err = dec.Decode(&brutx.Names.Remove[i])
				if err != nil {
					return err
				}
			}
		} else {
			brutx.Names.Remove = nil
		}
	} else {
		// explicitly set added/removed address to nil
		brutx.Names.Add, brutx.Names.Remove = nil, nil
	}

	// encode TxFee and CoinInputs
	err = dec.DecodeAll(&brutx.TransactionFee, &brutx.CoinInputs)
	if err != nil {
		return err
	}
	// decode refund coin output, if defined
	if maf.HasRefund {
		brutx.RefundCoinOutput = new(types.CoinOutput)
		err = dec.Decode(brutx.RefundCoinOutput)
		if err != nil {
			return err
		}
	}
	// decode the signature at the end
	return dec.Decode(&brutx.Signature)
}

type (
	// BotNameTransferTransaction defines the Transaction (with version 0x92)
	// used to transfer one or multiple names from the active
	// 3bot that up to the point of the Tx to another 3bot.
	BotNameTransferTransaction struct {
		// Sender is in this context the 3bot that owns and transfers the names
		// defined in this Tx to the 3bot defined in this Tx as the Receiver.
		// The Sender has to be different from the Receiver.
		Sender struct {
			Identifier BotID           `json:"id"`
			Signature  types.ByteSlice `json:"signature"`
		} `json:"sender"`
		// Receiver is in this context the 3bot that receives the names
		// defined in this Tx from the 3bot defined in this Tx as the Sender.
		// The Receiver has to be different from the Sender.
		Receiver struct {
			Identifier BotID           `json:"id"`
			Signature  types.ByteSlice `json:"signature"`
		} `json:"receiver"`

		// Addresses can be used to add and/or remove network addresses
		// to/from the existing 3bot record. Note that after each Tx,
		// no more than 10 addresses can be linked to a single 3bot record.
		Addresses struct {
			Add    []NetworkAddress `json:"add"`
			Remove []NetworkAddress `json:"remove"`
		} `json:"addresses"`

		// Names can be used to add and/or remove names
		// to/from the existing 3bot record. Note that after each Tx,
		// no more than 5 names can be linked to a single 3bot record.
		Names struct {
			Add    []BotName `json:"add"`
			Remove []BotName `json:"remove"`
		} `json:"names"`

		// NrOfMonths defines the optional amount of months that
		// is desired to be paid upfront in this update. Note that the amount of
		// months defined here defines how much additional fees are to be paid.
		// The NrOfMonths has to be within this inclusive range [0,24].
		NrOfMonths uint8 `json:"nrofmonths"`

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency `json:"txfee"`

		// CoinInputs are only used for the required fees,
		// which contains the regular Tx fee as well as the additional fees,
		// to be paid for a 3bot record update. At least one CoinInput is required.
		// If this 3bot record update is only to pay for extending the 3bot activity,
		// than no fees are required other than the monthly fees as defined by this bots usage.
		CoinInputs []types.CoinInput `json:"coininputs"`
		// RefundCoinOutput is an optional coin output that can be used
		// to refund coins paid as inputs for the required fees.
		RefundCoinOutput *types.CoinOutput `json:"refundcoinoutput"`
	}
	// BotNameTransferTransactionExtension defines the BotNameTransferTransaction Extension Data
	BotNameTransferTransactionExtension struct {
		SenderIdentifier        BotID
		SenderSignature         types.ByteSlice
		ReceiverIdentifier      BotID
		ReceiverSenderSignature types.ByteSlice
		AddressesAdded          []NetworkAddress
		AddressesRemoved        []NetworkAddress
		NamesAdded              []BotName
		NamesRemoved            []BotName
		NrOfMonths              uint8
	}
)

// RequiredBotFee computes the required Bot Fee, that is to be applied as a required
// additional fee on top of the regular required (minimum) Tx fee.
func (bnttx *BotNameTransferTransaction) RequiredBotFee(oneCoin types.Currency) (fee types.Currency) {
	// all additional months have to be paid
	if bnttx.NrOfMonths > 0 {
		fee = fee.Add(ComputeMonthlyBotFees(bnttx.NrOfMonths, oneCoin))
	}
	// a Tx that modifies the network address info of a 3bot record also has to be paid
	if len(bnttx.Addresses.Add) > 0 || len(bnttx.Addresses.Remove) > 0 {
		fee = fee.Add(oneCoin.Mul64(BotFeeForNetworkAddressInfoChangeMultiplier))
	}
	// each additional name has to be paid as well
	// (regardless of the fact that the 3bot has a name or not)
	if n := len(bnttx.Names.Add); n > 0 {
		fee = fee.Add(oneCoin.Mul64(BotFeePerAdditionalNameMultiplier * uint64(n)))
	}
	// return the total fees
	return fee
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (bnttx BotNameTransferTransaction) MarshalSia(w io.Writer) error {
	addrAddLen, addrRemoveLen := len(bnttx.Addresses.Add), len(bnttx.Addresses.Remove)
	nameAddLen, nameRemoveLen := len(bnttx.Names.Add), len(bnttx.Names.Remove)

	// the tfchain binary encoder used for this implementation
	enc := tfencoding.NewEncoder(w)

	// encode the sender, receiver, nr of months, flags and paired lenghts
	maf := BotMonthsAndFlagsData{
		NrOfMonths:   bnttx.NrOfMonths,
		HasAddresses: addrAddLen > 0 || addrRemoveLen > 0,
		HasNames:     nameAddLen > 0 || nameRemoveLen > 0,
		HasRefund:    bnttx.RefundCoinOutput != nil,
	}
	err := enc.EncodeAll(bnttx.Sender, bnttx.Receiver, maf)
	if err != nil {
		return err
	}

	// encode addressed added and removed, if defined
	if maf.HasAddresses {
		err = enc.Encode(uint8(addrAddLen) | (uint8(addrRemoveLen) << 4))
		if err != nil {
			return err
		}
		for _, addr := range bnttx.Addresses.Add {
			err = enc.Encode(addr)
			if err != nil {
				return err
			}
		}
		for _, addr := range bnttx.Addresses.Remove {
			err = enc.Encode(addr)
			if err != nil {
				return err
			}
		}
	}

	// encode names added and removed, if defined
	if maf.HasNames {
		err = enc.Encode(uint8(nameAddLen) | (uint8(nameRemoveLen) << 4))
		if err != nil {
			return err
		}
		for _, name := range bnttx.Names.Add {
			err = enc.Encode(name)
			if err != nil {
				return err
			}
		}
		for _, name := range bnttx.Names.Remove {
			err = enc.Encode(name)
			if err != nil {
				return err
			}
		}
	}

	// encode TxFee and CoinInputs
	err = enc.EncodeAll(bnttx.TransactionFee, bnttx.CoinInputs)
	if err != nil {
		return err
	}
	// encode refund coin output, if given
	if maf.HasRefund {
		// deref to ensure we do not also encode one byte
		// for the pointer indication
		err = enc.Encode(*bnttx.RefundCoinOutput)
		if err != nil {
			return err
		}
	}
	// nothing more to do
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (bnttx *BotNameTransferTransaction) UnmarshalSia(r io.Reader) error {
	dec := tfencoding.NewDecoder(r)

	// unmarshal sender, receiver, NrOfMonths and flags
	var maf BotMonthsAndFlagsData
	err := dec.DecodeAll(&bnttx.Sender, &bnttx.Receiver, &maf)
	if err != nil {
		return err
	}

	// assign number of months
	bnttx.NrOfMonths = maf.NrOfMonths

	// decode addressed added and removed, if defined
	if maf.HasAddresses {
		var pairLength uint8
		err = dec.Decode(&pairLength)
		if err != nil {
			return err
		}
		addrAddLen, addrRemoveLen := pairLength&15, pairLength>>4
		if addrAddLen > 0 {
			bnttx.Addresses.Add = make([]NetworkAddress, addrAddLen)
			for i := range bnttx.Addresses.Add {
				err = dec.Decode(&bnttx.Addresses.Add[i])
				if err != nil {
					return err
				}
			}
		} else {
			bnttx.Addresses.Add = nil
		}
		if addrRemoveLen > 0 {
			bnttx.Addresses.Remove = make([]NetworkAddress, addrRemoveLen)
			for i := range bnttx.Addresses.Remove {
				err = dec.Decode(&bnttx.Addresses.Remove[i])
				if err != nil {
					return err
				}
			}
		} else {
			bnttx.Addresses.Remove = nil
		}
	} else {
		// explicitly set added/removed address to nil
		bnttx.Addresses.Add, bnttx.Addresses.Remove = nil, nil
	}

	if maf.HasNames {
		var pairLength uint8
		err = dec.Decode(&pairLength)
		if err != nil {
			return err
		}
		nameAddLen, nameRemoveLen := pairLength&15, pairLength>>4
		if nameAddLen > 0 {
			bnttx.Names.Add = make([]BotName, nameAddLen)
			for i := range bnttx.Names.Add {
				err = dec.Decode(&bnttx.Names.Add[i])
				if err != nil {
					return err
				}
			}
		} else {
			bnttx.Names.Add = nil
		}
		if nameRemoveLen > 0 {
			bnttx.Names.Remove = make([]BotName, nameRemoveLen)
			for i := range bnttx.Names.Remove {
				err = dec.Decode(&bnttx.Names.Remove[i])
				if err != nil {
					return err
				}
			}
		} else {
			bnttx.Names.Remove = nil
		}
	} else {
		// explicitly set added/removed address to nil
		bnttx.Names.Add, bnttx.Names.Remove = nil, nil
	}

	// encode TxFee and CoinInputs
	err = dec.DecodeAll(&bnttx.TransactionFee, &bnttx.CoinInputs)
	if err != nil {
		return err
	}
	// decode refund coin output, if defined
	if maf.HasRefund {
		bnttx.RefundCoinOutput = new(types.CoinOutput)
		err = dec.Decode(bnttx.RefundCoinOutput)
		if err != nil {
			return err
		}
	}
	// nothing more to do
	return nil
}

type (
	// BotRecordReadRegistry defines the public READ API expected from a bot record Read-Only registry.
	BotRecordReadRegistry interface {
		// GetRecordForID returns the record mapped to the given BotID.
		GetRecordForID(id BotID) (*BotRecord, error)
		// GetRecordForKey returns the record mapped to the given Key.
		GetRecordForKey(key PublicKey) (*BotRecord, error)
		// GetRecordForName returns the record mapped to the given Name.
		GetRecordForName(name BotName) (*BotRecord, error)
	}
)

// public BotRecordReadRegistry errors
var (
	ErrBotNotFound     = errors.New("3bot not found")
	ErrBotKeyNotFound  = errors.New("3bot public key not found")
	ErrBotNameNotFound = errors.New("3bot name not found")
)

// 3bot Tx controllers

type (
	// BotRegistrationTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 0x90. It allows the registration of a new3bot.
	BotRegistrationTransactionController struct {
		Registry              BotRecordReadRegistry
		RegistryPoolCondition types.UnlockConditionProxy
		OneCoin               types.Currency
	}
)

var (
	// ensure at compile time that BotRegistrationTransactionController
	// implements the desired interfaces
	_ types.TransactionController      = BotRegistrationTransactionController{}
	_ types.TransactionExtensionSigner = BotRegistrationTransactionController{}
	_ types.TransactionValidator       = BotRegistrationTransactionController{}
	_ types.BlockStakeOutputValidator  = BotRegistrationTransactionController{}
	_ types.InputSigHasher             = BotRegistrationTransactionController{}
	_ types.TransactionIDEncoder       = BotRegistrationTransactionController{}
)

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (brtc BotRegistrationTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	brtx, err := BotRegistrationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a BotRegistrationTx: %v", err)
	}
	return tfencoding.NewEncoder(w).Encode(brtx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (brtc BotRegistrationTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var brtx BotRegistrationTransaction
	err := tfencoding.NewDecoder(r).Decode(&brtx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to binary-decode tx as a BotRegistrationTx: %v", err)
	}
	// return minter definition tx as regular tfchain tx data
	return brtx.TransactionData(brtc.OneCoin, brtc.RegistryPoolCondition), nil
}

// JSONEncodeTransactionData implements TransactionController.JSONEncodeTransactionData
func (brtc BotRegistrationTransactionController) JSONEncodeTransactionData(txData types.TransactionData) ([]byte, error) {
	brtx, err := BotRegistrationTransactionFromTransactionData(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert txData to a BotRegistrationTx: %v", err)
	}
	return json.Marshal(brtx)
}

// JSONDecodeTransactionData implements TransactionController.JSONDecodeTransactionData
func (brtc BotRegistrationTransactionController) JSONDecodeTransactionData(data []byte) (types.TransactionData, error) {
	var brtx BotRegistrationTransaction
	err := json.Unmarshal(data, &brtx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to json-decode tx as a BotRegistrationTx: %v", err)
	}
	// return minter definition tx as regular tfchain tx data
	return brtx.TransactionData(brtc.OneCoin, brtc.RegistryPoolCondition), nil
}

// SignExtension implements TransactionExtensionSigner.SignExtension
func (brtc BotRegistrationTransactionController) SignExtension(extension interface{}, sign func(*types.UnlockFulfillmentProxy, types.UnlockConditionProxy) error) (interface{}, error) {
	// (tx) extension (data) is expected to be a pointer to a valid BotRegistrationTransactionExtension
	brtxExtension, ok := extension.(*BotRegistrationTransactionExtension)
	if !ok {
		return nil, errors.New("invalid extension data for a BotRegistrationTx")
	}

	// create a publicKeyUnlockHashCondition
	spk, err := brtxExtension.Identification.PublicKey.SiaPublicKey()
	if err != nil {
		return nil, errors.New("invalid public key in extension data for a BotRegistrationTx")
	}
	condition := types.NewCondition(types.NewUnlockHashCondition(types.NewPubKeyUnlockHash(spk)))
	// and a matching single-signature fulfillment
	fulfillment := types.NewFulfillment(types.NewSingleSignatureFulfillment(spk))

	// sign the fulfillment
	err = sign(&fulfillment, condition)
	if err != nil {
		return nil, fmt.Errorf("failed to sign BotRegistrationTx: %v", err)
	}

	// extract signature
	signature := fulfillment.Fulfillment.(*types.SingleSignatureFulfillment).Signature
	brtxExtension.Identification.Signature = signature
	// and return the signed extension
	return brtxExtension, nil
}

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (brtc BotRegistrationTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) error {
	err := types.TransactionFitsInABlock(t, constants.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get BotRegistrationTx
	brtx, err := BotRegistrationTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to use tx as a bot registration tx: %v", err)
	}

	// look up the public key, to ensure it is not registered yet
	_, err = brtc.Registry.GetRecordForKey(brtx.Identification.PublicKey)
	if err == nil {
		return ErrBotNameAlreadyRegistered
	}
	if err != ErrBotKeyNotFound {
		return fmt.Errorf("unexpected error while validating non-existence of bot's public key: %v", err)
	}

	// create a publicKeyUnlockHashCondition
	err = validateBotSignature(t, brtx.Identification.PublicKey, brtx.Identification.Signature, ctx)
	if err != nil {
		return fmt.Errorf("failed to fulfill bot registration condition: %v", err)
	}

	// ensure the NrOfMonths is in the inclusive range of [1, 24]
	if brtx.NrOfMonths == 0 {
		return errors.New("bot registration requires at least one month to be paid already")
	}
	if brtx.NrOfMonths > MaxBotPrepaidMonths {
		return ErrBotExpirationExtendOverflow
	}

	// validate the lengths,
	// and ensure that at least one name or one addr is registered
	addrLen := len(brtx.Addresses)
	if addrLen > MaxAddressesPerBot {
		return ErrTooManyBotAddresses
	}
	nameLen := len(brtx.Names)
	if nameLen > MaxNamesPerBot {
		return ErrTooManyBotNames
	}
	if addrLen == 0 && nameLen == 0 {
		return errors.New("bot registration requires a name or address to be defined")
	}

	// validate that all network addresses are unique
	err = validateUniquenessOfNetworkAddresses(brtx.Addresses)
	if err != nil {
		return fmt.Errorf("invalid bot registration Tx: %v", err)
	}

	// validate that all names are unique
	err = validateUniquenessOfBotNames(brtx.Names)
	if err != nil {
		return fmt.Errorf("invalid bot registration Tx: %v", err)
	}

	// validate that the names are not registered yet
	for _, name := range brtx.Names {
		_, err = brtc.Registry.GetRecordForName(name)
		if err == nil {
			return ErrBotNameAlreadyRegistered
		}
		if err != ErrBotNameNotFound {
			return fmt.Errorf(
				"unexpected error while validating non-existence of bot's name %v: %v",
				name, err)
		}
	}

	// validate the miner fee
	if brtx.TransactionFee.Cmp(constants.MinimumMinerFee) == -1 {
		return types.ErrTooSmallMinerFee
	}
	return nil
}

// Rivine handles ValidateCoinOutputs,
// which is possible as all our coin inputs are standard,
// the (single) miner fee is standard as well, and
// the additional (bot) fee is seen by Rivine as a coin output to a hardcoded condition.

// ValidateBlockStakeOutputs implements BlockStakeOutputValidator.ValidateBlockStakeOutputs
func (brtc BotRegistrationTransactionController) ValidateBlockStakeOutputs(t types.Transaction, ctx types.FundValidationContext, blockStakeInputs map[types.BlockStakeOutputID]types.BlockStakeOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within a bot registration transaction
}

// InputSigHash implements InputSigHasher.InputSigHash
func (brtc BotRegistrationTransactionController) InputSigHash(t types.Transaction, _ uint64, extraObjects ...interface{}) (crypto.Hash, error) {
	brtx, err := BotRegistrationTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a BotRegistrationTx: %v", err)
	}

	h := crypto.NewHash()
	enc := tfencoding.NewEncoder(h)

	enc.EncodeAll(
		t.Version,
		SpecifierBotRegistrationTransaction,
	)

	if len(extraObjects) > 0 {
		enc.EncodeAll(extraObjects...)
	}

	enc.EncodeAll(
		brtx.Addresses,
		brtx.Names,
		brtx.NrOfMonths,
	)

	enc.Encode(len(brtx.CoinInputs))
	for _, ci := range brtx.CoinInputs {
		enc.Encode(ci.ParentID)
	}

	enc.EncodeAll(
		brtx.TransactionFee,
		brtx.RefundCoinOutput,
		brtx.Identification.PublicKey,
	)

	var hash crypto.Hash
	h.Sum(hash[:0])
	return hash, nil
}

// EncodeTransactionIDInput implements TransactionIDEncoder.EncodeTransactionIDInput
func (brtc BotRegistrationTransactionController) EncodeTransactionIDInput(w io.Writer, txData types.TransactionData) error {
	brtx, err := BotRegistrationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a BotRegistrationTx: %v", err)
	}
	return tfencoding.NewEncoder(w).EncodeAll(SpecifierBotRegistrationTransaction, brtx)
}

type (
	// BotUpdateRecordTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 0x91. It allows the update of the record of an existing 3bot.
	BotUpdateRecordTransactionController struct {
		registry BotRecordReadRegistry
		oneCoin  types.Currency
	}
)

// TODO: implement BotUpdateRecordTransactionController controller

type (
	// BotNameTransferTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 0x92. It allows the transfer of names and update of the record
	// of the two existing 3bot that participate in this transfer.
	BotNameTransferTransactionController struct {
		registry BotRecordReadRegistry
		oneCoin  types.Currency
	}
)

// TODO: implement BotNameTransferTransactionController controller

// ComputeMonthlyBotFees computes the total monthly fees required for the given months,
// using the given oneCoin value as the currency's unit value.
func ComputeMonthlyBotFees(months uint8, oneCoin types.Currency) types.Currency {
	multiplier := uint64(months) * BotMonthlyFeeMultiplier
	if months < 12 {
		// return plain monthly fees without discounts
		return oneCoin.Mul64(multiplier)
	}
	fees := big.NewFloat(float64(multiplier))
	fees.Mul(fees, new(big.Float).SetInt(oneCoin.Big()))
	if months < 24 {
		// return plain monthly fees with 30% discount applied to the total
		i, _ := fees.Mul(fees, big.NewFloat(0.7)).Int(nil)
		return types.NewCurrency(i)
	}
	// return plain monthly fees with 50% discount applied to the total
	i, _ := fees.Mul(fees, big.NewFloat(0.5)).Int(nil)
	return types.NewCurrency(i)
}

// BotMonthsAndFlagsData is a utility structure that is used to encode
// the NrOfMonths (paid up front for a 3bot) as well as several flags
// in a single byte.
type BotMonthsAndFlagsData struct {
	NrOfMonths   uint8
	HasAddresses bool
	HasNames     bool
	HasRefund    bool
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (maf BotMonthsAndFlagsData) MarshalSia(w io.Writer) error {
	x := uint8(maf.NrOfMonths)
	if maf.HasAddresses {
		x |= 32
	}
	if maf.HasNames {
		x |= 64
	}
	if maf.HasRefund {
		x |= 128
	}
	return tfencoding.MarshalUint8(w, x)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (maf *BotMonthsAndFlagsData) UnmarshalSia(r io.Reader) error {
	x, err := tfencoding.UnmarshalUint8(r)
	if err != nil {
		return err
	}
	maf.NrOfMonths = x & 31
	maf.HasAddresses = ((x & 32) != 0)
	maf.HasNames = ((x & 64) != 0)
	maf.HasRefund = ((x & 128) != 0)
	return nil
}

// TransactionNonce is a nonce
// used to ensure the uniqueness of an otherwise potentially non-unique Tx
type TransactionNonce [TransactionNonceLength]byte

// TransactionNonceLength defines the length of a TransactionNonce
const TransactionNonceLength = 8

// RandomTransactionNonce creates a random Transaction nonce
func RandomTransactionNonce() (nonce TransactionNonce) {
	for nonce == (TransactionNonce{}) {
		// generate non-nil crypto-Random TransactionNonce
		rand.Read(nonce[:])
	}
	return
}

// MarshalJSON implements JSON.Marshaller.MarshalJSON
// encodes the Nonce as a base64-encoded string
func (tn TransactionNonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(tn[:])
}

// UnmarshalJSON implements JSON.Unmarshaller.UnmarshalJSON
// piggy-backing on the base64-decoding used for byte slices in the std JSON lib
func (tn *TransactionNonce) UnmarshalJSON(in []byte) error {
	var out []byte
	err := json.Unmarshal(in, &out)
	if err != nil {
		return err
	}
	if len(out) != TransactionNonceLength {
		return errors.New("invalid tx nonce length")
	}
	copy(tn[:], out[:])
	return nil
}

func unlockHashFromHex(hstr string) (uh types.UnlockHash) {
	err := uh.LoadString(hstr)
	if err != nil {
		panic(fmt.Sprintf("func unlockHashFromHex(%s) failed: %v", hstr, err))
	}
	return
}
