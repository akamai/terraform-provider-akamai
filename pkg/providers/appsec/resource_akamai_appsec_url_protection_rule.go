package appsec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &urlProtectionRuleResource{}
	_ resource.ResourceWithConfigure      = &urlProtectionRuleResource{}
	_ resource.ResourceWithImportState    = &urlProtectionRuleResource{}
	_ resource.ResourceWithValidateConfig = &urlProtectionRuleResource{}
)

// urlProtectionRuleResource represents akamai_appsec_url_protection_rule resource.
// Import format: "<config_id>:<url_protection_rule_id>".
type urlProtectionRuleResource struct {
	meta meta.Meta
}

type urlProtectionRuleResourceModel struct {
	ConfigID                types.Int64                   `tfsdk:"config_id"`
	URLProtectionID         types.Int64                   `tfsdk:"url_protection_rule_id"`
	Name                    types.String                  `tfsdk:"name"`
	Description             types.String                  `tfsdk:"description"`
	BypassConditions        types.List                    `tfsdk:"bypass_conditions"`
	MaxRateThreshold        types.Int64                   `tfsdk:"max_rate_threshold"`
	APIDefinitions          types.List                    `tfsdk:"api_definitions"`
	HostnamePaths           types.List                    `tfsdk:"hostname_paths"`
	IntelligentLoadShedding *intelligentLoadSheddingModel `tfsdk:"intelligent_load_shedding"`
	CreateDate              types.String                  `tfsdk:"create_date"`
	CreatedBy               types.String                  `tfsdk:"created_by"`
	UpdateDate              types.String                  `tfsdk:"update_date"`
	UpdatedBy               types.String                  `tfsdk:"updated_by"`
}

const urlProtectionRuleResourceName = "urlProtectionRule"

// NewURLProtectionRuleResource returns a new instance of urlProtectionRuleResource.
func NewURLProtectionRuleResource() resource.Resource { return &urlProtectionRuleResource{} }

func (r *urlProtectionRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_url_protection_rule"
}

func (r *urlProtectionRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL Protection Rule resource.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"url_protection_rule_id": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Unique identifier of the URL protection rule",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the URL protection rule",
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the URL protection rule",
			},
			"bypass_conditions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of bypass conditions for the URL protection rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Type of condition (e.g., RequestHeaderCondition, NetworkListCondition)",
						},
						"names": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of header names for RequestHeaderCondition",
						},
						"name_wildcard": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to use wildcard matching for header names",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of values for the condition",
						},
						"value_case_sensitive": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether the value matching is case sensitive",
						},
						"value_wildcard": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to use wildcard matching for values",
						},
					},
				},
			},
			"max_rate_threshold": schema.Int64Attribute{
				Required:    true,
				Description: "Maximum rate threshold for the URL protection rule",
			},
			"api_definitions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of API definitions associated with the URL protection rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_definition_id": schema.Int64Attribute{
							Required:    true,
							Description: "Unique identifier of the API definition",
						},
						"defined_resources": schema.BoolAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Whether defined resources are included",
						},
						"resource_ids": schema.ListAttribute{
							ElementType: types.Int64Type,
							Computed:    true,
							Optional:    true,
							Description: "List of resource IDs",
						},
						"undefined_resources": schema.BoolAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Whether undefined resources are included",
						},
					},
				},
			},
			"hostname_paths": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Validators:  []validator.List{listvalidator.SizeAtMost(5)},
				Description: "List of hostname and path configurations",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hostname": schema.StringAttribute{
							Required:    true,
							Description: "Hostname for the URL protection rule",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 256),
							},
						},
						"paths": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Validators:  []validator.List{listvalidator.SizeAtMost(5), listvalidator.UniqueValues()},
							Description: "List of paths associated with the hostname",
						},
					},
				},
			},
			"intelligent_load_shedding": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Intelligent load shedding configuration",
				Attributes: map[string]schema.Attribute{
					"hits_per_sec": schema.Int64Attribute{
						Required:    true,
						Description: "Number of hits per second threshold",
					},
					"categories": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "List of categories for intelligent load shedding",
					},
					"custom_criteria": schema.ListNestedAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Custom criteria for intelligent load shedding",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "Type of custom criteria (e.g., CLIENT_LIST)",
								},
								"list_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Required:    true,
									Description: "List of client list IDs",
								},
								"positive_match": schema.BoolAttribute{
									Required:    true,
									Description: "Whether this is a positive match condition",
								},
							},
						},
					},
				},
			},
			"create_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the URL protection rule was created",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who created the URL protection rule",
			},
			"update_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the URL protection rule was last updated",
			},
			"updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who last updated the URL protection rule",
			},
		},
	}
}

