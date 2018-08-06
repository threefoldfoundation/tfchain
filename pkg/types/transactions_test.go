package types

import (
	"testing"

	"github.com/rivine/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/config"
)

var (
	validJSONEncodedTransactions = []string{
		`{
	"version": 1,
	"data": {
		"coininputs": [
			{
				"parentid": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "ed25519:def123def123def123def123def123def123def123def123def123def123def1",
						"signature": "ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"
					}
				}
			}
		],
		"minerfees": ["1000000000"],
		"arbitrarydata": "SGVsbG8sIFdvcmxkIQ=="
	}
}`,
		`{
	"version": 0,
	"data": {
		"coininputs": [
			{
				"parentid": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				"unlocker": {
					"type": 1,
					"condition": {
						"publickey": "ed25519:def123def123def123def123def123def123def123def123def123def123def1"
					},
					"fulfillment": {
						"signature": "ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"
					}
				}
			}
		],
		"minerfees": ["1000000000"],
		"arbitrarydata": "SGVsbG8sIFdvcmxkIQ=="
	}
}`,
	}
)

func TestMinimumFeeValidationForTransactions(t *testing.T) {
	defer func() {
		types.RegisterTransactionVersion(types.TransactionVersionZero, types.LegacyTransactionController{})
		types.RegisterTransactionVersion(types.TransactionVersionOne, types.DefaultTransactionController{})
	}()
	constants := config.GetStandardnetGenesis()
	validationConstants := types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForStandardNetwork()
	testMinimumFeeValidationForTransactions(t, "standard", validationConstants)
	constants = config.GetTestnetGenesis()
	validationConstants = types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForTestNetwork()
	testMinimumFeeValidationForTransactions(t, "test", validationConstants)
	constants = config.GetDevnetGenesis()
	validationConstants = types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForDevNetwork()
	testMinimumFeeValidationForTransactions(t, "dev", validationConstants)
}

func testMinimumFeeValidationForTransactions(t *testing.T, name string, validationConstants types.TransactionValidationConstants) {
	for idx, validJSONEncodedTransaction := range validJSONEncodedTransactions {
		var txn types.Transaction
		err := txn.UnmarshalJSON([]byte(validJSONEncodedTransaction))
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee is OK
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 100000,
		}, validationConstants)
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee is OK
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 0,
		}, validationConstants)
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee is OK
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   false,
			BlockHeight: 0,
		}, validationConstants)
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee is OK
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   false,
			BlockHeight: 100000,
		}, validationConstants)
		if err != nil {
			t.Fatal(name, idx, err)
		}
		txn.MinerFees[0] = types.NewCurrency64(1)
		// should be invalid, as miner fee isn't OK
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 100000,
		}, validationConstants)
		if err == nil {
			t.Fatal(name, idx, "expected error, but no error received")
		}
		// should be valid, as miner fee isn't OK, but block height is low enough
		// except for devnet, in that case it isn't OK
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 0,
		}, validationConstants)
		if name == "dev" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
		// should be invalid, as miner fee isn't OK, and in unconfirmed state the block height doesn't matter
		err = txn.ValidateTransaction(types.ValidationContext{
			Confirmed:   false,
			BlockHeight: 0,
		}, validationConstants)
		if err == nil {
			t.Fatal(name, idx, "expected error, but no error received")
		}

	}
}
