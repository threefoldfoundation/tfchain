package types

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
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
	RegisterTransactionTypesForStandardNetwork(nil, types.Currency{}, config.DaemonNetworkConfig{}) // no MintConditionGetter is required for this test
	testMinimumFeeValidationForTransactions(t, "standard", validationConstants)
	constants = config.GetTestnetGenesis()
	validationConstants = types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForTestNetwork(nil, types.Currency{}, config.DaemonNetworkConfig{}) // no MintConditionGetter is required for this test
	testMinimumFeeValidationForTransactions(t, "test", validationConstants)
	constants = config.GetDevnetGenesis()
	validationConstants = types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForDevNetwork(nil, types.Currency{}, config.DaemonNetworkConfig{}) // no MintConditionGetter is required for this test
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

// cctx -> txData -> cctx
func TestCoinCreationTransactionToAndFromTransactionData(t *testing.T) {
	for i, testCase := range testCoinCreationTransactions {
		txData := testCase.TransactionData()
		cctx, err := CoinCreationTransactionFromTransactionData(txData)
		if err != nil {
			t.Error(i, "failed to create cctx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// cctx -> tx -> cctx
func TestCoinCreationTransactionToAndFromTransaction(t *testing.T) {
	for i, testCase := range testCoinCreationTransactions {
		tx := testCase.Transaction()
		cctx, err := CoinCreationTransactionFromTransaction(tx)
		if err != nil {
			t.Error(i, "failed to create tx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// cctx -> JSON -> cctx
func TestCoinCreationTransactionToAndFromJSON(t *testing.T) {
	for i, testCase := range testCoinCreationTransactions {
		b, err := json.Marshal(testCase)
		if err != nil {
			t.Error(i, "failed to JSON-marshal", err)
			continue
		}
		if len(b) == 0 {
			t.Error(i, "JSON-marshal output is empty")
		}
		var cctx CoinCreationTransaction
		err = json.Unmarshal(b, &cctx)
		if err != nil {
			t.Error(i, "failed to JSON-unmarshal tx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// tx(cctx) -> JSON -> tx(cctx)
func TestCoinCreationTransactionAsTransactionToAndFromJSON(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	for i, testCase := range testCoinCreationTransactions {
		b, err := json.Marshal(testCase.Transaction())
		if err != nil {
			t.Error(i, "failed to JSON-marshal", err)
			continue
		}
		if len(b) == 0 {
			t.Error(i, "JSON-marshal output is empty")
		}
		var tx types.Transaction
		err = json.Unmarshal(b, &tx)
		if err != nil {
			t.Error(i, "failed to JSON-unmarshal tx", err)
			continue
		}
		cctx, err := CoinCreationTransactionFromTransaction(tx)
		if err != nil {
			t.Error(i, "failed to transform tx->cctx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// cctx -> Binary -> cctx
func TestCoinCreationTransactionToAndFromBinary(t *testing.T) {
	for i, testCase := range testCoinCreationTransactions {
		b := siabin.Marshal(testCase)
		if len(b) == 0 {
			t.Error(i, "Binary-marshal output is empty")
		}
		var cctx CoinCreationTransaction
		err := siabin.Unmarshal(b, &cctx)
		if err != nil {
			t.Error(i, "failed to Binary-unmarshal tx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// tx(cctx) -> Binary -> tx(cctx)
func TestCoinCreationTransactionAsTransactionToAndFromBinary(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	for i, testCase := range testCoinCreationTransactions {
		b := siabin.Marshal(testCase.Transaction())
		if len(b) == 0 {
			t.Error(i, "Binary-marshal output is empty")
		}
		var tx types.Transaction
		err := siabin.Unmarshal(b, &tx)
		if err != nil {
			t.Error(i, "failed to Binary-unmarshal tx", err)
			continue
		}
		cctx, err := CoinCreationTransactionFromTransaction(tx)
		if err != nil {
			t.Error(i, "failed to transform tx->cctx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// mdtx -> txData -> mdtx
func TestMinterDefinitionTransactionToAndFromTransactionData(t *testing.T) {
	for i, testCase := range testMinterDefinitionTransactions {
		txData := testCase.TransactionData()
		mdtx, err := MinterDefinitionTransactionFromTransactionData(txData)
		if err != nil {
			t.Error(i, "failed to create mdtx", err)
			continue
		}
		testCompareTwoMinterDefinitionTransactions(t, i, mdtx, testCase)
	}
}

// mdtx -> tx -> mdtx
func TestMinterDefinitionTransactionToAndFromTransaction(t *testing.T) {
	for i, testCase := range testMinterDefinitionTransactions {
		tx := testCase.Transaction()
		mdtx, err := MinterDefinitionTransactionFromTransaction(tx)
		if err != nil {
			t.Error(i, "failed to create mdtx", err)
			continue
		}
		testCompareTwoMinterDefinitionTransactions(t, i, mdtx, testCase)
	}
}

// mdtx -> JSON -> mdtx
func TestMinterDefinitionTransactionToAndFromJSON(t *testing.T) {
	for i, testCase := range testMinterDefinitionTransactions {
		b, err := json.Marshal(testCase)
		if err != nil {
			t.Error(i, "failed to JSON-marshal", err)
			continue
		}
		if len(b) == 0 {
			t.Error(i, "JSON-marshal output is empty")
		}
		var mdtx MinterDefinitionTransaction
		err = json.Unmarshal(b, &mdtx)
		if err != nil {
			t.Error(i, "failed to JSON-unmarshal tx", err)
			continue
		}
		testCompareTwoMinterDefinitionTransactions(t, i, mdtx, testCase)
	}
}

// tx(mdtx) -> JSON -> tx(mdtx)
func TestMinterDefinitionTransactionAsTransactionToAndFromJSON(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	for i, testCase := range testMinterDefinitionTransactions {
		b, err := json.Marshal(testCase.Transaction())
		if err != nil {
			t.Error(i, "failed to JSON-marshal", err)
			continue
		}
		if len(b) == 0 {
			t.Error(i, "JSON-marshal output is empty")
		}
		var tx types.Transaction
		err = json.Unmarshal(b, &tx)
		if err != nil {
			t.Error(i, "failed to JSON-unmarshal tx", err)
			continue
		}
		mdtx, err := MinterDefinitionTransactionFromTransaction(tx)
		if err != nil {
			t.Error(i, "failed to transform tx->mdtx", err)
			continue
		}
		testCompareTwoMinterDefinitionTransactions(t, i, mdtx, testCase)
	}
}

// mdtx -> Binary -> mdtx
func TestMinterDefinitionTransactionToAndFromBinary(t *testing.T) {
	for i, testCase := range testMinterDefinitionTransactions {
		b := siabin.Marshal(testCase)
		if len(b) == 0 {
			t.Error(i, "Binary-marshal output is empty")
		}
		var mdtx MinterDefinitionTransaction
		err := siabin.Unmarshal(b, &mdtx)
		if err != nil {
			t.Error(i, "failed to Binary-unmarshal tx", err)
			continue
		}
		testCompareTwoMinterDefinitionTransactions(t, i, mdtx, testCase)
	}
}

// tx(mdtx) -> Binary -> tx(mdtx)
func TestMinterDefinitionTransactionAsTransactionToAndFromBinary(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	for i, testCase := range testMinterDefinitionTransactions {
		b := siabin.Marshal(testCase.Transaction())
		if len(b) == 0 {
			t.Error(i, "Binary-marshal output is empty")
		}
		var tx types.Transaction
		err := siabin.Unmarshal(b, &tx)
		if err != nil {
			t.Error(i, "failed to Binary-unmarshal tx", err)
			continue
		}
		mdtx, err := MinterDefinitionTransactionFromTransaction(tx)
		if err != nil {
			t.Error(i, "failed to transform tx->mdtx", err)
			continue
		}
		testCompareTwoMinterDefinitionTransactions(t, i, mdtx, testCase)
	}
}

// tests if we can unmarshal a JSON example inspired by
// originally proposed structure given in spec at
// https://github.com/threefoldfoundation/tfchain/issues/155#issuecomment-408029100
func TestJSONUnmarshalSpecCoinCreationTransactionExample(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(`
{
	"version": 129,
	"data": {
		"nonce": "MTIzNDU2Nzg=",
		"mintfulfillment": {
			"type": 3,
			"data": {
				"pairs": [
					{
						"publickey": "ed25519:def123def123def123def123def123def123def123def123def123def123def1",
						"signature": "ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"
					},
					{
						"publickey": "ed25519:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
						"signature": "abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"
					}
				]
			}
		},
		"coinoutputs": [
			{
				"value": "10000221",
				"condition": {
					"type": 4,
					"data": {
						"unlockhashes": [
							"01e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70b1ccc65e2105",
							"01a6a6c5584b2bfbd08738996cd7930831f958b9a5ed1595525236e861c1a0dc353bdcf54be7d8"
						],
						"minimumsignaturecount": 2
					}
				}
			},
			{
				"value": "5000000000000000",
				"condition": {
					"type": 3,
					"data": {
						"locktime": 42,
						"condition": {
							"type": 1,
							"data": {
								"unlockhash": "01e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70b1ccc65e2105"
							}
						}
					}
				}
			}
		],
		"minerfees": [
			"3000000",
			"1230000000"
		],
		"arbitrarydata": "ZGF0YQ=="
	}
}
`))
	if err != nil {
		t.Fatalf("failed to JSON-unmarshal coin creation transaction (129): %v", err)
	}
	cctx, err := CoinCreationTransactionFromTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	testCompareTwoCoinCreationTransactions(t, 0, cctx, CoinCreationTransaction{
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
					},
					Signature: hbs("abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
				},
			},
		}),
		CoinOutputs: []types.CoinOutput{
			{
				Value: types.NewCurrency64(10000221),
				Condition: types.NewCondition(&types.MultiSignatureCondition{
					MinimumSignatureCount: 2,
					UnlockHashes: types.UnlockHashSlice{
						types.UnlockHash{
							Type: types.UnlockTypePubKey,
							Hash: hs("e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
						},
						types.UnlockHash{
							Type: types.UnlockTypePubKey,
							Hash: hs("a6a6c5584b2bfbd08738996cd7930831f958b9a5ed1595525236e861c1a0dc35"),
						},
					},
				}),
			},
			{
				Value: types.NewCurrency64(5000000000000000),
				Condition: types.NewCondition(&types.TimeLockCondition{
					LockTime: 42,
					Condition: &types.UnlockHashCondition{
						TargetUnlockHash: types.UnlockHash{
							Type: types.UnlockTypePubKey,
							Hash: hs("e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
						},
					},
				}),
			},
		},
		MinerFees: []types.Currency{
			types.NewCurrency64(3000000),
			types.NewCurrency64(1230000000),
		},
		ArbitraryData: types.ArbitraryData{Data: []byte("data")},
	})
}

// tests if we can unmarshal a JSON example inspired by
// originally proposed structure given in spec at
// https://github.com/threefoldfoundation/tfchain/issues/165#issue-349622350
func TestJSONUnmarshalSpecMinterDefinitionTransactionExample(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(`
{
	"version": 128,
	"data": {
		"nonce": "MTIzNDU2Nzg=",
		"mintfulfillment": {
			"type": 3,
			"data": {
				"pairs": [
					{
						"publickey": "ed25519:def123def123def123def123def123def123def123def123def123def123def1",
						"signature": "ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"
					},
					{
						"publickey": "ed25519:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
						"signature": "abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"
					}
				]
			}
		},
		"mintcondition": {
			"type": 4,
			"data": {
				"unlockhashes": [
					"01e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70b1ccc65e2105",
					"01a6a6c5584b2bfbd08738996cd7930831f958b9a5ed1595525236e861c1a0dc353bdcf54be7d8"
				],
				"minimumsignaturecount": 2
			}
		},
		"minerfees": [
			"3000000",
			"1230000000"
		],
		"arbitrarydata": "ZGF0YQ=="
	}
}
`))
	if err != nil {
		t.Fatalf("failed to JSON-unmarshal minter definition transaction (128): %v", err)
	}
	mdtx, err := MinterDefinitionTransactionFromTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	testCompareTwoMinterDefinitionTransactions(t, 0, mdtx, MinterDefinitionTransaction{
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
					},
					Signature: hbs("abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
				},
			},
		}),
		MintCondition: types.NewCondition(&types.MultiSignatureCondition{
			MinimumSignatureCount: 2,
			UnlockHashes: types.UnlockHashSlice{
				types.UnlockHash{
					Type: types.UnlockTypePubKey,
					Hash: hs("e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
				},
				types.UnlockHash{
					Type: types.UnlockTypePubKey,
					Hash: hs("a6a6c5584b2bfbd08738996cd7930831f958b9a5ed1595525236e861c1a0dc35"),
				},
			},
		}),
		MinerFees: []types.Currency{
			types.NewCurrency64(3000000),
			types.NewCurrency64(1230000000),
		},
		ArbitraryData: types.ArbitraryData{Data: []byte("data")},
	})
}

func testCompareTwoCoinCreationTransactions(t *testing.T, i int, a, b CoinCreationTransaction) {
	// compare mint fulfillment
	if !a.MintFulfillment.Equal(b.MintFulfillment) {
		t.Error(i, "mint fulfillment not equal")
	}
	// compare coin outputs
	if len(a.CoinOutputs) != len(b.CoinOutputs) {
		t.Error(i, "length coin outputs not equal")
	} else {
		for u, co := range a.CoinOutputs {
			if !co.Value.Equals(b.CoinOutputs[u].Value) {
				t.Error(i, u, "coin out value not equal",
					co.Value.String(), "!=", b.CoinOutputs[u].Value.String())
			}
			if !co.Condition.Equal(b.CoinOutputs[u].Condition) {
				t.Error(i, u, "coin out condition not equal")
			}
		}
	}
	// compare miner fees
	if len(a.MinerFees) != len(b.MinerFees) {
		t.Error(i, "length miner fees not equal")
	} else {
		for u, mf := range a.MinerFees {
			if !mf.Equals(b.MinerFees[u]) {
				t.Error(i, u, "miner fees not equal",
					mf.String(), "!=", b.MinerFees[u].String())
			}
		}
	}
	// compare arbitrary data
	if bytes.Compare(a.ArbitraryData.Data, b.ArbitraryData.Data) != 0 {
		t.Error(i, "arbitrary not equal",
			string(a.ArbitraryData.Data), "!=", string(b.ArbitraryData.Data))
	}
	if a.ArbitraryData.Type != b.ArbitraryData.Type {
		t.Error(i, "arbitrary data type not equal", a.ArbitraryData.Type, "!=", b.ArbitraryData.Type)
	}
}

var testCoinCreationTransactions = []CoinCreationTransaction{
	// most minimalistic Coin Creation Transaction
	{
		Nonce: RandomTransactionNonce(),
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
			},
			Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
		}),
		// A single coin output, with the lowest value possible,
		// locked by a NilCondition, meaning it is free for anyone to take
		CoinOutputs: []types.CoinOutput{
			{
				Value: types.NewCurrency64(1),
			},
		},
		// smallest tx fee
		MinerFees: []types.Currency{
			config.GetTestnetGenesis().MinimumTransactionFee,
		},
	},
	// a more complex Coin Creation Transaction
	{
		Nonce: RandomTransactionNonce(),
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
					},
					Signature: hbs("abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
				},
			},
		}),
		// Give money to a pool, future pool and the minter multisig wallet
		CoinOutputs: []types.CoinOutput{
			{
				Value: config.GetTestnetGenesis().CurrencyUnits.OneCoin.Mul64(50000000),
				Condition: types.NewCondition(&types.MultiSignatureCondition{
					MinimumSignatureCount: 2,
					UnlockHashes: types.UnlockHashSlice{
						types.UnlockHash{
							Type: types.UnlockTypePubKey,
							Hash: hs("e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
						},
						types.UnlockHash{
							Type: types.UnlockTypePubKey,
							Hash: hs("a6a6c5584b2bfbd08738996cd7930831f958b9a5ed1595525236e861c1a0dc35"),
						},
					},
				}),
			},
			{
				Value: config.GetTestnetGenesis().CurrencyUnits.OneCoin.Mul64(100000),
				Condition: types.NewCondition(types.NewUnlockHashCondition(types.UnlockHash{
					Type: types.UnlockTypePubKey,
					Hash: hs("b39baa9a58319fa47f78ed542a733a7198d106caeabf0a231b91ea3e4e222ffd"),
				})),
			},
			{
				Value: config.GetTestnetGenesis().CurrencyUnits.OneCoin.Mul64(20000000),
				Condition: types.NewCondition(types.NewTimeLockCondition(uint64(time.Now().AddDate(1, 6, 0).Unix()),
					types.NewUnlockHashCondition(types.UnlockHash{
						Type: types.UnlockTypePubKey,
						Hash: hs("def123def123def123def123def123def123def123def123def123def123def1"),
					}))),
			},
		},
		// smallest tx fee
		MinerFees: []types.Currency{config.GetTestnetGenesis().MinimumTransactionFee},
		// with a message
		ArbitraryData: types.ArbitraryData{Data: []byte("2300202+e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70")},
	},
}

func testCompareTwoMinterDefinitionTransactions(t *testing.T, i int, a, b MinterDefinitionTransaction) {
	// compare mint fulfillment
	if !a.MintFulfillment.Equal(b.MintFulfillment) {
		t.Error(i, "mint fulfillment not equal")
	}
	// compare mint condition
	if !a.MintCondition.Equal(b.MintCondition) {
		t.Error(i, "mint condition not equal")
	}
	// compare miner fees
	if len(a.MinerFees) != len(b.MinerFees) {
		t.Error(i, "length miner fees not equal")
	} else {
		for u, mf := range a.MinerFees {
			if !mf.Equals(b.MinerFees[u]) {
				t.Error(i, u, "miner fees not equal",
					mf.String(), "!=", b.MinerFees[u].String())
			}
		}
	}
	// compare arbitrary data
	if bytes.Compare(a.ArbitraryData.Data, b.ArbitraryData.Data) != 0 {
		t.Error(i, "arbitrary not equal",
			string(a.ArbitraryData.Data), "!=", string(b.ArbitraryData.Data))
	}
	if a.ArbitraryData.Type != b.ArbitraryData.Type {
		t.Error(i, "arbitrary data type not equal", a.ArbitraryData.Type, "!=", b.ArbitraryData.Type)
	}
}

var testMinterDefinitionTransactions = []MinterDefinitionTransaction{
	// most minimalistic Minter Definition Transaction
	{
		Nonce: RandomTransactionNonce(),
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
			},
			Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
		}),
		MintCondition: types.NewCondition(types.NewUnlockHashCondition(types.UnlockHash{
			Type: types.UnlockTypePubKey,
			Hash: hs("def123def123def123def123def123def123def123def123def123def123def1"),
		})),
		// smallest tx fee
		MinerFees: []types.Currency{
			config.GetTestnetGenesis().MinimumTransactionFee,
		},
	},
	// a more complex Minter Definition Transaction
	{
		Nonce: RandomTransactionNonce(),
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
					},
					Signature: hbs("abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
				},
			},
		}),
		MintCondition: types.NewCondition(types.NewTimeLockCondition(uint64(time.Now().AddDate(1, 6, 0).Unix()),
			&types.MultiSignatureCondition{
				MinimumSignatureCount: 2,
				UnlockHashes: types.UnlockHashSlice{
					types.UnlockHash{
						Type: types.UnlockTypePubKey,
						Hash: hs("e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
					},
					types.UnlockHash{
						Type: types.UnlockTypePubKey,
						Hash: hs("a6a6c5584b2bfbd08738996cd7930831f958b9a5ed1595525236e861c1a0dc35"),
					},
				},
			})),
		// smallest tx fee
		MinerFees: []types.Currency{config.GetTestnetGenesis().MinimumTransactionFee},
		// with a message
		ArbitraryData: types.ArbitraryData{Data: []byte("2300202+e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70")},
	},
}

// tests the patch which fixes https://github.com/threefoldfoundation/tfchain/issues/164
func TestCoinCreationTransactionIDUniqueness(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	for i, testCCTX := range testCoinCreationTransactions {
		a := testCCTX
		b := testCCTX

		// if two cctxs use the same nonce, outputs, miner fees and fulfillment,
		// the ID should be the same
		idA := a.Transaction().ID()
		idB := b.Transaction().ID()
		if bytes.Compare(idA[:], idB[:]) != 0 {
			t.Error(i,
				"expected the ID of two coin creation txs to be equal when same nonce is used, but:",
				idA.String(), " == ", idB.String(), " ; txA = ", a, " and tx B = ", b)
		}

		// if however at least the nonce is different, the ID will be different,
		// no matter the rest of the other fields
		b.Nonce = RandomTransactionNonce()
		idA = a.Transaction().ID()
		idB = b.Transaction().ID()
		if bytes.Compare(idA[:], idB[:]) == 0 {
			t.Error(i,
				"expected the ID of two coin creation txs to be different when a different nonce is used, but:",
				idA.String(), " == ", idB.String(), " ; txA = ", a, " and tx B = ", b)
		}
	}
}

// test inspired by TestCoinCreationTransactionIDUniqueness
func TestMinterDefinitionTransactionIDUniqueness(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	for i, testMDTX := range testMinterDefinitionTransactions {
		a := testMDTX
		b := testMDTX

		// if two cctxs use the same nonce, outputs, miner fees and fulfillment,
		// the ID should be the same
		idA := a.Transaction().ID()
		idB := b.Transaction().ID()
		if bytes.Compare(idA[:], idB[:]) != 0 {
			t.Error(i,
				"expected the ID of two minter definition txs to be equal when same nonce is used, but:",
				idA.String(), " == ", idB.String(), " ; txA = ", a, " and tx B = ", b)
		}

		// if however at least the nonce is different, the ID will be different,
		// no matter the rest of the other fields
		b.Nonce = RandomTransactionNonce()
		idA = a.Transaction().ID()
		idB = b.Transaction().ID()
		if bytes.Compare(idA[:], idB[:]) == 0 {
			t.Error(i,
				"expected the ID of two minter definition txs to be different when a different nonce is used, but:",
				idA.String(), " == ", idB.String(), " ; txA = ", a, " and tx B = ", b)
		}
	}
}

const validDevnetJSONEncodedMinterDefinitionTx = `{
	"version": 128,
	"data": {
		"nonce": "FoAiO8vN2eU=",
		"mintfulfillment": {
			"type": 1,
			"data": {
				"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
				"signature": "bdf023fbe7e0efec584d254b111655e1c2f81b9488943c3a712b91d9ad3a140cb0949a8868c5f72e08ccded337b79479114bdb4ed05f94dfddb359e1a6124602"
			}
		},
		"mintcondition": {
			"type": 1,
			"data": {
				"unlockhash": "01e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73e047fe6a0703"
			}
		},
		"minerfees": ["1000000000"],
		"arbitrarydata": "YSBtaW50ZXIgZGVmaW5pdGlvbiB0ZXN0"
	}
}`

func TestSignMinterDefinitionTransactionExtension(t *testing.T) {
	inMemoryMintConditionGetter := newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
		unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))))
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: inMemoryMintConditionGetter,
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(validDevnetJSONEncodedMinterDefinitionTx))
	if err != nil {
		t.Fatal("failed to decode valid minter definition tx:", err)
	}

	// util function to sign tx
	type testKeyPair struct {
		types.KeyPair
		SignCount int
	}
	signTxAndValidate := func(key interface{}) error {
		t.Helper()

		signCount := 1
		mintCondition, _ := inMemoryMintConditionGetter.GetActiveMintCondition()

		// validate condition first (free validateMintCondition test)
		err := validateMintCondition(mintCondition)
		if err != nil {
			return fmt.Errorf("invalid mint condition cannot be signed: %v", err)
		}

		// redefine fulfillment, as signing an already signed fulfillment is not possible
		mdtxExtension := tx.Extension.(*MinterDefinitionTransactionExtension)

		var condition types.UnlockCondition = mintCondition.Condition
	signSwitch:
		switch condition.ConditionType() {
		case types.ConditionTypeUnlockHash:
			tx.Extension = &MinterDefinitionTransactionExtension{
				Nonce: mdtxExtension.Nonce,
				MintFulfillment: types.NewFulfillment(types.NewSingleSignatureFulfillment(
					mdtxExtension.MintFulfillment.Fulfillment.(*types.SingleSignatureFulfillment).PublicKey,
				)),
				MintCondition: mdtxExtension.MintCondition,
			}

		case types.ConditionTypeMultiSignature:
			k, ok := key.(testKeyPair)
			if !ok {
				panic("sign key should be a testKeyPair")
			}
			signCount = k.SignCount
			key = k.KeyPair
			tx.Extension = &MinterDefinitionTransactionExtension{
				Nonce:           mdtxExtension.Nonce,
				MintFulfillment: types.NewFulfillment(types.NewMultiSignatureFulfillment(nil)),
				MintCondition:   mdtxExtension.MintCondition,
			}

		case types.ConditionTypeTimeLock:
			cg, ok := condition.(types.MarshalableUnlockConditionGetter)
			if !ok {
				panic(fmt.Errorf("unexpected Go-type for TimeLockCondition: %T", condition))
			}
			condition = cg.GetMarshalableUnlockCondition()
			goto signSwitch

		default:
			panic("unsupported condition type")
		}
		for i := 0; i < signCount; i++ {
			// sign fulfillment, which lives in the tx extension data
			err := tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, eo ...interface{}) error {
				return fulfillment.Sign(types.FulfillmentSignContext{
					ExtraObjects: eo,
					Transaction:  tx,
					Key:          key,
				})
			})
			if err != nil {
				return fmt.Errorf("sign #%d failed: %v", i, err)
			}
		}
		validationCtx := types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 4072,
			BlockTime:   1534271219,
		}
		chainConstants := config.GetDevnetGenesis()
		txValidationConstants := types.TransactionValidationConstants{
			BlockSizeLimit:         chainConstants.BlockSizeLimit,
			ArbitraryDataSizeLimit: chainConstants.ArbitraryDataSizeLimit,
			MinimumMinerFee:        chainConstants.MinimumTransactionFee,
		}
		return tx.ValidateTransaction(validationCtx, txValidationConstants)
	}

	// sign as the signer did on devnet, should succeed
	err = signTxAndValidate(hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"))
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	// sign as the signer did on devnet, should succeed
	err = signTxAndValidate(func() crypto.SecretKey { sk, _ := crypto.GenerateKeyPair(); return sk }())
	if err == nil {
		t.Fatalf("succeeded to sign, while it should fail")
	}

	// test that we can use a time lock condition
	inMemoryMintConditionGetter.applyMintCondition(0, types.NewCondition(types.NewTimeLockCondition(types.LockTimeMinTimestampValue, types.NewUnlockHashCondition(
		unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))))

	// sign as the signer did on devnet, should still succeed
	err = signTxAndValidate(hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"))
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	// overwrite coin tx multisig condition to be a multisig condition instead
	inMemoryMintConditionGetter.applyMintCondition(0, types.NewCondition(types.NewMultiSignatureCondition(
		types.UnlockHashSlice{
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("016438a548b6d377e87b08e8eae5ef641a4e70cc861b85b54b0921330e03084ffe0a8d9a38e3a8"),
		},
		2,
	)))

	// sign multisig condition, should fail as we didn't sign enough
	err = signTxAndValidate(testKeyPair{
		KeyPair: types.KeyPair{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			},
			PrivateKey: hbs("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		},
		SignCount: 1,
	})
	if err == nil {
		t.Fatal("should fail to validate as we didn't sign twice, but it succeeded")
	}

	// sign multisig condition, should succeed
	err = signTxAndValidate(testKeyPair{
		KeyPair: types.KeyPair{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			},
			PrivateKey: hbs("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		},
		SignCount: 2,
	})
	if err != nil {
		t.Fatalf("failed to sign multisig: %v", err)
	}

	// test that we can use a time lock condition with multisig
	// overwrite coin tx mint condition to be a timelocked unlock hash condition instead
	inMemoryMintConditionGetter.applyMintCondition(0, types.NewCondition(types.NewTimeLockCondition(types.LockTimeMinTimestampValue, types.NewMultiSignatureCondition(
		types.UnlockHashSlice{
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("016438a548b6d377e87b08e8eae5ef641a4e70cc861b85b54b0921330e03084ffe0a8d9a38e3a8"),
		},
		2,
	))))

	// sign multisig condition, should succeed
	err = signTxAndValidate(testKeyPair{
		KeyPair: types.KeyPair{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			},
			PrivateKey: hbs("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		},
		SignCount: 2,
	})
	if err != nil {
		t.Fatalf("failed to sign multisig: %v", err)
	}
}

func TestMinterDefinitionTransactionValidation(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(validDevnetJSONEncodedMinterDefinitionTx))
	if err != nil {
		t.Fatal("failed to decode valid minter definition tx:", err)
	}

	// util function to resign tx when needed,
	// as to ensure we fail because of the reason we want it to fail
	removeTxSignature := func() {
		mdtxExtension := tx.Extension.(*MinterDefinitionTransactionExtension)
		tx.Extension = &MinterDefinitionTransactionExtension{
			Nonce: mdtxExtension.Nonce,
			MintFulfillment: types.NewFulfillment(types.NewSingleSignatureFulfillment(
				mdtxExtension.MintFulfillment.Fulfillment.(*types.SingleSignatureFulfillment).PublicKey,
			)),
			MintCondition: mdtxExtension.MintCondition,
		}
	}
	resignTx := func(description string) {
		t.Helper()
		// redefine fulfillment, as signing an already signed fulfillment is not possible
		removeTxSignature()
		// sign fulfillment, which lives in the tx extension data
		err := tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, eo ...interface{}) error {
			return fulfillment.Sign(types.FulfillmentSignContext{
				ExtraObjects: eo,
				Transaction:  tx,
				Key:          hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			})
		})
		if err != nil {
			t.Fatalf("failed to resign after %q: %v", description, err)
		}
	}

	validationCtx := types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 4072,
		BlockTime:   1534271219,
	}
	chainConstants := config.GetDevnetGenesis()
	txValidationConstants := types.TransactionValidationConstants{
		BlockSizeLimit:         chainConstants.BlockSizeLimit,
		ArbitraryDataSizeLimit: chainConstants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        chainConstants.MinimumTransactionFee,
	}

	// should be valid, as it was published on a local devnet
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate minter definition tx, while it is supposed to be valid:", err)
	}

	origExtension := tx.Extension
	origMDExtension := origExtension.(*MinterDefinitionTransactionExtension)

	// make the arbitrary data too big, should fail
	origArbitraryData := tx.ArbitraryData
	tx.ArbitraryData.Data = make([]byte, chainConstants.ArbitraryDataSizeLimit+1)
	resignTx("changed arbitrary data")
	// should fail now
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of exceeded arbitrary data byte size limit")
	}
	// restore to valid arbitrary data
	tx.ArbitraryData = origArbitraryData
	// use orig extension as well, as we signed
	tx.Extension = origExtension

	// mess with nonce
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce: TransactionNonce{}, // nil-nonce is not allowed
		MintFulfillment: types.NewFulfillment(types.NewSingleSignatureFulfillment(
			origMDExtension.MintFulfillment.Fulfillment.(*types.SingleSignatureFulfillment).PublicKey,
		)),
		MintCondition: origMDExtension.MintCondition,
	}
	// sign fulfillment, which lives in the tx extension data
	err = tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, eo ...interface{}) error {
		return fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: eo,
			Transaction:  tx,
			Key:          hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		})
	})
	if err != nil {
		t.Fatalf("failed to resign after modifying Nonce to use NilTransactionNonce: %v", err)
	}
	// validate, should fail
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of nil TransactionNonce")
	}
	// use orig extension as well, as we modified and re-signed
	tx.Extension = origExtension

	// mess with the extension to make it fail
	// should all fail
	tx.Extension = nil
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of nil extension")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce:           RandomTransactionNonce(),
		MintFulfillment: origMDExtension.MintFulfillment,
		MintCondition:   origMDExtension.MintCondition,
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of incompatible nonce with signature")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce:         origMDExtension.Nonce,
		MintCondition: origMDExtension.MintCondition,
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of nil fulfillment")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce: origMDExtension.Nonce,
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
					},
					Signature: hbs("abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
				},
			},
		}),
		MintCondition: origMDExtension.MintCondition,
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of wrong fulfillment")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce: origMDExtension.Nonce,
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d"),
			},
			Signature: hbs("3e2ed4e893f66ffd57e26afe83d570ca4b8ba873f8236a60c018cde4852de1027256d088b2253ec061ae973f961f26cde8fa42f5d3c0ce1316560ceb25786f03"),
		}),
		MintCondition: origMDExtension.MintCondition,
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of wrong fulfillment")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce:           origMDExtension.Nonce,
		MintFulfillment: origMDExtension.MintFulfillment,
	}
	resignTx("changed mint condition to use nil-condition")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of nil condition")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce:           origMDExtension.Nonce,
		MintFulfillment: origMDExtension.MintFulfillment,
		MintCondition: types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("0282bbab17110c5e3556a9ce8ef9b243cdacde2c92d2f13283501d84b920bf48fc630b7cbab96d"))),
	}
	resignTx("changed mint condition to use atomic-swap-unlock-hash condition")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of non-supported unlock hash condition")
	}
	tx.Extension = &MinterDefinitionTransactionExtension{
		Nonce:           origMDExtension.Nonce,
		MintFulfillment: origMDExtension.MintFulfillment,
		MintCondition: types.NewCondition(&types.AtomicSwapCondition{
			Sender: types.UnlockHash{
				Type: types.UnlockTypePubKey,
				Hash: hs("1234567891234567891234567891234567891234567891234567891234567891"),
			},
			Receiver: types.UnlockHash{
				Type: types.UnlockTypePubKey,
				Hash: hs("6363636363636363636363636363636363636363636363636363636363636363"),
			},
			HashedSecret: types.AtomicSwapHashedSecret(hs("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")),
			TimeLock:     1522068743,
		}),
	}
	resignTx("changed mint condition to use atomic-swap condition")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of non-supported atomic swap condition")
	}
	// restore to valid extension
	tx.Extension = origExtension

	// at least one miner fee is given,
	// and each miner fee has to be at least the minimum defined miner fee amount
	origMinerFees := tx.MinerFees
	tx.MinerFees = nil
	removeTxSignature()
	// should all fail now
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of no miner fees defined")
	}
	tx.MinerFees = []types.Currency{types.ZeroCurrency}
	resignTx("use zero miner fee")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of a too low miner fee defined")
	}
	tx.MinerFees = []types.Currency{chainConstants.MinimumTransactionFee, types.ZeroCurrency}
	resignTx("use minimum and zero miner fees")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of a too low miner fee defined")
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee.Sub(types.NewCurrency64(1)),
		types.NewCurrency64(1),
	}
	resignTx("total miner fees is good enough")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of a too low miner fee defined")
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee,
	}
	resignTx("use minimum miner fee")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate minter definition tx, while it is supposed to be valid:", err)
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee,
		chainConstants.MinimumTransactionFee,
	}
	resignTx("use more minir fees than required")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate minter definition tx, while it is supposed to be valid:", err)
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee.Mul64(1000),
		chainConstants.MinimumTransactionFee,
	}
	resignTx("use more minir fees than required")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate minter definition tx, while it is supposed to be valid:", err)
	}
	// these should all work though
	// restore to valid miner fees
	tx.MinerFees = origMinerFees
	// replace extension as well with orig, as we resigned
	tx.Extension = origExtension

	// restored as it is, it should be valid
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate minter definition tx, while it is supposed to be valid:", err)
	}
}

