package erc20

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	tfeth "github.com/threefoldfoundation/tfchain/pkg/eth"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethstats"
	"github.com/ethereum/go-ethereum/les"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/params"
)

// LightClient creates a light client that can be used to interact with the Ethereum network,
// for ERC20 purposes. By default it is read-only, in order to also write to the network,
// you'll need load an account using the LoadAccount method.
type LightClient struct {
	*ethclient.Client // Client connection to the Ethereum chain
	stack             *node.Node
	lesc              *les.LightEthereum

	datadir string

	// optional account info
	accountLock sync.RWMutex
	account     *clientAccountInfo
}

type clientAccountInfo struct {
	keystore *keystore.KeyStore // Keystore containing the signing info
	account  accounts.Account   // Account funding the bridge requests
}

// LightClientConfig combines all configuration required for
// creating and configuring a LightClient.
type LightClientConfig struct {
	Port    int
	DataDir string

	BootstrapNodes []*discv5.Node
	NetworkName    string
	NetworkID      uint64
	GenesisBlock   *core.Genesis
}

func (lccfg *LightClientConfig) validate() error {
	if lccfg.Port == 0 {
		return errors.New("invalid LightClientConfig: no port defined")
	}
	if lccfg.DataDir == "" {
		return errors.New("invalid LightClientConfig: no data directory defined")
	}
	if len(lccfg.BootstrapNodes) == 0 {
		return errors.New("invalid LightClientConfig: no bootstrap nodes defined")
	}
	if lccfg.NetworkName == "" {
		return errors.New("invalid LightClientConfig: no network name defined")
	}
	if lccfg.NetworkID == 0 {
		return errors.New("invalid LightClientConfig: no network ID defined")
	}
	if lccfg.GenesisBlock == nil {
		return errors.New("invalid LightClientConfig: no genesis block defined")
	}
	return nil
}

// NewLightClient creates a new light client that can be used to interact with the ETH network.
// See `LightClient` for more information.
func NewLightClient(lccfg LightClientConfig, cancel <-chan struct{}) (*LightClient, error) {
	// validate the cfg, as to provide better error reporting for obvious errors
	err := lccfg.validate()
	if err != nil {
		return nil, err
	}

	// separate saved data per network
	datadir := filepath.Join(lccfg.DataDir, lccfg.NetworkName)

	// Assemble the raw devp2p protocol stack
	stack, err := node.New(&node.Config{
		Name:    "chain",
		Version: params.VersionWithMeta,
		DataDir: datadir,
		P2P: p2p.Config{
			NAT:              nat.Any(),
			NoDiscovery:      true,
			DiscoveryV5:      true,
			ListenAddr:       fmt.Sprintf(":%d", lccfg.Port),
			MaxPeers:         25,
			BootstrapNodesV5: lccfg.BootstrapNodes,
		},
	})
	if err != nil {
		return nil, err
	}

	// Assemble the Ethereum light client protocol
	var lesc *les.LightEthereum
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		cfg := eth.DefaultConfig
		cfg.Ethash.DatasetDir = filepath.Join(datadir, "ethash")
		cfg.SyncMode = downloader.LightSync
		cfg.NetworkId = lccfg.NetworkID
		cfg.Genesis = lccfg.GenesisBlock
		var err error
		lesc, err = les.New(ctx, &cfg)
		return lesc, err
	}); err != nil {
		return nil, err
	}

	stats := "" // Todo: should this stay in here?
	// Assemble the ethstats monitoring and reporting service'
	if stats != "" {
		if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			var serv *les.LightEthereum
			ctx.Service(&serv)
			return ethstats.New(stats, nil, serv)
		}); err != nil {
			return nil, err
		}
	}

	// Boot up the client and ensure it connects to bootnodes
	if err := stack.Start(); err != nil {
		return nil, err
	}

	// add all bootstrap peers
	for _, boot := range lccfg.BootstrapNodes {
		old, err := enode.ParseV4(boot.String())
		if err != nil {
			stack.Server().AddPeer(old)
		}
	}

	// Attach to the client and retrieve any interesting metadata
	api, err := stack.Attach()
	if err != nil {
		stack.Stop()
		return nil, err
	}

	// create a client for the stack
	client := ethclient.NewClient(api)

	// wait until (light) client is fully synced
	downloader := lesc.Downloader()
	for {
		progress := downloader.Progress()
		if progress.HighestBlock == 0 {
			log.Info(
				"LightClient's downloader needs to start to sync, waiting 10 seconds...",
				"current_block", progress.CurrentBlock)
		} else if downloader.Synchronising() {
			log.Info(
				"LightClient's downloader is still syncing, waiting 10 seconds...",
				"current_block", progress.CurrentBlock, "highest_block", progress.HighestBlock)
		} else {
			log.Info(
				"LightClient's downloader is synced",
				"current_block", progress.CurrentBlock, "highest_block", progress.HighestBlock, "starting_block", progress.StartingBlock)
			break
		}
		select {
		case <-time.After(time.Second * 10):
		case <-cancel:
			return nil, errors.New("failed to create light client, call got cancelled")
		}
	}

	// return created light client
	return &LightClient{
		Client:  client,
		stack:   stack,
		lesc:    lesc,
		datadir: datadir,
	}, nil
}

