package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/retry"
	"github.com/akamai/terraform-provider-akamai/v11/internal/text"
	fwdate "github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &uploadResource{}
	_ resource.ResourceWithConfigure   = &uploadResource{}
	_ resource.ResourceWithImportState = &uploadResource{}
	_ resource.ResourceWithModifyPlan  = &uploadResource{}
)

const (
	defaultUploadPollInterval = 5 * time.Second
	defaultUploadPollTimeout  = 2 * time.Minute
)

type (
	uploadResource struct {
		meta.Resource
		uploadResourceConfig
	}

	uploadResourceConfig struct {
		pollInterval time.Duration
		pollTimeout  time.Duration
	}

	uploadResourceModel struct {
		LineageID                   types.Int64    `tfsdk:"lineage_id"`
		GenerationID                types.Int64    `tfsdk:"generation_id"`
		AcknowledgeWarnings         types.Bool     `tfsdk:"acknowledge_warnings"`
		GenerationStatus            types.String   `tfsdk:"generation_status"`
		FirstPromotedToProductionAt types.String   `tfsdk:"first_promoted_to_production_at"`
		GenerationCreatedBy         types.String   `tfsdk:"generation_created_by"`
		GenerationCreatedTime       types.String   `tfsdk:"generation_created_time"`
		GenerationModifiedBy        types.String   `tfsdk:"generation_modified_by"`
		GenerationModifiedTime      types.String   `tfsdk:"generation_modified_time"`
		Algorithms                  types.Map      `tfsdk:"algorithms"`
		Timeouts                    timeouts.Value `tfsdk:"timeouts"`
	}
)

// NewUploadResource returns a new Cloud Certificates lineage upload resource.
func NewUploadResource(config uploadResourceConfig) func() resource.Resource {
	return func() resource.Resource {
		return &uploadResource{uploadResourceConfig: config}
	}
}

// defaultUploadResourceConfig returns the production configuration for the lineage upload resource.
func defaultUploadResourceConfig() uploadResourceConfig {
	return uploadResourceConfig{pollInterval: defaultUploadPollInterval, pollTimeout: defaultUploadPollTimeout}
}

// Metadata configures the resource's type name.
func (r *uploadResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_upload"
}

// Schema defines the Terraform schema for the lineage upload resource.
func (r *uploadResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Uploads signed certificates to a certificate lineage generation. Create binds the resource to " +
			"the generation to which the certificates belong: either the lineage's head generation or, for a " +
			"MULTIPLE_STACK lineage already fully promoted to production with one algorithm still missing, the new head " +
			"generation that CompleteLineage creates with the completed algorithm set. The resource never moves to a " +
			"different generation: Update only adds a missing algorithm to that same generation, and any change " +
			"that would require binding to a different generation replaces the resource.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the lineage to upload the certificate to.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"generation_id": schema.Int64Attribute{
				Computed: true,
				Description: "Unique identifier of the generation this resource is bound to. Set to the lineage's head " +
					"generation when created or to the new head generation created by a complete lineage operation when " +
					"adding a missing algorithm to a fully promoted MULTIPLE_STACK lineage; never changes for the " +
					"lifetime of this resource instance: if the lineage's head generation moves on to a different " +
					"generation, this resource is destroyed and recreated to bind to it.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"algorithms": schema.MapNestedAttribute{
				Required:    true,
				Description: "Signed certificate material to upload, keyed by key_type (at most one RSA and one ECDSA entry).",
				Validators: []validator.Map{
					mapvalidator.SizeBetween(1, 2),
					mapvalidator.KeysAre(stringvalidator.OneOf(string(cloudcertificates.CryptographicAlgorithmRSA), string(cloudcertificates.CryptographicAlgorithmECDSA))),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: mergeResourceAttributes(algorithmResourceAttributes(), map[string]schema.Attribute{
						"signed_certificate_pem": schema.StringAttribute{
							Required:    true,
							Description: "PEM-encoded signed certificate to upload for this key type.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.RegexMatches(validators.CertificatePEMRegex,
									fmt.Sprintf("must be in PEM format: '%s'", validators.CertificatePEMRegex)),
							},
						},
						"trust_chain_pem": schema.StringAttribute{
							Optional:    true,
							Description: "Optional PEM-encoded trust chain to upload alongside the signed certificate.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(validators.ToolchainPEMRegex,
									fmt.Sprintf("must be in PEM format: '%s'", validators.ToolchainPEMRegex)),
							},
						},
					}),
				},
			},
			"generation_status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the generation.",
			},
			"first_promoted_to_production_at": schema.StringAttribute{
				Computed:    true,
				Description: "Time the generation was first promoted to the production network, in RFC3339 format. Null if never promoted.",
			},
			"generation_created_by": schema.StringAttribute{
				Computed:    true,
				Description: "Username of the person who created the generation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"generation_created_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time the generation was created, in RFC3339 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"generation_modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "Username of the person who last modified the generation.",
			},
			"generation_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time the generation was last modified, in RFC3339 format.",
			},
			"acknowledge_warnings": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Whether to automatically acknowledge non-fatal warnings when uploading or completing a certificate - " +
					"whichever operation this resource actually performs, depending on whether the lineage already has a head " +
					"generation. It's a plain request flag, not a durable property of an already-uploaded certificate, so changing " +
					"it never forces a replace or re-upload. Defaults to `false`.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "Optional configurable timeout for polling until the uploaded certificate(s) finish processing server-side. By default it's 2m.",
				Update:            true,
				UpdateDescription: "Optional configurable timeout for polling until the uploaded certificate(s) finish processing server-side. By default it's 2m.",
			}),
		},
	}
}

