package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/stretchr/testify/mock"
)

type mockiam struct {
	mock.Mock
}

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
