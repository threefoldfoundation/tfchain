package erc20

import (
	"errors"
	"fmt"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/persist"
	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/types"

	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"

	bolt "github.com/rivine/bbolt"
)

const (
	pluginDBVersion = "1.0.0.0"
	pluginDBHeader  = "erc20BridgePlugin"
)

var (
	bucketERC20ToTFTAddresses = []byte("addresses_erc20_to_tft") // erc20 => TFT
	bucketTFTToERC20Addresses = []byte("addresses_tft_to_erc20") // TFT => erc20
	bucketERC20TransactionIDs = []byte("erc20_transactionids")   // stores all unique ERC20 transaction ids used for erc20=>TFT exchanges

	bucketSlice = [][]byte{
		bucketERC20ToTFTAddresses,
		bucketTFTToERC20Addresses,
		bucketERC20TransactionIDs,
	}
)

type (
	// Plugin is a struct defines the ERC20 plugin
	Plugin struct {
		storage            modules.PluginViewStorage
		unregisterCallback modules.PluginUnregisterCallback
		txValidator        erc20types.ERC20TransactionValidator
		oneCoin            types.Currency
		txVersions         erc20types.TransactionVersions
	}
)

// NewPlugin creates a new ERC20 Plugin.
func NewPlugin(feePoolAddress types.UnlockHash, oneCoin types.Currency, txValidator erc20types.ERC20TransactionValidator, txVersions erc20types.TransactionVersions) *Plugin {
	p := &Plugin{
		txValidator: txValidator,
		oneCoin:     oneCoin,
	}
	types.RegisterTransactionVersion(txVersions.ERC20Conversion, erc20types.ERC20ConvertTransactionController{
		TransactionVersion: txVersions.ERC20Conversion,
	})
	types.RegisterTransactionVersion(txVersions.ERC20AddressRegistration, erc20types.ERC20AddressRegistrationTransactionController{
		TransactionVersion:   txVersions.ERC20AddressRegistration,
		Registry:             p,
		BridgeFeePoolAddress: feePoolAddress,
		OneCoin:              oneCoin,
	})
	types.RegisterTransactionVersion(txVersions.ERC20CoinCreation, erc20types.ERC20CoinCreationTransactionController{
		TransactionVersion: txVersions.ERC20CoinCreation,
		Registry:           p,
		OneCoin:            oneCoin,
	})
	return p
}

// GetERC20AddressForTFTAddress returns the mapped ERC20 address for the given TFT Address,
// iff the TFT Address has registered an ERC20 address explicitly.
func (p *Plugin) GetERC20AddressForTFTAddress(uh types.UnlockHash) (addr erc20types.ERC20Address, found bool, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) (err error) {
		addr, found, err = getERC20AddressForTFTAddress(bucket, uh)
		return
	})
	return
}

// GetTFTAddressForERC20Address returns the mapped TFT address for the given ERC20 Address,
// iff the TFT Address has registered an ERC20 address explicitly.
func (p *Plugin) GetTFTAddressForERC20Address(addr erc20types.ERC20Address) (uh types.UnlockHash, found bool, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) (err error) {
		uh, found, err = getTFTAddressForERC20Address(bucket, addr)
		return
	})
	return
}

// GetTFTTransactionIDForERC20TransactionID returns the mapped TFT TransactionID for the given ERC20 TransactionID,
// iff the ERC20 TransactionID has been used to fund an ERC20 CoinCreation Tx and has been registered as such, a nil TransactionID is returned otherwise.
func (p *Plugin) GetTFTTransactionIDForERC20TransactionID(id erc20types.ERC20Hash) (txid types.TransactionID, found bool, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) (err error) {
		txid, found, err = getTfchainTransactionIDForERC20TransactionID(bucket, id)
		return
	})
	return
}

