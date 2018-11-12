package types

import (
	"encoding/hex"
	"testing"

	"github.com/threefoldtech/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/config"
)

func TestValidateUniquenessOfNetworkAddresses_Correct(t *testing.T) {
	testCases := [][]NetworkAddress{
		{},
		{mustNewNetworkAddress(t, "mybot.io")},
		{mustNewNetworkAddress(t, "127.0.0.1")},
		{mustNewNetworkAddress(t, "127.0.0.1"), mustNewNetworkAddress(t, "mybot.io")},
		{mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "127.0.0.1")},
		{mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "127.0.0.1"), mustNewNetworkAddress(t, "100.1.2.3")},
		{
			mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "127.0.0.1"),
			mustNewNetworkAddress(t, "100.1.2.3"), mustNewNetworkAddress(t, "example.org"),
		}, {
			mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "127.0.0.1"),
			mustNewNetworkAddress(t, "2001:db8:85a3::8a2e:370:7334"),
			mustNewNetworkAddress(t, "100.1.2.3"), mustNewNetworkAddress(t, "example.org"),
		},
	}
	for idx, testCase := range testCases {
		err := validateUniquenessOfNetworkAddresses(testCase)
		if err != nil {
			t.Error(idx, "unexpected error", err)
		}
	}
}

func TestValidateUniquenessOfNetworkAddresses_Error(t *testing.T) {
	testCases := [][]NetworkAddress{
		{mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "mybot.io")},
		{mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "127.0.0.1"), mustNewNetworkAddress(t, "mybot.io")},
		{mustNewNetworkAddress(t, "127.0.0.1"), mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "mybot.io")},
		{mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "mybot.io"), mustNewNetworkAddress(t, "127.0.0.1")},
		{mustNewNetworkAddress(t, "127.0.0.1"), mustNewNetworkAddress(t, "127.0.0.1")},
		{mustNewNetworkAddress(t, "127.0.0.1"), mustNewNetworkAddress(t, "2001:db8:85a3::8a2e:370:7334"), mustNewNetworkAddress(t, "127.0.0.1")},
		{mustNewNetworkAddress(t, "2001:db8:85a3::8a2e:370:7334"), mustNewNetworkAddress(t, "2001:db8:85a3::8a2e:370:7334")},
	}
	for idx, testCase := range testCases {
		err := validateUniquenessOfNetworkAddresses(testCase)
		if err == nil {
			t.Error(idx, "error expected but none was received", testCase)
		}
	}
}

func TestValidateUniquenessOfBotNames_Correct(t *testing.T) {
	testCases := [][]BotName{
		{},
		{mustNewBotName(t, "aaaaa.bbbbb")},
		{mustNewBotName(t, "ccccc")},
		{mustNewBotName(t, "ccccc"), mustNewBotName(t, "aaaaa.bbbbb")},
		{mustNewBotName(t, "aaaaa.bbbbb"), mustNewBotName(t, "ccccc")},
		{mustNewBotName(t, "aaaaa.bbbbb"), mustNewBotName(t, "ccccc"), mustNewBotName(t, "ddddd.eeeee.fffff")},
		{
			mustNewBotName(t, "aaaaa.bbbbb"), mustNewBotName(t, "ccccc"),
			mustNewBotName(t, "ddddd.eeeee.fffff"), mustNewBotName(t, "ggggg"),
		}, {
			mustNewBotName(t, "aaaaa.bbbbb"), mustNewBotName(t, "ccccc"),
			mustNewBotName(t, "hhhhh"),
			mustNewBotName(t, "ddddd.eeeee.fffff"), mustNewBotName(t, "ggggg"),
		},
	}
	for idx, testCase := range testCases {
		err := validateUniquenessOfBotNames(testCase)
		if err != nil {
			t.Error(idx, "unexpected error", err)
		}
	}
}

func TestValidateUniquenessOfBotNames_Error(t *testing.T) {
	testCases := [][]BotName{
		{mustNewBotName(t, "aaaaa"), mustNewBotName(t, "aaaaa")},
		{mustNewBotName(t, "aaaaa"), mustNewBotName(t, "bbbbb"), mustNewBotName(t, "aaaaa")},
		{mustNewBotName(t, "bbbbb"), mustNewBotName(t, "aaaaa"), mustNewBotName(t, "aaaaa")},
		{mustNewBotName(t, "aaaaa"), mustNewBotName(t, "aaaaa"), mustNewBotName(t, "bbbbb")},
		{mustNewBotName(t, "bbbbb"), mustNewBotName(t, "bbbbb")},
		{mustNewBotName(t, "bbbbb"), mustNewBotName(t, "ccccc"), mustNewBotName(t, "bbbbb")},
		{mustNewBotName(t, "ccccc"), mustNewBotName(t, "ccccc")},
	}
	for idx, testCase := range testCases {
		err := validateUniquenessOfBotNames(testCase)
		if err == nil {
			t.Error(idx, "error expected but none was received", testCase)
		}
	}
}

