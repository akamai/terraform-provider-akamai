package property

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResCPCode(t *testing.T) {
	expectGetCPCode := func(m *papi.Mock, contractID, groupID string, CPCodeID int, CPCodeName string, CPCodeProductIDs []string, err error) *mock.Call {
		req := papi.GetCPCodeRequest{CPCodeID: strconv.Itoa(CPCodeID), ContractID: contractID, GroupID: groupID}
		var res *papi.GetCPCodesResponse
		if err == nil {
			res = &papi.GetCPCodesResponse{
				CPCode: papi.CPCode{
					ID:         fmt.Sprintf("%d", CPCodeID),
					Name:       CPCodeName,
					ProductIDs: CPCodeProductIDs,
				},
			}
		}
		return m.On("GetCPCode", testutils.MockContext, req).Return(res, nil)
	}

	expectGetCPCodes := func(m *papi.Mock, ContractID, GroupID string, CPCodes []papi.CPCode) *mock.Call {
		req := papi.GetCPCodesRequest{ContractID: ContractID, GroupID: GroupID}
		res := &papi.GetCPCodesResponse{
			ContractID: req.ContractID,
			GroupID:    req.GroupID,
			CPCodes:    papi.CPCodeItems{Items: CPCodes},
		}

		return m.On("GetCPCodes", testutils.MockContext, req).Return(res, nil)
	}

	expectCreateCPCode := func(m *papi.Mock, CPCName, Product, Contract, Group string) *mock.Call {
		req := papi.CreateCPCodeRequest{
			ContractID: Contract,
			GroupID:    Group,
			CPCode: papi.CreateCPCode{
				ProductID:  Product,
				CPCodeName: CPCName,
			},
		}
		cpc := papi.CPCode{
			ID:         "cpc_0",
			Name:       req.CPCode.CPCodeName,
			ProductIDs: []string{req.CPCode.ProductID},
		}

		res := &papi.CreateCPCodeResponse{CPCodeID: cpc.ID}

		return m.On("CreateCPCode", testutils.MockContext, req).Return(res, nil)
	}

	expectUpdateCPCode := func(m *papi.Mock, CPCodeID int, name string, err error) *mock.Call {
		var res *papi.CPCodeDetailResponse

		f := false
		req := papi.UpdateCPCodeRequest{
			ID:               CPCodeID,
			Name:             name,
			Purgeable:        &f,
			OverrideTimeZone: &papi.CPCodeTimeZone{},
		}

		if err == nil {
			res = &papi.CPCodeDetailResponse{
				ID:   req.ID,
				Name: req.Name,
			}

		}

		return m.On("UpdateCPCode", testutils.MockContext, req).Return(res, err)
	}

	expectGetCPCodeDetail := func(m *papi.Mock, CPCodeID int, CPCodeName string, err error) *mock.Call {
		var res *papi.CPCodeDetailResponse
		if err == nil {
			res = &papi.CPCodeDetailResponse{
				ID:   CPCodeID,
				Name: CPCodeName,
			}
		}
		return m.On("GetCPCodeDetail", testutils.MockContext, CPCodeID).Return(res, err)
	}

	// redefining times to accelerate tests
	updatePollMinimum = time.Millisecond * 1
	updatePollInterval = updatePollMinimum

	t.Run("create new CP Code", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		expectGetCPCodes(client, "ctr_1", "grp_1", nil).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1")
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(2)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/create_new_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "timeouts.0.update", "1h"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code with multiple products", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		CPCodes := []papi.CPCode{
			{ID: "0", Name: "test cpcode", ProductIDs: []string{"prd_test", "prd_wrong", "another_wrong"}},
		}

		expectGetCPCodes(client, "ctr_test", "grp_test", CPCodes).Once()
		// No mock behavior for create because we're using an existing CP code

		// Read and plan
		expectGetCPCode(client, "ctr_test", "grp_test", 0, "test cpcode", []string{"prd_test", "prd_wrong", "another_wrong"}, nil).Times(2)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/use_existing_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "timeouts.#", "0"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		CPCodes := []papi.CPCode{
			{ID: "0", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_1", Name: "test cpcode", ProductIDs: []string{"prd_test"}},
		}

		expectGetCPCodes(client, "ctr_test", "grp_test", CPCodes).Once()
		// No mock behavior for create because we're using an existing CP code

		// Read and plan
		expectGetCPCode(client, "ctr_test", "grp_test", 1, "test cpcode", []string{"prd_test"}, nil).Times(2)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/use_existing_cp_code.tf"),
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
		expectGetCPCodes(client, "ctr_test", "grp_test", CPCodes).Once()
		// No mock behavior for create because we're using an existing CP code

		// Read and plan
		expectGetCPCode(client, "ctr_test", "grp_test", 1, "test cpcode", nil, nil).Times(1)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/use_existing_cp_code.tf"),
					ExpectError: regexp.MustCompile("Couldn't find product id on the CP Code"),
				}},
			})
		})
	})

	t.Run("change name", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		expectGetCPCodes(client, "ctr_1", "grp_1", nil).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1").Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(3)

		expectGetCPCodeDetail(client, 0, "test cpcode", nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", nil).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "renamed cpcode", []string{"prd_1"}, nil).Times(3)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
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
		expectGetCPCodes(client, "ctr_1", "grp_2", CPCodes)
		expectGetCPCode(client, "ctr_1", "grp_2", 0, "test cpcode", []string{"prd_Web_Accel"}, nil).Times(4)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
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
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:        testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
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
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:        testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
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

		expectGetCPCodes(client, "ctr_1", "grp_1", nil).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1").Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(5)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_immutable.tf"),
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
						Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_immutable.tf"),
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
						Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_immutable.tf"),
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

		expectGetCPCodes(client, "ctr_1", "grp_1", nil).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1").Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(3)

		expectGetCPCodeDetail(client, 0, "test cpcode", fmt.Errorf("oops")).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
	})

	t.Run("error updating cpCode", func(t *testing.T) {
		client := &papi.Mock{}
		defer client.AssertExpectations(t)

		expectGetCPCodes(client, "ctr_1", "grp_1", nil).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1").Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(3)

		expectGetCPCodeDetail(client, 0, "test cpcode", nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", fmt.Errorf("oops")).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
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

		expectGetCPCodes(client, "ctr_1", "grp_1", nil).Once()
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1").Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(3)

		expectGetCPCodeDetail(client, 0, "test cpcode", nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", nil).Once()
		expectGetCPCode(client, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil).Times(3)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
	})

	t.Run("error when no product and product_id provided", func(t *testing.T) {
		expectedErr := regexp.MustCompile("`product_id` must be specified for creation")

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/missing_product.tf"),
					ExpectError: expectedErr,
				},
			},
		})
	})
}
