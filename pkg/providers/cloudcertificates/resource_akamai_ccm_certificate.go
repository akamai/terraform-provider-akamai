package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	"github.com/akamai/terraform-provider-akamai/v9/internal/text"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &certificateResource{}
	_ resource.ResourceWithConfigure      = &certificateResource{}
	_ resource.ResourceWithImportState    = &certificateResource{}
	_ resource.ResourceWithModifyPlan     = &certificateResource{}
	_ resource.ResourceWithValidateConfig = &certificateResource{}
)

// The suffix added by CCM to the certificate name when a certificate is renewed.
const renewedNameSuffix = ".renewed."

// The date format used in the renewed certificate name in format <base_name>.renewed.YYYY-MM-DD.
const renewedNameDateLayout = "2006-01-02"

// The regex pattern used by API to validate domain names: commonName and SANs.
// Modified to disallow uppercase letters, as API will lowercase them automatically.
// Original pattern from API: ^(\*\.)?([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$
var domainNameRegex = regexp.MustCompile(`^(\*\.)?([a-z0-9]([a-z0-9-]*[a-z0-9])?\.)+[a-z]{2,}$`)

type certificateResourceModel struct {
	ContractID    types.String `tfsdk:"contract_id"`
	GroupID       types.String `tfsdk:"group_id"`
	BaseName      types.String `tfsdk:"base_name"`
	Name          types.String `tfsdk:"name"`
	KeyType       types.String `tfsdk:"key_type"`
	KeySize       types.String `tfsdk:"key_size"`
	SecureNetwork types.String `tfsdk:"secure_network"`
	SANs          types.Set    `tfsdk:"sans"`
	Subject       types.Object `tfsdk:"subject"`
	// TODO: For next story - implement renew_before_expiration_days logic.
	// RenewBeforeExpirationDays types.Int64  `tfsdk:"renew_before_expiration_days"`
	// NeedsRenewal              types.Bool   `tfsdk:"needs_renewal"`
	CertificateID     types.String `tfsdk:"certificate_id"`
	CertificateType   types.String `tfsdk:"certificate_type"`
	AccountID         types.String `tfsdk:"account_id"`
	CreatedDate       types.String `tfsdk:"created_date"`
	CreatedBy         types.String `tfsdk:"created_by"`
	ModifiedDate      types.String `tfsdk:"modified_date"`
	ModifiedBy        types.String `tfsdk:"modified_by"`
	CertificateStatus types.String `tfsdk:"certificate_status"`
	CSRPEM            types.String `tfsdk:"csr_pem"`
	CSRExpirationDate types.String `tfsdk:"csr_expiration_date"`
}

type subjectModel struct {
	CommonName   types.String `tfsdk:"common_name"`
	Organization types.String `tfsdk:"organization"`
	Country      types.String `tfsdk:"country"`
	State        types.String `tfsdk:"state"`
	Locality     types.String `tfsdk:"locality"`
}

func (m *certificateResourceModel) validateKeyTypeAndSize() diag.Diagnostics {
	validKeyCombinations := map[string][]string{
		"RSA":   {"2048"},
		"ECDSA": {"P-256"},
	}

	var diags diag.Diagnostics
	if m.KeyType.IsNull() || m.KeyType.IsUnknown() || m.KeySize.IsNull() || m.KeySize.IsUnknown() {
		return diags
	}

	validSizes, ok := validKeyCombinations[m.KeyType.ValueString()]
	if ok && !slices.Contains(validSizes, m.KeySize.ValueString()) {
		diags.AddAttributeError(
			path.Root("key_size"),
			fmt.Sprintf("Invalid key size for %s.", m.KeyType.ValueString()),
			fmt.Sprintf("The specified value '%s' for the %s key type is invalid. Valid values are '%s'.",
				m.KeySize.ValueString(), m.KeyType.ValueString(), strings.Join(validSizes, "', '")),
		)
	}
	return diags
}

