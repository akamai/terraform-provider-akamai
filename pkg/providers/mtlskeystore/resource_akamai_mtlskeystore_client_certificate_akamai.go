package mtlskeystore

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const missedContractAndGroupError = "you need to provide an importID in the format 'certificateID,groupID,contractID'. Where certificate, groupID and contractID are required"

var (
	_ resource.Resource                = &clientCertificateAkamaiResource{}
	_ resource.ResourceWithConfigure   = &clientCertificateAkamaiResource{}
	_ resource.ResourceWithModifyPlan  = &clientCertificateAkamaiResource{}
	_ resource.ResourceWithImportState = &clientCertificateAkamaiResource{}
)

type clientCertificateAkamaiResource struct {
	meta meta.Meta
}

type clientCertificateAkamaiResourceModel struct {
	CertificateName    types.String `tfsdk:"certificate_name"`
	CertificateID      types.Int64  `tfsdk:"certificate_id"`
	ContractID         types.String `tfsdk:"contract_id"`
	Geography          types.String `tfsdk:"geography"`
	GroupID            types.Int64  `tfsdk:"group_id"`
	KeyAlgorithm       types.String `tfsdk:"key_algorithm"`
	NotificationEmails types.List   `tfsdk:"notification_emails"`
	SecureNetwork      types.String `tfsdk:"secure_network"`
	Subject            types.String `tfsdk:"subject"`
	CreatedBy          types.String `tfsdk:"created_by"`
	CreatedDate        types.String `tfsdk:"created_date"`
	Versions           types.List   `tfsdk:"versions"`
	CurrentGUID        types.String `tfsdk:"current_guid"`
	PreviousGUID       types.String `tfsdk:"previous_guid"`
}

type clientCertificateAkamaiVersionModel struct {
	Version             types.Int64  `tfsdk:"version"`
	Status              types.String `tfsdk:"status"`
	ExpiryDate          types.String `tfsdk:"expiry_date"`
	Issuer              types.String `tfsdk:"issuer"`
	KeyAlgorithm        types.String `tfsdk:"key_algorithm"`
	CreatedBy           types.String `tfsdk:"created_by"`
	CreatedDate         types.String `tfsdk:"created_date"`
	DeleteRequestedDate types.String `tfsdk:"delete_requested_date"`
	IssuedDate          types.String `tfsdk:"issued_date"`
	EllipticCurve       types.String `tfsdk:"elliptic_curve"`
	KeySizeInBytes      types.String `tfsdk:"key_size_in_bytes"`
	ScheduledDeleteDate types.String `tfsdk:"scheduled_delete_date"`
	SignatureAlgorithm  types.String `tfsdk:"signature_algorithm"`
	Subject             types.String `tfsdk:"subject"`
	VersionGUID         types.String `tfsdk:"version_guid"`
	CertificateBlock    types.Object `tfsdk:"certificate_block"`
}

type certificateAkamaiResourceBlockModel struct {
	Certificate types.String `tfsdk:"certificate"`
	TrustChain  types.String `tfsdk:"trust_chain"`
}

var (
	versionObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"version":               types.Int64Type,
			"status":                types.StringType,
			"expiry_date":           types.StringType,
			"issuer":                types.StringType,
			"key_algorithm":         types.StringType,
			"created_by":            types.StringType,
			"created_date":          types.StringType,
			"delete_requested_date": types.StringType,
			"issued_date":           types.StringType,
			"elliptic_curve":        types.StringType,
			"key_size_in_bytes":     types.StringType,
			"scheduled_delete_date": types.StringType,
			"signature_algorithm":   types.StringType,
			"subject":               types.StringType,
			"version_guid":          types.StringType,
			"certificate_block": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"certificate": types.StringType,
					"trust_chain": types.StringType,
				},
			},
		},
	}

	pollingDuration = 30 * time.Second
)

// NewClientCertificateAkamaiResource returns a new mtls keystore client certificate akamai resource.
func NewClientCertificateAkamaiResource() resource.Resource {
	return &clientCertificateAkamaiResource{}
}

