package appsec

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &urlProtectionActionResource{}
	_ resource.ResourceWithConfigure      = &urlProtectionActionResource{}
	_ resource.ResourceWithImportState    = &urlProtectionActionResource{}
	_ resource.ResourceWithValidateConfig = &urlProtectionActionResource{}
)

// urlProtectionActionResource represents akamai_appsec_url_protection_action resource
type urlProtectionActionResource struct {
	meta.Resource
}

// urlProtectionActionResourceModel is a model for akamai_appsec_url_protection_action resource
type urlProtectionActionResourceModel struct {
	ConfigID               types.Int64  `tfsdk:"config_id"`
	SecurityPolicyID       types.String `tfsdk:"security_policy_id"`
	URLProtectionPolicyID  types.Int64  `tfsdk:"url_protection_policy_id"`
	MaxRateThresholdAction types.String `tfsdk:"max_rate_threshold_action"`
	LoadSheddingAction     types.String `tfsdk:"load_shedding_action"`
}

const (
	urlProtectionActionResourceName = "akamai_appsec_url_protection_action"
)

// NewURLProtectionActionResource returns a new URL Protection Action resource.
func NewURLProtectionActionResource() resource.Resource {
	return &urlProtectionActionResource{}
}

// Metadata implements resource's Metadata method.
func (r *urlProtectionActionResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = urlProtectionActionResourceName
}

// Schema implements resource's Schema method.
func (r *urlProtectionActionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage URL Protection Actions for Application Security",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy",
				Validators: []validator.String{
					validators.NotEmptyString(),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"url_protection_policy_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the URL protection policy",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"max_rate_threshold_action": schema.StringAttribute{
				Required:    true,
				Description: "Action to apply when max rate threshold is exceeded (e.g., alert, deny, none, deny_custom_{custom_deny_id}, challenge_{challenge_id})",
				Validators: []validator.String{
					validators.NotEmptyString(),
					stringvalidator.Any(
						stringvalidator.OneOf("alert", "deny", "none"),
						stringvalidator.RegexMatches(regexp.MustCompile("^deny_custom_.+"), "must start with 'deny_custom_'"),
						stringvalidator.RegexMatches(regexp.MustCompile("^challenge_.+"), "must start with 'challenge_'"),
					),
				},
			},
			"load_shedding_action": schema.StringAttribute{
				Optional:    true,
				Description: "Load shedding action to apply (e.g., alert, deny, none, deny_custom_{custom_deny_id}, challenge_{challenge_id})",
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.OneOf("alert", "deny", "none"),
						stringvalidator.RegexMatches(regexp.MustCompile("^deny_custom_.+"), "must start with 'deny_custom_'"),
						stringvalidator.RegexMatches(regexp.MustCompile("^challenge_.+"), "must start with 'challenge_'"),
					),
				},
			},
		},
	}
}

// ValidateConfig implements resource's ValidateConfig method.
func (r *urlProtectionActionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

	tflog.Debug(ctx, "Validating URL Protection Action resource configuration")

	if r.Client == nil {
		return
	}

	var config urlProtectionActionResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip validation if any required fields are unknown (during plan with computed values)
	if config.ConfigID.IsUnknown() || config.URLProtectionPolicyID.IsUnknown() {
		return
	}

	client := r.Client.GetAPPSEC()

	version, err := getLatestConfigVersion(ctx, int(config.ConfigID.ValueInt64()), r.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	// Get URL Protection Policy to check intelligentLoadShedding flag
	getPolicyRequest := appsec.GetURLProtectionPolicyRequest{
		ConfigID:              config.ConfigID.ValueInt64(),
		ConfigVersion:         int64(version), // Use the modifiable version
		URLProtectionPolicyID: config.URLProtectionPolicyID.ValueInt64(),
	}

	policyResponse, err := client.GetURLProtectionPolicy(ctx, getPolicyRequest)
	if err != nil {
		// If we can't fetch the policy, we'll throw error
		resp.Diagnostics.AddError(
			"Failed to prepare URL protection action",
			fmt.Sprintf("URL Protection Policy could not be found: %v", err),
		)
		return
	}

	// Check if intelligentLoadShedding is enabled and load_shedding_action is required
	if policyResponse.IntelligentLoadShedding {
		if config.LoadSheddingAction.IsNull() || config.LoadSheddingAction.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("load_shedding_action"),
				"Missing required attribute",
				"load_shedding_action is required when intelligent load shedding is enabled for this URL protection policy",
			)
		}
	} else {
		if !config.LoadSheddingAction.IsNull() && config.LoadSheddingAction.ValueString() != "" && config.LoadSheddingAction.ValueString() != "none" {
			resp.Diagnostics.AddAttributeError(
				path.Root("load_shedding_action"),
				"Attribute not allowed",
				"load_shedding_action is not allowed when intelligent load shedding is disabled for this URL protection policy",
			)
		}
	}
}

