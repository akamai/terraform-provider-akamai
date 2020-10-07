package property

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
)

func TestResourceEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*mockpapi)
		expectedAttributes map[string]string
		withError          *regexp.Regexp
	}{
		"edge hostname with .edgesuite.net, create edge hostname": {
			givenTF: "new_edgesuite_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "edgesuite.net",
						SecureNetwork:     "STANDARD_TLS",
						IPVersionBehavior: "IPV6_COMPLIANCE",
						CertEnrollmentID:  123,
						SlotNumber:        123,
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV6_COMPLIANCE",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.edgesuite.net",
			},
		},
		"edge hostname with .edgekey.net, create edge hostname": {
			givenTF: "new_edgekey_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "edgekey.net",
						SecureNetwork:     "ENHANCED_TLS",
						IPVersionBehavior: "IPV6",
						CertEnrollmentID:  123,
						SlotNumber:        123,
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV6",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.edgekey.net",
			},
		},
		"edge hostname with .akamaized.net, create edge hostname": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV4",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV4",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.akamaized.net",
			},
		},
		"different edge hostname, create": {
			givenTF: "new.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test.aka.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "aka.net.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test.aka",
						DomainSuffix:      "edgesuite.net",
						SecureNetwork:     "",
						IPVersionBehavior: "IPV4",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV4",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.aka",
			},
		},
		"edge hostname exists": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_1",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.akamaized.net",
			},
		},
		"group not found": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
					}},
				}, nil)
			},
			withError: regexp.MustCompile("group not found: grp_2"),
		},
		"contract not found": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
					}},
				}, nil)
			},
			withError: regexp.MustCompile("contract not found: ctr_2"),
		},
		"product not found": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
						}},
					}, nil)
			},
			withError: regexp.MustCompile("product not found: prd_2"),
		},
		"error fetching edge hostnames": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
					}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
		"invalid IP version behavior": {
			givenTF: "invalid_ip.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
					}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{}, nil)
			},
			withError: regexp.MustCompile("ipv4 or ipv6 must be specified to create a new Edge Hostname"),
		},
		"certificate required for ENHANCED_TLS": {
			givenTF: "missing_certificate.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
					}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{}, nil)
			},
			withError: regexp.MustCompile("A certificate enrollment ID is required for Enhanced TLS \\(edgekey.net\\) edge hostnames"),
		},
		"error creating edge hostname": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV4",
					},
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", k, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{
						{
							ExpectNonEmptyPlan: true,
							Config:             loadFixtureString(fmt.Sprintf("testdata/TestResourceEdgeHostname/%s", test.givenTF)),
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

func TestResourceEdgeHostnames_WithImport(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*mockpapi)
		expectedAttributes map[string]string
	}{
		"edge hostname with .akamaized.net, create edge hostname": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetGroups", mock.Anything).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{Items: []*papi.Group{
						{GroupID: "grp_1"},
						{GroupID: "grp_2"},
					}},
				}, nil)
				m.On("GetContracts", mock.Anything).Return(&papi.GetContractsResponse{
					Contracts: papi.ContractsItems{Items: []*papi.Contract{
						{ContractID: "ctr_1"},
						{ContractID: "ctr_2"},
					}},
				}, nil)
				m.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "ctr_2"}).
					Return(&papi.GetProductsResponse{
						Products: papi.ProductsItems{Items: []papi.ProductItem{
							{ProductID: "prd_1"},
							{ProductID: "prd_2"},
						}},
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
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV4",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "eh_1",
				}).Return(nil, fmt.Errorf("oops")).Once()
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyHostname,
					Value: "eh_1",
				}).Return(nil, fmt.Errorf("oops")).Once()
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyEdgeHostname,
					Value: "eh_1",
				}).Return(&papi.SearchResponse{
					Versions: papi.SearchItems{Items: []papi.SearchItem{
						{
							PropertyID: "prp_1",
						},
					}},
				}, nil).Once()
				m.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
					PropertyID: "prp_1",
				}).Return(&papi.GetPropertyResponse{
					Property: &papi.Property{
						AccountID:     "acc_1",
						ContractID:    "ctr_2",
						GroupID:       "grp_2",
						LatestVersion: 1,
						PropertyName:  "property 1",
						PropertyID:    "prp_1",
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id": "prp_1",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", k, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{
						{
							ExpectNonEmptyPlan: true,
							Config:             loadFixtureString(fmt.Sprintf("testdata/TestResourceEdgeHostname/%s", test.givenTF)),
						},
						{
							ImportState:       true,
							ImportStateVerify: true,
							ResourceName:      "akamai_edge_hostname.edgehostname",
							ImportStateId:     "eh_1",
							Check:             resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