// Create resolves which concrete generation this resource instance is bound to, then uploads to it. That's the
// lineage's head generation if one exists, whether empty or already carrying one algorithm. Otherwise it's the
// new head generation CompleteLineage creates with the completed algorithm set, which is only possible once the
// lineage has no head left. Once bound, the resource never moves to another generation. Update only ever adds
// algorithms in place. ModifyPlan forces a replace, a fresh Create, for any change that would otherwise require
// binding to a different generation.
func (r *uploadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Upload Resource Create")

	var plan uploadResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var planAlgorithms map[string]algorithmModel
	if resp.Diagnostics.Append(plan.Algorithms.ElementsAs(ctx, &planAlgorithms, false)...); resp.Diagnostics.HasError() {
		return
	}

	lineageID := plan.LineageID.ValueInt64()
	client := r.Client.GetCloudCertificates()

	lineage, err := client.GetLineage(ctx, cloudcertificates.GetLineageRequest{
		LineageID: lineageID,
		ExpandGenerations: []cloudcertificates.ExpandGenerations{
			cloudcertificates.ExpandGenerationsHead,
			cloudcertificates.ExpandGenerationsCurrentProduction,
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to look up lineage head generation", err.Error())
		return
	}

	var generationID int64
	var uploadedAlgorithms map[string]algorithmModel
	switch {
	case lineage.Head != nil:
		// A head generation targeted by a fresh resource must be empty; if it already has a signed certificate, it's someone else's.
		generationID = lineage.Head.HeadGenerationID
		headGeneration, err := client.GetGeneration(ctx, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: generationID})
		if err != nil {
			resp.Diagnostics.AddError("Failed to read head generation", err.Error())
			return
		}
		for _, a := range headGeneration.Algorithms {
			if a.SignedCertificatePEM == nil {
				continue
			}
			resp.Diagnostics.AddError("Head Generation Already Has A Signed Certificate",
				fmt.Sprintf("Generation %d of lineage %d already has an accepted signed certificate for the %s algorithm. "+
					"Import the existing upload instead: terraform import <resource address> %d,%d", generationID, lineageID, a.KeyType, lineageID, generationID))
			return
		}
		_, err = client.UploadSignedCertificate(ctx, cloudcertificates.UploadSignedCertificateRequest{
			LineageID:           lineageID,
			GenerationID:        generationID,
			AcknowledgeWarnings: plan.AcknowledgeWarnings.ValueBool(),
			Body:                cloudcertificates.UploadSignedCertificateRequestBody{Algorithms: algorithmsRequestBody(planAlgorithms)},
		})
		if err != nil {
			resp.Diagnostics.Append(uploadCertificateErrorDiagnostics(err, lineageID, generationID)...)
			return
		}
		uploadedAlgorithms = planAlgorithms
	case lineage.CurrentProduction != nil:
		// No head means the lineage is fully promoted. CompleteLineage leaves production untouched and creates a
		// new head generation with the full, signed algorithm set, and the resource binds to that new head. An
		// already-signed algorithm on production that isn't configured here means the configuration is out of
		// sync, so the resource returns an error and suggests importing it.
		generationID = lineage.CurrentProduction.ProductionGenerationID
		currentGeneration, err := client.GetGeneration(ctx, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: generationID})
		if err != nil {
			resp.Diagnostics.AddError("Failed to read current production generation", err.Error())
			return
		}
		diff := diffAlgorithms(planAlgorithms, mapLineageAlgorithms(currentGeneration.Algorithms))
		for keyType := range diff.Modified {
			resp.Diagnostics.AddError("Algorithm Change Not Supported",
				fmt.Sprintf("the %s algorithm instance already has an accepted signed certificate; it cannot be replaced or have its trust chain changed", keyType))
			return
		}
		for keyType := range diff.Removed {
			resp.Diagnostics.AddError("Current Production Generation Already Has A Signed Certificate",
				fmt.Sprintf("Generation %d of lineage %d already has an accepted signed certificate for the %s algorithm, which isn't configured here. "+
					"Import the existing upload instead: terraform import <resource address> %d,%d", generationID, lineageID, keyType, lineageID, generationID))
			return
		}
		completed, err := client.CompleteLineage(ctx, cloudcertificates.CompleteLineageRequest{
			LineageID:           lineageID,
			AcknowledgeWarnings: plan.AcknowledgeWarnings.ValueBool(),
			Body: cloudcertificates.CompleteLineageRequestBody{
				SourceGenerationID: ptr.To(generationID),
				Algorithms:         algorithmsRequestBody(diff.New),
			},
		})
		if err != nil {
			resp.Diagnostics.Append(completeLineageErrorDiagnostics(err, lineageID)...)
			return
		}
		if completed.Head == nil {
			resp.Diagnostics.AddError("Failed to complete lineage",
				fmt.Sprintf("Lineage %d has no head generation after completion.", lineageID))
			return
		}
		generationID = completed.Head.HeadGenerationID
		uploadedAlgorithms = diff.New
	default:
		resp.Diagnostics.AddError("No Head Or Production Generation",
			fmt.Sprintf("Lineage %d has no head or current production generation to upload a certificate to.", lineageID))
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pollTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	generation, err := r.waitForSignedAlgorithms(ctx, client, lineageID, generationID, uploadedAlgorithms, timeout)
	if err != nil {
		// The certificate is already uploaded or completed on the server at this point. Advise an import as the
		// recovery mechanism.
		resp.Diagnostics.AddError("Failed to read generation after upload",
			fmt.Sprintf("The certificate was uploaded to generation %d of lineage %d, but reading it back failed: %s. "+
				"To recover, import the newly created upload resource: terraform import <resource address> %d,%d",
				generationID, lineageID, err.Error(), lineageID, generationID))
		return
	}

	data, diags := mapGenerationToUploadModel(ctx, lineageID, generation, planAlgorithms)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AcknowledgeWarnings = plan.AcknowledgeWarnings
	data.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ModifyPlan replaces the resource whenever its head generation has moved on to a different generation. A new
// head is always empty, so it's safe to rebind to without question. If instead the lineage no longer recognizes
// the bound generation as either its head or its current production generation, ModifyPlan errors instead of
// replacing. This typically happens after an out-of-band rollback repoints production straight to a previous
// generation, skipping any intermediate head state this resource could have caught. Unlike a head simply moving
// on, which is always a fresh, empty canvas safe to rebind to, a generation reached this way may already carry
// real signed content unrelated to this resource instance, so silently binding a new instance to it is never
// safe. Write-once violations, see diffAlgorithms, are only checked once the bound generation is confirmed
// still current. A promoted generation gaining a new algorithm also forces a replace, since only Create or
// CompleteLineage can do that.
func (r *uploadResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if modifiers.IsCreate(req) || modifiers.IsDelete(req) {
		// create or destroy: nothing to compare against.
		return
	}

	var plan, state uploadResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	if plan.Algorithms.IsUnknown() {
		return
	}

	var planAlgorithms, stateAlgorithms map[string]algorithmModel
	resp.Diagnostics.Append(plan.Algorithms.ElementsAs(ctx, &planAlgorithms, false)...)
	resp.Diagnostics.Append(state.Algorithms.ElementsAs(ctx, &stateAlgorithms, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	diff := diffAlgorithms(planAlgorithms, stateAlgorithms)

	lineageID := state.LineageID.ValueInt64()
	generationID := state.GenerationID.ValueInt64()
	lineage, err := r.Client.GetCloudCertificates().GetLineage(ctx, cloudcertificates.GetLineageRequest{
		LineageID: lineageID,
		ExpandGenerations: []cloudcertificates.ExpandGenerations{
			cloudcertificates.ExpandGenerationsHead,
			cloudcertificates.ExpandGenerationsCurrentProduction,
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed To Check Lineage Head Generation",
			fmt.Sprintf("Could not determine whether lineage %d's head generation has moved on: %s.", lineageID, err.Error()))
		return
	} else if needsReplace, replaceErr := generationNeedsReplace(lineage, generationID, len(diff.New) > 0); replaceErr != nil {
		resp.Diagnostics.AddError("Generation No Longer Current",
			fmt.Sprintf("Generation %d of lineage %d is no longer the lineage's head or current production generation. "+
				"This can happen after an out-of-band rollback moved production onto a different generation. Import the "+
				"generation that's now current instead: terraform import <resource address> %d,<current generation ID>", generationID, lineageID, lineageID))
		return
	} else if needsReplace {
		// mark generation_id unknown and force a replace: the resource can no longer stay bound to it.
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("generation_id"), types.Int64Unknown())...)
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("generation_id"))
		return
	}

	for keyType := range diff.Modified {
		resp.Diagnostics.AddError("Algorithm Change Not Supported",
			fmt.Sprintf("The %s algorithm instance already has an accepted signed certificate; it cannot be replaced or have its trust chain changed.", keyType))
		return
	}
	for keyType := range diff.Removed {
		resp.Diagnostics.AddError("Algorithm Change Not Supported",
			fmt.Sprintf("Removing the %s algorithm from an existing upload is not supported.", keyType))
		return
	}

	// An already-signed algorithm untouched by this plan keeps its exact prior state. Without this, adding or
	// removing a sibling key in the algorithms map makes the framework replan every other key's Computed
	// sub-attributes as unknown, purely because the map itself changed shape. A UseStateForUnknown plan modifier
	// on the sub-attributes can't fix this. A key that's brand new to the map has no prior state to fall back
	// to, so the modifier would resolve straight to null instead of staying unknown.
	for keyType, existing := range stateAlgorithms {
		if existing.SignedCertificatePEM.IsNull() {
			continue
		}
		if _, modified := diff.Modified[keyType]; modified {
			continue
		}
		if _, stillPlanned := planAlgorithms[keyType]; !stillPlanned {
			continue
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("algorithms").AtMapKey(keyType), existing)...)
	}
}

// errGenerationUnrecognized indicates the lineage no longer recognizes the bound generation as either its head
// or its current production generation. This typically happens after an out-of-band rollback repoints
// production straight to a previous generation, skipping any intermediate head state this resource could have
// caught. Unlike a head simply moving on, which is always a fresh, empty canvas safe to rebind to, a generation
// reached this way may already carry real signed content unrelated to this resource instance. It's never safe
// to silently replace onto it. The caller must ask for an explicit import instead.
var errGenerationUnrecognized = errors.New("generation is no longer the lineage's head or current production generation")

// generationNeedsReplace reports whether an upload resource bound to generationID must be replaced to stay
// bound to the lineage's live state. This is true when the head has moved on to a different generation, which
// is always empty and safe to rebind to, or when a still-current, no-head production generation is gaining a
// new algorithm, which is only possible via a fresh Create or CompleteLineage. It returns
// errGenerationUnrecognized if generationID is no longer the lineage's head or current production generation at
// all.
func generationNeedsReplace(lineage *cloudcertificates.GetLineageResponse, generationID int64, hasNewAlgorithm bool) (bool, error) {
	switch {
	case lineage.Head != nil:
		return lineage.Head.HeadGenerationID != generationID, nil
	case lineage.CurrentProduction != nil && lineage.CurrentProduction.ProductionGenerationID == generationID:
		return hasNewAlgorithm, nil
	default:
		return false, errGenerationUnrecognized
	}
}

// Read refreshes the upload state from the generation's current metadata. If the generation is gone, it delegates
// to handleGenerationGone rather than assuming the resource is safe to drop.
func (r *uploadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Upload Resource Read")

	var state uploadResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var stateAlgorithms map[string]algorithmModel
	if resp.Diagnostics.Append(state.Algorithms.ElementsAs(ctx, &stateAlgorithms, false)...); resp.Diagnostics.HasError() {
		return
	}

	lineageID := state.LineageID.ValueInt64()
	generationID := state.GenerationID.ValueInt64()

	generation, err := r.Client.GetCloudCertificates().GetGeneration(ctx, cloudcertificates.GetGenerationRequest{
		LineageID:    lineageID,
		GenerationID: generationID,
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddWarning("Lineage Not Found",
				fmt.Sprintf("Removing upload for generation %d of lineage %d since the lineage is gone: %s.",
					generationID, lineageID, err.Error()))
			resp.State.RemoveResource(ctx)
			return
		}
		if errors.Is(err, cloudcertificates.ErrGenerationNotFound) {
			r.handleGenerationGone(ctx, lineageID, generationID, "it no longer exists", resp)
			return
		}
		resp.Diagnostics.AddError("Failed to read generation", err.Error())
		return
	}

	if generation.GenerationStatus == string(cloudcertificates.GenerationStatusArchived) ||
		generation.GenerationStatus == string(cloudcertificates.GenerationStatusAbandoned) {
		r.handleGenerationGone(ctx, lineageID, generationID, fmt.Sprintf("the generation has been %s", generation.GenerationStatus), resp)
		return
	}

	data, diags := mapGenerationToUploadModel(ctx, lineageID, generation, stateAlgorithms)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AcknowledgeWarnings = state.AcknowledgeWarnings
	data.Timeouts = state.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *uploadResource) handleGenerationGone(ctx context.Context, lineageID, generationID int64, reason string, resp *resource.ReadResponse) {
	lineage, err := r.Client.GetCloudCertificates().GetLineage(ctx, cloudcertificates.GetLineageRequest{
		LineageID: lineageID,
		ExpandGenerations: []cloudcertificates.ExpandGenerations{
			cloudcertificates.ExpandGenerationsHead,
			cloudcertificates.ExpandGenerationsCurrentProduction,
		},
	})
	if err != nil && !errors.Is(err, cloudcertificates.ErrLineageNotFound) {
		resp.Diagnostics.AddWarning("Failed To Check Lineage For Import Target",
			fmt.Sprintf("Could not determine whether lineage %d has a current generation to import instead: %s.", lineageID, err.Error()))
	}

	// A still-existing lineage always has a head or current production generation, since the API refuses to
	// delete the last one. So reaching neither here, as opposed to the lineage itself being gone, should never
	// happen.
	var currentID int64
	switch {
	case err == nil && lineage.Head != nil:
		currentID = lineage.Head.HeadGenerationID
	case err == nil && lineage.CurrentProduction != nil:
		currentID = lineage.CurrentProduction.ProductionGenerationID
	default:
		resp.Diagnostics.AddWarning("Generation Not Found",
			fmt.Sprintf("Removing upload for lineage %d, generation %d from state; %s.", lineageID, generationID, reason))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.AddWarning("Generation No Longer Current",
		fmt.Sprintf("Removing upload for generation %d of lineage %d since the generation is gone: %s. "+
			"Lineage %d now has generation %d as its head or current "+
			"production generation - import that instead: terraform import <resource address> %d,%d",
			generationID, lineageID, reason, lineageID, currentID, lineageID, currentID))
	resp.State.RemoveResource(ctx)
}

// Update only ever adds a newly-configured algorithm to the same, still-draft generation this resource is
// already bound to. It never moves the resource to another generation. Modifying or removing an algorithm that
// already has an accepted signed certificate is never supported. This write-once rule applies regardless of
// promotion status, see diffAlgorithms. A generation that's since been promoted to production never reaches
// Update at all: ModifyPlan catches that case and replaces the resource instead, see ModifyPlan and Create.
func (r *uploadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Upload Resource Update")

	var plan, state uploadResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var planAlgorithms, stateAlgorithms map[string]algorithmModel
	resp.Diagnostics.Append(plan.Algorithms.ElementsAs(ctx, &planAlgorithms, false)...)
	resp.Diagnostics.Append(state.Algorithms.ElementsAs(ctx, &stateAlgorithms, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lineageID := state.LineageID.ValueInt64()
	generationID := state.GenerationID.ValueInt64()
	client := r.Client.GetCloudCertificates()

	diff := diffAlgorithms(planAlgorithms, stateAlgorithms)
	for keyType := range diff.Modified {
		resp.Diagnostics.AddError("Algorithm Change Not Supported",
			fmt.Sprintf("The %s algorithm instance already has an accepted signed certificate; it cannot be replaced or have its trust chain changed.", keyType))
		return
	}
	for keyType := range diff.Removed {
		resp.Diagnostics.AddError("Algorithm Change Not Supported",
			fmt.Sprintf("Removing the %s algorithm from an existing upload is not supported.", keyType))
		return
	}

	var generation *cloudcertificates.GetGenerationResponse
	if len(diff.New) > 0 {
		_, err := client.UploadSignedCertificate(ctx, cloudcertificates.UploadSignedCertificateRequest{
			LineageID:           lineageID,
			GenerationID:        generationID,
			AcknowledgeWarnings: plan.AcknowledgeWarnings.ValueBool(),
			Body:                cloudcertificates.UploadSignedCertificateRequestBody{Algorithms: algorithmsRequestBody(diff.New)},
		})
		if err != nil {
			resp.Diagnostics.Append(uploadCertificateErrorDiagnostics(err, lineageID, generationID)...)
			return
		}

		timeout, diags := plan.Timeouts.Update(ctx, r.pollTimeout)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		generation, err = r.waitForSignedAlgorithms(ctx, client, lineageID, generationID, diff.New, timeout)
		if err != nil {
			resp.Diagnostics.AddError("Failed to read generation after update",
				fmt.Sprintf("The certificate was uploaded to generation %d of lineage %d, but reading it back failed: %s. "+
					"To recover, import the upload resource: terraform import <resource address> %d,%d",
					generationID, lineageID, err.Error(), lineageID, generationID))
			return
		}
	} else {
		resp.Diagnostics.AddWarning("Local Update Only",
			"Only acknowledge_warnings or timeouts has changed, no API call will be made.")

		var err error
		generation, err = client.GetGeneration(ctx, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: generationID})
		if err != nil {
			resp.Diagnostics.AddError("Failed to read generation after update", err.Error())
			return
		}
	}

	data, diags := mapGenerationToUploadModel(ctx, lineageID, generation, planAlgorithms)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AcknowledgeWarnings = plan.AcknowledgeWarnings
	data.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// waitForSignedAlgorithms polls GetGeneration until every key type in uploaded shows its signed-certificate
// details, or until timeout elapses. A GET immediately following a successful upload or complete can still
// briefly report CERT_UPLOAD_PROCESSING, or omit signed-certificate fields, while the upload finishes
// processing server-side.
func (r *uploadResource) waitForSignedAlgorithms(ctx context.Context, client cloudcertificates.CloudCertificates, lineageID, generationID int64, uploaded map[string]algorithmModel, timeout time.Duration) (*cloudcertificates.GetGenerationResponse, error) {
	return retry.Poll(ctx, retry.PollingOpts[cloudcertificates.GetGenerationResponse]{
		Fn: func(ctx context.Context) (*cloudcertificates.GetGenerationResponse, error) {
			return client.GetGeneration(ctx, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: generationID})
		},
		ShouldRetryData: func(generation cloudcertificates.GetGenerationResponse) bool {
			return !allAlgorithmsSigned(generation.Algorithms, uploaded)
		},
		Interval: r.pollInterval,
		Deadline: timeout,
	})
}

// allAlgorithmsSigned reports whether every key type in uploaded has a matching entry in algorithms that already
// exposes its signed-certificate details.
func allAlgorithmsSigned(algorithms []cloudcertificates.Algorithm, uploaded map[string]algorithmModel) bool {
	for keyType := range uploaded {
		signed := false
		for _, a := range algorithms {
			if a.KeyType == keyType && a.SignedCertificatePEM != nil {
				signed = true
				break
			}
		}
		if !signed {
			return false
		}
	}
	return true
}

// Delete removes the resource from Terraform state only. Certificate uploads are integral to the generation, which
// is managed by the lineage's own lifecycle (inactivity timers or explicit lineage deletion).
func (r *uploadResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Upload Resource Delete")
}

// ImportState imports an upload by its lineage_id and generation_id, plus an optional trailing
// acknowledge_warnings. Config can't be reached from here, and the value isn't part of the API's response, so
// without it the attribute would just fall back to its schema default.
func (r *uploadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := text.ImportIDSplitter("lineageID,generationID[,acknowledgeWarnings]").AcceptLen(2).AcceptLen(3).Split(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Incorrect import ID", err.Error())
		return
	}

	lineageID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("expected a numeric lineage_id, got: %q", parts[0]))
		return
	}
	if lineageID <= 0 {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("lineage_id must be greater than 0, got: %d", lineageID))
		return
	}

	generationID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("expected a numeric generation_id, got: %q", parts[1]))
		return
	}
	if generationID <= 0 {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("generation_id must be greater than 0, got: %d", generationID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("lineage_id"), lineageID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("generation_id"), generationID)...)

	if len(parts) == 3 {
		acknowledgeWarnings, err := strconv.ParseBool(parts[2])
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("expected a boolean acknowledge_warnings, got: %q", parts[2]))
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("acknowledge_warnings"), acknowledgeWarnings)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("acknowledge_warnings"), false)...)
	}
}

