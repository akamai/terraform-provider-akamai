package property

import (
	"regexp"
	"slices"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v12/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestNewEdgeHostnamesDataSource(t *testing.T) {
	t.Parallel()

	testDir := "testdata/TestDataEdgeHostnames"

	defaultEdgeHostnames := papi.EdgeHostnameItems{
		Items: []papi.EdgeHostnameGetItem{
			{
				ID:                "ehn_123",
				Domain:            "test1.edgekey.net",
				ProductID:         "prd_123",
				DomainPrefix:      "test1",
				DomainSuffix:      "edgekey.net",
				Status:            "active",
				Secure:            true,
				IPVersionBehavior: "IPV4",
				UseCases: []papi.UseCase{
					{
						Option:  "LIVE",
						Type:    "GLOBAL",
						UseCase: "Segmented_Media_Mode",
					},
				},
			},
			{
				ID:                "ehn_234",
				Domain:            "test2.edgekey.net",
				DomainPrefix:      "test1.test2",
				DomainSuffix:      "edgesuite.net",
				Status:            "inactive",
				Secure:            false,
				IPVersionBehavior: "IPV6_COMPLIANCE",
			},
			{
				ID:                  "ehn_345",
				Domain:              "test3.edgekey.net",
				ProductID:           "prd_345",
				DomainPrefix:        "test2.test3",
				DomainSuffix:        "edgesuite.net",
				Status:              "active",
				Secure:              false,
				IPVersionBehavior:   "IPV4",
				HTTPSServiceBinding: ptr.To("H3"),
			},
		},
	}

	checker := test.NewStateChecker("data.akamai_edge_hostnames.test").
		CheckMissing("options").
		CheckEqual("contract_id", "ctr_123").
		CheckEqual("group_id", "grp_123").
		CheckEqual("account_id", "act_123").
		CheckEqual("edge_hostnames.#", "3").
		CheckEqual("edge_hostnames.0.domain_prefix", "test1").
		CheckEqual("edge_hostnames.0.domain_suffix", "edgekey.net").
		CheckEqual("edge_hostnames.0.edge_hostname_domain", "test1.edgekey.net").
		CheckEqual("edge_hostnames.0.edge_hostname_id", "ehn_123").
		CheckEqual("edge_hostnames.0.ip_version_behaviour", "IPV4").
		CheckEqual("edge_hostnames.0.product_id", "prd_123").
		CheckEqual("edge_hostnames.0.secure", "true").
		CheckEqual("edge_hostnames.0.status", "active").
		CheckEqual("edge_hostnames.0.use_cases.#", "1").
		CheckEqual("edge_hostnames.0.use_cases.0.option", "LIVE").
		CheckEqual("edge_hostnames.0.use_cases.0.type", "GLOBAL").
		CheckEqual("edge_hostnames.0.use_cases.0.use_case", "Segmented_Media_Mode").
		CheckEqual("edge_hostnames.1.domain_prefix", "test1.test2").
		CheckEqual("edge_hostnames.1.domain_suffix", "edgesuite.net").
		CheckEqual("edge_hostnames.1.edge_hostname_domain", "test2.edgekey.net").
		CheckEqual("edge_hostnames.1.edge_hostname_id", "ehn_234").
		CheckEqual("edge_hostnames.1.ip_version_behaviour", "IPV6_COMPLIANCE").
		CheckMissing("edge_hostnames.1.product_id").
		CheckEqual("edge_hostnames.1.secure", "false").
		CheckEqual("edge_hostnames.1.status", "inactive").
		CheckEqual("edge_hostnames.1.use_cases.#", "0").
		CheckMissing("edge_hostnames.1.use_cases").
		CheckMissing("edge_hostnames.1.use_cases.0.option").
		CheckMissing("edge_hostnames.1.use_cases.0.type").
		CheckMissing("edge_hostnames.1.use_cases.0.use_case").
		CheckEqual("edge_hostnames.2.domain_prefix", "test2.test3").
		CheckEqual("edge_hostnames.2.domain_suffix", "edgesuite.net").
		CheckEqual("edge_hostnames.2.edge_hostname_domain", "test3.edgekey.net").
		CheckEqual("edge_hostnames.2.edge_hostname_id", "ehn_345").
		CheckEqual("edge_hostnames.2.ip_version_behaviour", "IPV4").
		CheckEqual("edge_hostnames.2.product_id", "prd_345").
		CheckEqual("edge_hostnames.2.secure", "false").
		CheckEqual("edge_hostnames.2.status", "active").
		CheckEqual("edge_hostnames.2.use_cases.#", "0").
		CheckMissing("edge_hostnames.2.use_cases").
		CheckMissing("edge_hostnames.2.use_cases.0.option").
		CheckMissing("edge_hostnames.2.use_cases.0.type").
		CheckMissing("edge_hostnames.2.use_cases.0.use_case").
		CheckEqual("edge_hostnames.2.https_service_binding", "H3")

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"happy path - list edge hostnames": {
			init: func(m *papi.Mock) {
				mockGetEdgeHostnames(m, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_123",
					GroupID:    "grp_123",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:     "act_123",
					ContractID:    "ctr_123",
					GroupID:       "grp_123",
					EdgeHostnames: defaultEdgeHostnames,
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostnames.tf", testDir),
					Check:  checker.Build(),
				},
			},
		},
		"happy path - list edge hostnames with delivery policy": {
			init: func(m *papi.Mock) {
				withDeliveryPolicy := defaultEdgeHostnames
				withDeliveryPolicy.Items = slices.Clone(defaultEdgeHostnames.Items)
				for i := range withDeliveryPolicy.Items {
					if i%2 == 0 {
						withDeliveryPolicy.Items[i].DeliveryPolicy = &papi.DeliveryPolicy{
							PolicyType: "GEO_OPT_IN",
							PolicyAttributes: []papi.PolicyAttribute{
								{
									AttributeType: "GEO",
									Code:          "US",
									Name:          "USA",
								},
							},
						}
					}
				}
				mockGetEdgeHostnames(m, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_123",
					GroupID:    "grp_123",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:     "act_123",
					ContractID:    "ctr_123",
					GroupID:       "grp_123",
					EdgeHostnames: withDeliveryPolicy,
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostnames.tf", testDir),
					Check: checker.
						CheckEqual("edge_hostnames.#", "3").
						CheckEqual("edge_hostnames.0.delivery_policy.policy_type", "GEO_OPT_IN").
						CheckEqual("edge_hostnames.0.delivery_policy.policy_attributes.0.attribute_type", "GEO").
						CheckEqual("edge_hostnames.0.delivery_policy.policy_attributes.0.code", "US").
						CheckEqual("edge_hostnames.0.delivery_policy.policy_attributes.0.name", "USA").
						CheckEqual("edge_hostnames.2.delivery_policy.policy_type", "GEO_OPT_IN").
						CheckEqual("edge_hostnames.2.delivery_policy.policy_attributes.0.attribute_type", "GEO").
						CheckEqual("edge_hostnames.2.delivery_policy.policy_attributes.0.code", "US").
						CheckEqual("edge_hostnames.2.delivery_policy.policy_attributes.0.name", "USA").
						CheckMissing("edge_hostnames.1.delivery_policy.policy_type").
						CheckMissing("edge_hostnames.1.delivery_policy.policy_attributes.0.attribute_type").
						CheckMissing("edge_hostnames.1.delivery_policy.policy_attributes.0.code").
						CheckMissing("edge_hostnames.1.delivery_policy.policy_attributes.0.name").
						Build(),
				},
			},
		},
		"happy path - list edge hostnames with options": {
			init: func(m *papi.Mock) {
				withMapDeatails := defaultEdgeHostnames
				withMapDeatails.Items = slices.Clone(defaultEdgeHostnames.Items)
				for i := range withMapDeatails.Items {
					if i%2 == 0 {
						withMapDeatails.Items[i].MapDetailsSlotNumber = ptr.To(int64(4321 - i))
					} else {
						withMapDeatails.Items[i].MapDetailsMapDomain = ptr.To("example.com.edgekey.net")
						withMapDeatails.Items[i].MapDetailsSerialNumber = ptr.To(int64(1234 + i))
					}
				}
				mockGetEdgeHostnames(m, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_123",
					GroupID:    "grp_123",
					Options:    "mapDetails",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:     "act_123",
					ContractID:    "ctr_123",
					GroupID:       "grp_123",
					EdgeHostnames: withMapDeatails,
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostnames_with_options.tf", testDir),
					Check: checker.
						CheckEqual("options", "mapDetails").
						CheckEqual("edge_hostnames.#", "3").
						CheckEqual("edge_hostnames.0.map_details_slot_number", "4321").
						CheckMissing("edge_hostnames.0.map_details_serial_number").
						CheckMissing("edge_hostnames.0.map_details_map_domain").
						CheckEqual("edge_hostnames.1.map_details_serial_number", "1235").
						CheckMissing("edge_hostnames.1.map_details_slot_number").
						CheckEqual("edge_hostnames.1.map_details_map_domain", "example.com.edgekey.net").
						CheckEqual("edge_hostnames.2.map_details_slot_number", "4319").
						CheckMissing("edge_hostnames.2.map_details_serial_number").
						CheckMissing("edge_hostnames.2.map_details_map_domain").
						Build(),
				},
			},
		},
		"happy path - API returns empty list": {
			init: func(m *papi.Mock) {
				mockGetEdgeHostnames(m, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_123",
					GroupID:    "grp_123",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:  "act_123",
					ContractID: "ctr_123",
					GroupID:    "grp_123",
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostnames.tf", testDir),
					Check: test.NewStateChecker("data.akamai_edge_hostnames.test").
						CheckEqual("contract_id", "ctr_123").
						CheckEqual("group_id", "grp_123").
						CheckEqual("account_id", "act_123").
						CheckEqual("edge_hostnames.#", "0").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *papi.Mock) {
				mockGetEdgeHostnames(m, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_123",
					GroupID:    "grp_123",
				}, papi.GetEdgeHostnamesResponse{}, &papi.Error{
					StatusCode: 500,
					Detail:     "oops",
				})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgehostnames.tf", testDir),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - missing contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgehostnames_missing_contract_id.tf", testDir),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
			},
		},
		"error - missing group_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgehostnames_missing_group_id.tf", testDir),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.PAPI)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps:                    tc.steps,
			})

			client.PAPI.AssertExpectations(t)
		})
	}
}

func mockGetEdgeHostnames(m *papi.Mock, req papi.GetEdgeHostnamesRequest, resp papi.GetEdgeHostnamesResponse,
	err error) *mock.Call {
	if err != nil {
		return m.On("GetEdgeHostnames", testutils.MockContext, req).Return((*papi.GetEdgeHostnamesResponse)(nil), err).Once()
	}

	return m.On("GetEdgeHostnames", testutils.MockContext, req).Return(&resp, nil).Times(3)
}
