package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRatePolicy_res_basic(t *testing.T) {
	t.Run("match by RatePolicy ID", func(t *testing.T) {
		client := &mockappsec{}

		configResponse := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &configResponse)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createResponse := appsec.CreateRatePolicyResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicy.json")), &createResponse)
		createRatePolicyJSON := loadFixtureBytes("testdata/TestResRatePolicy/CreateRatePolicy.json")
		client.On("CreateRatePolicy",
			mock.Anything,
			appsec.CreateRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createRatePolicyJSON},
		).Return(&createResponse, nil)

		getResponse := appsec.GetRatePolicyResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicy.json")), &getResponse)
		client.On("GetRatePolicy",
			mock.Anything,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Once()
		client.On("GetRatePolicy",
			mock.Anything,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Once()
		client.On("GetRatePolicy",
			mock.Anything,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Once()

		updateResponse := appsec.UpdateRatePolicyResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyUpdated.json")), &updateResponse)
		updateRatePolicyJSON := loadFixtureBytes("testdata/TestResRatePolicy/UpdateRatePolicy.json")
		client.On("UpdateRatePolicy",
			mock.Anything,
			appsec.UpdateRatePolicyRequest{RatePolicyID: 134644, PolicyID: 0, ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: updateRatePolicyJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdate := appsec.GetRatePolicyResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyUpdated.json")), &getResponseAfterUpdate)
		client.On("GetRatePolicy",
			mock.Anything,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponse := appsec.RemoveRatePolicyResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyEmpty.json")), &removeResponse)
		client.On("RemoveRatePolicy",
			mock.Anything,
			appsec.RemoveRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRatePolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "43253:134644"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResRatePolicy/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "43253:134644"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
