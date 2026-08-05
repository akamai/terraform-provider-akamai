package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/retry"
	"github.com/akamai/terraform-provider-akamai/v11/internal/text"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &activationResource{}
	_ resource.ResourceWithConfigure      = &activationResource{}
	_ resource.ResourceWithImportState    = &activationResource{}
	_ resource.ResourceWithModifyPlan     = &activationResource{}
	_ resource.ResourceWithValidateConfig = &activationResource{}
)

type (
	activationResource struct {
		meta.Resource
		activationResourceConfig
	}

	activationResourceConfig struct {
		defaultTimeout time.Duration
		pollInterval   time.Duration
	}

	// activationResourceModel manages a lineage's activation state on both networks together, so their
	// PromoteLineage calls can be ordered (or combined into one) deliberately instead of racing across two
	// independent resource instances. staging_generation_id and production_generation_id are independent and
	// optional: omitting one leaves that network's last-known activation exactly as tracked in Staging/Production,
	// without ever touching it - the only way to affect a network is to set (or change) its generation_id.
	//
	// Staging/Production are typed as types.Object (rather than a plain *networkActivationModel pointer): they're
	// Computed, so during Create/Update the incoming plan represents them as unknown (not yet determined) rather
	// than null - and only types.Object (not a Go struct) can represent that. See staging/production/setStaging/
	// setProduction for converting to/from the richer networkActivationModel used everywhere else.
	activationResourceModel struct {
		LineageID              types.Int64    `tfsdk:"lineage_id"`
		StagingGenerationID    types.Int64    `tfsdk:"staging_generation_id"`
		ProductionGenerationID types.Int64    `tfsdk:"production_generation_id"`
		Staging                types.Object   `tfsdk:"staging"`
		Production             types.Object   `tfsdk:"production"`
		Timeouts               timeouts.Value `tfsdk:"timeouts"`
	}

	// networkActivationModel is the last-known activation tracked for a single network (staging or production).
	// generation_id here is the actual, tracked generation - it can differ from the top-level
	// staging_generation_id/production_generation_id input when that input has been removed from config (frozen).
	networkActivationModel struct {
		commonActivationModel
		ActivationID types.Int64 `tfsdk:"activation_id"`
	}
)

