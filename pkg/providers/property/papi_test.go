package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/stretchr/testify/mock"
)

type mockpapi struct {
	mock.Mock
}

func (p *mockpapi) GetGroups(ctx context.Context) (*papi.GetGroupsResponse, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetGroupsResponse), args.Error(1)
}

func (p *mockpapi) GetContracts(ctx context.Context) (*papi.GetContractsResponse, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetContractsResponse), args.Error(1)
}

func (p *mockpapi) CreateActivation(ctx context.Context, r papi.CreateActivationRequest) (*papi.CreateActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreateActivationResponse), args.Error(1)
}

func (p *mockpapi) GetActivations(ctx context.Context, r papi.GetActivationsRequest) (*papi.GetActivationsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetActivationsResponse), args.Error(1)
}

func (p *mockpapi) GetActivation(ctx context.Context, r papi.GetActivationRequest) (*papi.GetActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetActivationResponse), args.Error(1)
}

func (p *mockpapi) CancelActivation(ctx context.Context, r papi.CancelActivationRequest) (*papi.CancelActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CancelActivationResponse), args.Error(1)
}

func (p *mockpapi) GetCPCodes(ctx context.Context, r papi.GetCPCodesRequest) (*papi.GetCPCodesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetCPCodesResponse), args.Error(1)
}

func (p *mockpapi) GetCPCode(ctx context.Context, r papi.GetCPCodeRequest) (*papi.GetCPCodesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetCPCodesResponse), args.Error(1)
}

func (p *mockpapi) CreateCPCode(ctx context.Context, r papi.CreateCPCodeRequest) (*papi.CreateCPCodeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreateCPCodeResponse), args.Error(1)
}

func (p *mockpapi) GetProperties(ctx context.Context, r papi.GetPropertiesRequest) (*papi.GetPropertiesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetPropertiesResponse), args.Error(1)
}

func (p *mockpapi) CreateProperty(ctx context.Context, r papi.CreatePropertyRequest) (*papi.CreatePropertyResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreatePropertyResponse), args.Error(1)
}

func (p *mockpapi) GetProperty(ctx context.Context, r papi.GetPropertyRequest) (*papi.GetPropertyResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetPropertyResponse), args.Error(1)
}

func (p *mockpapi) RemoveProperty(ctx context.Context, r papi.RemovePropertyRequest) (*papi.RemovePropertyResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.RemovePropertyResponse), args.Error(1)
}

func (p *mockpapi) GetPropertyVersions(ctx context.Context, r papi.GetPropertyVersionsRequest) (*papi.GetPropertyVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetPropertyVersionsResponse), args.Error(1)
}

func (p *mockpapi) GetPropertyVersion(ctx context.Context, r papi.GetPropertyVersionRequest) (*papi.GetPropertyVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetPropertyVersionsResponse), args.Error(1)
}

func (p *mockpapi) CreatePropertyVersion(ctx context.Context, r papi.CreatePropertyVersionRequest) (*papi.CreatePropertyVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreatePropertyVersionResponse), args.Error(1)
}

func (p *mockpapi) GetLatestVersion(ctx context.Context, r papi.GetLatestVersionRequest) (*papi.GetPropertyVersionsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetPropertyVersionsResponse), args.Error(1)
}

func (p *mockpapi) GetAvailableBehaviors(ctx context.Context, r papi.GetFeaturesRequest) (*papi.GetFeaturesCriteriaResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetFeaturesCriteriaResponse), args.Error(1)
}

func (p *mockpapi) GetAvailableCriteria(ctx context.Context, r papi.GetFeaturesRequest) (*papi.GetFeaturesCriteriaResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetFeaturesCriteriaResponse), args.Error(1)
}

func (p *mockpapi) GetEdgeHostnames(ctx context.Context, r papi.GetEdgeHostnamesRequest) (*papi.GetEdgeHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetEdgeHostnamesResponse), args.Error(1)
}

func (p *mockpapi) GetEdgeHostname(ctx context.Context, r papi.GetEdgeHostnameRequest) (*papi.GetEdgeHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetEdgeHostnamesResponse), args.Error(1)
}

func (p *mockpapi) CreateEdgeHostname(ctx context.Context, r papi.CreateEdgeHostnameRequest) (*papi.CreateEdgeHostnameResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreateEdgeHostnameResponse), args.Error(1)
}

func (p *mockpapi) GetProducts(ctx context.Context, r papi.GetProductsRequest) (*papi.GetProductsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetProductsResponse), args.Error(1)
}

func (p *mockpapi) SearchProperties(ctx context.Context, r papi.SearchRequest) (*papi.SearchResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.SearchResponse), args.Error(1)
}

func (p *mockpapi) GetPropertyVersionHostnames(ctx context.Context, r papi.GetPropertyVersionHostnamesRequest) (*papi.GetPropertyVersionHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetPropertyVersionHostnamesResponse), args.Error(1)
}

func (p *mockpapi) UpdatePropertyVersionHostnames(ctx context.Context, r papi.UpdatePropertyVersionHostnamesRequest) (*papi.UpdatePropertyVersionHostnamesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.UpdatePropertyVersionHostnamesResponse), args.Error(1)
}

func (p *mockpapi) GetClientSettings(ctx context.Context) (*papi.ClientSettingsBody, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.ClientSettingsBody), args.Error(1)
}

func (p *mockpapi) UpdateClientSettings(ctx context.Context, r papi.ClientSettingsBody) (*papi.ClientSettingsBody, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.ClientSettingsBody), args.Error(1)
}

func (p *mockpapi) GetRuleTree(ctx context.Context, r papi.GetRuleTreeRequest) (*papi.GetRuleTreeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetRuleTreeResponse), args.Error(1)
}

func (p *mockpapi) UpdateRuleTree(ctx context.Context, r papi.UpdateRulesRequest) (*papi.UpdateRulesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.UpdateRulesResponse), args.Error(1)
}

func (p *mockpapi) GetRuleFormats(ctx context.Context) (*papi.GetRuleFormatsResponse, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetRuleFormatsResponse), args.Error(1)
}
