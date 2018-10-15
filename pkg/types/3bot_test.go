package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/threefoldfoundation/tfchain/pkg/encoding"
)

func TestBotIDLoadEmptyString(t *testing.T) {
	var id BotID
	err := id.LoadString("")
	if err == nil {
		t.Fatal("expected error for loading empty string, but received none")
	}
}

func TestBotIDLoadInvalidStrings(t *testing.T) {
	testCases := []string{
		"0", // has to be at least 1
		"a",
		"foo",
		"0f",
	}
	for idx, str := range testCases {
		var id BotID
		err := id.LoadString(str)
		if err == nil {
			t.Errorf("%d — expected error for loading invalid string %s, but received none", idx, str)
		}
	}
}

func TestBotIDLoadStrings(t *testing.T) {
	testCases := []string{
		"1",
		"42",
		fmt.Sprintf("%d", MinBotID),
		fmt.Sprintf("%d", MaxBotID),
		fmt.Sprintf("%d", MaxBotID-MinBotID),
	}
	for idx, str := range testCases {
		var id BotID
		err := id.LoadString(str)
		if err != nil {
			t.Errorf("%d — unexpected error for loading valid string %s: %v", idx, str, err)
		}
	}
}

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
	"addresses": [
		"network.address.com",
		"83.200.201.201",
		"2001:db8:85a3::8a2e:370:7334"
	],
	"names": [
		"threefold.token",
		"trading.botzone"
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

// a total of 46 bytes, for the most minimalistic record in the current context
const minimalHexEncodedBinaryBotRecord = `00000000` + // first bot, index 0
	`01` + // 1 => 0 names and 1 addr
	`117F000001` + // IPv4 => 127.0.0.1
	`004683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d` + // ed25519 pub key
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
	if record.Names.Len() != 0 {
		t.Error("unexpected names", record.Names)
	}
	if record.Addresses.Len() != 1 {
		t.Error("unexpected addresses", record.Addresses)
	}
	if str := record.Addresses.slice[0].String(); str != "127.0.0.1" {
		t.Error("unexpected address", str, "!= 127.0.0.1")
	}
	if str := record.PublicKey.String(); str != "ed25519:4683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d" {
		t.Error("unexpected public key", str, "!= ed25519:4683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d")
	}
	if record.Expiration != 1538492760 {
		t.Error("unexpected expiration time", record.Expiration, "!=", 1538492760)
	}
}

// TODO:
// TestAddNames
// TestRemoveNames
// TestAddNetworkAddresses
// TestRemoveNetworkAddresses
// TestExtendExpirationDate

