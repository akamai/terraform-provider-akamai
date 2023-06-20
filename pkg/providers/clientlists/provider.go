// Package clientlists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of client lists
package clientlists

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/apex/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client clientlists.ClientLists
	}
	// Option is a clientlists provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider(opts ...Option) akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_clientlist_lists": dataSourceClientLists(),
		},
		ResourcesMap: map[string]*schema.Resource{},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c clientlists.ClientLists) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the CLIENTLISTS interface
func (p *provider) Client(meta akamai.OperationMeta) clientlists.ClientLists {
	if p.client != nil {
		return p.client
	}
	return clientlists.Client(meta.Session())
}

// Name return provider name
func (p *provider) Name() string {
	return "clientlists"
}

// Version returns sub-provider version
func (p *provider) Version() string {
	return "v1.0.0"
}

func (p *provider) Schema() map[string]*schema.Schema {
	return p.Provider.Schema
}

func (p *provider) Resources() map[string]*schema.Resource {
	return p.Provider.ResourcesMap
}

func (p *provider) DataSources() map[string]*schema.Resource {
	return p.Provider.DataSourcesMap
}

func (p *provider) Configure(log log.Interface, _ *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")

	return nil
}
