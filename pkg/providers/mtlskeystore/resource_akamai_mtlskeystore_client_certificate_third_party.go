package mtlskeystore

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &clientCertificateThirdPartyResource{}
	_ resource.ResourceWithImportState = &clientCertificateThirdPartyResource{}
	_ resource.ResourceWithModifyPlan  = &clientCertificateThirdPartyResource{}
	_ resource.ResourceWithConfigure   = &clientCertificateThirdPartyResource{}
)

type clientCertificateThirdPartyResource struct {
	meta meta.Meta
}

// NewClientCertificateThirdPartyResource returns new akamai_mtlskeystore_client_certificate_third_party resource.
func NewClientCertificateThirdPartyResource() resource.Resource {
	return &clientCertificateThirdPartyResource{}
}

type (
	clientCertificateThirdPartyResourceModel struct {
		CertificateName    types.String `tfsdk:"certificate_name"`
		CertificateID      types.Int64  `tfsdk:"certificate_id"`
		ContractID         types.String `tfsdk:"contract_id"`
		Geography          types.String `tfsdk:"geography"`
		GroupID            types.Int64  `tfsdk:"group_id"`
		KeyAlgorithm       types.String `tfsdk:"key_algorithm"`
		NotificationEmails types.List   `tfsdk:"notification_emails"`
		SecureNetwork      types.String `tfsdk:"secure_network"`
		Subject            types.String `tfsdk:"subject"`
		Versions           types.Map    `tfsdk:"versions"`
	}

	clientCertificateVersionModel struct {
		Version                  types.Int64  `tfsdk:"version"`
		Status                   types.String `tfsdk:"status"`
		ExpiryDate               types.String `tfsdk:"expiry_date"`
		Issuer                   types.String `tfsdk:"issuer"`
		KeyAlgorithm             types.String `tfsdk:"key_algorithm"`
		CertificateSubmittedBy   types.String `tfsdk:"certificate_submitted_by"`
		CertificateSubmittedDate types.String `tfsdk:"certificate_submitted_date"`
		CreatedBy                types.String `tfsdk:"created_by"`
		CreatedDate              types.String `tfsdk:"created_date"`
		DeleteRequestedDate      types.String `tfsdk:"delete_requested_date"`
		IssuedDate               types.String `tfsdk:"issued_date"`
		EllipticCurve            types.String `tfsdk:"elliptic_curve"`
		KeySizeInBytes           types.String `tfsdk:"key_size_in_bytes"`
		ScheduledDeleteDate      types.String `tfsdk:"scheduled_delete_date"`
		SignatureAlgorithm       types.String `tfsdk:"signature_algorithm"`
		Subject                  types.String `tfsdk:"subject"`
		VersionGUID              types.String `tfsdk:"version_guid"`
		CertificateBlock         types.Object `tfsdk:"certificate_block"`
		CSRBlock                 types.Object `tfsdk:"csr_block"`
	}

	certificateBlockResourceModel struct {
		Certificate types.String `tfsdk:"certificate"`
		TrustChain  types.String `tfsdk:"trust_chain"`
	}

	csrBlockResourceModel struct {
		CSR          types.String `tfsdk:"csr"`
		KeyAlgorithm types.String `tfsdk:"key_algorithm"`
	}
)

var (
	errorOneVersionWithDeletePendingStatus = errors.New("one version with `DELETE_PENDING` status found")
	cnRegex                                = regexp.MustCompile(`/CN=[^/]{1,64}/`)
	versionsObjectType                     = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"version":                    types.Int64Type,
			"status":                     types.StringType,
			"expiry_date":                types.StringType,
			"issuer":                     types.StringType,
			"key_algorithm":              types.StringType,
			"certificate_submitted_by":   types.StringType,
			"certificate_submitted_date": types.StringType,
			"created_by":                 types.StringType,
			"created_date":               types.StringType,
			"delete_requested_date":      types.StringType,
			"issued_date":                types.StringType,
			"elliptic_curve":             types.StringType,
			"key_size_in_bytes":          types.StringType,
			"scheduled_delete_date":      types.StringType,
			"signature_algorithm":        types.StringType,
			"subject":                    types.StringType,
			"version_guid":               types.StringType,
			"certificate_block": types.ObjectType{
				AttrTypes: certificateBlockType(),
			},
			"csr_block": types.ObjectType{
				AttrTypes: csrBlockType(),
			},
		},
	}
)