// InitPlugin initializes the Bucket for the first time
func (p *Plugin) InitPlugin(metadata *persist.Metadata, bucket *bolt.Bucket, storage modules.PluginViewStorage, unregisterCallback modules.PluginUnregisterCallback) (persist.Metadata, error) {
	p.storage = storage
	p.unregisterCallback = unregisterCallback
	if metadata == nil {
		for _, bucketName := range bucketSlice {
			b := bucket.Bucket([]byte(bucketName))
			if b == nil {
				var err error
				_, err = bucket.CreateBucket([]byte(bucketName))
				if err != nil {
					return persist.Metadata{}, fmt.Errorf("failed to create bucket %s: %v", string(bucketName), err)
				}
			}
		}

		metadata = &persist.Metadata{
			Version: pluginDBVersion,
			Header:  pluginDBHeader,
		}
	} else if metadata.Version != pluginDBVersion {
		return persist.Metadata{}, errors.New("There is only 1 version of this plugin, version mismatch")
	}
	return *metadata, nil
}

// ApplyBlock applies a block's ERC20 transactions to the ERC20 bucket.
func (p *Plugin) ApplyBlock(block modules.ConsensusBlock, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("ERC20 bucket does not exist")
	}
	var err error
	for idx, txn := range block.Transactions {
		cTxn := modules.ConsensusTransaction{
			Transaction:            txn,
			BlockHeight:            block.Height,
			BlockTime:              block.Timestamp,
			SequenceID:             uint16(idx),
			SpentCoinOutputs:       block.SpentCoinOutputs,
			SpentBlockStakeOutputs: block.SpentBlockStakeOutputs,
		}
		err = p.ApplyTransaction(cTxn, bucket)
		if err != nil {
			return err
		}
	}
	return nil
}

// ApplyTransaction applies a ERC20 transactions to the ERC20 bucket.
func (p *Plugin) ApplyTransaction(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("ERC20 bucket does not exist")
	}
	// check the version and handle the ones we care about
	var err error
	switch txn.Version {
	case p.txVersions.ERC20AddressRegistration:
		err = p.applyERC20AddressRegistrationTx(txn, bucket)
	case p.txVersions.ERC20CoinCreation:
		err = p.applyERC20CoinCreationTx(txn, bucket)
	}
	return err
}

func (p *Plugin) applyERC20AddressRegistrationTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	etartx, err := erc20types.ERC20AddressRegistrationTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20AddressRegistration)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Address Registration tx type: %v", err)
	}

	tftaddr, err := types.NewPubKeyUnlockHash(etartx.PublicKey)
	if err != nil {
		return err
	}
	erc20addr, err := erc20types.ERC20AddressFromUnlockHash(tftaddr)
	if err != nil {
		return err
	}

	return applyERC20AddressMapping(bucket, tftaddr, erc20addr)
}

func (p *Plugin) applyERC20CoinCreationTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	etcctx, err := erc20types.ERC20CoinCreationTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20CoinCreation)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Coin Creation Tx type: %v", err)
	}
	return applyERC20TransactionID(bucket, etcctx.TransactionID, txn.ID())
}

// RevertBlock reverts a block's ERC20 transaction from the ERC20 bucket
func (p *Plugin) RevertBlock(block modules.ConsensusBlock, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("mint conditions bucket does not exist")
	}
	// collect all one-per-block mint conditions
	var err error
	for idx, txn := range block.Transactions {
		cTxn := modules.ConsensusTransaction{
			Transaction:            txn,
			BlockHeight:            block.Height,
			BlockTime:              block.Timestamp,
			SequenceID:             uint16(idx),
			SpentCoinOutputs:       block.SpentCoinOutputs,
			SpentBlockStakeOutputs: block.SpentBlockStakeOutputs,
		}
		err = p.RevertTransaction(cTxn, bucket)
		if err != nil {
			return err
		}
	}
	return nil
}

