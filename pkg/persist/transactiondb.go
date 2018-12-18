package persist

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/threefoldfoundation/tfchain/pkg/persist/internal"
	"github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/threefoldtech/rivine/build"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/persist"
	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	rivinesync "github.com/threefoldtech/rivine/sync"
	rivinetypes "github.com/threefoldtech/rivine/types"

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

// TODO:
//   add an in-memory cache (layer), such that we do not constantly have to look up the same
//   names/records/identifiers/publickeys
//   (for mint condition I do not think it is required, as interaction with minting is minimal)

// internal bucket database keys used for the transactionDB
var (
	bucketInternal         = []byte("internal")
	bucketInternalKeyStats = []byte("stats") // stored as a single struct, see `transactionDBStats`

	// getBucketMintConditionPerHeightRangeKey is used to compute the keys
	// of the values in this bucket
	bucketMintConditions = []byte("mintconditions")

	// buckets for the 3bot feature
	bucketBotRecords               = []byte("botrecords")      // ID => name
	bucketBotKeyToIDMapping        = []byte("botkeys")         // Key => ID
	bucketBotNameToIDMapping       = []byte("botnames")        // Name => ID
	bucketBotRecordImplicitUpdates = []byte("botimplupdates")  // txID => implicitBotRecordUpdate
	bucketBotTransactions          = []byte("bottransactions") // ID => []txID

	// buckets for the ERC20-bridge feature
	bucketERC20ToTFTAddresses = []byte("addresses_erc20_to_tft") // erc20 => TFT
	bucketTFTToERC20Addresses = []byte("addresses_tft_to_erc20") // TFT => erc20
	bucketERC20TransactionIDs = []byte("erc20_transactionids")   // stores all unique ERC20 transaction ids used for erc20=>TFT exchanges
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
		ChainTime         rivinetypes.Timestamp
		Synced            bool
	}
)

var (
	// ensure TransactionDB implements the MintConditionGetter interface
	_ types.MintConditionGetter = (*TransactionDB)(nil)
	// ensure TransactionDB implements the BotRecordReadRegistry interface
	_ types.BotRecordReadRegistry = (*TransactionDB)(nil)
	// ensure TransactionDB implements the ERC20Registry interface
	_ types.ERC20Registry = (*TransactionDB)(nil)
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

// Retrieves the Last ConsensusChangeID stored.
func (txdb *TransactionDB) GetLastConsensusChangeID() modules.ConsensusChangeID {
	return txdb.stats.ConsensusChangeID
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
	err = siabin.Unmarshal(b, &mintCondition)
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
	err = siabin.Unmarshal(b, &mintCondition)
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

	b := recordBucket.Get(rivbin.Marshal(id))
	if len(b) == 0 {
		return nil, types.ErrBotNotFound
	}

	record := new(types.BotRecord)
	err := rivbin.Unmarshal(b, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// GetRecordForKey returns the record mapped to the given Key.
func (txdb *TransactionDB) GetRecordForKey(key rivinetypes.PublicKey) (record *types.BotRecord, err error) {
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

func getBotIDForPublicKey(tx *bolt.Tx, key rivinetypes.PublicKey) (types.BotID, error) {
	keyBucket := tx.Bucket(bucketBotKeyToIDMapping)
	if keyBucket == nil {
		return 0, errors.New("corrupt transaction DB: bot key bucket does not exist")
	}

	b := keyBucket.Get(rivbin.Marshal(key))
	if len(b) == 0 {
		return 0, types.ErrBotKeyNotFound
	}

	var id types.BotID
	err := rivbin.Unmarshal(b, &id)
	return id, err
}

// GetRecordForName returns the record mapped to the given Name.
func (txdb *TransactionDB) GetRecordForName(name types.BotName) (record *types.BotRecord, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) error {
		nameBucket := tx.Bucket(bucketBotNameToIDMapping)
		if nameBucket == nil {
			return errors.New("corrupt transaction DB: bot name bucket does not exist")
		}

		b := nameBucket.Get(rivbin.Marshal(name))
		if len(b) == 0 {
			return types.ErrBotNameNotFound
		}

		var id types.BotID
		err := rivbin.Unmarshal(b, &id)
		if err != nil {
			return err
		}

		record, err = getRecordForID(tx, id)
		if err != nil {
			return err
		}
		if record.Expiration.SiaTimestamp() <= txdb.stats.ChainTime {
			// a botname automatically expires as soon as the last 3bot that owned it expired as well
			return types.ErrBotNameExpired
		}
		return nil
	})
	return
}

// GetBotTransactionIdentifiers returns the identifiers of all transactions that created and updated the given bot's record.
//
// The transaction identifiers are returned in the (stable) order as defined by the blockchain.
func (txdb *TransactionDB) GetBotTransactionIdentifiers(id types.BotID) (ids []rivinetypes.TransactionID, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) (err error) {
		ids, err = getBotTransactions(tx, id)
		return
	})
	return
}

