package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiReputationProfile_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by ReputationProfile ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		getReputationProfileResponse := appsec.GetReputationProfileResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfiles.json"), &getReputationProfileResponse)
		require.NoError(t, err)

		createReputationProfileResponse := appsec.CreateReputationProfileResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfileCreated.json"), &createReputationProfileResponse)
		require.NoError(t, err)

		removeReputationProfileResponse := appsec.RemoveReputationProfileResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/ReputationProfileCreated.json"), &removeReputationProfileResponse)
		require.NoError(t, err)

		createReputationProfileJSON := testutils.LoadFixtureBytes(t, "testdata/TestResReputationProfile/CreateReputationProfile.json")

		client.APPSEC.On("CreateReputationProfile",
			testutils.MockContext,
			appsec.CreateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createReputationProfileJSON},
		).Return(&createReputationProfileResponse, nil)

		client.APPSEC.On("GetReputationProfile",
			testutils.MockContext,
			appsec.GetReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&getReputationProfileResponse, nil)

		client.APPSEC.On("RemoveReputationProfile",
			testutils.MockContext,
			appsec.RemoveReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&removeReputationProfileResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReputationProfile/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_reputation_profile.test", "id", "43253:12345"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
