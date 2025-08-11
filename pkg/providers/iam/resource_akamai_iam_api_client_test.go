package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/internal/test"
	tst "github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type testData struct {
	createAPIClientRequest                   iam.CreateAPIClientRequest
	createAPIClientResponse                  iam.CreateAPIClientResponse
	updateCredentialRequest                  iam.UpdateCredentialRequest
	updateAPIClientNotificationEmailsRequest iam.UpdateAPIClientRequest
	getAPIClientRequest                      iam.GetAPIClientRequest
	getAPIClientResponse                     iam.GetAPIClientResponse
	updateAPIClientRequest                   iam.UpdateAPIClientRequest
	updateAPIClientResponse                  iam.UpdateAPIClientResponse
}

func TestResourceAPIClient(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		configPath string
		init       func(*iam.Mock, testData, testData)
		createData testData
		updateData testData
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"happy path - create with all fields set": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
			},
		},
		"happy path - create with all fields set but unknown at plan": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(3)
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all_unknown.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
			},
		},
		"happy path - create with all fields set but cp_code_access unknown at plan": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(3)
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all_cp_code_access_unknown.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
			},
		},
		"happy path - create with all fields set, all_accessible_apis is true": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				// Create
				purgeOptions := iam.PurgeOptions{
					CanPurgeByCacheTag: true,
					CanPurgeByCPCode:   true,
					CPCodeAccess: iam.CPCodeAccess{
						AllCurrentAndNewCPCodes: false,
						CPCodes:                 []int64{},
					},
				}
				createData.createAPIClientRequest.APIAccess.AllAccessibleAPIs = true
				createData.createAPIClientRequest.APIAccess.APIs = nil
				createData.createAPIClientRequest.PurgeOptions = &purgeOptions
				createData.createAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				createData.createAPIClientResponse.APIAccess.APIs = nil
				createData.createAPIClientResponse.PurgeOptions = &purgeOptions
				createData.getAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				createData.getAPIClientResponse.APIAccess.APIs = nil
				createData.getAPIClientResponse.PurgeOptions = &purgeOptions

				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all_all_accessible_apis_true.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("api_access.all_accessible_apis", "true").
						CheckEqual("api_access.apis.#", "0").
						CheckMissing("api_access.apis.0.access_level").
						CheckMissing("api_access.apis.0.api_id").
						CheckMissing("api_access.apis.0.api_name").
						CheckMissing("api_access.apis.0.description").
						CheckMissing("api_access.apis.0.documentation_url").
						CheckMissing("api_access.apis.0.endpoint").
						CheckMissing("api_access.apis.1.access_level").
						CheckMissing("api_access.apis.1.api_id").
						CheckMissing("api_access.apis.1.api_name").
						CheckMissing("api_access.apis.1.description").
						CheckMissing("api_access.apis.1.documentation_url").
						CheckMissing("api_access.apis.1.endpoint").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - create with all fields set, no cpcodes": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.createAPIClientRequest.PurgeOptions = &purgeOptionsNoCPCodes
				createData.createAPIClientResponse.PurgeOptions = &purgeOptionsNoCPCodes
				createData.getAPIClientResponse.PurgeOptions = &purgeOptionsNoCPCodes
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all_no_cpcodes.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("purge_options.cp_code_access.all_current_and_new_cp_codes", "true").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - create with all fields set, missing cpcodes": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.createAPIClientRequest.PurgeOptions = &purgeOptionsNoCPCodes
				createData.createAPIClientResponse.PurgeOptions = &purgeOptionsNoCPCodes
				createData.getAPIClientResponse.PurgeOptions = &purgeOptionsNoCPCodes
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all_missing_cpcodes.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("purge_options.cp_code_access.all_current_and_new_cp_codes", "true").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - create with min set of fields": {
			init: func(m *iam.Mock, createData, _ testData) {
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockLockAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: minData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - create with min set of fields, api_access is unknown at plan": {
			init: func(m *iam.Mock, createData, _ testData) {
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockLockAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: minData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min_api_access_unknown.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - create with all fields set and custom credential details": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateCredential(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullDataWithCredential,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.expires_on", "2026-06-13T14:48:07Z").
						CheckEqual("credential.description", "Test API Client Credential").
						Build(),
				},
			},
		},
		"happy path - create with all fields set and credential inactive status": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateCredential(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeleteAPIClient(m, createData)
			},
			createData: fullDataWithInactiveCredential,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_inactive.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.expires_on", "2025-06-13T14:48:07Z").
						CheckEqual("credential.description", "Test API Client Credential").
						CheckEqual("credential.status", "INACTIVE").
						CheckEqual("credential.actions.deactivate", "false").
						CheckEqual("credential.actions.activate", "true").
						CheckEqual("credential.actions.delete", "true").
						CheckEqual("active_credential_count", "0").
						Build(),
				},
			},
		},
		"happy path - create with credential, encounter diff in read which returns credential as DELETED": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateCredential(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeleteAPIClient(m, createData)
			},
			createData: fullDataWithCredentialDrift,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.expires_on", "2026-06-13T14:48:07Z").
						CheckEqual("credential.description", "Test API Client Credential").
						CheckEqual("credential.status", "DELETED").
						CheckEqual("credential.actions.deactivate", "false").
						CheckEqual("credential.actions.activate", "false").
						CheckEqual("credential.actions.delete", "false").
						CheckEqual("credential.actions.edit_description", "false").
						CheckEqual("credential.actions.edit_expiration", "false").
						CheckEqual("active_credential_count", "1").
						Build(),
				},
			},
		},
		"happy path - create without credential, update by changing description and expires_on": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(8)
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				mockUpdateCredential(m, updateData)
				mockGetAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, updateData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: fullData,
			updateData: fullDataWithCredential,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.description", "Test API Client Credential").
						CheckEqual("credential.expires_on", "2026-06-13T14:48:07Z").
						Build(),
				},
			},
		},
		"happy path - create without credential, update by changing description, expires_on and status": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(8)
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				mockUpdateCredential(m, updateData)
				mockGetAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeleteAPIClient(m, updateData)
			},
			createData: fullData,
			updateData: fullDataWithInactiveCredential,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_inactive.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.description", "Test API Client Credential").
						CheckEqual("credential.status", "INACTIVE").
						CheckEqual("credential.expires_on", "2025-06-13T14:48:07Z").
						CheckEqual("credential.actions.deactivate", "false").
						CheckEqual("credential.actions.activate", "true").
						CheckEqual("credential.actions.delete", "true").
						CheckEqual("active_credential_count", "0").
						Build(),
				},
			},
		},
		"happy path - create with credential, update by changing description and expires_on": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(8)
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateCredential(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				mockUpdateCredential(m, updateData)
				mockGetAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, updateData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: fullDataWithCredential,
			updateData: fullDataWithCredential2,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.expires_on", "2026-06-13T14:48:07Z").
						CheckEqual("credential.description", "Test API Client Credential").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail_2.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.description", "Test API Client Credential 2").
						CheckEqual("credential.expires_on", "2027-06-13T14:48:07Z").
						Build(),
				},
			},
		},
		"happy path - create with credential, update by changing description, expires_on and status": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(8)
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateCredential(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				mockUpdateCredential(m, updateData)
				mockGetAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeleteAPIClient(m, updateData)
			},
			createData: fullDataWithCredential,
			updateData: fullDataWithInactiveCredential2,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.expires_on", "2026-06-13T14:48:07Z").
						CheckEqual("credential.description", "Test API Client Credential").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_inactive_2.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.description", "Test API Client Credential 2").
						CheckEqual("credential.status", "INACTIVE").
						CheckEqual("credential.expires_on", "2027-06-13T14:48:07Z").
						CheckEqual("credential.actions.deactivate", "false").
						CheckEqual("credential.actions.activate", "true").
						CheckEqual("credential.actions.delete", "true").
						CheckEqual("active_credential_count", "0").
						Build(),
				},
			},
		},
		"happy path - create with credential, update by changing description, expires_on and api client description - expect update credential and api client": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(8)
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateCredential(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				mockUpdateCredential(m, updateData)
				mockUpdateAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, updateData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: fullDataWithCredential,
			updateData: fullDataWithCredential2AndClientDescription,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.expires_on", "2026-06-13T14:48:07Z").
						CheckEqual("credential.description", "Test API Client Credential").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/with_credential/create_all_with_credential_detail_2_and_api_client_description.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("credential.description", "Test API Client Credential 2").
						CheckEqual("credential.expires_on", "2027-06-13T14:48:07Z").
						CheckEqual("client_description", "Test API Client 2").
						Build(),
				},
			},
		},
		"happy path - create with all fields set but `ip_acl.enable` as false no `cidr`": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.createAPIClientRequest.IPACL = &ipACLNoCidr
				createData.createAPIClientResponse.IPACL = &ipACLNoCidr
				createData.getAPIClientResponse.IPACL = &ipACLNoCidr
				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_no_cidr.tf"),
					Check: fullDataChecker.
						CheckMissing("ip_acl.cidr.0").
						CheckEqual("ip_acl.enable", "false").
						Build(),
				},
			},
		},
		"happy path - purge_options.cp_code_access.cp_codes can be provided when all_accessible_apis is true": {
			init: func(m *iam.Mock, createData, _ testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.createAPIClientRequest.APIAccess = iam.APIAccessRequest{
					AllAccessibleAPIs: true,
				}
				createData.createAPIClientResponse.APIAccess = apiAccess
				createData.getAPIClientResponse.APIAccess = apiAccess

				createData.createAPIClientRequest.PurgeOptions = &purgeOptionsCPCodes
				createData.createAPIClientResponse.PurgeOptions = &purgeOptionsCPCodes
				createData.getAPIClientResponse.PurgeOptions = &purgeOptionsCPCodes

				// Create
				mockCreateAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, createData)
			},
			createData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/cp_codes_can_be_provided_when_all_accessible_apis_true.tf"),
					Check: tst.NewStateChecker("akamai_iam_api_client.test").
						CheckEqual("api_access.all_accessible_apis", "true").
						CheckMissing("api_access.all_accessible_apis.0").
						CheckEqual("purge_options.can_purge_by_cache_tag", "false").
						CheckEqual("purge_options.can_purge_by_cp_code", "false").
						CheckEqual("purge_options.cp_code_access.all_current_and_new_cp_codes", "false").
						CheckEqual("purge_options.cp_code_access.cp_codes.0", "101").
						Build(),
				},
			},
		},
		"happy path - client does not have access to the cp code - expect error": {
			init: func(m *iam.Mock, _, _ testData) {
				m.On("ListAllowedCPCodes", testutils.MockContext, iam.ListAllowedCPCodesRequest{
					UserName: "mw+2",
					Body: iam.ListAllowedCPCodesRequestBody{
						ClientType: "CLIENT",
						Groups: []iam.ClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
				}).Return(iam.ListAllowedCPCodesResponse{
					{},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/vanished_cp_codes_after_creation.tf"),
					ExpectError: regexp.MustCompile(`Error: provided invalid data`),
				},
			},
		},
		"happy path - update from min to full and from locked to unlocked": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.getAPIClientResponse.IsLocked = true
				createData.getAPIClientResponse.Actions = &lockedClientActions
				createData.createAPIClientResponse.IsLocked = true
				createData.createAPIClientResponse.Actions = &lockedClientActions
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockLockAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				mockUpdateAPIClient(m, updateData)
				mockUnlockAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: minData,
			updateData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "true").
						CheckEqual("actions.lock", "false").
						CheckEqual("actions.unlock", "true").
						CheckEqual("lock", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/update_all.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "false").
						CheckEqual("actions.lock", "true").
						CheckEqual("actions.unlock", "false").
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
			},
		},
		"happy path - update from min to full and all_accessible_apis from false to true ": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.getAPIClientResponse.IsLocked = true
				createData.getAPIClientResponse.Actions = &lockedClientActions
				createData.createAPIClientResponse.IsLocked = true
				createData.createAPIClientResponse.Actions = &lockedClientActions
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockLockAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)
				// Update
				updateData.updateAPIClientRequest.Body.ClientDescription = ""
				updateData.updateAPIClientRequest.Body.NotificationEmails = []string{}
				updateData.updateAPIClientRequest.Body.APIAccess.AllAccessibleAPIs = true
				updateData.updateAPIClientRequest.Body.APIAccess.APIs = nil
				updateData.updateAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				updateData.updateAPIClientResponse.APIAccess.APIs = apisAllGet
				updateData.updateAPIClientResponse.ClientDescription = ""
				updateData.updateAPIClientResponse.NotificationEmails = []string{}
				updateData.getAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				updateData.getAPIClientResponse.APIAccess.APIs = apisAllGet
				updateData.getAPIClientResponse.ClientDescription = ""
				updateData.getAPIClientResponse.NotificationEmails = []string{}

				mockUpdateAPIClient(m, updateData)
				mockUnlockAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: minData,
			updateData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "true").
						CheckEqual("actions.lock", "false").
						CheckEqual("actions.unlock", "true").
						CheckEqual("lock", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/update_all_all_accessible_apis_true.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "false").
						CheckEqual("actions.lock", "true").
						CheckEqual("actions.unlock", "false").
						CheckEqual("api_access.all_accessible_apis", "true").
						CheckEqual("api_access.apis.#", "5").
						CheckEqual("api_access.apis.0.access_level", "READ-ONLY").
						CheckEqual("api_access.apis.0.api_id", "5580").
						CheckEqual("api_access.apis.0.api_name", "Search Data Feed").
						CheckEqual("api_access.apis.0.description", "Search Data Feed").
						CheckEqual("api_access.apis.0.documentation_url", "/").
						CheckEqual("api_access.apis.0.endpoint", "/search-portal-data-feed-api/").
						CheckEqual("api_access.apis.1.access_level", "READ-ONLY").
						CheckEqual("api_access.apis.1.api_id", "6681").
						CheckEqual("api_access.apis.1.api_name", "Test").
						CheckEqual("api_access.apis.1.description", "Test").
						CheckEqual("api_access.apis.1.documentation_url", "-").
						CheckEqual("api_access.apis.1.endpoint", "/test").
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						Build(),
				},
			},
		},
		"happy path - create with clone_authorized_user_groups false and one group, update to clone_authorized_user_groups true, no groups": {
			init: func(m *iam.Mock, createData, updateData testData) {

				createData.createAPIClientRequest.GroupAccess.Groups = []iam.ClientGroupRequestItem{
					{
						GroupID: 578,
						RoleID:  341,
					},
				}
				createData.createAPIClientResponse.IsLocked = false
				createData.createAPIClientResponse.GroupAccess.Groups = []iam.ClientGroup{
					{
						GroupID: 578,
						RoleID:  341,
					},
				}

				createData.getAPIClientResponse.IsLocked = false
				createData.getAPIClientResponse.GroupAccess.Groups = []iam.ClientGroup{
					{
						GroupID:         578,
						RoleID:          341,
						GroupName:       "group3",
						IsBlocked:       false,
						ParentGroupID:   0,
						RoleDescription: "group description 3",
						RoleName:        "role 3",
						Subgroups:       nil,
					},
				}

				createData.updateAPIClientNotificationEmailsRequest.Body.GroupAccess.Groups = []iam.ClientGroupRequestItem{
					{
						GroupID: 578,
						RoleID:  341,
					},
				}

				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)

				// Update
				updateData.updateAPIClientRequest.Body.GroupAccess.CloneAuthorizedUserGroups = true
				updateData.updateAPIClientRequest.Body.GroupAccess.Groups = nil
				updateData.updateAPIClientResponse.IsLocked = false
				updateData.updateAPIClientResponse.GroupAccess.CloneAuthorizedUserGroups = true
				updateData.getAPIClientResponse.IsLocked = false
				updateData.getAPIClientResponse.GroupAccess.CloneAuthorizedUserGroups = true

				mockUpdateAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: minData,
			updateData: minData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min_clone_authorized_user_groups_false.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.clone_authorized_user_groups", "false").
						CheckEqual("group_access.groups.#", "1").
						CheckEqual("group_access.groups.0.group_id", "578").
						CheckEqual("group_access.groups.0.group_name", "group3").
						CheckEqual("group_access.groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.role_description", "group description 3").
						CheckEqual("group_access.groups.0.role_id", "341").
						CheckEqual("group_access.groups.0.role_name", "role 3").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/update_min_clone_authorized_user_groups_true.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.clone_authorized_user_groups", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - create with clone_authorized_user_groups true and no explicit group, update to clone_authorized_user_groups false, one group": {
			init: func(m *iam.Mock, createData, updateData testData) {
				createData.createAPIClientRequest.GroupAccess.CloneAuthorizedUserGroups = true
				createData.createAPIClientRequest.GroupAccess.Groups = nil
				createData.createAPIClientResponse.IsLocked = false
				createData.createAPIClientResponse.GroupAccess.CloneAuthorizedUserGroups = true
				createData.updateAPIClientNotificationEmailsRequest.Body.GroupAccess.CloneAuthorizedUserGroups = true
				createData.updateAPIClientNotificationEmailsRequest.Body.GroupAccess.Groups = nil
				createData.getAPIClientResponse.IsLocked = false
				createData.getAPIClientResponse.GroupAccess.CloneAuthorizedUserGroups = true

				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)

				// Update
				updateData.updateAPIClientRequest.Body.GroupAccess.Groups = []iam.ClientGroupRequestItem{
					{
						GroupID: 578,
						RoleID:  341,
					},
				}
				updateData.updateAPIClientResponse.IsLocked = false
				updateData.updateAPIClientResponse.GroupAccess.Groups = []iam.ClientGroup{
					{
						GroupID:         578,
						RoleID:          341,
						GroupName:       "group3",
						IsBlocked:       false,
						ParentGroupID:   0,
						RoleDescription: "group description 3",
						RoleName:        "role 3",
						Subgroups:       nil,
					},
				}

				updateData.getAPIClientResponse.IsLocked = false
				updateData.getAPIClientResponse.GroupAccess.Groups = []iam.ClientGroup{
					{
						GroupID:         578,
						RoleID:          341,
						GroupName:       "group3",
						IsBlocked:       false,
						ParentGroupID:   0,
						RoleDescription: "group description 3",
						RoleName:        "role 3",
						Subgroups:       nil,
					},
				}

				mockUpdateAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: minData,
			updateData: minData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min_clone_authorized_user_groups_true.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.clone_authorized_user_groups", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/update_min_clone_authorized_user_groups_false.tf"),
					Check: fullDataChecker.
						CheckEqual("group_access.clone_authorized_user_groups", "false").
						CheckEqual("group_access.groups.#", "1").
						CheckEqual("group_access.groups.0.group_id", "578").
						CheckEqual("group_access.groups.0.group_name", "group3").
						CheckEqual("group_access.groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.role_description", "group description 3").
						CheckEqual("group_access.groups.0.role_id", "341").
						CheckEqual("group_access.groups.0.role_name", "role 3").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"happy path - update from min to full and from unlocked to locked": {
			init: func(m *iam.Mock, createData, updateData testData) {
				mockListAllowedCPCodes(m).Times(4)
				createData.getAPIClientResponse.IsLocked = false
				createData.getAPIClientResponse.Actions = &unlockedClientActions
				createData.createAPIClientResponse.IsLocked = false
				createData.createAPIClientResponse.Actions = &unlockedClientActions
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)

				updateData.updateAPIClientResponse.IsLocked = true
				updateData.updateAPIClientResponse.Actions = &lockedClientActions
				updateData.getAPIClientResponse.IsLocked = true
				updateData.getAPIClientResponse.Actions = &lockedClientActions

				// Update
				mockUpdateAPIClient(m, updateData)
				mockLockAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: minData,
			updateData: fullData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min_unlock.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "false").
						CheckEqual("actions.lock", "true").
						CheckEqual("actions.unlock", "false").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_all_lock.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "true").
						CheckEqual("actions.lock", "false").
						CheckEqual("actions.unlock", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
			},
		},
		"happy path - update only from unlocked to locked": {
			init: func(m *iam.Mock, createData, updateData testData) {
				createData.getAPIClientResponse.IsLocked = false
				createData.getAPIClientResponse.Actions = &unlockedClientActions
				createData.createAPIClientResponse.IsLocked = false
				createData.createAPIClientResponse.Actions = &unlockedClientActions
				// Create
				mockCreateAPIClient(m, createData)
				mockUpdateAPIClientNotificationEmails(m, createData)
				mockGetAPIClient(m, createData)
				// Read
				mockGetAPIClient(m, createData)
				mockGetAPIClient(m, createData)

				updateData.updateAPIClientResponse.IsLocked = true
				updateData.updateAPIClientResponse.Actions = &lockedClientActions
				updateData.getAPIClientResponse.IsLocked = true
				updateData.getAPIClientResponse.Actions = &lockedClientActions

				// Update
				mockUpdateAPIClient(m, updateData)
				mockLockAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, createData)
				mockDeleteAPIClient(m, updateData)
			},
			createData: minData,
			updateData: minData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min_unlock.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "false").
						CheckEqual("actions.lock", "true").
						CheckEqual("actions.unlock", "false").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "true").
						CheckEqual("actions.lock", "false").
						CheckEqual("actions.unlock", "true").
						CheckEqual("group_access.groups.0.sub_groups.#", "0").
						CheckEqual("client_description", "").
						CheckMissing("notification_emails.0").
						CheckMissing("ip_acl.enable").
						CheckMissing("ip_acl.cidr.0").
						CheckMissing("purge_options.can_purge_by_cache_tag").
						CheckMissing("purge_options.can_purge_by_cp_code").
						CheckMissing("purge_options.cp_code_access.all_current_and_new_cp_codes").
						CheckMissing("purge_options.cp_code_access.cp_codes.0").
						Build(),
				},
			},
		},
		"validation error - 'groups' should be provided when 'clone_authorized_user_groups' is true": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/invalid_groups.tf"),
					ExpectError: regexp.MustCompile(`You cannot specify any group when 'clone_authorized_user_groups' is true`),
				},
			},
		},
		"validation error - 'apis' should be required when 'all_accessible_apis' is false": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/missing_apis.tf"),
					ExpectError: regexp.MustCompile(`You must specify at least one API when 'all_accessible_apis' is false`),
				},
			},
		},
		"validation error - 'apis' should not be provided when 'all_accessible_apis' is true": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/invalid_all_accessible_apis.tf"),
					ExpectError: regexp.MustCompile(`You cannot specify any API when 'all_accessible_apis' is true`),
				},
			},
		},
		"validation error - 'purge_options' are required when 'all_accessible_apis' is true": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/missing_purge_options_when_all_accessible_apis.tf"),
					ExpectError: regexp.MustCompile(`You must specify 'purge_options' when 'all_accessible_apis' is true`),
				},
			},
		},
		"validation error - 'purge_options' are mandatory when one of the 'apis' is 'CCU'": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/invalid_all_accessible_apis.tf"),
					ExpectError: regexp.MustCompile(`You cannot specify any API when 'all_accessible_apis' is true`),
				},
			},
		},
		"validation error - 'groups' should be required if 'cloneAuthorizedUserGroups' is false": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/missing_groups.tf"),
					ExpectError: regexp.MustCompile(`You must specify at least one group when 'clone_authorized_user_groups' is\s+false`),
				},
			},
		},
		"validation error - 'cp_codes' should be empty when 'all_current_and_new_cp_codes' are true": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/invalid_cp_codes.tf"),
					ExpectError: regexp.MustCompile(`You cannot specify any CP Code when 'all_current_and_new_cp_codes' is true`),
				},
			},
		},
		"validation error - 'ip_acl.cidr' must not be empty when `ip_acl.enable` is true": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/missing_cidr.tf"),
					ExpectError: regexp.MustCompile(`You should specify 'cidr' when 'enable' is true`),
				},
			},
		},
		"expect error - create": {
			init: func(m *iam.Mock, _, _ testData) {
				m.On("CreateAPIClient", testutils.MockContext, createAPIClientRequestMin).Return(nil, fmt.Errorf("create failed")).Once()
			},
			createData: minData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/create_min.tf"),
					ExpectError: regexp.MustCompile(`create failed`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client, tc.createData, tc.updateData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ExternalProviders: map[string]resource.ExternalProvider{
						"random": {
							Source:            "registry.terraform.io/hashicorp/random",
							VersionConstraint: "3.1.0",
						},
					},
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}

}

