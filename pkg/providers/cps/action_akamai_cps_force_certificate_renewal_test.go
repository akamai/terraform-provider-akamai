package cps

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const commonVerificationFailureDetail = ". The renewal might still be pending. Resolve the reported issue if possible, " +
	"then retry the action with cancel_pending_changes set to true. This instructs CPS to cancel this pending renewal" +
	" before starting a new one."

type invokeInput struct {
	pollInterval                       time.Duration
	contextTimeout                     time.Duration
	enrollmentID                       int64
	cancelPendingChanges               *bool
	acknowledgePreVerificationWarnings *bool
}

type invokeExpected struct {
	errorSummary   string
	errorDetail    string
	warningSummary string
	warningDetail  string
	progress       []string
}

type invokeTestCase struct {
	init     func(client *cps.Mock)
	input    invokeInput
	expected invokeExpected
}

// newTestForceCertificateRenewalActionWithInterval builds an action wired to a mock CPS client via the real
// Configure flow with a custom polling interval.
func newTestForceCertificateRenewalActionWithInterval(t *testing.T, pollInterval time.Duration) (*ForceCertificateRenewalAction, *edgegrid.TestClient) {
	t.Helper()

	client := edgegrid.NewTestClient()
	m, err := meta.New(session.Must(session.New()), hclog.NewNullLogger(), "test")
	require.NoError(t, err)
	m.SetClient(client)

	a := &ForceCertificateRenewalAction{
		pollChangeStatusInterval: pollInterval,
	}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), action.ConfigureRequest{ProviderData: m}, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	return a, client
}

// invokeForceCertificateRenewalWithContext builds the action config from raw tftypes values, invokes the action
// with the given context, and returns the collected progress messages alongside the response.
func invokeForceCertificateRenewalWithContext(ctx context.Context, t *testing.T, a *ForceCertificateRenewalAction, enrollmentID int64, cancelPendingChanges, acknowledgePreVerificationWarnings *bool) (*action.InvokeResponse, []string) {
	t.Helper()

	schemaResp := &action.SchemaResponse{}
	a.Schema(ctx, action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	raw := tftypes.NewValue(objectType, map[string]tftypes.Value{
		"enrollment_id":                         tftypes.NewValue(tftypes.Number, enrollmentID),
		"cancel_pending_changes":                tftypes.NewValue(tftypes.Bool, cancelPendingChanges),
		"acknowledge_pre_verification_warnings": tftypes.NewValue(tftypes.Bool, acknowledgePreVerificationWarnings),
	})

	progress := []string{}
	resp := &action.InvokeResponse{
		SendProgress: func(e action.InvokeProgressEvent) {
			progress = append(progress, e.Message)
		},
	}
	a.Invoke(ctx, action.InvokeRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: raw}}, resp)
	return resp, progress
}

func testDVEnrollment() *cps.GetEnrollmentResponse {
	return &cps.GetEnrollmentResponse{
		AdminContact:         &cps.Contact{FirstName: "R1", LastName: "D1", Email: "r1d1@akamai.com", Phone: "123123123"},
		TechContact:          &cps.Contact{FirstName: "R2", LastName: "D2", Email: "r2d2@akamai.com", Phone: "321321321"},
		CertificateChainType: "default",
		CertificateType:      "san",
		CSR:                  &cps.CSR{CN: "test.example.com", SANS: []string{"test.example.com"}},
		NetworkConfiguration: &cps.NetworkConfiguration{Geography: "core", SecureNetwork: "enhanced-tls"},
		Org:                  &cps.Org{Name: "Akamai"},
		OrgID:                ptr.To(12345),
		RA:                   "lets-encrypt",
		SignatureAlgorithm:   "SHA-256",
		ValidationType:       "dv",
	}
}

