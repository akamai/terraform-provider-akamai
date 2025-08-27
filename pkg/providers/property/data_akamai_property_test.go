package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataProperty(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*papi.Mock)
		expectedAttributes map[string]string
		withError          *regexp.Regexp
	}{
		"valid rules, no version provided": {
			givenTF: "no_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Properties: papi.PropertiesItems{Items: []*papi.Property{
						{
							AssetID:           "aid_555",
							ContractID:        "ctr_1",
							GroupID:           "grp_1",
							LatestVersion:     1,
							ProductionVersion: ptr.To(1),
							PropertyID:        "prp_123",
							StagingVersion:    ptr.To(1),
							PropertyType:      ptr.To(""),
						},
					}},
				}, nil)
				m.On("GetRuleTree", testutils.MockContext, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetRuleTreeResponse{
					Response: papi.Response{
						ContractID: "ctr_1",
						GroupID:    "grp_1",
					},
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "beh 1",
							},
						},
						Name:                "rule 1",
						CriteriaMustSatisfy: "all",
					},
				}, nil)
				m.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetPropertyVersionsResponse{
					Version: papi.PropertyVersionGetItem{
						Note:       "note",
						ProductID:  "prd_1",
						RuleFormat: "latest",
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"asset_id":           "aid_555",
				"name":               "property_name",
				"rules":              compactJSON(testutils.LoadFixtureBytes(t, "testdata/TestDataProperty/no_version_rules.json")),
				"contract_id":        "ctr_1",
				"group_id":           "grp_1",
				"latest_version":     "1",
				"note":               "note",
				"product_id":         "prd_1",
				"production_version": "1",
				"property_id":        "prp_123",
				"rule_format":        "latest",
				"staging_version":    "1",
				"property_type":      "",
			},
		},
		"valid rules, with version provided": {
			givenTF: "with_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Properties: papi.PropertiesItems{Items: []*papi.Property{
						{
							AssetID:           "aid_555",
							ContractID:        "ctr_1",
							GroupID:           "grp_1",
							LatestVersion:     1,
							ProductionVersion: ptr.To(2),
							PropertyID:        "prp_123",
							StagingVersion:    ptr.To(3),
							PropertyType:      ptr.To(""),
						},
					}},
				}, nil)
				m.On("GetRuleTree", testutils.MockContext, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetRuleTreeResponse{
					Response: papi.Response{
						ContractID: "ctr_1",
						GroupID:    "grp_1",
					},
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "beh 1",
							},
						},
						Name:                "rule 1",
						CriteriaMustSatisfy: "all",
					},
				}, nil)
				m.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetPropertyVersionsResponse{
					Version: papi.PropertyVersionGetItem{
						Note:       "note",
						ProductID:  "prd_1",
						RuleFormat: "latest",
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"asset_id":           "aid_555",
				"name":               "property_name",
				"rules":              compactJSON(testutils.LoadFixtureBytes(t, "testdata/TestDataProperty/with_version_rules.json")),
				"contract_id":        "ctr_1",
				"group_id":           "grp_1",
				"latest_version":     "2",
				"note":               "note",
				"product_id":         "prd_1",
				"production_version": "2",
				"property_id":        "prp_123",
				"rule_format":        "latest",
				"staging_version":    "3",
				"property_type":      "",
			},
		},
		"valid rules, hostname bucket support enabled": {
			givenTF: "no_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Properties: papi.PropertiesItems{Items: []*papi.Property{
						{
							AssetID:           "aid_555",
							ContractID:        "ctr_1",
							GroupID:           "grp_1",
							LatestVersion:     1,
							ProductionVersion: ptr.To(1),
							PropertyID:        "prp_123",
							StagingVersion:    ptr.To(1),
							PropertyType:      ptr.To("HOSTNAME_BUCKET"),
						},
					}},
				}, nil)
				m.On("GetRuleTree", testutils.MockContext, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetRuleTreeResponse{
					Response: papi.Response{
						ContractID: "ctr_1",
						GroupID:    "grp_1",
					},
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "beh 1",
							},
						},
						Name:                "rule 1",
						CriteriaMustSatisfy: "all",
					},
				}, nil)
				m.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetPropertyVersionsResponse{
					Version: papi.PropertyVersionGetItem{
						Note:       "note",
						ProductID:  "prd_1",
						RuleFormat: "latest",
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"asset_id":           "aid_555",
				"name":               "property_name",
				"rules":              compactJSON(testutils.LoadFixtureBytes(t, "testdata/TestDataProperty/no_version_rules.json")),
				"contract_id":        "ctr_1",
				"group_id":           "grp_1",
				"latest_version":     "1",
				"note":               "note",
				"product_id":         "prd_1",
				"production_version": "1",
				"property_id":        "prp_123",
				"rule_format":        "latest",
				"staging_version":    "1",
				"property_type":      "HOSTNAME_BUCKET",
			},
		},
		"valid rules, no version provided, no staging & production version returned": {
			givenTF: "no_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Properties: papi.PropertiesItems{Items: []*papi.Property{
						{
							AssetID:       "aid_555",
							ContractID:    "ctr_1",
							GroupID:       "grp_1",
							LatestVersion: 1,
							PropertyID:    "prp_123",
							PropertyType:  ptr.To(""),
						},
					}},
				}, nil)
				m.On("GetRuleTree", testutils.MockContext, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetRuleTreeResponse{
					Response: papi.Response{
						ContractID: "ctr_1",
						GroupID:    "grp_1",
					},
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "beh 1",
							},
						},
						Name:                "rule 1",
						CriteriaMustSatisfy: "all",
					},
				}, nil)
				m.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetPropertyVersionsResponse{
					Version: papi.PropertyVersionGetItem{
						Note:       "note",
						ProductID:  "prd_1",
						RuleFormat: "latest",
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"asset_id":           "aid_555",
				"name":               "property_name",
				"rules":              compactJSON(testutils.LoadFixtureBytes(t, "testdata/TestDataProperty/no_version_rules.json")),
				"contract_id":        "ctr_1",
				"group_id":           "grp_1",
				"latest_version":     "1",
				"note":               "note",
				"product_id":         "prd_1",
				"production_version": "0",
				"property_id":        "prp_123",
				"rule_format":        "latest",
				"staging_version":    "0",
				"property_type":      "",
			},
		},
		"error searching for property": {
			givenTF: "with_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
		"error fetching property": {
			givenTF: "with_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
		"property not found": {
			givenTF: "with_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{},
					},
				}, nil)
			},
			withError: regexp.MustCompile("property not found"),
		},
		"error fetching rules": {
			givenTF: "with_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Properties: papi.PropertiesItems{Items: []*papi.Property{
						{
							PropertyID:    "prp_123",
							LatestVersion: 1,
							ContractID:    "ctr_1",
							GroupID:       "grp_1",
						},
					}},
				}, nil)
				m.On("GetRuleTree", testutils.MockContext, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("property rules not found"),
		},
		"error name not provided": {
			givenTF:   "no_name.tf",
			withError: regexp.MustCompile("Missing required argument"),
		},
		"error property version not found": {
			givenTF: "no_version.tf",
			init: func(m *papi.Mock) {
				m.On("SearchProperties", testutils.MockContext, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{
							{
								ContractID: "ctr_1",
								GroupID:    "grp_1",
								PropertyID: "prp_123",
							},
						},
					},
				}, nil)
				m.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Properties: papi.PropertiesItems{Items: []*papi.Property{
						{
							ContractID:        "ctr_1",
							GroupID:           "grp_1",
							LatestVersion:     1,
							ProductionVersion: ptr.To(1),
							PropertyID:        "prp_123",
							StagingVersion:    ptr.To(1),
						},
					}},
				}, nil)
				m.On("GetRuleTree", testutils.MockContext, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(&papi.GetRuleTreeResponse{
					Response: papi.Response{
						ContractID: "ctr_1",
						GroupID:    "grp_1",
					},
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "beh 1",
							},
						},
						Name:                "rule 1",
						CriteriaMustSatisfy: "all",
					},
				}, nil)
				m.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_property.prop", k, v))
			}
			useClient(client, nil, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataProperty/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.withError,
					},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