// RevertTransaction reverts a ERC20 transactions to the ERC20 bucket.
func (p *Plugin) RevertTransaction(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("ERC20 bucket does not exist")
	}
	var err error
	switch txn.Version {
	case p.txVersions.ERC20AddressRegistration:
		err = p.revertERC20AddressRegistrationTx(txn, bucket)
	case p.txVersions.ERC20CoinCreation:
		err = p.revertERC20CoinCreationTx(txn, bucket)
	}
	return err
}

func (p *Plugin) revertERC20AddressRegistrationTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	etartx, err := erc20types.ERC20AddressRegistrationTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20AddressRegistration)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Address Registration tx type: %v", err)
	}

	tftaddr, err := types.NewPubKeyUnlockHash(etartx.PublicKey)
	if err != nil {
		return err
	}
	erc20addr, err := erc20types.ERC20AddressFromUnlockHash(tftaddr)
	if err != nil {
		return err
	}

	return revertERC20AddressMapping(bucket, tftaddr, erc20addr)
}

func (p *Plugin) revertERC20CoinCreationTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	etcctx, err := erc20types.ERC20CoinCreationTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20CoinCreation)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Coin Creation Tx type: %v", err)
	}
	return revertERC20TransactionID(bucket, etcctx.TransactionID)
}

// TransactionValidators returns all tx validators linked to this plugin
func (p *Plugin) TransactionValidators() []modules.PluginTransactionValidationFunction {
	return nil
}

// TransactionValidatorVersionFunctionMapping returns all tx validators linked to this plugin
func (p *Plugin) TransactionValidatorVersionFunctionMapping() map[types.TransactionVersion][]modules.PluginTransactionValidationFunction {
	return map[types.TransactionVersion][]modules.PluginTransactionValidationFunction{
		p.txVersions.ERC20Conversion: {
			p.validateERC20ConvertTx,
		},
		p.txVersions.ERC20CoinCreation: {
			p.validateERC20CoinCreationTx,
		},
		p.txVersions.ERC20AddressRegistration: {
			p.validateERC20CAddressRegistrationTx,
		},
	}
}

func (p *Plugin) validateERC20ConvertTx(txn modules.ConsensusTransaction, ctx types.TransactionValidationContext, bucket *persist.LazyBoltBucket) error {
	// check tx fits within a block
	err := types.TransactionFitsInABlock(txn.Transaction, ctx.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get ERC20ConvertTx
	etctx, err := erc20types.ERC20ConvertTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20Conversion)
	if err != nil {
		return fmt.Errorf("failed to use tx as an ERC20 convert tx: %v", err)
	}

	// ensure the value is a valid minimum
	if etctx.Value.Cmp(p.oneCoin.Mul64(erc20types.ERC20ConversionMinimumValue)) < 0 {
		return fmt.Errorf("ERC20 requires a minimum value of %d coins to be converted", erc20types.ERC20ConversionMinimumValue)
	}

	// validate the miner fee
	if etctx.TransactionFee.Cmp(ctx.MinimumMinerFee) < 0 {
		return types.ErrTooSmallMinerFee
	}

	// prevent double spending
	spendCoins := make(map[types.CoinOutputID]struct{})
	for _, ci := range etctx.CoinInputs {
		if _, found := spendCoins[ci.ParentID]; found {
			return types.ErrDoubleSpend
		}
		spendCoins[ci.ParentID] = struct{}{}
	}

	// TODO: what to do with this one

	// check if optional coin output is using standard condition
	if etctx.RefundCoinOutput != nil {
		err = etctx.RefundCoinOutput.Condition.IsStandardCondition(ctx.ValidationContext)
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
		err = sci.Fulfillment.IsStandardFulfillment(ctx.ValidationContext)
		if err != nil {
			return err
		}
	}

	// validate coin outputs
	// ... collect the coin input sum
	var coinInputSum types.Currency
	for _, ci := range txn.CoinInputs {
		co, ok := txn.SpentCoinOutputs[ci.ParentID]
		if !ok {
			return fmt.Errorf(
				"unable to find parent ID %s as an unspent coin output in the current consensus transaction at block height %d",
				ci.ParentID.String(), ctx.BlockHeight)
		}
		coinInputSum = coinInputSum.Add(co.Value)
	}
	// ... compute the expected total output (fee + total fee)
	totalOutput := etctx.TransactionFee.Add(etctx.Value)
	if etctx.RefundCoinOutput != nil {
		totalOutput = totalOutput.Add(etctx.RefundCoinOutput.Value)
	}
	// ... compare output/input
	if !coinInputSum.Equals(totalOutput) {
		return types.ErrCoinInputOutputMismatch
	}
	return nil
}

