package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationProfile_res_basic(t *testing.T) {
	t.Run("match by ReputationProfile ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateReputationProfileResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetReputationProfileResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfiles.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crp := appsec.CreateReputationProfileResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileCreated.json"))
		json.Unmarshal([]byte(expectJSC), &crp)

		crd := appsec.RemoveReputationProfileResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileCreated.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&cr, nil)

		client.On("RemoveReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&crd, nil)

		client.On("CreateReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&crp, nil)

		client.On("UpdateReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&cu, nil)

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
