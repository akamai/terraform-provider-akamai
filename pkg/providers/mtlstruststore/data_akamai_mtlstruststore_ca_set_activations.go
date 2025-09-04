package mtlstruststore

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &caSetDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetDataSource{}
)

type (
	caSetActivationsDataSource struct {
		meta meta.Meta
	}

	caSetActivationsDataSourceModel struct {
		CASetID     types.String      `tfsdk:"ca_set_id"`
		CASetName   types.String      `tfsdk:"ca_set_name"`
		Network     types.String      `tfsdk:"network"`
		Version     types.Int64       `tfsdk:"version"`
		Status      types.String      `tfsdk:"status"`
		Type        types.String      `tfsdk:"type"`
		Activations []activationModel `tfsdk:"activations"`
	}

	activationModel struct {
		ID           types.Int64  `tfsdk:"id"`
		Version      types.Int64  `tfsdk:"version"`
		Network      types.String `tfsdk:"network"`
		CreatedBy    types.String `tfsdk:"created_by"`
		CreatedDate  types.String `tfsdk:"created_date"`
		ModifiedBy   types.String `tfsdk:"modified_by"`
		ModifiedDate types.String `tfsdk:"modified_date"`
		Status       types.String `tfsdk:"status"`
		Type         types.String `tfsdk:"type"`
	}
)

// NewCASetActivationsDataSource returns a new mtls truststore ca set activations data source.
func NewCASetActivationsDataSource() datasource.DataSource {
	return &caSetActivationsDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetActivationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_activations"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetActivationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	d.meta = meta.Must(req.ProviderData)
}

func (d *caSetActivationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve list of MTLS Truststore CA Set activations.",
		Attributes: map[string]schema.Attribute{
			"ca_set_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "CA set identifier that filters out the activations or deactivations. Either 'ca_set_id' or 'ca_set_name' must be provided.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("ca_set_name"), path.MatchRoot("ca_set_id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ca_set_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the CA set. Either 'ca_set_id' or 'ca_set_name' must be provided.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("ca_set_name"), path.MatchRoot("ca_set_id")),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S{3,}`), "must not be empty or only whitespace"),
				},
			},
			"network": schema.StringAttribute{
				Optional:    true,
				Description: "If provided, filters the results to return only activations or deactivations from the specified network, either 'STAGING' or 'PRODUCTION'.",
			},
			"version": schema.Int64Attribute{
				Optional:    true,
				Description: "If provided, filters the results to return only activities for the specified CA set version.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "If provided, filters the results to return only activities with the specified status, either 'IN_PROGRESS', 'COMPLETE', or 'FAILED'.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "If provided, filters the results to return only activities of the specified type, either 'ACTIVATE', 'DEACTIVATE', or 'DELETE'.",
			},
			"activations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of CA set activations or deactivations.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "Uniquely identifies the activation.",
						},
						"network": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates the network for any activation-related activities, either 'STAGING' or 'PRODUCTION'.",
						},
						"version": schema.Int64Attribute{
							Computed:    true,
							Description: "CA set version identifier.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the current activity, either 'IN_PROGRESS', 'COMPLETE', or 'FAILED'.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of requested activity, either 'ACTIVATE', 'DEACTIVATE', or 'DELETE'.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user who requested the activity.",
						},
						"created_date": schema.StringAttribute{
							Computed:    true,
							Description: "When the activity was requested.",
						},
						"modified_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user who completed the activity.",
						},
						"modified_date": schema.StringAttribute{
							Computed:    true,
							Description: "When the request was last modified, or null` if not yet modified.",
						},
					},
				},
			},
		},
	}
}

func (d *caSetActivationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set Activations DataSource Read")

	var data caSetActivationsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(d.meta)

	if !data.CASetName.IsNull() {
		tflog.Debug(ctx, "'name' provided, attempting to find CA set ID")
		setID, err := findCASetID(ctx, client, data.CASetName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read CA set activations failed", err.Error())
			return
		}
		data.CASetID = types.StringValue(setID)
	}

	response, err := client.ListCASetActivations(ctx, mtlstruststore.ListCASetActivationsRequest{
		CASetID: data.CASetID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read CA set activations failed", err.Error())
		return
	}

	var filteredActivations []activationModel
	actFilter := activationFilter{
		network:        data.Network.ValueString(),
		status:         data.Status.ValueString(),
		activationType: data.Type.ValueString(),
		version:        data.Version.ValueInt64(),
	}
	for _, act := range response.Activations {
		if actFilter.matches(act) {
			filteredActivations = append(filteredActivations, convertCASetActivationDataToModel(act))
		}
	}

	data.Activations = filteredActivations
	if data.CASetName.IsNull() && len(response.Activations) > 0 {
		data.CASetName = types.StringValue(response.Activations[0].CASetName)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func convertCASetActivationDataToModel(act mtlstruststore.ActivateCASetVersionResponse) activationModel {
	return activationModel{
		ID:           types.Int64Value(act.ActivationID),
		Version:      types.Int64Value(act.Version),
		Network:      types.StringValue(act.Network),
		CreatedBy:    types.StringValue(act.CreatedBy),
		CreatedDate:  date.TimeRFC3339NanoValue(act.CreatedDate),
		ModifiedBy:   types.StringPointerValue(act.ModifiedBy),
		ModifiedDate: date.TimeRFC3339NanoPointerValue(act.ModifiedDate),
		Status:       types.StringValue(act.ActivationStatus),
		Type:         types.StringValue(act.ActivationType),
	}
}

type activationFilter struct {
	network        string
	status         string
	activationType string
	version        int64
}

func (f activationFilter) matches(act mtlstruststore.ActivateCASetVersionResponse) bool {
	if f.network != "" && act.Network != strings.ToUpper(f.network) {
		return false
	}
	if f.version != 0 && act.Version != f.version {
		return false
	}
	if f.status != "" && act.ActivationStatus != strings.ToUpper(f.status) {
		return false
	}
	if f.activationType != "" && act.ActivationType != strings.ToUpper(f.activationType) {
		return false
	}
	return true
}
