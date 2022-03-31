package imaging

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/stretchr/testify/mock"
)

type mockimaging struct {
	mock.Mock
}

func (m *mockimaging) GetPolicy(ctx context.Context, req imaging.GetPolicyRequest) (imaging.PolicyOutput, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(imaging.PolicyOutput), args.Error(1)
}

func (m *mockimaging) UpsertPolicy(ctx context.Context, req imaging.UpsertPolicyRequest) (*imaging.PolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.PolicyResponse), args.Error(1)
}

func (m *mockimaging) DeletePolicy(ctx context.Context, req imaging.DeletePolicyRequest) (*imaging.PolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.PolicyResponse), args.Error(1)
}

func (m *mockimaging) RollbackPolicy(ctx context.Context, req imaging.RollbackPolicyRequest) (*imaging.PolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.PolicyResponse), args.Error(1)
}

func (m *mockimaging) GetPolicyHistory(ctx context.Context, req imaging.GetPolicyHistoryRequest) (*imaging.GetPolicyHistoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.GetPolicyHistoryResponse), args.Error(1)
}

func (m *mockimaging) ListPolicies(ctx context.Context, req imaging.ListPoliciesRequest) (*imaging.ListPoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.ListPoliciesResponse), args.Error(1)
}

func (m *mockimaging) ListPolicySets(ctx context.Context, req imaging.ListPolicySetsRequest) ([]imaging.PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]imaging.PolicySet), args.Error(1)
}

func (m *mockimaging) GetPolicySet(ctx context.Context, req imaging.GetPolicySetRequest) (*imaging.PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.PolicySet), args.Error(1)
}

func (m *mockimaging) CreatePolicySet(ctx context.Context, req imaging.CreatePolicySetRequest) (*imaging.PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.PolicySet), args.Error(1)
}

func (m *mockimaging) UpdatePolicySet(ctx context.Context, req imaging.UpdatePolicySetRequest) (*imaging.PolicySet, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*imaging.PolicySet), args.Error(1)
}

func (m *mockimaging) DeletePolicySet(ctx context.Context, req imaging.DeletePolicySetRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
