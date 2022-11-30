package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/hapi"
	"github.com/stretchr/testify/mock"
)

type mockhapi struct {
	mock.Mock
}

func (m *mockhapi) DeleteEdgeHostname(ctx context.Context, request hapi.DeleteEdgeHostnameRequest) (*hapi.DeleteEdgeHostnameResponse, error) {
	args := m.Called(ctx, request)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*hapi.DeleteEdgeHostnameResponse), nil
}

func (m *mockhapi) GetEdgeHostname(ctx context.Context, id int) (*hapi.GetEdgeHostnameResponse, error) {
	args := m.Called(ctx, id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*hapi.GetEdgeHostnameResponse), nil
}

func (m *mockhapi) UpdateEdgeHostname(ctx context.Context, request hapi.UpdateEdgeHostnameRequest) (*hapi.UpdateEdgeHostnameResponse, error) {
	args := m.Called(ctx, request)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*hapi.UpdateEdgeHostnameResponse), nil
}