func TestImportAPIClientResource(t *testing.T) {
	t.Parallel()
	importChecker := tst.NewImportChecker().
		CheckEqual("access_token", "access_token").
		CheckEqual("actions.deactivate_all", "false").
		CheckEqual("actions.delete", "true").
		CheckEqual("actions.edit", "true").
		CheckEqual("actions.edit_apis", "true").
		CheckEqual("actions.edit_auth", "true").
		CheckEqual("actions.edit_groups", "true").
		CheckEqual("actions.edit_ip_acl", "true").
		CheckEqual("actions.edit_switch_account", "false").
		CheckEqual("actions.lock", "true").
		CheckEqual("actions.transfer", "true").
		CheckEqual("actions.unlock", "false").
		CheckEqual("active_credential_count", "1").
		CheckEqual("allow_account_switch", "false").
		CheckEqual("api_access.all_accessible_apis", "false").
		CheckEqual("api_access.apis.#", "2").
		CheckEqual("api_access.apis.0.access_level", "READ-ONLY").
		CheckEqual("api_access.apis.0.api_id", "5580").
		CheckEqual("api_access.apis.0.api_name", "Search Data Feed").
		CheckEqual("api_access.apis.0.description", "Search Data Feed").
		CheckEqual("api_access.apis.0.documentation_url", "/").
		CheckEqual("api_access.apis.0.endpoint", "/search-portal-data-feed-api/").
		CheckEqual("api_access.apis.1.access_level", "READ-WRITE").
		CheckEqual("api_access.apis.1.api_id", "5801").
		CheckEqual("api_access.apis.1.api_name", "EdgeWorkers").
		CheckEqual("api_access.apis.1.description", "EdgeWorkers").
		CheckEqual("api_access.apis.1.documentation_url", "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html").
		CheckEqual("api_access.apis.1.endpoint", "/edgeworkers/").
		CheckEqual("authorized_users.0", "mw+2").
		CheckEqual("base_url", "base_url").
		CheckEqual("can_auto_create_credential", "false").
		CheckEqual("client_description", "Test API Client").
		CheckEqual("client_id", "c1ien41d").
		CheckEqual("client_name", "mw+2_1").
		CheckEqual("client_type", "CLIENT").
		CheckEqual("created_by", "someuser").
		CheckEqual("created_date", "2024-06-13T14:48:07Z").
		CheckEqual("credential.description", "Update this credential").
		CheckEqual("credential.status", "ACTIVE").
		CheckEqual("credential.expires_on", "2025-06-13T14:48:08Z").
		CheckEqual("credential.actions.deactivate", "true").
		CheckEqual("credential.actions.activate", "false").
		CheckEqual("credential.actions.delete", "false").
		CheckEqual("credential.actions.edit_expiration", "true").
		CheckEqual("credential.actions.edit_description", "true").
		CheckEqual("group_access.clone_authorized_user_groups", "false").
		CheckEqual("group_access.groups.#", "1").
		CheckEqual("group_access.groups.0.group_id", "123").
		CheckEqual("group_access.groups.0.group_name", "group2").
		CheckEqual("group_access.groups.0.is_blocked", "false").
		CheckEqual("group_access.groups.0.parent_group_id", "0").
		CheckEqual("group_access.groups.0.role_description", "group description").
		CheckEqual("group_access.groups.0.role_id", "340").
		CheckEqual("group_access.groups.0.role_name", "role").
		CheckEqual("group_access.groups.0.sub_groups.#", "1").
		CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
		CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
		CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
		CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
		CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
		CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
		CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
		CheckEqual("id", "c1ien41d").
		CheckEqual("ip_acl.cidr.0", "128.5.6.5/24").
		CheckEqual("ip_acl.enable", "true").
		CheckEqual("lock", "false").
		CheckEqual("notification_emails.0", "mw+2@example.com").
		CheckEqual("purge_options.can_purge_by_cache_tag", "true").
		CheckEqual("purge_options.can_purge_by_cp_code", "true").
		CheckEqual("purge_options.cp_code_access.all_current_and_new_cp_codes", "false").
		CheckEqual("purge_options.cp_code_access.cp_codes.0", "101")
	tests := map[string]struct {
		importID    string
		init        func(*iam.Mock, testData, testData)
		importData  testData
		updateData  testData
		expectError *regexp.Regexp
		steps       []resource.TestStep
	}{
		"happy path with import id": {
			importID: "c1ien41d",
			importData: testData{
				getAPIClientRequest:  getAPIClientRequest,
				getAPIClientResponse: getAPIClientResponse,
			},
			init: func(m *iam.Mock, data, _ testData) {
				// Import
				mockGetAPIClient(m, data)
				// Read
				mockGetAPIClient(m, data)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.Build(),
					ImportStateId:    "c1ien41d",
					ImportState:      true,
					ResourceName:     "akamai_iam_api_client.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable.tf"),
				},
			},
		},
		"happy path with import id and allAccessibleApis true": {
			importID: "c1ien41d",
			importData: testData{
				getAPIClientRequest:  getAPIClientRequest,
				getAPIClientResponse: getAPIClientResponseAAA,
			},
			init: func(m *iam.Mock, data, _ testData) {
				// Import
				mockGetAPIClient(m, data)
				// Read
				mockGetAPIClient(m, data)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("api_access.all_accessible_apis", "true").
						CheckEqual("lock", "true").
						Build(),
					ImportStateId: "c1ien41d",
					ImportState:   true,
					ResourceName:  "akamai_iam_api_client.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable_all_accessible_apis_true.tf"),
				},
			},
		},
		"happy path without import id (self)": {
			importID: " ",
			importData: testData{
				getAPIClientRequest:  getAPIClientRequestSelf,
				getAPIClientResponse: getAPIClientResponse,
			},
			init: func(m *iam.Mock, data, _ testData) {
				// Import
				mockGetAPIClient(m, data)
				// Read
				mockGetAPIClient(m, data)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.Build(),
					ImportStateId:    " ",
					ImportState:      true,
					ResourceName:     "akamai_iam_api_client.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable.tf"),
				},
			},
		},
		"happy path with import with lock, update with unlock": {
			importID: "c1ien41d",
			importData: testData{
				getAPIClientRequest:  getAPIClientRequest,
				getAPIClientResponse: getAPIClientResponse,
			},
			updateData: fullData,
			init: func(m *iam.Mock, data, updateData testData) {
				mockListAllowedCPCodes(m).Times(4)
				data.getAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				data.getAPIClientResponse.IsLocked = true
				data.getAPIClientResponse.Actions = &lockedClientActions
				// Import
				mockGetAPIClient(m, data)
				// Read
				mockGetAPIClient(m, data)

				// Update
				mockGetAPIClient(m, data)

				// Update
				updateData.updateAPIClientRequest.Body.APIAccess.AllAccessibleAPIs = true
				updateData.updateAPIClientRequest.Body.APIAccess.APIs = nil
				updateData.updateAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				updateData.updateAPIClientResponse.IsLocked = false
				updateData.updateAPIClientResponse.Actions = &unlockedClientActions
				updateData.getAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				updateData.getAPIClientResponse.IsLocked = false
				updateData.getAPIClientResponse.Actions = &unlockedClientActions

				mockUnlockAPIClient(m, updateData)
				mockUpdateAPIClient(m, updateData)
				// Read
				mockGetAPIClient(m, updateData)
				// Delete
				mockDeactivateCredential(m, updateData)
				mockDeleteAPIClient(m, updateData)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("lock", "true").
						CheckEqual("actions.lock", "false").
						CheckEqual("actions.unlock", "true").
						CheckEqual("api_access.all_accessible_apis", "true").
						Build(),
					ImportStateId:      "c1ien41d",
					ImportState:        true,
					ImportStatePersist: true,
					ResourceName:       "akamai_iam_api_client.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable_lock.tf"),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable_unlock.tf"),
					Check: fullDataChecker.
						CheckEqual("lock", "false").
						CheckEqual("api_access.all_accessible_apis", "true").
						CheckEqual("actions.lock", "true").
						CheckEqual("actions.unlock", "false").
						CheckMissing("credential.client_secret").
						CheckEqual("group_access.groups.0.sub_groups.#", "1").
						CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
						CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
						CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
						CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
						CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
						CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
						CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
						Build(),
				},
			},
		},
		"import with lock, update without unlock should fail": {
			importID: "c1ien41d",
			importData: testData{
				getAPIClientRequest:  getAPIClientRequest,
				getAPIClientResponse: getAPIClientResponse,
			},
			updateData: fullData,
			init: func(m *iam.Mock, data, updateData testData) {
				mockListAllowedCPCodes(m).Once()
				data.getAPIClientResponse.APIAccess.AllAccessibleAPIs = true
				data.getAPIClientResponse.IsLocked = true
				data.getAPIClientResponse.Actions = &lockedClientActions
				// Import
				mockGetAPIClient(m, data)
				// Read after import
				mockGetAPIClient(m, data)
				// Update attempt
				mockGetAPIClient(m, data)
				// Delete
				mockDeactivateCredential(m, updateData)
				mockDeleteAPIClient(m, updateData)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("lock", "true").
						CheckEqual("actions.lock", "false").
						CheckEqual("actions.unlock", "true").
						CheckEqual("api_access.all_accessible_apis", "true").
						Build(),
					ImportStateId:      "c1ien41d",
					ImportState:        true,
					ImportStatePersist: true,
					ResourceName:       "akamai_iam_api_client.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable_lock.tf"),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable_lock_description.tf"),
					ExpectError: regexp.MustCompile("You cannot change API client without unlocking it first."),
				},
			},
		},
		"expect error - unknown import id": {
			importID: "unknown",
			init: func(m *iam.Mock, _, _ testData) {
				m.On("GetAPIClient", testutils.MockContext, iam.GetAPIClientRequest{
					ClientID:    "unknown",
					Actions:     true,
					GroupAccess: true,
					APIAccess:   true,
					Credentials: true,
					IPACL:       true,
				}).Return(nil, fmt.Errorf("get failed")).Once()
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "unknown",
					ImportState:   true,
					ResourceName:  "akamai_iam_api_client.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable.tf"),
					ExpectError:   regexp.MustCompile(`Importing API Client Resource failed`),
				},
			},
		},
		"expect error - import client with no credentials": {
			importID: "c1ien41d",
			init: func(m *iam.Mock, _, _ testData) {
				m.On("GetAPIClient", testutils.MockContext, iam.GetAPIClientRequest{
					ClientID:    "c1ien41d",
					Actions:     true,
					GroupAccess: true,
					APIAccess:   true,
					Credentials: true,
					IPACL:       true,
				}).Return(&iam.GetAPIClientResponse{
					Credentials: []iam.APIClientCredential{},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "c1ien41d",
					ImportState:   true,
					ResourceName:  "akamai_iam_api_client.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceAPIClient/importable.tf"),
					ExpectError:   regexp.MustCompile(`API client has no credentials`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client, tc.importData, tc.updateData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockCreateAPIClient(m *iam.Mock, testData testData) *mock.Call {
	return m.On("CreateAPIClient", testutils.MockContext, testData.createAPIClientRequest).
		Return(&testData.createAPIClientResponse, nil).Once()
}

func mockUpdateCredential(m *iam.Mock, testData testData) *mock.Call {
	return m.On("UpdateCredential", testutils.MockContext, testData.updateCredentialRequest).
		Return(nil, nil).Once()
}

func mockUpdateAPIClient(m *iam.Mock, testData testData) *mock.Call {
	return m.On("UpdateAPIClient", testutils.MockContext, testData.updateAPIClientRequest).
		Return(&testData.updateAPIClientResponse, nil).Once()
}

// For updating notification emails to an empty list, the response in the production code is disregarded.
func mockUpdateAPIClientNotificationEmails(m *iam.Mock, testData testData) *mock.Call {
	return m.On("UpdateAPIClient", testutils.MockContext, testData.updateAPIClientNotificationEmailsRequest).
		Return(nil, nil).Once()
}

func mockGetAPIClient(m *iam.Mock, testData testData) *mock.Call {
	return m.On("GetAPIClient", testutils.MockContext, testData.getAPIClientRequest).Return(&testData.getAPIClientResponse, nil).Once()
}

func mockListAllowedCPCodes(m *iam.Mock) *mock.Call {
	return m.On("ListAllowedCPCodes", testutils.MockContext, iam.ListAllowedCPCodesRequest{
		UserName: "mw+2",
		Body: iam.ListAllowedCPCodesRequestBody{
			ClientType: "CLIENT",
			Groups: []iam.ClientGroupRequestItem{
				{
					GroupID: 123,
					RoleID:  340,
					Subgroups: []iam.ClientGroupRequestItem{
						{
							GroupID: 333,
							RoleID:  540,
						},
					},
				},
			},
		},
	}).Return(iam.ListAllowedCPCodesResponse{
		{
			Name:  "test",
			Value: 101,
		},
	}, nil).Once()
}

func mockLockAPIClient(m *iam.Mock, data testData) *mock.Call {
	return m.On("LockAPIClient", testutils.MockContext, iam.LockAPIClientRequest{
		ClientID: data.createAPIClientResponse.ClientID,
	}).Return(nil, nil).Once() // Response is ignored in the code
}

func mockUnlockAPIClient(m *iam.Mock, data testData) *mock.Call {
	return m.On("UnlockAPIClient", testutils.MockContext, iam.UnlockAPIClientRequest{
		ClientID: data.createAPIClientResponse.ClientID,
	}).Return(nil, nil).Once() // Response is ignored in the code
}

func mockDeleteAPIClient(m *iam.Mock, testData testData) *mock.Call {
	return m.On("DeleteAPIClient", testutils.MockContext, iam.DeleteAPIClientRequest{
		ClientID: testData.createAPIClientResponse.ClientID,
	}).Return(nil).Once()
}

func mockDeactivateCredential(m *iam.Mock, testData testData) *mock.Call {
	return m.On("DeactivateCredential", testutils.MockContext, iam.DeactivateCredentialRequest{
		ClientID:     testData.createAPIClientResponse.ClientID,
		CredentialID: testData.createAPIClientResponse.Credentials[0].CredentialID,
	}).Return(nil).Once()
}

var (
	apisCreate = []iam.APIRequestItem{
		{
			APIID:       5580,
			AccessLevel: "READ-ONLY",
		},
		{
			APIID:       5801,
			AccessLevel: "READ-WRITE",
		},
	}

	apisGet = []iam.API{
		{
			APIID:            5801,
			APIName:          "EdgeWorkers",
			Description:      "EdgeWorkers",
			Endpoint:         "/edgeworkers/",
			DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
			AccessLevel:      "READ-WRITE",
		},
		{
			APIID:            5580,
			APIName:          "Search Data Feed",
			Description:      "Search Data Feed",
			Endpoint:         "/search-portal-data-feed-api/",
			DocumentationURL: "/",
			AccessLevel:      "READ-ONLY",
		},
	}

	// real list of all APIs is way longer, but we don't need to test all of them
	apisAllGet = []iam.API{
		{
			APIID:            5801,
			APIName:          "EdgeWorkers",
			Description:      "EdgeWorkers",
			Endpoint:         "/edgeworkers/",
			DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
			AccessLevel:      "READ-WRITE",
		},
		{
			APIID:            5580,
			APIName:          "Search Data Feed",
			Description:      "Search Data Feed",
			Endpoint:         "/search-portal-data-feed-api/",
			DocumentationURL: "/",
			AccessLevel:      "READ-ONLY",
		},
		{
			APIID:            6681,
			APIName:          "Test",
			Description:      "Test",
			Endpoint:         "/test",
			DocumentationURL: "-",
			AccessLevel:      "READ-ONLY",
		},
		{
			APIID:            6307,
			APIName:          "CCU APIs",
			Description:      "Content control utility APIs",
			Endpoint:         "/ccu",
			DocumentationURL: "https://developer.akamai.com/",
			AccessLevel:      "READ-WRITE",
		},
		{
			APIID:            6781,
			APIName:          "IDM APIs",
			Description:      "IDM APIs",
			Endpoint:         "/idm",
			DocumentationURL: "blah",
			AccessLevel:      "READ-WRITE",
		},
	}

	purgeOptions = iam.PurgeOptions{
		CanPurgeByCacheTag: true,
		CanPurgeByCPCode:   true,
		CPCodeAccess: iam.CPCodeAccess{
			AllCurrentAndNewCPCodes: false,
			CPCodes:                 []int64{101},
		},
	}

	purgeOptionsNoCPCodes = iam.PurgeOptions{
		CanPurgeByCacheTag: true,
		CanPurgeByCPCode:   true,
		CPCodeAccess: iam.CPCodeAccess{
			AllCurrentAndNewCPCodes: true,
			CPCodes:                 []int64{},
		},
	}

	purgeOptionsCPCodes = iam.PurgeOptions{
		CanPurgeByCacheTag: false,
		CanPurgeByCPCode:   false,
		CPCodeAccess: iam.CPCodeAccess{
			AllCurrentAndNewCPCodes: false,
			CPCodes:                 []int64{101},
		},
	}

	apiAccess = iam.APIAccess{
		AllAccessibleAPIs: true,
		APIs:              nil,
	}

	ipACL = iam.IPACL{
		CIDR:   []string{"128.5.6.5/24"},
		Enable: true,
	}

	ipACLNoCidr = iam.IPACL{
		Enable: false,
	}

	unlockedClientActions = iam.APIClientActions{
		EditGroups:        true,
		EditAPIs:          true,
		Lock:              true,
		Unlock:            false,
		EditAuth:          true,
		Edit:              true,
		EditSwitchAccount: false,
		Transfer:          true,
		EditIPACL:         true,
		Delete:            true,
		DeactivateAll:     false,
	}

	lockedClientActions = iam.APIClientActions{
		EditGroups:        true,
		EditAPIs:          true,
		Lock:              false,
		Unlock:            true,
		EditAuth:          true,
		Edit:              true,
		EditSwitchAccount: false,
		Transfer:          true,
		EditIPACL:         true,
		Delete:            true,
		DeactivateAll:     false,
	}

	singleGroup = []iam.ClientGroup{
		{
			GroupID:         123,
			GroupName:       "group2",
			IsBlocked:       false,
			ParentGroupID:   0,
			RoleDescription: "group description",
			RoleID:          340,
			RoleName:        "role",
			Subgroups:       nil,
		},
	}

	nestedGroups = []iam.ClientGroup{
		{
			GroupID:         123,
			GroupName:       "group2",
			IsBlocked:       false,
			ParentGroupID:   0,
			RoleDescription: "group description",
			RoleID:          340,
			RoleName:        "role",
			Subgroups: []iam.ClientGroup{
				{
					GroupID:         333,
					GroupName:       "group2_1",
					IsBlocked:       false,
					ParentGroupID:   0,
					RoleDescription: "group description",
					RoleID:          540,
					RoleName:        "role 2",
				},
			},
		},
	}

	credentials = []iam.APIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			Status:       "ACTIVE",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Update this credential",
			ExpiresOn:    test.NewTimeFromStringMust("2025-06-13T14:48:08.000Z"),
			Actions: iam.CredentialActions{
				Deactivate:      true,
				Delete:          false,
				Activate:        false,
				EditDescription: true,
				EditExpiration:  true,
			},
		},
	}

	customCredentials = []iam.APIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			Status:       "ACTIVE",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Test API Client Credential",
			ExpiresOn:    test.NewTimeFromStringMust("2026-06-13T14:48:07.000Z"),
			Actions: iam.CredentialActions{
				Deactivate:      true,
				Delete:          false,
				Activate:        false,
				EditDescription: true,
				EditExpiration:  true,
			},
		},
	}

	customCredentials2 = []iam.APIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			Status:       "ACTIVE",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Test API Client Credential 2",
			ExpiresOn:    test.NewTimeFromStringMust("2027-06-13T14:48:07.000Z"),
			Actions: iam.CredentialActions{
				Deactivate:      true,
				Delete:          false,
				Activate:        false,
				EditDescription: true,
				EditExpiration:  true,
			},
		},
	}

	customCredentialsDrift = []iam.APIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			Status:       "DELETED",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Test API Client Credential",
			ExpiresOn:    test.NewTimeFromStringMust("2026-06-13T14:48:07.000Z"),
			Actions: iam.CredentialActions{
				Deactivate:      false,
				Delete:          false,
				Activate:        false,
				EditDescription: false,
				EditExpiration:  false,
			},
		},
	}

	inactiveCredentials = []iam.APIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			Status:       "INACTIVE",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Test API Client Credential",
			ExpiresOn:    test.NewTimeFromStringMust("2025-06-13T14:48:07.000Z"),
			Actions: iam.CredentialActions{
				Deactivate:      false,
				Delete:          true,
				Activate:        true,
				EditDescription: true,
				EditExpiration:  true,
			},
		},
	}

	inactiveCredentials2 = []iam.APIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			Status:       "INACTIVE",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Test API Client Credential 2",
			ExpiresOn:    test.NewTimeFromStringMust("2027-06-13T14:48:07.000Z"),
			Actions: iam.CredentialActions{
				Deactivate:      false,
				Delete:          true,
				Activate:        true,
				EditDescription: true,
				EditExpiration:  true,
			},
		},
	}

	createCredentials = []iam.CreateAPIClientCredential{
		{
			CredentialID: 4444,
			ClientToken:  "token",
			ClientSecret: "secret",
			CreatedOn:    test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
			Description:  "Update this credential",
			ExpiresOn:    test.NewTimeFromStringMust("2025-06-13T14:48:08.000Z"),
			Status:       "ACTIVE",
			Actions: iam.CredentialActions{
				Deactivate:      true,
				Delete:          false,
				Activate:        false,
				EditDescription: true,
				EditExpiration:  true,
			},
		},
	}

	createAPIClientRequest = iam.CreateAPIClientRequest{
		AllowAccountSwitch: false,
		APIAccess: iam.APIAccessRequest{
			AllAccessibleAPIs: false,
			APIs:              apisCreate,
		},
		AuthorizedUsers:         []string{"mw+2"},
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreateCredential:        true,
		GroupAccess: iam.GroupAccessRequest{
			CloneAuthorizedUserGroups: false,
			Groups: []iam.ClientGroupRequestItem{
				{
					GroupID: 123,
					RoleID:  340,
					Subgroups: []iam.ClientGroupRequestItem{
						{
							GroupID: 333,
							RoleID:  540,
						},
					},
				},
			},
		},
		IPACL:              &ipACL,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	createAPIClientResponse = iam.CreateAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             createCredentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	updateCredentialRequest = iam.UpdateCredentialRequest{
		CredentialID: 4444,
		ClientID:     "c1ien41d",
		Body: iam.UpdateCredentialRequestBody{
			Description: "Test API Client Credential",
			ExpiresOn:   test.NewTimeFromStringMust("2026-06-13T14:48:07.000Z"),
			Status:      iam.CredentialActive,
		},
	}

	updateCredentialRequest2 = iam.UpdateCredentialRequest{
		CredentialID: 4444,
		ClientID:     "c1ien41d",
		Body: iam.UpdateCredentialRequestBody{
			Description: "Test API Client Credential 2",
			ExpiresOn:   test.NewTimeFromStringMust("2027-06-13T14:48:07.000Z"),
			Status:      iam.CredentialActive,
		},
	}

	updateInactiveCredentialRequest = iam.UpdateCredentialRequest{
		CredentialID: 4444,
		ClientID:     "c1ien41d",
		Body: iam.UpdateCredentialRequestBody{
			Description: "Test API Client Credential",
			ExpiresOn:   test.NewTimeFromStringMust("2025-06-13T14:48:07.000Z"),
			Status:      iam.CredentialInactive,
		},
	}

	updateInactiveCredentialRequest2 = iam.UpdateCredentialRequest{
		CredentialID: 4444,
		ClientID:     "c1ien41d",
		Body: iam.UpdateCredentialRequestBody{
			Description: "Test API Client Credential 2",
			ExpiresOn:   test.NewTimeFromStringMust("2027-06-13T14:48:07.000Z"),
			Status:      iam.CredentialInactive,
		},
	}

	createAPIClientRequestMin = iam.CreateAPIClientRequest{
		AllowAccountSwitch: false,
		APIAccess: iam.APIAccessRequest{
			AllAccessibleAPIs: false,
			APIs:              apisCreate,
		},
		AuthorizedUsers:         []string{"mw+2"},
		CanAutoCreateCredential: false,
		ClientDescription:       "",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreateCredential:        true,
		GroupAccess: iam.GroupAccessRequest{
			CloneAuthorizedUserGroups: false,
			Groups: []iam.ClientGroupRequestItem{
				{
					GroupID: 123,
					RoleID:  340,
				},
			},
		},
		IPACL:              nil,
		NotificationEmails: []string{},
		PurgeOptions:       nil,
	}

	updateAPIClientRequestMin = iam.UpdateAPIClientRequest{
		ClientID: "c1ien41d",
		Body: iam.UpdateAPIClientRequestBody{
			AllowAccountSwitch: false,
			APIAccess: iam.APIAccessRequest{
				AllAccessibleAPIs: false,
				APIs:              apisCreate,
			},
			AuthorizedUsers:         []string{"mw+2"},
			CanAutoCreateCredential: false,
			ClientDescription:       "",
			ClientName:              "mw+2_1",
			ClientType:              "CLIENT",
			GroupAccess: iam.GroupAccessRequest{
				CloneAuthorizedUserGroups: false,
				Groups: []iam.ClientGroupRequestItem{
					{
						GroupID: 123,
						RoleID:  340,
					},
				},
			},
			IPACL:              nil,
			NotificationEmails: []string{},
			PurgeOptions:       nil,
		},
	}

	updateAPIClientNotificationEmails = iam.UpdateAPIClientRequest{
		ClientID: "c1ien41d",
		Body: iam.UpdateAPIClientRequestBody{
			AllowAccountSwitch: false,
			APIAccess: iam.APIAccessRequest{
				AllAccessibleAPIs: false,
				APIs:              apisCreate,
			},
			AuthorizedUsers:         []string{"mw+2"},
			CanAutoCreateCredential: false,
			ClientDescription:       "",
			ClientName:              "mw+2_1",
			ClientType:              "CLIENT",
			GroupAccess: iam.GroupAccessRequest{
				CloneAuthorizedUserGroups: false,
				Groups: []iam.ClientGroupRequestItem{
					{
						GroupID: 123,
						RoleID:  340,
					},
				},
			},
			IPACL:              nil,
			NotificationEmails: []string{},
			PurgeOptions:       nil,
		},
	}

	createAPIClientResponseMin = iam.CreateAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             createCredentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    singleGroup,
		},
		IPACL:              nil,
		IsLocked:           true,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       nil,
	}

	getAPIClientRequest = iam.GetAPIClientRequest{
		ClientID:    "c1ien41d",
		Actions:     true,
		GroupAccess: true,
		APIAccess:   true,
		Credentials: true,
		IPACL:       true,
	}

	getAPIClientRequestSelf = iam.GetAPIClientRequest{
		ClientID:    "",
		Actions:     true,
		GroupAccess: true,
		APIAccess:   true,
		Credentials: true,
		IPACL:       true,
	}

	getAPIClientResponse = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             credentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseWithCustomCredential = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             customCredentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseWithCustomCredential2 = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             customCredentials2,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseWithCustomCredentialDrift = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             customCredentialsDrift,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseWithCustomCredential2AndAPIClientDescription = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client 2",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             customCredentials2,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseWithInactiveCredential = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 0,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             inactiveCredentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseWithInactiveCredential2 = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 0,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             inactiveCredentials2,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	getAPIClientResponseMin = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             credentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    singleGroup,
		},
		IPACL:              nil,
		IsLocked:           true,
		NotificationEmails: []string{},
		PurgeOptions:       nil,
	}

	updateAPIClientResponseMin = iam.UpdateAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             credentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    singleGroup,
		},
		IPACL:              nil,
		IsLocked:           true,
		NotificationEmails: []string{},
		PurgeOptions:       nil,
	}

	getAPIClientResponseAAA = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: true,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             credentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           true,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	updateAPIClientRequest = iam.UpdateAPIClientRequest{
		ClientID: "c1ien41d",
		Body: iam.UpdateAPIClientRequestBody{
			AllowAccountSwitch: false,
			APIAccess: iam.APIAccessRequest{
				AllAccessibleAPIs: false,
				APIs:              apisCreate,
			},
			AuthorizedUsers:         []string{"mw+2"},
			CanAutoCreateCredential: false,
			ClientDescription:       "Test API Client",
			ClientName:              "mw+2_1",
			ClientType:              "CLIENT",
			GroupAccess: iam.GroupAccessRequest{
				CloneAuthorizedUserGroups: false,
				Groups: []iam.ClientGroupRequestItem{
					{
						GroupID: 123,
						RoleID:  340,
						Subgroups: []iam.ClientGroupRequestItem{
							{
								GroupID: 333,
								RoleID:  540,
							},
						},
					},
				},
			},
			IPACL:              &ipACL,
			NotificationEmails: []string{"mw+2@example.com"},
			PurgeOptions:       &purgeOptions,
		},
	}

	updateAPIClientRequestWithClientDescription = iam.UpdateAPIClientRequest{
		ClientID: "c1ien41d",
		Body: iam.UpdateAPIClientRequestBody{
			AllowAccountSwitch: false,
			APIAccess: iam.APIAccessRequest{
				AllAccessibleAPIs: false,
				APIs:              apisCreate,
			},
			AuthorizedUsers:         []string{"mw+2"},
			CanAutoCreateCredential: false,
			ClientDescription:       "Test API Client 2",
			ClientName:              "mw+2_1",
			ClientType:              "CLIENT",
			GroupAccess: iam.GroupAccessRequest{
				CloneAuthorizedUserGroups: false,
				Groups: []iam.ClientGroupRequestItem{
					{
						GroupID: 123,
						RoleID:  340,
						Subgroups: []iam.ClientGroupRequestItem{
							{
								GroupID: 333,
								RoleID:  540,
							},
						},
					},
				},
			},
			IPACL:              &ipACL,
			NotificationEmails: []string{"mw+2@example.com"},
			PurgeOptions:       &purgeOptions,
		},
	}

	updateAPIClientResponse = iam.UpdateAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             credentials,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	updateAPIClientResponseWithClientDescription = iam.UpdateAPIClientResponse{
		AccessToken:           "access_token",
		Actions:               &unlockedClientActions,
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client 2",
		ClientID:                "c1ien41d",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             test.NewTimeFromStringMust("2024-06-13T14:48:07.000Z"),
		Credentials:             customCredentials2,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    nestedGroups,
		},
		IPACL:              &ipACL,
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		PurgeOptions:       &purgeOptions,
	}

	fullData = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponse,
		updateAPIClientRequest:  updateAPIClientRequest,
		updateAPIClientResponse: updateAPIClientResponse,
	}

	fullDataWithCredential = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		updateCredentialRequest: updateCredentialRequest,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponseWithCustomCredential,
		updateAPIClientRequest:  updateAPIClientRequest,
		updateAPIClientResponse: updateAPIClientResponse,
	}

	fullDataWithCredential2 = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		updateCredentialRequest: updateCredentialRequest2,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponseWithCustomCredential2,
		updateAPIClientRequest:  updateAPIClientRequest,
		updateAPIClientResponse: updateAPIClientResponse,
	}

	fullDataWithCredentialDrift = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		updateCredentialRequest: updateCredentialRequest,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponseWithCustomCredentialDrift,
		updateAPIClientRequest:  updateAPIClientRequest,
		updateAPIClientResponse: updateAPIClientResponse,
	}

	fullDataWithCredential2AndClientDescription = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		updateCredentialRequest: updateCredentialRequest2,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponseWithCustomCredential2AndAPIClientDescription,
		updateAPIClientRequest:  updateAPIClientRequestWithClientDescription,
		updateAPIClientResponse: updateAPIClientResponseWithClientDescription,
	}

	fullDataWithInactiveCredential = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		updateCredentialRequest: updateInactiveCredentialRequest,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponseWithInactiveCredential,
		updateAPIClientRequest:  updateAPIClientRequest,
		updateAPIClientResponse: updateAPIClientResponse,
	}

	fullDataWithInactiveCredential2 = testData{
		createAPIClientRequest:  createAPIClientRequest,
		createAPIClientResponse: createAPIClientResponse,
		updateCredentialRequest: updateInactiveCredentialRequest2,
		getAPIClientRequest:     getAPIClientRequest,
		getAPIClientResponse:    getAPIClientResponseWithInactiveCredential2,
		updateAPIClientRequest:  updateAPIClientRequest,
		updateAPIClientResponse: updateAPIClientResponse,
	}

	minData = testData{
		createAPIClientRequest:                   createAPIClientRequestMin,
		createAPIClientResponse:                  createAPIClientResponseMin,
		updateAPIClientNotificationEmailsRequest: updateAPIClientNotificationEmails,
		getAPIClientRequest:                      getAPIClientRequest,
		getAPIClientResponse:                     getAPIClientResponseMin,
		updateAPIClientRequest:                   updateAPIClientRequestMin,
		updateAPIClientResponse:                  updateAPIClientResponseMin,
	}
)

