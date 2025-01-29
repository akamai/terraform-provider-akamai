package property

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourcePropertyInclude(t *testing.T) {
	type testData struct {
		assetID           string
		groupID           string
		rulesPath         string
		includeID         string
		productID         string
		contractID        string
		ruleFormat        string
		includeName       string
		latestVersion     int
		stagingVersion    *int
		productionVersion *int
		stagingStatus     papi.VersionStatus
		productionStatus  papi.VersionStatus
		includeType       papi.IncludeType
		rules             papi.RulesUpdate
	}

	workdir := "./testdata/TestResPropertyInclude"
	includeID := "inc_123"

	newCreateIncludeResp := func(testData *testData) *papi.CreateIncludeResponse {
		return &papi.CreateIncludeResponse{
			IncludeLink: fmt.Sprintf("/papi/v1/includes/%s?contractId=%s&groupId=%s", testData.includeID, testData.contractID, testData.groupID),
			IncludeID:   testData.includeID,
		}
	}

	newGetIncludeResp := func(testData *testData) *papi.GetIncludeResponse {
		return &papi.GetIncludeResponse{
			Include: papi.Include{
				AssetID:           testData.assetID,
				GroupID:           testData.groupID,
				IncludeID:         testData.includeID,
				ContractID:        testData.contractID,
				IncludeName:       testData.includeName,
				IncludeType:       testData.includeType,
				LatestVersion:     testData.latestVersion,
				StagingVersion:    testData.stagingVersion,
				ProductionVersion: testData.productionVersion,
			},
		}
	}

	newGetIncludeRuleTreeResp := func(testData *testData) *papi.GetIncludeRuleTreeResponse {
		resp := &papi.GetIncludeRuleTreeResponse{
			RuleFormat: testData.ruleFormat,
		}

		var rulesResp papi.GetIncludeRuleTreeResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, path.Join(workdir, "rules_response", testData.rulesPath)), &rulesResp)
		assert.NoError(t, err)

		resp.Rules = rulesResp.Rules
		resp.Errors = rulesResp.Errors
		resp.Warnings = rulesResp.Warnings
		resp.Comments = rulesResp.Comments

		return resp
	}

	simpleRules := papi.RulesUpdate{
		Rules: papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{Name: "caching", Options: papi.RuleOptionsMap{
					"behavior":       "MAX_AGE",
					"mustRevalidate": false, "ttl": "13d"}},
				{Name: "cpCode", Options: papi.RuleOptionsMap{
					"value": map[string]interface{}{
						"id": 1.013931e+06}}},
				{Name: "origin", Options: papi.RuleOptionsMap{
					"cacheKeyHostname":   "ORIGIN_HOSTNAME",
					"compress":           true,
					"enableTrueClientIp": false,
					"forwardHostHeader":  "REQUEST_HOST_HEADER",
					"hostname":           "terraform.prov.test.net",
					"httpPort":           float64(80),
					"httpsPort":          float64(443),
					"originCertificate":  "",
					"originSni":          true,
					"originType":         "CUSTOMER",
					"ports":              "",
					"verificationMode":   "PLATFORM_SETTINGS"}},
			},
			Name: "default",
			Children: []papi.Rules{
				{
					Behaviors: []papi.RuleBehavior{
						{Name: "caching", Options: papi.RuleOptionsMap{
							"behavior":       "MAX_AGE",
							"mustRevalidate": false, "ttl": "13d"}},
						{Name: "cpCode", Options: papi.RuleOptionsMap{
							"value": map[string]interface{}{
								"id": 1.013931e+06}}},
						{Name: "origin", Options: papi.RuleOptionsMap{
							"cacheKeyHostname":   "ORIGIN_HOSTNAME",
							"compress":           true,
							"enableTrueClientIp": false,
							"forwardHostHeader":  "REQUEST_HOST_HEADER",
							"hostname":           "terraform.prov.test.net",
							"httpPort":           float64(80),
							"httpsPort":          float64(443),
							"originCertificate":  "",
							"originSni":          true,
							"originType":         "CUSTOMER",
							"ports":              "",
							"verificationMode":   "PLATFORM_SETTINGS"}},
					},
				},
			},
		},
	}

	simpleRulesWithComment := papi.RulesUpdate{
		Rules:    simpleRules.Rules,
		Comments: "test comment",
	}

	nullRules := papi.RulesUpdate{
		Rules: papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{Name: "cpCode", Options: papi.RuleOptionsMap{
					"value": map[string]interface{}{
						"id":          1.047836e+06,
						"description": "CliTerraformCPCode",
						"name":        "DevExpCliTerraformPapiTest",
						"products":    []interface{}{"Web_App_Accel"},
					}}}},
			Name: "default",
		},
	}

	newUpdateIncludeRuleTreeReq := func(testData *testData) papi.UpdateIncludeRuleTreeRequest {
		unifiedRules := testutils.LoadFixtureString(t, path.Join(workdir, "property-snippets", testData.rulesPath))

		var rules papi.RulesUpdate
		err := json.Unmarshal([]byte(unifiedRules), &rules)
		assert.NoError(t, err)

		return papi.UpdateIncludeRuleTreeRequest{
			ContractID:     testData.contractID,
			GroupID:        testData.groupID,
			IncludeID:      testData.includeID,
			IncludeVersion: testData.latestVersion,
			Rules:          testData.rules,
		}
	}

	newGetIncludeVersionResp := func(testData *testData) *papi.GetIncludeVersionResponse {
		return &papi.GetIncludeVersionResponse{
			IncludeVersion: papi.IncludeVersion{
				StagingStatus:    testData.stagingStatus,
				ProductionStatus: testData.productionStatus,
			},
		}
	}

	expectCreate := func(m *papi.Mock, testData *testData) testutils.MockCalls {
		testData.latestVersion++

		createIncludeCall := m.On("CreateInclude", testutils.MockContext, papi.CreateIncludeRequest{
			GroupID:     testData.groupID,
			ContractID:  testData.contractID,
			ProductID:   testData.productID,
			IncludeName: testData.includeName,
			RuleFormat:  testData.ruleFormat,
			IncludeType: testData.includeType,
		}).Return(newCreateIncludeResp(testData), nil)

		if len(testData.rulesPath) == 0 {
			testData.rulesPath = "default_rules.json"
			return testutils.MockCalls{createIncludeCall}
		}

		updateIncludeRuleTreeCall := m.On("UpdateIncludeRuleTree", testutils.MockContext,
			newUpdateIncludeRuleTreeReq(testData)).Return(&papi.UpdateIncludeRuleTreeResponse{}, nil) // Return argument is ignored

		return testutils.MockCalls{createIncludeCall, updateIncludeRuleTreeCall}
	}

	expectGetIncludeRuleTree := func(m *papi.Mock, testData *testData) *mock.Call {
		call := m.On("GetIncludeRuleTree", testutils.MockContext, papi.GetIncludeRuleTreeRequest{
			ContractID:     testData.contractID,
			GroupID:        testData.groupID,
			IncludeID:      testData.includeID,
			IncludeVersion: testData.latestVersion,
			ValidateRules:  true,
		}).Return(newGetIncludeRuleTreeResp(testData), nil)
		return call
	}

	expectRead := func(m *papi.Mock, testData *testData) testutils.MockCalls {
		getIncludeCall := m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
			ContractID: testData.contractID,
			GroupID:    testData.groupID,
			IncludeID:  includeID,
		}).Return(newGetIncludeResp(testData), nil)

		getIncludeRuleTreeCall := expectGetIncludeRuleTree(m, testData)

		getIncludeVersionCall := m.On("GetIncludeVersion", testutils.MockContext, papi.GetIncludeVersionRequest{
			IncludeID:  includeID,
			Version:    testData.latestVersion,
			ContractID: testData.contractID,
			GroupID:    testData.groupID,
		}).Return(&papi.GetIncludeVersionResponse{
			IncludeVersion: papi.IncludeVersion{
				ProductID: testData.productID,
			}}, nil)

		return testutils.MockCalls{getIncludeCall, getIncludeRuleTreeCall, getIncludeVersionCall}
	}

	expectUpdate := func(m *papi.Mock, testData *testData) testutils.MockCalls {
		getIncludeVersionCall := m.On("GetIncludeVersion", testutils.MockContext, papi.GetIncludeVersionRequest{
			Version:    testData.latestVersion,
			GroupID:    testData.groupID,
			IncludeID:  testData.includeID,
			ContractID: testData.contractID,
		}).Return(newGetIncludeVersionResp(testData), nil)

		calls := testutils.MockCalls{getIncludeVersionCall}

		if testData.stagingVersion != nil || testData.productionVersion != nil {
			version := testData.latestVersion
			testData.latestVersion++

			createIncludeVersionCall := m.On("CreateIncludeVersion", testutils.MockContext, papi.CreateIncludeVersionRequest{
				IncludeID: includeID,
				IncludeVersionRequest: papi.IncludeVersionRequest{
					CreateFromVersion: version,
				},
			}).Return(&papi.CreateIncludeVersionResponse{Version: testData.latestVersion}, nil)

			calls = append(calls, createIncludeVersionCall)
		}

		updateIncludeRuleTreeCall := m.On("UpdateIncludeRuleTree", testutils.MockContext,
			newUpdateIncludeRuleTreeReq(testData)).Return(&papi.UpdateIncludeRuleTreeResponse{}, nil)

		return append(calls, updateIncludeRuleTreeCall)
	}

	expectDelete := func(m *papi.Mock, testData *testData) testutils.MockCalls {
		getIncludeCall := m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
			ContractID: testData.contractID,
			GroupID:    testData.groupID,
			IncludeID:  includeID,
		}).Return(newGetIncludeResp(testData), nil)

		deleteCall := m.On("DeleteInclude", testutils.MockContext, papi.DeleteIncludeRequest{
			GroupID:    testData.groupID,
			IncludeID:  testData.includeID,
			ContractID: testData.contractID,
		}).Return(&papi.DeleteIncludeResponse{}, nil)

		return testutils.MockCalls{getIncludeCall, deleteCall}
	}

	simulateActivation := func(testData *testData, version int, network papi.ActivationNetwork) {
		if network == papi.ActivationNetworkStaging {
			testData.stagingVersion = &version
			testData.stagingStatus = papi.VersionStatusActive
		} else {
			testData.productionVersion = &version
			testData.productionStatus = papi.VersionStatusActive
		}
	}

	simulateDeactivation := func(testData *testData, network papi.ActivationNetwork) {
		if network == papi.ActivationNetworkStaging {
			testData.stagingVersion = nil
			testData.stagingStatus = papi.VersionStatusDeactivated
		} else {
			testData.productionVersion = nil
			testData.productionStatus = papi.VersionStatusDeactivated
		}
	}

	tests := map[string]struct {
		init     func(*papi.Mock, *testData)
		steps    []resource.TestStep
		testData testData
	}{
		"create include - no rules and no warnings": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/default_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
			},
		},
		"create include - with rules": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
			},
		},
		"create include - with rules with comment": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules_with_comment.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRulesWithComment,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_with_comment.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_with_comment.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
			},
		},
		"create include - rules with validation errors": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				testData.rulesPath = "simple_rules_with_errors.json"
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_errors.json", workdir)),
					),
				},
			},
		},
		"create include - rules with validation warnings": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				testData.rulesPath = "simple_rules_with_warnings.json"
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_warnings.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
					),
				},
			},
		},
		"create include - rules with validation errors and warnings": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				testData.rulesPath = "simple_rules_with_errors_and_warnings.json"
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_errors.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_warnings.json", workdir)),
					),
				},
			},
		},
		"create include - server returns no warnings on second apply": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				// first step when server returns validation warnings
				expectCreate(m, testData).Once()
				testData.rulesPath = "simple_rules_with_warnings.json"
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()

				// second step when server returns no validation warnings
				testData.rulesPath = "simple_rules.json"
				expectRead(m, testData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_warnings.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
					),
				},
			},
		},
		"create include - server returns no errors on second apply": {
			testData: testData{
				assetID:     "aid_555",
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				// first step when server returns validation errors
				expectCreate(m, testData).Once()
				testData.rulesPath = "simple_rules_with_errors.json"
				expectRead(m, testData).Times(2)

				// second step when server returns no validation errors
				testData.rulesPath = "simple_rules.json"
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules_errors.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
			},
		},
		"update include - editable version": {
			testData: testData{
				assetID:          "aid_555",
				groupID:          "grp_123",
				productID:        "prd_test",
				includeID:        includeID,
				ruleFormat:       "v2022-06-28",
				contractID:       "ctr_123",
				includeName:      "test_include",
				includeType:      papi.IncludeTypeMicroServices,
				stagingStatus:    papi.VersionStatusInactive,
				productionStatus: papi.VersionStatusInactive,
				rules:            simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)

				expectRead(m, testData).Once()

				testData.rulesPath = "simple_rules.json"
				expectUpdate(m, testData).Once()
				expectRead(m, testData).Times(2)

				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/default_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
			},
		},
		"update include - create new version": {
			testData: testData{
				assetID:          "aid_555",
				groupID:          "grp_123",
				productID:        "prd_test",
				includeID:        includeID,
				ruleFormat:       "v2022-06-28",
				contractID:       "ctr_123",
				includeName:      "test_include",
				includeType:      papi.IncludeTypeMicroServices,
				stagingStatus:    papi.VersionStatusInactive,
				productionStatus: papi.VersionStatusInactive,
				rules:            simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)

				simulateActivation(testData, 1, papi.ActivationNetworkStaging)
				expectRead(m, testData).Once()

				testData.rulesPath = "simple_rules.json"
				expectUpdate(m, testData).Once()

				expectRead(m, testData).Times(1)

				simulateDeactivation(testData, papi.ActivationNetworkStaging)
				expectRead(m, testData).Times(1)

				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/default_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "asset_id", "aid_555"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "2"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_warnings", ""),
					),
				},
			},
		},
		"update include - version is computed": {
			testData: testData{
				groupID:          "grp_123",
				productID:        "prd_test",
				includeID:        includeID,
				ruleFormat:       "v2022-06-28",
				contractID:       "ctr_123",
				includeName:      "test_include",
				includeType:      papi.IncludeTypeMicroServices,
				stagingStatus:    papi.VersionStatusInactive,
				productionStatus: papi.VersionStatusInactive,
				rules:            simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				// Resource create & post-create plan calls
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)

				// Data source create & post-create plan calls
				expectGetIncludeRuleTree(m, testData).Times(2)

				// Resource refresh calls
				simulateActivation(testData, 1, papi.ActivationNetworkStaging)
				expectRead(m, testData).Once()

				// Data source refresh call
				expectGetIncludeRuleTree(m, testData).Times(1)

				// Resource update calls
				testData.rulesPath = "simple_rules.json"
				expectUpdate(m, testData).Once()
				expectRead(m, testData).Once()

				// Data source update call
				expectGetIncludeRuleTree(m, testData).Times(2)

				// Resource post-update plan calls
				simulateDeactivation(testData, papi.ActivationNetworkStaging)
				expectRead(m, testData).Once()

				// Data source post-update call
				expectGetIncludeRuleTree(m, testData).Times(1)

				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_with_ds_create.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/default_rules.json", workdir)),
						resource.TestCheckResourceAttrPair("akamai_property_include.test", "latest_version", "data.akamai_property_include_rules.rules", "version"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_with_ds_update.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "2"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttrPair("akamai_property_include.test", "latest_version", "data.akamai_property_include_rules.rules", "version"),
					),
				},
			},
		},
		"update include - verify staging_version and production_version are known at plan": {
			testData: testData{
				groupID:          "grp_123",
				productID:        "prd_test",
				includeID:        includeID,
				ruleFormat:       "v2022-06-28",
				contractID:       "ctr_123",
				includeName:      "test_include",
				includeType:      papi.IncludeTypeMicroServices,
				stagingStatus:    papi.VersionStatusInactive,
				productionStatus: papi.VersionStatusInactive,
				rules:            simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)

				expectRead(m, testData).Once()

				testData.rulesPath = "simple_rules.json"
				expectUpdate(m, testData).Once()
				expectRead(m, testData).Times(2)

				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("akamai_property_include.test", tfjsonpath.New("staging_version"), knownvalue.StringExact("")),
						plancheck.ExpectKnownValue("akamai_property_include.test", tfjsonpath.New("production_version"), knownvalue.StringExact("")),
						plancheck.ExpectUnknownValue("akamai_property_include.test", tfjsonpath.New("latest_version")),
					},
					},
				},
			},
		},
		"lifecycle with null in returned cpcode - expect no diff on refresh": {
			testData: testData{
				groupID:     "grp_123",
				rulesPath:   "null_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2023-01-05",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       nullRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_null_cpcode.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test_include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", testutils.LoadFixtureStringf(t, "%s/expected/rules_cpcode_null.json", workdir)),
					),
				},
			},
		},
		"import include": {
			testData: testData{
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(3)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_import.tf", workdir),
				},
				{
					ImportState:       true,
					ImportStateId:     "ctr_123:grp_123:inc_123",
					ResourceName:      "akamai_property_include.test",
					ImportStateVerify: true,
				},
			},
		},
		"error importing - incorrect import id": {
			testData: testData{
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test_include",
				includeType: papi.IncludeTypeMicroServices,
				rules:       simpleRules,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/property_include_import.tf", workdir),
				},
				{
					ImportState:   true,
					ImportStateId: "invalid:importID",
					ResourceName:  "akamai_property_include.test",
					ExpectError:   regexp.MustCompile("property include import: invalid import id 'invalid:importID'"),
				},
			},
		},
		"validation errors": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "type" is required, but no definition was found`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/custom_validation_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`Error: expected type to be one of \["MICROSERVICES" "COMMON_SETTINGS"]`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/custom_validation_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`Error: "rule_format" must be of the form vYYYY-MM-DD \(with a leading "v"\)`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/custom_validation_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`Error: "rules" contains an invalid JSON`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/product_id_error.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "product_id" is required during create, but no definition was found`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/rule_format_latest.tf", workdir),
					ExpectError: regexp.MustCompile(`"rule_format" 'latest' is not valid, must be of the form vYYYY-MM-DD`),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/rule_format_blank.tf", workdir),
					ExpectError: regexp.MustCompile(`provided value cannot be blank`),
				},
			},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			if testCase.init != nil {
				testCase.init(client, &testCase.testData)
			}

			useClient(client, &hapi.Mock{}, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    testCase.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestValidatePropertyIncludeName(t *testing.T) {
	invalidNameCharacters := diag.Errorf("a name must only contain letters, numbers, and these characters: . _ -")
	invalidNameLength := diag.Errorf("a name must be longer than 2 characters and shorter than 86 characters")

	tests := map[string]struct {
		propertyName   string
		expectedReturn diag.Diagnostics
	}{
		"name contains only valid characters": {
			propertyName:   "Test_Name.With_Valid-Chars.123",
			expectedReturn: nil,
		},
		"name contains only numbers": {
			propertyName:   "123",
			expectedReturn: nil,
		},
		"name contains only letters": {
			propertyName:   "TestName",
			expectedReturn: nil,
		},
		"name contains invalid char !": {
			propertyName:   "Invalid_Char_!",
			expectedReturn: invalidNameCharacters,
		},
		"name contains invalid char @": {
			propertyName:   "@_Invalid_Char",
			expectedReturn: invalidNameCharacters,
		},
		"name contains invalid spaces": {
			propertyName:   "test name",
			expectedReturn: invalidNameCharacters,
		},
		"name too long (86 chars)": {
			propertyName:   strings.Repeat("a", 86),
			expectedReturn: invalidNameLength,
		},
		"name of max length (85 chars)": {
			propertyName:   strings.Repeat("a", 85),
			expectedReturn: nil,
		},
		"name too short (2 chars)": {
			propertyName:   strings.Repeat("a", 2),
			expectedReturn: invalidNameLength,
		},
		"name of min length (3 char)": {
			propertyName:   strings.Repeat("a", 3),
			expectedReturn: nil,
		},
		"name empty": {
			propertyName:   "",
			expectedReturn: invalidNameLength,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ret := validateNameWithBound(3)(test.propertyName, cty.Path{})

			assert.Equal(t, test.expectedReturn, ret)

		})
	}
}