func (r *urlProtectionRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
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

func (r *urlProtectionRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data urlProtectionRuleResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.ConfigID.IsUnknown() || data.ConfigID.IsNull() {
		return
	}

	r.validateMaxRateThreshold(ctx, &data, resp)
	r.validateHostnameOrAPIDefinitions(ctx, &data, resp)
	r.validateIntelligentLoadShedding(ctx, &data, resp)
	r.validateBypassConditions(ctx, &data, resp)
}

// --- Helper methods for ValidateConfig ---
func (r *urlProtectionRuleResource) validateMaxRateThreshold(_ context.Context, data *urlProtectionRuleResourceModel, resp *resource.ValidateConfigResponse) {
	if data.MaxRateThreshold.IsUnknown() || data.MaxRateThreshold.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("max_rate_threshold"), "Invalid configuration attribute", "max_rate_threshold is required")
	} else {
		maxRate := data.MaxRateThreshold.ValueInt64()
		if maxRate < 10 {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_rate_threshold"),
				"Invalid configuration attribute",
				"max_rate_threshold must be at least 10.",
			)
		}
	}
}

func (r *urlProtectionRuleResource) validateHostnameOrAPIDefinitions(_ context.Context, data *urlProtectionRuleResourceModel, resp *resource.ValidateConfigResponse) {
	hasHostnamePaths := !data.HostnamePaths.IsNull() && !data.HostnamePaths.IsUnknown()
	hasAPIDefinitions := !data.APIDefinitions.IsNull() && !data.APIDefinitions.IsUnknown()

	if hasHostnamePaths && hasAPIDefinitions {
		resp.Diagnostics.AddError(
			"Invalid Resource",
			"Only one of 'hostname_paths' or 'api_definitions' can be specified, not both.",
		)
	}
	if !hasHostnamePaths && !hasAPIDefinitions {
		resp.Diagnostics.AddError(
			"Invalid Resource",
			"Either 'hostname_paths' or 'api_definitions' must be specified.",
		)
	}
}

func (r *urlProtectionRuleResource) validateIntelligentLoadShedding(ctx context.Context, data *urlProtectionRuleResourceModel, resp *resource.ValidateConfigResponse) {
	ils := data.IntelligentLoadShedding
	if ils == nil {
		return
	}
	if ils.HitsPerSec.IsUnknown() || ils.HitsPerSec.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"hits_per_sec is required when intelligent load shedding is configured.",
		)
	}
	if !ils.HitsPerSec.IsNull() && !ils.HitsPerSec.IsUnknown() &&
		!data.MaxRateThreshold.IsNull() && !data.MaxRateThreshold.IsUnknown() {
		maxRate := data.MaxRateThreshold.ValueInt64()
		hitsPerSec := ils.HitsPerSec.ValueInt64()
		minAllowed := int64(7)
		calculatedMin := int64((float64(maxRate) * 0.25) + 0.5)
		if calculatedMin > minAllowed {
			minAllowed = calculatedMin
		}
		maxAllowed := int64((float64(maxRate) * 0.9) + 0.5)
		if hitsPerSec < minAllowed {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				fmt.Sprintf("hits_per_sec must be at least 25%% of max_rate_threshold (rounded, min 7). Calculated minimum: %d, got: %d", minAllowed, hitsPerSec),
			)
		}
		if hitsPerSec > maxAllowed {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				fmt.Sprintf("hits_per_sec must be less than or equal to 90%% of max_rate_threshold. Calculated maximum: %d, got: %d", maxAllowed, hitsPerSec),
			)
		}
	}
	if !ils.CustomCriteria.IsNull() && !ils.CustomCriteria.IsUnknown() {
		allowedCustomCriteria := map[string]bool{
			"CLIENT_LIST": true,
		}
		var customCriteria []customCriteriaModel
		diags := ils.CustomCriteria.ElementsAs(ctx, &customCriteria, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i, custCriteria := range customCriteria {
			criteriaType := custCriteria.Type.ValueString()
			if !allowedCustomCriteria[criteriaType] {
				resp.Diagnostics.AddAttributeError(
					path.Root("intelligent_load_shedding").AtName("custom_criteria").AtListIndex(i),
					"Invalid Custom Criteria Type",
					fmt.Sprintf("Custom criteria type '%s' is not valid. Allowed values are: CLIENT_LIST", criteriaType),
				)
			}
			if custCriteria.ListIDs.IsNull() || custCriteria.ListIDs.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("intelligent_load_shedding").AtName("custom_criteria").AtListIndex(i).AtName("list_ids"),
					"Invalid Custom Criteria List IDs",
					"list_ids must not be null or unknown.",
				)
			} else {
				var listIDs []string
				diags := custCriteria.ListIDs.ElementsAs(ctx, &listIDs, false)
				resp.Diagnostics.Append(diags...)
				if !resp.Diagnostics.HasError() && len(listIDs) == 0 {
					resp.Diagnostics.AddAttributeError(
						path.Root("intelligent_load_shedding").AtName("custom_criteria").AtListIndex(i).AtName("list_ids"),
						"Invalid Custom Criteria List IDs",
						"list_ids must not be empty.",
					)
				}
			}
		}
	}
	if ils.Categories.IsNull() || ils.Categories.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"intelligent_load_shedding.categories list must not be null or unknown when intelligent load shedding is configured.",
		)
	} else {
		var categories []string
		diags := ils.Categories.ElementsAs(ctx, &categories, false)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() && len(categories) == 0 {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"intelligent_load_shedding.categories list must not be empty when intelligent load shedding is configured.",
			)
		}
	}
}

