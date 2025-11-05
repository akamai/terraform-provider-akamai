package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v9/internal/text"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &uploadSignedCertificateResource{}
	_ resource.ResourceWithConfigure   = &uploadSignedCertificateResource{}
	_ resource.ResourceWithImportState = &uploadSignedCertificateResource{}
	_ resource.ResourceWithModifyPlan  = &uploadSignedCertificateResource{}
)

type uploadSignedCertificateResource struct {
	meta meta.Meta
}

// NewUploadSignedCertificateResource returns a new CloudCertificates Certificate resource.
func NewUploadSignedCertificateResource() resource.Resource {
	return &uploadSignedCertificateResource{}
}

func (c *uploadSignedCertificateResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_upload_signed_certificate"
}

// Configure implements resource.ResourceWithConfigure.
func (c *uploadSignedCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"unexpected resource configure type",
				fmt.Sprintf("expected meta.Meta, got: %T. please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	c.meta = meta.Must(req.ProviderData)
}

func (c *uploadSignedCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"certificate_id": schema.StringAttribute{
				Required:    true,
				Description: "Certificate identifier on which to perform the upload operation.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"signed_certificate_pem": schema.StringAttribute{
				Required:    true,
				Description: "PEM-encoded signed certificate to upload.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(validators.CertificatePEMRegex,
						fmt.Sprintf("must be in PEM format: '%s'", validators.CertificatePEMRegex))},
			},
			"trust_chain_pem": schema.StringAttribute{
				Optional:    true,
				Description: "PEM-encoded trust chain for the signed certificate to upload.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(validators.ToolchainPEMRegex,
						fmt.Sprintf("must be in PEM format: '%s'", validators.ToolchainPEMRegex))},
			},
			"acknowledge_warnings": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Acknowledges warnings and retries certificate upload when " +
					"the returned response contains warnings for the uploaded certificate.",
			},
			"modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date the certificate was last updated.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who last modified the certificate.",
			},
			"certificate_status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the certificate. Can be one of 'CSR_READY', 'READY_FOR_USE', 'ACTIVE'.",
			},
			"signed_certificate_not_valid_after_date": schema.StringAttribute{
				Computed:    true,
				Description: "This marks the end of the signed certificate's valid period.",
			},
			"signed_certificate_not_valid_before_date": schema.StringAttribute{
				Computed:    true,
				Description: "This marks the start of the signed certificate's valid period.",
			},
			"signed_certificate_serial_number": schema.StringAttribute{
				Computed:    true,
				Description: "Signed certificate serial number in hex format.",
			},
			"signed_certificate_sha256_fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 fingerprint of the signed certificate.",
			},
			"signed_certificate_issuer": schema.StringAttribute{
				Computed:    true,
				Description: "Issuer field of the signed certificate.",
			},
		},
	}
}

