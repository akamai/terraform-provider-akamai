package cps

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceThirdPartyEnrollment(t *testing.T) {
	t.Parallel()
	t.Run("lifecycle test", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.AdminContact.FirstName = "R5"
				e.AdminContact.LastName = "D5"
				e.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = "SHA-1"
			}),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(
			WithBase(&enrollmentUpdate),
			WithPendingChangeID(2),
		)
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.0.default", "2h"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/update_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.0.default", "2h"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test, remove san, returns 'wait-review-cert-warning' status", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.AdminContact.FirstName = "R1"
				e.AdminContact.LastName = "D1"
				e.CSR.SANS = nil
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = "SHA-256"
			}),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(
			WithBase(&enrollmentUpdate),
			WithPendingChangeID(2),
		)
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewCertWarning,
			},
		}, nil).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.0.default", "2h"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/empty_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.#", "0"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test update sans add cn", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithEmptySans,
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithSans(commonName),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGetUpdate := newEnrollment(WithBase(&enrollmentUpdate), WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGetUpdate, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle_no_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle_no_sans/update_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment, empty sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithEmptySans, WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/empty_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with empty sans and waiting for deletion", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithEmptySans, WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the first get enrollment call still returns the enrollment.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// Mock that the second get enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/empty_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment, MTLS", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			withMTLS(cps.ClientMutualAuthentication{
				AuthenticationOptions: &cps.AuthenticationOptions{
					OCSP: &cps.OCSP{
						Enabled: ptr.To(true),
					},
					SendCAListToClient: ptr.To(false),
				},
				SetID: "12345",
			}),
		)
		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)

		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentID:              1,
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				AllowCancelPendingChanges: ptr.To(true),
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/3"},
		}, nil).Once()

		enrollmentGet := newEnrollment(
			WithBase(&enrollmentUpdate),
			WithPendingChangeID(3),
		)
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/client_mutual_auth/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test with common name not empty, present in sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithSans(commonName, "san.test.akamai.com"),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test with common name not empty, not present in sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithSans("san.test.akamai.com"),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(
			WithBase(&enrollment),
			WithPendingChangeID(2),
			WithSans(commonName, "san.test.akamai.com"),
		)

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle_no_cn_in_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle_no_cn_in_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("set challenges arrays to empty if no allowedInput found", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getSimpleEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("update with acknowledge warnings change, no enrollment update", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test exclude_sans update", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
		enrollmentReqBody.ThirdParty.ExcludeSANS = false

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.ThirdParty.ExcludeSANS = true
			}),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(
			WithBase(&enrollmentUpdate),
			WithPendingChangeID(2),
		)
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/exclude_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "exclude_sans", "false"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/exclude_sans/update_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "exclude_sans", "true"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)

	})

	t.Run("acknowledge warnings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns 404", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		// GetChangePreVerificationWarnings returns 404: the state suggests warnings exist but the API disagrees.
		// The provider should continue polling rather than fail.
		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(nil, cps.ErrNotFound).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("auto_approve_warnings.#", "0").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("change_management", "false").
						CheckEqual("exclude_sans", "false").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns 404, later actual warns to acknowledge", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		// GetChangePreVerificationWarnings returns 404: the state suggests warnings exist but the API disagrees.
		// The provider should continue polling rather than fail.
		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(nil, cps.ErrNotFound).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Once()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("auto_approve_warnings.#", "0").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("change_management", "false").
						CheckEqual("exclude_sans", "false").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns empty warnings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: ""}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("auto_approve_warnings.#", "0").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("change_management", "false").
						CheckEqual("exclude_sans", "false").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns empty warnings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: ""}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Once()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("auto_approve_warnings.#", "0").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("change_management", "false").
						CheckEqual("exclude_sans", "false").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment, allow duplicate common name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
				AllowDuplicateCN:      true,
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/allow_duplicate_cn/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "allow_duplicate_common_name", "true"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("verification failed with warnings, no acknowledgement", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
					ExpectError: regexp.MustCompile(`enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_verification_warnings flag to true then run 'terraform apply' again: some warning`),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment returns an error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(nil, fmt.Errorf("error creating enrollment")).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
					ExpectError: regexp.MustCompile(`error creating enrollment`),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("auto approve warnings - all warnings on the list to auto approve", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("auto approve warnings - some warnings not on the list to auto approve", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nError parsing expected trust chains.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
					ExpectError: regexp.MustCompile(`warnings cannot be approved: "FIXED_TRUST_CHAIN_PARSING_ERROR"`),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("auto approve warnings - some warnings are unknown", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThis is unknown warning.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
					ExpectError: regexp.MustCompile(`received warning\(s\) does not match any known warning: 'This is unknown warning.'`),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans and clone_dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings_conflict/create_enrollment.tf"),
					ExpectError: regexp.MustCompile("'enable_for_all_sans' and 'clone_dns_names' cannot both be provided at the same time"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans and dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings_conflict/create_enrollment_with_dns_names.tf"),
					ExpectError: regexp.MustCompile("'dns_names' cannot be provided when 'enable_for_all_sans' is true"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans unknown at plan and dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					Source:            "registry.terraform.io/hashicorp/random",
					VersionConstraint: "3.1.0",
				},
			},
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings_conflict/create_enrollment_with_unknown_enable_for_all_sans.tf"),
					ExpectError: regexp.MustCompile("'dns_names' cannot be provided when 'enable_for_all_sans' is true"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with clone_dns_names and dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings_conflict/create_enrollment_with_clone_dns_names.tf"),
					ExpectError: regexp.MustCompile("'dns_names' cannot be provided when 'clone_dns_names' is true"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with clone_dns_names false without dns_names", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings_conflict/create_enrollment_with_false_without_dns_names.tf"),
					ExpectError: regexp.MustCompile("'dns_names' is required when 'enable_for_all_sans' or 'clone_dns_names' is false"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test transitions DNS names settings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		allowCancel := true

		enrollment := newEnrollment(
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: true,
					DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
		assert.Nil(t, enrollmentReqBody.NetworkConfiguration.DNSNameSettings.DNSNames)
		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		currentEnrollment := enrollmentGet
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&currentEnrollment, nil).Times(15)
		mockThirdPartyTransitionChangeStatus(client).Times(4)

		enrollmentWithExplicitDNSNames := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
				}
			}),
		)
		enrollmentWithExplicitDNSNamesReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentWithExplicitDNSNames)
		assert.Equal(t, []string{"test.akamai.com"}, enrollmentWithExplicitDNSNamesReqBody.NetworkConfiguration.DNSNameSettings.DNSNames)
		enrollmentWithExplicitDNSNamesGet := newEnrollment(WithBase(&enrollmentWithExplicitDNSNames), WithPendingChangeID(2))
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentWithExplicitDNSNamesReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Run(func(mock.Arguments) {
			currentEnrollment = enrollmentWithExplicitDNSNamesGet
		}).Once()

		enrollmentWithAdditionalExplicitDNSName := newEnrollment(
			WithBase(&enrollmentWithExplicitDNSNames),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com", "san.test.akamai.com"}
			}),
		)
		enrollmentWithAdditionalExplicitDNSNameReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentWithAdditionalExplicitDNSName)
		enrollmentWithAdditionalExplicitDNSNameGet := newEnrollment(WithBase(&enrollmentWithAdditionalExplicitDNSName), WithPendingChangeID(2))
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentWithAdditionalExplicitDNSNameReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Run(func(mock.Arguments) {
			currentEnrollment = enrollmentWithAdditionalExplicitDNSNameGet
		}).Once()

		enrollmentWithAllSANs := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: true,
					DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
				}
			}),
		)
		enrollmentWithAllSANsReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentWithAllSANs)
		assert.Nil(t, enrollmentWithAllSANsReqBody.NetworkConfiguration.DNSNameSettings.DNSNames)
		enrollmentWithAllSANsGet := newEnrollment(WithBase(&enrollmentWithAllSANs), WithPendingChangeID(2))
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentWithAllSANsReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Run(func(mock.Arguments) {
			currentEnrollment = enrollmentWithAllSANsGet
		}).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_names_transition/enable_for_all_sans.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("network_configuration.0.enable_for_all_sans", "true").
						CheckEqual("network_configuration.0.dns_names.#", "2").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "san.test.akamai.com").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.dns_names.#", "1").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/update_explicit_dns_names.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.dns_names.#", "2").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "san.test.akamai.com").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_names_transition/enable_for_all_sans.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("network_configuration.0.enable_for_all_sans", "true").
						CheckEqual("network_configuration.0.dns_names.#", "2").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "san.test.akamai.com").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("reject enable_for_all_sans mode transition without changing dns_names", func(t *testing.T) {
		t.Parallel()
		testThirdPartyEnrollmentDNSModeTransitionRejected(t,
			testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/update_explicit_dns_names.tf"),
		)
	})

	t.Run("reject clone_dns_names mode transition without changing dns_names", func(t *testing.T) {
		t.Parallel()
		testThirdPartyEnrollmentDNSModeTransitionRejected(t,
			testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/clone_dns_names_unchanged.tf"),
		)
	})

	t.Run("reject DNS name settings for non-SNI enrollment", func(t *testing.T) {
		t.Parallel()
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(edgegrid.NewTestClient(), NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings_conflict/non_sni_with_enable_for_all_sans.tf"),
				ExpectError: regexp.MustCompile("'enable_for_all_sans', 'clone_dns_names', and 'dns_names' cannot be provided when 'sni_only' is false"),
			}},
		})
	})

	t.Run("create enrollment with explicit dns_names", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("network_configuration.#", "1").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.dns_names.#", "1").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("migration: create with clone_dns_names false - no conflict with enable_for_all_sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		enrollment := newEnrollment(
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/clone_dns_names_migration/with_clone_dns_names_false.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.clone_dns_names", "false").
						Build(),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("default: create without dns name settings defaults enable_for_all_sans to true", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		enrollment := newEnrollment(
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: true,
					DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
		assert.Nil(t, enrollmentReqBody.NetworkConfiguration.DNSNameSettings.DNSNames)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/clone_dns_names_migration/without_dns_fields.tf"),
					Check: test.NewStateChecker("akamai_cps_third_party_enrollment.third_party").
						CheckEqual("network_configuration.0.enable_for_all_sans", "true").
						Build(),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})
	t.Run("non-SNI enrollment without DNS name settings has no drift", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := newEnrollment(WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
			e.NetworkConfiguration.DNSNameSettings = nil
			e.NetworkConfiguration.SNIOnly = false
			e.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
			e.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
			e.CSR.SANS = nil
		}))
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
		assert.Nil(t, enrollmentReqBody.NetworkConfiguration.DNSNameSettings)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{ID: 1, Enrollment: "/cps/v2/enrollments/1"}, nil).Once()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		config := strings.Replace(
			testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/clone_dns_names_migration/without_dns_fields.tf"),
			"  sans = [\n    \"san.test.akamai.com\",\n  ]\n",
			"",
			1,
		)
		config = strings.Replace(config, "sni_only       = true", "sni_only       = false", 1)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{Config: config},
				{Config: config, PlanOnly: true},
			},
		})
		client.CPS.AssertExpectations(t)
	})
}

