package cloudlets

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/stretchr/testify/mock"
)

type mockcloudlets struct {
	mock.Mock
}

func (m *mockcloudlets) CreateOrigin(_ context.Context, _ cloudlets.LoadBalancerOriginRequest) (*cloudlets.Origin, error) {
	panic("implement me")
}

func (m *mockcloudlets) UpdateOrigin(_ context.Context, _ cloudlets.LoadBalancerOriginRequest) (*cloudlets.Origin, error) {
	panic("implement me")
}

func (m *mockcloudlets) CreateLoadBalancerVersion(_ context.Context, _ cloudlets.CreateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetLoadBalancerVersion(_ context.Context, _ cloudlets.GetLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) UpdateLoadBalancerVersion(_ context.Context, _ cloudlets.UpdateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetLoadBalancerActivations(_ context.Context, _ string) (cloudlets.ActivationsList, error) {
	panic("implement me")
}

func (m *mockcloudlets) ActivateLoadBalancerVersion(_ context.Context, _ cloudlets.ActivateLoadBalancerVersionRequest) (*cloudlets.ActivationResponse, error) {
	panic("implement me")
}

func (m *mockcloudlets) ListPolicyActivations(_ context.Context, _ cloudlets.ListPolicyActivationsRequest) ([]cloudlets.PolicyActivation, error) {
	panic("implement me")
}

func (m *mockcloudlets) ActivatePolicyVersion(_ context.Context, _ cloudlets.ActivatePolicyVersionRequest) error {
	panic("implement me")
}

func (m *mockcloudlets) ListOrigins(_ context.Context, _ cloudlets.ListOriginsRequest) ([]cloudlets.OriginResponse, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetOrigin(_ context.Context, _ string) (*cloudlets.Origin, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetPolicy(_ context.Context, _ int64) (*cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) CreatePolicy(_ context.Context, _ cloudlets.CreatePolicyRequest) (*cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) RemovePolicy(_ context.Context, _ int64) error {
	panic("implement me")
}

func (m *mockcloudlets) UpdatePolicy(_ context.Context, _ cloudlets.UpdatePolicyRequest) (*cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetPolicyVersion(_ context.Context, _ cloudlets.GetPolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) CreatePolicyVersion(_ context.Context, _ cloudlets.CreatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) DeletePolicyVersion(_ context.Context, _ cloudlets.DeletePolicyVersionRequest) error {
	panic("implement me")
}

func (m *mockcloudlets) UpdatePolicyVersion(_ context.Context, _ cloudlets.UpdatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) ListPolicies(_ context.Context, _ cloudlets.ListPoliciesRequest) ([]cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetPolicyProperties(_ context.Context, _ int64) (cloudlets.GetPolicyPropertiesResponse, error) {
	panic("implement me")
}

func (m *mockcloudlets) ListPolicyVersions(_ context.Context, _ cloudlets.ListPolicyVersionsRequest) ([]cloudlets.PolicyVersion, error) {
	panic("implement me")
}
