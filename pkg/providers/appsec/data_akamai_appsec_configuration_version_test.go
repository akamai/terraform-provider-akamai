package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfigurationVersion_data_basic(t *testing.T) {
	t.Run("match by ConfigurationVersion ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfigurationVersion/ConfigurationVersion.json"), &getConfigurationVersionsResponse)
		require.NoError(t, err)

		client.On("GetConfigurationVersions",
			testutils.MockContext,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSConfigurationVersion/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_configuration_version.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
