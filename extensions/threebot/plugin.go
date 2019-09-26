package threebot

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/persist"
	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/types"

	tbtypes "github.com/threefoldfoundation/tfchain/extensions/threebot/types"

	bolt "github.com/rivine/bbolt"
)

const (
	pluginDBVersion = "1.0.0.0"
	pluginDBHeader  = "threebotPlugin"
)

var (
	bucketBotRecords               = []byte("botrecords")      // ID => name
	bucketBotKeyToIDMapping        = []byte("botkeys")         // Key => ID
	bucketBotNameToIDMapping       = []byte("botnames")        // Name => ID
	bucketBotRecordImplicitUpdates = []byte("botimplupdates")  // txID => implicitBotRecordUpdate
	bucketBotTransactions          = []byte("bottransactions") // ID => []txID

	bucketBlockTime = []byte("blockTimes") // block times

	bucketSlice = [][]byte{
		bucketBotRecords,
		bucketBotKeyToIDMapping,
		bucketBotNameToIDMapping,
		bucketBotRecordImplicitUpdates,
		bucketBotTransactions,
		bucketBlockTime,
	}
)

type (
	// Plugin is a struct defines the 3bot plugin
	Plugin struct {
		storage            modules.PluginViewStorage
		unregisterCallback modules.PluginUnregisterCallback
	}
)

// NewPlugin creates a new 3bot Plugin.
func NewPlugin(registryPool types.UnlockHash, oneCoin types.Currency) *Plugin {
	p := new(Plugin)
	types.RegisterTransactionVersion(tbtypes.TransactionVersionBotRegistration, tbtypes.BotRegistrationTransactionController{
		Registry:            p,
		RegistryPoolAddress: registryPool,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(tbtypes.TransactionVersionBotRecordUpdate, tbtypes.BotUpdateRecordTransactionController{
		Registry:            p,
		RegistryPoolAddress: registryPool,
		OneCoin:             oneCoin,
	})
	types.RegisterTransactionVersion(tbtypes.TransactionVersionBotNameTransfer, tbtypes.BotNameTransferTransactionController{
		Registry:            p,
		RegistryPoolAddress: registryPool,
		OneCoin:             oneCoin,
	})
	return p
}

// GetRecordForID returns the record mapped to the given BotID.
func (p *Plugin) GetRecordForID(id tbtypes.BotID) (record *tbtypes.BotRecord, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) (err error) {
		record, err = getRecordForID(bucket, id)
		return
	})
	return
}

// internal function to get a record from the TxDB
func getRecordForID(bucket *bolt.Bucket, id tbtypes.BotID) (*tbtypes.BotRecord, error) {
	recordBucket := bucket.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return nil, errors.New("corrupt 3bot Plugin DB: bot record bucket does not exist")
	}

	bid, err := rivbin.Marshal(id)
	if err != nil {
		return nil, err
	}
	b := recordBucket.Get(bid)
	if len(b) == 0 {
		return nil, tbtypes.ErrBotNotFound
	}

	record := new(tbtypes.BotRecord)
	err = rivbin.Unmarshal(b, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// GetRecordForKey returns the record mapped to the given Key.
func (p *Plugin) GetRecordForKey(key types.PublicKey) (record *tbtypes.BotRecord, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) error {
		id, err := getBotIDForPublicKey(bucket, key)
		if err != nil {
			return err
		}

		record, err = getRecordForID(bucket, id)
		return err
	})
	return
}

func getBotIDForPublicKey(bucket *bolt.Bucket, key types.PublicKey) (tbtypes.BotID, error) {
	keyBucket := bucket.Bucket(bucketBotKeyToIDMapping)
	if keyBucket == nil {
		return 0, errors.New("corrupt 3bot plugin DB: bot key bucket does not exist")
	}

	bkey, err := rivbin.Marshal(key)
	if err != nil {
		return 0, err
	}
	b := keyBucket.Get(bkey)
	if len(b) == 0 {
		return 0, tbtypes.ErrBotKeyNotFound
	}

	var id tbtypes.BotID
	err = rivbin.Unmarshal(b, &id)
	return id, err
}

// GetRecordForName returns the record mapped to the given Name.
func (p *Plugin) GetRecordForName(name tbtypes.BotName) (record *tbtypes.BotRecord, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) error {
		nameBucket := bucket.Bucket(bucketBotNameToIDMapping)
		if nameBucket == nil {
			return errors.New("corrupt 3bot plugin DB: bot name bucket does not exist")
		}

		bname, err := rivbin.Marshal(name)
		if err != nil {
			return err
		}
		b := nameBucket.Get(bname)
		if len(b) == 0 {
			return tbtypes.ErrBotNameNotFound
		}

		var id tbtypes.BotID
		err = rivbin.Unmarshal(b, &id)
		if err != nil {
			return err
		}

		record, err = getRecordForID(bucket, id)
		if err != nil {
			return err
		}

		blockTimeBucket := bucket.Bucket(bucketBlockTime)
		if blockTimeBucket == nil {
			return fmt.Errorf("corrupt 3bot plugin DB: bucket %s not foumd", string(bucketBlockTime))
		}
		_, chainTime, err := getCurrentBlockHeightAndTime(blockTimeBucket)
		if err != nil {
			return err
		}

		if record.Expiration.SiaTimestamp() <= chainTime {
			// a botname automatically expires as soon as the last 3bot that owned it expired as well
			return tbtypes.ErrBotNameExpired
		}
		return nil
	})
	return
}

