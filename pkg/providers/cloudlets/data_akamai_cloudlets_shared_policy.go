package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type sharedPolicyModel struct {
	ID                 types.String     `tfsdk:"id"`
	PolicyID           types.Int64      `tfsdk:"policy_id"`
	Version            types.Int64      `tfsdk:"version"`
	VersionDescription types.String     `tfsdk:"version_description"`
	GroupID            types.Int64      `tfsdk:"group_id"`
	Name               types.String     `tfsdk:"name"`
	CloudletType       types.String     `tfsdk:"cloudlet_type"`
	Description        types.String     `tfsdk:"description"`
	MatchRules         types.String     `tfsdk:"match_rules"`
	Warnings           types.String     `tfsdk:"warnings"`
	Activations        *activationModel `tfsdk:"activations"`
}

type activationModel struct {
	Production activationInfoModel `tfsdk:"production"`
	Staging    activationInfoModel `tfsdk:"staging"`
}

type activationInfoModel struct {
	Effective *policyActivationModel `tfsdk:"effective"`
	Latest    *policyActivationModel `tfsdk:"latest"`
}

type policyActivationModel struct {
	ActivationID         types.Int64  `tfsdk:"activation_id"`
	CreatedBy            types.String `tfsdk:"created_by"`
	CreatedDate          types.String `tfsdk:"created_date"`
	FinishDate           types.String `tfsdk:"finish_date"`
	Network              types.String `tfsdk:"network"`
	Operation            types.String `tfsdk:"operation"`
	PolicyID             types.Int64  `tfsdk:"policy_id"`
	PolicyVersion        types.Int64  `tfsdk:"policy_version"`
	Status               types.String `tfsdk:"status"`
	PolicyVersionDeleted types.Bool   `tfsdk:"policy_version_deleted"`
}

var (
	_ datasource.DataSource              = &sharedPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &sharedPolicyDataSource{}
)

type sharedPolicyDataSource struct {
	meta meta.Meta
}

// NewSharedPolicyDataSource returns a new cloudlets shared policy data source
func NewSharedPolicyDataSource() datasource.DataSource {
	return &sharedPolicyDataSource{}
}

func (d *sharedPolicyDataSource) name() string {
	return "akamai_cloudlets_shared_policy"
}

// Metadata configures data source's meta information
func (d *sharedPolicyDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.name()
}

// Configure configures data source at the beginning of the lifecycle
func (d *sharedPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	d.meta = meta.Must(req.ProviderData)
}

