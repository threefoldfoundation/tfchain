package types

import (
	"testing"

	"github.com/rivine/rivine/encoding"
)

var exampleNetworkAddresses = []string{
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
		b := encoding.Marshal(na)
		if len(b) == 0 {
			t.Error(idx, "encoding.Marshal=><nil>", example)
			continue
		}
		err = encoding.Unmarshal(b, &na)
		if err != nil {
			t.Error(idx, "encoding.Unmarshal", example, err)
			continue
		}
		// end binary marshal

		str := na.String()
		if example != str {
			t.Error(idx, example, "!=", str)
		}
	}
}
