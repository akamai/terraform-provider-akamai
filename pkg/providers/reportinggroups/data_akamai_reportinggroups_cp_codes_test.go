package reportinggroups

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCPCodesDataSource(t *testing.T) {
	testDir := "testdata/TestDataCPCodes/"
	t.Parallel()

	cpCode1 := reportinggroups.CPCodeDetails{
		CPCodeID:        12345,
		CPCodeName:      "My CP Code",
		Purgeable:       true,
		AccountID:       "act_1234",
		DefaultTimeZone: "GMT 0 (Greenwich Mean Time)",
		OverrideTimeZone: reportinggroups.CPCodeTimeZone{
			TimeZoneID:    "0",
			TimeZoneValue: "GMT 0 (Greenwich Mean Time)",
		},
		Type: "Regular",
		Contracts: []reportinggroups.CPCodeContract{
			{
				ContractID: "1-2ABCDE",
				Status:     "ongoing",
			},
		},
		Products: []reportinggroups.Product{
			{
				ProductID:   "Web_Exp::Ion_Na",
				ProductName: "Ion Standard",
			},
		},
		AccessGroup: reportinggroups.AccessGroup{
			ContractID: "1-2ABCDE",
			GroupID:    ptr.To(int64(67890)),
		},
	}

	cpCode2 := reportinggroups.CPCodeDetails{
		CPCodeID:        23456,
		CPCodeName:      "Other CP Code",
		Purgeable:       false,
		AccountID:       "act_5678",
		DefaultTimeZone: "GMT+1 (Central European Time)",
		OverrideTimeZone: reportinggroups.CPCodeTimeZone{
			TimeZoneID:    "1",
			TimeZoneValue: "GMT+1 (Central European Time)",
		},
		Type: "Regular",
		Contracts: []reportinggroups.CPCodeContract{
			{
				ContractID: "1-3FGHIJ",
				Status:     "ongoing",
			},
		},
		Products: []reportinggroups.Product{
			{
				ProductID:   "Site_Accel::Site_Accel",
				ProductName: "Site Accelerator",
			},
		},
		AccessGroup: reportinggroups.AccessGroup{
			ContractID: "1-3FGHIJ",
			GroupID:    ptr.To(int64(78901)),
		},
	}

	defaultListCPCodesResp := &reportinggroups.ListCPCodesResponse{
		CPCodes: []reportinggroups.CPCodeDetails{cpCode1, cpCode2},
	}

	stateChecker := test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
		CheckEqual("cp_codes.#", "2").
		CheckEqual("cp_codes.0.cp_code_id", "12345").
		CheckEqual("cp_codes.0.name", "My CP Code").
		CheckEqual("cp_codes.0.purgeable", "true").
		CheckEqual("cp_codes.0.account_id", "act_1234").
		CheckEqual("cp_codes.0.default_time_zone", "GMT 0 (Greenwich Mean Time)").
		CheckEqual("cp_codes.0.override_time_zone.time_zone_id", "0").
		CheckEqual("cp_codes.0.override_time_zone.time_zone_value", "GMT 0 (Greenwich Mean Time)").
		CheckEqual("cp_codes.0.type", "Regular").
		CheckEqual("cp_codes.0.contracts.#", "1").
		CheckEqual("cp_codes.0.contracts.0.contract_id", "1-2ABCDE").
		CheckEqual("cp_codes.0.contracts.0.status", "ongoing").
		CheckEqual("cp_codes.0.products.#", "1").
		CheckEqual("cp_codes.0.products.0.product_id", "Web_Exp::Ion_Na").
		CheckEqual("cp_codes.0.products.0.product_name", "Ion Standard").
		CheckEqual("cp_codes.0.access_group.contract_id", "1-2ABCDE").
		CheckEqual("cp_codes.0.access_group.group_id", "67890").
		CheckEqual("cp_codes.1.cp_code_id", "23456").
		CheckEqual("cp_codes.1.name", "Other CP Code").
		CheckEqual("cp_codes.1.purgeable", "false").
		CheckEqual("cp_codes.1.account_id", "act_5678").
		CheckEqual("cp_codes.1.default_time_zone", "GMT+1 (Central European Time)").
		CheckEqual("cp_codes.1.override_time_zone.time_zone_id", "1").
		CheckEqual("cp_codes.1.override_time_zone.time_zone_value", "GMT+1 (Central European Time)").
		CheckEqual("cp_codes.1.type", "Regular").
		CheckEqual("cp_codes.1.contracts.#", "1").
		CheckEqual("cp_codes.1.contracts.0.contract_id", "1-3FGHIJ").
		CheckEqual("cp_codes.1.contracts.0.status", "ongoing").
		CheckEqual("cp_codes.1.products.#", "1").
		CheckEqual("cp_codes.1.products.0.product_id", "Site_Accel::Site_Accel").
		CheckEqual("cp_codes.1.products.0.product_name", "Site Accelerator").
		CheckEqual("cp_codes.1.access_group.contract_id", "1-3FGHIJ").
		CheckEqual("cp_codes.1.access_group.group_id", "78901")

	tests := map[string]struct {
		init  func(*reportinggroups.Mock)
		steps []resource.TestStep
	}{
		"happy path - list CP codes without filters": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{}).
					Return(defaultListCPCodesResp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes.tf"),
					Check:  stateChecker.Build(),
				},
			},
		},
		"happy path - list CP codes with all filters": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{
					ContractID: "1-2ABCDE",
					GroupID:    "grp_67890",
					ProductID:  "Web_Exp::Ion_Na",
					CPCodeName: "My CP Code",
				}).Return(&reportinggroups.ListCPCodesResponse{
					CPCodes: []reportinggroups.CPCodeDetails{cpCode1},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes_with_filters.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "12345").
						CheckEqual("contract_id", "1-2ABCDE").
						CheckEqual("group_id", "grp_67890").
						CheckEqual("product_id", "Web_Exp::Ion_Na").
						CheckEqual("cp_code_name", "My CP Code").
						Build(),
				},
			},
		},
		"happy path - null group_id in access_group": {
			init: func(m *reportinggroups.Mock) {
				nullGroupResp := &reportinggroups.ListCPCodesResponse{
					CPCodes: []reportinggroups.CPCodeDetails{cpCode1, cpCode2},
				}
				nullGroupResp.CPCodes[0].AccessGroup.GroupID = nil
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{}).
					Return(nullGroupResp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes.tf"),
					Check:  stateChecker.CheckMissing("cp_codes.0.access_group.group_id").Build(),
				},
			},
		},
		"happy path - empty list": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{}).
					Return(&reportinggroups.ListCPCodesResponse{CPCodes: []reportinggroups.CPCodeDetails{}}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
						CheckEqual("cp_codes.#", "0").
						Build(),
				},
			},
		},
		"happy path - filter by contract_id only": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{
					ContractID: "1-2ABCDE",
				}).Return(&reportinggroups.ListCPCodesResponse{
					CPCodes: []reportinggroups.CPCodeDetails{cpCode1},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes_contract_id_only.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "12345").
						CheckEqual("contract_id", "1-2ABCDE").
						CheckMissing("group_id").
						CheckMissing("product_id").
						CheckMissing("cp_code_name").
						Build(),
				},
			},
		},
		"happy path - filter by group_id only": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{
					GroupID: "grp_67890",
				}).Return(&reportinggroups.ListCPCodesResponse{
					CPCodes: []reportinggroups.CPCodeDetails{cpCode1},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes_group_id_only.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "12345").
						CheckEqual("group_id", "grp_67890").
						CheckMissing("contract_id").
						CheckMissing("product_id").
						CheckMissing("cp_code_name").
						Build(),
				},
			},
		},
		"happy path - filter by product_id only": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{
					ProductID: "Site_Accel::Site_Accel",
				}).Return(&reportinggroups.ListCPCodesResponse{
					CPCodes: []reportinggroups.CPCodeDetails{cpCode2},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes_product_id_only.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "23456").
						CheckEqual("product_id", "Site_Accel::Site_Accel").
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("cp_code_name").
						Build(),
				},
			},
		},
		"happy path - filter by cp_code_name only": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{
					CPCodeName: "Other CP Code",
				}).Return(&reportinggroups.ListCPCodesResponse{
					CPCodes: []reportinggroups.CPCodeDetails{cpCode2},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_codes_cp_code_name_only.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_cp_codes.test").
						CheckEqual("cp_codes.#", "1").
						CheckEqual("cp_codes.0.cp_code_id", "23456").
						CheckEqual("cp_code_name", "Other CP Code").
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("product_id").
						Build(),
				},
			},
		},
		"error - API error on list CP codes": {
			init: func(m *reportinggroups.Mock) {
				m.On("ListCPCodes", testutils.MockContext, reportinggroups.ListCPCodesRequest{}).
					Return(nil, fmt.Errorf("API failure")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"cp_codes.tf"),
					ExpectError: regexp.MustCompile("Failed to retrieve CP codes detail"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.ReportingGroups)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps:                    tc.steps,
			})

			client.ReportingGroups.AssertExpectations(t)
		})
	}
}