func TestBotNameSortedSet(t *testing.T) {
	var bnss BotNameSortedSet
	if s := bnss.Len(); s != 0 {
		t.Fatal("unexpected set length:", s)
	}

	// adding names should preserve order in an efficient way
	stringsToAdd := []string{"aaaaa.aaaaa", "ccccc", "aaaaa", "bbbbb", "ccccc.ccccc", "aaaaa.ccccc.ddddd"}
	expectedStrings := [][]string{
		{"aaaaa.aaaaa"},
		{"aaaaa.aaaaa", "ccccc"},
		{"aaaaa", "aaaaa.aaaaa", "ccccc"},
		{"aaaaa", "aaaaa.aaaaa", "bbbbb", "ccccc"},
		{"aaaaa", "aaaaa.aaaaa", "bbbbb", "ccccc", "ccccc.ccccc"},
		{"aaaaa", "aaaaa.aaaaa", "aaaaa.ccccc.ddddd", "bbbbb", "ccccc", "ccccc.ccccc"},
	}
	for i, str := range stringsToAdd {
		// add element
		err := bnss.AddName(mustNewBotName(t, str))
		if err != nil {
			t.Fatal(i, "error while adding string", str, err)
		}
		// ensure slice length is as expected
		if s, e := bnss.Len(), len(expectedStrings[i]); s != e {
			t.Error(i, "unexpected set length:", s, "!=", e)
		}
		// ensure all is in order
		for idx, name := range bnss.slice {
			str := name.String()
			if expectedStrings[i][idx] != str {
				t.Error(i, idx, "unexpected stringified name", expectedStrings[i][idx], "!=", str)
			}
		}
	}

	// removing names should preserve order as well
	stringsToRemove := []string{"ccccc.ccccc", "aaaaa", "ccccc", "aaaaa.ccccc.ddddd", "bbbbb", "aaaaa.aaaaa"}
	expectedStrings = [][]string{
		{"aaaaa", "aaaaa.aaaaa", "aaaaa.ccccc.ddddd", "bbbbb", "ccccc"},
		{"aaaaa.aaaaa", "aaaaa.ccccc.ddddd", "bbbbb", "ccccc"},
		{"aaaaa.aaaaa", "aaaaa.ccccc.ddddd", "bbbbb"},
		{"aaaaa.aaaaa", "bbbbb"},
		{"aaaaa.aaaaa"},
		{},
	}
	for i, str := range stringsToRemove {
		// remove bot name
		err := bnss.RemoveName(mustNewBotName(t, str))
		if err != nil {
			t.Fatal(i, "error while removing string", str, err)
		}
		// ensure slice length is as expected
		if s, e := bnss.Len(), len(expectedStrings[i]); s != e {
			t.Error(i, "unexpected set length:", s, "!=", e)
		}
		// ensure all is in order
		for idx, name := range bnss.slice {
			str := name.String()
			if expectedStrings[i][idx] != str {
				t.Error(i, idx, "unexpected stringified name", expectedStrings[i][idx], "!=", str)
			}
		}
	}

	// test to ensure we do not panic AND return an error
	err := bnss.RemoveName(mustNewBotName(t, "zzzzz.xxxxx"))
	if err == nil {
		t.Fatal("removing a bot name from an empty BotNameSortedSet should return an error, but none was received")
	}

	if s := bnss.Len(); s != 0 {
		t.Fatal("unexpected set length:", s)
	}
}

func TestBotNamesSortedSetJSON(t *testing.T) {
	const input = `["aaaaa.ddddd","ccccc","aaaae","bbbbb","aaaaa"]`
	var bnss BotNameSortedSet
	err := json.Unmarshal([]byte(input), &bnss)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{"aaaaa", "aaaaa.ddddd", "aaaae", "bbbbb", "ccccc"}
	// ensure slice length is as expected
	if s, e := bnss.Len(), len(expectedStrings); s != e {
		t.Error("unexpected set length:", s, "!=", e)
	}
	// ensure all is in order
	for idx, name := range bnss.slice {
		str := name.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified name", expectedStrings[idx], "!=", str)
		}
	}
	// ensure JSON output is correct as well
	expectedJSON := "["
	for _, str := range expectedStrings {
		expectedJSON += fmt.Sprintf(`"%s",`, str)
	}
	expectedJSON = expectedJSON[:len(expectedJSON)-1]
	expectedJSON += "]"
	b, err := json.Marshal(bnss)
	if err != nil {
		t.Fatal(err)
	}
	result := string(b)
	if expectedJSON != result {
		t.Fatal(expectedJSON, "!=", result)
	}
}

func TestBotNameSortedSetSiaMarshaling(t *testing.T) {
	const input = `061661616161612e61616161611e7468726565666f6c642e746f6b656e0a6161616161`
	b, err := hex.DecodeString(input)
	if err != nil {
		t.Fatal(err)
	}
	var bnss BotNameSortedSet
	err = encoding.Unmarshal(b, &bnss)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{"aaaaa", "aaaaa.aaaaa", "threefold.token"}
	// ensure slice length is as expected
	if s, e := bnss.Len(), len(expectedStrings); s != e {
		t.Error("unexpected set length:", s, "!=", e)
	}
	// ensure all is in order
	for idx, name := range bnss.slice {
		str := name.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified name", expectedStrings[idx], "!=", str)
		}
	}
	const expectedHex = `060a61616161611661616161612e61616161611e7468726565666f6c642e746f6b656e`
	b = encoding.Marshal(bnss)
	result := hex.EncodeToString(b)
	if expectedHex != result {
		t.Fatal(expectedHex, "!=", result)
	}
}

