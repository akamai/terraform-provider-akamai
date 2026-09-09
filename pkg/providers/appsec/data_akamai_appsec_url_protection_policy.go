package appsec

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	urlProtectionPolicyDataSource struct {
		meta.DataSource
	}

	// urlProtectionPolicyDataSourceModel maps the API response to a structure that can be used in the Terraform schema
	urlProtectionPolicyDataSourceModel struct {
		ConfigID                types.Int64                   `tfsdk:"config_id"`
		URLProtectionID         types.Int64                   `tfsdk:"url_protection_policy_id"`
		Name                    types.String                  `tfsdk:"name"`
		Description             types.String                  `tfsdk:"description"`
		BypassConditions        types.List                    `tfsdk:"bypass_conditions"`
		MaxRateThreshold        types.Int64                   `tfsdk:"max_rate_threshold"`
		APIDefinitions          types.List                    `tfsdk:"api_definitions"`
		HostnamePaths           types.List                    `tfsdk:"hostname_paths"`
		IntelligentLoadShedding *intelligentLoadSheddingModel `tfsdk:"intelligent_load_shedding"`
		Used                    types.Bool                    `tfsdk:"used"`
		CreateDate              types.String                  `tfsdk:"create_date"`
		CreatedBy               types.String                  `tfsdk:"created_by"`
		UpdateDate              types.String                  `tfsdk:"update_date"`
		UpdatedBy               types.String                  `tfsdk:"updated_by"`
	}

	bypassConditionModel struct {
		Type               types.String `tfsdk:"type"`
		Names              types.List   `tfsdk:"names"`
		NameWildcard       types.Bool   `tfsdk:"name_wildcard"`
		Values             types.List   `tfsdk:"values"`
		ValueCaseSensitive types.Bool   `tfsdk:"value_case_sensitive"`
		ValueWildcard      types.Bool   `tfsdk:"value_wildcard"`
	}

	apiDefinitionModel struct {
		APIDefinitionID    types.Int64 `tfsdk:"api_definition_id"`
		DefinedResources   types.Bool  `tfsdk:"defined_resources"`
		ResourceIDs        types.List  `tfsdk:"resource_ids"`
		UndefinedResources types.Bool  `tfsdk:"undefined_resources"`
	}

	hostnamePathModel struct {
		Hostname types.String `tfsdk:"hostname"`
		Paths    types.List   `tfsdk:"paths"`
	}

	intelligentLoadSheddingModel struct {
		HitsPerSec     types.Int64 `tfsdk:"hits_per_sec"`
		Categories     types.List  `tfsdk:"categories"`
		CustomCriteria types.List  `tfsdk:"custom_criteria"`
	}

	customCriteriaModel struct {
		Type          types.String `tfsdk:"type"`
		ListIDs       types.List   `tfsdk:"list_ids"`
		PositiveMatch types.Bool   `tfsdk:"positive_match"`
	}
)

var (
	_ datasource.DataSource              = &urlProtectionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &urlProtectionPolicyDataSource{}
)

// NewURLProtectionPolicyDataSource returns a new URL protection policy data source.
func NewURLProtectionPolicyDataSource() datasource.DataSource {
	return &urlProtectionPolicyDataSource{}
}

// Metadata configures data source's meta information.
func (d *urlProtectionPolicyDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_url_protection_policy"
}

// Schema is used to define data source's terraform schema.
func (d *urlProtectionPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL Protection Policy data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"url_protection_policy_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the URL protection policy",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the URL protection policy",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the URL protection",
			},
			"bypass_conditions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of bypass conditions for the URL protection policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of condition (e.g., RequestHeaderCondition, NetworkListCondition)",
						},
						"names": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "List of header names for RequestHeaderCondition",
						},
						"name_wildcard": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to use wildcard matching for header names",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "List of values for the condition",
						},
						"value_case_sensitive": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the value matching is case sensitive",
						},
						"value_wildcard": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to use wildcard matching for values",
						},
					},
				},
			},
			"max_rate_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: "Maximum rate threshold for the URL protection",
			},
			"api_definitions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of API definitions associated with the URL protection policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_definition_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the API definition",
						},
						"defined_resources": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether defined resources are included",
						},
						"resource_ids": schema.ListAttribute{
							ElementType: types.Int64Type,
							Computed:    true,
							Description: "List of resource IDs",
						},
						"undefined_resources": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether undefined resources are included",
						},
					},
				},
			},
			"hostname_paths": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of hostname and path configurations",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hostname": schema.StringAttribute{
							Computed:    true,
							Description: "Hostname for the URL protection",
						},
						"paths": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "List of paths associated with the hostname",
						},
					},
				},
			},
			"intelligent_load_shedding": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Intelligent load shedding configuration",
				Attributes: map[string]schema.Attribute{
					"hits_per_sec": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of hits per second threshold",
					},
					"categories": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "List of categories for intelligent load shedding",
					},
					"custom_criteria": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Custom criteria for intelligent load shedding",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "Type of custom criteria (e.g., CLIENT_LIST)",
								},
								"list_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "List of client list IDs",
								},
								"positive_match": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether this is a positive match condition",
								},
							},
						},
					},
				},
			},
			"used": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether you're currently using the URL protection policy",
			},
			"create_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the URL protection policy was created",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who created the URL protection policy",
			},
			"update_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the URL protection policy was last updated",
			},
			"updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who last updated the URL protection policy",
			},
		},
	}
}

