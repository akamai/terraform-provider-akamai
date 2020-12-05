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

		client.On("GetReputationProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyReputationControls: false},
		).Return(&cr, nil)

		client.On("UpdateReputationProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyReputationControls: false},
		).Return(&cu, nil)

		client.On("RemoveReputationProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveReputationProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyReputationControls: false},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "enabled", "false"),
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResReputationProtection/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_reputation_protection.test", "enabled", "false"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
