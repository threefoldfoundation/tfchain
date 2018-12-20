package bridge

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethstats"
	"github.com/ethereum/go-ethereum/les"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/params"

	tfeth "github.com/threefoldfoundation/tfchain/pkg/eth"

	"github.com/threefoldfoundation/tfchain/pkg/erc20/contract"
)

var (
	// OneToken is the exact value of one token
	OneToken = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

var (
	ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

// bridgeContract exposes a higher lvl api for specific contract bindings. In case of proxy contracts,
// the bridge needs to use the bindings of the implementation contract, but the address of the proxy.
type bridgeContract struct {
	networkConfig tfeth.NetworkConfiguration // Ethereum network
	stack         *node.Node                 // Ethereum protocol stack
	client        *ethclient.Client          // Client connection to the Ethereum chain
	keystore      *keystore.KeyStore         // Keystore containing the signing info
	account       accounts.Account           // Account funding the bridge requests

	filter     *contract.TTFT20Filterer
	transactor *contract.TTFT20Transactor

	// cache some stats in case they might be usefull
	head    *types.Header // Current head header of the bridge
	balance *big.Int      // The current balance of the bridge (note: ethers only!)
	nonce   uint64        // Current pending nonce of the bridge
	price   *big.Int      // Current gas price to issue funds with

	lock sync.RWMutex // Lock protecting the bridge's internals
}

func (bridge *bridgeContract) GetContractAdress() common.Address {
	return bridge.networkConfig.ContractAddress

}

func newBridgeContract(networkName string, port int, accountJSON, accountPass string, datadir string) (*bridgeContract, error) {
	// load correct network config
	networkConfig, err := tfeth.GetEthNetworkConfiguration(networkName)
	if err != nil {
		return nil, err
	}

	// load bootstrap enodes
	enodes, err := networkConfig.GetBootnodes()
	if err != nil {
		return nil, err
	}

	// separate saved data per network
	datadir = filepath.Join(datadir, networkName)

	// create keystore
	ks, err := tfeth.InitializeKeystore(datadir, accountJSON, accountPass)
	if err != nil {
		return nil, err
	}

	// Assemble the raw devp2p protocol stack
	stack, err := node.New(&node.Config{
		Name:    "chain",
		Version: params.VersionWithMeta,
		DataDir: datadir,
		P2P: p2p.Config{
			NAT:              nat.Any(),
			NoDiscovery:      true,
			DiscoveryV5:      true,
			ListenAddr:       fmt.Sprintf(":%d", port),
			MaxPeers:         25,
			BootstrapNodesV5: enodes,
		},
	})
	if err != nil {
		return nil, err
	}

	// Assemble the Ethereum light client protocol
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		cfg := eth.DefaultConfig
		cfg.Ethash.DatasetDir = filepath.Join(datadir, "ethash")
		cfg.SyncMode = downloader.LightSync
		cfg.NetworkId = networkConfig.NetworkID
		cfg.Genesis = networkConfig.GenesisBlock
		return les.New(ctx, &cfg)
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

	for _, boot := range enodes {
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

	filter, err := contract.NewTTFT20Filterer(networkConfig.ContractAddress, client)
	if err != nil {
		return nil, err
	}

	transactor, err := contract.NewTTFT20Transactor(networkConfig.ContractAddress, client)
	if err != nil {
		return nil, err
	}

	return &bridgeContract{
		networkConfig: networkConfig,
		stack:         stack,
		client:        client,
		keystore:      ks,
		account:       ks.Accounts()[0],

		filter:     filter,
		transactor: transactor,
	}, nil
}

// close terminates the Ethereum connection and tears down the stack.
func (bridge *bridgeContract) close() error {
	return bridge.stack.Stop()
}

// refresh attempts to retrieve the latest header from the chain and extract the
// associated bridge balance and nonce for connectivity caching.
func (bridge *bridgeContract) refresh(head *types.Header) error {
	// Ensure a state update does not run for too long
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// If no header was specified, use the current chain head
	var err error
	if head == nil {
		if head, err = bridge.client.HeaderByNumber(ctx, nil); err != nil {
			return err
		}
	}
	// Retrieve the balance, nonce and gas price from the current head
	var (
		nonce   uint64
		price   *big.Int
		balance *big.Int
	)
	if price, err = bridge.client.SuggestGasPrice(ctx); err != nil {
		return err
	}
	if balance, err = bridge.client.BalanceAt(ctx, bridge.account.Address, head.Number); err != nil {
		return err
	}
	// Everything succeeded, update the cached stats
	bridge.lock.Lock()
	bridge.head, bridge.balance = head, balance
	bridge.price, bridge.nonce = price, nonce
	bridge.lock.Unlock()
	return nil
}

func (bridge *bridgeContract) loop() {
	log.Info("Subscribing to eth headers")
	// channel to receive head updates from client on
	heads := make(chan *types.Header, 16)
	// subscribe to head upates
	sub, err := bridge.client.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		log.Error("Failed to subscribe to head events", "err", err)
	}
	defer sub.Unsubscribe()
	// channel so we can update the internal state from the heads
	update := make(chan *types.Header)
	go func() {
		for head := range update {
			// old heads should be ignored during a chain sync after some downtime
			if err := bridge.refresh(head); err != nil {
				log.Warn("Failed to update state", "block", head.Number, "err", err)
			}
			log.Debug("Internal stats updated", "block", head.Number, "account balance", bridge.balance, "gas price", bridge.price, "nonce", bridge.nonce)
		}
	}()
	for head := range heads {
		select {
		// only process new head if another isn't being processed yet
		case update <- head:
			log.Debug("Processing new head")
		default:
			log.Debug("Ignoring current head, update already in progress")
		}
	}
}

// SubscribeTransfers subscribes to new Transfer events on the given contract. This call blocks
// and prints out info about any transfer as it happened
func (bridge *bridgeContract) subscribeTransfers() error {
	sink := make(chan *contract.TTFT20Transfer)
	opts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	sub, err := bridge.filter.WatchTransfer(opts, sink, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err = <-sub.Err():
			return err
		case transfer := <-sink:
			log.Info("Noticed transfer event", "from", transfer.From, "to", transfer.To, "amount", transfer.Tokens)
		}
	}
}

// SubscribeMint subscribes to new Mint events on the given contract. This call blocks
// and prints out info about any mint as it happened
func (bridge *bridgeContract) subscribeMint() error {
	sink := make(chan *contract.TTFT20Mint)
	opts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	sub, err := bridge.filter.WatchMint(opts, sink, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err = <-sub.Err():
			return err
		case mint := <-sink:
			log.Info("Noticed mint event", "receiver", mint.Receiver, "amount", mint.Tokens, "TFT tx id", mint.Txid)
		}
	}
}