func certificateBlockType() map[string]attr.Type {
	return certificateBlockSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func csrBlockType() map[string]attr.Type {
	return csrBlockSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func (r *clientCertificateThirdPartyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_client_certificate_third_party"
}

func (r *clientCertificateThirdPartyResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"certificate_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the client certificate. Must be between 1 and 64 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"certificate_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier of the client certificate.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "The contract assigned to the client certificate. Must have a length of at least 1.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"geography": schema.StringAttribute{
				Required:    true,
				Description: "Specifies the type of network to deploy the client certificate. Possible values: `CORE`, `RUSSIA_AND_CORE`, or `CHINA_AND_CORE`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"CORE",
						"RUSSIA_AND_CORE",
						"CHINA_AND_CORE",
					),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "The group assigned to the client certificate. Must be greater than or equal to 0.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"key_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The cryptographic algorithm used for key generation. Possible values: `RSA` or `ECDSA`. The default is `RSA`.",
				Default:     stringdefault.StaticString(string(mtlskeystore.KeyAlgorithmRSA)),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"RSA",
						"ECDSA",
					),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"notification_emails": schema.ListAttribute{
				Required:    true,
				Description: "The email addresses to notify for client certificate-related issues. Must have at least one email address.",
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"secure_network": schema.StringAttribute{
				Required:    true,
				Description: "Identifies the network deployment type. Possible values: `STANDARD_TLS` or `ENHANCED_TLS`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"STANDARD_TLS",
						"ENHANCED_TLS",
					),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"subject": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Specifies the client certificate. If provided, the `CN` attribute is required and cannot exceed 64 characters. When `null`, the subject is constructed with the following format: `/C=US/O=Akamai Technologies, Inc./OU={vcdId} {contractId} {groupId}/CN={certificateName}/`.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(cnRegex, "The `subject` must contain a valid `CN` attribute with a maximum length of 64 characters."),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"versions": versionsSchema(),
		},
	}
}

func versionsSchema() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Description: "A map of client certificate versions as a value and user defined identifier as a key. Each version represents a specific iteration of the client certificate.",
		Required:    true,
		Validators: []validator.Map{
			mapvalidator.SizeBetween(1, 5),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"version": schema.Int64Attribute{
					Description: "The unique identifier of the client certificate version.",
					Computed:    true,
				},
				"status": schema.StringAttribute{
					Description: "The client certificate version status. Possible values: `AWAITING_SIGNED_CERTIFICATE`, `DEPLOYMENT_PENDING`, `DEPLOYED`, or `DELETE_PENDING`.",
					Computed:    true,
				},
				"expiry_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating when the client certificate version expires.",
					Computed:    true,
				},
				"issuer": schema.StringAttribute{
					Description: "The signing entity of the client certificate version.",
					Computed:    true,
				},
				"key_algorithm": schema.StringAttribute{
					Description: "Identifies the client certificate version's encryption algorithm. Supported values are `RSA` and `ECDSA`.",
					Computed:    true,
				},
				"certificate_submitted_by": schema.StringAttribute{
					Description: "The user who uploaded the THIRD_PARTY client certificate version. Appears as null if not specified.",
					Computed:    true,
				},
				"certificate_submitted_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating when the THIRD_PARTY signer client certificate version was uploaded. Appears as null if not specified.",
					Computed:    true,
				},
				"created_by": schema.StringAttribute{
					Description: "The user who created the client certificate version.",
					Computed:    true,
				},
				"created_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating the client certificate version's creation.",
					Computed:    true,
				},
				"delete_requested_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating the client certificate version's deletion request. Appears as null if there's no request.",
					Computed:    true,
				},
				"issued_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating the client certificate version's availability.",
					Computed:    true,
				},
				"elliptic_curve": schema.StringAttribute{
					Description: "Specifies the key elliptic curve when key algorithm `ECDSA` is used.",
					Computed:    true,
				},
				"key_size_in_bytes": schema.StringAttribute{
					Description: "The private key length of the client certificate version when key algorithm `RSA` is used.",
					Computed:    true,
				},
				"scheduled_delete_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating the client certificate version's deletion. Appears as null if there's no request.",
					Computed:    true,
				},
				"signature_algorithm": schema.StringAttribute{
					Description: "Specifies the algorithm that secures the data exchange between the edge server and origin.",
					Computed:    true,
				},
				"subject": schema.StringAttribute{
					Description: "The public key's entity stored in the client certificate version's subject public key field.",
					Computed:    true,
				},
				"version_guid": schema.StringAttribute{
					Description: "Unique identifier for the client certificate version. Use it to configure mutual authentication (mTLS) sessions between the origin and edge servers in Property Manager's Mutual TLS Origin Keystore behavior.",
					Computed:    true,
				},
				"certificate_block": certificateBlockSchema(),
				"csr_block":         csrBlockSchema(),
			},
		},
	}
}

func certificateBlockSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Details of the certificate block for the client certificate version.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"certificate": schema.StringAttribute{
				Description: "A text representation of the client certificate in PEM format.",
				Computed:    true,
			},
			"trust_chain": schema.StringAttribute{
				Description: "A text representation of the trust chain in PEM format.",
				Computed:    true,
			},
		},
	}
}

func csrBlockSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Details of the Certificate Signing Request (CSR) for the client certificate version.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"csr": schema.StringAttribute{
				Description: "Text of the certificate signing request.",
				Computed:    true,
			},
			"key_algorithm": schema.StringAttribute{
				Description: "Identifies the client certificate's encryption algorithm.",
				Computed:    true,
			},
		},
	}
}

