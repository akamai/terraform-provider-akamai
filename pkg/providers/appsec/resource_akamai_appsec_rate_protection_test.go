package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRateProtection_res_basic(t *testing.T) {
	t.Run("match by RateProtection ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		tempJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(tempJSON), &config)

		updateResponseAllProtectionsFalse := appsec.UpdateRateProtectionResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResRateProtection/PolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &updateResponseAllProtectionsFalse)

		getResponseAllProtectionsFalse := appsec.GetRateProtectionResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResRateProtection/PolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &getResponseAllProtectionsFalse)

		updateResponseOneProtectionTrue := appsec.UpdateRateProtectionResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResRateProtection/UpdatedPolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &updateResponseOneProtectionTrue)

		getResponseOneProtectionTrue := appsec.GetRateProtectionResponse{}
		json.Unmarshal([]byte(tempJSON), &getResponseOneProtectionTrue)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.On("UpdateRateProtection",
			mock.Anything,
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()
		client.On("GetRateProtection",
			mock.Anything,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.On("GetRateProtection",
			mock.Anything,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()
		client.On("GetRateProtection",
			mock.Anything,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Update, with terminal Read
		client.On("UpdateRateProtection",
			mock.Anything,
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230",
				ApplyRateControls: true},
		).Return(&updateResponseOneProtectionTrue, nil).Once()
		client.On("GetRateProtection",
			mock.Anything,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Read, performed as part of "id" check.
		client.On("GetRateProtection",
			mock.Anything,
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Delete, performed automatically to clean up
		client.On("UpdateRateProtection",
			mock.Anything,
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				PreCheck:   func() { testAccPreCheck(t) },
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRateProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "enabled", "false"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResRateProtection/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "enabled", "true"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
