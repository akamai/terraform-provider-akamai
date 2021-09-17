package cloudlets

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/stretchr/testify/mock"
)

type mockcloudlets struct {
	mock.Mock
}

func (m *mockcloudlets) CreateLoadBalancerVersion(ctx context.Context, req cloudlets.CreateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockcloudlets) GetLoadBalancerVersion(ctx context.Context, req cloudlets.GetLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockcloudlets) UpdateLoadBalancerVersion(ctx context.Context, req cloudlets.UpdateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockcloudlets) GetLoadBalancerActivations(ctx context.Context, req string) (cloudlets.ActivationsList, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(cloudlets.ActivationsList), args.Error(1)
}

func (m *mockcloudlets) ActivateLoadBalancerVersion(ctx context.Context, req cloudlets.ActivateLoadBalancerVersionRequest) (*cloudlets.ActivationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.ActivationResponse), args.Error(1)
}

func (m *mockcloudlets) ListPolicyActivations(ctx context.Context, req cloudlets.ListPolicyActivationsRequest) ([]cloudlets.PolicyActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.PolicyActivation), args.Error(1)
}

func (m *mockcloudlets) ActivatePolicyVersion(ctx context.Context, req cloudlets.ActivatePolicyVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockcloudlets) ListOrigins(ctx context.Context, req cloudlets.ListOriginsRequest) ([]cloudlets.OriginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.OriginResponse), args.Error(1)
}

func (m *mockcloudlets) GetOrigin(ctx context.Context, req string) (*cloudlets.Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Origin), args.Error(1)
}

func (m *mockcloudlets) CreateOrigin(ctx context.Context, req cloudlets.LoadBalancerOriginCreateRequest) (*cloudlets.Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Origin), args.Error(1)
}

func (m *mockcloudlets) UpdateOrigin(ctx context.Context, req cloudlets.LoadBalancerOriginUpdateRequest) (*cloudlets.Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Origin), args.Error(1)
}

func (m *mockcloudlets) ListPolicies(ctx context.Context, request cloudlets.ListPoliciesRequest) ([]cloudlets.Policy, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.Policy), args.Error(1)
}

func (m *mockcloudlets) GetPolicy(ctx context.Context, policyID int64) (*cloudlets.Policy, error) {
	args := m.Called(ctx, policyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Policy), args.Error(1)
}

func (m *mockcloudlets) CreatePolicy(ctx context.Context, req cloudlets.CreatePolicyRequest) (*cloudlets.Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Policy), args.Error(1)
}

func (m *mockcloudlets) RemovePolicy(ctx context.Context, policyID int64) error {
	args := m.Called(ctx, policyID)
	return args.Error(0)
}

func (m *mockcloudlets) UpdatePolicy(ctx context.Context, req cloudlets.UpdatePolicyRequest) (*cloudlets.Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Policy), args.Error(1)
}

func (m *mockcloudlets) ListPolicyVersions(ctx context.Context, request cloudlets.ListPolicyVersionsRequest) ([]cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockcloudlets) GetPolicyVersion(ctx context.Context, req cloudlets.GetPolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockcloudlets) CreatePolicyVersion(ctx context.Context, req cloudlets.CreatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockcloudlets) DeletePolicyVersion(ctx context.Context, req cloudlets.DeletePolicyVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockcloudlets) UpdatePolicyVersion(ctx context.Context, req cloudlets.UpdatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockcloudlets) GetPolicyProperties(ctx context.Context, req int64) (cloudlets.GetPolicyPropertiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(cloudlets.GetPolicyPropertiesResponse), args.Error(1)
}

func (m *mockcloudlets) ListLoadBalancerVersions(ctx context.Context, req cloudlets.ListLoadBalancerVersionsRequest) ([]cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.LoadBalancerVersion), args.Error(1)
}
