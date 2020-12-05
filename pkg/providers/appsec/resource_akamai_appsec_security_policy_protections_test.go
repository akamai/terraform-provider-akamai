package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiPolicyProtections_res_basic(t *testing.T) {
	t.Run("match by PolicyProtections ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdatePolicyProtectionsResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResPolicyProtections/PolicyProtectionsUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetPolicyProtectionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResPolicyProtections/PolicyProtections.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemovePolicyProtectionsResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResPolicyProtections/PolicyProtections.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetPolicyProtections",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyAPIConstraints: false, ApplyApplicationLayerControls: false, ApplyBotmanControls: false, ApplyNetworkLayerControls: false, ApplyRateControls: false, ApplyReputationControls: false, ApplySlowPostControls: false},
		).Return(&cr, nil)

		client.On("UpdatePolicyProtections",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdatePolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyAPIConstraints: false, ApplyApplicationLayerControls: false, ApplyBotmanControls: false, ApplyNetworkLayerControls: false, ApplyRateControls: false, ApplyReputationControls: false, ApplySlowPostControls: false},
		).Return(&cu, nil)

		client.On("RemovePolicyProtections",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemovePolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyAPIConstraints: false, ApplyApplicationLayerControls: false, ApplyBotmanControls: false, ApplyNetworkLayerControls: false, ApplyRateControls: false, ApplyReputationControls: false, ApplySlowPostControls: false},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResPolicyProtections/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_security_policy_protections.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
