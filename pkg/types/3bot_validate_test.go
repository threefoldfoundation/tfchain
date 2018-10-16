package types

import (
	"encoding/hex"
	"testing"

	"github.com/rivine/rivine/types"
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
	exampleBotTransactionPublicKey            = `ed25519:7bcf68f4b120aee2ad4ed4b4f69f0239137aa8c6b9c8ef2f6b4abd6b1a56f48c`
	exampleBotTransactionSignature            = `8ef94865b9e29ef0bf36221cd180d8f3783a65b65ecbeec6e79d162f7ec946ead3cce54adc9bfd47003ff30029407611ba3e8c3ebd468d1737d235204d6f7f07`
	exampleUnsignedBotTransactionInJSONFormat = `{"version":144,"data":{"addresses":null,"names":["thisis.bot02"],"nrofmonths":1,"txfee":"1000000000","coininputs":[{"parentid":"29ccd3551afce39581d826f8614df1f54b8796137e115b613de6376b6cdb3bff","fulfillment":{"type":1,"data":{"publickey":"` +
		exampleBotTransactionPublicKey + `","signature":""}}}],"refundcoinoutput":{"value":"399000000000","condition":{"type":1,"data":{"unlockhash":"01117e4beb05cd067a72ec474256f455545f5fe4fdd188b46228636982bd9935cc44f923c39928"}}},"identification":{"publickey":"ed25519:8ebcb01222f7c72d78f4a282fdc1ca1ec66be354690ec6f5dc308904164b8495","signature":"9c09865af6cdf55d97aea015810fe90deaaaadc31b0c988e8b21102bfd94661d540e7133af9989e55f3f1908bd6a0ab051bfef83e9f7fb7cbe7eeaef5d0dd50f"}}}`
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
