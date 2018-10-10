package persist

import (
	"errors"
	"fmt"
	"os"
	"path"

	tfencoding "github.com/threefoldfoundation/tfchain/pkg/encoding"
	"github.com/threefoldfoundation/tfchain/pkg/persist/internal"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/rivine/rivine/build"
	"github.com/rivine/rivine/encoding"
	"github.com/rivine/rivine/modules"
	"github.com/rivine/rivine/persist"
	rivinesync "github.com/rivine/rivine/sync"
	rivinetypes "github.com/rivine/rivine/types"

	bolt "github.com/rivine/bbolt"
)

// TransactionDB I/O constants
const (
	TransactionDBDir      = "transactiondb"
	TransactionDBFilename = TransactionDBDir + ".db"
)

// TODO:
//   - modify (function godoc) comments to take into account that now we also store/delete/manage 3bots, not just mintconditions

// TODO:
//  ensure that we cannot add 2 (or more) (bot) Tx's in the same block,
//  that use the same pubKey or name
//     for double spending, it seems to ensure no double spending happens, even when
//     spending an uncomfirmed spended coin output, as can be seen from following output:
//        Could not publish transaction: HTTP 400 error: error after call to
//        /wallet/transactions: consensus conflict: transaction spends a nonexisting coin output
//     ... need to figure out how this works, and also make such logic possible
//      for the checking of unique names and unique addresses
//     see: https://github.com/threefoldfoundation/tfchain/issues/195

// internal bucket database keys used for the transactionDB
var (
	bucketInternal         = []byte("internal")
	bucketInternalKeyStats = []byte("stats") // stored as a single struct, see `transactionDBStats`

	// getBucketMintConditionPerHeightRangeKey is used to compute the keys
	// of the values in this bucket
	bucketMintConditions = []byte("mintconditions")

	// buckets for the 3bot feature
	bucketBotRecords         = []byte("botrecords") // ID => name
	bucketBotKeyToIDMapping  = []byte("botkeys")    // Key => ID
	bucketBotNameToIDMapping = []byte("botnames")   // Name => ID
)

type (
	// TransactionDB extends Rivine's ConsensusSet module,
	// allowing us to track transactions (and specifically parts of it) that we care about,
	// and for which Rivine does not implement any logic.
	//
	// The initial motivation (and currently only use case) is to track MintConditions,
	// as to be able to know for any given block height what the active MintCondition is,
	// but other use cases can be supported in future updates should they appear.
	TransactionDB struct {
		// The DB's ThreadGroup tells tracked functions to shut down and
		// blocks until they have all exited before returning from Close.
		tg rivinesync.ThreadGroup

		db    *persist.BoltDatabase
		stats transactionDBStats

		subscriber *transactionDBCSSubscriber
	}

	// implements modules.ConsensusSetSubscriber,
	// such that the TransactionDB does not have to publicly implement
	// the ConsensusSetSubscriber interface, allowing us to "force"
	// the user to register to the consensus set using our provided
	// (*TransactionDB).SubscribeToConsensusSet method
	transactionDBCSSubscriber struct {
		txdb *TransactionDB
		cs   modules.ConsensusSet
	}
	transactionDBStats struct {
		ConsensusChangeID modules.ConsensusChangeID
		BlockHeight       rivinetypes.BlockHeight
		Synced            bool
	}
)

var (
	// ensure TransactionDB implements the MintConditionGetter interface
	_ types.MintConditionGetter = (*TransactionDB)(nil)
	// enssure TransactionDB implements the BotRecordReadRegistry interface
	_ types.BotRecordReadRegistry = (*TransactionDB)(nil)
)

// NewTransactionDB creates a new TransactionDB, using the given file (path) to store the (single) persistent BoltDB file.
// A new db will be created if it doesn't exist yet, if it does exist it should be ensured that the given genesis mint condition
// equals the already stored genesis mint condition.
func NewTransactionDB(rootDir string, genesisMintCondition rivinetypes.UnlockConditionProxy) (*TransactionDB, error) {
	persistDir := path.Join(rootDir, TransactionDBDir)
	// Create the directory if it doesn't exist.
	err := os.MkdirAll(persistDir, 0700)
	if err != nil {
		return nil, err
	}

	txdb := new(TransactionDB)
	err = txdb.openDB(path.Join(persistDir, TransactionDBFilename), genesisMintCondition)
	if err != nil {
		return nil, fmt.Errorf("failed to open the transaction DB: %v", err)
	}
	return txdb, nil
}