// GetBotTransactionIdentifiers returns the identifiers of all transactions that created and updated the given bot's record.
//
// The transaction identifiers are returned in the (stable) order as defined by the blockchain.
func (p *Plugin) GetBotTransactionIdentifiers(id tbtypes.BotID) (ids []types.TransactionID, err error) {
	err = p.storage.View(func(bucket *bolt.Bucket) (err error) {
		ids, err = getBotTransactions(bucket, id)
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

// ApplyBlock applies a block's 3bot transactions to the 3Bot bucket.
func (p *Plugin) ApplyBlock(block modules.ConsensusBlock, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("3Bot bucket does not exist")
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
	blockTimeBucket, err := bucket.Bucket(bucketBlockTime)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	return setStatsBlockTime(blockTimeBucket, block.Height, block.Timestamp)
}

// ApplyTransaction applies a 3Bot transactions to the 3Bot bucket.
func (p *Plugin) ApplyTransaction(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("3Bot bucket does not exist")
	}
	// check the version and handle the ones we care about
	var err error
	switch txn.Version {
	case tbtypes.TransactionVersionBotRegistration:
		err = p.applyBotRegistrationTx(txn, bucket)
	case tbtypes.TransactionVersionBotRecordUpdate:
		err = p.applyRecordUpdateTx(txn, bucket)
	case tbtypes.TransactionVersionBotNameTransfer:
		err = p.applyBotNameTransferTx(txn, bucket)
	}
	return err
}

func (p *Plugin) applyBotRegistrationTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	recordBucket, err := bucket.Bucket(bucketBotRecords)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: bot record bucket error: %v", err)
	}
	brtx, err := tbtypes.BotRegistrationTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot registration tx type: %v", err)
	}
	// get the unique ID for the 3bot, using bolt's auto incrementing feature
	sequenceIndex, err := recordBucket.NextSequence()
	if err != nil {
		return fmt.Errorf("error while getting auto incrementing sequence bot ID: %v", err)
	}
	if sequenceIndex > tbtypes.MaxBotID {
		return errors.New("error while getting auto incrementing sequence bot ID: value exceeds 32 bit")
	}
	id := tbtypes.BotID(sequenceIndex)
	// create the record
	record := tbtypes.BotRecord{
		ID:         id,
		PublicKey:  brtx.Identification.PublicKey,
		Expiration: tbtypes.SiaTimestampAsCompactTimestamp(txn.BlockTime) + tbtypes.CompactTimestamp(brtx.NrOfMonths)*tbtypes.BotMonth,
	}
	err = record.AddNetworkAddresses(brtx.Addresses...)
	if err != nil {
		return fmt.Errorf("error while adding network addresses to bot (%v): %v", brtx.Identification.PublicKey, err)
	}
	err = record.AddNames(brtx.Names...)
	if err != nil {
		return fmt.Errorf("error while adding bot names to bot (%v): %v", brtx.Identification.PublicKey, err)
	}
	// store the record, and the other mappings, assuming the consensus validated that
	// the registration Tx is completely valid
	bid, err := rivbin.Marshal(id)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	brecord, err := rivbin.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while storing record for bot %d with public key %v: %v", id, brtx.Identification.PublicKey, err)
	}
	// store pubkey to ID mapping
	err = applyKeyToIDMapping(bucket, brtx.Identification.PublicKey, id)
	if err != nil {
		return fmt.Errorf("error while storing pubKey %s to bot id %d mapping: %v", brtx.Identification.PublicKey, id, err)
	}
	// store all name mappings
	for _, name := range brtx.Names {
		err = applyNameToIDMapping(bucket, name, id)
		if err != nil {
			return fmt.Errorf("error while storing name %s to bot id %d mapping: %v", name.String(), id, err)
		}
	}
	// apply the transactionID to the list of transactionIDs for the given bot
	err = applyBotTransaction(bucket, id, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID), txn.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for bot %d: %v", id, err)
	}
	// all information is applied
	return nil
}

func (p *Plugin) applyRecordUpdateTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	recordBucket, err := bucket.Bucket(bucketBotRecords)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: bot record bucket error: %v", err)
	}
	brutx, err := tbtypes.BotRecordUpdateTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot record update tx type: %v", err)
	}

	// get the bot record
	bid, err := rivbin.Marshal(brutx.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot ID: %v", err)
	}
	b := recordBucket.Get(bid)
	if len(b) == 0 {
		return errors.New("no bot record found for the specified identifier")
	}
	var record tbtypes.BotRecord
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found bot record: %v", err)
	}

	// check if bot is still active,
	// if the bot is still active we expect that the Tx defines the names to remove
	// otherwise we require that all names should be removed
	var namesInRecordRemovedImplicitly []tbtypes.BotName
	if record.IsExpired(txn.BlockTime) {
		namesInRecordRemovedImplicitly = record.Names.Difference(tbtypes.BotNameSortedSet{}) // A \ {} = A
		// store the implicit update that will happen due to the invalid period prior to this Tx,
		// this will help is in reverting the record back to its original state,
		// as such implicit updates cannot be easily reverted otherwise
		err = applyImplicitBotRecordUpdate(bucket, txn.ID(), implicitBotRecordUpdate{
			PreviousExpirationTime: record.Expiration,
			InactiveNamesRemoved:   namesInRecordRemovedImplicitly,
		})
		if err != nil {
			return fmt.Errorf("failed to apply implicit record update: %v", err)
		}
	}

	// update it (will also reset names of an inactive bot)
	err = brutx.UpdateBotRecord(txn.BlockTime, &record)
	if err != nil {
		return fmt.Errorf("failed to update bot record: %v", err)
	}

	// save it
	bid, err = rivbin.Marshal(brutx.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot ID: %v", err)
	}
	brecord, err := rivbin.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while updating record for bot %d: %v", brutx.Identifier, err)
	}

	if len(namesInRecordRemovedImplicitly) > 0 {
		// otherwise remove all names that previously active,
		// as we can assume that an update of a record update HAS to make it active again
		for _, name := range namesInRecordRemovedImplicitly {
			err = revertNameToIDMappingIfOwnedByBot(bucket, name, record.ID)
			if err != nil {
				return fmt.Errorf("failed to update bot record: error while tx-removing mapping of name %v: %v", name, err)
			}
		}
	} else {
		// if the bot was active, we apply the removals as defined by the Tx
		for _, name := range brutx.Names.Remove {
			err = revertNameToIDMapping(bucket, name)
			if err != nil {
				return fmt.Errorf("failed to update bot record: error while record-removing mapping of name %v: %v", name, err)
			}
		}
	}

	// add mapping for all the added names
	for _, name := range brutx.Names.Add {
		err = applyNameToIDMapping(bucket, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// apply the transactionID to the list of transactionIDs for the given bot
	err = applyBotTransaction(bucket, record.ID, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID), txn.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for bot %d: %v", record.ID, err)
	}

	// all information is applied
	return nil
}

