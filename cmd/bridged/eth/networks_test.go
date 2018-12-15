package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

		assert.Nil(t, err, "An existing networkname should not return an error")
		assert.Equal(t, testcase.networkname, conf.NetworkName)
		assert.Equal(t, testcase.networkID, conf.NetworkID)
		bootnodes, err := conf.GetBootnodes()
		assert.Nil(t, err)
		assert.NotEmpty(t, bootnodes)
	}
	//Unexisting networkname should return an error
	_, err := GetEthNetworkConfiguration("unexisting_blablabla")
	assert.Error(t, err, "An unexisting networkname should return an error")
}
