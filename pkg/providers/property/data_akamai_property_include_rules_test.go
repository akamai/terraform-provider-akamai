package property

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testDataPropertyIncludeRules struct {
	GroupID     string
	ContractID  string
	Version     int
	IncludeID   string
	Name        string
	RuleFormat  string
	Rules       string
	RuleErrors  string
	IncludeType papi.IncludeType
}

var (
	workdir = "./testdata/TestDSPropertyIncludeRules"

	propertyIncludeRulesWithoutRuleErrors = testDataPropertyIncludeRules{
		ContractID:  "ctr_1",
		GroupID:     "grp_2",
		IncludeID:   "12345",
		Version:     1,
		RuleFormat:  "v2022-06-28",
		Name:        "TestIncludeName",
		IncludeType: "MICROSERVICES",
		Rules:       loadFixtureString("%s/property-snippets/rules_without_errors.json", workdir),
	}

	propertyIncludeRulesWithRuleErrors = testDataPropertyIncludeRules{
		ContractID:  "ctr_1",
		GroupID:     "grp_2",
		IncludeID:   "12345",
		RuleFormat:  "v2022-06-28",
		Version:     1,
		Name:        "TestIncludeName",
		IncludeType: "MICROSERVICES",
		Rules:       loadFixtureString("%s/property-snippets/rules_with_errors.json", workdir),
		RuleErrors:  loadFixtureString("%s/property-snippets/rule_errors.json", workdir),
	}

	expectReadPropertyRulesInclude = func(t *testing.T, client *papi.Mock, data testDataPropertyIncludeRules, timesToRun int, withRuleErrors bool) {
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
		err := json.Unmarshal(loadFixtureBytes(path.Join(workdir, "expected-response", provideIncludesRulesJSONFileName(withRuleErrors))), &ruleResp)
		assert.NoError(t, err)

		getIncludeRuleTreeResponse.Rules = ruleResp.Rules
		getIncludeRuleTreeResponse.Comments = ruleResp.Comments
		if withRuleErrors && ruleResp.Errors != nil {
			getIncludeRuleTreeResponse.Errors = ruleResp.Errors
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
				expectReadPropertyRulesInclude(t, m, testData, 5, true)
			},
			mockData:   propertyIncludeRulesWithRuleErrors,
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_with_errors.tf",
		},
		"happy path include rules without rules errors": {
			init: func(t *testing.T, m *papi.Mock, testData testDataPropertyIncludeRules) {
				expectReadPropertyRulesInclude(t, m, testData, 5, false)
			},
			mockData:   propertyIncludeRulesWithoutRuleErrors,
			configPath: "./testdata/TestDSPropertyIncludeRules/property_include_rules_without_errors.tf",
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
			mockData:   propertyIncludeRulesWithRuleErrors,
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
							Check:       checkPropertyIncludeRulesAttrs(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkPropertyIncludeRulesAttrs(data testDataPropertyIncludeRules) resource.TestCheckFunc {
	var testCheckFuncs []resource.TestCheckFunc
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "contract_id", data.ContractID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "group_id", data.GroupID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "include_id", data.IncludeID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "version", strconv.Itoa(data.Version)))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "name", data.Name))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "rule_format", data.RuleFormat))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_rules.test", "type", string(data.IncludeType)))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttrSet("data.akamai_property_include_rules.test", "rules"))
	if len(data.RuleErrors) > 0 {
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttrWith("data.akamai_property_include_rules.test", "rule_errors", func(value string) error {
			var providedRuleErrors, expectedRuleErrors []papi.Error
			err := json.Unmarshal([]byte(data.RuleErrors), &providedRuleErrors)
			if err != nil {
				return fmt.Errorf("problem with unmarshal JSON")
			}
			err = json.Unmarshal([]byte(value), &expectedRuleErrors)
			if err != nil {
				return fmt.Errorf("problem with unmarshal JSON")
			}
			if !reflect.DeepEqual(providedRuleErrors, expectedRuleErrors) {
				return fmt.Errorf("two JSONs not equal")
			}
			return nil
		}))
	}
	return resource.ComposeAggregateTestCheckFunc(testCheckFuncs...)
}

func provideIncludesRulesJSONFileName(withRuleErrors bool) string {
	if withRuleErrors {
		return "rules_with_errors.json"
	}
	return "rules_without_errors.json"
}