func (p *Plugin) applyBotNameTransferTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	recordBucket, err := bucket.Bucket(bucketBotRecords)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: bot record bucket error: %v", err)
	}
	bnttx, err := tbtypes.BotNameTransferTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot name transfer tx type: %v", err)
	}

	// get the sender bot record
	bid, err := rivbin.Marshal(bnttx.Sender.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot ID: %v", err)
	}
	b := recordBucket.Get(bid)
	if len(b) == 0 {
		return errors.New("no bot record found for the specified sender bot identifier")
	}
	var record tbtypes.BotRecord
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found record of sender bot %d: %v", bnttx.Sender.Identifier, err)
	}
	// update sender bot (this also ensures the sender bot isn't expired)
	err = bnttx.UpdateSenderBotRecord(txn.BlockTime, &record)
	if err != nil { // automatically checks also if at least one name is transferred, returning an error if not
		return fmt.Errorf("failed to update record of sender bot %d: %v", record.ID, err)
	}
	// save the record of the sender bot
	bid, err = rivbin.Marshal(record.ID)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot ID: %v", err)
	}
	brecord, err := rivbin.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while saving the updated record for sender bot %d: %v", record.ID, err)
	}

	// apply the transactionID to the list of transactionIDs for the sender bot
	err = applyBotTransaction(bucket, record.ID, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID), txn.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for sender bot %d: %v", record.ID, err)
	}

	// get the receiver bot record
	bid, err = rivbin.Marshal(bnttx.Receiver.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot ID: %v", err)
	}
	b = recordBucket.Get(bid)
	if len(b) == 0 {
		return errors.New("no bot record found for the specified receiver bot identifier")
	}
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found record of receiver bot %d: %v", bnttx.Receiver.Identifier, err)
	}
	// update receiver bot (this also ensures the receiver bot isn't expired)
	err = bnttx.UpdateReceiverBotRecord(txn.BlockTime, &record)
	if err != nil { // automatically checks also if at least one name is transferred, returning an error if not
		return fmt.Errorf("failed to update record of receiver bot %d: %v", record.ID, err)
	}
	// save the record of the receiver bot
	bid, err = rivbin.Marshal(record.ID)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot ID: %v", err)
	}
	brecord, err = rivbin.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal 3bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while saving the updated record for receiver bot %d: %v", record.ID, err)
	}

	// update mapping for all the transferred names
	for _, name := range bnttx.Names {
		err = applyNameToIDMapping(bucket, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while mapping name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// apply the transactionID to the list of transactionIDs for the receiver bot
	err = applyBotTransaction(bucket, record.ID, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID), txn.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for receiver bot %d: %v", record.ID, err)
	}

	// update went fine
	return nil
}

// RevertBlock reverts a block's 3Bot transaction from the 3Bot bucket
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

	blockTimeBucket, err := bucket.Bucket(bucketBlockTime)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	return deleteStatsBlockTime(blockTimeBucket, block.Height)
}

// RevertTransaction reverts a 3Bot transactions to the 3Bot bucket.
func (p *Plugin) RevertTransaction(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	if bucket == nil {
		return errors.New("3Bot bucket does not exist")
	}
	var err error
	switch txn.Version {
	case tbtypes.TransactionVersionBotRegistration:
		err = p.revertBotRegistrationTx(txn, bucket)
	case tbtypes.TransactionVersionBotRecordUpdate:
		err = p.revertRecordUpdateTx(txn, bucket)
	case tbtypes.TransactionVersionBotNameTransfer:
		err = p.revertBotNameTransferTx(txn, bucket)
	}
	return err
}

func (p *Plugin) revertBotRegistrationTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	recordBucket, err := bucket.Bucket(bucketBotRecords)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	brtx, err := tbtypes.BotRegistrationTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot registration tx type: %v", err)
	}
	// the ID should be equal to the current bucket sequence, given it was incremented by the registration process
	rbSequence := recordBucket.Sequence()
	id := tbtypes.BotID(rbSequence)
	// delete the record
	bid, err := rivbin.Marshal(id)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	err = recordBucket.Delete(bid)
	if err != nil {
		return fmt.Errorf("error while deleting record for bot %d with public key %v: %v",
			id, brtx.Identification.PublicKey, err)
	}
	// delete the name->ID mappings
	for _, name := range brtx.Names {
		err = revertNameToIDMapping(bucket, name)
		if err != nil {
			return fmt.Errorf("error while deleting name %s to bot id %d mapping: %v", name.String(), id, err)
		}
	}
	// delete the publicKey->ID mapping,
	// doing it last as this is the initial check that happens when registering a bot,
	// as to ensure we only have one bot per public key
	err = revertKeyToIDMapping(bucket, brtx.Identification.PublicKey)
	if err != nil {
		return fmt.Errorf("error while deleting pubKey %s to bot id %d mapping: %v", brtx.Identification.PublicKey, id, err)
	}
	// revert the transactionID from the list of transactionIDs for the given bot
	err = revertBotTransaction(bucket, id, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID))
	if err != nil {
		return fmt.Errorf("error while reverting transaction for bot %d: %v", id, err)
	}
	// decrease the sequence counter of the bucket
	err = recordBucket.SetSequence(rbSequence - 1)
	if err != nil {
		return fmt.Errorf("error while decrementing the sequence counter of bot record bucket: %v", err)
	}
	// all information is reverted
	return nil
}

