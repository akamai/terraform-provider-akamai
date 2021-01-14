package property

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tj/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
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
				ID:         fmt.Sprintf("cpc_%d", len(*CPCodes)),
				Name:       req.CPCode.CPCodeName,
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
		expectGet(client, "ctr_1", "grp_1", &CPCodes)
		expectCreate(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/create_new_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_1"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code with multiple products", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test", "prd_wrong", "another_wrong"}}, // Matches name from fixture
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
					Config:      loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					ExpectError: regexp.MustCompile("Couldn't find product id on the CP Code"),
				}},
			})
		})
	})

	t.Run("change name", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGet(client, "ctr_1", "grp_1", &CPCodes)
		expectCreate(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectCreate(client, "renamed cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check:  resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						Check:  resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_1"),
					},
				},
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
		expectGet(client, "ctr_1", "grp_1", &CPCodes)
		expectCreate(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectCreate(client, "test cpcode", "prd_2", "ctr_1", "grp_1", &CPCodes).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_product_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_product_step1.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_2"),
						),
					},
				},
			})
		})
	})

	t.Run("import existing cp code", func(t *testing.T) {
		client := &mockpapi{}
		id := "123,1,2"

		cpCodes := []papi.CPCode{{ID: "cpc_123", Name: "test cpcode", ProductIDs: []string{"prd_2"}}}
		expectGet(client, "ctr_1", "grp_2", &cpCodes)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:             loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
						ExpectNonEmptyPlan: true,
					},
					{
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cp_code.test",
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "prd_2", rs.Attributes["product"])
							assert.Equal(t, "cpc_123", rs.Attributes["id"])
							assert.Equal(t, "test cpcode", rs.Attributes["name"])
							return nil
						},
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("invalid import ID passed", func(t *testing.T) {
		client := &mockpapi{}
		id := "123"

		TODO(t, "error assertion in import is impossible using provider testing framework as it only checks for errors in `apply`")
		cpCodes := []papi.CPCode{{ID: "cpc_123", Name: "test cpcode", ProductIDs: []string{"prd_2"}}}
		expectGet(client, "ctr_1", "grp_2", &cpCodes)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:        loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cp_code.test",
						ExpectError:   regexp.MustCompile("comma-separated list of CP code ID, contract ID and group ID has to be supplied in import"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}
