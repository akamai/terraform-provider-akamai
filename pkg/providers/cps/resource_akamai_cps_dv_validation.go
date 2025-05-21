package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	cpstools "github.com/akamai/terraform-provider-akamai/v8/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	changeAckRetryInterval = 10 * time.Second
)

func resourceCPSDVValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPSDVValidationCreate,
		ReadContext:   resourceCPSDVValidationRead,
		UpdateContext: resourceCPSDVValidationUpdate,
		DeleteContext: resourceCPSDVValidationDelete,

		Schema: map[string]*schema.Schema{
			"enrollment_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of enrollment",
			},
			"sans": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "List of SANs",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of validation",
			},
			"acknowledge_post_verification_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to acknowledge all post-verification warnings",
			},
			"timeouts": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables to set timeout for processing",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: timeouts.ValidateDurationFormat,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &timeouts.SDKDefaultTimeout,
		},
	}
}

func resourceCPSDVValidationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVValidationCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating dv validation")
	enrollmentID, err := tf.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.FromErr(err)
	}

	changeID, err := cpstools.GetChangeIDFromPendingChanges(res.PendingChanges)
	if err != nil {
		if errors.Is(err, cpstools.ErrNoPendingChanges) {
			logger.Debug("No pending changes found on the enrollment")
			d.SetId(strconv.Itoa(enrollmentID))
			return nil
		}
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// if status is `coordinate-domain-validation` or `wait-review-cert-warning` proceed further
	status, err := waitForChangeStatus(ctx, client, enrollmentID, changeID, coodinateDomainValidation, coordinateDomainValidation, waitReviewCertWarning)
	if err != nil {
		return diag.FromErr(err)
	}

	ackPostVerification, err := tf.GetBoolValue("acknowledge_post_verification_warnings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// if the status is `wait-review-cert-warning`, handle post warnings
	if status.StatusInfo != nil && status.StatusInfo.Status == waitReviewCertWarning && ackPostVerification {
		if err = sendPostVerificationAcknowledgement(ctx, client, enrollmentID, changeID); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(strconv.Itoa(enrollmentID))
		return resourceCPSDVValidationRead(ctx, d, m)
	}

	// if the status is `coordinate-domain-validation`, send ack for DV challenges
	err = client.AcknowledgeDVChallenges(ctx, cps.AcknowledgementRequest{
		Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
		EnrollmentID:    enrollmentID,
		ChangeID:        changeID,
	})
	if err == nil {
		status, err = waitForChangeStatus(ctx, client, enrollmentID, changeID, waitReviewCertWarning, complete, coordinateDomainValidation, coodinateDomainValidation)
		if err != nil {
			return diag.FromErr(err)
		}

		if status.StatusInfo != nil && status.StatusInfo.Status == waitReviewCertWarning && ackPostVerification {
			if err = sendPostVerificationAcknowledgement(ctx, client, enrollmentID, changeID); err != nil {
				return diag.FromErr(err)
			}
		}

		// for other statuses: `coordinate-domain-validation` and `complete`, go to read
		d.SetId(strconv.Itoa(enrollmentID))
		return resourceCPSDVValidationRead(ctx, d, m)
	}

	// in case of error, attempt retry
	logger.Debugf("error sending acknowledgement request: %s", err)
	for {
		select {
		case <-time.After(changeAckRetryInterval):
			err = client.AcknowledgeDVChallenges(ctx, cps.AcknowledgementRequest{
				Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
				EnrollmentID:    enrollmentID,
				ChangeID:        changeID,
			})
			if err == nil {
				status, err = waitForChangeStatus(ctx, client, enrollmentID, changeID, waitReviewCertWarning, complete, coordinateDomainValidation, coodinateDomainValidation)
				if err != nil {
					return diag.FromErr(err)
				}

				if status.StatusInfo != nil && status.StatusInfo.Status == waitReviewCertWarning && ackPostVerification {
					if err = sendPostVerificationAcknowledgement(ctx, client, enrollmentID, changeID); err != nil {
						return diag.FromErr(err)
					}
				}

				// for other statuses: `coordinate-domain-validation` and `complete`, go to read
				d.SetId(strconv.Itoa(enrollmentID))
				return resourceCPSDVValidationRead(ctx, d, m)
			}
		case <-ctx.Done():
			return diag.Errorf("retry timeout reached - error sending acknowledgement request: %s", err)
		}
	}
}

func resourceCPSDVValidationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVValidationRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Reading dv validation")
	enrollmentID, err := tf.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.FromErr(err)
	}

	changeID, err := cpstools.GetChangeIDFromPendingChanges(res.PendingChanges)
	if err != nil {
		if errors.Is(err, cpstools.ErrNoPendingChanges) {
			logger.Debug("No pending changes found on the enrollment")
			return nil
		}
		return diag.FromErr(err)
	}
	changeStatusReq := cps.GetChangeStatusRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	}
	status, err := client.GetChangeStatus(ctx, changeStatusReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(enrollmentID))
	if status.StatusInfo != nil {
		if err := d.Set("status", status.StatusInfo.Status); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourceCPSDVValidationUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVValidationUpdate")

	if !d.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	return diag.Errorf("Update in this resource is not allowed") //all fields are force new - it should never reach here
}

func resourceCPSDVValidationDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func sendPostVerificationAcknowledgement(ctx context.Context, client cps.CPS, enrollmentID, changeID int) error {
	if err := client.AcknowledgePostVerificationWarnings(ctx, cps.AcknowledgementRequest{
		Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
		EnrollmentID:    enrollmentID,
		ChangeID:        changeID,
	}); err != nil {
		return fmt.Errorf("could not acknowledge post-verification warnings: %s", err)
	}

	return nil
}
