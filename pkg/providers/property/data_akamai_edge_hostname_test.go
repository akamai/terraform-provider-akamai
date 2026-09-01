package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestNewEdgeHostnameDataSource(t *testing.T) {
	t.Parallel()

	testDir := "testdata/TestDataEdgeHostname"

	defaultEdgeHostname := papi.EdgeHostnameGetItem{
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
		HTTPSServiceBinding: ptr.To("H3"),
	}

	checker := test.NewStateChecker("data.akamai_edge_hostname.test").
		CheckEqual("id", "ehn_123").
		CheckEqual("contract_id", "ctr_123").
		CheckEqual("group_id", "grp_123").
		CheckEqual("account_id", "act_123").
		CheckEqual("edge_hostname.domain_prefix", "test1").
		CheckEqual("edge_hostname.domain_suffix", "edgekey.net").
		CheckEqual("edge_hostname.edge_hostname_domain", "test1.edgekey.net").
		CheckEqual("edge_hostname.edge_hostname_id", "ehn_123").
		CheckEqual("edge_hostname.ip_version_behaviour", "IPV4").
		CheckEqual("edge_hostname.product_id", "prd_123").
		CheckEqual("edge_hostname.secure", "true").
		CheckEqual("edge_hostname.status", "active").
		CheckEqual("edge_hostname.use_cases.#", "1").
		CheckEqual("edge_hostname.use_cases.0.option", "LIVE").
		CheckEqual("edge_hostname.use_cases.0.type", "GLOBAL").
		CheckEqual("edge_hostname.use_cases.0.use_case", "Segmented_Media_Mode").
		CheckEqual("edge_hostname.https_service_binding", "H3")

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"happy path - get edge hostname": {
			init: func(m *papi.Mock) {
				mockGetEdgeHostname(m, papi.GetEdgeHostnameRequest{
					EdgeHostnameID: "ehn_123",
					ContractID:     "ctr_123",
					GroupID:        "grp_123",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:  "act_123",
					ContractID: "ctr_123",
					GroupID:    "grp_123",
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{defaultEdgeHostname},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostname.tf", testDir),
					Check:  checker.Build(),
				},
			},
		},
		"happy path - get edge hostname with options": {
			init: func(m *papi.Mock) {
				mockGetEdgeHostname(m, papi.GetEdgeHostnameRequest{
					EdgeHostnameID: "ehn_123",
					ContractID:     "ctr_123",
					GroupID:        "grp_123",
					Options:        []string{"mapDetails"},
				}, papi.GetEdgeHostnamesResponse{
					AccountID:  "act_123",
					ContractID: "ctr_123",
					GroupID:    "grp_123",
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{defaultEdgeHostname},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostname_with_options.tf", testDir),
					Check: checker.
						CheckEqual("options.#", "1").
						CheckEqual("options.0", "mapDetails").
						Build(),
				},
			},
		},
		"happy path - edge hostname without product_id": {
			init: func(m *papi.Mock) {
				noProductID := defaultEdgeHostname
				noProductID.ProductID = ""
				mockGetEdgeHostname(m, papi.GetEdgeHostnameRequest{
					EdgeHostnameID: "ehn_123",
					ContractID:     "ctr_123",
					GroupID:        "grp_123",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:  "act_123",
					ContractID: "ctr_123",
					GroupID:    "grp_123",
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{noProductID},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostname.tf", testDir),
					Check: checker.
						CheckMissing("edge_hostname.product_id").
						Build(),
				},
			},
		},
		"happy path - edge hostname without use_cases and https_service_binding": {
			init: func(m *papi.Mock) {
				minimal := defaultEdgeHostname
				minimal.UseCases = nil
				minimal.HTTPSServiceBinding = nil
				mockGetEdgeHostname(m, papi.GetEdgeHostnameRequest{
					EdgeHostnameID: "ehn_123",
					ContractID:     "ctr_123",
					GroupID:        "grp_123",
				}, papi.GetEdgeHostnamesResponse{
					AccountID:  "act_123",
					ContractID: "ctr_123",
					GroupID:    "grp_123",
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{minimal},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgehostname.tf", testDir),
					Check: checker.
						CheckEqual("edge_hostname.use_cases.#", "0").
						CheckMissing("edge_hostname.use_cases.0.option").
						CheckMissing("edge_hostname.use_cases.0.type").
						CheckMissing("edge_hostname.use_cases.0.use_case").
						CheckMissing("edge_hostname.https_service_binding").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *papi.Mock) {
				mockGetEdgeHostname(m, papi.GetEdgeHostnameRequest{
					EdgeHostnameID: "ehn_123",
					ContractID:     "ctr_123",
					GroupID:        "grp_123",
				}, papi.GetEdgeHostnamesResponse{}, &papi.Error{
					StatusCode: 500,
					Detail:     "oops",
				})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgehostname.tf", testDir),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - missing id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/missing_id.tf", testDir),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
				},
			},
		},
		"error - missing contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/missing_contract_id.tf", testDir),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
			},
		},
		"error - missing group_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/missing_group_id.tf", testDir),
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

func mockGetEdgeHostname(m *papi.Mock, req papi.GetEdgeHostnameRequest, resp papi.GetEdgeHostnamesResponse,
	err error) *mock.Call {
	if err != nil {
		return m.On("GetEdgeHostname", testutils.MockContext, req).Return((*papi.GetEdgeHostnamesResponse)(nil), err).Once()
	}

	return m.On("GetEdgeHostname", testutils.MockContext, req).Return(&resp, nil).Times(3)
}
