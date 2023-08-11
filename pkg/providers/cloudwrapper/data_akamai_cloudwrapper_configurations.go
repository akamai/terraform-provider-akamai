package cloudwrapper

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &configurationsDataSource{}
	_ datasource.DataSourceWithConfigure = &configurationsDataSource{}
)

type (
	configurationsDataSource struct {
		client cloudwrapper.CloudWrapper
	}

	configurationsDataSourceModel struct {
		ID             types.String                   `tfsdk:"id"`
		Configurations []configurationDataSourceModel `tfsdk:"configurations"`
	}
)

// NewConfigurationsDataSource returns configurations data source
func NewConfigurationsDataSource() datasource.DataSource {
	return &configurationsDataSource{}
}

// setClient assigns given client to properties data source
func (d *configurationsDataSource) setClient(client cloudwrapper.CloudWrapper) {
	d.client = client
}

// Metadata configures data source meta information
func (d *configurationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_configurations"
}

// Configure configures data source at the beginning of the lifecycle
func (d *configurationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	if d.client != nil {
		return
	}

	m, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	d.client = cloudwrapper.Client(m.Session())
}

// Schema is used to define data source terraform schema
func (d *configurationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudWrapper configurations",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "It will be removed after migration to the new testing framework",
				Description:        "ID of the data source.",
			},
		},
		Blocks: map[string]schema.Block{
			"configurations": schema.ListNestedBlock{
				Description: "List of the configurations on the contract.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of a Cloud Wrapper configuration.",
						},
						"capacity_alerts_threshold": schema.Int64Attribute{
							Computed:    true,
							Description: "Represents the threshold for sending alerts.",
						},
						"comments": schema.StringAttribute{
							Computed:    true,
							Description: "Additional information provided by user which can help to differentiate or track changes of the configuration.",
						},
						"config_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the configuration.",
						},
						"contract_id": schema.StringAttribute{
							Computed:    true,
							Description: "Contract ID with Cloud Wrapper entitlement.",
						},
						"last_activated_by": schema.StringAttribute{
							Computed:    true,
							Description: "User to last activate the configuration.",
						},
						"last_activated_date": schema.StringAttribute{
							Computed:    true,
							Description: "ISO format date that represents when the configuration was last activated successfully.",
						},
						"last_updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "User to last modify the configuration.",
						},
						"last_updated_date": schema.StringAttribute{
							Computed:    true,
							Description: "ISO format date that represents when the configuration was last edited.",
						},
						"notification_emails": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Email addresses to receive notifications.",
						},
						"property_ids": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "List of properties belonging to media delivery products. Properties need to be unique across configurations.",
						},
						"retain_idle_objects": schema.BoolAttribute{
							Computed:    true,
							Description: "Retain idle objects beyond their max idle lifetime.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current state of the provisioning of the configuration, either SAVED, IN_PROGRESS, ACTIVE, DELETE_IN_PROGRESS, or FAILED.",
						},
					},
					Blocks: configBlock,
				},
			},
		},
	}
}

func (d *configurationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudWrapper Configurations DataSource Read")

	var data configurationsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	configs, err := d.client.ListConfigurations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Reading CloudWrapper Configurations", err.Error())
		return
	}

	for _, cfg := range configs.Configurations {
		var cfgModel configurationDataSourceModel
		if resp.Diagnostics.Append(cfgModel.populate(ctx, &cfg)...); resp.Diagnostics.HasError() {
			return
		}
		cfgModel.ID = types.Int64Value(cfg.ConfigID)
		data.Configurations = append(data.Configurations, cfgModel)
	}

	data.ID = types.StringValue("akamai_cloudwrapper_configurations")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