func (c *uploadSignedCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CCM Upload Signed Certificate resource Create")

	var plan uploadSignedCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := c.uploadSignedCertificate(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error uploading signed certificate during resource creation",
			err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

func (c *uploadSignedCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "CCM Upload Signed Certificate resource Read")

	var state uploadSignedCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "certificate_id", state.CertificateID.ValueString())

	client := Client(c.meta)
	cert, err := client.GetCertificate(ctx, cloudcertificates.GetCertificateRequest{
		CertificateID: state.CertificateID.ValueString(),
	})

	if err != nil {
		if errors.Is(err, cloudcertificates.ErrCertificateNotFound) {
			resp.Diagnostics.AddError("CCM Certificate not found",
				fmt.Sprintf("The certificate '%s' was not found on the server: %s",
					state.CertificateID.ValueString(), err.Error()))
		} else {
			resp.Diagnostics.AddError("Error reading CCM Certificate",
				fmt.Sprintf("Error retrieving certificate '%s': %s",
					state.CertificateID.ValueString(), err.Error()))
		}
		return
	}

	state.populateCertFields(cert.Certificate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *uploadSignedCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "CCM Upload Signed Certificate resource Update")

	var plan uploadSignedCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := c.uploadSignedCertificate(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error uploading signed certificate during resource update",
			err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

func (c *uploadSignedCertificateResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "CCM Upload Signed Certificate resource Delete")
}

func (c *uploadSignedCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx = tflog.SetField(ctx, "import_id", req.ID)
	tflog.Debug(ctx, "CCM Upload Signed Certificate resource ImportState")

	parts, err := text.ImportIDSplitter("certificateID[,acknowledge_warnings]").
		AcceptLen(1).
		AcceptLen(2).
		Split(req.ID)

	if err != nil {
		resp.Diagnostics.AddError("Incorrect import ID", err.Error())
		return
	}
	certificateID := parts[0]

	state := uploadSignedCertificateResourceModel{
		CertificateID: types.StringValue(certificateID),
	}
	if len(parts) == 2 {
		b, err := strconv.ParseBool(parts[1])
		if err != nil {
			resp.Diagnostics.AddError("Incorrect import ID", "acknowledge_warnings must be 'true' or 'false'")
			return
		}
		state.AcknowledgeWarnings = types.BoolValue(b)
	} else {
		// If not provided, set to false to match default value
		state.AcknowledgeWarnings = types.BoolValue(false)
	}

	client := Client(c.meta)
	cert, err := client.GetCertificate(ctx, cloudcertificates.GetCertificateRequest{
		CertificateID: certificateID,
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrCertificateNotFound) {
			resp.Diagnostics.AddError("CCM Certificate for import not found",
				fmt.Sprintf("The certificate '%s' was not found on the server: %s",
					state.CertificateID.ValueString(), err.Error()))
		} else {
			resp.Diagnostics.AddError("Error reading CCM Certificate for import",
				fmt.Sprintf("Error retrieving certificate '%s': %s",
					state.CertificateID.ValueString(), err.Error()))
		}
		return
	}

	if cert.Certificate.CertificateStatus == string(cloudcertificates.StatusCSRReady) {
		resp.Diagnostics.AddError("Cannot import CCM Certificate in 'CSR_READY' status",
			fmt.Sprintf("The certificate '%s' has status '%s' and does not support importing "+
				"as the signed certificate PEM has not been uploaded yet.",
				state.CertificateID.ValueString(), cert.Certificate.CertificateStatus))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *uploadSignedCertificateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "CCM Upload Signed Certificate resource ModifyPlan")

	// Do not perform checks on delete
	if modifiers.IsDelete(req) {
		return
	}

	var plan uploadSignedCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CertificateID.IsNull() || plan.CertificateID.IsUnknown() {
		return
	}

	// We need to fetch the current certificate to check following conditions at plan time:
	// - If the certificate exists (the ID may not come from the certificate resource)
	// - If the certificate is already in READY_FOR_USE or ACTIVE status, in which case
	//   we cannot upload the signed certificate (and optionally the trust chain) again.
	// Note that during creation, the state will be empty, so we decided to always fetch
	// the status from the API for simplicity and up-to-date information.
	client := Client(c.meta)
	cert, err := client.GetCertificate(ctx, cloudcertificates.GetCertificateRequest{
		CertificateID: plan.CertificateID.ValueString(),
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrCertificateNotFound) {
			resp.Diagnostics.AddAttributeError(
				path.Root("certificate_id"),
				"Cannot upload signed certificate to a non-existent CCM Certificate object",
				fmt.Sprintf("The certificate '%s' was not found on the server. "+
					"Please verify certificate_id is correct.\n\n%s",
					plan.CertificateID.ValueString(), err.Error()))
		} else {
			resp.Diagnostics.AddError("Unable to get CCM Certificate for signed certificate upload",
				fmt.Sprintf("Error retrieving certificate '%s': %s",
					plan.CertificateID.ValueString(), err.Error()))
		}
		return
	}

	if cert.Certificate.CertificateStatus == string(cloudcertificates.StatusReadyForUse) ||
		cert.Certificate.CertificateStatus == string(cloudcertificates.StatusActive) {
		// Use pointer - state can be null if resource is being created
		var state *uploadSignedCertificateResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if state == nil || plan.hasDifferentInput(*state) {
			resp.Diagnostics.AddError("Cannot upload signed certificate",
				fmt.Sprintf("The certificate '%s' has status '%s' and does not support uploading a signed certificate "+
					"as it has been already uploaded.", cert.Certificate.CertificateID, cert.Certificate.CertificateStatus))
		}
	}
}

func (c *uploadSignedCertificateResource) uploadSignedCertificate(ctx context.Context,
	m uploadSignedCertificateResourceModel) (uploadSignedCertificateResourceModel, error) {

	patchReq := cloudcertificates.PatchCertificateRequest{
		CertificateID:        m.CertificateID.ValueString(),
		SignedCertificatePEM: m.SignedCertificatePEM.ValueString(),
		TrustChainPEM:        m.TrustChainPEM.ValueString(),
		AcknowledgeWarnings:  m.AcknowledgeWarnings.ValueBool(),
	}

	client := Client(c.meta)
	cert, err := client.PatchCertificate(ctx, patchReq)
	if err != nil {
		return uploadSignedCertificateResourceModel{}, err
	}
	m.populateCertFields(cert.Certificate)
	return m, nil
}

type uploadSignedCertificateResourceModel struct {
	CertificateID                       types.String `tfsdk:"certificate_id"`
	SignedCertificatePEM                types.String `tfsdk:"signed_certificate_pem"`
	TrustChainPEM                       types.String `tfsdk:"trust_chain_pem"`
	AcknowledgeWarnings                 types.Bool   `tfsdk:"acknowledge_warnings"`
	ModifiedDate                        types.String `tfsdk:"modified_date"`
	ModifiedBy                          types.String `tfsdk:"modified_by"`
	CertificateStatus                   types.String `tfsdk:"certificate_status"`
	SignedCertificateNotValidAfterDate  types.String `tfsdk:"signed_certificate_not_valid_after_date"`
	SignedCertificateNotValidBeforeDate types.String `tfsdk:"signed_certificate_not_valid_before_date"`
	SignedCertificateSerialNumber       types.String `tfsdk:"signed_certificate_serial_number"`
	SignedCertificateSHA256Fingerprint  types.String `tfsdk:"signed_certificate_sha256_fingerprint"`
	SignedCertificateIssuer             types.String `tfsdk:"signed_certificate_issuer"`
}

func (m *uploadSignedCertificateResourceModel) populateCertFields(cert cloudcertificates.Certificate) {
	m.SignedCertificatePEM = types.StringPointerValue(cert.SignedCertificatePEM)
	m.TrustChainPEM = types.StringPointerValue(cert.TrustChainPEM)
	m.ModifiedDate = date.TimeRFC3339NanoValue(cert.ModifiedDate)
	m.ModifiedBy = types.StringValue(cert.ModifiedBy)
	m.CertificateStatus = types.StringValue(cert.CertificateStatus)
	m.SignedCertificateNotValidAfterDate = date.TimeRFC3339PointerValue(cert.SignedCertificateNotValidAfterDate)
	m.SignedCertificateNotValidBeforeDate = date.TimeRFC3339PointerValue(cert.SignedCertificateNotValidBeforeDate)
	m.SignedCertificateSerialNumber = types.StringPointerValue(cert.SignedCertificateSerialNumber)
	m.SignedCertificateSHA256Fingerprint = types.StringPointerValue(cert.SignedCertificateSHA256Fingerprint)
	m.SignedCertificateIssuer = types.StringPointerValue(cert.SignedCertificateIssuer)
}

func (m *uploadSignedCertificateResourceModel) hasDifferentInput(other uploadSignedCertificateResourceModel) bool {
	return m.CertificateID != other.CertificateID ||
		m.SignedCertificatePEM != other.SignedCertificatePEM ||
		m.TrustChainPEM != other.TrustChainPEM ||
		m.AcknowledgeWarnings != other.AcknowledgeWarnings
}