func (d *urlProtectionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLProtectionPolicyDataSource Read")

	var data urlProtectionPolicyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetAPPSEC()
	configID := data.ConfigID.ValueInt64()
	urlProtectionID := data.URLProtectionID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), d.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("invalid config version", err.Error())
		return
	}

	getURLProtectionPolicyRequest := appsec.GetURLProtectionPolicyRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		URLProtectionPolicyID: urlProtectionID,
	}

	urlProtectionPolicy, err := client.GetURLProtectionPolicy(ctx, getURLProtectionPolicyRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetURLProtectionPolicy'", err.Error())
		return
	}

	if urlProtectionPolicy == nil {
		resp.Diagnostics.AddError("URL Protection Policy not found", fmt.Sprintf("URL Protection Policy with ID %d not found in config %d version %d", urlProtectionID, configID, version))
		return
	}

	urlProtectionPolicyModel, diags := createURLProtectionPolicyModel(ctx, urlProtectionPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &urlProtectionPolicyModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func createURLProtectionPolicyModel(ctx context.Context, response *appsec.GetURLProtectionPolicyResponse) (urlProtectionPolicyDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model urlProtectionPolicyDataSourceModel

	if response == nil {
		diags.AddError("Invalid response", "GetURLProtectionPolicyResponse is nil")
		return model, diags
	}

	// Set basic fields
	model.ConfigID = types.Int64Value(response.ConfigID)
	model.URLProtectionID = types.Int64Value(response.URLProtectionPolicyID)
	model.Name = types.StringValue(response.Name)
	if response.Description != nil && *response.Description != "" {
		model.Description = types.StringValue(*response.Description)
	}
	model.MaxRateThreshold = types.Int64Value(response.MaxRateThreshold)
	model.Used = types.BoolValue(response.Used)
	model.CreateDate = types.StringValue(response.CreateDate)
	model.CreatedBy = types.StringValue(response.CreatedBy)
	model.UpdateDate = types.StringValue(response.UpdateDate)
	model.UpdatedBy = types.StringValue(response.UpdatedBy)

	// Convert bypass conditions
	if response.BypassCondition != nil && response.BypassCondition.AtomicConditions != nil && len(response.BypassCondition.AtomicConditions) > 0 {
		diags = populateBypassConditions(ctx, &model, response, diags)
	} else {
		// Initialize with empty list if no bypass conditions
		model.BypassConditions = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":                 types.StringType,
				"names":                types.ListType{ElemType: types.StringType},
				"name_wildcard":        types.BoolType,
				"values":               types.ListType{ElemType: types.StringType},
				"value_case_sensitive": types.BoolType,
				"value_wildcard":       types.BoolType,
			},
		})
	}

	if response.APIDefinitions != nil && response.HostnamePaths != nil {
		diags.AddError("Invalid response", "URLProtectionPolicy cannot have both APIDefinitions and HostnamePaths defined at the same time")
		return model, diags
	} else if len(response.APIDefinitions) > 0 {
		diags = populateAPIDefinitions(ctx, &model, response, diags)
		// Initialize HostnamePaths as null since we have APIDefinitions
		model.HostnamePaths = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"hostname": types.StringType,
				"paths":    types.ListType{ElemType: types.StringType},
			},
		})
	} else if len(response.HostnamePaths) > 0 {
		diags = populateHostnamePaths(ctx, &model, response, diags)
		// Initialize APIDefinitions as null since we have HostnamePaths
		model.APIDefinitions = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"api_definition_id":   types.Int64Type,
				"defined_resources":   types.BoolType,
				"resource_ids":        types.ListType{ElemType: types.Int64Type},
				"undefined_resources": types.BoolType,
			},
		})
	} else {
		diags.AddError("Invalid response", "There is some problem with the API response, either APIDefinitions or HostnamePaths should be defined and not both. Please report this issue to the provider developers.")
		return model, diags
	}

	// Convert intelligent load shedding
	if response.IntelligentLoadShedding {
		ilsModel := intelligentLoadSheddingModel{}
		if response.SheddingThresholdHitsPerSec != nil {
			ilsModel.HitsPerSec = types.Int64Value(*response.SheddingThresholdHitsPerSec)
		} else {
			diags.AddError("Invalid response", "HitsPerSec value is missing for intelligent load shedding configuration")
			return model, diags
		}

		// Parse categories and custom criteria
		if len(response.Categories) > 0 {
			diags = populateCategoriesAndCustomCriteria(ctx, &ilsModel, response, diags)
		} else {
			// Set empty values if no categories
			ilsModel.Categories = types.ListNull(types.StringType)
			ilsModel.CustomCriteria = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type":           types.StringType,
					"list_ids":       types.ListType{ElemType: types.StringType},
					"positive_match": types.BoolType,
				},
			})
		}

		model.IntelligentLoadShedding = &ilsModel

	} else {
		model.IntelligentLoadShedding = nil
	}

	return model, diags
}