func (r *clientCertificateThirdPartyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

func (r *clientCertificateThirdPartyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "Modifying Client Certificate Third Party Resource Plan")
	var state, plan *clientCertificateThirdPartyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if modifiers.IsDelete(req) {
		var stateVersions map[string]clientCertificateVersionModel
		if diags := state.Versions.ElementsAs(ctx, &stateVersions, false); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		var versionsToRemove []int64
		for _, version := range stateVersions {
			versionsToRemove = append(versionsToRemove, version.Version.ValueInt64())
		}

		err := checkStatus(ctx, Client(r.meta), state.CertificateID.ValueInt64(), versionsToRemove)
		if err != nil {
			resp.Diagnostics.AddWarning("Versions deletion", err.Error())
			return
		}
	}

	// Updating 'subject' is not allowed and must be handled
	// on plan modification level as it is optional and PreventStringUpdate() modifier is not enough.
	if modifiers.IsUpdate(req) {
		if !state.Subject.Equal(plan.Subject) {
			resp.Diagnostics.AddAttributeError(path.Root("subject"), "Cannot Update 'subject'",
				"The `subject` attribute cannot be updated after the resource has been created.")
			return
		}
	}

	// If the only attributes that have changed are 'notification_emails' and/or 'certificate_name',
	// suppress changes to the 'versions' attribute by using the state value.
	if modifiers.IsUpdate(req) {
		if state.haveNotificationEmailsChanged(plan) || state.hasNameChanged(plan) {
			tflog.Debug(ctx, "only 'notification_emails' or 'certificate_name' changed, using state value for versions instead")
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("versions"), state.Versions)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	// Copy state versions to plan versions for the keys that are the same.
	// This will suppress the changes for existing versions if other versions are being removed or created.
	var numberOfNewVersions int
	if modifiers.IsUpdate(req) {
		var stateVersions map[string]clientCertificateVersionModel
		if diags := state.Versions.ElementsAs(ctx, &stateVersions, false); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		var planVersions map[string]clientCertificateVersionModel
		if diags := plan.Versions.ElementsAs(ctx, &planVersions, false); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		var versionsToRemove []int64
		for k := range stateVersions {
			if _, ok := planVersions[k]; !ok {
				versionsToRemove = append(versionsToRemove, stateVersions[k].Version.ValueInt64())
			}
		}
		if len(versionsToRemove) > 0 {
			err := checkStatus(ctx, Client(r.meta), state.CertificateID.ValueInt64(), versionsToRemove)
			if err != nil {
				resp.Diagnostics.AddWarning("Versions deletion", err.Error())
				return
			}
		}

		for k := range planVersions {
			if stateVersion, ok := stateVersions[k]; ok {
				tflog.Debug(ctx, fmt.Sprintf("Using state values for version with key '%s'", k))
				planVersions[k] = stateVersion
			} else {
				numberOfNewVersions++
			}
		}

		versionsValue, diags := types.MapValueFrom(ctx, versionsObjectType, planVersions)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("versions"), versionsValue)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Adding a warning if there are multiple versions to create or update.
	var planVersionsToCreate map[string]clientCertificateVersionModel
	if modifiers.IsCreate(req) {
		if diags := plan.Versions.ElementsAs(ctx, &planVersionsToCreate, false); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	if len(planVersionsToCreate) > 1 || numberOfNewVersions > 1 {
		tflog.Warn(ctx, "Versions will be sorted by key names in ascending order.")
		resp.Diagnostics.AddWarning("Adding multiple versions", "While adding multiple new versions, they will be sorted by key names in ascending order.")
	}
}

func (m *clientCertificateThirdPartyResourceModel) haveNotificationEmailsChanged(plan *clientCertificateThirdPartyResourceModel) bool {
	return !m.NotificationEmails.Equal(plan.NotificationEmails)
}

func (m *clientCertificateThirdPartyResourceModel) hasNameChanged(plan *clientCertificateThirdPartyResourceModel) bool {
	return !m.CertificateName.Equal(plan.CertificateName)
}

func (r *clientCertificateThirdPartyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Client Certificate Third Party Resource")

	var plan clientCertificateThirdPartyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.create(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Creating Client Certificate Third Party failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("failed to set data", fmt.Sprintf("Client Certificate has been created successfully, but setting data. "+
			"If you wish to manage the same client certificate, please import it and apply your configuration again. "+
			"Client Certificate ID: %d,%d,%s", plan.CertificateID, plan.GroupID.ValueInt64(), strings.TrimPrefix(plan.ContractID.ValueString(), "ctr_")))
	}
}

func (r *clientCertificateThirdPartyResource) create(ctx context.Context, plan *clientCertificateThirdPartyResourceModel) error {
	client := Client(r.meta)

	var notificationEmails []string
	diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
	if diags.HasError() {
		return fmt.Errorf("failed to get notification emails: %v", diags)
	}

	var versions map[string]clientCertificateVersionModel

	diags = plan.Versions.ElementsAs(ctx, &versions, false)
	if diags.HasError() {
		return fmt.Errorf("failed to get versions: %v", diags)
	}

	contract := strings.TrimPrefix(plan.ContractID.ValueString(), "ctr_")
	request := mtlskeystore.CreateClientCertificateRequest{
		CertificateName:    plan.CertificateName.ValueString(),
		ContractID:         contract,
		Geography:          mtlskeystore.Geography(plan.Geography.ValueString()),
		GroupID:            plan.GroupID.ValueInt64(),
		NotificationEmails: notificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(plan.SecureNetwork.ValueString()),
		Subject:            plan.Subject.ValueStringPointer(),
		Signer:             mtlskeystore.SignerThirdParty,
	}

	keyAlgorithm := mtlskeystore.CryptographicAlgorithm(plan.KeyAlgorithm.ValueString())
	if keyAlgorithm != "" {
		request.KeyAlgorithm = &keyAlgorithm
	}

	createClientCertificateResponse, err := client.CreateClientCertificate(ctx, request)
	if err != nil {
		return fmt.Errorf("creating client certificate third party: %w", err)
	}

	versionsKeys := extractVersionKeys(versions)
	if len(versionsKeys) > 1 {
		slices.Sort(versionsKeys)
		tflog.Debug(ctx, fmt.Sprintf("Sorted versions keys for creation: %v", versionsKeys))
		err = createMissingVersionsExplicitly(ctx, client, createClientCertificateResponse.CertificateID, versionsKeys[1:])
		if err != nil {
			return fmt.Errorf("failed to create missing versions: %w. "+
				"Client Certificate has been created successfully, but other operations failed. "+
				"If you wish to manage the same client certificate, please import it and apply your configuration again. "+
				"Client Certificate ID: %d,%d,%s", err, createClientCertificateResponse.CertificateID, plan.GroupID.ValueInt64(), contract)
		}
	}

	clientCertificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: createClientCertificateResponse.CertificateID,
	})
	if err != nil {
		return fmt.Errorf("failed to get client certificate: %w. "+
			"Client Certificate has been created successfully, but other operations failed. "+
			"If you wish to manage the same client certificate, please import it and apply your configuration again. "+
			"Client Certificate ID: %d,%d,%s", err, createClientCertificateResponse.CertificateID, plan.GroupID.ValueInt64(), contract)
	}

	diags = plan.setClientCertificateData(ctx, clientCertificate)
	if diags.HasError() {
		return fmt.Errorf("failed to set data: %v. "+
			"Client Certificate has been created successfully, but other operations failed. "+
			"If you wish to manage the same client certificate, please import it and apply your configuration again. "+
			"Client Certificate ID: %d,%d,%s", diags, createClientCertificateResponse.CertificateID, plan.GroupID.ValueInt64(), contract)
	}

	clientCertificateVersions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               createClientCertificateResponse.CertificateID,
		IncludeAssociatedProperties: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get client certificate versions: %w. "+
			"Client Certificate has been created successfully, but other operations failed. "+
			"If you wish to manage the same client certificate, please import it and apply your configuration again. "+
			"Client Certificate ID: %d,%d,%s", err, createClientCertificateResponse.CertificateID, plan.GroupID.ValueInt64(), contract)
	}

	diags = plan.setVersionsData(ctx, clientCertificateVersions, versionsKeys, nil)
	if diags.HasError() {
		return fmt.Errorf("failed to set versions data: %v. "+
			"Client Certificate has been created successfully, but other operations failed. "+
			"If you wish to manage the same client certificate, please import it and apply your configuration again. "+
			"Client Certificate ID: %d,%d,%s", diags, createClientCertificateResponse.CertificateID, plan.GroupID.ValueInt64(), contract)
	}

	return nil
}

