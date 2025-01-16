package cloudlets

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type policyActivationDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	PolicyID             types.Int64  `tfsdk:"policy_id"`
	Network              types.String `tfsdk:"network"`
	Version              types.Int64  `tfsdk:"version"`
	Status               types.String `tfsdk:"status"`
	AssociatedProperties types.Set    `tfsdk:"associated_properties"`
}

var (
	_ datasource.DataSource              = &policyActivationDataSource{}
	_ datasource.DataSourceWithConfigure = &policyActivationDataSource{}
)

type policyActivationDataSource struct {
	meta meta.Meta
}

// NewPolicyActivationDataSource returns a new capacity data source
func NewPolicyActivationDataSource() datasource.DataSource {
	return &policyActivationDataSource{}
}

// Metadata configures data source's meta information
func (d *policyActivationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.name()
}

func (d *policyActivationDataSource) name() string {
	return "akamai_cloudlets_policy_activation"
}

// Configure configures data source at the beginning of the lifecycle
func (d *policyActivationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	m, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	d.meta = m
}

// Schema is used to define data source's terraform schema
func (d *policyActivationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cloudlets Policy Activation",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "Required by the terraform plugin testing framework.",
				Description:        "ID of the data source.",
			},
			"policy_id": schema.Int64Attribute{
				Required:    true,
				Description: "Identifies the policy.",
			},
			"network": schema.StringAttribute{
				Required:    true,
				Description: "The networks where you can get activated policy version (options are Staging and Production).",
				Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive("production", "prod", "p", "staging", "stag", "s")},
			},
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "Policy version that is activated on provided network.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Activation status for this Cloudlets policy.",
			},
			"associated_properties": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of associated properties for non-shared cloudlets activation policy.",
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *policyActivationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	logger := d.meta.Log("Cloudlets", "Read")
	logger.Debug("Cloudlets Policy Activation DataSource Read")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))

	var data policyActivationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := data.PolicyID.ValueInt64()
	network := data.Network.ValueString()

	strategy, _, err := discoverActivationStrategy(ctx, policyID, d.meta, logger)
	if err != nil {
		resp.Diagnostics.AddError("Reading Policy Failed", err.Error())
		return
	}

	dat, err := strategy.getPolicyActivation(ctx, policyID, network)
	if err != nil {
		resp.Diagnostics.AddError("Reading Policy Failed", err.Error())
		return
	}
	data = *dat
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", data.PolicyID, data.Network))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (strategy *v2ActivationStrategy) getPolicyActivation(ctx context.Context, policyID int64, network string) (*policyActivationDataSourceModel, error) {
	listPolicyActivationsRequest := cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  cloudlets.PolicyActivationNetwork(network),
	}

	activations, err := strategy.client.ListPolicyActivations(ctx, listPolicyActivationsRequest)
	if err != nil {
		return nil, fmt.Errorf("%v read: reading list of policy activations failed. %s", ErrPolicyActivation, err.Error())
	}

	if len(activations) == 0 {
		return nil, fmt.Errorf("%v read: cannot find any activation for the given policy '%d' and network '%s'", ErrPolicyActivation, policyID, network)
	}

	activations = sortPolicyActivationsByDate(activations)
	associatedProperties := getActiveProperties(activations)

	ap, d := types.SetValueFrom(ctx, types.StringType, associatedProperties)
	if d.HasError() {
		return nil, errors.New(d.Errors()[0].Summary())
	}
	data := policyActivationDataSourceModel{
		PolicyID:             types.Int64Value(policyID),
		Network:              types.StringValue(network),
		Version:              types.Int64Value(activations[0].PolicyInfo.Version),
		Status:               types.StringValue(string(activations[0].PolicyInfo.Status)),
		AssociatedProperties: ap,
	}
	return &data, nil
}

func (strategy *v3ActivationStrategy) getPolicyActivation(ctx context.Context, policyID int64, network string) (*policyActivationDataSourceModel, error) {
	getPolicyRequest := v3.GetPolicyRequest{
		PolicyID: policyID,
	}

	policy, err := strategy.client.GetPolicy(ctx, getPolicyRequest)
	if err != nil {
		return nil, fmt.Errorf("%v read: reading policy failed. %s", ErrPolicyActivation, err.Error())
	}

	var effective *v3.PolicyActivation
	switch tf.StateNetwork(network) {
	case "staging":
		effective = policy.CurrentActivations.Staging.Effective
	case "production":
		effective = policy.CurrentActivations.Production.Effective
	}
	if effective == nil || effective.Operation != v3.OperationActivation {
		return nil, fmt.Errorf("%v read: cannot find any activation for the given policy '%d' and network '%s'", ErrPolicyActivation, policyID, network)
	}
	data := policyActivationDataSourceModel{
		PolicyID:             types.Int64Value(policyID),
		Network:              types.StringValue(network),
		Version:              types.Int64Value(effective.PolicyVersion),
		Status:               types.StringValue(string(effective.Status)),
		AssociatedProperties: types.SetNull(types.StringType),
	}
	return &data, nil
}
