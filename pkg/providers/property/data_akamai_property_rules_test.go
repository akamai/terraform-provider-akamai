package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDSPropertyRulesRead(t *testing.T) {
	t.Run("get datasource property rules", func(t *testing.T) {
		client := &papi.Mock{}
		mockImpl := func(m *papi.Mock) {
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(&papi.GetPropertyVersionsResponse{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				Version: papi.PropertyVersionGetItem{
					PropertyVersion: 1,
				},
			}, nil)
			m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				ContractID:      "ctr_2",
				GroupID:         "grp_2",
				PropertyID:      "prp_2",
				PropertyVersion: 1,
				ValidateRules:   true,
				ValidateMode:    papi.RuleValidateModeFull,
			}).Return(&papi.GetRuleTreeResponse{
				Rules: papi.Rules{
					Name: "some rule tree",
				},
				RuleFormat: "latest",
				Response: papi.Response{
					Errors: []*papi.Error{
						{
							Title: "some error",
						},
					},
				},
			}, nil)
		}
		mockImpl(client)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/ds_property_rules.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "id", "prp_2"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "property_id", "prp_2"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "group_id", "grp_2"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "contract_id", "ctr_2"),
							resource.TestCheckResourceAttrSet("data.akamai_property_rules.rules", "rules"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "rule_format", "latest"),
							resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "errors", `[{"type":"","title":"some error","detail":""}]`),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("get datasource property rules with rule format", func(t *testing.T) {
		tests := map[string]struct {
			configFile         string
			expectedRuleFormat string
		}{
			"versioned": {
				configFile:         "testdata/TestDSPropertyRules/with_versioned_rule_format.tf",
				expectedRuleFormat: "v2015-08-17",
			},
			"latest": {
				configFile:         "testdata/TestDSPropertyRules/with_latest_rule_format.tf",
				expectedRuleFormat: "latest",
			},
		}

		mockImpl := func(m *papi.Mock, ruleFormat string) {
			m.On("GetRuleFormats", mock.Anything).Return(&papi.GetRuleFormatsResponse{
				RuleFormats: papi.RuleFormatItems{
					Items: []string{
						"latest",
						"v2021-09-22",
						"v2016-11-15",
						"v2015-08-17",
					},
				},
			}, nil)
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(&papi.GetPropertyVersionsResponse{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				Version: papi.PropertyVersionGetItem{
					PropertyVersion: 1,
				},
			}, nil)
			m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				ContractID:      "ctr_2",
				GroupID:         "grp_2",
				PropertyID:      "prp_2",
				PropertyVersion: 1,
				ValidateRules:   true,
				ValidateMode:    papi.RuleValidateModeFull,
				RuleFormat:      ruleFormat,
			}).Return(&papi.GetRuleTreeResponse{
				RuleFormat: ruleFormat,
			}, nil)
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				client := &papi.Mock{}
				mockImpl(client, test.expectedRuleFormat)

				useClient(client, nil, func() {
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
						Steps: []resource.TestStep{
							{
								Config: testutils.LoadFixtureString(t, test.configFile),
								Check: resource.ComposeAggregateTestCheckFunc(
									resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "rule_format", test.expectedRuleFormat),
								),
							},
						},
					})
				})
				client.AssertExpectations(t)
			})
		}
	})
	t.Run("error getting datasource property rules with invalid rule format", func(t *testing.T) {
		client := &papi.Mock{}
		mockImpl := func(m *papi.Mock) {
			m.On("GetRuleFormats", mock.Anything).Return(&papi.GetRuleFormatsResponse{
				RuleFormats: papi.RuleFormatItems{
					Items: []string{
						"latest",
						"v2021-09-22",
					},
				},
			}, nil)
		}
		mockImpl(client)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/with_versioned_rule_format.tf"),
						ExpectError: regexp.MustCompile("given 'rule_format' is not supported: \"v2015-08-17\""),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error getting rule formats", func(t *testing.T) {
		client := &papi.Mock{}
		mockImpl := func(m *papi.Mock) {
			m.On("GetRuleFormats", mock.Anything).Return(nil, fmt.Errorf("oops"))
		}
		mockImpl(client)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/with_versioned_rule_format.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("group_id is required with contract_id", func(t *testing.T) {
		client := &papi.Mock{}
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/missing_group_id.tf"),
						ExpectError: regexp.MustCompile("\"contract_id\": all of `contract_id,group_id` must be specified"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("contract_id is required with group_id", func(t *testing.T) {
		client := &papi.Mock{}
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/missing_contract_id.tf"),
						ExpectError: regexp.MustCompile("\"group_id\": all of `contract_id,group_id` must be specified"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("contract_id cannot be empty", func(t *testing.T) {
		client := &papi.Mock{}
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/empty_contract_id.tf"),
						ExpectError: regexp.MustCompile(`provided value cannot be blank((.|\n)*)contract_id = ""`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("group_id cannot be empty", func(t *testing.T) {
		client := &papi.Mock{}
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/empty_group_id.tf"),
						ExpectError: regexp.MustCompile(`provided value cannot be blank((.|\n)*)group_id += ""`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error fetching latest version", func(t *testing.T) {
		client := &papi.Mock{}
		mockImpl := func(m *papi.Mock) {
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(nil, fmt.Errorf("fetching latest version")).Once()
		}
		mockImpl(client)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/ds_property_rules.tf"),
						ExpectError: regexp.MustCompile("fetching latest version"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error fetching rules", func(t *testing.T) {
		client := &papi.Mock{}
		mockImpl := func(m *papi.Mock) {
			m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				PropertyID: "prp_2",
			}).Return(&papi.GetPropertyVersionsResponse{
				ContractID: "ctr_2",
				GroupID:    "grp_2",
				Version: papi.PropertyVersionGetItem{
					PropertyVersion: 1,
				},
			}, nil).Once()
			m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				ContractID:      "ctr_2",
				GroupID:         "grp_2",
				PropertyID:      "prp_2",
				PropertyVersion: 1,
				ValidateRules:   true,
				ValidateMode:    papi.RuleValidateModeFull,
			}).Return(nil, fmt.Errorf("fetching rule tree")).Once()
		}
		mockImpl(client)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/ds_property_rules.tf"),
						ExpectError: regexp.MustCompile("fetching rule tree"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestDSPropertyRulesRead_Fail(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
		Steps: []resource.TestStep{{
			Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRules/always_fails.tf"),
			ExpectError: regexp.MustCompile(`Error: provided value cannot be blank`),
		}},
	})
}