func populateBypassConditions(ctx context.Context, model *urlProtectionPolicyDataSourceModel, response *appsec.GetURLProtectionPolicyResponse, diags diag.Diagnostics) diag.Diagnostics {
	bypassConditions := make([]bypassConditionModel, 0, len(response.BypassCondition.AtomicConditions))
	for _, atomicCond := range response.BypassCondition.AtomicConditions {
		var bypassCond bypassConditionModel
		bypassCond.Type = types.StringValue(atomicCond.Type)
		if atomicCond.NameWildcard != nil {
			bypassCond.NameWildcard = types.BoolValue(*atomicCond.NameWildcard)
		}
		if atomicCond.ValueCase != nil {
			bypassCond.ValueCaseSensitive = types.BoolValue(*atomicCond.ValueCase)
		}
		if atomicCond.ValueWildcard != nil {
			bypassCond.ValueWildcard = types.BoolValue(*atomicCond.ValueWildcard)
		}

		// Convert Names slice to types.List
		if len(atomicCond.Names) > 0 {
			namesList, diagsTemp := types.ListValueFrom(ctx, types.StringType, atomicCond.Names)
			diags.Append(diagsTemp...)
			bypassCond.Names = namesList
		} else {
			bypassCond.Names = types.ListNull(types.StringType)
		}

		// Convert Values slice to types.List
		if len(atomicCond.Values) > 0 {
			valuesList, diagsTemp := types.ListValueFrom(ctx, types.StringType, atomicCond.Values)
			diags.Append(diagsTemp...)
			bypassCond.Values = valuesList
		} else {
			bypassCond.Values = types.ListNull(types.StringType)
		}

		bypassConditions = append(bypassConditions, bypassCond)
	}

	// Convert slice to types.List
	bypassConditionsList, diagsTemp := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":                 types.StringType,
			"names":                types.ListType{ElemType: types.StringType},
			"name_wildcard":        types.BoolType,
			"values":               types.ListType{ElemType: types.StringType},
			"value_case_sensitive": types.BoolType,
			"value_wildcard":       types.BoolType,
		},
	}, bypassConditions)
	diags.Append(diagsTemp...)
	model.BypassConditions = bypassConditionsList
	return diags
}

func populateAPIDefinitions(ctx context.Context, model *urlProtectionPolicyDataSourceModel, response *appsec.GetURLProtectionPolicyResponse, diags diag.Diagnostics) diag.Diagnostics {
	apiDefinitions := make([]apiDefinitionModel, 0, len(response.APIDefinitions))
	for _, apiDef := range response.APIDefinitions {
		var apiDefModel apiDefinitionModel
		apiDefModel.APIDefinitionID = types.Int64Value(apiDef.APIDefinitionID)
		apiDefModel.DefinedResources = types.BoolValue(apiDef.DefinedResources)
		apiDefModel.UndefinedResources = types.BoolValue(apiDef.UndefinedResources)

		// Convert ResourceIDs slice to types.List
		if len(apiDef.ResourceIDs) > 0 {
			resourceIDsList, diagsTemp := types.ListValueFrom(ctx, types.Int64Type, apiDef.ResourceIDs)
			diags.Append(diagsTemp...)
			apiDefModel.ResourceIDs = resourceIDsList
		} else {
			var diagsTemp diag.Diagnostics
			apiDefModel.ResourceIDs, diagsTemp = types.ListValueFrom(ctx, types.Int64Type, []int64{})
			if diagsTemp.HasError() {
				diags.Append(diagsTemp...)
				return diags
			}
		}

		apiDefinitions = append(apiDefinitions, apiDefModel)
	}

	// Convert slice to types.List
	apiDefinitionsList, diagsTemp := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"api_definition_id":   types.Int64Type,
			"defined_resources":   types.BoolType,
			"resource_ids":        types.ListType{ElemType: types.Int64Type},
			"undefined_resources": types.BoolType,
		},
	}, apiDefinitions)
	diags.Append(diagsTemp...)
	model.APIDefinitions = apiDefinitionsList
	return diags
}

