package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
)

const (
	// TransactionVersionERC20Conversion defines the Transaction version
	// for an ERC20ConvertTransaction, used to convert TFT into ERC20 funds.
	TransactionVersionERC20Conversion types.TransactionVersion = iota + 208
	// TransactionVersionERC20CoinCreation defines the Transaction version
	// for an ERC20CoinCreationTransaction, used to convert ERC20 funds into TFT.
	TransactionVersionERC20CoinCreation
	// TransactionVersionERC20AddressRegistration defines the Transaction version
	// for an TransactionVersionERC20AddressRegistration, used to register an ERC20 address,
	// linked to an TFT address.
	TransactionVersionERC20AddressRegistration
)

// These Specifiers are used internally when calculating a Transaction's ID.
// See Rivine's Specifier for more details.
var (
	SpecifierERC20ConvertTransaction             = types.Specifier{'e', 'r', 'c', '2', '0', ' ', 'c', 'o', 'n', 'v', 'e', 'r', 't', ' ', 't', 'x'}
	SpecifierERC20CoinCreationTransaction        = types.Specifier{'e', 'r', 'c', '2', '0', ' ', 'c', 'o', 'i', 'n', 'g', 'e', 'n', ' ', 't', 'x'}
	SpecifierERC20AddressRegistrationTransaction = types.Specifier{'e', 'r', 'c', '2', '0', ' ', 'a', 'd', 'd', 'r', 'r', 'e', 'g', ' ', 't', 'x'}
)

var (
	// ERC20ConversionMinimumValue defines the minimum value of TFT
	// you can convert to ERC20 funds using the ERC20ConvertTransaction
	ERC20ConversionMinimumValue = config.GetCurrencyUnits().OneCoin.Mul64(1000)
)

// ERC20AddressLength defines the length of the fixed-sized ERC20Address type explicitly.
const ERC20AddressLength = 20

// ERC20Address defines an ERC20 address as a fixed-sized byte array of length 20,
// and is used in order to be able to convert TFT into tradeable tfchain ERC20 funds.
type ERC20Address [ERC20AddressLength]byte

// String returns this ERC20Address as a string.
func (address ERC20Address) String() string {
	return hex.EncodeToString(address[:])
}

// LoadString loads this ERC20Address from a hex-encoded string of length 40.
func (address *ERC20Address) LoadString(str string) error {
	if str == "" {
		*address = ERC20Address{}
		return nil
	}
	if len(str) != ERC20AddressLength*2 {
		return errors.New("passed string cannot be loaded as an ERC20Address: invalid length")
	}
	n, err := hex.Decode(address[:], []byte(str))
	if err != nil {
		return err
	}
	if n != ERC20AddressLength {
		return io.ErrShortWrite
	}
	return nil
}

// MarshalJSON implements json.Marshaler.MarshalJSON,
// and returns this ERC20Address as a hex-encoded JSON string.
func (address ERC20Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(address.String())
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON,
// and decodes the given byte slice as a hex-encoded JSON string into the
// 20 bytes that make up this ERC20Address.
func (address *ERC20Address) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return address.LoadString(str)
}

type (
	// ERC20Registry defines the public READ API expected from an ERC20 Read-Only registry.
	ERC20Registry interface {
		GetERC20AddressForTFTAddress(types.UnlockHash) (ERC20Address, bool, error)
		GetTFTTransactionIDForERC20TransactionID(ERC20Hash) (types.TransactionID, bool, error)
	}

	// ERC20TransactionValidator is the validation API used by the ERC20 CoinCreation Tx Controller,
	// in order to validate the attached ERC20 Tx. Use the NopERC20TransactionValidator if no such validation is required.
	ERC20TransactionValidator interface {
		ValidateWithdrawTx(blockID, txID ERC20Hash, expectedAddress ERC20Address, expecedAmount types.Currency) error
	}
)

// NopERC20TransactionValidator provides a NOP-implementation of the ERC20TransactionValidator interface,
// allowing you to disable any extra validation on ERC20 Transactions.
type NopERC20TransactionValidator struct{}

// ValidateWithdrawTx implements ERC20TransactionValidator.ValidateWithdrawTx,
// returning nil for every call.
func (nop NopERC20TransactionValidator) ValidateWithdrawTx(ERC20Hash, ERC20Hash, ERC20Address, types.Currency) error {
	return nil
}

type (
	// ERC20ConvertTransaction defines the Transaction (with version 0xD1)
	// used to convert TFT into ERC20 funds paid to the defined ERC20 address.
	ERC20ConvertTransaction struct {
		// The address to send the TFT-converted tfchain ERC20 funds into.
		Address ERC20Address `json:"address"`

		// Amount of TFT to be paid towards buying ERC20 funds,
		// note that the bridge will take part of this amount towards
		// paying for the transaction costs, prior to sending the ERC20 funds to
		// the defined target address.
		Value types.Currency `json:"value"`

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency `json:"txfee"`

		// CoinInputs are only used for the required fees,
		// which contains the regular Tx fee as well as the additional fees,
		// to be paid for the address registration. At least one CoinInput is required.
		CoinInputs []types.CoinInput `json:"coininputs"`
		// RefundCoinOutput is an optional coin output that can be used
		// to refund coins paid as inputs for the required fees.
		RefundCoinOutput *types.CoinOutput `json:"refundcoinoutput,omitempty"`
	}

	// ERC20ConvertTransactionExtension defines the ERC20ConvertTransaction Extension Data
	ERC20ConvertTransactionExtension struct {
		// The address to send the TFT-converted tfchain ERC20 funds into.
		Address ERC20Address
		// Amount of TFT to be paid towards buying ERC20 funds.
		Value types.Currency
	}
)

// ERC20ConvertTransactionFromTransaction creates an ERC20ConvertTransaction,
// using a regular in-memory tfchain transaction.
//
// Past the (tx) Version validation it piggy-backs onto the
// `ERC20ConvertTransactionFromTransactionData` constructor.
func ERC20ConvertTransactionFromTransaction(tx types.Transaction) (ERC20ConvertTransaction, error) {
	if tx.Version != TransactionVersionERC20Conversion {
		return ERC20ConvertTransaction{}, fmt.Errorf(
			"an ERC20 convert transaction requires tx version %d",
			TransactionVersionERC20Conversion)
	}
	return ERC20ConvertTransactionFromTransactionData(types.TransactionData{
		CoinInputs:        tx.CoinInputs,
		CoinOutputs:       tx.CoinOutputs,
		BlockStakeInputs:  tx.BlockStakeInputs,
		BlockStakeOutputs: tx.BlockStakeOutputs,
		MinerFees:         tx.MinerFees,
		ArbitraryData:     tx.ArbitraryData,
		Extension:         tx.Extension,
	})
}

