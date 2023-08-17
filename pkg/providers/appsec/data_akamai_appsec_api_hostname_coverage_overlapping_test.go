package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiApiHostnameCoverageOverlapping_data_basic(t *testing.T) {
	t.Run("match by ApiHostnameCoverageOverlapping ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getOverlapResponse := appsec.GetApiHostnameCoverageOverlappingResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSApiHostnameCoverageOverlapping/ApiHostnameCoverageOverlapping.json"), &getOverlapResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetApiHostnameCoverageOverlapping",
			mock.Anything,
			appsec.GetApiHostnameCoverageOverlappingRequest{ConfigID: 43253, Version: 7, Hostname: "example.com"},
		).Return(&getOverlapResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSApiHostnameCoverageOverlapping/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_hostname_coverage_overlapping.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
