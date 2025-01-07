package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDSPropertiesSearch(t *testing.T) {
	t.Run("match by hostname", func(t *testing.T) {
		client := &papi.Mock{}

		search := papi.SearchItems{Items: []papi.SearchItem{
			{
				AccountID:        "acc_test",
				AssetID:          "aid_test",
				GroupID:          "grp_test",
				ContractID:       "ctr_test",
				EdgeHostname:     "www.example.com.edgekey.net",
				Hostname:         "www.example.com",
				ProductionStatus: "INACTIVE",
				StagingStatus:    "ACTIVE",
				PropertyID:       "prp_test",
				PropertyName:     "test_www.example.com",
				PropertyVersion:  1,
				UpdatedByUser:    "test_user@example.com",
				UpdatedDate:      "2021-11-22T07:24:56Z",
			},
			{
				AccountID:        "acc_test",
				AssetID:          "aid_test",
				GroupID:          "grp_test",
				ContractID:       "ctr_test",
				EdgeHostname:     "www.example.com.edgekey.net",
				Hostname:         "www.example.com",
				ProductionStatus: "ACTIVE",
				StagingStatus:    "INACTIVE",
				PropertyID:       "prp_test1",
				PropertyName:     "test1_www.example.com",
				PropertyVersion:  1,
				UpdatedByUser:    "test_user@example.com",
				UpdatedDate:      "2021-11-22T07:24:56Z",
			},
		}}

		client.On("SearchProperties",
			testutils.MockContext,
			papi.SearchRequest{Key: "hostname", Value: "www.example.com"},
		).Return(&papi.SearchResponse{Versions: search}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertiesSearch/match_by_hostname.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "id", "hostname:www.example.com"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.#", "2"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.account_id", "acc_test"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.asset_id", "aid_test"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.contract_id", "ctr_test"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.group_id", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.property_id", "prp_test"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.1.property_id", "prp_test1"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.property_version", "1"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.property_name", "test_www.example.com"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.1.property_name", "test1_www.example.com"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.edge_hostname", "www.example.com.edgekey.net"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.hostname", "www.example.com"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.production_status", "INACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.1.production_status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.staging_status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.1.staging_status", "INACTIVE"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.updated_by_user", "test_user@example.com"),
						resource.TestCheckResourceAttr("data.akamai_properties_search.test", "properties.0.updated_date", "2021-11-22T07:24:56Z"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("search error", func(t *testing.T) {
		client := &papi.Mock{}

		client.On("SearchProperties",
			testutils.MockContext,
			papi.SearchRequest{Key: "hostname", Value: "www.example.com"},
		).Return(nil, papi.ErrSearchProperties)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertiesSearch/match_by_hostname.tf"),
					ExpectError: regexp.MustCompile("searching for properties"),
				}},
			})
		})

		client.AssertExpectations(t)
	})
}
