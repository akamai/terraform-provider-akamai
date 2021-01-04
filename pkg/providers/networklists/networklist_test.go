package networklists

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/stretchr/testify/mock"
)

type mocknetworklists struct {
	mock.Mock
}

func (p *mocknetworklists) CreateActivations(ctx context.Context, params networklists.CreateActivationsRequest, acknowledgeWarnings bool) (*networklists.CreateActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.CreateActivationsResponse), args.Error(1)
}

func (p *mocknetworklists) GetActivations(ctx context.Context, params networklists.GetActivationsRequest) (*networklists.GetActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.GetActivationsResponse), args.Error(1)
}

func (p *mocknetworklists) RemoveActivations(ctx context.Context, params networklists.RemoveActivationsRequest) (*networklists.RemoveActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.RemoveActivationsResponse), args.Error(1)
}

func (p *mocknetworklists) CreateNetworkList(ctx context.Context, params networklists.CreateNetworkListRequest) (*networklists.CreateNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.CreateNetworkListResponse), args.Error(1)
}

func (p *mocknetworklists) RemoveNetworkList(ctx context.Context, params networklists.RemoveNetworkListRequest) (*networklists.RemoveNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.RemoveNetworkListResponse), args.Error(1)
}

func (p *mocknetworklists) UpdateNetworkList(ctx context.Context, params networklists.UpdateNetworkListRequest) (*networklists.UpdateNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.UpdateNetworkListResponse), args.Error(1)
}

func (p *mocknetworklists) GetNetworkList(ctx context.Context, params networklists.GetNetworkListRequest) (*networklists.GetNetworkListResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.GetNetworkListResponse), args.Error(1)
}

func (p *mocknetworklists) GetNetworkLists(ctx context.Context, params networklists.GetNetworkListsRequest) (*networklists.GetNetworkListsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.GetNetworkListsResponse), args.Error(1)
}

func (p *mocknetworklists) GetNetworkListDescription(ctx context.Context, params networklists.GetNetworkListDescriptionRequest) (*networklists.GetNetworkListDescriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.GetNetworkListDescriptionResponse), args.Error(1)
}

func (p *mocknetworklists) UpdateNetworkListDescription(ctx context.Context, params networklists.UpdateNetworkListDescriptionRequest) (*networklists.UpdateNetworkListDescriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.UpdateNetworkListDescriptionResponse), args.Error(1)
}

func (p *mocknetworklists) GetNetworkListSubscription(ctx context.Context, params networklists.GetNetworkListSubscriptionRequest) (*networklists.GetNetworkListSubscriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.GetNetworkListSubscriptionResponse), args.Error(1)
}

func (p *mocknetworklists) RemoveNetworkListSubscription(ctx context.Context, params networklists.RemoveNetworkListSubscriptionRequest) (*networklists.RemoveNetworkListSubscriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.RemoveNetworkListSubscriptionResponse), args.Error(1)
}

func (p *mocknetworklists) UpdateNetworkListSubscription(ctx context.Context, params networklists.UpdateNetworkListSubscriptionRequest) (*networklists.UpdateNetworkListSubscriptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*networklists.UpdateNetworkListSubscriptionResponse), args.Error(1)
}