// GetERC20AddressForTFTAddress returns the mapped ERC20 address for the given TFT Address,
// iff the TFT Address has registered an ERC20 address explicitly.
func (txdb *TransactionDB) GetERC20AddressForTFTAddress(uh rivinetypes.UnlockHash) (addr types.ERC20Address, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) (err error) {
		addr, err = getERC20AddressForTFTAddress(tx, uh)
		return
	})
	return
}

// GetTFTAddressForERC20Address returns the mapped TFT address for the given ERC20 Address,
// iff the TFT Address has registered an ERC20 address explicitly.
func (txdb *TransactionDB) GetTFTAddressForERC20Address(addr types.ERC20Address) (uh rivinetypes.UnlockHash, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) (err error) {
		uh, err = getTFTAddressForERC20Address(tx, addr)
		return
	})
	return
}

// GetTFTTransactionIDForERC20TransactionID returns the mapped TFT TransactionID for the given ERC20 TransactionID,
// iff the ERC20 TransactionID has been used to fund an ERC20 CoinCreation Tx and has been registered as such, a nil TransactionID is returned otherwise.
func (txdb *TransactionDB) GetTFTTransactionIDForERC20TransactionID(id types.ERC20TransactionID) (txid rivinetypes.TransactionID, err error) {
	err = txdb.db.View(func(tx *bolt.Tx) (err error) {
		txid, err = getTfchainTransactionIDForERC20TransactionID(tx, id)
		return
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
			Version: "1.1.2.1",
		}
	)

	txdb.db, err = persist.OpenDatabase(dbMetadata, filename)
	if err != nil {
		if err != persist.ErrBadVersion {
			return fmt.Errorf("error opening tfchain transaction database: %v", err)
		}
		// try to migrate the DB
		err = txdb.migrateDB(filename)
		if err != nil {
			return err
		}
		// save the new metadata
		txdb.db.Metadata = dbMetadata
		err = txdb.db.SaveMetadata()
		if err != nil {
			return fmt.Errorf("error while saving the v1.1.2 metadata in the tfchain transaction database: %v", err)
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
			err = siabin.Unmarshal(b, &txdb.stats)
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
			err = siabin.Unmarshal(b, &storedMintCondition)
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

func (txdb *TransactionDB) migrateDB(filename string) error {
	// try to open the DB using the original version
	dbMetadata := persist.Metadata{
		Header:  "TFChain Transaction Database",
		Version: "1.1.0",
	}
	var err error
	txdb.db, err = persist.OpenDatabase(dbMetadata, filename)
	if err == nil {
		// migrate from a v1.1.0 DB
		return txdb.db.Update(txdb.migrateV110DB)
	}
	if err != persist.ErrBadVersion {
		return fmt.Errorf("error opening tfchain transaction v1.1.0 database: %v", err)
	}

	// try to open the initial v1.2.0 DB (never released, but already out in field for dev purposes)
	dbMetadata.Version = "1.2.0"
	txdb.db, err = persist.OpenDatabase(dbMetadata, filename)
	if err == nil {
		// migrate from a v1.2.0 DB
		return txdb.db.Update(txdb.migrateV120DB)
	}
	if err == persist.ErrBadVersion {
		return fmt.Errorf("error opening tfchain transaction database with unknown version: %v", err)
	}
	return fmt.Errorf("error opening tfchain transaction v1.2.0 database: %v", err)
}

func (txdb *TransactionDB) migrateV110DB(tx *bolt.Tx) error {
	// Enumerate and create the new database buckets.
	buckets := [][]byte{
		bucketBotRecords,
		bucketBotKeyToIDMapping,
		bucketBotNameToIDMapping,
		bucketBotRecordImplicitUpdates,
		bucketBotTransactions,
	}
	var err error
	for _, bucket := range buckets {
		_, err = tx.CreateBucket(bucket)
		if err != nil {
			return err
		}
	}
	// update the stats bucket
	var oldStats struct {
		ConsensusChangeID modules.ConsensusChangeID
		BlockHeight       rivinetypes.BlockHeight
		Synced            bool
	}
	internalBucket := tx.Bucket(bucketInternal)
	b := internalBucket.Get(bucketInternalKeyStats)
	if len(b) == 0 {
		return errors.New("structured stats value could not be found in existing transaction db")
	}
	err = siabin.Unmarshal(b, &oldStats)
	if err != nil {
		return fmt.Errorf("failed to unmarshal structured stats value from existing transaction db: %v", err)
	}
	err = internalBucket.Put(bucketInternalKeyStats, siabin.Marshal(transactionDBStats{
		ConsensusChangeID: oldStats.ConsensusChangeID,
		BlockHeight:       oldStats.BlockHeight,
		ChainTime:         0, // will fix itself on the first block it receives
		Synced:            oldStats.Synced,
	}))
	if err != nil {
		return err
	}

	// Continue the migration process towards the newest version
	return txdb.migrateV120DB(tx)
}

func (txdb *TransactionDB) migrateV120DB(tx *bolt.Tx) error {
	// Enumerate and create the new database buckets.
	buckets := [][]byte{
		bucketERC20ToTFTAddresses,
		bucketTFTToERC20Addresses,
		bucketERC20TransactionIDs,
	}
	var err error
	for _, bucket := range buckets {
		_, err = tx.CreateBucket(bucket)
		if err != nil {
			return err
		}
	}

	// migration process is finished
	return nil
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
		bucketBotRecordImplicitUpdates,
		bucketBotTransactions,
		bucketERC20ToTFTAddresses,
		bucketTFTToERC20Addresses,
		bucketERC20TransactionIDs,
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
	err = internalBucket.Put(bucketInternalKeyStats, siabin.Marshal(txdb.stats))
	if err != nil {
		return fmt.Errorf("failed to store transaction db (height=%d; changeID=%x) as a stat: %v",
			txdb.stats.BlockHeight, txdb.stats.ConsensusChangeID, err)
	}

	// store the genesis mint condition
	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	err = mintConditionsBucket.Put(internal.EncodeBlockheight(0), siabin.Marshal(genesisMintCondition))
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

	err := txdb.db.Update(func(tx *bolt.Tx) (err error) {
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
		err = internalBucket.Put(bucketInternalKeyStats, siabin.Marshal(txdb.stats))
		if err != nil {
			return fmt.Errorf("failed to store transaction db (height=%d; changeID=%x; synced=%v) as a stat: %v",
				txdb.stats.BlockHeight, txdb.stats.ConsensusChangeID, txdb.stats.Synced, err)
		}

		return nil // all good
	})
	if err != nil {
		build.Critical("transactionDB update failed:", err)
	}
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
			if rtx.Version == rivinetypes.TransactionVersionOne {
				continue // ignore most common Tx
			}
			ctx := transactionContext{
				BlockHeight:  txdb.stats.BlockHeight,
				BlockTime:    block.Timestamp,
				TxSequenceID: uint16(i),
			}
			// check the version and handle the ones we care about
			switch rtx.Version {
			case types.TransactionVersionBotRegistration:
				err = txdb.revertBotRegistrationTx(tx, ctx, rtx)
			case types.TransactionVersionBotRecordUpdate:
				err = txdb.revertRecordUpdateTx(tx, ctx, rtx)
			case types.TransactionVersionBotNameTransfer:
				err = txdb.revertBotNameTransferTx(tx, ctx, rtx)

			case types.TransactionVersionERC20CoinCreation:
				err = txdb.revertERC20CoinCreationTx(tx, ctx, rtx)
			case types.TransactionVersionERC20AddressRegistration:
				err = txdb.revertERC20AddressRegistrationTx(tx, ctx, rtx)

			case types.TransactionVersionMinterDefinition:
				err = txdb.revertMintConditionTx(tx, rtx)
			}
			if err != nil {
				return err
			}
		}

		// decrease block height (store later)
		txdb.stats.BlockHeight--
		// not super accurate, should be accurate enough and will fix itself when new blocks get applied
		txdb.stats.ChainTime = block.Timestamp
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
		txdb.stats.ChainTime = block.Timestamp

		for i := range block.Transactions {
			rtx = &block.Transactions[i]
			if rtx.Version == rivinetypes.TransactionVersionOne {
				continue // ignore most common Tx
			}
			ctx := transactionContext{
				BlockHeight:  txdb.stats.BlockHeight,
				BlockTime:    block.Timestamp,
				TxSequenceID: uint16(i),
			}
			// check the version and handle the ones we care about
			switch rtx.Version {
			case types.TransactionVersionBotRegistration:
				err = txdb.applyBotRegistrationTx(tx, ctx, rtx)
			case types.TransactionVersionBotRecordUpdate:
				err = txdb.applyRecordUpdateTx(tx, ctx, rtx)
			case types.TransactionVersionBotNameTransfer:
				err = txdb.applyBotNameTransferTx(tx, ctx, rtx)

			case types.TransactionVersionERC20CoinCreation:
				err = txdb.applyERC20CoinCreationTx(tx, ctx, rtx)
			case types.TransactionVersionERC20AddressRegistration:
				err = txdb.applyERC20AddressRegistrationTx(tx, ctx, rtx)

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
	err = mintConditionsBucket.Put(internal.EncodeBlockheight(txdb.stats.BlockHeight), siabin.Marshal(mdtx.MintCondition))
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

type transactionContext struct {
	BlockHeight  rivinetypes.BlockHeight
	BlockTime    rivinetypes.Timestamp
	TxSequenceID uint16
}

func (tctx transactionContext) TransactionShortID() sortableTransactionShortID {
	return newSortableTransactionShortID(tctx.BlockHeight, tctx.TxSequenceID)
}

func (txdb *TransactionDB) applyBotRegistrationTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
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
	// create the record
	record := types.BotRecord{
		ID:         id,
		PublicKey:  brtx.Identification.PublicKey,
		Expiration: types.SiaTimestampAsCompactTimestamp(ctx.BlockTime) + types.CompactTimestamp(brtx.NrOfMonths)*types.BotMonth,
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
	err = recordBucket.Put(rivbin.Marshal(id), rivbin.Marshal(record))
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
	// apply the transactionID to the list of transactionIDs for the given bot
	err = applyBotTransaction(tx, id, ctx.TransactionShortID(), rtx.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for bot %d: %v", id, err)
	}
	// all information is applied
	return nil
}

func (txdb *TransactionDB) revertBotRegistrationTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	brtx, err := types.BotRegistrationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot registration tx type: %v", err)
	}
	// the ID should be equal to the current bucket sequence, given it was incremented by the registration process
	rbSequence := recordBucket.Sequence()
	id := types.BotID(rbSequence)
	// delete the record
	err = recordBucket.Delete(rivbin.Marshal(id))
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
	// revert the transactionID from the list of transactionIDs for the given bot
	err = revertBotTransaction(tx, id, ctx.TransactionShortID())
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

func (txdb *TransactionDB) applyRecordUpdateTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	brutx, err := types.BotRecordUpdateTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot record update tx type: %v", err)
	}

	// get the bot record
	b := recordBucket.Get(rivbin.Marshal(brutx.Identifier))
	if len(b) == 0 {
		return errors.New("no bot record found for the specified identifier")
	}
	var record types.BotRecord
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found bot record: %v", err)
	}

	// check if bot is still active,
	// if the bot is still active we expect that the Tx defines the names to remove
	// otherwise we require that all names should be removed
	var namesInRecordRemovedImplicitly []types.BotName
	if record.IsExpired(ctx.BlockTime) {
		namesInRecordRemovedImplicitly = record.Names.Difference(types.BotNameSortedSet{}) // A \ {} = A
		// store the implicit update that will happen due to the invalid period prior to this Tx,
		// this will help is in reverting the record back to its original state,
		// as such implicit updates cannot be easily reverted otherwise
		err = applyImplicitBotRecordUpdate(tx, rtx.ID(), implicitBotRecordUpdate{
			PreviousExpirationTime: record.Expiration,
			InactiveNamesRemoved:   namesInRecordRemovedImplicitly,
		})
		if err != nil {
			return fmt.Errorf("failed to apply implicit record update: %v", err)
		}
	}

	// update it (will also reset names of an inactive bot)
	err = brutx.UpdateBotRecord(ctx.BlockTime, &record)
	if err != nil {
		return fmt.Errorf("failed to update bot record: %v", err)
	}

	// save it
	err = recordBucket.Put(rivbin.Marshal(brutx.Identifier), rivbin.Marshal(record))
	if err != nil {
		return fmt.Errorf("error while updating record for bot %d: %v", brutx.Identifier, err)
	}

	if len(namesInRecordRemovedImplicitly) > 0 {
		// otherwise remove all names that previously active,
		// as we can assume that an update of a record update HAS to make it active again
		for _, name := range namesInRecordRemovedImplicitly {
			err = revertNameToIDMappingIfOwnedByBot(tx, name, record.ID)
			if err != nil {
				return fmt.Errorf("failed to update bot record: error while tx-removing mapping of name %v: %v", name, err)
			}
		}
	} else {
		// if the bot was active, we apply the removals as defined by the Tx
		for _, name := range brutx.Names.Remove {
			err = revertNameToIDMapping(tx, name)
			if err != nil {
				return fmt.Errorf("failed to update bot record: error while record-removing mapping of name %v: %v", name, err)
			}
		}
	}

	// add mapping for all the added names
	for _, name := range brutx.Names.Add {
		err = applyNameToIDMapping(tx, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// apply the transactionID to the list of transactionIDs for the given bot
	err = applyBotTransaction(tx, record.ID, ctx.TransactionShortID(), rtx.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for bot %d: %v", record.ID, err)
	}

	// all information is applied
	return nil
}

func (txdb *TransactionDB) revertRecordUpdateTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	brutx, err := types.BotRecordUpdateTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot record update tx type: %v", err)
	}

	// get the bot record
	b := recordBucket.Get(rivbin.Marshal(brutx.Identifier))
	if len(b) == 0 {
		return errors.New("no bot record found for the specified identifier")
	}
	var record types.BotRecord
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
	if record.IsExpired(ctx.BlockTime) {
		txID := rtx.ID()

		// if the record is expired, there is a big chance that
		// there was an implicit update, as such let's try to get it,
		// if an implicit update did indeed take place, let's restore that info
		update, err := getImplicitBotRecordUpdate(tx, txID)
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
				err = applyNameToIDMappingIfAvailable(tx, name, record.ID)
				if err != nil {
					return fmt.Errorf("failed to revert bot record: :"+
						"failed to add back mapping of expired bot's name %v to its ID %d: %v", name, record.ID, err)
				}
			}

			// delete the implicit record update, it is no longer required
			err = revertImplicitBotRecordUpdate(tx, txID)
			if err != nil {
				return fmt.Errorf("failed to revert bot record: failed to revert implicit bot record update content:%v", err)
			}
		}
	}

	// save it
	err = recordBucket.Put(rivbin.Marshal(brutx.Identifier), rivbin.Marshal(record))
	if err != nil {
		return fmt.Errorf("error while updating record for bot %d: %v", brutx.Identifier, err)
	}

	// revert all names that were added
	for _, name := range brutx.Names.Add {
		err = revertNameToIDMapping(tx, name)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// apply all names again that were removed,
	// which can only be in case the bot was active
	for _, name := range brutx.Names.Remove {
		err = applyNameToIDMapping(tx, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to revert update bot record: error while revert mapping of name %v that was removed: %v", name, err)
		}
	}

	// revert the transactionID from the list of transactionIDs for the given bot
	err = revertBotTransaction(tx, record.ID, ctx.TransactionShortID())
	if err != nil {
		return fmt.Errorf("error while reverting transaction for bot %d: %v", record.ID, err)
	}

	// all information is applied
	return nil
}

