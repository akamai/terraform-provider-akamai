package mtlskeystore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource               = &clientCertificateUploadResource{}
	_ resource.ResourceWithConfigure  = &clientCertificateUploadResource{}
	_ resource.ResourceWithModifyPlan = &clientCertificateUploadResource{}

	pollingInterval = 1 * time.Minute
	defaultTimeout  = 30 * time.Minute
)

// clientCertificateUploadResource represents akamai_mtlskeystore_client_certificate_upload resource.
type clientCertificateUploadResource struct {
	meta meta.Meta
}

// NewClientCertificateUploadResource creates a new instance of the Akamai MTLS Keystore Client Certificate Upload resource.
func NewClientCertificateUploadResource() resource.Resource {
	return &clientCertificateUploadResource{}
}

// Metadata returns the type name for the Akamai MTLS Keystore Client Certificate Upload resource.
func (r *clientCertificateUploadResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_client_certificate_upload"
}

// Schema defines the schema for the Akamai MTLS Keystore Client Certificate Upload resource.
func (r *clientCertificateUploadResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_certificate_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the client certificate to which the signed certificate will be uploaded.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"version_number": schema.Int64Attribute{
				Required:    true,
				Description: "The version number of the client certificate to upload the signed certificate to.",
			},
			"signed_certificate": schema.StringAttribute{
				Required:    true,
				Description: "The signed certificate to upload.",
			},
			"trust_chain": schema.StringAttribute{
				Optional:    true,
				Description: "The optional trust chain associated with the signed certificate.",
			},
			"wait_for_deployment": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Indicates whether to wait for the deployment of the uploaded certificate. Defaults to `true`.",
			},
			"version_guid": schema.StringAttribute{
				Computed:    true,
				Description: "A unique identifier for the client certificate version.",
			},
			"auto_acknowledge_warnings": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If set to true, all warnings will be acknowledged automatically. Defaults to `false`.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "Optional configurable resource create timeout. By default it's 30m.",
				Update:            true,
				UpdateDescription: "Optional configurable resource update timeout. By default it's 30m.",
			}),
		},
	}
}

// clientCertificateUploadModel is a model for akamai_mtlskeystore_client_certificate_upload resource.
type clientCertificateUploadModel struct {
	ClientCertificateID     types.Int64    `tfsdk:"client_certificate_id"`
	VersionNumber           types.Int64    `tfsdk:"version_number"`
	SignedCertificate       types.String   `tfsdk:"signed_certificate"`
	TrustChain              types.String   `tfsdk:"trust_chain"`
	WaitForDeployment       types.Bool     `tfsdk:"wait_for_deployment"`
	VersionGUID             types.String   `tfsdk:"version_guid"`
	AutoAcknowledgeWarnings types.Bool     `tfsdk:"auto_acknowledge_warnings"`
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
}

// Configure implements the resource.ResourceWithConfigure interface.
func (r *clientCertificateUploadResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	r.meta = meta.Must(req.ProviderData)
}

func (r *clientCertificateUploadResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Upload ModifyPlan")

	if modifiers.IsUpdate(request) {
		var plan, state clientCertificateUploadModel

		response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
		if response.Diagnostics.HasError() {
			return
		}

		response.Diagnostics.Append(request.State.Get(ctx, &state)...)
		if response.Diagnostics.HasError() {
			return
		}

		versionChanged := state.VersionNumber != plan.VersionNumber

		if !versionChanged && state.TrustChain != plan.TrustChain {
			response.Diagnostics.AddError("Only updates with a different version_number are supported", "The trust_chain attribute cannot be updated after the initial upload. Please create a new version with the updated trust chain.")
		}

		if !versionChanged && state.SignedCertificate != plan.SignedCertificate {
			response.Diagnostics.AddError("Only updates with a different version_number are supported", "The signed_certificate attribute cannot be updated after the initial upload. Please create a new version with the updated signed certificate.")
		}

		if !versionChanged && state.AutoAcknowledgeWarnings != plan.AutoAcknowledgeWarnings {
			response.Diagnostics.AddError("Only updates with a different version_number are supported", "The auto_acknowledge_warnings attribute cannot be updated after the initial upload. Please create a new version with the updated auto_acknowledge_warnings.")
		}

		if !versionChanged && state.WaitForDeployment != plan.WaitForDeployment {
			response.Diagnostics.AddError("Only updates with a different version_number are supported", "The wait_for_deployment attribute cannot be updated after the initial upload. Please create a new version with the updated wait_for_deployment.")
		}

		if !versionChanged && !state.Timeouts.Equal(plan.Timeouts) {
			response.Diagnostics.AddError("Only updates with a different version_number are supported", "The timeouts attribute cannot be updated after the initial upload. Please create a new version with the updated timeout.")
		}

		if versionChanged && state.SignedCertificate == plan.SignedCertificate {
			response.Diagnostics.AddError("No change in signed certificate", "Updating version_number requires a change in the signed_certificate attribute. Please provide a new signed certificate for the new version.")
		}

		if response.Diagnostics.HasError() {
			return
		}
	}
}

