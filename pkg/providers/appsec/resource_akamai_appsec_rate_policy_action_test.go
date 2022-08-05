package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiRatePolicyAction_res_basic(t *testing.T) {
	t.Run("match by RatePolicyAction ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		actionAfterCreate := appsec.UpdateRatePolicyActionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRatePolicyAction/ActionAfterCreate.json"), &actionAfterCreate)

		allActionsAfterCreate := appsec.GetRatePolicyActionsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRatePolicyAction/AllActionsAfterCreate.json"), &allActionsAfterCreate)

		actionAfterUpdate := appsec.UpdateRatePolicyActionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRatePolicyAction/ActionAfterUpdate.json"), &actionAfterUpdate)

		allActionsAfterUpdate := appsec.GetRatePolicyActionsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRatePolicyAction/AllActionsAfterUpdate.json"), &allActionsAfterUpdate)

		actionAfterDelete := appsec.UpdateRatePolicyActionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRatePolicyAction/ActionAfterDelete.json"), &actionAfterDelete)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create ("none", "none")
		client.On("UpdateRatePolicyAction",
			mock.Anything,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "none", Ipv6Action: "none"},
		).Return(&actionAfterCreate, nil).Once()

		// Read called from Create, returns ("none, "none") for this policy
		client.On("GetRatePolicyActions",
			mock.Anything,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterCreate, nil).Once()

		// Read called for TestStep 1, returns ("none", "none") for this policy
		client.On("GetRatePolicyActions",
			mock.Anything,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterCreate, nil).Once()

		// Read called by TF to determine whether Update will be called (diff check), returns ("none", "none") for this policy
		client.On("GetRatePolicyActions",
			mock.Anything,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterCreate, nil).Once()

		// Update with ("alert", "deny")
		client.On("UpdateRatePolicyAction",
			mock.Anything,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "alert", Ipv6Action: "deny"},
		).Return(&actionAfterUpdate, nil).Once()

		// Read called from Update, returns ("alert", "deny") for this policy
		client.On("GetRatePolicyActions",
			mock.Anything,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterUpdate, nil).Once()

		// Read called for TestStep 2, returns ("alert", "deny") for this policy
		client.On("GetRatePolicyActions",
			mock.Anything,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterUpdate, nil).Once()

		// Delete, returns ("none", "none") and does NOT generate a Read
		client.On("UpdateRatePolicyAction",
			mock.Anything,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "none", Ipv6Action: "none"},
		).Return(&actionAfterDelete, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRatePolicyAction/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "43253:AAAA_81230:135355"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "none"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv6_action", "none"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResRatePolicyAction/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "43253:AAAA_81230:135355"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "alert"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv6_action", "deny"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