// algorithmDiff categorizes planAlgorithms against a set of existing algorithms, purely by key type and
// content. New covers keys present only in plan, or present in existing but not yet signed. Modified covers
// keys already signed in both but with different content. Removed covers keys already signed in existing only.
// It makes no judgment on whether a given modification or removal is actually allowed. That's a business rule
// left to each caller.
type algorithmDiff struct {
	New      map[string]algorithmModel
	Modified map[string]algorithmModel
	Removed  map[string]algorithmModel
}

// diffAlgorithms computes the algorithmDiff of planAlgorithms against existing. existing needn't be
// pre-filtered to already-signed algorithms. An unsigned entry, CSR_READY with a null SignedCertificatePEM, and
// a missing key, whose zero-value algorithmModel is likewise null, are both treated as "nothing to compare
// against".
func diffAlgorithms(planAlgorithms, existing map[string]algorithmModel) algorithmDiff {
	diff := algorithmDiff{
		New:      make(map[string]algorithmModel),
		Modified: make(map[string]algorithmModel),
		Removed:  make(map[string]algorithmModel),
	}
	for keyType, a := range planAlgorithms {
		other, exists := existing[keyType]
		switch {
		case !exists || other.SignedCertificatePEM.IsNull():
			diff.New[keyType] = a
		case a.SignedCertificatePEM.IsUnknown() || a.TrustChainPEM.IsUnknown():
			// Defer content comparisons for unknown nested values during planning
		case a.SignedCertificatePEM.ValueString() != other.SignedCertificatePEM.ValueString() ||
			a.TrustChainPEM.ValueString() != other.TrustChainPEM.ValueString():

			diff.Modified[keyType] = a
		}
	}
	for keyType, a := range existing {
		if a.SignedCertificatePEM.IsNull() {
			continue
		}
		if _, exists := planAlgorithms[keyType]; !exists {
			diff.Removed[keyType] = a
		}
	}
	return diff
}