func TestMinterDefinitionTransactionValidationWithUnknownMintCondition(t *testing.T) {
	mintConditionGetter := newInMemoryMintConditionGetter()
	mintConditionGetter.applyMintCondition(1, types.NewCondition(types.NewUnlockHashCondition(
		unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))))

	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionMinterDefinition, MinterDefinitionTransactionController{
		MintConditionGetter: mintConditionGetter,
	})
	defer types.RegisterTransactionVersion(TransactionVersionMinterDefinition, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(validDevnetJSONEncodedMinterDefinitionTx))
	if err != nil {
		t.Fatal("failed to decode valid minter definition tx:", err)
	}

	validationCtx := types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 0,
		BlockTime:   1534271219,
	}
	chainConstants := config.GetDevnetGenesis()
	txValidationConstants := types.TransactionValidationConstants{
		BlockSizeLimit:         chainConstants.BlockSizeLimit,
		ArbitraryDataSizeLimit: chainConstants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        chainConstants.MinimumTransactionFee,
	}

	// should be invalid, as no mint condition exists for that height (0)
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition transaction, " +
			"while it was expected to fail due to an unknown mint condition")
	}
}

const validDevnetJSONEncodedCoinCreationTx = `{
	"version": 129,
	"data": {
		"nonce": "1oQFzIwsLs8=",
		"mintfulfillment": {
			"type": 1,
			"data": {
				"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
				"signature": "ad59389329ed01c5ee14ce25ae38634c2b3ef694a2bdfa714f73b175f979ba6613025f9123d68c0f11e8f0a7114833c0aab4c8596d4c31671ec8a73923f02305"
			}
		},
		"coinoutputs": [{
			"value": "500000000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01e3cbc41bd3cdfec9e01a6be46a35099ba0e1e1b793904fce6aa5a444496c6d815f5e3e981ccf"
				}
			}
		}],
		"minerfees": ["1000000000"],
		"arbitrarydata": "dGVzdC4uLiAxLCAyLi4uIDM="
	}
}`

