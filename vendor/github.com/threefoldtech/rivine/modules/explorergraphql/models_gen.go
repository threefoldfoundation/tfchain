// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package explorergraphql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/types"
)

// Contract represents a contract object (also called smart contracts) that can be looked up by a unique identifier,
// see UnlockHash for more information about the identifier type used for contracts.
// See the used types in this union for more information about the possible contracts.
type Contract interface {
	IsContract()
}

// Object represents an object that can be looked up by a unique identifier,
// see ObjectID for more information about the identifier type used for objects.
// See the used types in this union for more information about the possible objects.
type Object interface {
	IsObject()
}

// The different `Object` types that can contain an `Output`.
type OutputParent interface {
	IsOutputParent()
}

// An unlock condition is used to define "who"
// can spent a coin or block stake output.
//
// See the different `UnlockCondition` implementations
// to kow the different conditions that are possible.
type UnlockCondition interface {
	IsUnlockCondition()
}

type UnlockFulfillment interface {
	IsUnlockFulfillment()
}

// A wallet is identified by an `UnlockHash` and can be sent
// coins and block stakes to, as well as spent those
// coins and block stakes received. In practise
// it is nothing more than a storage of a private/public key pair
// a public key which can be converted to (and exposed as) an UnlockHash
// to look up its balance, in the case of a non-full client wallet.
//
// See the Wallet implementations for the different wallets that are possible.
type Wallet interface {
	IsWallet()
}

type AtomicSwapCondition struct {
	Version      ByteVersion              `json:"Version"`
	UnlockHash   types.UnlockHash         `json:"UnlockHash"`
	Sender       *UnlockHashPublicKeyPair `json:"Sender"`
	Receiver     *UnlockHashPublicKeyPair `json:"Receiver"`
	HashedSecret BinaryData               `json:"HashedSecret"`
	TimeLock     LockTime                 `json:"TimeLock"`
}

func (AtomicSwapCondition) IsUnlockCondition() {}

type AtomicSwapFulfillment struct {
	Version         ByteVersion     `json:"Version"`
	ParentCondition UnlockCondition `json:"ParentCondition"`
	PublicKey       types.PublicKey `json:"PublicKey"`
	Signature       Signature       `json:"Signature"`
	Secret          *BinaryData     `json:"Secret"`
}

func (AtomicSwapFulfillment) IsUnlockFulfillment() {}

// The balance contains aggregated asset values,
// and is updated for each block that affect's a wallet's
// coin or block stake balance.
type Balance struct {
	Unlocked BigInt `json:"Unlocked"`
	Locked   BigInt `json:"Locked"`
}

// The API of the chainshots facts collected for a block.
type BlockChainSnapshotFacts struct {
	TotalCoins                 *BigInt `json:"TotalCoins"`
	TotalLockedCoins           *BigInt `json:"TotalLockedCoins"`
	TotalBlockStakes           *BigInt `json:"TotalBlockStakes"`
	TotalLockedBlockStakes     *BigInt `json:"TotalLockedBlockStakes"`
	EstimatedActiveBlockStakes *BigInt `json:"EstimatedActiveBlockStakes"`
}

// The API of the facts collected for a block.
type BlockFacts struct {
	// The difficulty used, in the consensus algorithm,
	// for creating this block.
	Difficulty *BigInt `json:"Difficulty"`
	// The target hash used, in the consensus algorithm,
	// for creating this block.
	Target *crypto.Hash `json:"Target"`
	// The aggregated chain data as a snapshot taken,
	// after this fact's block was applied.
	ChainSnapshot *BlockChainSnapshotFacts `json:"ChainSnapshot"`
}

// The API for the block-specific "header" data of a block.
// Containing information such as the ID, the ID of its parent block,
// block time and height as well as (miner) payout information.
// The Parent and Child block can be queried recursively as well.
type BlockHeader struct {
	ID          crypto.Hash        `json:"ID"`
	ParentID    *crypto.Hash       `json:"ParentID"`
	Parent      *Block             `json:"Parent"`
	Child       *Block             `json:"Child"`
	BlockTime   *types.Timestamp   `json:"BlockTime"`
	BlockHeight *types.BlockHeight `json:"BlockHeight"`
	Payouts     []*BlockPayout     `json:"Payouts"`
}

// The aggregated chain data,
// updated for every block is that is applied and reverted.
type ChainAggregatedData struct {
	TotalCoins                 *BigInt `json:"TotalCoins"`
	TotalLockedCoins           *BigInt `json:"TotalLockedCoins"`
	TotalBlockStakes           *BigInt `json:"TotalBlockStakes"`
	TotalLockedBlockStakes     *BigInt `json:"TotalLockedBlockStakes"`
	EstimatedActiveBlockStakes *BigInt `json:"EstimatedActiveBlockStakes"`
}