func (p *Plugin) revertRecordUpdateTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	recordBucket, err := bucket.Bucket(bucketBotRecords)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	brutx, err := tbtypes.BotRecordUpdateTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot record update tx type: %v", err)
	}

	// get the bot record
	bid, err := rivbin.Marshal(brutx.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	b := recordBucket.Get(bid)
	if len(b) == 0 {
		return errors.New("no bot record found for the specified identifier")
	}
	var record tbtypes.BotRecord
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found bot record: %v", err)
	}

	// first revert the Tx
	err = brutx.RevertBotRecordUpdate(&record)
	if err != nil {
		return fmt.Errorf("failed to revert bot record: %v", err)
	}

	// now check if we're expired, if so,
	// we might need to revert an implicit update that happened during the apply phase of this Tx
	if record.IsExpired(txn.BlockTime) {
		txID := txn.ID()

		// if the record is expired, there is a big chance that
		// there was an implicit update, as such let's try to get it,
		// if an implicit update did indeed take place, let's restore that info
		update, err := getImplicitBotRecordUpdate(bucket, txID)
		if err != nil {
			return fmt.Errorf("failed to revert bot record: failed to fetch implicit record update: %v", err)
		}

		if update.PreviousExpirationTime != 0 {
			// only if the previous expiration time was set,
			// could there have been an update (in the current setup)
			record.Expiration = update.PreviousExpirationTime

			// add all names to the record again (they were removed as well, so no error is expected)
			err = record.AddNames(update.InactiveNamesRemoved...)
			if err != nil {
				return fmt.Errorf("failed to revert bot record: unexpected error while adding implicitly removed names back to record: %v", err)
			}

			// apply all removed names again in the mapping bucket as well (for all those names that aren't taken in the meanwhile)
			for _, name := range update.InactiveNamesRemoved {
				err = applyNameToIDMappingIfAvailable(bucket, name, record.ID)
				if err != nil {
					return fmt.Errorf("failed to revert bot record: :"+
						"failed to add back mapping of expired bot's name %v to its ID %d: %v", name, record.ID, err)
				}
			}

			// delete the implicit record update, it is no longer required
			err = revertImplicitBotRecordUpdate(bucket, txID)
			if err != nil {
				return fmt.Errorf("failed to revert bot record: failed to revert implicit bot record update content:%v", err)
			}
		}
	}

	// save it
	bid, err = rivbin.Marshal(brutx.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	brecord, err := rivbin.Marshal(brutx.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while updating record for bot %d: %v", brutx.Identifier, err)
	}

	// revert all names that were added
	for _, name := range brutx.Names.Add {
		err = revertNameToIDMapping(bucket, name)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// apply all names again that were removed,
	// which can only be in case the bot was active
	for _, name := range brutx.Names.Remove {
		err = applyNameToIDMapping(bucket, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to revert update bot record: error while revert mapping of name %v that was removed: %v", name, err)
		}
	}

	// revert the transactionID from the list of transactionIDs for the given bot
	err = revertBotTransaction(bucket, record.ID, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID))
	if err != nil {
		return fmt.Errorf("error while reverting transaction for bot %d: %v", record.ID, err)
	}

	// all information is applied
	return nil
}

