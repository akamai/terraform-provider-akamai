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

func TestAkamaiReputationAnalysis_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by ReputationAnalysis ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateReputationAnalysisResponse := appsec.UpdateReputationAnalysisResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationAnalysis/ReputationAnalysisUpdated.json"), &updateReputationAnalysisResponse)
		require.NoError(t, err)

		getReputationAnalysisResponse := appsec.GetReputationAnalysisResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationAnalysis/ReputationAnalysis.json"), &getReputationAnalysisResponse)
		require.NoError(t, err)

		removeReputationAnalysisResponse := appsec.RemoveReputationAnalysisResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationAnalysis/ReputationAnalysisDelete.json"), &removeReputationAnalysisResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetReputationAnalysis",
			testutils.MockContext,
			appsec.GetReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getReputationAnalysisResponse, nil)

		client.APPSEC.On("UpdateReputationAnalysis",
			testutils.MockContext,
			appsec.UpdateReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ForwardToHTTPHeader: true, ForwardSharedIPToHTTPHeaderAndSIEM: true},
		).Return(&updateReputationAnalysisResponse, nil)

		client.APPSEC.On("RemoveReputationAnalysis",
			testutils.MockContext,
			appsec.RemoveReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ForwardToHTTPHeader: false, ForwardSharedIPToHTTPHeaderAndSIEM: false},
		).Return(&removeReputationAnalysisResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReputationAnalysis/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_analysis.test", "id", "43253:AAAA_81230"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReputationAnalysis/update_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_analysis.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
