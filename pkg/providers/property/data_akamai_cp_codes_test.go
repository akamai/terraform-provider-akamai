package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDSCPCodes(t *testing.T) {
	workdir := "testdata/TestDSCPCodes"
	commonStateChecker := test.NewStateChecker("data.akamai_cp_codes.test").
		CheckEqual("contract_id", "11").
		CheckEqual("group_id", "22")

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - full datasource": {
			init: func(m *papi.Mock) {
				m.On("GetCPCodes", testutils.MockContext, papi.GetCPCodesRequest{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
				}).Return(&papi.GetCPCodesResponse{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
					AccountID:  "act_33",
					CPCodes: papi.CPCodeItems{
						Items: []papi.CPCode{
							{
								ID:          "cpc_code_1",
								Name:        "test_name",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_1"},
							},
							{
								ID:          "cpc_code_2",
								Name:        "test_name_2",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_2"},
							},
						},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/full.tf", workdir),
					Check: commonStateChecker.
						CheckEqual("account_id", "33").
						CheckEqual("cp_codes.#", "2").
						CheckEqual("cp_codes.0.cp_code_id", "code_1").
						CheckEqual("cp_codes.0.name", "test_name").
						CheckEqual("cp_codes.0.created_date", "2021-11-11T11:22:33Z").
						CheckEqual("cp_codes.0.product_ids.#", "1").
						CheckEqual("cp_codes.0.product_ids.0", "prd_1").
						CheckEqual("cp_codes.1.cp_code_id", "code_2").
						CheckEqual("cp_codes.1.name", "test_name_2").
						CheckEqual("cp_codes.1.created_date", "2021-11-11T11:22:33Z").
						CheckEqual("cp_codes.1.product_ids.#", "1").
						CheckEqual("cp_codes.1.product_ids.0", "prd_2").
						Build(),
				},
			},
		},
		"happy path - filter by name": {
			init: func(m *papi.Mock) {
				m.On("GetCPCodes", testutils.MockContext, papi.GetCPCodesRequest{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
				}).Return(&papi.GetCPCodesResponse{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
					AccountID:  "act_33",
					CPCodes: papi.CPCodeItems{
						Items: []papi.CPCode{
							{
								ID:          "cpc_code_1",
								Name:        "test_name",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_1"},
							},
							{
								ID:          "cpc_code_2",
								Name:        "test_name_2",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_2"},
							},
						},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/filter_by_name.tf", workdir),
					Check: commonStateChecker.
						CheckEqual("account_id", "33").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "code_1").
						CheckEqual("cp_codes.0.name", "test_name").
						CheckEqual("cp_codes.0.created_date", "2021-11-11T11:22:33Z").
						CheckEqual("cp_codes.0.product_ids.#", "1").
						CheckEqual("cp_codes.0.product_ids.0", "prd_1").
						Build(),
				},
			},
		},
		"happy path - filter by product": {
			init: func(m *papi.Mock) {
				m.On("GetCPCodes", testutils.MockContext, papi.GetCPCodesRequest{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
				}).Return(&papi.GetCPCodesResponse{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
					AccountID:  "act_33",
					CPCodes: papi.CPCodeItems{
						Items: []papi.CPCode{
							{
								ID:          "cpc_code_1",
								Name:        "test_name",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_1"},
							},
							{
								ID:          "cpc_code_2",
								Name:        "test_name_2",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_2"},
							},
						},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/filter_by_product_id.tf", workdir),
					Check: commonStateChecker.
						CheckEqual("account_id", "33").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "code_2").
						CheckEqual("cp_codes.0.name", "test_name_2").
						CheckEqual("cp_codes.0.created_date", "2021-11-11T11:22:33Z").
						CheckEqual("cp_codes.0.product_ids.#", "1").
						CheckEqual("cp_codes.0.product_ids.0", "prd_2").
						Build(),
				},
			},
		},
		"happy path - filter by product id and name": {
			init: func(m *papi.Mock) {
				m.On("GetCPCodes", testutils.MockContext, papi.GetCPCodesRequest{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
				}).Return(&papi.GetCPCodesResponse{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
					AccountID:  "act_33",
					CPCodes: papi.CPCodeItems{
						Items: []papi.CPCode{
							{
								ID:          "cpc_code_1",
								Name:        "test_name",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_1"},
							},
							{
								ID:          "cpc_code_1",
								Name:        "test_name",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_2"},
							},
							{
								ID:          "cpc_code_2",
								Name:        "test_name_2",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_1"},
							},
							{
								ID:          "cpc_code_2",
								Name:        "test_name_2",
								CreatedDate: "2021-11-11T11:22:33Z",
								ProductIDs:  []string{"prd_2"},
							},
						},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/filter_by_product_id_and_name.tf", workdir),
					Check: commonStateChecker.
						CheckEqual("account_id", "33").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "code_1").
						CheckEqual("cp_codes.0.name", "test_name").
						CheckEqual("cp_codes.0.created_date", "2021-11-11T11:22:33Z").
						CheckEqual("cp_codes.0.product_ids.#", "1").
						CheckEqual("cp_codes.0.product_ids.0", "prd_2").
						Build(),
				},
			},
		},
		"missing required argument contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/missing_contract.tf", workdir),
					ExpectError: regexp.MustCompile("The argument \"contract_id\" is required, but no definition was found"),
				},
			},
		},
		"missing required argument group_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/missing_group.tf", workdir),
					ExpectError: regexp.MustCompile("The argument \"group_id\" is required, but no definition was found"),
				},
			},
		},
		"error API response": {
			init: func(m *papi.Mock) {
				m.On("GetCPCodes", testutils.MockContext, papi.GetCPCodesRequest{
					GroupID:    "grp_22",
					ContractID: "ctr_11",
				}).Return(nil, fmt.Errorf("oops")).Times(1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/full.tf", workdir),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			hapiClient := &hapi.Mock{}
			if test.init != nil {
				test.init(client)
			}

			useClient(client, hapiClient, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
