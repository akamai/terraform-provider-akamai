package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestReportingGroupDataSource(t *testing.T) {
	t.Parallel()

	getResp := reportinggroups.GetReportingGroupResponse{
		ReportingGroupID:   12345,
		ReportingGroupName: "Test Reporting Group",
		AccessGroup: reportinggroups.AccessGroupModel{
			ContractID: "ctr_456",
			GroupID:    ptr.To(int64(456)),
		},
		Contracts: []reportinggroups.ContractModel{
			{
				ContractID: "ctr_123",
				CpCodes: []reportinggroups.CpCodeModel{
					{CpCodeID: 111, CpCodeName: "CP Code One"},
					{CpCodeID: 222, CpCodeName: "CP Code Two"},
				},
			},
		},
	}

	commonStateChecker := test.NewStateChecker("data.akamai_property_reportinggroups_group.test").
		CheckEqual("reporting_group_id", "12345").
		CheckEqual("reporting_group_name", "Test Reporting Group").
		CheckEqual("access_group.contract_id", "ctr_456").
		CheckEqual("access_group.group_id", "456").
		CheckEqual("contracts.#", "1").
		CheckEqual("contracts.0.contract_id", "ctr_123").
		CheckEqual("contracts.0.cp_codes.#", "2").
		CheckEqual("contracts.0.cp_codes.0.cp_code_id", "111").
		CheckEqual("contracts.0.cp_codes.0.cp_code_name", "CP Code One").
		CheckEqual("contracts.0.cp_codes.1.cp_code_id", "222").
		CheckEqual("contracts.0.cp_codes.1.cp_code_name", "CP Code Two")

	tests := map[string]struct {
		init  func(*reportinggroups.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - get reporting group": {
			init: func(m *reportinggroups.Mock) {
				mockGetReportingGroup(m, 12345, getResp, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroup/reporting_group.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - group_id is null": {
			init: func(m *reportinggroups.Mock) {
				respNilGroup := getResp
				respNilGroup.AccessGroup = reportinggroups.AccessGroupModel{
					ContractID: "ctr_456",
					GroupID:    nil,
				}
				mockGetReportingGroup(m, 12345, respNilGroup, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroup/reporting_group.tf"),
					Check: test.NewStateChecker("data.akamai_property_reportinggroups_group.test").
						CheckEqual("reporting_group_id", "12345").
						CheckEqual("access_group.contract_id", "ctr_456").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *reportinggroups.Mock) {
				mockGetReportingGroup(m, 12345, reportinggroups.GetReportingGroupResponse{}, fmt.Errorf("internal server error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroup/reporting_group.tf"),
					ExpectError: regexp.MustCompile("internal server error"),
				},
			},
		},
		"validation error - reporting_group_id missing": {
			init: func(_ *reportinggroups.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataReportingGroupsGroup/reporting_group_id_missing.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument(\n|.)+` + `The argument "reporting_group_id" is required, but no definition was found.`),
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

func mockGetReportingGroup(m *reportinggroups.Mock, reportingGroupID int64, resp reportinggroups.GetReportingGroupResponse, err error) *mock.Call {
	if err != nil {
		return m.On("GetReportingGroup", mock.Anything, reportinggroups.GetReportingGroupsRequest{
			ReportingGroupID: reportingGroupID,
		}).Return((*reportinggroups.GetReportingGroupResponse)(nil), err).Once()
	}

	return m.On("GetReportingGroup", mock.Anything, reportinggroups.GetReportingGroupsRequest{
		ReportingGroupID: reportingGroupID,
	}).Return(&resp, nil).Times(3)
}
