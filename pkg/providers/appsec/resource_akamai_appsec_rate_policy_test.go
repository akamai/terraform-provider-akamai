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

		cu := appsec.UpdateRatePolicyResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetRatePolicyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicy.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crpol := appsec.CreateRatePolicyResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicy.json"))
		json.Unmarshal([]byte(expectJSC), &crpol)

		crpolr := appsec.RemoveRatePolicyResponse{}
		expectJSCR := compactJSON(loadFixtureBytes("testdata/TestResRatePolicy/RatePolicyEmpty.json"))
		json.Unmarshal([]byte(expectJSCR), &crpolr)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&cr, nil)

		updateRatePolicyJSON := loadFixtureBytes("testdata/TestResRatePolicy/UpdateRatePolicy.json")
		client.On("UpdateRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRatePolicyRequest{RatePolicyID: 134644, PolicyID: 0, ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: updateRatePolicyJSON},
		).Return(&cu, nil)

		createRatePolicyJSON := loadFixtureBytes("testdata/TestResRatePolicy/CreateRatePolicy.json")
		client.On("CreateRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createRatePolicyJSON},
		).Return(&crpol, nil)

		client.On("RemoveRatePolicy",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&crpolr, nil)

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
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