// SubscribeToConsensusSet subscribes the TransactionDB to the given ConsensusSet,
// allowing it to stay in sync with the blockchain, and also making it automatically unsubscribe
// from the consensus set when the TransactionDB is closed (using (*TransactionDB).Close).
func (txdb *TransactionDB) SubscribeToConsensusSet(cs modules.ConsensusSet) error {
	if txdb.subscriber != nil {
		return errors.New("transactionDB is already subscribed to a consensus set")
	}

	subscriber := &transactionDBCSSubscriber{txdb: txdb, cs: cs}
	err := cs.ConsensusSetSubscribe(
		subscriber,
		txdb.stats.ConsensusChangeID,
		txdb.tg.StopChan(),
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to consensus set: %v", err)
	}
	txdb.subscriber = subscriber
	return nil
}

// GetActiveMintCondition implements types.MintConditionGetter.GetActiveMintCondition
func (txdb *TransactionDB) GetActiveMintCondition() (rivinetypes.UnlockConditionProxy, error) {
	var b []byte
	err := txdb.db.View(func(tx *bolt.Tx) (err error) {
		mintConditionsBucket := tx.Bucket(bucketMintConditions)
		if mintConditionsBucket == nil {
			return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
		}

		// return the last cursor
		cursor := mintConditionsBucket.Cursor()

		var k []byte
		k, b = cursor.Last()
		if len(k) == 0 {
			return errors.New("corrupt transaction DB: no matching mint condition could be found")
		}
		return nil
	})
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, err
	}

	var mintCondition rivinetypes.UnlockConditionProxy
	err = encoding.Unmarshal(b, &mintCondition)
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf("corrupt transaction DB: failed to decode found mint condition: %v", err)
	}
	// mint condition found, return it
	return mintCondition, nil
}

// GetMintConditionAt implements types.MintConditionGetter.GetMintConditionAt
func (txdb *TransactionDB) GetMintConditionAt(height rivinetypes.BlockHeight) (rivinetypes.UnlockConditionProxy, error) {
	var b []byte
	err := txdb.db.View(func(tx *bolt.Tx) (err error) {
		mintConditionsBucket := tx.Bucket(bucketMintConditions)
		if mintConditionsBucket == nil {
			return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
		}

		cursor := mintConditionsBucket.Cursor()

		var k []byte
		k, b = cursor.Seek(internal.EncodeBlockheight(height))
		if len(k) == 0 {
			// could be that we're past the last key, let's try the last key first
			k, b = cursor.Last()
			if len(k) == 0 {
				return errors.New("corrupt transaction DB: no matching mint condition could be found")
			}
			return nil
		}
		foundHeight := internal.DecodeBlockheight(k)
		if foundHeight <= height {
			return nil
		}
		k, b = cursor.Prev()
		if len(k) == 0 {
			return errors.New("corrupt transaction DB: no matching mint condition could be found")
		}
		return nil

	})
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, err
	}

	var mintCondition rivinetypes.UnlockConditionProxy
	err = encoding.Unmarshal(b, &mintCondition)
	if err != nil {
		return rivinetypes.UnlockConditionProxy{}, fmt.Errorf("corrupt transaction DB: failed to decode found mint condition: %v", err)
	}
	// mint condition found, return it
	return mintCondition, nil
}

// GetRecordForID returns the record mapped to the given BotID.
func (txdb *TransactionDB) GetRecordForID(id types.BotID) (record *types.BotRecord, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) (err error) {
		record, err = getRecordForID(tx, id)
		return
	})
	return
}

// internal function to get a record from the TxDB
func getRecordForID(tx *bolt.Tx, id types.BotID) (*types.BotRecord, error) {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return nil, errors.New("corrupt transaction DB: bot record bucket does not exist")
	}

	b := recordBucket.Get(tfencoding.Marshal(id))
	if len(b) == 0 {
		return nil, types.ErrBotNotFound
	}

	record := new(types.BotRecord)
	err := tfencoding.Unmarshal(b, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// GetRecordForKey returns the record mapped to the given Key.
func (txdb *TransactionDB) GetRecordForKey(key types.PublicKey) (record *types.BotRecord, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) error {
		id, err := getBotIDForPublicKey(tx, key)
		if err != nil {
			return err
		}

		record, err = getRecordForID(tx, id)
		return err
	})
	return
}

