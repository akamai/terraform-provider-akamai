package iam

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceIAMRole(t *testing.T) {
	type roleAttributes struct {
		name, description string
		grantedRoles      []int
	}

	var (
		updateAPIError = "{\n       \"type\": \"/useradmin-api/error-types/1301\",\n       \"title\": \"Validation Exception\",\n       \"detail\": \"There is a role already with this name\",\n       \"statusCode\": 400,\n       \"httpStatus\": 400\n }\n"
		readAPIError   = "{\n    \"instance\": \"\",\n    \"httpStatus\": 404,\n    \"detail\": \"Role ID not found\",\n    \"title\": \"Role ID not found\",\n    \"type\": \"/useradmin-api/error-types/1311\"\n}"

		expectCreateRole = func(t *testing.T, client *iam.Mock, name, description string, rolesIDs []int) *iam.Role {
			rolesIDsToGrant := getRolesIDsToGrant(rolesIDs)

			roleCreateReq := iam.CreateRoleRequest{
				Name:         name,
				Description:  description,
				GrantedRoles: rolesIDsToGrant,
			}

			grantedRoles := getGrantedRoles(rolesIDsToGrant)

			createdRole := iam.Role{
				RoleID:          123,
				RoleName:        name,
				RoleDescription: description,
				GrantedRoles:    grantedRoles,
			}
			client.On("CreateRole", mock.Anything, roleCreateReq).Return(&createdRole, nil).Once()
			return &createdRole
		}

		expectUpdateRole = func(t *testing.T, client *iam.Mock, id int64, name, description string, rolesIDs []int) *iam.Role {
			rolesIDsToGrant := getRolesIDsToGrant(rolesIDs)

			roleUpdateReq := iam.UpdateRoleRequest{
				ID: id,
				RoleRequest: iam.RoleRequest{
					Name:         name,
					Description:  description,
					GrantedRoles: rolesIDsToGrant,
				},
			}

			grantedRoles := getGrantedRoles(rolesIDsToGrant)

			updatedRole := iam.Role{
				RoleID:          id,
				RoleName:        name,
				RoleDescription: description,
				GrantedRoles:    grantedRoles,
			}
			client.On("UpdateRole", mock.Anything, roleUpdateReq).Return(&updatedRole, nil).Once()
			return &updatedRole
		}

		expectAPIErrorWithUpdateRole = func(t *testing.T, client *iam.Mock, id int64, name, description string, rolesIDs []int) {
			rolesIDsToGrant := getRolesIDsToGrant(rolesIDs)

			roleUpdateReq := iam.UpdateRoleRequest{
				ID: id,
				RoleRequest: iam.RoleRequest{
					Name:         name,
					Description:  description,
					GrantedRoles: rolesIDsToGrant,
				},
			}

			err := fmt.Errorf(updateAPIError)

			client.On("UpdateRole", mock.Anything, roleUpdateReq).Return(nil, err).Once()
		}

		expectReadRole = func(t *testing.T, client *iam.Mock, roleID int64, name, description string, grantedRoles []iam.RoleGrantedRole, numberOfExecutions int) {
			roleGetReq := iam.GetRoleRequest{
				ID:           roleID,
				GrantedRoles: true,
			}

			createdRole := iam.Role{
				RoleID:          roleID,
				RoleName:        name,
				RoleDescription: description,
				GrantedRoles:    grantedRoles,
			}
			client.On("GetRole", mock.Anything, roleGetReq).Return(&createdRole, nil).Times(numberOfExecutions)
		}

		expectReadRoleAPIError = func(t *testing.T, client *iam.Mock, roleID int64) {
			roleGetReq := iam.GetRoleRequest{
				ID:           roleID,
				GrantedRoles: true,
			}
			err := fmt.Errorf(readAPIError)

			client.On("GetRole", mock.Anything, roleGetReq).Return(nil, err).Once()
		}

		expectDeleteRole = func(t *testing.T, client *iam.Mock, roleID int64) {
			roleDeleteReq := iam.DeleteRoleRequest{
				ID: roleID,
			}
			client.On("DeleteRole", mock.Anything, roleDeleteReq).Return(nil, nil).Once()
		}

		checkAttributes = func(attrs roleAttributes) resource.TestCheckFunc {
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_iam_role.role", "name", attrs.name),
				resource.TestCheckResourceAttr("akamai_iam_role.role", "description", attrs.description),
				resource.TestCheckResourceAttr("akamai_iam_role.role", "granted_roles.#", strconv.Itoa(len(attrs.grantedRoles))),
			}
			return resource.ComposeAggregateTestCheckFunc(checks...)
		}
	)

	t.Run("create a new role lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(iam.Mock)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 2)
		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_create.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update a role lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(iam.Mock)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		updatedRole := expectUpdateRole(t, client, role.RoleID, "role name update", "role description update", []int{12345, 54321, 67890, 1000})
		expectReadRole(t, client, role.RoleID, updatedRole.RoleName, updatedRole.RoleDescription, updatedRole.GrantedRoles, 2)

		expectDeleteRole(t, client, updatedRole.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_create.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_update.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name update",
							description:  "role description update",
							grantedRoles: []int{12345, 54321, 67890, 1000},
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("role update is not expected if granted roles are reordered lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(iam.Mock)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 4)

		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_create.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 6789},
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_with_reordered_granted_roles.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 67890, 54321},
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update a role returns an API error lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(iam.Mock)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectAPIErrorWithUpdateRole(t, client, role.RoleID, "role name update", "role description update", []int{12345, 54321, 67890, 1000})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 1)

		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_create.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_update.tf", testDir)),
						ExpectError: regexp.MustCompile(updateAPIError),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update a role returns an API error lifecycle with error in Read", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(iam.Mock)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectAPIErrorWithUpdateRole(t, client, role.RoleID, "role name update", "role description update", []int{12345, 54321, 67890, 1000})
		expectReadRoleAPIError(t, client, role.RoleID)

		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_create.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_update.tf", testDir)),
						ExpectError: regexp.MustCompile(readAPIError),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(iam.Mock)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/role_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     fmt.Sprint(role.RoleID),
						ResourceName:      "akamai_iam_role.role",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func getRolesIDsToGrant(rolesIDs []int) []iam.GrantedRoleID {
	rolesIDsToGrant := make([]iam.GrantedRoleID, 0, len(rolesIDs))
	for _, id := range rolesIDs {
		rolesIDsToGrant = append(rolesIDsToGrant, iam.GrantedRoleID{ID: int64(id)})
	}
	return rolesIDsToGrant
}

func getGrantedRoles(rolesIDsToGrant []iam.GrantedRoleID) []iam.RoleGrantedRole {
	grantedRoles := make([]iam.RoleGrantedRole, 0, len(rolesIDsToGrant))
	for _, roleID := range rolesIDsToGrant {
		grantedRoles = append(grantedRoles, iam.RoleGrantedRole{RoleID: roleID.ID})
	}
	return grantedRoles
}
