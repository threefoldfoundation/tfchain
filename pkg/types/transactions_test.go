package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
		Nonce: RandomTransactionNonce(),
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

// tests the patch which fixes https://github.com/threefoldfoundation/tfchain/issues/164
func TestCoinCreationTransactionIDUniqueness(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintCondition: types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))),
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

const validDevnetJSONEncodedCoinCreationTx = `{
	"version": 129,
	"data": {
		"nonce": [51, 166, 67, 34, 32, 51, 73, 70],
		"mintfulfillment": {
			"type": 1,
			"data": {
				"publickey": "ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780",
				"signature": "a074b976556d6ea2e4ae8d51fbbb5ec99099f11918201abfa31cf80d415c8d5bdfda5a32d9cc167067b6b798e80c6c1a45f6fd9e0f01ac09053e767b15d31005"
			}
		},
		"coinoutputs": [{
			"value": "500000000000000",
			"condition": {
				"type": 1,
				"data": {
					"unlockhash": "01e78fd5af261e49643dba489b29566db53fa6e195fa0e6aad4430d4f06ce88b73e047fe6a0703"
				}
			}
		}],
		"minerfees": ["1000000000"],
		"arbitrarydata": "bW9uZXkgZnJvbSB0aGUgc2t5"
	}
}`

func TestSignCoinCreationTransactionExtension(t *testing.T) {
	mintCondition := types.NewCondition(types.NewUnlockHashCondition(
		unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f")))
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintCondition: mintCondition,
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
		switch mintCondition.ConditionType() {
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
			err := tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy) error {
				return fulfillment.Sign(types.FulfillmentSignContext{
					InputIndex:  0, // doesn't matter really for this extension
					Transaction: tx,
					Key:         key,
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

	mintCondition = types.NewCondition(types.NewMultiSignatureCondition(
		types.UnlockHashSlice{
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"),
			unlockHashFromHex("016438a548b6d377e87b08e8eae5ef641a4e70cc861b85b54b0921330e03084ffe0a8d9a38e3a8"),
		},
		2,
	))
	// overwrite coin tx multisig condition to be a multisig condition instead
	types.RegisterTransactionVersion(TransactionVersionCoinCreation, CoinCreationTransactionController{
		MintCondition: mintCondition,
	})

	// sign multisig condition, should fail as we didn't sign enough
	err = signTxAndValidate(testKeyPair{
		KeyPair: types.KeyPair{
			PublicKey: types.SiaPublicKey{
				Algorithm: types.SignatureEd25519,
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
			PublicKey: types.SiaPublicKey{
				Algorithm: types.SignatureEd25519,
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
		MintCondition: types.NewCondition(types.NewUnlockHashCondition(
			unlockHashFromHex("015a080a9259b9d4aaa550e2156f49b1a79a64c7ea463d810d4493e8242e6791584fbdac553e6f"))),
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
		err := tx.SignExtension(func(fulfillment *types.UnlockFulfillmentProxy, condition types.UnlockConditionProxy) error {
			return fulfillment.Sign(types.FulfillmentSignContext{
				InputIndex:  0, // doesn't matter really for this extension
				Transaction: tx,
				Key:         hsk("788c0aaeec8e0d916a712535826fa2d47d19fd7b341242f05de0d2e6e7e06104d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780"),
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
	tx.Extension = CoinCreationTransactionExtension{
		Nonce: RandomTransactionNonce(),
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of nil fulfillment")
	}
	tx.Extension = CoinCreationTransactionExtension{
		Nonce: RandomTransactionNonce(),
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
	}
	err = tx.ValidateTransaction(validationCtx, txValidationConstants)
	if err == nil {
		t.Fatal("succeeded to validate coin creation tx, " +
			"while it is supposed to fail because of wrong fulfillment")
	}
	tx.Extension = CoinCreationTransactionExtension{
		Nonce: RandomTransactionNonce(),
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.SiaPublicKey{
				Algorithm: types.SignatureEd25519,
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
	tx.Extension = CoinCreationTransactionExtension{
		Nonce: RandomTransactionNonce(),
		MintFulfillment: types.NewFulfillment(&types.SingleSignatureFulfillment{
			PublicKey: types.SiaPublicKey{
				Algorithm: types.SignatureEd25519,
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
