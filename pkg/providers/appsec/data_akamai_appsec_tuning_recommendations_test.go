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
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRecs := appsec.GetTuningRecommendationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSTuningRecommendations/Recommendations.json"))
		json.Unmarshal([]byte(expectJS), &getRecs)

		getGroupRecs := appsec.GetAttackGroupRecommendationsResponse{}
		expectJS = compactJSON(loadFixtureBytes("testdata/TestDSTuningRecommendations/AttackGroupRecommendations.json"))
		json.Unmarshal([]byte(expectJS), &getGroupRecs)

		client.On("GetTuningRecommendations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetTuningRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getRecs, nil)

		client.On("GetAttackGroupRecommendations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "XSS"},
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
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRecs := appsec.GetTuningRecommendationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSTuningRecommendations/Recommendations.json"))
		json.Unmarshal([]byte(expectJS), &getRecs)

		getGroupRecs := appsec.GetAttackGroupRecommendationsResponse{}
		expectJS = compactJSON(loadFixtureBytes("testdata/TestDSTuningRecommendations/AttackGroupRecommendations.json"))
		json.Unmarshal([]byte(expectJS), &getGroupRecs)

		client.On("GetTuningRecommendations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetTuningRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetTuningRecommendations failed"))

		client.On("GetAttackGroupRecommendations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupRecommendationsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "XSS"},
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
