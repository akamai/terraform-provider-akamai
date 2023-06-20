package property

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

// Alias of mock.Anything to use as a placeholder for any context.Context
var AnyCTX = mock.Anything

func TestResCPCode(t *testing.T) {
	// Helper to set up an expected call to mock papi.GetCPCode
	expectGetCPCode := func(m *papi.Mock, contractID, groupID string, CPCodeID int, CPCodes *[]papi.CPCode, err error) *mock.Call {
		var call *mock.Call
		req := papi.GetCPCodeRequest{CPCodeID: strconv.Itoa(CPCodeID), ContractID: contractID, GroupID: groupID}

		call = m.On("GetCPCode", AnyCTX, req).Run(func(args mock.Arguments) {
			if err != nil {
				call.Return(nil, err)
			} else {
				CPCode := (*CPCodes)[CPCodeID]
				res := &papi.GetCPCodesResponse{
					CPCode: papi.CPCode{
						ID:         CPCode.ID,
						Name:       CPCode.Name,
						ProductIDs: CPCode.ProductIDs,
					},
				}
				call.Return(res, nil)
			}
		})
		return call
	}

	// Helper to set up an expected call to mock papi.GetCPCodes with mock impl backed by the given slice
	expectGetCPCodes := func(m *papi.Mock, ContractID, GroupID string, CPCodes *[]papi.CPCode) *mock.Call {
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
	expectCreateCPCode := func(m *papi.Mock, CPCName, Product, Contract, Group string, CPCodes *[]papi.CPCode) *mock.Call {
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

	// Helper to set up an expected call to mock papi.UpdateCPCode with mock impl backed by the given slice
	expectUpdateCPCode := func(m *papi.Mock, CPCodeID int, name string, CPCodes, CPCodesCopy *[]papi.CPCode, err error) *mock.Call {
		mockImpl := func(_ context.Context, req papi.UpdateCPCodeRequest) (*papi.CPCodeDetailResponse, error) {
			if err != nil {
				return nil, err
			}
			copy(*CPCodesCopy, *CPCodes)
			(*CPCodes)[CPCodeID].Name = name

			res := &papi.CPCodeDetailResponse{
				ID:   req.ID,
				Name: req.Name,
			}

			return res, nil
		}

		f := false
		req := papi.UpdateCPCodeRequest{
			ID:               CPCodeID,
			Name:             name,
			Purgeable:        &f,
			OverrideTimeZone: &papi.CPCodeTimeZone{},
		}

		return m.OnUpdateCPCode(mockImpl, AnyCTX, req)
	}

	// // Helper to set up an expected call to mock papi.GetCPCodeDetail
	expectGetCPCodeDetail := func(m *papi.Mock, CPCodeID int, CPCodes *[]papi.CPCode, err error) *mock.Call {
		var call *mock.Call

		call = m.On("GetCPCodeDetail", AnyCTX, CPCodeID).Run(func(args mock.Arguments) {
			if err != nil {
				call.Return(nil, err)
			} else {
				CPCode := (*CPCodes)[CPCodeID]
				res := &papi.CPCodeDetailResponse{
					ID:   CPCodeID,
					Name: CPCode.Name,
				}
				call.Return(res, nil)
			}
		})
		return call
	}

	// redefining times to accelerate tests
	updatePollMinimum = time.Millisecond * 1
	updatePollInterval = updatePollMinimum

	t.Run("create new CP Code", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_1", "grp_1", &CPCodes).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes)
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(2)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/create_new_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code with multiple products", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "0", Name: "test cpcode", ProductIDs: []string{"prd_test", "prd_wrong", "another_wrong"}}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_test", "grp_test", &CPCodes).Once()
		// No mock behavior for create because we're using an existing CP code

		// Read and plan
		expectGetCPCode(client, "ctr_test", "grp_test", 0, &CPCodes, nil).Times(2)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_test"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "0", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_1", Name: "test cpcode", ProductIDs: []string{"prd_test"}}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_test", "grp_test", &CPCodes).Once()
		// No mock behavior for create because we're using an existing CP code

		// Read and plan
		expectGetCPCode(client, "ctr_test", "grp_test", 1, &CPCodes, nil).Times(2)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_test"),
					),
				}},
			})
		})
	})

	t.Run("product missing from CP Code", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "0", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_1", Name: "test cpcode"}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_test", "grp_test", &CPCodes).Once()
		// No mock behavior for create because we're using an existing CP code

		// Read and plan
		expectGetCPCode(client, "ctr_test", "grp_test", 1, &CPCodes, nil).Times(1)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					ExpectError: regexp.MustCompile("Couldn't find product id on the CP Code"),
				}},
			})
		})
	})

	t.Run("change name", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}
		CPCodesCopy := make([]papi.CPCode, 1)

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_1", "grp_1", &CPCodes).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(3)

		expectGetCPCodeDetail(client, 0, &CPCodes, nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", &CPCodes, &CPCodesCopy, nil).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodesCopy, nil).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(3)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "renamed cpcode"),
						),
					},
				},
			})
		})
	})

	t.Run("import existing cp code", func(t *testing.T) {
		client := &papi.Mock{}
		id := "0,1,2"

		CPCodes := []papi.CPCode{{ID: "0", Name: "test cpcode", ProductIDs: []string{"prd_Web_Accel"}}}
		expectGetCPCodes(client, "ctr_1", "grp_2", &CPCodes)
		expectGetCPCode(client, "ctr_1", "grp_2", 0, &CPCodes, nil).Times(4)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
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
							assert.Equal(t, "prd_Web_Accel", rs.Attributes["product_id"])
							assert.Equal(t, "0", rs.Attributes["id"])
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
		client := &papi.Mock{}
		id := "123"

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
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

	t.Run("empty CP code ID passed", func(t *testing.T) {
		client := &papi.Mock{}
		id := ",ctr_1-1NC95D,grp_194665"

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:        loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cp_code.test",
						ExpectError:   regexp.MustCompile("CP Code is a mandatory parameter"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("immutable attributes updated", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_1", "grp_1", &CPCodes).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(5)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_immutable.tf"),
						ExpectError: regexp.MustCompile(`cp code attribute 'contract_id' cannot be changed after creation \(immutable\)`),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_immutable.tf"),
						ExpectError: regexp.MustCompile(`cp code attribute 'product_id' cannot be changed after creation \(immutable\)`),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_immutable.tf"),
						ExpectError: regexp.MustCompile(`cp code attribute 'group_id' cannot be changed after creation \(immutable\)`),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
						),
					},
				},
			})
		})
	})

	t.Run("error fetching cpCode details", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_1", "grp_1", &CPCodes).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(3)

		expectGetCPCodeDetail(client, 0, &CPCodes, fmt.Errorf("oops")).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
	})

	t.Run("error updating cpCode", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_1", "grp_1", &CPCodes).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(3)

		expectGetCPCodeDetail(client, 0, &CPCodes, nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", &CPCodes, &[]papi.CPCode{}, fmt.Errorf("oops")).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
	})

	t.Run("timeout waiting for update", func(t *testing.T) {
		timeoutVal := cpCodeResourceUpdateTimeout
		oldInterval := cpCodeResourceUpdateTimeout

		cpCodeResourceUpdateTimeout = time.Millisecond * 6
		updatePollInterval = time.Millisecond * 4

		client := &papi.Mock{}

		defer func() {
			cpCodeResourceUpdateTimeout = timeoutVal
			updatePollInterval = oldInterval
			client.AssertExpectations(t)
		}()

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}
		CPCodesCopy := make([]papi.CPCode, 1)

		// Values are from fixture:
		expectGetCPCodes(client, "ctr_1", "grp_1", &CPCodes).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodes, nil).Times(3)

		expectGetCPCodeDetail(client, 0, &CPCodes, nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", &CPCodes, &CPCodesCopy, nil).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, &CPCodesCopy, nil).Times(3)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:             loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
	})

	t.Run("error when no product and product_id provided", func(t *testing.T) {
		expectedErr := regexp.MustCompile("`product_id` must be specified for creation")

		resource.UnitTest(t, resource.TestCase{
			ProtoV5ProviderFactories: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config:      loadFixtureString("testdata/TestResCPCode/missing_product.tf"),
					ExpectError: expectedErr,
				},
			},
		})
	})
}
