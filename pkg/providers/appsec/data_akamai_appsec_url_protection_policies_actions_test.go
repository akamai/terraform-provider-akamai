package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestDataURLProtectionPoliciesActions(t *testing.T) {

	listURLProtectionPoliciesActionsResponse := appsec.ListURLProtectionPoliciesActionsResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSURLProtectionPoliciesActions/URLProtectionPoliciesActions.json"), &listURLProtectionPoliciesActionsResponse)
	require.NoError(t, err)

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return URL protection policies actions": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionPoliciesActions(m, 3)
				mockListURLProtectionPoliciesActions(m, listURLProtectionPoliciesActionsResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPoliciesActions/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "config_id", "43253"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "security_policy_id", "AAAA_81230"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.#", "2"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.0.url_protection_policy_id", "135355"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.0.max_rate_threshold_action", "alert"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.0.load_shedding_action", "none"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.1.url_protection_policy_id", "135356"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.1.max_rate_threshold_action", "deny"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policies_actions.test", "url_protection_policies_actions.1.load_shedding_action", "alert"),
					),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPoliciesActions/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPoliciesActions/missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &appsec.Mock{}
			if test.init != nil {
				test.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func mockListURLProtectionPoliciesActions(m *appsec.Mock, response appsec.ListURLProtectionPoliciesActionsResponse, times int) {
	m.On("ListURLProtectionPoliciesActions", testutils.MockContext, appsec.ListURLProtectionPoliciesActionsRequest{ConfigID: 43253, ConfigVersion: 7, PolicyID: "AAAA_81230"}).
		Return(&response, nil).Times(times)
}

func mockGetConfigurationURLProtectionPoliciesActions(m *appsec.Mock, times int) {
	m.On("GetConfiguration", testutils.MockContext, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43253,
			LatestVersion:     7,
			Name:              "Test Config",
			ProductionVersion: 7,
			StagingVersion:    7,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}