func (p *Plugin) validateERC20CoinCreationTx(txn modules.ConsensusTransaction, ctx types.TransactionValidationContext, bucket *persist.LazyBoltBucket) error {
	// get CoinCreationTxn
	etctx, err := erc20types.ERC20CoinCreationTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20CoinCreation)
	if err != nil {
		return fmt.Errorf("failed to use Tx as a ERC20 CoinCreation Tx: %v", err)
	}
	// check if the miner fee has the required minimum miner fee
	if etctx.TransactionFee.Cmp(ctx.MinimumMinerFee) == -1 {
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
	txid, found, err := p.GetTFTTransactionIDForERC20TransactionID(etctx.TransactionID)
	if err != nil {
		return fmt.Errorf("internal error occured while checking if the ERC20 TransactionID %v was already registered: %v", etctx.TransactionID, err)
	}
	if found {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: ERC20 Tx ID %v already mapped to TFT Tx ID %v", etctx.TransactionID, txid)
	}

	// validate if the TFT Target Address is actually registered
	// as an ERC20 Withdrawal address
	_, found, err = p.GetERC20AddressForTFTAddress(etctx.Address)
	if err != nil {
		return fmt.Errorf("internal error occured while checking if the TFT address %v is registered as ERC20 withdrawal address: %v", etctx.Address, err)
	}
	if !found {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: Address %v is not registered as an ERC20 withdrawal address", etctx.Address)
	}

	// validate the ERC20 Tx using the used Validator
	erc20Address, err := erc20types.ERC20AddressFromUnlockHash(etctx.Address)
	if err != nil {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: %v", err)
	}
	// we need to validate the total amount in the transaction, since the contract does not know which part went to txfee and which part was actually received
	err = p.txValidator.ValidateWithdrawTx(etctx.BlockID, etctx.TransactionID, erc20Address, etctx.Value.Add(etctx.TransactionFee))
	if err != nil {
		return fmt.Errorf("invalid ERC20 CoinCreation Tx: invalid attached ERC20 Tx: %v", err)
	}

	return nil
}

