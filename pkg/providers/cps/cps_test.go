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

func (m *mockcps) GetChangeManagementInfo(ctx context.Context, r cps.GetChangeRequest) (*cps.ChangeManagementInfoResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.ChangeManagementInfoResponse), args.Error(1)
}

func (m *mockcps) GetChangeDeploymentInfo(ctx context.Context, r cps.GetChangeRequest) (*cps.ChangeDeploymentInfoResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.ChangeDeploymentInfoResponse), args.Error(1)
}

func (m *mockcps) ListDeployments(ctx context.Context, r cps.ListDeploymentsRequest) (*cps.ListDeploymentsResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.ListDeploymentsResponse), args.Error(1)
}

func (m *mockcps) GetProductionDeployment(ctx context.Context, r cps.GetDeploymentRequest) (*cps.GetProductionDeploymentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.GetProductionDeploymentResponse), args.Error(1)
}

func (m *mockcps) GetStagingDeployment(ctx context.Context, r cps.GetDeploymentRequest) (*cps.GetStagingDeploymentResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.GetStagingDeploymentResponse), args.Error(1)
}

func (m *mockcps) GetDeploymentSchedule(ctx context.Context, r cps.GetDeploymentScheduleRequest) (*cps.DeploymentSchedule, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.DeploymentSchedule), args.Error(1)
}

func (m *mockcps) UpdateDeploymentSchedule(ctx context.Context, r cps.UpdateDeploymentScheduleRequest) (*cps.UpdateDeploymentScheduleResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.UpdateDeploymentScheduleResponse), args.Error(1)
}

func (m *mockcps) GetDVHistory(ctx context.Context, r cps.GetDVHistoryRequest) (*cps.GetDVHistoryResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.GetDVHistoryResponse), args.Error(1)
}

func (m *mockcps) GetCertificateHistory(ctx context.Context, r cps.GetCertificateHistoryRequest) (*cps.GetCertificateHistoryResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.GetCertificateHistoryResponse), args.Error(1)
}

func (m *mockcps) GetChangeHistory(ctx context.Context, r cps.GetChangeHistoryRequest) (*cps.GetChangeHistoryResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.GetChangeHistoryResponse), args.Error(1)
}

func (m *mockcps) GetChangePostVerificationWarnings(ctx context.Context, r cps.GetChangeRequest) (*cps.PostVerificationWarnings, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.PostVerificationWarnings), args.Error(1)
}

func (m *mockcps) GetChangeThirdPartyCSR(ctx context.Context, r cps.GetChangeRequest) (*cps.ThirdPartyCSRResponse, error) {
	args := m.Called(ctx, r)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*cps.ThirdPartyCSRResponse), args.Error(1)
}
