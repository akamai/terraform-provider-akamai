package iam

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceIAMRole(t *testing.T) {
	type roleAttributes struct {
		name, description string
		grantedRoles      []int
	}

	var (
		expectCreateRole = func(t *testing.T, client *mockiam, name, description string, rolesIDs []int) *iam.Role {
			rolesIDsToGrant := getSortedGrantedRolesIDs(rolesIDs)

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

		expectUpdateRole = func(t *testing.T, client *mockiam, id int64, name, description string, rolesIDs []int) *iam.Role {
			rolesIDsToGrant := getSortedGrantedRolesIDs(rolesIDs)

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

		expectReadRole = func(t *testing.T, client *mockiam, roleID int64, name, description string, grantedRoles []iam.RoleGrantedRole, numberOfExecutions int) {
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

		expectDeleteRole = func(t *testing.T, client *mockiam, roleID int64) {
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
		client := new(mockiam)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 2)
		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/role_create.tf", testDir)),
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
		client := new(mockiam)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		updatedRole := expectUpdateRole(t, client, role.RoleID, "role name update", "role description update", []int{12345, 54321, 67890, 1000})
		expectReadRole(t, client, role.RoleID, updatedRole.RoleName, updatedRole.RoleDescription, updatedRole.GrantedRoles, 2)

		expectDeleteRole(t, client, updatedRole.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/role_create.tf", testDir)),
						Check: checkAttributes(roleAttributes{
							name:         "role name",
							description:  "role description",
							grantedRoles: []int{12345, 54321, 67890},
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/role_update.tf", testDir)),
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

	t.Run("import", func(t *testing.T) {
		testDir := "testdata/TestResourceRoleLifecycle"
		client := new(mockiam)
		role := expectCreateRole(t, client, "role name", "role description", []int{12345, 54321, 67890})
		expectReadRole(t, client, role.RoleID, role.RoleName, role.RoleDescription, role.GrantedRoles, 3)

		expectDeleteRole(t, client, role.RoleID)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/role_create.tf", testDir)),
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

func getGrantedRoles(rolesIDsToGrant []iam.GrantedRoleID) []iam.RoleGrantedRole {
	grantedRoles := make([]iam.RoleGrantedRole, 0, len(rolesIDsToGrant))
	for _, roleID := range rolesIDsToGrant {
		grantedRoles = append(grantedRoles, iam.RoleGrantedRole{RoleID: roleID.ID})
	}
	return grantedRoles
}
