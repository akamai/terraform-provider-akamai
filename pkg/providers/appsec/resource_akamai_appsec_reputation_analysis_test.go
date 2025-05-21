package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationAnalysis_res_basic(t *testing.T) {
	t.Run("match by ReputationAnalysis ID", func(t *testing.T) {
		client := &appsec.Mock{}

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

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetReputationAnalysis",
			testutils.MockContext,
			appsec.GetReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getReputationAnalysisResponse, nil)

		client.On("UpdateReputationAnalysis",
			testutils.MockContext,
			appsec.UpdateReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ForwardToHTTPHeader: true, ForwardSharedIPToHTTPHeaderAndSIEM: true},
		).Return(&updateReputationAnalysisResponse, nil)

		client.On("RemoveReputationAnalysis",
			testutils.MockContext,
			appsec.RemoveReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ForwardToHTTPHeader: false, ForwardSharedIPToHTTPHeaderAndSIEM: false},
		).Return(&removeReputationAnalysisResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

}