// ChainConstants collect all constant information known about
// a chain network and exposed via this API.
type ChainConstants struct {
	// The name of the chain that this explorer is connected to.
	Name string `json:"Name"`
	// The name of the network,
	// usually one of `"standard"`, `"testnet"`, "`devnet"`.
	// The name of a network is not restricted to these values however.
	NetworkName string `json:"NetworkName"`
	// The unit of coins.
	// For the Threefold Chain this is for example `"TFT"`.
	CoinUnit string `json:"CoinUnit"`
	// The amount of decimals that the coins can be expressed in.
	// The coin values are always exposed as the lowered unit,
	// see the `BigInt` type for more information about the encoding
	// within the context of this API.
	//
	// If for example the CoinPrecision is `2`,
	// than a currency value of `"104"` is actually `1.04`.
	CoinPecision int `json:"CoinPecision"`
	// The source code version this daemon is compiled on.
	ChainVersion string `json:"ChainVersion"`
	// The transaction version that clients should use as the default
	// transaction version.
	DefaultTransactionVersion ByteVersion `json:"DefaultTransactionVersion"`
	// The gateway Protocol Version used by this daemon's gateway module.
	GatewayProtocolVersion string `json:"GatewayProtocolVersion"`
	// ConsensusPlugins provide you with the names of all plugins used by
	// the consensus of this network's daemons and thus allows you to know
	// what extra features might be available for this network.
	ConsensusPlugins []string `json:"ConsensusPlugins"`
	// The (Unix Epoch, seconds) timestamp of the first block,
	// the so called genesis block.
	GenesisTimestamp types.Timestamp `json:"GenesisTimestamp"`
	// Defines the maximum size a block is allowed to be, in bytes.
	BlockSizeLimitInBytes int `json:"BlockSizeLimitInBytes"`
	// The average block creation time in seconds the consensus algorithm
	// aims to achieve. It does not mean that it will take exatly this amount of seconds for a
	// new block to be created, nor it is an upper limit. You will however notice that the average
	// block creation time of a sufficient amount of sequential blocks does come close
	// to this average block creation time.
	AverageBlockCreationTimeInSeconds int `json:"AverageBlockCreationTimeInSeconds"`
	// The total amount of block stakes available at the creation of this blockchain.
	// As blockchains can currently not create new block stakes it is also the final amount of block stakes on the chain.
	GenesisTotalBlockStakes BigInt `json:"GenesisTotalBlockStakes"`
	// Defines how many blocks a block stake have to age
	// prior to being able to use block stakes for creating blocks.
	// The age is calculated by computing the height when the stakes were
	// transfered until the current block height.
	// When transfering stakes to yourself as part of a block creation,
	// the constant required aging concept (using the amount as defined here) does not apply.
	BlockStakeAging int `json:"BlockStakeAging"`
	// The fee that a block creator recieves for the creation of a block.
	// Can be null in case the chain does not award fees for such creations,
	// a possibility for private chains where all nodes are owned by one organisation.
	BlockCreatorFee *BigInt `json:"BlockCreatorFee"`
	// The minimum fee that has to be spent by a wallet in order to make a coin or block stake transfer.
	// The fee does not apply for block creation transactions.
	// Can be null in case the network does not require transaction fees.
	MinimumTransactionFee *BigInt `json:"MinimumTransactionFee"`
	// Some networks collect all transaction fees in a single wallet,
	// if this is the case it will be available as condition in this field, for query purposes.
	TransactionFeeBeneficiary UnlockCondition `json:"TransactionFeeBeneficiary"`
	// This delay, in block amount, defines how long a miner payout (e.g. block creator or transaction fee)
	// is locked prior to being spendable.
	PayoutMaturityDelay types.BlockHeight `json:"PayoutMaturityDelay"`
}

// ChainFacts collects facts about the queried chain.
type ChainFacts struct {
	// Constants collects all constant (static) data known about the chain,
	// and is provided by the daemon network configuration.
	Constants *ChainConstants `json:"Constants"`
	// LastBlock allows you to look up the last block,
	// saving you a second query, in case you need it, to look up that block,
	// even if it is just the height or timestamp.
	LastBlock *Block `json:"LastBlock"`
	// Contains all aggregated global data,
	// updated for this chain for every applied and reverted block.
	Aggregated *ChainAggregatedData `json:"Aggregated"`
}

type LockTimeCondition struct {
	Version    ByteVersion       `json:"Version"`
	UnlockHash *types.UnlockHash `json:"UnlockHash"`
	LockValue  LockTime          `json:"LockValue"`
	LockType   LockType          `json:"LockType"`
	Condition  UnlockCondition   `json:"Condition"`
}

func (LockTimeCondition) IsUnlockCondition() {}

