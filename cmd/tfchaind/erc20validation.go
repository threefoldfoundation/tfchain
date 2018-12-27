package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/threefoldfoundation/tfchain/pkg/eth/erc20/contract"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/threefoldtech/rivine/types"
)

// ERC20NodeValidator implements the ERC20TransactionValidator,
// getting the transactions using the INFURA service.
//
// Goal is to move away from this ASAP and be able to fetch transactions using a Light client,
// problem is however that the LES/v2 protocol is not implemented yet by default in go-ethereum,
// it is however implemented server-side. So should we want, we can fork Ethereum,
// and contribute the client-side calls we care about as to be able to do it all from
// a light client.
type ERC20NodeValidator struct {
	ec  *ethclient.Client
	abi abi.ABI
}

// NewERC20NodeValidator creates a new INFURA-based ERC20NodeValidator.
// See the `ERC20NodeValidator` struct description for more information.
func NewERC20NodeValidator(network, apiKey string) (*ERC20NodeValidator, error) {
	// parse the ERC20 smart contract
	abi, err := abi.JSON(strings.NewReader(contract.TTFT20ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while parsing contract ABI: %v", err)
	}

	// create the ethereum client using the INFURA endpoint,
	// with the network and api-key filled into the relevant parts of the API URL.
	endpoint := fmt.Sprintf("https://%s.infura.io/v3/%s", network, apiKey)
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20NodeValidator: error while dialing the INFURA API: %v", err)
	}

	return &ERC20NodeValidator{
		ec:  client,
		abi: abi,
	}, nil
}

// ValidateWithdrawTx implements ERC20TransactionValidator.ValidateWithdrawTx
func (ev *ERC20NodeValidator) ValidateWithdrawTx(txID tftypes.ERC20TransactionID, expectedAddress tftypes.ERC20Address, expecedAmount types.Currency) error {
	// Get the transaction
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	tx, isPending, err := ev.ec.TransactionByHash(ctx, common.Hash(txID))
	if err != nil {
		return fmt.Errorf("failed to fetch ERC20 Tx: %v", err)
	}
	if isPending {
		return errors.New("ERC20 Tx is still pending")
	}

	// Extract the data
	// (and as a necessary step also validate if the Tx and input are of the correct type)
	txData := tx.Data()
	if len(txData) <= 4 {
		return fmt.Errorf("invalid ERC20 Tx: unexpected Tx data length: %v", len(txData))
	}

	// first 4 bytes contain the id, so let's get method using that ID
	method, err := ev.abi.MethodById(txData[:4])
	if err != nil {
		return fmt.Errorf("invalid ERC20 Tx: failed to get method using its parsed id: %v", err)
	}
	if method.Name != "transfer" {
		return fmt.Errorf("invalid ERC20 Tx: unexpected name for unpacked method ID: %s", method.Name)
	}
	// unpack the input into a struct we can work with
	params := struct {
		To     common.Address
		Tokens *big.Int
	}{}
	err = method.Inputs.Unpack(&params, txData[4:])
	if err != nil {
		return fmt.Errorf("error while unpacking transfer ERC20 Tx %v input: %v", txID, err)
	}

	// validate if the address is correct
	toAddress := params.To.String()
	expectedToAddress := common.Address(expectedAddress).String()
	if toAddress != expectedToAddress {
		return fmt.Errorf("unexpected to address %v: expected address %v", toAddress, expectedToAddress)
	}

	// validate if the amount of tokens withdrawn is correct
	amount := types.NewCurrency(params.Tokens)
	if amount.Equals(expecedAmount) {
		return fmt.Errorf("unexpected transferred TFT value %v: expected value %v",
			amount.String(), expecedAmount.String())
	}

	// all is good, return nil to indicate this
	return nil
}
