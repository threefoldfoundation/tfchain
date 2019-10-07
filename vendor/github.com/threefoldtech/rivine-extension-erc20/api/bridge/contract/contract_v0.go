package contract

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Code from old bindings being kept for legacy reasons

// TTFT20WithdrawV0 represents a Withdraw event raised by the TTFT20 contract.
//
// This type is kept to ensure we can still unpack withdraw events which happened before the
// withdraw event signature was updated
type TTFT20WithdrawV0 struct {
	From     common.Address
	Receiver common.Address
	Tokens   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}