func createMissingVersionsExplicitly(ctx context.Context, client mtlskeystore.MTLSKeystore, certificateID int64, versions []string) error {
	for _, version := range versions {
		_, err := client.RotateClientCertificateVersion(ctx, mtlskeystore.RotateClientCertificateVersionRequest{
			CertificateID: certificateID,
		})
		if err != nil {
			return fmt.Errorf("failed to create version %s: %w", version, err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Successfully created version %s", version))
	}
	return nil
}

func (r *clientCertificateThirdPartyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Client Certificate Third Party Resource")

	var state clientCertificateThirdPartyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.read(ctx, &state)
	if errors.Is(err, mtlskeystore.ErrClientCertificateNotFound) {
		tflog.Debug(ctx, "Client Certificate Third Party Resource not found, removing from state")
		resp.Diagnostics.AddWarning("Resource Removal", "The client certificate was not found on the server. The resource will be removed from the state.")
		resp.State.RemoveResource(ctx)
		return
	}
	if errors.Is(err, errorOneVersionWithDeletePendingStatus) {
		tflog.Debug(ctx, "Client Certificate Third Party Resource has one version with DELETE_PENDING status, removing from state")
		resp.Diagnostics.AddWarning("Resource Removal", "The last version of the Client Certificate is in `DELETE_PENDING` status. The resource will be removed from the state.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Reading Client Certificate Third Party Resource failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *clientCertificateThirdPartyResource) read(ctx context.Context, data *clientCertificateThirdPartyResourceModel) error {
	client := Client(r.meta)
	clientCertificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: data.CertificateID.ValueInt64(),
	})
	if err != nil {
		return err
	}

	diags := data.setClientCertificateData(ctx, clientCertificate)
	if diags.HasError() {
		return fmt.Errorf("failed to set data: %v", diags)
	}

	clientCertificateVersions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               data.CertificateID.ValueInt64(),
		IncludeAssociatedProperties: true,
	})
	if err != nil {
		return err
	}

	// Read the state versions only if they are known.
	// This excludes the case when the Read is invoked after the Import.
	// In case of Import, dataVersions are nil.
	var dataVersions map[string]clientCertificateVersionModel
	if !data.Versions.IsUnknown() {
		diags = data.Versions.ElementsAs(ctx, &dataVersions, false)
		if diags.HasError() {
			return fmt.Errorf("failed to get versions: %v", diags)
		}
	}
	stateVersions := mapKeyToVersionGUID(dataVersions)
	tflog.Debug(ctx, fmt.Sprintf("State versions keys: %v", stateVersions))

	diags = data.setVersionsData(ctx, clientCertificateVersions, nil, stateVersions)
	if diags.HasError() {
		return fmt.Errorf("failed to set versions data: %v", diags)
	}

	return nil
}

