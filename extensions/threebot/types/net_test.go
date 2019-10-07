package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
)

var exampleNetworkAddresses = []string{
	"0.0.0.1",
	"127.0.0.1",
	"network.address.com",
	"83.200.201.201",
	"2001:db8:85a3::8a2e:370:7334",
}

func TestNewNetworkAddress(t *testing.T) {
	for idx, example := range exampleNetworkAddresses {
		_, err := NewNetworkAddress(example)
		if err != nil {
			t.Error(idx, example, err)
		}
	}
}

var exampleInvalidNetworkAddresses = []string{
	"localhost",
	"foo",
	"",
}

func TestInvalidNetworkAddress(t *testing.T) {
	for idx, example := range exampleInvalidNetworkAddresses {
		ip, err := NewNetworkAddress(example)
		if err == nil {
			t.Error(idx, "parsed correctly, while an error was expected", ip, example)
		}
	}
}

func TestNetworkAddressLoadStringString(t *testing.T) {
	for idx, example := range exampleNetworkAddresses {
		var na NetworkAddress
		err := na.LoadString(example)
		if err != nil {
			t.Error(idx, example, err)
			continue
		}
		str := na.String()
		if example != str {
			t.Error(idx, example, "!=", str)
		}
	}
}

func TestNetworkAddressBinaryMarshalUnmarshal(t *testing.T) {
	for idx, example := range exampleNetworkAddresses {
		var na NetworkAddress
		err := na.LoadString(example)
		if err != nil {
			t.Error(idx, example, err)
			continue
		}

		// start binary marshal
		b, err := rivbin.Marshal(na)
		if err != nil {
			t.Error(err)
		}
		if len(b) == 0 {
			t.Error(idx, "rivbin.Marshal=><nil>", example)
			continue
		}
		err = rivbin.Unmarshal(b, &na)
		if err != nil {
			t.Error(idx, "rivbin.Unmarshal", example, err)
			continue
		}
		// end binary marshal

		str := na.String()
		if example != str {
			t.Error(idx, example, "!=", str)
		}
	}
}

