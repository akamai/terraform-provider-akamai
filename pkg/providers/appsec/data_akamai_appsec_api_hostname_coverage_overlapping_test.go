package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiApiHostnameCoverageOverlapping_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by ApiHostnameCoverageOverlapping ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getOverlapResponse := appsec.GetApiHostnameCoverageOverlappingResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSApiHostnameCoverageOverlapping/ApiHostnameCoverageOverlapping.json"), &getOverlapResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetApiHostnameCoverageOverlapping",
			testutils.MockContext,
			appsec.GetApiHostnameCoverageOverlappingRequest{ConfigID: 43253, Version: 7, Hostname: "example.com"},
		).Return(&getOverlapResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSApiHostnameCoverageOverlapping/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_hostname_coverage_overlapping.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
