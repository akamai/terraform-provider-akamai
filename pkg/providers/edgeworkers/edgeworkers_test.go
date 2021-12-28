package edgeworkers

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/stretchr/testify/mock"
)

type mockedgeworkers struct {
	mock.Mock
}

func (m *mockedgeworkers) ListActivations(ctx context.Context, req edgeworkers.ListActivationsRequest) (*edgeworkers.ListActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListActivationsResponse), args.Error(1)
}

func (m *mockedgeworkers) GetActivation(ctx context.Context, req edgeworkers.GetActivationRequest) (*edgeworkers.Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Activation), args.Error(1)
}

func (m *mockedgeworkers) ActivateVersion(ctx context.Context, req edgeworkers.ActivateVersionRequest) (*edgeworkers.Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Activation), args.Error(1)
}

func (m *mockedgeworkers) CancelPendingActivation(ctx context.Context, req edgeworkers.CancelActivationRequest) (*edgeworkers.Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Activation), args.Error(1)

}

func (m *mockedgeworkers) ListDeactivations(ctx context.Context, req edgeworkers.ListDeactivationsRequest) (*edgeworkers.ListDeactivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListDeactivationsResponse), args.Error(1)
}

func (m *mockedgeworkers) GetDeactivation(ctx context.Context, req edgeworkers.GetDeactivationRequest) (*edgeworkers.Deactivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Deactivation), args.Error(1)
}

func (m *mockedgeworkers) DeactivateVersion(ctx context.Context, req edgeworkers.DeactivateVersionRequest) (*edgeworkers.Deactivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Deactivation), args.Error(1)
}

func (m *mockedgeworkers) GetEdgeWorkerID(ctx context.Context, req edgeworkers.GetEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockedgeworkers) ListEdgeWorkersID(ctx context.Context, req edgeworkers.ListEdgeWorkersIDRequest) (*edgeworkers.ListEdgeWorkersIDResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListEdgeWorkersIDResponse), args.Error(1)
}

func (m *mockedgeworkers) CreateEdgeWorkerID(ctx context.Context, req edgeworkers.CreateEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockedgeworkers) UpdateEdgeWorkerID(ctx context.Context, req edgeworkers.UpdateEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockedgeworkers) CloneEdgeWorkerID(ctx context.Context, req edgeworkers.CloneEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockedgeworkers) DeleteEdgeWorkerID(ctx context.Context, req edgeworkers.DeleteEdgeWorkerIDRequest) error {
	args := m.Called(ctx, req)
	return args.Error(1)
}

func (m *mockedgeworkers) GetEdgeWorkerVersion(ctx context.Context, req edgeworkers.GetEdgeWorkerVersionRequest) (*edgeworkers.EdgeWorkerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerVersion), args.Error(1)
}

func (m *mockedgeworkers) ListEdgeWorkerVersions(ctx context.Context, req edgeworkers.ListEdgeWorkerVersionsRequest) (*edgeworkers.ListEdgeWorkerVersionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListEdgeWorkerVersionsResponse), args.Error(1)
}

func (m *mockedgeworkers) GetEdgeWorkerVersionContent(ctx context.Context, req edgeworkers.GetEdgeWorkerVersionContentRequest) (*edgeworkers.Bundle, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Bundle), args.Error(1)
}

func (m *mockedgeworkers) CreateEdgeWorkerVersion(ctx context.Context, req edgeworkers.CreateEdgeWorkerVersionRequest) (*edgeworkers.EdgeWorkerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerVersion), args.Error(1)
}

func (m *mockedgeworkers) DeleteEdgeWorkerVersion(ctx context.Context, req edgeworkers.DeleteEdgeWorkerVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(1)
}

func (m *mockedgeworkers) GetPermissionGroup(ctx context.Context, req edgeworkers.GetPermissionGroupRequest) (*edgeworkers.PermissionGroup, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.PermissionGroup), args.Error(1)
}

func (m *mockedgeworkers) ListPermissionGroups(ctx context.Context) (*edgeworkers.ListPermissionGroupsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListPermissionGroupsResponse), args.Error(1)
}

func (m *mockedgeworkers) ListProperties(ctx context.Context, req edgeworkers.ListPropertiesRequest) (*edgeworkers.ListPropertiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListPropertiesResponse), args.Error(1)
}

func (m *mockedgeworkers) ListResourceTiers(ctx context.Context, req edgeworkers.ListResourceTiersRequest) (*edgeworkers.ListResourceTiersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListResourceTiersResponse), args.Error(1)
}

func (m *mockedgeworkers) GetResourceTier(ctx context.Context, req edgeworkers.GetResourceTierRequest) (*edgeworkers.ResourceTier, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ResourceTier), args.Error(1)
}