// networkActivationType returns the attr.Type map for networkActivationModel, derived from the schema so it
// can't drift from the attributes actually declared there. Identical for both networks: only descriptions differ.
func networkActivationType() map[string]attr.Type {
	return networkActivationSchemaAttribute("STAGING").GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

// staging extracts the Staging object into a typed model, or nil if this resource has never managed STAGING.
func (m *activationResourceModel) staging(ctx context.Context) (*networkActivationModel, diag.Diagnostics) {
	if m.Staging.IsNull() || m.Staging.IsUnknown() {
		return nil, nil
	}
	var s networkActivationModel
	return &s, m.Staging.As(ctx, &s, basetypes.ObjectAsOptions{})
}

// production extracts the Production object into a typed model, or nil if this resource has never managed PRODUCTION.
func (m *activationResourceModel) production(ctx context.Context) (*networkActivationModel, diag.Diagnostics) {
	if m.Production.IsNull() || m.Production.IsUnknown() {
		return nil, nil
	}
	var p networkActivationModel
	return &p, m.Production.As(ctx, &p, basetypes.ObjectAsOptions{})
}

// setStaging stores s into the Staging object field - null when s is nil.
func (m *activationResourceModel) setStaging(ctx context.Context, s *networkActivationModel) diag.Diagnostics {
	if s == nil {
		m.Staging = types.ObjectNull(networkActivationType())
		return nil
	}
	var dd diag.Diagnostics
	m.Staging, dd = types.ObjectValueFrom(ctx, networkActivationType(), s)
	return dd
}

// setProduction stores p into the Production object field - null when p is nil.
func (m *activationResourceModel) setProduction(ctx context.Context, p *networkActivationModel) diag.Diagnostics {
	if p == nil {
		m.Production = types.ObjectNull(networkActivationType())
		return nil
	}
	var dd diag.Diagnostics
	m.Production, dd = types.ObjectValueFrom(ctx, networkActivationType(), p)
	return dd
}

// NewActivationResource returns a new Cloud Certificates lineage activation resource.
func NewActivationResource(config activationResourceConfig) func() resource.Resource {
	return func() resource.Resource {
		return &activationResource{activationResourceConfig: config}
	}
}

// defaultActivationResourceConfig returns the production configuration for the lineage activation resource.
func defaultActivationResourceConfig() activationResourceConfig {
	return activationResourceConfig{
		defaultTimeout: 30 * time.Minute,
		pollInterval:   15 * time.Second,
	}
}

// Metadata configures the resource's type name.
func (r *activationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_activation"
}

// Schema defines the Terraform schema for the lineage activation resource.
func (r *activationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Promotes a certificate lineage's head generation to the STAGING and/or PRODUCTION networks. " +
			"staging_generation_id and production_generation_id are independent and optional - at least one must be " +
			"set. When both target the same generation, a single API call activates both networks atomically, " +
			"avoiding the ordering issues that arise from managing them as two separate resources (PROMOTE always " +
			"targets the lineage's head generation, and promoting one network first can move the head on before the " +
			"other network's own request runs). Removing (or never setting) one of the two generation_id attributes " +
			"leaves that network's last-known activation exactly as tracked - Terraform simply stops managing further " +
			"changes to it, without ever deactivating or rolling it back automatically. Changing lineage_id replaces " +
			"the whole resource; changing just one network's generation_id promotes a fresh activation for that " +
			"network in place, without disturbing the other. There's no such thing as editing an activation - PROMOTE " +
			"always creates a brand new one. Activations are immutable historical records - Delete only removes the " +
			"resource from Terraform state, it never reverses a promotion.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the lineage whose head generation to promote.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"staging_generation_id": schema.Int64Attribute{
				Optional: true,
				Description: "Unique identifier of the generation to promote to the STAGING network. Must be the " +
					"lineage's head generation at the time it's first activated on this network. Omit (or remove " +
					"from config) to stop managing this network - its last-known activation is preserved in the " +
					"staging attribute and is never touched.",
				Validators: []validator.Int64{int64validator.AtLeast(1)},
			},
			"production_generation_id": schema.Int64Attribute{
				Optional: true,
				Description: "Unique identifier of the generation to promote to the PRODUCTION network. Must be the " +
					"lineage's head generation at the time it's first activated on this network. Omit (or remove " +
					"from config) to stop managing this network - its last-known activation is preserved in the " +
					"production attribute and is never touched.",
				Validators: []validator.Int64{int64validator.AtLeast(1)},
			},
			"staging":    networkActivationSchemaAttribute("STAGING"),
			"production": networkActivationSchemaAttribute("PRODUCTION"),
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx,
				timeouts.Opts{
					Create: true,
					CreateDescription: "Optional configurable timeout for waiting for a newly created activation to reach a " +
						"terminal status. By default it's 30m with a 15s polling interval.",
					Update: true,
					UpdateDescription: "Optional configurable timeout for waiting for a network's changed activation to reach " +
						"a terminal status. By default it's 30m with a 15s polling interval.",
				},
			),
		},
	}
}

// networkActivationSchemaAttribute returns the computed schema attributes tracked for a single network (staging
// or production) - identical shape for both, only the description differs.
func networkActivationSchemaAttribute(network string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Description: fmt.Sprintf("Last-known activation tracked for the %s network, or null if this resource has never "+
			"managed one. Reflects the actual, tracked generation—which can differ from %s_generation_id when that "+
			"attribute has been removed from config.", network, strings.ToLower(network)),
		Attributes: map[string]schema.Attribute{
			"generation_id": schema.Int64Attribute{
				Computed:    true,
				Description: fmt.Sprintf("Unique identifier of the generation actually tracked as active on the %s network.", network),
			},
			"activation_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier of the activation request.",
			},
			"activation_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the activation operation. Always `PROMOTE` for this resource.",
			},
			"activation_status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the activation request: `INIT`, `PENDING`, `IN_PROGRESS`, `COMPLETE`, `PARTIAL_SUCCESS`, `FAILED`, or `ABORTED`.",
			},
			"activation_created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time the activation request was created.",
			},
			"activation_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time the activation request was last modified.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the activation request.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last modified the activation request.",
			},
			"total_hostname_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of hostnames being deployed as part of this activation, or null if not yet known.",
			},
			"in_progress_hostname_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of hostnames still in progress for this activation, or null if not yet known.",
			},
			"pre_empted_by": schema.Int64Attribute{
				Computed:    true,
				Description: "The activation request that pre-empted (superseded) this one, or null if this activation was not pre-empted.",
			},
			"error_types": schema.StringAttribute{
				Computed:    true,
				Description: "Error type information when the activation failed, or null otherwise.",
			},
		},
	}
}