func TestResourceThirdPartyEnrollmentImport(t *testing.T) {
	t.Parallel()
	t.Run("import", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: true,
						DNSNames:      []string{"test.akamai.com"},
					},
					Geography:        "core",
					MustHaveCiphers:  "ak-akamai-2020q1",
					OCSPStapling:     "on",
					PreferredCiphers: "ak-akamai-2020q1",
					QuicEnabled:      false,
					SecureNetwork:    "enhanced-tls",
					SNIOnly:          true,
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_third_party_enrollment.third_party",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						require.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						assert.Equal(t, "true", s[0].Attributes["network_configuration.0.enable_for_all_sans"])
						assert.Equal(t, "1", s[0].Attributes["network_configuration.0.dns_names.#"])
						assert.ElementsMatch(t, []string{"test.akamai.com"}, dnsNamesFromState(rs.Attributes))

						return nil
					},
					// It looks that there bug in SDK that values for bool optional fields are not persisted on create
					ImportStateVerifyIgnore: []string{"network_configuration.0.clone_dns_names", "network_configuration.0.quic_enabled"},
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("import with explicit dns_names", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"
		enrollment := newEnrollment(
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/create_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/create_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_third_party_enrollment.third_party",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						require.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						assert.Equal(t, "false", rs.Attributes["network_configuration.0.enable_for_all_sans"])
						assert.Equal(t, "1", rs.Attributes["network_configuration.0.dns_names.#"])
						assert.ElementsMatch(t, []string{"test.akamai.com"}, dnsNamesFromState(rs.Attributes))
						return nil
					},
					ImportStateVerifyIgnore: []string{"network_configuration.0.quic_enabled"},
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("import with clone_dns_names false", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"
		enrollment := newEnrollment(
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
				}
			}),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/create_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_name_settings/create_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_third_party_enrollment.third_party",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						require.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						assert.Equal(t, "false", rs.Attributes["network_configuration.0.enable_for_all_sans"])
						assert.Equal(t, "false", rs.Attributes["network_configuration.0.clone_dns_names"])
						assert.ElementsMatch(t, []string{"test.akamai.com"}, dnsNamesFromState(rs.Attributes))
						return nil
					},
					ImportStateVerifyIgnore: []string{"network_configuration.0.quic_enabled"},
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("import error when validation type is not third_party", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"

		enrollment := cps.GetEnrollmentResponse{
			ValidationType: "dv",
		}

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(1)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
					ImportState:   true,
					ImportStateId: id,
					ResourceName:  "akamai_cps_third_party_enrollment.third_party",
					ExpectError:   regexp.MustCompile("unable to import: wrong validation type: expected 'third-party', got 'dv'"),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
}

