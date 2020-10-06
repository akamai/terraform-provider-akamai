package property

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
)

func TestResourceProperty(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*mockpapi)
		expectedAttributes map[string]string
		withError          *regexp.Regexp
	}{
		"simple property create, property does not exist, no rules provided": {
			givenTF: "property_create.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{
						Items: []*papi.Group{
							{GroupID: "grp_1"},
							{GroupID: "grp_2"},
						},
					},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{
						Items: []*papi.Contract{
							{ContractID: "ctr_1"},
							{ContractID: "ctr_2"},
						},
					},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{
							Items: []papi.ProductItem{
								{ProductID: "prd_1"},
								{ProductID: "prd_2"},
							},
						},
					}, nil)
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "property_name",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{
						Items: []papi.SearchItem{},
					},
				}, nil)
				m.On("CreateProperty", mock.Anything, papi.CreatePropertyRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					Property: papi.PropertyCreate{
						ProductID:    "prd_2",
						PropertyName: "property_name",
						RuleFormat:   "rule_format",
					},
				}).Return(&papi.CreatePropertyResponse{
					Response: papi.Response{
						ContractID: "ctr_2",
						GroupID:    "grp_2",
					},
					PropertyID: "prp_123",
				}, nil)
				m.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Response: papi.Response{
						ContractID: "ctr_2",
						GroupID:    "grp_2",
					},
					Property: &papi.Property{
						ContractID:    "ctr_2",
						GroupID:       "grp_2",
						LatestVersion: 1,
						ProductID:     "prd_2",
						PropertyID:    "prp_123",
						PropertyName:  "property_name",
						RuleFormat:    "rule_format",
					},
				}, nil).Once()
				m.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID: "prp_123",
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetPropertyVersionsResponse{
					PropertyID:   "prp_123",
					PropertyName: "property_name",
					ContractID:   "ctr_2",
					GroupID:      "grp_2",
					Version: papi.PropertyVersionGetItem{
						ProductID:        "prd_2",
						ProductionStatus: "ACTIVE",
						PropertyVersion:  1,
						RuleFormat:       "rule_format",
						StagingStatus:    "INACTIVE",
					},
				}, nil)
				m.On("CreatePropertyVersion", mock.Anything, papi.CreatePropertyVersionRequest{
					PropertyID: "prp_123",
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					Version: papi.PropertyVersionCreate{
						CreateFromVersion: 1,
					},
				}).Return(&papi.CreatePropertyVersionResponse{
					PropertyVersion: 2,
				}, nil)
				m.On("GetCPCode", mock.Anything, papi.GetCPCodeRequest{
					CPCodeID:   "cpc_1",
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetCPCodesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					CPCode: papi.CPCode{
						ID:        "cpc_1",
						Name:      "cpc_name",
						ProductID: "prd_2",
					},
				}, nil)
				m.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Response: papi.Response{
						ContractID: "ctr_2",
						GroupID:    "grp_2",
					},
					Property: &papi.Property{
						ContractID:    "ctr_2",
						GroupID:       "grp_2",
						LatestVersion: 2,
						ProductID:     "prd_2",
						PropertyID:    "prp_123",
						PropertyName:  "property_name",
						RuleFormat:    "rule_format",
					},
				}, nil).Once()
				m.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "cpCode",
								Options: papi.RuleOptionsMap{
									"value": papi.RuleOptionsMap{
										"id": "cpc_1",
									},
								},
							},
						},
						Name: "default",
					},
				}).Return(&papi.UpdateRulesResponse{
					ContractID:      "ctr_2",
					GroupID:         "grp_2",
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					RuleFormat:      "rule_format",
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "cpCode",
								Options: papi.RuleOptionsMap{
									"value": papi.RuleOptionsMap{
										"id": 1,
									},
								},
							},
						},
						Name: "default",
					},
				}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_1",
							Domain:       "akamai.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "akamai",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("UpdatePropertyVersionHostnames", mock.Anything, papi.UpdatePropertyVersionHostnamesRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					ContractID:      "ctr_2",
					GroupID:         "grp_2",
					Hostnames: papi.HostnameRequestItems{Items: []papi.Hostname{
						{
							CnameType:      papi.HostnameCnameTypeEdgeHostname,
							EdgeHostnameID: "eh_1",
							CnameFrom:      "cnamefrom",
							CnameTo:        "akamai.edgesuite.net",
						},
					}},
				}).Return(&papi.UpdatePropertyVersionHostnamesResponse{
					ContractID:      "ctr_2",
					GroupID:         "grp_2",
					PropertyID:      "prp_123",
					PropertyVersion: 2,
				}, nil)
				m.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					ContractID:      "ctr_2",
					GroupID:         "grp_2",
				}).Return(&papi.GetRuleTreeResponse{
					Response: papi.Response{
						ContractID: "ctr_2",
						GroupID:    "grp_2",
					},
					PropertyID:      "prp_123",
					PropertyVersion: 2,
					RuleFormat:      "rule_format",
					Rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "cpCode",
								Options: papi.RuleOptionsMap{
									"value": papi.RuleOptionsMap{
										"id": 1,
									},
								},
							},
						},
						Name: "default",
					},
				}, nil)
				m.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
					PropertyID: "prp_123",
				}).Return(&papi.GetPropertyResponse{
					Response: papi.Response{
						ContractID: "ctr_2",
						GroupID:    "grp_2",
					},
					Property: &papi.Property{
						ContractID:    "ctr_2",
						GroupID:       "grp_2",
						LatestVersion: 2,
						ProductID:     "prd_2",
						PropertyID:    "prp_123",
						PropertyName:  "property_name",
						RuleFormat:    "rule_format",
					},
				}, nil)
				m.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
					PropertyID: "prp_123",
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetPropertyResponse{
					Response: papi.Response{
						ContractID: "ctr_2",
						GroupID:    "grp_2",
					},
					Property: &papi.Property{
						ContractID:    "ctr_2",
						GroupID:       "grp_2",
						LatestVersion: 2,
						ProductID:     "prd_2",
						PropertyID:    "prp_123",
						PropertyName:  "property_name",
						RuleFormat:    "rule_format",
					},
				}, nil).Once()
				m.On("RemoveProperty", mock.Anything, papi.RemovePropertyRequest{
					PropertyID: "prp_123",
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.RemovePropertyResponse{}, nil)
			},
			expectedAttributes: map[string]string{
				"product":     "prd_2",
				"contract":    "ctr_2",
				"group":       "grp_2",
				"version":     "2",
				"rules":       compactJSON(loadFixtureBytes("testdata/TestResourceProperty/property_create_rules.json")),
				"id":          "prp_123",
				"cp_code":     "cpc_1",
				"name":        "property_name",
				"rule_format": "rule_format",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_property.prop", k, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{{
						ExpectNonEmptyPlan: true,
						Config:             loadFixtureString(fmt.Sprintf("testdata/TestResourceProperty/%s", test.givenTF)),
						Check:              resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError:        test.withError,
					},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
