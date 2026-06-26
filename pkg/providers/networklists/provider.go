// Package networklists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of network lists
package networklists

import (
	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers networklists resources and data sources
	Subprovider struct {
		config subproviderConfig
	}

	// subproviderConfig aggregates the configuration of all networklists resources
	// so that polling intervals and other timing values can be overridden,
	// in particular by tests.
	subproviderConfig struct {
		activations resourceActivationsConfig
	}
)

var _ subprovider.Subprovider = &Subprovider{}

// defaultSubproviderConfig returns the production defaults for the networklists subprovider.
func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activations: defaultResourceActivationsConfig(),
	}
}

// NewSubprovider returns a new networklists subprovider with the default configuration.
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// newSubproviderWithConfig returns a new networklists subprovider initialized with the given configuration.
func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
}

// SDKResources returns the networklists resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_networklist_activations":  resourceActivations(p.config.activations),
		"akamai_networklist_description":  resourceNetworkListDescription(),
		"akamai_networklist_subscription": resourceNetworkListSubscription(),
		"akamai_networklist_network_list": resourceNetworkList(),
	}
}

// SDKDataSources returns the networklists data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_networklist_network_lists": dataSourceNetworkList(),
	}
}

// FrameworkResources returns the networklists resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the networklists data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
