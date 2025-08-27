package cloudaccess

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &keyDataSource{}
	_ datasource.DataSourceWithConfigure = &keyDataSource{}
)

type (
	keyDataSource struct {
		meta meta.Meta
	}

	keyDataSourceModel struct {
		AccessKeyUID         types.Int64                `tfsdk:"access_key_uid"`
		AccessKeyName        types.String               `tfsdk:"access_key_name"`
		Groups               []groupModel               `tfsdk:"groups"`
		AuthenticationMethod types.String               `tfsdk:"authentication_method"`
		CreatedTime          types.String               `tfsdk:"created_time"`
		CreatedBy            types.String               `tfsdk:"created_by"`
		LatestVersion        types.Int64                `tfsdk:"latest_version"`
		NetworkConfiguration *networkConfigurationModel `tfsdk:"network_configuration"`
	}

	groupModel struct {
		GroupID      types.Int64    `tfsdk:"group_id"`
		GroupName    types.String   `tfsdk:"group_name"`
		ContractsIDs []types.String `tfsdk:"contracts_ids"`
	}

	networkConfigurationModel struct {
		AdditionalCDN   types.String `tfsdk:"additional_cdn"`
		SecurityNetwork types.String `tfsdk:"security_network"`
	}
)

// NewKeyDataSource returns a new location's data source
func NewKeyDataSource() datasource.DataSource {
	return &keyDataSource{}
}

// Metadata configures data source's meta information
func (d *keyDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudaccess_key"
}

// Configure configures data source at the beginning of the lifecycle
func (d *keyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema is used to define data source's terraform schema
func (d *keyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cloud Access key",
		Attributes: map[string]schema.Attribute{
			"access_key_uid": schema.Int64Attribute{
				Computed:    true,
				Description: "Identifier of the retrieved access key.",
			},
			"access_key_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the access key.",
			},
			"groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of groups to which the access key is assigned.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique identifier of Akamai group that's associated with the access key.",
						},
						"group_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of Akamai group that's associated with the access key.",
						},
						"contracts_ids": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The Akamai contracts that are associated with this access key for the group_id.",
						},
					},
				},
			},
			"authentication_method": schema.StringAttribute{
				Computed:    true,
				Description: "The type of signing process used to authenticate API requests: AWS4_HMAC_SHA256 for Amazon Web Services or GOOG4_HMAC_SHA256 for Google Cloud Services in interoperability mode.",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time the access key was created, in ISO 8601 format.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The username of the person who created the access key.",
			},
			"latest_version": schema.Int64Attribute{
				Computed:    true,
				Description: "The most recent version of the access key.",
			},
			"network_configuration": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The API deploys the access key to this secure network.",
				Attributes: map[string]schema.Attribute{
					"additional_cdn": schema.StringAttribute{
						Computed:    true,
						Description: "The access key can be deployed to the Akamai’s  additional networks. Available options are RUSSIA_CDN and CHINA_CDN.",
					},
					"security_network": schema.StringAttribute{
						Computed:    true,
						Description: "Attribute defines the type of secure network to which access key is deployed. Two options are available: STANDARD_TLS and ENHANCED_TLS.",
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *keyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudAccess Key DataSource Read")

	var data keyDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client = Client(d.meta)
	keys, err := client.ListAccessKeys(ctx, cloudaccess.ListAccessKeysRequest{})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("reading %s", ErrCloudAccessKey), err.Error())
		return
	}

	for _, key := range keys.AccessKeys {
		if key.AccessKeyName == data.AccessKeyName.ValueString() {
			data.AccessKeyName = types.StringValue(key.AccessKeyName)
			data.AccessKeyUID = types.Int64Value(key.AccessKeyUID)
			data.AuthenticationMethod = types.StringValue(key.AuthenticationMethod)
			data.CreatedBy = types.StringValue(key.CreatedBy)
			dateString, err := date.ToString(key.CreatedTime)
			if err != nil {
				resp.Diagnostics.AddError("error parsing date:", err.Error())
				return
			}
			data.CreatedTime = types.StringValue(dateString)
			for _, group := range key.Groups {
				contractIDs := make([]types.String, 0)
				for _, id := range group.ContractIDs {
					contractIDs = append(contractIDs, types.StringValue(id))
				}
				g := groupModel{
					GroupID:      types.Int64Value(group.GroupID),
					GroupName:    types.StringPointerValue(group.GroupName),
					ContractsIDs: contractIDs,
				}
				data.Groups = append(data.Groups, g)
			}
			data.LatestVersion = types.Int64Value(key.LatestVersion)
			if key.NetworkConfiguration.AdditionalCDN != nil {
				data.NetworkConfiguration = &networkConfigurationModel{
					AdditionalCDN:   types.StringValue(string(*key.NetworkConfiguration.AdditionalCDN)),
					SecurityNetwork: types.StringPointerValue(ptr.To(string(key.NetworkConfiguration.SecurityNetwork))),
				}
			} else {
				data.NetworkConfiguration = &networkConfigurationModel{
					SecurityNetwork: types.StringPointerValue(ptr.To(string(key.NetworkConfiguration.SecurityNetwork))),
				}
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return

		}
	}
	resp.Diagnostics.AddError("No matching key", "no key with given name")
}
