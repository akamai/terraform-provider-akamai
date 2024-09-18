// Package testprovider is a package containing a terraform provider for internal testing purposes
package testprovider

import (
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewMockSubprovider return the test subprovider used for internal testing
func NewMockSubprovider() subprovider.Subprovider {
	return &mockSubprovider{}
}

type (
	mockSubprovider struct{}
)

// SDKResources returns the test resources implemented using terraform-plugin-sdk
func (p *mockSubprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the test data sources implemented using terraform-plugin-sdk
func (p *mockSubprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the test resources implemented using terraform-plugin-framework
func (p *mockSubprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewTestResource,
	}
}

// FrameworkDataSources returns the test data sources implemented using terraform-plugin-framework
func (p *mockSubprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
