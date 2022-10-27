package datastream

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
	"github.com/stretchr/testify/mock"
)

type mockdatastream struct {
	mock.Mock
}

func (m *mockdatastream) CreateStream(ctx context.Context, r datastream.CreateStreamRequest) (*datastream.StreamUpdate, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*datastream.StreamUpdate), args.Error(1)
}

func (m *mockdatastream) GetStream(ctx context.Context, r datastream.GetStreamRequest) (*datastream.DetailedStreamVersion, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*datastream.DetailedStreamVersion), args.Error(1)
}

func (m *mockdatastream) UpdateStream(ctx context.Context, r datastream.UpdateStreamRequest) (*datastream.StreamUpdate, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*datastream.StreamUpdate), args.Error(1)
}

func (m *mockdatastream) DeleteStream(ctx context.Context, r datastream.DeleteStreamRequest) (*datastream.DeleteStreamResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*datastream.DeleteStreamResponse), args.Error(1)
}

func (m *mockdatastream) ListStreams(ctx context.Context, r datastream.ListStreamsRequest) ([]datastream.StreamDetails, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]datastream.StreamDetails), args.Error(1)
}

func (m *mockdatastream) ActivateStream(ctx context.Context, r datastream.ActivateStreamRequest) (*datastream.ActivateStreamResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*datastream.ActivateStreamResponse), args.Error(1)
}

func (m *mockdatastream) DeactivateStream(ctx context.Context, r datastream.DeactivateStreamRequest) (*datastream.DeactivateStreamResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*datastream.DeactivateStreamResponse), args.Error(1)
}

func (m *mockdatastream) GetActivationHistory(ctx context.Context, r datastream.GetActivationHistoryRequest) ([]datastream.ActivationHistoryEntry, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]datastream.ActivationHistoryEntry), args.Error(1)
}

func (m *mockdatastream) GetProperties(ctx context.Context, r datastream.GetPropertiesRequest) ([]datastream.Property, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]datastream.Property), args.Error(1)
}

func (m *mockdatastream) GetPropertiesByGroup(ctx context.Context, r datastream.GetPropertiesByGroupRequest) ([]datastream.Property, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]datastream.Property), args.Error(1)
}

func (m *mockdatastream) GetDatasetFields(ctx context.Context, r datastream.GetDatasetFieldsRequest) ([]datastream.DataSets, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]datastream.DataSets), args.Error(1)
}
