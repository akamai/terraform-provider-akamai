package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiThreatIntel_data_basic(t *testing.T) {
	t.Run("match by ThreatIntel ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getThreatIntelResponse := appsec.GetThreatIntelResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSThreatIntel/ThreatIntel.json"), &getThreatIntelResponse)
		require.NoError(t, err)

		client.On("GetThreatIntel",
			mock.Anything,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getThreatIntelResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSThreatIntel/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_threat_intel.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiThreatIntel_data_error_retrieving_threat_intel(t *testing.T) {
	t.Run("match by ThreatIntel ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		threatIntelResponse := appsec.GetThreatIntelResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSThreatIntel/ThreatIntel.json"), &threatIntelResponse)
		require.NoError(t, err)

		client.On("GetThreatIntel",
			mock.Anything,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetThreatIntel failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSThreatIntel/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_threat_intel.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`GetThreatIntel failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
