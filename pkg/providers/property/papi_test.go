package property

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

type mockpapi struct {
	mock.Mock
}

type (
	// GetGroupsFn is any function having the same signature as papi.GetGroups
	GetGroupsFn func(context.Context) (*papi.GetGroupsResponse, error)

	// GetCPCodesFn is any function having the same signature as papi.GetCPCodes
	GetCPCodesFn func(context.Context, papi.GetCPCodesRequest) (*papi.GetCPCodesResponse, error)

	// CreateCPCodeFn is any function having the same signature as papi.CreateCPCode
	CreateCPCodeFn func(context.Context, papi.CreateCPCodeRequest) (*papi.CreateCPCodeResponse, error)

	// UpdateCPCodeFn is any function having the same signature as papi.UpdateCPCode
	UpdateCPCodeFn func(context.Context, papi.UpdateCPCodeRequest) (*papi.CPCodeDetailResponse, error)

	// GetPropertyFunc is any function having the same signature as papi.GetProperty
	GetPropertyFunc func(context.Context, papi.GetPropertyRequest) (*papi.GetPropertyResponse, error)

	// GetPropertyVersionsFn is any function having the same signature as papi.GetPropertyVersions
	GetPropertyVersionsFn func(context.Context, papi.GetPropertyVersionsRequest) (*papi.GetPropertyVersionsResponse, error)

	// GetPropertyVersionHostnamesFn is any function having the same signature as papi.GetPropertyVersionHostnames
	GetPropertyVersionHostnamesFn func(context.Context, papi.GetPropertyVersionHostnamesRequest) (*papi.GetPropertyVersionHostnamesResponse, error)

	// GetRuleTreeFn is any function having the same signature as papi.GetRuleTree
	GetRuleTreeFn = func(context.Context, papi.GetRuleTreeRequest) (*papi.GetRuleTreeResponse, error)

	// UpdateRuleTreeFnv is any function having the same signature as papi.UpdateRuleTree
	UpdateRuleTreeFn func(context.Context, papi.UpdateRulesRequest) (*papi.UpdateRulesResponse, error)
)

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

func (p *mockpapi) UpdateCPCode(ctx context.Context, r papi.UpdateCPCodeRequest) (*papi.CPCodeDetailResponse, error) {
	args := p.Called(ctx, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*papi.CPCodeDetailResponse), args.Error(1)
}

func (p *mockpapi) GetCPCodeDetail(ctx context.Context, r int) (*papi.CPCodeDetailResponse, error) {
	args := p.Called(ctx, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*papi.CPCodeDetailResponse), args.Error(1)
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

// Expect a call to the mock's papi.GetGroups() where the return value is computed by the given function
func (p *mockpapi) OnGetGroups(ctx interface{}, impl GetGroupsFn) *mock.Call {
	call := p.On("GetGroups", ctx)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)

		call.Return(impl(callCtx))
	})

	return call
}

// Expect a call to the mock's papi.GetCPCodes() where the return value is computed by the given function. The args
// param are used to match calls on the mock as normal. If no args are given, then the expectation matches any calls
// to mock.GetCPCodes()
func (p *mockpapi) OnGetCPCodes(impl GetCPCodesFn, args ...interface{}) *mock.Call {
	var call *mock.Call

	runFn := func(callArgs mock.Arguments) {
		ctx := callArgs.Get(0).(context.Context)
		req := callArgs.Get(1).(papi.GetCPCodesRequest)

		call.Return(impl(ctx, req))
	}

	if len(args) == 0 {
		args = mock.Arguments{AnyCTX, mock.Anything}
	}

	call = p.On("GetCPCodes", args...).Run(runFn)
	return call
}

// Expect a call to the mock's papi.CreateCPCode() where the return value is computed by the given function. The args
// param are used to match calls on the mock as normal. If no args are given, then the expectation matches any calls
// to mock.GetCPCodes()
func (p *mockpapi) OnCreateCPCode(impl CreateCPCodeFn, args ...interface{}) *mock.Call {
	var call *mock.Call

	runFn := func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		req := args.Get(1).(papi.CreateCPCodeRequest)

		call.Return(impl(ctx, req))
	}

	if len(args) == 0 {
		args = mock.Arguments{AnyCTX, mock.Anything}
	}

	call = p.On("CreateCPCode", args...).Run(runFn)
	return call
}

// Expect a call to the mock's papi.UpdateCPCode() where the return value is computed by the given function. The args
// param are used to match calls on the mock as normal. If no args are given, then the expectation matches any calls
// to mock.UpdateCPCode()
func (p *mockpapi) OnUpdateCPCode(impl UpdateCPCodeFn, args ...interface{}) *mock.Call {
	var call *mock.Call

	runFn := func(callArgs mock.Arguments) {
		ctx := callArgs.Get(0).(context.Context)
		req := callArgs.Get(1).(papi.UpdateCPCodeRequest)

		call.Return(impl(ctx, req))
	}

	if len(args) == 0 {
		args = mock.Arguments{AnyCTX, mock.Anything}
	}

	call = p.On("UpdateCPCode", args...).Run(runFn)
	return call
}

