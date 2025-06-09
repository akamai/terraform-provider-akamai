package mtlskeystore

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &clientCertificateUploadResource{}
	_ resource.ResourceWithConfigure = &clientCertificateUploadResource{}

	pollingInterval = 1 * time.Minute
	numberOfRetries = 20
)

// clientCertificateUploadResource represents akamai_mtlskeystore_client_certificate_upload resource.
type clientCertificateUploadResource struct {
	meta meta.Meta
}

// NewAkamaiMTLSKeystoreClientCertificateUploadResource creates a new instance of the Akamai MTLS Keystore Client Certificate Upload resource.
func NewAkamaiMTLSKeystoreClientCertificateUploadResource() resource.Resource {
	return &clientCertificateUploadResource{}
}

// Metadata returns the type name for the Akamai MTLS Keystore Client Certificate Upload resource.
func (r *clientCertificateUploadResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_client_certificate_upload"
}

// Schema defines the schema for the Akamai MTLS Keystore Client Certificate Upload resource.
func (r *clientCertificateUploadResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "Indicates whether to wait for the deployment of the uploaded certificate. Defaults to `true`.",
			},
			"version_guid": schema.StringAttribute{
				Computed:    true,
				Description: "A unique identifier for the client certificate version.",
			},
			"auto_acknowledge_warnings": schema.BoolAttribute{
				Optional:    true,
				Description: "If set to true, all warnings will be acknowledged automatically. Defaults to `false`.",
			},
		},
	}
}

// clientCertificateUploadModel is a model for akamai_mtlskeystore_client_certificate_upload resource.
type clientCertificateUploadModel struct {
	ClientCertificateID     types.Int64  `tfsdk:"client_certificate_id"`
	VersionNumber           types.Int64  `tfsdk:"version_number"`
	SignedCertificate       types.String `tfsdk:"signed_certificate"`
	TrustChain              types.String `tfsdk:"trust_chain"`
	WaitForDeployment       types.Bool   `tfsdk:"wait_for_deployment"`
	GUID                    types.String `tfsdk:"version_guid"`
	AutoAcknowledgeWarnings types.Bool   `tfsdk:"auto_acknowledge_warnings"`
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

// Create implements resource's Create method.
func (r *clientCertificateUploadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state clientCertificateUploadModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	waitForDeployment := true
	if !state.WaitForDeployment.IsNull() {
		waitForDeployment = state.WaitForDeployment.ValueBool()
		state.WaitForDeployment = types.BoolValue(waitForDeployment)
	}

	client := Client(r.meta)
	uploadReq := mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: state.ClientCertificateID.ValueInt64(),
		Version:       state.VersionNumber.ValueInt64(),
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: state.SignedCertificate.ValueString(),
			TrustChain:  state.TrustChain.ValueStringPointer(),
		},
		AcknowledgeAllWarnings: state.AutoAcknowledgeWarnings.ValueBoolPointer(),
	}
	uploadedVersion, err := r.upsertClientCertificateUpload(ctx, client, uploadReq, state.VersionNumber.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error uploading signed certificate: ", err.Error())
		return
	}
	state.GUID = types.StringValue(uploadedVersion.VersionGUID)

	// Wait for deployment if needed
	if waitForDeployment && (uploadedVersion.Status != mtlskeystore.Deployed) {
		uploadedVersion, err = pollForCertificateDeployment(
			ctx,
			client,
			mtlskeystore.GetClientCertificateVersionsRequest{
				CertificateID: state.ClientCertificateID.ValueInt64(),
			},
			state.VersionNumber.ValueInt64(),
		)

		if err != nil {
			resp.Diagnostics.AddError("Error polling for client certificate deployment", err.Error())
			return
		}
		state.GUID = types.StringValue(uploadedVersion.VersionGUID)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read implements resource's Read method.
func (r *clientCertificateUploadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clientCertificateUploadModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(r.meta)
	getVersionsRequest := mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: state.ClientCertificateID.ValueInt64(),
	}
	clientCertificateVersionsResp, err := client.GetClientCertificateVersions(ctx, getVersionsRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving client certificate versions",
			err.Error(),
		)
		return
	}
	if clientCertificateVersionsResp == nil || len(clientCertificateVersionsResp.Versions) == 0 {
		// Resource drift: remove from state
		resp.State.RemoveResource(ctx)
		return
	}
	for _, version := range clientCertificateVersionsResp.Versions {
		if version.Version == state.VersionNumber.ValueInt64() {
			state.GUID = types.StringValue(version.VersionGUID)
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
			return
		}
	}
	resp.Diagnostics.AddError(
		"Uploaded version not found",
		"Could not find the uploaded version in the response from the client",
	)
}

// Delete implements resource's Delete method.
func (r *clientCertificateUploadResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No action needed on the server side, just remove from state
	resp.State.RemoveResource(ctx)
}