func getBotIDForPublicKey(tx *bolt.Tx, key types.PublicKey) (types.BotID, error) {
	keyBucket := tx.Bucket(bucketBotKeyToIDMapping)
	if keyBucket == nil {
		return 0, errors.New("corrupt transaction DB: bot key bucket does not exist")
	}

	b := keyBucket.Get(tfencoding.Marshal(key))
	if len(b) == 0 {
		return 0, types.ErrBotKeyNotFound
	}

	var id types.BotID
	err := tfencoding.Unmarshal(b, &id)
	return id, err
}

// GetRecordForName returns the record mapped to the given Name.
func (txdb *TransactionDB) GetRecordForName(name types.BotName) (record *types.BotRecord, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) error {
		nameBucket := tx.Bucket(bucketBotNameToIDMapping)
		if nameBucket == nil {
			return errors.New("corrupt transaction DB: bot name bucket does not exist")
		}

		b := nameBucket.Get(tfencoding.Marshal(name))
		if len(b) == 0 {
			return types.ErrBotNameNotFound
		}

		var id types.BotID
		err := tfencoding.Unmarshal(b, &id)
		if err != nil {
			return err
		}

		record, err = getRecordForID(tx, id)
		return err
	})
	return
}

// Close the transaction DB,
// meaning the db will be unsubscribed from the consensus set,
// as well the threadgroup will be stopped and the internal bolt db will be closed.
func (txdb *TransactionDB) Close() error {
	if txdb.db == nil {
		return errors.New("transactionDB is already closed or was never created")
	}

	// unsubscribe from the consensus set, if subscribed at all
	if txdb.subscriber != nil {
		txdb.subscriber.unsubscribe()
		txdb.subscriber = nil
	}
	// stop thread group
	tgErr := txdb.tg.Stop()
	if tgErr != nil {
		tgErr = fmt.Errorf("failed to stop the threadgroup of TransactionDB: %v", tgErr)
	}
	// close database
	dbErr := txdb.db.Close()
	if dbErr != nil {
		dbErr = fmt.Errorf("failed to close the internal bolt db of TransactionDB: %v", dbErr)
	}
	txdb.db = nil

	return build.ComposeErrors(tgErr, dbErr)
}

// openDB loads the set database and populates it with the necessary buckets
func (txdb *TransactionDB) openDB(filename string, genesisMintCondition rivinetypes.UnlockConditionProxy) (err error) {
	var (
		dbMetadata = persist.Metadata{
			Header:  "TFChain Transaction Database",
			Version: "1.1.1",
		}
	)

	txdb.db, err = persist.OpenDatabase(dbMetadata, filename)
	if err != nil {
		if err != persist.ErrBadVersion {
			return fmt.Errorf("error opening tfchain transaction database: %v", err)
		}
		// try to open the DB using the original version
		originalDBMetadata := persist.Metadata{
			Header:  "TFChain Transaction Database",
			Version: "1.1.0",
		}
		txdb.db, err = persist.OpenDatabase(originalDBMetadata, filename)
		if err != nil {
			return fmt.Errorf("error opening tfchain transaction database using v1.1.0: %v", err)
		}
		// create added buckets
		err = txdb.db.Update(func(tx *bolt.Tx) (err error) {
			// Enumerate and create the new database buckets.
			buckets := [][]byte{
				bucketBotRecords,
				bucketBotKeyToIDMapping,
				bucketBotNameToIDMapping,
			}
			for _, bucket := range buckets {
				_, err = tx.CreateBucket(bucket)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error adding v1.1.1 buckets to tfchain transaction database: %v", err)
		}
	}
	return txdb.db.Update(func(tx *bolt.Tx) (err error) {
		if txdb.dbInitialized(tx) {
			// db is already created, get the stored stats
			internalBucket := tx.Bucket(bucketInternal)
			b := internalBucket.Get(bucketInternalKeyStats)
			if len(b) == 0 {
				return errors.New("structured stats value could not be found in existing transaction db")
			}
			err = encoding.Unmarshal(b, &txdb.stats)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structured stats value from existing transaction db: %v", err)
			}

			// and ensure the genesis mint condition is the same as the given one
			mintConditionsBucket := tx.Bucket(bucketMintConditions)
			b = mintConditionsBucket.Get(internal.EncodeBlockheight(0))
			if len(b) == 0 {
				return errors.New("genesis mint condition could not be found in existing transaction db")
			}
			var storedMintCondition rivinetypes.UnlockConditionProxy
			err = encoding.Unmarshal(b, &storedMintCondition)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis mint condition from existing transaction db: %v", err)
			}
			if !storedMintCondition.Equal(genesisMintCondition) {
				return errors.New("stored genesis mint condition is different from the given genesis mint condition")
			}

			return nil // nothing to do
		}

		// successfully create the DB
		err = txdb.createDB(tx, genesisMintCondition)
		if err != nil {
			return fmt.Errorf("failed to create transactionDB: %v", err)
		}
		return nil
	})
}

