package cloudaccess

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &keyVersionsDataSource{}
	_ datasource.DataSourceWithConfigure = &keyVersionsDataSource{}
)

type keyVersionsDataSource struct {
	meta meta.Meta
}
type keyVersionsDataSourceModel struct {
	AccessKeyName     types.String            `tfsdk:"access_key_name"`
	AccessKeyUID      types.Int64             `tfsdk:"access_key_uid"`
	AccessKeyVersions []accessKeyVersionModel `tfsdk:"access_key_versions"`
}

type accessKeyVersionModel struct {
	CloudAccessKeyID types.String `tfsdk:"cloud_access_key_id"`
	CreatedTime      types.String `tfsdk:"created_time"`
	CreatedBy        types.String `tfsdk:"created_by"`
	DeploymentStatus types.String `tfsdk:"deployment_status"`
	Version          types.Int64  `tfsdk:"version"`
	VersionGUID      types.String `tfsdk:"version_guid"`
}

// NewKeyVersionsDataSource returns a new key versions data source
func NewKeyVersionsDataSource() datasource.DataSource {
	return &keyVersionsDataSource{}
}

// Metadata configures data source's meta information
func (d *keyVersionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudaccess_key_versions"
}

// Configure configures data source at the beginning of the lifecycle
func (d *keyVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *keyVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Cloud Access key versions",
		Attributes: map[string]schema.Attribute{
			"access_key_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the access key.",
			},
			"access_key_uid": schema.Int64Attribute{
				Computed:    true,
				Description: "Identifier of the access key to retrieve.",
			},
			"access_key_versions": schema.ListNestedAttribute{
				Description: "List of access key versions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cloud_access_key_id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier assigned to the access key assigned from AWS or GCS.",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "The time the access key was created, in ISO 8601 format.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The username of the person who created the access key.",
						},
						"deployment_status": schema.StringAttribute{
							Computed:    true,
							Description: "Indicates if the version has been activated to the Akamai networks. Available statuses are: PENDING_DELETION, ACTIVE and PENDING_ACTIVATION.",
						},
						"version": schema.Int64Attribute{
							Computed:    true,
							Description: "Version of the access key.",
						},
						"version_guid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier assigned to an access key version.",
						},
					},
				},
			},
		},
	}

}

// Read is called when the provider must read data source values in order to update state
func (d *keyVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudAccess Key Versions DataSource Read")
	client = Client(d.meta)

	var data keyVersionsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	keys, err := client.ListAccessKeys(ctx, cloudaccess.ListAccessKeysRequest{})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("reading %s", ErrCloudAccessKey), err.Error())
		return
	}

	var keyUID int64
	for _, key := range keys.AccessKeys {
		if key.AccessKeyName == data.AccessKeyName.ValueString() {
			keyUID = key.AccessKeyUID
			break
		}
	}
	if keyUID == 0 {
		resp.Diagnostics.AddError("No matching key", "no key with given name")
		return
	}

	keyVersions, err := client.ListAccessKeyVersions(ctx, cloudaccess.ListAccessKeyVersionsRequest{AccessKeyUID: keyUID})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("reading %s", ErrCloudAccessKeyVersions), err.Error())
		return
	}

	for _, version := range keyVersions.AccessKeyVersions {
		versionModel := accessKeyVersionModel{
			CreatedBy:        types.StringValue(version.CreatedBy),
			DeploymentStatus: types.StringValue(string(version.DeploymentStatus)),
			Version:          types.Int64Value(version.Version),
			VersionGUID:      types.StringValue(version.VersionGUID),
		}
		if version.CloudAccessKeyID != nil {
			versionModel.CloudAccessKeyID = types.StringValue(*version.CloudAccessKeyID)
		}

		dateString, err := date.ToString(version.CreatedTime)
		if err != nil {
			resp.Diagnostics.AddError("error parsing date:", err.Error())
			return
		}
		versionModel.CreatedTime = types.StringValue(dateString)

		data.AccessKeyVersions = append(data.AccessKeyVersions, versionModel)
	}
	data.AccessKeyUID = types.Int64Value(keyUID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
