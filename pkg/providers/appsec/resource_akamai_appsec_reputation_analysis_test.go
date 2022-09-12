package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationAnalysis_res_basic(t *testing.T) {
	t.Run("match by ReputationAnalysis ID", func(t *testing.T) {
		client := &mockappsec{}

		updateReputationAnalysisResponse := appsec.UpdateReputationAnalysisResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResReputationAnalysis/ReputationAnalysisUpdated.json"), &updateReputationAnalysisResponse)
		require.NoError(t, err)

		getReputationAnalysisResponse := appsec.GetReputationAnalysisResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResReputationAnalysis/ReputationAnalysis.json"), &getReputationAnalysisResponse)
		require.NoError(t, err)

		removeReputationAnalysisResponse := appsec.RemoveReputationAnalysisResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResReputationAnalysis/ReputationAnalysisDelete.json"), &removeReputationAnalysisResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetReputationAnalysis",
			mock.Anything,
			appsec.GetReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getReputationAnalysisResponse, nil)

		client.On("UpdateReputationAnalysis",
			mock.Anything,
			appsec.UpdateReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ForwardToHTTPHeader: true, ForwardSharedIPToHTTPHeaderAndSIEM: true},
		).Return(&updateReputationAnalysisResponse, nil)

		client.On("RemoveReputationAnalysis",
			mock.Anything,
			appsec.RemoveReputationAnalysisRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ForwardToHTTPHeader: false, ForwardSharedIPToHTTPHeaderAndSIEM: false},
		).Return(&removeReputationAnalysisResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationAnalysis/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile_analysis.test", "id", "43253:AAAA_81230"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResReputationAnalysis/update_by_id.tf"),
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