type MultiSignatureCondition struct {
	Version                ByteVersion                `json:"Version"`
	UnlockHash             types.UnlockHash           `json:"UnlockHash"`
	Owners                 []*UnlockHashPublicKeyPair `json:"Owners"`
	RequiredSignatureCount int                        `json:"RequiredSignatureCount"`
}

func (MultiSignatureCondition) IsUnlockCondition() {}

type MultiSignatureFulfillment struct {
	Version         ByteVersion               `json:"Version"`
	ParentCondition UnlockCondition           `json:"ParentCondition"`
	Pairs           []*PublicKeySignaturePair `json:"Pairs"`
}

func (MultiSignatureFulfillment) IsUnlockFulfillment() {}

type NilCondition struct {
	Version    ByteVersion      `json:"Version"`
	UnlockHash types.UnlockHash `json:"UnlockHash"`
}

func (NilCondition) IsUnlockCondition() {}

type PublicKeySignaturePair struct {
	PublicKey types.PublicKey `json:"PublicKey"`
	Signature Signature       `json:"Signature"`
}

type SingleSignatureFulfillment struct {
	Version         ByteVersion     `json:"Version"`
	ParentCondition UnlockCondition `json:"ParentCondition"`
	PublicKey       types.PublicKey `json:"PublicKey"`
	Signature       Signature       `json:"Signature"`
}

func (SingleSignatureFulfillment) IsUnlockFulfillment() {}

type TransactionFeePayout struct {
	BlockPayout *BlockPayout `json:"BlockPayout"`
	Value       BigInt       `json:"Value"`
}

type UnlockHashCondition struct {
	Version    ByteVersion      `json:"Version"`
	UnlockHash types.UnlockHash `json:"UnlockHash"`
	PublicKey  *types.PublicKey `json:"PublicKey"`
}

func (UnlockHashCondition) IsUnlockCondition() {}

// Each `01` prefixed `UnlockHash` (wallet address) is linked to a `PublicKey`.
// If it is known, and thus exposed on the chain at some point,
// it will be stored, and can be queried using the `PublicKey` field.
// That field will be `null` in case the `PublicKey` is not (yet) known.
type UnlockHashPublicKeyPair struct {
	UnlockHash types.UnlockHash `json:"UnlockHash"`
	PublicKey  *types.PublicKey `json:"PublicKey"`
}

// The different types of Payouts one can find in a block (header).
type BlockPayoutType string

const (
	BlockPayoutTypeBlockReward    BlockPayoutType = "BLOCK_REWARD"
	BlockPayoutTypeTransactionFee BlockPayoutType = "TRANSACTION_FEE"
)

var AllBlockPayoutType = []BlockPayoutType{
	BlockPayoutTypeBlockReward,
	BlockPayoutTypeTransactionFee,
}

func (e BlockPayoutType) IsValid() bool {
	switch e {
	case BlockPayoutTypeBlockReward, BlockPayoutTypeTransactionFee:
		return true
	}
	return false
}

func (e BlockPayoutType) String() string {
	return string(e)
}

func (e *BlockPayoutType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BlockPayoutType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BlockPayoutType", str)
	}
	return nil
}

func (e BlockPayoutType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type LockType string

const (
	LockTypeBlockHeight LockType = "BLOCK_HEIGHT"
	LockTypeTimestamp   LockType = "TIMESTAMP"
)

var AllLockType = []LockType{
	LockTypeBlockHeight,
	LockTypeTimestamp,
}

func (e LockType) IsValid() bool {
	switch e {
	case LockTypeBlockHeight, LockTypeTimestamp:
		return true
	}
	return false
}

func (e LockType) String() string {
	return string(e)
}

func (e *LockType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = LockType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid LockType", str)
	}
	return nil
}

func (e LockType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// The different types of `Output` possible within the context of this API.
type OutputType string

const (
	OutputTypeCoin                OutputType = "COIN"
	OutputTypeBlockStake          OutputType = "BLOCK_STAKE"
	OutputTypeBlockCreationReward OutputType = "BLOCK_CREATION_REWARD"
	OutputTypeTransactionFee      OutputType = "TRANSACTION_FEE"
)

var AllOutputType = []OutputType{
	OutputTypeCoin,
	OutputTypeBlockStake,
	OutputTypeBlockCreationReward,
	OutputTypeTransactionFee,
}

func (e OutputType) IsValid() bool {
	switch e {
	case OutputTypeCoin, OutputTypeBlockStake, OutputTypeBlockCreationReward, OutputTypeTransactionFee:
		return true
	}
	return false
}

func (e OutputType) String() string {
	return string(e)
}

func (e *OutputType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OutputType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OutputType", str)
	}
	return nil
}

func (e OutputType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