// Close terminates the Ethereum connection and tears down the stack.
func (lc *LightClient) Close() error {
	return lc.stack.Stop()
}

// FetchTransaction fetches a transaction from a remote peer using its block hash and tx index (within that block).
// Together with a found transactions it also returns the confirmations available for that Tx.
func (lc *LightClient) FetchTransaction(ctx context.Context, blockHash common.Hash, txHash common.Hash) (*types.Transaction, uint64, error) {
	block, err := lc.lesc.ApiBackend.GetBlock(ctx, blockHash)
	if err != nil {
		return nil, 0, err
	}
	chainHeight := lc.lesc.Downloader().Progress().HighestBlock
	blockHeight := block.Header().Number.Uint64()
	if blockHeight > chainHeight {
		return nil, 0, fmt.Errorf(
			"Tx %q is in block %d while the known chain height is only %d",
			txHash.String(), blockHeight, chainHeight)
	}
	tx := block.Transaction(txHash)
	if tx == nil {
		return nil, 0, errors.New("transaction could not be found")
	}
	confirmations := (chainHeight - blockHeight) + 1
	return tx, confirmations, nil
}

// LoadAccount loads an account into this light client,
// allowing writeable operations using the loaded account.
// An error is returned in case no account could be loaded.
func (lc *LightClient) LoadAccount(accountJSON, accountPass string) error {
	// create keystore
	ks, err := tfeth.InitializeKeystore(lc.datadir, accountJSON, accountPass)
	if err != nil {
		return err
	}
	lc.accountLock.Lock()
	lc.account = &clientAccountInfo{
		keystore: ks,
		account:  ks.Accounts()[0],
	}
	lc.accountLock.Unlock()
	return nil
}

var (
	// ErrNoAccountLoaded is an error returned for all Light Client methods
	// that require an account and for which no account is loaded.
	ErrNoAccountLoaded = errors.New("no account was loaded into the light client")
)

// AccountBalanceAt returns the balance for the account at the given block height.
func (lc *LightClient) AccountBalanceAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	lc.accountLock.RLock()
	defer lc.accountLock.RUnlock()
	if lc.account == nil {
		return nil, ErrNoAccountLoaded
	}
	return lc.Client.BalanceAt(ctx, lc.account.account.Address, blockNumber)
}

// SignTx signs a given traction with the loaded account, returning the signed transaction and no error on success.
func (lc *LightClient) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	lc.accountLock.RLock()
	defer lc.accountLock.RUnlock()
	if lc.account == nil {
		return nil, ErrNoAccountLoaded
	}
	return lc.account.keystore.SignTx(lc.account.account, tx, chainID)
}

// AccountAddress returns the address of the loaded account,
// returning an error only if no account was loaded.
func (lc *LightClient) AccountAddress() (common.Address, error) {
	lc.accountLock.RLock()
	defer lc.accountLock.RUnlock()
	var addr common.Address
	if lc.account == nil {
		return addr, ErrNoAccountLoaded
	}
	copy(addr[:], lc.account.account.Address[:])
	return addr, nil
}