func (c *clientCertificateAkamaiResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_client_certificate_akamai"
}

// Configure implements resource.ResourceWithConfigure.
func (c *clientCertificateAkamaiResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (c *clientCertificateAkamaiResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"certificate_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the client certificate. Must be between 1 and 64 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
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
					stringvalidator.OneOf(string(mtlskeystore.GeographyChinaAndCore),
						string(mtlskeystore.GeographyRussiaAndCore), string(mtlskeystore.GeographyCore)),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "The group assigned to the client certificate. Must be greater than or equal to 0.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
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
					stringvalidator.OneOf(string(mtlskeystore.SecureNetworkStandardTLS), string(mtlskeystore.SecureNetworkEnhancedTLS)),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"key_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The cryptographic algorithm used for key generation. Possible values: `RSA` or `ECDSA`.",
				Validators: []validator.String{
					stringvalidator.OneOf(string(mtlskeystore.KeyAlgorithmRSA), string(mtlskeystore.KeyAlgorithmECDSA)),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.PreventStringUpdate(),
				},
			},
			"subject": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The CA certificate’s key value details. The `CN` attribute is required and included in the subject. When not specified, the subject is constructed in this format: `/C=US/O=Akamai Technologies, Inc./OU={vcd_id} {contract_id} {group_id}/CN={certificate_name}/`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.PreventStringUpdate(),
				},
			},
			"certificate_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier of the client certificate.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the client certificate. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "An ISO 8601 timestamp indicating the client certificate's creation. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"current_guid": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the `current` client certificate version.",
				// Once GUID is established for the current version by API, it should not change.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"previous_guid": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the `previous` client certificate version.",
				// Once GUID is established for the previous version by API, it should not change.
			},
			"versions": versionSchema(),
		},
	}
}

func versionSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "A list of client certificate versions. Each version represents a specific iteration of the client certificate.",
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"version": schema.Int64Attribute{
					Computed:    true,
					Description: "The unique identifier of the client certificate version.",
				},
				"status": schema.StringAttribute{
					Description: "The client certificate version status. Possible values: `DEPLOYMENT_PENDING`, `DEPLOYED`, or `DELETE_PENDING`.",
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
			},
		},
	}
}

