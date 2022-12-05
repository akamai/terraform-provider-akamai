package property

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourcePropertyInclude(t *testing.T) {
	type testData struct {
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
		err := json.Unmarshal(loadFixtureBytes(path.Join(workdir, "rules_response", testData.rulesPath)), &rulesResp)
		assert.NoError(t, err)

		resp.Rules = rulesResp.Rules
		resp.Errors = rulesResp.Errors
		resp.Comments = rulesResp.Comments

		return resp
	}

	newUpdateIncludeRuleTreeReq := func(testData *testData) papi.UpdateIncludeRuleTreeRequest {
		unifiedRules := loadFixtureString(path.Join(workdir, "property-snippets", testData.rulesPath))

		var rules papi.RulesUpdate
		err := json.Unmarshal([]byte(unifiedRules), &rules)
		assert.NoError(t, err)

		return papi.UpdateIncludeRuleTreeRequest{
			ContractID:     testData.contractID,
			GroupID:        testData.groupID,
			IncludeID:      testData.includeID,
			IncludeVersion: testData.latestVersion,
			Rules:          rules,
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

	expectCreate := func(m *papi.Mock, testData *testData) test.MockCalls {
		testData.latestVersion++

		createIncludeCall := m.On("CreateInclude", mock.Anything, papi.CreateIncludeRequest{
			GroupID:     testData.groupID,
			ContractID:  testData.contractID,
			ProductID:   testData.productID,
			IncludeName: testData.includeName,
			RuleFormat:  testData.ruleFormat,
			IncludeType: testData.includeType,
		}).Return(newCreateIncludeResp(testData), nil)

		if len(testData.rulesPath) == 0 {
			testData.rulesPath = "default_rules.json"
			return test.MockCalls{createIncludeCall}
		}

		updateIncludeRuleTreeCall := m.On("UpdateIncludeRuleTree", mock.Anything,
			newUpdateIncludeRuleTreeReq(testData)).Return(&papi.UpdateIncludeRuleTreeResponse{}, nil) // Return argument is ignored

		return test.MockCalls{createIncludeCall, updateIncludeRuleTreeCall}
	}

	expectRead := func(m *papi.Mock, testData *testData) test.MockCalls {
		getIncludeCall := m.On("GetInclude", mock.Anything, papi.GetIncludeRequest{
			ContractID: testData.contractID,
			GroupID:    testData.groupID,
			IncludeID:  includeID,
		}).Return(newGetIncludeResp(testData), nil)

		getIncludeRuleTreeCall := m.On("GetIncludeRuleTree", mock.Anything, papi.GetIncludeRuleTreeRequest{
			ContractID:     testData.contractID,
			GroupID:        testData.groupID,
			IncludeID:      testData.includeID,
			IncludeVersion: testData.latestVersion,
			ValidateRules:  true,
		}).Return(newGetIncludeRuleTreeResp(testData), nil)

		return test.MockCalls{getIncludeCall, getIncludeRuleTreeCall}
	}

	expectUpdate := func(m *papi.Mock, testData *testData) test.MockCalls {
		getIncludeVersionCall := m.On("GetIncludeVersion", mock.Anything, papi.GetIncludeVersionRequest{
			Version:    testData.latestVersion,
			GroupID:    testData.groupID,
			IncludeID:  testData.includeID,
			ContractID: testData.contractID,
		}).Return(newGetIncludeVersionResp(testData), nil)

		calls := test.MockCalls{getIncludeVersionCall}

		if testData.stagingVersion != nil || testData.productionVersion != nil {
			version := testData.latestVersion
			testData.latestVersion++

			createIncludeVersionCall := m.On("CreateIncludeVersion", mock.Anything, papi.CreateIncludeVersionRequest{
				IncludeID: includeID,
				IncludeVersionRequest: papi.IncludeVersionRequest{
					CreateFromVersion: version,
				},
			}).Return(&papi.CreateIncludeVersionResponse{Version: testData.latestVersion}, nil)

			calls = append(calls, createIncludeVersionCall)
		}

		updateIncludeRuleTreeCall := m.On("UpdateIncludeRuleTree", mock.Anything,
			newUpdateIncludeRuleTreeReq(testData)).Return(&papi.UpdateIncludeRuleTreeResponse{}, nil)

		return append(calls, updateIncludeRuleTreeCall)
	}

	expectDelete := func(m *papi.Mock, testData *testData) test.MockCalls {
		getIncludeCall := m.On("GetInclude", mock.Anything, papi.GetIncludeRequest{
			ContractID: testData.contractID,
			GroupID:    testData.groupID,
			IncludeID:  includeID,
		}).Return(newGetIncludeResp(testData), nil)

		deleteCall := m.On("DeleteInclude", mock.Anything, papi.DeleteIncludeRequest{
			GroupID:    testData.groupID,
			IncludeID:  testData.includeID,
			ContractID: testData.contractID,
		}).Return(&papi.DeleteIncludeResponse{}, nil)

		return test.MockCalls{getIncludeCall, deleteCall}
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
		"create include - no rules": {
			testData: testData{
				groupID:     "grp_123",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test include",
				includeType: papi.IncludeTypeMicroServices,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/default_rules.json", workdir)),
						resource.TestCheckNoResourceAttr("akamai_property_include.test", "rule_errors"),
					),
				},
			},
		},
		"create include - with rules": {
			testData: testData{
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test include",
				includeType: papi.IncludeTypeMicroServices,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/simple_rules.json", workdir)),
						resource.TestCheckNoResourceAttr("akamai_property_include.test", "rule_errors"),
					),
				},
			},
		},
		"create include - rules with validation errors": {
			testData: testData{
				groupID:     "grp_123",
				rulesPath:   "simple_rules.json",
				productID:   "prd_test",
				includeID:   includeID,
				ruleFormat:  "v2022-06-28",
				contractID:  "ctr_123",
				includeName: "test include",
				includeType: papi.IncludeTypeMicroServices,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				testData.rulesPath = "simple_rules_with_errors.json"
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/simple_rules.json", workdir)),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_errors", loadFixtureString("%s/expected/simple_rules_errors.json", workdir)),
					),
				},
			},
		},
		"update include - editable version": {
			testData: testData{
				groupID:          "grp_123",
				productID:        "prd_test",
				includeID:        includeID,
				ruleFormat:       "v2022-06-28",
				contractID:       "ctr_123",
				includeName:      "test include",
				includeType:      papi.IncludeTypeMicroServices,
				stagingStatus:    papi.VersionStatusInactive,
				productionStatus: papi.VersionStatusInactive,
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
					Config: loadFixtureString("%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/default_rules.json", workdir)),
						resource.TestCheckNoResourceAttr("akamai_property_include.test", "rule_errors"),
					),
				},
				{
					Config: loadFixtureString("%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/simple_rules.json", workdir)),
						resource.TestCheckNoResourceAttr("akamai_property_include.test", "rule_errors"),
					),
				},
			},
		},
		"update incude - create new version": {
			testData: testData{
				groupID:          "grp_123",
				productID:        "prd_test",
				includeID:        includeID,
				ruleFormat:       "v2022-06-28",
				contractID:       "ctr_123",
				includeName:      "test include",
				includeType:      papi.IncludeTypeMicroServices,
				stagingStatus:    papi.VersionStatusInactive,
				productionStatus: papi.VersionStatusInactive,
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
					Config: loadFixtureString("%s/property_include_no_rules.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/default_rules.json", workdir)),
						resource.TestCheckNoResourceAttr("akamai_property_include.test", "rule_errors"),
					),
				},
				{
					Config: loadFixtureString("%s/property_include.tf", workdir),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_include.test", "group_id", "grp_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "name", "test include"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rule_format", "v2022-06-28"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "type", "MICROSERVICES"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "latest_version", "2"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "staging_version", "1"),
						resource.TestCheckResourceAttr("akamai_property_include.test", "production_version", ""),
						resource.TestCheckResourceAttr("akamai_property_include.test", "rules", loadFixtureString("%s/expected/simple_rules.json", workdir)),
						resource.TestCheckNoResourceAttr("akamai_property_include.test", "rule_errors"),
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
				includeName: "test include",
				includeType: papi.IncludeTypeMicroServices,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(3)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("%s/property_include_import.tf", workdir),
				},
				{
					ImportState:             true,
					ImportStateId:           "ctr_123:grp_123:inc_123",
					ResourceName:            "akamai_property_include.test",
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"product_id"},
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
				includeName: "test include",
				includeType: papi.IncludeTypeMicroServices,
			},
			init: func(m *papi.Mock, testData *testData) {
				expectCreate(m, testData).Once()
				expectRead(m, testData).Times(2)
				expectDelete(m, testData).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("%s/property_include_import.tf", workdir),
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
					Config:      loadFixtureString("%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
				},
				{
					Config:      loadFixtureString("%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
				},
				{
					Config:      loadFixtureString("%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
				{
					Config:      loadFixtureString("%s/validation_required_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "type" is required, but no definition was found`),
				},
				{
					Config:      loadFixtureString("%s/custom_validation_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`Error: expected type to be one of \[MICROSERVICES COMMON_SETTINGS\]`),
				},
				{
					Config:      loadFixtureString("%s/custom_validation_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`Error: "rule_format" must be of the form vYYYY-MM-DD \(with a leading "v"\)`),
				},
				{
					Config:      loadFixtureString("%s/custom_validation_errors.tf", workdir),
					ExpectError: regexp.MustCompile(`Error: "rules" contains an invalid JSON`),
				},
				{
					Config:      loadFixtureString("%s/product_id_error.tf", workdir),
					ExpectError: regexp.MustCompile(`The argument "product_id" is required during create, but no definition was found`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			if test.init != nil {
				test.init(client, &test.testData)
			}

			useClient(client, &hapi.Mock{}, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
