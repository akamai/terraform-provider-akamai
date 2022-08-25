package property

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func TestResourceEdgeHostname(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname"
	tests := map[string]struct {
		init      func(*mockpapi, *mockhapi)
		withError *regexp.Regexp
		steps     []resource.TestStep
	}{
		"edge hostname with .edgesuite.net, create edge hostname": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
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
				}, nil).Once()
				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test2",
						DomainSuffix:      "edgesuite.net",
						SecureNetwork:     "STANDARD_TLS",
						IPVersionBehavior: "IPV6_COMPLIANCE",
						CertEnrollmentID:  123,
						SlotNumber:        123,
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test2.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:           "eh_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_123",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_edgesuite_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test2.edgesuite.net"),
						resource.TestCheckOutput("edge_hostname", "test2.edgesuite.net"),
					),
				},
			},
		},
		"edge hostname with .edgekey.net, create edge hostname": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "edgekey.net",
						SecureNetwork:     "ENHANCED_TLS",
						IPVersionBehavior: "IPV6_PERFORMANCE",
						CertEnrollmentID:  123,
						SlotNumber:        123,
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.edgesuite.net",
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
						{
							ID:           "eh_123",
							Domain:       "test.edgekey.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "edgekey.net",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_edgekey_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_PERFORMANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.edgekey.net"),
						resource.TestCheckOutput("edge_hostname", "test.edgekey.net"),
					),
				},
			},
		},
		"edge hostname with .akamaized.net, create edge hostname": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
					),
				},
			},
		},
		"different edge hostname, create": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.aka.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "aka.net.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test.aka",
						DomainSuffix:      "edgesuite.net",
						SecureNetwork:     "",
						IPVersionBehavior: "IPV4",
						UseCases: []papi.UseCase{
							{
								UseCase: "Download_Mode",
								Option:  "BACKGROUND",
								Type:    "GLOBAL",
							},
						},
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
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
						{
							ID:           "eh_123",
							Domain:       "test.aka.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test.aka",
							DomainSuffix: "edgesuite.net",
							UseCases: []papi.UseCase{
								{
									UseCase: "Download_Mode",
									Option:  "BACKGROUND",
									Type:    "GLOBAL",
								},
							},
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.aka.edgesuite.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "use_cases", loadFixtureString(fmt.Sprintf("%s/use_cases/use_cases_new.json", testDir))),
						resource.TestCheckOutput("edge_hostname", "test.aka.edgesuite.net"),
					),
				},
			},
		},
		"edge hostname exists": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
					),
				},
			},
		},
		"edge hostname exists - update ip_behavior": {
			init: func(mp *mockpapi, mh *mockhapi) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil).Times(3)

				// refresh
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil).Once()

				// 2nd step
				// update
				mh.On("UpdateEdgeHostname", mock.Anything, hapi.UpdateEdgeHostnameRequest{
					DNSZone:           "akamaized.net",
					RecordName:        "test",
					Comments:          "change /ipVersionBehavior to IPV6_PERFORMANCE",
					StatusUpdateEmail: []string{"hello@akamai.com"},
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ipVersionBehavior",
							Value: "IPV6_PERFORMANCE",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{}, nil).Once()

				// read
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_update_ip_behavior.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_PERFORMANCE"),
					),
				},
			},
		},
		"error fetching edge hostnames": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"certificate required for ENHANCED_TLS": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/%s", testDir, "missing_certificate.tf")),
					ExpectError: regexp.MustCompile("a certificate enrollment ID is required for Enhanced TLS edge hostnames with 'edgekey.net' suffix"),
				},
			},
		},
		"error creating edge hostname": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
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
				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error edge hostname not found": {
			init: func(mp *mockpapi, mh *mockhapi) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
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
				}, nil).Twice()
				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("unable to find edge hostname"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			clientHapi := &mockhapi{}
			test.init(client, clientHapi)
			useClient(client, clientHapi, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps:     test.steps,
				})
			})
			client.AssertExpectations(t)
			clientHapi.AssertExpectations(t)
		})
	}
}

