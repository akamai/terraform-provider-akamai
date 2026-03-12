package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestURLProtectionRuleResource(t *testing.T) {
	t.Parallel()

	// Load test fixtures

	urlProtectionRuleResponse := appsec.GetURLProtectionRuleResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionRule/URLProtectionRule.json"), &urlProtectionRuleResponse)
	require.NoError(t, err)

	urlProtectionRuleUpdatedResponse := appsec.GetURLProtectionRuleResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionRule/URLProtectionRuleUpdated.json"), &urlProtectionRuleUpdatedResponse)
	require.NoError(t, err)

	urlProtectionRuleWithAPIResponse := appsec.GetURLProtectionRuleResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionRule/URLProtectionRuleWithAPI.json"), &urlProtectionRuleWithAPIResponse)
	require.NoError(t, err)

	urlProtectionRuleAPIUpdatedResponse := appsec.GetURLProtectionRuleResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionRule/URLProtectionRuleAPIUpdated.json"), &urlProtectionRuleAPIUpdatedResponse)
	require.NoError(t, err)

	createResponse := appsec.CreateURLProtectionRuleResponse{
		URLProtectionRuleID: 681,
	}

	var tests = map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required argument name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/missing_name.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required argument max_rate_threshold": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/missing_max_rate_threshold.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"validate config - both hostname_paths and api_definitions specified": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/both_hostname_and_api.tf"),
					ExpectError: regexp.MustCompile("Only one of 'hostname_paths' or 'api_definitions' can be specified"),
				},
			},
		},
		"validate config - neither hostname_paths nor api_definitions specified": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/missing_hostname_and_api.tf"),
					ExpectError: regexp.MustCompile("Either 'hostname_paths' or 'api_definitions' must be specified"),
				},
			},
		},
		"validate config - invalid custom criteria type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_custom_criteria_type.tf"),
					ExpectError: regexp.MustCompile("Invalid Custom Criteria Type"),
				},
			},
		},
		"error creating url protection rule": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionRule(m, 1)
				mockCreateURLProtectionRuleFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_hostname_paths.tf"),
					ExpectError: regexp.MustCompile("Error creating URL Protection Rule"),
				},
			},
		},
		"validate config - intelligent_load_shedding.categories is null": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_categories_null.tf"),
				ExpectError: regexp.MustCompile("categories list must not be null or unknown"),
			}},
		},
		"validate config - intelligent_load_shedding.categories is empty": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_categories_empty.tf"),
				ExpectError: regexp.MustCompile("categories list must not be empty"),
			}},
		},
		"validate config - hits_per_sec below allowed": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_hits_per_sec_low.tf"),
				ExpectError: regexp.MustCompile("hits_per_sec must be at least 25% of max_rate_threshold|min 7"),
			}},
		},
		"validate config - hits_per_sec above allowed": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_hits_per_sec_high.tf"),
				ExpectError: regexp.MustCompile("hits_per_sec must be less than or equal to 90% of max_rate_threshold"),
			}},
		},
		"validate config - custom_criteria.list_ids is null": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_custom_criteria_list_ids_null.tf"),
				ExpectError: regexp.MustCompile("list_ids must not be null or unknown"),
			}},
		},
		"validate config - custom_criteria.list_ids is empty": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_custom_criteria_list_ids_empty.tf"),
				ExpectError: regexp.MustCompile("list_ids must not be empty"),
			}},
		},
		"validate config - bypass_conditions names missing": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_bypass_condition_names_missing.tf"),
				ExpectError: regexp.MustCompile("'names' is required when type is 'RequestHeaderCondition'"),
			}},
		},
		"validate config - bypass_conditions names empty": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_bypass_condition_names_empty.tf"),
				ExpectError: regexp.MustCompile("'names' must not be empty when type is 'RequestHeaderCondition'"),
			}},
		},
		"validate config - max_rate_threshold below minimum": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_max_rate_threshold_low.tf"),
				ExpectError: regexp.MustCompile("max_rate_threshold must be at least 10"),
			}},
		},
		"validate config - intelligent_load_shedding.hits_per_sec missing": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/invalid_intelligent_load_shedding_hits_per_sec_missing.tf"),
				ExpectError: regexp.MustCompile("Incorrect attribute value type"),
			}},
		},
		"create url protection rule with hostname paths - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for Create phase
				mockGetConfigurationURLProtectionRule(m, 1)

				// Mock CreateURLProtectionRule
				mockCreateURLProtectionRuleSuccess(m, createResponse, 1)

				// Mock GetURLProtectionRule for reading after create
				mockGetURLProtectionRuleData(m, urlProtectionRuleResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionRule(m, 1)

				// Mock GetURLProtectionRule for reading after create
				mockGetURLProtectionRuleData(m, urlProtectionRuleResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionRule(m, 1)

				// Mock RemoveURLProtectionRule for cleanup
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_hostname_paths.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "url_protection_rule_id", "681"),
					),
				},
			},
		},
		"create url protection rule with api definitions - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for Create phase
				mockGetConfigurationURLProtectionRule(m, 1)

				// Mock CreateURLProtectionRule
				mockCreateURLProtectionRuleSuccess(m, createResponse, 1)

				// Mock GetURLProtectionRule for reading after create
				mockGetURLProtectionRuleData(m, urlProtectionRuleWithAPIResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionRule(m, 1)

				// Mock GetURLProtectionRule for reading after create
				mockGetURLProtectionRuleData(m, urlProtectionRuleWithAPIResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionRule(m, 1)

				// Mock RemoveURLProtectionRule for cleanup
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "url_protection_rule_id", "681"),
					),
				},
			},
		},
		"update url protection rule with hostname paths - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionRule(m, 6)
				// Mock CreateURLProtectionRule for initial resource creation
				mockCreateURLProtectionRuleSuccess(m, createResponse, 1)
				// Mock GetURLProtectionRule for reading after create (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleResponse, 3)
				// Mock UpdateURLProtectionRule
				m.On("UpdateURLProtectionRule", mock.Anything, mock.AnythingOfType("appsec.UpdateURLProtectionRuleRequest")).Return(&appsec.UpdateURLProtectionRuleResponse{
					URLProtectionRuleID: 681,
				}, nil).Once()
				// Mock GetURLProtectionRule for reading after update (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleUpdatedResponse, 2)
				// Mock RemoveURLProtectionRule for cleanup
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_hostname_paths.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "url_protection_rule_id", "681"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_hostname_paths.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "description", "Updated URL Protection"),
					),
				},
			},
		},
		"update url protection rule with api definitions - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionRule(m, 6)
				// Mock CreateURLProtectionRule for initial resource creation
				mockCreateURLProtectionRuleSuccess(m, createResponse, 1)
				// Mock GetURLProtectionRule for reading after create (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleWithAPIResponse, 3)
				// Mock UpdateURLProtectionRule
				m.On("UpdateURLProtectionRule", mock.Anything, mock.AnythingOfType("appsec.UpdateURLProtectionRuleRequest")).Return(&appsec.UpdateURLProtectionRuleResponse{
					URLProtectionRuleID: 681,
				}, nil).Once()
				// Mock GetURLProtectionRule for reading after update (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleAPIUpdatedResponse, 2)
				// Mock RemoveURLProtectionRule for cleanup
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "url_protection_rule_id", "681"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "description", "Updated API Protection"),
					),
				},
			},
		},
		"create and update with different config_id should fail": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionRule(m, 4)
				// Mock CreateURLProtectionRule for initial resource creation
				mockCreateURLProtectionRuleSuccess(m, createResponse, 1)
				// Mock GetURLProtectionRule for reading after create (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleWithAPIResponse, 2)
				// Mock GetURLProtectionRule for reading after update (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleAPIUpdatedResponse, 1)
				// Mock RemoveURLProtectionRule for cleanup
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "url_protection_rule_id", "681"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_api_definitions_different_configID.tf"),
					ExpectError: regexp.MustCompile("updating field `config_id` is not possible"),
				},
			},
		},
		"create and update with different name should fail": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionRule(m, 4)
				// Mock CreateURLProtectionRule for initial resource creation
				mockCreateURLProtectionRuleSuccess(m, createResponse, 1)
				// Mock GetURLProtectionRule for reading after create (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleWithAPIResponse, 2)
				// Mock GetURLProtectionRule for reading after update (may be called multiple times)
				mockGetURLProtectionRuleData(m, urlProtectionRuleAPIUpdatedResponse, 1)
				// Mock RemoveURLProtectionRule for cleanup
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_rule.test", "url_protection_rule_id", "681"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_api_definitions_different_name.tf"),
					ExpectError: regexp.MustCompile("updating field `name` is not possible"),
				},
			},
		},
		"import url protection": {
			init: func(m *appsec.Mock) {
				// Import and post-import refresh reads
				mockGetConfigurationURLProtectionRule(m, 1)

				mockGetURLProtectionRuleData(m, urlProtectionRuleWithAPIResponse, 1)

				mockGetConfigurationURLProtectionRule(m, 1)
				// Delete after test
				mockRemoveURLProtectionRuleSuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "43007:681",
					ResourceName:  "akamai_appsec_url_protection_rule.test",
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("config_id", "43007").
						CheckEqual("url_protection_rule_id", "681").
						CheckEqual("max_rate_threshold", "195").
						Build(),
					ImportStatePersist: true,
				},
			},
		},
		"import url protection - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_appsec_url_protection_rule.test",
					ExpectError:   regexp.MustCompile("Invalid Import ID"),
				},
			},
		},
		"import url protection - invalid config id value": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "abc:681",
					ResourceName:  "akamai_appsec_url_protection_rule.test",
					ExpectError:   regexp.MustCompile("Invalid Config ID"),
				},
			},
		},
		"import url protection - invalid url protection id value": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionRule/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "43007:xyz",
					ResourceName:  "akamai_appsec_url_protection_rule.test",
					ExpectError:   regexp.MustCompile("Invalid URL Protection ID"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &appsec.Mock{}
			if tc.init != nil {
				tc.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

// Mock functions for URL Protection Rule resource tests

func mockGetConfigurationURLProtectionRule(m *appsec.Mock, times int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43007}).
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

func mockCreateURLProtectionRuleFailure(m *appsec.Mock, times int) {
	m.On("CreateURLProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.CreateURLProtectionRuleRequest) bool {
		return req.ConfigID == 43007 && req.ConfigVersion == 40
	})).Return(nil, fmt.Errorf("create url protection rule failed")).Times(times)
}

func mockCreateURLProtectionRuleSuccess(m *appsec.Mock, response appsec.CreateURLProtectionRuleResponse, times int) {
	m.On("CreateURLProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.CreateURLProtectionRuleRequest) bool {
		return req.ConfigID == 43007 && req.ConfigVersion == 40
	})).Return(&response, nil).Times(times)
}

func mockRemoveURLProtectionRuleSuccess(m *appsec.Mock, times int) {
	m.On("RemoveURLProtectionRule", mock.Anything, appsec.RemoveURLProtectionRuleRequest{
		ConfigID:            43007,
		ConfigVersion:       40,
		URLProtectionRuleID: 681,
	}).Return(nil).Times(times)
}