func testThirdPartyEnrollment() *cps.GetEnrollmentResponse {
	return &cps.GetEnrollmentResponse{
		AdminContact:                   &cps.Contact{FirstName: "R1", LastName: "D1", Email: "r1d1@akamai.com", Phone: "123123123"},
		TechContact:                    &cps.Contact{FirstName: "R2", LastName: "D2", Email: "r2d2@akamai.com", Phone: "321321321"},
		CertificateChainType:           "default",
		CertificateType:                "third-party",
		CSR:                            &cps.CSR{CN: "third-party.example.com", SANS: []string{"third-party.example.com"}},
		EnableMultiStackedCertificates: true,
		NetworkConfiguration:           &cps.NetworkConfiguration{Geography: "core", SecureNetwork: "enhanced-tls"},
		Org:                            &cps.Org{Name: "Akamai"},
		OrgID:                          ptr.To(54321),
		RA:                             "third-party",
		SignatureAlgorithm:             "SHA-256",
		ThirdParty:                     &cps.ThirdParty{ExcludeSANS: false},
		ValidationType:                 "third-party",
	}
}

func testPendingChange(enrollmentID, changeID int) cps.PendingChange {
	return cps.PendingChange{
		Location:   fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID),
		ChangeType: "renewal",
	}
}

func testEnrollmentWithPendingChanges(en *cps.GetEnrollmentResponse, pendingChanges []cps.PendingChange) *cps.GetEnrollmentResponse {
	updated := *en
	updated.PendingChanges = pendingChanges
	return &updated
}

func testUpdateEnrollmentResponse(enrollmentID, changeID int) *cps.UpdateEnrollmentResponse {
	return &cps.UpdateEnrollmentResponse{
		ID:         enrollmentID,
		Enrollment: fmt.Sprintf("/cps/v2/enrollments/%d", enrollmentID),
		Changes:    []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)},
	}
}

// mockSuccessfulRenewalSequence mocks the GetEnrollment -> UpdateEnrollment -> GetEnrollment calls common to every
// successful force-renewal invocation; callers add any further GetChangeStatus/polling expectations on top.
func mockSuccessfulRenewalSequence(client *cps.Mock, enrollmentID int, enrollment, afterUpdate *cps.GetEnrollmentResponse, changeID int, cancelPendingChanges bool) {
	client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}).
		Return(enrollment, nil).Once()
	client.On("UpdateEnrollment", testutils.MockContext, cps.UpdateEnrollmentRequest{
		EnrollmentRequestBody:     createEnrollmentReqBodyFromEnrollment(*enrollment),
		EnrollmentID:              enrollmentID,
		ForceRenewal:              ptr.To(true),
		AllowCancelPendingChanges: ptr.To(cancelPendingChanges),
	}).Return(testUpdateEnrollmentResponse(enrollmentID, changeID), nil).Once()
	client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}).
		Return(afterUpdate, nil).Once()
}

