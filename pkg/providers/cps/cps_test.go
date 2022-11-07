package cps

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/stretchr/testify/mock"
)

type mockcps struct {
	mock.Mock
}

func (m *mockcps) ListEnrollments(ctx context.Context, r cps.ListEnrollmentsRequest) (*cps.ListEnrollmentsResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.ListEnrollmentsResponse), args.Error(1)
}

func (m *mockcps) GetEnrollment(ctx context.Context, r cps.GetEnrollmentRequest) (*cps.Enrollment, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.Enrollment), args.Error(1)
}

func (m *mockcps) CreateEnrollment(ctx context.Context, r cps.CreateEnrollmentRequest) (*cps.CreateEnrollmentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.CreateEnrollmentResponse), args.Error(1)
}

func (m *mockcps) UpdateEnrollment(ctx context.Context, r cps.UpdateEnrollmentRequest) (*cps.UpdateEnrollmentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.UpdateEnrollmentResponse), args.Error(1)
}

func (m *mockcps) RemoveEnrollment(ctx context.Context, r cps.RemoveEnrollmentRequest) (*cps.RemoveEnrollmentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.RemoveEnrollmentResponse), args.Error(1)
}

func (m *mockcps) GetChangeStatus(ctx context.Context, r cps.GetChangeStatusRequest) (*cps.Change, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.Change), args.Error(1)
}

func (m *mockcps) CancelChange(ctx context.Context, r cps.CancelChangeRequest) (*cps.CancelChangeResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.CancelChangeResponse), args.Error(1)
}

func (m *mockcps) UpdateChange(ctx context.Context, r cps.UpdateChangeRequest) (*cps.UpdateChangeResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.UpdateChangeResponse), args.Error(1)
}

func (m *mockcps) GetChangeLetsEncryptChallenges(ctx context.Context, r cps.GetChangeRequest) (*cps.DVArray, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.DVArray), args.Error(1)
}

func (m *mockcps) GetChangePreVerificationWarnings(ctx context.Context, r cps.GetChangeRequest) (*cps.PreVerificationWarnings, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.PreVerificationWarnings), args.Error(1)
}

func (m *mockcps) AcknowledgeDVChallenges(ctx context.Context, r cps.AcknowledgementRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}

func (m *mockcps) AcknowledgePreVerificationWarnings(ctx context.Context, r cps.AcknowledgementRequest) error {
	args := m.Called(ctx, r)

	return args.Error(0)
}
