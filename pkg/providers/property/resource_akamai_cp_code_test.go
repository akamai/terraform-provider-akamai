package property

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResCPCode(t *testing.T) {
	t.Parallel()
	expectGetCPCode := func(m *papi.Mock, contractID, groupID string, CPCodeID int, CPCodeName string, CPCodeProductIDs []string, err error, times int) *mock.Call {
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
		return m.On("GetCPCode", testutils.MockContext, req).Return(res, nil).Times(times)
	}

	expectGetCPCodes := func(m *papi.Mock, ContractID, GroupID string, CPCodes []papi.CPCode) *mock.Call {
		req := papi.GetCPCodesRequest{ContractID: ContractID, GroupID: GroupID}
		res := &papi.GetCPCodesResponse{
			ContractID: req.ContractID,
			GroupID:    req.GroupID,
			CPCodes:    papi.CPCodeItems{Items: CPCodes},
		}

		return m.On("GetCPCodes", testutils.MockContext, req).Return(res, nil).Once()
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

		return m.On("CreateCPCode", testutils.MockContext, req).Return(res, nil).Once()
	}

	expectUpdateCPCode := func(m *reportinggroups.Mock, CPCodeID int, name string, err error) *mock.Call {
		var res *reportinggroups.UpdateCPCodeResponse
		req := reportinggroups.UpdateCPCodeRequest{
			CPCodeID:         int64(CPCodeID),
			CPCodeName:       name,
			Purgeable:        ptr.To(false),
			OverrideTimeZone: &reportinggroups.CPCodeTimeZone{},
		}

		if err == nil {
			res = &reportinggroups.UpdateCPCodeResponse{
				CPCodeID:   req.CPCodeID,
				CPCodeName: req.CPCodeName,
			}

		}

		return m.On("UpdateCPCode", testutils.MockContext, req).Return(res, err).Once()
	}

	expectGetCPCodeDetail := func(m *reportinggroups.Mock, CPCodeID int, CPCodeName string, err error) *mock.Call {
		var res *reportinggroups.GetCPCodeResponse
		if err == nil {
			res = &reportinggroups.GetCPCodeResponse{
				CPCodeID:   int64(CPCodeID),
				CPCodeName: CPCodeName,
			}
		}
		return m.On("GetCPCode", testutils.MockContext, reportinggroups.GetCPCodeRequest{CPCodeID: int64(CPCodeID)}).Return(res, err).Once()
	}

	tests := map[string]struct {
		init                        func(*papi.Mock, *reportinggroups.Mock)
		steps                       []resource.TestStep
		updatePollInterval          time.Duration
		cpCodeResourceUpdateTimeout time.Duration
	}{
		"create new CP Code": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2", "prd_3"})
				expectGetCPCodes(p, "ctr_1", "grp_1", nil)
				expectCreateCPCode(p, "test cpcode", "prd_1", "ctr_1", "grp_1")
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 2)

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/create_new_cp_code.tf"),
				Check: test.NewStateChecker("akamai_cp_code.test").
					CheckEqual("id", "0").
					CheckEqual("name", "test cpcode").
					CheckEqual("group_id", "grp_1").
					CheckEqual("contract_id", "ctr_1").
					CheckEqual("product_id", "prd_1").
					CheckEqual("timeouts.#", "1").
					CheckEqual("timeouts.0.update", "1h").
					Build(),
			}},
		},
		"use existing CP Code with multiple products": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_test", []string{"prd_test", "prd_wrong", "another_wrong"})

				CPCodes := []papi.CPCode{
					{ID: "0", Name: "test cpcode", ProductIDs: []string{"prd_test", "prd_wrong", "another_wrong"}},
				}

				expectGetCPCodes(p, "ctr_test", "grp_test", CPCodes)
				// No mock behavior for create because we're using an existing CP code

				// Read and plan
				expectGetCPCode(p, "ctr_test", "grp_test", 0, "test cpcode", []string{"prd_test", "prd_wrong", "another_wrong"}, nil, 2)

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/use_existing_cp_code.tf"),
				Check: test.NewStateChecker("akamai_cp_code.test").
					CheckEqual("id", "0").
					CheckEqual("name", "test cpcode").
					CheckEqual("group_id", "grp_test").
					CheckEqual("contract_id", "ctr_test").
					CheckEqual("product_id", "prd_test").
					CheckEqual("timeouts.#", "0").
					Build(),
			}},
		},
		"use existing CP Code": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_test", []string{"prd_test", "prd_wrong", "another_wrong"})

				CPCodes := []papi.CPCode{
					{ID: "0", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
					{ID: "cpc_1", Name: "test cpcode", ProductIDs: []string{"prd_test"}},
				}

				expectGetCPCodes(p, "ctr_test", "grp_test", CPCodes)
				// No mock behavior for create because we're using an existing CP code

				// Read and plan
				expectGetCPCode(p, "ctr_test", "grp_test", 1, "test cpcode", []string{"prd_test"}, nil, 2)

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/use_existing_cp_code.tf"),
				Check: test.NewStateChecker("akamai_cp_code.test").
					CheckEqual("id", "1").
					CheckEqual("name", "test cpcode").
					CheckEqual("group_id", "grp_test").
					CheckEqual("contract_id", "ctr_test").
					CheckEqual("product_id", "prd_test").
					Build(),
			}},
		},
		"product missing from CP Code": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_test", []string{"prd_1", "prd_2", "prd_3"})

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/use_existing_cp_code.tf"),
				ExpectError: regexp.MustCompile("`product_id` `prd_test` does not exist under contract `ctr_test`, you need to provide a valid `product_id`"),
			}},
		},
		"change name": {
			init: func(p *papi.Mock, rg *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2", "prd_3"})

				expectGetCPCodes(p, "ctr_1", "grp_1", nil)
				expectCreateCPCode(p, "test cpcode", "prd_1", "ctr_1", "grp_1")
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 3)

				expectGetCPCodeDetail(rg, 0, "test cpcode", nil)
				expectUpdateCPCode(rg, 0, "renamed cpcode", nil)
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 1)
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "renamed cpcode", []string{"prd_1"}, nil, 3)

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "renamed cpcode").
						Build(),
				},
			},
		},
		"create CP Code but existing CPCode has no product IDs returns error": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1"})
				CPCodes := []papi.CPCode{
					{ID: "0", Name: "test cpcode", ProductIDs: []string{}},
				}
				expectGetCPCodes(p, "ctr_1", "grp_1", CPCodes)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/create_new_cp_code.tf"),
					ExpectError: regexp.MustCompile("the CP code named `test cpcode` already exists, but does not have a PAPI-supported product ID, so it cannot be managed by Terraform"),
				},
			},
		},
		"import existing cp code": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_test", "prd_Web_Accel"})

				CPCodes := []papi.CPCode{{ID: "0", Name: "test cpcode", ProductIDs: []string{"prd_Web_Accel"}}}
				expectGetCPCodes(p, "ctr_1", "grp_2", CPCodes)
				expectGetCPCode(p, "ctr_1", "grp_2", 0, "test cpcode", []string{"prd_Web_Accel"}, nil, 4)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
				},
				{
					ImportState:   true,
					ImportStateId: "0,1,2",
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
		},
		"invalid import ID passed": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
					ImportState:   true,
					ImportStateId: "123",
					ResourceName:  "akamai_cp_code.test",
					ExpectError:   regexp.MustCompile("comma-separated list of CP code ID, contract ID and group ID has to be supplied in import"),
				},
			},
		},
		"import cp code with no product IDs returns error": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				// Mock GetCPCode to return a CPCode with no ProductIDs
				req := papi.GetCPCodeRequest{CPCodeID: "123", ContractID: "ctr_1", GroupID: "grp_1"}
				p.On("GetCPCode", mock.Anything, req).Return(&papi.GetCPCodesResponse{
					CPCode: papi.CPCode{
						ID:         "123",
						Name:       "test cpcode",
						ProductIDs: []string{},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
					ImportState:   true,
					ImportStateId: "123,ctr_1,grp_1",
					ResourceName:  "akamai_cp_code.test",
					ExpectError:   regexp.MustCompile("the CP code named `test cpcode` already exists, but does not have a PAPI-supported product ID, so it cannot be managed by Terraform"),
				},
			},
		},
		"empty CP code ID passed": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCPCode/import_cp_code.tf"),
					ImportState:   true,
					ImportStateId: ",ctr_1-1NC95D,grp_194665",
					ResourceName:  "akamai_cp_code.test",
					ExpectError:   regexp.MustCompile("CP Code is a mandatory parameter"),
				},
			},
		},
		"immutable attributes updated": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2", "prd_3"})

				expectGetCPCodes(p, "ctr_1", "grp_1", nil)
				expectCreateCPCode(p, "test cpcode", "prd_1", "ctr_1", "grp_1")
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 5)

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						CheckEqual("group_id", "grp_1").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("product_id", "prd_1").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_immutable.tf"),
					ExpectError: regexp.MustCompile(`cp code attribute 'contract_id' cannot be changed after creation \(immutable\)`),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						CheckEqual("group_id", "grp_1").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("product_id", "prd_1").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_immutable.tf"),
					ExpectError: regexp.MustCompile(`cp code attribute 'product_id' cannot be changed after creation \(immutable\)`),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						CheckEqual("group_id", "grp_1").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("product_id", "prd_1").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_immutable.tf"),
					ExpectError: regexp.MustCompile(`cp code attribute 'group_id' cannot be changed after creation \(immutable\)`),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						CheckEqual("group_id", "grp_1").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("product_id", "prd_1").
						Build(),
				},
			},
		},
		"error fetching cpCode details": {
			init: func(p *papi.Mock, rg *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2", "prd_3"})

				expectGetCPCodes(p, "ctr_1", "grp_1", nil)
				expectCreateCPCode(p, "test cpcode", "prd_1", "ctr_1", "grp_1")
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 3)

				expectGetCPCodeDetail(rg, 0, "test cpcode", fmt.Errorf("oops"))

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error updating cpCode": {
			init: func(p *papi.Mock, rg *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2", "prd_3"})

				expectGetCPCodes(p, "ctr_1", "grp_1", nil)
				expectCreateCPCode(p, "test cpcode", "prd_1", "ctr_1", "grp_1")
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 3)

				expectGetCPCodeDetail(rg, 0, "test cpcode", nil)
				expectUpdateCPCode(rg, 0, "renamed cpcode", fmt.Errorf("oops"))

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"timeout waiting for update": {
			init: func(p *papi.Mock, rg *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2", "prd_3"})

				expectGetCPCodes(p, "ctr_1", "grp_1", nil)
				expectCreateCPCode(p, "test cpcode", "prd_1", "ctr_1", "grp_1")
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 3)

				expectGetCPCodeDetail(rg, 0, "test cpcode", nil)
				expectUpdateCPCode(rg, 0, "renamed cpcode", nil)
				expectGetCPCode(p, "ctr_1", "grp_1", 0, "test cpcode", []string{"prd_1"}, nil, 3)

				// No mock behavior for delete because there is no delete operation for CP Codes
			},
			updatePollInterval:          time.Millisecond * 40,
			cpCodeResourceUpdateTimeout: time.Millisecond * 60,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step0.tf"),
					Check: test.NewStateChecker("akamai_cp_code.test").
						CheckEqual("id", "0").
						CheckEqual("name", "test cpcode").
						Build(),
				},
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCPCode/change_name_step1.tf"),
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"error when no product and product_id provided": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCPCode/missing_product.tf"),
					ExpectError: regexp.MustCompile("`product_id` must be specified for creation"),
				},
			},
		},
		"use existing CP Code emits warning for associated product IDs": {
			init: func(p *papi.Mock, _ *reportinggroups.Mock) {
				expectGetProducts(p, "ctr_1", []string{"prd_1", "prd_2"})

				CPCodes := []papi.CPCode{
					{ID: "12", Name: "test cpcode", ProductIDs: []string{"prd_2"}},
				}

				expectGetCPCodes(p, "ctr_1", "grp_1", CPCodes)
				expectGetCPCode(p, "ctr_1", "grp_1", 12, "test cpcode", []string{"prd_2"}, nil, 2)
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResCPCode/existing_cp_code.tf"),
				// Note: The warning diagnostic is not directly assertable here, but this test exercises the code path.
				Check: test.NewStateChecker("akamai_cp_code.test").
					CheckEqual("id", "12").
					CheckEqual("name", "test cpcode").
					CheckEqual("group_id", "grp_1").
					CheckEqual("contract_id", "ctr_1").
					CheckEqual("product_id", "prd_2").
					Build(),
			}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()

			if tc.init != nil {
				tc.init(client.PAPI, client.ReportingGroups)
			}

			// redefining times to accelerate tests where possible
			config := defaultSubproviderConfig()
			config.cpCode.updatePollMinimum = time.Millisecond * 1
			config.cpCode.updatePollInterval = config.cpCode.updatePollMinimum

			if tc.updatePollInterval != 0 {
				config.cpCode.updatePollInterval = tc.updatePollInterval
			}
			if tc.cpCodeResourceUpdateTimeout != 0 {
				config.cpCode.cpCodeResourceUpdateTimeout = tc.cpCodeResourceUpdateTimeout
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(config)),
				Steps:                    tc.steps,
			})

			client.PAPI.AssertExpectations(t)
			client.ReportingGroups.AssertExpectations(t)
		})
	}
}

func expectGetProducts(m *papi.Mock, ContractID string, ProductIDs []string) *mock.Call {
	req := papi.GetProductsRequest{ContractID: ContractID}
	products := make([]papi.ProductItem, 0, len(ProductIDs))
	for _, pid := range ProductIDs {
		products = append(products, papi.ProductItem{
			ProductID:   pid,
			ProductName: "Product " + pid,
		})
	}
	res := &papi.GetProductsResponse{
		ContractID: req.ContractID,
		Products:   papi.ProductsItems{Items: products},
	}

	return m.On("GetProducts", testutils.MockContext, req).Return(res, nil).Once()
}