func (r *clientCertificateThirdPartyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Client Certificate Third Party Resource")
	var plan, oldState clientCertificateThirdPartyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.update(ctx, &plan, &oldState); err != nil {
		resp.Diagnostics.AddError("Updating Client Certificate Third Party failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clientCertificateThirdPartyResource) update(ctx context.Context, plan, oldState *clientCertificateThirdPartyResourceModel) error {
	client := Client(r.meta)

	if oldState.hasNameChanged(plan) || oldState.haveNotificationEmailsChanged(plan) {
		patchReq := mtlskeystore.PatchClientCertificateRequest{
			CertificateID: oldState.CertificateID.ValueInt64(),
			Body:          mtlskeystore.PatchClientCertificateRequestBody{},
		}
		if oldState.hasNameChanged(plan) {
			tflog.Debug(ctx, fmt.Sprintf("Updating certificate name from '%s' to '%s'", oldState.CertificateName.ValueString(), plan.CertificateName.ValueString()))
			patchReq.Body.CertificateName = plan.CertificateName.ValueStringPointer()
		}
		if oldState.haveNotificationEmailsChanged(plan) {
			tflog.Debug(ctx, "Updating notification emails")
			var notificationEmails []string
			diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
			if diags.HasError() {
				return fmt.Errorf("failed to get notification emails: %v", diags)
			}
			patchReq.Body.NotificationEmails = notificationEmails
		}

		if err := client.PatchClientCertificate(ctx, patchReq); err != nil {
			return fmt.Errorf("failed to update client certificate: %w", err)
		}
	}

	var planVersions map[string]clientCertificateVersionModel
	var oldStateVersions map[string]clientCertificateVersionModel

	if diags := plan.Versions.ElementsAs(ctx, &planVersions, false); diags.HasError() {
		return fmt.Errorf("failed to get plan versions: %v", diags)
	}
	if diags := oldState.Versions.ElementsAs(ctx, &oldStateVersions, false); diags.HasError() {
		return fmt.Errorf("failed to get old state versions: %v", diags)
	}

	planKeys := extractVersionKeys(planVersions)
	stateKeys := extractVersionKeys(oldStateVersions)
	slices.Sort(planKeys)

	versionsToRemove, versionsToAdd, stateVersions := filterVersionsToRemoveAndAdd(stateKeys, planKeys, oldStateVersions)
	tflog.Debug(ctx, "Changes to be performed during the update", map[string]any{
		"versionsToRemove": versionsToRemove,
		"versionsToAdd":    versionsToAdd,
		"stateVersions":    stateVersions,
	})

	if len(versionsToRemove) != 0 || len(versionsToAdd) != 0 {
		if len(versionsToRemove) < len(stateKeys) {
			err := someVersionDeletionAndRotation(ctx, client, versionsToAdd, versionsToRemove, stateVersions, oldState.CertificateID.ValueInt64())
			if err != nil {
				return err
			}
		} else {
			err := allVersionDeletionAndRotation(ctx, client, versionsToAdd, versionsToRemove, stateVersions, oldState.CertificateID.ValueInt64())
			if err != nil {
				return err
			}
		}
	}

	clientCertificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: plan.CertificateID.ValueInt64(),
	})
	if err != nil {
		return fmt.Errorf("failed to get updated client certificate: %w", err)
	}

	diags := plan.setClientCertificateData(ctx, clientCertificate)
	if diags.HasError() {
		return fmt.Errorf("failed to set updated client certificate data: %v", diags)
	}

	clientCertificateVersions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               plan.CertificateID.ValueInt64(),
		IncludeAssociatedProperties: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get client certificate versions: %w", err)
	}

	diags = plan.setVersionsData(ctx, clientCertificateVersions, planKeys, stateVersions)
	if diags.HasError() {
		return fmt.Errorf("failed to set versions data: %v", diags)
	}

	return nil
}

