// Package networklists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of network lists
package networklists

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers networklists resources and data sources
	Subprovider struct {
		client networklists.NTWRKLISTS
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Plugin = &Subprovider{}

func newSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(c networklists.NTWRKLISTS) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *Subprovider) Client(meta meta.Meta) networklists.NTWRKLISTS {
	if p.client != nil {
		return p.client
	}
	return networklists.Client(meta.Session())
}

// Resources returns terraform resources for networklists
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_networklist_activations":  resourceActivations(),
		"akamai_networklist_description":  resourceNetworkListDescription(),
		"akamai_networklist_subscription": resourceNetworkListSubscription(),
		"akamai_networklist_network_list": resourceNetworkList(),
	}
}

// DataSources returns terraform data sources for networklists
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_networklist_network_lists": dataSourceNetworkList(),
	}
}