func TestNetworkAddressSortedSet(t *testing.T) {
	var nass NetworkAddressSortedSet
	if s := nass.Len(); s != 0 {
		t.Fatal("unexpected set length:", s)
	}

	// adding addresses should preserve order in an efficient way
	stringsToAdd := []string{"127.0.0.1", "example.edu", "125.0.0.1", "2001:db8:85a3::8a2e:370:7334", "125.0.0.2", "123.0.0.3"}
	expectedStrings := [][]string{
		{"127.0.0.1"},
		{"example.edu", "127.0.0.1"},
		{"example.edu", "125.0.0.1", "127.0.0.1"},
		{"example.edu", "125.0.0.1", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"},
		{"example.edu", "125.0.0.1", "125.0.0.2", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"},
		{"example.edu", "123.0.0.3", "125.0.0.1", "125.0.0.2", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"},
	}
	for i, str := range stringsToAdd {
		// add element
		err := nass.AddAddress(mustNewNetworkAddress(t, str))
		if err != nil {
			t.Fatal(i, "error while adding string", str, err)
		}
		// ensure slice length is as expected
		if s, e := nass.Len(), len(expectedStrings[i]); s != e {
			t.Error(i, "unexpected set length:", s, "!=", e)
		}
		// ensure all is in order
		for idx, addr := range nass.slice {
			str := addr.String()
			if expectedStrings[i][idx] != str {
				t.Error(i, idx, "unexpected stringified addr", expectedStrings[i][idx], "!=", str)
			}
		}
	}

	// removing elements should preserve order as well
	stringsToRemove := []string{"125.0.0.1", "example.edu", "2001:db8:85a3::8a2e:370:7334", "127.0.0.1", "125.0.0.2", "123.0.0.3"}
	expectedStrings = [][]string{
		{"example.edu", "123.0.0.3", "125.0.0.2", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"},
		{"123.0.0.3", "125.0.0.2", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"},
		{"123.0.0.3", "125.0.0.2", "127.0.0.1"},
		{"123.0.0.3", "125.0.0.2"},
		{"123.0.0.3"},
		{},
	}
	for i, str := range stringsToRemove {
		// remove element
		err := nass.RemoveAddress(mustNewNetworkAddress(t, str))
		if err != nil {
			t.Fatal(i, "error while removing string", str, err)
		}
		// ensure slice length is as expected
		if s, e := nass.Len(), len(expectedStrings[i]); s != e {
			t.Error(i, "unexpected set length:", s, "!=", e)
		}
		// ensure all is in order
		for idx, addr := range nass.slice {
			str := addr.String()
			if expectedStrings[i][idx] != str {
				t.Error(i, idx, "unexpected stringified addr", expectedStrings[i][idx], "!=", str)
			}
		}
	}

	// test to ensure we do not panic AND return an error
	err := nass.RemoveAddress(mustNewNetworkAddress(t, "foo.bar"))
	if err == nil {
		t.Fatal("removing an address from an empty NetworkAddressSortedSet should return an error, but none was received")
	}

	if s := nass.Len(); s != 0 {
		t.Fatal("unexpected set length:", s)
	}
}

func TestNetworkAddressSortedSetJSON(t *testing.T) {
	const input = `["125.1.0.1","2001:db8:85a3::8a2e:370:7334","127.0.0.1","example.edu"]`
	var nass NetworkAddressSortedSet
	err := json.Unmarshal([]byte(input), &nass)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{"example.edu", "125.1.0.1", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"}
	// ensure slice length is as expected
	if s, e := nass.Len(), len(expectedStrings); s != e {
		t.Error("unexpected set length:", s, "!=", e)
	}
	// ensure all is in order
	for idx, addr := range nass.slice {
		str := addr.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified addr", expectedStrings[idx], "!=", str)
		}
	}
	// ensure JSON output is correct as well
	expectedJSON := "["
	for _, str := range expectedStrings {
		expectedJSON += fmt.Sprintf(`"%s",`, str)
	}
	expectedJSON = expectedJSON[:len(expectedJSON)-1]
	expectedJSON += "]"
	b, err := json.Marshal(nass)
	if err != nil {
		t.Fatal(err)
	}
	result := string(b)
	if expectedJSON != result {
		t.Fatal(expectedJSON, "!=", result)
	}
}

func TestNetworkAddressSortedSetSiaMarshaling(t *testing.T) {
	const input = `064220010db885a3000000008a2e03707334117f0000014c6e6574776f726b2e616464726573732e636f6d`
	b, err := hex.DecodeString(input)
	if err != nil {
		t.Fatal(err)
	}
	var nass NetworkAddressSortedSet
	err = rivbin.Unmarshal(b, &nass)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{"network.address.com", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"}
	// ensure slice length is as expected
	if s, e := nass.Len(), len(expectedStrings); s != e {
		t.Error("unexpected set length:", s, "!=", e)
	}
	// ensure all is in order
	for idx, addr := range nass.slice {
		str := addr.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified addr", expectedStrings[idx], "!=", str)
		}
	}
	const expectedHex = `064c6e6574776f726b2e616464726573732e636f6d117f0000014220010db885a3000000008a2e03707334`
	b, err = rivbin.Marshal(nass)
	if err != nil {
		t.Error(err)
	}
	result := hex.EncodeToString(b)
	if expectedHex != result {
		t.Fatal(expectedHex, "!=", result)
	}
}

func TestNetworkAddressSortedSetBinaryEncoding(t *testing.T) {
	const input = `117f0000014220010db885a3000000008a2e03707334117f0000014c6e6574776f726b2e616464726573732e636f6d`
	b, err := hex.DecodeString(input)
	if err != nil {
		t.Fatal(err)
	}
	var nass NetworkAddressSortedSet
	r := bytes.NewReader(b)
	// first just read one
	err = nass.BinaryDecode(r, 1)
	if err != nil {
		t.Fatal(err)
	}
	if s := nass.Len(); s != 1 {
		t.Fatal("unexpected length", s)
	}
	str := nass.slice[0].String()
	if str != "127.0.0.1" {
		t.Fatal("unexpected stringified address: ", str)
	}
	// read all of them (except first one of course, that one is gone)
	err = nass.BinaryDecode(r, 3)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{"network.address.com", "127.0.0.1", "2001:db8:85a3::8a2e:370:7334"}
	// ensure slice length is as expected
	if s, e := nass.Len(), len(expectedStrings); s != e {
		t.Error("unexpected set length:", s, "!=", e)
	}
	// ensure all is in order
	for idx, addr := range nass.slice {
		str := addr.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified addr", expectedStrings[idx], "!=", str)
		}
	}
	const expectedHex = `4c6e6574776f726b2e616464726573732e636f6d117f0000014220010db885a3000000008a2e03707334`
	buf := bytes.NewBuffer(nil)
	_, err = nass.BinaryEncode(buf)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(buf.Bytes())
	if expectedHex != result {
		t.Fatal(expectedHex, "!=", result)
	}
}

func TestNetworkAddressSliceSort(t *testing.T) {
	var slice networkAddressSlice
	if s := slice.Len(); s != 0 {
		t.Fatal("unexpected slice length:", s)
	}

	slice = append(slice,
		mustNewNetworkAddress(t, "127.0.0.1"),
		mustNewNetworkAddress(t, "example.edu"),
		mustNewNetworkAddress(t, "125.0.0.1"),
		mustNewNetworkAddress(t, "125.0.0.2"),
		mustNewNetworkAddress(t, "123.0.0.3"))
	if s := slice.Len(); s != 5 {
		t.Fatal("unexpected slice length:", s)
	}

	sort.Sort(slice)
	expectedStrings := []string{"example.edu", "123.0.0.3", "125.0.0.1", "125.0.0.2", "127.0.0.1"}
	for idx, addr := range slice {
		str := addr.String()
		if expectedStrings[idx] != str {
			t.Error(idx, "unexpected stringified addr", expectedStrings[idx], "!=", str)
		}
	}
}

func mustNewNetworkAddress(t *testing.T, str string) NetworkAddress {
	t.Helper()
	addr, err := NewNetworkAddress(str)
	if err != nil {
		t.Fatal(err, str)
	}
	return addr
}
