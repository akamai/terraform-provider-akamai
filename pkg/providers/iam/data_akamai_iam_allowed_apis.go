package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &allowedAPIsDataSource{}
	_ datasource.DataSourceWithConfigure = &allowedAPIsDataSource{}

	// ErrIAMListAllowedAPIs is returned when ListAllowedAPIs fails.
	ErrIAMListAllowedAPIs = errors.New("IAM list allowed APIs failed")
)

type (
	allowedAPIsDataSource struct {
		meta meta.Meta
	}

	allowedAPIsModel struct {
		Username           types.String `tfsdk:"username"`
		ClientType         types.String `tfsdk:"client_type"`
		AllowAccountSwitch types.Bool   `tfsdk:"allow_account_switch"`
		AllowedAPIs        []apiModel   `tfsdk:"allowed_apis"`
	}

	apiModel struct {
		AccessLevels      []types.String `tfsdk:"access_levels"`
		APIID             types.Int64    `tfsdk:"api_id"`
		APIName           types.String   `tfsdk:"api_name"`
		Description       types.String   `tfsdk:"description"`
		DocumentationURL  types.String   `tfsdk:"documentation_url"`
		Endpoint          types.String   `tfsdk:"endpoint"`
		HasAccess         types.Bool     `tfsdk:"has_access"`
		ServiceProviderID types.Int64    `tfsdk:"service_provider_id"`
	}
)

// NewAllowedAPIsDataSource returns a new iam allowed APIs data source
func NewAllowedAPIsDataSource() datasource.DataSource {
	return &allowedAPIsDataSource{}
}

func (d *allowedAPIsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_allowed_apis"
}

func (d *allowedAPIsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *allowedAPIsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management allowed APIs",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Unique username for each user.",
			},
			"client_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filters data by client type, either USER_CLIENT, SERVICE_ACCOUNT, or default CLIENT.",
			},
			"allow_account_switch": schema.BoolAttribute{
				Optional:    true,
				Description: "Includes account a user can switch to, false by default.",
			},
			"allowed_apis": schema.ListNestedAttribute{
				Description: "List of available APIs for the user.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_levels": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "API access levels, possible values are READ-ONLY, READ-WRITE, CREDENTIAL-READ-ONLY and CREDENTIAL-READ-WRITE.",
						},
						"api_id": schema.Int64Attribute{
							Computed:    true,
							Description: "A unique identifier for each API.",
						},
						"api_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the API.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A human-readable name for the API.",
						},
						"documentation_url": schema.StringAttribute{
							Computed:    true,
							Description: "A link to more information about the API.",
						},
						"endpoint": schema.StringAttribute{
							Computed:    true,
							Description: "Specifies where the API can access resources.",
						},
						"has_access": schema.BoolAttribute{
							Computed:    true,
							Description: "Confirms access to the API.",
						},
						"service_provider_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier for the API's service provider.",
						},
					},
				},
			},
		},
	}
}

func (d *allowedAPIsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Allowed APIs DataSource Read")

	var data allowedAPIsModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(d.meta)

	apis, err := client.ListAllowedAPIs(ctx, iam.ListAllowedAPIsRequest{
		UserName:           data.Username.ValueString(),
		ClientType:         iam.ClientType(data.ClientType.ValueString()),
		AllowAccountSwitch: data.AllowAccountSwitch.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s:", ErrIAMListAllowedAPIs), err.Error())
		return
	}

	data.read(apis)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *allowedAPIsModel) read(allowedAPIs iam.ListAllowedAPIsResponse) {
	var apis []apiModel
	for _, api := range allowedAPIs {
		accessLevels := make([]types.String, 0, len(api.AccessLevels))
		for _, accessLevel := range api.AccessLevels {
			accessLevels = append(accessLevels, types.StringValue(string(accessLevel)))
		}

		a := apiModel{
			AccessLevels:      accessLevels,
			APIID:             types.Int64Value(api.APIID),
			APIName:           types.StringValue(api.APIName),
			Description:       types.StringValue(api.Description),
			DocumentationURL:  types.StringValue(api.DocumentationURL),
			Endpoint:          types.StringValue(api.Endpoint),
			HasAccess:         types.BoolValue(api.HasAccess),
			ServiceProviderID: types.Int64Value(api.ServiceProviderID),
		}
		apis = append(apis, a)
	}

	d.AllowedAPIs = apis
}