func (p *Plugin) validateERC20CAddressRegistrationTx(txn modules.ConsensusTransaction, ctx types.TransactionValidationContext, bucket *persist.LazyBoltBucket) error {
	// check tx fits within a block
	err := types.TransactionFitsInABlock(txn.Transaction, ctx.BlockSizeLimit)
	if err != nil {
		return err
	}

	// get ERC20AddressRegistration Tx
	eartx, err := erc20types.ERC20AddressRegistrationTransactionFromTransaction(txn.Transaction, p.txVersions.ERC20AddressRegistration)
	if err != nil {
		return fmt.Errorf("failed to use tx as an ERC20 AddressRegistration tx: %v", err)
	}

	// validate the signature
	// > create condition
	uh, err := types.NewPubKeyUnlockHash(eartx.PublicKey)
	if err != nil {
		return err
	}
	condition := types.NewCondition(types.NewUnlockHashCondition(uh))
	// > and a matching single-signature fulfillment
	fulfillment := types.NewFulfillment(&types.SingleSignatureFulfillment{
		PublicKey: eartx.PublicKey,
		Signature: eartx.Signature,
	})
	// > validate the signature is correct
	err = condition.Fulfill(fulfillment, types.FulfillContext{
		ExtraObjects: []interface{}{erc20types.ERC20AdddressRegistrationSignatureSpecifier},
		BlockHeight:  ctx.BlockHeight,
		BlockTime:    ctx.BlockTime,
		Transaction:  txn.Transaction,
	})
	if err != nil {
		return fmt.Errorf("unauthorized ERC20 AddressRegistration tx: %v", err)
	}

	// validate the public key is not registered yet
	_, found, err := p.GetERC20AddressForTFTAddress(uh)
	if err == nil && found {
		return errors.New("invalid ERC20 AddressRegistration tx: public key has already registered an ERC20 address")
	}
	if err != nil {
		return fmt.Errorf("error while validating ERC20 AddressRegistration tx: error originating from TransactiondB: %v", err)
	}

	// validate the registration fee
	if eartx.RegistrationFee.Cmp(p.oneCoin.Mul64(erc20types.HardcodedERC20AddressRegistrationFeeOneCoinMultiplier)) != 0 {
		return errors.New("invalid ERC20 Address Registration fee")
	}

	// validate the miner fee
	if eartx.TransactionFee.Cmp(ctx.MinimumMinerFee) < 0 {
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
		err = eartx.RefundCoinOutput.Condition.IsStandardCondition(ctx.ValidationContext)
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
		err = sci.Fulfillment.IsStandardFulfillment(ctx.ValidationContext)
		if err != nil {
			return err
		}
	}

	// Tx is valid
	return nil
}

// Close unregisters the plugin from the consensus
func (p *Plugin) Close() error {
	return p.storage.Close()
}

func applyERC20AddressMapping(rb *persist.LazyBoltBucket, tftaddr types.UnlockHash, erc20addr erc20types.ERC20Address) error {
	btft, err := rivbin.Marshal(tftaddr)
	if err != nil {
		return err
	}
	berc20, err := rivbin.Marshal(erc20addr)
	if err != nil {
		return err
	}

	// store ERC20->TFT mapping
	bucket, err := rb.Bucket(bucketERC20ToTFTAddresses)
	if err != nil {
		return fmt.Errorf("corrupt ERC20 Plugin: %v", err)
	}
	err = bucket.Put(berc20, btft)
	if err != nil {
		return fmt.Errorf("error while storing ERC20->TFT address mapping: %v", err)
	}

	// store TFT->ERC20 mapping
	bucket, err = rb.Bucket(bucketTFTToERC20Addresses)
	if err != nil {
		return fmt.Errorf("corrupt ERC20 Plugin: %v", err)
	}
	err = bucket.Put(btft, berc20)
	if err != nil {
		return fmt.Errorf("error while storing TFT->ERC20 address mapping: %v", err)
	}

	// done
	return nil
}
func revertERC20AddressMapping(rb *persist.LazyBoltBucket, tftaddr types.UnlockHash, erc20addr erc20types.ERC20Address) error {
	btft, err := rivbin.Marshal(tftaddr)
	if err != nil {
		return err
	}
	berc20, err := rivbin.Marshal(erc20addr)
	if err != nil {
		return err
	}

	// delete ERC20->TFT mapping
	bucket, err := rb.Bucket(bucketERC20ToTFTAddresses)
	if err != nil {
		return fmt.Errorf("corrupt ERC20 Plugin: %v", err)
	}
	err = bucket.Delete(berc20)
	if err != nil {
		return fmt.Errorf("error while deleting ERC20->TFT address mapping: %v", err)
	}

	// delete TFT->ERC20 mapping
	bucket, err = rb.Bucket(bucketTFTToERC20Addresses)
	if err != nil {
		return fmt.Errorf("corrupt ERC20 Plugin: %v", err)
	}
	err = bucket.Delete(btft)
	if err != nil {
		return fmt.Errorf("error while deleting TFT->ERC20 address mapping: %v", err)
	}

	// done
	return nil
}

