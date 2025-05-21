package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfigurationRename_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateConfigurationResponse := appsec.UpdateConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfigurationRename/ConfigurationUpdate.json"), &updateConfigurationResponse)
		require.NoError(t, err)

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfigurationRename/Configuration.json"), &getConfigurationResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 432531},
		).Return(&getConfigurationResponse, nil)

		client.On("UpdateConfiguration",
			testutils.MockContext,
			appsec.UpdateConfigurationRequest{ConfigID: 432531, Name: "Akamai Tools New", Description: "TF Tools"},
		).Return(&updateConfigurationResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfigurationRename/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_rename.test", "id", "432531"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfigurationRename/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_rename.test", "id", "432531"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
