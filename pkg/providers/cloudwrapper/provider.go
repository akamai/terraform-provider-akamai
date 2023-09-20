// Package cloudwrapper contains implementation for Akamai Terraform sub-provider responsible for cloud wrapper
package cloudwrapper

import (
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers cloudwrapper resources and data sources
	Subprovider struct{}
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

// NewSubprovider returns a new cloudwrapper subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// SDKResources returns the cloudwrapper resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the cloudwrapper data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the cloudwrapper resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewActivationResource,
		NewConfigurationResource,
	}
}

// FrameworkDataSources returns the cloudwrapper data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCapacitiesDataSource,
		NewConfigurationDataSource,
		NewConfigurationsDataSource,
		NewLocationDataSource,
		NewLocationsDataSource,
		NewPropertiesDataSource,
	}
}
