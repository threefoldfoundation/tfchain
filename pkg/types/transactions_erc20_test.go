package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
)

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
		"bridgefee": "50000000000",
		"blockid": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"txid": "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
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

	const hexEncodedExample = `d101f68299b26a89efdb4351a61c3a062321d23edbc1399c8499947c1313375609ad0a174876e800083b9aca000a0ba43b74000123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdefabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789`

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
