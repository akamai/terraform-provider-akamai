package akamai

import (
	"context"

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

// Configure sets a meta.OperationMeta wired to the mock client directly, skipping Provider.Configure entirely -
// parsing the provider config block and building a real session/client there would only be thrown away here.
func (p *TestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Set",
			"The Akamai client was not set in the test provider.",
		)
		return
	}

	m := meta.NewTestMeta(p.client)
	resp.DataSourceData = m
	resp.ResourceData = m
	resp.ActionData = m
}
