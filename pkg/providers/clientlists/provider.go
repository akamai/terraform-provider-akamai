// Package clientlists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of client lists
package clientlists

import (
	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers clientlists resources and data sources
	Subprovider struct {
		config subproviderConfig
	}

	// subproviderConfig aggregates the configuration of all clientlists resources
	// so that polling intervals and other timing values can be overridden,
	// in particular by tests.
	subproviderConfig struct {
		activation clientListActivationConfig
	}
)

var _ subprovider.Subprovider = &Subprovider{}

// defaultSubproviderConfig returns the production defaults for the clientlists subprovider.
func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activation: defaultClientListActivationConfig(),
	}
}

// NewSubprovider returns a new clientlists subprovider
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// newSubproviderWithConfig returns a clientlists subprovider initialized with the given configuration.
func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
}

// SDKResources returns the clientlists resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_clientlist_activation": resourceClientListActivation(p.config.activation),
		"akamai_clientlist_list":       resourceClientList(),
	}
}

// SDKDataSources returns the clientlists data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the clientlists resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the clientlists data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClientListDataSource,
		NewClientListsDataSource,
	}
}
