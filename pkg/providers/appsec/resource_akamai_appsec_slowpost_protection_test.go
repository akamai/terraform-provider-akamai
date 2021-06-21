package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSlowPostProtection_res_basic(t *testing.T) {
	t.Run("match by SlowPostProtection ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSlowPostProtectionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtection/SlowPostProtection.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSlowPostProtectionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtection/SlowPostProtection.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSlowPostProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateSlowPostProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSlowPostProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_slowpost_protection.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