func (c *clientCertificateAkamaiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificates Akamai Resource Create")
	var plan clientCertificateAkamaiResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	clientCertificateRequest, diags := createClientCertificateRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	client := Client(c.meta)
	certificate, err := client.CreateClientCertificate(ctx, *clientCertificateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Client Certificate", err.Error())
		return
	}
	if resp.Diagnostics.Append(plan.populateCertModelFromResponse(ctx, mtlskeystore.Certificate(*certificate))...); resp.Diagnostics.HasError() {
		return
	}

	version, err := plan.waitUntilVersionDeployed(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for client certificate version deployment", err.Error())
		return
	}
	if resp.Diagnostics.Append(plan.populateVersionModelFromResponse(ctx, []mtlskeystore.ClientCertificateVersion{*version})...); resp.Diagnostics.HasError() {
		return
	}
	plan.CurrentGUID = types.StringValue(version.VersionGUID)
	plan.PreviousGUID = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func createClientCertificateRequest(ctx context.Context, plan *clientCertificateAkamaiResourceModel) (*mtlskeystore.CreateClientCertificateRequest, diag.Diagnostics) {
	var notificationEmails []string
	diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
	if diags.HasError() {
		return nil, diags
	}
	req := mtlskeystore.CreateClientCertificateRequest{
		CertificateName:    plan.CertificateName.ValueString(),
		ContractID:         strings.TrimPrefix(plan.ContractID.ValueString(), "ctr_"),
		Geography:          mtlskeystore.Geography(plan.Geography.ValueString()),
		GroupID:            plan.GroupID.ValueInt64(),
		NotificationEmails: notificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(plan.SecureNetwork.ValueString()),
		Signer:             mtlskeystore.SignerAkamai,
		Subject:            plan.Subject.ValueStringPointer(),
	}

	if plan.KeyAlgorithm.ValueString() != "" {
		req.KeyAlgorithm = (*mtlskeystore.CryptographicAlgorithm)(plan.KeyAlgorithm.ValueStringPointer())
	}
	return &req, nil
}

func (c *clientCertificateAkamaiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificates Akamai Resource Read")
	var state clientCertificateAkamaiResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(c.meta)
	certificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: state.CertificateID.ValueInt64(),
	})
	if err != nil {
		if errors.Is(err, mtlskeystore.ErrClientCertificateNotFound) {
			tflog.Debug(ctx, "Client Certificate Akamai Resource not found, removing from state")
			resp.Diagnostics.AddWarning("Resource Removal", "The client certificate was not found on the server. The resource will be removed from the state.")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to Get Client Certificate", err.Error())
		return
	}

	if resp.Diagnostics.Append(state.populateCertModelFromResponse(ctx, mtlskeystore.Certificate(*certificate))...); resp.Diagnostics.HasError() {
		return
	}

	versions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: state.CertificateID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Get Client Certificate Versions", err.Error())
		return
	}

	if len(versions.Versions) == 1 && versions.Versions[0].Status == string(mtlskeystore.DeletePending) {
		tflog.Debug(ctx, "Client Certificate Akamai Resource's last version is in pending delete status, removing from state")
		resp.Diagnostics.AddWarning("Resource Removal", "The last version of the Client Certificate is in `DELETE_PENDING` status. The resource will be removed from the state.")
		resp.State.RemoveResource(ctx)
		return
	}
	state.populateVersionModelFromResponse(ctx, versions.Versions)
	state.CurrentGUID = types.StringValue(versions.Versions[0].VersionGUID)
	if len(versions.Versions) == 2 {
		state.PreviousGUID = types.StringValue(versions.Versions[1].VersionGUID)
	} else {
		state.PreviousGUID = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *clientCertificateAkamaiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Akamai Resource Update")
	var plan, state clientCertificateAkamaiResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	isCertNameChanged := !plan.CertificateName.Equal(state.CertificateName)
	areNotificationEmailsChanged := !plan.NotificationEmails.Equal(state.NotificationEmails)

	if isCertNameChanged || areNotificationEmailsChanged {
		updateRequest := mtlskeystore.PatchClientCertificateRequest{
			CertificateID: state.CertificateID.ValueInt64(),
		}
		if isCertNameChanged {
			tflog.Debug(ctx, "Updating Client Certificate Name")
			updateRequest.Body.CertificateName = plan.CertificateName.ValueStringPointer()
		}
		var planEmails []string
		if areNotificationEmailsChanged {
			tflog.Debug(ctx, "Updating Client Certificate Notification Emails")
			diags := plan.NotificationEmails.ElementsAs(ctx, &planEmails, false)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
			updateRequest.Body.NotificationEmails = planEmails
		}

		client := Client(c.meta)
		err := client.PatchClientCertificate(ctx, updateRequest)
		if err != nil {
			resp.Diagnostics.AddError("Unable to Patch Client Certificate", err.Error())
			return
		}

		certResponse, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
			CertificateID: state.CertificateID.ValueInt64(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Unable to Fetch Client Certificate", err.Error())
			return
		}

		if resp.Diagnostics.Append(plan.populateCertModelFromResponse(ctx, mtlskeystore.Certificate(*certResponse))...); resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "Client Certificate Name or Notification Emails updated successfully")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *clientCertificateAkamaiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Akamai Resource Import")

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

	data := clientCertificateAkamaiResourceModel{
		CertificateID:      types.Int64Value(certificateID),
		CertificateName:    types.StringUnknown(),
		Geography:          types.StringUnknown(),
		NotificationEmails: types.ListUnknown(types.StringType),
		SecureNetwork:      types.StringUnknown(),
		Versions:           types.ListUnknown(versionObjectType),
	}

	// API call is needed to populate subject from server, and extract contract and group ID from it
	client := Client(c.meta)
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

	if diags := data.assignGroupAndContract(contractID, groupID); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if onlyOneVersionPendingDelete(ctx, client, certificateID, false, mtlskeystore.SignerAkamai) {
		resp.Diagnostics.AddError("Certificate in Delete Pending State", fmt.Sprintf("The client certificate %d has only one version and it's in `DELETE_PENDING` state. In order to import this resource, rotate this client certificate first", certificateID))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func onlyOneVersionPendingDelete(ctx context.Context, client mtlskeystore.MTLSKeystore, clientID int64, withAssociatedProperties bool, signer mtlskeystore.Signer) bool {
	versions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{CertificateID: clientID, IncludeAssociatedProperties: withAssociatedProperties})
	if err != nil {
		return false
	}
	numberOfActualVersions := 0
	for _, version := range versions.Versions {
		// In the third party case, sometimes versions are duplicated and those duplicates have alias. We don't count them.
		if signer == mtlskeystore.SignerThirdParty && version.VersionAlias != nil {
			continue
		}
		numberOfActualVersions++
	}

	return numberOfActualVersions == 1 && versions.Versions[0].Status == string(mtlskeystore.DeletePending)
}

func (c *clientCertificateAkamaiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Akamai Resource Delete")
	var state clientCertificateAkamaiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := Client(c.meta)

	versions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: state.CertificateID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Get Client Certificate Versions", err.Error())
		return
	}
	for _, version := range versions.Versions {
		if version.Status == string(mtlskeystore.DeletePending) {
			tflog.Debug(ctx, fmt.Sprintf("Client Certificate Version %d is already in delete pending state, skipping deletion", version.Version))
			continue
		}

		err := client.DeleteClientCertificateVersion(ctx, mtlskeystore.DeleteClientCertificateVersionRequest{
			CertificateID: state.CertificateID.ValueInt64(),
			Version:       version.Version,
		})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete Client Certificate Version %d", version.Version), err.Error())
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Deleted Client Certificate Version %d", version.Version))
	}
}

