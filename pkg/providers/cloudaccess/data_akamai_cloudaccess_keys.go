package cloudaccess

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &keysDataSource{}
	_ datasource.DataSourceWithConfigure = &keysDataSource{}
)

type (
	keysDataSource struct {
		meta meta.Meta
	}

	keysDataSourceModel struct {
		AccessKeys []keyModel `tfsdk:"access_keys"`
	}

	keyModel struct {
		AccessKeyUID         types.Int64                `tfsdk:"access_key_uid"`
		AccessKeyName        types.String               `tfsdk:"access_key_name"`
		Groups               []groupModel               `tfsdk:"groups"`
		AuthenticationMethod types.String               `tfsdk:"authentication_method"`
		CreatedTime          types.String               `tfsdk:"created_time"`
		CreatedBy            types.String               `tfsdk:"created_by"`
		LatestVersion        types.Int64                `tfsdk:"latest_version"`
		NetworkConfiguration *networkConfigurationModel `tfsdk:"network_configuration"`
	}
)

// NewKeysDataSource returns a new cloudaccess keys data source
func NewKeysDataSource() datasource.DataSource {
	return &keysDataSource{}
}

// Metadata configures data source's meta information
func (d *keysDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudaccess_keys"
}

// Configure configures data source at the beginning of the lifecycle
func (d *keysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *keysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cloud Access keys",
		Attributes: map[string]schema.Attribute{
			"access_keys": schema.ListNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
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
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *keysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudAccess Keys DataSource Read")

	var data keysDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client = Client(d.meta)
	keys, err := client.ListAccessKeys(ctx, cloudaccess.ListAccessKeysRequest{})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s: ListAccessKeys failed:", ErrCloudAccessKeys), err.Error())
		return
	}

	if resp.Diagnostics.Append(data.read(keys)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *keysDataSourceModel) read(keys *cloudaccess.ListAccessKeysResponse) diag.Diagnostics {
	var accessKeys []keyModel
	var diags diag.Diagnostics

	for _, key := range keys.AccessKeys {
		var groups []groupModel
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
			groups = append(groups, g)
		}

		var netConf *networkConfigurationModel
		if key.NetworkConfiguration != nil {
			netConf = &networkConfigurationModel{
				SecurityNetwork: types.StringValue(string(key.NetworkConfiguration.SecurityNetwork)),
			}
			if key.NetworkConfiguration.AdditionalCDN != nil {
				netConf.AdditionalCDN = types.StringValue(string(*key.NetworkConfiguration.AdditionalCDN))
			}
		}

		stringDate, err := date.ToString(key.CreatedTime)
		if err != nil {
			diags.AddError("error parsing date:", err.Error())
			return diags
		}
		accessKey := keyModel{
			AccessKeyUID:         types.Int64Value(key.AccessKeyUID),
			AccessKeyName:        types.StringValue(key.AccessKeyName),
			Groups:               groups,
			AuthenticationMethod: types.StringValue(key.AuthenticationMethod),
			CreatedTime:          types.StringValue(stringDate),
			CreatedBy:            types.StringValue(key.CreatedBy),
			LatestVersion:        types.Int64Value(key.LatestVersion),
			NetworkConfiguration: netConf,
		}
		accessKeys = append(accessKeys, accessKey)
	}

	d.AccessKeys = accessKeys
	return nil
}
