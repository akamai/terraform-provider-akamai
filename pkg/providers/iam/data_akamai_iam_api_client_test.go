package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/internal/test"
	tst "github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var res = iam.GetAPIClientResponse{
	AccessToken: "access_token",
	Actions: &iam.APIClientActions{
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
	},
	ActiveCredentialCount: 1,
	AllowAccountSwitch:    false,
	APIAccess: iam.APIAccess{
		AllAccessibleAPIs: false,
		APIs:              apisGet,
	},
	AuthorizedUsers:         []string{"ts+2"},
	BaseURL:                 "base_url",
	CanAutoCreateCredential: false,
	ClientDescription:       "Test API Client",
	ClientID:                "c1ien41d",
	ClientName:              "ts+2_1",
	ClientType:              "CLIENT",
	CreatedBy:               "someUser",
	CreatedDate:             test.NewTimeFromStringMust("2025-05-13T14:48:07.000Z"),
	Credentials:             credentials,
	GroupAccess: iam.GroupAccess{
		CloneAuthorizedUserGroups: false,
		Groups:                    nestedGroups,
	},
	IPACL:              &ipACL,
	IsLocked:           false,
	NotificationEmails: []string{"ts+2@example.com"},
	PurgeOptions: &iam.PurgeOptions{
		CanPurgeByCacheTag: true,
		CanPurgeByCPCode:   true,
		CPCodeAccess: iam.CPCodeAccess{
			AllCurrentAndNewCPCodes: false,
			CPCodes:                 []int64{101, 202},
		},
	},
}

func TestDataAPIClient(t *testing.T) {
	tests := map[string]struct {
		init  func(*iam.Mock)
		steps []resource.TestStep
	}{
		"self API client not provided client id": {
			init: func(m *iam.Mock) {
				m.On("GetAPIClient", testutils.MockContext, iam.GetAPIClientRequest{
					ClientID:    "",
					Actions:     true,
					GroupAccess: true,
					APIAccess:   true,
					Credentials: true,
					IPACL:       true,
				}).Return(&res, nil).Times(3)
			},
			steps: getTestStep(t, "testdata/TestDataAPIClient/self.tf"),
		},
		"API client with provided client id": {
			init: func(m *iam.Mock) {
				m.On("GetAPIClient", testutils.MockContext, iam.GetAPIClientRequest{
					ClientID:    "c1ien41d",
					Actions:     true,
					GroupAccess: true,
					APIAccess:   true,
					Credentials: true,
					IPACL:       true,
				}).Return(&res, nil).Times(3)
			},
			steps: getTestStep(t, "testdata/TestDataAPIClient/client.tf"),
		},
		"error - GetAPIClient call failed": {
			init: func(m *iam.Mock) {
				m.On("GetAPIClient", testutils.MockContext, iam.GetAPIClientRequest{
					ClientID:    "c1ien41d",
					Actions:     true,
					GroupAccess: true,
					APIAccess:   true,
					Credentials: true,
					IPACL:       true,
				}).Return(nil, errors.New("test error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataAPIClient/client.tf"),
					ExpectError: regexp.MustCompile("test error"),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func getTestStep(t *testing.T, path string) []resource.TestStep {
	return []resource.TestStep{
		{
			Config: testutils.LoadFixtureString(t, path),
			Check: tst.NewStateChecker("data.akamai_iam_api_client.test").
				CheckEqual("client_id", "c1ien41d").
				CheckEqual("access_token", "access_token").
				CheckEqual("active_credential_count", "1").
				CheckEqual("allow_account_switch", "false").
				CheckEqual("authorized_users.#", "1").
				CheckEqual("authorized_users.0", "ts+2").
				CheckEqual("base_url", "base_url").
				CheckEqual("can_auto_create_credential", "false").
				CheckEqual("client_description", "Test API Client").
				CheckEqual("client_name", "ts+2_1").
				CheckEqual("client_type", "CLIENT").
				CheckEqual("created_by", "someUser").
				CheckEqual("created_date", "2025-05-13T14:48:07Z").
				CheckEqual("is_locked", "false").
				CheckEqual("notification_emails.#", "1").
				CheckEqual("notification_emails.0", "ts+2@example.com").

				// Nested: actions
				CheckEqual("actions.edit_groups", "true").
				CheckEqual("actions.edit_apis", "true").
				CheckEqual("actions.lock", "true").
				CheckEqual("actions.unlock", "false").
				CheckEqual("actions.edit_auth", "true").
				CheckEqual("actions.edit", "true").
				CheckEqual("actions.edit_switch_account", "false").
				CheckEqual("actions.transfer", "true").
				CheckEqual("actions.edit_ip_acl", "true").
				CheckEqual("actions.delete", "true").
				CheckEqual("actions.deactivate_all", "false").

				// Nested: api_access
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

				// Nested: ip_acl
				CheckEqual("ip_acl.enable", "true").
				CheckEqual("ip_acl.cidr.#", "1").
				CheckEqual("ip_acl.cidr.0", "128.5.6.5/24").

				// Nested: purge_options
				CheckEqual("purge_options.can_purge_by_cache_tag", "true").
				CheckEqual("purge_options.can_purge_by_cp_code", "true").
				CheckEqual("purge_options.cp_code_access.all_current_and_new_cp_codes", "false").
				CheckEqual("purge_options.cp_code_access.cp_codes.#", "2").
				CheckEqual("purge_options.cp_code_access.cp_codes.0", "101").
				CheckEqual("purge_options.cp_code_access.cp_codes.1", "202").

				// Nested: credentials
				CheckEqual("credentials.#", "1").
				CheckEqual("credentials.0.credential_id", "4444").
				CheckEqual("credentials.0.client_token", "token").
				CheckEqual("credentials.0.status", "ACTIVE").
				CheckEqual("credentials.0.created_on", "2024-06-13T14:48:07Z").
				CheckEqual("credentials.0.description", "Update this credential").
				CheckEqual("credentials.0.expires_on", "2025-06-13T14:48:08Z").
				CheckEqual("credentials.0.actions.deactivate", "true").
				CheckEqual("credentials.0.actions.delete", "false").
				CheckEqual("credentials.0.actions.activate", "false").
				CheckEqual("credentials.0.actions.edit_description", "true").
				CheckEqual("credentials.0.actions.edit_expiration", "true").

				// Nested: group_access
				CheckEqual("group_access.clone_authorized_user_groups", "false").
				CheckEqual("group_access.groups.0.group_id", "123").
				CheckEqual("group_access.groups.0.group_name", "group2").
				CheckEqual("group_access.groups.0.is_blocked", "false").
				CheckEqual("group_access.groups.0.parent_group_id", "0").
				CheckEqual("group_access.groups.0.role_description", "group description").
				CheckEqual("group_access.groups.0.role_id", "340").
				CheckEqual("group_access.groups.0.role_name", "role").
				CheckEqual("group_access.groups.0.sub_groups.0.group_id", "333").
				CheckEqual("group_access.groups.0.sub_groups.0.group_name", "group2_1").
				CheckEqual("group_access.groups.0.sub_groups.0.is_blocked", "false").
				CheckEqual("group_access.groups.0.sub_groups.0.parent_group_id", "0").
				CheckEqual("group_access.groups.0.sub_groups.0.role_description", "group description").
				CheckEqual("group_access.groups.0.sub_groups.0.role_id", "540").
				CheckEqual("group_access.groups.0.sub_groups.0.role_name", "role 2").
				Build(),
		},
	}
}
