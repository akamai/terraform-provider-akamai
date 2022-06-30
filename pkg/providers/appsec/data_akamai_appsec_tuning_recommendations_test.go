package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiTuningRecommendationsDataBasic(t *testing.T) {
	t.Run(" Recommendations basic", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRecs := appsec.GetTuningRecommendationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/Recommendations.json"), &getRecs)

		getGroupRecs := appsec.GetAttackGroupRecommendationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/AttackGroupRecommendations.json"), &getGroupRecs)

		client.On("GetTuningRecommendations",
			mock.Anything,
			appsec.GetTuningRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RulesetType: "active"},
		).Return(&getRecs, nil)

		client.On("GetAttackGroupRecommendations",
			mock.Anything,
			appsec.GetAttackGroupRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "XSS", RulesetType: "evaluation"},
		).Return(&getGroupRecs, nil)

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

func TestAccAkamaiTuningRecommenadationsDataErrorRetrievingTuningRecommenadations(t *testing.T) {
	t.Run("Tuning Recommendations Error", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRecs := appsec.GetTuningRecommendationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/Recommendations.json"), &getRecs)

		getGroupRecs := appsec.GetAttackGroupRecommendationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSTuningRecommendations/AttackGroupRecommendations.json"), &getGroupRecs)

		client.On("GetTuningRecommendations",
			mock.Anything,
			appsec.GetTuningRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RulesetType: "active"},
		).Return(nil, fmt.Errorf("GetTuningRecommendations failed"))

		client.On("GetAttackGroupRecommendations",
			mock.Anything,
			appsec.GetAttackGroupRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "XSS", RulesetType: "evaluation"},
		).Return(nil, fmt.Errorf("GetAttackGroupRecommendations failed"))

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