// Create implements resource's Create method.
func (r *clientCertificateUploadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Upload Create")

	var plan clientCertificateUploadModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(r.meta)

	err := plan.validateCertificateVersion(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Error validating client certificate version", err.Error())
		return
	}

	uploadedVersion, err := plan.upsertClientCertificateUpload(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Error uploading signed certificate: ", err.Error())
		return
	}
	plan.VersionGUID = types.StringValue(uploadedVersion.VersionGUID)

	// Wait for deployment if needed
	if plan.WaitForDeployment.ValueBool() && (uploadedVersion.Status != string(mtlskeystore.CertificateVersionStatusDeployed)) {
		timeout, diag := plan.Timeouts.Create(ctx, defaultTimeout)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		if err = plan.waitForDeployment(ctx, client, timeout); err != nil {
			resp.Diagnostics.AddError("Error polling for client certificate deployment", err.Error())
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read implements resource's Read method.
func (r *clientCertificateUploadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Upload Read")

	var state clientCertificateUploadModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(r.meta)
	listVersionsRequest := mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: state.ClientCertificateID.ValueInt64(),
	}
	clientCertificateVersionsResp, err := client.ListClientCertificateVersions(ctx, listVersionsRequest)
	if err != nil {
		if errors.Is(err, mtlskeystore.ErrClientCertificateNotFound) {
			tflog.Debug(ctx, "Certificate not found, removing resource from state")
			resp.Diagnostics.AddWarning("Client Certificate Not Found; removing resource from state", fmt.Sprintf("Client certificate %d not found", state.ClientCertificateID.ValueInt64()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving client certificate versions", err.Error())
		return
	}
	if clientCertificateVersionsResp == nil || len(clientCertificateVersionsResp.Versions) == 0 {
		// Resource drift: remove from state
		tflog.Debug(ctx, "Certificate versions not found, removing resource from state")
		resp.Diagnostics.AddWarning("No versions found for Client Certificate; removing resource from state", fmt.Sprintf("No versions found for client certificate %d", state.ClientCertificateID.ValueInt64()))
		resp.State.RemoveResource(ctx)
		return
	}
	for _, version := range clientCertificateVersionsResp.Versions {
		if isCorrectNonAliasedVersion(version, state.VersionNumber.ValueInt64()) {
			state.VersionGUID = types.StringValue(version.VersionGUID)
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
			return
		}
	}
	resp.Diagnostics.AddWarning("Client Certificate Version Not Found; removing resource from state", fmt.Sprintf("Uploaded version %d not found for client certificate %d", state.VersionNumber.ValueInt64(), state.ClientCertificateID.ValueInt64()))
	tflog.Debug(ctx, fmt.Sprintf("Could not find the uploaded version %d in the client certificate %d, removing resource from state", state.VersionNumber.ValueInt64(), state.ClientCertificateID.ValueInt64()))
	resp.State.RemoveResource(ctx)
}

// Delete implements resource's Delete method.
func (r *clientCertificateUploadResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No action needed on the server side, just remove from state
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Upload Delete")
	tflog.Debug(ctx, "Removing client certificate upload resource from state")
	resp.State.RemoveResource(ctx)
}

// Update implements resource's Update method.
func (r *clientCertificateUploadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Upload Update")

	var plan, state clientCertificateUploadModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(r.meta)

	err := plan.validateCertificateVersion(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Error validating client certificate version", err.Error())
		return
	}

	uploadedVersion, err := plan.upsertClientCertificateUpload(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Error updating client certificate upload", err.Error())
		return
	}
	plan.VersionGUID = types.StringValue(uploadedVersion.VersionGUID)

	// Wait for deployment if needed
	if plan.WaitForDeployment.ValueBool() && (uploadedVersion.Status != string(mtlskeystore.CertificateVersionStatusDeployed)) {
		timeout, diag := plan.Timeouts.Update(ctx, defaultTimeout)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		if err = plan.waitForDeployment(ctx, client, timeout); err != nil {
			resp.Diagnostics.AddError("Error waiting for client certificate deployment", err.Error())
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (data *clientCertificateUploadModel) upsertClientCertificateUpload(ctx context.Context, client mtlskeystore.MTLSKeystore) (*mtlskeystore.ClientCertificateVersion, error) {
	request := mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: data.ClientCertificateID.ValueInt64(),
		Version:       data.VersionNumber.ValueInt64(),
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: data.SignedCertificate.ValueString(),
			TrustChain:  data.TrustChain.ValueStringPointer(),
		},
		AcknowledgeAllWarnings: data.AutoAcknowledgeWarnings.ValueBoolPointer(),
	}

	if err := client.UploadSignedClientCertificate(ctx, request); err != nil {
		return nil, fmt.Errorf("error uploading client certificate: %w", err)
	}

	listVersionsRequest := mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: data.ClientCertificateID.ValueInt64(),
	}
	clientCertificateVersionsResp, err := client.ListClientCertificateVersions(ctx, listVersionsRequest)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client certificate versions: %w", err)
	}
	if clientCertificateVersionsResp == nil || len(clientCertificateVersionsResp.Versions) == 0 {
		return nil, fmt.Errorf("no client certificate versions found: received nil or empty response from client")
	}
	var uploadedVersion *mtlskeystore.ClientCertificateVersion
	for _, version := range clientCertificateVersionsResp.Versions {
		if isCorrectNonAliasedVersion(version, data.VersionNumber.ValueInt64()) {
			uploadedVersion = &version
			break
		}
	}
	if uploadedVersion == nil || uploadedVersion.VersionGUID == "" {
		return nil, fmt.Errorf("uploaded version not found: could not find the uploaded version in the response from the client")
	}

	if uploadedVersion.Status != string(mtlskeystore.CertificateVersionStatusAwaitingSigned) &&
		uploadedVersion.Status != string(mtlskeystore.CertificateVersionStatusDeploymentPending) && uploadedVersion.Status != string(mtlskeystore.CertificateVersionStatusDeployed) {
		return nil, fmt.Errorf(
			"unexpected client certificate version status: expected status to be either 'AWAITING_SIGNED_CERTIFICATE', 'DEPLOYMENT_PENDING' or 'DEPLOYED', but got: %s",
			uploadedVersion.Status,
		)
	}
	return uploadedVersion, nil
}

func (data *clientCertificateUploadModel) validateCertificateVersion(ctx context.Context, client mtlskeystore.MTLSKeystore) error {
	versions, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: data.ClientCertificateID.ValueInt64(),
	})
	if err != nil {
		return fmt.Errorf("could not retrieve client certificate versions: %w", err)
	}
	for _, version := range versions.Versions {
		if isCorrectNonAliasedVersion(version, data.VersionNumber.ValueInt64()) {
			return nil
		}
	}
	return fmt.Errorf("could not find client certificate version %d for certificate ID %d", data.VersionNumber.ValueInt64(), data.ClientCertificateID.ValueInt64())
}

func (data *clientCertificateUploadModel) waitForDeployment(ctx context.Context, client mtlskeystore.MTLSKeystore, timeout time.Duration) error {
	tflog.Debug(ctx, "Waiting for client certificate deployment")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-time.After(pollingInterval):
			clientCertificateVersionsResp, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
				CertificateID: data.ClientCertificateID.ValueInt64(),
			})
			if err != nil {
				return fmt.Errorf("error retrieving client certificate versions: %w", err)
			}
			if clientCertificateVersionsResp == nil || len(clientCertificateVersionsResp.Versions) == 0 {
				return fmt.Errorf("no client certificate versions found")
			}
			for _, version := range clientCertificateVersionsResp.Versions {
				if isCorrectNonAliasedVersion(version, data.VersionNumber.ValueInt64()) {
					if version.Status == string(mtlskeystore.CertificateVersionStatusDeployed) {
						tflog.Debug(ctx, fmt.Sprintf("Client certificate %d version %d is deployed", data.ClientCertificateID.ValueInt64(), data.VersionNumber.ValueInt64()))
						return nil
					}
					tflog.Debug(ctx, fmt.Sprintf("Client certificate %d version %d is in %s status", data.ClientCertificateID.ValueInt64(), data.VersionNumber.ValueInt64(), version.Status))
				}
			}
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for client certificate deployment: exceeded %s retries", timeout.String())
		}
	}
}

func isCorrectNonAliasedVersion(version mtlskeystore.ClientCertificateVersion, expectedVersion int64) bool {
	return version.Version == expectedVersion && version.VersionAlias == nil
}
