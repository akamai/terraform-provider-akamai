package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceIAMRole(t *testing.T) {
	type roleAttributes struct {
		name, description string
		grantedRoles      []int
	}

	var (
		updateAPIError = "{\n       \"type\": \"/useradmin-api/error-types/1301\",\n       \"title\": \"Validation Exception\",\n       \"detail\": \"There is a role already with this name\",\n       \"statusCode\": 400,\n       \"httpStatus\": 400\n }\n"
		readAPIError   = "{\n    \"instance\": \"\",\n    \"httpStatus\": 404,\n    \"detail\": \"Role ID not found\",\n    \"title\": \"Role ID not found\",\n    \"type\": \"/useradmin-api/error-types/1311\"\n}"

		expectCreateRole = func(client *iam.Mock, name, description string, rolesIDs []int) *iam.Role {
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
			client.On("CreateRole", testutils.MockContext, roleCreateReq).Return(&createdRole, nil).Once()
			return &createdRole
		}

		expectUpdateRole = func(client *iam.Mock, id int64, name, description string, rolesIDs []int) *iam.Role {
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
			client.On("UpdateRole", testutils.MockContext, roleUpdateReq).Return(&updatedRole, nil).Once()
			return &updatedRole
		}

		expectAPIErrorWithUpdateRole = func(client *iam.Mock, id int64, name, description string, rolesIDs []int) {
			rolesIDsToGrant := getRolesIDsToGrant(rolesIDs)

			roleUpdateReq := iam.UpdateRoleRequest{
				ID: id,
				RoleRequest: iam.RoleRequest{
					Name:         name,
					Description:  description,
					GrantedRoles: rolesIDsToGrant,
				},
			}

			err := errors.New(updateAPIError)

			client.On("UpdateRole", testutils.MockContext, roleUpdateReq).Return(nil, err).Once()
		}

		expectReadRole = func(client *iam.Mock, roleID int64, name, description string, grantedRoles []iam.RoleGrantedRole, numberOfExecutions int) {
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
			client.On("GetRole", testutils.MockContext, roleGetReq).Return(&createdRole, nil).Times(numberOfExecutions)
		}

		expectReadRoleAPIError = func(client *iam.Mock, roleID int64) {
			roleGetReq := iam.GetRoleRequest{
				ID:           roleID,
				GrantedRoles: true,
			}
			err := errors.New(readAPIError)

			client.On("GetRole", testutils.MockContext, roleGetReq).Return(nil, err).Once()
		}

		expectDeleteRole = func(client *iam.Mock, roleID int64) {
			roleDeleteReq := iam.DeleteRoleRequest{
				ID: roleID,
			}
			client.On("DeleteRole", testutils.MockContext, roleDeleteReq).Return(nil, nil).Once()
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
		role := expectCreateRole(client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 2)
		expectDeleteRole(client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_create.tf", testDir),
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
		role := expectCreateRole(client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		updatedRole := expectUpdateRole(client, role.RoleID, "role name update", "role description update", []int{12345, 1000, 54321, 67890})
		expectReadRole(client, role.RoleID, updatedRole.RoleName, updatedRole.RoleDescription, updatedRole.GrantedRoles, 2)

		expectDeleteRole(client, updatedRole.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_create.tf", testDir),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_update.tf", testDir),
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
		role := expectCreateRole(client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 4)

		expectDeleteRole(client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_create.tf", testDir),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 6789},
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_with_reordered_granted_roles.tf", testDir),
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
		role := expectCreateRole(client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectAPIErrorWithUpdateRole(client, role.RoleID, "role name update", "role description update", []int{12345, 1000, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 1)

		expectDeleteRole(client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_create.tf", testDir),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/role_update.tf", testDir),
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
		role := expectCreateRole(client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectAPIErrorWithUpdateRole(client, role.RoleID, "role name update", "role description update", []int{12345, 1000, 54321, 67890})
		expectReadRoleAPIError(client, role.RoleID)

		expectDeleteRole(client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_create.tf", testDir),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/role_update.tf", testDir),
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
		role := expectCreateRole(client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectDeleteRole(client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/role_create.tf", testDir),
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
