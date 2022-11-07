package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiTuningRecommendationsDataBasic(t *testing.T) {
	t.Run(" Recommendations basic", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRecs := appsec.GetTuningRecommendationsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/Recommendations.json"), &getRecs)
		require.NoError(t, err)

		getGroupRecs := appsec.GetAttackGroupRecommendationsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/AttackGroupRecommendations.json"), &getGroupRecs)
		require.NoError(t, err)

		getRuleRecs := appsec.GetRuleRecommendationsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/RuleRecommendations.json"), &getRuleRecs)
		require.NoError(t, err)

		client.On("GetTuningRecommendations",
			mock.Anything,
			appsec.GetTuningRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RulesetType: "active"},
		).Return(&getRecs, nil)

		client.On("GetAttackGroupRecommendations",
			mock.Anything,
			appsec.GetAttackGroupRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "XSS", RulesetType: "evaluation"},
		).Return(&getGroupRecs, nil)

		client.On("GetRuleRecommendations",
			mock.Anything,
			appsec.GetRuleRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 958008, RulesetType: "active"},
		).Return(&getRuleRecs, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSTuningRecommendations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_tuning_recommendations.recommendations", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiTuningRecommenadationsDataErrorRetrievingTuningRecommenadations(t *testing.T) {
	t.Run("Tuning Recommendations Error", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRecs := appsec.GetTuningRecommendationsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/Recommendations.json"), &getRecs)
		require.NoError(t, err)

		getGroupRecs := appsec.GetAttackGroupRecommendationsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/AttackGroupRecommendations.json"), &getGroupRecs)
		require.NoError(t, err)

		getRuleRecs := appsec.GetAttackGroupRecommendationsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/RuleRecommendations.json"), &getRuleRecs)
		require.NoError(t, err)

		client.On("GetTuningRecommendations",
			mock.Anything,
			appsec.GetTuningRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RulesetType: "active"},
		).Return(nil, fmt.Errorf("GetTuningRecommendations failed"))

		client.On("GetAttackGroupRecommendations",
			mock.Anything,
			appsec.GetAttackGroupRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "XSS", RulesetType: "evaluation"},
		).Return(nil, fmt.Errorf("GetAttackGroupRecommendations failed"))

		client.On("GetRuleRecommendations",
			mock.Anything,
			appsec.GetRuleRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 958008, RulesetType: "active"},
		).Return(nil, fmt.Errorf("GetRuleRecommendations failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSTuningRecommendations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_tuning_recommendations.recommendations", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`GetTuningRecommendations failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
