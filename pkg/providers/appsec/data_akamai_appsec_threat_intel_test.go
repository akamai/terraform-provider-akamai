package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiThreatIntel_data_basic(t *testing.T) {
	t.Run("match by ThreatIntel ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getThreatIntelResponse := appsec.GetThreatIntelResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSThreatIntel/ThreatIntel.json"), &getThreatIntelResponse)
		require.NoError(t, err)

		client.On("GetThreatIntel",
			testutils.MockContext,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getThreatIntelResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSThreatIntel/match_by_id.tf"),
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
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		threatIntelResponse := appsec.GetThreatIntelResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSThreatIntel/ThreatIntel.json"), &threatIntelResponse)
		require.NoError(t, err)

		client.On("GetThreatIntel",
			testutils.MockContext,
			appsec.GetThreatIntelRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetThreatIntel failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSThreatIntel/match_by_id.tf"),
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
