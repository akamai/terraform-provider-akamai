package iam

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataAPIClients(t *testing.T) {
	tests := map[string]struct {
		init  func(*iam.Mock)
		steps []resource.TestStep
	}{
		"happy path": {
			init: func(m *iam.Mock) {
				listAPIClientsResp := iam.ListAPIClientsResponse{
					{
						AccessToken: "test_token1",
						Actions: &iam.ListAPIClientsActions{
							Delete:        true,
							DeactivateAll: false,
							Edit:          true,
							Lock:          false,
							Transfer:      true,
							Unlock:        false,
						},
						ActiveCredentialCount: 123,
						AllowAccountSwitch:    true,
						AuthorizedUsers: []string{
							"jdoe",
						},
						CanAutoCreateCredential: true,
						ClientDescription:       "test",
						ClientID:                "1234",
						ClientName:              "test_name",
						ClientType:              "test_type",
						CreatedBy:               "jdoe",
						CreatedDate:             time.Date(2017, 7, 27, 18, 11, 25, 0, time.UTC),
						IsLocked:                true,
						NotificationEmails: []string{
							"jdoe@example.com",
						},
						ServiceConsumerToken: "token123",
					},
					{
						AccessToken: "test_token2",
						Actions: &iam.ListAPIClientsActions{
							Delete:        false,
							DeactivateAll: true,
							Edit:          false,
							Lock:          true,
							Transfer:      false,
							Unlock:        true,
						},
						ActiveCredentialCount: 123,
						AllowAccountSwitch:    false,
						AuthorizedUsers: []string{
							"jkowalski",
							"jdoe",
						},
						CanAutoCreateCredential: false,
						ClientDescription:       "test",
						ClientID:                "1234",
						ClientName:              "test_name",
						ClientType:              "test_type",
						CreatedBy:               "jkowalski",
						CreatedDate:             time.Date(2017, 7, 27, 18, 11, 25, 0, time.UTC),
						IsLocked:                false,
						NotificationEmails: []string{
							"jkowalski@example.com",
							"jdoe@example.com",
						},
						ServiceConsumerToken: "token1234",
					},
				}
				m.On("ListAPIClients", testutils.MockContext, iam.ListAPIClientsRequest{
					Actions: true,
				}).Return(listAPIClientsResp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataAPIClients/default.tf"),
					Check: test.NewStateChecker("data.akamai_iam_api_clients.test").
						CheckEqual("api_clients.0.access_token", "test_token1").
						CheckEqual("api_clients.0.actions.delete", "true").
						CheckEqual("api_clients.0.actions.deactivate_all", "false").
						CheckEqual("api_clients.0.actions.edit", "true").
						CheckEqual("api_clients.0.actions.lock", "false").
						CheckEqual("api_clients.0.actions.transfer", "true").
						CheckEqual("api_clients.0.actions.unlock", "false").
						CheckEqual("api_clients.0.active_credential_count", "123").
						CheckEqual("api_clients.0.allow_account_switch", "true").
						CheckEqual("api_clients.0.authorized_users.#", "1").
						CheckEqual("api_clients.0.authorized_users.0", "jdoe").
						CheckEqual("api_clients.0.can_auto_create_credential", "true").
						CheckEqual("api_clients.0.client_description", "test").
						CheckEqual("api_clients.0.client_id", "1234").
						CheckEqual("api_clients.0.client_name", "test_name").
						CheckEqual("api_clients.0.client_type", "test_type").
						CheckEqual("api_clients.0.created_by", "jdoe").
						CheckEqual("api_clients.0.created_date", "2017-07-27T18:11:25Z").
						CheckEqual("api_clients.0.is_locked", "true").
						CheckEqual("api_clients.0.notification_emails.#", "1").
						CheckEqual("api_clients.0.notification_emails.0", "jdoe@example.com").
						CheckEqual("api_clients.0.service_consumer_token", "token123").
						CheckEqual("api_clients.1.access_token", "test_token2").
						CheckEqual("api_clients.1.actions.delete", "false").
						CheckEqual("api_clients.1.actions.deactivate_all", "true").
						CheckEqual("api_clients.1.actions.edit", "false").
						CheckEqual("api_clients.1.actions.lock", "true").
						CheckEqual("api_clients.1.actions.transfer", "false").
						CheckEqual("api_clients.1.actions.unlock", "true").
						CheckEqual("api_clients.1.active_credential_count", "123").
						CheckEqual("api_clients.1.allow_account_switch", "false").
						CheckEqual("api_clients.1.authorized_users.#", "2").
						CheckEqual("api_clients.1.authorized_users.0", "jkowalski").
						CheckEqual("api_clients.1.authorized_users.1", "jdoe").
						CheckEqual("api_clients.1.can_auto_create_credential", "false").
						CheckEqual("api_clients.1.client_description", "test").
						CheckEqual("api_clients.1.client_id", "1234").
						CheckEqual("api_clients.1.client_name", "test_name").
						CheckEqual("api_clients.1.client_type", "test_type").
						CheckEqual("api_clients.1.created_by", "jkowalski").
						CheckEqual("api_clients.1.created_date", "2017-07-27T18:11:25Z").
						CheckEqual("api_clients.1.is_locked", "false").
						CheckEqual("api_clients.1.notification_emails.#", "2").
						CheckEqual("api_clients.1.notification_emails.0", "jkowalski@example.com").
						CheckEqual("api_clients.1.notification_emails.1", "jdoe@example.com").
						CheckEqual("api_clients.1.service_consumer_token", "token1234").
						Build(),
				},
			},
		},
		"error - ListAPIClients call failed ": {
			init: func(m *iam.Mock) {
				m.On("ListAPIClients", testutils.MockContext, iam.ListAPIClientsRequest{
					Actions: true,
				}).Return(nil, errors.New("test error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataAPIClients/default.tf"),
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
