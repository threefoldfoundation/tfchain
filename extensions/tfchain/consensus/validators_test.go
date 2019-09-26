package consensus

import (
	"encoding/hex"
	"testing"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/types"
)

var (
	minerFeesOkJSONEncodedTransactions = []string{
		`{
	"version": 0,
	"data": {
		"coininputs": [{
			"parentid": "6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7",
			"unlocker": {
				"type": 1,
				"condition": {
					"publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614"
				},
				"fulfillment": {
					"signature": "c69eb703637fda3d398fc770fd2a77c3d346e0f3e0f07d6652852a8d94a9cbcac5b119726f2356277866323194be2fb5b62e9d8ff455f69dfa7899cd2a39a003"
				}
			}
		}],
		"coinoutputs": [{
			"value": "49999000000000",
			"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
		}],
		"minerfees": ["1000000000"]
	}
}`, `{
	"version": 1,
	"data": {
		"coininputs": [{
			"parentid": "6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614",
					"signature": "3d137d1d6ca8bb997156c2def5cc30012063f845e128cbcf186e11fe77c5ae603d8f200bc5b92e627414bdd32bac64c2b4dda2131e458b8ba41700a8bf598b03"
				}
			}
		}],
		"coinoutputs": [{
			"value": "49999000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
				}
			}
		}],
		"minerfees": ["1000000000"]
	}
}`,
	}

	minerFeesLowJSONEncodedTransactions = []string{
		`{
	"version": 0,
	"data": {
		"coininputs": [{
			"parentid": "6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7",
			"unlocker": {
				"type": 1,
				"condition": {
					"publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614"
				},
				"fulfillment": {
					"signature": "0be97865ee1c9b7e5b302e8cb5d1b56dbd9d5b479f510cf05edddea08b3d6b5960c3cce23fb8b32342ca68919633da3e051dadec809a3529b44dbf1819d60005"
				}
			}
		}],
		"coinoutputs": [{
			"value": "49999999999999",
			"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
		}],
		"minerfees": ["1"]
	}
}`, `{
	"version": 1,
	"data": {
		"coininputs": [{
			"parentid": "6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614",
					"signature": "bc5e00022886af25bd2af890ad31808f5af3762524b9c2b54f70b1522eb80fda0a374bfc04a14e64be879ea42754f4f157891ceec639487dbebb79ed2963830d"
				}
			}
		}],
		"coinoutputs": [{
			"value": "49999999999999",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
				}
			}
		}],
		"minerfees": ["1"]
	}
}`,
	}

	minerFeesMissingJSONEncodedTransactions = []string{
		`{
	"version": 0,
	"data": {
		"coininputs": [{
			"parentid": "6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7",
			"unlocker": {
				"type": 1,
				"condition": {
					"publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614"
				},
				"fulfillment": {
					"signature": "a29c10920d9c19d8a07f5f72473f0832d5c73e97f6d184bf031bce70b340d1099d77fcf7519f90946e0b1d30449889d9bd45da6d78d871e7ea0a173a13982604"
				}
			}
		}],
		"coinoutputs": [{
			"value": "50000000000000",
			"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
		}]
	}
}`, `{
	"version": 1,
	"data": {
		"coininputs": [{
			"parentid": "6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614",
					"signature": "075ec24cc2e3e64c18284759e346fd0c123e6aed84487142f1463b778fb7096b6715134fa0d40235217ba6aad13c109d6d435ddff64dfe3d2154143aa2227c03"
				}
			}
		}],
		"coinoutputs": [{
			"value": "50000000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
				}
			}
		}]
	}
}`,
	}
)

