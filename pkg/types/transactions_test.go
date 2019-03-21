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
	RegisterTransactionTypesForStandardNetwork(nil, NopERC20TransactionValidator{}, types.Currency{}, config.DaemonNetworkConfig{}) // no MintConditionGetter is required for this test
	testMinimumFeeValidationForTransactions(t, "standard", validationConstants)
	constants = config.GetTestnetGenesis()
	validationConstants = types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForTestNetwork(nil, NopERC20TransactionValidator{}, types.Currency{}, config.DaemonNetworkConfig{}) // no MintConditionGetter is required for this test
	testMinimumFeeValidationForTransactions(t, "test", validationConstants)
	constants = config.GetDevnetGenesis()
	validationConstants = types.TransactionValidationConstants{
		BlockSizeLimit:         constants.BlockSizeLimit,
		ArbitraryDataSizeLimit: constants.ArbitraryDataSizeLimit,
		MinimumMinerFee:        constants.MinimumTransactionFee,
	}
	RegisterTransactionTypesForDevNetwork(nil, NopERC20TransactionValidator{}, types.Currency{}, config.DaemonNetworkConfig{}) // no MintConditionGetter is required for this test
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
		ArbitraryData: []byte("data"),
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
		ArbitraryData: []byte("data"),
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
	if bytes.Compare(a.ArbitraryData, b.ArbitraryData) != 0 {
		t.Error(i, "arbitrary not equal",
			string(a.ArbitraryData), "!=", string(b.ArbitraryData))
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
		ArbitraryData: []byte("2300202+e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
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
	if bytes.Compare(a.ArbitraryData, b.ArbitraryData) != 0 {
		t.Error(i, "arbitrary not equal",
			string(a.ArbitraryData), "!=", string(b.ArbitraryData))
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
		ArbitraryData: []byte("2300202+e89843e4b8231a01ba18b254d530110364432aafab8206bea72e5a20eaa55f70"),
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
	tx.ArbitraryData = make([]byte, chainConstants.ArbitraryDataSizeLimit+1)
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
	tx.ArbitraryData = make([]byte, chainConstants.ArbitraryDataSizeLimit+1)
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
