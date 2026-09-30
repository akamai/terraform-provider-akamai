package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/text"
	fwdate "github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/validators"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &lineageResource{}
	_ resource.ResourceWithConfigure   = &lineageResource{}
	_ resource.ResourceWithImportState = &lineageResource{}
	_ resource.ResourceWithModifyPlan  = &lineageResource{}
)

// lineageActivationReminderWarningType is the ValidationResults warning CreateLineage always returns on a
// brand-new lineage, regardless of whether the same apply goes on to activate it right after: nothing has been
// promoted yet at the moment CreateLineage itself returns.
const lineageActivationReminderWarningType = "/error-types/lineage-activation-warning"

// notBlank matches EG's Subject validation, which rejects whitespace-only values.
var notBlank = regexp.MustCompile(`\S`)

// noSurroundingWhitespace rejects leading/trailing whitespace (and whitespace-only values), while still
// allowing whitespace between other characters.
var noSurroundingWhitespace = regexp.MustCompile(`^\S(.*\S)?$`)

// validSAN matches a domain name with an optional leading wildcard label, e.g. example.com or *.example.com.
// This mirrors the API's own JSON schema regex for sans, confirmed via a real
// 400 schema-validation-failure response. Uppercase letters are rejected: the API lowercases every SAN it
// stores, so accepting them here would silently diverge from what ends up in state after the next refresh.
var validSAN = regexp.MustCompile(`^(\*\.)?([a-z0-9]([a-z0-9-]*[a-z0-9])?\.)+[a-z]{2,}$`)

// noLeadingWhitespace mirrors the API's subject field pattern, which rejects a
// leading whitespace character but otherwise allows any content.
var noLeadingWhitespace = regexp.MustCompile(`^\S.*$`)

// validCountryCode mirrors the API's subject.country pattern: exactly two letters.
var validCountryCode = regexp.MustCompile(`^[A-Za-z]{2}$`)

// validLineageName mirrors the API's lineageName pattern.
var validLineageName = regexp.MustCompile(`^[0-9a-zA-Z _.-]*$`)

type (
	lineageResource struct {
		meta.Resource
	}

	lineageResourceModel struct {
		LineageID           types.Int64  `tfsdk:"lineage_id"`
		AccountID           types.String `tfsdk:"account_id"`
		ContractID          types.String `tfsdk:"contract_id"`
		GroupID             types.Int64  `tfsdk:"group_id"`
		GeoClass            types.String `tfsdk:"geo_class"`
		LineageName         types.String `tfsdk:"lineage_name"`
		LineageType         types.String `tfsdk:"lineage_type"`
		SecureNetwork       types.String `tfsdk:"secure_network"`
		StackMode           types.String `tfsdk:"stack_mode"`
		KeySpecs            types.Map    `tfsdk:"key_specs"`
		SANs                types.Set    `tfsdk:"sans"`
		Subject             types.Object `tfsdk:"subject"`
		LineageCreatedBy    types.String `tfsdk:"lineage_created_by"`
		LineageCreatedTime  types.String `tfsdk:"lineage_created_time"`
		LineageModifiedBy   types.String `tfsdk:"lineage_modified_by"`
		LineageModifiedTime types.String `tfsdk:"lineage_modified_time"`
		Head                types.Object `tfsdk:"head"`
		CurrentProduction   types.Object `tfsdk:"current_production"`
		PreviousProduction  types.Object `tfsdk:"previous_production"`
		CurrentStaging      types.Object `tfsdk:"current_staging"`
		SigningTarget       types.Object `tfsdk:"signing_target"`
	}
)

// NewLineageResource returns a new Cloud Certificates lineage resource.
func NewLineageResource() resource.Resource {
	return &lineageResource{}
}

// Metadata configures the resource's type name.
func (r *lineageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_lineage"
}