// allVersionDeletionAndRotation is responsible for managing the deletion and rotation of client certificate versions.
// It ensures that all versions except the first one are deleted to prevent deletion of the whole client certificate, rotates first version to add, deletes the first version and then rotates the remaining versions to add.
func allVersionDeletionAndRotation(ctx context.Context, client mtlskeystore.MTLSKeystore, versionsToAdd []string, versionsToRemove []int64, stateVersions map[string]string, certificateID int64) error {
	slices.Sort(versionsToRemove)
	for _, version := range versionsToRemove[1:] {
		err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
			CertificateID: certificateID,
			Version:       version,
		})
		if err != nil {
			return fmt.Errorf("failed to delete version %d: %w", version, err)
		}
	}

	firstVersionToRotate, err := client.RotateClientCertificateVersion(ctx, mtlskeystore.RotateClientCertificateVersionRequest{
		CertificateID: certificateID,
	})
	if err != nil {
		return fmt.Errorf("failed to create version %s: %w", versionsToAdd[0], err)
	}
	stateVersions[firstVersionToRotate.VersionGUID] = versionsToAdd[0]

	err = client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
		CertificateID: certificateID,
		Version:       versionsToRemove[0],
	})
	if err != nil {
		return fmt.Errorf("failed to delete version %d: %w", versionsToRemove[0], err)
	}

	for _, version := range versionsToAdd[1:] {
		newVersion, err := client.RotateClientCertificateVersion(ctx, mtlskeystore.RotateClientCertificateVersionRequest{
			CertificateID: certificateID,
		})
		if err != nil {
			return fmt.Errorf("failed to create version %s: %w", version, err)
		}
		stateVersions[newVersion.VersionGUID] = version
	}

	return nil
}

// someVersionDeletionAndRotation is responsible for managing the deletion and rotation of client certificate versions.
func someVersionDeletionAndRotation(ctx context.Context, client mtlskeystore.MTLSKeystore, versionsToAdd []string, versionsToRemove []int64, stateVersions map[string]string, certificateID int64) error {
	for _, version := range versionsToRemove {
		err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
			CertificateID: certificateID,
			Version:       version,
		})
		if err != nil {
			return fmt.Errorf("failed to delete version %d: %w", version, err)
		}
	}
	for _, version := range versionsToAdd {
		newVersion, err := client.RotateClientCertificateVersion(ctx, mtlskeystore.RotateClientCertificateVersionRequest{
			CertificateID: certificateID,
		})
		if err != nil {
			return fmt.Errorf("failed to create version %s: %w", version, err)
		}
		stateVersions[newVersion.VersionGUID] = version
	}

	return nil
}

func (r *clientCertificateThirdPartyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Client Certificate Third Party Resource")

	var state clientCertificateThirdPartyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(r.meta)

	var stateVersions map[string]clientCertificateVersionModel
	diags := state.Versions.ElementsAs(ctx, &stateVersions, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	for version := range stateVersions {
		if stateVersions[version].Status.ValueString() != "DELETE_PENDING" {
			err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
				CertificateID: state.CertificateID.ValueInt64(),
				Version:       stateVersions[version].Version.ValueInt64(),
			})
			if err != nil {
				resp.Diagnostics.AddError(
					"Deleting Client Certificate Version failed",
					fmt.Sprintf("Failed to delete version %s: %v. ", version, err),
				)
				return
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Version %s is already in DELETE_PENDING status, skipping deletion", version))
		}
	}
}