func (c *clientCertificateAkamaiResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var state, plan clientCertificateAkamaiResourceModel
	if modifiers.IsUpdate(req) {
		if !state.Subject.Equal(plan.Subject) {
			resp.Diagnostics.AddAttributeError(path.Root("subject"), "Cannot Update 'subject'",
				"The `subject` attribute cannot be updated after the resource has been created.")
			return
		}
	}
	// Update
	if modifiers.IsUpdate(req) {
		var state, plan clientCertificateAkamaiResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// PreviousGUID is null in state after create, and stays null until the first automatic rotation.
		// Therefore we need to copy the state value also for null.
		if plan.PreviousGUID.IsUnknown() {
			plan.PreviousGUID = state.PreviousGUID
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

func (m *clientCertificateAkamaiResourceModel) populateCertModelFromResponse(ctx context.Context, cert mtlskeystore.Certificate) diag.Diagnostics {
	m.CertificateID = types.Int64Value(cert.CertificateID)
	m.CertificateName = types.StringValue(cert.CertificateName)
	m.Geography = types.StringValue(cert.Geography)
	m.KeyAlgorithm = types.StringValue(cert.KeyAlgorithm)
	emails, diags := types.ListValueFrom(ctx, types.StringType, cert.NotificationEmails)
	if diags.HasError() {
		return diags
	}
	m.NotificationEmails = emails
	m.SecureNetwork = types.StringValue(cert.SecureNetwork)
	m.Subject = types.StringValue(cert.Subject)
	m.CreatedBy = types.StringValue(cert.CreatedBy)
	m.CreatedDate = date.TimeRFC3339Value(cert.CreatedDate)
	return nil
}

func (m *clientCertificateAkamaiResourceModel) populateVersionModelFromResponse(
	ctx context.Context, versions []mtlskeystore.ClientCertificateVersion) diag.Diagnostics {

	var diagnostics diag.Diagnostics
	var versionsModel []clientCertificateAkamaiVersionModel
	for _, version := range versions {
		var certificateBlock certificateAkamaiResourceBlockModel
		certificateBlock.Certificate = types.StringValue(version.CertificateBlock.Certificate)
		certificateBlock.TrustChain = types.StringValue(version.CertificateBlock.TrustChain)

		certificateObject, diags := types.ObjectValueFrom(ctx, certificateBlockSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes(), certificateBlock)
		diagnostics.Append(diags...)
		if diags.HasError() {
			return diagnostics
		}
		versionsModel = append(versionsModel, clientCertificateAkamaiVersionModel{
			Version:             types.Int64Value(version.Version),
			Status:              types.StringValue(version.Status),
			ExpiryDate:          date.TimeRFC3339PointerValue(version.ExpiryDate),
			Issuer:              types.StringPointerValue(version.Issuer),
			KeyAlgorithm:        types.StringValue(version.KeyAlgorithm),
			CreatedBy:           types.StringValue(version.CreatedBy),
			CreatedDate:         date.TimeRFC3339Value(version.CreatedDate),
			DeleteRequestedDate: date.TimeRFC3339PointerValue(version.DeleteRequestedDate),
			IssuedDate:          date.TimeRFC3339PointerValue(version.IssuedDate),
			EllipticCurve:       types.StringPointerValue(version.EllipticCurve),
			KeySizeInBytes:      types.StringPointerValue(version.KeySizeInBytes),
			ScheduledDeleteDate: date.TimeRFC3339PointerValue(version.ScheduledDeleteDate),
			SignatureAlgorithm:  types.StringPointerValue(version.SignatureAlgorithm),
			Subject:             types.StringPointerValue(version.Subject),
			VersionGUID:         types.StringValue(version.VersionGUID),
			CertificateBlock:    certificateObject,
		})
	}

	versionList, diags := types.ListValueFrom(ctx, versionSchema().NestedObject.Type().(types.ObjectType), versionsModel)
	diagnostics.Append(diags...)
	if diags.HasError() {
		return diagnostics
	}
	m.Versions = versionList
	return nil
}

func (m *clientCertificateAkamaiResourceModel) waitUntilVersionDeployed(
	ctx context.Context, client mtlskeystore.MTLSKeystore) (*mtlskeystore.ClientCertificateVersion, error) {

	for {
		versions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
			CertificateID: m.CertificateID.ValueInt64(),
		})
		if err != nil {
			return nil,
				fmt.Errorf("error getting client certificate versions: %w", err)
		}
		if len(versions.Versions) == 1 {
			tflog.Debug(ctx, fmt.Sprintf("Client Certificate Version status: %s", versions.Versions[0].Status))
			if versions.Versions[0].Status == string(mtlskeystore.Deployed) {
				return &versions.Versions[0], nil
			}
		} else {
			return nil,
				fmt.Errorf("unexpected number of client certificate versions: %d", len(versions.Versions))
		}

		select {
		case <-time.After(pollingDuration):
			continue
		case <-ctx.Done():
			return nil,
				fmt.Errorf("context cancelled while waiting for client certificate version deployment: %w", ctx.Err())
		}
	}
}

func parseCertificateID(idStr string) (int64, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse certificate ID as an integer: %v", idStr)
	}
	return int64(id), nil
}

func subjectContainsContractAndGroup(subject string) (bool, []string) {
	// Capture the part before required '/CN=' label.
	re := regexp.MustCompile(`\/([^\/]+)\/CN=`)
	matches := re.FindStringSubmatch(subject)
	if len(matches) < 2 {
		return false, nil
	}
	parts := strings.Fields(matches[1])
	return len(parts) >= 2, parts
}

func (m *clientCertificateAkamaiResourceModel) assignGroupAndContract(ctrID, grpID string) diag.Diagnostics {
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
