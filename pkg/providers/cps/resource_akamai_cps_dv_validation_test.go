package cps

import (
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDVValidation(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		client := &mockcps{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []string{"/cps/v2/enrollments/1/changes/2"}}, nil).Once()

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
			Acknowledgement: cps.Acknowledgement{"acknowledge"},
			EnrollmentID:    1,
			ChangeID:        2,
		}).Return(nil)

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&cps.Enrollment{PendingChanges: []string{"/cps/v2/enrollments/1/changes/2"}}, nil).Twice()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}).
			Return(&cps.Change{StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			}}, nil).Twice()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVValidation/create_validation.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "id", "1"),
							resource.TestCheckResourceAttr("akamai_cps_dv_validation.dv_validation", "status", "coodinate-domain-validation"),
						),
					},
				},
			})
			mock.AssertExpectationsForObjects(t)
		})
	})
}