// ValidateConfig requires at least one of staging_generation_id/production_generation_id to be set - a resource
// managing neither network does nothing at all.
func (r *activationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data activationResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.StagingGenerationID.IsUnknown() || data.ProductionGenerationID.IsUnknown() {
		// can't yet tell whether at least one will end up set. Terraform Core guarantees ValidateConfig is called
		// again once this resolves, so deferring here is safe even when the value comes from another resource
		// created in the same apply.
		return
	}

	if data.StagingGenerationID.IsNull() && data.ProductionGenerationID.IsNull() {
		resp.Diagnostics.AddError("Missing Activation Target",
			"At least one of staging_generation_id or production_generation_id must be set.")
	}
}

// ModifyPlan decides, per network, whether staging/production is about to be freshly promoted this apply - either
// because the plan asks for a new generation, or because the last promotion attempt for the currently-requested
// one never actually completed (e.g. it ended FAILED/ABORTED, so tracked still lags behind what was asked for).
// Only then is that network's Computed object marked unknown, forcing Update to run and recompute it (RES-07): a
// blanket objectplanmodifier.UseStateForUnknown() can't do this, since it can't tell a stale value needing a
// fresh promotion apart from one that's simply not changing - it would keep the stale value in both cases. When
// a network isn't being freshly promoted, its last-known state value is copied into the plan explicitly, since
// the framework otherwise marks every Computed attribute unknown on any update by default. Conflicting promotion
// targets (RES-07-adjacent) are also caught here - on both Create and Update - rather than only at apply time in
// reconcileActivations, so the user sees the error in plan output before Create/Update ever runs.
func (r *activationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if modifiers.IsDelete(req) {
		return
	}

	var plan activationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if modifiers.IsCreate(req) {
		// nothing tracked yet: both networks need promotion whenever their generation_id is configured at all.
		resp.Diagnostics.Append(conflictingPromotionTargetsDiagnostics(plan.StagingGenerationID, plan.ProductionGenerationID,
			networkNeedsPromotion(plan.StagingGenerationID, nil), networkNeedsPromotion(plan.ProductionGenerationID, nil))...)
		return
	}

	var state activationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stagingModel, d := state.staging(ctx)
	resp.Diagnostics.Append(d...)
	productionModel, d := state.production(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	needStaging := networkNeedsPromotion(plan.StagingGenerationID, stagingModel)
	needProduction := networkNeedsPromotion(plan.ProductionGenerationID, productionModel)
	resp.Diagnostics.Append(conflictingPromotionTargetsDiagnostics(plan.StagingGenerationID, plan.ProductionGenerationID, needStaging, needProduction)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if needStaging {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("staging"), types.ObjectUnknown(networkActivationType()))...)
	} else {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("staging"), state.Staging)...)
	}
	if needProduction {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("production"), types.ObjectUnknown(networkActivationType()))...)
	} else {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("production"), state.Production)...)
	}
}

// networkNeedsPromotion reports whether want - a network's desired generation, from either a plan or config -
// requires a fresh PromoteLineage call given tracked, that network's last-known activation: true when want is
// set and either nothing is tracked yet or tracked's own generation differs from it. This covers both a
// genuinely new target and a retry of one whose last attempt never completed (tracked still lags behind want).
// An unknown want can't yet be compared, so it's conservatively treated as needing one too. A null want means
// the network was never requested, or is intentionally frozen - not an error.
func networkNeedsPromotion(want types.Int64, tracked *networkActivationModel) bool {
	if want.IsUnknown() {
		return true
	}
	if want.IsNull() {
		return false
	}
	if tracked == nil {
		return true
	}
	return tracked.GenerationID.ValueInt64() != want.ValueInt64()
}

