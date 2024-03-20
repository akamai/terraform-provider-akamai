package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &ConfigurationResource{}
	_ resource.ResourceWithConfigure   = &ConfigurationResource{}
	_ resource.ResourceWithModifyPlan  = &ConfigurationResource{}
	_ resource.ResourceWithImportState = &ConfigurationResource{}
)

// ConfigurationResource represents akamai_cloudwrapper_configuration resource
type ConfigurationResource struct {
	client        cloudwrapper.CloudWrapper
	deleteTimeout time.Duration
	pollInterval  time.Duration
}

func (r *ConfigurationResource) setClient(client cloudwrapper.CloudWrapper) {
	r.client = client
}

func (r *ConfigurationResource) setPollInterval(duration time.Duration) {
	r.pollInterval = duration
}

// NewConfigurationResource returns new cloud wrapper configuration resource
func NewConfigurationResource() resource.Resource {
	return &ConfigurationResource{
		deleteTimeout: 2 * time.Hour,
		pollInterval:  30 * time.Second,
	}
}

// Metadata implements resource.Resource.
func (r *ConfigurationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_configuration"
}

// Schema implements resource.Resource.
func (r *ConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Contract ID having Cloud Wrapper entitlement.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"config_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"property_ids": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of properties belonging to eligible products.",
				PlanModifiers: []planmodifier.Set{
					modifiers.SetUseStateIf(modifiers.EqualUpToPrefixFunc("prp_")),
				},
			},
			"comments": schema.StringAttribute{
				Required:    true,
				Description: "Additional information you provide to differentiate or track changes of the configuration.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"retain_idle_objects": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Retain idle objects beyond their max idle lifetime.",
				Default:     booldefault.StaticBool(false),
			},
			"capacity_alerts_threshold": schema.Int64Attribute{
				Optional: true,
				Description: "Capacity Alerts enablement information for the configuration. " +
					"The Alert Threshold should be between 50 and 100.",
				Validators: []validator.Int64{
					int64validator.Between(50, 100),
				},
			},
			"notification_emails": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Email addresses to use for notifications.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"revision": schema.StringAttribute{
				Computed:    true,
				Description: "Unique hash value of the configuration.",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Resource's unique identifier.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"location": schema.SetNestedBlock{
				Description: "List of locations to use with the configuration.",
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"comments": schema.StringAttribute{
							Required:    true,
							Description: "Additional comments provided by the user.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"traffic_type_id": schema.Int64Attribute{
							Required:    true,
							Description: "Unique identifier for the location and traffic type combination",
						},
					},
					Blocks: map[string]schema.Block{
						"capacity": schema.SingleNestedBlock{
							Description: "The capacity assigned to this configuration's location",
							Validators: []validator.Object{
								objectvalidator.IsRequired(),
							},
							Attributes: map[string]schema.Attribute{
								"value": schema.Int64Attribute{
									Required:    true,
									Description: "Value of capacity.",
								},
								"unit": schema.StringAttribute{
									Required:    true,
									Description: "Unit of capacity. Can be either 'GB' or 'TB'.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											string(cloudwrapper.UnitGB),
											string(cloudwrapper.UnitTB),
										),
									},
								},
							},
						},
					},
				},
			},
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Delete:            true,
				CreateDescription: "Optional configurable resource delete timeout. By default it's 2h with 30s pooling interval.",
			}),
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *ConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	if r.client != nil {
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

	r.client = cloudwrapper.Client(meta.Session())
}

// ModifyPlan implements resource.ResourceWithModifyPlan.
func (*ConfigurationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// config will be deleted
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning("Deletion May Not Succeed",
			"Only Akamai internal users can delete configurations. I you are not internal user, "+
				"the configuration will only be removed from state")
		return
	}

	var plan *ConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.hasUnknown() {
		return
	}

	plan.Revision = types.StringValue(plan.revision(ctx))
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)

	var state *ConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if onlyTimeoutChanged(state, plan) {
		resp.Diagnostics.Append(onlyTimeoutChangeWarn)
	}
}

func onlyTimeoutChanged(state, plan *ConfigurationResourceModel) bool {
	return state != nil && plan != nil &&
		state.Revision == plan.Revision &&
		!state.Timeouts.Equal(plan.Timeouts)
}

