package reportinggroups

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

// reportingGroupTestData bundles all input and expected output fields used to
// configure mock expectations and state checkers for a single test scenario.
type reportingGroupTestData struct {
	// request fields
	reportingGroupName  string
	accessGroupContract string
	contractID          string
	groupID             string
	cpCodes             []string
	// response fields
	reportingGroupID int64
	cpCodeNames      []string
}

var (
	minReportingGroup = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group",
			accessGroupContract: "test_contract",
			groupID:             "12345",
			contractID:          "test_contract_2",
			cpCodes:             []string{"111111"},
			cpCodeNames:         []string{"test-cp_code"},
			reportingGroupID:    999,
		}
	}

	noGroupIDReportingGroup = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group",
			accessGroupContract: "test_contract",
			contractID:          "test_contract_2",
			cpCodes:             []string{"111111"},
			cpCodeNames:         []string{"test-cp_code"},
			reportingGroupID:    999,
		}
	}

	reportingGroupWithPrefixes = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group",
			accessGroupContract: "ctr_test_contract",
			groupID:             "grp_12345",
			contractID:          "ctr_test_contract_2",
			cpCodes:             []string{"cpc_111111"},
			cpCodeNames:         []string{"test-cp_code"},
			reportingGroupID:    999,
		}
	}

	updatedNameReportingGroupWithPrefixes = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group-updated",
			accessGroupContract: "ctr_test_contract",
			groupID:             "grp_12345",
			contractID:          "ctr_test_contract_2",
			cpCodes:             []string{"cpc_111111"},
			cpCodeNames:         []string{"test-cp_code"},
			reportingGroupID:    999,
		}
	}

	updatedNameReportingGroup = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group-updated",
			accessGroupContract: "test_contract",
			groupID:             "12345",
			contractID:          "test_contract_2",
			cpCodes:             []string{"111111"},
			cpCodeNames:         []string{"test-cp_code"},
			reportingGroupID:    999,
		}
	}

	updatedCPCodesReportingGroup = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group",
			accessGroupContract: "test_contract",
			groupID:             "12345",
			contractID:          "test_contract_2",
			cpCodes:             []string{"111111", "333333"},
			cpCodeNames:         []string{"test-cp_code", "test-cp_code-new"},
			reportingGroupID:    999,
		}
	}

	multipleCPCodesReportingGroup = func() reportingGroupTestData {
		return reportingGroupTestData{
			reportingGroupName:  "test-reporting-group",
			accessGroupContract: "test_contract",
			groupID:             "12345",
			contractID:          "test_contract_2",
			cpCodes:             []string{"111111", "222222"},
			cpCodeNames:         []string{"test-cp_code", "test-cp_code2"},
			reportingGroupID:    999,
		}
	}
)

