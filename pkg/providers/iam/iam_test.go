package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/iam"
	"github.com/stretchr/testify/mock"
)

type mockiam struct {
	mock.Mock
}

var _ iam.IAM = &mockiam{}

func (m *mockiam) ListGroups(ctx context.Context, request iam.ListGroupsRequest) ([]iam.Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.Group), args.Error(1)
}

func (m *mockiam) ListRoles(ctx context.Context, request iam.ListRolesRequest) ([]iam.Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.Role), args.Error(1)
}

func (m *mockiam) SupportedCountries(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *mockiam) SupportedContactTypes(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *mockiam) SupportedLanguages(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *mockiam) SupportedTimezones(ctx context.Context) ([]iam.Timezone, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.Timezone), args.Error(1)
}

func (m *mockiam) ListProducts(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *mockiam) ListTimeoutPolicies(ctx context.Context) ([]iam.TimeoutPolicy, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.TimeoutPolicy), args.Error(1)
}

func (m *mockiam) ListStates(ctx context.Context, request iam.ListStatesRequest) ([]string, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (m *mockiam) CreateUser(ctx context.Context, request iam.CreateUserRequest) (*iam.User, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.User), args.Error(1)
}

func (m *mockiam) GetUser(ctx context.Context, request iam.GetUserRequest) (*iam.User, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.User), args.Error(1)
}

func (m *mockiam) UpdateUserInfo(ctx context.Context, request iam.UpdateUserInfoRequest) (*iam.UserBasicInfo, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.UserBasicInfo), args.Error(1)
}

func (m *mockiam) UpdateUserNotifications(ctx context.Context, request iam.UpdateUserNotificationsRequest) (*iam.UserNotifications, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.UserNotifications), args.Error(1)
}

func (m *mockiam) UpdateUserAuthGrants(ctx context.Context, request iam.UpdateUserAuthGrantsRequest) ([]iam.AuthGrant, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.AuthGrant), args.Error(1)
}

func (m *mockiam) RemoveUser(ctx context.Context, request iam.RemoveUserRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) ListBlockedProperties(ctx context.Context, request iam.ListBlockedPropertiesRequest) ([]int64, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]int64), args.Error(1)
}

func (m *mockiam) UpdateBlockedProperties(ctx context.Context, request iam.UpdateBlockedPropertiesRequest) ([]int64, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]int64), args.Error(1)
}

func (m *mockiam) CreateGroup(ctx context.Context, request iam.GroupRequest) (*iam.Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.Group), args.Error(1)
}

func (m *mockiam) GetGroup(ctx context.Context, request iam.GetGroupRequest) (*iam.Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.Group), args.Error(1)
}

func (m *mockiam) ListAffectedUsers(ctx context.Context, request iam.ListAffectedUsersRequest) ([]iam.GroupUser, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.GroupUser), args.Error(1)
}

func (m *mockiam) RemoveGroup(ctx context.Context, request iam.RemoveGroupRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) UpdateGroupName(ctx context.Context, request iam.GroupRequest) (*iam.Group, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.Group), args.Error(1)
}

func (m *mockiam) MoveGroup(ctx context.Context, request iam.MoveGroupRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) CreateRole(ctx context.Context, request iam.CreateRoleRequest) (*iam.Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.Role), args.Error(1)
}

func (m *mockiam) GetRole(ctx context.Context, request iam.GetRoleRequest) (*iam.Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.Role), args.Error(1)
}

func (m *mockiam) UpdateRole(ctx context.Context, request iam.UpdateRoleRequest) (*iam.Role, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.Role), args.Error(1)
}

func (m *mockiam) DeleteRole(ctx context.Context, request iam.DeleteRoleRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) ListGrantableRoles(ctx context.Context) ([]iam.RoleGrantedRole, error) {
	args := m.Called(ctx)
	return args.Get(0).([]iam.RoleGrantedRole), args.Error(1)
}

func (m *mockiam) LockUser(ctx context.Context, request iam.LockUserRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) UnlockUser(ctx context.Context, request iam.UnlockUserRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) ListUsers(ctx context.Context, request iam.ListUsersRequest) ([]iam.UserListItem, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]iam.UserListItem), args.Error(1)
}

func (m *mockiam) UpdateTFA(ctx context.Context, request iam.UpdateTFARequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}

func (m *mockiam) ResetUserPassword(ctx context.Context, request iam.ResetUserPasswordRequest) (*iam.ResetUserPasswordResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*iam.ResetUserPasswordResponse), args.Error(1)
}

func (m *mockiam) SetUserPassword(ctx context.Context, request iam.SetUserPasswordRequest) error {
	args := m.Called(ctx, request)

	return args.Error(0)
}