func TestSignCoinCreationTransactionExtension(t *testing.T) {
	inMemoryMintConditionGetter := newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
		unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))))
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: inMemoryMintConditionGetter,
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(validDevnetJSONEncodedCoinCreationTx))
	if err != nil {
		t.Fatal("failed to decode valid coin creation tx:", err)
	}

	// util function to sign tx
	type testKeyPair struct {
		types.KeyPair
		SignCount int
	}
	signTxAndValidate := func(key interface{}) error {
		t.Helper()

		signCount := 1

		// redefine fulfillment, as signing an already signed fulfillment is not possible
		cctxExtension := tx.Extension.(*CoinCreationTransactionExtension)
		switch mc, _ := inMemoryMintConditionGetter.GetActiveMintCondition(); mc.ConditionType() {
		case types.ConditionTypeUnlockHash:
			tx.Extension = &CoinCreationTransactionExtension{
				Nonce: cctxExtension.Nonce,
				MintFulfillment: types.NewFulfillment(types.NewSingleSignatureFulfillment(
					cctxExtension.MintFulfillment.Fulfillment.(*types.SingleSignatureFulfillment).PublicKey,
				)),
			}
		case types.ConditionTypeMultiSignature:
			k, ok := key.(testKeyPair)
			if !ok {
				panic("sign key should be a testKeyPair")
			}
			signCount = k.SignCount
			key = k.KeyPair
			tx.Extension = &CoinCreationTransactionExtension{
				Nonce:           cctxExtension.Nonce,
				MintFulfillment: types.NewFulfillment(types.NewMultiSignatureFulfillment(nil)),
			}
		default:
			panic("unsupported condition type")
		}
		for i := 0; i < signCount; i++ {
			// sign fulfillment, which lives in the tx extension data
			err := tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, eo ...interface{}) error {
				return fulfillment.Sign(types.FulfillmentSignContext{
					ExtraObjects: eo,
					Transaction:  tx,
					Key:          key,
				})
			})
			if err != nil {
				return fmt.Errorf("sign #%d failed: %v", i, err)
			}
		}
		validationCtx := types.ValidationContext{
			Confirmed:   true,
			BlockHeight: 4072,
			BlockTime:   1534271219,
		}
		chainConstants := config.GetDevnetGenesis()
		txValidationConstants := types.TransactionValidationConstants{
			BlockSizeLimit:         chainConstants.BlockSizeLimit,
			ArbitraryDataSizeLimit: chainConstants.ArbitraryDataSizeLimit,
			MinimumMinerFee:        chainConstants.MinimumTransactionFee,
		}
		return tx.ValidateTransaction(validationCtx, txValidationConstants)
	}

	// sign as the signer did on devnet, should succeed
	err = signTxAndValidate(hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"))
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	// sign as the signer did on devnet, should succeed
	err = signTxAndValidate(func() crypto.SecretKey { sk, _ := crypto.GenerateKeyPair(); return sk }())
	if err == nil {
		t.Fatalf("succeeded to sign, while it should fail")
	}

	// overwrite coin tx multisig condition to be a multisig condition instead
	inMemoryMintConditionGetter.applyMintCondition(0, types.NewCondition(types.NewMultiSignatureCondition(
		types.UnlockHashSlice{
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("016438a548b6d377e87b08e8eae5ef641a4e70cc861b85b54b0921330e03084ffe0a8d9a38e3a8"),
		},
		2,
	)))

	// sign multisig condition, should fail as we didn't sign enough
	err = signTxAndValidate(testKeyPair{
		KeyPair: types.KeyPair{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			},
			PrivateKey: hbs("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		},
		SignCount: 1,
	})
	if err == nil {
		t.Fatal("should fail to validate as we didn't sign twice, but it succeeded")
	}

	// sign multisig condition, should succeed
	err = signTxAndValidate(testKeyPair{
		KeyPair: types.KeyPair{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			},
			PrivateKey: hbs("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		},
		SignCount: 2,
	})
	if err != nil {
		t.Fatalf("failed to sign multisig: %v", err)
	}
}

func TestCoinCreationTransactionValidation(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: newInMemoryMintConditionGetter(types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(validDevnetJSONEncodedCoinCreationTx))
	if err != nil {
		t.Fatal("failed to decode valid coin creation tx:", err)
	}

	// util function to resign tx when needed,
	// as to ensure we fail because of the reason we want it to fail
	removeTxSignature := func() {
		cctxExtension := tx.Extension.(*CoinCreationTransactionExtension)
		tx.Extension = &CoinCreationTransactionExtension{
			Nonce: cctxExtension.Nonce,
			MintFulfillment: types.NewFulfillment(types.NewSingleSignatureFulfillment(
				cctxExtension.MintFulfillment.Fulfillment.(*types.SingleSignatureFulfillment).PublicKey,
			)),
		}
	}
	resignTx := func(description string) {
		t.Helper()
		// redefine fulfillment, as signing an already signed fulfillment is not possible
		removeTxSignature()
		// sign fulfillment, which lives in the tx extension data
		err := tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, eo ...interface{}) error {
			return fulfillment.Sign(types.FulfillmentSignContext{
				ExtraObjects: eo,
				Transaction:  tx,
				Key:          hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			})
		})
		if err != nil {
			t.Fatalf("failed to resign after %q: %v", description, err)
		}
	}

	validationCtx := types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 4072,
		BlockTime:   1534271219,
	}
	chainConstants := config.GetDevnetGenesis()
	txValidationConstants := types.TransactionValidationConstants{
		BlockSizeLimit:         chainConstants.BlockSizeLimit,
		ArbitraryDataSizeLimit: chainConstants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        chainConstants.MinimumTransactionFee,
	}

	// should be valid, as it was published on a local devnet
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate coin creation tx, while it is supposed to be valid:", err)
	}

	origExtension := tx.Extension
	origCCExtension := origExtension.(*CoinCreationTransactionExtension)

	// mess with nonce
	tx.Extension = &CoinCreationTransactionExtension{
		Nonce: TransactionNonce{}, // nil-nonce is not allowed
		MintFulfillment: types.NewFulfillment(types.NewSingleSignatureFulfillment(
			origCCExtension.MintFulfillment.Fulfillment.(*types.SingleSignatureFulfillment).PublicKey,
		)),
	}
	// sign fulfillment, which lives in the tx extension data
	err = tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, eo ...interface{}) error {
		return fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: eo,
			Transaction:  tx,
			Key:          hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
		})
	})
	if err != nil {
		t.Fatalf("failed to resign after modifying Nonce to use NilTransactionNonce: %v", err)
	}
	// validate, should fail
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of nil TransactionNonce")
	}
	// use orig extension as well, as we modified and re-signed
	tx.Extension = origExtension

	// make the arbitrary data too big, should fail
	origArbitraryData := tx.ArbitraryData
	tx.ArbitraryData.Data = make([]byte, chainConstants.ArbitraryDataSizeLimit+1)
	resignTx("changed arbitrary data")
	// should fail now
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of exceeded arbitrary data byte size limit")
	}
	// restore to valid arbitrary data
	tx.ArbitraryData = origArbitraryData
	// use orig extension as well, as we signed
	tx.Extension = origExtension

	// mess with the extension to make it fail
	// should all fail
	tx.Extension = nil
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of nil extension")
	}
	tx.Extension = &CoinCreationTransactionExtension{
		Nonce:           RandomTransactionNonce(),
		MintFulfillment: origCCExtension.MintFulfillment,
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate minter definition tx, " +
			"while it is supposed to fail because of incompatible nonce with signature")
	}
	tx.Extension = &CoinCreationTransactionExtension{
		Nonce: origCCExtension.Nonce,
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of nil fulfillment")
	}
	tx.Extension = &CoinCreationTransactionExtension{
		Nonce: origCCExtension.Nonce,
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.PublicKey{
						Algorithm: types.SignatureAlgoEd25519,
						Key:       hbs("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
					},
					Signature: hbs("abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab"),
				},
			},
		}),
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of wrong fulfillment")
	}
	tx.Extension = &CoinCreationTransactionExtension{
		Nonce: origCCExtension.Nonce,
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d"),
			},
			Signature: hbs("3e2ed4e893f66ffd57e26afe83d570ca4b8ba873f8236a60c018cde4852de1027256d088b2253ec061ae973f961f26cde8fa42f5d3c0ce1316560ceb25786f03"),
		}),
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of wrong fulfillment")
	}
	tx.Extension = &CoinCreationTransactionExtension{
		Nonce: origCCExtension.Nonce,
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.PublicKey{
				Algorithm: types.SignatureAlgoEd25519,
				Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
			},
			Signature: nil,
		}),
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of wrong fulfillment")
	}
	// restore to valid extension
	tx.Extension = origExtension

	// at least one coin output is required
	origCoinOutputs := tx.CoinOutputs
	tx.CoinOutputs = nil
	removeTxSignature()
	// should fail now
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of no coin outputs defined")
	}
	// restore to valid coin outputs
	tx.CoinOutputs = origCoinOutputs
	// replace extension as well with orig, as we resigned
	tx.Extension = origExtension

	// at least one miner fee is given,
	// and each miner fee has to be at least the minimum defined miner fee amount
	origMinerFees := tx.MinerFees
	tx.MinerFees = nil
	removeTxSignature()
	// should all fail now
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of no miner fees defined")
	}
	tx.MinerFees = []types.Currency{types.ZeroCurrency}
	resignTx("use zero miner fee")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of a too low miner fee defined")
	}
	tx.MinerFees = []types.Currency{chainConstants.MinimumTransactionFee, types.ZeroCurrency}
	resignTx("use minimum and zero miner fees")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of a too low miner fee defined")
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee.Sub(types.NewCurrency64(1)),
		types.NewCurrency64(1),
	}
	resignTx("total miner fees is good enough")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of a too low miner fee defined")
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee,
	}
	resignTx("use minimum miner fee")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate coin creation tx, while it is supposed to be valid:", err)
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee,
		chainConstants.MinimumTransactionFee,
	}
	resignTx("use more minir fees than required")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate coin creation tx, while it is supposed to be valid:", err)
	}
	tx.MinerFees = []types.Currency{
		chainConstants.MinimumTransactionFee.Mul64(1000),
		chainConstants.MinimumTransactionFee,
	}
	resignTx("use more minir fees than required")
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate coin creation tx, while it is supposed to be valid:", err)
	}
	// these should all work though
	// restore to valid miner fees
	tx.MinerFees = origMinerFees
	// replace extension as well with orig, as we resigned
	tx.Extension = origExtension

	// restored as it is, it should be valid
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err != nil {
		t.Fatal("failed to validate coin creation tx, while it is supposed to be valid:", err)
	}
}

