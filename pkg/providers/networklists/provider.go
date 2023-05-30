// Package networklists contains implementation for Akamai Terraform sub-provider responsible for creation, deployment, and management of network lists
package networklists

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/config"
	"github.com/apex/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client networklists.NTWRKLISTS
	}
	// Option is a networklist provider option
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
		Schema: map[string]*schema.Schema{
			"networklist_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("networklist_section"),
			},
			"network": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("network"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_networklist_network_lists": dataSourceNetworkList(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_networklist_activations":  resourceActivations(),
			"akamai_networklist_description":  resourceNetworkListDescription(),
			"akamai_networklist_subscription": resourceNetworkListSubscription(),
			"akamai_networklist_network_list": resourceNetworkList(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c networklists.NTWRKLISTS) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *provider) Client(meta akamai.OperationMeta) networklists.NTWRKLISTS {
	if p.client != nil {
		return p.client
	}
	return networklists.Client(meta.Session())
}

func getNetworkListV1Service(d *schema.ResourceData) error {
	var section string

	for _, s := range tf.FindStringValues(d, "networklist_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		if err := d.Set("config_section", section); err != nil {
			return err
		}
	}

	return nil
}

func (p *provider) Name() string {
	return "networklists"
}

// NetworkProviderVersion update version string anytime provider adds new features
const NetworkProviderVersion string = "v1.0.0"

func (p *provider) Version() string {
	return NetworkProviderVersion
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

func (p *provider) Configure(log log.Interface, d *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")

	err := getNetworkListV1Service(d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
