package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/threefoldfoundation/tfchain/pkg/encoding"
)

var exampleBotNames = []string{
	"aaaaa",
	"aaaaa.bbbbb",
	"aaaaa.aaaaa",
	"aaaaa.bbbbb.ccccc.ddddd.eeeee.fffff.ggggg.hhhhhh.jjjjj.kkkkkkkk",
	"threefold.token",
	"trading.botzone",
}

func TestNewBotNames(t *testing.T) {
	for idx, example := range exampleBotNames {
		_, err := NewBotName(example)
		if err != nil {
			t.Error(idx, example, err)
		}
	}
}

func TestBotNamesLoadStringString(t *testing.T) {
	for idx, example := range exampleBotNames {
		var bn BotName
		err := bn.LoadString(example)
		if err != nil {
			t.Error(idx, example, err)
			continue
		}
		str := bn.String()
		if example != str {
			t.Error(idx, example, "!=", str)
		}
	}
}

func TestBotNameBinaryMarshalUnmarshal(t *testing.T) {
	for idx, example := range exampleBotNames {
		var bn BotName
		err := bn.LoadString(example)
		if err != nil {
			t.Error(idx, example, err)
			continue
		}

		// start binary marshal
		b := encoding.Marshal(bn)
		if len(b) == 0 {
			t.Error(idx, "encoding.Marshal=><nil>", example)
			continue
		}
		err = encoding.Unmarshal(b, &bn)
		if err != nil {
			t.Error(idx, "encoding.Unmarshal", example, err)
			continue
		}
		// end binary marshal

		str := bn.String()
		if example != str {
			t.Error(idx, example, "!=", str)
		}
	}
}

var (
	ExampleBotJSONRecord = []byte(`{
	"id": 42,
	"names": [
		"threefold.token",
		"trading.botzone"
	],
	"addresses": [
		"network.address.com",
		"83.200.201.201",
		"2001:db8:85a3::8a2e:370:7334"
	],
	"publickey": "ed25519:28c1edd4c35f662cccfa7fc02194959d75855c02d342c1131b110c9e96764d9b",
	"expiration": 1538484360
}
`)
)

func TestExampleBotRecordJSON(t *testing.T) {
	var record BotRecord
	err := json.Unmarshal(ExampleBotJSONRecord, &record)
	if err != nil {
		t.Fatal("json.Unmarshal", err)
	}
	b, err := json.Marshal(record)
	if err != nil {
		t.Fatal("json.Marshal", err)
	}
	result := string(b)
	buffer := bytes.NewBuffer(nil)
	err = json.Compact(buffer, ExampleBotJSONRecord)
	if err != nil {
		t.Fatal("json.Compact", err)
	}
	expected := buffer.String()
	if expected != result {
		t.Fatal("unexpected result:", expected, "!=", result)
	}
}

func TestExampleBotRecordBinarySia(t *testing.T) {
	var record BotRecord
	err := json.Unmarshal(ExampleBotJSONRecord, &record)
	if err != nil {
		t.Fatal("json.Unmarshal", err)
	}

	// sia marshal start
	b := encoding.Marshal(record)
	if len(b) == 0 {
		t.Fatal("encoding.Marshal: <nil>")
	}
	err = encoding.Unmarshal(b, &record)
	if err != nil {
		t.Fatal("encoding.Unmarshal", err)
	}
	// sia marshal end

	b, err = json.Marshal(record)
	if err != nil {
		t.Fatal("json.Marshal", err)
	}
	result := string(b)
	buffer := bytes.NewBuffer(nil)
	err = json.Compact(buffer, ExampleBotJSONRecord)
	if err != nil {
		t.Fatal("json.Compact", err)
	}
	expected := buffer.String()
	if expected != result {
		t.Fatal("unexpected result:", expected, "!=", result)
	}
}

// a total of 42 bytes, for the most minimalistic record in the current context
const minimalHexEncodedBinaryBotRecord = `00000000` + // first bot, index 0
	`10` + // 16 => 0 names and 1 addr
	`117F000001` + // IPv4 => 127.0.0.1
	`00404683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d` + // ed25519 pub key
	`7AF905` // some data

func TestMinimalHexEncodedBinaryBotRecord(t *testing.T) {
	b, err := hex.DecodeString(minimalHexEncodedBinaryBotRecord)
	if err != nil {
		t.Fatal(err)
	}
	var record BotRecord
	err = encoding.Unmarshal(b, &record)
	if err != nil {
		t.Fatal(err)
	}
	if record.ID != 0 {
		t.Error("unexpected ID", record.ID)
	}
	if len(record.Names) != 0 {
		t.Error("unexpected names", record.Names)
	}
	if len(record.Addresses) != 1 {
		t.Error("unexpected addresses", record.Addresses)
	}
	if str := record.Addresses[0].String(); str != "127.0.0.1" {
		t.Error("unexpected address", str, "!= 127.0.0.1")
	}
	if str := record.PublicKey.String(); str != "ed25519:4683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d" {
		t.Error("unexpected public key", str, "!= ed25519:4683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d")
	}
	if record.Expiration != 1538492760 {
		t.Error("unexpected expiration time", record.Expiration, "!=", 1538492760)
	}
}