func (p *Plugin) revertBotNameTransferTx(txn modules.ConsensusTransaction, bucket *persist.LazyBoltBucket) error {
	recordBucket, err := bucket.Bucket(bucketBotRecords)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bnttx, err := tbtypes.BotNameTransferTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot name transfer tx type: %v", err)
	}

	// get the receiver bot record
	bid, err := rivbin.Marshal(bnttx.Receiver.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	b := recordBucket.Get(bid)
	if len(b) == 0 {
		return errors.New("no bot record found for the specified receiver bot identifier")
	}
	var record tbtypes.BotRecord
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found record of receiver bot %d: %v", bnttx.Receiver.Identifier, err)
	}
	// revert receiver bot
	err = bnttx.RevertReceiverBotRecordUpdate(&record)
	if err != nil {
		return fmt.Errorf("failed to update record of receiver bot %d: %v", record.ID, err)
	}
	// save the record of the receiver bot
	bid, err = rivbin.Marshal(record.ID)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	brecord, err := rivbin.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while saving the reverted (update of the) receiver bot %d: %v", record.ID, err)
	}

	// revert the transactionID from the list of transactionIDs for the receiver bot
	err = revertBotTransaction(bucket, record.ID, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID))
	if err != nil {
		return fmt.Errorf("error while reverting transaction for receiver bot %d: %v", record.ID, err)
	}

	// get the sender bot record
	bid, err = rivbin.Marshal(bnttx.Sender.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	b = recordBucket.Get(bid)
	if len(b) == 0 {
		return errors.New("no bot record found for the specified sender bot identifier")
	}
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found record of sender bot %d: %v", bnttx.Sender.Identifier, err)
	}
	// revert sender bot
	err = bnttx.RevertSenderBotRecordUpdate(&record)
	if err != nil {
		return fmt.Errorf("failed to revert record of sender bot %d: %v", record.ID, err)
	}
	// save the record of the sender bot
	bid, err = rivbin.Marshal(bnttx.Receiver.Identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal bot ID: %v", err)
	}
	brecord, err = rivbin.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal bot Record: %v", err)
	}
	err = recordBucket.Put(bid, brecord)
	if err != nil {
		return fmt.Errorf("error while saving the reverted record for sender bot %d: %v", record.ID, err)
	}

	// update mapping for all the transferred names, reverting them back to the sender
	for _, name := range bnttx.Names {
		err = applyNameToIDMapping(bucket, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while mapping name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// revert the transactionID from the list of transactionIDs for the sender bot
	err = revertBotTransaction(bucket, record.ID, newSortableTransactionShortID(txn.BlockHeight, txn.SequenceID))
	if err != nil {
		return fmt.Errorf("error while reverting transaction for sender bot %d: %v", record.ID, err)
	}

	// revert went fine
	return nil
}

// TransactionValidators returns all tx validators linked to this plugin
func (p *Plugin) TransactionValidators() []modules.PluginTransactionValidationFunction {
	return nil
}

// TransactionValidatorVersionFunctionMapping returns all tx validators linked to this plugin
func (p *Plugin) TransactionValidatorVersionFunctionMapping() map[types.TransactionVersion][]modules.PluginTransactionValidationFunction {
	return map[types.TransactionVersion][]modules.PluginTransactionValidationFunction{
		tbtypes.TransactionVersionBotRegistration: {
			p.validateBotRegistrationTx,
		},
		tbtypes.TransactionVersionBotRecordUpdate: {
			p.validateBotUpdateTx,
		},
		tbtypes.TransactionVersionBotNameTransfer: {
			p.validateBotNameTransferTx,
		},
	}
}

func (p *Plugin) validateBotRegistrationTx(txn modules.ConsensusTransaction, ctx types.TransactionValidationContext, bucket *persist.LazyBoltBucket) error {
	// get BotRegistrationTx
	brtx, err := tbtypes.BotRegistrationTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("failed to use tx as a bot registration tx: %v", err)
	}

	// look up the public key, to ensure it is not registered yet
	_, err = p.GetRecordForKey(brtx.Identification.PublicKey)
	if err == nil {
		return tbtypes.ErrBotKeyAlreadyRegistered
	}
	if err != tbtypes.ErrBotKeyNotFound {
		return fmt.Errorf("unexpected error while validating non-existence of bot's public key: %v", err)
	}

	// validate the signature of the to-be-registered bot
	err = validateBotSignature(txn.Transaction, brtx.Identification.PublicKey, brtx.Identification.Signature, ctx, tbtypes.BotSignatureSpecifierSender)
	if err != nil {
		return fmt.Errorf("failed to fulfill bot registration condition: %v", err)
	}

	// ensure the NrOfMonths is in the inclusive range of [1, 24]
	if brtx.NrOfMonths == 0 {
		return errors.New("bot registration requires at least one month to be paid already")
	}
	if brtx.NrOfMonths > tbtypes.MaxBotPrepaidMonths {
		return tbtypes.ErrBotExpirationExtendOverflow
	}

	// validate the lengths,
	// and ensure that at least one name or one addr is registered
	addrLen := len(brtx.Addresses)
	if addrLen > tbtypes.MaxAddressesPerBot {
		return tbtypes.ErrTooManyBotAddresses
	}
	nameLen := len(brtx.Names)
	if nameLen > tbtypes.MaxNamesPerBot {
		return tbtypes.ErrTooManyBotNames
	}
	if addrLen == 0 && nameLen == 0 {
		return errors.New("bot registration requires a name or address to be defined")
	}

	// validate that all network addresses are unique
	err = validateUniquenessOfNetworkAddresses(brtx.Addresses)
	if err != nil {
		return fmt.Errorf("invalid bot registration Tx: validateUniquenessOfNetworkAddresses: %v", err)
	}

	// validate that all names are unique
	err = validateUniquenessOfBotNames(brtx.Names)
	if err != nil {
		return fmt.Errorf("invalid bot registration Tx: validateUniquenessOfBotNames: %v", err)
	}

	// validate that the names are not registered yet
	for _, name := range brtx.Names {
		_, err = p.GetRecordForName(name)
		if err == nil {
			return tbtypes.ErrBotNameAlreadyRegistered
		}
		if err != tbtypes.ErrBotNameNotFound {
			return fmt.Errorf(
				"unexpected error while validating non-existence of bot's name %v: %v",
				name, err)
		}
	}

	// validate the miner fee
	if brtx.TransactionFee.Cmp(ctx.MinimumMinerFee) == -1 {
		return types.ErrTooSmallMinerFee
	}
	return nil
}

func (p *Plugin) validateBotUpdateTx(txn modules.ConsensusTransaction, ctx types.TransactionValidationContext, bucket *persist.LazyBoltBucket) error {
	// get BotRecordUpdateTx
	brutx, err := tbtypes.BotRecordUpdateTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("failed to use tx as a bot record update tx: %v", err)
	}

	// validate the miner fee
	if brutx.TransactionFee.Cmp(ctx.MinimumMinerFee) == -1 {
		return types.ErrTooSmallMinerFee
	}

	// look up the record, using the given ID, to ensure it is registered
	record, err := p.GetRecordForID(brutx.Identifier)
	if err != nil {
		return fmt.Errorf("bot cannot be updated: GetRecordForID(%v): %v", brutx.Identifier, err)
	}

	// validate the signature of the to-be-updated bot
	err = validateBotSignature(txn.Transaction, record.PublicKey, brutx.Signature, ctx, tbtypes.BotSignatureSpecifierSender)
	if err != nil {
		return fmt.Errorf("failed to fulfill bot record update condition: %v", err)
	}

	// at least something has to be updated, a nop-update is not allowed
	if brutx.NrOfMonths == 0 &&
		len(brutx.Addresses.Add) == 0 && len(brutx.Addresses.Remove) == 0 &&
		len(brutx.Names.Add) == 0 && len(brutx.Names.Remove) == 0 {
		return errors.New("bot record updates requires nrOfMonths, a name or address to be defined")
	}

	// ensure all to-be-added names are available
	err = areBotNamesAvailable(p, brutx.Names.Add...)
	if err != nil {
		return fmt.Errorf("bot cannot be updated: areBotNamesAvailable: %v", err)
	}

	// try to update the record, to spot any errors should that happen for real
	err = brutx.UpdateBotRecord(ctx.BlockTime, record)
	if err != nil {
		return fmt.Errorf("bot cannot be updated: UpdateBotRecord: %v", err)
	}

	// update Tx is valid
	return nil
}

