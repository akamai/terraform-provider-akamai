package iam

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/internal/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestUsersAffectedByMovingGroup(t *testing.T) {

	mockListAffectedUsers := func(client *iam.Mock, userType string, users []iam.GroupUser) *mock.Call {
		return client.On("ListAffectedUsers", testutils.MockContext, iam.ListAffectedUsersRequest{SourceGroupID: 123, DestinationGroupID: 321, UserType: userType}).
			Return(users, nil).Times(3)
	}

	generateExpectedGroup := func(path, accountID, email, firstName, lastName, userName, identityID string, lastLoginDate *string) resource.TestCheckFunc {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.account_id", path), accountID),
			resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.email", path), email),
			resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.first_name", path), firstName),
			resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.last_name", path), lastName),
			resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.ui_username", path), userName),
			resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.ui_identity_id", path), identityID),
		}

		if lastLoginDate != nil {
			checks = append(checks, resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.last_login_date", path), *lastLoginDate))
		} else {
			checks = append(checks, resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", fmt.Sprintf("%s.last_login_date", path), ""))
		}

		return resource.ComposeAggregateTestCheckFunc(checks...)
	}

	tests := map[string]struct {
		init           func(mock *iam.Mock)
		config         string
		expectedError  *regexp.Regexp
		expectedChecks resource.TestCheckFunc
	}{
		"normal case - no filter": {
			init: func(client *iam.Mock) {
				mockListAffectedUsers(client, "", []iam.GroupUser{
					{
						AccountID:  "1-ABCD",
						Email:      "ab.cd@test.com",
						FirstName:  "ab",
						LastName:   "cd",
						UserName:   "abcd",
						IdentityID: "F-CO-abcd",
					},
					{
						AccountID:     "1-EFGH",
						Email:         "ef.gh@test.com",
						FirstName:     "ef",
						LastName:      "gh",
						UserName:      "efgh",
						IdentityID:    "F-CO-efgh",
						LastLoginDate: test.NewTimeFromString(t, "2024-02-06T15:00:00.000Z"),
					},
				})
			},
			config: "testdata/TestDataUsersAffected/basic.tf",
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", "users.#", "2"),
				generateExpectedGroup("users.0", "1-ABCD", "ab.cd@test.com", "ab", "cd", "abcd", "F-CO-abcd", nil),
				generateExpectedGroup("users.1", "1-EFGH", "ef.gh@test.com", "ef", "gh", "efgh", "F-CO-efgh", ptr.To("2024-02-06T15:00:00Z")),
			),
		},
		"normal case - only gained": {
			init: func(client *iam.Mock) {
				mockListAffectedUsers(client, "gainAccess", []iam.GroupUser{
					{
						AccountID:  "1-ABCD",
						Email:      "ab.cd@test.com",
						FirstName:  "ab",
						LastName:   "cd",
						UserName:   "abcd",
						IdentityID: "F-CO-abcd",
					},
				})
			},
			config: "testdata/TestDataUsersAffected/gained.tf",
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", "users.#", "1"),
				generateExpectedGroup("users.0", "1-ABCD", "ab.cd@test.com", "ab", "cd", "abcd", "F-CO-abcd", nil),
			),
		},
		"normal case - only lost": {
			init: func(client *iam.Mock) {
				mockListAffectedUsers(client, "lostAccess", []iam.GroupUser{
					{
						AccountID:     "1-EFGH",
						Email:         "ef.gh@test.com",
						FirstName:     "ef",
						LastName:      "gh",
						UserName:      "efgh",
						IdentityID:    "F-CO-efgh",
						LastLoginDate: test.NewTimeFromString(t, "2024-02-06T15:00:00.000Z"),
					},
				})
			},
			config: "testdata/TestDataUsersAffected/lost.tf",
			expectedChecks: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_iam_users_affected_by_moving_group.test", "users.#", "1"),
				generateExpectedGroup("users.0", "1-EFGH", "ef.gh@test.com", "ef", "gh", "efgh", "F-CO-efgh", ptr.To("2024-02-06T15:00:00Z")),
			),
		},
		"validation - missing source": {
			config:        "testdata/TestDataUsersAffected/missing-source.tf",
			expectedError: regexp.MustCompile("Missing required argument(.|\n)*The argument \"source_group_id\" is required, but no definition was found."),
		},
		"validation - missing destination": {
			config:        "testdata/TestDataUsersAffected/missing-destination.tf",
			expectedError: regexp.MustCompile("Missing required argument(.|\n)*The argument \"destination_group_id\" is required, but no definition was found."),
		},
		"validation - incorrect user_type": {
			config:        "testdata/TestDataUsersAffected/incorrect-usertype.tf",
			expectedError: regexp.MustCompile("Invalid Attribute Value Match(.|\n)*Attribute user_type value must be one of: \\[\"lostAccess\" \"gainAccess\" \"\"]"),
		},
		"api failed": {
			init: func(client *iam.Mock) {
				client.On("ListAffectedUsers", testutils.MockContext, iam.ListAffectedUsersRequest{SourceGroupID: 123, DestinationGroupID: 321, UserType: ""}).
					Return(nil, errors.New("api failed")).Once()
			},
			config:        "testdata/TestDataUsersAffected/basic.tf",
			expectedError: regexp.MustCompile("api failed"),
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
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.config),
							Check:       tc.expectedChecks,
							ExpectError: tc.expectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