// it is defined outside due to linter issues
var fullDataChecker = tst.NewStateChecker("akamai_iam_api_client.test").
	CheckEqual("access_token", "access_token").
	CheckEqual("actions.deactivate_all", "false").
	CheckEqual("actions.delete", "true").
	CheckEqual("actions.edit", "true").
	CheckEqual("actions.edit_apis", "true").
	CheckEqual("actions.edit_auth", "true").
	CheckEqual("actions.edit_auth", "true").
	CheckEqual("actions.edit_groups", "true").
	CheckEqual("actions.edit_ip_acl", "true").
	CheckEqual("actions.edit_switch_account", "false").
	CheckEqual("actions.lock", "true").
	CheckEqual("actions.transfer", "true").
	CheckEqual("actions.unlock", "false").
	CheckEqual("active_credential_count", "1").
	CheckEqual("allow_account_switch", "false").
	CheckEqual("api_access.all_accessible_apis", "false").
	CheckEqual("api_access.apis.#", "2").
	CheckEqual("api_access.apis.0.access_level", "READ-ONLY").
	CheckEqual("api_access.apis.0.api_id", "5580").
	CheckEqual("api_access.apis.0.api_name", "Search Data Feed").
	CheckEqual("api_access.apis.0.description", "Search Data Feed").
	CheckEqual("api_access.apis.0.documentation_url", "/").
	CheckEqual("api_access.apis.0.endpoint", "/search-portal-data-feed-api/").
	CheckEqual("api_access.apis.1.access_level", "READ-WRITE").
	CheckEqual("api_access.apis.1.api_id", "5801").
	CheckEqual("api_access.apis.1.api_name", "EdgeWorkers").
	CheckEqual("api_access.apis.1.description", "EdgeWorkers").
	CheckEqual("api_access.apis.1.documentation_url", "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html").
	CheckEqual("api_access.apis.1.endpoint", "/edgeworkers/").
	CheckEqual("authorized_users.0", "mw+2").
	CheckEqual("base_url", "base_url").
	CheckEqual("can_auto_create_credential", "false").
	CheckEqual("client_description", "Test API Client").
	CheckEqual("client_id", "c1ien41d").
	CheckEqual("client_name", "mw+2_1").
	CheckEqual("client_type", "CLIENT").
	CheckEqual("created_by", "someuser").
	CheckEqual("created_date", "2024-06-13T14:48:07Z").
	CheckEqual("credential.actions.activate", "false").
	CheckEqual("credential.actions.deactivate", "true").
	CheckEqual("credential.actions.delete", "false").
	CheckEqual("credential.actions.edit_description", "true").
	CheckEqual("credential.actions.edit_expiration", "true").
	CheckEqual("credential.client_secret", "secret").
	CheckEqual("credential.client_secret", "secret").
	CheckEqual("credential.client_token", "token").
	CheckEqual("credential.created_on", "2024-06-13T14:48:07Z").
	CheckEqual("credential.credential_id", "4444").
	CheckEqual("credential.description", "Update this credential").
	CheckEqual("credential.expires_on", "2025-06-13T14:48:08Z").
	CheckEqual("credential.status", "ACTIVE").
	CheckEqual("group_access.clone_authorized_user_groups", "false").
	CheckEqual("group_access.groups.#", "1").
	CheckEqual("group_access.groups.0.group_id", "123").
	CheckEqual("group_access.groups.0.group_name", "group2").
	CheckEqual("group_access.groups.0.is_blocked", "false").
	CheckEqual("group_access.groups.0.parent_group_id", "0").
	CheckEqual("group_access.groups.0.role_description", "group description").
	CheckEqual("group_access.groups.0.role_id", "340").
	CheckEqual("group_access.groups.0.role_name", "role").
	CheckEqual("id", "c1ien41d").
	CheckEqual("ip_acl.cidr.0", "128.5.6.5/24").
	CheckEqual("ip_acl.enable", "true").
	CheckEqual("lock", "false").
	CheckEqual("notification_emails.0", "mw+2@example.com").
	CheckEqual("purge_options.can_purge_by_cache_tag", "true").
	CheckEqual("purge_options.can_purge_by_cp_code", "true").
	CheckEqual("purge_options.cp_code_access.all_current_and_new_cp_codes", "false").
	CheckEqual("purge_options.cp_code_access.cp_codes.0", "101")