func TestCoinCreationTransactionValidationWithUnknownMintCondition(t *testing.T) {
	mintConditionGetter := newInMemoryMintConditionGetter()
	mintConditionGetter.applyMintCondition(1, types.NewCondition(types.NewUnlockHashCondition(
		unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))))

	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintConditionGetter: mintConditionGetter,
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(validDevnetJSONEncodedCoinCreationTx))
	if err != nil {
		t.Fatal("failed to decode valid coin creation tx:", err)
	}

	validationCtx := types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 0,
		BlockTime:   1534271219,
	}
	chainConstants := config.GetDevnetGenesis()
	txValidationConstants := types.TransactionValidationConstants{
		BlockSizeLimit:         chainConstants.BlockSizeLimit,
		ArbitraryDataSizeLimit: chainConstants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        chainConstants.MinimumTransactionFee,
	}

	// should be invalid, as no mint condition exists for that height (0)
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation transaction, " +
			"while it was expected to fail due to an unknown mint condition")
	}
}

// test to ensure we json-encode by default the tx nonce as a base64-encoded string
func TestBase64DecodingOfJSONEncodedNonce(t *testing.T) {
	nonce := RandomTransactionNonce()
	b, err := json.Marshal(nonce)
	if err != nil {
		t.Fatal("failed to json-encode random tx nonce:", err)
	}
	db, err := base64.StdEncoding.DecodeString(string(b[1 : len(b)-1])) // remove quotes
	if err != nil {
		t.Fatalf("failed to base64-decode random tx nonce %s: %v", string(b), err)
	}
	if bytes.Compare(nonce[:], db[:]) != 0 {
		t.Fatal("unexpected result:", hex.EncodeToString(nonce[:]), "!=", hex.EncodeToString(db))
	}
}

