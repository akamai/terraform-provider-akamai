// Package clientlists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of client lists
package clientlists

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers clientlists resources and data sources
	Subprovider struct {
		client clientlists.ClientLists
	}
	// Option is a clientlists provider option
	Option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

// NewSubprovider returns a core sub provider
func NewSubprovider(opts ...Option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c clientlists.ClientLists) Option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the CLIENTLISTS interface
func (p *Subprovider) Client(meta meta.Meta) clientlists.ClientLists {
	if p.client != nil {
		return p.client
	}
	return clientlists.Client(meta.Session())
}

// Resources returns terraform resources for clientlists
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_clientlist_list": resourceClientList(),
	}
}

// DataSources returns terraform data sources for clientlists
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_clientlist_lists": dataSourceClientLists(),
	}
}