// Create implements resource's Create method.
func (r *urlProtectionActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating URL Protection Action Resource")

	var data urlProtectionActionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	version, err := getModifiableConfigVersion(ctx, int(configID), "urlProtectionAction", r.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	client := r.Client.GetAPPSEC()

	loadSheddingAction := data.LoadSheddingAction.ValueString()
	if data.LoadSheddingAction.IsNull() || data.LoadSheddingAction.ValueString() == "" {
		loadSheddingAction = "none"
	}

	updateRequest := appsec.UpdateURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              data.SecurityPolicyID.ValueString(),
		URLProtectionPolicyID: data.URLProtectionPolicyID.ValueInt64(),
		Body: appsec.URLProtectionPolicyActions{
			MaxRateThresholdAction: data.MaxRateThresholdAction.ValueString(),
			LoadSheddingAction:     loadSheddingAction,
		},
	}

	_, err = client.UpdateURLProtectionPolicyActions(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create URL protection action", err.Error())
		return
	}

	// Read back the created resource to populate computed fields
	readRequest := appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              data.SecurityPolicyID.ValueString(),
		URLProtectionPolicyID: data.URLProtectionPolicyID.ValueInt64(),
	}

	actionResponse, err := client.GetURLProtectionPolicyActions(ctx, readRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read URL protection action after creation", err.Error())
		return
	}

	data.MaxRateThresholdAction = types.StringValue(actionResponse.MaxRateThresholdAction)
	// Handle LoadSheddingAction: only update it if loadSheddingAction is not null or empty in request
	if !data.LoadSheddingAction.IsNull() && loadSheddingAction != "" {
		data.LoadSheddingAction = types.StringValue(actionResponse.LoadSheddingAction)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource's Read method.
func (r *urlProtectionActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading URL Protection Action Resource")

	var data urlProtectionActionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	version, err := getLatestConfigVersion(ctx, int(configID), r.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read config version", err.Error())
		return
	}

	readRequest := appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              data.SecurityPolicyID.ValueString(),
		URLProtectionPolicyID: data.URLProtectionPolicyID.ValueInt64(),
	}

	client := r.Client.GetAPPSEC()
	actionResponse, err := client.GetURLProtectionPolicyActions(ctx, readRequest)
	if err != nil {
		// If the URL Protection Policy or its actions are not found, remove the resource from state. May happen if url protection policy is not present in latest config version.
		if strings.Contains(err.Error(), "incorrect URL Protection Policy ID") ||
			strings.Contains(err.Error(), "no actions found for the specified policy") {
			tflog.Warn(ctx, "URL Protection Policy not found or has no actions, removing from state", map[string]interface{}{
				"config_id":                configID,
				"security_policy_id":       data.SecurityPolicyID.ValueString(),
				"url_protection_policy_id": data.URLProtectionPolicyID.ValueInt64(),
			})
			resp.Diagnostics.AddWarning("Could not find URL Protection Policy in config version", err.Error())
			// Remove resource from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read URL protection action", err.Error())
		return
	}

	data.MaxRateThresholdAction = types.StringValue(actionResponse.MaxRateThresholdAction)

	data.LoadSheddingAction = types.StringValue(actionResponse.LoadSheddingAction)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource's Update method.
func (r *urlProtectionActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating URL Protection Action Resource")

	var data urlProtectionActionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	version, err := getModifiableConfigVersion(ctx, int(configID), "urlProtectionAction", r.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	client := r.Client.GetAPPSEC()

	loadSheddingAction := data.LoadSheddingAction.ValueString()
	if data.LoadSheddingAction.IsNull() || data.LoadSheddingAction.ValueString() == "" {
		loadSheddingAction = "none"
	}

	updateRequest := appsec.UpdateURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              data.SecurityPolicyID.ValueString(),
		URLProtectionPolicyID: data.URLProtectionPolicyID.ValueInt64(),
		Body: appsec.URLProtectionPolicyActions{
			MaxRateThresholdAction: data.MaxRateThresholdAction.ValueString(),
			LoadSheddingAction:     loadSheddingAction,
		},
	}

	if err != nil {
		if strings.Contains(err.Error(), "load_shedding_action is required") {
			resp.Diagnostics.AddAttributeError(
				path.Root("load_shedding_action"),
				"Missing required attribute",
				err.Error(),
			)
		} else {
			resp.Diagnostics.AddError("Failed to prepare URL protection action", err.Error())
		}
		return
	}

	_, err = client.UpdateURLProtectionPolicyActions(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update URL protection action", err.Error())
		return
	}

	// Read back the updated resource to populate state correctly
	readRequest := appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              data.SecurityPolicyID.ValueString(),
		URLProtectionPolicyID: data.URLProtectionPolicyID.ValueInt64(),
	}

	actionResponse, err := client.GetURLProtectionPolicyActions(ctx, readRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read URL protection action after update", err.Error())
		return
	}

	data.MaxRateThresholdAction = types.StringValue(actionResponse.MaxRateThresholdAction)
	// Handle LoadSheddingAction: only update it if loadSheddingAction is not null or empty in request
	if !data.LoadSheddingAction.IsNull() && data.LoadSheddingAction.ValueString() != "" {
		data.LoadSheddingAction = types.StringValue(actionResponse.LoadSheddingAction)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource's Delete method.
func (r *urlProtectionActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, "Deleting URL Protection Action Resource")

	var data urlProtectionActionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	version, err := getModifiableConfigVersion(ctx, int(configID), "urlProtectionAction", r.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	// Delete sets action to "none"
	deleteRequest := appsec.UpdateURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              data.SecurityPolicyID.ValueString(),
		URLProtectionPolicyID: data.URLProtectionPolicyID.ValueInt64(),
		Body: appsec.URLProtectionPolicyActions{
			MaxRateThresholdAction: "none",
			LoadSheddingAction:     "none",
		},
	}

	client := r.Client.GetAPPSEC()
	_, err = client.UpdateURLProtectionPolicyActions(ctx, deleteRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete URL protection action", err.Error())
		return
	}
}

// ImportState implements resource's ImportState method.
func (r *urlProtectionActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing URL Protection Action resource")

	parts := strings.Split(req.ID, ":")

	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ID '%s' incorrectly formatted: should be 'CONFIG_ID:SECURITY_POLICY_ID:URL_PROTECTION_POLICY_ID'", req.ID),
			"",
		)
		return
	}

	configID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid configuration id '%v'", parts[0]), "")
		return
	}

	securityPolicyID := parts[1]
	if securityPolicyID == "" {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid security policy id '%v'", parts[1]), "")
		return
	}

	urlProtectionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid url protection id '%v'", parts[2]), "")
		return
	}

	// id is not in the schema, do not set it in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("config_id"), configID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("security_policy_id"), securityPolicyID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("url_protection_policy_id"), urlProtectionID)...)
}
