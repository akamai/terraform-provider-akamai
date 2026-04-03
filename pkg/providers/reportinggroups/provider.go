// Package reportinggroups contains implementation for Akamai Terraform sub-provider responsible for Reporting Groups
package reportinggroups

import (
	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers reportinggroups resources and data sources
	Subprovider struct {
		config subproviderConfig
	}

	subproviderConfig struct {
	}
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{}
}

func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{
		config: config,
	}
}

// NewSubprovider returns a new reportinggroups subprovider
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// SDKResources returns the reportinggroups resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the reportinggroups data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the reportinggroups resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the reportinggroups data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCPCodeDataSource,
	}
}
