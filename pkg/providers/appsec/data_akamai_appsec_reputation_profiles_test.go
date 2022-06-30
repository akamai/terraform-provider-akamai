package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationProfiles_data_basic(t *testing.T) {
	t.Run("match by ReputationProfiles ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getReputationProfilesResponse := appsec.GetReputationProfilesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSReputationProfiles/ReputationProfiles.json"), &getReputationProfilesResponse)

		client.On("GetReputationProfiles",
			mock.Anything,
			appsec.GetReputationProfilesRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&getReputationProfilesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSReputationProfiles/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_reputation_profiles.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