func (m *certificateResourceModel) validateSubjectAndSANs(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	if m.Subject.IsNull() || m.Subject.IsUnknown() {
		return diags
	}

	var subject subjectModel
	diags.Append(m.Subject.As(ctx, &subject, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	if subject.CommonName.IsNull() && subject.Organization.IsNull() && subject.Country.IsNull() &&
		subject.State.IsNull() && subject.Locality.IsNull() {
		diags.AddAttributeError(
			path.Root("subject"),
			"Missing Required Attribute",
			"At least one of the subject fields (common_name, organization, country, state, locality) must be specified.",
		)
	}

	if m.SANs.IsUnknown() {
		return diags
	}

	var sans []string
	diags.Append(m.SANs.ElementsAs(ctx, &sans, false)...)
	if diags.HasError() {
		return diags
	}

	if subject.CommonName.IsNull() || subject.CommonName.IsUnknown() {
		return diags
	}
	if !slices.Contains(sans, subject.CommonName.ValueString()) {
		diags.AddAttributeError(
			path.Root("sans"),
			"SANs missing common name",
			fmt.Sprintf("The specified common name '%s' must be included in the SANs list.",
				subject.CommonName.ValueString()),
		)
	}
	return diags
}

func (m *certificateResourceModel) populateCertificateFields(ctx context.Context, cert ccm.Certificate, setContractID bool) diag.Diagnostics {
	var diags diag.Diagnostics

	m.KeySize = types.StringValue(string(cert.KeySize))
	m.KeyType = types.StringValue(string(cert.KeyType))
	m.SecureNetwork = types.StringValue(string(cert.SecureNetwork))
	m.AccountID = types.StringValue(cert.AccountID)
	m.Name = types.StringValue(cert.CertificateName)
	m.CertificateID = types.StringValue(cert.CertificateID)
	m.CertificateStatus = types.StringValue(cert.CertificateStatus)
	m.CertificateType = types.StringValue(cert.CertificateType)
	m.CreatedBy = types.StringValue(cert.CreatedBy)
	m.CreatedDate = date.TimeRFC3339NanoValue(cert.CreatedDate)
	m.ModifiedBy = types.StringValue(cert.ModifiedBy)
	m.ModifiedDate = date.TimeRFC3339NanoValue(cert.ModifiedDate)
	m.CSRExpirationDate = date.TimeRFC3339Value(cert.CSRExpirationDate)
	m.CSRPEM = types.StringPointerValue(cert.CSRPEM)

	sans, dd := types.SetValueFrom(ctx, types.StringType, cert.SANs)
	diags.Append(dd...)
	if diags.HasError() {
		return diags
	}
	m.SANs = sans

	if cert.Subject != nil && !isEmptySubject(*cert.Subject) {
		sm := subjectModel{}
		sm.CommonName = tf.StringValueOrNullIfEmpty(cert.Subject.CommonName)
		sm.Organization = tf.StringValueOrNullIfEmpty(cert.Subject.Organization)
		sm.Country = tf.StringValueOrNullIfEmpty(cert.Subject.Country)
		sm.State = tf.StringValueOrNullIfEmpty(cert.Subject.State)
		sm.Locality = tf.StringValueOrNullIfEmpty(cert.Subject.Locality)

		m.Subject, dd = types.ObjectValueFrom(ctx, subjectType(), sm)
		diags.Append(dd...)
		if diags.HasError() {
			return diags
		}
	} else {
		m.Subject = types.ObjectNull(subjectType())
	}

	if setContractID {
		tflog.Debug(ctx, fmt.Sprintf("Setting contract_id to %s", cert.ContractID))
		m.ContractID = types.StringValue(cert.ContractID)
	}

	return diags
}

type certificateResource struct {
	meta meta.Meta
}

// NewCertificateResource returns a new CloudCertificates Certificate resource.
func NewCertificateResource() resource.Resource {
	return &certificateResource{}
}

func (c *certificateResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_certificate"
}

// Configure implements resource.ResourceWithConfigure.
func (c *certificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"unexpected resource configure type",
				fmt.Sprintf("expected meta.Meta, got: %T. please report this issue to the provider developers.",
					req.ProviderData),
			)
		}
	}()

	c.meta = meta.Must(req.ProviderData)
}