func TestResourceEdgeHostnames_WithImport(t *testing.T) {
	expectGetEdgeHostname := func(m *mockpapi, edgehostID, ContractID, GroupID string) *mock.Call {
		return m.On("GetEdgeHostname", mock.Anything, papi.GetEdgeHostnameRequest{
			EdgeHostnameID: edgehostID,
			ContractID:     ContractID,
			GroupID:        GroupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameGetItem{
				ID:                "eh_1",
				Domain:            "test.akamaized.net",
				ProductID:         "prd_2",
				DomainPrefix:      "test",
				DomainSuffix:      "akamaized.net",
				IPVersionBehavior: "IPV4",
			},
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "eh_1",
					Domain:            "test.akamaized.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "eh_2",
					Domain:            "test3.edgesuite.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}

	expectGetEdgeHostnames := func(m *mockpapi, ContractID, GroupID string) *mock.Call {
		return m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
			ContractID: ContractID,
			GroupID:    GroupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "eh_1",
					Domain:            "test.akamaized.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "eh_2",
					Domain:            "test3.edgesuite.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}

	t.Run("import existing edgehostname code", func(t *testing.T) {
		client := &mockpapi{}
		id := "eh_1,1,2"

		expectGetEdgeHostname(client, "eh_1", "ctr_1", "grp_2")
		expectGetEdgeHostnames(client, "ctr_1", "grp_2")
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResourceEdgeHostname/import_edgehostname.tf"),
					},
					{
						Config:      loadFixtureString("testdata/TestResourceEdgeHostname/import_edgehostname.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "eh_1", rs.Attributes["id"])
							return nil
						},
						ImportStateId:           id,
						ResourceName:            "akamai_edge_hostname.importedgehostname",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"ip_behavior"},
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestFindEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		hostnames papi.EdgeHostnameItems
		domain    string
		expected  *papi.EdgeHostnameGetItem
		withError error
	}{
		"edge hostname found, different domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.edgesuite.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "edgesuite.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.edgesuite.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "edgesuite.net",
			},
		},
		"edge hostname found, edgesuite domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.edgesuite.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "edgesuite.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain.edgesuite.net",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.edgesuite.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "edgesuite.net",
			},
		},
		"edge hostname found, edgekey domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.edgekey.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "edgekey.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain.edgekey.net",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.edgekey.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "edgekey.net",
			},
		},
		"edge hostname found, akamaized domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.akamaized.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "akamaized.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain.akamaized.net",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.akamaized.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "akamaized.net",
			},
		},
		"edge hostname not found": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.akamaized.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "akamaized.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain:    "other.domain.akamaized.net",
			withError: ErrEdgeHostnameNotFound,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := findEdgeHostname(test.hostnames, test.domain)
			if test.withError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, test.withError), "want: %v; got: %v", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestDiffSuppressEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		old, new string
		expected bool
	}{
		"equal domains": {
			old:      "test.com",
			new:      "test.com",
			expected: true,
		},
		"domain defaulting to edgesuite.net": {
			old:      "test.com.edgesuite.net",
			new:      "test.com",
			expected: true,
		},
		"different domains": {
			old:      "test.com.akamaized.net",
			new:      "test1.com.akamaized.net",
			expected: false,
		},
		"case insensitive domains": {
			old:      "test.com.akamaized.net",
			new:      "Test.com.akamaized.net",
			expected: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, diffSuppressEdgeHostname("", test.old, test.new, nil))
		})
	}
}

func TestSuppressEdgeHostnameUseCases(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname/use_cases"
	tests := map[string]struct {
		oldPath, newPath string
		expected         bool
	}{
		"equal use cases in the same order": {
			oldPath:  "use_cases1.json",
			newPath:  "use_cases2.json",
			expected: true,
		},
		"equal use cases in different order": {
			oldPath:  "use_cases1.json",
			newPath:  "use_cases1_mixed.json",
			expected: true,
		},
		"not equal use cases": {
			oldPath:  "use_cases1.json",
			newPath:  "use_cases3.json",
			expected: false,
		},
		"error unmarshalling new": {
			oldPath:  "use_cases1.json",
			newPath:  "invalid_use_cases.json",
			expected: false,
		},
		"error unmarshalling old": {
			oldPath:  "invalid_use_cases.json",
			newPath:  "use_cases1.json",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oldFixture := loadFixtureString(fmt.Sprintf("%s/%s", testDir, test.oldPath))
			newFixture := loadFixtureString(fmt.Sprintf("%s/%s", testDir, test.newPath))

			assert.Equal(t, test.expected, suppressEdgeHostnameUseCases("", oldFixture, newFixture, nil))
		})
	}
}

func TestConvertingUseCases2JSON(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname/use_cases"
	tests := map[string]struct {
		useCases []papi.UseCase
		expected []byte
	}{
		"no use cases": {
			useCases: []papi.UseCase{},
			expected: []byte{},
		},
		"two use cases": {
			useCases: []papi.UseCase{
				{
					Option:  "BACKGROUND",
					Type:    "GLOBAL",
					UseCase: "Download_Mode",
				},
				{
					Option:  "FOREGROUND",
					Type:    "GLOBAL",
					UseCase: "Download_Mode",
				},
			},
			expected: loadFixtureBytes(fmt.Sprintf("%s/%s", testDir, "use_cases1.json")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useCasesJSON, err := useCases2JSON(test.useCases)
			assert.NoError(t, err)

			if len(useCasesJSON) > 0 {
				expected := new(bytes.Buffer)
				err = json.Compact(expected, test.expected)
				assert.NoError(t, err)

				actual := new(bytes.Buffer)
				err = json.Compact(actual, useCasesJSON)
				assert.NoError(t, err)

				assert.Equal(t, expected.String(), actual.String())
			} else {
				assert.Equal(t, test.expected, useCasesJSON)
			}
		})
	}
}
