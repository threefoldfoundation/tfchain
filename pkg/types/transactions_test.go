package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/rivine/rivine/crypto"
	"github.com/rivine/rivine/encoding"
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

// cctx -> Binary -> cctx
func TestCoinCreationTransactionToAndFromBinary(t *testing.T) {
	for i, testCase := range testCoinCreationTransactions {
		b := encoding.Marshal(testCase)
		if len(b) == 0 {
			t.Error(i, "Binary-marshal output is empty")
		}
		var cctx CoinCreationTransaction
		err := encoding.Unmarshal(b, &cctx)
		if err != nil {
			t.Error(i, "failed to Binary-unmarshal tx", err)
			continue
		}
		testCompareTwoCoinCreationTransactions(t, i, cctx, testCase)
	}
}

// tests if we can unmarshal a JSON example inspired by
// originally proposed structure given in spec at
// https://github.com/threefoldfoundation/tfchain/issues/155#issuecomment-408029100
func TestJSONUnmarshalSpecCoinCreationTransactionExample(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintCondition: types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))),
	})
	defer types.RegisterTransactionVersion(TransactionVersionCoinCreation, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(`
{
	"version": 129,
	"data": {
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
					PublicKey: types.SiaPublicKey{
						Algorithm: types.SignatureEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.SiaPublicKey{
						Algorithm: types.SignatureEd25519,
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

func testCompareTwoCoinCreationTransactions(t *testing.T, i int, a, b CoinCreationTransaction) {
	// compare mint fulfillment
	if !a.MintFulfillment.Equal(b.MintFulfillment) {
		t.Error(i, "munt fulfillment not equal")
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
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.SiaPublicKey{
				Algorithm: types.SignatureEd25519,
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
			config.GetTestnetGenesis().MinimumTransactionFee.Mul64(2),
		},
	},
	// a more complex Coin Creation Transaction
	{
		MintFulfillment: types.NewFulfillment(&types.MultiSignatureFulfillment{
			Pairs: []types.PublicKeySignaturePair{
				{
					PublicKey: types.SiaPublicKey{
						Algorithm: types.SignatureEd25519,
						Key:       hbs("def123def123def123def123def123def123def123def123def123def123def1"),
					},
					Signature: hbs("ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef12345ef"),
				},
				{
					PublicKey: types.SiaPublicKey{
						Algorithm: types.SignatureEd25519,
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

// utility funcs
func hbs(str string) []byte { // hexStr -> byte slice
	bs, _ := hex.DecodeString(str)
	return bs
}
func hs(str string) (hash crypto.Hash) { // hbs -> crypto.Hash
	copy(hash[:], hbs(str))
	return
}
