package cloudlets

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/stretchr/testify/mock"
)

type mockcloudlets struct {
	mock.Mock
}

func (m *mockcloudlets) CreateLoadBalancerVersion(context.Context, cloudlets.CreateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetLoadBalancerVersion(context.Context, cloudlets.GetLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) UpdateLoadBalancerVersion(context.Context, cloudlets.UpdateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetLoadBalancerActivations(context.Context, string) (cloudlets.ActivationsList, error) {
	panic("implement me")
}

func (m *mockcloudlets) ActivateLoadBalancerVersion(context.Context, cloudlets.ActivateLoadBalancerVersionRequest) (*cloudlets.ActivationResponse, error) {
	panic("implement me")
}

func (m *mockcloudlets) ListPolicyActivations(context.Context, cloudlets.ListPolicyActivationsRequest) ([]cloudlets.PolicyActivation, error) {
	panic("implement me")
}

func (m *mockcloudlets) ActivatePolicyVersion(context.Context, cloudlets.ActivatePolicyVersionRequest) error {
	panic("implement me")
}

func (m *mockcloudlets) ListOrigins(context.Context, cloudlets.ListOriginsRequest) (cloudlets.Origins, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetOrigin(context.Context, string) (*cloudlets.Origin, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetPolicy(context.Context, int64) (*cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) CreatePolicy(context.Context, cloudlets.CreatePolicyRequest) (*cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) RemovePolicy(context.Context, int64) error {
	panic("implement me")
}

func (m *mockcloudlets) UpdatePolicy(context.Context, cloudlets.UpdatePolicyRequest) (*cloudlets.Policy, error) {
	panic("implement me")
}

func (m *mockcloudlets) GetPolicyVersion(context.Context, cloudlets.GetPolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) CreatePolicyVersion(context.Context, cloudlets.CreatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	panic("implement me")
}

func (m *mockcloudlets) DeletePolicyVersion(context.Context, cloudlets.DeletePolicyVersionRequest) error {
	panic("implement me")
}

func (m *mockcloudlets) UpdatePolicyVersion(context.Context, cloudlets.UpdatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	panic("implement me")
}
