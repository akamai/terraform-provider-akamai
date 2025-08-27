package cloudaccess

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &keyPropertiesDataSource{}
	_ datasource.DataSourceWithConfigure = &keyPropertiesDataSource{}
)

type (
	keyPropertiesDataSource struct {
		meta meta.Meta
	}

	keyPropertiesDataSourceModel struct {
		AccessKeyName types.String    `tfsdk:"access_key_name"`
		Properties    []propertyModel `tfsdk:"properties"`
		AccessKeyUID  types.Int64     `tfsdk:"access_key_uid"`
	}

	propertyModel struct {
		AccessKeyVersion  types.Int64  `tfsdk:"access_key_version"`
		PropertyID        types.String `tfsdk:"property_id"`
		PropertyName      types.String `tfsdk:"property_name"`
		StagingVersion    types.Int64  `tfsdk:"staging_version"`
		ProductionVersion types.Int64  `tfsdk:"production_version"`
	}
)

// NewKeyPropertiesDataSource returns a new cloud access key properties data source
func NewKeyPropertiesDataSource() datasource.DataSource {
	return &keyPropertiesDataSource{}
}

// Metadata configures data source's meta information
func (d *keyPropertiesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudaccess_key_properties"
}

// Configure configures data source at the beginning of the lifecycle
func (d *keyPropertiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *keyPropertiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cloud Access key properties",
		Attributes: map[string]schema.Attribute{
			"access_key_name": schema.StringAttribute{
				Description: "Name of the access key.",
				Required:    true,
			},
			"access_key_uid": schema.Int64Attribute{
				Description: "Uniquely identifies the access key",
				Computed:    true,
			},
			"properties": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Property lookup results, one per property.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_key_version": schema.Int64Attribute{
							Computed:    true,
							Description: "Version of the access key.",
						},
						"property_id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier Akamai assigned to the matching property.",
						},
						"property_name": schema.StringAttribute{
							Computed:    true,
							Description: "Identifies the specific property name whose Origin Characteristics behavior uses the access key version.",
						},
						"staging_version": schema.Int64Attribute{
							Computed:    true,
							Description: "Identifies the specific property version whose staging status is either active or activating.",
						},
						"production_version": schema.Int64Attribute{
							Computed:    true,
							Description: "Identifies the specific property version whose production status is either active or activating.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *keyPropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudAccess Key Properties DataSource Read")

	var data keyPropertiesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client = Client(d.meta)

	if resp.Diagnostics.Append(data.getAccessKey(ctx)...); resp.Diagnostics.HasError() {
		return
	}

	versions, err := client.ListAccessKeyVersions(ctx, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: data.AccessKeyUID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s: ListAccessKeyVersions failed:", ErrCloudAccessKeyProperties), err.Error())
		return
	}

	if resp.Diagnostics.Append(data.read(ctx, versions, client)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (data *keyPropertiesDataSourceModel) getAccessKey(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	accessKeys, err := client.ListAccessKeys(ctx, cloudaccess.ListAccessKeysRequest{})
	if err != nil {
		diags.AddError(fmt.Sprintf("%s: ListAccessKeys failed:", ErrCloudAccessKeyProperties), err.Error())
		return diags
	}

	for _, key := range accessKeys.AccessKeys {
		if key.AccessKeyName == data.AccessKeyName.ValueString() {
			data.AccessKeyUID = types.Int64Value(key.AccessKeyUID)
			return nil
		}
	}
	diags.AddError(fmt.Sprintf("%s:", ErrCloudAccessKeyProperties), fmt.Sprintf("access key with name: '%s' does not exist", data.AccessKeyName.ValueString()))
	return diags
}

func (data *keyPropertiesDataSourceModel) read(ctx context.Context, versions *cloudaccess.ListAccessKeyVersionsResponse, client cloudaccess.CloudAccess) diag.Diagnostics {
	var diags diag.Diagnostics
	var propertiesModel []propertyModel
	for _, ver := range versions.AccessKeyVersions {
		properties, err := client.LookupProperties(ctx, cloudaccess.LookupPropertiesRequest{
			AccessKeyUID: data.AccessKeyUID.ValueInt64(),
			Version:      ver.Version,
		})
		if err != nil {
			diags.AddError(fmt.Sprintf("%s: LookupProperties failed:", ErrCloudAccessKeyProperties), err.Error())
			return diags
		}

		for _, prp := range properties.Properties {
			if prp.ProductionVersion != nil || prp.StagingVersion != nil {
				model := propertyModel{
					AccessKeyVersion: types.Int64Value(ver.Version),
					PropertyID:       types.StringValue(prp.PropertyID),
					PropertyName:     types.StringValue(prp.PropertyName),
				}
				if prp.ProductionVersion != nil {
					model.ProductionVersion = types.Int64Value(*prp.ProductionVersion)
				}
				if prp.StagingVersion != nil {
					model.StagingVersion = types.Int64Value(*prp.StagingVersion)
				}
				propertiesModel = append(propertiesModel, model)
			}
		}
	}
	data.Properties = propertiesModel

	return nil
}
