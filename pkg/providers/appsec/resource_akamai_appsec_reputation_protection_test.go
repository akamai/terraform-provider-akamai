package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationProtection_res_basic(t *testing.T) {
	t.Run("match by ReputationProtection ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		tempJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(tempJSON), &config)

		allProtectionsFalse := appsec.GetPolicyProtectionsResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResReputationProtection/PolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &allProtectionsFalse)

		oneProtectionTrue := appsec.GetPolicyProtectionsResponse{}
		tempJSON = compactJSON(loadFixtureBytes("testdata/TestResReputationProtection/UpdatedPolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &oneProtectionTrue)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()
		client.On("UpdatePolicyProtections",
			mock.Anything,
			appsec.UpdatePolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()

		// Update, with terminal Read
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()
		client.On("UpdatePolicyProtections",
			mock.Anything,
			appsec.UpdatePolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230",
				ApplyReputationControls: true},
		).Return(&oneProtectionTrue, nil).Once()
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&oneProtectionTrue, nil).Once()

		// Read, performed as part of "id" check.
		// Question: shouldn't there be another one of these for the "enabled" check?
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&oneProtectionTrue, nil).Once()

		// Delete, performed automatically to clean up
		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&oneProtectionTrue, nil).Once()
		client.On("UpdatePolicyProtections",
			mock.Anything,
			appsec.UpdatePolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				PreCheck:   func() { testAccPreCheck(t) },
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "enabled", "false"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResReputationProtection/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "enabled", "true"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
