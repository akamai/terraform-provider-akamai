package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &blockedPropertiesDataSource{}
var _ datasource.DataSourceWithConfigure = &blockedPropertiesDataSource{}

// NewBlockedPropertiesDataSource returns all the properties that are blocked for a certain user in a group
func NewBlockedPropertiesDataSource() datasource.DataSource {
	return &blockedPropertiesDataSource{}
}

// blockedPropertiesDataSource defines the data source implementation for fetching Blocked Properties information
type blockedPropertiesDataSource struct {
	meta meta.Meta
}

// blockedPropertiesDataSource describes the data source data model for BlockedPropertiesDataSource
type blockedPropertiesDataSourceModel struct {
	GroupID           types.Int64  `tfsdk:"group_id"`
	ContractID        types.String `tfsdk:"contract_id"`
	UIIdentityID      types.String `tfsdk:"ui_identity_id"`
	BlockedProperties types.List   `tfsdk:"blocked_properties"`
}

// Metadata configures data source's meta information
func (d *blockedPropertiesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_blocked_properties"
}

// Schema is used to define data source's terraform schema
func (d *blockedPropertiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Blocked Properties data source",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier for each group.",
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Contract ID for which block properties are retrieved.",
			},
			"ui_identity_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for each user.",
			},
			"blocked_properties": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The list of blocked properties.",
			},
		},
	}
}

// Configure  configures data source at the beginning of the lifecycle
func (d *blockedPropertiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read is called when the provider must read data source values in order to update state
func (d *blockedPropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "BlockedPropertiesDataSource Read")

	var data blockedPropertiesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	papiClient := inst.PapiClient(d.meta)

	listBlockedPropertiesResp, err := client.ListBlockedProperties(ctx, iam.ListBlockedPropertiesRequest{
		IdentityID: data.UIIdentityID.ValueString(),
		GroupID:    data.GroupID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching iam blocked properties failed", err.Error())
		return
	}

	var blockedPropertiesPAPI []string
	groupID := str.AddPrefix(data.GroupID.String(), "grp_")

	for _, prop := range listBlockedPropertiesResp {
		req := iam.MapPropertyIDToNameRequest{
			GroupID:    data.GroupID.ValueInt64(),
			PropertyID: prop,
		}
		papiPropertyName, err := client.MapPropertyIDToName(ctx, req)
		if err != nil {
			resp.Diagnostics.AddError("fetching PAPI propertyName failed", err.Error())
			return
		}
		papiPropertyID, err := papiClient.MapPropertyNameToID(ctx, papi.MapPropertyNameToIDRequest{
			GroupID:    groupID,
			ContractID: data.ContractID.ValueString(),
			Name:       *papiPropertyName,
		})
		if err != nil {
			resp.Diagnostics.AddError("fetching PAPI propertyID failed", err.Error())
			return
		}

		blockedPropertiesPAPI = append(blockedPropertiesPAPI, *papiPropertyID)
	}

	blockedPropertyIDsPAPI, diags := types.ListValueFrom(ctx, types.StringType, blockedPropertiesPAPI)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.BlockedProperties = blockedPropertyIDsPAPI

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