const (
	exampleBotTransactionPublicKey            = `ed25519:00bde9571b30e1742c41fcca8c730183402d967df5b17b5f4ced22c677806614`
	exampleBotTransactionSignature            = `12bb912737dbd572a5c6695537cbf9d72654264b8b98d2929f5b829abbc682749a3a93c83c545315f5c15ee895e136abc023bb58f691010899b7a1d9d222340f`
	exampleUnsignedBotTransactionInJSONFormat = `{"version":144,"data":{"addresses":["91.198.174.192","example.org"],"names":["chatbot.example"],"nrofmonths":1,"txfee":"1000000000","coininputs":[{"parentid":"a3c8f44d64c0636018a929d2caeec09fb9698bfdcbfa3a8225585a51e09ee563","fulfillment":{"type":1,"data":{"publickey":"ed25519:d285f92d6d449d9abb27f4c6cf82713cec0696d62b8c123f1627e054dc6d7780","signature":"78168863933e533c4686ad9749933a02db79c2dd49fc44e46984990e59df704c48e61b8ba845eb781367a55ea49d14ca51d4994315e451fd90f9a3760513bd0b"}}}],"refundcoinoutput":{"value":"99999899000000000","condition":{"type":1,"data":{"unlockhash":"01b49da2ff193f46ee0fc684d7a6121a8b8e324144dffc7327471a4da79f1730960edcb2ce737f"}}},"identification":{"publickey":"` +
		exampleBotTransactionPublicKey + `","signature":""}}}`
)

func TestValidateBotSignature_Correct(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry:              nil,
		RegistryPoolCondition: types.UnlockConditionProxy{},
		OneCoin:               config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(exampleUnsignedBotTransactionInJSONFormat))
	if err != nil {
		t.Fatal(err)
	}

	var publicKey PublicKey
	err = publicKey.LoadString(exampleBotTransactionPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	var signature types.ByteSlice
	err = signature.LoadString(exampleBotTransactionSignature)
	if err != nil {
		t.Fatal(err)
	}

	// taken from devnet
	err = validateBotSignature(tx, publicKey, signature, types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 784,
		BlockTime:   1539634805,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateBotSignature_Error(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry:              nil,
		RegistryPoolCondition: types.UnlockConditionProxy{},
		OneCoin:               config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(exampleUnsignedBotTransactionInJSONFormat))
	if err != nil {
		t.Fatal(err)
	}

	var publicKey PublicKey
	err = publicKey.LoadString(exampleBotTransactionPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	var signature types.ByteSlice
	err = signature.LoadString(exampleBotTransactionSignature)
	if err != nil {
		t.Fatal(err)
	}
	// change one letter in the signature
	if signature[0] == '0' {
		signature[0] = '1'
	} else {
		signature[0] = '0'
	}

	// taken from devnet
	err = validateBotSignature(tx, publicKey, signature, types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 784,
		BlockTime:   1539634805,
	})
	if err == nil {
		t.Fatal("validateBotSignature should be invalid, given signature is invalid, but doesn't seem to be invalid")
	}
}

func TestValidateBotSignature_InvalidPublicKey(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry:              nil,
		RegistryPoolCondition: types.UnlockConditionProxy{},
		OneCoin:               config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(exampleUnsignedBotTransactionInJSONFormat))
	if err != nil {
		t.Fatal(err)
	}

	// use invalid public key
	publicKey := PublicKey{
		Algorithm: 42,
	}
	publicKey.Key, _ = hex.DecodeString(`69f0239137aa8c6b9c8ef269f0239137aa8c6b9c8ef269f0239137aa8c6b9c8e`)

	var signature types.ByteSlice
	err = signature.LoadString(exampleBotTransactionSignature)
	if err != nil {
		t.Fatal(err)
	}

	// taken from devnet
	err = validateBotSignature(tx, publicKey, signature, types.ValidationContext{
		Confirmed:   true,
		BlockHeight: 784,
		BlockTime:   1539634805,
	})
	if err == nil {
		t.Fatal("validateBotSignature should be invalid, given signature is invalid, but doesn't seem to be invalid")
	}
}
