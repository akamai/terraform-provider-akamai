package cps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	cpstools "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	changeAckDeadline      = 5 * time.Minute
	changeAckRetryInterval = 10 * time.Second
)

func resourceCPSDVValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPSDVValidationCreate,
		ReadContext:   resourceCPSDVValidationRead,
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
		},
	}
}

func resourceCPSDVValidationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceCPSDVValidationCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating dv validation")
	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
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
	changeStatusReq := cps.GetChangeStatusRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	}
	status, err := client.GetChangeStatus(ctx, changeStatusReq)
	if err != nil {
		return diag.FromErr(err)
	}
	for status.StatusInfo.Status != statusCoordinateDomainValidation {
		select {
		case <-time.After(PollForChangeStatusInterval):
			status, err = client.GetChangeStatus(ctx, changeStatusReq)
			if err != nil {
				return diag.FromErr(err)
			}
			changeStatusJSON, err := json.MarshalIndent(status, "", "\t")
			if err != nil {
				return diag.FromErr(err)
			}
			logger.Debugf("Change status: %s", changeStatusJSON)
			if status.StatusInfo != nil && status.StatusInfo.Error != nil && status.StatusInfo.Error.Description != "" {
				return diag.Errorf(status.StatusInfo.Error.Description)
			}
		case <-ctx.Done():
			return diag.Errorf("change status context terminated: %s", ctx.Err())
		}
	}
	err = client.AcknowledgeDVChallenges(ctx, cps.AcknowledgementRequest{
		Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
		EnrollmentID:    enrollmentID,
		ChangeID:        changeID,
	})
	if err == nil {
		d.SetId(strconv.Itoa(enrollmentID))
		return resourceCPSDVValidationRead(ctx, d, m)
	}

	// in case of error, attempt retry
	logger.Debugf("error sending acknowledgement request: %s", err)
	ackCtx, cancel := context.WithTimeout(ctx, changeAckDeadline)
	defer cancel()
	for {
		select {
		case <-time.After(changeAckRetryInterval):
			err = client.AcknowledgeDVChallenges(ctx, cps.AcknowledgementRequest{
				Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
				EnrollmentID:    enrollmentID,
				ChangeID:        changeID,
			})
			if err == nil {
				d.SetId(strconv.Itoa(enrollmentID))
				return resourceCPSDVValidationRead(ctx, d, m)
			}
		case <-ackCtx.Done():
			return diag.Errorf("retry timeout reached - error sending acknowledgement request: %s", err)
		}
	}
}

func resourceCPSDVValidationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceCPSDVValidationRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Reading dv validation")
	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
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
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourceCPSDVValidationDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
