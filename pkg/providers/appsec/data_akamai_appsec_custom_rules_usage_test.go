package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataCustomRulesUsage(t *testing.T) {
	getCustomRulesUsageResponse := appsec.GetCustomRulesUsageResponse{}
	getCustomRulesUsageRequest := appsec.GetCustomRulesUsageRequest{
		ConfigID: 111111,
		Version:  2,
		RequestBody: appsec.RuleIDs{
			IDs: []int64{12345, 67890},
		},
	}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSCustomRulesUsage/CustomRulesUsage.json"), &getCustomRulesUsageResponse)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_custom_rules_usage.test")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return custom rules usage": {
			init: func(m *appsec.Mock) {
				mockGetLatestConfiguration(m, 111111, 3)
				mockGetCustomRulesUsage(m, getCustomRulesUsageRequest, getCustomRulesUsageResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSCustomRulesUsage/custom_rules_usage.tf"),
					Check: baseChecker.
						CheckEqual("rules.0.rule_id", "12345").
						CheckEqual("rules.0.policies.0.policy_id", "POLICY_1").
						CheckEqual("rules.0.policies.0.policy_name", "Policy One").
						CheckEqual("rules.1.rule_id", "67890").
						CheckEqual("rules.1.policies.0.policy_id", "POLICY_2").
						CheckEqual("rules.1.policies.0.policy_name", "Policy Two").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSCustomRulesUsage/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required`),
				},
			},
		},
		"missing required argument rule_ids": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSCustomRulesUsage/missing_rule_ids.tf"),
					ExpectError: regexp.MustCompile(`The argument "rule_ids" is required`),
				},
			},
		},
		"error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetLatestConfigurationFailure(m, 111111, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSCustomRulesUsage/custom_rules_usage.tf"),
					ExpectError: regexp.MustCompile(`Error: get latest config version error`),
				},
			},
		},
		"error response from GetCustomRuleUsage api": {
			init: func(m *appsec.Mock) {
				mockGetLatestConfiguration(m, 111111, 1)
				mockGetCustomRulesUsageFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSCustomRulesUsage/custom_rules_usage.tf"),
					ExpectError: regexp.MustCompile(`Error: calling 'GetCustomRuleUsage'`),
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

func mockGetCustomRulesUsage(m *appsec.Mock, request appsec.GetCustomRulesUsageRequest, response appsec.GetCustomRulesUsageResponse, times int) {
	m.On("GetCustomRulesUsage", mock.Anything, request).Return(&response, nil).Times(times)
}

func mockGetCustomRulesUsageFailure(m *appsec.Mock) {
	m.On("GetCustomRulesUsage", mock.Anything, mock.Anything).
		Return(nil, &serverError).Once()
}

func mockGetLatestConfiguration(client *appsec.Mock, configID int, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: configID}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43253,
			LatestVersion:     2,
			Name:              "Akamai Tools",
			ProductionVersion: 1,
			StagingVersion:    1,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}

func mockGetLatestConfigurationFailure(client *appsec.Mock, configID int, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: configID}).
		Return(nil, &serverError).Times(times)
}
