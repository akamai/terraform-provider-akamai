package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type (
	testDataForAuthorizedUsers struct {
		authorizedUsers []userData
	}

	userData struct {
		Email        string
		Username     string
		FirstName    string
		LastName     string
		UIIdentityID string
	}
)

var (
	basicTestDataForAuthorizedUsers = testDataForAuthorizedUsers{
		authorizedUsers: []userData{
			{
				Email:        "user1@example.com",
				Username:     "user1",
				FirstName:    "John",
				LastName:     "Doe",
				UIIdentityID: "U-I-DYS45",
			},
			{
				Email:        "user2@example.com",
				Username:     "user2",
				FirstName:    "Jane",
				LastName:     "Smith",
				UIIdentityID: "B-P-2XYC01",
			},
		},
	}
)

func TestAuthorizedUsers(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*iam.Mock, testDataForAuthorizedUsers)
		mockData   testDataForAuthorizedUsers
		error      *regexp.Regexp
	}{
		"success path": {
			configPath: "testdata/TestAuthorizedUsers/default.tf",
			init: func(m *iam.Mock, testData testDataForAuthorizedUsers) {
				expectFullListAuthorizedUsers(m, testData, 3)
			},
			mockData: basicTestDataForAuthorizedUsers,
		},
		"fail path": {
			configPath: "testdata/TestAuthorizedUsers/default.tf",
			init: func(m *iam.Mock, _ testDataForAuthorizedUsers) {
				m.On("ListAuthorizedUsers", testutils.MockContext).Return(iam.ListAuthorizedUsersResponse{}, errors.New("could not get authorized users"))
			},
			error:    regexp.MustCompile(`could not get authorized users`),
			mockData: basicTestDataForAuthorizedUsers,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockData)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.configPath),
							Check:       checkAuthorizedUsersAttrs(tc.mockData),
							ExpectError: tc.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectFullListAuthorizedUsers(client *iam.Mock, data testDataForAuthorizedUsers, timesToRun int) {
	listAuthorizedUsersRes := iam.ListAuthorizedUsersResponse{}

	for _, user := range data.authorizedUsers {
		listAuthorizedUsersRes = append(listAuthorizedUsersRes, iam.AuthorizedUser{
			Email:        user.Email,
			Username:     user.Username,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			UIIdentityID: user.UIIdentityID,
		})
	}

	client.On("ListAuthorizedUsers", testutils.MockContext).Return(listAuthorizedUsersRes, nil).Times(timesToRun)
}

func checkAuthorizedUsersAttrs(data testDataForAuthorizedUsers) resource.TestCheckFunc {
	name := "data.akamai_iam_authorized_users.test"
	checksFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "authorized_users.#", strconv.Itoa(len(data.authorizedUsers))),
	}

	for i, user := range data.authorizedUsers {
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("authorized_users.%d.ui_identity_id", i), user.UIIdentityID))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("authorized_users.%d.email", i), user.Email))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("authorized_users.%d.first_name", i), user.FirstName))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("authorized_users.%d.last_name", i), user.LastName))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("authorized_users.%d.username", i), user.Username))
	}

	return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
}
