package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationProfile_res_basic(t *testing.T) {
	t.Run("match by ReputationProfile ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateReputationProfileResponse := appsec.UpdateReputationProfileResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfileUpdated.json"), &updateReputationProfileResponse)
		require.NoError(t, err)

		getReputationProfileResponse := appsec.GetReputationProfileResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfiles.json"), &getReputationProfileResponse)
		require.NoError(t, err)

		createReputationProfileResponse := appsec.CreateReputationProfileResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfileCreated.json"), &createReputationProfileResponse)
		require.NoError(t, err)

		removeReputationProfileResponse := appsec.RemoveReputationProfileResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfileCreated.json"), &removeReputationProfileResponse)
		require.NoError(t, err)

		client.On("GetReputationProfile",
			testutils.MockContext,
			appsec.GetReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&getReputationProfileResponse, nil)

		client.On("RemoveReputationProfile",
			testutils.MockContext,
			appsec.RemoveReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&removeReputationProfileResponse, nil)

		client.On("CreateReputationProfile",
			testutils.MockContext,
			appsec.CreateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&createReputationProfileResponse, nil)

		client.On("UpdateReputationProfile",
			testutils.MockContext,
			appsec.UpdateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&updateReputationProfileResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               false,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResReputationProfile/match_by_id.tf"),
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
