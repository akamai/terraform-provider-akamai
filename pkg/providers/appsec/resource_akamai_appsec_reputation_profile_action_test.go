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

func TestAkamaiReputationProfileAction_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by ReputationProfileAction ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateReputationProfileActionResponse := appsec.UpdateReputationProfileActionResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfileAction/ReputationProfileAction.json"), &updateReputationProfileActionResponse)
		require.NoError(t, err)

		getReputationProfileActionResponse := appsec.GetReputationProfileActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfileAction/ReputationProfileAction.json"), &getReputationProfileActionResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetReputationProfileAction",
			testutils.MockContext,
			appsec.GetReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099},
		).Return(&getReputationProfileActionResponse, nil)

		client.APPSEC.On("UpdateReputationProfileAction",
			testutils.MockContext,
			appsec.UpdateReputationProfileActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 1685099, Action: "none"},
		).Return(&updateReputationProfileActionResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReputationProfileAction/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_action.test", "id", "43253:AAAA_81230:1685099"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