// Schema defines the Terraform schema for the lineage resource.
func (r *lineageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	generationPointerAttrs := func(description string) schema.SingleNestedAttribute {
		return schema.SingleNestedAttribute{
			Computed:    true,
			Description: description,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"generation_id": schema.Int64Attribute{
					Computed:    true,
					Description: "Unique identifier of this generation.",
				},
				"generation_status": schema.StringAttribute{
					Computed:    true,
					Description: "Status of this generation.",
				},
				"algorithms": schema.MapNestedAttribute{
					Computed:    true,
					Description: "Per key-type (RSA or ECDSA) certificate details for this generation, keyed by key_type.",
					NestedObject: schema.NestedAttributeObject{
						Attributes: algorithmResourceAttributes(),
					},
				},
				"first_promoted_to_production_time": schema.StringAttribute{
					Computed:    true,
					Description: "Time the generation was first promoted to production, in RFC3339 format. Null if never promoted.",
				},
				"generation_created_by": schema.StringAttribute{
					Computed:    true,
					Description: "Username of the person who created this generation.",
				},
				"generation_created_time": schema.StringAttribute{
					Computed:    true,
					Description: "Time the generation was created, in RFC3339 format.",
				},
				"generation_modified_by": schema.StringAttribute{
					Computed:    true,
					Description: "Username of the person who last modified this generation.",
				},
				"generation_modified_time": schema.StringAttribute{
					Computed:    true,
					Description: "Time the generation was last modified, in RFC3339 format. Null if never modified since creation.",
				},
			},
		}
	}

	resp.Schema = schema.Schema{
		Description: "Manages the lifecycle of a Cloud Certificate Manager certificate lineage.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier of the lineage, assigned by the API on creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Account identifier associated with the contract.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Akamai contract ID to associate with this lineage.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(notBlank, "must not be blank"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "Akamai group ID under the contract.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"geo_class": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Geographic traffic class. Defaults to STANDARD_WORLDWIDE if not set.",
				Validators: []validator.String{
					stringvalidator.OneOf(string(cloudcertificates.GeoClassStandardWorldwide)),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lineage_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Human-readable name of the lineage. If not set, the API generates one.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 270),
					stringvalidator.RegexMatches(noSurroundingWhitespace, "must not have leading or trailing whitespace"),
					stringvalidator.RegexMatches(validLineageName, "must contain only letters, digits, spaces, underscores, periods, and hyphens"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lineage_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Controls future renewals. Allowed values are `MULTIPLE_GENERATION` and " +
					"`SINGLE_GENERATION`. Defaults to `MULTIPLE_GENERATION`.",
				Validators: []validator.String{
					stringvalidator.OneOf(string(cloudcertificates.LineageTypeMultipleGeneration), string(cloudcertificates.LineageTypeSingleGeneration)),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secure_network": schema.StringAttribute{
				Required:    true,
				Description: "Secure network deployment tier. Allowed values are `STANDARD_TLS` and `ENHANCED_TLS`.",
				Validators: []validator.String{
					stringvalidator.OneOf(string(cloudcertificates.SecureNetworkStandardTLS), string(cloudcertificates.SecureNetworkEnhancedTLS)),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_specs": schema.MapAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Algorithm slots for each generation, mapping key_type (`RSA` and/or `ECDSA`) to key_size. Provide " +
					"one entry for a single-stack lineage, or two (one `RSA`, one `ECDSA`) for a multiple-stack lineage.",
				Validators: []validator.Map{
					mapvalidator.SizeBetween(1, 2),
					mapvalidator.NoNullValues(),
					KeySpecsValidator(),
				},
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"sans": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Subject Alternative Name domains for the certificate. Allows wildcards such as `*.example.com`. " +
					"The API does not preserve input order and lowercases every value, so this is a set rather than a list.",
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 100),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(253),
						stringvalidator.RegexMatches(validSAN, "must be a valid lowercase domain name or wildcard domain, e.g. example.com or *.example.com"),
					),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"subject": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional X.509 subject fields of the certificate.",
				Validators: []validator.Object{
					validators.NonEmptyObject(),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"common_name": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Common name (CN).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							stringvalidator.RegexMatches(noLeadingWhitespace, "must not have leading whitespace"),
						},
					},
					"organization": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Organization (O).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							stringvalidator.RegexMatches(noLeadingWhitespace, "must not have leading whitespace"),
						},
					},
					"organizational_unit": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Organizational unit (OU).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							stringvalidator.RegexMatches(noLeadingWhitespace, "must not have leading whitespace"),
						},
					},
					"country": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Two-letter ISO 3166 country code (C).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(2, 2),
							stringvalidator.RegexMatches(validCountryCode, "must contain only two letters"),
						},
					},
					"state": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "State or province name (ST).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							stringvalidator.RegexMatches(noLeadingWhitespace, "must not have leading whitespace"),
						},
					},
					"locality": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Locality or city name (L).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							stringvalidator.RegexMatches(noLeadingWhitespace, "must not have leading whitespace"),
						},
					},
				},
			},
			"stack_mode": schema.StringAttribute{
				Computed: true,
				Description: "Stack mode of the lineage, derived from the number of key_specs entries: " +
					"`SINGLE_STACK` for one entry, `MULTIPLE_STACK` for two.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lineage_created_by": schema.StringAttribute{
				Computed:    true,
				Description: "Username of the person who created the lineage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lineage_created_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time the lineage was created, in RFC3339 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lineage_modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "Username of the person who last modified the lineage.",
			},
			"lineage_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time the lineage was last modified, in RFC3339 format.",
			},
			"head":                generationPointerAttrs("Head generation of the lineage: the generation currently being prepared (CSR generation, certificate upload) ahead of activation. Null if none exists (e.g. once fully promoted to production with no renewal in progress)."),
			"current_production":  generationPointerAttrs("Generation currently deployed to the production network. Null if none has been promoted yet."),
			"previous_production": generationPointerAttrs("Generation previously deployed to production, the rollback candidate. Null if none exists yet."),
			"current_staging":     generationPointerAttrs("Generation currently deployed to the staging network. Null if none has been promoted yet."),
			"signing_target": generationPointerAttrs("Derived convenience field, not a distinct API concept: the generation whose CSR currently needs to be " +
				"signed and uploaded - the lineage's head generation if one exists, otherwise its current production generation " +
				"(e.g. while completing a MULTIPLE_STACK lineage's second algorithm). It is null if neither exists."),
		},
	}
}

