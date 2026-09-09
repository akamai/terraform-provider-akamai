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

func TestAkamaiMatchTargetSequence_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by MatchTargetSequence ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateSequenceResponse := appsec.UpdateMatchTargetSequenceResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTargetSequence/MatchTargetSequenceResp.json"), &updateSequenceResponse)
		require.NoError(t, err)

		getSequenceResponse := appsec.GetMatchTargetSequenceResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTargetSequence/MatchTargetSequence.json"), &getSequenceResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetMatchTargetSequence",
			testutils.MockContext,
			appsec.GetMatchTargetSequenceRequest{ConfigID: 43253, ConfigVersion: 7, Type: "website"},
		).Return(&getSequenceResponse, nil)

		client.APPSEC.On("UpdateMatchTargetSequence",
			testutils.MockContext,
			appsec.UpdateMatchTargetSequenceRequest{ConfigID: 43253, ConfigVersion: 7, Type: "website", TargetSequence: []appsec.MatchTargetItem{{Sequence: 1, TargetID: 2052813}, {Sequence: 2, TargetID: 2971336}}},
		).Return(&updateSequenceResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTargetSequence/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_match_target_sequence.test", "id", "43253:website"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTargetSequence/update_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_match_target_sequence.test", "id", "43253:website"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
