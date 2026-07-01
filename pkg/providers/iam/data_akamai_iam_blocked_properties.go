package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &blockedPropertiesDataSource{}
	_ datasource.DataSourceWithConfigure = &blockedPropertiesDataSource{}
)

// NewBlockedPropertiesDataSource returns all the properties that are blocked for a certain user in a group.
func NewBlockedPropertiesDataSource() datasource.DataSource {
	return &blockedPropertiesDataSource{}
}

type (
	blockedPropertiesDataSource struct {
		meta.DataSource
	}

	blockedPropertiesModel struct {
		GroupID           types.Int64               `tfsdk:"group_id"`
		ContractID        types.String              `tfsdk:"contract_id"`
		UIIdentityID      types.String              `tfsdk:"ui_identity_id"`
		BlockedProperties []blockedPropertyIDsModel `tfsdk:"blocked_properties"`
	}

	blockedPropertyIDsModel struct {
		PropertyID types.String `tfsdk:"property_id"`
		AssetID    types.Int64  `tfsdk:"asset_id"`
	}
)

func (d *blockedPropertiesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_blocked_properties"
}

func (d *blockedPropertiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Blocked Properties data source.",
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
			"blocked_properties": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of blocked properties.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property_id": schema.StringAttribute{
							Computed:    true,
							Description: "PAPI's blocked property ID.",
						},
						"asset_id": schema.Int64Attribute{
							Computed:    true,
							Description: "IAM's blocked property ID.",
						},
					},
				},
			},
		},
	}
}

func (d *blockedPropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Blocked Properties DataSource Read")

	var data blockedPropertiesModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetIAM()
	papiClient := d.Client.GetPAPI()

	listBlockedPropertiesResp, err := client.ListBlockedProperties(ctx, iam.ListBlockedPropertiesRequest{
		IdentityID: data.UIIdentityID.ValueString(),
		GroupID:    data.GroupID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching iam blocked properties failed", err.Error())
		return
	}

	var blockedProperties []blockedPropertyIDsModel
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

		blockedProperties = append(blockedProperties, blockedPropertyIDsModel{PropertyID: types.StringValue(*papiPropertyID), AssetID: types.Int64Value(prop)})
	}

	data.BlockedProperties = blockedProperties

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