func (txdb *TransactionDB) applyBotNameTransferTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	bnttx, err := types.BotNameTransferTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot name transfer tx type: %v", err)
	}

	// get the sender bot record
	b := recordBucket.Get(rivbin.Marshal(bnttx.Sender.Identifier))
	if len(b) == 0 {
		return errors.New("no bot record found for the specified sender bot identifier")
	}
	var record types.BotRecord
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found record of sender bot %d: %v", bnttx.Sender.Identifier, err)
	}
	// update sender bot (this also ensures the sender bot isn't expired)
	err = bnttx.UpdateSenderBotRecord(ctx.BlockTime, &record)
	if err != nil { // automatically checks also if at least one name is transferred, returning an error if not
		return fmt.Errorf("failed to update record of sender bot %d: %v", record.ID, err)
	}
	// save the record of the sender bot
	err = recordBucket.Put(rivbin.Marshal(record.ID), rivbin.Marshal(record))
	if err != nil {
		return fmt.Errorf("error while saving the updated record for sender bot %d: %v", record.ID, err)
	}

	// apply the transactionID to the list of transactionIDs for the sender bot
	err = applyBotTransaction(tx, record.ID, ctx.TransactionShortID(), rtx.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for sender bot %d: %v", record.ID, err)
	}

	// get the receiver bot record
	b = recordBucket.Get(rivbin.Marshal(bnttx.Receiver.Identifier))
	if len(b) == 0 {
		return errors.New("no bot record found for the specified receiver bot identifier")
	}
	err = rivbin.Unmarshal(b, &record)
	if err != nil {
		return fmt.Errorf("failed to unmarshal found record of receiver bot %d: %v", bnttx.Receiver.Identifier, err)
	}
	// update receiver bot (this also ensures the receiver bot isn't expired)
	err = bnttx.UpdateReceiverBotRecord(ctx.BlockTime, &record)
	if err != nil { // automatically checks also if at least one name is transferred, returning an error if not
		return fmt.Errorf("failed to update record of receiver bot %d: %v", record.ID, err)
	}
	// save the record of the receiver bot
	err = recordBucket.Put(rivbin.Marshal(record.ID), rivbin.Marshal(record))
	if err != nil {
		return fmt.Errorf("error while saving the updated record for receiver bot %d: %v", record.ID, err)
	}

	// update mapping for all the transferred names
	for _, name := range bnttx.Names {
		err = applyNameToIDMapping(tx, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while mapping name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// apply the transactionID to the list of transactionIDs for the receiver bot
	err = applyBotTransaction(tx, record.ID, ctx.TransactionShortID(), rtx.ID())
	if err != nil {
		return fmt.Errorf("error while applying transaction for receiver bot %d: %v", record.ID, err)
	}

	// update went fine
	return nil
}
func (txdb *TransactionDB) revertBotNameTransferTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	recordBucket := tx.Bucket(bucketBotRecords)
	if recordBucket == nil {
		return errors.New("corrupt transaction DB: bot record bucket does not exist")
	}
	bnttx, err := types.BotNameTransferTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the bot name transfer tx type: %v", err)
	}

	// get the receiver bot record
	b := recordBucket.Get(rivbin.Marshal(bnttx.Receiver.Identifier))
	if len(b) == 0 {
		return errors.New("no bot record found for the specified receiver bot identifier")
	}
	var record types.BotRecord
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
	err = recordBucket.Put(rivbin.Marshal(record.ID), rivbin.Marshal(record))
	if err != nil {
		return fmt.Errorf("error while saving the reverted (update of the) receiver bot %d: %v", record.ID, err)
	}

	// revert the transactionID from the list of transactionIDs for the receiver bot
	err = revertBotTransaction(tx, record.ID, ctx.TransactionShortID())
	if err != nil {
		return fmt.Errorf("error while reverting transaction for receiver bot %d: %v", record.ID, err)
	}

	// get the sender bot record
	b = recordBucket.Get(rivbin.Marshal(bnttx.Sender.Identifier))
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
	err = recordBucket.Put(rivbin.Marshal(record.ID), rivbin.Marshal(record))
	if err != nil {
		return fmt.Errorf("error while saving the reverted record for sender bot %d: %v", record.ID, err)
	}

	// update mapping for all the transferred names, reverting them back to the sender
	for _, name := range bnttx.Names {
		err = applyNameToIDMapping(tx, name, record.ID)
		if err != nil {
			return fmt.Errorf("failed to update bot record: error while mapping name %v to ID %v: %v", name, record.ID, err)
		}
	}

	// revert the transactionID from the list of transactionIDs for the sender bot
	err = revertBotTransaction(tx, record.ID, ctx.TransactionShortID())
	if err != nil {
		return fmt.Errorf("error while reverting transaction for sender bot %d: %v", record.ID, err)
	}

	// revert went fine
	return nil
}

