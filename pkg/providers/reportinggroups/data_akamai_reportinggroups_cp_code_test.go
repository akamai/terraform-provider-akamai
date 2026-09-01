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

func TestCPCodeDataSource(t *testing.T) {
	testDir := "testdata/TestDataCPCode/"
	t.Parallel()

	defaultCPCodeDetailResp := func(groupID *int64) *reportinggroups.GetCPCodeResponse {
		return &reportinggroups.GetCPCodeResponse{
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
					ProductID:   "Site_Accel::Site_Accel",
					ProductName: "Site Accelerator",
				},
			},
			AccessGroup: reportinggroups.AccessGroup{
				ContractID: "1-2ABCDE",
				GroupID:    groupID,
			},
		}
	}

	stateChecker := test.NewStateChecker("data.akamai_reportinggroups_cp_code.test").
		CheckEqual("cp_code_id", "12345").
		CheckEqual("name", "My CP Code").
		CheckEqual("purgeable", "true").
		CheckEqual("account_id", "act_1234").
		CheckEqual("default_time_zone", "GMT 0 (Greenwich Mean Time)").
		CheckEqual("override_time_zone.time_zone_id", "0").
		CheckEqual("override_time_zone.time_zone_value", "GMT 0 (Greenwich Mean Time)").
		CheckEqual("type", "Regular").
		CheckEqual("contracts.#", "1").
		CheckEqual("contracts.0.contract_id", "1-2ABCDE").
		CheckEqual("contracts.0.status", "ongoing").
		CheckEqual("products.#", "1").
		CheckEqual("products.0.product_id", "Site_Accel::Site_Accel").
		CheckEqual("products.0.product_name", "Site Accelerator").
		CheckEqual("access_group.contract_id", "1-2ABCDE").
		CheckEqual("access_group.group_id", "67890")

	tests := map[string]struct {
		init  func(*reportinggroups.Mock)
		steps []resource.TestStep
	}{
		"happy path - get CP code detail": {
			init: func(m *reportinggroups.Mock) {
				m.On("GetCPCode", testutils.MockContext, reportinggroups.GetCPCodeRequest{CPCodeID: 12345}).Return(defaultCPCodeDetailResp(ptr.To(int64(67890))), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_code.tf"),
					Check:  stateChecker.Build(),
				},
			},
		},
		"happy path - null group_id in access_group": {
			init: func(m *reportinggroups.Mock) {
				m.On("GetCPCode", testutils.MockContext, reportinggroups.GetCPCodeRequest{CPCodeID: 12345}).Return(defaultCPCodeDetailResp(nil), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"cp_code.tf"),
					Check:  stateChecker.CheckMissing("access_group.group_id").Build(),
				},
			},
		},
		"error - API error on get CP code detail": {
			init: func(m *reportinggroups.Mock) {
				m.On("GetCPCode", testutils.MockContext, reportinggroups.GetCPCodeRequest{CPCodeID: 12345}).Return(nil, fmt.Errorf("API failure")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"cp_code.tf"),
					ExpectError: regexp.MustCompile("Failed to retrieve CP code detail"),
				},
			},
		},
		"error - no cp_code_id provided": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"no_cp_code.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
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