// Create provisions a new certificate lineage.
func (r *lineageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Resource Create")

	var plan lineageResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	createReq, diags := planToCreateRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.Client.GetCloudCertificates().CreateLineage(ctx, createReq)
	if err != nil {
		switch {
		case errors.Is(err, cloudcertificates.ErrLineageAccountNotAllowed):
			resp.Diagnostics.AddError("Account Not Allowed", "This account is not permitted to use Certificate Lineage.")
		default:
			resp.Diagnostics.AddError("Failed to create lineage", err.Error())
		}
		return
	}

	// Surface DOM-validation warnings (e.g. unvalidated SANs) directly as diagnostics: GetLineage never returns
	// them on subsequent reads, so storing them in state would just go stale after the first refresh.
	if created.ValidationResults != nil {
		for _, w := range created.ValidationResults.Warnings {
			detail := w.Detail
			if w.Type == lineageActivationReminderWarningType {
				detail += " This is expected for every newly created lineage - if your configuration also activates it " +
					"(e.g. via akamai_cloudcertificates_activation) and that succeeds, you can safely ignore this warning."
			}
			resp.Diagnostics.AddWarning(w.Title, detail)
		}
	}

	data, diags := mapLineageToResourceModel(ctx, (*cloudcertificates.Lineage)(created))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the lineage state from the API.
func (r *lineageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Resource Read")

	var state lineageResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	lineage, err := r.Client.GetCloudCertificates().GetLineage(ctx, cloudcertificates.GetLineageRequest{
		LineageID:         state.LineageID.ValueInt64(),
		ExpandGenerations: allExpandGenerations,
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddWarning("Lineage Not Found",
				fmt.Sprintf("Removing lineage %d from state; it no longer exists.", state.LineageID.ValueInt64()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read lineage", err.Error())
		return
	}

	data, diags := mapLineageToResourceModel(ctx, (*cloudcertificates.Lineage)(lineage))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update renames the lineage. All other attributes require replacement and never reach this method.
func (r *lineageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Resource Update")

	var plan, state lineageResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	lineageID := state.LineageID.ValueInt64()
	client := r.Client.GetCloudCertificates()
	renamed, err := client.RenameLineage(ctx, cloudcertificates.RenameLineageRequest{
		LineageID:   lineageID,
		LineageName: plan.LineageName.ValueString(),
	})
	if err != nil {
		switch {
		case errors.Is(err, cloudcertificates.ErrLineageNotFound):
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", lineageID))
		case errors.Is(err, cloudcertificates.ErrLineageNameConflict):
			resp.Diagnostics.AddError("Lineage Name Conflict", fmt.Sprintf("A lineage named %q already exists for this account.", plan.LineageName.ValueString()))
		default:
			resp.Diagnostics.AddError("Failed to rename lineage", err.Error())
		}
		return
	}

	// RenameLineage's response does not include expanded generation detail (unlike GetLineage/CreateLineage), so
	// re-fetch via GetLineage to avoid wiping out algorithm/generation detail already present in state.
	lineage, err := client.GetLineage(ctx, cloudcertificates.GetLineageRequest{
		LineageID:         lineageID,
		ExpandGenerations: allExpandGenerations,
	})

	var data lineageResourceModel
	var diags diag.Diagnostics
	if err != nil {
		// The rename succeeded even though the refresh failed: fall back to RenameLineage's own response, which
		// carries the renamed name/modified fields, and keep the previously known generation pointers, which a
		// rename cannot affect.
		resp.Diagnostics.AddWarning("Lineage Renamed, But Refresh Failed",
			fmt.Sprintf("Lineage %d was renamed successfully, but re-fetching generation detail failed: %s. Run terraform plan/apply again to refresh it.", lineageID, err.Error()))

		data, diags = mapLineageToResourceModel(ctx, (*cloudcertificates.Lineage)(renamed))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Head, data.CurrentProduction, data.PreviousProduction, data.CurrentStaging =
			state.Head, state.CurrentProduction, state.PreviousProduction, state.CurrentStaging
		data.SigningTarget = state.SigningTarget
	} else {
		data, diags = mapLineageToResourceModel(ctx, (*cloudcertificates.Lineage)(lineage))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes the lineage.
func (r *lineageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Resource Delete")

	var state lineageResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	lineageID := state.LineageID.ValueInt64()
	err := r.Client.GetCloudCertificates().DeleteLineage(ctx, cloudcertificates.DeleteLineageRequest{LineageID: lineageID})
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		// Already gone; nothing left to do.
	case errors.Is(err, cloudcertificates.ErrLineageHasActiveProduction):
		resp.Diagnostics.AddError("Lineage Has Active Production",
			fmt.Sprintf("Lineage %d has an active production certificate. Deactivate before deleting.", lineageID))
	case errors.Is(err, cloudcertificates.ErrLineageHasActiveStaging):
		resp.Diagnostics.AddError("Lineage Has Active Staging",
			fmt.Sprintf("Lineage %d has an active staging certificate. Deactivate before deleting.", lineageID))
	default:
		resp.Diagnostics.AddError("Failed to delete lineage", err.Error())
	}
}

// ImportState imports a lineage by its numeric lineage_id.
func (r *lineageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Cloud Certificates Lineage Resource")

	parts, err := text.ImportIDSplitter("lineageID").AcceptLen(1).Split(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Incorrect import ID", err.Error())
		return
	}

	lineageID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("expected a numeric lineage_id, got: %q", req.ID))
		return
	}
	if lineageID <= 0 {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("lineage_id must be greater than 0, got: %d", lineageID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("lineage_id"), lineageID)...)
}

// ModifyPlan checks for a lineage name conflict ahead of apply, so the user gets an early diagnostic instead of
// discovering it only when CreateLineage/RenameLineage fails. This is a best-effort check only: the name could
// still be taken (or freed) between plan and apply, so CreateLineage/RenameLineage's own error handling remains
// the authoritative guard.
func (r *lineageResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if modifiers.IsDelete(req) {
		return // destroy plan
	}

	var plan lineageResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if !tf.IsKnown(plan.LineageName) {
		return
	}

	if modifiers.IsUpdate(req) {
		var state lineageResourceModel
		if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
			return
		}
		if state.LineageName.ValueString() == plan.LineageName.ValueString() {
			return // name unchanged, nothing to check
		}
	}

	name := plan.LineageName.ValueString()
	existing, err := r.Client.GetCloudCertificates().ListLineages(ctx, cloudcertificates.ListLineagesRequest{
		LineageName: name,
		PageSize:    cloudcertificates.MaxListLineagesPageSize,
	})
	if err != nil {
		resp.Diagnostics.AddWarning("Lineage Name Conflict Check Skipped",
			fmt.Sprintf("Could not check whether the name %q is already in use: %s. "+
				"Any conflict will still be caught when the change is applied.", name, err.Error()))
		return
	}
	for _, l := range existing.Lineages {
		// Uniqueness is enforced server-side with an exact-case match (ListLineages' substring filter is only
		// for searching), so an exact match here can never be this lineage's own pre-rename name.
		if l.LineageName == name {
			resp.Diagnostics.AddAttributeError(path.Root("lineage_name"),
				"Lineage Name Conflict", fmt.Sprintf("A lineage named %q already exists for this account.", name))
			return
		}
	}
}

var allExpandGenerations = []cloudcertificates.ExpandGenerations{
	cloudcertificates.ExpandGenerationsHead,
	cloudcertificates.ExpandGenerationsCurrentProduction,
	cloudcertificates.ExpandGenerationsCurrentStaging,
	cloudcertificates.ExpandGenerationsPreviousProduction,
}

// planToCreateRequest extracts a CreateLineageRequest from the resource plan.
func planToCreateRequest(ctx context.Context, plan lineageResourceModel) (cloudcertificates.CreateLineageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var keySpecs map[string]string
	diags.Append(plan.KeySpecs.ElementsAs(ctx, &keySpecs, false)...)

	var sans []string
	diags.Append(plan.SANs.ElementsAs(ctx, &sans, false)...)
	// sans is a set: sort before sending so the request is deterministic instead of depending on the set's
	// internal iteration order (the API doesn't preserve order either way).
	sort.Strings(sans)

	body := cloudcertificates.CreateLineageRequestBody{
		ContractID:    plan.ContractID.ValueString(),
		GroupID:       plan.GroupID.ValueInt64(),
		SecureNetwork: cloudcertificates.SecureNetwork(plan.SecureNetwork.ValueString()),
		SANs:          sans,
	}
	// key_specs is a map: sort its keys before sending so the request is deterministic instead of depending on
	// Go's randomized map iteration order.
	keyTypes := make([]string, 0, len(keySpecs))
	for keyType := range keySpecs {
		keyTypes = append(keyTypes, keyType)
	}
	sort.Strings(keyTypes)
	for _, keyType := range keyTypes {
		body.KeySpecs = append(body.KeySpecs, cloudcertificates.KeySpec{
			KeyType: cloudcertificates.CryptographicAlgorithm(keyType),
			KeySize: cloudcertificates.KeySize(keySpecs[keyType]),
		})
	}
	if tf.IsKnown(plan.GeoClass) {
		body.GeoClass = cloudcertificates.GeoClass(plan.GeoClass.ValueString())
	}
	if tf.IsKnown(plan.LineageName) {
		body.LineageName = plan.LineageName.ValueString()
	}
	if tf.IsKnown(plan.LineageType) {
		body.LineageType = cloudcertificates.LineageType(plan.LineageType.ValueString())
	}
	if tf.IsKnown(plan.Subject) && !plan.Subject.IsNull() {
		var subj lineageSubjectModel
		diags.Append(plan.Subject.As(ctx, &subj, basetypes.ObjectAsOptions{})...)
		body.Subject = &cloudcertificates.Subject{
			CommonName:         subj.CommonName.ValueString(),
			Organization:       subj.Organization.ValueString(),
			OrganizationalUnit: subj.OrganizationalUnit.ValueString(),
			Country:            subj.Country.ValueString(),
			State:              subj.State.ValueString(),
			Locality:           subj.Locality.ValueString(),
		}
	}

	return cloudcertificates.CreateLineageRequest{Body: body}, diags
}

// mapLineageToResourceModel converts an API lineage response into the resource model.
func mapLineageToResourceModel(ctx context.Context, lineage *cloudcertificates.Lineage) (lineageResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := lineageResourceModel{
		LineageID:           types.Int64Value(lineage.LineageID),
		AccountID:           types.StringValue(lineage.AccountID),
		ContractID:          types.StringValue(lineage.ContractID),
		GroupID:             types.Int64Value(lineage.GroupID),
		GeoClass:            types.StringValue(lineage.GeoClass),
		LineageName:         types.StringValue(lineage.LineageName),
		LineageType:         types.StringValue(lineage.LineageType),
		SecureNetwork:       types.StringValue(lineage.SecureNetwork),
		StackMode:           types.StringValue(lineage.StackMode),
		LineageCreatedBy:    types.StringValue(lineage.LineageCreatedBy),
		LineageCreatedTime:  fwdate.TimeRFC3339Value(lineage.LineageCreatedTime),
		LineageModifiedBy:   types.StringValue(lineage.LineageModifiedBy),
		LineageModifiedTime: fwdate.TimeRFC3339Value(lineage.LineageModifiedTime),
	}

	keySpecs := make(map[string]string, len(lineage.KeySpecs))
	for _, ks := range lineage.KeySpecs {
		keySpecs[ks.KeyType] = ks.KeySize
	}
	var d diag.Diagnostics
	data.KeySpecs, d = types.MapValueFrom(ctx, types.StringType, keySpecs)
	diags.Append(d...)

	data.SANs, d = types.SetValueFrom(ctx, types.StringType, lineage.SANs)
	diags.Append(d...)

	// mirrors Head/CurrentProduction/etc.: collapse to a null object instead of one with every field null, so
	// an API response with no subject at all is indistinguishable from one that was never set.
	if lineage.Subject == (cloudcertificates.Subject{}) {
		data.Subject = types.ObjectNull(subjectObjectAttrTypes)
	} else {
		subject := lineageSubjectModel{
			CommonName:         tf.StringValueOrNullIfEmpty(lineage.Subject.CommonName),
			Organization:       tf.StringValueOrNullIfEmpty(lineage.Subject.Organization),
			OrganizationalUnit: tf.StringValueOrNullIfEmpty(lineage.Subject.OrganizationalUnit),
			Country:            tf.StringValueOrNullIfEmpty(lineage.Subject.Country),
			State:              tf.StringValueOrNullIfEmpty(lineage.Subject.State),
			Locality:           tf.StringValueOrNullIfEmpty(lineage.Subject.Locality),
		}
		data.Subject, d = types.ObjectValueFrom(ctx, subjectObjectAttrTypes, subject)
		diags.Append(d...)
	}

	if lineage.Head != nil {
		data.Head, d = types.ObjectValueFrom(ctx, generationPointerObjectAttrTypes, mapGenerationPointer(lineage.Head))
		diags.Append(d...)
	} else {
		data.Head = types.ObjectNull(generationPointerObjectAttrTypes)
	}

	if lineage.CurrentProduction != nil {
		data.CurrentProduction, d = types.ObjectValueFrom(ctx, generationPointerObjectAttrTypes, mapGenerationPointer(lineage.CurrentProduction))
		diags.Append(d...)
	} else {
		data.CurrentProduction = types.ObjectNull(generationPointerObjectAttrTypes)
	}

	if lineage.PreviousProduction != nil {
		data.PreviousProduction, d = types.ObjectValueFrom(ctx, generationPointerObjectAttrTypes, mapGenerationPointer(lineage.PreviousProduction))
		diags.Append(d...)
	} else {
		data.PreviousProduction = types.ObjectNull(generationPointerObjectAttrTypes)
	}

	if lineage.CurrentStaging != nil {
		data.CurrentStaging, d = types.ObjectValueFrom(ctx, generationPointerObjectAttrTypes, mapGenerationPointer(lineage.CurrentStaging))
		diags.Append(d...)
	} else {
		data.CurrentStaging = types.ObjectNull(generationPointerObjectAttrTypes)
	}

	switch signingTarget(*lineage) {
	case lineage.Head:
		data.SigningTarget = data.Head
	case lineage.CurrentProduction:
		data.SigningTarget = data.CurrentProduction
	default:
		data.SigningTarget = types.ObjectNull(generationPointerObjectAttrTypes)
	}

	return data, diags
}