// algorithmsRequestBody builds the wire request body from the resource's algorithms map.
func algorithmsRequestBody(algorithms map[string]algorithmModel) map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate {
	body := make(map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate, len(algorithms))
	for keyType, a := range algorithms {
		body[cloudcertificates.CryptographicAlgorithm(keyType)] = cloudcertificates.SignedCertificate{
			SignedCertificatePEM: a.SignedCertificatePEM.ValueString(),
			TrustChainPEM:        a.TrustChainPEM.ValueString(),
		}
	}
	return body
}

// mapGenerationToUploadModel maps a generation's signed algorithms into the upload resource model. When
// trackedAlgorithms is non-nil, only the resource's own configured or planned key types are included.
// signed_certificate_pem is Required, so an algorithm this instance doesn't configure can't be represented.
// Including it would diverge the state's algorithm keys from the plan. Terraform reports that as "provider
// produced inconsistent result". trackedAlgorithms is nil on the Read that follows ImportState, since
// ImportState itself never sets the algorithms attribute; that Read includes every signed algorithm instead.
func mapGenerationToUploadModel(ctx context.Context, lineageID int64, generation *cloudcertificates.GetGenerationResponse, trackedAlgorithms map[string]algorithmModel) (uploadResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := uploadResourceModel{
		LineageID:                   types.Int64Value(lineageID),
		GenerationID:                types.Int64Value(generation.GenerationID),
		GenerationStatus:            types.StringValue(generation.GenerationStatus),
		FirstPromotedToProductionAt: fwdate.TimeRFC3339PointerValue(generation.FirstPromotedToProductionTime),
		GenerationCreatedBy:         types.StringPointerValue(generation.GenerationCreatedBy),
		GenerationCreatedTime:       fwdate.TimeRFC3339PointerValue(generation.GenerationCreatedTime),
		GenerationModifiedBy:        types.StringPointerValue(generation.GenerationModifiedBy),
		GenerationModifiedTime:      fwdate.TimeRFC3339PointerValue(generation.GenerationModifiedTime),
	}

	var signed []cloudcertificates.Algorithm
	for _, a := range generation.Algorithms {
		if a.SignedCertificatePEM == nil {
			continue
		}
		if trackedAlgorithms != nil {
			if _, tracked := trackedAlgorithms[a.KeyType]; !tracked {
				continue
			}
		}
		signed = append(signed, a)
	}

	algorithmsMap, d := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: algorithmObjectAttrTypes}, mapLineageAlgorithms(signed))
	diags.Append(d...)
	data.Algorithms = algorithmsMap

	return data, diags
}