// Schema is used to define data source's terraform schema
func (d *sharedPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	policyActivationAttributes := map[string]schema.Attribute{
		"activation_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Identifies the activation.",
		},
		"created_by": schema.StringAttribute{
			Computed:    true,
			Description: "The username who created the activation.",
		},
		"created_date": schema.StringAttribute{
			Computed:    true,
			Description: "ISO 8601 timestamp indicating when the activation was created.",
		},
		"finish_date": schema.StringAttribute{
			Computed:    true,
			Description: "ISO 8601 timestamp indicating when the activation ended, either successfully or unsuccessfully. You can check details of unsuccessful attempts in 'failureDetails'.",
		},
		"network": schema.StringAttribute{
			Computed:    true,
			Description: "The networks where you can activate or deactivate the policy version, either 'PRODUCTION' or 'STAGING'.",
		},
		"operation": schema.StringAttribute{
			Computed:    true,
			Description: "The operations that you can perform on a policy version, either 'ACTIVATION' or 'DEACTIVATION'.",
		},
		"policy_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Identifies the shared policy.",
		},
		"policy_version": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of the policy version.",
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "The status of the operation, either 'IN_PROGRESS', 'SUCCESS', or 'FAILED'.",
		},
		"policy_version_deleted": schema.BoolAttribute{
			Computed:    true,
			Description: "Indicates if the policy version is deleted.",
		},
	}
	activationBlock := map[string]schema.Block{
		"effective": schema.SingleNestedBlock{
			Description: "The status of the activation that's currently in use on this network, or null if the policy has no activations.",
			Attributes:  policyActivationAttributes,
		},
		"latest": schema.SingleNestedBlock{
			Description: "The status of the latest activation or null if the policy has no activations.",
			Attributes:  policyActivationAttributes,
		},
	}

	resp.Schema = schema.Schema{
		Description: "Cloudlets Shared Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "Required by the terraform plugin testing framework, always set to `akamai_cloudlets_shared_policy`.",
				Description:        "ID of the data source.",
			},
			"policy_id": schema.Int64Attribute{
				Required:    true,
				Description: "An integer ID that is associated with a policy.",
			},
			"version": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of the policy version.",
			},
			"version_description": schema.StringAttribute{
				Computed:    true,
				Description: "A human-readable label for the policy version.",
			},
			"group_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Identifies the group where to which policy is assigned.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the policy.",
			},
			"cloudlet_type": schema.StringAttribute{
				Computed:    true,
				Description: "The two- or three- letter code of the Cloudlet that the shared policy is for (AP, CD, ER, FR, IG, AS, VWR).",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A human-readable label for the policy.",
			},
			"match_rules": schema.StringAttribute{
				Computed:    true,
				Description: "A list of Cloudlet-specific match rules for this shared policy as a JSON",
			},
			"warnings": schema.StringAttribute{
				Computed:    true,
				Description: "A JSON encoded list of warnings.",
			},
		},
		Blocks: map[string]schema.Block{
			"activations": schema.SingleNestedBlock{
				Description: "Information about the active policy version that's currently in use and the status of the most recent activation or deactivation operation on the policy's versions for the production and staging networks.",
				Blocks: map[string]schema.Block{
					"production": schema.SingleNestedBlock{
						Description: "The policy version number that's currently in use on this network and the status of the most recent activation or deactivation operation for this policy's versions.",
						Blocks:      activationBlock,
					},
					"staging": schema.SingleNestedBlock{
						Description: "The policy version number that's currently in use on this network and the status of the most recent activation or deactivation operation for this policy's versions.",
						Blocks:      activationBlock,
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *sharedPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloudlets Shared Policy DataSource Read")

	var data sharedPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := ClientV3(d.meta)
	policy, err := client.GetPolicy(ctx, v3.GetPolicyRequest{
		PolicyID: data.PolicyID.ValueInt64(),
	})
	if err != nil {
		if errors.Is(err, v3.ErrPolicyNotFound) {
			resp.Diagnostics.AddError("Policy does not exist or is not of 'SHARED' type", err.Error())
			return
		}
		resp.Diagnostics.AddError("Reading Cloudlets Shared Policy Failed", err.Error())
		return
	}

	version := data.Version.ValueInt64()
	if version == 0 {
		policyVersions, err := client.ListPolicyVersions(ctx, v3.ListPolicyVersionsRequest{
			PolicyID: data.PolicyID.ValueInt64(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Reading Cloudlets Shared Policy Failed", err.Error())
			return
		}

		if len(policyVersions.PolicyVersions) != 0 {
			version = policyVersions.PolicyVersions[0].PolicyVersion
		}
	}

	var policyVersion *v3.PolicyVersion
	if version != 0 {
		policyVersion, err = client.GetPolicyVersion(ctx, v3.GetPolicyVersionRequest{
			PolicyID:      data.PolicyID.ValueInt64(),
			PolicyVersion: version,
		})
		if err != nil {
			resp.Diagnostics.AddError("Reading Cloudlets Shared Policy Failed", err.Error())
			return
		}

		matchRulesWarnings, err := json.Marshal(policyVersion.MatchRulesWarnings)
		if err != nil {
			resp.Diagnostics.AddError("Reading Cloudlets Shared Policy Failed", err.Error())
			return
		}
		data.Warnings = types.StringValue(string(matchRulesWarnings))
		matchRules, err := json.Marshal(policyVersion.MatchRules)
		if err != nil {
			resp.Diagnostics.AddError("Reading Cloudlets Shared Policy Failed", err.Error())
			return
		}
		data.MatchRules = types.StringValue(string(matchRules))
		data.Version = types.Int64Value(version)
		if policyVersion.Description != nil {
			data.VersionDescription = types.StringValue(*policyVersion.Description)
		}
	}

	data.setPolicyData(policy)
	data.setActivations(policy.CurrentActivations)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *sharedPolicyModel) setPolicyData(policy *v3.Policy) {
	m.CloudletType = types.StringValue(string(policy.CloudletType))
	m.Name = types.StringValue(policy.Name)
	m.ID = types.StringValue("akamai_cloudlets_shared_policy")
	m.GroupID = types.Int64Value(policy.GroupID)
	if policy.Description != nil {
		m.Description = types.StringValue(*policy.Description)
	}
}

func (m *sharedPolicyModel) setActivations(activations v3.CurrentActivations) {
	actMod := &activationModel{
		Production: activationInfoModel{},
		Staging:    activationInfoModel{},
	}
	if activations.Production.Effective != nil {
		actMod.Production.Effective = getActivationModel(activations.Production.Effective)
	}
	if activations.Production.Latest != nil {
		actMod.Production.Latest = getActivationModel(activations.Production.Latest)
	}
	if activations.Staging.Effective != nil {
		actMod.Staging.Effective = getActivationModel(activations.Staging.Effective)
	}
	if activations.Staging.Latest != nil {
		actMod.Staging.Latest = getActivationModel(activations.Staging.Latest)
	}
	m.Activations = actMod
}

func getActivationModel(activation *v3.PolicyActivation) *policyActivationModel {
	return &policyActivationModel{
		ActivationID:         types.Int64Value(activation.ID),
		CreatedBy:            types.StringValue(activation.CreatedBy),
		CreatedDate:          types.StringValue(activation.CreatedDate.String()),
		FinishDate:           types.StringValue(activation.FinishDate.String()),
		Network:              types.StringValue(string(activation.Network)),
		Operation:            types.StringValue(string(activation.Operation)),
		PolicyID:             types.Int64Value(activation.PolicyID),
		PolicyVersion:        types.Int64Value(activation.PolicyVersion),
		Status:               types.StringValue(string(activation.Status)),
		PolicyVersionDeleted: types.BoolValue(activation.PolicyVersionDeleted),
	}
}
