package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
)

// Test to ensure that the initial bot reg tx binary format is no longer valid when decoding.
func TestOutdatedBotRegisterationTransactionBinaryFormat(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	b, err := hex.DecodeString(`90e1113833626f742e7a6169626f6e2e62651a746633626f742e7a6169626f6e040000000000000005f5e1000201c4de6f0b9dbac6e9a398c313378733b98da930ab4accb403efc24c5267bb1a0180000000000000006564323535313900000000000000000020000000000000009e095c02584a5b042dfcf679837c88be924c40c95f173fe24d96852f6fd8c1934000000000000000b0f127cd3d85bf81fe354738deeaad6f8e570ecdb84a4ad435d24e4340718668fe22d0cdf3779e82257a85de18499755a35ac104240d47523940d244fc43140605000000000000005c98c4690001210000000000000001e9e4ab0970a899d02588d002cecb67ff942c81737501ddefe8aaf14bbdf722790172ebed8fd8b75fce87485ebe7184cf28b838d9e9ff55bbb23b8508f60fdede9edd0388c935f463a4d58f0c2e961aef54f69ca09c2ecff3a962bfebbae19648bf813cfd7b6830260ba9d464b8bbe11a064426bcd5f93675c72c67ddfedf15c404`)
	if err != nil {
		t.Fatal(err)
	}
	var tx types.Transaction
	err = siabin.Unmarshal(b, &tx)
	if err == nil {
		t.Fatal("expected error, but managed to unmarshal Tx:", tx)
	}
	msg := err.Error()
	if !strings.Contains(msg, "could not decode type types.Transaction") {
		t.Fatal("unexpected error, expected type-decode error, but received:", err)
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
	b, err := rivbin.Marshal(tx)
	if err != nil {
		t.Error(err)
	}

	// go to 3bot Tx and back
	botRegistrationTx, err := BotRegistrationTransactionFromTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}
	oTx := botRegistrationTx.Transaction(config.GetCurrencyUnits().OneCoin)
	oID := oTx.ID()
	oB, err := rivbin.Marshal(oTx)
	if err != nil {
		t.Error(err)
	}
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
		uh, err := types.NewPubKeyUnlockHash(cryptoKeyPair.PublicKey)
		if err != nil {
			return err
		}
		if condition.UnlockHash().Cmp(uh) != 0 {
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
		uh, err := types.NewPubKeyUnlockHash(cryptoKeyPair.PublicKey)
		if err != nil {
			return err
		}
		if condition.UnlockHash().Cmp(uh) != 0 {
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
		uh, err := types.NewPubKeyUnlockHash(cryptoKeyPair.PublicKey)
		if err != nil {
			return err
		}
		if condition.UnlockHash().Cmp(uh) != 0 {
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
