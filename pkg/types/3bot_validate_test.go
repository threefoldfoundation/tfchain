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
	exampleBotTransactionPublicKey            = `ed25519:880ee50bd7efa4c8b2b5949688a09818a652727fd3c0cb406013be442df68b34`
	exampleBotTransactionSignature            = `625c2db62790a2be025ba72356bb5f0539ada3d2feb923eaeda4aa798845dd71c08f0e669479087b3c59f828e3abe38f75690c443188f4cadcdad0b539e5dc0e`
	exampleUnsignedBotTransactionInJSONFormat = `{"version":144,"data":{"addresses":["91.198.174.192","example.org"],"names":["chatbot.example"],"nrofmonths":1,"txfee":"1000000000","coininputs":[{"parentid":"e6239feadc465055e17ab9a3111836e82ad35e7bb1559da3317e6f2cc624582c","fulfillment":{"type":1,"data":{"publickey":"ed25519:a271b9d4c1258f070e1e8d95250e6d29f683649829c2227564edd5ddeb75819d","signature":"d4e5d23929151fe511be963dde8b221f314a4decb9dc8ddcd34a0bc969a5d129dc32127d753054b80837192eead9a353bce4841de3d911e1de3d05ba8ae30102"}}}],"refundcoinoutput":{"value":"99999798000000000","condition":{"type":1,"data":{"unlockhash":"01972837ee396f22f96846a0c700f9cf7c8fa83ab4110da91a1c7d02f94f28ff03e45f1470df82"}}},"identification":{"publickey":"` +
		exampleBotTransactionPublicKey + `","signature":"` +
		exampleBotTransactionSignature + `"}}}`
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
