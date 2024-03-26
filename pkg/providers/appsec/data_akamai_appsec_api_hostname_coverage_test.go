package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiApiHostnameCoverage_data_basic(t *testing.T) {
	t.Run("match by ApiHostnameCoverage ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getCoverageResponse := appsec.GetApiHostnameCoverageResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSApiHostnameCoverage/ApiHostnameCoverage.json"), &getCoverageResponse)
		require.NoError(t, err)

		client.On("GetApiHostnameCoverage",
			mock.Anything,
			appsec.GetApiHostnameCoverageRequest{ConfigID: 0, Version: 0},
		).Return(&getCoverageResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSApiHostnameCoverage/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_hostname_coverage.test", "id", "3644"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