// conflictingPromotionTargetsDiagnostics reports an error when both networks need a fresh promotion this apply
// but target different generations: PROMOTE always targets the lineage's single current head, so two different
// generations can never both be the head at once - promoting them separately in the same operation would
// inevitably fail one of them against the real API. Deferred (no diagnostics) while either target is still
// unknown, since the values can't yet be compared; Terraform Core's second, fully-resolved ModifyPlan pass (or
// reconcileActivations at apply time) still catches it once resolved, if it truly is a conflict.
func conflictingPromotionTargetsDiagnostics(wantStaging, wantProduction types.Int64, needStaging, needProduction bool) diag.Diagnostics {
	var diags diag.Diagnostics
	if !needStaging || !needProduction || wantStaging.IsUnknown() || wantProduction.IsUnknown() {
		return diags
	}
	if wantStaging.ValueInt64() == wantProduction.ValueInt64() {
		return diags
	}
	diags.AddError("Conflicting Promotion Targets",
		fmt.Sprintf("staging_generation_id (%d) and production_generation_id (%d) both need a fresh promotion in this "+
			"apply, but PROMOTE always targets the lineage's current head generation - only one of them can actually be "+
			"the head at a time. Apply one network's change first, then the other's in a later apply, or set both to the "+
			"same generation_id to promote them together.", wantStaging.ValueInt64(), wantProduction.ValueInt64()))
	return diags
}

// Create issues the promotion(s) for whichever of staging_generation_id/production_generation_id are set, and
// waits for each to reach a terminal status.
func (r *activationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Activation Resource Create")

	var plan activationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, r.defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	staging, production, diags := r.reconcileActivations(ctx, plan.LineageID.ValueInt64(), plan.StagingGenerationID, plan.ProductionGenerationID, nil, nil, createTimeout)
	resp.Diagnostics.Append(diags...)
	if staging == nil && production == nil {
		return
	}

	data := &activationResourceModel{
		LineageID:              plan.LineageID,
		StagingGenerationID:    plan.StagingGenerationID,
		ProductionGenerationID: plan.ProductionGenerationID,
		Timeouts:               plan.Timeouts,
	}
	resp.Diagnostics.Append(data.setStaging(ctx, staging)...)
	resp.Diagnostics.Append(data.setProduction(ctx, production)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// reconcileActivations promotes staging and/or production to their desired generation, calling PromoteLineage
// only for whichever network actually needs it. It returns the resulting staging/production models to persist.
// A network needs a fresh call when its desired generation differs from what's tracked, or when nothing is
// tracked yet. It issues a single combined call when both networks need one and target the same generation.
// When both need one but target different generations, this is rejected outright: PROMOTE always targets the
// lineage's single current head, so two different generations can never both be the head at once - promoting
// them separately in the same operation would inevitably fail one of them against the real API. A null desired
// generation freezes that network at whatever was already tracked. Removing an Optional generation_id attribute
// from config never touches the real, already-promoted activation on that network. This matches
// property_activation's own pollActivation: a failed promotion attempt is simply discarded rather than
// persisted, so that network keeps whatever was tracked for it before this call (nil if nothing was), and its
// generation_id is never marked as satisfied - the next apply automatically retries it fresh. The other
// network's successful result is still recorded.
func (r *activationResource) reconcileActivations(ctx context.Context, lineageID int64, wantStaging, wantProduction types.Int64, trackedStaging, trackedProduction *networkActivationModel, timeout time.Duration) (*networkActivationModel, *networkActivationModel, diag.Diagnostics) {
	staging := trackedStaging
	production := trackedProduction
	var diags diag.Diagnostics

	needStaging := networkNeedsPromotion(wantStaging, trackedStaging)
	needProduction := networkNeedsPromotion(wantProduction, trackedProduction)

	if d := conflictingPromotionTargetsDiagnostics(wantStaging, wantProduction, needStaging, needProduction); d.HasError() {
		diags.Append(d...)
		return staging, production, diags
	}

	if needStaging && needProduction && wantStaging.ValueInt64() == wantProduction.ValueInt64() {
		results, d := r.promoteNetworksAndWait(ctx, lineageID, wantStaging.ValueInt64(),
			[]cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging, cloudcertificates.TargetNetworkProduction}, timeout)
		diags.Append(d...)
		if results[0] != nil {
			staging = results[0]
		}
		if results[1] != nil {
			production = results[1]
		}
		return staging, production, diags
	}

	if needStaging {
		results, d := r.promoteNetworksAndWait(ctx, lineageID, wantStaging.ValueInt64(), []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging}, timeout)
		diags.Append(d...)
		if results[0] != nil {
			staging = results[0]
		}
	}
	if needProduction {
		results, d := r.promoteNetworksAndWait(ctx, lineageID, wantProduction.ValueInt64(), []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction}, timeout)
		diags.Append(d...)
		if results[0] != nil {
			production = results[0]
		}
	}
	return staging, production, diags
}

