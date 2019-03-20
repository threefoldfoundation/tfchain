package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/threefoldfoundation/tfchain/pkg/config"

	"github.com/threefoldtech/rivine/build"
	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
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
	ERC20Registry
}

// RegisterTransactionTypesForStandardNetwork registers he transaction controllers
// for all transaction versions supported on the standard network.
func RegisterTransactionTypesForStandardNetwork(db TFChainReadDB, erc20TxValidator ERC20TransactionValidator, oneCoin types.Currency, cfg config.DaemonNetworkConfig) {
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
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(TransactionVersionBotRecordUpdate, BotUpdateRecordTransactionController{
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(TransactionVersionBotNameTransfer, BotNameTransferTransactionController{
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})

	types.RegisterTransactionVersion(TransactionVersionERC20Conversion, ERC20ConvertTransactionController{})
	types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, ERC20CoinCreationTransactionController{
		Registry:    db,
		OneCoin:     oneCoin,
		TxValidator: erc20TxValidator,
	})
	types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, ERC20AddressRegistrationTransactionController{
		Registry:             db,
		OneCoin:              oneCoin,
		BridgeFeePoolAddress: cfg.ERC20FeePoolAddress,
	})
}

// RegisterTransactionTypesForTestNetwork registers he transaction controllers
// for all transaction versions supported on the test network.
func RegisterTransactionTypesForTestNetwork(db TFChainReadDB, erc20TxValidator ERC20TransactionValidator, oneCoin types.Currency, cfg config.DaemonNetworkConfig) {
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
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(TransactionVersionBotRecordUpdate, BotUpdateRecordTransactionController{
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(TransactionVersionBotNameTransfer, BotNameTransferTransactionController{
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})

	types.RegisterTransactionVersion(TransactionVersionERC20Conversion, ERC20ConvertTransactionController{})
	types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, ERC20CoinCreationTransactionController{
		Registry:    db,
		OneCoin:     oneCoin,
		TxValidator: erc20TxValidator,
	})
	types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, ERC20AddressRegistrationTransactionController{
		Registry:             db,
		OneCoin:              oneCoin,
		BridgeFeePoolAddress: cfg.ERC20FeePoolAddress,
	})
}