// uploadCertificateErrorDiagnostics maps well-documented UploadSignedCertificate sentinel errors to clear diagnostics.
func uploadCertificateErrorDiagnostics(err error, lineageID, generationID int64) diag.Diagnostics {
	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrGenerationImmutable):
		diags.AddError("Generation Immutable",
			fmt.Sprintf("Generation %d of lineage %d has already been promoted to production and can no longer be modified directly.", generationID, lineageID))
	case errors.Is(err, cloudcertificates.ErrCertAlreadyUploaded):
		diags.AddError("Certificate Already Uploaded", "A signed certificate was already uploaded for this algorithm instance.")
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		diags.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d.", lineageID))
	case errors.Is(err, cloudcertificates.ErrCertificateNotFound):
		diags.AddError("Generation Not Found", fmt.Sprintf("No generation %d found for lineage %d.", generationID, lineageID))
	case errors.Is(err, cloudcertificates.ErrCertExpiryInvalid):
		diags.AddError("Invalid Certificate Validity", "The uploaded certificate's validity dates are invalid.")
	case errors.Is(err, cloudcertificates.ErrCertParseError):
		diags.AddError("Certificate Parse Error", "The uploaded certificate is malformed or could not be parsed.")
	case errors.Is(err, cloudcertificates.ErrCertCSRMismatch):
		diags.AddError("Certificate/CSR Mismatch", "The uploaded certificate's public key does not match the stored CSR.")
	case errors.Is(err, cloudcertificates.ErrUnknownKeyType):
		diags.AddError("Unknown Key Type", "The given key type is not part of the lineage's key specs.")
	case errors.Is(err, cloudcertificates.ErrDomainNotValidatedUploadFailed):
		diags.AddError("Domain Not Validated", "One or more SANs are not yet Domain Validated.")
	default:
		diags.AddError("Failed to upload signed certificate", err.Error())
	}
	return diags
}

// completeLineageErrorDiagnostics maps well-documented CompleteLineage sentinel errors to clear diagnostics.
func completeLineageErrorDiagnostics(err error, lineageID int64) diag.Diagnostics {
	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		diags.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d.", lineageID))
	case errors.Is(err, cloudcertificates.ErrNoCurrentProduction):
		diags.AddError("No Current Production", fmt.Sprintf("Lineage %d has no current production generation to complete.", lineageID))
	case errors.Is(err, cloudcertificates.ErrCompletePreconditionFailed):
		diags.AddError("Complete Precondition Failed",
			"A pending head generation already exists, or no unique CSR_READY algorithm matching the given key type was found on production.")
	case errors.Is(err, cloudcertificates.ErrCertParseError):
		diags.AddError("Certificate Parse Error", "The uploaded certificate is malformed or could not be parsed.")
	case errors.Is(err, cloudcertificates.ErrCertCSRMismatch):
		diags.AddError("Certificate/CSR Mismatch", "The uploaded certificate's public key does not match the stored CSR.")
	default:
		diags.AddError("Failed to complete lineage", err.Error())
	}
	return diags
}