func (r *clientCertificateThirdPartyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Client Certificate Third Party Resource")

	importID := strings.TrimSpace(req.ID)
	parts := strings.Split(importID, ",")

	if len(parts) != 3 && len(parts) != 1 {
		resp.Diagnostics.AddError("Incorrect import ID: ", "you need to provide an importID in the format 'certificateID,[groupID,contractID]'. Where certificateID is required and groupID and contractID are optional")
		return
	}

	certificateID, err := parseCertificateID(parts[0])
	if err != nil {
		resp.Diagnostics.AddError("Invalid Certificate ID", err.Error())
		return
	}

	data := clientCertificateThirdPartyResourceModel{
		CertificateName:    types.StringUnknown(),
		CertificateID:      types.Int64Value(certificateID),
		Geography:          types.StringUnknown(),
		NotificationEmails: types.ListUnknown(types.StringType),
		SecureNetwork:      types.StringUnknown(),
		Versions:           types.MapUnknown(versionsObjectType),
	}

	// API call is needed to populate subject from server, and extract contract and group ID from it
	client := Client(r.meta)
	certificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: certificateID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Get Client Certificate", err.Error())
		return
	}

	var contractID, groupID string
	if len(parts) == 3 {
		contractID, groupID = parts[2], parts[1]
	} else {
		ok, subjectParts := subjectContainsContractAndGroup(certificate.Subject)
		if !ok {
			resp.Diagnostics.AddError("Incorrect import ID: ", fmt.Sprintf("since it is not possible to extract contract and group from certificate subject, "+missedContractAndGroupError))
			return
		}
		contractID, groupID, err = extractContractAndGroup(certificate.Subject, subjectParts)
		if err != nil {
			resp.Diagnostics.AddError("Unable to extract contract ID or group ID", err.Error())
			return
		}
	}

	if diags := data.assignGroupAndContractThirdParty(contractID, groupID); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if onlyOneVersionPendingDelete(ctx, client, certificateID, true, mtlskeystore.SignerThirdParty) {
		resp.Diagnostics.AddError("Certificate in Delete Pending State", fmt.Sprintf("The client certificate %d has only one version and it's in `DELETE_PENDING` state. In order to import this resource, rotate this client certificate first", certificateID))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *clientCertificateThirdPartyResourceModel) setVersionsData(ctx context.Context, clientCertificateVersions *mtlskeystore.ListClientCertificateVersionsResponse, planVersions []string, stateVersions map[string]string) diag.Diagnostics {
	slices.Reverse(planVersions)
	alreadySetVersions := make(map[string]clientCertificateVersionModel)
	for v, version := range clientCertificateVersions.Versions {
		if version.VersionAlias != nil {
			// Skip the version with an alias
			continue
		}
		versionsModel, diags := createClientVersionModel(ctx, version)
		if diags.HasError() {
			return diags
		}

		if versionKey, exists := stateVersions[version.VersionGUID]; exists {
			alreadySetVersions[versionKey] = versionsModel
		} else {
			// If no version keys are provided, use the index as the key.
			if v >= len(planVersions) {
				alreadySetVersions[fmt.Sprintf("%s_v%d", strings.TrimSuffix(version.CreatedDate.Format(time.RFC3339), "Z"), version.Version)] = versionsModel
			} else {
				alreadySetVersions[planVersions[v]] = versionsModel
			}
		}
	}

	versionsValue, diags := types.MapValueFrom(ctx, versionsObjectType, alreadySetVersions)
	if diags.HasError() {
		return diags
	}
	m.Versions = versionsValue

	return nil
}

func (m *clientCertificateThirdPartyResourceModel) setClientCertificateData(ctx context.Context, clientCertificate *mtlskeystore.GetClientCertificateResponse) diag.Diagnostics {
	m.CertificateName = types.StringValue(clientCertificate.CertificateName)
	m.CertificateID = types.Int64Value(clientCertificate.CertificateID)
	m.Geography = types.StringValue(clientCertificate.Geography)
	m.KeyAlgorithm = types.StringValue(clientCertificate.KeyAlgorithm)
	m.SecureNetwork = types.StringValue(clientCertificate.SecureNetwork)
	m.Subject = types.StringValue(clientCertificate.Subject)

	notificationEmailsObject, diags := types.ListValueFrom(ctx, types.StringType, clientCertificate.NotificationEmails)
	if diags.HasError() {
		return diags
	}
	m.NotificationEmails = notificationEmailsObject

	return nil
}

func createClientVersionModel(ctx context.Context, version mtlskeystore.ClientCertificateVersion) (clientCertificateVersionModel, diag.Diagnostics) {
	versionsModel := clientCertificateVersionModel{
		Version:                  types.Int64Value(version.Version),
		Status:                   types.StringValue(version.Status),
		ExpiryDate:               types.StringPointerValue(formatOptionalRFC3339(version.ExpiryDate)),
		Issuer:                   types.StringPointerValue(version.Issuer),
		KeyAlgorithm:             types.StringValue(version.KeyAlgorithm),
		CertificateSubmittedBy:   types.StringPointerValue(version.CertificateSubmittedBy),
		CertificateSubmittedDate: types.StringPointerValue(formatOptionalRFC3339(version.CertificateSubmittedDate)),
		CreatedBy:                types.StringValue(version.CreatedBy),
		CreatedDate:              types.StringValue(version.CreatedDate.Format(time.RFC3339)),
		DeleteRequestedDate:      types.StringPointerValue(formatOptionalRFC3339(version.DeleteRequestedDate)),
		IssuedDate:               types.StringPointerValue(formatOptionalRFC3339(version.IssuedDate)),
		EllipticCurve:            types.StringPointerValue(version.EllipticCurve),
		KeySizeInBytes:           types.StringPointerValue(version.KeySizeInBytes),
		ScheduledDeleteDate:      types.StringPointerValue(formatOptionalRFC3339(version.ScheduledDeleteDate)),
		SignatureAlgorithm:       types.StringPointerValue(version.SignatureAlgorithm),
		Subject:                  types.StringPointerValue(version.Subject),
		VersionGUID:              types.StringValue(version.VersionGUID),
	}

	if version.CertificateBlock != nil {
		certificate := &certificateBlockResourceModel{
			Certificate: types.StringPointerValue(&version.CertificateBlock.Certificate),
			TrustChain:  types.StringPointerValue(&version.CertificateBlock.TrustChain),
		}
		certificateObject, diags := types.ObjectValueFrom(ctx, certificateBlockType(), certificate)
		if diags.HasError() {
			return clientCertificateVersionModel{}, diags
		}
		versionsModel.CertificateBlock = certificateObject
	} else {
		certificate := &certificateBlockResourceModel{
			Certificate: types.StringNull(),
			TrustChain:  types.StringNull(),
		}
		certificateObject, diags := types.ObjectValueFrom(ctx, certificateBlockType(), certificate)
		if diags.HasError() {
			return clientCertificateVersionModel{}, diags
		}
		versionsModel.CertificateBlock = certificateObject
	}

	if version.CSRBlock != nil {
		csr := &csrBlockResourceModel{
			CSR:          types.StringPointerValue(&version.CSRBlock.CSR),
			KeyAlgorithm: types.StringPointerValue(&version.CSRBlock.KeyAlgorithm),
		}
		csrObject, diags := types.ObjectValueFrom(ctx, csrBlockType(), csr)
		if diags.HasError() {
			return clientCertificateVersionModel{}, diags
		}
		versionsModel.CSRBlock = csrObject
	}

	return versionsModel, nil
}