// RegisterTransactionTypesForDevNetwork registers he transaction controllers
// for all transaction versions supported on the dev network.
func RegisterTransactionTypesForDevNetwork(db TFChainReadDB, erc20TxValidator ERC20TransactionValidator, oneCoin types.Currency, cfg config.DaemonNetworkConfig) {
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
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(TransactionVersionBotRecordUpdate, BotUpdateRecordTransactionController{
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(TransactionVersionBotNameTransfer, BotNameTransferTransactionController{
		Registry:            db,
		RegistryPoolAddress: cfg.FoundationPoolAddress,
		OneCoin:             oneCoin,
	})

	types.RegisterTransactionVersion(TransactionVersionERC20Conversion, ERC20ConvertTransactionController{})
	types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, ERC20CoinCreationTransactionController{
		Registry:    db,
		OneCoin:     oneCoin,
		TxValidator: erc20TxValidator,
	})
	types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, ERC20AddressRegistrationTransactionController{
		Registry:             db,
		OneCoin:              oneCoin,
		BridgeFeePoolAddress: cfg.ERC20FeePoolAddress,
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
	_ types.TransactionController      = LegacyTransactionController{}
	_ types.TransactionValidator       = LegacyTransactionController{}
	_ types.TransactionSignatureHasher = LegacyTransactionController{}
	_ types.TransactionIDEncoder       = LegacyTransactionController{}

	// ensure at compile time that CoinCreationTransactionController
	// implements the desired interfaces
	_ types.TransactionController      = CoinCreationTransactionController{}
	_ types.TransactionExtensionSigner = CoinCreationTransactionController{}
	_ types.TransactionValidator       = CoinCreationTransactionController{}
	_ types.CoinOutputValidator        = CoinCreationTransactionController{}
	_ types.BlockStakeOutputValidator  = CoinCreationTransactionController{}
	_ types.TransactionSignatureHasher = CoinCreationTransactionController{}
	_ types.TransactionIDEncoder       = CoinCreationTransactionController{}

	// ensure at compile time that MinterDefinitionTransactionController
	// implements the desired interfaces
	_ types.TransactionController                = MinterDefinitionTransactionController{}
	_ types.TransactionExtensionSigner           = MinterDefinitionTransactionController{}
	_ types.TransactionValidator                 = MinterDefinitionTransactionController{}
	_ types.CoinOutputValidator                  = MinterDefinitionTransactionController{}
	_ types.BlockStakeOutputValidator            = MinterDefinitionTransactionController{}
	_ types.TransactionSignatureHasher           = MinterDefinitionTransactionController{}
	_ types.TransactionIDEncoder                 = MinterDefinitionTransactionController{}
	_ types.TransactionCommonExtensionDataGetter = MinterDefinitionTransactionController{}
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
	return siabin.NewEncoder(w).Encode(cctx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (cctc CoinCreationTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var cctx CoinCreationTransaction
	err := siabin.NewDecoder(r).Decode(&cctx)
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
func (cctc CoinCreationTransactionController) SignExtension(extension interface{}, sign func(*types.UnlockFulfillmentProxy, types.UnlockConditionProxy, ...interface{}) error) (interface{}, error) {
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

// SignatureHash implements TransactionSignatureHasher.SignatureHash
func (cctc CoinCreationTransactionController) SignatureHash(t types.Transaction, extraObjects ...interface{}) (crypto.Hash, error) {
	cctx, err := CoinCreationTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a coin creation tx: %v", err)
	}

	h := crypto.NewHash()
	enc := siabin.NewEncoder(h)

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
	return siabin.NewEncoder(w).EncodeAll(SpecifierCoinCreationTransaction, cctx)
}

// MinterDefinitionTransactionController

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (mdtc MinterDefinitionTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	mdtx, err := MinterDefinitionTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a MinterDefinitionTx: %v", err)
	}
	return siabin.NewEncoder(w).Encode(mdtx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (mdtc MinterDefinitionTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var mdtx MinterDefinitionTransaction
	err := siabin.NewDecoder(r).Decode(&mdtx)
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
func (mdtc MinterDefinitionTransactionController) SignExtension(extension interface{}, sign func(*types.UnlockFulfillmentProxy, types.UnlockConditionProxy, ...interface{}) error) (interface{}, error) {
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

// SignatureHash implements TransactionSignatureHasher.SignatureHash
func (mdtc MinterDefinitionTransactionController) SignatureHash(t types.Transaction, extraObjects ...interface{}) (crypto.Hash, error) {
	mdtx, err := MinterDefinitionTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a MinterDefinitionTx: %v", err)
	}

	h := crypto.NewHash()
	enc := siabin.NewEncoder(h)

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
	return siabin.NewEncoder(w).EncodeAll(SpecifierMintDefinitionTransaction, mdtx)
}

// GetCommonExtensionData implements TransactionCommonExtensionDataGetter.GetCommonExtensionData
func (mdtc MinterDefinitionTransactionController) GetCommonExtensionData(extension interface{}) (types.CommonTransactionExtensionData, error) {
	mdext, ok := extension.(*MinterDefinitionTransactionExtension)
	if !ok {
		return types.CommonTransactionExtensionData{}, errors.New("invalid extension data for a MinterDefinitionTx")
	}
	return types.CommonTransactionExtensionData{
		UnlockConditions: []types.UnlockConditionProxy{mdext.MintCondition},
	}, nil
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

func unlockHashFromHex(hstr string) (uh types.UnlockHash) {
	err := uh.LoadString(hstr)
	if err != nil {
		panic(fmt.Sprintf("func unlockHashFromHex(%s) failed: %v", hstr, err))
	}
	return
}
