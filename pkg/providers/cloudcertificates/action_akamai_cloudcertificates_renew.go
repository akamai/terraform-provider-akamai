package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = &RenewAction{}
	_ action.ActionWithConfigure = &RenewAction{}
)

// RenewAction represents the akamai_cloudcertificates_renew action.
type RenewAction struct {
	meta.Action
}

// NewRenewAction returns a new instance of the akamai_cloudcertificates_renew action.
func NewRenewAction() action.Action {
	return &RenewAction{}
}

// RenewActionModel represents the model for akamai_cloudcertificates_renew action configuration.
type RenewActionModel struct {
	LineageID          types.Int64 `tfsdk:"lineage_id"`
	ConfirmAbandonHead types.Bool  `tfsdk:"confirm_abandon_head"`
}

// Metadata sets the full name of the akamai_cloudcertificates_renew action.
func (a *RenewAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_renew"
}

// Schema sets the schema of the akamai_cloudcertificates_renew action.
func (a *RenewAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Renews a certificate lineage: creates a new head generation with fresh private keys and CSRs " +
			"for every algorithm declared in the lineage's key specs.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the lineage to renew.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"confirm_abandon_head": schema.BoolAttribute{
				Optional: true,
				Description: "Whether to confirm abandoning the existing head generation when it already has a " +
					"still-valid, unexpired certificate. Defaults to false.",
			},
		},
	}
}

// Invoke renews the lineage's head generation by calling RenewLineage.
func (a *RenewAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RenewActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cloud Certificates Renew Action Invoke")

	lineageID := data.LineageID.ValueInt64()
	client := a.Client.GetCloudCertificates()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Renewing lineage %d...", lineageID),
	})

	lineage, err := client.RenewLineage(ctx, cloudcertificates.RenewLineageRequest{
		LineageID:          lineageID,
		ConfirmAbandonHead: data.ConfirmAbandonHead.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.Append(renewLineageErrorDiagnostics(err, lineageID)...)
		return
	}

	if lineage.Head == nil {
		resp.Diagnostics.AddError("Unexpected Response",
			fmt.Sprintf("Lineage %d was renewed, but the response did not include a new head generation.", lineageID))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Lineage %d renewed: new head generation %d created with fresh CSRs for every algorithm.",
			lineageID, lineage.Head.HeadGenerationID),
	})

	detail := fmt.Sprintf("Lineage %d's new head generation %d still needs to be signed before it can be deployed. "+
		"Retrieve its CSRs (for example via the signing_target attribute of the akamai_cloudcertificates_lineage "+
		"resource or data source), get them signed by your CA, then update akamai_cloudcertificates_upload with the "+
		"newly signed certificates and run terraform apply - the upload resource will detect that the head moved and "+
		"recreate itself for generation %d. If akamai_cloudcertificates_activation's staging_generation_id/"+
		"production_generation_id reference the upload resource's generation_id (the recommended wiring), promotion "+
		"to staging and/or production then happens automatically on that same apply; otherwise, update those IDs "+
		"manually to %d.", lineageID, lineage.Head.HeadGenerationID, lineage.Head.HeadGenerationID, lineage.Head.HeadGenerationID)

	if data.ConfirmAbandonHead.ValueBool() {
		detail += " If a previous head generation existed, it was abandoned (along with any uploaded certificates) and cannot be recovered."
	}
	resp.Diagnostics.AddWarning("Manual Follow-Up Required", detail)
}

// renewLineageErrorDiagnostics maps well-documented RenewLineage sentinel errors to clear diagnostics.
func renewLineageErrorDiagnostics(err error, lineageID int64) diag.Diagnostics {
	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		diags.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d.", lineageID))
	case errors.Is(err, cloudcertificates.ErrRenewBlockedHeadOnStaging):
		diags.AddError("Head Deployed On Staging",
			fmt.Sprintf("Lineage %d's head generation is currently deployed on staging. Use the "+
				"akamai_cloudcertificates_replace_staging action to free it up, then renew again.", lineageID))
	case errors.Is(err, cloudcertificates.ErrConfirmationRequired):
		diags.AddError("Confirmation Required",
			fmt.Sprintf("Lineage %d's head generation has a still-valid, unexpired certificate. Set "+
				"confirm_abandon_head to true and invoke the action again to abandon it and renew.", lineageID))
	default:
		diags.AddError("Failed To Renew Lineage", fmt.Sprintf("Failed to renew lineage %d: %s", lineageID, err))
	}
	return diags
}
