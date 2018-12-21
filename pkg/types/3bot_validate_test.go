package types

import (
	"encoding/hex"
	"testing"

	"github.com/threefoldfoundation/tfchain/pkg/config"
	"github.com/threefoldtech/rivine/types"
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
	exampleBotTransactionPublicKey            = `ed25519:29eb219e2a943325c2ce4bad26e464e24c9eed50f3b5acdd8a772e0947a0db4f`
	exampleBotTransactionSignature            = `73fd66f07ec9df121dcacc48ef00d332509210e689fc33d56a12e990c669bb86c59f33e276669ddaf187ac3238b623035ecb91a8edc7ffaf1131b802b6947807`
	exampleUnsignedBotTransactionInJSONFormat = `{"version":144,"data":{"addresses":["91.198.174.192","example.org"],"names":["chatbot.example"],"nrofmonths":1,"txfee":"1000000000","coininputs":[{"parentid":"73b691ee623eac8cadb20407aef944eda82cce3778840892f56ea1c78eb9cd60","fulfillment":{"type":1,"data":{"publickey":"ed25519:cb7db3934c904fb70e50c063fc54a20d7cad375e101a9ef21f7b7d1f7ad23cd8","signature":"2c46f922629032aab4053d57c46febce1ca7b5e597f81b55784ddad85da03c00c2db8aea0817b453c0e24134ba6eb5b738eb17f891eee6ebb65b890dd990e406"}}}],"refundcoinoutput":{"value":"99999565000000000","condition":{"type":1,"data":{"unlockhash":"01a006599af1155f43d687635e9680650003a6c506934996b90ae84d07648927414046f9f0e936"}}},"identification":{"publickey":"` +
		exampleBotTransactionPublicKey + `","signature":"` +
		exampleBotTransactionSignature + `"}}}`
)

func TestValidateBotSignature_Correct(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(exampleUnsignedBotTransactionInJSONFormat))
	if err != nil {
		t.Fatal(err)
	}

	var publicKey types.PublicKey
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
		BlockHeight: 575,
		BlockTime:   1545236350,
	}, BotSignatureSpecifierSender)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateBotSignature_Error(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(exampleUnsignedBotTransactionInJSONFormat))
	if err != nil {
		t.Fatal(err)
	}

	var publicKey types.PublicKey
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
		BlockHeight: 575,
		BlockTime:   1545236350,
	}, BotSignatureSpecifierSender)
	if err == nil {
		t.Fatal("validateBotSignature should be invalid, given signature is invalid, but doesn't seem to be invalid")
	}
}

func TestValidateBotSignature_InvalidPublicKey(t *testing.T) {
	// define tfchain-specific transaction versions
	types.RegisterTransactionVersion(TransactionVersionBotRegistration, BotRegistrationTransactionController{
		Registry: nil,
		OneCoin:  config.GetCurrencyUnits().OneCoin,
	})
	defer types.RegisterTransactionVersion(TransactionVersionBotRegistration, nil)

	var tx types.Transaction
	err := tx.UnmarshalJSON([]byte(exampleUnsignedBotTransactionInJSONFormat))
	if err != nil {
		t.Fatal(err)
	}

	// use invalid public key
	publicKey := types.PublicKey{
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
		BlockHeight: 575,
		BlockTime:   1545236350,
	}, BotSignatureSpecifierSender)
	if err == nil {
		t.Fatal("validateBotSignature should be invalid, given signature is invalid, but doesn't seem to be invalid")
	}
}