// test to ensure that we can still accept a JSON byte array,
// as a nonce, should the input have used that.
func TestJSONDecodeArrayTxNonce(t *testing.T) {
	var nonce TransactionNonce
	err := json.Unmarshal([]byte("[52,82,198,39,242,116,81,220]"), &nonce)
	if err != nil {
		t.Fatal("failed to json-decode tx nonce:", err)
	}
	expectedNonce := TransactionNonce{52, 82, 198, 39, 242, 116, 81, 220}
	if bytes.Compare(expectedNonce[:], nonce[:]) != 0 {
		t.Fatal("unexpected result:", hex.EncodeToString(expectedNonce[:]), "!=", hex.EncodeToString(nonce[:]))
	}
}

func TestNonNilRandomNonceCheck(t *testing.T) {
	for i := 0; i < 1000; i++ {
		nonce := RandomTransactionNonce()
		if nonce == (TransactionNonce{}) {
			panic("nil TransactionNonce crypto-rand created")
		}
	}
}

// utility funcs
func hbs(str string) []byte { // hexStr -> byte slice
	bs, _ := hex.DecodeString(str)
	return bs
}
func hs(str string) (hash crypto.Hash) { // hbs -> crypto.Hash
	copy(hash[:], hbs(str))
	return
}
func hsk(str string) (pk crypto.SecretKey) { // hbs -> crypto.SecretKey
	copy(pk[:], hbs(str))
	return
}

// tests for utility types

