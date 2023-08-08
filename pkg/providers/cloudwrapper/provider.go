// Package cloudwrapper contains implementation for Akamai Terraform sub-provider responsible for cloud wrapper
package cloudwrapper

import (
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type (
	// Subprovider gathers cloud wrapper resources and data sources written using terraform-plugin-framework
	Subprovider struct{}
)

var (
	_ subprovider.Framework = &Subprovider{}
)

// NewSubprovider returns a core Framework based sub provider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// Resources returns terraform resources for cloudwrapper
func (p *Subprovider) Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewConfigurationResource,
	}
}

// DataSources returns terraform data sources for cloudwrapper
func (p *Subprovider) DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConfigurationDataSource,
		NewLocationDataSource,
		NewLocationsDataSource,
		NewPropertiesDataSource,
	}
}
