package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationAnalysis_data_basic(t *testing.T) {
	t.Run("match by ReputationAnalysis ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getReputationAnalysisResponse := appsec.GetReputationAnalysisResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSReputationAnalysis/ReputationAnalysis.json"), &getReputationAnalysisResponse)
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

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSReputationAnalysis/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_reputation_profile_analysis.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
