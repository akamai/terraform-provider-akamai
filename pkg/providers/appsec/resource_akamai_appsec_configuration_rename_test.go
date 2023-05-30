package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfigurationRename_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateConfigurationResponse := appsec.UpdateConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfigurationRename/ConfigurationUpdate.json"), &updateConfigurationResponse)
		require.NoError(t, err)

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfigurationRename/Configuration.json"), &getConfigurationResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 432531},
		).Return(&getConfigurationResponse, nil)

		client.On("UpdateConfiguration",
			mock.Anything,
			appsec.UpdateConfigurationRequest{ConfigID: 432531, Name: "Akamai Tools New", Description: "TF Tools"},
		).Return(&updateConfigurationResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfigurationRename/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration_rename.test", "id", "432531"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResConfigurationRename/update_by_id.tf"),
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
