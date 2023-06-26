// Package edgeworkers contains implementation for Akamai Terraform sub-provider responsible for maintaining EdgeWorkers functions to dynamically manage web traffic
package edgeworkers

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
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

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
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

// Resources returns terraform resources for edgeworkers
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_edgekv":                 resourceEdgeKV(),
		"akamai_edgekv_group_items":     resourceEdgeKVGroupItems(),
		"akamai_edgeworkers_activation": resourceEdgeworkersActivation(),
		"akamai_edgeworker":             resourceEdgeWorker(),
	}
}

// DataSources returns terraform data sources for edgeworkers
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_edgekv_group_items":         dataSourceEdgeKVGroupItems(),
		"akamai_edgekv_groups":              dataSourceEdgeKVGroups(),
		"akamai_edgeworkers_resource_tier":  dataSourceEdgeworkersResourceTier(),
		"akamai_edgeworkers_property_rules": dataSourceEdgeworkersPropertyRules(),
		"akamai_edgeworker":                 dataSourceEdgeWorker(),
		"akamai_edgeworker_activation":      dataSourceEdgeWorkerActivation(),
	}
}