func TestSuppressingSignatureAlgorithm(t *testing.T) {
	t.Parallel()
	t.Run("suppress signature algorithm", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getSimpleEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		enrollmentGet.SignatureAlgorithm = ""

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		enrollmentUpdate := newEnrollment(WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.AdminContact.FirstName = "R5"
				e.AdminContact.LastName = "D5"
				e.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = ""
			}),
		)
		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentUpdateGet := newEnrollment(WithBase(&enrollmentUpdate), WithPendingChangeID(2))
		enrollmentUpdateGet.SignatureAlgorithm = ""

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdateGet, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/update_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
}

func getSimpleEnrollment() cps.GetEnrollmentResponse {
	return cps.GetEnrollmentResponse{
		AdminContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r1d1@akamai.com",
			FirstName:        "R1",
			LastName:         "D1",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		CertificateChainType: "default",
		CertificateType:      "third-party",
		ChangeManagement:     false,
		CSR: &cps.CSR{
			C:    "US",
			CN:   "test.akamai.com",
			L:    "Cambridge",
			O:    "Akamai",
			OU:   "WebEx",
			SANS: []string{"san.test.akamai.com"},
			ST:   "MA",
		},
		EnableMultiStackedCertificates: true,
		NetworkConfiguration: &cps.NetworkConfiguration{
			DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
			DNSNameSettings: &cps.DNSNameSettings{
				CloneDNSNames: false,
				DNSNames:      []string{"test.akamai.com"},
			},
			Geography:        "core",
			MustHaveCiphers:  "ak-akamai-default",
			OCSPStapling:     "on",
			PreferredCiphers: "ak-akamai-default",
			QuicEnabled:      false,
			SecureNetwork:    "enhanced-tls",
			SNIOnly:          true,
		},
		Org: &cps.Org{
			AddressLineOne: "150 Broadway",
			City:           "Cambridge",
			Country:        "US",
			Name:           "Akamai",
			Phone:          "321321321",
			PostalCode:     "12345",
			Region:         "MA",
		},
		RA:                 "third-party",
		SignatureAlgorithm: "SHA-256",
		TechContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r2d2@akamai.com",
			FirstName:        "R2",
			LastName:         "D2",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		ValidationType: "third-party",
		ThirdParty: &cps.ThirdParty{
			ExcludeSANS: false,
		},
	}
}