func TestCheckCPCodesAllowed(t *testing.T) {
	tests := []struct {
		name     string
		cpCodes  []int64
		allowed  []iam.ListAllowedCPCodesResponseItem
		expected bool
	}{
		{
			name:     "All codes allowed",
			cpCodes:  []int64{1001, 1002},
			allowed:  []iam.ListAllowedCPCodesResponseItem{{Name: "A", Value: 1001}, {Name: "B", Value: 1002}},
			expected: true,
		},
		{
			name:     "One code not allowed",
			cpCodes:  []int64{1001, 9999},
			allowed:  []iam.ListAllowedCPCodesResponseItem{{Name: "A", Value: 1001}, {Name: "B", Value: 1002}},
			expected: false,
		},
		{
			name:     "Empty cpCodes list",
			cpCodes:  []int64{},
			allowed:  []iam.ListAllowedCPCodesResponseItem{{Name: "A", Value: 1001}},
			expected: true,
		},
		{
			name:     "Empty allowed list",
			cpCodes:  []int64{1001},
			allowed:  []iam.ListAllowedCPCodesResponseItem{},
			expected: false,
		},
		{
			name:     "All codes allowed with more allowed values",
			cpCodes:  []int64{1002},
			allowed:  []iam.ListAllowedCPCodesResponseItem{{Name: "A", Value: 1001}, {Name: "B", Value: 1002}, {Name: "C", Value: 1003}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkAllowedCPCodes(tt.cpCodes, tt.allowed)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
