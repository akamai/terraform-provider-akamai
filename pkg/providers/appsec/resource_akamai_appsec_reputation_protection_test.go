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
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		updateResponseAllProtectionsFalse := appsec.UpdateReputationProtectionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProtection/PolicyProtections.json"), &updateResponseAllProtectionsFalse)

		getResponseAllProtectionsFalse := appsec.GetReputationProtectionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProtection/PolicyProtections.json"), &getResponseAllProtectionsFalse)

		updateResponseOneProtectionTrue := appsec.UpdateReputationProtectionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProtection/UpdatedPolicyProtections.json"), &updateResponseOneProtectionTrue)

		getResponseOneProtectionTrue := appsec.GetReputationProtectionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResReputationProtection/UpdatedPolicyProtections.json"), &getResponseOneProtectionTrue)

		// Mock each call to the EdgeGrid library. With the exception of GetConfiguration, each call
		// is mocked individually because calls with the same parameters may have different return values.

		// All calls to GetConfiguration have same parameters and return value
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create, with terminal Read
		client.On("UpdateReputationProtection",
			mock.Anything,
			appsec.UpdateReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()
		client.On("GetReputationProtection",
			mock.Anything,
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Reads performed via "id" and "enabled" checks
		client.On("GetReputationProtection",
			mock.Anything,
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()
		client.On("GetReputationProtection",
			mock.Anything,
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseAllProtectionsFalse, nil).Once()

		// Update, with terminal Read
		client.On("UpdateReputationProtection",
			mock.Anything,
			appsec.UpdateReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230",
				ApplyReputationControls: true},
		).Return(&updateResponseOneProtectionTrue, nil).Once()
		client.On("GetReputationProtection",
			mock.Anything,
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Read, performed as part of "id" check.
		client.On("GetReputationProtection",
			mock.Anything,
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResponseOneProtectionTrue, nil).Once()

		// Delete, performed automatically to clean up
		client.On("UpdateReputationProtection",
			mock.Anything,
			appsec.UpdateReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateResponseAllProtectionsFalse, nil).Once()

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
