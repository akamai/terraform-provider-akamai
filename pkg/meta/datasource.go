package meta

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSource is the base struct for all data sources.
type DataSource struct {
	Client edgegrid.Client
}

// Configure implements datasource.DataSourceWithConfigure.
func (r *DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig
		// in framework provider.
		return
	}

	m, ok := req.ProviderData.(Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}
	r.Client = m.Client()
}
