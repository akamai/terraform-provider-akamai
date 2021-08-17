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

		cu := appsec.UpdateReputationProtectionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResReputationProtection/ReputationProtectionUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetReputationProtectionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResReputationProtection/ReputationProtection.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveReputationProtectionResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResReputationProtection/ReputationProtectionDelete.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetReputationProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateReputationProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cu, nil)

		client.On("RemoveReputationProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
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
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "enabled", "false"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