// Update implements resource's Update method.
func (r *clientCertificateUploadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state clientCertificateUploadModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	waitForDeployment := true
	if !plan.WaitForDeployment.IsNull() {
		waitForDeployment = state.WaitForDeployment.ValueBool()
		state.WaitForDeployment = plan.WaitForDeployment

	}

	// Only allow update if version_number changes
	if plan.VersionNumber.ValueInt64() == state.VersionNumber.ValueInt64() {
		resp.Diagnostics.AddError(
			"Update Not Supported",
			"Only updates with a different version_number are supported.",
		)
		return
	}

	client := Client(r.meta)
	uploadReq := mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: plan.ClientCertificateID.ValueInt64(),
		Version:       plan.VersionNumber.ValueInt64(),
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: plan.SignedCertificate.ValueString(),
			TrustChain:  plan.TrustChain.ValueStringPointer(),
		},
		AcknowledgeAllWarnings: plan.AutoAcknowledgeWarnings.ValueBoolPointer(),
	}
	uploadedVersion, err := r.upsertClientCertificateUpload(ctx, client, uploadReq, plan.VersionNumber.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error updating client certificate upload", err.Error())
		return
	}
	plan.GUID = types.StringValue(uploadedVersion.VersionGUID)

	// Wait for deployment if needed
	if waitForDeployment && (uploadedVersion.Status != mtlskeystore.Deployed) {
		uploadedVersion, err = pollForCertificateDeployment(
			ctx,
			client,
			mtlskeystore.GetClientCertificateVersionsRequest{
				CertificateID: plan.ClientCertificateID.ValueInt64(),
			},
			plan.VersionNumber.ValueInt64(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Error polling for client certificate deployment", err.Error())
			return
		}
		plan.GUID = types.StringValue(uploadedVersion.VersionGUID)

	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *clientCertificateUploadResource) upsertClientCertificateUpload(
	ctx context.Context,
	client mtlskeystore.MTLSKeystore,
	request mtlskeystore.UploadSignedClientCertificateRequest,
	versionNumber int64,
) (*mtlskeystore.ClientCertificateVersion, error) {
	if err := client.UploadSignedClientCertificate(ctx, request); err != nil {
		return nil, fmt.Errorf("error uploading client certificate: %w", err)
	}

	getVersionsRequest := mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: request.CertificateID,
	}
	clientCertificateVersionsResp, err := client.GetClientCertificateVersions(ctx, getVersionsRequest)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client certificate versions: %w", err)
	}
	if clientCertificateVersionsResp == nil || len(clientCertificateVersionsResp.Versions) == 0 {
		return nil, fmt.Errorf("no client certificate versions found: received nil or empty response from client")
	}
	var uploadedVersion *mtlskeystore.ClientCertificateVersion
	for _, version := range clientCertificateVersionsResp.Versions {
		if version.Version == versionNumber {
			uploadedVersion = &version
			break
		}
	}
	if uploadedVersion == nil || uploadedVersion.VersionGUID == "" {
		return nil, fmt.Errorf("uploaded version not found: could not find the uploaded version in the response from the client")
	}

	if uploadedVersion.Status != mtlskeystore.AwaitingSigned &&
		uploadedVersion.Status != mtlskeystore.DeploymentPending && uploadedVersion.Status != mtlskeystore.Deployed {
		return nil, fmt.Errorf(
			"unexpected client certificate version status: expected status to be either 'AWAITING_SIGNED_CERTIFICATE', 'DEPLOYMENT_PENDING' or 'DEPLOYED', but got: %s",
			uploadedVersion.Status,
		)
	}
	return uploadedVersion, nil
}

func pollForCertificateDeployment(
	ctx context.Context,
	client mtlskeystore.MTLSKeystore,
	getVersionsRequest mtlskeystore.GetClientCertificateVersionsRequest,
	versionNumber int64,
) (*mtlskeystore.ClientCertificateVersion, error) {
	for i := 0; i < numberOfRetries; i++ {
		clientCertificateVersionsResp, err := client.GetClientCertificateVersions(ctx, getVersionsRequest)
		if err != nil {
			return nil, fmt.Errorf("error retrieving client certificate versions: %w", err)
		}
		if clientCertificateVersionsResp == nil || len(clientCertificateVersionsResp.Versions) == 0 {
			return nil, fmt.Errorf("no client certificate versions found")
		}
		for _, version := range clientCertificateVersionsResp.Versions {
			if version.Version == versionNumber {
				if version.Status == mtlskeystore.Deployed {
					return &version, nil
				}
			}
		}
		time.Sleep(pollingInterval)
	}
	return nil, fmt.Errorf("timeout waiting for client certificate deployment: exceeded %d retries", numberOfRetries)
}
