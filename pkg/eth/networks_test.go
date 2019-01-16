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
		bootnodes, err := conf.GetBootnodes(nil)
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
	nodes, err := cfg.GetBootnodes(nil)
	if err == nil {
		t.Fatal("expected error, but received", nodes)
	}
}
func TestGetBootnodes(t *testing.T) {
	conf, _ := GetEthNetworkConfiguration("ropsten")

	bootnodes, err := conf.GetBootnodes(nil)
	if err != nil {
		t.Error("error while getting BootNodes:", err)
	}
	if len(bootnodes) == 0 {
		t.Error("unexpected empty bootnodes list")
	}

	bootnodes, err = conf.GetBootnodes([]string{})
	if err != nil {
		t.Error("error while getting BootNodes:", err)
	}
	if len(bootnodes) == 0 {
		t.Error("unexpected empty bootnodes list")
	}

	bootnodes, err = conf.GetBootnodes([]string{"enode://2b4ae2d7ece11acffa1ce1ceac14e55e68f0c43a5f35b58e89f55d6de7c06ab98777b85d7f1f15eb1f6de0c39e3f35a3917bf647abbdc22f65ea2c73056162ca@127.0.0.1:3003"})
	if err != nil {
		t.Error("error while getting BootNodes:", err)
	}
	if len(bootnodes) != 1 {
		t.Error("unexpected empty bootnodes list")
	}
}
