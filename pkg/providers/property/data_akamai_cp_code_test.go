package property

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
)

func TestDSCPCode(t *testing.T) {
	t.Run("match by name", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "contract", "ctr_test"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by name output products", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_name_output_products.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "contract", "ctr_test"),
						resource.TestCheckOutput("product1", "prd_test1"),
						resource.TestCheckOutput("product2", "prd_test2"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by full ID", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "cpc_test2"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_full_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "name", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "contract", "ctr_test"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by unprefixed ID", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test2"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_unprefixed_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("no matches", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test3", Name: "Also wrong CP code"},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSCPCode/match_by_unprefixed_id.tf"),
					ExpectError: regexp.MustCompile(`cp code not found`),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("contract collides with contract ID", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			Providers:  testAccProviders,
			IsUnitTest: true,
			Steps: []resource.TestStep{{
				Config:      loadFixtureString("testdata/TestDSCPCode/contract_collides_with_id.tf"),
				ExpectError: regexp.MustCompile("only one of `contract,contract_id` can be specified"),
			}},
		})
	})

	t.Run("group collides with group ID", func(t *testing.T) {
		client := &mockpapi{}
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSCPCode/group_collides_with_id.tf"),
					ExpectError: regexp.MustCompile("only one of `group,group_id` can be specified"),
				}},
			})
		})
	})

	t.Run("group not found in state", func(t *testing.T) {
		client := &mockpapi{}
		client.On("GetCPCodes",
			AnyCTX, mock.Anything,
		).Return(&papi.GetCPCodesResponse{CPCodes: papi.CPCodeItems{Items: []papi.CPCode{{
			ID: "cpc_test-ft-cp-code", Name: "test-ft-cp-code", CreatedDate: "", ProductIDs: []string{"prd_prod1"},
		}}}}, nil)
		client.On("CreateCPCode", AnyCTX, mock.Anything).Return(&papi.CreateCPCodeResponse{}, nil)
		client.On("GetCPCode", AnyCTX, mock.Anything).Return(&papi.GetCPCodesResponse{CPCode: papi.CPCode{
			ID: "cpc_test-ft-cp-code", Name: "test-ft-cp-code", CreatedDate: "", ProductIDs: []string{"prd_prod1"},
		}}, nil).Times(3)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSGroupNotFound/cp_code.tf"),
				}},
			})
		})
	})
}
