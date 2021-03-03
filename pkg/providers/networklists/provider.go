package networklists

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	network "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client networklists.NETWORKLISTS
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
			"akamai_networklist_activations":              resourceActivations(),
			"akamai_networklist_network_list_description": resourceNetworkListDescription(),
			"akamai_networklist_subscription":             resourceNetworkListSubscription(),
			"akamai_networklist_network_list":             resourceNetworkList(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c networklists.NETWORKLISTS) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the PAPI interface
func (p *provider) Client(meta akamai.OperationMeta) network.NETWORKLISTS {
	if p.client != nil {
		return p.client
	}
	return networklists.Client(meta.Session())
}

func getNETWORKLISTV1Service(d *schema.ResourceData) (interface{}, error) {
	var section string

	for _, s := range tools.FindStringValues(d, "networklist_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	if section != "" {
		d.Set("config_section", section)
	}

	return nil, nil
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

	_, err := getNETWORKLISTV1Service(d)
	if err != nil {
		return nil
	}

	return nil
}