func getERC20AddressForTFTAddress(rb *bolt.Bucket, uh types.UnlockHash) (erc20types.ERC20Address, bool, error) {
	bucket := rb.Bucket(bucketTFTToERC20Addresses)
	if bucket == nil {
		return erc20types.ERC20Address{}, false, errors.New("corrupt transaction DB: TFT->ERC20 bucket does not exist")
	}
	bUH, err := rivbin.Marshal(uh)
	if err != nil {
		return erc20types.ERC20Address{}, false, err
	}
	b := bucket.Get(bUH)
	if len(b) == 0 {
		return erc20types.ERC20Address{}, false, nil
	}
	var addr erc20types.ERC20Address
	err = rivbin.Unmarshal(b, &addr)
	if err != nil {
		return erc20types.ERC20Address{}, false, fmt.Errorf("failed to fetch ERC20 Address for TFT address %v: %v", uh, err)
	}
	return addr, true, nil
}

func getTFTAddressForERC20Address(rb *bolt.Bucket, addr erc20types.ERC20Address) (types.UnlockHash, bool, error) {
	bucket := rb.Bucket(bucketERC20ToTFTAddresses)
	if bucket == nil {
		return types.UnlockHash{}, false, errors.New("corrupt transaction DB: ERC20->TFT bucket does not exist")
	}
	bAddr, err := rivbin.Marshal(addr)
	if err != nil {
		return types.UnlockHash{}, false, err
	}
	b := bucket.Get(bAddr)
	if len(b) == 0 {
		return types.UnlockHash{}, false, nil
	}
	var uh types.UnlockHash
	err = rivbin.Unmarshal(b, &uh)
	if err != nil {
		return types.UnlockHash{}, false, fmt.Errorf("failed to fetch TFT Address for ERC20 address %v: %v", addr, err)
	}
	return uh, true, nil
}

func applyERC20TransactionID(rb *persist.LazyBoltBucket, erc20id erc20types.ERC20Hash, tftid types.TransactionID) error {
	bucket, err := rb.Bucket(bucketERC20TransactionIDs)
	if err != nil {
		return fmt.Errorf("corrupt ERC20 Plugin: %v", err)
	}
	bERC20id, err := rivbin.Marshal(erc20id)
	if err != nil {
		return err
	}
	bTFTid, err := rivbin.Marshal(tftid)
	if err != nil {
		return err
	}
	err = bucket.Put(bERC20id, bTFTid)
	if err != nil {
		return fmt.Errorf("error while storing ERC20 TransactionID %v: %v", erc20id, err)
	}
	return nil
}
func revertERC20TransactionID(rb *persist.LazyBoltBucket, id erc20types.ERC20Hash) error {
	bucket, err := rb.Bucket(bucketERC20TransactionIDs)
	if err != nil {
		return fmt.Errorf("corrupt ERC20 Plugin: %v", err)
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return err
	}
	err = bucket.Delete(bID)
	if err != nil {
		return fmt.Errorf("error while deleting ERC20 TransactionID %v: %v", id, err)
	}
	return nil
}
func getTfchainTransactionIDForERC20TransactionID(rb *bolt.Bucket, id erc20types.ERC20Hash) (types.TransactionID, bool, error) {
	bucket := rb.Bucket(bucketERC20TransactionIDs)
	if bucket == nil {
		return types.TransactionID{}, false, errors.New("corrupt transaction DB: ERC20 TransactionIDs bucket does not exist")
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return types.TransactionID{}, false, err
	}
	b := bucket.Get(bID)
	if len(b) == 0 {
		return types.TransactionID{}, false, nil
	}
	var txid types.TransactionID
	err = rivbin.Unmarshal(b, &txid)
	if err != nil {
		return types.TransactionID{}, false, fmt.Errorf("corrupt ERC20 Plugin: invalid tfchain TransactionID fetched for ERC20 TxID %v: %v", id, err)
	}
	return txid, true, nil
}