// dbInitialized returns true if the database appears to be initialized, false
// if not. Checking for the existence of the siafund pool bucket is typically
// sufficient to determine whether the database has gone through the
// initialization process.
func (txdb *TransactionDB) dbInitialized(tx *bolt.Tx) bool {
	return tx.Bucket(bucketInternal) != nil
}

// createConsensusObjects initialzes the consensus portions of the database.
func (txdb *TransactionDB) createDB(tx *bolt.Tx, genesisMintCondition rivinetypes.UnlockConditionProxy) (err error) {
	// Enumerate and create the database buckets.
	buckets := [][]byte{
		bucketInternal,
		bucketMintConditions,
		bucketBotRecords,
		bucketBotKeyToIDMapping,
		bucketBotNameToIDMapping,
	}
	for _, bucket := range buckets {
		_, err = tx.CreateBucket(bucket)
		if err != nil {
			return err
		}
	}

	// set the initial block height and initial consensus change iD
	txdb.stats.BlockHeight = 0
	txdb.stats.ConsensusChangeID = modules.ConsensusChangeBeginning
	internalBucket := tx.Bucket(bucketInternal)
	err = internalBucket.Put(bucketInternalKeyStats, encoding.Marshal(txdb.stats))
	if err != nil {
		return fmt.Errorf("failed to store transaction db (height=%d; changeID=%x) as a stat: %v",
			txdb.stats.BlockHeight, txdb.stats.ConsensusChangeID, err)
	}

	// store the genesis mint condition
	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	err = mintConditionsBucket.Put(internal.EncodeBlockheight(0), encoding.Marshal(genesisMintCondition))
	if err != nil {
		return fmt.Errorf("failed to store genesis mint condition: %v", err)
	}

	// all buckets created, and populated with initial content
	return nil
}

// ProcessConsensusChange implements modules.ConsensusSetSubscriber,
// calling txdb.processConsensusChange, so that the TransactionDB
// does not expose its interface implementation outside this package,
// given that we want the user to subscribe using the (*TransactionDB).SubscribeToConsensusSet method.
func (sub *transactionDBCSSubscriber) ProcessConsensusChange(css modules.ConsensusChange) {
	sub.txdb.processConsensusChange(css)
}

func (sub *transactionDBCSSubscriber) unsubscribe() {
	sub.cs.Unsubscribe(sub)
}

// processConsensusChange implements modules.ConsensusSetSubscriber,
// used to apply/revert transactions we care about in the internal persistent storage.
func (txdb *TransactionDB) processConsensusChange(css modules.ConsensusChange) {
	if err := txdb.tg.Add(); err != nil {
		// The TransactionDB should gracefully reject updates from the consensus set
		// that are sent after the wallet's Close method has closed the wallet's ThreadGroup.
		return
	}
	defer txdb.tg.Done()

	txdb.db.Update(func(tx *bolt.Tx) (err error) {
		// update reverted transactions in a block-defined order
		err = txdb.revertBlocks(tx, css.RevertedBlocks)
		if err != nil {
			return fmt.Errorf("failed to revert blocks: %v", err)
		}

		// update applied transactions in a block-defined order
		err = txdb.applyBlocks(tx, css.AppliedBlocks)
		if err != nil {
			return fmt.Errorf("failed to apply blocks: %v", err)
		}

		// update the consensus change ID and synced status
		txdb.stats.ConsensusChangeID, txdb.stats.Synced = css.ID, css.Synced

		// store stats
		internalBucket := tx.Bucket(bucketInternal)
		err = internalBucket.Put(bucketInternalKeyStats, encoding.Marshal(txdb.stats))
		if err != nil {
			return fmt.Errorf("failed to store transaction db (height=%d; changeID=%x; synced=%v) as a stat: %v",
				txdb.stats.BlockHeight, txdb.stats.ConsensusChangeID, txdb.stats.Synced, err)
		}

		return nil // all good
	})
}

