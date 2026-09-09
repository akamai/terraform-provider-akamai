package cps

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ action.Action              = &ForceCertificateRenewalAction{}
	_ action.ActionWithConfigure = &ForceCertificateRenewalAction{}
)

const forceRenewalGenericError = "Forcing certificate renewal failed"

// ForceCertificateRenewalAction represents the akamai_cps_force_certificate_renewal action
type ForceCertificateRenewalAction struct {
	meta.Action
	pollChangeStatusInterval time.Duration
}

// NewForceCertificateRenewalAction returns a new instance of akamai_cps_force_certificate_renewal action
func NewForceCertificateRenewalAction(pollChangeStatusInterval time.Duration) func() action.Action {
	return func() action.Action {
		return &ForceCertificateRenewalAction{
			pollChangeStatusInterval: pollChangeStatusInterval,
		}
	}
}

// ForceCertificateRenewalActionModel represents the model for akamai_cps_force_certificate_renewal action configuration
type ForceCertificateRenewalActionModel struct {
	EnrollmentID                       types.Int64 `tfsdk:"enrollment_id"`
	CancelPendingChanges               types.Bool  `tfsdk:"cancel_pending_changes"`
	AcknowledgePreVerificationWarnings types.Bool  `tfsdk:"acknowledge_pre_verification_warnings"`
}

// Metadata sets the full name of the akamai_cps_force_certificate_renewal action
func (a *ForceCertificateRenewalAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cps_force_certificate_renewal"
}

// Schema sets the schema of the akamai_cps_force_certificate_renewal action
func (a *ForceCertificateRenewalAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forces the renewal of the certificate for the given CPS enrollment.",
		Attributes: map[string]schema.Attribute{
			"enrollment_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the CPS enrollment for which to force certificate renewal.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"cancel_pending_changes": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to cancel any pending changes for the enrollment. Defaults to false.",
			},
			"acknowledge_pre_verification_warnings": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to acknowledge all pre-verification warnings. Defaults to false.",
			},
		},
	}
}

// Invoke forces the renewal of the certificate for the given CPS enrollment by calling UpdateEnrollment with
// force-renewal set to true; it can optionally cancel any pending changes and acknowledge pre-verification warnings.
// It waits for the certificate verification to complete before returning.
func (a *ForceCertificateRenewalAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ForceCertificateRenewalActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logger := a.Log("CPS", "ForceCertificateRenewalAction")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := a.Client.GetCPS()
	enrollmentID := int(data.EnrollmentID.ValueInt64())
	acknowledgePreVerificationWarnings := data.AcknowledgePreVerificationWarnings.ValueBool()
	cancelPendingChanges := data.CancelPendingChanges.ValueBool()

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		resp.Diagnostics.AddError("Fetching enrollment failed", err.Error())
		return
	}
	if err := validateEnrollmentType(enrollment.ValidationType); err != nil {
		resp.Diagnostics.AddError(forceRenewalGenericError, err.Error())
		return
	}

	if !cancelPendingChanges && len(enrollment.PendingChanges) > 0 {
		resp.Diagnostics.AddError(forceRenewalGenericError, fmt.Sprintf("Enrollment %d has pending changes. "+
			"To proceed set cancel_pending_changes to true and invoke the action again.", enrollmentID))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Forcing certificate renewal for enrollment %d...", enrollmentID),
	})

	_, err = client.UpdateEnrollment(ctx, cps.UpdateEnrollmentRequest{
		EnrollmentRequestBody:     enrollmentRequestBodyFromEnrollmentResponse(*enrollment),
		EnrollmentID:              enrollmentID,
		ForceRenewal:              ptr.To(true),
		AllowCancelPendingChanges: ptr.To(cancelPendingChanges),
	})
	if err != nil {
		resp.Diagnostics.AddError(forceRenewalGenericError, err.Error())
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Successfully sent force renewal request. Waiting for verification for enrollment %d...",
			enrollmentID),
	})

	if err := waitForVerification(ctx, logger, client, enrollmentID, acknowledgePreVerificationWarnings,
		nil, a.pollChangeStatusInterval); err != nil {
		detail := err.Error()
		if !cancelPendingChanges {
			detail += ". The renewal might still be pending. Resolve the reported issue if possible, then retry the action with " +
				"cancel_pending_changes set to true. This instructs CPS to cancel this pending renewal before starting a new one."
		}
		resp.Diagnostics.AddError("Waiting for certificate verification failed", detail)
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Certificate verification completed for enrollment %d.", enrollmentID),
	})

	addPostRenewalWarning(resp, enrollment.ValidationType, enrollmentID)
}

func validateEnrollmentType(validationType string) error {
	if validationType != "dv" && validationType != "third-party" {
		return fmt.Errorf("unsupported enrollment validation type %q: force renewal supports only %q and %q", validationType, "dv", "third-party")
	}
	return nil
}

func enrollmentRequestBodyFromEnrollmentResponse(e cps.GetEnrollmentResponse) cps.EnrollmentRequestBody {
	return cps.EnrollmentRequestBody{
		AdminContact:                   e.AdminContact,
		AutoRenewalStartTime:           e.AutoRenewalStartTime,
		CertificateChainType:           e.CertificateChainType,
		CertificateType:                e.CertificateType,
		ChangeManagement:               e.ChangeManagement,
		CSR:                            e.CSR,
		EnableMultiStackedCertificates: e.EnableMultiStackedCertificates,
		NetworkConfiguration:           e.NetworkConfiguration,
		Org:                            e.Org,
		OrgID:                          e.OrgID,
		RA:                             e.RA,
		SignatureAlgorithm:             e.SignatureAlgorithm,
		TechContact:                    e.TechContact,
		ThirdParty:                     e.ThirdParty,
		ValidationType:                 e.ValidationType,
	}
}

// addPostRenewalWarning adds a warning to the response indicating that further action is required to complete the renewal.
func addPostRenewalWarning(resp *action.InvokeResponse, validationType string, enrollmentID int) {
	if validationType == "third-party" {
		resp.Diagnostics.AddWarning(
			"Action required: third-party certificate",
			fmt.Sprintf("Enrollment %d is a third-party enrollment. To complete the renewal, fetch the CSR via the "+
				"akamai_cps_csr data source, sign it with your CA, and upload the signed certificate via the "+
				"akamai_cps_upload_certificate resource", enrollmentID),
		)
	} else {
		resp.Diagnostics.AddWarning(
			"Action required: DV certificate",
			fmt.Sprintf("Enrollment %d is a DV enrollment. To complete the renewal, you might need to complete domain "+
				"validation. You can fetch challenges by using the akamai_cps_dv_enrollment resource or the "+
				"akamai_cps_enrollment data source", enrollmentID),
		)
	}
}
