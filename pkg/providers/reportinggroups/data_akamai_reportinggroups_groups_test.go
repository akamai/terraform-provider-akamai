package reportinggroups

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestReportingGroupsDataSource(t *testing.T) {
	t.Parallel()

	listResp := reportinggroups.ListReportingGroupsResponse{
		Groups: []reportinggroups.ReportingGroup{
			{
				ReportingGroupID:   12345,
				ReportingGroupName: "First Reporting Group",
				AccessGroup: reportinggroups.AccessGroup{
					ContractID: "ctr_123",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []reportinggroups.Contract{
					{
						ContractID: "ctr_123",
						CPCodes: []reportinggroups.CPCode{
							{CPCodeID: 111, CPCodeName: "CP Code One"},
							{CPCodeID: 222, CPCodeName: "CP Code Two"},
						},
					},
				},
			},
			{
				ReportingGroupID:   12347,
				ReportingGroupName: "Third Reporting Group",
				AccessGroup: reportinggroups.AccessGroup{
					ContractID: "ctr_456",
					GroupID:    ptr.To(int64(456)),
				},
				Contracts: []reportinggroups.Contract{
					{
						ContractID: "ctr_456",
						CPCodes: []reportinggroups.CPCode{
							{CPCodeID: 333, CPCodeName: "CP Code Three"},
						},
					},
				},
			},
			{
				ReportingGroupID:   12346,
				ReportingGroupName: "Second Reporting Group",
				AccessGroup: reportinggroups.AccessGroup{
					ContractID: "ctr_456",
					GroupID:    nil,
				},
				Contracts: []reportinggroups.Contract{
					{
						ContractID: "ctr_456",
						CPCodes: []reportinggroups.CPCode{
							{CPCodeID: 333, CPCodeName: "CP Code Three"},
							{CPCodeID: 222, CPCodeName: "CP Code Two"},
						},
					},
				},
			},
		},
	}

	firstGroupWithOneCPCode := listResp.Groups[0]
	firstGroupWithOneCPCode.Contracts = []reportinggroups.Contract{
		{
			ContractID: "ctr_123",
			CPCodes: []reportinggroups.CPCode{
				{
					CPCodeID:   111,
					CPCodeName: "CP Code One",
				},
			},
		},
	}

	firstReportingGroup := test.AttributeBatch{
		"reporting_group_id":               "12345",
		"reporting_group_name":             "First Reporting Group",
		"access_group.contract_id":         "ctr_123",
		"access_group.group_id":            "456",
		"contract.contract_id":             "ctr_123",
		"contract.cp_codes.#":              "2",
		"contract.cp_codes.0.cp_code_id":   "111",
		"contract.cp_codes.0.cp_code_name": "CP Code One",
		"contract.cp_codes.1.cp_code_id":   "222",
		"contract.cp_codes.1.cp_code_name": "CP Code Two",
	}

	secondReportingGroup := test.AttributeBatch{
		"reporting_group_id":               "12346",
		"reporting_group_name":             "Second Reporting Group",
		"access_group.contract_id":         "ctr_456",
		"contract.contract_id":             "ctr_456",
		"contract.cp_codes.#":              "2",
		"contract.cp_codes.0.cp_code_id":   "222",
		"contract.cp_codes.0.cp_code_name": "CP Code Two",
		"contract.cp_codes.1.cp_code_id":   "333",
		"contract.cp_codes.1.cp_code_name": "CP Code Three",
	}

	thirdReportingGroup := test.AttributeBatch{
		"reporting_group_id":               "12347",
		"reporting_group_name":             "Third Reporting Group",
		"access_group.contract_id":         "ctr_456",
		"access_group.group_id":            "456",
		"contract.contract_id":             "ctr_456",
		"contract.cp_codes.#":              "1",
		"contract.cp_codes.0.cp_code_id":   "333",
		"contract.cp_codes.0.cp_code_name": "CP Code Three",
	}

	commonStateChecker := test.NewStateChecker("data.akamai_reportinggroups_groups.test").
		CheckEqual("groups.#", "3").
		CheckEqualBatch("groups.0.", firstReportingGroup).
		CheckEqualBatch("groups.1.", secondReportingGroup).
		CheckMissing("groups.1.access_group.group_id").
		CheckEqualBatch("groups.2.", thirdReportingGroup)

	firstGroupWithOneCPCodeChecker := test.NewStateChecker("data.akamai_reportinggroups_groups.test").
		CheckEqual("cp_code_id", "111").
		CheckEqual("groups.#", "1").
		CheckEqualBatch("groups.0.", firstReportingGroup).
		CheckEqual("groups.0.contract.cp_codes.#", "1").
		CheckMissing("groups.0.contract.cp_codes.1.cp_code_id").
		CheckMissing("groups.0.contract.cp_codes.1.cp_code_name")

	tests := map[string]struct {
		init  func(*reportinggroups.Mock)
		steps []resource.TestStep
	}{
		"happy path - list reporting groups no filters": {
			init: func(m *reportinggroups.Mock) {
				mockListReportingGroups(m, reportinggroups.ListReportingGroupsRequest{}, listResp, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroups/reporting_groups.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - list reporting groups filtered by cp_code_id": {
			init: func(m *reportinggroups.Mock) {
				mockListReportingGroups(m, reportinggroups.ListReportingGroupsRequest{
					CPCodeID: "111",
				}, reportinggroups.ListReportingGroupsResponse{
					Groups: []reportinggroups.ReportingGroup{firstGroupWithOneCPCode},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroups/reporting_groups_with_cp_code_id.tf"),
					Check: firstGroupWithOneCPCodeChecker.
						Build(),
				},
			},
		},
		"happy path - list reporting groups filtered by contract_id": {
			init: func(m *reportinggroups.Mock) {
				mockListReportingGroups(m, reportinggroups.ListReportingGroupsRequest{
					ContractID: "ctr_456",
				}, reportinggroups.ListReportingGroupsResponse{
					Groups: listResp.Groups[1:],
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroups/reporting_groups_with_contract_id.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_groups.test").
						CheckEqual("groups.#", "2").
						CheckEqualBatch("groups.0.", secondReportingGroup).
						CheckMissing("groups.0.access_group.group_id").
						CheckEqualBatch("groups.1.", thirdReportingGroup).
						Build(),
				},
			},
		},
		"happy path - list reporting groups with all filters": {
			init: func(m *reportinggroups.Mock) {
				mockListReportingGroups(m, reportinggroups.ListReportingGroupsRequest{
					ContractID:         "ctr_123",
					GroupID:            "456",
					ReportingGroupName: "First Reporting Group",
					CPCodeID:           "111",
				}, reportinggroups.ListReportingGroupsResponse{
					Groups: []reportinggroups.ReportingGroup{firstGroupWithOneCPCode},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroups/reporting_groups_with_all_filters.tf"),
					Check: firstGroupWithOneCPCodeChecker.
						CheckEqual("contract_id", "ctr_123").
						CheckEqual("group_id", "456").
						CheckEqual("reporting_group_name", "First Reporting Group").
						Build(),
				},
			},
		},
		"happy path - empty list": {
			init: func(m *reportinggroups.Mock) {
				mockListReportingGroups(m, reportinggroups.ListReportingGroupsRequest{}, reportinggroups.ListReportingGroupsResponse{
					Groups: []reportinggroups.ReportingGroup{},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroups/reporting_groups.tf"),
					Check: test.NewStateChecker("data.akamai_reportinggroups_groups.test").
						CheckEqual("groups.#", "0").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *reportinggroups.Mock) {
				mockListReportingGroups(m, reportinggroups.ListReportingGroupsRequest{}, reportinggroups.ListReportingGroupsResponse{}, fmt.Errorf("internal server error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroups/reporting_groups.tf"),
					ExpectError: regexp.MustCompile("internal server error"),
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

func mockListReportingGroups(m *reportinggroups.Mock, req reportinggroups.ListReportingGroupsRequest, resp reportinggroups.ListReportingGroupsResponse, err error) *mock.Call {
	if err != nil {
		return m.On("ListReportingGroups", testutils.MockContext, req).
			Return((*reportinggroups.ListReportingGroupsResponse)(nil), err).Once()
	}

	return m.On("ListReportingGroups", testutils.MockContext, req).
		Return(&resp, nil).Times(3)
}