// promoteNetworksAndWait promotes generationID to networks (one or two) in a single PromoteLineage call, then finishes
// each network's activation. It returns one model per network, in the same order as networks; a network has a nil
// model in its position only when PromoteLineage's response omitted it, or its activation ended FAILED/ABORTED.
// timeout bounds the whole call; when networks has two entries (a combined promotion), they share this single
// budget sequentially, since they're really one atomic operation.
func (r *activationResource) promoteNetworksAndWait(ctx context.Context, lineageID, generationID int64, networks []cloudcertificates.TargetNetwork, timeout time.Duration) ([]*networkActivationModel, diag.Diagnostics) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var diags diag.Diagnostics

	results := make([]*networkActivationModel, len(networks))

	result, err := r.Client.GetCloudCertificates().PromoteLineage(ctx, cloudcertificates.PromoteLineageRequest{
		LineageID:    lineageID,
		GenerationID: generationID,
		Networks:     networks,
	})
	if err != nil {
		diags.Append(promoteLineageErrorDiagnostics(err, lineageID, generationID)...)
		return results, diags
	}
	if result == nil {
		diags.AddError("Generation Already Active",
			fmt.Sprintf("Generation %d of lineage %d is already active on %s; there is nothing new to activate. If you "+
				"want to manage the existing activation(s) in Terraform, find their activation IDs (e.g. via the "+
				"akamai_cloudcertificates_activation_status data source) and import them.",
				generationID, lineageID, text.JoinStringBased(networks, " and ")))
		return results, diags
	}

	for i, network := range networks {
		model, d := r.finishNetworkActivation(ctx, lineageID, network, result.Items)
		diags.Append(d...)
		results[i] = model
	}
	return results, diags
}

// finishNetworkActivation extracts the activation for network from items, polls it to a terminal status, and
// builds the per-network model to persist. It returns nil on a transient read error, on the configured timeout
// elapsing, or when the activation itself ends FAILED/ABORTED. In every such case the promotion may have already
// succeeded, but its outcome can't be confirmed here, so nothing is tracked for it. Recovery is a fresh
// PromoteLineage call for the same generation on the next apply. The API's own response then determines whether
// that's actually needed. For example, it returns ErrPendingActivationInProgress if the original attempt is
// still in flight, or a nil "already active" result if it already completed. Polling is bounded by ctx's own
// deadline rather than a fresh one of its own, so callers polling more than one network in turn share a single
// overall budget.
func (r *activationResource) finishNetworkActivation(ctx context.Context, lineageID int64, network cloudcertificates.TargetNetwork, items []cloudcertificates.GetActivationStatusResponse) (*networkActivationModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	activation := findActivationForNetwork(items, network)
	if activation == nil {
		diags.AddError("Unexpected Response", fmt.Sprintf("Did not receive an activation for the %s network.", network))
		return nil, diags
	}

	// Poll GetActivationStatus until the activation reaches a terminal status (COMPLETE, PARTIAL_SUCCESS, FAILED,
	// or ABORTED) or ctx's deadline elapses.
	final, err := retry.Poll(ctx, retry.PollingOpts[cloudcertificates.GetActivationStatusResponse]{
		Fn: func(ctx context.Context) (*cloudcertificates.GetActivationStatusResponse, error) {
			return r.Client.GetCloudCertificates().GetActivationStatus(ctx, cloudcertificates.GetActivationStatusRequest{
				LineageID:    lineageID,
				ActivationID: activation.ActivationID,
			})
		},
		ShouldRetryData: func(status cloudcertificates.GetActivationStatusResponse) bool {
			return !isTerminalActivationStatus(status.ActivationStatus)
		},
		Interval: r.pollInterval,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			diags.AddError("Reached Activation Timeout",
				fmt.Sprintf("Activation %d for lineage %d did not reach a terminal status within the configured timeout.",
					activation.ActivationID, lineageID))
			return nil, diags
		}
		diags.AddError("Failed To Read Activation Status",
			fmt.Sprintf("Failed to read activation %d for lineage %d: %s", activation.ActivationID, lineageID, err.Error()))
		return nil, diags
	}

	diags.Append(activationTerminalStatusDiagnostics(*final)...)
	if final.ActivationStatus == "FAILED" || final.ActivationStatus == "ABORTED" {
		return nil, diags
	}

	return mapNetworkActivationModel(*final), diags
}

