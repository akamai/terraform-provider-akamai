package mtlskeystore

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &clientCertificateThirdPartyResource{}
	_ resource.ResourceWithImportState = &clientCertificateThirdPartyResource{}
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
		PreferredCA        types.String `tfsdk:"preferred_ca"`
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
		DeployedDate             types.String `tfsdk:"deployed_date"`
		IssuedDate               types.String `tfsdk:"issued_date"`
		KeyEllipticCurve         types.String `tfsdk:"key_elliptic_curve"`
		KeySizeInBytes           types.String `tfsdk:"key_size_in_bytes"`
		ScheduledDeleteDate      types.String `tfsdk:"scheduled_delete_date"`
		SignatureAlgorithm       types.String `tfsdk:"signature_algorithm"`
		Subject                  types.String `tfsdk:"subject"`
		VersionGUID              types.String `tfsdk:"version_guid"`
		CertificateBlock         types.Object `tfsdk:"certificate_block"`
		CSRBlock                 types.Object `tfsdk:"csr_block"`
		Validation               types.Object `tfsdk:"validation"`
		AssociatedProperties     types.List   `tfsdk:"associated_properties"`
	}

	certificateBlockResourceModel struct {
		Certificate types.String `tfsdk:"certificate"`
		TrustChain  types.String `tfsdk:"trust_chain"`
	}

	csrBlockResourceModel struct {
		CSR          types.String `tfsdk:"csr"`
		KeyAlgorithm types.String `tfsdk:"key_algorithm"`
	}

	associatedPropertiesModel struct {
		AssetID         types.Int64  `tfsdk:"asset_id"`
		GroupID         types.Int64  `tfsdk:"group_id"`
		PropertyName    types.String `tfsdk:"property_name"`
		PropertyVersion types.Int64  `tfsdk:"property_version"`
	}
)

var (
	cnRegex            = regexp.MustCompile(`/CN=[^/]{1,64}/`)
	versionsObjectType = types.ObjectType{
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
			"deployed_date":              types.StringType,
			"issued_date":                types.StringType,
			"key_elliptic_curve":         types.StringType,
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
			"validation": types.ObjectType{
				AttrTypes: validationType(),
			},
			"associated_properties": types.ListType{
				ElemType: associatedPropertiesType(),
			},
		},
	}
)

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
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "The group assigned to the client certificate. Must be greater than or equal to 0.",
			},
			"key_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The cryptographic algorithm used for key generation. Possible values: `RSA` or `ECDSA`.",
				Default:     stringdefault.StaticString(string(mtlskeystore.KeyAlgorithmRSA)),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"RSA",
						"ECDSA",
					),
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
			"preferred_ca": schema.StringAttribute{
				Optional:    true,
				Description: "The common name of the account CA certificate selected to sign the client certificate. Specify `null` if you want to add this later.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 64),
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
			},
			"subject": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Specifies the client certificate. The `CN` attribute is required and cannot exceed 64 characters. When `null`, the subject is constructed with the following format: `/C=US/O=Akamai Technologies, Inc./OU={vcdId} {contractId} {groupId}/CN={certificateName}/`.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(cnRegex, "The `subject` must contain a valid `CN` attribute with a maximum length of 64 characters."),
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
				"deployed_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating the client certificate version's activation. Appears as null if not specified.",
					Computed:    true,
				},
				"issued_date": schema.StringAttribute{
					Description: "An ISO 8601 timestamp indicating the client certificate version's availability.",
					Computed:    true,
				},
				"key_elliptic_curve": schema.StringAttribute{
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
				"certificate_block":     certificateBlockSchema(),
				"csr_block":             csrBlockSchema(),
				"validation":            validationSchema(),
				"associated_properties": associatedPropertiesSchema(),
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
				Description: "Identifies the client certificate's encryption algorithm. The only currently supported value is `RSA`.",
				Computed:    true,
			},
		},
	}
}

func validationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Details of the validation for the client certificate version.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"errors": schema.ListNestedAttribute{
				Description: "A list of validation errors for the client certificate version.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"message": schema.StringAttribute{
							Description: "The validation error message.",
							Computed:    true,
						},
						"reason": schema.StringAttribute{
							Description: "The validation error reason.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The validation error type.",
							Computed:    true,
						},
					},
				},
			},
			"warnings": schema.ListNestedAttribute{
				Description: "A list of validation warnings for the client certificate version.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"message": schema.StringAttribute{
							Description: "The validation warning message.",
							Computed:    true,
						},
						"reason": schema.StringAttribute{
							Description: "The validation warning reason.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The validation warning type.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func associatedPropertiesSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "A list of properties associated with the client certificate version.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"asset_id": schema.Int64Attribute{
					Description: "The unique identifier of the property associated with the client certificate version.",
					Computed:    true,
				},
				"group_id": schema.Int64Attribute{
					Description: "The group ID of the property associated with the client certificate version.",
					Computed:    true,
				},
				"property_name": schema.StringAttribute{
					Description: "The name of the property associated with the client certificate version.",
					Computed:    true,
				},
				"property_version": schema.Int64Attribute{
					Description: "The version of the property associated with the client certificate version.",
					Computed:    true,
				},
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

func (r *clientCertificateThirdPartyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Client Certificate Third Party Resource")
	var plan clientCertificateThirdPartyResourceModel

	// Read the plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.create(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Creating Client Certificate Third Party failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clientCertificateThirdPartyResource) create(ctx context.Context, plan *clientCertificateThirdPartyResourceModel) error {
	client := Client(r.meta)

	var notificationEmails []string
	if isKnown(plan.NotificationEmails) {
		diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
		if diags.HasError() {
			return fmt.Errorf("failed to get notification emails: %v", diags)
		}
	}

	var versions map[string]clientCertificateVersionModel

	diags := plan.Versions.ElementsAs(ctx, &versions, false)
	if diags.HasError() {
		return fmt.Errorf("failed to get versions: %v", diags)
	}
	versionsKeys := extractVersionKeys(versions)
	slices.Sort(versionsKeys)

	request := mtlskeystore.CreateClientCertificateRequest{
		CertificateName:    plan.CertificateName.ValueString(),
		ContractID:         plan.ContractID.ValueString(),
		Geography:          mtlskeystore.Geography(plan.Geography.ValueString()),
		GroupID:            plan.GroupID.ValueInt64(),
		NotificationEmails: notificationEmails,
		PreferredCA:        plan.PreferredCA.ValueStringPointer(),
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

	err = createMissingVersionsExplicitly(ctx, client, createClientCertificateResponse.CertificateID, versionsKeys[1:])
	if err != nil {
		return fmt.Errorf("failed to create missing versions: %w", err)
	}

	clientCertificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: createClientCertificateResponse.CertificateID,
	})
	if err != nil {
		return fmt.Errorf("failed to get client certificate: %w", err)
	}

	diags = plan.setClientCertificateData(ctx, clientCertificate)
	if diags.HasError() {
		return fmt.Errorf("failed to set data: %v", diags)
	}

	clientCertificateVersions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               createClientCertificateResponse.CertificateID,
		IncludeAssociatedProperties: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get client certificate versions: %w", err)
	}

	diags = plan.setVersionsData(ctx, clientCertificateVersions, versionsKeys, nil)
	if diags.HasError() {
		return fmt.Errorf("failed to set versions data: %v", diags)
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

	if err := r.read(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Reading API Client Resource failed", err.Error())
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
	var dataVersions map[string]clientCertificateVersionModel

	data.Versions.ElementsAs(ctx, &dataVersions, false)
	stateVersions := mapKeyToVersionGUID(dataVersions)

	diags = data.setVersionsData(ctx, clientCertificateVersions, nil, stateVersions)
	if diags.HasError() {
		return fmt.Errorf("failed to set versions data: %v", diags)
	}

	return nil
}

func (r *clientCertificateThirdPartyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Client Certificate Third Party Resource")

	certificateID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("could not convert import ID to int", err.Error())
		return
	}

	data := &clientCertificateThirdPartyResourceModel{
		CertificateName:    types.StringUnknown(),
		CertificateID:      types.Int64Value(certificateID),
		Geography:          types.StringUnknown(),
		NotificationEmails: types.ListUnknown(types.StringType),
		SecureNetwork:      types.StringUnknown(),
		Versions:           types.MapUnknown(versionsObjectType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	if plan.CertificateName != oldState.CertificateName || !plan.NotificationEmails.Equal(oldState.NotificationEmails) {
		var notificationEmails []string
		if isKnown(plan.NotificationEmails) {
			diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
			if diags.HasError() {
				return fmt.Errorf("failed to get notification emails: %v", diags)
			}
		}
		err := client.PatchClientCertificate(ctx, mtlskeystore.PatchClientCertificateRequest{
			CertificateID: oldState.CertificateID.ValueInt64(),
			Body: mtlskeystore.PatchClientCertificateRequestBody{
				CertificateName:    plan.CertificateName.ValueStringPointer(),
				NotificationEmails: notificationEmails,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update client certificate: %w", err)
		}
	}

	var planVersions map[string]clientCertificateVersionModel
	var oldStateVersions map[string]clientCertificateVersionModel

	plan.Versions.ElementsAs(ctx, &planVersions, false)
	oldState.Versions.ElementsAs(ctx, &oldStateVersions, false)

	planKeys := extractVersionKeys(planVersions)
	stateKeys := extractVersionKeys(oldStateVersions)
	slices.Sort(planKeys)

	versionsToRemove, versionsToAdd, stateVersions := filterVersionsToRemoveAndAdd(stateKeys, planKeys, oldStateVersions)

	if len(versionsToRemove) > 0 {
		err := checkStatus(ctx, client, oldState.CertificateID.ValueInt64(), versionsToRemove)
		if err != nil {
			return err
		}
	}

	if len(versionsToRemove) != 0 || len(versionsToAdd) != 0 {
		if len(versionsToRemove) < len(stateKeys) {
			for _, version := range versionsToRemove {
				_, err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
					CertificateID: oldState.CertificateID.ValueInt64(),
					Version:       version,
				})
				if err != nil {
					return fmt.Errorf("failed to delete version %d: %w", version, err)
				}
			}
			for _, version := range versionsToAdd {
				newVersion, err := client.RotateClientCertificateVersion(ctx, mtlskeystore.RotateClientCertificateVersionRequest{
					CertificateID: oldState.CertificateID.ValueInt64(),
				})
				if err != nil {
					return fmt.Errorf("failed to create version %s: %w", version, err)
				}
				stateVersions[newVersion.VersionGUID] = version
			}
		} else {
			slices.Sort(versionsToRemove)
			for _, version := range versionsToRemove[1:] {
				_, err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
					CertificateID: oldState.CertificateID.ValueInt64(),
					Version:       version,
				})
				if err != nil {
					return fmt.Errorf("failed to delete version %d: %w", version, err)
				}
			}
			for _, version := range versionsToAdd {
				newVersion, err := client.RotateClientCertificateVersion(ctx, mtlskeystore.RotateClientCertificateVersionRequest{
					CertificateID: oldState.CertificateID.ValueInt64(),
				})
				if err != nil {
					return fmt.Errorf("failed to create version %s: %w", version, err)
				}
				stateVersions[newVersion.VersionGUID] = version
			}
			_, err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
				CertificateID: oldState.CertificateID.ValueInt64(),
				Version:       versionsToRemove[0],
			})
			if err != nil {
				return fmt.Errorf("failed to delete version %d: %w", versionsToRemove[0], err)
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

func (r *clientCertificateThirdPartyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Client Certificate Third Party Resource")

	var state clientCertificateThirdPartyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(r.meta)

	err := checkStatus(ctx, client, state.CertificateID.ValueInt64(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Deleting Client Certificate Versions failed", err.Error())
		return
	}

	var stateVersions map[string]clientCertificateVersionModel
	state.Versions.ElementsAs(ctx, &stateVersions, false)

	for version := range stateVersions {
		deleteMessage, err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
			CertificateID: state.CertificateID.ValueInt64(),
			Version:       stateVersions[version].Version.ValueInt64(),
		})
		if err != nil {
			message := ""
			if deleteMessage != nil {
				message += fmt.Sprintf("Delete message: %s", *deleteMessage)
			}
			resp.Diagnostics.AddError(
				"Deleting Client Certificate Version failed",
				fmt.Sprintf("Failed to delete version %s: %v. ", version, err)+message,
			)
			return
		}
	}
}

func isKnown(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func (m *clientCertificateThirdPartyResourceModel) setVersionsData(ctx context.Context, clientCertificateVersions *mtlskeystore.ListClientCertificateVersionsResponse, versionKeys []string, remainingVersions map[string]string) diag.Diagnostics {
	slices.Reverse(versionKeys)
	versions := make(map[string]clientCertificateVersionModel)
	for v, version := range clientCertificateVersions.Versions {
		if version.VersionAlias != nil {
			// Skip the version with an alias
			continue
		}
		versionsModel := clientCertificateVersionModel{
			Version:                  types.Int64Value(version.Version),
			Status:                   types.StringValue(string(version.Status)),
			ExpiryDate:               types.StringValue(version.ExpiryDate),
			Issuer:                   types.StringValue(version.Issuer),
			KeyAlgorithm:             types.StringValue(string(version.KeyAlgorithm)),
			CertificateSubmittedBy:   types.StringPointerValue(version.CertificateSubmittedBy),
			CertificateSubmittedDate: types.StringPointerValue(version.CertificateSubmittedDate),
			CreatedBy:                types.StringValue(version.CreatedBy),
			CreatedDate:              types.StringValue(version.CreatedDate),
			DeleteRequestedDate:      types.StringPointerValue(version.DeleteRequestedDate),
			DeployedDate:             types.StringPointerValue(version.DeployedDate),
			IssuedDate:               types.StringValue(version.IssuedDate),
			KeyEllipticCurve:         types.StringValue(version.KeyEllipticCurve),
			KeySizeInBytes:           types.StringValue(version.KeySizeInBytes),
			ScheduledDeleteDate:      types.StringPointerValue(version.ScheduledDeleteDate),
			SignatureAlgorithm:       types.StringValue(version.SignatureAlgorithm),
			Subject:                  types.StringValue(version.Subject),
			VersionGUID:              types.StringValue(version.VersionGUID),
		}

		if version.CertificateBlock != nil {
			certificate := &certificateBlockResourceModel{
				Certificate: types.StringPointerValue(&version.CertificateBlock.Certificate),
				TrustChain:  types.StringPointerValue(&version.CertificateBlock.TrustChain),
			}
			certificateObject, diags := types.ObjectValueFrom(ctx, certificateBlockType(), certificate)
			if diags.HasError() {
				return diags
			}
			versionsModel.CertificateBlock = certificateObject
		} else {
			certificate := &certificateBlockResourceModel{
				Certificate: types.StringNull(),
				TrustChain:  types.StringNull(),
			}
			certificateObject, diags := types.ObjectValueFrom(ctx, certificateBlockType(), certificate)
			if diags.HasError() {
				return diags
			}
			versionsModel.CertificateBlock = certificateObject
		}

		if version.CSRBlock != nil {
			csr := &csrBlockResourceModel{
				CSR:          types.StringPointerValue(&version.CSRBlock.CSR),
				KeyAlgorithm: types.StringPointerValue((*string)(&version.CSRBlock.KeyAlgorithm)),
			}
			csrObject, diags := types.ObjectValueFrom(ctx, csrBlockType(), csr)
			if diags.HasError() {
				return diags
			}
			versionsModel.CSRBlock = csrObject
		}

		var errors []validationErrorModel
		for _, err := range version.Validation.Errors {
			errors = append(errors, validationErrorModel{
				Message: err.Message,
				Reason:  err.Reason,
				Type:    err.Type,
			})
		}

		var warnings []validationErrorModel
		for _, warning := range version.Validation.Warnings {
			warnings = append(warnings, validationErrorModel{
				Message: warning.Message,
				Reason:  warning.Reason,
				Type:    warning.Type,
			})
		}

		validation := validationModel{
			Errors:   errors,
			Warnings: warnings,
		}

		validationObject, diags := types.ObjectValueFrom(ctx, validationType(), validation)
		if diags.HasError() {
			return diags
		}

		versionsModel.Validation = validationObject

		var associatedProperties []associatedPropertiesModel
		for _, property := range version.AssociatedProperties {
			associatedProperties = append(associatedProperties, associatedPropertiesModel{
				AssetID:         types.Int64Value(property.AssetID),
				GroupID:         types.Int64Value(property.GroupID),
				PropertyName:    types.StringValue(property.PropertyName),
				PropertyVersion: types.Int64Value(property.PropertyVersion),
			})
		}
		associatedPropertiesObject, diags := types.ListValueFrom(ctx, associatedPropertiesType(), associatedProperties)
		if diags.HasError() {
			return diags
		}

		versionsModel.AssociatedProperties = associatedPropertiesObject

		if versionKey, exists := remainingVersions[version.VersionGUID]; exists {
			versions[versionKey] = versionsModel
		} else {
			if v >= len(versionKeys) {
				versions[fmt.Sprintf("%s_v%d", strings.TrimSuffix(version.CreatedDate, "Z"), version.Version)] = versionsModel
			} else {
				versions[versionKeys[v]] = versionsModel
			}
		}
	}

	versionsValue, diags := types.MapValueFrom(ctx, versionsObjectType, versions)
	if diags.HasError() {
		return diags
	}
	m.Versions = versionsValue

	return nil
}

func (m *clientCertificateThirdPartyResourceModel) setClientCertificateData(ctx context.Context, clientCertificate *mtlskeystore.GetClientCertificateResponse) diag.Diagnostics {
	m.CertificateName = types.StringValue(clientCertificate.CertificateName)
	m.CertificateID = types.Int64Value(clientCertificate.CertificateID)
	m.Geography = types.StringValue(string(clientCertificate.Geography))
	m.KeyAlgorithm = types.StringValue(string(clientCertificate.KeyAlgorithm))

	notificationEmailsObject, diags := types.ListValueFrom(ctx, types.StringType, clientCertificate.NotificationEmails)
	if diags.HasError() {
		return diags
	}
	m.NotificationEmails = notificationEmailsObject

	m.SecureNetwork = types.StringValue(string(clientCertificate.SecureNetwork))
	m.Subject = types.StringValue(clientCertificate.Subject)

	return nil
}

func extractVersionKeys(versions map[string]clientCertificateVersionModel) []string {
	keys := make([]string, 0, len(versions))
	for key := range versions {
		keys = append(keys, key)
	}
	return keys
}

func mapKeyToVersionGUID(versions map[string]clientCertificateVersionModel) map[string]string {
	var stateVersions = make(map[string]string)

	for key, version := range versions {
		stateVersions[version.VersionGUID.ValueString()] = key
	}

	return stateVersions
}

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
		if version.Status == mtlskeystore.DeletePending {
			return fmt.Errorf("cannot delete client certificate version with status %s", version.Status)
		}
		if len(version.AssociatedProperties) > 0 {
			return fmt.Errorf("cannot delete client certificate version %d with associated properties", version.Version)
		}
	}
	return nil
}

func certificateBlockType() map[string]attr.Type {
	return certificateBlockSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func csrBlockType() map[string]attr.Type {
	return csrBlockSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func validationType() map[string]attr.Type {
	return validationSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func associatedPropertiesType() types.ObjectType {
	return associatedPropertiesSchema().NestedObject.Type().(types.ObjectType)
}
