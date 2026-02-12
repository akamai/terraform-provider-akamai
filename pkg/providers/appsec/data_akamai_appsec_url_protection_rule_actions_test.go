package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestDataURLProtectionRuleActions(t *testing.T) {

	getURLProtectionRuleActionsResponse := appsec.GetURLProtectionRuleActionsResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSURLProtectionRuleActions/URLProtectionRuleActions.json"), &getURLProtectionRuleActionsResponse)
	require.NoError(t, err)

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return URL protection rule actions": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionRuleActions(m, 3)
				mockGetURLProtectionRuleActionsDS(m, getURLProtectionRuleActionsResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionRuleActions/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_rule_actions.test", "max_rate_threshold_action", "alert"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_rule_actions.test", "load_shedding_action", "none"),
					),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionRuleActions/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionRuleActions/missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument url_protection_rule_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionRuleActions/missing_url_protection_rule_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "url_protection_rule_id" is required`),
				},
			},
		},
		"error when url_protection_rule_id not found": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionRuleActions(m, 1)
				m.On("GetURLProtectionRuleActions", testutils.MockContext, appsec.GetURLProtectionRuleActionsRequest{
					ConfigID:            43253,
					ConfigVersion:       7,
					PolicyID:            "AAAA_81230",
					URLProtectionRuleID: 999999,
				}).Return(nil, &appsec.Error{StatusCode: 404, Title: "Not Found"}).Times(1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionRuleActions/not_found.tf"),
					ExpectError: regexp.MustCompile(`Read URL Protection Rule Actions failed`),
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

func mockGetURLProtectionRuleActionsDS(m *appsec.Mock, response appsec.GetURLProtectionRuleActionsResponse, times int) {
	m.On("GetURLProtectionRuleActions", testutils.MockContext, appsec.GetURLProtectionRuleActionsRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		PolicyID:            "AAAA_81230",
		URLProtectionRuleID: 135355,
	}).Return(&response, nil).Times(times)
}

func mockGetConfigurationURLProtectionRuleActions(m *appsec.Mock, times int) {
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
