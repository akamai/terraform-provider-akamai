package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiApiHostnameCoverage_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by ApiHostnameCoverage ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getCoverageResponse := appsec.GetApiHostnameCoverageResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSApiHostnameCoverage/ApiHostnameCoverage.json"), &getCoverageResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetApiHostnameCoverage",
			testutils.MockContext,
			appsec.GetApiHostnameCoverageRequest{ConfigID: 0, Version: 0},
		).Return(&getCoverageResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSApiHostnameCoverage/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_hostname_coverage.test", "id", "3644"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
