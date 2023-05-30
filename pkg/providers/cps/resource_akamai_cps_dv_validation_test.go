package cps

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDVValidation(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State: "running",
			}}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.On("AcknowledgeDVChallenges", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil)

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Times(3)

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State: "running",
			}}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.On("AcknowledgeDVChallenges", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil)

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Twice()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVValidation/create_validation.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResDVValidation/update_validation.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "2"),
						),
					},
				},
			})
			mock.AssertExpectationsForObjects(t)
		})
	})
	t.Run("retry acknowledgement", func(t *testing.T) {
		client := &cps.Mock{}
		changeAckRetryInterval = 1 * time.Millisecond
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.On("AcknowledgeDVChallenges", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(fmt.Errorf("oops")).Once()

		client.On("AcknowledgeDVChallenges", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Twice()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVValidation/create_validation.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "sans.#", "1"),
						),
					},
				},
			})
			mock.AssertExpectationsForObjects(t)
		})
	})
	t.Run("retry acknowledgement with timeout", func(t *testing.T) {
		client := &cps.Mock{}
		changeAckRetryInterval = 1 * time.Millisecond
		changeAckDeadline = 2 * time.Millisecond
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []cps.PendingChange{
				{
					Location:   "/cps/v2/enrollments/1/changes/2",
					ChangeType: "new-certificate",
				},
			}}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Once()

		client.On("AcknowledgeDVChallenges", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(fmt.Errorf("oops"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResDVValidation/create_validation.tf"),
						ExpectError: regexp.MustCompile("retry timeout reached - error sending acknowledgement request: oops"),
					},
				},
			})
			mock.AssertExpectationsForObjects(t)
		})
	})
}