func (c *certificateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "Modifying Client Certificate Third Party Resource Plan")
	var state, plan *certificateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !modifiers.IsUpdate(req) {
		return
	}

	if state.GroupID == types.StringNull() && plan.GroupID != types.StringNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("group_id"),
			"The resource was imported without a group_id",
			"To fix this, you need to first remove the state and then re-import it with the group_id specified in the import ID.",
		)
		return
	} else if !state.GroupID.Equal(plan.GroupID) {
		resp.Diagnostics.AddError(
			"Update not Supported",
			"updating field `group_id` is not possible")
	}

	// Changing only 'base_name' results in a change for 'name', 'modified_by' and 'modified_date' fields.
	if state.BaseName.Equal(plan.BaseName) {
		plan.Name = state.Name
		plan.ModifiedBy = state.ModifiedBy
		plan.ModifiedDate = state.ModifiedDate
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (c *certificateResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	tflog.Debug(ctx, "CCM Certificate resource ValidateConfig")

	var config certificateResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(config.validateKeyTypeAndSize()...)
	resp.Diagnostics.Append(config.validateSubjectAndSANs(ctx)...)
}

func (c *certificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Contract ID under which this certificate will be created.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Group that will be associated with the new certificate when it gets created. Required for creation.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("grp_")),
					modifiers.StringRequiredForCreate(),
				},
			},
			"base_name": schema.StringAttribute{
				Optional:    true,
				Description: "The base name for the certificate. If not provided, the name will be auto-generated by the CCM API.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The certificate name.",
			},
			"key_type": schema.StringAttribute{
				Required:    true,
				Description: "The key type for a certificate. Valid values are 'RSA' or 'ECDSA'",
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"RSA", "ECDSA"}...),
				},
			},
			"key_size": schema.StringAttribute{
				Required:    true,
				Description: "The key size for a certificate. Valid value for key type RSA: '2048'. Valid value for key type ECDSA: 'P-256'.",
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"2048", "P-256"}...),
				},
			},
			"secure_network": schema.StringAttribute{
				Required:    true,
				Description: "Secure network type to use for the certificate. The only valid value is 'ENHANCED_TLS'",
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ENHANCED_TLS"}...),
				},
			},
			"sans": schema.SetAttribute{
				Required:    true,
				Description: "The list of Subject Alternative Names (SANs) for the certificate.",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					modifiers.PreventSetUpdate(),
				},
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 100),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(domainNameRegex,
							"must be a valid domain name with all letters lowercase, optionally starting with '*.' for wildcard")),
				},
			},
			"subject": subjectSchema(),
			// TODO: For next story - implement renew_before_expiration_days logic.
			// "renew_before_expiration_days": schema.Int64Attribute{
			// 	Optional:    true,
			// },
			// "needs_renewal": schema.BoolAttribute{
			// 	Computed:    true,
			// },
			"certificate_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier assigned to the newly created CCM certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"certificate_type": schema.StringAttribute{
				Computed:    true,
				Description: "Certificate type. Defaults to 'THIRD_PARTY'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Account associated with 'contract_id'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date the certificate was created in UTC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who created the certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Description: "The status of the certificate. Can be one of 'ACTIVE', 'CSR_READY', 'READY_FOR_USE'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"csr_pem": schema.StringAttribute{
				Computed:    true,
				Description: "CSR PEM content generated by Akamai.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"csr_expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when CSR will expire and a signed certificate uploaded based on that CSR will NOT be accepted beyond this date.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func subjectType() map[string]attr.Type {
	return subjectSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func subjectSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Subject fields as defined in X.509 certificates (RFC 5280). At least one of the inner fields must be specified.",
		PlanModifiers: []planmodifier.Object{
			modifiers.PreventObjectUpdate(),
		},
		Attributes: map[string]schema.Attribute{
			"common_name": schema.StringAttribute{
				Optional:    true,
				Description: "Fully qualified domain name (FQDN) or other name associated with the subject. If specified, this value must also be included in the SANs list.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(domainNameRegex,
						"must be a valid domain name with all letters lowercase, optionally starting with '*.' for wildcard"),
				},
			},
			"organization": schema.StringAttribute{
				Optional:    true,
				Description: "Legal name of the organization.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "cannot be empty or whitespace"),
				},
			},
			"country": schema.StringAttribute{
				Optional:    true,
				Description: "Two-letter ISO 3166 country code.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 2),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "cannot be empty or whitespace"),
				},
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Description: "Full name of the state or province.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "cannot be empty or whitespace"),
				},
			},
			"locality": schema.StringAttribute{
				Optional:    true,
				Description: "City or locality name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "cannot be empty or whitespace"),
				},
			},
		},
	}
}

