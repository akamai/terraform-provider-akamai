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

		config := appsec.GetConfigurationResponse{}
		tempJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(tempJSON), &config)

		updateResponseAllProtectionsFalse := appsec.UpdateSlowPostProtectionResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtection/PolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &updateResponseAllProtectionsFalse)

		getResponseAllProtectionsFalse := appsec.GetSlowPostProtectionResponse{}
		json.Unmarshal([]byte(tempJSON), &getResponseAllProtectionsFalse)

		updateResponseOneProtectionTrue := appsec.UpdateSlowPostProtectionResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtection/UpdatedPolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &updateResponseOneProtectionTrue)

		getResponseOneProtectionTrue := appsec.GetSlowPostProtectionResponse{}
		json.Unmarshal([]byte(tempJSON), &getResponseOneProtectionTrue)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.On("UpdateSlowPostProtection",
			mock.Anything,
			appsec.UpdateSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()
		client.On("GetSlowPostProtection",
			mock.Anything,
			appsec.GetSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.On("GetSlowPostProtection",
			mock.Anything,
			appsec.GetSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Delete, performed automatically to clean up
		client.On("UpdateSlowPostProtection",
			mock.Anything,
			appsec.UpdateSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				PreCheck:   func() { testAccPreCheck(t) },
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSlowPostProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_slowpost_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_slowpost_protection.test", "enabled", "false"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
