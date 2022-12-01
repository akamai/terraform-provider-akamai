package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationProfile_res_basic(t *testing.T) {
	t.Run("match by ReputationProfile ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateReputationProfileResponse := appsec.UpdateReputationProfileResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileUpdated.json"), &updateReputationProfileResponse)
		require.NoError(t, err)

		getReputationProfileResponse := appsec.GetReputationProfileResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfiles.json"), &getReputationProfileResponse)
		require.NoError(t, err)

		createReputationProfileResponse := appsec.CreateReputationProfileResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileCreated.json"), &createReputationProfileResponse)
		require.NoError(t, err)

		removeReputationProfileResponse := appsec.RemoveReputationProfileResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileCreated.json"), &removeReputationProfileResponse)
		require.NoError(t, err)

		client.On("GetReputationProfile",
			mock.Anything,
			appsec.GetReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&getReputationProfileResponse, nil)

		client.On("RemoveReputationProfile",
			mock.Anything,
			appsec.RemoveReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&removeReputationProfileResponse, nil)

		client.On("CreateReputationProfile",
			mock.Anything,
			appsec.CreateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&createReputationProfileResponse, nil)

		client.On("UpdateReputationProfile",
			mock.Anything,
			appsec.UpdateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&updateReputationProfileResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationProfile/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile.test", "id", "12345"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
