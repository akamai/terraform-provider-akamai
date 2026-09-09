package akamai

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

var _ provider.Provider = &TestProvider{}

// TestProvider is a version of Provider used for testing purposes
type TestProvider struct {
	Provider
	client edgegrid.Client
}

// NewTestFrameworkProvider returns a test version of Provider
func NewTestFrameworkProvider(client edgegrid.Client, subproviders ...subprovider.Subprovider) func() provider.Provider {
	return func() provider.Provider {
		return &TestProvider{
			Provider: Provider{
				subproviders: subproviders,
			},
			client: client,
		}
	}
}

// Configure configures provider context at the beginning of the lifecycle
// based on the values user specified in the provider configuration block
func (p *TestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.Provider.Configure(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	if p.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Set",
			"The Akamai client was not set in the test provider.",
		)
		return
	}

	m, ok := resp.DataSourceData.(*meta.OperationMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected resp.DataSourceData to be of type meta.OperationMeta, got %T", resp.DataSourceData),
		)
		return
	}
	m.SetClient(p.client)

	m, ok = resp.ResourceData.(*meta.OperationMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected resp.ResourceData to be of type meta.OperationMeta, got %T", resp.ResourceData),
		)
		return
	}
	m.SetClient(p.client)
}