// ERC20ConvertTransactionFromTransactionData creates an ERC20ConvertTransaction,
// using the TransactionData from a regular in-memory tfchain transaction.
func ERC20ConvertTransactionFromTransactionData(txData types.TransactionData) (ERC20ConvertTransaction, error) {
	// validate the Transaction Data

	// at least one coin input as well as one miner fee is required
	if len(txData.CoinInputs) == 0 || len(txData.MinerFees) != 1 {
		return ERC20ConvertTransaction{}, errors.New("at least one coin input and exactly one miner fee is required for an ERC20 Convert Transaction")
	}
	// no block stake inputs or block stake outputs are allowed
	if len(txData.BlockStakeInputs) != 0 || len(txData.BlockStakeOutputs) != 0 {
		return ERC20ConvertTransaction{}, errors.New("no block stake inputs/outputs are allowed in an ERC20 Convert Transaction")
	}
	// no arbitrary data is allowed
	if len(txData.ArbitraryData.Data) > 0 {
		return ERC20ConvertTransaction{}, errors.New("no arbitrary data is allowed in an ERC20 Convert Transaction")
	}
	// validate that the coin outputs is within the expected range
	if len(txData.CoinOutputs) > 1 {
		return ERC20ConvertTransaction{}, errors.New("an ERC20 Convert Transaction can only have one coin output")
	}

	// (tx) extension (data) is expected to be a pointer to a valid ERC20ConvertTransaction,
	// which contains all the properties unique to a 3bot (name transfer) Tx
	extensionData, ok := txData.Extension.(*ERC20ConvertTransactionExtension)
	if !ok {
		return ERC20ConvertTransaction{}, errors.New("invalid extension data for an ERC20 Convert Transaction")
	}

	// create the ERC20ConvertTransaction and return it,
	// further validation will/has-to be done using the Transaction Type, if required
	tx := ERC20ConvertTransaction{
		Address:        extensionData.Address,
		Value:          extensionData.Value,
		TransactionFee: txData.MinerFees[0],
		CoinInputs:     txData.CoinInputs,
	}
	if len(txData.CoinOutputs) == 1 {
		// take refund coin output if it exists
		tx.RefundCoinOutput = &txData.CoinOutputs[0]
	}
	return tx, nil
}

// TransactionData returns this ERC20ConvertTransaction
// as regular tfchain transaction data.
func (etctx *ERC20ConvertTransaction) TransactionData() types.TransactionData {
	txData := types.TransactionData{
		CoinInputs: etctx.CoinInputs,
		MinerFees:  []types.Currency{etctx.TransactionFee},
		Extension: &ERC20ConvertTransactionExtension{
			Address: etctx.Address,
			Value:   etctx.Value,
		},
	}
	if etctx.RefundCoinOutput != nil {
		txData.CoinOutputs = append(txData.CoinOutputs, *etctx.RefundCoinOutput)
	}
	return txData
}

// Transaction returns this ERC20ConvertTransaction
// as regular tfchain transaction, using TransactionVersionBotNameTransfer as the type.
func (etctx *ERC20ConvertTransaction) Transaction() types.Transaction {
	tx := types.Transaction{
		Version:    TransactionVersionERC20Conversion,
		CoinInputs: etctx.CoinInputs,
		MinerFees:  []types.Currency{etctx.TransactionFee},
		Extension: &ERC20ConvertTransactionExtension{
			Address: etctx.Address,
			Value:   etctx.Value,
		},
	}
	if etctx.RefundCoinOutput != nil {
		tx.CoinOutputs = append(tx.CoinOutputs, *etctx.RefundCoinOutput)
	}
	return tx
}