func (c *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CCM Certificate resource Create")

	var plan certificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sans []string
	resp.Diagnostics.Append(plan.SANs.ElementsAs(ctx, &sans, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := ccm.CreateCertificateRequest{
		ContractID: strings.TrimPrefix(plan.ContractID.ValueString(), "ctr_"),
		GroupID:    strings.TrimPrefix(plan.GroupID.ValueString(), "grp_"),
		Body: ccm.CreateCertificateRequestBody{
			CertificateName: plan.BaseName.ValueString(),
			KeyType:         ccm.CryptographicAlgorithm(plan.KeyType.ValueString()),
			KeySize:         ccm.KeySize(plan.KeySize.ValueString()),
			SecureNetwork:   ccm.SecureNetwork(plan.SecureNetwork.ValueString()),
			SANs:            sans,
		},
	}

	if !plan.Subject.IsNull() {
		var subject subjectModel
		resp.Diagnostics.Append(plan.Subject.As(ctx, &subject, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Body.Subject = &ccm.Subject{
			CommonName:   subject.CommonName.ValueString(),
			Organization: subject.Organization.ValueString(),
			Country:      subject.Country.ValueString(),
			State:        subject.State.ValueString(),
			Locality:     subject.Locality.ValueString(),
		}
	}

	client := Client(c.meta)
	cert, err := client.CreateCertificate(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create CCM Certificate", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("CCM Certificate with ID %s created", cert.Certificate.CertificateID))

	resp.Diagnostics.Append(plan.populateCertificateFields(ctx, cert.Certificate, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "CCM Certificate resource Read")

	var state certificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "certificate_id", state.CertificateID.ValueString())

	client := Client(c.meta)
	cert, err := client.GetCertificate(ctx, ccm.GetCertificateRequest{
		CertificateID: state.CertificateID.ValueString(),
	})
	if err != nil && errors.Is(err, ccm.ErrCertificateNotFound) {
		tflog.Debug(ctx, fmt.Sprintf("CCM Certificate with ID %s not found, removing from state",
			state.CertificateID.ValueString()))
		resp.Diagnostics.AddWarning("Certificate not found",
			"The certificate was not found on the server. The resource will be removed from the state.")
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to get CCM Certificate", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateCertificateFields(ctx, *cert, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "CCM Certificate resource Update")

	var plan certificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "certificate_id", plan.CertificateID.ValueString())

	// TODO: For next story - update only if 'renew_before_expiration_days' is false.
	// Add support for 'renew_before_expiration_days' logic.
	tflog.Debug(ctx, "'base_name' change detected, updating the certificate name")
	client := Client(c.meta)
	cert, err := client.PatchCertificate(ctx, ccm.PatchCertificateRequest{
		CertificateID: plan.CertificateID.ValueString(),
		// If base_name is Null, it must be used as empty string to reset the name to the default value.
		CertificateName: ptr.To(plan.BaseName.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to update CCM Certificate", err.Error())
		return
	}

	tflog.Debug(ctx, "'base_name' updated to "+plan.BaseName.ValueString())

	resp.Diagnostics.Append(plan.populateCertificateFields(ctx, *cert, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *certificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "CCM Certificate resource Delete")

	var state certificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "certificate_id", state.CertificateID.ValueString())

	client := Client(c.meta)
	if err := client.DeleteCertificate(ctx, ccm.DeleteCertificateRequest{
		CertificateID: state.CertificateID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to delete CCM Certificate", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("CCM Certificate with ID %s deleted", state.CertificateID.ValueString()))
}

// ImportState implements resource's ImportState method. The import ID has a format of <certificateID[,groupID]>
func (c *certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "CCM Certificate resource ImportState")

	parts, err := text.ImportIDSplitter("certificateID[,groupID]").
		AcceptLen(1).
		AcceptLen(2).
		Split(req.ID)

	if err != nil {
		resp.Diagnostics.AddError("Incorrect import ID: ", err.Error())
		return
	}
	certificateID := parts[0]

	client := Client(c.meta)
	cert, err := client.GetCertificate(ctx, ccm.GetCertificateRequest{
		CertificateID: certificateID,
	})
	if err != nil && errors.Is(err, ccm.ErrCertificateNotFound) {
		tflog.Debug(ctx, fmt.Sprintf("CCM Certificate with ID %s not found, cannot be imported", certificateID))
		resp.Diagnostics.AddError("Cannot import non-existent remote object",
			"The certificate was not found on the server. Please verify the Certificate ID is correct.")
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to get CCM Certificate", err.Error())
		return
	}

	// Check if CertificateName has a suffix of format ".rotated.{YYYY-MM-DD}".
	// If so, strip that part to set the base_name.
	baseName := extractBaseName(cert.CertificateName)
	tflog.Debug(ctx, fmt.Sprintf("Setting base_name to %s", baseName))

	state := certificateResourceModel{
		CertificateID: types.StringValue(certificateID),
		BaseName:      types.StringValue(baseName),
		Subject:       types.ObjectNull(subjectType()),
		SANs:          types.SetNull(types.StringType),
	}
	if len(parts) == 2 {
		state.GroupID = types.StringValue(parts[1])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func extractBaseName(name string) string {
	parts := strings.Split(name, renewedNameSuffix)
	if len(parts) != 2 || parts[0] == "" {
		// The name does not include the rotated suffix or starts with the suffix.
		return name
	}
	if _, err := time.Parse(renewedNameDateLayout, parts[1]); err == nil {
		// Valid date part, return the base name.
		return parts[0]
	}
	// Invalid date part, return the original name.
	return name
}

func isEmptySubject(subject ccm.Subject) bool {
	return subject.CommonName == "" && subject.Organization == "" && subject.Country == "" &&
		subject.State == "" && subject.Locality == ""
}
