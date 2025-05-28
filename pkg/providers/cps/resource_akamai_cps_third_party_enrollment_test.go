package cps

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func TestResourceThirdPartyEnrollment(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := newEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.AdminContact.FirstName = "R5"
				e.AdminContact.LastName = "D5"
				e.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = "SHA-1"
			}),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.On("UpdateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("lifecycle test, remove san, returns 'wait-review-cert-warning' status", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := newEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.AdminContact.FirstName = "R1"
				e.AdminContact.LastName = "D1"
				e.CSR.SANS = nil
				e.NetworkConfiguration.DNSNameSettings.DNSNames = nil
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = "SHA-256"
			}),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.On("UpdateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewCertWarning,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("lifecycle test update sans add cn", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithEmptySans,
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithSans(commonName),
		)

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.On("UpdateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGetUpdate, nil).Times(3)

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("create enrollment, empty sans", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}

		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/empty_sans/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("create enrollment, MTLS", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("UpdateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/client_mutual_auth/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("lifecycle test with common name not empty, present in sans", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithSans(commonName, "san.test.akamai.com"),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})
	t.Run("lifecycle test with common name not empty, not present in sans", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithSans("san.test.akamai.com"),
		)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("set challenges arrays to empty if no allowedInput found", func(t *testing.T) {
		client := &cps.Mock{}
		enrollment := getSimpleEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update with acknowledge warnings change, no enrollment update", func(t *testing.T) {
		client := &cps.Mock{}
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("acknowledge warnings", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("create enrollment, allow duplicate common name", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := newEnrollment(WithEmptySans)
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})

	t.Run("verification failed with warnings, no acknowledgement", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_verification_warnings flag to true then run 'terraform apply' again: some warning`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create enrollment returns an error", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(nil, fmt.Errorf("error creating enrollment")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`error creating enrollment`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("auto approve warnings - all warnings on the list to auto approve", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		client.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("auto approve warnings - some warnings not on the list to auto approve", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nError parsing expected trust chains.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`warnings cannot be approved: "FIXED_TRUST_CHAIN_PARSING_ERROR"`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("auto approve warnings - some warnings are unknown", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThis is unknown warning.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`received warning\(s\) does not match any known warning: 'This is unknown warning.'`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestResourceThirdPartyEnrollmentImport(t *testing.T) {
	t.Run("import", func(t *testing.T) {
		client := &cps.Mock{}
		id := "1,ctr_1"
		enrollment := newEnrollment(
			WithEmptySans,
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
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

		client.On("CreateEnrollment",
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
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Times(3)

		allowCancel := true
		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "1", rs.Attributes["id"])
							return nil
						},
						// It looks that there bug in SDK that values for bool optional fields are not persisted on create
						ImportStateVerifyIgnore: []string{"network_configuration.0.clone_dns_names", "network_configuration.0.quic_enabled"},
					},
				},
			})
		})
	})

	t.Run("import error when validation type is not third_party", func(t *testing.T) {
		client := &cps.Mock{}
		id := "1,ctr_1"

		enrollment := cps.GetEnrollmentResponse{
			ValidationType: "dv",
		}

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})
	})
}

func TestSuppressingSignatureAlgorithm(t *testing.T) {
	t.Run("suppress signature algorithm", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := getSimpleEnrollment()
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.On("CreateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		enrollmentUpdate := newEnrollment(WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.GetEnrollmentResponse) {
				e.AdminContact.FirstName = "R5"
				e.AdminContact.LastName = "D5"
				e.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = ""
			}),
		)
		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.On("UpdateEnrollment",
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

		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdateGet, nil).Times(3)

		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
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
				DNSNames:      []string{"san.test.akamai.com"},
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
	e.NetworkConfiguration.DNSNameSettings.DNSNames = w
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
