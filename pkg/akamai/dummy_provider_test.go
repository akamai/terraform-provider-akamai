package akamai_test

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dummy subprovider has only one datasource: "akamai_dummy"
// that itself does nothing and is used only for test purposes
// to trigger provider configuration
type dummy struct {
}

var _ subprovider.Framework = dummy{}

// DataSources implements subprovider.Framework.
func (dummy) DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return dummyDataSource{}
		},
	}
}

// Resources implements subprovider.Framework.
func (dummy) Resources() []func() resource.Resource {
	return nil
}

type dummyDataSource struct{}

type dummyDataSourceModel struct {
	ID types.String `tfsdk:"id"`
}

var _ datasource.DataSource = dummyDataSource{}

// Metadata implements datasource.DataSource.
func (dummyDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_dummy"
}

// Read implements datasource.DataSource.
func (dummyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dummyDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue("1")
	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSource.
func (dummyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
		},
	}
}