// Create implements resource.Resource.
func (r *ConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Configuration Resource")

	var data *ConfigurationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.create(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) create(ctx context.Context, data *ConfigurationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	resp, err := r.client.CreateConfiguration(ctx, data.buildCreateRequest(ctx))
	if err != nil {
		diags.AddError("Create Failed", err.Error())
		return diags
	}

	return data.populateFrom(ctx, resp)
}

var diagErrConfigurationNotFound = diag.NewErrorDiagnostic("Cannot Find Configuration", "")

// Read implements resource.Resource.
func (r *ConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Configuration Resource")

	var data *ConfigurationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.read(ctx, data)
	if diags.Contains(diagErrConfigurationNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) read(ctx context.Context, data *ConfigurationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	result, err := r.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{
		ConfigID: data.ID.ValueInt64(),
	})
	if errors.Is(err, cloudwrapper.ErrConfigurationNotFound) {
		diags.Append(diagErrConfigurationNotFound)
		return diags
	}
	if err != nil {
		diags.AddError("Reading Configuration Failed", err.Error())
		return diags
	}

	if result.MultiCDNSettings != nil {
		diags.AddError("Configuration Contains Multi CDN Settings",
			"Cloud Wrapper Configuration resource does not currently support Mutli CDN settings. "+
				"This error is caused by a configuration drift. "+
				"Make sure to remove Mutli CDN settings before continuing.")
		return diags
	}

	return data.populateFrom(ctx, result)
}

// Update implements resource.Resource.
func (r *ConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Configuration Resource")

	var data *ConfigurationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var oldState *ConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if onlyTimeoutChanged(oldState, data) {
		oldState.Timeouts = data.Timeouts
		resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
		return
	}

	resp.Diagnostics.Append(r.update(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) update(ctx context.Context, data *ConfigurationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	resp, err := r.client.UpdateConfiguration(ctx, data.buildUpdateRequest(ctx))
	if err != nil {
		diags.AddError("Update Failed", err.Error())
		return diags
	}

	return data.populateFrom(ctx, resp)
}

// Delete implements resource.Resource.
func (r *ConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Configuration Resource")

	var data *ConfigurationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, r.deleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	isPending, diags := r.isPendingDelete(ctx, data.ID.ValueInt64())
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if isPending {
		resp.Diagnostics.Append(r.waitForDelete(ctx, data.ID.ValueInt64())...)
		return
	}

	err := r.client.DeleteConfiguration(ctx, cloudwrapper.DeleteConfigurationRequest{
		ConfigID: data.ID.ValueInt64(),
	})
	if errors.Is(err, cloudwrapper.ErrDeletionNotAllowed) {
		resp.Diagnostics.AddWarning("Deletion Unsuccessful",
			"Configuration only removed from state. "+
				fmt.Sprintf("To completely remove this configuration [id=%d], contact your akamai representative.", data.ID.ValueInt64()))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Deletion Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(r.waitForDelete(ctx, data.ID.ValueInt64())...)
}

func (r *ConfigurationResource) isPendingDelete(ctx context.Context, id int64) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	resp, err := r.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{
		ConfigID: id,
	})
	if err != nil {
		diags.AddError("Error Retreiving Configutation", err.Error())
		return false, diags
	}

	return resp.Status == cloudwrapper.StatusDeleteInProgress, diags
}

func (r *ConfigurationResource) waitForDelete(ctx context.Context, id int64) diag.Diagnostics {
	var diags diag.Diagnostics
	for {
		_, err := r.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{
			ConfigID: id,
		})
		if errors.Is(err, cloudwrapper.ErrConfigurationNotFound) {
			return diags
		}

		select {
		case <-time.Tick(r.pollInterval):
			continue
		case <-ctx.Done():
			diags.AddError("Deletion Terminated",
				"context terminated the wait for deletion to finish")
			return diags
		}
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *ConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Configuration Resource")

	var data = &ConfigurationResourceModel{}

	configID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Incorrect ID", err.Error())
		return
	}

	result, err := r.client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot Find Configuration", err.Error())
		return
	}

	if result.MultiCDNSettings != nil {
		resp.Diagnostics.AddError("Cannot Import",
			"Importing configuration with Multi CDN is not supported")
		return
	}

	data.Timeouts = timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"delete": types.StringType,
		}),
	}

	resp.Diagnostics.Append(data.populateFrom(ctx, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
