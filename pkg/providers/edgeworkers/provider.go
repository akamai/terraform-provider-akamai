// Package edgeworkers contains implementation for Akamai Terraform sub-provider responsible for maintaining EdgeWorkers functions to dynamically manage web traffic
package edgeworkers

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers edgeworkers resources and data sources
	Subprovider struct {
		client edgeworkers.Edgeworkers
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new edgeworkers subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(c edgeworkers.Edgeworkers) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the edgeworkers interface
func (p *Subprovider) Client(meta meta.Meta) edgeworkers.Edgeworkers {
	if p.client != nil {
		return p.client
	}
	return edgeworkers.Client(meta.Session())
}

// SDKResources returns the edgeworkers resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_edgekv":                 resourceEdgeKV(),
		"akamai_edgekv_group_items":     resourceEdgeKVGroupItems(),
		"akamai_edgeworkers_activation": resourceEdgeworkersActivation(),
		"akamai_edgeworker":             resourceEdgeWorker(),
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
		"akamai_edgeworker_activation":      dataSourceEdgeWorkerActivation(),
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
