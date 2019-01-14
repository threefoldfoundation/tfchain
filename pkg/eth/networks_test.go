package eth

import (
	"testing"
)

//TestGetEthNetworkConfiguration is not an extensive test,as it would just duplicate  the constans defined in the networks file.
func TestGetEthNetworkConfiguration(t *testing.T) {
	testcases := []struct {
		networkname string
		networkID   uint64
	}{
		{"rinkeby", 4},
		{"main", 1},
		{"ropsten", 3},
	}
	for _, testcase := range testcases {
		conf, err := GetEthNetworkConfiguration(testcase.networkname)
		if err != nil {
			t.Error("An existing networkname should not return an error:", err)
		}
		if testcase.networkname != conf.NetworkName {
			t.Error(testcase.networkname, "!=", conf.NetworkName)
		}
		if testcase.networkID != conf.NetworkID {
			t.Error(testcase.networkID, "!=", conf.NetworkID)
		}
		bootnodes, err := conf.GetBootnodes()
		if err != nil {
			t.Error("error while getting BootNodes:", err)
		}
		if len(bootnodes) == 0 {
			t.Error("unexpected empty bootnodes list")
		}
	}
	//Unexisting networkname should return an error
	_, err := GetEthNetworkConfiguration("foo")
	if err == nil {
		t.Error("An unexisting networkname should return an error")
	}
}

func TestGetEthNetworkConfigurationInvalid(t *testing.T) {
	cfg := NetworkConfiguration{bootnodes: []string{"foo"}}
	nodes, err := cfg.GetBootnodes()
	if err == nil {
		t.Fatal("expected error, but received", nodes)
	}
}
