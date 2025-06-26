package mtlstruststore

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/schema/nullstringdefault"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &caSetResource{}
	_ resource.ResourceWithConfigure   = &caSetResource{}
	_ resource.ResourceWithModifyPlan  = &caSetResource{}
	_ resource.ResourceWithImportState = &caSetResource{}
)

type caSetResource struct {
	meta          meta.Meta
	deleteTimeout time.Duration
}

// NewCASetResource returns a new akamai_mtlstruststore_ca_set resource.
func NewCASetResource() resource.Resource {
	return &caSetResource{
		deleteTimeout: 1 * time.Hour,
	}
}

type caSetResourceModel struct {
	Name                types.String   `tfsdk:"name"`
	Description         types.String   `tfsdk:"description"`
	AccountID           types.String   `tfsdk:"account_id"`
	ID                  types.String   `tfsdk:"id"`
	CreatedBy           types.String   `tfsdk:"created_by"`
	CreatedDate         types.String   `tfsdk:"created_date"`
	VersionCreatedBy    types.String   `tfsdk:"version_created_by"`
	VersionCreatedDate  types.String   `tfsdk:"version_created_date"`
	VersionModifiedBy   types.String   `tfsdk:"version_modified_by"`
	VersionModifiedDate types.String   `tfsdk:"version_modified_date"`
	AllowInsecureSHA1   types.Bool     `tfsdk:"allow_insecure_sha1"`
	VersionDescription  types.String   `tfsdk:"version_description"`
	LatestVersion       types.Int64    `tfsdk:"latest_version"`
	StagingVersion      types.Int64    `tfsdk:"staging_version"`
	ProductionVersion   types.Int64    `tfsdk:"production_version"`
	Certificates        types.Set      `tfsdk:"certificates"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

type certificateModel struct {
	CertificatePEM     types.String `tfsdk:"certificate_pem"`
	Description        types.String `tfsdk:"description"`
	CreatedBy          types.String `tfsdk:"created_by"`
	CreatedDate        types.String `tfsdk:"created_date"`
	StartDate          types.String `tfsdk:"start_date"`
	EndDate            types.String `tfsdk:"end_date"`
	Fingerprint        types.String `tfsdk:"fingerprint"`
	Issuer             types.String `tfsdk:"issuer"`
	SerialNumber       types.String `tfsdk:"serial_number"`
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	Subject            types.String `tfsdk:"subject"`
}

func (r *caSetResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set"
}

func (r *caSetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the CA set.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(regexp.MustCompile(mtlstruststore.CASetNamePattern), "allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.)"),
				},
				PlanModifiers: []planmodifier.String{modifiers.PreventStringUpdate()},
			},
			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Validators:    []validator.String{stringvalidator.LengthAtMost(255)},
				PlanModifiers: []planmodifier.String{modifiers.PreventStringUpdate()},
				Default:       nullstringdefault.NullString(),
				Description:   "Any additional comments you can add to the CA set.",
			},
			"account_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Identifies the account the CA set belongs to.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Uniquely identifies the CA set.",
			},
			"created_by": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The user who created the CA set.",
			},
			"created_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "When the CA set was created.",
			},
			"version_created_by": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The user who created the CA set version.",
			},
			"version_created_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "When the CA set version was created.",
			},
			"version_modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who modified the CA set version.",
			},
			"version_modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "When the CA set version was modified.",
			},
			"allow_insecure_sha1": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Allows certificates with SHA-1 signatures if enabled.",
			},
			"version_description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.LengthAtMost(255)},
				Default:     nullstringdefault.NullString(),
				Description: "Additional description for the CA set version.",
			},
			"latest_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version number for newly created or cloned version in a CA set.",
			},
			"staging_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version number of the CA set that is active on staging.",
			},
			"production_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version of the CA set that is active on production.",
			},
			"certificates": certificatesSchema(),
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Delete:            true,
				CreateDescription: "Optional configurable resource delete timeout. By default it's 1h with 15m polling interval.",
			}),
		},
	}
}

func certificatesSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Required:    true,
		Validators:  []validator.Set{setvalidator.SizeBetween(1, 300)},
		Description: "The certificates that are valid, non-expired, root, or intermediate.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"certificate_pem": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`-----BEGIN CERTIFICATE-----\n[0-9A-Za-z+/=\s]+\n-----END CERTIFICATE-----`), "Certificate must be in PEM format")},
					Description: "The certificate in PEM format, as found in a Base64 ASCII encoded file.",
				},
				"description": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Validators:  []validator.String{stringvalidator.LengthAtMost(255)},
					Description: "Optional description for the certificate.",
				},
				"created_by": schema.StringAttribute{
					Computed:    true,
					Description: "The user who created this CA certificate.",
				},
				"created_date": schema.StringAttribute{
					Computed:    true,
					Description: "When the CA certificate was created.",
				},
				"start_date": schema.StringAttribute{
					Computed:    true,
					Description: "The start date of the certificate.",
				},
				"end_date": schema.StringAttribute{
					Computed:    true,
					Description: "The certificate's ISO 8601 formatted expiration date.",
				},
				"fingerprint": schema.StringAttribute{
					Computed:    true,
					Description: "The fingerprint of the certificate.",
				},
				"issuer": schema.StringAttribute{
					Computed:    true,
					Description: "The certificate's issuer.",
				},
				"serial_number": schema.StringAttribute{
					Computed:    true,
					Description: "The unique serial number of the certificate.",
				},
				"signature_algorithm": schema.StringAttribute{
					Computed:    true,
					Description: "The signature algorithm of the CA certificate.",
				},
				"subject": schema.StringAttribute{
					Computed:    true,
					Description: "The certificate's subject field.",
				},
			},
		},
	}
}

func (r *caSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

func (r *caSetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	tflog.Debug(ctx, "Validating CA set resource configuration")

	if r.meta == nil {
		return
	}
	var config caSetResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if strings.Contains(config.Name.ValueString(), "...") {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid CA set name", "CA set name cannot contain three consecutive periods (...)")
		return
	}
	client = Client(r.meta)

	if diags := validateCerts(ctx, client, &config); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func generateDiagnosticForValidationErrors(err error) diag.Diagnostics {
	var diags diag.Diagnostics
	var e *mtlstruststore.Error
	if errors.As(err, &e) {
		// If there are no details on every certificate, we can just add the general error.
		if len(e.Errors) == 0 {
			diags.AddAttributeError(path.Root("certificates"), "Certificates are invalid", e.Title)
			return diags
		}
		for _, ee := range e.Errors {
			index, err := strconv.Atoi(strings.TrimPrefix(ee.Pointer, "/certificates/"))
			// If we cannot parse the index, we can just add the general error
			if err != nil {
				diags.AddAttributeError(path.Root("certificates"), "Certificates are invalid. Unknown pointer: "+ee.Pointer, err.Error())
			} else {
				diags.AddAttributeError(path.Root("certificates").AtListIndex(index), "Certificate is invalid", ee.Detail)
			}
		}
	} else {
		diags.AddAttributeError(path.Root("certificates"), "Certificates are invalid", err.Error())
	}
	return diags
}

func (r *caSetResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		// Verify if ca set is in use before deleting.
		var state *caSetResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		client = Client(r.meta)

		associationsResponse, err := client.ListCASetAssociations(ctx, mtlstruststore.ListCASetAssociationsRequest{
			CASetID: state.ID.ValueString(),
		})
		if err != nil {
			if errors.Is(err, mtlstruststore.ErrGetCASetNotFound) {
				tflog.Debug(ctx, "CA set is not found, we can mark the resource as deleted")
				return
			}
			resp.Diagnostics.AddError("ca set resource ModifyPlan failed", err.Error())
			return
		}
		if len(associationsResponse.Associations.Enrollments) > 0 || len(associationsResponse.Associations.Properties) > 0 {
			tflog.Warn(ctx, "CA set is in use, cannot delete it", map[string]interface{}{
				"ca_set_id":   state.ID.ValueString(),
				"enrollments": len(associationsResponse.Associations.Enrollments),
				"properties":  len(associationsResponse.Associations.Properties),
			})
			resp.Diagnostics.AddWarning("CA set is in use and cannot be deleted", getAssociationDetails(associationsResponse))
			return
		}
	}
	// Empty state means that the resource is being created.
	if req.State.Raw.IsNull() {
		var plan *caSetResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// If the certificate description is not set, we have to set it to null.
		var planCerts []certificateModel
		diags := plan.Certificates.ElementsAs(ctx, &planCerts, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		for i, pc := range planCerts {
			if !pc.CertificatePEM.IsUnknown() && !pc.CertificatePEM.IsNull() {
				if pc.Description.IsNull() || pc.Description.IsUnknown() {
					planCerts[i].Description = types.StringNull()
				}
			}
		}
		updatedCerts, diags := types.SetValueFrom(ctx, certificatesType(), planCerts)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		diags = resp.Plan.SetAttribute(ctx, path.Root("certificates"), updatedCerts)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	// Update
	if !req.Plan.Raw.IsNull() && !req.State.Raw.IsNull() {
		var state, plan *caSetResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
		// if only timeouts are changed, we can just ignore the plan
		if state.onlyTimeoutChanged(plan) {
			plan.Timeouts = state.Timeouts
		}
		return
	}
}

func getAssociationDetails(associationsResponse *mtlstruststore.ListCASetAssociationsResponse) string {
	var details string
	if len(associationsResponse.Associations.Enrollments) > 0 {
		details = fmt.Sprintf("CA set is in use by %d enrollments: ", len(associationsResponse.Associations.Enrollments))
		for i, enrollment := range associationsResponse.Associations.Enrollments {
			if i > 0 {
				details += ", "
			}
			details += fmt.Sprintf("%s (%d)", enrollment.CN, enrollment.EnrollmentID)
		}
		details += "\n"
	} else {
		details = fmt.Sprintf("CA set is in use by %d properties: ", len(associationsResponse.Associations.Properties))
		for i, property := range associationsResponse.Associations.Properties {
			if i > 0 {
				details += ", "
			}
			var propertyName string
			if property.PropertyName != nil {
				propertyName = *property.PropertyName
			}
			details += fmt.Sprintf("%s (%s)", propertyName, property.PropertyID)
		}
	}
	return details
}

func (m *caSetResourceModel) onlyTimeoutChanged(plan *caSetResourceModel) bool {
	if m == nil || plan == nil {
		return false
	}
	if !m.Name.Equal(plan.Name) {
		return false
	}
	if !m.Description.Equal(plan.Description) {
		return false
	}
	if !m.AllowInsecureSHA1.Equal(plan.AllowInsecureSHA1) {
		return false
	}
	if !m.VersionDescription.Equal(plan.VersionDescription) {
		return false
	}
	if !m.Certificates.Equal(plan.Certificates) {
		return false
	}
	if !m.Timeouts.Equal(plan.Timeouts) {
		return true
	}
	return false
}

func validateCerts(ctx context.Context, client mtlstruststore.MTLSTruststore, config *caSetResourceModel) diag.Diagnostics {
	var certificates []certificateModel
	diags := config.Certificates.ElementsAs(ctx, &certificates, false)
	if diags.HasError() {
		return diags
	}
	var certs []mtlstruststore.ValidateCertificate
	for _, cert := range certificates {
		// Certificate content can be provided from external sources and may not be known at plan time
		if !cert.CertificatePEM.IsUnknown() && !cert.CertificatePEM.IsNull() {
			certs = append(certs, mtlstruststore.ValidateCertificate{
				CertificatePEM: cert.CertificatePEM.ValueString(),
				Description:    cert.Description.ValueStringPointer(),
			})
		}
	}

	// If there are no certificates, we can skip validation
	if len(certs) == 0 {
		return nil
	}
	_, err := client.ValidateCertificates(ctx, mtlstruststore.ValidateCertificatesRequest{
		AllowInsecureSHA1: config.AllowInsecureSHA1.ValueBool(),
		Certificates:      certs,
	})
	if err != nil {
		return generateDiagnosticForValidationErrors(err)
	}
	return nil
}

func (r *caSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating CA set resource")
	var plan caSetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client = Client(r.meta)
	if diags := validateCerts(ctx, client, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if diags := r.create(ctx, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *caSetResource) create(ctx context.Context, plan *caSetResourceModel) diag.Diagnostics {
	client = Client(r.meta)

	createCASetResponse, err := client.CreateCASet(ctx, mtlstruststore.CreateCASetRequest{
		CASetName:   plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
	})
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("create ca set failed", err.Error())}
	}
	tflog.Debug(ctx, "Created CA set", map[string]interface{}{
		"ca_set_id": createCASetResponse.CASetID,
	})

	certs, diags := createCertificatesRequest(ctx, plan.Certificates)
	if diags.HasError() {
		return diags
	}
	createCASetVersionResponse, err := client.CreateCASetVersion(ctx, mtlstruststore.CreateCASetVersionRequest{
		CASetID: createCASetResponse.CASetID,
		Body: mtlstruststore.CreateCASetVersionRequestBody{
			AllowInsecureSHA1: plan.AllowInsecureSHA1.ValueBool(),
			Description:       plan.VersionDescription.ValueStringPointer(),
			Certificates:      certs,
		},
	})
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("create ca set version failed", err.Error())}
	}
	tflog.Debug(ctx, "Created CA set version", map[string]interface{}{
		"current_version": strconv.FormatInt(createCASetVersionResponse.Version, 10),
	})
	// After creating the CA set, version is empty. Version appears after creating the CA set version.
	// To get the CA set version, we need to fetch it again.
	getCASetResponse, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
		CASetID: createCASetResponse.CASetID,
	})
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("get ca set failed", err.Error())}
	}

	plan.setCASetData((*mtlstruststore.CASetResponse)(getCASetResponse))
	diags = plan.setCASetVersionData(ctx, (*mtlstruststore.CASetVersion)(createCASetVersionResponse))

	return diags
}

func createCertificatesRequest(ctx context.Context, c types.Set) ([]mtlstruststore.CertificateRequest, diag.Diagnostics) {
	var certificates []certificateModel
	diags := c.ElementsAs(ctx, &certificates, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]mtlstruststore.CertificateRequest, 0, len(certificates))
	for _, cert := range certificates {
		result = append(result, mtlstruststore.CertificateRequest{
			CertificatePEM: cert.CertificatePEM.ValueString(),
			Description:    cert.Description.ValueStringPointer(),
		})
	}
	return result, nil
}

func (r *caSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading CA Set Resource")
	var state caSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client = Client(r.meta)

	caSetResp, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
		CASetID: state.ID.ValueString(),
	})
	if err != nil {
		if errors.Is(err, mtlstruststore.ErrGetCASetNotFound) {
			tflog.Debug(ctx, "CA set is not found, we can mark the resource as deleted")
			resp.Diagnostics.AddWarning("CA set is not found, we can mark the resource as deleted", fmt.Sprintf("CA set with ID %s is not found. It may have been deleted outside of Terraform.", state.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("read ca set error", err.Error())
		return
	}
	state.setCASetData((*mtlstruststore.CASetResponse)(caSetResp))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if caSetResp.LatestVersion == nil {
		tflog.Debug(ctx, "CA set or version is not found")
		resp.Diagnostics.AddWarning("CA set version is not found, we can mark the resource as deleted", fmt.Sprintf("CA set with ID %s has no version. It may have been deleted outside of Terraform.", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	caSetVersionResp, err := client.GetCASetVersion(ctx, mtlstruststore.GetCASetVersionRequest{
		CASetID: state.ID.ValueString(),
		Version: *caSetResp.LatestVersion,
	})
	if err != nil {
		if errors.Is(err, mtlstruststore.ErrGetCASetVersionNotFound) || errors.Is(err, mtlstruststore.ErrMissingCASetVersion) || errors.Is(err, mtlstruststore.ErrGetCASetNotFound) {
			tflog.Debug(ctx, "CA set version is not found")
			resp.Diagnostics.AddWarning("CA set version is not found, we can mark the resource as deleted", fmt.Sprintf("CA set with ID %s has no version. It may have been deleted outside of Terraform.", state.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("read ca set error", err.Error())
		return
	}

	diags := state.setCASetVersionData(ctx, (*mtlstruststore.CASetVersion)(caSetVersionResp))
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *caSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating CA Set Resource")
	var plan, state caSetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.onlyTimeoutChanged(&plan) {
		state.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	client = Client(r.meta)
	if diags := validateCerts(ctx, client, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if diags := r.update(ctx, &plan, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *caSetResource) update(ctx context.Context, plan, state *caSetResourceModel) diag.Diagnostics {
	client = Client(r.meta)

	// We cannot update CA Set if it is active, or it was activated in the past.
	// In such a situation we need to clone the current version and update the new one.
	activations, err := client.ListCASetActivations(ctx, mtlstruststore.ListCASetActivationsRequest{
		CASetID: state.ID.ValueString(),
	})
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("list ca set activations failed", err.Error())}
	}
	latestVersionNeverActivated := true
	for _, activation := range activations.Activations {
		if activation.Version == state.LatestVersion.ValueInt64() {
			latestVersionNeverActivated = false
			break
		}
	}

	if latestVersionNeverActivated {
		// Current version is not activated in staging or production, so we can just update the current version
		tflog.Debug(ctx, "Current version is not activated in staging or production, updating the current version.", map[string]interface{}{
			"ca_set_id":       state.ID.ValueString(),
			"current_version": strconv.FormatInt(state.LatestVersion.ValueInt64(), 10),
		})
		plan.LatestVersion = state.LatestVersion
	} else {
		// Current version is/was already activated in staging or production, so we need to create a new version by cloning the current one
		clonedCASetVersionResp, err := client.CloneCASetVersion(ctx, mtlstruststore.CloneCASetVersionRequest{
			CASetID: state.ID.ValueString(),
			Version: state.LatestVersion.ValueInt64(),
		})
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("clone ca set version failed", err.Error())}
		}
		plan.LatestVersion = types.Int64Value(clonedCASetVersionResp.Version)
		tflog.Debug(ctx, "Current version is or was activated on staging or production, creating a new version by cloning current one.", map[string]interface{}{
			"ca_set_id":       state.ID.ValueString(),
			"current_version": strconv.FormatInt(state.LatestVersion.ValueInt64(), 10),
			"new_version":     strconv.FormatInt(plan.LatestVersion.ValueInt64(), 10),
		})
	}

	certs, diags := createCertificatesRequest(ctx, plan.Certificates)
	if diags.HasError() {
		return diags
	}
	caSetVersionResp, err := client.UpdateCASetVersion(ctx, mtlstruststore.UpdateCASetVersionRequest{
		CASetID: plan.ID.ValueString(),
		Version: plan.LatestVersion.ValueInt64(),
		Body: mtlstruststore.UpdateCASetVersionRequestBody{
			AllowInsecureSHA1: plan.AllowInsecureSHA1.ValueBool(),
			Description:       plan.VersionDescription.ValueString(),
			Certificates:      certs,
		},
	})
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("update ca set version failed", err.Error())}
	}

	plan.StagingVersion = state.StagingVersion
	plan.ProductionVersion = state.ProductionVersion
	return plan.setCASetVersionData(ctx, (*mtlstruststore.CASetVersion)(caSetVersionResp))
}

func (r *caSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting CA Set Resource")

	var state *caSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client = Client(r.meta)
	// Check if the CA set is in use before deleting
	associationsResponse, err := client.ListCASetAssociations(ctx, mtlstruststore.ListCASetAssociationsRequest{
		CASetID: state.ID.ValueString(),
	})
	if err != nil {
		if errors.Is(err, mtlstruststore.ErrGetCASetNotFound) {
			tflog.Debug(ctx, "CA set is not found, we can mark the resource as deleted")
			return
		}
		resp.Diagnostics.AddError("delete ca set resource failed", err.Error())
		return
	}
	if len(associationsResponse.Associations.Enrollments) > 0 || len(associationsResponse.Associations.Properties) > 0 {
		resp.Diagnostics.AddError("CA set is in use and cannot be deleted", getAssociationDetails(associationsResponse))
		return
	}

	// initiate deletion process
	if err := client.DeleteCASet(ctx, mtlstruststore.DeleteCASetRequest{
		CASetID: state.ID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("delete ca set %s failed", state.ID), err.Error())
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, r.deleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	// Initial value, to be changed based on the response from GetCASetDeletionStatus later
	const initialInterval = 10 * time.Millisecond
	const defaultInterval = 10 * time.Second
	// get status as soon as possible
	var CASetDeleteStatusPollInterval = initialInterval
	for deleteInProgress := true; deleteInProgress; {
		select {
		case <-time.After(CASetDeleteStatusPollInterval):
			status, err := client.GetCASetDeletionStatus(ctx, mtlstruststore.GetCASetDeletionStatusRequest{
				CASetID: state.ID.ValueString(),
			})
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("get ca set %s deletion status failed", state.ID.ValueString()), err.Error())
				return
			}

			switch status.Status {
			case mtlstruststore.DeletionStatusComplete:
				deleteInProgress = false
			case mtlstruststore.DeletionStatusInProgress:
				tflog.Debug(ctx, fmt.Sprintf("delete ca set %s in progress", state.ID))
				if !status.RetryAfter.IsZero() {
					CASetDeleteStatusPollInterval = time.Until(status.RetryAfter)
				} else {
					CASetDeleteStatusPollInterval = defaultInterval
				}
			case mtlstruststore.DeletionStatusFailed:
				// In case of the failure, there is no point to retry deletion, we can just return the error
				// Any failure has to be handled manually by the support team
				var reason string
				if status.FailureReason != nil {
					reason = *status.FailureReason
				}

				resp.Diagnostics.AddError(fmt.Sprintf("delete ca set %s failed", state.ID), "contact support team to resolve the issue. "+reason)
				return
			}
		case <-ctx.Done():
			resp.Diagnostics.AddError("delete ca set context terminated", ctx.Err().Error())
			return
		}
	}
}

func (r *caSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing CA Set Resource")

	caSetID := req.ID
	if caSetID == "" {
		resp.Diagnostics.AddError("empty import ID", "Import ID cannot be empty")
		return
	}

	client = Client(r.meta)
	caSetResp, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
		CASetID: caSetID,
	})
	if err != nil {
		resp.Diagnostics.AddError("import ca set resource failed", err.Error())
		return
	}
	if caSetResp.LatestVersion == nil {
		resp.Diagnostics.AddError("It is not possible to import ca set without version", "The CA set does not have any version")
		return
	}

	data := &caSetResourceModel{
		ID:            types.StringValue(caSetID),
		LatestVersion: types.Int64Value(*caSetResp.LatestVersion),
		Certificates:  types.SetNull(certificatesType()),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"delete": types.StringType,
			}),
		},
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *caSetResourceModel) setCASetData(caSetResp *mtlstruststore.CASetResponse) {
	m.ID = types.StringValue(caSetResp.CASetID)
	m.Name = types.StringValue(caSetResp.CASetName)
	m.CreatedBy = types.StringValue(caSetResp.CreatedBy)
	m.CreatedDate = types.StringValue(caSetResp.CreatedDate.Format(time.RFC3339Nano))
	m.AccountID = types.StringValue(caSetResp.AccountID)
	m.Description = types.StringPointerValue(caSetResp.Description)
	m.LatestVersion = types.Int64PointerValue(caSetResp.LatestVersion)
	m.StagingVersion = types.Int64PointerValue(caSetResp.StagingVersion)
	m.ProductionVersion = types.Int64PointerValue(caSetResp.ProductionVersion)
}

func (m *caSetResourceModel) setCASetVersionData(ctx context.Context, caSetVersionResp *mtlstruststore.CASetVersion) diag.Diagnostics {
	m.AllowInsecureSHA1 = types.BoolValue(caSetVersionResp.AllowInsecureSHA1)
	m.VersionDescription = types.StringPointerValue(caSetVersionResp.Description)
	m.VersionModifiedBy = types.StringPointerValue(caSetVersionResp.ModifiedBy)
	if caSetVersionResp.ModifiedDate == nil {
		m.VersionModifiedDate = types.StringNull()
	} else {
		m.VersionModifiedDate = types.StringValue(caSetVersionResp.ModifiedDate.Format(time.RFC3339Nano))
	}
	m.VersionCreatedBy = types.StringValue(caSetVersionResp.CreatedBy)
	m.VersionCreatedDate = types.StringValue(caSetVersionResp.CreatedDate.Format(time.RFC3339Nano))
	certificates := make([]certificateModel, 0, len(caSetVersionResp.Certificates))
	for _, cert := range caSetVersionResp.Certificates {
		certificates = append(certificates, certificateModel{
			CertificatePEM:     types.StringValue(cert.CertificatePEM),
			Description:        types.StringPointerValue(cert.Description),
			CreatedBy:          types.StringValue(cert.CreatedBy),
			CreatedDate:        types.StringValue(cert.CreatedDate.Format(time.RFC3339Nano)),
			StartDate:          types.StringValue(cert.StartDate.Format(time.RFC3339Nano)),
			EndDate:            types.StringValue(cert.EndDate.Format(time.RFC3339Nano)),
			Fingerprint:        types.StringValue(cert.Fingerprint),
			Issuer:             types.StringValue(cert.Issuer),
			SerialNumber:       types.StringValue(cert.SerialNumber),
			SignatureAlgorithm: types.StringValue(cert.SignatureAlgorithm),
			Subject:            types.StringValue(cert.Subject),
		})
	}
	certs, diags := types.SetValueFrom(ctx, certificatesType(), certificates)
	if diags.HasError() {
		return diags
	}
	m.Certificates = certs
	return nil
}

func certificatesType() types.ObjectType {
	return certificatesSchema().NestedObject.Type().(types.ObjectType)
}