func (txdb *TransactionDB) applyERC20AddressRegistrationTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	etartx, err := types.ERC20AddressRegistrationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Address Registration tx type: %v", err)
	}

	tftaddr := rivinetypes.NewPubKeyUnlockHash(etartx.PublicKey)
	erc20addr := types.ERC20AddressFromUnlockHash(tftaddr)

	return applyERC20AddressMapping(tx, tftaddr, erc20addr)
}

func (txdb *TransactionDB) revertERC20AddressRegistrationTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	etartx, err := types.ERC20AddressRegistrationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Address Registration tx type: %v", err)
	}

	tftaddr := rivinetypes.NewPubKeyUnlockHash(etartx.PublicKey)
	erc20addr := types.ERC20AddressFromUnlockHash(tftaddr)

	return revertERC20AddressMapping(tx, tftaddr, erc20addr)
}

func (txdb *TransactionDB) applyERC20CoinCreationTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	etcctx, err := types.ERC20CoinCreationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Coin Creation Tx type: %v", err)
	}
	return applyERC20TransactionID(tx, etcctx.TransactionID, rtx.ID())
}

func (txdb *TransactionDB) revertERC20CoinCreationTx(tx *bolt.Tx, ctx transactionContext, rtx *rivinetypes.Transaction) error {
	etcctx, err := types.ERC20CoinCreationTransactionFromTransaction(*rtx)
	if err != nil {
		return fmt.Errorf("unexpected error while unpacking the ERC20 Coin Creation Tx type: %v", err)
	}
	return revertERC20TransactionID(tx, etcctx.TransactionID)
}

