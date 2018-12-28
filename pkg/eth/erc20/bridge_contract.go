package erc20

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	tfeth "github.com/threefoldfoundation/tfchain/pkg/eth"
	"github.com/threefoldfoundation/tfchain/pkg/eth/erc20/contract"
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

	lc *LightClient

	filter     *contract.TTFT20Filterer
	transactor *contract.TTFT20Transactor
	caller     *contract.TTFT20Caller

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

	bootstrapNodes, err := networkConfig.GetBootnodes()
	if err != nil {
		return nil, err
	}
	lc, err := NewLightClient(LightClientConfig{
		Port:           port,
		DataDir:        datadir,
		BootstrapNodes: bootstrapNodes,
		NetworkName:    networkConfig.NetworkName,
		NetworkID:      networkConfig.NetworkID,
		GenesisBlock:   networkConfig.GenesisBlock,
	})
	if err != nil {
		return nil, err
	}
	err = lc.LoadAccount(accountJSON, accountPass)
	if err != nil {
		return nil, err
	}

	filter, err := contract.NewTTFT20Filterer(networkConfig.ContractAddress, lc.Client)
	if err != nil {
		return nil, err
	}

	transactor, err := contract.NewTTFT20Transactor(networkConfig.ContractAddress, lc.Client)
	if err != nil {
		return nil, err
	}

	caller, err := contract.NewTTFT20Caller(networkConfig.ContractAddress, lc.Client)
	if err != nil {
		return nil, err
	}

	return &bridgeContract{
		lc:         lc,
		filter:     filter,
		transactor: transactor,
		caller:     caller,
	}, nil
}

// close terminates the Ethereum connection and tears down the stack.
func (bridge *bridgeContract) close() error {
	return bridge.lc.Close()
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
		if head, err = bridge.lc.HeaderByNumber(ctx, nil); err != nil {
			return err
		}
	}
	// Retrieve the balance, nonce and gas price from the current head
	var (
		nonce   uint64
		price   *big.Int
		balance *big.Int
	)
	if price, err = bridge.lc.SuggestGasPrice(ctx); err != nil {
		return err
	}
	if balance, err = bridge.lc.AccountBalanceAt(ctx, head.Number); err != nil {
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
	sub, err := bridge.lc.SubscribeNewHead(context.Background(), heads)
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
	accountAddress, err := bridge.lc.AccountAddress()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	opts := &bind.TransactOpts{
		Context: ctx, From: accountAddress,
		Signer: bridge.getSignerFunc(),
		Value:  nil, Nonce: nil, GasLimit: 0, GasPrice: nil,
	}
	_, err = bridge.transactor.Transfer(opts, recipient, amount)
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
	accountAddress, err := bridge.lc.AccountAddress()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	opts := &bind.TransactOpts{
		Context: ctx, From: accountAddress,
		Signer: bridge.getSignerFunc(),
		Value:  nil, Nonce: nil, GasLimit: 0, GasPrice: nil,
	}
	_, err = bridge.transactor.MintTokens(opts, receiver, amount, txID)
	if err != nil {
		return err
	}
	return nil
}

func (bridge *bridgeContract) registerWithdrawalAddress(address common.Address) error {
	log.Info("Calling register withdrawel address function in contract")
	accountAddress, err := bridge.lc.AccountAddress()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	opts := &bind.TransactOpts{
		Context: ctx, From: accountAddress,
		Signer: bridge.getSignerFunc(),
		Value:  nil, Nonce: nil, GasLimit: 0, GasPrice: nil,
	}
	_, err = bridge.transactor.RegisterWithdrawalAddress(opts, address)
	if err != nil {
		return err
	}
	return nil
}

func (bridge *bridgeContract) getSignerFunc() bind.SignerFn {
	return func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		accountAddress, err := bridge.lc.AccountAddress()
		if err != nil {
			return nil, err
		}
		if address != accountAddress {
			return nil, errors.New("not authorized to sign this account")
		}
		networkID := int64(bridge.networkConfig.NetworkID)
		return bridge.lc.SignTx(tx, big.NewInt(networkID))
	}
}