type enrolOpt interface {
	apply(response *cps.GetEnrollmentResponse)
}

type withPendingChangeID int

func (w withPendingChangeID) apply(e *cps.GetEnrollmentResponse) {
	e.Location = "/cps/v2/enrollments/1"
	e.PendingChanges = []cps.PendingChange{
		{
			Location:   fmt.Sprintf("/cps/v2/enrollments/1/changes/%d", w),
			ChangeType: "new-certificate",
		},
	}
}
func WithPendingChangeID(id int) enrolOpt {
	return withPendingChangeID(id)
}

type withCN string

func (w withCN) apply(e *cps.GetEnrollmentResponse) {
	e.CSR.CN = string(w)
}
func WithCN(cn string) enrolOpt {
	return withCN(cn)
}

type withFunc func(response *cps.GetEnrollmentResponse)

func (w withFunc) apply(e *cps.GetEnrollmentResponse) {
	w(e)
}
func WithUpdateFunc(f func(response *cps.GetEnrollmentResponse)) enrolOpt {
	return withFunc(f)
}

type withMTLS cps.ClientMutualAuthentication

func (w withMTLS) apply(e *cps.GetEnrollmentResponse) {
	e.NetworkConfiguration.ClientMutualAuthentication = (*cps.ClientMutualAuthentication)(&w)
}
func WithMTLS(mtls cps.ClientMutualAuthentication) enrolOpt {
	return withMTLS(mtls)
}