// findActivationForNetwork returns the item in items whose TargetEnvironment matches network, or nil if none does.
func findActivationForNetwork(items []cloudcertificates.GetActivationStatusResponse, network cloudcertificates.TargetNetwork) *cloudcertificates.GetActivationStatusResponse {
	for i := range items {
		if items[i].TargetEnvironment == string(network) {
			return &items[i]
		}
	}
	return nil
}

func isTerminalActivationStatus(status string) bool {
	switch status {
	case "COMPLETE", "PARTIAL_SUCCESS", "FAILED", "ABORTED":
		return true
	default:
		return false
	}
}

// activationTerminalStatusDiagnostics reports FAILED/ABORTED activations as errors and PARTIAL_SUCCESS as a
// warning. In all three cases the activation is still a real, immutable record, so the caller should persist it
// to state regardless of the diagnostic added here.
func activationTerminalStatusDiagnostics(activation cloudcertificates.GetActivationStatusResponse) diag.Diagnostics {
	var diags diag.Diagnostics
	detail := fmt.Sprintf("Activation %d for lineage %d ended with status %s.", activation.ActivationID, activation.LineageID, activation.ActivationStatus)
	if activation.ErrorTypes != "" {
		detail += fmt.Sprintf(" Error types: %s.", activation.ErrorTypes)
	}
	switch activation.ActivationStatus {
	case "FAILED":
		diags.AddError("Activation Failed", detail)
	case "ABORTED":
		diags.AddError("Activation Aborted", detail)
	case "PARTIAL_SUCCESS":
		diags.AddWarning("Activation Partially Succeeded", detail)
	}
	return diags
}

// promoteLineageErrorDiagnostics maps well-documented PromoteLineage sentinel errors to clear diagnostics.
func promoteLineageErrorDiagnostics(err error, lineageID, generationID int64) diag.Diagnostics {
	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		diags.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d.", lineageID))
	case errors.Is(err, cloudcertificates.ErrLineageNoHeadGeneration):
		diags.AddError("No Head Generation", fmt.Sprintf("Lineage %d has no head generation to activate.", lineageID))
	case errors.Is(err, cloudcertificates.ErrIncompleteCertMaterial):
		diags.AddError("Incomplete Certificate Material",
			fmt.Sprintf("Generation %d of lineage %d has no algorithm instance ready for use; upload a signed certificate before activating.", generationID, lineageID))
	case errors.Is(err, cloudcertificates.ErrSingleGenerationActivationNotSupported):
		diags.AddError("Activation Not Supported", fmt.Sprintf("Lineage %d is a SINGLE_GENERATION lineage; activation operations are not supported.", lineageID))
	case errors.Is(err, cloudcertificates.ErrActivationCooldownInEffect):
		diags.AddError("Activation Cooldown In Effect", fmt.Sprintf("Lineage %d had a recent activation that is still propagating; wait before starting another.", lineageID))
	case errors.Is(err, cloudcertificates.ErrPendingActivationInProgress):
		diags.AddError("Pending Activation In Progress", fmt.Sprintf("Lineage %d already has a pending activation in progress; wait for it to complete before starting another.", lineageID))
	case errors.Is(err, cloudcertificates.ErrUpstreamActivationError):
		diags.AddError("Upstream Activation Error", "An upstream service error occurred while processing the activation; this is typically transient and retrying may succeed.")
	case errors.Is(err, cloudcertificates.ErrLineageBadRequest):
		diags.AddError("Invalid Activation Request",
			fmt.Sprintf("Generation %d is not a valid promotion target for lineage %d - it must be the lineage's head generation. %s", generationID, lineageID, err.Error()))
	default:
		diags.AddError("Failed to promote lineage", fmt.Sprintf("Failed to promote generation %d of lineage %d: %s", generationID, lineageID, err))
	}
	return diags
}