// apply/revert the Key->ID mapping for a 3bot
func applyKeyToIDMapping(tx *bolt.Tx, key rivinetypes.PublicKey, id types.BotID) error {
	mappingBucket := tx.Bucket(bucketBotKeyToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot key bucket does not exist")
	}
	return mappingBucket.Put(rivbin.Marshal(key), rivbin.Marshal(id))
}
func revertKeyToIDMapping(tx *bolt.Tx, key rivinetypes.PublicKey) error {
	mappingBucket := tx.Bucket(bucketBotKeyToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot key bucket does not exist")
	}
	return mappingBucket.Delete(rivbin.Marshal(key))
}

// apply/revert the Name->ID mapping for a 3bot
func applyNameToIDMapping(tx *bolt.Tx, name types.BotName, id types.BotID) error {
	mappingBucket := tx.Bucket(bucketBotNameToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot name bucket does not exist")
	}
	return mappingBucket.Put(rivbin.Marshal(name), rivbin.Marshal(id))
}
func revertNameToIDMapping(tx *bolt.Tx, name types.BotName) error {
	mappingBucket := tx.Bucket(bucketBotNameToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot name bucket does not exist")
	}
	return mappingBucket.Delete(rivbin.Marshal(name))
}
func revertNameToIDMappingIfOwnedByBot(tx *bolt.Tx, name types.BotName, id types.BotID) error {
	mappingBucket := tx.Bucket(bucketBotNameToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot name bucket does not exist")
	}
	b := mappingBucket.Get(rivbin.Marshal(name))
	if len(b) == 0 {
		return nil // might be deleted by another bot, who took over ownership
	}
	var mappedID types.BotID
	err := rivbin.Unmarshal(b, &mappedID)
	if err != nil {
		return fmt.Errorf("corrupt BotID used as key in mapping of bot name %v", name)
	}
	if mappedID != id {
		return nil // ID no longer owned by this bot, ignore removal request in the mapping context
	}
	// delete name (mapping), as it was still owned by this bot
	return mappingBucket.Delete(rivbin.Marshal(name))
}
func applyNameToIDMappingIfAvailable(tx *bolt.Tx, name types.BotName, id types.BotID) error {
	mappingBucket := tx.Bucket(bucketBotNameToIDMapping)
	if mappingBucket == nil {
		return errors.New("corrupt transaction DB: bot name bucket does not exist")
	}
	b := mappingBucket.Get(rivbin.Marshal(name))
	if len(b) != 0 {
		return nil // already taken
	}
	return mappingBucket.Put(rivbin.Marshal(name), rivbin.Marshal(id))
}