func populateHostnamePaths(ctx context.Context, model *urlProtectionPolicyDataSourceModel, response *appsec.GetURLProtectionPolicyResponse, diags diag.Diagnostics) diag.Diagnostics {
	hostnamePaths := make([]hostnamePathModel, 0, len(response.HostnamePaths))
	for _, hostPath := range response.HostnamePaths {
		var hostPathModel hostnamePathModel
		hostPathModel.Hostname = types.StringValue(hostPath.Hostname)

		// Convert Paths slice to types.List
		if len(hostPath.Paths) > 0 {
			pathsList, diagsTemp := types.ListValueFrom(ctx, types.StringType, hostPath.Paths)
			diags.Append(diagsTemp...)
			hostPathModel.Paths = pathsList
		} else {
			hostPathModel.Paths = types.ListNull(types.StringType)
		}

		hostnamePaths = append(hostnamePaths, hostPathModel)
	}

	// Convert slice to types.List
	hostnamePathsList, diagsTemp := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"hostname": types.StringType,
			"paths":    types.ListType{ElemType: types.StringType},
		},
	}, hostnamePaths)
	diags.Append(diagsTemp...)
	model.HostnamePaths = hostnamePathsList
	return diags
}

func populateCategoriesAndCustomCriteria(ctx context.Context, ilsModel *intelligentLoadSheddingModel, response *appsec.GetURLProtectionPolicyResponse, diags diag.Diagnostics) diag.Diagnostics {
	var categories []string
	var customCriteria []customCriteriaModel

	if len(response.Categories) > 0 {
		for _, category := range response.Categories {
			if category.Type == "CLIENT_LIST" {
				// Convert ListIDs to types.List
				var listIDsList types.List
				if len(category.ListIDs) > 0 {
					listIDsListTemp, diagsTemp := types.ListValueFrom(ctx, types.StringType, category.ListIDs)
					diags.Append(diagsTemp...)
					listIDsList = listIDsListTemp
				} else {
					var diagsTemp diag.Diagnostics
					listIDsList, diagsTemp = types.ListValueFrom(ctx, types.Int64Type, []int64{})
					if diagsTemp.HasError() {
						diags.Append(diagsTemp...)
						return diags
					}
				}

				customCriteria = append(customCriteria, customCriteriaModel{
					Type:          types.StringValue(category.Type),
					PositiveMatch: types.BoolValue(category.PositiveMatch != nil && *category.PositiveMatch),
					ListIDs:       listIDsList,
				})
			} else {
				categories = append(categories, category.Type)
			}
		}
	}

	// Convert categories to types.List
	if len(categories) > 0 {
		categoriesList, diagsTemp := types.ListValueFrom(ctx, types.StringType, categories)
		diags.Append(diagsTemp...)
		ilsModel.Categories = categoriesList
	} else {
		var diagsTemp diag.Diagnostics
		ilsModel.Categories, diagsTemp = types.ListValueFrom(ctx, types.StringType, []string{})
		if diagsTemp.HasError() {
			diags.Append(diagsTemp...)
			return diags
		}
	}

	// Convert customCriteria to types.List
	if len(customCriteria) > 0 {
		customCriteriaList, diagsTemp := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":           types.StringType,
				"list_ids":       types.ListType{ElemType: types.StringType},
				"positive_match": types.BoolType,
			},
		}, customCriteria)
		diags.Append(diagsTemp...)
		ilsModel.CustomCriteria = customCriteriaList
	} else {
		var diagsTemp diag.Diagnostics
		ilsModel.CustomCriteria, diagsTemp = types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":           types.StringType,
				"list_ids":       types.ListType{ElemType: types.StringType},
				"positive_match": types.BoolType,
			},
		}, []customCriteriaModel{})
		if diagsTemp.HasError() {
			diags.Append(diagsTemp...)
			return diags
		}

	}
	return diags
}