// Read refreshes the tracked activation status for whichever of staging/production this resource is managing.
// Activations are immutable historical records that should never disappear once created, so a genuinely missing
// activation (ErrActivationNotFound) is treated as drift: that network's tracked object is cleared with a
// warning rather than the whole resource being removed, since the other network's activation is still valid.
// Clearing it (while leaving staging_generation_id/production_generation_id untouched) is exactly the signal
// ModifyPlan/networkNeedsPromotion already looks for, so the very next plan automatically re-promotes that network -
// no manual intervention needed. Any other read failure remains a hard error, and leaves the last-known value
// untouched, since it doesn't necessarily mean the activation is actually gone. The two networks are refreshed
// independently: one failing (or disappearing) never skips the other.
func (r *activationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Activation Resource Read")

	var state activationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lineageID := state.LineageID.ValueInt64()

	stagingModel, d := state.staging(ctx)
	resp.Diagnostics.Append(d...)
	productionModel, d := state.production(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if stagingModel != nil {
		refreshed, diags, gone := r.refreshActivation(ctx, lineageID, stagingModel.ActivationID.ValueInt64(), "STAGING")
		resp.Diagnostics.Append(diags...)
		if refreshed != nil {
			resp.Diagnostics.Append(state.setStaging(ctx, refreshed)...)
		} else if gone {
			resp.Diagnostics.Append(state.setStaging(ctx, nil)...)
		}
	}

	if productionModel != nil {
		refreshed, diags, gone := r.refreshActivation(ctx, lineageID, productionModel.ActivationID.ValueInt64(), "PRODUCTION")
		resp.Diagnostics.Append(diags...)
		if refreshed != nil {
			resp.Diagnostics.Append(state.setProduction(ctx, refreshed)...)
		} else if gone {
			resp.Diagnostics.Append(state.setProduction(ctx, nil)...)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// refreshActivation re-fetches a tracked activation's current status. If it has genuinely disappeared
// (ErrActivationNotFound), it returns (nil, a warning, true) instead of a hard error, so the caller clears just
// that network's tracked object - the next plan then detects the gap and automatically re-promotes it. Any
// other error returns (nil, an error, false), leaving the network's last-known value untouched, since it
// doesn't necessarily mean the activation is actually gone. A tracked activation can still be non-terminal here
// (e.g. persisted mid-flight after Create/Update's own poll timed out): reapplies the same terminal-status
// diagnostics as the polling path, so a later FAILED/ABORTED/PARTIAL_SUCCESS transition is never reported back
// as a silent success just because it happened after the resource was already created/updated.
func (r *activationResource) refreshActivation(ctx context.Context, lineageID, activationID int64, network string) (*networkActivationModel, diag.Diagnostics, bool) {
	var diags diag.Diagnostics

	activation, err := r.Client.GetCloudCertificates().GetActivationStatus(ctx, cloudcertificates.GetActivationStatusRequest{
		LineageID:    lineageID,
		ActivationID: activationID,
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrActivationNotFound) {
			diags.AddWarning(fmt.Sprintf("%s Activation No Longer Found", network),
				fmt.Sprintf("Activation %d for lineage %d (network %s) no longer exists; it may have been deleted or "+
					"rolled back outside Terraform. This network's tracked activation has been cleared from state; "+
					"the next plan will automatically re-promote it.", activationID, lineageID, network))
			return nil, diags, true
		}
		diags.AddError(fmt.Sprintf("Failed To Read %s Activation Status", network),
			fmt.Sprintf("Activation %d for lineage %d (network %s) could not be read; this shouldn't normally happen "+
				"since activations are immutable historical records, so it's likely transient: %s. State is unchanged; "+
				"running plan or apply again will simply retry the read.", activationID, lineageID, network, err.Error()))
		return nil, diags, false
	}

	diags.Append(activationTerminalStatusDiagnostics(*activation)...)
	return mapNetworkActivationModel(*activation), diags, false
}

// Update reconciles staging/production against whichever of staging_generation_id/production_generation_id
// changed. See reconcileActivations for the exact rules (combined call, freeze-on-null, etc.).
func (r *activationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Activation Resource Update")

	var plan, state activationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, r.defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trackedStaging, d := state.staging(ctx)
	resp.Diagnostics.Append(d...)
	trackedProduction, d := state.production(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	staging, production, diags := r.reconcileActivations(ctx, plan.LineageID.ValueInt64(), plan.StagingGenerationID, plan.ProductionGenerationID, trackedStaging, trackedProduction, updateTimeout)
	resp.Diagnostics.Append(diags...)

	data := &activationResourceModel{
		LineageID:              plan.LineageID,
		StagingGenerationID:    plan.StagingGenerationID,
		ProductionGenerationID: plan.ProductionGenerationID,
		Timeouts:               plan.Timeouts,
	}
	resp.Diagnostics.Append(data.setStaging(ctx, staging)...)
	resp.Diagnostics.Append(data.setProduction(ctx, production)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Delete removes the resource from Terraform state only. Activations are immutable historical records: they
// can't be reversed, and Cloud Certificate Manager doesn't support deleting them directly.
func (r *activationResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage Activation Resource Delete")
}

// ImportState imports by lineage_id plus one or two activation IDs: `lineageID,activationID` manages a single
// network, self-detected from the activation's own target_environment; `lineageID,activationID,activationID`
// manages both (order-independent, since each is self-describing too). Attributes are set individually via
// SetAttribute (rather than a single whole-model State.Set) so that untouched attributes - notably timeouts -
// are left at the framework's own default null value of the correct type, instead of our own model's zero value.
func (r *activationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := text.ImportIDSplitter("lineageID,activationID[,activationID]").AcceptLen(2).AcceptLen(3).Split(req.ID)
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("lineage_id"), lineageID)...)

	seenNetworks := map[string]bool{}
	for _, part := range parts[1:] {
		activationID, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("expected a numeric activation_id, got: %q", part))
			return
		}
		if activationID <= 0 {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("activation_id must be greater than 0, got: %d", activationID))
			return
		}

		activation, err := r.Client.GetCloudCertificates().GetActivationStatus(ctx, cloudcertificates.GetActivationStatusRequest{
			LineageID:    lineageID,
			ActivationID: activationID,
		})
		if err != nil {
			resp.Diagnostics.AddError("Failed to read activation status", err.Error())
			return
		}

		if seenNetworks[activation.TargetEnvironment] {
			resp.Diagnostics.AddError("Duplicate Network In Import ID",
				fmt.Sprintf("Activation %d is also on the %s network, which was already imported from an earlier activation ID.",
					activationID, activation.TargetEnvironment))
			return
		}
		seenNetworks[activation.TargetEnvironment] = true

		model := mapNetworkActivationModel(*activation)
		switch activation.TargetEnvironment {
		case string(cloudcertificates.TargetNetworkStaging):
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("staging_generation_id"), activation.GenerationID)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("staging"), model)...)
		case string(cloudcertificates.TargetNetworkProduction):
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("production_generation_id"), activation.GenerationID)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("production"), model)...)
		default:
			resp.Diagnostics.AddError("Unexpected Network",
				fmt.Sprintf("Activation %d has an unrecognized target network %q.", activationID, activation.TargetEnvironment))
			return
		}
	}
}

// mapNetworkActivationModel maps a GetActivationStatusResponse onto a networkActivationModel.
func mapNetworkActivationModel(activation cloudcertificates.GetActivationStatusResponse) *networkActivationModel {
	return &networkActivationModel{
		commonActivationModel: mapActivationCore(activation),
		ActivationID:          types.Int64Value(activation.ActivationID),
	}
}