func (p *Plugin) validateBotNameTransferTx(txn modules.ConsensusTransaction, ctx types.TransactionValidationContext, bucket *persist.LazyBoltBucket) error {
	// get BotRecordUpdateTx
	bnttx, err := tbtypes.BotNameTransferTransactionFromTransaction(txn.Transaction)
	if err != nil {
		return fmt.Errorf("failed to use tx as a bot name transfer tx: %v", err)
	}

	// validate the miner fee
	if bnttx.TransactionFee.Cmp(ctx.MinimumMinerFee) == -1 {
		return types.ErrTooSmallMinerFee
	}

	// validate the sender/receiver ID is different
	if bnttx.Sender.Identifier == bnttx.Receiver.Identifier {
		return errors.New("the identifiers of the sender and receiver bot have to be different")
	}

	// look up the record of the sender, using the given (sender) ID, to ensure it is registered,
	// as well as for validation checks that follow
	recordSender, err := p.GetRecordForID(bnttx.Sender.Identifier)
	if err != nil {
		return fmt.Errorf("invalid sender (%d) of bot name transfer: %v", bnttx.Sender.Identifier, err)
	}

	// look up the record of the sender, using the given (sender) ID, to ensure it is registered,
	// as well as for validation checks that follow
	recordReceiver, err := p.GetRecordForID(bnttx.Receiver.Identifier)
	if err != nil {
		return fmt.Errorf("invalid sender (%d) of bot name transfer: %v", bnttx.Receiver.Identifier, err)
	}

	// validate the signature of the sender
	err = validateBotSignature(txn.Transaction, recordSender.PublicKey, bnttx.Sender.Signature, ctx, tbtypes.BotSignatureSpecifierSender)
	if err != nil {
		return fmt.Errorf("failed to fulfill bot record name transfer condition of the sender: %v", err)
	}
	// validate the signature of the receiver
	err = validateBotSignature(txn.Transaction, recordReceiver.PublicKey, bnttx.Receiver.Signature, ctx, tbtypes.BotSignatureSpecifierReceiver)
	if err != nil {
		return fmt.Errorf("failed to fulfill bot record name transfer condition of the receiver: %v", err)
	}

	// at least one name has to be transferred
	if len(bnttx.Names) == 0 {
		return errors.New("a bot name transfer transaction has to transfer at least one name")
	}

	// try to update the sender bot (if the sender bot is expired, an error is returned as well)
	err = bnttx.UpdateSenderBotRecord(ctx.BlockTime, recordSender)
	if err != nil {
		return fmt.Errorf("sender bot (%v) cannot be updated by name transfer: %v", bnttx.Sender.Identifier, err)
	}

	// try to update the receiver bot
	// (the sender bot doesn't need this validation,
	// as we already checked that it owns the address, the only update to that bot)
	err = bnttx.UpdateReceiverBotRecord(ctx.BlockTime, recordReceiver)
	if err != nil {
		return fmt.Errorf("receiver bot (%v) cannot be updated by name transfer: %v", bnttx.Receiver.Identifier, err)
	}

	// given all names originate from the sender,
	// we do not require availability checks of names, as no names will be available at this point

	// name transfer Tx is valid
	return nil
}

// Close unregisters the plugin from the consensus
func (p *Plugin) Close() error {
	return p.storage.Close()
}

// encodeBlockheight encodes the given blockheight as a sortable key
func encodeBlockheight(height types.BlockHeight) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key[:], uint64(height))
	return key
}

// eecodeBlockheight decodes the given sortable key as a blockheight
func decodeBlockheight(key []byte) types.BlockHeight {
	return types.BlockHeight(binary.BigEndian.Uint64(key))
}

// apply/revert the Key->ID mapping for a 3bot
func applyKeyToIDMapping(bucket *persist.LazyBoltBucket, key types.PublicKey, id tbtypes.BotID) error {
	mappingBucket, err := bucket.Bucket(bucketBotKeyToIDMapping)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bKey, err := rivbin.Marshal(key)
	if err != nil {
		return err
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return err
	}
	return mappingBucket.Put(bKey, bID)
}
func revertKeyToIDMapping(bucket *persist.LazyBoltBucket, key types.PublicKey) error {
	mappingBucket, err := bucket.Bucket(bucketBotKeyToIDMapping)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bKey, err := rivbin.Marshal(key)
	if err != nil {
		return err
	}
	return mappingBucket.Delete(bKey)
}

// apply/revert the Name->ID mapping for a 3bot
func applyNameToIDMapping(bucket *persist.LazyBoltBucket, name tbtypes.BotName, id tbtypes.BotID) error {
	mappingBucket, err := bucket.Bucket(bucketBotNameToIDMapping)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bName, err := rivbin.Marshal(name)
	if err != nil {
		return err
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return err
	}
	return mappingBucket.Put(bName, bID)
}
func revertNameToIDMapping(bucket *persist.LazyBoltBucket, name tbtypes.BotName) error {
	mappingBucket, err := bucket.Bucket(bucketBotNameToIDMapping)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bName, err := rivbin.Marshal(name)
	if err != nil {
		return err
	}
	return mappingBucket.Delete(bName)
}
func revertNameToIDMappingIfOwnedByBot(bucket *persist.LazyBoltBucket, name tbtypes.BotName, id tbtypes.BotID) error {
	mappingBucket, err := bucket.Bucket(bucketBotNameToIDMapping)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bName, err := rivbin.Marshal(name)
	if err != nil {
		return err
	}
	b := mappingBucket.Get(bName)
	if len(b) == 0 {
		return nil // might be deleted by another bot, who took over ownership
	}
	var mappedID tbtypes.BotID
	err = rivbin.Unmarshal(b, &mappedID)
	if err != nil {
		return fmt.Errorf("corrupt BotID used as key in mapping of bot name %v", name)
	}
	if mappedID != id {
		return nil // ID no longer owned by this bot, ignore removal request in the mapping context
	}
	// delete name (mapping), as it was still owned by this bot
	bName, err = rivbin.Marshal(name)
	if err != nil {
		return err
	}
	return mappingBucket.Delete(bName)
}
func applyNameToIDMappingIfAvailable(bucket *persist.LazyBoltBucket, name tbtypes.BotName, id tbtypes.BotID) error {
	mappingBucket, err := bucket.Bucket(bucketBotNameToIDMapping)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bName, err := rivbin.Marshal(name)
	if err != nil {
		return err
	}
	b := mappingBucket.Get(bName)
	if len(b) != 0 {
		return nil // already taken
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return err
	}
	return mappingBucket.Put(bName, bID)
}