func TestInMemoryMintConditionGetter(t *testing.T) {
	getter := newInMemoryMintConditionGetter()
	if getter == nil {
		t.Fatal("nil getter")
	}
	mcAsStr := func(cp types.UnlockConditionProxy) string {
		b, err := cp.MarshalJSON()
		if err != nil {
			panic("failed to JSON-marshal given UnlockConditionProxy: " + err.Error())
		}
		return string(b)
	}

	// no mint condition should be found in nil-InMemoryMintConditionGetter
	mc, err := getter.GetActiveMintCondition()
	if err == nil {
		t.Fatal("no active mint condition was expected, but found one: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(0)
	if err == nil {
		t.Fatal("no mint condition was expected at height 0, but found one: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(42)
	if err == nil {
		t.Fatal("no mint condition was expected at height 42, but found one: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(types.BlockHeight(math.MaxUint64))
	if err == nil {
		t.Fatal("no mint condition was expected at height math.MaxUint64, but found one: ", mcAsStr(mc))
	}

	// apply at height 1  & check
	conditionOne := types.NewCondition(types.NewUnlockHashCondition(unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))
	getter.applyMintCondition(1, conditionOne)
	mc, err = getter.GetActiveMintCondition()
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected active mint condition: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(1)
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected mint condition at height 1: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(42)
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected mint condition at height 42: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(types.BlockHeight(math.MaxUint64))
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected mint condition at height 1: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(0)
	if err == nil {
		t.Fatal("no mint condition was expected, but found one at height 0: ", mcAsStr(mc))
	}

	// add at height 0 & check
	conditionTwo := types.NewCondition(types.NewUnlockHashCondition(unlockHashFromHex("01e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70b1ccc65e2105")))
	getter.applyMintCondition(0, conditionTwo)
	mc, err = getter.GetActiveMintCondition()
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected active mint condition: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(1)
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected mint condition at height 1: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(42)
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected mint condition at height 42: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(types.BlockHeight(math.MaxUint64))
	if err != nil {
		t.Fatal("an active mint condition was expected, but none was found: ", err)
	}
	if !conditionOne.Equal(mc) {
		t.Fatal("unexpected mint condition at height 1: ", mcAsStr(mc))
	}
	mc, err = getter.GetMintConditionAt(0)
	if err != nil {
		t.Fatal("a mint condition was expected at height, but received err: ", err)
	}
	if !conditionTwo.Equal(mc) {
		t.Fatal("unexpected mint condition at height 0: ", mcAsStr(mc))
	}

	for _, heightToRevert := range []types.BlockHeight{1, 42, 3} {
		// revert && check
		getter.revertMintCondition(heightToRevert)

		mc, err = getter.GetActiveMintCondition()
		if err != nil {
			t.Fatal("an active mint condition was expected, but none was found: ", err)
		}
		if !conditionTwo.Equal(mc) {
			t.Fatal("unexpected active mint condition: ", mcAsStr(mc))
		}
		mc, err = getter.GetMintConditionAt(1)
		if err != nil {
			t.Fatal("an active mint condition was expected, but none was found: ", err)
		}
		if !conditionTwo.Equal(mc) {
			t.Fatal("unexpected mint condition at height 1: ", mcAsStr(mc))
		}
		mc, err = getter.GetMintConditionAt(42)
		if err != nil {
			t.Fatal("an active mint condition was expected, but none was found: ", err)
		}
		if !conditionTwo.Equal(mc) {
			t.Fatal("unexpected mint condition at height 42: ", mcAsStr(mc))
		}
		mc, err = getter.GetMintConditionAt(types.BlockHeight(math.MaxUint64))
		if err != nil {
			t.Fatal("an active mint condition was expected, but none was found: ", err)
		}
		if !conditionTwo.Equal(mc) {
			t.Fatal("unexpected mint condition at height 1: ", mcAsStr(mc))
		}
		mc, err = getter.GetMintConditionAt(0)
		if err != nil {
			t.Fatal("a mint condition was expected at height, but received err: ", err)
		}
		if !conditionTwo.Equal(mc) {
			t.Fatal("unexpected mint condition at height 0: ", mcAsStr(mc))
		}
	}
}

// utility types

type (
	inMemoryMintConditionGetter struct {
		mintConditions []conditionHeightPair
	}
	conditionHeightPair struct {
		Height        types.BlockHeight
		MintCondition types.UnlockConditionProxy
	}
)

// newInMemoryMintConditionGetter creates a new inMemoryMintConditionGetterr,
// applying all given mint conditions using their index-as-given-in-order as the represenative block height.
func newInMemoryMintConditionGetter(mintConditions ...types.UnlockConditionProxy) *inMemoryMintConditionGetter {
	mem := new(inMemoryMintConditionGetter)
	for index, mintCondition := range mintConditions {
		mem.mintConditions = append(mem.mintConditions,
			conditionHeightPair{types.BlockHeight(index), mintCondition})
	}
	return mem
}

// GetActiveMintCondition implements MintConditionGetter.GetActiveMintCondition
func (mem *inMemoryMintConditionGetter) GetActiveMintCondition() (types.UnlockConditionProxy, error) {
	n := len(mem.mintConditions)
	if n == 0 {
		return types.UnlockConditionProxy{}, errors.New("no mint condition is applied")
	}
	return mem.mintConditions[n-1].MintCondition, nil
}

// GetMintConditionAt implements MintConditionGetter.GetMintConditionAt
func (mem *inMemoryMintConditionGetter) GetMintConditionAt(height types.BlockHeight) (types.UnlockConditionProxy, error) {
	var (
		maxHeight     types.BlockHeight
		mintCondition types.UnlockConditionProxy
	)
	for _, pair := range mem.mintConditions {
		switch {
		case pair.Height == height:
			mintCondition, maxHeight = pair.MintCondition, pair.Height
			break
		case pair.Height > height:
			break
		case pair.Height >= maxHeight:
			mintCondition, maxHeight = pair.MintCondition, pair.Height
		}
	}
	if mintCondition.ConditionType() == types.ConditionTypeNil {
		return types.UnlockConditionProxy{}, errors.New("no mint condition is applied")
	}
	return mintCondition, nil
}

// apply/revert mint conditions for a given block height,
// also keeping track of the highest
func (mem *inMemoryMintConditionGetter) applyMintCondition(height types.BlockHeight, mintCondition types.UnlockConditionProxy) {
	for index, pair := range mem.mintConditions {
		switch {
		case height > pair.Height:
			// continue searching '<=' case
			continue
		case height == pair.Height:
			// ovewrite
			mem.mintConditions[index] = conditionHeightPair{height, mintCondition}
			return // overwritten and finished
		default:
			// push to the front or insert somewhere in the middle
			mem.mintConditions = append(mem.mintConditions[:index],
				append([]conditionHeightPair{{height, mintCondition}}, mem.mintConditions[index:]...)...)
			return // pushed and finished
		}
	}
	// append to the back
	mem.mintConditions = append(mem.mintConditions, conditionHeightPair{height, mintCondition})
}

// fint the condition for the given height and cut it out
func (mem *inMemoryMintConditionGetter) revertMintCondition(height types.BlockHeight) {
	for index, pair := range mem.mintConditions {
		if pair.Height == height {
			mem.mintConditions = append(mem.mintConditions[:index], mem.mintConditions[index+1:]...)
			return
		}
	}
}

func TestBotTransactionVersionConstants(t *testing.T) {
	if TransactionVersionBotRegistration != 0x90 {
		t.Errorf("unexpected bot registration Tx version: %x", TransactionVersionBotRegistration)
	}
	if TransactionVersionBotRecordUpdate != 0x91 {
		t.Errorf("unexpected bot record update Tx version: %x", TransactionVersionBotRecordUpdate)
	}
	if TransactionVersionBotNameTransfer != 0x92 {
		t.Errorf("unexpected bot name transfer Tx version: %x", TransactionVersionBotNameTransfer)
	}
}

func TestBotMonthsAndFlagsData(t *testing.T) {
	testCases := []struct {
		Input    uint8
		Expected BotMonthsAndFlagsData
	}{
		{0, BotMonthsAndFlagsData{}},
		{4 | 32, BotMonthsAndFlagsData{NrOfMonths: 4, HasAddresses: true}},
		{1 | 64, BotMonthsAndFlagsData{NrOfMonths: 1, HasNames: true}},
		{29 | 128, BotMonthsAndFlagsData{NrOfMonths: 29, HasRefund: true}},
		{2 | 32 | 128, BotMonthsAndFlagsData{NrOfMonths: 2, HasAddresses: true, HasRefund: true}},
		{3 | 32 | 64, BotMonthsAndFlagsData{NrOfMonths: 3, HasAddresses: true, HasNames: true}},
		{5 | 64 | 128, BotMonthsAndFlagsData{NrOfMonths: 5, HasNames: true, HasRefund: true}},
		{31 | 32 | 64 | 128, BotMonthsAndFlagsData{NrOfMonths: 31, HasAddresses: true, HasNames: true, HasRefund: true}},
	}
	for idx, testCase := range testCases {
		var result BotMonthsAndFlagsData
		err := rivbin.Unmarshal(rivbin.Marshal(testCase.Input), &result)
		if err != nil {
			t.Error(idx, "error(Unmarshal:BotMonthsAndFlagData)", err)
			continue
		}
		if result != testCase.Expected {
			t.Error(idx, "unexpected result", result, "!=", testCase.Expected)
			continue
		}
		var number uint8
		err = rivbin.Unmarshal(rivbin.Marshal(result), &number)
		if err != nil {
			t.Error(idx, "error(Unmarshal:uint8)", err)
			continue
		}
		if number != testCase.Input {
			t.Error(idx, "unexpected number result", number, "!=", testCase.Input)
		}
	}
}

func TestComputeMonthlyBotFees(t *testing.T) {
	oneCoin := types.NewCurrency64(10)
	testCases := []struct {
		NrOfMonths  uint8
		ExpectedFee types.Currency
	}{
		{0, types.Currency{}},
		{1, types.NewCurrency64(100)},
		{2, types.NewCurrency64(200)},
		{8, types.NewCurrency64(800)},
		{11, types.NewCurrency64(1100)},
		{12, types.NewCurrency64(840)},
		{16, types.NewCurrency64(1120)},
		{23, types.NewCurrency64(1610)},
		{24, types.NewCurrency64(1200)},
		{31, types.NewCurrency64(1550)},
		{32, types.NewCurrency64(1600)},
		{math.MaxUint8, types.NewCurrency64(12750)},
	}
	for idx, testCase := range testCases {
		fee := ComputeMonthlyBotFees(testCase.NrOfMonths, oneCoin)
		if !fee.Equals(testCase.ExpectedFee) {
			t.Error(idx, testCase.NrOfMonths, "unexpected result", fee, "!=", testCase.ExpectedFee)
		}
	}
}

func TestBotRegistrationTransactionBinaryEncodingAndID(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	const input = `{"version":144,"data":{"addresses":null,"names":["crazybot.foobar"],"nrofmonths":1,"txfee":"1000000000","coininputs":[{"parentid":"6678e3a75da2026da76753a60ac44f7e7737784015676b37cc2cdcf670dce2e5","fulfillment":{"type":1,"data":{"publickey":"ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780","signature":"cd07fbfd78be0edd1c9ca46bc18f91cde1ed05848083828c5d3848cd9671054527b630af72f7d95c0ddcd3a0f0c940eb8cfe4b085cb00efc8338b28f39155809"}}}],"refundcoinoutput":{"value":"99979897000000000","condition":{"type":1,"data":{"unlockhash":"017fda17489854109399aa8c1bfa6bdef40f93606744d95cc5055270d78b465e6acd263c96ab2b"}}},"identification":{"publickey":"ed25519:adc4090edbe28e3628f08a85d20b5055ea301cdb080d3b65a337a326e2e3556d","signature":"5211f813fb4e34ae348e2e746846bc72255512dc246ccafbb3bd3b916aac738bfe2737308d87cced4f9476be8715983cc6000e37f8e82e7b83f120776a358105"}}}`
	var tx types.Transaction
	err := json.Unmarshal([]byte(input), &tx)
	if err != nil {
		t.Fatal(err)
	}
	id := tx.ID()
	b := rivbin.Marshal(tx)

	// go to 3bot Tx and back
	botRegistrationTx, err := BotRegistrationTransactionFromTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	oTx := botRegistrationTx.Transaction(config.GetCurrencyUnits().OneCoin)
	oID := oTx.ID()
	oB := rivbin.Marshal(oTx)
	if id != oID {
		t.Fatal(id, "!=", oID)
	}
	if !bytes.Equal(b, oB) {
		t.Fatal(hex.EncodeToString(b), "!=", hex.EncodeToString(oB))
	}
}

func TestBotRegistrationExtractedFromBlockConsensusDB(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	const (
		hexBlock       = `0d3a8d36b50c3325044b5d994e52f00ce86b43ff84bdc0e7a1347c9b7621624ccf5af45b000000000100000000000000000000000000000000000000000000000200000000000000050000000000000002540be400015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e67915804000000000000003b9aca00015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791580200000000000000010d010000000000000000000000000000000000000000000001000000000000001d7f4ac218a2f360dd802843a0003443f77d151ba9329fdecbd8da37519b3419018000000000000000656432353531390000000000000000002000000000000000d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d77804000000000000000b82990bcbdd96acb14a877f8b0364abbd8ceab232ce9caa3f8f3a15f7277978484a390d928cce671e9829d780715a6aaf8c686cc7074f7d558b03a4a73f96b07010000000000000002000000000000000bb8012100000000000000015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791580000000000000000000000000000000090e112115bc6aec02c6578616d706c652e6f72671e63686174626f742e6578616d706c65083b9aca0002a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee56301c401d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d778080909a7df820ec3cee1c99bd2c297b938f830da891439ef7d78452e29efb0c7e593683274c356f72d3b627c2954a24b2bc2276fed47b24cd62816c540c88f13d051001634560d9784e00014201b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960100bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c67780661498e71668dfe7726a357039d7c0e871b6c0ca8fa49dc1fcdccb5f23f5f0a5cab95cfcfd72a9fd2c5045ba899ecb0207ff01125a0151f3e35e3c6e13a7538b340a`
		expectedJSONTx = `{"version":144,"data":{"addresses":["91.198.174.192","example.org"],"names":["chatbot.example"],"nrofmonths":1,"txfee":"1000000000","coininputs":[{"parentid":"a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee563","fulfillment":{"type":1,"data":{"publickey":"ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780","signature":"909a7df820ec3cee1c99bd2c297b938f830da891439ef7d78452e29efb0c7e593683274c356f72d3b627c2954a24b2bc2276fed47b24cd62816c540c88f13d05"}}}],"refundcoinoutput":{"value":"99999899000000000","condition":{"type":1,"data":{"unlockhash":"01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"}}},"identification":{"publickey":"ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614","signature":"98e71668dfe7726a357039d7c0e871b6c0ca8fa49dc1fcdccb5f23f5f0a5cab95cfcfd72a9fd2c5045ba899ecb0207ff01125a0151f3e35e3c6e13a7538b340a"}}}`
	)

	b, err := hex.DecodeString(hexBlock)
	if err != nil {
		t.Fatal(err)
	}
	var block types.Block
	err = siabin.Unmarshal(b, &block)
	if err != nil {
		t.Fatal(err)
	}
	tx := block.Transactions[1]
	b, err = json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	jsonTx := string(b)
	if expectedJSONTx != jsonTx {
		t.Fatal(expectedJSONTx, "!=", jsonTx)
	}
}

var cryptoKeyPair = types.KeyPair{
	PublicKey: types.PublicKey{
		Algorithm: types.SignatureAlgoEd25519,
		Key:       hbs("d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
	},
	PrivateKey: hbs("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
}

func TestBotRegistrationTransactionUniqueSignatures(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(fmt.Sprintf(`{
	"version": 144,
	"data": {
		"names": ["foobar"],
		"nrofmonths": 1,
		"txfee": "1000000000",
		"coininputs": [
			{
				"parentid": "a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee563",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "%[1]s",
						"signature": ""
					}
				}
			},
			{
				"parentid": "91431da29b53669cdaecf5e31d9ae4d47fe4ebbd02e12fec185e28b7db6960dd",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "%[1]s",
						"signature": ""
					}
				}
			}
		],
		"refundcoinoutput": {
			"value": "99999899000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"
				}
			}
		},
		"identification": {
			"publickey": "%[1]s",
			"signature": ""
		}
	}
}`, cryptoKeyPair.PublicKey.String())))
	if err != nil {
		t.Fatal(err)
	}

	signatures := map[string]struct{}{}
	// sign coin inputs, validate a signature is defined and ensure they are unique
	for cindex, ci := range tx.CoinInputs {
		err = ci.Fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: []interface{}{uint64(cindex)},
			Transaction:  tx,
			Key:          cryptoKeyPair.PrivateKey,
		})
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}

		b, err := json.Marshal(ci.Fulfillment)
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}
		var rawFulfillment map[string]interface{}
		err = json.Unmarshal(b, &rawFulfillment)
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}
		signature := rawFulfillment["data"].(map[string]interface{})["signature"].(string)
		if signature == "" {
			t.Error(cindex, "coin input: signature is empty")
			continue
		}
		if _, ok := signatures[signature]; ok {
			t.Error(cindex, "coin input: signature exists already:", signature)
			continue
		}
		signatures[signature] = struct{}{}
	}

	// sign extension (the actual signature)
	err = tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, extraObjects ...interface{}) error {
		if condition.UnlockHash().Cmp(types.NewPubKeyUnlockHash(cryptoKeyPair.PublicKey)) != 0 {
			b, _ := json.Marshal(condition)
			t.Fatalf("unexpected extension fulfill condition: %v", string(b))
		}
		return fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: extraObjects,
			Transaction:  tx,
			Key:          cryptoKeyPair.PrivateKey,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	signature := tx.Extension.(*BotRegistrationTransactionExtension).Identification.Signature.String()
	if signature == "" {
		t.Fatal("extension (Sender): signature is empty")
	}
	if _, ok := signatures[signature]; ok {
		t.Fatal("extension (Sender): signature exists already:", signature)
	}
	signatures[signature] = struct{}{}
}

func TestBotRecordUpdateTransactionUniqueSignatures(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRecordUpdate, BotUpdateRecordTransactionController{
		Registry: &inMemoryBotRegistry{
			idMapping: map[BotID]BotRecord{
				1: botRecordFromJSON(t, `{
	"id": 1,
	"addresses": ["93.184.216.34"],
	"names": ["example"],
	"publickey": "`+cryptoKeyPair.PublicKey.String()+`",
	"expiration": 1538484360
}`),
			},
		},
		OneCoin: config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRecordUpdate, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(fmt.Sprintf(`{
	"version": 145,
	"data": {
		"id": 1,
		"addresses": {
			"add": ["127.0.0.1", "api.mybot.io", "0:0:0:0:0:ffff:5db8:d822"],
			"remove": ["93.184.216.34"]
		},
		"names": {
			"add": ["mybot"],
			"remove": ["example"]
		},
		"nrofmonths": 5,
		"txfee": "1000000000",
		"coininputs": [
			{
				"parentid": "c6b161d192d8095efd4d9946f7d154bf335f51fdfdeca4bb0cb990b25ffd7e95",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "%[1]s",
						"signature": ""
					}
				}
			},
			{
				"parentid": "91431da29b53669cdaecf5e31d9ae4d47fe4ebbd02e12fec185e28b7db6960dd",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "%[1]s",
						"signature": ""
					}
				}
			}
		],
		"refundcoinoutput": {
			"value": "99999778000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "0161fbcf58efaeba8813150e88fc33405b3a77d51277a2cdf3f4d2ab770de287c7af9d456c4e68"
				}
			}
		},
		"signature": ""
	}
}`, cryptoKeyPair.PublicKey.String())))
	if err != nil {
		t.Fatal(err)
	}

	signatures := map[string]struct{}{}
	// sign coin inputs, validate a signature is defined and ensure they are unique
	for cindex, ci := range tx.CoinInputs {
		err = ci.Fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: []interface{}{uint64(cindex)},
			Transaction:  tx,
			Key:          cryptoKeyPair.PrivateKey,
		})
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}

		b, err := json.Marshal(ci.Fulfillment)
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}
		var rawFulfillment map[string]interface{}
		err = json.Unmarshal(b, &rawFulfillment)
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}
		signature := rawFulfillment["data"].(map[string]interface{})["signature"].(string)
		if signature == "" {
			t.Error(cindex, "coin input: signature is empty")
			continue
		}
		if _, ok := signatures[signature]; ok {
			t.Error(cindex, "coin input: signature exists already:", signature)
			continue
		}
		signatures[signature] = struct{}{}
	}

	// sign extension (the actual signature)
	err = tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, extraObjects ...interface{}) error {
		if condition.UnlockHash().Cmp(types.NewPubKeyUnlockHash(cryptoKeyPair.PublicKey)) != 0 {
			b, _ := json.Marshal(condition)
			t.Fatalf("unexpected extension fulfill condition: %v", string(b))
		}
		return fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: extraObjects,
			Transaction:  tx,
			Key:          cryptoKeyPair.PrivateKey,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	signature := tx.Extension.(*BotRecordUpdateTransactionExtension).Signature.String()
	if signature == "" {
		t.Fatal("extension (Sender): signature is empty")
	}
	if _, ok := signatures[signature]; ok {
		t.Fatal("extension (Sender): signature exists already:", signature)
	}
	signatures[signature] = struct{}{}
}