// MarshalSia implements SiaMarshaler.MarshalSia,
// alias of MarshalRivine for backwards-compatibility reasons.
func (etctx ERC20ConvertTransaction) MarshalSia(w io.Writer) error {
	return etctx.MarshalRivine(w)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// alias of UnmarshalRivine for backwards-compatibility reasons.
func (etctx *ERC20ConvertTransaction) UnmarshalSia(r io.Reader) error {
	return etctx.UnmarshalRivine(r)
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (etctx ERC20ConvertTransaction) MarshalRivine(w io.Writer) error {
	return rivbin.NewEncoder(w).EncodeAll(
		etctx.Address,
		etctx.Value,
		etctx.TransactionFee,
		etctx.CoinInputs,
		etctx.RefundCoinOutput,
	)
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (etctx *ERC20ConvertTransaction) UnmarshalRivine(r io.Reader) error {
	return rivbin.NewDecoder(r).DecodeAll(
		&etctx.Address,
		&etctx.Value,
		&etctx.TransactionFee,
		&etctx.CoinInputs,
		&etctx.RefundCoinOutput,
	)
}

type (
	// ERC20ConvertTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 0xD0. It allows the conversion of TFT to ERC20-funds.
	ERC20ConvertTransactionController struct{}
)

var (
	// ensure at compile time that ERC20ConvertTransactionController
	// implements the desired interfaces
	_ types.TransactionController      = ERC20ConvertTransactionController{}
	_ types.TransactionValidator       = ERC20ConvertTransactionController{}
	_ types.CoinOutputValidator        = ERC20ConvertTransactionController{}
	_ types.BlockStakeOutputValidator  = ERC20ConvertTransactionController{}
	_ types.TransactionSignatureHasher = ERC20ConvertTransactionController{}
	_ types.TransactionIDEncoder       = ERC20ConvertTransactionController{}
)

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (etctc ERC20ConvertTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	etctx, err := ERC20ConvertTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to an ERC20ConvertTx: %v", err)
	}
	return rivbin.NewEncoder(w).Encode(etctx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (etctc ERC20ConvertTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var etctx ERC20ConvertTransaction
	err := rivbin.NewDecoder(r).Decode(&etctx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to binary-decode tx as an ERC20ConvertTx: %v", err)
	}
	// return ERC20 convert tx as regular tfchain tx data
	return etctx.TransactionData(), nil
}

// JSONEncodeTransactionData implements TransactionController.JSONEncodeTransactionData
func (etctc ERC20ConvertTransactionController) JSONEncodeTransactionData(txData types.TransactionData) ([]byte, error) {
	etctx, err := ERC20ConvertTransactionFromTransactionData(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert txData to an ERC20ConvertTx: %v", err)
	}
	return json.Marshal(etctx)
}

// JSONDecodeTransactionData implements TransactionController.JSONDecodeTransactionData
func (etctc ERC20ConvertTransactionController) JSONDecodeTransactionData(data []byte) (types.TransactionData, error) {
	var etctx ERC20ConvertTransaction
	err := json.Unmarshal(data, &etctx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to json-decode tx as an ERC20ConvertTx: %v", err)
	}
	// return ERC20 convert tx as regular tfchain tx data
	return etctx.TransactionData(), nil
}

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (etctc ERC20ConvertTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) error {
	// check tx fits within a block
	err := types.TransactionFitsInABlock(t, constants.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get ERC20ConvertTx
	etctx, err := ERC20ConvertTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to use tx as an ERC20 convert tx: %v", err)
	}

	// ensure the value is a valid minimum
	if etctx.Value.Cmp(ERC20ConversionMinimumValue) < 0 {
		return errors.New("ERC20 requires a minimum value of 1000 TFT to be converted")
	}

	// validate the miner fee
	if etctx.TransactionFee.Cmp(constants.MinimumMinerFee) < 0 {
		return types.ErrTooSmallMinerFee
	}

	// prevent double spending
	spendCoins := make(map[types.CoinOutputID]struct{})
	for _, ci := range t.CoinInputs {
		if _, found := spendCoins[ci.ParentID]; found {
			return types.ErrDoubleSpend
		}
		spendCoins[ci.ParentID] = struct{}{}
	}

	// check if optional coin output is using standard condition
	if etctx.RefundCoinOutput != nil {
		err = etctx.RefundCoinOutput.Condition.IsStandardCondition(ctx)
		if err != nil {
			return err
		}
		// ensure the value is not 0
		if etctx.RefundCoinOutput.Value.IsZero() {
			return types.ErrZeroOutput
		}
	}
	// check if all fulfillments are standard
	for _, sci := range etctx.CoinInputs {
		err = sci.Fulfillment.IsStandardFulfillment(ctx)
		if err != nil {
			return err
		}
	}

	// Tx is valid
	return nil
}

// ValidateCoinOutputs implements CoinOutputValidator.ValidateCoinOutputs,
// implemented here, overwriting the default logic, as the Tx value is not registered as a coin output,
// instead those TFT are "burned"
func (etctc ERC20ConvertTransactionController) ValidateCoinOutputs(t types.Transaction, ctx types.FundValidationContext, coinInputs map[types.CoinOutputID]types.CoinOutput) error {
	etctx, err := ERC20ConvertTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to convert Tx to an ERC20ConvertTx: %v", err)
	}

	var inputSum types.Currency
	for index, sci := range etctx.CoinInputs {
		sco, ok := coinInputs[sci.ParentID]
		if !ok {
			return types.MissingCoinOutputError{ID: sci.ParentID}
		}
		// check if the referenced output's condition has been fulfilled
		err = sco.Condition.Fulfill(sci.Fulfillment, types.FulfillContext{
			ExtraObjects: []interface{}{uint64(index)},
			BlockHeight:  ctx.BlockHeight,
			BlockTime:    ctx.BlockTime,
			Transaction:  t,
		})
		if err != nil {
			return err
		}
		inputSum = inputSum.Add(sco.Value)
	}

	expectedTotalFee := etctx.TransactionFee.Add(etctx.Value)
	if etctx.RefundCoinOutput != nil {
		expectedTotalFee = expectedTotalFee.Add(etctx.RefundCoinOutput.Value)
	}
	if !inputSum.Equals(expectedTotalFee) {
		return types.ErrCoinInputOutputMismatch
	}
	return nil
}

// ValidateBlockStakeOutputs implements BlockStakeOutputValidator.ValidateBlockStakeOutputs
func (etctc ERC20ConvertTransactionController) ValidateBlockStakeOutputs(t types.Transaction, ctx types.FundValidationContext, blockStakeInputs map[types.BlockStakeOutputID]types.BlockStakeOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within an ERC20 convert tx transaction
}

// SignatureHash implements TransactionSignatureHasher.SignatureHash
func (etctc ERC20ConvertTransactionController) SignatureHash(t types.Transaction, extraObjects ...interface{}) (crypto.Hash, error) {
	etctx, err := ERC20ConvertTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a ERC20ConvertTx: %v", err)
	}

	h := crypto.NewHash()
	enc := rivbin.NewEncoder(h)

	enc.EncodeAll(
		t.Version,
		SpecifierERC20ConvertTransaction,
		etctx.Address,
		etctx.Value,
	)

	if len(extraObjects) > 0 {
		enc.EncodeAll(extraObjects...)
	}

	enc.Encode(len(etctx.CoinInputs))
	for _, ci := range etctx.CoinInputs {
		enc.Encode(ci.ParentID)
	}

	enc.EncodeAll(
		etctx.TransactionFee,
		etctx.RefundCoinOutput,
	)

	var hash crypto.Hash
	h.Sum(hash[:0])
	return hash, nil
}

// EncodeTransactionIDInput implements TransactionIDEncoder.EncodeTransactionIDInput
func (etctc ERC20ConvertTransactionController) EncodeTransactionIDInput(w io.Writer, txData types.TransactionData) error {
	etctx, err := ERC20ConvertTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a ERC20ConvertTx: %v", err)
	}
	return rivbin.NewEncoder(w).EncodeAll(SpecifierERC20ConvertTransaction, etctx)
}

// ERC20HashLength defines the length of the fixed-sized ERC20 Hash type explicitly,
// used for Transaction and Block hashes alike.
const ERC20HashLength = 32

// ERC20Hash defines an ERC20 Hash as a fixed-sized byte array of length 32.
type ERC20Hash [ERC20HashLength]byte

// String returns this TransactionID as a string.
func (eh ERC20Hash) String() string {
	return hex.EncodeToString(eh[:])
}

// LoadString loads this TransactionID from a hex-encoded string of length 40.
func (eh *ERC20Hash) LoadString(str string) error {
	if len(str) != ERC20HashLength*2 {
		return errors.New("passed string cannot be loaded as an ERC20 Hash: invalid length")
	}
	n, err := hex.Decode(eh[:], []byte(str))
	if err != nil {
		return err
	}
	if n != ERC20HashLength {
		return io.ErrShortWrite
	}
	return nil
}

// MarshalJSON implements json.Marshaler.MarshalJSON,
// and returns this Hash as a hex-encoded JSON string.
func (eh ERC20Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(eh.String())
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON,
// and decodes the given byte slice as a hex-encoded JSON string into the
// 20 bytes that make up this Hash.
func (eh *ERC20Hash) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return eh.LoadString(str)
}

type (
	// ERC20CoinCreationTransaction defines the Transaction (with version 0xD1)
	// used to convert ERC20 funds into TFT (the reverse of the ERC20ConvertTransaction).
	ERC20CoinCreationTransaction struct {
		// The address to send the TFT-converted tfchain ERC20 funds into.
		Address types.UnlockHash `json:"address"`

		// Amount of TFT to be paid towards buying ERC20 funds,
		// note that the bridge will take part of this amount towards
		// paying for the transaction costs, prior to sending the ERC20 funds to
		// the defined target address.
		Value types.Currency `json:"value"`

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency `json:"txfee"`

		// ERC20 BlockID (Sending ERC20 Funds to TFT) used as to identify
		// the parent block of the source of this coin creation.
		BlockID ERC20Hash `json:"blockid"`

		// ERC20 TransactionID (Sending ERC20 Funds to TFT) used as the source of this coin creation.
		TransactionID ERC20Hash `json:"txid"`
	}

	// ERC20CoinCreationTransactionExtension defines the ERC20CoinCreationTransaction Extension Data
	ERC20CoinCreationTransactionExtension struct {
		BlockID       ERC20Hash
		TransactionID ERC20Hash
	}
)

// ERC20CoinCreationTransactionFromTransaction creates an ERC20CoinCreationTransaction,
// using a regular in-memory tfchain transaction.
//
// Past the (tx) Version validation it piggy-backs onto the
// `ERC20CoinCreationTransactionFromTransactionData` constructor.
func ERC20CoinCreationTransactionFromTransaction(tx types.Transaction) (ERC20CoinCreationTransaction, error) {
	if tx.Version != TransactionVersionERC20CoinCreation {
		return ERC20CoinCreationTransaction{}, fmt.Errorf(
			"an ERC20 CoinCreation transaction requires tx version %d",
			TransactionVersionERC20CoinCreation)
	}
	return ERC20CoinCreationTransactionFromTransactionData(types.TransactionData{
		CoinInputs:        tx.CoinInputs,
		CoinOutputs:       tx.CoinOutputs,
		BlockStakeInputs:  tx.BlockStakeInputs,
		BlockStakeOutputs: tx.BlockStakeOutputs,
		MinerFees:         tx.MinerFees,
		ArbitraryData:     tx.ArbitraryData,
		Extension:         tx.Extension,
	})
}

// ERC20CoinCreationTransactionFromTransactionData creates an ERC20CoinCreationTransaction,
// using the TransactionData from a regular in-memory tfchain transaction.
func ERC20CoinCreationTransactionFromTransactionData(txData types.TransactionData) (ERC20CoinCreationTransaction, error) {
	// validate the Transaction Data

	// no coin inputs are allowed
	if len(txData.CoinInputs) != 0 {
		return ERC20CoinCreationTransaction{}, errors.New("no coin inputs are allowed in an ERC20 CoinCreation Tx")
	}
	// exactly one miner fee is required and expected
	if len(txData.MinerFees) != 1 {
		return ERC20CoinCreationTransaction{}, errors.New("exactly one miner fee is required for an ERC20 CoinCreation Transaction")
	}
	// no block stake inputs or block stake outputs are allowed
	if len(txData.BlockStakeInputs) != 0 || len(txData.BlockStakeOutputs) != 0 {
		return ERC20CoinCreationTransaction{}, errors.New("no block stake inputs/outputs are allowed in an ERC20 CoinCreation Transaction")
	}
	// no arbitrary data is allowed
	if len(txData.ArbitraryData.Data) > 0 {
		return ERC20CoinCreationTransaction{}, errors.New("no arbitrary data is allowed in an ERC20 CoinCreation Transaction")
	}
	// validate that we only have one coin output
	if len(txData.CoinOutputs) != 1 {
		return ERC20CoinCreationTransaction{}, errors.New("an ERC20 CoinCreation Transaction has to have exactlyone coin output")
	}

	// (tx) extension (data) is expected to be a pointer to a valid BotNameTransferTransaction,
	// which contains all the properties unique to a 3bot (name transfer) Tx
	extensionData, ok := txData.Extension.(*ERC20CoinCreationTransactionExtension)
	if !ok {
		return ERC20CoinCreationTransaction{}, errors.New("invalid extension data for an ERC20 CoinCreation Transaction")
	}

	// create the ERC20CoinCreationTransaction and return it,
	// further validation will/has-to be done using the Transaction Type, if required
	co := txData.CoinOutputs[0]
	return ERC20CoinCreationTransaction{
		Address:        co.Condition.UnlockHash(),
		Value:          co.Value,
		TransactionFee: txData.MinerFees[0],
		BlockID:        extensionData.BlockID,
		TransactionID:  extensionData.TransactionID,
	}, nil
}

// TransactionData returns this ERC20CoinCreationTransaction
// as regular tfchain transaction data.
func (etctx *ERC20CoinCreationTransaction) TransactionData() types.TransactionData {
	return types.TransactionData{
		CoinOutputs: []types.CoinOutput{
			{
				Condition: types.NewCondition(types.NewUnlockHashCondition(etctx.Address)),
				Value:     etctx.Value,
			},
		},
		MinerFees: []types.Currency{etctx.TransactionFee},
		Extension: &ERC20CoinCreationTransactionExtension{
			BlockID:       etctx.BlockID,
			TransactionID: etctx.TransactionID,
		},
	}
}

// Transaction returns this ERC20CoinCreationTransaction
// as regular tfchain transaction, using TransactionVersionERC20CoinCreation as the type.
func (etctx *ERC20CoinCreationTransaction) Transaction() types.Transaction {
	return types.Transaction{
		Version: TransactionVersionERC20CoinCreation,
		CoinOutputs: []types.CoinOutput{
			{
				Condition: types.NewCondition(types.NewUnlockHashCondition(etctx.Address)),
				Value:     etctx.Value,
			},
		},
		MinerFees: []types.Currency{etctx.TransactionFee},
		Extension: &ERC20CoinCreationTransactionExtension{
			BlockID:       etctx.BlockID,
			TransactionID: etctx.TransactionID,
		},
	}
}

// MarshalSia implements SiaMarshaler.MarshalSia,
// alias of MarshalRivine for backwards-compatibility reasons.
func (etctx ERC20CoinCreationTransaction) MarshalSia(w io.Writer) error {
	return etctx.MarshalRivine(w)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// alias of UnmarshalRivine for backwards-compatibility reasons.
func (etctx *ERC20CoinCreationTransaction) UnmarshalSia(r io.Reader) error {
	return etctx.UnmarshalRivine(r)
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (etctx ERC20CoinCreationTransaction) MarshalRivine(w io.Writer) error {
	return rivbin.NewEncoder(w).EncodeAll(
		etctx.Address,
		etctx.Value,
		etctx.TransactionFee,
		etctx.BlockID,
		etctx.TransactionID,
	)
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (etctx *ERC20CoinCreationTransaction) UnmarshalRivine(r io.Reader) error {
	return rivbin.NewDecoder(r).DecodeAll(
		&etctx.Address,
		&etctx.Value,
		&etctx.TransactionFee,
		&etctx.BlockID,
		&etctx.TransactionID,
	)
}

type (
	// ERC20CoinCreationTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 0xD1. It allows the conversion of ERC20-funds to TFT.
	ERC20CoinCreationTransactionController struct {
		Registry    ERC20Registry
		OneCoin     types.Currency
		TxValidator ERC20TransactionValidator
	}
)

// ensure at compile time that CoinCreationTransactionController
// implements the desired interfaces
var (
	_ types.TransactionController      = ERC20CoinCreationTransactionController{}
	_ types.TransactionValidator       = ERC20CoinCreationTransactionController{}
	_ types.CoinOutputValidator        = ERC20CoinCreationTransactionController{}
	_ types.BlockStakeOutputValidator  = ERC20CoinCreationTransactionController{}
	_ types.TransactionSignatureHasher = ERC20CoinCreationTransactionController{}
	_ types.TransactionIDEncoder       = ERC20CoinCreationTransactionController{}
)

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (etctc ERC20CoinCreationTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	etctx, err := ERC20CoinCreationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a ERC20CoinCreationTx: %v", err)
	}
	return siabin.NewEncoder(w).Encode(etctx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (etctc ERC20CoinCreationTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var etctx ERC20CoinCreationTransaction
	err := siabin.NewDecoder(r).Decode(&etctx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to binary-decode tx as a ERC20CoinCreationTx: %v", err)
	}
	// return ERC20 CoinCreation tx as regular tfchain tx data
	return etctx.TransactionData(), nil
}

// JSONEncodeTransactionData implements TransactionController.JSONEncodeTransactionData
func (etctc ERC20CoinCreationTransactionController) JSONEncodeTransactionData(txData types.TransactionData) ([]byte, error) {
	etctx, err := ERC20CoinCreationTransactionFromTransactionData(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert txData to a ERC20 CoinCreation Tx: %v", err)
	}
	return json.Marshal(etctx)
}

// JSONDecodeTransactionData implements TransactionController.JSONDecodeTransactionData
func (etctc ERC20CoinCreationTransactionController) JSONDecodeTransactionData(data []byte) (types.TransactionData, error) {
	var etctx ERC20CoinCreationTransaction
	err := json.Unmarshal(data, &etctx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to json-decode tx as a ERC20 CoinCreation Tx: %v", err)
	}
	// return ERC20 CoinCreation tx as regular tfchain tx data
	return etctx.TransactionData(), nil
}

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (etctc ERC20CoinCreationTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) error {
	// the content of the Tx ensures it will always fit in a block, due to how little data can be put in it

	// get CoinCreationTxn
	etctx, err := ERC20CoinCreationTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to use Tx as a ERC20 CoinCreation Tx: %v", err)
	}
	// check if the miner fee has the required minimum miner fee
	if etctx.TransactionFee.Cmp(constants.MinimumMinerFee) == -1 {
		return types.ErrTooSmallMinerFee
	}
	// check if the unlock hash is not the nil hash, and that the sent (created) value is not zero
	if etctx.Address.Type == types.UnlockTypeNil {
		return errors.New("ERC20 CoinCreation Tx is not allowed to send to the Nil UnlockHash")
	}
	if etctx.Value.IsZero() {
		return types.ErrZeroOutput
	}

	// validate if the ERC20 Transaction ID isn't already used
	txid, found, err := etctc.Registry.GetTFTTransactionIDForERC20TransactionID(etctx.TransactionID)
	if err != nil {
		return fmt.Errorf("internal error occured while checking if the ERC20 TransactionID %v was already registered: %v", etctx.TransactionID, err)
	}
	if found {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: ERC20 Tx ID %v already mapped to TFT Tx ID %v", etctx.TransactionID, txid)
	}

	// validate if the TFT Target Address is actually registered
	// as an ERC20 Withdrawal address
	_, found, err = etctc.Registry.GetERC20AddressForTFTAddress(etctx.Address)
	if err != nil {
		return fmt.Errorf("internal error occured while checking if the TFT address %v is registered as ERC20 withdrawal address: %v", etctx.Address, err)
	}
	if !found {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: Address %v is not registered as an ERC20 withdrawal address", etctx.Address)
	}

	// validate the ERC20 Tx using the used Validator
	erc20Address := ERC20AddressFromUnlockHash(etctx.Address)
	err = etctc.TxValidator.ValidateWithdrawTx(etctx.BlockID, etctx.TransactionID, erc20Address, etctx.Value)
	if err != nil {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: invalid attached ERC20 Tx: %v", err)
	}

	return nil
}

// ValidateCoinOutputs implements CoinOutputValidator.ValidateCoinOutputs
func (etctc ERC20CoinCreationTransactionController) ValidateCoinOutputs(t types.Transaction, ctx types.FundValidationContext, coinInputs map[types.CoinOutputID]types.CoinOutput) (err error) {
	return nil // always valid, coin outputs (and miner fees) are created not backed within an ERC20 CoinCreation transaction
}

// ValidateBlockStakeOutputs implements BlockStakeOutputValidator.ValidateBlockStakeOutputs
func (etctc ERC20CoinCreationTransactionController) ValidateBlockStakeOutputs(t types.Transaction, ctx types.FundValidationContext, blockStakeInputs map[types.BlockStakeOutputID]types.BlockStakeOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within an ERC20 CoinCreation transaction
}

// SignatureHash implements TransactionSignatureHasher.SignatureHash
func (etctc ERC20CoinCreationTransactionController) SignatureHash(t types.Transaction, extraObjects ...interface{}) (crypto.Hash, error) {
	etctx, err := ERC20CoinCreationTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as a ERC20 CoinCreation tx: %v", err)
	}

	h := crypto.NewHash()
	enc := siabin.NewEncoder(h)

	enc.EncodeAll(
		t.Version,
		SpecifierERC20CoinCreationTransaction,
	)

	if len(extraObjects) > 0 {
		enc.EncodeAll(extraObjects...)
	}

	enc.EncodeAll(
		etctx.Address,
		etctx.Value,
		etctx.TransactionFee,
		etctx.BlockID,
		etctx.TransactionID, // this ID has to ensure the TxSig and Hash is unique per transaction
	)

	var hash crypto.Hash
	h.Sum(hash[:0])
	return hash, nil
}

// EncodeTransactionIDInput implements TransactionIDEncoder.EncodeTransactionIDInput
func (etctc ERC20CoinCreationTransactionController) EncodeTransactionIDInput(w io.Writer, txData types.TransactionData) error {
	cctx, err := ERC20CoinCreationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a ERC20 CoinCreation Tx: %v", err)
	}
	return siabin.NewEncoder(w).EncodeAll(SpecifierERC20CoinCreationTransaction, cctx)
}

const (
	// HardcodedERC20AddressRegistrationFeeOneCoinMultiplier defines the hardcoded multiplier
	// (to be multiplied with the OneCoin Currency Value of the network), that defines the constant (hardcoded)
	// Registration Fee to be paid for the Registration of an ERC20 Withdrawal address.
	HardcodedERC20AddressRegistrationFeeOneCoinMultiplier = 10
)

type (
	// ERC20AddressRegistrationTransaction defines the Transaction (with version 0xD2)
	// used to register an ERC20 address linked to a regular TFT address (derived from the given public key).
	// This is required as to be able to convert ERC20 Funds back into TFT.
	ERC20AddressRegistrationTransaction struct {
		// The public key of which a TFT address can be derived, and thus also an ERC20 Address
		PublicKey types.PublicKey

		// Signature that proofs the ownership of the attached Public Key.
		Signature types.ByteSlice

		// RegistrationFee defines the Registration fee to be paid for the
		// registration on top of the regular Transaction fee.
		// TODO: integrate it into the parent block, for now it is ignored.
		RegistrationFee types.Currency

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency

		// CoinInputs are only used for the required fees,
		// which contains the regular Tx fee as well as the additional fees,
		// to be paid for the address registration. At least one CoinInput is required.
		CoinInputs []types.CoinInput
		// RefundCoinOutput is an optional coin output that can be used
		// to refund coins paid as inputs for the required fees.
		RefundCoinOutput *types.CoinOutput
	}

	// ERC20AddressRegistrationTransactionJSON defines the JSON structure of an ERC20AddressRegistrationTransaction,
	// which is an extended data structure when compared to the binary structure of an ERC20AddressRegistrationTransaction
	ERC20AddressRegistrationTransactionJSON struct {
		// The public key of which a TFT address can be derived, and thus also an ERC20 Address
		PublicKey types.PublicKey `json:"pubkey"`

		// TFTAddresses can be derived from the PublicKey,
		// if defined however it will be validated that the public key matches the given PublicKey.
		// Can be omitted as well, given that the raw tx does not contain this duplicate data.
		TFTAddress types.UnlockHash `json:"tftaddress,omitempty"`
		// ERC20Address can be derived from the PublicKey,
		// if defined however it will be validated that the public key matches the given PublicKey.
		// Can be omitted as well, given that the raw tx does not contain this duplicate data.
		ERC20Address ERC20Address `json:"erc20address,omitempty"`

		// Signature that proofs the ownership of the attached Public Key.
		Signature types.ByteSlice `json:"signature"`

		// RegistrationFee defines the Registration fee to be paid for the
		// registration on top of the regular Transaction fee.
		RegistrationFee types.Currency `json:"regfee"`

		// TransactionFee defines the regular Tx fee.
		TransactionFee types.Currency `json:"txfee"`

		// CoinInputs are only used for the required fees,
		// which contains the regular Tx fee as well as the additional fees,
		// to be paid for the address registration. At least one CoinInput is required.
		CoinInputs []types.CoinInput `json:"coininputs"`
		// RefundCoinOutput is an optional coin output that can be used
		// to refund coins paid as inputs for the required fees.
		RefundCoinOutput *types.CoinOutput `json:"refundcoinoutput,omitempty"`
	}

	// ERC20AddressRegistrationTransactionExtension defines the ERC20AddressRegistrationTransaction Extension Data
	ERC20AddressRegistrationTransactionExtension struct {
		RegistrationFee types.Currency
		PublicKey       types.PublicKey
		Signature       types.ByteSlice
	}
)

// ERC20AddressRegistrationTransactionFromTransaction creates an ERC20AddressRegistrationTransaction,
// using a regular in-memory tfchain transaction.
//
// Past the (tx) Version validation it piggy-backs onto the
// `ERC20AddressRegistrationTransactionFromTransactionData` constructor.
func ERC20AddressRegistrationTransactionFromTransaction(tx types.Transaction) (ERC20AddressRegistrationTransaction, error) {
	if tx.Version != TransactionVersionERC20AddressRegistration {
		return ERC20AddressRegistrationTransaction{}, fmt.Errorf(
			"an ERC20 address registration requires tx version %d",
			TransactionVersionERC20AddressRegistration)
	}
	return ERC20AddressRegistrationTransactionFromTransactionData(types.TransactionData{
		CoinInputs:        tx.CoinInputs,
		CoinOutputs:       tx.CoinOutputs,
		BlockStakeInputs:  tx.BlockStakeInputs,
		BlockStakeOutputs: tx.BlockStakeOutputs,
		MinerFees:         tx.MinerFees,
		ArbitraryData:     tx.ArbitraryData,
		Extension:         tx.Extension,
	})
}

// ERC20AddressRegistrationTransactionFromTransactionData creates an ERC20ConvertTransaction,
// using the TransactionData from a regular in-memory tfchain transaction.
func ERC20AddressRegistrationTransactionFromTransactionData(txData types.TransactionData) (ERC20AddressRegistrationTransaction, error) {
	// validate the Transaction Data

	// at least one coin input as well as one miner fee is required
	if len(txData.CoinInputs) == 0 || len(txData.MinerFees) != 1 {
		return ERC20AddressRegistrationTransaction{}, errors.New("at least one coin input and exactly one miner fee is required for an ERC20 Address Registration Transaction")
	}
	// no block stake inputs or block stake outputs are allowed
	if len(txData.BlockStakeInputs) != 0 || len(txData.BlockStakeOutputs) != 0 {
		return ERC20AddressRegistrationTransaction{}, errors.New("no block stake inputs/outputs are allowed in an ERC20 Address Registration Transaction")
	}
	// no arbitrary data is allowed
	if len(txData.ArbitraryData.Data) > 0 {
		return ERC20AddressRegistrationTransaction{}, errors.New("no arbitrary data is allowed in an ERC20 Address Registration Transaction")
	}
	// validate that the coin outputs is within the expected range
	if len(txData.CoinOutputs) > 1 {
		return ERC20AddressRegistrationTransaction{}, errors.New("an ERC20 Address Registration Transaction can only have one coin output")
	}

	// (tx) extension (data) is expected to be a pointer to a valid ERC20AddressRegistrationTransaction,
	// which contains all the properties unique to a 3bot (name transfer) Tx
	extensionData, ok := txData.Extension.(*ERC20AddressRegistrationTransactionExtension)
	if !ok {
		return ERC20AddressRegistrationTransaction{}, errors.New("invalid extension data for an ERC20 Address Registration Transaction")
	}

	// create the ERC20AddressRegistrationTransaction and return it,
	// further validation will/has-to be done using the Transaction Type, if required
	tx := ERC20AddressRegistrationTransaction{
		PublicKey:       extensionData.PublicKey,
		Signature:       extensionData.Signature,
		RegistrationFee: extensionData.RegistrationFee,
		TransactionFee:  txData.MinerFees[0],
		CoinInputs:      txData.CoinInputs,
	}
	if len(txData.CoinOutputs) == 1 {
		// take refund coin output if it exists
		tx.RefundCoinOutput = &txData.CoinOutputs[0]
	}
	return tx, nil
}

// TransactionData returns this ERC20AddressRegistrationTransaction
// as regular tfchain transaction data.
func (eartx *ERC20AddressRegistrationTransaction) TransactionData() types.TransactionData {
	txData := types.TransactionData{
		CoinInputs: eartx.CoinInputs,
		MinerFees:  []types.Currency{eartx.TransactionFee},
		Extension: &ERC20AddressRegistrationTransactionExtension{
			PublicKey:       eartx.PublicKey,
			Signature:       eartx.Signature,
			RegistrationFee: eartx.RegistrationFee,
		},
	}
	if eartx.RefundCoinOutput != nil {
		txData.CoinOutputs = append(txData.CoinOutputs, *eartx.RefundCoinOutput)
	}
	return txData
}

// Transaction returns this ERC20AddressRegistrationTransaction
// as regular tfchain transaction, using TransactionVersionERC20AddressRegistration as the type.
func (eartx *ERC20AddressRegistrationTransaction) Transaction() types.Transaction {
	tx := types.Transaction{
		Version:    TransactionVersionERC20AddressRegistration,
		CoinInputs: eartx.CoinInputs,
		MinerFees:  []types.Currency{eartx.TransactionFee},
		Extension: &ERC20AddressRegistrationTransactionExtension{
			PublicKey:       eartx.PublicKey,
			Signature:       eartx.Signature,
			RegistrationFee: eartx.RegistrationFee,
		},
	}
	if eartx.RefundCoinOutput != nil {
		tx.CoinOutputs = append(tx.CoinOutputs, *eartx.RefundCoinOutput)
	}
	return tx
}

// MarshalSia implements SiaMarshaler.MarshalSia,
// alias of MarshalRivine for backwards-compatibility reasons.
func (eartx ERC20AddressRegistrationTransaction) MarshalSia(w io.Writer) error {
	return eartx.MarshalRivine(w)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// alias of UnmarshalRivine for backwards-compatibility reasons.
func (eartx *ERC20AddressRegistrationTransaction) UnmarshalSia(r io.Reader) error {
	return eartx.UnmarshalRivine(r)
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (eartx ERC20AddressRegistrationTransaction) MarshalRivine(w io.Writer) error {
	return rivbin.NewEncoder(w).EncodeAll(
		eartx.PublicKey,
		eartx.Signature,
		eartx.RegistrationFee,
		eartx.TransactionFee,
		eartx.CoinInputs,
		eartx.RefundCoinOutput,
	)
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (eartx *ERC20AddressRegistrationTransaction) UnmarshalRivine(r io.Reader) error {
	return rivbin.NewDecoder(r).DecodeAll(
		&eartx.PublicKey,
		&eartx.Signature,
		&eartx.RegistrationFee,
		&eartx.TransactionFee,
		&eartx.CoinInputs,
		&eartx.RefundCoinOutput,
	)
}

// ERC20AddressFromUnlockHash creates an ERC20Address using as input
// for a new blake2b hash an UnlockHash (TFT Address), and taking the last 20 bytes of that.
func ERC20AddressFromUnlockHash(uh types.UnlockHash) (addr ERC20Address) {
	hash := crypto.HashObject(uh)
	offset := crypto.HashSize - ERC20AddressLength
	copy(addr[:], hash[offset:])
	return
}

// MarshalJSON implements json.Marshaler.MarshalRivine
func (eartx ERC20AddressRegistrationTransaction) MarshalJSON() ([]byte, error) {
	uh := types.NewPubKeyUnlockHash(eartx.PublicKey)
	addr := ERC20AddressFromUnlockHash(uh)
	return json.Marshal(ERC20AddressRegistrationTransactionJSON{
		PublicKey:        eartx.PublicKey,
		TFTAddress:       uh,
		ERC20Address:     addr,
		Signature:        eartx.Signature,
		RegistrationFee:  eartx.RegistrationFee,
		TransactionFee:   eartx.TransactionFee,
		CoinInputs:       eartx.CoinInputs,
		RefundCoinOutput: eartx.RefundCoinOutput,
	})
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON
func (eartx *ERC20AddressRegistrationTransaction) UnmarshalJSON(data []byte) error {
	var tx ERC20AddressRegistrationTransactionJSON
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return err
	}
	// validate the TFT/ERC20 address if given
	tftAddressDefined := tx.TFTAddress.Cmp(types.NilUnlockHash) != 0
	erc20AddressDefined := tx.ERC20Address != (ERC20Address{})
	if tftAddressDefined || erc20AddressDefined {
		uh := types.NewPubKeyUnlockHash(tx.PublicKey)
		if tftAddressDefined && tx.TFTAddress.Cmp(uh) != 0 {
			return errors.New("non-matching public key and TFT Address defined")
		}
		if erc20AddressDefined {
			addr := ERC20AddressFromUnlockHash(uh)
			if tx.ERC20Address != addr {
				return errors.New("non-matching public key and ERC20 Address defined")
			}
		}
	}
	// copy the in-memory (binary) format properties over and call it done
	eartx.PublicKey = tx.PublicKey
	eartx.Signature = tx.Signature
	eartx.RegistrationFee = tx.RegistrationFee
	eartx.TransactionFee = tx.TransactionFee
	eartx.CoinInputs = tx.CoinInputs
	eartx.RefundCoinOutput = tx.RefundCoinOutput
	return nil
}

type (
	// ERC20AddressRegistrationTransactionController defines a tfchain-specific transaction controller,
	// for a transaction type reserved at type 0xD2. It allows the registration of an ERC20 Address.
	ERC20AddressRegistrationTransactionController struct {
		Registry             ERC20Registry
		OneCoin              types.Currency
		BridgeFeePoolAddress types.UnlockHash
	}
)

var (
	// ensure at compile time that ERC20AddressRegistrationTransactionController
	// implements the desired interfaces
	_ types.TransactionController                = ERC20AddressRegistrationTransactionController{}
	_ types.TransactionValidator                 = ERC20AddressRegistrationTransactionController{}
	_ types.BlockStakeOutputValidator            = ERC20AddressRegistrationTransactionController{}
	_ types.TransactionSignatureHasher           = ERC20AddressRegistrationTransactionController{}
	_ types.TransactionExtensionSigner           = ERC20AddressRegistrationTransactionController{}
	_ types.TransactionIDEncoder                 = ERC20AddressRegistrationTransactionController{}
	_ types.TransactionCustomMinerPayoutGetter   = ERC20AddressRegistrationTransactionController{}
	_ types.TransactionCommonExtensionDataGetter = ERC20AddressRegistrationTransactionController{}
)

// EncodeTransactionData implements TransactionController.EncodeTransactionData
func (eartc ERC20AddressRegistrationTransactionController) EncodeTransactionData(w io.Writer, txData types.TransactionData) error {
	eartx, err := ERC20AddressRegistrationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to an ERC20AddressRegistrationTx: %v", err)
	}
	return rivbin.NewEncoder(w).Encode(eartx)
}

// DecodeTransactionData implements TransactionController.DecodeTransactionData
func (eartc ERC20AddressRegistrationTransactionController) DecodeTransactionData(r io.Reader) (types.TransactionData, error) {
	var eartx ERC20AddressRegistrationTransaction
	err := rivbin.NewDecoder(r).Decode(&eartx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to binary-decode tx as an ERC20AddressRegistrationTx: %v", err)
	}
	// return ERC20 Address Registration tx as regular tfchain tx data
	return eartx.TransactionData(), nil
}

// JSONEncodeTransactionData implements TransactionController.JSONEncodeTransactionData
func (eartc ERC20AddressRegistrationTransactionController) JSONEncodeTransactionData(txData types.TransactionData) ([]byte, error) {
	eartx, err := ERC20AddressRegistrationTransactionFromTransactionData(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert txData to an ERC20AddressRegistrationTx: %v", err)
	}
	return json.Marshal(eartx)
}

// JSONDecodeTransactionData implements TransactionController.JSONDecodeTransactionData
func (eartc ERC20AddressRegistrationTransactionController) JSONDecodeTransactionData(data []byte) (types.TransactionData, error) {
	var eartx ERC20AddressRegistrationTransaction
	err := json.Unmarshal(data, &eartx)
	if err != nil {
		return types.TransactionData{}, fmt.Errorf(
			"failed to json-decode tx as an ERC20AddressRegistrationTx: %v", err)
	}
	// return bot record update tx as regular tfchain tx data
	return eartx.TransactionData(), nil
}

// Specifiers used to ensure the bot-signatures are unique within each Tx.
var (
	ERC20AdddressRegistrationSignatureSpecifier = [...]byte{'r', 'e', 'g', 'i', 's', 't', 'r', 'a', 't', 'i', 'o', 'n'}
)

// ValidateTransaction implements TransactionValidator.ValidateTransaction
func (eartc ERC20AddressRegistrationTransactionController) ValidateTransaction(t types.Transaction, ctx types.ValidationContext, constants types.TransactionValidationConstants) error {
	// check tx fits within a block
	err := types.TransactionFitsInABlock(t, constants.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get ERC20AddressRegistration Tx
	eartx, err := ERC20AddressRegistrationTransactionFromTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to use tx as an ERC20 AddressRegistration tx: %v", err)
	}

	// validate the signature
	// > create condition
	uh := types.NewPubKeyUnlockHash(eartx.PublicKey)
	condition := types.NewCondition(types.NewUnlockHashCondition(uh))
	// > and a matching single-signature fulfillment
	fulfillment := types.NewFulfillment(&types.SingleSignatureFulfillment{
		PublicKey: eartx.PublicKey,
		Signature: eartx.Signature,
	})
	// > validate the signature is correct
	err = condition.Fulfill(fulfillment, types.FulfillContext{
		ExtraObjects: []interface{}{ERC20AdddressRegistrationSignatureSpecifier},
		BlockHeight:  ctx.BlockHeight,
		BlockTime:    ctx.BlockTime,
		Transaction:  t,
	})
	if err != nil {
		return fmt.Errorf("unauthorized ERC20 AddressRegistration tx: %v", err)
	}

	// validate the public key is not registered yet
	_, found, err := eartc.Registry.GetERC20AddressForTFTAddress(uh)
	if err == nil && found {
		return errors.New("invalid ERC20 AddressRegistration tx: public key has already registered an ERC20 address")
	}
	if err != nil {
		return fmt.Errorf("error while validating ERC20 AddressRegistration tx: error originating from TransactiondB: %v", err)
	}

	// validate the registration fee
	if eartx.RegistrationFee.Cmp(eartc.OneCoin.Mul64(HardcodedERC20AddressRegistrationFeeOneCoinMultiplier)) != 0 {
		return errors.New("invalid ERC20 Address Registration fee")
	}

	// validate the miner fee
	if eartx.TransactionFee.Cmp(constants.MinimumMinerFee) < 0 {
		return types.ErrTooSmallMinerFee
	}

	// prevent double spending
	spendCoins := make(map[types.CoinOutputID]struct{})
	for _, ci := range eartx.CoinInputs {
		if _, found := spendCoins[ci.ParentID]; found {
			return types.ErrDoubleSpend
		}
		spendCoins[ci.ParentID] = struct{}{}
	}

	// check if optional coin output is using standard condition
	if eartx.RefundCoinOutput != nil {
		err = eartx.RefundCoinOutput.Condition.IsStandardCondition(ctx)
		if err != nil {
			return err
		}
		// ensure the value is not 0
		if eartx.RefundCoinOutput.Value.IsZero() {
			return types.ErrZeroOutput
		}
	}
	// check if all fulfillments are standard
	for _, sci := range eartx.CoinInputs {
		err = sci.Fulfillment.IsStandardFulfillment(ctx)
		if err != nil {
			return err
		}
	}

	// Tx is valid
	return nil
}

// ValidateCoinOutputs is not implemented here for ERC20AddressRegistrationTransactionController,
// instead we can rely on the default ValidateCoinOutputs logic provided by Rivine.

// ValidateBlockStakeOutputs implements BlockStakeOutputValidator.ValidateBlockStakeOutputs
func (eartc ERC20AddressRegistrationTransactionController) ValidateBlockStakeOutputs(t types.Transaction, ctx types.FundValidationContext, blockStakeInputs map[types.BlockStakeOutputID]types.BlockStakeOutput) (err error) {
	return nil // always valid, no block stake inputs/outputs exist within an ERC20 AddressRegistration tx transaction
}

// SignatureHash implements TransactionSignatureHasher.SignatureHash
func (eartc ERC20AddressRegistrationTransactionController) SignatureHash(t types.Transaction, extraObjects ...interface{}) (crypto.Hash, error) {
	eartx, err := ERC20AddressRegistrationTransactionFromTransaction(t)
	if err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to use tx as an ERC20AddressRegistrationTx: %v", err)
	}

	h := crypto.NewHash()
	enc := rivbin.NewEncoder(h)

	enc.EncodeAll(
		t.Version,
		SpecifierERC20AddressRegistrationTransaction,
		eartx.PublicKey,
	)

	if len(extraObjects) > 0 {
		enc.EncodeAll(extraObjects...)
	}

	enc.Encode(len(eartx.CoinInputs))
	for _, ci := range eartx.CoinInputs {
		enc.Encode(ci.ParentID)
	}

	enc.EncodeAll(
		eartx.RegistrationFee,
		eartx.TransactionFee,
		eartx.RefundCoinOutput,
	)

	var hash crypto.Hash
	h.Sum(hash[:0])
	return hash, nil
}

// SignExtension implements TransactionExtensionSigner.SignExtension
func (eartc ERC20AddressRegistrationTransactionController) SignExtension(extension interface{}, sign func(*types.UnlockFulfillmentProxy, types.UnlockConditionProxy, ...interface{}) error) (interface{}, error) {
	// (tx) extension (data) is expected to be a pointer to a valid ERC20AddressRegistrationTransactionExtension
	eartxExtension, ok := extension.(*ERC20AddressRegistrationTransactionExtension)
	if !ok {
		return nil, errors.New("invalid extension data for a ERC20 AddressRegistration Transaction")
	}

	// > create condition
	condition := types.NewCondition(types.NewUnlockHashCondition(types.NewPubKeyUnlockHash(eartxExtension.PublicKey)))
	// > and a matching single-signature fulfillment
	fulfillment := types.NewFulfillment(&types.SingleSignatureFulfillment{
		PublicKey: eartxExtension.PublicKey,
		Signature: eartxExtension.Signature,
	})
	// sign the tx
	err := sign(&fulfillment, condition, ERC20AdddressRegistrationSignatureSpecifier)
	if err != nil {
		return nil, fmt.Errorf("failed to sign ERC20 Address Registration Tx: %v", err)
	}

	// copy over the tx signature
	// TODO: be able to get this directly from the fulfillment somehnow
	var ffData map[string]interface{}
	b, _ := json.Marshal(fulfillment)
	err = json.Unmarshal(b, &ffData)
	if err != nil {
		return nil, fmt.Errorf("invalid signed ERC20 Address Registration Tx: %v", err)
	}
	rawsig := ffData["data"].(map[string]interface{})["signature"].(string)
	err = eartxExtension.Signature.LoadString(rawsig)
	if err != nil {
		return nil, fmt.Errorf("invalid signed ERC20 Address Registration Tx: %v", err)
	}

	return eartxExtension, nil
}

// GetCustomMinerPayouts implements TransactionCustomMinerPayoutGetter.GetCustomMinerPayouts
func (eartc ERC20AddressRegistrationTransactionController) GetCustomMinerPayouts(extension interface{}) ([]types.MinerPayout, error) {
	// (tx) extension (data) is expected to be a pointer to a valid ERC20AddressRegistrationTransactionExtension
	eartxExtension, ok := extension.(*ERC20AddressRegistrationTransactionExtension)
	if !ok {
		return nil, errors.New("invalid extension data for a ERC20 AddressRegistration Transaction")
	}
	return []types.MinerPayout{
		{
			Value:      eartxExtension.RegistrationFee,
			UnlockHash: eartc.BridgeFeePoolAddress,
		},
	}, nil
}

// EncodeTransactionIDInput implements TransactionIDEncoder.EncodeTransactionIDInput
func (eartc ERC20AddressRegistrationTransactionController) EncodeTransactionIDInput(w io.Writer, txData types.TransactionData) error {
	eartx, err := ERC20AddressRegistrationTransactionFromTransactionData(txData)
	if err != nil {
		return fmt.Errorf("failed to convert txData to a ERC20AddressRegistration: %v", err)
	}
	return rivbin.NewEncoder(w).EncodeAll(SpecifierERC20AddressRegistrationTransaction, eartx)
}

// GetCommonExtensionData implements TransactionCommonExtensionDataGetter.GetCommonExtensionData
func (eartc ERC20AddressRegistrationTransactionController) GetCommonExtensionData(extension interface{}) (types.CommonTransactionExtensionData, error) {
	// (tx) extension (data) is expected to be a pointer to a valid ERC20AddressRegistrationTransactionExtension
	eartxExtension, ok := extension.(*ERC20AddressRegistrationTransactionExtension)
	if !ok {
		return types.CommonTransactionExtensionData{}, errors.New("invalid extension data for a ERC20 AddressRegistration Transaction")
	}
	return types.CommonTransactionExtensionData{
		UnlockConditions: []types.UnlockConditionProxy{
			types.NewCondition(types.NewUnlockHashCondition(types.NewPubKeyUnlockHash(eartxExtension.PublicKey))),
		},
	}, nil
}