// Expect a call to the mock's papi.GetProperty() where the return value is computed by the given function
func (p *mockpapi) OnGetProperty(ctx, req interface{}, impl GetPropertyFunc) *mock.Call {
	call := p.On("GetProperty", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(papi.GetPropertyRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

// Expect a call to the mock's papi.GetPropertyVersions() where the return value is computed by the given
// function
func (p *mockpapi) OnGetPropertyVersions(ctx, req interface{}, impl GetPropertyVersionsFn) *mock.Call {
	call := p.On("GetPropertyVersions", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(papi.GetPropertyVersionsRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

// Expect a call to the mock's papi.GetPropertyVersionHostnames() where the return value is computed by the given
// function
func (p *mockpapi) OnGetPropertyVersionHostnames(ctx, req interface{}, impl GetPropertyVersionHostnamesFn) *mock.Call {
	call := p.On("GetPropertyVersionHostnames", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(papi.GetPropertyVersionHostnamesRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

// Expect a call to the mock's papi.GetPropertyVersionHostnames() where the return value is computed by the given
// function
func (p *mockpapi) OnGetRuleTree(ctx, req interface{}, impl GetRuleTreeFn) *mock.Call {
	call := p.On("GetRuleTree", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(papi.GetRuleTreeRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *mockpapi) OnUpdateRuleTree(ctx, req interface{}, impl UpdateRuleTreeFn) *mock.Call {
	call := p.On("UpdateRuleTree", ctx, req)
	call.Run(func(CallArgs mock.Arguments) {
		callCtx := CallArgs.Get(0).(context.Context)
		callReq := CallArgs.Get(1).(papi.UpdateRulesRequest)

		call.Return(impl(callCtx, callReq))
	})

	return call
}

func (p *mockpapi) ListIncludes(ctx context.Context, r papi.ListIncludesRequest) (*papi.ListIncludesResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.ListIncludesResponse), args.Error(1)
}

func (p *mockpapi) ListIncludeParents(ctx context.Context, r papi.ListIncludeParentsRequest) (*papi.ListIncludeParentsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.ListIncludeParentsResponse), args.Error(1)
}

func (p *mockpapi) GetInclude(ctx context.Context, r papi.GetIncludeRequest) (*papi.GetIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetIncludeResponse), args.Error(1)
}

func (p *mockpapi) CreateInclude(ctx context.Context, r papi.CreateIncludeRequest) (*papi.CreateIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreateIncludeResponse), args.Error(1)
}

func (p *mockpapi) DeleteInclude(ctx context.Context, r papi.DeleteIncludeRequest) (*papi.DeleteIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.DeleteIncludeResponse), args.Error(1)
}

func (p *mockpapi) GetIncludeRuleTree(ctx context.Context, r papi.GetIncludeRuleTreeRequest) (*papi.GetIncludeRuleTreeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.GetIncludeRuleTreeResponse), args.Error(1)
}

func (p *mockpapi) UpdateIncludeRuleTree(ctx context.Context, r papi.UpdateIncludeRuleTreeRequest) (*papi.UpdateIncludeRuleTreeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.UpdateIncludeRuleTreeResponse), args.Error(1)
}

func (p *mockpapi) ActivateInclude(ctx context.Context, r papi.ActivateIncludeRequest) (*papi.ActivationIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.ActivationIncludeResponse), args.Error(1)
}

func (p *mockpapi) DeactivateInclude(ctx context.Context, r papi.DeactivateIncludeRequest) (*papi.DeactivationIncludeResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.DeactivationIncludeResponse), args.Error(1)
}

func (p *mockpapi) GetIncludeActivation(ctx context.Context, r papi.GetIncludeActivationRequest) (*papi.IncludeActivationResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.IncludeActivationResponse), args.Error(1)
}

func (p *mockpapi) ListIncludeActivations(ctx context.Context, r papi.ListIncludeActivationsRequest) (*papi.IncludeActivationsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.IncludeActivationsResponse), args.Error(1)
}

func (p *mockpapi) CreateIncludeVersion(ctx context.Context, r papi.CreateIncludeVersionRequest) (*papi.CreateIncludeVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.CreateIncludeVersionResponse), args.Error(1)
}

func (p *mockpapi) GetIncludeVersion(ctx context.Context, r papi.GetIncludeVersionRequest) (*papi.IncludeVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.IncludeVersionResponse), args.Error(1)
}

func (p *mockpapi) ListIncludeVersions(ctx context.Context, r papi.ListIncludeVersionsRequest) (*papi.IncludeVersionResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.IncludeVersionResponse), args.Error(1)
}

func (p *mockpapi) ListIncludeVersionAvailableCriteria(ctx context.Context, r papi.ListAvailableCriteriaRequest) (*papi.AvailableCriteriaResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.AvailableCriteriaResponse), args.Error(1)
}

func (p *mockpapi) ListIncludeVersionAvailableBehaviors(ctx context.Context, r papi.ListAvailableBehaviorsRequest) (*papi.AvailableBehaviorsResponse, error) {
	args := p.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*papi.AvailableBehaviorsResponse), args.Error(1)
}
