package property

import (
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/stretchr/testify/assert"
)

func TestNetworkAlias(t *testing.T) {
	t.Skip()
	tests := map[string]struct {
		hasNetwork  bool
		addNetwork  papi.ActivationNetwork
		networkTest papi.ActivationNetwork
		withError   error
	}{
		"ok production": {
			hasNetwork:  true,
			addNetwork:  papi.ActivationNetworkProduction,
			networkTest: papi.ActivationNetworkProduction,
		},
		"ok p": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkProduction,
			addNetwork:  "P",
		},
		"ok prod": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkProduction,
			addNetwork:  "PROD",
		},
		"ok stag": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkStaging,
			addNetwork:  "STAG",
		},
		"ok stage": {
			hasNetwork:  true,
			networkTest: papi.ActivationNetworkStaging,
			addNetwork:  "STAGE",
		},
		"ok default staging": {
			hasNetwork:  false,
			addNetwork:  papi.ActivationNetworkStaging,
			networkTest: papi.ActivationNetworkStaging,
		},
		"nok malformed input": {
			hasNetwork: true,
			addNetwork: "other",
			withError:  fmt.Errorf("network not recognized"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			net, err := NetworkAlias(string(test.addNetwork))
			resultNetwork := papi.ActivationNetwork(net)

			if test.withError != nil {
				assert.Error(t, test.withError, err)
			} else {
				assert.Equal(t, test.networkTest, resultNetwork)
				assert.NoError(t, err)
			}
		})
	}
}