func applyBotTransaction(bucket *persist.LazyBoltBucket, id tbtypes.BotID, shortTxID sortableTransactionShortID, txID types.TransactionID) error {
	txBucket, err := bucket.Bucket(bucketBotTransactions)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return err
	}
	botBucket, err := txBucket.CreateBucketIfNotExists(bID)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: failed to create/get bot %d inner bucket: %v", id, err)
	}
	bShortTxID, err := rivbin.Marshal(shortTxID)
	if err != nil {
		return err
	}
	bTxID, err := rivbin.Marshal(txID)
	if err != nil {
		return err
	}
	return botBucket.Put(bShortTxID, bTxID)
}
func revertBotTransaction(bucket *persist.LazyBoltBucket, id tbtypes.BotID, shortTxID sortableTransactionShortID) error {
	txBucket, err := bucket.Bucket(bucketBotTransactions)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return err
	}
	botBucket := txBucket.Bucket(bID)
	if botBucket != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: bot %d inner bucket does not exist", id)
	}
	bShortTxID, err := rivbin.Marshal(shortTxID)
	if err != nil {
		return err
	}
	return botBucket.Delete(bShortTxID)
}
func getBotTransactions(bucket *bolt.Bucket, id tbtypes.BotID) ([]types.TransactionID, error) {
	txBucket := bucket.Bucket(bucketBotTransactions)
	if txBucket == nil {
		return nil, fmt.Errorf("corrupt 3bot plugin DB: no bucket with name %s found", string(bucketBotTransactions))
	}
	bID, err := rivbin.Marshal(id)
	if err != nil {
		return nil, err
	}
	botBucket := txBucket.Bucket(bID)
	if botBucket == nil {
		return nil, nil // no transactions is acceptable
	}
	var txIDs []types.TransactionID
	err = botBucket.ForEach(func(_, v []byte) (err error) {
		var txID types.TransactionID
		err = rivbin.Unmarshal(v, &txID)
		txIDs = append(txIDs, txID)
		return
	})
	if err != nil {
		return nil, fmt.Errorf("corrupt 3bot plugin DB: error while parsing stored txID for bot %d: %v", id, err)
	}
	return txIDs, nil
}

func setStatsBlockTime(blockTimeBucket *bolt.Bucket, height types.BlockHeight, time types.Timestamp) error {
	// validate blockheight
	expectedHeight := types.BlockHeight(blockTimeBucket.Sequence())
	if expectedHeight != height {
		return fmt.Errorf("corrupt 3bot plugin DB: unexpected block height %d, expected %d", height, expectedHeight)
	}
	// store time using the height as ID
	bHeight, err := rivbin.Marshal(height)
	if err != nil {
		return err
	}
	bTime, err := rivbin.Marshal(time)
	if err != nil {
		return err
	}
	err = blockTimeBucket.Put(bHeight, bTime)
	if err != nil {
		return err
	}
	// increase the bucket's sequence
	_, err = blockTimeBucket.NextSequence()
	return err
}
func getStatsBlockTime(blockTimeBucket *bolt.Bucket, height types.BlockHeight) (types.Timestamp, error) {
	// get timestamp and return it
	bHeight, err := rivbin.Marshal(height)
	if err != nil {
		return 0, err
	}
	b := blockTimeBucket.Get(bHeight)
	if len(b) == 0 {
		return 0, fmt.Errorf("no timestamp found for block %d", height)
	}
	var ts types.Timestamp
	err = rivbin.Unmarshal(b, &ts)
	return ts, err
}
func getCurrentBlockHeightAndTime(blockTimeBucket *bolt.Bucket) (types.BlockHeight, types.Timestamp, error) {
	// get current blockheight
	seq := blockTimeBucket.Sequence()
	if seq == 0 {
		return 0, 0, errors.New("bucket contains no block info")
	}
	height := types.BlockHeight(seq - 1)
	ts, err := getStatsBlockTime(blockTimeBucket, height)
	return height, ts, err

}
func deleteStatsBlockTime(blockTimeBucket *bolt.Bucket, height types.BlockHeight) error {
	// validate the given height
	seq := blockTimeBucket.Sequence()
	if seq == 0 {
		return errors.New("bucket contains no block info")
	}
	expectedHeight := types.BlockHeight(seq - 1)
	if expectedHeight != height {
		return fmt.Errorf("corrupt 3bot plugin DB: unexpected block height %d, expected %d", height, expectedHeight)
	}
	// delete time for height
	bHeight, err := rivbin.Marshal(height)
	if err != nil {
		return err
	}
	err = blockTimeBucket.Delete(bHeight)
	if err != nil {
		return err
	}
	// decrease the bucket's sequence
	return blockTimeBucket.SetSequence(uint64(height))
}

