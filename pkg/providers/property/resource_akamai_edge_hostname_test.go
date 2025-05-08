package property

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/internal/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceEdgeHostname(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname"

	EgdeHostnameCreatePollInterval = time.Microsecond

	tests := map[string]struct {
		init      func(*papi.Mock, *hapi.Mock)
		withError *regexp.Regexp
		steps     []resource.TestStep
	}{
		"edge hostname with .edgesuite.net, create edge hostname": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "ehn_123",
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "ehn_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_456",
							Domain:            "test2.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:           "ehn_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "ehn_123",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_edgesuite_net.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
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
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test3.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_edgekey_net.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
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
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
					),
				},
			},
		},
		"edge hostname with .akamaized.net, create edge hostname with ttl": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil)
				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
					UseDefaultTTL:  true,
					TTL:            300,
				}, nil).Times(3)
				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
					UseDefaultTTL:  false,
					TTL:            600,
				}, nil).Times(3)
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
					DNSZone:    "akamaized.net",
					RecordName: "test",
					Comments:   "change /ttl to 600",
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ttl",
							Value: "600",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net_with_ttl.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ttl", "600"),
					),
				},
			},
		},
		"different edge hostname, create": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.aka.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "aka.net.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
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
							ID:                "ehn_2",
							Domain:            "test3.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.aka.edgesuite.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "use_cases", testutils.LoadFixtureStringf(t, "%s/use_cases/use_cases_new.json", testDir)),
						resource.TestCheckOutput("edge_hostname", "test.aka.edgesuite.net"),
					),
				},
			},
		},
		"edge hostname exists": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "ehn_2",
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
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net.tf"),
					ExpectError: regexp.MustCompile("edgehostname 'test.akamaized.net' already exists"),
				},
			},
		},
		"edge hostname - update ip_behavior": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID:    "ctr_2",
					GroupID:       "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)

				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_456",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
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
				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 456,
				}, nil).Once()
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
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

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				// read
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_456",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "ehn_2",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_update_ip_behavior.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
					),
				},
			},
		},
		"edge hostname - update ttl": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil)
				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 456,
					UseDefaultTTL:  true,
					TTL:            300,
				}, nil).Times(3)
				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 456,
					UseDefaultTTL:  false,
					TTL:            600,
				}, nil).Times(3)
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
					DNSZone:    "akamaized.net",
					RecordName: "test",
					Comments:   "change /ttl to 600",
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ttl",
							Value: "600",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				// 2nd step
				// update
				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 456,
					UseDefaultTTL:  false,
					TTL:            600,
				}, nil).Times(3)
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
					DNSZone:           "akamaized.net",
					RecordName:        "test",
					Comments:          "change /ttl to 800",
					StatusUpdateEmail: []string{"hello@akamai.com"},
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ttl",
							Value: "800",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				mh.On("GetEdgeHostname", testutils.MockContext, 456).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 456,
					UseDefaultTTL:  false,
					TTL:            800,
				}, nil).Times(4)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net_with_ttl.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ttl", "600"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_update_ttl.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ttl", "800"),
					),
				},
			},
		},
		"edge hostname - update ip_behavior and ttl": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)

				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
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
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
					UseDefaultTTL:  true,
					TTL:            300,
				}, nil).Twice()
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
					DNSZone:           "akamaized.net",
					RecordName:        "test",
					Comments:          "change /ipVersionBehavior to IPV6_IPV4_DUALSTACK; change /ttl to 600",
					StatusUpdateEmail: []string{"hello@akamai.com"},
					Body: []hapi.UpdateEdgeHostnameRequestBody{
						{
							Op:    "replace",
							Path:  "/ipVersionBehavior",
							Value: "IPV6_IPV4_DUALSTACK",
						},
						{
							Op:    "replace",
							Path:  "/ttl",
							Value: "600",
						},
					},
				}).Return(&hapi.UpdateEdgeHostnameResponse{
					ChangeID: 123,
				}, nil).Once()

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				// read
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Twice()

				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
					UseDefaultTTL:  false,
					TTL:            600,
				}, nil).Times(3)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
						resource.TestCheckNoResourceAttr("akamai_edge_hostname.edgehostname", "ttl"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_update_ip_behavior_and_ttl.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ttl", "600"),
					),
				},
			},
		},
		"edge hostname - ip_behavior drift": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)

				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Twice()

				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE", // drift happens here
						},
						{
							ID:                "ehn_2",
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
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
				}, nil).Once()
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
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

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "PENDING",
				}, nil).Once()
				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil).Once()

				// read
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4_with_email.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
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
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)

				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Times(3)

				// 2nd step - update
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID:    123,
					IPVersionBehavior: "IPV4",
				}, nil).Once()
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_error_update_ipv6_performance.tf"),
					ExpectError: regexp.MustCompile("invalid IP version behavior: valid values are IPV4 and IPV6_IPV4_DUALSTACK; IPV6 and other values aren't currently supported"),
				},
			},
		},
		"edge hostname - update right after create - should wait for HAPI propagation": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)

				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil).Twice()

				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE", // drift happens here
						},
						{
							ID:                "ehn_2",
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
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(nil, &hapi.Error{
					Status: 404,
					Title:  "Invalid Record Name/DNS Zone",
				}).Times(15)
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
				}, nil).Once()
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
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

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil).Once()

				// read
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4_with_email.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
			},
		},
		"error - update retry reached retriesMax": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1. call from create method and refresh 2. update ip_behvior to improper value
				// 1st step - create
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)

				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Times(3)

				// 2nd step - update
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(nil, &hapi.Error{
					Status: 404,
					Title:  "Invalid Record Name/DNS Zone",
				}).Times(15)
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(nil, &hapi.Error{
					Status: 404,
					Title:  "reached max number of retries: 15",
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_error_update_ipv6_performance.tf"),
					ExpectError: regexp.MustCompile("reached max number of retries: 15"),
				},
			},
		},
		"error fetching edge hostnames": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"certificate required for ENHANCED_TLS": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "missing_certificate.tf"),
					ExpectError: regexp.MustCompile("a certificate enrollment ID is required for Enhanced TLS edge hostnames with 'edgekey.net' suffix"),
				},
			},
		},
		"error creating edge hostname": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
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
							ID:           "ehn_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error edge hostname not found": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
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
							ID:           "ehn_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Twice()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net.tf"),
					ExpectError: regexp.MustCompile("unable to find edge hostname"),
				},
			},
		},
		"error on empty product id for creation": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net_without_product_id.tf"),
					ExpectError: regexp.MustCompile("`product_id` must be specified for creation"),
				},
			},
		},
		"error on updating product id": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net.tf"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_net_different_product_id.tf"),
					ExpectError: regexp.MustCompile(`Changes to non-updatable fields 'product_id' and 'certificate' are not permitted`),
				},
			},
		},
		"error on updating certificate": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
						},
						{
							ID:                "ehn_2",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_PERFORMANCE",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test3.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test3",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_456",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_edgekey_net.tf"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_edgekey_net_different_certificate.tf"),
					ExpectError: regexp.MustCompile(`Changes to non-updatable fields 'product_id' and 'certificate' are not permitted`),
				},
			},
		},
		"update edge hostname without status email": {
			init: func(mp *papi.Mock, mh *hapi.Mock) {
				// 1st step
				// 1. call from create method
				// 2. and 3. call from read method
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()

				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
					EdgeHostnameID: "ehn_123",
				}, nil)
				// refresh
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
						{
							ID:                "ehn_2",
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
				mh.On("GetEdgeHostname", testutils.MockContext, 123).Return(&hapi.GetEdgeHostnameResponse{
					EdgeHostnameID: 123,
				}, nil).Once()
				mh.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
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

				mh.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
					Status: "SUCCEEDED",
				}, nil)

				// read
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:                "ehn_2",
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
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "new_akamaized_ipv4.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV4"),
					),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "update_no_status_update_email.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_123"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test.akamaized.net"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckNoResourceAttr("akamai_edge_hostname.edgehostname", "status_update_email"),
					),
				},
			},
		},
		"invalid edge hostname domain prefix for the akamaized.net domain suffix - The character '#' isn't allowed in the domain prefix.": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "invalid_char_#_in_edgehostname_domainprefix.tf"),
					ExpectError: regexp.MustCompile("A prefix for the edge hostname with the \"akamaized.net\" suffix must begin with a letter, end with a letter or digit, and contain only letters, digits, and hyphens, for example, abc-def, or abc-123"),
				},
			},
		},
		"invalid edge hostname domain prefix for the `akamaized.net` domain suffix. The domain prefix contains non-UTF-8 characters ('t中esãt').": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "incorrect_edgehostname_domainprefix_3.tf"),
					ExpectError: regexp.MustCompile("A prefix for the edge hostname with the \"akamaized.net\" suffix must begin with a letter, end with a letter or digit, and contain only letters, digits, and hyphens, for example, abc-def, or abc-123"),
				},
			},
		},
		"invalid edge hostname domain prefix for `akamaized.net`. The domain prefix can't end with a hyphen": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "edgehostname_domainprefix_ends_with_invalid_char_-.tf"),
					ExpectError: regexp.MustCompile("A prefix for the edge hostname with the \"akamaized.net\" suffix must begin with a letter, end with a letter or digit, and contain only letters, digits, and hyphens, for example, abc-def, or abc-123"),
				},
			},
		},
		"invalid edge hostname domain prefix for `edgesuite.net`. The domain prefix can't end with two consecutive dots": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "edgehostname_domainprefix_ends_with_consecutive_dots.tf"),
					ExpectError: regexp.MustCompile("A prefix for the edge hostname with the \"edgesuite.net\" suffix must begin with a letter, end with a letter, digit, or dot, and contain only letters, digits, dots, and hyphens, for example, abc-def.123.456., or abc.123-def"),
				},
			},
		},
		"invalid edge hostname domain prefix. The domain prefix exceeds the maximum allowed length of 63 characters": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "edgehostname_domainprefix_more_than_allowed_length.tf"),
					ExpectError: regexp.MustCompile(regexp.QuoteMeta(`The edge hostname prefix must be at least 1 character(s) and no more than 63 characters for "edgesuite.net" suffix; you provided 64 character(s).`)),
				},
			},
		},
		"invalid edge hostname domain prefix. The domain prefix less the minimum required length of 4 characters for akamaized.net": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "edgehostname_domainprefix_for_akamaized_dot_net_less_than_minimum_required_length.tf"),
					ExpectError: regexp.MustCompile(regexp.QuoteMeta(`The edge hostname prefix must be at least 4 character(s) and no more than 63 characters for "akamaized.net" suffix; you provided 2 character(s).`)),
				},
			},
		},
		"invalid edge hostname domain prefix. The domain prefix less the minimum required length of 1 characters": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/%s", testDir, "edgehostname_domainprefix_less_than_minimum_required_length.tf"),
					ExpectError: regexp.MustCompile(regexp.QuoteMeta(`The edge hostname prefix must be at least 1 character(s) and no more than 63 characters for "edgesuite.net" suffix; you provided 0 character(s).`)),
				},
			},
		},
		"valid edge hostname with hyphen in domain prefix name for akamaized.net, create edge hostname": {
			init: func(mp *papi.Mock, _ *hapi.Mock) {
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_123",
							Domain:            "test2.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV4",
						},
					}},
				}, nil).Once()
				mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test-prefix1",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "ehn_456",
				}, nil)
				mp.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "ehn_456",
							Domain:            "test-prefix1.akamaized.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test-prefix1",
							DomainSuffix:      "akamaized.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
					}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/%s", testDir, "valid_domain_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "id", "ehn_456"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "ip_behavior", "IPV6_COMPLIANCE"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "contract_id", "ctr_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "group_id", "grp_2"),
						resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", "edge_hostname", "test-prefix1.akamaized.net"),
					),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			clientHapi := &hapi.Mock{}
			if test.init != nil {
				test.init(client, clientHapi)
			}
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
	expectGetEdgeHostname := func(m *papi.Mock, edgehostID, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostname", testutils.MockContext, papi.GetEdgeHostnameRequest{
			EdgeHostnameID: edgehostID,
			ContractID:     contractID,
			GroupID:        groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameGetItem{
				ID:                "ehn_1",
				Domain:            "test.edgekey.net",
				ProductID:         "prd_2",
				DomainPrefix:      "test",
				DomainSuffix:      "edgekey.net",
				IPVersionBehavior: "IPV4",
				Secure:            true,
			},
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test2.edgekey.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test2",
					DomainSuffix:      "edgekey.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	expectGetEdgeHostnameAkamaized := func(m *papi.Mock, edgehostID, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostname", testutils.MockContext, papi.GetEdgeHostnameRequest{
			EdgeHostnameID: edgehostID,
			ContractID:     contractID,
			GroupID:        groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameGetItem{
				ID:                "ehn_1",
				Domain:            "test.akamaized.net",
				ProductID:         "prd_2",
				DomainPrefix:      "test",
				DomainSuffix:      "akamaized.net",
				IPVersionBehavior: "IPV4",
				Secure:            true,
			},
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test2.edgekey.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test2",
					DomainSuffix:      "edgekey.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	expectGetEdgeHostnameAkamaizedNoProductID := func(m *papi.Mock, edgehostID, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostname", testutils.MockContext, papi.GetEdgeHostnameRequest{
			EdgeHostnameID: edgehostID,
			ContractID:     contractID,
			GroupID:        groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameGetItem{
				ID:                "ehn_1",
				Domain:            "test.akamaized.net",
				DomainPrefix:      "test",
				DomainSuffix:      "akamaized.net",
				IPVersionBehavior: "IPV4",
				Secure:            true,
			},
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test2.edgekey.net",
					DomainPrefix:      "test2",
					DomainSuffix:      "edgekey.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	expectGetEdgeHostnames := func(m *papi.Mock, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
			ContractID: contractID,
			GroupID:    groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test1.edgekey.net",
					DomainPrefix:      "test1",
					DomainSuffix:      "edgekey.net",
					IPVersionBehavior: "IPV4",
					Secure:            true,
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	expectGetEdgeHostnamesAkamaized := func(m *papi.Mock, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
			ContractID: contractID,
			GroupID:    groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test1.akamaized.net",
					DomainPrefix:      "test1",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
					Secure:            true,
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	createEdgeHostnames := func(mp *papi.Mock) *mock.Call {
		return mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameCreate{
				ProductID:         "prd_2",
				DomainPrefix:      "test",
				DomainSuffix:      "edgekey.net",
				SecureNetwork:     "ENHANCED_TLS",
				IPVersionBehavior: "IPV4",
				CertEnrollmentID:  123456,
				SlotNumber:        123456,
			},
		}).Return(&papi.CreateEdgeHostnameResponse{
			EdgeHostnameID: "ehn_1",
		}, nil)
	}
	createEdgeHostnamesAkamaized := func(mp *papi.Mock) *mock.Call {
		return mp.On("CreateEdgeHostname", testutils.MockContext, papi.CreateEdgeHostnameRequest{
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
			EdgeHostnameID: "ehn_1",
		}, nil)
	}

	expectGetEdgeHostnamesAfterCreate := func(m *papi.Mock, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
			ContractID: contractID,
			GroupID:    groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test.edgekey.net",
					DomainPrefix:      "test",
					DomainSuffix:      "edgekey.net",
					IPVersionBehavior: "IPV4",
					Secure:            true,
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}
	expectGetEdgeHostnamesAkamaizedAfterCreate := func(m *papi.Mock, contractID, groupID string) *mock.Call {
		return m.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
			ContractID: contractID,
			GroupID:    groupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "test.akamaized.net",
					DomainPrefix:      "test",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
					Secure:            true,
				},
				{
					ID:                "ehn_2",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}

	expectGetEdgeHostnameHAPIByID := func(m *hapi.Mock, edgeHostnameID int) *mock.Call {
		return m.On("GetEdgeHostname", testutils.MockContext, edgeHostnameID).Return(&hapi.GetEdgeHostnameResponse{
			EdgeHostnameID:         edgeHostnameID,
			RecordName:             "test",
			DNSZone:                "edgekey.net",
			SecurityType:           "ENHANCED-TLS",
			UseDefaultTTL:          true,
			UseDefaultMap:          true,
			IPVersionBehavior:      "IPV4",
			ProductID:              "DSA",
			TTL:                    21600,
			Map:                    "a;dscb.akamai.net",
			SlotNumber:             3250,
			Comments:               "Created by Property-Manager/PAPI on Tue Feb 21 14:03:36 UTC 2023",
			SerialNumber:           1520,
			CustomTarget:           "",
			ChinaCdn:               hapi.ChinaCDN{},
			IsEdgeIPBindingEnabled: false,
			MapAlias:               "",
			UseCases:               nil,
		}, nil)
	}
	expectGetEdgeHostnameHAPIByIDWithCustomTTL := func(m *hapi.Mock, edgeHostnameID int) *mock.Call {
		return m.On("GetEdgeHostname", testutils.MockContext, edgeHostnameID).Return(&hapi.GetEdgeHostnameResponse{
			EdgeHostnameID:         edgeHostnameID,
			RecordName:             "test",
			DNSZone:                "edgekey.net",
			SecurityType:           "ENHANCED-TLS",
			UseDefaultTTL:          false,
			UseDefaultMap:          true,
			IPVersionBehavior:      "IPV4",
			ProductID:              "DSA",
			TTL:                    56789,
			Map:                    "a;dscb.akamai.net",
			SlotNumber:             3250,
			Comments:               "Created by Property-Manager/PAPI on Tue Feb 21 14:03:36 UTC 2023",
			SerialNumber:           1520,
			CustomTarget:           "",
			ChinaCdn:               hapi.ChinaCDN{},
			IsEdgeIPBindingEnabled: false,
			MapAlias:               "",
			UseCases:               nil,
		}, nil)
	}
	expectGetEdgeHostnameHAPIByIDAkamaized := func(m *hapi.Mock, edgeHostnameID int) *mock.Call {
		return m.On("GetEdgeHostname", testutils.MockContext, edgeHostnameID).Return(&hapi.GetEdgeHostnameResponse{
			EdgeHostnameID:         edgeHostnameID,
			RecordName:             "test",
			DNSZone:                "akamaized.net",
			SecurityType:           "STANDARD-TLS",
			UseDefaultTTL:          true,
			UseDefaultMap:          true,
			IPVersionBehavior:      "IPV4",
			ProductID:              "DSA",
			TTL:                    21600,
			Map:                    "a;dscb.akamai.net",
			SlotNumber:             3250,
			Comments:               "Created by Property-Manager/PAPI on Tue Feb 21 14:03:36 UTC 2023",
			SerialNumber:           1520,
			CustomTarget:           "",
			ChinaCdn:               hapi.ChinaCDN{},
			IsEdgeIPBindingEnabled: false,
			MapAlias:               "",
			UseCases:               nil,
		}, nil)
	}
	expectGetCertificate := func(m *hapi.Mock) *mock.Call {
		return m.On("GetCertificate", testutils.MockContext, hapi.GetCertificateRequest{
			RecordName: "test",
			DNSZone:    "edgekey.net",
		}).Return(&hapi.GetCertificateResponse{
			AvailableDomains: []string{"DevExpAutomatedTest.2dUlc9.com"},
			CertificateID:    "123456",
			CertificateType:  "THIRD_PARTY",
			CommonName:       "DevExpAutomatedTest.2dUlc9.com",
			ExpirationDate:   test.NewTimeFromString(t, "2026-06-26T15:34:04.000+00:00"),
			SerialNumber:     "fa:ke:5a:e5:a8:c9:8f:0c:28:48:af:db:fa:78:cc:db",
			SlotNumber:       3250,
			Status:           "DEPLOYED",
			ValidationType:   "THIRD_PARTY",
		}, nil)
	}

	t.Run("import existing akamaized edgehostname without certificate - no product id provided by user", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2"

		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once()
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostnameAkamaized(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_akamaized.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_2", rs.Attributes["product_id"])
							assert.Equal(t, "", rs.Attributes["certificate"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname without domain validation", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,2,2"

		client.On("GetEdgeHostname", testutils.MockContext, papi.GetEdgeHostnameRequest{
			EdgeHostnameID: "ehn_1",
			ContractID:     "ctr_2",
			GroupID:        "grp_2",
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_2",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameGetItem{
				ID:                "ehn_1",
				Domain:            "t1.akamaized.net",
				DomainPrefix:      "t1",
				DomainSuffix:      "akamaized.net",
				IPVersionBehavior: "IPV4",
				Secure:            true,
			},
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_2",
					Domain:            "test2.edgekey.net",
					DomainPrefix:      "test2",
					DomainSuffix:      "edgekey.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "ehn_3",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
		clientHapi.On("GetEdgeHostname", testutils.MockContext, 1).Return(&hapi.GetEdgeHostnameResponse{
			EdgeHostnameID:         1,
			RecordName:             "t1",
			DNSZone:                "akamaized.net",
			SecurityType:           "STANDARD-TLS",
			UseDefaultTTL:          true,
			UseDefaultMap:          true,
			IPVersionBehavior:      "IPV4",
			ProductID:              "DSA",
			TTL:                    21600,
			Map:                    "a;dscb.akamai.net",
			SlotNumber:             3250,
			Comments:               "Created by Property-Manager/PAPI on Tue Feb 21 14:03:36 UTC 2023",
			SerialNumber:           1520,
			CustomTarget:           "",
			ChinaCdn:               hapi.ChinaCDN{},
			IsEdgeIPBindingEnabled: false,
			MapAlias:               "",
			UseCases:               nil,
		}, nil)
		client.On("GetEdgeHostnames", testutils.MockContext, papi.GetEdgeHostnamesRequest{
			ContractID: "ctr_2",
			GroupID:    "grp_2",
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_2",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "ehn_1",
					Domain:            "t1.akamaized.net",
					DomainPrefix:      "t1",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
					Secure:            true,
				},
				{
					ID:                "ehn_3",
					Domain:            "test3.edgesuite.net",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/edgehostname_domainprefix_for_akamaized_dot_net_less_than_minimum_required_length.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_2", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							// edgehostname bellow is not validated, because flow does not use create context
							assert.Equal(t, "t1.akamaized.net", rs.Attributes["edge_hostname"])
							assert.Equal(t, "", rs.Attributes["certificate"])
							return nil
						},
						ImportStateId: id,
						ResourceName:  "akamai_edge_hostname.edgehostname",
						// ImportStateVerify is set to false. Because of validation it can't use create context, which means it doesn't have old state.
						ImportStateVerify: false,
					},
				},
			})
		})
		client.AssertExpectations(t)
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname without certificate - product provided by user, different product id returned by api", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2,prd_10"

		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once()
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostnameAkamaized(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_akamaized.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_10", rs.Attributes["product_id"])
							assert.Equal(t, "", rs.Attributes["certificate"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname without certificate - product id provided by user, product id not returned by api", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2,prd_2"

		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once()
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostnameAkamaizedNoProductID(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_akamaized.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_2", rs.Attributes["product_id"])
							assert.Equal(t, "", rs.Attributes["certificate"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname without certificate - product id without prefix provided by user, product id not returned by api", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2,2"

		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once()
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostnameAkamaizedNoProductID(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_akamaized.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_2", rs.Attributes["product_id"])
							assert.Equal(t, "", rs.Attributes["certificate"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname without certificate - product id not provided by user, product id not returned by api", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2"

		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once()
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostnameAkamaizedNoProductID(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_akamaized.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "", rs.Attributes["product_id"])
							assert.Equal(t, "", rs.Attributes["certificate"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing edgehostname with certificate - no product id provided by user", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2"

		expectGetEdgeHostnames(client, "ctr_1", "grp_2").Once()
		createEdgeHostnames(client).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostname(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByID(clientHapi, 1).Once()
		expectGetCertificate(clientHapi).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
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
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_2", rs.Attributes["product_id"])
							assert.Equal(t, "123456", rs.Attributes["certificate"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing edgehostname with missing certificate - no product id provided by user", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2"

		expectGetEdgeHostnames(client, "ctr_1", "grp_2").Once()
		createEdgeHostnames(client).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2")
		expectGetEdgeHostname(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByID(clientHapi, 1).Once()

		resp := hapi.Error{
			Type:            "CERTIFICATE_NOT_FOUND",
			Title:           "Certificate Not Found",
			Status:          404,
			Detail:          "Details are not available for this certificate; the certificate is missing or access is denied",
			Instance:        "/hapi/error-instances/a30f67cc-df20-4e02-bbc3-cf7c204a4aab",
			RequestInstance: "http://origin.pulsar.akamai.com/hapi/open/v1/dns-zones/edgekey.net/edge-hostnames/example.com/certificate?depth=ALL&accountSwitchKey=F-AC-1937217#d7aa7348",
			Method:          "GET",
			RequestTime:     "2022-11-30T18:51:43.482982Z",
		}
		err := fmt.Errorf("%s: %s: %w", hapi.ErrGetCertificate, hapi.ErrNotFound, &resp)
		clientHapi.On("GetCertificate", testutils.MockContext, hapi.GetCertificateRequest{
			RecordName: "test",
			DNSZone:    "edgekey.net",
		}).Return(nil, err)

		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_missing_certificate.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_missing_certificate.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_2", rs.Attributes["product_id"])
							assert.Equal(t, "", rs.Attributes["certificate"])
							return nil
						},
						ImportStateId:           id,
						ResourceName:            "akamai_edge_hostname.importedgehostname",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"product_id", "certificate"},
						// TODO Remove "certificate" here after implementing wait for EdgeHostname propagation to HAPI in DXE-3821,
					},
				},
			})
		})
		client.AssertExpectations(t)
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing edgehostname with custom ttl - no product id provided by user", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}
		id := "ehn_1,1,2"
		EgdeHostnameCreatePollInterval = time.Microsecond

		expectGetEdgeHostnames(client, "ctr_1", "grp_2").Once()
		createEdgeHostnames(client).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2")

		clientHapi.On("GetEdgeHostname", testutils.MockContext, 1).Return(&hapi.GetEdgeHostnameResponse{
			EdgeHostnameID: 1,
			UseDefaultTTL:  true,
			TTL:            300,
		}, nil).Times(3)
		clientHapi.On("GetEdgeHostname", testutils.MockContext, 1).Return(&hapi.GetEdgeHostnameResponse{
			EdgeHostnameID: 1,
			UseDefaultTTL:  false,
			TTL:            56789,
		}, nil).Times(3)
		clientHapi.On("UpdateEdgeHostname", testutils.MockContext, hapi.UpdateEdgeHostnameRequest{
			DNSZone:    "edgekey.net",
			RecordName: "test",
			Comments:   "change /ttl to 56789",
			Body: []hapi.UpdateEdgeHostnameRequestBody{
				{
					Op:    "replace",
					Path:  "/ttl",
					Value: "56789",
				},
			},
		}).Return(&hapi.UpdateEdgeHostnameResponse{
			ChangeID: 123,
		}, nil).Once()

		clientHapi.On("GetChangeRequest", testutils.MockContext, hapi.GetChangeRequest{ChangeID: 123}).Return(&hapi.ChangeRequest{
			Status: "SUCCEEDED",
		}, nil)

		expectGetEdgeHostname(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDWithCustomTTL(clientHapi, 1).Times(3)
		expectGetCertificate(clientHapi).Once()
		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/creation_before_import_edgehostname_with_ttl.tf"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_with_ttl.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ehn_1", rs.Attributes["id"])
							assert.Equal(t, "prd_2", rs.Attributes["product_id"])
							assert.Equal(t, "123456", rs.Attributes["certificate"])
							assert.Equal(t, "56789", rs.Attributes["ttl"])
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
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname without certificate - product_id supplied", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}

		// create
		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once() // check if exists: not found
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2").Twice()

		// import
		expectGetEdgeHostnameAkamaized(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2").Once()

		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized_product_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateId:     "ehn_1,1,2,prd_2",
						ResourceName:      "akamai_edge_hostname.importedgehostname",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing akamaized edgehostname - product_id without prefix supplied ", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}

		// create
		expectGetEdgeHostnamesAkamaized(client, "ctr_1", "grp_2").Once() // check if exists: not found
		createEdgeHostnamesAkamaized(client).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2").Twice()

		// import
		expectGetEdgeHostnameAkamaized(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByIDAkamaized(clientHapi, 1).Once()
		expectGetEdgeHostnamesAkamaizedAfterCreate(client, "ctr_1", "grp_2").Once()

		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized_product_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateId:     "ehn_1,1,2,2",
						ResourceName:      "akamai_edge_hostname.importedgehostname",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing edgehostname with certificate", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}

		// create
		expectGetEdgeHostnames(client, "ctr_1", "grp_2").Once() // check if exists: not found
		createEdgeHostnames(client).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2").Twice()

		// import
		expectGetEdgeHostname(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByID(clientHapi, 1).Once()
		expectGetCertificate(clientHapi).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2").Once()

		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_product_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateId:     "ehn_1,1,2,prd_2",
						ResourceName:      "akamai_edge_hostname.importedgehostname",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
		clientHapi.AssertExpectations(t)
	})
	t.Run("import existing edgehostname with missing certificate", func(t *testing.T) {
		client := &papi.Mock{}
		clientHapi := &hapi.Mock{}

		// create
		expectGetEdgeHostnames(client, "ctr_1", "grp_2").Once() // check if exists: not found
		createEdgeHostnames(client).Once()
		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2").Twice()

		// import
		expectGetEdgeHostname(client, "ehn_1", "ctr_1", "grp_2").Once()
		expectGetEdgeHostnameHAPIByID(clientHapi, 1).Once()

		resp := hapi.Error{
			Type:            "CERTIFICATE_NOT_FOUND",
			Title:           "Certificate Not Found",
			Status:          404,
			Detail:          "Details are not available for this certificate; the certificate is missing or access is denied",
			Instance:        "/hapi/error-instances/a30f67cc-df20-4e02-bbc3-cf7c204a4aab",
			RequestInstance: "http://origin.pulsar.akamai.com/hapi/open/v1/dns-zones/edgekey.net/edge-hostnames/example.com/certificate?depth=ALL&accountSwitchKey=F-AC-1937217#d7aa7348",
			Method:          "GET",
			RequestTime:     "2022-11-30T18:51:43.482982Z",
		}
		err := fmt.Errorf("%s: %s: %w", hapi.ErrGetCertificate, hapi.ErrNotFound, &resp)
		clientHapi.On("GetCertificate", testutils.MockContext, hapi.GetCertificateRequest{
			RecordName: "test",
			DNSZone:    "edgekey.net",
		}).Return(nil, err)

		expectGetEdgeHostnamesAfterCreate(client, "ctr_1", "grp_2").Once()

		useClient(client, clientHapi, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_product_id.tf"),
					},
					{
						ImportState:             true,
						ImportStateId:           "ehn_1,1,2,prd_2",
						ResourceName:            "akamai_edge_hostname.importedgehostname",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"certificate"},
						// TODO Remove "certificate" here after implementing wait for EdgeHostname propagation to HAPI in DXE-3821,
					},
				},
			})
		})
		client.AssertExpectations(t)
		clientHapi.AssertExpectations(t)
	})
	t.Run("import error - too few parts of id", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized_product_id.tf"),
					ImportState:   true,
					ImportStateId: "ehn_1,1",
					ResourceName:  "akamai_edge_hostname.importedgehostname",
					ExpectError: regexp.MustCompile("expected import identifier with format: " +
						`"EdgehostNameID,contractID,groupID\[,productID]". Got: "ehn_1,1"`),
				},
			},
		})
	})
	t.Run("import error - too many parts of id", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized_product_id.tf"),
					ImportState:   true,
					ImportStateId: "ehn_1,1,2,prd_2,foo",
					ResourceName:  "akamai_edge_hostname.importedgehostname",
					ExpectError: regexp.MustCompile("expected import identifier with format: " +
						`"EdgehostNameID,contractID,groupID\[,productID]". Got: "ehn_1,1,2,prd_2,foo"`),
				},
			},
		})
	})
	t.Run("import error - empty product id", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceEdgeHostname/import_edgehostname_akamaized_product_id.tf"),
					ImportState:   true,
					ImportStateId: "ehn_1,1,2,",
					ResourceName:  "akamai_edge_hostname.importedgehostname",
					ExpectError:   regexp.MustCompile(`productID is empty for the import ID="ehn_1,1,2,"`),
				},
			},
		})
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
					ID:           "ehn_2",
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
					ID:           "ehn_2",
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
					ID:           "ehn_2",
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
					ID:           "ehn_2",
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
					ID:           "ehn_2",
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
			oldFixture := testutils.LoadFixtureStringf(t, "%s/%s", testDir, test.oldPath)
			newFixture := testutils.LoadFixtureStringf(t, "%s/%s", testDir, test.newPath)

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