func TestMinimumFeeValidationForTransactions(t *testing.T) {
	constants := config.GetStandardnetGenesis()
	validationConstants := types.TransactionValidationContext{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	validators := GetStandardTransactionValidators()
	txMappedValidators := GetStandardTransactionVersionMappedValidators()
	testMinimumFeeValidationForTransactions(t, "standard", validationConstants, validators, txMappedValidators)

	constants = config.GetTestnetGenesis()
	validationConstants = types.TransactionValidationContext{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	validators = GetTestnetTransactionValidators()
	txMappedValidators = GetTestnetTransactionVersionMappedValidators()
	testMinimumFeeValidationForTransactions(t, "testnet", validationConstants, validators, txMappedValidators)

	constants = config.GetDevnetGenesis()
	validationConstants = types.TransactionValidationContext{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	validators = GetDevnetTransactionValidators()
	txMappedValidators = GetDevnetTransactionVersionMappedValidators()
	testMinimumFeeValidationForTransactions(t, "devnet", validationConstants, validators, txMappedValidators)
}

func testMinimumFeeValidationForTransactions(t *testing.T, name string, validationConstants types.TransactionValidationContext, validators []modules.TransactionValidationFunction, txMappedValidators map[types.TransactionVersion][]modules.TransactionValidationFunction) {
	testMinimumFeeOkValidationForTransactions(t, name, validationConstants, validators, txMappedValidators)
	testMinimumFeeLowValidationForTransactions(t, name, validationConstants, validators, txMappedValidators)
	testMinimumFeeMissingValidationForTransactions(t, name, validationConstants, validators, txMappedValidators)
}

func testMinimumFeeOkValidationForTransactions(t *testing.T, name string, validationConstants types.TransactionValidationContext, validators []modules.TransactionValidationFunction, txMappedValidators map[types.TransactionVersion][]modules.TransactionValidationFunction) {
	for idx, validJSONEncodedTransaction := range minerFeesOkJSONEncodedTransactions {
		var txn types.Transaction
		err := txn.UnmarshalJSON([]byte(validJSONEncodedTransaction))
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee is OK
		err = testTransaction(txn, types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 100000,
		}, validationConstants, validators, txMappedValidators)
		if txn.Version == 0 && name == "devnet" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
		// should be valid, as miner fee is OK
		err = testTransaction(txn, types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 0,
		}, validationConstants, validators, txMappedValidators)
		if txn.Version == 0 && name == "devnet" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
		// should be valid, as miner fee is OK
		err = testTransaction(txn, types.ValidationContext{
			Confirmed:   false,
			BlockHeight: 0,
		}, validationConstants, validators, txMappedValidators)
		if txn.Version == 0 && name == "devnet" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
		// should be valid, as miner fee is OK
		err = testTransaction(txn, types.ValidationContext{
			Confirmed:   false,
			BlockHeight: 100000,
		}, validationConstants, validators, txMappedValidators)
		if txn.Version == 0 && name == "devnet" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
	}
}

func testMinimumFeeLowValidationForTransactions(t *testing.T, name string, validationConstants types.TransactionValidationContext, validators []modules.TransactionValidationFunction, txMappedValidators map[types.TransactionVersion][]modules.TransactionValidationFunction) {
	for idx, validJSONEncodedTransaction := range minerFeesLowJSONEncodedTransactions {
		var txn types.Transaction
		err := txn.UnmarshalJSON([]byte(validJSONEncodedTransaction))
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee isn't OK, but block height is low enough
		// except for devnet, in that case it isn't OK
		err = testTransaction(txn, types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 0,
		}, validationConstants, validators, txMappedValidators)
		if name == "devnet" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
	}
}

func testMinimumFeeMissingValidationForTransactions(t *testing.T, name string, validationConstants types.TransactionValidationContext, validators []modules.TransactionValidationFunction, txMappedValidators map[types.TransactionVersion][]modules.TransactionValidationFunction) {
	for idx, validJSONEncodedTransaction := range minerFeesMissingJSONEncodedTransactions {
		var txn types.Transaction
		err := txn.UnmarshalJSON([]byte(validJSONEncodedTransaction))
		if err != nil {
			t.Fatal(name, idx, err)
		}
		// should be valid, as miner fee isn't OK, but block height is low enough
		// except for devnet, in that case it isn't OK
		err = testTransaction(txn, types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 0,
		}, validationConstants, validators, txMappedValidators)
		if name == "devnet" {
			if err == nil {
				t.Fatal(name, idx, "expected error, but none received")
			}
		} else {
			if err != nil {
				t.Fatal(name, idx, err)
			}
		}
	}
}

func testTransaction(txn types.Transaction, validationConstants types.ValidationContext, transactionValidationContext types.TransactionValidationContext, validators []modules.TransactionValidationFunction, txMappedValidators map[types.TransactionVersion][]modules.TransactionValidationFunction) error {
	for _, validator := range validators {
		err := validator(modules.ConsensusTransaction{
			Transaction: txn,
			BlockHeight: validationConstants.BlockHeight,
			BlockTime:   validationConstants.BlockTime,
			SequenceID:  1,
			SpentCoinOutputs: map[types.CoinOutputID]types.CoinOutput{
				hcoid(`6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7`): {
					Value:     types.NewCurrency64(50000).Mul64(1000000000),
					Condition: types.NewCondition(types.NewUnlockHashCondition(huh(`01f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609adbbcd3977363c`))),
				},
			},
		}, transactionValidationContext)
		if err != nil {
			return err
		}
	}
	if validators, ok := txMappedValidators[txn.Version]; ok {
		for _, validator := range validators {
			err := validator(modules.ConsensusTransaction{
				Transaction: txn,
				BlockHeight: validationConstants.BlockHeight,
				BlockTime:   validationConstants.BlockTime,
				SequenceID:  1,
				SpentCoinOutputs: map[types.CoinOutputID]types.CoinOutput{
					hcoid(`6b7e26eb0938e667a8b32a0781e7618f8590d05fd4ddf8b09f124bd20abe8ad7`): {
						Value:     types.NewCurrency64(50000).Mul64(1000000000),
						Condition: types.NewCondition(types.NewUnlockHashCondition(huh(`01f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609adbbcd3977363c`))),
					},
				},
			}, transactionValidationContext)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func hbs(str string) []byte { // hexStr -> byte slice
	bs, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return bs
}
func hcoid(str string) (hash types.CoinOutputID) { // hbs -> types.CoinOutputID
	copy(hash[:], hbs(str))
	return
}
func huh(str string) (uh types.UnlockHash) {
	err := uh.LoadString(str)
	if err != nil {
		panic(err)
	}
	return
}
