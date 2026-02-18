package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type accountHostnamesTestData struct {
	req       papi.ListActiveAccountHostnamesRequest
	hostnames papi.ActiveAccountHostnames
}

func TestDataPropertyAccountHostnames(t *testing.T) {
	t.Parallel()

	defaultHostnames := papi.ActiveAccountHostnames{
		Items: []papi.ActiveAccountHostnameItem{
			{
				CnameFrom:                "api.example.com",
				ContractID:               "ctr_1-AB123",
				GroupID:                  "grp_100",
				LatestVersion:            5,
				ProductionCertType:       ptr.To("CPS_MANAGED"),
				ProductionCnameTo:        ptr.To("api.example.com.edgesuite.net"),
				ProductionCnameType:      ptr.To("EDGE_HOSTNAME"),
				ProductionEdgeHostnameID: ptr.To("ehn_1001"),
				ProductionProductID:      ptr.To("prd_SPM"),
				PropertyID:               "prp_api",
				PropertyName:             "api-property",
				PropertyType:             "STANDARD",
				StagingCertType:          ptr.To("CPS_MANAGED"),
				StagingCnameTo:           ptr.To("api.example.com.edgesuite-staging.net"),
				StagingCnameType:         ptr.To("EDGE_HOSTNAME"),
				StagingEdgeHostnameID:    ptr.To("ehn_2001"),
				StagingProductID:         ptr.To("prd_SPM"),
			},
			{
				CnameFrom:                "www.example.com",
				ContractID:               "ctr_2-CD456",
				GroupID:                  "grp_200",
				LatestVersion:            3,
				ProductionCertType:       ptr.To("THIRD_PARTY"),
				ProductionCnameTo:        ptr.To("www.example.com.akamaized.net"),
				ProductionCnameType:      ptr.To("EDGE_HOSTNAME"),
				ProductionEdgeHostnameID: ptr.To("ehn_1002"),
				ProductionProductID:      ptr.To("prd_FRESCA"),
				PropertyID:               "prp_www",
				PropertyName:             "www-property",
				PropertyType:             "STANDARD",
				StagingCertType:          ptr.To("THIRD_PARTY"),
				StagingCnameTo:           ptr.To("www.example.com.akamaized-staging.net"),
				StagingCnameType:         ptr.To("EDGE_HOSTNAME"),
				StagingEdgeHostnameID:    ptr.To("ehn_2002"),
				StagingProductID:         ptr.To("prd_FRESCA"),
			},
			{
				CnameFrom:                "cdn.example.com",
				ContractID:               "ctr_3-EF789",
				GroupID:                  "grp_300",
				LatestVersion:            10,
				ProductionCertType:       ptr.To("CPS_MANAGED"),
				ProductionCnameTo:        ptr.To("cdn.example.com.edgekey.net"),
				ProductionCnameType:      ptr.To("EDGE_HOSTNAME"),
				ProductionEdgeHostnameID: ptr.To("ehn_1003"),
				ProductionProductID:      ptr.To("prd_DOWNLOAD_DELIVERY"),
				PropertyID:               "prp_cdn",
				PropertyName:             "cdn-property",
				PropertyType:             "MEDIA",
				StagingCertType:          ptr.To("CPS_MANAGED"),
				StagingCnameTo:           ptr.To("cdn.example.com.edgekey-staging.net"),
				StagingCnameType:         ptr.To("EDGE_HOSTNAME"),
				StagingEdgeHostnameID:    ptr.To("ehn_2003"),
				StagingProductID:         ptr.To("prd_DOWNLOAD_DELIVERY"),
			},
		},
	}

	checker := test.NewStateChecker("data.akamai_property_account_hostnames.test").
		CheckEqual("account_id", "act_1").
		CheckEqual("current_sort", "hostname:a").
		CheckEqual("default_sort", "hostname:a").
		CheckEqual("available_sort.#", "2").
		CheckEqual("hostnames.#", "3").
		CheckEqual("hostnames.0.cname_from", "api.example.com").
		CheckEqual("hostnames.0.contract_id", "ctr_1-AB123").
		CheckEqual("hostnames.0.group_id", "grp_100").
		CheckEqual("hostnames.0.latest_version", "5").
		CheckEqual("hostnames.0.production_cert_type", "CPS_MANAGED").
		CheckEqual("hostnames.0.production_cname_to", "api.example.com.edgesuite.net").
		CheckEqual("hostnames.0.production_cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.0.production_edge_hostname_id", "ehn_1001").
		CheckEqual("hostnames.0.production_product_id", "prd_SPM").
		CheckEqual("hostnames.0.property_id", "prp_api").
		CheckEqual("hostnames.0.property_name", "api-property").
		CheckEqual("hostnames.0.property_type", "STANDARD").
		CheckEqual("hostnames.0.staging_cert_type", "CPS_MANAGED").
		CheckEqual("hostnames.0.staging_cname_to", "api.example.com.edgesuite-staging.net").
		CheckEqual("hostnames.0.staging_cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.0.staging_edge_hostname_id", "ehn_2001").
		CheckEqual("hostnames.0.staging_product_id", "prd_SPM").
		CheckEqual("hostnames.1.cname_from", "www.example.com").
		CheckEqual("hostnames.1.contract_id", "ctr_2-CD456").
		CheckEqual("hostnames.1.group_id", "grp_200").
		CheckEqual("hostnames.1.latest_version", "3").
		CheckEqual("hostnames.1.production_cert_type", "THIRD_PARTY").
		CheckEqual("hostnames.1.production_cname_to", "www.example.com.akamaized.net").
		CheckEqual("hostnames.1.production_cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.1.production_edge_hostname_id", "ehn_1002").
		CheckEqual("hostnames.1.production_product_id", "prd_FRESCA").
		CheckEqual("hostnames.1.property_id", "prp_www").
		CheckEqual("hostnames.1.property_name", "www-property").
		CheckEqual("hostnames.1.property_type", "STANDARD").
		CheckEqual("hostnames.1.staging_cert_type", "THIRD_PARTY").
		CheckEqual("hostnames.1.staging_cname_to", "www.example.com.akamaized-staging.net").
		CheckEqual("hostnames.1.staging_cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.1.staging_edge_hostname_id", "ehn_2002").
		CheckEqual("hostnames.1.staging_product_id", "prd_FRESCA").
		CheckEqual("hostnames.2.cname_from", "cdn.example.com").
		CheckEqual("hostnames.2.contract_id", "ctr_3-EF789").
		CheckEqual("hostnames.2.group_id", "grp_300").
		CheckEqual("hostnames.2.latest_version", "10").
		CheckEqual("hostnames.2.production_cert_type", "CPS_MANAGED").
		CheckEqual("hostnames.2.production_cname_to", "cdn.example.com.edgekey.net").
		CheckEqual("hostnames.2.production_cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.2.production_edge_hostname_id", "ehn_1003").
		CheckEqual("hostnames.2.production_product_id", "prd_DOWNLOAD_DELIVERY").
		CheckEqual("hostnames.2.property_id", "prp_cdn").
		CheckEqual("hostnames.2.property_name", "cdn-property").
		CheckEqual("hostnames.2.property_type", "MEDIA").
		CheckEqual("hostnames.2.staging_cert_type", "CPS_MANAGED").
		CheckEqual("hostnames.2.staging_cname_to", "cdn.example.com.edgekey-staging.net").
		CheckEqual("hostnames.2.staging_cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.2.staging_edge_hostname_id", "ehn_2003").
		CheckEqual("hostnames.2.staging_product_id", "prd_DOWNLOAD_DELIVERY")

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"happy path - get account hostnames with all filters": {
			init: func(m *papi.Mock) {
				mockListAccountHostnames(m, accountHostnamesTestData{
					req: papi.ListActiveAccountHostnamesRequest{
						ContractID: "ctr_1",
						GroupID:    "grp_1",
						Hostname:   "example.com",
						CnameTo:    "example.com.edgesuite.net",
						Network:    papi.ActivationNetworkProduction,
						Sort:       papi.SortAscending,
					},
					hostnames: defaultHostnames,
				})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/valid_with_filters.tf"),
					Check:  checker.Build(),
				},
			},
		},
		"happy path - get account hostnames with hostname filter": {
			init: func(m *papi.Mock) {
				mockListAccountHostnames(m, accountHostnamesTestData{
					req: papi.ListActiveAccountHostnamesRequest{
						Hostname: "example.com",
					},
					hostnames: defaultHostnames,
				})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/valid_with_hostname_filter.tf"),
					Check:  checker.Build(),
				},
			},
		},
		"happy path - staging only attributes": {
			init: func(m *papi.Mock) {
				stagingOnlyHostnames := papi.ActiveAccountHostnames{
					Items: []papi.ActiveAccountHostnameItem{
						{
							CnameFrom:             "staging-only.com",
							ContractID:            "ctr_1",
							GroupID:               "grp_1",
							LatestVersion:         1,
							PropertyID:            "prp_staging",
							PropertyName:          "staging-property",
							PropertyType:          "STANDARD",
							StagingCertType:       ptr.To("CPS_MANAGED"),
							StagingCnameTo:        ptr.To("staging-only.com.edgesuite-staging.net"),
							StagingCnameType:      ptr.To("EDGE_HOSTNAME"),
							StagingEdgeHostnameID: ptr.To("ehn_staging_789"),
							StagingProductID:      ptr.To("prd_SPM"),
						},
					},
					CurrentItemCount: 1,
					TotalItems:       1,
				}
				mockListAccountHostnames(m, accountHostnamesTestData{
					hostnames: stagingOnlyHostnames,
				})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/valid_minimal.tf"),
					Check: test.NewStateChecker("data.akamai_property_account_hostnames.test").
						CheckEqual("account_id", "act_1").
						CheckEqual("hostnames.#", "1").
						CheckEqual("hostnames.0.cname_from", "staging-only.com").
						CheckEqual("hostnames.0.contract_id", "ctr_1").
						CheckEqual("hostnames.0.group_id", "grp_1").
						CheckEqual("hostnames.0.latest_version", "1").
						CheckEqual("hostnames.0.property_id", "prp_staging").
						CheckEqual("hostnames.0.property_name", "staging-property").
						CheckEqual("hostnames.0.property_type", "STANDARD").
						CheckEqual("hostnames.0.staging_cert_type", "CPS_MANAGED").
						CheckEqual("hostnames.0.staging_cname_to", "staging-only.com.edgesuite-staging.net").
						CheckEqual("hostnames.0.staging_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.staging_edge_hostname_id", "ehn_staging_789").
						CheckEqual("hostnames.0.staging_product_id", "prd_SPM").
						CheckMissing("hostnames.0.production_cert_type").
						CheckMissing("hostnames.0.production_cname_to").
						CheckMissing("hostnames.0.production_cname_type").
						CheckMissing("hostnames.0.production_edge_hostname_id").
						CheckMissing("hostnames.0.production_product_id").
						Build(),
				},
			},
		},
		"happy path - production only attributes": {
			init: func(m *papi.Mock) {
				productionOnlyHostnames := papi.ActiveAccountHostnames{
					Items: []papi.ActiveAccountHostnameItem{
						{
							CnameFrom:                "production-only.com",
							ContractID:               "ctr_1",
							GroupID:                  "grp_1",
							LatestVersion:            1,
							PropertyID:               "prp_production",
							PropertyName:             "production-property",
							PropertyType:             "STANDARD",
							ProductionCertType:       ptr.To("CPS_MANAGED"),
							ProductionCnameTo:        ptr.To("production-only.com.edgesuite.net"),
							ProductionCnameType:      ptr.To("EDGE_HOSTNAME"),
							ProductionEdgeHostnameID: ptr.To("ehn_production_789"),
							ProductionProductID:      ptr.To("prd_SPM"),
						},
					},
					CurrentItemCount: 1,
					TotalItems:       1,
				}
				mockListAccountHostnames(m, accountHostnamesTestData{
					hostnames: productionOnlyHostnames,
				})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/valid_minimal.tf"),
					Check: test.NewStateChecker("data.akamai_property_account_hostnames.test").
						CheckEqual("account_id", "act_1").
						CheckEqual("hostnames.#", "1").
						CheckEqual("hostnames.0.cname_from", "production-only.com").
						CheckEqual("hostnames.0.contract_id", "ctr_1").
						CheckEqual("hostnames.0.group_id", "grp_1").
						CheckEqual("hostnames.0.latest_version", "1").
						CheckEqual("hostnames.0.property_id", "prp_production").
						CheckEqual("hostnames.0.property_name", "production-property").
						CheckEqual("hostnames.0.property_type", "STANDARD").
						CheckEqual("hostnames.0.production_cert_type", "CPS_MANAGED").
						CheckEqual("hostnames.0.production_cname_to", "production-only.com.edgesuite.net").
						CheckEqual("hostnames.0.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.production_edge_hostname_id", "ehn_production_789").
						CheckEqual("hostnames.0.production_product_id", "prd_SPM").
						CheckMissing("hostnames.0.staging_cert_type").
						CheckMissing("hostnames.0.staging_cname_to").
						CheckMissing("hostnames.0.staging_cname_type").
						CheckMissing("hostnames.0.staging_edge_hostname_id").
						CheckMissing("hostnames.0.staging_product_id").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *papi.Mock) {
				m.On("ListActiveAccountHostnames", mock.Anything, papi.ListActiveAccountHostnamesRequest{}).
					Return(nil, fmt.Errorf("API error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/valid_minimal.tf"),
					ExpectError: regexp.MustCompile("API error"),
				},
			},
		},
		"validation error - invalid network": {
			init: func(_ *papi.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/invalid_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute network value must be one of: \["STAGING" "PRODUCTION"\], got:(\n|.)"INVALID"`),
				},
			},
		},
		"validation error - invalid sort": {
			init: func(_ *papi.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyAccountHostnames/invalid_sort.tf"),
					ExpectError: regexp.MustCompile(`Attribute sort value must be one of: \["hostname:a" "hostname:d"\], got:(\n|.)"invalid:sort"`),
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

func mockListAccountHostnames(m *papi.Mock, data accountHostnamesTestData) *mock.Call {
	return m.On("ListActiveAccountHostnames", mock.Anything, data.req).
		Return(&papi.ListActiveAccountHostnamesResponse{
			AccountID:     "act_1",
			AvailableSort: []string{"hostname:a", "hostname:d"},
			CurrentSort:   "hostname:a",
			DefaultSort:   "hostname:a",
			Hostnames:     data.hostnames,
		}, nil).Times(3)
}
