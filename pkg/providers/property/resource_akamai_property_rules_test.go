package property

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
)

func TestResourcePropertyRules(t *testing.T) {

	t.Run("new rules with latest version with update", func(t *testing.T) {
		newRules := papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{
					Name: "beh_1",
				},
			},
			Name:                "default",
			Options:             papi.RuleOptions{IsSecure: true},
			CriteriaMustSatisfy: "all",
		}
		client := new(mockpapi)
		updatedRules := papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{
					Name: "beh_2",
				},
			},
			Name:                "updated",
			Options:             papi.RuleOptions{IsSecure: true},
			CriteriaMustSatisfy: "all",
		}
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  1,
								},
							}, nil).Times(4)
							client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
								Rules:           papi.RulesUpdate{Rules: newRules},
							}).Return(&papi.UpdateRulesResponse{
								AccountID:       "acc_1",
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
								Errors: []papi.RuleError{
									{
										Type:  "generic",
										Title: "some error",
									},
								},
							}, nil).Once()
							client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
							}).Return(&papi.GetRuleTreeResponse{
								Response: papi.Response{
									AccountID:  "acc_1",
									ContractID: "ctr_1",
									GroupID:    "grp_1",
									Errors: []*papi.Error{
										{
											Type:  "generic",
											Title: "some error",
										},
									},
								},
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
							}, nil).Times(3)
						},
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_create.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  1,
								},
							}, nil).Once()
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  2,
								},
							}, nil).Twice()
							client.On("CreatePropertyVersion", mock.Anything, papi.CreatePropertyVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								Version: papi.PropertyVersionCreate{
									CreateFromVersion: 1,
								},
							}).Return(&papi.CreatePropertyVersionResponse{
								PropertyVersion: 2,
							}, nil).Once()
							client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 2,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
								Rules:           papi.RulesUpdate{Rules: updatedRules},
							}).Return(&papi.UpdateRulesResponse{
								AccountID:       "acc_1",
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								PropertyID:      "prp_1",
								PropertyVersion: 2,
								Rules:           newRules,
								Errors: []papi.RuleError{
									{
										Type:  "generic",
										Title: "some error",
									},
								},
							}, nil).Once()
							client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 2,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
							}).Return(&papi.GetRuleTreeResponse{
								Response: papi.Response{
									AccountID:  "acc_1",
									ContractID: "ctr_1",
									GroupID:    "grp_1",
								},
								PropertyID:      "prp_1",
								PropertyVersion: 2,
								Rules:           updatedRules,
							}, nil).Twice()
						},
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "2"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_update.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("new rules, diff suppressed", func(t *testing.T) {
		newRules := papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{
					Name: "beh_1",
				},
			},
			Name:                "default",
			Options:             papi.RuleOptions{IsSecure: true},
			CriteriaMustSatisfy: "all",
		}
		client := new(mockpapi)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  1,
								},
							}, nil)
							client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
								Rules:           papi.RulesUpdate{Rules: newRules},
							}).Return(&papi.UpdateRulesResponse{
								AccountID:       "acc_1",
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
								Errors: []papi.RuleError{
									{
										Type:  "generic",
										Title: "some error",
									},
								},
							}, nil).Once()
							client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
							}).Return(&papi.GetRuleTreeResponse{
								Response: papi.Response{
									AccountID:  "acc_1",
									ContractID: "ctr_1",
									GroupID:    "grp_1",
									Errors: []*papi.Error{
										{
											Type:  "generic",
											Title: "some error",
										},
									},
								},
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
							}, nil)
						},
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_create.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_create.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching latest version on create", func(t *testing.T) {
		client := new(mockpapi)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(nil, fmt.Errorf("oops")).Once()
						},
						Config:      loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("invalid rules JSON", func(t *testing.T) {
		client := new(mockpapi)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResourcePropertyRules/invalid_rules_json.tf"),
						ExpectError: regexp.MustCompile("invalid JSON"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("new contract ID on update", func(t *testing.T) {
		client := new(mockpapi)
		newRules := papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{
					Name: "beh_1",
				},
			},
			Name:                "default",
			Options:             papi.RuleOptions{IsSecure: true},
			CriteriaMustSatisfy: "all",
		}
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  1,
								},
							}, nil).Times(4)
							client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
								Rules:           papi.RulesUpdate{Rules: newRules},
							}).Return(&papi.UpdateRulesResponse{
								AccountID:       "acc_1",
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
								Errors: []papi.RuleError{
									{
										Type:  "generic",
										Title: "some error",
									},
								},
							}, nil).Once()
							client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
							}).Return(&papi.GetRuleTreeResponse{
								Response: papi.Response{
									AccountID:  "acc_1",
									ContractID: "ctr_1",
									GroupID:    "grp_1",
								},
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
							}, nil).Times(3)
						},
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_create.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResourcePropertyRules/update_new_contract.tf"),
						ExpectError: regexp.MustCompile("contract_id field is immutable and cannot be updated"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("new group ID on update", func(t *testing.T) {
		client := new(mockpapi)
		newRules := papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{
					Name: "beh_1",
				},
			},
			Name:                "default",
			Options:             papi.RuleOptions{IsSecure: true},
			CriteriaMustSatisfy: "all",
		}
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  1,
								},
							}, nil).Times(4)
							client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
								Rules:           papi.RulesUpdate{Rules: newRules},
							}).Return(&papi.UpdateRulesResponse{
								AccountID:       "acc_1",
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
								Errors: []papi.RuleError{
									{
										Type:  "generic",
										Title: "some error",
									},
								},
							}, nil).Once()
							client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
							}).Return(&papi.GetRuleTreeResponse{
								Response: papi.Response{
									AccountID:  "acc_1",
									ContractID: "ctr_1",
									GroupID:    "grp_1",
								},
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
							}, nil).Times(3)
						},
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_create.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResourcePropertyRules/update_new_group.tf"),
						ExpectError: regexp.MustCompile("group_id field is immutable and cannot be updated"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("new property ID on update", func(t *testing.T) {
		client := new(mockpapi)
		newRules := papi.Rules{
			Behaviors: []papi.RuleBehavior{
				{
					Name: "beh_1",
				},
			},
			Name:                "default",
			Options:             papi.RuleOptions{IsSecure: true},
			CriteriaMustSatisfy: "all",
		}
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
								PropertyID: "prp_1",
								ContractID: "ctr_1",
								GroupID:    "grp_1",
							}).Return(&papi.GetPropertyVersionsResponse{
								PropertyID:   "prp_1",
								PropertyName: "property 1",
								AccountID:    "acc_1",
								ContractID:   "ctr_1",
								GroupID:      "grp_1",
								Version: papi.PropertyVersionGetItem{
									ProductID:        "prd_1",
									ProductionStatus: "ACTIVE",
									PropertyVersion:  1,
								},
							}, nil).Times(4)
							client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
								Rules:           papi.RulesUpdate{Rules: newRules},
							}).Return(&papi.UpdateRulesResponse{
								AccountID:       "acc_1",
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
								Errors: []papi.RuleError{
									{
										Type:  "generic",
										Title: "some error",
									},
								},
							}, nil).Once()
							client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								ContractID:      "ctr_1",
								GroupID:         "grp_1",
								ValidateRules:   true,
							}).Return(&papi.GetRuleTreeResponse{
								Response: papi.Response{
									AccountID:  "acc_1",
									ContractID: "ctr_1",
									GroupID:    "grp_1",
								},
								PropertyID:      "prp_1",
								PropertyVersion: 1,
								Rules:           newRules,
							}, nil).Times(3)
						},
						Config: loadFixtureString("testdata/TestResourcePropertyRules/latest_version_create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "rules", compactJSON(loadFixtureBytes("testdata/TestResourcePropertyRules/rules_create.json"))),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "property_id", "prp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_property_rules.rules", "contract_id", "ctr_1"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResourcePropertyRules/update_new_property.tf"),
						ExpectError: regexp.MustCompile("property_id field is immutable and cannot be updated"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

}