// revert all the given blocks using the given writable bolt Transaction,
// meaning the block height will be decreased per reverted block and
// all reverted mint conditions will be deleted as well
func (txdb *TransactionDB) revertBlocks(tx *bolt.Tx, blocks []rivinetypes.Block) error {
	var (
		err error
		rtx *rivinetypes.Transaction
	)

	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	if mintConditionsBucket == nil {
		return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
	}

	// collect all one-per-block mint conditions
	for _, block := range blocks {
		for i := range block.Transactions {
			rtx = &block.Transactions[i]
			// check the version and handle the ones we care about
			switch rtx.Version {
			case rivinetypes.TransactionVersionOne:
				// ignore most common Tx
				continue
			case types.TransactionVersionBotRegistration:
				err = txdb.revertBotRegistrationTx(tx, rtx)
			case types.TransactionVersionMinterDefinition:
				err = txdb.revertMintConditionTx(tx, rtx)
			}
			if err != nil {
				return err
			}
		}

		// decrease block height (store later)
		txdb.stats.BlockHeight--
	}

	// all good
	return nil
}

// apply all the given blocks using the given writable bolt Transaction,
// meaning the block height will be increased per applied block and
// all applied mint conditions will be stored linked to their block height as well
//
// if a block contains multiple transactions with a mint condition,
// only the mint condition of the last transaction in the block's transaction list will be stored
func (txdb *TransactionDB) applyBlocks(tx *bolt.Tx, blocks []rivinetypes.Block) error {
	var (
		err error
		rtx *rivinetypes.Transaction
	)

	// collect all one-per-block mint conditions
	for _, block := range blocks {
		// increase block height (store later)
		txdb.stats.BlockHeight++

		for i := range block.Transactions {
			rtx = &block.Transactions[i]
			// check the version and handle the ones we care about
			switch rtx.Version {
			case rivinetypes.TransactionVersionOne:
				// ignore most common Tx
				continue
			case types.TransactionVersionBotRegistration:
				err = txdb.applyBotRegistrationTx(tx, block.Timestamp, rtx)
			case types.TransactionVersionMinterDefinition:
				err = txdb.applyMintConditionTx(tx, rtx)
			}
			if err != nil {
				return err
			}
		}
	}

	// all good
	return nil
}

func (txdb *TransactionDB) applyMintConditionTx(tx *bolt.Tx, rtx *rivinetypes.Transaction) error {
	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	if mintConditionsBucket == nil {
		return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
	}
	mdtx, err := types.MinterDefinitionTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the minter def. tx type: %v" + err.Error())
	}
	err = mintConditionsBucket.Put(internal.EncodeBlockheight(txdb.stats.BlockHeight), encoding.Marshal(mdtx.MintCondition))
	if err != nil {
		return fmt.Errorf(
			"failed to put mint condition for block height %d: %v",
			txdb.stats.BlockHeight, err)
	}
	return nil
}

func (txdb *TransactionDB) revertMintConditionTx(tx *bolt.Tx, rtx *rivinetypes.Transaction) error {
	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	if mintConditionsBucket == nil {
		return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
	}
	err := mintConditionsBucket.Delete(internal.EncodeBlockheight(txdb.stats.BlockHeight))
	if err != nil {
		return fmt.Errorf(
			"failed to delete mint condition for block height %d: %v",
			txdb.stats.BlockHeight, err)
	}
	return nil
}

