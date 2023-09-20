package cps

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourceThirdPartyEnrollment(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := newEnrollment()

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.AdminContact.FirstName = "R5"
				e.AdminContact.LastName = "D5"
				e.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = "SHA-1"
			}),
		)

		allowCancel := true
		client.On("UpdateEnrollment",
			mock.Anything,
			cps.UpdateEnrollmentRequest{
				Enrollment:                enrollmentUpdate,
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
	t.Run("lifecycle test update sans add cn", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		commonName := "test.akamai.com"
		enrollment := newEnrollment(
			WithCN(commonName),
			WithEmptySans,
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		enrollmentUpdate := newEnrollment(
			WithBase(&enrollment),
			WithSans(commonName),
		)

		allowCancel := true
		client.On("UpdateEnrollment",
			mock.Anything,
			cps.UpdateEnrollmentRequest{
				Enrollment:                enrollmentUpdate,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGetUpdate := newEnrollment(WithBase(&enrollmentUpdate), WithPendingChangeID(2))

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGetUpdate, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithEmptySans, WithPendingChangeID(2))

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
						Enabled: tools.BoolPtr(true),
					},
					SendCAListToClient: tools.BoolPtr(false),
				},
				SetID: "12345",
			}),
		)

		client.On("UpdateEnrollment",
			mock.Anything,
			cps.UpdateEnrollmentRequest{
				EnrollmentID:              1,
				Enrollment:                enrollmentUpdate,
				AllowCancelPendingChanges: tools.BoolPtr(true),
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.On("AcknowledgePreVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment:       enrollment,
				ContractID:       "1",
				AllowDuplicateCN: true,
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		allowCancel := true

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(nil, fmt.Errorf("error creating enrollment")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "Certificate data is blank or missing for 'RSA'.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		client.On("AcknowledgePreVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Twice()

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "Certificate data is blank or missing for 'RSA'.\nError parsing expected trust chains.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "Certificate data is blank or missing for 'RSA'.\nThis is unknown warning.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.NetworkConfiguration = &cps.NetworkConfiguration{
					DNSNameSettings: &cps.DNSNameSettings{
						CloneDNSNames: false,
					},
					Geography:     "core",
					SecureNetwork: "enhanced-tls",
					SNIOnly:       true,
				}
			}),
		)

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(4)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "1", rs.Attributes["id"])
							return nil
						},
					},
					{
						Config:            testutils.LoadFixtureString(t, "testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
						ImportState:       true,
						ImportStateId:     id,
						ResourceName:      "akamai_cps_third_party_enrollment.third_party",
						ImportStateVerify: true,
					},
				},
			})
		})
	})

	t.Run("import error when validation type is not third_party", func(t *testing.T) {
		client := &cps.Mock{}
		id := "1,ctr_1"

		enrollment := cps.Enrollment{
			ValidationType: "dv",
		}

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentGet := newEnrollment(WithBase(&enrollment), WithPendingChangeID(2))
		enrollmentGet.SignatureAlgorithm = ""

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		enrollmentUpdate := newEnrollment(WithBase(&enrollment),
			WithUpdateFunc(func(e *cps.Enrollment) {
				e.AdminContact.FirstName = "R5"
				e.AdminContact.LastName = "D5"
				e.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
				e.Location = ""
				e.PendingChanges = nil
				e.SignatureAlgorithm = ""
			}),
		)
		allowCancel := true
		client.On("UpdateEnrollment",
			mock.Anything,
			cps.UpdateEnrollmentRequest{
				Enrollment:                enrollmentUpdate,
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

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdateGet, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
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

func getSimpleEnrollment() cps.Enrollment {
	return cps.Enrollment{
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
	apply(*cps.Enrollment)
}

type withPendingChangeID int

func (w withPendingChangeID) apply(e *cps.Enrollment) {
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

func (w withCN) apply(e *cps.Enrollment) {
	e.CSR.CN = string(w)
}
func WithCN(cn string) enrolOpt {
	return withCN(cn)
}

type withFunc func(*cps.Enrollment)

func (w withFunc) apply(e *cps.Enrollment) {
	w(e)
}
func WithUpdateFunc(f func(*cps.Enrollment)) enrolOpt {
	return withFunc(f)
}

type withMTLS cps.ClientMutualAuthentication

func (w withMTLS) apply(e *cps.Enrollment) {
	e.NetworkConfiguration.ClientMutualAuthentication = (*cps.ClientMutualAuthentication)(&w)
}
func WithMTLS(mtls cps.ClientMutualAuthentication) enrolOpt {
	return withMTLS(mtls)
}

type withBase cps.Enrollment

func (w withBase) apply(e *cps.Enrollment) {
	*e = (cps.Enrollment)(w)
}
func WithBase(e *cps.Enrollment) enrolOpt {
	var newEn cps.Enrollment
	err := copier.CopyWithOption(&newEn, e, copier.Option{DeepCopy: true, IgnoreEmpty: true})
	if err != nil {
		panic(fmt.Sprintln("copier.CopyWithOption failed: ", err))
	}
	return withBase(newEn)
}

type withSans []string

func (w withSans) apply(e *cps.Enrollment) {
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

func newEnrollment(opts ...enrolOpt) cps.Enrollment {
	enrollment := getSimpleEnrollment()
	for _, o := range opts {
		o.apply(&enrollment)
	}
	return enrollment
}
