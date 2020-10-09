package property

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

// Alias of mock.Anything to use as a placeholder for any context.Context
var AnyCTX = mock.Anything

func TestResCPCode(t *testing.T) {
	// Helper to set up an expected call to mock papi.GetCPCodes with mock impl backed by the given slice
	expectGet := func(m *mockpapi, ContractID, GroupID string, CPCodes *[]papi.CPCode) *mock.Call {
		mockImpl := func(_ context.Context, req papi.GetCPCodesRequest) (*papi.GetCPCodesResponse, error) {
			res := &papi.GetCPCodesResponse{
				ContractID: req.ContractID,
				GroupID:    req.GroupID,
				CPCodes:    papi.CPCodeItems{Items: *CPCodes},
			}
			return res, nil
		}

		req := papi.GetCPCodesRequest{ContractID: ContractID, GroupID: GroupID}

		return m.OnGetCPCodes(mockImpl, AnyCTX, req)
	}

	// Helper to set up an expected call to mock papi.CreateCPCode with mock impl backed by the given slice
	expectCreate := func(m *mockpapi, CPCName, Product, Contract, Group string, CPCodes *[]papi.CPCode) *mock.Call {
		mockImpl := func(_ context.Context, req papi.CreateCPCodeRequest) (*papi.CreateCPCodeResponse, error) {
			cpc := papi.CPCode{
				ID:        fmt.Sprintf("cpc_%s:%s:%d", req.ContractID, req.GroupID, len(*CPCodes)),
				Name:      req.CPCode.CPCodeName,
				ProductIDs: []string{req.CPCode.ProductID},
			}

			*CPCodes = append(*CPCodes, cpc)
			res := &papi.CreateCPCodeResponse{CPCodeID: cpc.ID}

			return res, nil
		}

		req := papi.CreateCPCodeRequest{
			ContractID: Contract,
			GroupID:    Group,
			CPCode: papi.CreateCPCode{
				ProductID:  Product,
				CPCodeName: CPCName,
			},
		}

		return m.OnCreateCPCode(mockImpl, AnyCTX, req)
	}

	t.Run("create new CP Code", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGet(client, "ctr1", "grp1", &CPCodes)
		expectCreate(client, "test cpcode", "prd1", "ctr1", "grp1", &CPCodes)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/create_new_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd1"),
					),
				}},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})

	t.Run("use existing CP Code", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test"}}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGet(client, "ctr_test", "grp_test", &CPCodes)
		// No mock behavior for create because we're using an existing CP code

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_test"),
					),
				}},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})

	t.Run("product missing from CP Code", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_test2", Name: "test cpcode"}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGet(client, "ctr_test", "grp_test", &CPCodes)
		// No mock behavior for create because we're using an existing CP code

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					ExpectError:regexp.MustCompile("Couldn't find product id on the CP Code"),

				}},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})

	t.Run("change name", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGet(client, "ctr1", "grp1", &CPCodes)
		expectCreate(client, "test cpcode", "prd1", "ctr1", "grp1", &CPCodes).Once()
		expectCreate(client, "renamed cpcode", "prd1", "ctr1", "grp1", &CPCodes).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check:  resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:0"),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						Check:  resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:1"),
					},
				},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})

	t.Run("change group", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI for first group
		CPCodes1 := []papi.CPCode{}

		// Contains CP Codes known to mock PAPI for second group
		CPCodes2 := []papi.CPCode{}

		// Values are from fixture:
		expectGet(client, "ctr1", "grp1", &CPCodes1)
		expectCreate(client, "test cpcode", "prd1", "ctr1", "grp1", &CPCodes1).Once()

		expectGet(client, "ctr1", "grp2", &CPCodes2)
		expectCreate(client, "test cpcode", "prd1", "ctr1", "grp2", &CPCodes2).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_group_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_group_step1.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp2:0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp2"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd1"),
						),
					},
				},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})

	t.Run("change contract", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock API for first contract
		CPCodes1 := []papi.CPCode{}

		// Contains CP Codes known to mock API for second contract
		CPCodes2 := []papi.CPCode{}

		// Values are from fixture:
		expectGet(client, "ctr1", "grp1", &CPCodes1)
		expectCreate(client, "test cpcode", "prd1", "ctr1", "grp1", &CPCodes1).Once()

		expectGet(client, "ctr2", "grp1", &CPCodes2)
		expectCreate(client, "test cpcode", "prd1", "ctr2", "grp1", &CPCodes2).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_contract_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_contract_step1.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr2:grp1:0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr2"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd1"),
						),
					},
				},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})

	t.Run("change product", func(t *testing.T) {
		t.Log("eelmore: My expectation is that we create a new CP code with a new CP Code ID and new Product ID when")
		t.Log("         only the Product ID changesInstead of that, the provider is attempting to re-use the old CP Code")
		t.Log("         which has the old Product ID.")

		TODO(t, "Is this case fundamentally broken by current logic?")

		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock API for first group
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGet(client, "ctr1", "grp1", &CPCodes)
		expectCreate(client, "test cpcode", "prd1", "ctr1", "grp1", &CPCodes).Once()
		expectCreate(client, "test cpcode", "prd2", "ctr1", "grp1", &CPCodes).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_product_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_product_step1.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_ctr1:grp1:1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd2"),
						),
					},
				},
				CheckDestroy: resource.TestCheckNoResourceAttr("akamai_cp_code.test", "id"),
			})
		})
	})
}
