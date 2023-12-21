package property

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceEdgeHostname(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname"
	tests := map[string]struct {
		init      func(*papi.Mock, *hapi.Mock)
		withError *regexp.Regexp
		steps     []resource.TestStep
	}{
		"edge hostname with .edgesuite.net, create edge hostname": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
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
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_edgesuite_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test2.edgesuite.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "timeouts.0.default", "55m"),
						resource.TestCheckOutput("edge_hostname", "test2.edgesuite.net"),
					),
				},
			},
		},
		"edge hostname with .edgekey.net, create edge hostname": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
						},
						{
							ID:                "eh_2",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
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
							ID:                "eh_123",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
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
						{
							ID:                "eh_123",
							Domain:            "test.edgekey.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "edgekey.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_edgekey_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_PERFORMANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.edgekey.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "timeouts.#", "0"),
						resource.TestCheckOutput("edge_hostname", "test.edgekey.net"),
					),
				},
			},
		},
		"edge hostname with .akamaized.net, create edge hostname": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
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
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
					),
				},
			},
		},
		"different edge hostname, create": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.aka.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "aka.net.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
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
							ID:                "eh_1",
							Domain:            "test2.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
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
						{
							ID:                "eh_123",
							Domain:            "test.aka.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test.aka",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
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
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.aka.edgesuite.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "use_cases", testutils.LoadFixtureString(t, fmt.Sprintf("%s/use_cases/use_cases_new.json", testDir))),
						resource.TestCheckOutput("edge_hostname", "test.aka.edgesuite.net"),
					),
				},
			},
		},
		"edge hostname exists": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("edgehostname 'test.akamaized.net' already exists"),
				},
			},
		},
		"edge hostname exists - update ip_behavior": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID:    "ctr_2",
					GroupID:       "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
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

				// refresh
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Times(3)

				// 2nd step
				// update
				mh.On("UpdateEdgeHostname", mock.Anything, hapi.UpdateEdgeHostnameRequest{
					DNSZone:           "akamaized.net",
					RecordName:        "test",
					Comments:          "change /ipVersionBehavior to IPV6_IPV4_DUALSTACK",
					StatusUpdateEmail: []string{"hello@akamai.com"},
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ipVersionBehavior",
							Value: "IPV6_IPV4_DUALSTACK",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", mock.Anything, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				// read
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_ipv4.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_update_ip_behavior.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
					),
				},
			},
		},
		"edge hostname - ip_behavior drift": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
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
							ID:                "eh_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
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
						IPVersionBehavior: "IPV4",
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
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Twice()

				// refresh
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE", // drift happens here
						},
						{
							ID:                "eh_2",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Once()

				// 2nd step
				// update
				mh.On("UpdateEdgeHostname", mock.Anything, hapi.UpdateEdgeHostnameRequest{
					DNSZone:           "akamaized.net",
					RecordName:        "test",
					Comments:          "change /ipVersionBehavior to IPV4",
					StatusUpdateEmail: []string{"hello@akamai.com"},
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ipVersionBehavior",
							Value: "IPV4",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", mock.Anything, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "PENDING",
				}, nil).Once()
				mh.On("GetChangeRequest", mock.Anything, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil).Once()

				// read
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_ipv4.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_ipv4_with_email.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
			},
		},
		"error - update ip_behavior to ipv6_performance": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1. call from create method and refresh 2. update ip_behvior to improper value
				// 1st step - create
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
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
						IPVersionBehavior: "IPV4",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)

				// refresh
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Times(3)

				// 2nd step - update
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
				}).Return(nil, errors.New("invalid IP version behavior: valid values are IPV4 and IPV6_IPV4_DUALSTACK; IPV6 and other values aren't currently supported")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_ipv4.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_error_update_ipv6_performance.tf")),
					ExpectError: regexp.MustCompile("invalid IP version behavior: valid values are IPV4 and IPV6_IPV4_DUALSTACK; IPV6 and other values aren't currently supported"),
				},
			},
		},
		"error fetching edge hostnames": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"certificate required for ENHANCED_TLS": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "missing_certificate.tf")),
					ExpectError: regexp.MustCompile("a certificate enrollment ID is required for Enhanced TLS edge hostnames with 'edgekey.net' suffix"),
				},
			},
		},
		"error creating edge hostname": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
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
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error edge hostname not found": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
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
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_net.tf")),
					ExpectError: regexp.MustCompile("unable to find edge hostname"),
				},
			},
		},
		"error on empty product id for creation": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_net_without_product_id.tf")),
					ExpectError: regexp.MustCompile("`product_id` must be specified for creation"),
				},
			},
		},
		"update edge hostname without status email": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
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
							ID:                "eh_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
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
						IPVersionBehavior: "IPV4",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				// refresh
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "eh_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Times(3)

				// 2nd step
				// update
				mh.On("UpdateEdgeHostname", mock.Anything, hapi.UpdateEdgeHostnameRequest{
					DNSZone:    "akamaized.net",
					RecordName: "test",
					Comments:   "change /ipVersionBehavior to IPV6_IPV4_DUALSTACK",
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ipVersionBehavior",
							Value: "IPV6_IPV4_DUALSTACK",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", mock.Anything, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				// read
				mp.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "eh_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "new_akamaized_ipv4.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, "update_no_status_update_email.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "eh_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckNoResourceAttr("akamai_edge_hostname.edgehostname", "status_update_email"),
					),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			clientHapi := &hapi.Mock{}
			test.init(client, clientHapi)
			useClient(client, clientHapi, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
			clientHapi.AssertExpectations(t)
		})
	}
}

func TestResourceEdgeHostnames_WithImport(t *testing.T) {
	expectGetEdgeHostname := func(m *papi.Mock, edgehostID, ContractID, GroupID string) *mock.Call {
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
					Domain:            "test2.akamaized.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test2",
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

	expectGetEdgeHostnames := func(m *papi.Mock, ContractID, GroupID string) *mock.Call {
		return m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
			ContractID: ContractID,
			GroupID:    GroupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "eh_1",
					Domain:            "test1.akamaized.net",
					DomainPrefix:      "test1",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "eh_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	createEdgeHostnames := func(mp *papi.Mock) *mock.Call {
		return mp.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameCreate{
				ProductID:         "prd_2",
				DomainPrefix:      "test",
				DomainSuffix:      "akamaized.net",
				SecureNetwork:     "SHARED_CERT",
				IPVersionBehavior: "IPV4",
			},
		}).Return(&papi.CreateEdgeHostnameResponse{
			EdgeHostnameID: "eh_1",
		}, nil)
	}

	expectGetEdgeHostnamesAfterCreate := func(m *papi.Mock, ContractID, GroupID string) *mock.Call {
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
					DomainPrefix:      "test",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "eh_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}

	t.Run("import existing edgehostname code", func(t *testing.T) {
		client := &papi.Mock{}
		id := "eh_1,1,2"

		expectGetEdgeHostnames(client, "ctr_1", "grp_2").Once()
		createEdgeHostnames(client).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostname(client, "eh_1", "ctr_1", "grp_2").Once()
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						// please note that import does not support product id, that why we only define it in config for creation
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname.tf"),
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
						ImportStateVerifyIgnore: []string{"product_id"},
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
			oldFixture := testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, test.oldPath))
			newFixture := testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", testDir, test.newPath))

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
			expected: testutils.LoadFixtureBytes(t, fmt.Sprintf("%s/%s", testDir, "use_cases1.json")),
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