func TestForceCertificateRenewalActionUnitTests(t *testing.T) {
	t.Parallel()
	t.Run("Invoke", func(t *testing.T) {
		t.Parallel()
		tests := map[string]invokeTestCase{
			"force renewal for a DV enrollment and wait for coordinate-domain-validation": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(100, 5)})
					mockSuccessfulRenewalSequence(client, 100, enrollment, afterUpdate, 5, false)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 100, ChangeID: 5}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: coordinateDomainValidation},
						}, nil).Once()
				},
				input: invokeInput{enrollmentID: 100},
				expected: invokeExpected{
					warningSummary: "Action required: DV certificate",
					warningDetail: "Enrollment 100 is a DV enrollment. To complete the renewal, you might need to complete domain " +
						"validation. You can fetch challenges by using the akamai_cps_dv_enrollment resource or the " +
						"akamai_cps_enrollment data source",
					progress: []string{
						"Forcing certificate renewal for enrollment 100...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 100...",
						"Certificate verification completed for enrollment 100.",
					},
				},
			},
			"force renewal for a third-party enrollment and warn that manual upload is required": {
				init: func(client *cps.Mock) {
					enrollment := testThirdPartyEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(200, 9)})
					mockSuccessfulRenewalSequence(client, 200, enrollment, afterUpdate, 9, false)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 200, ChangeID: 9}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: waitUploadThirdParty},
						}, nil).Once()
				},
				input: invokeInput{enrollmentID: 200},
				expected: invokeExpected{
					warningSummary: "Action required: third-party certificate",
					warningDetail: "Enrollment 200 is a third-party enrollment. To complete the renewal, fetch the CSR via the " +
						"akamai_cps_csr data source, sign it with your CA, and upload the signed certificate via the " +
						"akamai_cps_upload_certificate resource",
					progress: []string{
						"Forcing certificate renewal for enrollment 200...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 200...",
						"Certificate verification completed for enrollment 200.",
					},
				},
			},
			"force renewal for a DV enrollment and acknowledge pre-verification warnings while polling for verification": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(300, 7)})
					mockSuccessfulRenewalSequence(client, 300, enrollment, afterUpdate, 7, false)

					// first poll: non-terminal, does not yet carry the warnings-acknowledgement input
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 300, ChangeID: 7}).
						Return(&cps.Change{StatusInfo: &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks}}, nil).Once()
					// second poll: carries the warnings-acknowledgement input, triggers fetch + acknowledge
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 300, ChangeID: 7}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks},
						}, nil).Once()
					client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{EnrollmentID: 300, ChangeID: 7}).
						Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()
					client.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
						Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
						EnrollmentID:    300,
						ChangeID:        7,
					}).Return(nil).Once()
					// third poll: terminal
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 300, ChangeID: 7}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: coordinateDomainValidation},
						}, nil).Once()
				},
				input: invokeInput{
					enrollmentID:                       300,
					acknowledgePreVerificationWarnings: ptr.To(true),
				},
				expected: invokeExpected{
					warningSummary: "Action required: DV certificate",
					warningDetail: "Enrollment 300 is a DV enrollment. To complete the renewal, you might need to complete domain " +
						"validation. You can fetch challenges by using the akamai_cps_dv_enrollment resource or the " +
						"akamai_cps_enrollment data source",
					progress: []string{
						"Forcing certificate renewal for enrollment 300...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 300...",
						"Certificate verification completed for enrollment 300.",
					},
				},
			},
			"force renewal for a DV enrollment and cancel pending changes when cancel_pending_changes is true": {
				init: func(client *cps.Mock) {
					enrollment := testEnrollmentWithPendingChanges(testDVEnrollment(), []cps.PendingChange{testPendingChange(500, 13)})
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(500, 20)})
					mockSuccessfulRenewalSequence(client, 500, enrollment, afterUpdate, 20, true)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 500, ChangeID: 20}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: coordinateDomainValidation},
						}, nil).Once()
				},
				input: invokeInput{
					enrollmentID:         500,
					cancelPendingChanges: ptr.To(true),
				},
				expected: invokeExpected{
					warningSummary: "Action required: DV certificate",
					warningDetail: "Enrollment 500 is a DV enrollment. To complete the renewal, you might need to complete domain " +
						"validation. You can fetch challenges by using the akamai_cps_dv_enrollment resource or the " +
						"akamai_cps_enrollment data source",
					progress: []string{
						"Forcing certificate renewal for enrollment 500...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 500...",
						"Certificate verification completed for enrollment 500.",
					},
				},
			},
			"fail when enrollment has pending changes and cancel_pending_changes is false": {
				init: func(client *cps.Mock) {
					enrollment := testEnrollmentWithPendingChanges(testDVEnrollment(), []cps.PendingChange{testPendingChange(400, 11)})
					client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 400}).
						Return(enrollment, nil).Once()
				},
				input: invokeInput{
					enrollmentID:         400,
					cancelPendingChanges: ptr.To(false),
				},
				expected: invokeExpected{
					errorSummary: "Forcing certificate renewal failed",
					errorDetail: "Enrollment 400 has pending changes. To proceed set cancel_pending_changes to true " +
						"and invoke the action again.",
				},
			},
			"fail when fetching the enrollment fails": {
				init: func(client *cps.Mock) {
					client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 600}).
						Return(nil, errors.New("enrollment not found")).Once()
				},
				input: invokeInput{enrollmentID: 600},
				expected: invokeExpected{
					errorSummary: "Fetching enrollment failed",
					errorDetail:  "enrollment not found",
				},
			},
			"fail when enrollment validation type is unsupported": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					enrollment.ValidationType = "ov"
					client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 650}).
						Return(enrollment, nil).Once()
				},
				input: invokeInput{enrollmentID: 650},
				expected: invokeExpected{
					errorSummary: "Forcing certificate renewal failed",
					errorDetail:  "unsupported enrollment validation type \"ov\": force renewal supports only \"dv\" and \"third-party\"",
				},
			},
			"fail when the force renewal update request fails": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 700}).
						Return(enrollment, nil).Once()
					client.On("UpdateEnrollment", testutils.MockContext, cps.UpdateEnrollmentRequest{
						EnrollmentRequestBody:     createEnrollmentReqBodyFromEnrollment(*enrollment),
						EnrollmentID:              700,
						ForceRenewal:              ptr.To(true),
						AllowCancelPendingChanges: ptr.To(false),
					}).Return(nil, errors.New("update failed")).Once()
				},
				input: invokeInput{enrollmentID: 700},
				expected: invokeExpected{
					errorSummary: "Forcing certificate renewal failed",
					errorDetail:  "update failed",
					progress:     []string{"Forcing certificate renewal for enrollment 700..."},
				},
			},
			"fail when waiting for certificate verification errors": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(800, 21)})
					mockSuccessfulRenewalSequence(client, 800, enrollment, afterUpdate, 21, false)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 800, ChangeID: 21}).
						Return(nil, errors.New("status failed")).Once()
				},
				input: invokeInput{enrollmentID: 800},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "status failed" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 800...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 800...",
					},
				},
			},
			"fail without cancel_pending_changes guidance when it is enabled": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(850, 23)})
					mockSuccessfulRenewalSequence(client, 850, enrollment, afterUpdate, 23, true)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 850, ChangeID: 23}).
						Return(nil, errors.New("status failed")).Once()
				},
				input: invokeInput{
					enrollmentID:         850,
					cancelPendingChanges: ptr.To(true),
				},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "status failed",
					progress: []string{
						"Forcing certificate renewal for enrollment 850...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 850...",
					},
				},
			},
			"fail when fetching the enrollment for verification fails": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 750}).
						Return(enrollment, nil).Once()
					client.On("UpdateEnrollment", testutils.MockContext, cps.UpdateEnrollmentRequest{
						EnrollmentRequestBody:     createEnrollmentReqBodyFromEnrollment(*enrollment),
						EnrollmentID:              750,
						ForceRenewal:              ptr.To(true),
						AllowCancelPendingChanges: ptr.To(false),
					}).Return(testUpdateEnrollmentResponse(750, 22), nil).Once()
					client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 750}).
						Return(nil, errors.New("verification enrollment lookup failed")).Once()
				},
				input: invokeInput{enrollmentID: 750},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "verification enrollment lookup failed" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 750...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 750...",
					},
				},
			},
			"fail when pre-verification warnings are present and acknowledge_pre_verification_warnings is false": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(1000, 30)})
					mockSuccessfulRenewalSequence(client, 1000, enrollment, afterUpdate, 30, false)

					// first poll: non-terminal, does not yet carry the warnings-acknowledgement input
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1000, ChangeID: 30}).
						Return(&cps.Change{StatusInfo: &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks}}, nil).Once()
					// second poll: carries the warnings-acknowledgement input, but the warnings are never acknowledged
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1000, ChangeID: 30}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks},
						}, nil).Once()
					client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{EnrollmentID: 1000, ChangeID: 30}).
						Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()
				},
				input: invokeInput{
					enrollmentID:                       1000,
					acknowledgePreVerificationWarnings: ptr.To(false),
				},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail: "enrollment pre-verification returned warnings and the enrollment cannot be validated. " +
						"Please fix the issues or set acknowledge_pre_verification_warnings flag to true then run " +
						"'terraform apply' again: some warning" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 1000...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 1000...",
					},
				},
			},
			"fail when acknowledging pre-verification warnings fails": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(1100, 31)})
					mockSuccessfulRenewalSequence(client, 1100, enrollment, afterUpdate, 31, false)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1100, ChangeID: 31}).
						Return(&cps.Change{StatusInfo: &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks}}, nil).Once()
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1100, ChangeID: 31}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks},
						}, nil).Once()
					client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{EnrollmentID: 1100, ChangeID: 31}).
						Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()
					client.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
						Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
						EnrollmentID:    1100,
						ChangeID:        31,
					}).Return(errors.New("acknowledge failed")).Once()
				},
				input: invokeInput{
					enrollmentID:                       1100,
					acknowledgePreVerificationWarnings: ptr.To(true),
				},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "acknowledge failed" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 1100...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 1100...",
					},
				},
			},
			"fail when fetching pre-verification warnings fails": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(1200, 32)})
					mockSuccessfulRenewalSequence(client, 1200, enrollment, afterUpdate, 32, false)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1200, ChangeID: 32}).
						Return(&cps.Change{StatusInfo: &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks}}, nil).Once()
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1200, ChangeID: 32}).
						Return(&cps.Change{
							AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
							StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: waitReviewPreVerificationSafetyChecks},
						}, nil).Once()
					client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{EnrollmentID: 1200, ChangeID: 32}).
						Return(nil, errors.New("could not fetch warnings")).Once()
				},
				input: invokeInput{
					enrollmentID:                       1200,
					acknowledgePreVerificationWarnings: ptr.To(true),
				},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "could not fetch warnings" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 1200...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 1200...",
					},
				},
			},
			"fail when CPS reports a change error while polling": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(1300, 33)})
					mockSuccessfulRenewalSequence(client, 1300, enrollment, afterUpdate, 33, false)

					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1300, ChangeID: 33}).
						Return(&cps.Change{StatusInfo: &cps.StatusInfo{State: "awaiting-input", Status: "running"}}, nil).Once()
					// CPS can report an explicit change error mid-poll even while the overall status is still non-terminal.
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1300, ChangeID: 33}).
						Return(&cps.Change{
							StatusInfo: &cps.StatusInfo{
								State:  "awaiting-input",
								Status: "running",
								Error:  &cps.StatusInfoError{Description: "domain validation failed"},
							},
						}, nil).Once()
				},
				input: invokeInput{enrollmentID: 1300},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "domain validation failed" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 1300...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 1300...",
					},
				},
			},
			"fails when the context deadline is exceeded while polling": {
				init: func(client *cps.Mock) {
					enrollment := testDVEnrollment()
					afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(1500, 35)})
					mockSuccessfulRenewalSequence(client, 1500, enrollment, afterUpdate, 35, false)

					// the poll interval (200ms) is far longer than the context timeout below, so ctx.Done() wins the race.
					client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1500, ChangeID: 35}).
						Return(&cps.Change{StatusInfo: &cps.StatusInfo{State: "awaiting-input", Status: "running"}}, nil).Once()
				},
				input: invokeInput{
					pollInterval:   200 * time.Millisecond,
					contextTimeout: 10 * time.Millisecond,
					enrollmentID:   1500,
				},
				expected: invokeExpected{
					errorSummary: "Waiting for certificate verification failed",
					errorDetail:  "change status context terminated: context deadline exceeded" + commonVerificationFailureDetail,
					progress: []string{
						"Forcing certificate renewal for enrollment 1500...",
						"Successfully sent force renewal request. Waiting for verification for enrollment 1500...",
					},
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				interval := tc.input.pollInterval
				if interval == 0 {
					interval = testPollChangeStatusInterval
				}
				a, client := newTestForceCertificateRenewalActionWithInterval(t, interval)
				if tc.init != nil {
					tc.init(client.CPS)
				}

				ctx := context.Background()
				if tc.input.contextTimeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tc.input.contextTimeout)
					t.Cleanup(cancel)
				}

				resp, progress := invokeForceCertificateRenewalWithContext(ctx, t, a,
					tc.input.enrollmentID, tc.input.cancelPendingChanges, tc.input.acknowledgePreVerificationWarnings)

				if tc.expected.errorSummary != "" {
					require.Len(t, resp.Diagnostics.Errors(), 1)
					diagErr := resp.Diagnostics.Errors()[0]
					assert.Equal(t, tc.expected.errorSummary, diagErr.Summary())
					assert.Equal(t, tc.expected.errorDetail, diagErr.Detail())
				} else {
					assert.False(t, resp.Diagnostics.HasError())
				}

				if tc.expected.warningSummary != "" {
					require.Len(t, resp.Diagnostics.Warnings(), 1)
					warning := resp.Diagnostics.Warnings()[0]
					assert.Equal(t, tc.expected.warningSummary, warning.Summary())
					assert.Equal(t, tc.expected.warningDetail, warning.Detail())
				} else {
					assert.Empty(t, resp.Diagnostics.Warnings())
				}

				if len(tc.expected.progress) == 0 {
					assert.Empty(t, progress)
				} else {
					assert.Equal(t, tc.expected.progress, progress)
				}

				client.CPS.AssertExpectations(t)
			})
		}
	})
}