type withdrawEvent struct {
	sender   common.Address
	receiver common.Address
	amount   *big.Int
	txHash   common.Hash
}

// SubscribeWithdraw subscribes to new Withdraw events on the given contract. This call blocks
// and prints out info about any withdraw as it happened
func (bridge *bridgeContract) subscribeWithdraw(wc chan<- withdrawEvent) error {
	sink := make(chan *contract.TTFT20Withdraw)
	opts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	sub, err := bridge.filter.WatchWithdraw(opts, sink, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err = <-sub.Err():
			return err
		case withdraw := <-sink:
			log.Info("Noticed withdraw event", "receiver", withdraw.Receiver, "amount", withdraw.Tokens)
			wc <- withdrawEvent{sender: withdraw.From, receiver: withdraw.Receiver, amount: withdraw.Tokens, txHash: withdraw.Raw.TxHash}
		}
	}
}

// SubscribeRegisterWithdrawAddress subscribes to new RegisterWithdrawalAddress events on the given contract. This call blocks
// and prints out info about any RegisterWithdrawalAddress event as it happened
func (bridge *bridgeContract) subscribeRegisterWithdrawAddress() error {
	sink := make(chan *contract.TTFT20RegisterWithdrawalAddress)
	opts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	sub, err := bridge.filter.WatchRegisterWithdrawalAddress(opts, sink, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err = <-sub.Err():
			return err
		case withdraw := <-sink:
			log.Info("Noticed withadraw address registration event", "address", withdraw.Addr)
		}
	}
}

// TransferFunds transfers funds from one address to another
func (bridge *bridgeContract) transferFunds(recipient common.Address, amount *big.Int) error {
	if amount == nil {
		return errors.New("invalid amount")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	opts := &bind.TransactOpts{Context: ctx, From: bridge.account.Address, Signer: bridge.getSignerFunc(), Value: nil, Nonce: nil, GasLimit: 0, GasPrice: nil}
	_, err := bridge.transactor.Transfer(opts, recipient, amount)
	if err != nil {
		return err
	}
	return nil
}

//
func (bridge *bridgeContract) mint(receiver common.Address, amount *big.Int, txID string) error {
	log.Info("Calling mint function in contract")
	if amount == nil {
		return errors.New("invalid amount")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	opts := &bind.TransactOpts{Context: ctx, From: bridge.account.Address, Signer: bridge.getSignerFunc(), Value: nil, Nonce: nil, GasLimit: 0, GasPrice: nil}
	_, err := bridge.transactor.MintTokens(opts, receiver, amount, txID)
	if err != nil {
		return err
	}
	return nil
}

func (bridge *bridgeContract) registerWithdrawalAddress(address common.Address) error {
	log.Info("Calling register withdrawel address function in contract")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	opts := &bind.TransactOpts{Context: ctx, From: bridge.account.Address, Signer: bridge.getSignerFunc(), Value: nil, Nonce: nil, GasLimit: 0, GasPrice: nil}
	_, err := bridge.transactor.RegisterWithdrawalAddress(opts, address)
	if err != nil {
		return err
	}
	return nil
}

func (bridge *bridgeContract) getSignerFunc() bind.SignerFn {
	return func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != bridge.account.Address {
			return nil, errors.New("not authorized to sign this account")
		}
		networkID := int64(bridge.networkConfig.NetworkID)
		return bridge.keystore.SignTx(bridge.account, tx, big.NewInt(networkID))
	}
}
