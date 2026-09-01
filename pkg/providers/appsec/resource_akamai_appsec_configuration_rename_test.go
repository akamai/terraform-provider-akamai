package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfigurationRename_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by Configuration ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateConfigurationResponse := appsec.UpdateConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfigurationRename/ConfigurationUpdate.json"), &updateConfigurationResponse)
		require.NoError(t, err)

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfigurationRename/Configuration.json"), &getConfigurationResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 432531},
		).Return(&getConfigurationResponse, nil)

		client.APPSEC.On("UpdateConfiguration",
			testutils.MockContext,
			appsec.UpdateConfigurationRequest{ConfigID: 432531, Name: "Akamai Tools New", Description: "TF Tools"},
		).Return(&updateConfigurationResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
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

		client.APPSEC.AssertExpectations(t)
	})

}
