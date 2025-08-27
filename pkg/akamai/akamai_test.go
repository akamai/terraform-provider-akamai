package akamai_test

import (
	"context"
	"os"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMain(m *testing.M) {
	const testDefaultEdgercPath = "./testdata/edgerc"

	oldConfigFilePath := akamai.DefaultConfigFilePath
	akamai.DefaultConfigFilePath = testDefaultEdgercPath

	exitCode := m.Run()

	akamai.DefaultConfigFilePath = oldConfigFilePath

	os.Exit(exitCode)
}

// dummy subprovider has only one datasource: "akamai_dummy"
// that itself does nothing and is used only for test purposes
// to trigger provider configuration
type dummy struct {
}

// SDKDataSources implements subprovider.Subprovider.
func (dummy) SDKDataSources() map[string]*schema.Resource {
	return nil
}

// SDKResources implements subprovider.Subprovider.
func (dummy) SDKResources() map[string]*schema.Resource {
	return nil
}

var _ subprovider.Subprovider = dummy{}

// FrameworkDataSources implements subprovider.Subprovider.
func (dummy) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return dummyDataSource{}
		},
	}
}

// FrameworkResources implements subprovider.Subprovider.
func (dummy) FrameworkResources() []func() resource.Resource {
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
	resp.Schema = dataschema.Schema{
		Attributes: map[string]dataschema.Attribute{
			"id": dataschema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
		},
	}
}
