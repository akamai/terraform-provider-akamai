package mtlstruststore

import (
	"context"
	"fmt"
	"regexp"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &caSetActivationDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetActivationDataSource{}
)

type (
	caSetActivationDataSource struct {
		meta meta.Meta
	}

	caSetActivationDataSourceModel struct {
		ID           types.Int64  `tfsdk:"id"`
		CASetID      types.String `tfsdk:"ca_set_id"`
		CASetName    types.String `tfsdk:"ca_set_name"`
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

// NewCASetActivationDataSource returns a new mtls truststore ca set activation data source.
func NewCASetActivationDataSource() datasource.DataSource {
	return &caSetActivationDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetActivationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_activation"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetActivationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema is used to define data source's terraform schema.
func (d *caSetActivationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a specific MTLS Truststore CA Set Activation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Uniquely identifies the activation.",
				Required:    true,
			},
			"ca_set_id": schema.StringAttribute{
				Description: "CA set identifier that filters out the activations or deactivations. Either 'ca_set_id' or 'ca_set_name' must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("ca_set_id"), path.MatchRoot("ca_set_name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ca_set_name": schema.StringAttribute{
				Description: "The name of the CA set. Either 'ca_set_id' or 'ca_set_name' must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("ca_set_id"), path.MatchRoot("ca_set_name")),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S{3,}`), "must not be empty or only whitespace"),
				},
			},
			"version": schema.Int64Attribute{
				Description: "CA set version identifier.",
				Computed:    true,
			},
			"network": schema.StringAttribute{
				Description: "Indicates the network for any activation-related activities, either 'STAGING' or 'PRODUCTION'.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user who requested the activity.",
				Computed:    true,
			},
			"created_date": schema.StringAttribute{
				Description: "When the activity was requested.",
				Computed:    true,
			},
			"modified_by": schema.StringAttribute{
				Description: "The user who completed the activity.",
				Computed:    true,
			},
			"modified_date": schema.StringAttribute{
				Description: "When the request was last modified, or null if not yet modified.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the current activity, either 'IN_PROGRESS', 'COMPLETE', or 'FAILED'.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of requested activity, either 'ACTIVATE', 'DEACTIVATE', or 'DELETE'.",
				Computed:    true,
			},
		},
	}
}

func (d *caSetActivationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set Activation DataSource Read")

	var data caSetActivationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client = Client(d.meta)
	if !data.CASetName.IsNull() {
		tflog.Debug(ctx, "'ca_set_name' provided, attempting to find CA set ID")
		caSetID, err := findCASetID(ctx, client, data.CASetName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read CA Set Activation failed", err.Error())
			return
		}
		data.CASetID = types.StringValue(caSetID)
	}

	activations, err := client.ListCASetActivations(ctx, mtlstruststore.ListCASetActivationsRequest{
		CASetID: data.CASetID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read CA Set Activation failed", err.Error())
		return
	}
	act, err := findActivationByID(activations.Activations, data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Read CA Set Activation failed", err.Error())
		return
	}

	data.setFromActivationResponse(act)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func findActivationByID(activations []mtlstruststore.ActivateCASetVersionResponse, id int64) (mtlstruststore.ActivateCASetVersionResponse, error) {
	for _, a := range activations {
		if a.ActivationID == id {
			return a, nil
		}
	}
	return mtlstruststore.ActivateCASetVersionResponse{}, fmt.Errorf("activation with ID %d not found", id)
}

func (m *caSetActivationDataSourceModel) setFromActivationResponse(act mtlstruststore.ActivateCASetVersionResponse) {
	m.ID = types.Int64Value(act.ActivationID)
	m.CASetID = types.StringValue(act.CASetID)
	m.CASetName = types.StringValue(act.CASetName)
	m.Version = types.Int64Value(act.Version)
	m.Network = types.StringValue(act.Network)
	m.CreatedBy = types.StringValue(act.CreatedBy)
	m.CreatedDate = date.TimeRFC3339NanoValue(act.CreatedDate)
	m.ModifiedBy = types.StringPointerValue(act.ModifiedBy)
	m.ModifiedDate = date.TimeRFC3339NanoPointerValue(act.ModifiedDate)
	m.Status = types.StringValue(act.ActivationStatus)
	m.Type = types.StringValue(act.ActivationType)
}
