package cli

import (
	"testing"

	"github.com/threefoldtech/rivine/modules"
)

func TestNetAddressArrayFlag(t *testing.T) {
	var s []modules.NetAddress
	flag := netAddressArray{array: &s}

	str := flag.String()
	if len(str) != 0 {
		t.Fatal("unexpected stringified NetAddressArrayFlag: ", str)
	}

	flag.Set("127.0.0.1:23112")
	expectedStr := "127.0.0.1:23112"
	str = flag.String()
	if expectedStr != str {
		t.Fatal("stringified NetAddressArrayFlag unexpected:", expectedStr, "!=", str)
	}

	flag.Set("localhost:23122")
	expectedStr = "127.0.0.1:23112,localhost:23122"
	str = flag.String()
	if expectedStr != str {
		t.Fatal("stringified NetAddressArrayFlag unexpected:", expectedStr, "!=", str)
	}

	expectedStrs := []modules.NetAddress{"127.0.0.1:23112", "localhost:23122"}
	if len(s) != len(expectedStrs) {
		t.Fatal("unexpected NetAddressArrayFlag result:", s)
	}
	for i, addr := range s {
		if addr != expectedStrs[i] {
			t.Errorf("unexpected NetAddressArrayFlag result value #%d: %s != %s", i+1, addr, expectedStrs[i])
		}
	}
}