func applyImplicitBotRecordUpdate(bucket *persist.LazyBoltBucket, txID types.TransactionID, update implicitBotRecordUpdate) error {
	updateBucket, err := bucket.Bucket(bucketBotRecordImplicitUpdates)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bTxID, err := rivbin.Marshal(txID)
	if err != nil {
		return err
	}
	bUpdate, err := rivbin.Marshal(update)
	if err != nil {
		return err
	}
	return updateBucket.Put(bTxID, bUpdate)
}
func revertImplicitBotRecordUpdate(bucket *persist.LazyBoltBucket, txID types.TransactionID) error {
	updateBucket, err := bucket.Bucket(bucketBotRecordImplicitUpdates)
	if err != nil {
		return fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bTxID, err := rivbin.Marshal(txID)
	if err != nil {
		return err
	}
	return updateBucket.Delete(bTxID)
}
func getImplicitBotRecordUpdate(bucket *persist.LazyBoltBucket, txID types.TransactionID) (implicitBotRecordUpdate, error) {
	var update implicitBotRecordUpdate

	updateBucket, err := bucket.Bucket(bucketBotRecordImplicitUpdates)
	if err != nil {
		return implicitBotRecordUpdate{}, fmt.Errorf("corrupt 3bot plugin DB: %v", err)
	}
	bTxID, err := rivbin.Marshal(txID)
	if err != nil {
		return implicitBotRecordUpdate{}, err
	}
	b := updateBucket.Get(bTxID)
	if len(b) == 0 {
		return update, nil
	}
	err = rivbin.Unmarshal(b, &update)
	if err != nil {
		return update, fmt.Errorf("failed to fetch implicit record update for tx %v: %v", txID, err)
	}
	return update, nil
}

// implicitBotRecordUpdate collects all info that was erased/changed due to
// implicit updates to a bot record as part of a record update Tx.
// Such an implicit update is possible in case the bot was made active again by the update Tx,
// after being inactive prior to that.
type implicitBotRecordUpdate struct {
	PreviousExpirationTime tbtypes.CompactTimestamp
	InactiveNamesRemoved   []tbtypes.BotName
}

// MarshalSia implements SiaMarshaler.MarshalSia,
// alias of MarshalRivine for backwards-compatibility reasons.
func (ibru implicitBotRecordUpdate) MarshalSia(w io.Writer) error {
	return ibru.MarshalRivine(w)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// alias of UnmarshalRivine for backwards-compatibility reasons.
func (ibru *implicitBotRecordUpdate) UnmarshalSia(r io.Reader) error {
	return ibru.UnmarshalRivine(r)
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (ibru implicitBotRecordUpdate) MarshalRivine(w io.Writer) error {
	// write the info prefix byte first, containing
	// flags and length information
	var infoPrefixByte uint8
	if ibru.PreviousExpirationTime != 0 {
		infoPrefixByte |= 1 // first bit indicates if a previous expiration time is to be added
	}
	lenNamesRemoved := len(ibru.InactiveNamesRemoved)
	if lenNamesRemoved > tbtypes.MaxNamesPerBot {
		return errors.New("too many bot names were implicitly removed")
	}
	// next 3 bits indicate the length of names that were implicitly removed
	infoPrefixByte |= uint8(lenNamesRemoved) << 1
	// write the info prefix byte
	err := rivbin.MarshalUint8(w, infoPrefixByte)
	if err != nil {
		return fmt.Errorf("failed to write the info (implicit bot record update) prefix byte: %v", err)
	}

	// write the previous expiration time if non-0
	if ibru.PreviousExpirationTime != 0 {
		err = ibru.PreviousExpirationTime.MarshalRivine(w)
		if err != nil {
			return fmt.Errorf("implicitBotRecordUpdate: %v", err)
		}
	}

	// write all names one by one
	for _, name := range ibru.InactiveNamesRemoved {
		err = name.MarshalRivine(w)
		if err != nil {
			return fmt.Errorf("implicitBotRecordUpdate: %v", err)
		}
	}

	// all written succesfully
	return nil
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (ibru *implicitBotRecordUpdate) UnmarshalRivine(r io.Reader) error {
	infoPrefixByte, err := rivbin.UnmarshalUint8(r)
	if err != nil {
		return fmt.Errorf("implicitBotRecordUpdate: %v", err)
	}
	// read or reset the PreviousExpirationTime
	if infoPrefixByte&1 != 0 {
		err = ibru.PreviousExpirationTime.UnmarshalRivine(r)
		if err != nil {
			return fmt.Errorf("implicitBotRecordUpdate: %v", err)
		}
	} else {
		ibru.PreviousExpirationTime = 0
	}
	// read or reset the implicitly removed names
	length := (infoPrefixByte >> 1) & 7
	if length == 0 {
		ibru.InactiveNamesRemoved = nil
		return nil
	}
	if length > tbtypes.MaxNamesPerBot {
		return errors.New("too many bot names were implicitly removed")
	}
	ibru.InactiveNamesRemoved = make([]tbtypes.BotName, 0, length)
	for i := uint8(0); i < length; i++ {
		var name tbtypes.BotName
		err = name.UnmarshalRivine(r)
		if err != nil {
			return fmt.Errorf("implicitBotRecordUpdate: %v", err)
		}
		ibru.InactiveNamesRemoved = append(ibru.InactiveNamesRemoved, name)
	}

	// all read succesfully
	return nil
}

// sortableTransactionShortID wraps around the types.TransactionShortID,
// as to ensure it is encoded in a way that allows boltdb use it for natural ordering.
type sortableTransactionShortID types.TransactionShortID

func newSortableTransactionShortID(height types.BlockHeight, txSequenceID uint16) sortableTransactionShortID {
	return sortableTransactionShortID(types.NewTransactionShortID(height, txSequenceID))
}

// MarshalSia implements SiaMarshaler.MarshalSia,
// alias of MarshalRivine for backwards-compatibility reasons.
func (sid sortableTransactionShortID) MarshalSia(w io.Writer) error {
	return sid.MarshalRivine(w)
}

// UnmarshalSia implements SiaMarshaler.UnmarshalSia,
// alias of UnmarshalRivine for backwards-compatibility reasons.
func (sid *sortableTransactionShortID) UnmarshalSia(r io.Reader) error {
	return sid.UnmarshalRivine(r)
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (sid sortableTransactionShortID) MarshalRivine(w io.Writer) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(sid))
	n, err := w.Write(b[:])
	if err != nil {
		return err
	}
	if n != 8 {
		return io.ErrShortWrite
	}
	return nil
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (sid *sortableTransactionShortID) UnmarshalRivine(r io.Reader) error {
	var b [8]byte
	n, err := r.Read(b[:])
	if err != nil {
		return err
	}
	if n != 8 {
		return io.ErrUnexpectedEOF
	}
	*sid = sortableTransactionShortID(binary.BigEndian.Uint64(b[:]))
	return nil
}

func areBotNamesAvailable(registry tbtypes.BotRecordReadRegistry, names ...tbtypes.BotName) error {
	var err error
	for _, name := range names {
		_, err = registry.GetRecordForName(name)
		switch err {
		case tbtypes.ErrBotNameNotFound, tbtypes.ErrBotNameExpired:
			continue // name is available, check the others
		case nil:
			// when no error is returned a record is returned,
			// meaning the name is linked to a non-expired 3bot,
			// and consequently the name is not available
			return tbtypes.ErrBotNameAlreadyRegistered
		default:
			return err // unexpected
		}
	}
	return nil
}
