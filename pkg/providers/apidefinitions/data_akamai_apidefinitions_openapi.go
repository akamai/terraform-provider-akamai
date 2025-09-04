package apidefinitions

import (
	"context"
	"fmt"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &openAPIDataSource{}
	_ datasource.DataSourceWithConfigure = &openAPIDataSource{}
)

type (
	openAPIDataSource struct{}

	openAPIModel struct {
		FilePath    types.String `tfsdk:"file_path"`
		APIFileName types.String `tfsdk:"api_file_name"`
		API         types.String `tfsdk:"api"`
	}
)

// NewOpenAPIDataSource returns new api definition endpoint openapi data source
func NewOpenAPIDataSource() datasource.DataSource {
	return &openAPIDataSource{}
}

// Metadata configures data source's meta information
func (d *openAPIDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_apidefinitions_openapi"
}

// Configure configures data source at the beginning of the lifecycle
func (d *openAPIDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	metaConfig, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	if client == nil {
		client = apidefinitions.Client(metaConfig.Session())
	}
	if clientV0 == nil {
		clientV0 = v0.Client(metaConfig.Session())
	}
}

func (d *openAPIDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Map OpenAPI/Swagger file to API configuration",
		Attributes: map[string]schema.Attribute{
			"file_path": schema.StringAttribute{
				Required:    true,
				Description: "Path to OpenAPI/Swagger file",
				Validators:  []validator.String{validators.NotEmptyString(), validators.FileReadable()},
			},
			"api_file_name": schema.StringAttribute{
				Optional:    true,
				Description: "Main API file name in case of zip archive",
				Validators:  []validator.String{validators.NotEmptyString()},
			},
			"api": schema.StringAttribute{
				Optional:    true,
				Description: "JSON-formatted information about the API configuration",
			},
		},
	}
}

func (d *openAPIDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "API Definitions OpenAPI DataSource Read")

	var data openAPIModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	content, err := os.ReadFile(data.FilePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read file content", err.Error())
		return
	}

	response, err := clientV0.FromOpenAPIFile(ctx, v0.FromOpenAPIFileRequest{
		Content:  content,
		RootFile: data.APIFileName.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Mapping OpenAPI File Failed", err.Error())
		return
	}

	toJSON, err := serializeIndent(response.API)
	if err != nil {
		resp.Diagnostics.AddError("Unable to serialize API state", err.Error())
		return
	}
	data.API = types.StringValue(*toJSON)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
