package property

import (
	"encoding/json"
	"errors"
	"path"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testDataPropertyIncludeRules struct {
	GroupID      string
	ContractID   string
	Version      int
	IncludeID    string
	Name         string
	RuleFormat   string
	Rules        string
	RuleErrors   string
	RuleWarnings string
	IncludeType  papi.IncludeType
}

var (
	workdir = "./testdata/TestDSPropertyIncludeRules"

	propertyIncludeRules = testDataPropertyIncludeRules{
		ContractID:  "ctr_1",
		GroupID:     "grp_2",
		IncludeID:   "12345",
		Version:     1,
		RuleFormat:  "v2022-06-28",
		Name:        "TestIncludeName",
		IncludeType: "MICROSERVICES",
		Rules:       loadFixtureString("%s/property-snippets/rules_without_errors.json", workdir),
	}

	propertyIncludeRulesWithRuleErrors = func(propertyIncludeRules testDataPropertyIncludeRules, rulesErrors string) testDataPropertyIncludeRules {
		propertyIncludeRules.RuleErrors = rulesErrors
		return propertyIncludeRules
	}

	propertyIncludeRulesWithRuleWarnings = func(propertyIncludeRules testDataPropertyIncludeRules, rulesWarnings string) testDataPropertyIncludeRules {
		propertyIncludeRules.RuleWarnings = rulesWarnings
		return propertyIncludeRules
	}

	propertyIncludeRulesWithRuleWarningsAndErrors = func(propertyIncludeRules testDataPropertyIncludeRules, rulesWarnings, rulesErrors string) testDataPropertyIncludeRules {
		propertyIncludeRules.RuleWarnings = rulesWarnings
		propertyIncludeRules.RuleErrors = rulesErrors
		return propertyIncludeRules
	}

	expectReadPropertyRulesInclude = func(t *testing.T, client *papi.Mock, data testDataPropertyIncludeRules, timesToRun int, withRuleErrors, withRuleWarnings bool, rulesFileName string) {
		getIncludeRuleTreeRequest := papi.GetIncludeRuleTreeRequest{
			ContractID:     data.ContractID,
			GroupID:        data.GroupID,
			IncludeVersion: data.Version,
			IncludeID:      data.IncludeID,
			ValidateRules:  true,
		}
		getIncludeRuleTreeResponse := papi.GetIncludeRuleTreeResponse{
			IncludeName: data.Name,
			RuleFormat:  data.RuleFormat,
			IncludeType: data.IncludeType,
		}

		var ruleResp papi.GetIncludeRuleTreeResponse
		err := json.Unmarshal(loadFixtureBytes(path.Join(workdir, "expected-response", rulesFileName)), &ruleResp)
		assert.NoError(t, err)

		getIncludeRuleTreeResponse.Rules = ruleResp.Rules
		getIncludeRuleTreeResponse.Comments = ruleResp.Comments
		if withRuleErrors && ruleResp.Errors != nil {
			getIncludeRuleTreeResponse.Errors = ruleResp.Errors
		}
		if withRuleWarnings && ruleResp.Warnings != nil {
			getIncludeRuleTreeResponse.Warnings = ruleResp.Warnings
		}

		client.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeRequest).Return(&getIncludeRuleTreeResponse, nil).Times(timesToRun)
	}

	expectGetIncludeRuleTreeError = func(t *testing.T, client *papi.Mock, data testDataPropertyIncludeRules) {
		getIncludeRuleTreeRequest := papi.GetIncludeRuleTreeRequest{
			ContractID:     data.ContractID,
			GroupID:        data.GroupID,
			IncludeVersion: data.Version,
			IncludeID:      data.IncludeID,
			ValidateRules:  true,
		}
		client.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeRequest).Return(nil,
			errors.New("GetIncludeRuleTree response error"))
	}
)

func TestDataPropertyIncludeRules(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *papi.Mock, testDataPropertyIncludeRules)
		mockData   testDataPropertyIncludeRules
		configPath string
		error      *regexp.Regexp
	}{
		"happy path include rules with rule errors": {
			init: func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {
				expectReadPropertyRulesInclude(t, m, testData, 5, true, false, "rules_with_errors.json")
			},
			mockData:   propertyIncludeRulesWithRuleErrors(propertyIncludeRules, loadFixtureString("%s/property-snippets/rule_errors.json", workdir)),
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules.tf",
		},
		"happy path include rules with rules warnings": {
			init: func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {
				expectReadPropertyRulesInclude(t, m, testData, 5, false, true, "rules_with_warnings.json")
			},
			mockData:   propertyIncludeRulesWithRuleWarnings(propertyIncludeRules, loadFixtureString("%s/property-snippets/rule_warnings.json", workdir)),
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules.tf",
		},
		"happy path include rules with rules warnings and errors": {
			init: func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {
				expectReadPropertyRulesInclude(t, m, testData, 5, true, true, "rules_with_errors_and_warnings.json")
			},
			mockData: propertyIncludeRulesWithRuleWarningsAndErrors(propertyIncludeRules, loadFixtureString("%s/property-snippets/rule_warnings.json", workdir),
				loadFixtureString("%s/property-snippets/rule_errors.json", workdir)),
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules.tf",
		},
		"happy path include rules": {
			init: func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {
				expectReadPropertyRulesInclude(t, m, testData, 5, false, false, "rules_without_errors.json")
			},
			mockData:   propertyIncludeRules,
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules.tf",
		},
		"groupID not provided": {
			init:       func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {},
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_no_group_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
		"contractID not provided": {
			init:       func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {},
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_no_contract_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
		"includeID not provided": {
			init:       func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {},
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_no_include_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
		"version not provided": {
			init:       func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {},
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_no_version.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
		"GetIncludeRuleTree response error": {
			init: func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {
				expectGetIncludeRuleTreeError(t, m, testData)
			},
			mockData:   propertyIncludeRulesWithRuleErrors(propertyIncludeRules, loadFixtureString("%s/property-snippets/rule_errors.json", workdir)),
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_api_error.tf",
			error:      regexp.MustCompile("GetIncludeRuleTree response error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			test.init(t, client, test.mockData)
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPath),
							Check:       checkPropertyIncludeRulesAttrs(test.mockData, t),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkPropertyIncludeRulesAttrs(data testDataPropertyIncludeRules, t *testing.T) resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "contract_id", data.ContractID),
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "group_id", data.GroupID),
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "include_id", data.IncludeID),
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "version", strconv.Itoa(data.Version)),
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "name", data.Name),
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "rule_format", data.RuleFormat),
		resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "type", string(data.IncludeType)),
		resource.TestCheckResourceAttrSet("data.akamai_property_include_rules.test", "rules"),
	}

	if len(data.RuleErrors) > 0 {
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttrWith("data.akamai_property_include_rules.test", "rule_errors", func(value string) error {
			assert.JSONEq(t, data.RuleErrors, value)
			return nil
		}))
	}

	if len(data.RuleWarnings) > 0 {
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttrWith("data.akamai_property_include_rules.test", "rule_warnings", func(value string) error {
			assert.JSONEq(t, data.RuleWarnings, value)
			return nil
		}))
	}
	return resource.ComposeAggregateTestCheckFunc(testCheckFuncs...)
}