func TestBotNameSortedSetBinaryEncoding(t *testing.T) {
	const input = `0a61616161611661616161612e61616161611e7468726565666f6c642e746f6b656e0a6161616161`
	b, err := hex.DecodeString(input)
	if err != nil {
		t.Fatal(err)
	}
	var bnss BotNameSortedSet
	r := bytes.NewReader(b)
	// first just read one
	err = bnss.BinaryDecode(r, 1)
	if err != nil {
		t.Fatal(err)
	}
	if s := bnss.Len(); s != 1 {
		t.Fatal("unexpected length", s)
	}
	str := bnss.slice[0].String()
	if str != "aaaaa" {
		t.Fatal("unexpected stringified name: ", str)
	}
	// read all of them (except first one of course, that one is gone)
	err = bnss.BinaryDecode(r, 3)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{"aaaaa", "aaaaa.aaaaa", "threefold.token"}
	// ensure slice length is as expected
	if s, e := bnss.Len(), len(expectedStrings); s != e {
		t.Error("unexpected set length:", s, "!=", e)
	}
	// ensure all is in order
	for idx, name := range bnss.slice {
		str := name.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified name", expectedStrings[idx], "!=", str)
		}
	}
	const expectedHex = `0a61616161611661616161612e61616161611e7468726565666f6c642e746f6b656e`
	buf := bytes.NewBuffer(nil)
	_, err = bnss.BinaryEncode(buf)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(buf.Bytes())
	if expectedHex != result {
		t.Fatal(expectedHex, "!=", result)
	}
}

func TestBotNameSortedSetDifference(t *testing.T) {
	a := BotNameSortedSet{slice: []BotName{
		mustNewBotName(t, "aaaaa"), mustNewBotName(t, "bbbbb"),
		mustNewBotName(t, "ddddd"), mustNewBotName(t, "fffff"),
	}}
	b := BotNameSortedSet{slice: []BotName{
		mustNewBotName(t, "bbbbb"), mustNewBotName(t, "ccccc"),
		mustNewBotName(t, "eeeee"), mustNewBotName(t, "ggggg"),
		mustNewBotName(t, "hhhhh"),
	}}
	result := a.Difference(b)
	expected := []BotName{mustNewBotName(t, "aaaaa"), mustNewBotName(t, "ddddd"), mustNewBotName(t, "fffff")}
	if !reflect.DeepEqual(expected, result) {
		t.Error("unexpected result for 'a \\ b':", expected, "!=", result)
	}
	result = b.Difference(a)
	expected = []BotName{mustNewBotName(t, "ccccc"), mustNewBotName(t, "eeeee"), mustNewBotName(t, "ggggg"), mustNewBotName(t, "hhhhh")}
	if !reflect.DeepEqual(expected, result) {
		t.Error("unexpected result for 'b \\ a':", expected, "!=", result)
	}
}

func TestBotNameSliceSort(t *testing.T) {
	var slice botNameSlice
	if s := slice.Len(); s != 0 {
		t.Fatal("unexpected slice length:", s)
	}

	slice = append(slice,
		mustNewBotName(t, "bbbbb.aaaaa"),
		mustNewBotName(t, "aaaaa.bbbbb"),
		mustNewBotName(t, "acaaa.bbbbb"),
		mustNewBotName(t, "deaaa.bbbbb"),
		mustNewBotName(t, "aaaaa.abbbb"),
	)
	if s := slice.Len(); s != 5 {
		t.Fatal("unexpected slice length:", s)
	}

	sort.Sort(slice)
	expectedStrings := []string{"aaaaa.abbbb", "aaaaa.bbbbb", "acaaa.bbbbb", "bbbbb.aaaaa", "deaaa.bbbbb"}
	for idx, name := range slice {
		str := name.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified name", expectedStrings[idx], "!=", str)
		}
	}
}

func mustNewBotName(t *testing.T, str string) BotName {
	t.Helper()
	name, err := NewBotName(str)
	if err != nil {
		t.Fatal(err)
	}
	return name
}