// extractVersionKeys extracts the keys from the versions map.
func extractVersionKeys(versions map[string]clientCertificateVersionModel) []string {
	keys := make([]string, 0, len(versions))
	for key := range versions {
		keys = append(keys, key)
	}
	return keys
}

// mapKeyToVersionGUID maps the version GUIDs to their corresponding keys in the versions map.
func mapKeyToVersionGUID(versions map[string]clientCertificateVersionModel) map[string]string {
	var stateVersions = make(map[string]string)

	for key, version := range versions {
		stateVersions[version.VersionGUID.ValueString()] = key
	}

	return stateVersions
}

// filterVersionsToRemoveAndAdd filters the versions to remove, add and leave based on the state and plan keys.
// It returns three slices:
// 1. versionsToRemove: the versions that need to be removed from the state.
// 2. versionsToAdd: the versions that need to be added to the state.
// 3. remainingVersions: a map of version GUIDs to their corresponding keys in the state.
func filterVersionsToRemoveAndAdd(stateKeys, planKeys []string, versions map[string]clientCertificateVersionModel) ([]int64, []string, map[string]string) {
	versionsToRemove := make([]int64, 0)
	remainingVersions := make(map[string]string)
	versionsToAdd := slices.Clone(planKeys)

	for _, stateKey := range stateKeys {
		if slices.Contains(planKeys, stateKey) {
			versionsToAdd = slices.Delete(versionsToAdd, slices.Index(versionsToAdd, stateKey), slices.Index(versionsToAdd, stateKey)+1)
			remainingVersions[versions[stateKey].VersionGUID.ValueString()] = stateKey
		} else {
			versionsToRemove = append(versionsToRemove, versions[stateKey].Version.ValueInt64())
		}
	}

	return versionsToRemove, versionsToAdd, remainingVersions
}

// checkStatus checks the status of the client certificate versions to ensure they can be deleted.
func checkStatus(ctx context.Context, client mtlskeystore.MTLSKeystore, certificateID int64, versionsToRemove []int64) error {
	clientCertificateVersions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               certificateID,
		IncludeAssociatedProperties: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get client certificate versions: %w", err)
	}

	for _, version := range clientCertificateVersions.Versions {
		if versionsToRemove != nil {
			if !slices.Contains(versionsToRemove, version.Version) {
				continue
			}
		}
		if len(version.AssociatedProperties) > 0 {
			return fmt.Errorf("cannot delete client certificate version %d with associated properties", version.Version)
		}
	}
	return nil
}

// extractContractAndGroup extracts the contract and group from the subject string.
func extractContractAndGroup(subject string, parts []string) (string, string, error) {
	ctr := parts[len(parts)-2]
	grp := parts[len(parts)-1]
	// If groupID cannot be parsed as an integer, return an error.
	_, err := strconv.ParseInt(grp, 10, 64) // Ensure grp is a valid integer
	if err != nil {
		return "", "", fmt.Errorf("unable to extract group and contract from subject: '%s'", subject)
	}
	for _, label := range subjectLabels() {
		if strings.Contains(ctr, label) || strings.Contains(grp, label) {
			return "", "", fmt.Errorf("unable to extract group and contract from subject: '%s'", subject)
		}
	}

	return ctr, grp, nil
}

func subjectLabels() []string {
	return []string{
		"CN=", "O=", "OU=", "L=", "ST=", "C=",
	}
}

func (m *clientCertificateThirdPartyResourceModel) assignGroupAndContractThirdParty(ctrID, grpID string) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ContractID = types.StringValue(ctrID)

	gr, err := strconv.ParseInt(grpID, 10, 64)
	if err != nil {
		diags.AddError("Unable to parse group ID", err.Error())
		return diags
	}

	m.GroupID = types.Int64Value(gr)
	return diags
}