func TestReportingGroupsResource(t *testing.T) {
	t.Parallel()

	minChecker := test.NewStateChecker("akamai_reportinggroups_group.test").
		CheckEqual("reporting_group_id", "999").
		CheckEqual("reporting_group_name", "test-reporting-group").
		CheckEqual("access_group.contract_id", "test_contract").
		CheckEqual("access_group.group_id", "12345").
		CheckEqual("contract.contract_id", "test_contract_2").
		CheckEqual("contract.cp_codes.#", "1").
		CheckEqual("contract.cp_codes.0.cp_code_id", "111111").
		CheckEqual("contract.cp_codes.0.cp_code_name", "test-cp_code")

	importChecker := test.NewImportChecker().
		CheckEqual("reporting_group_id", "999").
		CheckEqual("reporting_group_name", "test-reporting-group").
		CheckEqual("access_group.contract_id", "test_contract").
		CheckEqual("contract.contract_id", "test_contract_2").
		CheckEqual("contract.cp_codes.#", "1").
		CheckEqual("contract.cp_codes.0.cp_code_id", "111111").
		CheckEqual("contract.cp_codes.0.cp_code_name", "test-cp_code")

	tests := map[string]struct {
		init           func(*reportinggroups.Mock, reportingGroupTestData, reportingGroupTestData)
		createMockData reportingGroupTestData
		updateMockData reportingGroupTestData
		steps          []resource.TestStep
	}{
		// ── Create ──────────────────────────────────────────────────────────────
		"happy path - create": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
			},
		},
		"happy path - create with multiple CP codes": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: multipleCPCodesReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/multiple_cp_codes.tf"),
					Check: minChecker.
						CheckEqual("contract.cp_codes.#", "2").
						CheckEqual("contract.cp_codes.1.cp_code_id", "222222").
						CheckEqual("contract.cp_codes.1.cp_code_name", "test-cp_code2").
						Build(),
				},
			},
		},
		"happy path - create with ctr_ and grp_ prefixes": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// The resource strips ctr_ and grp_ prefixes before sending to API
				mockCreateReportingGroup(m, createData)
				// After create, read returns prefixed values to match state
				mockGetRepGroup(m, createData)
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: reportingGroupWithPrefixes(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/with_prefixes.tf"),
					Check: minChecker.
						CheckEqual("access_group.group_id", "grp_12345").
						CheckEqual("access_group.contract_id", "ctr_test_contract").
						CheckEqual("contract.contract_id", "ctr_test_contract_2").
						CheckEqual("contract.cp_codes.0.cp_code_id", "cpc_111111").
						Build(),
				},
			},
		},

		// ── Update ──────────────────────────────────────────────────────────────
		"happy path - create then update reporting group name": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, updateData reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read before update
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Update
				mockUpdateReportingGroup(m, createData, updateData)
				// Read after update
				mockGetRepGroup(m, updateData)
				// Delete
				mockDeleteReportingGroup(m, updateData)
			},
			createMockData: minReportingGroup(),
			updateMockData: updatedNameReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/name.tf"),
					Check: minChecker.
						CheckEqual("reporting_group_name", "test-reporting-group-updated").
						Build(),
				},
			},
		},
		"happy path - create then update CP codes": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, updateData reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read before update
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Update
				mockUpdateReportingGroup(m, createData, updateData)
				// Read after update
				mockGetRepGroup(m, updateData)
				// Delete
				mockDeleteReportingGroup(m, updateData)
			},
			createMockData: minReportingGroup(),
			updateMockData: updatedCPCodesReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/add_cp_code.tf"),
					Check: minChecker.
						CheckEqual("contract.cp_codes.#", "2").
						CheckEqual("contract.cp_codes.1.cp_code_id", "333333").
						CheckEqual("contract.cp_codes.1.cp_code_name", "test-cp_code-new").
						Build(),
				},
			},
		},

		"happy path - create then update name with prefixes (prefixes stripped for API)": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, updateData reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read before update
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Update - must strip ctr_ prefix before sending to API
				mockUpdateReportingGroup(m, createData, updateData)
				// Read after update
				mockGetRepGroup(m, updateData)
				// Delete
				mockDeleteReportingGroup(m, updateData)
			},
			createMockData: reportingGroupWithPrefixes(),
			updateMockData: updatedNameReportingGroupWithPrefixes(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/with_prefixes.tf"),
					Check: minChecker.
						CheckEqual("access_group.group_id", "grp_12345").
						CheckEqual("access_group.contract_id", "ctr_test_contract").
						CheckEqual("contract.contract_id", "ctr_test_contract_2").
						CheckEqual("contract.cp_codes.0.cp_code_id", "cpc_111111").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/name_with_prefixes.tf"),
					Check: minChecker.
						CheckEqual("reporting_group_name", "test-reporting-group-updated").
						CheckEqual("access_group.group_id", "grp_12345").
						CheckEqual("access_group.contract_id", "ctr_test_contract").
						CheckEqual("contract.contract_id", "ctr_test_contract_2").
						CheckEqual("contract.cp_codes.0.cp_code_id", "cpc_111111").
						Build(),
				},
			},
		},

		"happy path - create with multiple CP codes then remove one": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, updateData reportingGroupTestData) {
				// Create with multiple CP codes
				mockCreateReportingGroup(m, createData)
				// Read before update
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Update (remove a CP code)
				mockUpdateReportingGroup(m, createData, updateData)
				// Read after update
				mockGetRepGroup(m, updateData)
				// Delete
				mockDeleteReportingGroup(m, updateData)
			},
			createMockData: multipleCPCodesReportingGroup(),
			updateMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/multiple_cp_codes.tf"),
					Check: minChecker.
						CheckEqual("contract.cp_codes.#", "2").
						CheckEqual("contract.cp_codes.1.cp_code_id", "222222").
						CheckEqual("contract.cp_codes.1.cp_code_name", "test-cp_code2").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check: minChecker.
						Build(),
				},
			},
		},
		"happy path - adding ctr_ prefix to contract.contract_id results in empty plan": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read for plan refresh (PlanOnly step)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					// Adding ctr_ prefix to contract.contract_id is a no-op: the
					// IgnorePrefixType semantic equals considers it unchanged.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/contract_id_ctr_prefix.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: false,
				},
			},
		},
		"happy path - adding ctr_ prefix to access_group.contract_id results in empty plan": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read for plan refresh (PlanOnly step)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					// Adding ctr_ prefix to access_group.contract_id is a no-op: the
					// IgnorePrefixType semantic equals considers it unchanged.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/access_group_contract_id_ctr_prefix.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: false,
				},
			},
		},
		"happy path - adding grp_ prefix to group_id results in empty plan": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read for plan refresh (PlanOnly step)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					// Adding grp_ prefix to group_id is a no-op: StringUseStateIf
					// normalises the plan to the stored bare value before PreventStringUpdate fires.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/group_id_grp_prefix.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: false,
				},
			},
		},
		"expect error - access_group.contract_id cannot be updated after creation": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read before update (plan phase)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/access_group_contract_id.tf"),
					ExpectError: regexp.MustCompile("updating field `access_group.contract_id` is not possible"),
				},
			},
		},
		"expect error - access_group.group_id cannot be updated after creation": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read before update (plan phase)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/access_group_group_id.tf"),
					ExpectError: regexp.MustCompile("Cannot update access_group.group_id after creation"),
				},
			},
		},
		"expect error - access_group.group_id cannot be removed after creation": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read before update (plan phase)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					ExpectError: regexp.MustCompile("Cannot update access_group.group_id after creation"),
				},
			},
		},
		"expect error - contract.contract_id cannot be updated after creation": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/update_contract.tf"),
					ExpectError: regexp.MustCompile("updating field `contract.contract_id` is not possible"),
				},
			},
		},
		// ── Read / Drift ─────────────────────────────────────────────────────────
		"happy path - resource removed outside terraform": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create (normal)
				mockGetRepGroup(m, createData)
				// Read from refresh step — returns 404 (deleted outside TF)
				m.On("GetReportingGroup", testutils.MockContext,
					reportinggroups.GetReportingGroupsRequest{ReportingGroupID: createData.reportingGroupID},
				).Return(nil, &reportinggroups.Error{HTTPStatus: http.StatusNotFound}).Once()
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					// Resource was removed outside TF — Terraform removes it from state (no error)
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		},
		// ── Import ──────────────────────────────────────────────────────────────
		"happy path - import by reporting_group_id": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// ImportState
				mockGetRepGroup(m, createData)
				// Read after import (plan check + state set)
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Delete after plan-only step
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: noGroupIDReportingGroup(),
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckMissing("access_group.group_id").
						Build(),
					ImportStateId:      "999",
					ImportState:        true,
					ResourceName:       "akamai_reportinggroups_group.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - import with reporting_group_id and group_id": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// ImportState
				mockGetRepGroup(m, createData)
				// Read after import
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Delete after plan-only step
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("access_group.group_id", "12345").
						Build(),
					ImportStateId:      "999,12345",
					ImportState:        true,
					ResourceName:       "akamai_reportinggroups_group.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - import with reporting_group_id and grp_ prefixed group_id": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// ImportState
				mockGetRepGroup(m, createData)
				// Read after import
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Delete after plan-only step
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					// The grp_ prefix is stripped when parsing the import ID; the stored
					// value matches the bare "12345" in with_group_id.tf.
					ImportStateCheck: importChecker.
						CheckEqual("access_group.group_id", "12345").
						Build(),
					ImportStateId:      "999,grp_12345",
					ImportState:        true,
					ResourceName:       "akamai_reportinggroups_group.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - import reporting group with multiple CP codes": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// ImportState
				mockGetRepGroup(m, createData)
				// Read after import
				mockGetRepGroup(m, createData)
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: multipleCPCodesReportingGroup(),
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("contract.cp_codes.#", "2").
						CheckEqual("access_group.group_id", "12345").
						CheckEqual("contract.cp_codes.1.cp_code_id", "222222").
						CheckEqual("contract.cp_codes.1.cp_code_name", "test-cp_code2").
						Build(),
					ImportStateId:      "999,12345",
					ImportState:        true,
					ResourceName:       "akamai_reportinggroups_group.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/multiple_cp_codes.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/multiple_cp_codes.tf"),
					PlanOnly: true,
				},
			},
		},
		"expect error - import without group_id but config specifies group_id": {
			// If group_id is omitted from the import ID, the imported state has
			// group_id = null. ModifyPlan detects null -> non-null and returns a clear
			// error directing the user to re-import with "reportingGroupID,groupID".
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// ImportState
				mockGetRepGroup(m, createData)
				// Read after import
				mockGetRepGroup(m, createData)
				// Delete (PlanOnly+ExpectError skips refresh, goes straight to cleanup)
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					ImportStateId:      "999",
					ImportState:        true,
					ResourceName:       "akamai_reportinggroups_group.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					ImportStatePersist: true,
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile("Resource was imported without a group_id"),
				},
			},
		},
		"expect error - too many parts in import ID": {
			steps: []resource.TestStep{
				{
					ImportStateId: "999,12345,extra",
					ImportState:   true,
					ResourceName:  "akamai_reportinggroups_group.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					ExpectError:   regexp.MustCompile(`Incorrect import ID`),
				},
			},
		},
		"expect error - non-existent reporting group": {
			init: func(m *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {
				m.On("GetReportingGroup", testutils.MockContext,
					reportinggroups.GetReportingGroupsRequest{ReportingGroupID: 99999},
				).Return(nil, &reportinggroups.Error{HTTPStatus: http.StatusNotFound}).Once()
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "99999",
					ImportState:   true,
					ResourceName:  "akamai_reportinggroups_group.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					ExpectError:   regexp.MustCompile(`Error: Cannot import non-existent remote object`),
				},
			},
		},
		"expect error - non-numeric reporting_group_id in import": {
			steps: []resource.TestStep{
				{
					ImportStateId: "not-a-number",
					ImportState:   true,
					ResourceName:  "akamai_reportinggroups_group.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					ExpectError:   regexp.MustCompile(`Invalid Reporting Group ID in import ID`),
				},
			},
		},
		"expect error - non-numeric group_id in import": {
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					ImportStateId: "999,not-a-number",
					ImportState:   true,
					ResourceName:  "akamai_reportinggroups_group.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/with_group_id.tf"),
					ExpectError:   regexp.MustCompile(`Invalid Group ID in import ID`),
				},
			},
		},
		"expect error - import with group_id API error (non-404)": {
			init: func(m *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {
				m.On("GetReportingGroup", testutils.MockContext,
					reportinggroups.GetReportingGroupsRequest{ReportingGroupID: int64(999)},
				).Return(nil, fmt.Errorf("server error")).Once()
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "999",
					ImportState:   true,
					ResourceName:  "akamai_reportinggroups_group.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/import/min.tf"),
					ExpectError:   regexp.MustCompile(`Unable to import Reporting Group`),
				},
			},
		},
		// ── API error paths ──────────────────────────────────────────────────────
		"expect error - CreateReportingGroup API fails": {
			init: func(m *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {
				m.On("CreateReportingGroup", testutils.MockContext,
					buildCreateRequest(minReportingGroup()),
				).Return(nil, fmt.Errorf("API error")).Once()
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to create Reporting Group(.|\n)*API error`),
				},
			},
		},
		"expect error - GetReportingGroup API fails on read after create": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create succeeds
				mockCreateReportingGroup(m, createData)
				// Read after create fails
				m.On("GetReportingGroup", testutils.MockContext,
					reportinggroups.GetReportingGroupsRequest{ReportingGroupID: createData.reportingGroupID},
				).Return(nil, fmt.Errorf("API read error")).Once()
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to read Reporting Group(.|\n)*API read error`),
				},
			},
		},
		"expect error - UpdateReportingGroup API fails": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, updateData reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read before update
				mockGetRepGroup(m, createData)
				// Update fails
				m.On("UpdateReportingGroup", testutils.MockContext,
					buildUpdateRequest(createData, updateData),
				).Return(nil, fmt.Errorf("API update error")).Once()
				// Read before destroy
				mockGetRepGroup(m, createData)
				// Delete
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			updateMockData: updatedNameReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/update/name.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to update Reporting Group(.|\n)*API update error`),
				},
			},
		},
		"expect error - DeleteReportingGroup API fails": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create
				mockCreateReportingGroup(m, createData)
				// Read after create
				mockGetRepGroup(m, createData)
				// Read before destroy
				mockGetRepGroup(m, createData)
				// Delete fails
				m.On("DeleteReportingGroup", testutils.MockContext,
					reportinggroups.DeleteReportingGroupRequest{ReportingGroupID: createData.reportingGroupID},
				).Return(fmt.Errorf("API delete error")).Once()
				// Delete success
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Check:  minChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					Destroy:     true,
					ExpectError: regexp.MustCompile(`Error: Unable to delete Reporting Group(.|\n)*API delete error`),
				},
			},
		},
		"expect error - Get response returns empty contracts list": {
			init: func(m *reportinggroups.Mock, createData reportingGroupTestData, _ reportingGroupTestData) {
				// Create succeeds
				mockCreateReportingGroup(m, createData)
				// Read after create returns an empty contracts slice
				resp := buildGetResponse(createData)
				resp.Contracts = []reportinggroups.Contract{}
				m.On("GetReportingGroup", testutils.MockContext,
					reportinggroups.GetReportingGroupsRequest{ReportingGroupID: createData.reportingGroupID},
				).Return(resp, nil).Once()
				// Delete (state was set by Create, so Terraform cleans up)
				mockDeleteReportingGroup(m, createData)
			},
			createMockData: minReportingGroup(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/min.tf"),
					ExpectError: regexp.MustCompile(`No Contract Data`),
				},
			},
		},
		// ── Schema validation ────────────────────────────────────────────────────
		"expect error - empty CP codes list": {
			init: func(_ *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/create/empty_cp_codes.tf"),
					ExpectError: regexp.MustCompile(`set must contain at least 1 element`),
				},
			},
		},
		"expect error - missing reporting_group_name": {
			init: func(_ *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/error/missing_reporting_group_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "reporting_group_name" is required`),
				},
			},
		},
		"expect error - missing access_group": {
			init: func(_ *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/error/missing_access_group.tf"),
					ExpectError: regexp.MustCompile(`The argument "access_group" is required`),
				},
			},
		},
		"expect error - missing contract": {
			init: func(_ *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/error/missing_contract.tf"),
					ExpectError: regexp.MustCompile(`The argument "contract" is required`),
				},
			},
		},
		"expect error - group_id is required during creation": {
			init: func(_ *reportinggroups.Mock, _ reportingGroupTestData, _ reportingGroupTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResReportingGroups/error/no_group_id_for_create.tf"),
					ExpectError: regexp.MustCompile("field `access_group.group_id` is required during creation"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()

			if tc.init != nil {
				tc.init(client.ReportingGroups, tc.createMockData, tc.updateMockData)
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

// ── Request / response builders ──────────────────────────────────────────────

// buildCreateRequest constructs the CreateReportingGroupRequest expected by the mock.
func buildCreateRequest(data reportingGroupTestData) reportinggroups.CreateReportingGroupRequest {
	cpCodes := make([]reportinggroups.CPCodeCreate, 0, len(data.cpCodes))
	for _, id := range data.cpCodes {
		cpCodeID, _ := str.GetInt64ID(id, "cpc_") // CP code IDs are sent to API without cpc_ prefix
		cpCodes = append(cpCodes, reportinggroups.CPCodeCreate{CPCodeID: cpCodeID})
	}
	grpID, _ := str.GetInt64ID(data.groupID, "grp_")
	return reportinggroups.CreateReportingGroupRequest{
		ReportingGroupName: data.reportingGroupName,
		AccessGroup: reportinggroups.AccessGroup{
			ContractID: strings.TrimPrefix(data.accessGroupContract, "ctr_"),
			GroupID:    &grpID,
		},
		Contracts: []reportinggroups.ContractCreate{
			{
				ContractID: strings.TrimPrefix(data.contractID, "ctr_"),
				CPCodes:    cpCodes,
			},
		},
	}
}

// buildUpdateRequest constructs the UpdateReportingGroupRequest expected by the mock.
// It uses the ID from stateData and the desired values from planData.
func buildUpdateRequest(stateData, planData reportingGroupTestData) reportinggroups.UpdateReportingGroupRequest {
	cpCodes := make([]reportinggroups.CPCode, 0, len(planData.cpCodes))
	for _, id := range planData.cpCodes {
		cpCodeID, _ := str.GetInt64ID(id, "cpc_")
		cpCodes = append(cpCodes, reportinggroups.CPCode{
			CPCodeID: cpCodeID,
		})
	}

	return reportinggroups.UpdateReportingGroupRequest{
		ReportingGroupID:   stateData.reportingGroupID,
		ReportingGroupName: planData.reportingGroupName,
		Contracts: []reportinggroups.Contract{
			{
				ContractID: strings.TrimPrefix(planData.contractID, "ctr_"),
				CPCodes:    cpCodes,
			},
		},
	}
}

// buildCreateResponse constructs the CreateReportingGroupResponse returned by the mock.
func buildCreateResponse(data reportingGroupTestData) *reportinggroups.CreateReportingGroupResponse {
	return &reportinggroups.CreateReportingGroupResponse{
		ReportingGroup: reportinggroups.ReportingGroup{
			ReportingGroupID:   data.reportingGroupID,
			ReportingGroupName: data.reportingGroupName,
			AccessGroup: reportinggroups.AccessGroup{
				ContractID: strings.TrimPrefix(data.accessGroupContract, "ctr_"),
				// GroupID is always null in API responses
			},
			Contracts: buildContractModels(data),
		},
	}
}

// buildUpdateResponse constructs the UpdateReportingGroupResponse returned by the mock.
func buildUpdateResponse(data reportingGroupTestData) *reportinggroups.UpdateReportingGroupResponse {
	return &reportinggroups.UpdateReportingGroupResponse{
		ReportingGroupID:   data.reportingGroupID,
		ReportingGroupName: data.reportingGroupName,
		AccessGroup: reportinggroups.AccessGroup{
			ContractID: strings.TrimPrefix(data.accessGroupContract, "ctr_"),
			// GroupID is always null in API responses
		},
		Contracts: buildContractModels(data),
	}
}

// buildGetResponse constructs the GetReportingGroupResponse returned by the mock.
func buildGetResponse(data reportingGroupTestData) *reportinggroups.GetReportingGroupResponse {
	return &reportinggroups.GetReportingGroupResponse{
		ReportingGroupID:   data.reportingGroupID,
		ReportingGroupName: data.reportingGroupName,
		AccessGroup: reportinggroups.AccessGroup{
			ContractID: strings.TrimPrefix(data.accessGroupContract, "ctr_"),
			// GroupID is always null in API responses
		},
		Contracts: buildContractModels(data),
	}
}

// buildContractModels pairs each CP code ID with its name by index.
func buildContractModels(data reportingGroupTestData) []reportinggroups.Contract {
	cpCodes := make([]reportinggroups.CPCode, 0, len(data.cpCodes))
	for i, id := range data.cpCodes {
		name := ""
		if i < len(data.cpCodeNames) {
			name = data.cpCodeNames[i]
		}
		cpCodeID, _ := str.GetInt64ID(id, "cpc_")
		cpCodes = append(cpCodes, reportinggroups.CPCode{
			CPCodeID:   cpCodeID,
			CPCodeName: name,
		})
	}
	return []reportinggroups.Contract{
		{
			ContractID: strings.TrimPrefix(data.contractID, "ctr_"),
			CPCodes:    cpCodes,
		},
	}
}

// ── Mock helpers ─────────────────────────────────────────────────────────────

func mockCreateReportingGroup(m *reportinggroups.Mock, data reportingGroupTestData) *mock.Call {
	return m.On("CreateReportingGroup", testutils.MockContext,
		buildCreateRequest(data),
	).Return(buildCreateResponse(data), nil).Once()
}

func mockGetRepGroup(m *reportinggroups.Mock, data reportingGroupTestData) *mock.Call {
	return m.On("GetReportingGroup", testutils.MockContext,
		reportinggroups.GetReportingGroupsRequest{ReportingGroupID: data.reportingGroupID},
	).Return(buildGetResponse(data), nil).Once()
}

// mockUpdateReportingGroup sets up the Update mock using stateData for the ID
// and planData for the desired values in the request / response.
func mockUpdateReportingGroup(m *reportinggroups.Mock, stateData, planData reportingGroupTestData) *mock.Call {
	return m.On("UpdateReportingGroup", testutils.MockContext,
		buildUpdateRequest(stateData, planData),
	).Return(buildUpdateResponse(planData), nil).Once()
}

func mockDeleteReportingGroup(m *reportinggroups.Mock, data reportingGroupTestData) *mock.Call {
	return m.On("DeleteReportingGroup", testutils.MockContext,
		reportinggroups.DeleteReportingGroupRequest{ReportingGroupID: data.reportingGroupID},
	).Return(nil).Once()
}
