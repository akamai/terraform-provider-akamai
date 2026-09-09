// Package edgeworkers contains implementation for Akamai Terraform sub-provider responsible for maintaining EdgeWorkers functions to dynamically manage web traffic
package edgeworkers

import (
	"github.com/akamai/terraform-provider-akamai/v11/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers edgeworkers resources and data sources
	Subprovider struct {
		config subproviderConfig
	}

	subproviderConfig struct {
		activation       edgeworkersActivationResourceConfig
		edgekv           edgeKVResourceConfig
		edgekvGroupItems edgeKVGroupItemsResourceConfig
		edgeworker       edgeworkerResourceConfig
	}
)

var _ subprovider.Subprovider = &Subprovider{}

func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activation:       defaultEdgeworkersActivationResourceConfig(),
		edgekv:           defaultEdgeKVResourceConfig(),
		edgekvGroupItems: defaultEdgeKVGroupItemsResourceConfig(),
		edgeworker:       defaultEdgeworkerResourceConfig(),
	}
}

func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
}

// NewSubprovider returns a new edgeworkers subprovider
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// SDKResources returns the edgeworkers resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_edgekv":                 resourceEdgeKV(p.config.edgekv),
		"akamai_edgekv_group_items":     resourceEdgeKVGroupItems(p.config.edgekvGroupItems),
		"akamai_edgeworkers_activation": resourceEdgeworkersActivation(p.config.activation),
		"akamai_edgeworker":             resourceEdgeWorker(p.config.edgeworker, p.config.activation),
	}
}

// SDKDataSources returns the edgeworkers data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_edgekv_group_items":         dataSourceEdgeKVGroupItems(),
		"akamai_edgekv_groups":              dataSourceEdgeKVGroups(),
		"akamai_edgeworkers_resource_tier":  dataSourceEdgeworkersResourceTier(),
		"akamai_edgeworkers_property_rules": dataSourceEdgeworkersPropertyRules(),
		"akamai_edgeworker":                 dataSourceEdgeWorker(),
		"akamai_edgeworker_activation":      dataSourceEdgeWorkerActivation(p.config.activation),
	}
}

// FrameworkResources returns the edgeworkers resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the edgeworkers data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