func TestBotNameTransferTransactionUniqueSignatures(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotNameTransfer, BotNameTransferTransactionController{
		Registry: &inMemoryBotRegistry{
			idMapping: map[BotID]BotRecord{
				1: botRecordFromJSON(t, `{
	"id": 1,
	"addresses": ["93.184.216.34"],
	"names": ["example"],
	"publickey": "`+cryptoKeyPair.PublicKey.String()+`",
	"expiration": 1538484360
}`),
			},
		},
		OneCoin: config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotNameTransfer, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(fmt.Sprintf(`{
	"version": 146,
	"data": {
		"sender": {
			"id": 1,
			"signature": ""
		},
		"receiver": {
			"id": 1,
			"signature": ""
		},
		"names": [
			"mybot"
		],
		"txfee": "1000000000",
		"coininputs": [
			{
				"parentid": "c6b161d192d8095efd4d9946f7d154bf335f51fdfdeca4bb0cb990b25ffd7e95",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "%[1]s",
						"signature": ""
					}
				}
			},
			{
				"parentid": "91431da29b53669cdaecf5e31d9ae4d47fe4ebbd02e12fec185e28b7db6960dd",
				"fulfillment": {
					"type": 1,
					"data": {
						"publickey": "%[1]s",
						"signature": ""
					}
				}
			}
		],
		"refundcoinoutput": {
			"value": "99999626000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01822fd5fefd2748972ea828a5c56044dec9a2b2275229ce5b212f926cd52fba015846451e4e46"
				}
			}
		}
	}
}`, cryptoKeyPair.PublicKey.String())))
	if err != nil {
		t.Fatal(err)
	}

	signatures := map[string]struct{}{}
	// sign coin inputs, validate a signature is defined and ensure they are unique
	for cindex, ci := range tx.CoinInputs {
		err = ci.Fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: []interface{}{uint64(cindex)},
			Transaction:  tx,
			Key:          cryptoKeyPair.PrivateKey,
		})
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}

		b, err := json.Marshal(ci.Fulfillment)
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}
		var rawFulfillment map[string]interface{}
		err = json.Unmarshal(b, &rawFulfillment)
		if err != nil {
			t.Error(cindex, "coin input", err)
			continue
		}
		signature := rawFulfillment["data"].(map[string]interface{})["signature"].(string)
		if signature == "" {
			t.Error(cindex, "coin input: signature is empty")
			continue
		}
		if _, ok := signatures[signature]; ok {
			t.Error(cindex, "coin input: signature exists already:", signature)
			continue
		}
		signatures[signature] = struct{}{}
	}

	// sign extension (the actual signature)
	err = tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy, extraObjects ...interface{}) error {
		if condition.UnlockHash().Cmp(types.NewPubKeyUnlockHash(cryptoKeyPair.PublicKey)) != 0 {
			b, _ := json.Marshal(condition)
			t.Fatalf("unexpected extension fulfill condition: %v", string(b))
		}
		return fulfillment.Sign(types.FulfillmentSignContext{
			ExtraObjects: extraObjects,
			Transaction:  tx,
			Key:          cryptoKeyPair.PrivateKey,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	ext := tx.Extension.(*BotNameTransferTransactionExtension)
	extSignatures := []string{ext.Sender.Signature.String(), ext.Receiver.Signature.String()}
	for index, signature := range extSignatures {
		if signature == "" {
			t.Fatalf("extension (%d): signature is empty", index)
		}
		if _, ok := signatures[signature]; ok {
			t.Fatalf("extension (%d): signature exists already: %v", index, signature)
		}
		signatures[signature] = struct{}{}
	}
}

type inMemoryBotRegistry struct {
	idMapping map[BotID]BotRecord
}

func botRecordFromJSON(t *testing.T, str string) BotRecord {
	var record BotRecord
	err := json.Unmarshal([]byte(str), &record)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func (reg *inMemoryBotRegistry) GetRecordForID(id BotID) (*BotRecord, error) {
	if len(reg.idMapping) == 0 {
		return nil, errors.New("no records available")
	}
	record, ok := reg.idMapping[id]
	if !ok {
		return nil, fmt.Errorf("no record available for id %v", id)
	}
	return &record, nil
}

func (reg *inMemoryBotRegistry) GetRecordForKey(key types.PublicKey) (*BotRecord, error) {
	panic("NOT IMPLEMENTED")
}

func (reg *inMemoryBotRegistry) GetRecordForName(name BotName) (*BotRecord, error) {
	panic("NOT IMPLEMENTED")
}

func (reg *inMemoryBotRegistry) GetBotTransactionIdentifiers(id BotID) ([]types.TransactionID, error) {
	panic("NOT IMPLEMENTED")
}

func TestJSONExampleERC20ConvertTransaction(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionERC20Conversion, ERC20ConvertTransactionController{})
	defer types.RegisterTransactionVersion(TransactionVersionERC20Conversion, nil)

	const jsonEncodedExample = `{
	"version": 208,
	"data": {
		"address": "0123456789012345678901234567890123456789",
		"value": "200000000000",
		"txfee": "1000000000",
		"coininputs": [{
			"parentid": "9c61ec964105ec48bc95ffc0ac820ada600a2914a8dd4ef511ed7f218a3bf469",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:7469d51063cdb690cc8025db7d28faadc71ff69f7c372779bf3a1e801a923e02",
					"signature": "a0c683e8728710b4d3cd7eed4e1bd38a4be8145a2cf91b875986870aa98c6265d76cbb637d78500010e3ab1b651e31ab26b05de79938d7d0aee01f8566d08b09"
				}
			}
		}],
		"refundcoinoutput": {
			"value": "99999476000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "011c17aaf2d54f63644f9ce91c06ff984182483d1b943e96b5e77cc36fdb887c846b60460bceb0"
				}
			}
		}
	}
}`

	var tx types.Transaction
	err := json.Unmarshal([]byte(jsonEncodedExample), &tx)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	output := string(b)
	buffer := bytes.NewBuffer(nil)
	err = json.Compact(buffer, []byte(jsonEncodedExample))
	if err != nil {
		t.Fatal(err)
	}
	expectedOutput := string(buffer.Bytes())
	if expectedOutput != output {
		t.Fatal(expectedOutput, "!=", output)
	}
}

func TestBinaryExampleERC20ConvertTransaction(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionERC20Conversion, ERC20ConvertTransactionController{})
	defer types.RegisterTransactionVersion(TransactionVersionERC20Conversion, nil)

	const hexEncodedExample = `d001234567890123456789012345678901234567890a2e90edd000083b9aca00029c61ec964105ec48bc95ffc0ac820ada600a2914a8dd4ef511ed7f218a3bf46901c4017469d51063cdb690cc8025db7d28faadc71ff69f7c372779bf3a1e801a923e0280a0c683e8728710b4d3cd7eed4e1bd38a4be8145a2cf91b875986870aa98c6265d76cbb637d78500010e3ab1b651e31ab26b05de79938d7d0aee01f8566d08b090110016344fe5cb488000142011c17aaf2d54f63644f9ce91c06ff984182483d1b943e96b5e77cc36fdb887c84`

	b, err := hex.DecodeString(hexEncodedExample)
	if err != nil {
		t.Fatal(err)
	}
	var tx types.Transaction
	err = siabin.Unmarshal(b, &tx)
	if err != nil {
		t.Fatal(err)
	}

	b = siabin.Marshal(tx)
	output := hex.EncodeToString(b)
	if hexEncodedExample != output {
		t.Fatal(hexEncodedExample, "!=", output)
	}
}

func TestJSONExampleERC20CoinCreationTransaction(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, ERC20CoinCreationTransactionController{})
	defer types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, nil)

	const jsonEncodedExample = `{
	"version": 209,
	"data": {
		"address": "01f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609adbbcd3977363c",
		"value": "100000000000",
		"txfee": "1000000000",
		"txid": "0000000000000000000000000000000000000000000000000000000000000000"
	}
}`

	var tx types.Transaction
	err := json.Unmarshal([]byte(jsonEncodedExample), &tx)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	output := string(b)
	buffer := bytes.NewBuffer(nil)
	err = json.Compact(buffer, []byte(jsonEncodedExample))
	if err != nil {
		t.Fatal(err)
	}
	expectedOutput := string(buffer.Bytes())
	if expectedOutput != output {
		t.Fatal(expectedOutput, "!=", output)
	}
}

