package cps

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v9/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDVValidation(t *testing.T) {
	t.Parallel()
	t.Run("lifecycle test", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State: "running",
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("AcknowledgeDVChallenges", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Times(3)

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State: "running",
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("AcknowledgeDVChallenges", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVValidation/create_validation.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVValidation/update_validation.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "2"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.0.default", "1h"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
	t.Run("lifecycle test with ack post verification warnings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State: "running",
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("AcknowledgeDVChallenges", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "wait-review-cert-warning",
			}}, nil).Once()

		client.CPS.On("AcknowledgePostVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{
				Acknowledgement: cps.AcknowledgementAcknowledge,
			},
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Times(3)

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State: "running",
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVValidation/create_validation_with_ack_post_verification.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "acknowledge_post_verification_warnings", strconv.FormatBool(true)),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVValidation/update_validation.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "2"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "acknowledge_post_verification_warnings", strconv.FormatBool(false)),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.0.default", "1h"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
	t.Run("receive `wait-review-cert-warning` early", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "running",
				Status: waitReviewCertWarning,
			}}, nil).Once()

		client.CPS.On("AcknowledgePostVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{
				Acknowledgement: cps.AcknowledgementAcknowledge,
			},
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "complete",
			}}, nil).Times(2)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVValidation/create_validation_with_ack_post_verification.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "complete"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "acknowledge_post_verification_warnings", strconv.FormatBool(true)),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "timeouts.#", "0"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
	t.Run("retry acknowledgement", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("AcknowledgeDVChallenges", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(fmt.Errorf("oops")).Once()

		client.CPS.On("AcknowledgeDVChallenges", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVValidation/create_validation.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
						resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
	t.Run("retry acknowledgement with timeout", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.CPS.On("AcknowledgeDVChallenges", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(fmt.Errorf("oops"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVValidation/create_validation_with_timeout.tf"),
					ExpectError: regexp.MustCompile("retry timeout reached - error sending acknowledgement request: oops"),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
}
