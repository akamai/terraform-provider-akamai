package cps

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func resourceCPSDVValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPSDVValidationCreate,
		ReadContext:   resourceCPSDVValidationRead,
		DeleteContext: resourceCPSDVValidationDelete,

		Schema: map[string]*schema.Schema{
			"enrollment_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCPSDVValidationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("DEBUG: enter resourceCPSDVValidationCreate")
	meta := akamai.Meta(m)
	log := meta.Log("CPS", "resourceDVEnrollment")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	client := inst.Client(meta)
	log.Debug("Creating dv validation")
	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.FromErr(err)
	}

	changeURL, err := url.Parse(res.PendingChanges[0])
	if err != nil {
		return diag.FromErr(err)
	}
	pathSplit := strings.Split(changeURL.Path, "/")
	changeIDStr := pathSplit[len(pathSplit)-1]
	changeID, err := strconv.Atoi(changeIDStr)
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
	for status.StatusInfo.State == "running" {
		select {
		case <-time.After(10 * time.Second):
			status, err = client.GetChangeStatus(ctx, changeStatusReq)
			if err != nil {
				return diag.FromErr(err)
			}
			log.Debugf("Change status: %s", status.StatusInfo.Status)
			if status.StatusInfo != nil && status.StatusInfo.Error != nil && status.StatusInfo.Error.Description != "" {
				return diag.Errorf(status.StatusInfo.Error.Description)
			}
			if status.StatusInfo.Status != "coodinate-domain-validation" {
				return diag.Errorf("invalid validation status received: %s", status.StatusInfo.Status)
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
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(enrollmentID))

	return resourceCPSDVValidationRead(ctx, d, m)
}

func resourceCPSDVValidationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("DEBUG: enter resourceCPSDVValidationCreate")
	meta := akamai.Meta(m)
	log := meta.Log("CPS", "resourceDVEnrollment")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	client := inst.Client(meta)
	log.Debug("Reading dv validation")
	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.FromErr(err)
	}

	changeURL, err := url.Parse(res.PendingChanges[0])
	if err != nil {
		return diag.FromErr(err)
	}
	pathSplit := strings.Split(changeURL.Path, "/")
	changeIDStr := pathSplit[len(pathSplit)-1]
	changeID, err := strconv.Atoi(changeIDStr)
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

	if status != nil {
		if err := d.Set("status", status.StatusInfo.Status); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		d.SetId(strconv.Itoa(enrollmentID))
		return nil
	}

	d.SetId("")
	return nil
}

func resourceCPSDVValidationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