func TestNewForceCertificateRenewalAction(t *testing.T) {
	t.Parallel()

	pollInterval := 2 * time.Second
	a, ok := NewForceCertificateRenewalAction(pollInterval)().(*ForceCertificateRenewalAction)
	require.True(t, ok)
	assert.Equal(t, pollInterval, a.pollChangeStatusInterval)
}

func TestValidateForceRenewalEnrollmentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		validationType string
		wantError      string
	}{
		{name: "DV enrollment", validationType: "dv"},
		{name: "third-party enrollment", validationType: "third-party"},
		{
			name:           "unsupported EV enrollment",
			validationType: "ev",
			wantError:      "unsupported enrollment validation type \"ev\": force renewal supports only \"dv\" and \"third-party\"",
		},
		{
			name:      "missing validation type",
			wantError: "unsupported enrollment validation type \"\": force renewal supports only \"dv\" and \"third-party\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validateEnrollmentType(test.validationType)
			if test.wantError != "" {
				require.EqualError(t, err, test.wantError)
				return
			}
			require.NoError(t, err)
		})
	}
}

// TestForceCertificateRenewalAction exercises the action end-to-end through the real Terraform CLI, using
// an HCL action block and a lifecycle action_trigger, to complement the Invoke-level unit tests above.
func TestForceCertificateRenewalAction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		configPath  string
		init        func(client *cps.Mock)
		expectError *regexp.Regexp
	}{
		"minimal: only enrollment_id is set": {
			configPath: "testdata/TestForceCertificateRenewalAction/minimal.tf",
			init: func(client *cps.Mock) {
				enrollment := testDVEnrollment()
				afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(2100, 40)})
				mockSuccessfulRenewalSequence(client, 2100, enrollment, afterUpdate, 40, false)
				client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 2100, ChangeID: 40}).
					Return(&cps.Change{
						AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
						StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: coordinateDomainValidation},
					}, nil).Once()
			},
		},
		"full: enrollment_id, cancel_pending_changes and acknowledge_pre_verification_warnings are set": {
			configPath: "testdata/TestForceCertificateRenewalAction/full.tf",
			init: func(client *cps.Mock) {
				enrollment := testEnrollmentWithPendingChanges(testDVEnrollment(), []cps.PendingChange{testPendingChange(2300, 44)})
				afterUpdate := testEnrollmentWithPendingChanges(enrollment, []cps.PendingChange{testPendingChange(2300, 45)})
				mockSuccessfulRenewalSequence(client, 2300, enrollment, afterUpdate, 45, true)
				client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 2300, ChangeID: 45}).
					Return(&cps.Change{
						AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
						StatusInfo:   &cps.StatusInfo{State: "awaiting-input", Status: coordinateDomainValidation},
					}, nil).Once()
			},
		},
		"missing enrollment_id fails plan-time validation": {
			configPath:  "testdata/TestForceCertificateRenewalAction/missing_enrollment_id.tf",
			expectError: regexp.MustCompile(`"enrollment_id" is required`),
		},
		"negative enrollment_id fails plan-time validation": {
			configPath:  "testdata/TestForceCertificateRenewalAction/negative_enrollment_id.tf",
			expectError: regexp.MustCompile(`must be at least 1`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.CPS)
			}

			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_14_0),
				},
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, tc.configPath),
						ExpectError: tc.expectError,
					},
				},
			})

			client.CPS.AssertExpectations(t)
		})
	}
}
