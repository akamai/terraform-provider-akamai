package cloudwrapper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &activationResource{}
	_ resource.ResourceWithConfigure   = &activationResource{}
	_ resource.ResourceWithModifyPlan  = &activationResource{}
	_ resource.ResourceWithImportState = &activationResource{}
)

var (
	activationTimeout     = 4 * time.Hour
	onlyTimeoutChangeWarn = diag.NewWarningDiagnostic("Update with no API calls", "requested only timeout change; API won't be called")
)

const readError = "could not read Config from API"

type activationResource struct {
	client                 cloudwrapper.CloudWrapper
	activationPollInterval time.Duration
}

// ModifyPlan implements resource.ResourceWithModifyPlan
func (a *activationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning("Deactivation is not Available", "currently it's not possible to deactivate configuration; removing only local state")
		return
	}

	var state, plan *activationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if onlyChangeInTimeout(state, plan) {
		resp.Diagnostics.Append(onlyTimeoutChangeWarn)
	}
}

// NewActivationResource returns new cloud wrapper activation resource
func NewActivationResource() resource.Resource {
	return &activationResource{}
}

// Metadata implements resource.Resource
func (a *activationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_activation"
}

func (a *activationResource) setClient(client cloudwrapper.CloudWrapper) {
	a.client = client
}

func (a *activationResource) setPollInterval(interval time.Duration) {
	a.activationPollInterval = interval
}

// Schema implements resource.Resource
func (a *activationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "The configuration you want to activate.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"revision": schema.StringAttribute{
				Required:    true,
				Description: "Unique hash value of the configuration.",
			},
			"id": schema.StringAttribute{
				Computed:           true,
				Description:        "ID of the resource.",
				DeprecationMessage: "Required by the terraform plugin testing framework, always set to `akamai_cloudwrapper_activation`.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx,
				timeouts.Opts{
					Create:            true,
					CreateDescription: "Optional configurable activation timeout to be used on resource create. By default it's 4h with 1m pooling interval.",
					Update:            true,
					UpdateDescription: "Optional configurable activation timeout to be used on resource update. By default it's 4h with 1m pooling interval.",
				},
			),
		},
	}
}

// Configure implements implements resource.ResourceWithConfigure
func (a *activationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	a.configureResource(req, resp)
	if a.activationPollInterval == 0 {
		a.activationPollInterval = time.Minute
	}
}

func (a *activationResource) configureResource(req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if a.client != nil {
		return
	}

	meta, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.client = cloudwrapper.Client(meta.Session())
}

type activationResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	ConfigID types.Int64    `tfsdk:"config_id"`
	Revision types.String   `tfsdk:"revision"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// Create implements resource.Resource
func (a *activationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Activation Resource")

	var data activationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, activationTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := a.upsert(ctx, data, createTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState.Timeouts = data.Timeouts
	newState.ID = types.StringValue("akamai_cloudwrapper_activation")
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (a *activationResource) upsert(ctx context.Context, data activationResourceModel, activationTimeout time.Duration) (*activationResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	configID := int(data.ConfigID.ValueInt64())
	err := a.client.ActivateConfiguration(ctx, cloudwrapper.ActivateConfigurationRequest{ConfigurationIDs: []int{configID}})
	if err != nil {
		diags.AddError("Activating Configuration Failed", err.Error())
		return nil, diags
	}

	diags.Append(a.waitUntilActivationCompleted(ctx, configID, activationTimeout)...)
	if diags.HasError() {
		return nil, diags
	}

	newState, err := a.readStateFromAPI(ctx, data, int64(configID))
	if err != nil {
		diags.AddError(readError, err.Error())
		return nil, diags
	}
	return newState, diags
}

func (a *activationResource) readStateFromAPI(ctx context.Context, model activationResourceModel, configID int64) (*activationResourceModel, error) {
	configuration, err := a.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{ConfigID: configID})
	if err != nil {
		return nil, err
	}
	if configuration.Status == cloudwrapper.StatusActive {
		model.ConfigID = types.Int64Value(configuration.ConfigID)
		model.Revision = types.StringValue(calculateRevision(configuration))
	} else {
		model.ConfigID = types.Int64Null()
		model.Revision = types.StringNull()
	}
	return &model, nil
}

func (a *activationResource) waitUntilActivationCompleted(ctx context.Context, configID int, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		configuration, err := a.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{ConfigID: int64(configID)})
		if err != nil {
			diags.AddError(readError, err.Error())
			return diags
		}
		if configuration.Status == cloudwrapper.StatusActive {
			return diags
		}

		select {
		case <-time.After(a.activationPollInterval):
			continue
		case <-timeoutCtx.Done():
			diags.AddError("Reached Activation Timeout", timeoutCtx.Err().Error())
			return diags
		}
	}
}

// Read implements resource.Resource
func (a *activationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Activation Resource")

	var data activationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := a.readStateFromAPI(ctx, data, data.ConfigID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(readError, err.Error())
		return
	}

	if newState.ConfigID.IsNull() {
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
	}
}

// Update implements resource.Resource
func (a *activationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Activation Resource")

	var plan, oldState *activationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if onlyChangeInTimeout(oldState, plan) {
		oldState.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(resp.State.Set(ctx, oldState)...)
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, activationTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := a.upsert(ctx, *plan, updateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func onlyChangeInTimeout(state, plan *activationResourceModel) bool {
	return state != nil && plan != nil &&
		plan.ConfigID == state.ConfigID &&
		plan.Revision == state.Revision &&
		!plan.Timeouts.Equal(state.Timeouts)
}

// Delete implements resource.Resource
func (a *activationResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Activation Resource")
	resp.State.RemoveResource(ctx)
}

// ImportState implements resource.ResourceWithImportState
func (a *activationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Activation Resource")

	configID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Activation Resource ID has incorrect Value", err.Error())
		return
	}
	config := activationResourceModel{
		Timeouts: getDefaultTimeoutValue(),
	}
	newState, err := a.readStateFromAPI(ctx, config, int64(configID))
	if err != nil {
		resp.Diagnostics.AddError(readError, err.Error())
		return
	}

	if newState.ConfigID.IsNull() {
		resp.Diagnostics.AddError("Import Failed", "configuration must be active prior to import; activate configuration instead")
		return
	}

	newState.ID = types.StringValue("akamai_cloudwrapper_activation")
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func getDefaultTimeoutValue() timeouts.Value {
	return timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"update": types.StringType,
		}),
	}
}