func TestBinaryExampleERC20CoinCreationTransaction(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, ERC20CoinCreationTransactionController{})
	defer types.RegisterTransactionVersion(TransactionVersionERC20CoinCreation, nil)

	const hexEncodedExample = `d101f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609ad0a174876e800083b9aca000000000000000000000000000000000000000000000000000000000000000000`

	b, err := hex.DecodeString(hexEncodedExample)
	if err != nil {
		t.Fatal(err)
	}
	var tx types.Transaction
	err = siabin.Unmarshal(b, &tx)
	if err != nil {
		t.Fatal(err)
	}

	b = siabin.Marshal(tx)
	output := hex.EncodeToString(b)
	if hexEncodedExample != output {
		t.Fatal(hexEncodedExample, "!=", output)
	}
}

func TestJSONExampleERC20AddressRegistrationTransaction(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, ERC20AddressRegistrationTransactionController{})
	defer types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, nil)

	const jsonEncodedExample = `{
	"version": 210,
	"data": {
		"pubkey": "ed25519:a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d",
		"tftaddress": "01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f",
		"erc20address": "828de486adc50aa52dab52a2ec284bcac75be211",
		"signature": "fe13823a96928a573f20a63f3b8d3cde08c506fa535d458120fdaa5f1c78f6939c81bf91e53393130fbfee32ff4e9cb6022f14ae7750d126a7b6c0202c674b02",
		"regfee": "10000000000",
		"txfee": "1000000000",
		"coininputs": [{
			"parentid": "a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee563",
			"fulfillment": {
				"type": 1,
				"data": {
					"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
					"signature": "4fe14adcbded85476680bfd4fa8ff35d51ac34bb8a9b3f4904eac6eee4f53e19b6a39c698463499b9961524f026db2fb5c8173307f483c6458d401ecec2e7a0c"
				}
			}
		}],
		"refundcoinoutput": {
			"value": "99999999000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01370af706b547dd4e562a047e6265d7e7750771f9bff633b1a12dbd59b11712c6ef65edb1690d"
				}
			}
		}
	}
}`

	var tx types.Transaction
	err := json.Unmarshal([]byte(jsonEncodedExample), &tx)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	output := string(b)
	buffer := bytes.NewBuffer(nil)
	err = json.Compact(buffer, []byte(jsonEncodedExample))
	if err != nil {
		t.Fatal(err)
	}
	expectedOutput := string(buffer.Bytes())
	if expectedOutput != output {
		t.Fatal(expectedOutput, "!=", output)
	}
}

func TestBinaryExampleERC20AddressRegistrationTransaction(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, ERC20AddressRegistrationTransactionController{})
	defer types.RegisterTransactionVersion(TransactionVersionERC20AddressRegistration, nil)

	const hexEncodedExample = `d201a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d80fe13823a96928a573f20a63f3b8d3cde08c506fa535d458120fdaa5f1c78f6939c81bf91e53393130fbfee32ff4e9cb6022f14ae7750d126a7b6c0202c674b020a02540be400083b9aca0002a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee56301c401d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780804fe14adcbded85476680bfd4fa8ff35d51ac34bb8a9b3f4904eac6eee4f53e19b6a39c698463499b9961524f026db2fb5c8173307f483c6458d401ecec2e7a0c01100163457821ef3600014201370af706b547dd4e562a047e6265d7e7750771f9bff633b1a12dbd59b11712c6`

	b, err := hex.DecodeString(hexEncodedExample)
	if err != nil {
		t.Fatal(err)
	}
	var tx types.Transaction
	err = siabin.Unmarshal(b, &tx)
	if err != nil {
		t.Fatal(err)
	}

	b = siabin.Marshal(tx)
	output := hex.EncodeToString(b)
	if hexEncodedExample != output {
		t.Fatal(hexEncodedExample, "!=", output)
	}
}