type withBase cps.GetEnrollmentResponse

func (w withBase) apply(e *cps.GetEnrollmentResponse) {
	*e = (cps.GetEnrollmentResponse)(w)
}
func WithBase(e *cps.GetEnrollmentResponse) enrolOpt {
	var newEn cps.Enrollment
	err := copier.CopyWithOption(&newEn, e, copier.Option{DeepCopy: true, IgnoreEmpty: true})
	if err != nil {
		panic(fmt.Sprintln("copier.CopyWithOption failed: ", err))
	}
	return withBase(newEn)
}

type withSans []string

func (w withSans) apply(e *cps.GetEnrollmentResponse) {
	e.CSR.SANS = w
	if e.NetworkConfiguration.DNSNameSettings.CloneDNSNames {
		e.NetworkConfiguration.DNSNameSettings.DNSNames = append([]string{e.CSR.CN}, w...)
	}
}
func WithSans(sans ...string) enrolOpt {
	if len(sans) == 0 {
		return withSans(nil)
	}
	return withSans(sans)
}

var WithEmptySans = WithSans()

func newEnrollment(opts ...enrolOpt) cps.GetEnrollmentResponse {
	enrollment := getSimpleEnrollment()
	for _, o := range opts {
		o.apply(&enrollment)
	}
	return enrollment
}

func testThirdPartyEnrollmentDNSModeTransitionRejected(t *testing.T, config string) {
	t.Helper()

	client := edgegrid.NewTestClient()
	allowCancel := true
	enrollment := newEnrollment(
		WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
			e.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
				CloneDNSNames: true,
				DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
			}
		}),
	)
	enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
	client.CPS.On("CreateEnrollment",
		testutils.MockContext,
		cps.CreateEnrollmentRequest{
			EnrollmentRequestBody: enrollmentReqBody,
			ContractID:            "1",
		},
	).Return(&cps.CreateEnrollmentResponse{
		ID:         1,
		Enrollment: "/cps/v2/enrollments/1",
		Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
	}, nil).Once()

	enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
	client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
		Return(&enrollmentGet, nil).Times(10)
	client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
		Return(nil, cps.ErrNotFound).Once()
	mockThirdPartyTransitionChangeStatus(client).Maybe()
	client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
		EnrollmentID:              1,
		AllowCancelPendingChanges: &allowCancel,
	}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Maybe()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
		Steps: []resource.TestStep{
			{
				Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/dns_names_transition/enable_for_all_sans.tf"),
			},
			{
				Config:      config,
				ExpectError: regexp.MustCompile("'dns_names' must change when switching 'enable_for_all_sans' or 'clone_dns_names' from true to false"),
			},
		},
	})
	client.CPS.AssertExpectations(t)
}

func mockThirdPartyTransitionChangeStatus(client *edgegrid.TestClient) *mock.Call {
	return client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.Change{
		AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
		StatusInfo: &cps.StatusInfo{
			State:  "awaiting-input",
			Status: waitUploadThirdParty,
		},
	}, nil)
}