// sortableTransactionShortID wraps around the rivinetypes.TransactionShortID,
// as to ensure it is encoded in a way that allows boltdb use it for natural ordering.
type sortableTransactionShortID rivinetypes.TransactionShortID

func newSortableTransactionShortID(height rivinetypes.BlockHeight, txSequenceID uint16) sortableTransactionShortID {
	return sortableTransactionShortID(rivinetypes.NewTransactionShortID(height, txSequenceID))
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

func applyBotTransaction(tx *bolt.Tx, id types.BotID, shortTxID sortableTransactionShortID, txID rivinetypes.TransactionID) error {
	txBucket := tx.Bucket(bucketBotTransactions)
	if txBucket == nil {
		return errors.New("corrupt transaction DB: implicit bot transactions bucket does not exist")
	}
	botBucket, err := txBucket.CreateBucketIfNotExists(rivbin.Marshal(id))
	if err != nil {
		return fmt.Errorf("corrupt transaction DB: failed to create/get bot %d inner bucket: %v", id, err)
	}
	return botBucket.Put(rivbin.Marshal(shortTxID), rivbin.Marshal(txID))
}
func revertBotTransaction(tx *bolt.Tx, id types.BotID, shortTxID sortableTransactionShortID) error {
	txBucket := tx.Bucket(bucketBotTransactions)
	if txBucket == nil {
		return errors.New("corrupt transaction DB: implicit bot transactions bucket does not exist")
	}
	botBucket := txBucket.Bucket(rivbin.Marshal(id))
	if botBucket != nil {
		return fmt.Errorf("corrupt transaction DB: bot %d inner bucket does not exist", id)
	}
	return botBucket.Delete(rivbin.Marshal(shortTxID))
}
func getBotTransactions(tx *bolt.Tx, id types.BotID) ([]rivinetypes.TransactionID, error) {
	txBucket := tx.Bucket(bucketBotTransactions)
	if txBucket == nil {
		return nil, errors.New("corrupt transaction DB: implicit bot transactions bucket does not exist")
	}
	botBucket := txBucket.Bucket(rivbin.Marshal(id))
	if botBucket == nil {
		return nil, nil // no transactions is acceptable
	}
	var txIDs []rivinetypes.TransactionID
	err := botBucket.ForEach(func(_, v []byte) (err error) {
		var txID rivinetypes.TransactionID
		err = rivbin.Unmarshal(v, &txID)
		txIDs = append(txIDs, txID)
		return
	})
	if err != nil {
		return nil, fmt.Errorf("corrupt transaction DB: error while parsing stored txID for bot %d: %v", id, err)
	}
	return txIDs, nil
}

func applyImplicitBotRecordUpdate(tx *bolt.Tx, txID rivinetypes.TransactionID, update implicitBotRecordUpdate) error {
	updateBucket := tx.Bucket(bucketBotRecordImplicitUpdates)
	if updateBucket == nil {
		return errors.New("corrupt transaction DB: implicit bot record update bucket does not exist")
	}
	return updateBucket.Put(rivbin.Marshal(txID), rivbin.Marshal(update))
}
func revertImplicitBotRecordUpdate(tx *bolt.Tx, txID rivinetypes.TransactionID) error {
	updateBucket := tx.Bucket(bucketBotRecordImplicitUpdates)
	if updateBucket == nil {
		return errors.New("corrupt transaction DB: implicit bot record update bucket does not exist")
	}
	return updateBucket.Delete(rivbin.Marshal(txID))
}
func getImplicitBotRecordUpdate(tx *bolt.Tx, txID rivinetypes.TransactionID) (implicitBotRecordUpdate, error) {
	var update implicitBotRecordUpdate

	updateBucket := tx.Bucket(bucketBotRecordImplicitUpdates)
	if updateBucket == nil {
		return update, errors.New("corrupt transaction DB: implicit bot record update bucket does not exist")
	}
	b := updateBucket.Get(rivbin.Marshal(txID))
	if len(b) == 0 {
		return update, nil
	}
	err := rivbin.Unmarshal(b, &update)
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
	PreviousExpirationTime types.CompactTimestamp
	InactiveNamesRemoved   []types.BotName
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
	if lenNamesRemoved > types.MaxNamesPerBot {
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
	if length > types.MaxNamesPerBot {
		return errors.New("too many bot names were implicitly removed")
	}
	ibru.InactiveNamesRemoved = make([]types.BotName, 0, length)
	for i := uint8(0); i < length; i++ {
		var name types.BotName
		err = name.UnmarshalRivine(r)
		if err != nil {
			return fmt.Errorf("implicitBotRecordUpdate: %v", err)
		}
		ibru.InactiveNamesRemoved = append(ibru.InactiveNamesRemoved, name)
	}

	// all read succesfully
	return nil
}

func applyERC20AddressMapping(tx *bolt.Tx, tftaddr rivinetypes.UnlockHash, erc20addr types.ERC20Address) error {
	btft, berc20 := rivbin.Marshal(tftaddr), rivbin.Marshal(erc20addr)

	// store ERC20->TFT mapping
	bucket := tx.Bucket(bucketERC20ToTFTAddresses)
	if bucket == nil {
		return errors.New("corrupt transaction DB: ERC20->TFT bucket does not exist")
	}
	err := bucket.Put(berc20, btft)
	if err != nil {
		return fmt.Errorf("error while storing ERC20->TFT address mapping: %v", err)
	}

	// store TFT->ERC20 mapping
	bucket = tx.Bucket(bucketTFTToERC20Addresses)
	if bucket == nil {
		return errors.New("corrupt transaction DB: TFT->ERC20 bucket does not exist")
	}
	err = bucket.Put(btft, berc20)
	if err != nil {
		return fmt.Errorf("error while storing TFT->ERC20 address mapping: %v", err)
	}

	// done
	return nil
}
func revertERC20AddressMapping(tx *bolt.Tx, tftaddr rivinetypes.UnlockHash, erc20addr types.ERC20Address) error {
	btft, berc20 := rivbin.Marshal(tftaddr), rivbin.Marshal(erc20addr)

	// delete ERC20->TFT mapping
	bucket := tx.Bucket(bucketERC20ToTFTAddresses)
	if bucket == nil {
		return errors.New("corrupt transaction DB: ERC20->TFT bucket does not exist")
	}
	err := bucket.Delete(berc20)
	if err != nil {
		return fmt.Errorf("error while deleting ERC20->TFT address mapping: %v", err)
	}

	// delete TFT->ERC20 mapping
	bucket = tx.Bucket(bucketTFTToERC20Addresses)
	if bucket == nil {
		return errors.New("corrupt transaction DB: TFT->ERC20 bucket does not exist")
	}
	err = bucket.Delete(btft)
	if err != nil {
		return fmt.Errorf("error while deleting TFT->ERC20 address mapping: %v", err)
	}

	// done
	return nil
}

func getERC20AddressForTFTAddress(tx *bolt.Tx, uh rivinetypes.UnlockHash) (types.ERC20Address, error) {
	bucket := tx.Bucket(bucketTFTToERC20Addresses)
	if bucket == nil {
		return types.ERC20Address{}, errors.New("corrupt transaction DB: TFT->ERC20 bucket does not exist")
	}
	b := bucket.Get(rivbin.Marshal(uh))
	if len(b) == 0 {
		return types.ERC20Address{}, nil
	}
	var addr types.ERC20Address
	err := rivbin.Unmarshal(b, &addr)
	if err != nil {
		return types.ERC20Address{}, fmt.Errorf("failed to fetch ERC20 Address for TFT address %v: %v", uh, err)
	}
	return addr, nil
}

func getTFTAddressForERC20Address(tx *bolt.Tx, addr types.ERC20Address) (rivinetypes.UnlockHash, error) {
	bucket := tx.Bucket(bucketERC20ToTFTAddresses)
	if bucket == nil {
		return rivinetypes.UnlockHash{}, errors.New("corrupt transaction DB: ERC20->TFT bucket does not exist")
	}
	b := bucket.Get(rivbin.Marshal(addr))
	if len(b) == 0 {
		return rivinetypes.UnlockHash{}, nil
	}
	var uh rivinetypes.UnlockHash
	err := rivbin.Unmarshal(b, &uh)
	if err != nil {
		return rivinetypes.UnlockHash{}, fmt.Errorf("failed to fetch TFT Address for ERC20 address %v: %v", addr, err)
	}
	return uh, nil
}

func applyERC20TransactionID(tx *bolt.Tx, erc20id types.ERC20TransactionID, tftid rivinetypes.TransactionID) error {
	bucket := tx.Bucket(bucketERC20TransactionIDs)
	if bucket == nil {
		return errors.New("corrupt transaction DB: ERC20 TransactionIDs bucket does not exist")
	}
	err := bucket.Put(rivbin.Marshal(erc20id), rivbin.Marshal(tftid))
	if err != nil {
		return fmt.Errorf("error while storing ERC20 TransactionID %v: %v", erc20id, err)
	}
	return nil
}
func revertERC20TransactionID(tx *bolt.Tx, id types.ERC20TransactionID) error {
	bucket := tx.Bucket(bucketERC20TransactionIDs)
	if bucket == nil {
		return errors.New("corrupt transaction DB: ERC20 TransactionIDs bucket does not exist")
	}
	err := bucket.Delete(rivbin.Marshal(id))
	if err != nil {
		return fmt.Errorf("error while deleting ERC20 TransactionID %v: %v", id, err)
	}
	return nil
}
func getTfchainTransactionIDForERC20TransactionID(tx *bolt.Tx, id types.ERC20TransactionID) (rivinetypes.TransactionID, error) {
	bucket := tx.Bucket(bucketERC20TransactionIDs)
	if bucket == nil {
		return rivinetypes.TransactionID{}, errors.New("corrupt transaction DB: ERC20 TransactionIDs bucket does not exist")
	}
	b := bucket.Get(rivbin.Marshal(id))
	if len(b) == 0 {
		return rivinetypes.TransactionID{}, nil
	}
	var txid rivinetypes.TransactionID
	err := rivbin.Unmarshal(b, &txid)
	if err != nil {
		return rivinetypes.TransactionID{}, fmt.Errorf("corrupt transaction DB: invalid tfchain TransactionID fetched for ERC20 TxID %v: %v", id, err)
	}
	return txid, nil
}
