package persist

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path"

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

// internal bucket database keys used for the transactionDB
var (
	bucketInternal         = []byte("internal")
	bucketInternalKeyStats = []byte("stats") // stored as a single struct, see `transactionDBStats`

	// getBucketMintConditionPerHeightRangeKey is used to compute the keys
	// of the values in this bucket
	bucketMintConditions = []byte("mintconditions")
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
		k, b = cursor.Seek(encodeBlockheight(height))
		if len(k) == 0 {
			// could be that we're past the last key, let's try the last key first
			k, b = cursor.Last()
			if len(k) == 0 {
				return errors.New("corrupt transaction DB: no matching mint condition could be found")
			}
			return nil
		}
		foundHeight := decodeBlockheight(k)
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
			Version: "1.1.0",
		}
	)

	txdb.db, err = persist.OpenDatabase(dbMetadata, filename)
	if err != nil {
		return fmt.Errorf("error opening tfchain transaction database: %v", err)
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
			b = mintConditionsBucket.Get(encodeBlockheight(0))
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
	err = mintConditionsBucket.Put(encodeBlockheight(0), encoding.Marshal(genesisMintCondition))
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
func (txdb *TransactionDB) revertBlocks(tx *bolt.Tx, blocks []rivinetypes.Block) (err error) {
	var (
		rtx      *rivinetypes.Transaction
		txLength int
	)

	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	if mintConditionsBucket == nil {
		return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
	}

	// collect all one-per-block mint conditions
	for _, block := range blocks {
		txLength = len(block.Transactions)
		// reverse check, as we only care about
		// the last registered mint condition of a block
		for i := txLength - 1; i >= 0; i-- {
			rtx = &block.Transactions[i]
			if rtx.Version != types.TransactionVersionMinterDefinition {
				continue // for now this txdb only cares about minter definitions
			}
			err = mintConditionsBucket.Delete(encodeBlockheight(txdb.stats.BlockHeight))
			if err != nil {
				return fmt.Errorf(
					"failed to delete mint condition for block height %d: %v",
					txdb.stats.BlockHeight, err)
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
func (txdb *TransactionDB) applyBlocks(tx *bolt.Tx, blocks []rivinetypes.Block) (err error) {
	var (
		rtx      *rivinetypes.Transaction
		txLength int
	)

	mintConditionsBucket := tx.Bucket(bucketMintConditions)
	if mintConditionsBucket == nil {
		return errors.New("corrupt transaction DB: mint conditions bucket does not exist")
	}

	// collect all one-per-block mint conditions
	for _, block := range blocks {
		// increase block height (store later)
		txdb.stats.BlockHeight++

		txLength = len(block.Transactions)
		// reverse check, as we only care about
		// the last registered mint condition of a block
		for i := txLength - 1; i >= 0; i-- {
			rtx = &block.Transactions[i]
			if rtx.Version != types.TransactionVersionMinterDefinition {
				continue // for now this txdb only cares about minter definitions
			}
			mdtx, err := types.MinterDefinitionTransactionFromTransaction(*rtx)
			if err != nil {
				return fmt.Errorf("unexpected error while unpacking the minter def. tx type: %v" + err.Error())
			}
			err = mintConditionsBucket.Put(encodeBlockheight(txdb.stats.BlockHeight), encoding.Marshal(mdtx.MintCondition))
			if err != nil {
				return fmt.Errorf(
					"failed to put mint condition for block height %d: %v",
					txdb.stats.BlockHeight, err)
			}
			break // only the last occurance matters for us
		}
	}

	// all good
	return nil
}

// encodeBlockheight encodes the given blockheight as a sortable key
func encodeBlockheight(height rivinetypes.BlockHeight) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key[:], uint64(height))
	return key
}

// decodeBlockheight decodes the given sortable key as a blockheight
func decodeBlockheight(key []byte) rivinetypes.BlockHeight {
	return rivinetypes.BlockHeight(binary.BigEndian.Uint64(key))
}