func (txdb *TransactionDB) applyBotRegistrationTx(tx *bolt.Tx, blockTime rivinetypes.Timestamp, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	brtx, err := types.BotRegistrationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot registration tx type: %v", err)
	}
	// get the unique ID for the 3bot, using bolt's auto incrementing feature
	sequenceIndex, err := recordBucket.NextSequence()
	if err != nil {
		return fmt.Errorf("error while getting auto incrementing sequence bot ID: %v", err)
	}
	if sequenceIndex > types.MaxBotID {
		return errors.New("error while getting auto incrementing sequence bot ID: value exceeds 32 bit")
	}
	id := types.BotID(sequenceIndex)
	// store the record, and the other mappings, assuming the consensus validated that
	// the registration Tx is completely valid
	err = recordBucket.Put(tfencoding.Marshal(id), tfencoding.Marshal(types.BotRecord{
		ID:         id,
		Addresses:  brtx.Addresses,
		Names:      brtx.Names,
		PublicKey:  brtx.Identification.PublicKey,
		Expiration: types.SiaTimestampAsCompactTimestamp(blockTime) + types.CompactTimestamp(brtx.NrOfMonths)*types.BotMonth,
	}))
	if err != nil {
		return fmt.Errorf("error while storing record for bot %d with public key %v: %v", id, brtx.Identification.PublicKey, err)
	}
	// store pubkey to ID mapping
	err = applyKeyToIDMapping(tx, brtx.Identification.PublicKey, id)
	if err != nil {
		return fmt.Errorf("error while storing pubKey %s to bot id %d mapping: %v", brtx.Identification.PublicKey, id, err)
	}
	// store all name mappings
	for _, name := range brtx.Names {
		err = applyNameToIDMapping(tx, name, id)
		if err != nil {
			return fmt.Errorf("error while storing name %s to bot id %d mapping: %v", name.String(), id, err)
		}
	}
	// all information is applied
	return nil
}

func (txdb *TransactionDB) revertBotRegistrationTx(tx *bolt.Tx, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	brtx, err := types.BotRegistrationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot registration tx type: %v", err)
	}
	// get the ID for the given public key
	id, err := getBotIDForPublicKey(tx, brtx.Identification.PublicKey)
	if err != nil {
		return fmt.Errorf("error while fetching ID mapped to public key %v: %v",
			brtx.Identification.PublicKey, err)
	}
	// delete the record
	err = recordBucket.Delete(tfencoding.Marshal(id))
	if err != nil {
		return fmt.Errorf("error while deleting record for bot %d with public key %v: %v",
			id, brtx.Identification.PublicKey, err)
	}
	// delete the name->ID mappings
	for _, name := range brtx.Names {
		err = revertNameToIDMapping(tx, name)
		if err != nil {
			return fmt.Errorf("error while deleting name %s to bot id %d mapping: %v", name.String(), id, err)
		}
	}
	// delete the publicKey->ID mapping,
	// doing it last as this is the initial check that happens when registering a bot,
	// as to ensure we only have one bot per public key
	err = revertKeyToIDMapping(tx, brtx.Identification.PublicKey)
	if err != nil {
		return fmt.Errorf("error while deleting pubKey %s to bot id %d mapping: %v", brtx.Identification.PublicKey, id, err)
	}
	// all information is reverted
	return nil
}

// apply/revert the Key->ID mapping for a 3bot
func applyKeyToIDMapping(tx *bolt.Tx, key types.PublicKey, id types.BotID) error {
	mappingBucket := tx.Bucket(bucketBotKeyToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot key bucket does not exist")
	}
	return mappingBucket.Put(tfencoding.Marshal(key), tfencoding.Marshal(id))
}
func revertKeyToIDMapping(tx *bolt.Tx, key types.PublicKey) error {
	mappingBucket := tx.Bucket(bucketBotKeyToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot key bucket does not exist")
	}
	return mappingBucket.Delete(tfencoding.Marshal(key))
}

// apply/revert the Name->ID mapping for a 3bot
func applyNameToIDMapping(tx *bolt.Tx, name types.BotName, id types.BotID) error {
	mappingBucket := tx.Bucket(bucketBotNameToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot name bucket does not exist")
	}
	return mappingBucket.Put(tfencoding.Marshal(name), tfencoding.Marshal(id))
}
func revertNameToIDMapping(tx *bolt.Tx, name types.BotName) error {
	mappingBucket := tx.Bucket(bucketBotNameToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot name bucket does not exist")
	}
	return mappingBucket.Delete(tfencoding.Marshal(name))
}