func (r *urlProtectionRuleResource) validateBypassConditions(ctx context.Context, data *urlProtectionRuleResourceModel, resp *resource.ValidateConfigResponse) {
	if data.BypassConditions.IsNull() || data.BypassConditions.IsUnknown() {
		return
	}
	var bypassConditions []bypassConditionModel
	diags := data.BypassConditions.ElementsAs(ctx, &bypassConditions, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	for i, cond := range bypassConditions {
		if cond.Type.ValueString() == "RequestHeaderCondition" {
			if cond.Names.IsNull() || cond.Names.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("bypass_conditions").AtListIndex(i).AtName("names"),
					"Missing Required Field",
					"'names' is required when type is 'RequestHeaderCondition'",
				)
			} else {
				var names []string
				diags := cond.Names.ElementsAs(ctx, &names, false)
				resp.Diagnostics.Append(diags...)
				if !resp.Diagnostics.HasError() && len(names) == 0 {
					resp.Diagnostics.AddAttributeError(
						path.Root("bypass_conditions").AtListIndex(i).AtName("names"),
						"Empty Field",
						"'names' must not be empty when type is 'RequestHeaderCondition'",
					)
				}
			}
		}
	}
}

func (r *urlProtectionRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating URL Protection Rule resource")

	var plan urlProtectionRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := plan.ConfigID.ValueInt64()

	version, err := getModifiableConfigVersion(ctx, int(configID), urlProtectionRuleResourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	createReq, diags := buildCreateURLProtectionRuleRequest(ctx, &plan, version)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(r.meta)
	out, err := client.CreateURLProtectionRule(ctx, *createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating URL Protection Rule", err.Error())
		return
	}

	// Read back the full resource state from the API to ensure all computed fields are populated
	getReq := appsec.GetURLProtectionRuleRequest{
		ConfigID:            configID,
		ConfigVersion:       int64(version),
		URLProtectionRuleID: out.URLProtectionRuleID,
	}

	fullResponse, err := client.GetURLProtectionRule(ctx, getReq)
	if err != nil {
		resp.Diagnostics.AddError("Error reading URL Protection Rule after creation", err.Error())
		return
	}

	// Map the full response to state using the data source model
	dsModel, diags := createURLProtectionRuleModel(ctx, fullResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start with the plan values for user-provided inputs
	state := plan
	state.URLProtectionID = types.Int64Value(out.URLProtectionRuleID)

	// Update with values from the API response
	state.Name = dsModel.Name
	state.Description = dsModel.Description
	state.MaxRateThreshold = dsModel.MaxRateThreshold
	state.CreateDate = dsModel.CreateDate
	state.CreatedBy = dsModel.CreatedBy
	state.UpdateDate = dsModel.UpdateDate
	state.UpdatedBy = dsModel.UpdatedBy

	// Preserve all nested attributes from plan to avoid inconsistent results
	// The API may return default values for optional fields that were not set
	if !plan.BypassConditions.IsNull() && !plan.BypassConditions.IsUnknown() {
		state.BypassConditions = plan.BypassConditions
	} else {
		state.BypassConditions = dsModel.BypassConditions
	}

	if !plan.APIDefinitions.IsNull() && !plan.APIDefinitions.IsUnknown() {
		state.APIDefinitions = plan.APIDefinitions
	} else {
		state.APIDefinitions = dsModel.APIDefinitions
	}

	if !plan.HostnamePaths.IsNull() && !plan.HostnamePaths.IsUnknown() {
		state.HostnamePaths = plan.HostnamePaths
	} else {
		state.HostnamePaths = dsModel.HostnamePaths
	}

	// Always use the API response for intelligent_load_shedding to ensure computed fields match
	state.IntelligentLoadShedding = dsModel.IntelligentLoadShedding

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *urlProtectionRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading URL Protection Rule resource")

	var state urlProtectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := state.ConfigID.ValueInt64()
	urlProtectionID := state.URLProtectionID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), r.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version", err.Error())
		return
	}

	getReq := appsec.GetURLProtectionRuleRequest{
		ConfigID:            configID,
		ConfigVersion:       int64(version),
		URLProtectionRuleID: urlProtectionID,
	}

	client := inst.Client(r.meta)
	out, err := client.GetURLProtectionRule(ctx, getReq)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetURLProtectionRule'", err.Error())
		return
	}

	if out == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Re-use existing mapping from the data source and keep input fields from current state.
	dsModel, diags := createURLProtectionRuleModel(ctx, out)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState := state
	newState.Name = dsModel.Name
	newState.Description = dsModel.Description
	newState.MaxRateThreshold = dsModel.MaxRateThreshold
	newState.CreateDate = dsModel.CreateDate
	newState.CreatedBy = dsModel.CreatedBy
	newState.UpdateDate = dsModel.UpdateDate
	newState.UpdatedBy = dsModel.UpdatedBy

	// The API may return default values for optional fields that were not set
	if !state.BypassConditions.IsNull() && !state.BypassConditions.IsUnknown() {
		newState.BypassConditions = state.BypassConditions
	} else if !dsModel.BypassConditions.IsNull() && !dsModel.BypassConditions.IsUnknown() {
		newState.BypassConditions = dsModel.BypassConditions
	}

	if !state.APIDefinitions.IsNull() && !state.APIDefinitions.IsUnknown() {
		newState.APIDefinitions = state.APIDefinitions
	} else if !dsModel.APIDefinitions.IsNull() && !dsModel.APIDefinitions.IsUnknown() {
		newState.APIDefinitions = dsModel.APIDefinitions
	}

	// For hostname_paths, preserve the order from state to avoid inconsistent results
	// The API may return items in a different order than what was sent
	if !state.HostnamePaths.IsNull() && !state.HostnamePaths.IsUnknown() {
		// Keep the plan/state order to avoid inconsistent results
		newState.HostnamePaths = state.HostnamePaths
	} else if !dsModel.HostnamePaths.IsNull() && !dsModel.HostnamePaths.IsUnknown() {
		// Only use API order if there was nothing in state (e.g., during import)
		newState.HostnamePaths = dsModel.HostnamePaths
	}

	if state.IntelligentLoadShedding != nil {
		newState.IntelligentLoadShedding = state.IntelligentLoadShedding
	} else if dsModel.IntelligentLoadShedding != nil {
		newState.IntelligentLoadShedding = dsModel.IntelligentLoadShedding
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *urlProtectionRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating URL Protection Rule resource")

	var plan urlProtectionRuleResourceModel
	var state urlProtectionRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := plan.ConfigID.ValueInt64()
	urlProtectionID := state.URLProtectionID.ValueInt64()

	version, err := getModifiableConfigVersion(ctx, int(configID), urlProtectionRuleResourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	updateReq, diags := buildUpdateURLProtectionRuleRequest(ctx, &plan, version, urlProtectionID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(r.meta)
	out, err := client.UpdateURLProtectionRule(ctx, *updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating URL Protection Rule", err.Error())
		return
	}

	// Read back the full resource state from the API to ensure all computed fields are populated
	getReq := appsec.GetURLProtectionRuleRequest{
		ConfigID:            configID,
		ConfigVersion:       int64(version),
		URLProtectionRuleID: out.URLProtectionRuleID,
	}

	fullResponse, err := client.GetURLProtectionRule(ctx, getReq)
	if err != nil {
		resp.Diagnostics.AddError("Error reading URL Protection Rule after update", err.Error())
		return
	}

	// Map the full response to state using the data source model
	dsModel, diags := createURLProtectionRuleModel(ctx, fullResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start with the plan values for user-provided inputs
	newState := plan
	newState.URLProtectionID = types.Int64Value(out.URLProtectionRuleID)

	// Always use the API response for computed fields
	newState.Name = dsModel.Name
	newState.Description = dsModel.Description
	newState.MaxRateThreshold = dsModel.MaxRateThreshold
	newState.CreateDate = dsModel.CreateDate
	newState.CreatedBy = dsModel.CreatedBy
	newState.UpdateDate = dsModel.UpdateDate
	newState.UpdatedBy = dsModel.UpdatedBy

	// Preserve all nested attributes from plan to avoid inconsistent results
	// The API may return default values for optional fields that were not set
	if !plan.BypassConditions.IsNull() && !plan.BypassConditions.IsUnknown() {
		newState.BypassConditions = plan.BypassConditions
	} else {
		newState.BypassConditions = dsModel.BypassConditions
	}

	if !plan.APIDefinitions.IsNull() && !plan.APIDefinitions.IsUnknown() {
		newState.APIDefinitions = plan.APIDefinitions
	} else {
		newState.APIDefinitions = dsModel.APIDefinitions
	}

	if !plan.HostnamePaths.IsNull() && !plan.HostnamePaths.IsUnknown() {
		newState.HostnamePaths = plan.HostnamePaths
	} else {
		newState.HostnamePaths = dsModel.HostnamePaths
	}

	// Always use the API response for intelligent_load_shedding to ensure computed fields match
	newState.IntelligentLoadShedding = dsModel.IntelligentLoadShedding

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *urlProtectionRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting URL Protection Rule resource")

	var state urlProtectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := state.ConfigID.ValueInt64()
	urlProtectionID := state.URLProtectionID.ValueInt64()

	version, err := getModifiableConfigVersion(ctx, int(configID), urlProtectionRuleResourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	deleteReq := appsec.RemoveURLProtectionRuleRequest{
		ConfigID:            configID,
		ConfigVersion:       int64(version),
		URLProtectionRuleID: urlProtectionID,
	}

	client := inst.Client(r.meta)
	if err := client.RemoveURLProtectionRule(ctx, deleteReq); err != nil {
		resp.Diagnostics.AddError("Error deleting URL Protection Rule", err.Error())
		return
	}
}

func (r *urlProtectionRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing URL Protection Rule resource")

	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID", "Import ID must be in the form '<config_id>:<url_protection_id>'")
		return
	}

	configID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Config ID", fmt.Sprintf("%q is not a valid int64", parts[0]))
		return
	}

	urlProtectionID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid URL Protection ID", fmt.Sprintf("%q is not a valid int64", parts[1]))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("config_id"), configID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("url_protection_rule_id"), urlProtectionID)...)

}

func buildCreateURLProtectionRuleRequest(ctx context.Context, plan *urlProtectionRuleResourceModel, version int) (*appsec.CreateURLProtectionRuleRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	body, bd := buildURLProtectionRuleRequestBody(ctx, plan)
	diags.Append(bd...)
	if diags.HasError() {
		return nil, diags
	}

	return &appsec.CreateURLProtectionRuleRequest{
		ConfigID:      plan.ConfigID.ValueInt64(),
		ConfigVersion: int64(version),
		Body:          body,
	}, diags
}

func buildUpdateURLProtectionRuleRequest(ctx context.Context, plan *urlProtectionRuleResourceModel, version int, urlProtectionID int64) (*appsec.UpdateURLProtectionRuleRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	body, bd := buildURLProtectionRuleRequestBody(ctx, plan)
	diags.Append(bd...)
	if diags.HasError() {
		return nil, diags
	}

	return &appsec.UpdateURLProtectionRuleRequest{
		ConfigID:            plan.ConfigID.ValueInt64(),
		ConfigVersion:       int64(version),
		URLProtectionRuleID: urlProtectionID,
		Body:                body,
	}, diags
}

func buildURLProtectionRuleRequestBody(ctx context.Context, plan *urlProtectionRuleResourceModel) (appsec.URLProtectionRuleRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := appsec.URLProtectionRuleRequestBody{
		Name:             plan.Name.ValueString(),
		MaxRateThreshold: plan.MaxRateThreshold.ValueInt64(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		d := plan.Description.ValueString()
		body.Description = &d
	}

	// Build bypass conditions
	if err := buildBypassConditions(ctx, plan, &body, &diags); err != nil {
		return body, diags
	}

	// Build API definitions
	if err := buildAPIDefinitions(ctx, plan, &body, &diags); err != nil {
		return body, diags
	}

	// Build hostname paths
	if err := buildHostnamePaths(ctx, plan, &body, &diags); err != nil {
		return body, diags
	}

	// Build intelligent load shedding
	if err := buildIntelligentLoadShedding(ctx, plan, &body, &diags); err != nil {
		return body, diags
	}

	tflog.Debug(ctx, "URL Protection Rule request body", map[string]interface{}{
		"body": fmt.Sprintf("%+v", body),
	})
	return body, diags
}

// buildBypassConditions processes bypass conditions from the plan and adds them to the request body
func buildBypassConditions(ctx context.Context, plan *urlProtectionRuleResourceModel, body *appsec.URLProtectionRuleRequestBody, diags *diag.Diagnostics) error {
	if plan.BypassConditions.IsNull() || plan.BypassConditions.IsUnknown() {
		return nil
	}

	var bypassConditions []bypassConditionModel
	diags.Append(plan.BypassConditions.ElementsAs(ctx, &bypassConditions, false)...)
	if diags.HasError() {
		return fmt.Errorf("failed to parse bypass conditions")
	}

	atomic := make([]appsec.AtomicCondition, 0, len(bypassConditions))
	for _, cond := range bypassConditions {
		ac := appsec.AtomicCondition{
			Type: cond.Type.ValueString(),
		}

		if !cond.NameWildcard.IsNull() && !cond.NameWildcard.IsUnknown() {
			v := cond.NameWildcard.ValueBool()
			ac.NameWildcard = &v
		}
		if !cond.ValueCaseSensitive.IsNull() && !cond.ValueCaseSensitive.IsUnknown() {
			v := cond.ValueCaseSensitive.ValueBool()
			ac.ValueCase = &v
		}
		if !cond.ValueWildcard.IsNull() && !cond.ValueWildcard.IsUnknown() {
			v := cond.ValueWildcard.ValueBool()
			ac.ValueWildcard = &v
		}

		if !cond.Names.IsNull() && !cond.Names.IsUnknown() {
			var names []string
			diags.Append(cond.Names.ElementsAs(ctx, &names, false)...)
			if diags.HasError() {
				return fmt.Errorf("failed to parse bypass_conditions names")
			}
			ac.Names = names
		}
		if !cond.Values.IsNull() && !cond.Values.IsUnknown() {
			var values []string
			diags.Append(cond.Values.ElementsAs(ctx, &values, false)...)
			if diags.HasError() {
				return fmt.Errorf("failed to parse bypass_conditions values")
			}
			ac.Values = values
		}

		atomic = append(atomic, ac)
	}

	body.BypassCondition = &appsec.BypassCondition{AtomicConditions: atomic}
	return nil
}

// buildAPIDefinitions processes API definitions from the plan and adds them to the request body
func buildAPIDefinitions(ctx context.Context, plan *urlProtectionRuleResourceModel, body *appsec.URLProtectionRuleRequestBody, diags *diag.Diagnostics) error {
	if plan.APIDefinitions.IsNull() || plan.APIDefinitions.IsUnknown() {
		return nil
	}

	var apiDefinitions []apiDefinitionModel
	diags.Append(plan.APIDefinitions.ElementsAs(ctx, &apiDefinitions, false)...)
	if diags.HasError() {
		return fmt.Errorf("failed to parse API definitions")
	}

	defs := make([]appsec.APIDefinition, 0, len(apiDefinitions))
	for _, d := range apiDefinitions {
		def := appsec.APIDefinition{
			APIDefinitionID:    d.APIDefinitionID.ValueInt64(),
			DefinedResources:   d.DefinedResources.ValueBool(),
			UndefinedResources: d.UndefinedResources.ValueBool(),
		}
		if !d.ResourceIDs.IsNull() && !d.ResourceIDs.IsUnknown() {
			var ids []int64
			diags.Append(d.ResourceIDs.ElementsAs(ctx, &ids, false)...)
			if diags.HasError() {
				return fmt.Errorf("failed to parse api_definitions ids")
			}
			def.ResourceIDs = ids
		}
		defs = append(defs, def)
	}
	body.APIDefinitions = defs
	return nil
}

// buildHostnamePaths processes hostname paths from the plan and adds them to the request body
func buildHostnamePaths(ctx context.Context, plan *urlProtectionRuleResourceModel, body *appsec.URLProtectionRuleRequestBody, diags *diag.Diagnostics) error {

	if plan.HostnamePaths.IsNull() || plan.HostnamePaths.IsUnknown() {
		return nil
	}

	var hostnamePaths []hostnamePathModel
	diags.Append(plan.HostnamePaths.ElementsAs(ctx, &hostnamePaths, false)...)
	if diags.HasError() {
		return fmt.Errorf("failed to parse hostname paths")
	}

	hp := make([]appsec.HostnamePath, 0, len(hostnamePaths))
	for _, h := range hostnamePaths {
		p := appsec.HostnamePath{Hostname: h.Hostname.ValueString()}
		if !h.Paths.IsNull() && !h.Paths.IsUnknown() {
			var paths []string
			diags.Append(h.Paths.ElementsAs(ctx, &paths, false)...)
			if diags.HasError() {
				return fmt.Errorf("failed to parse hostname_paths paths")
			}
			p.Paths = paths
		}
		hp = append(hp, p)
	}
	body.HostnamePaths = hp

	// Set protection type based on the number of hostname paths and paths
	if len(hostnamePaths) > 1 || len(hostnamePaths[0].Paths.Elements()) > 1 {
		multiple := "MULTIPLE"
		body.ProtectionType = &multiple
	} else {
		single := "SINGLE"
		body.ProtectionType = &single
	}
	return nil
}

// buildIntelligentLoadShedding processes intelligent load shedding configuration from the plan and adds it to the request body
func buildIntelligentLoadShedding(ctx context.Context, plan *urlProtectionRuleResourceModel, body *appsec.URLProtectionRuleRequestBody, diags *diag.Diagnostics) error {
	if plan.IntelligentLoadShedding == nil {
		return nil
	}

	body.IntelligentLoadShedding = true
	if !plan.IntelligentLoadShedding.HitsPerSec.IsNull() && !plan.IntelligentLoadShedding.HitsPerSec.IsUnknown() {
		v := plan.IntelligentLoadShedding.HitsPerSec.ValueInt64()
		body.SheddingThresholdHitsPerSec = &v
	}

	var loadSheddingCategories []appsec.Category

	// Process standard categories
	if !plan.IntelligentLoadShedding.Categories.IsNull() && !plan.IntelligentLoadShedding.Categories.IsUnknown() {
		var categories []string
		diags.Append(plan.IntelligentLoadShedding.Categories.ElementsAs(ctx, &categories, false)...)
		for _, category := range categories {
			loadSheddingCategories = append(loadSheddingCategories, appsec.Category{Type: category})
		}
	}

	// Process custom criteria
	if err := buildCustomCriteria(ctx, plan, &loadSheddingCategories, diags); err != nil {
		return err
	}

	body.Categories = loadSheddingCategories
	return nil
}

// buildCustomCriteria processes custom criteria for intelligent load shedding
func buildCustomCriteria(ctx context.Context, plan *urlProtectionRuleResourceModel, loadSheddingCategories *[]appsec.Category, diags *diag.Diagnostics) error {
	if plan.IntelligentLoadShedding.CustomCriteria.IsNull() || plan.IntelligentLoadShedding.CustomCriteria.IsUnknown() {
		return nil
	}

	var customCriteria []customCriteriaModel
	diags.Append(plan.IntelligentLoadShedding.CustomCriteria.ElementsAs(ctx, &customCriteria, false)...)
	if diags.HasError() {
		return fmt.Errorf("failed to parse custom criteria")
	}

	for _, customCriteria := range customCriteria {
		var listIDs []string
		diags.Append(customCriteria.ListIDs.ElementsAs(ctx, &listIDs, false)...)
		if diags.HasError() {
			return fmt.Errorf("failed to parse list IDs")
		}
		positiveMatch := customCriteria.PositiveMatch.ValueBool()

		*loadSheddingCategories = append(*loadSheddingCategories, appsec.Category{
			Type:          customCriteria.Type.ValueString(),
			ListIDs:       listIDs,
			PositiveMatch: &positiveMatch,
		})
	}
	return nil
}
