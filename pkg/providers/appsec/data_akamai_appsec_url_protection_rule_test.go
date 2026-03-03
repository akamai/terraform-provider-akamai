package appsec

import (
	"encoding/json"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataUrlProtectionRule(t *testing.T) {

	getURLProtectionRule := appsec.GetURLProtectionRuleResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSURLProtection/URLProtectionRule.json"), &getURLProtectionRule)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_url_protection_rule.test")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return a url protection rule": {
			init: func(m *appsec.Mock) {
				mockGetConfig(m, 3)
				mockGetURLProtectionRuleData(m, getURLProtectionRule, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/match_by_id.tf"),
					Check: baseChecker.
						CheckEqual("url_protection_rule_id", strconv.Itoa(681)).
						CheckEqual("config_id", strconv.Itoa(43007)).
						CheckEqual("description", "URL Protection").
						CheckEqual("name", "URL Protection").
						CheckEqual("used", strconv.FormatBool(true)).
						CheckEqual("max_rate_threshold", strconv.Itoa(195)).
						CheckEqual("hostname_paths.0.hostname", "custom.com").
						CheckEqual("hostname_paths.0.paths.0", "/asd").
						CheckEqual("hostname_paths.0.paths.1", "/my-test-path").
						CheckEqual("bypass_conditions.0.type", "NetworkListCondition").
						CheckEqual("bypass_conditions.1.type", "RequestHeaderCondition").
						CheckMissing("api_definitions").
						CheckEqual("intelligent_load_shedding.hits_per_sec", strconv.Itoa(150)).
						CheckEqual("intelligent_load_shedding.categories.0", "BOTS").
						CheckEqual("intelligent_load_shedding.categories.1", "CLOUD_PROVIDERS").
						CheckEqual("intelligent_load_shedding.categories.2", "PROXIES").
						CheckEqual("intelligent_load_shedding.categories.3", "TOR_EXIT_NODES").
						CheckEqual("intelligent_load_shedding.categories.4", "PLATFORM_DDOS_INTELLIGENCE").
						CheckEqual("intelligent_load_shedding.custom_criteria.0.type", "CLIENT_LIST").
						CheckEqual("intelligent_load_shedding.custom_criteria.0.positive_match", strconv.FormatBool(true)).
						CheckEqual("intelligent_load_shedding.custom_criteria.0.list_ids.0", "12345_10CLIENTLIST").
						CheckEqual("intelligent_load_shedding.custom_criteria.0.list_ids.1", "54321_123").
						CheckEqual("intelligent_load_shedding.custom_criteria.1.type", "CLIENT_LIST").
						CheckEqual("intelligent_load_shedding.custom_criteria.1.positive_match", strconv.FormatBool(true)).
						CheckEqual("intelligent_load_shedding.custom_criteria.1.list_ids.0", "16656_CPISERVERS").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/config_missing_match_by_id.tf"),
					ExpectError: regexp.MustCompile(`(?s)The argument "config_id" is required, but no definition was.*found`),
				},
			},
		},
		"missing required argument url_protection_rule_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/url_protection_rule_id_missing_match_by_id.tf"),
					ExpectError: regexp.MustCompile(`(?s)The argument "url_protection_rule_id" is required, but no definition was.*found`),
				},
			},
		},
		"error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error: invalid config version"),
				},
			},
		},
		"error response from GetURLProtectionRule api": {
			init: func(m *appsec.Mock) {
				mockGetConfig(m, 1)
				mockGetURLProtectionRuleFailureData(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'GetURLProtectionRule'"),
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

func mockGetConfig(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43007}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43007,
			LatestVersion:     40,
			Name:              "Akamai Tools",
			ProductionVersion: 1,
			StagingVersion:    1,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}

func mockGetURLProtectionRuleData(m *appsec.Mock, response appsec.GetURLProtectionRuleResponse, times int) {
	m.On("GetURLProtectionRule", mock.Anything, appsec.GetURLProtectionRuleRequest{ConfigID: 43007, ConfigVersion: 40, URLProtectionRuleID: 681}).
		Return(&response, nil).Times(times)
}

func mockGetURLProtectionRuleFailureData(m *appsec.Mock, times int) {
	m.On("GetURLProtectionRule", mock.Anything, appsec.GetURLProtectionRuleRequest{ConfigID: 43007, ConfigVersion: 40, URLProtectionRuleID: 681}).
		Return(nil, &serverError).Times(times)
}

func mockGetConfigFailure(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43007}).
		Return(nil, &serverError).Times(times)
}
