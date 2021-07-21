package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAPICoProtection_res_basic(t *testing.T) {
	t.Run("match by APIConstraintsProtection ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAPIConstraintsProtectionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAPIConstraintsProtection/APIConstraintsProtectionUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAPIConstraintsProtectionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAPIConstraintsProtection/APIConstraintsProtection.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAPIConstraintsProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateAPIConstraintsProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAPIConstraintsProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyAPIConstraints: false},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				PreCheck:   func() { testAccPreCheck(t) },
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAPIConstraintsProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "enabled", "false"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResAPIConstraintsProtection/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_api_constraints_protection.test", "enabled", "false"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
