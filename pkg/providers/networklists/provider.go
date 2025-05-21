// Package networklists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of network lists
package networklists

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers networklists resources and data sources
	Subprovider struct {
		client networklists.NetworkList
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns new networklists subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Client returns the NetworkList interface
func (p *Subprovider) Client(meta meta.Meta) networklists.NetworkList {
	if p.client != nil {
		return p.client
	}
	return networklists.Client(meta.Session())
}

// SDKResources returns the networklists resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_networklist_activations":  resourceActivations(),
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
