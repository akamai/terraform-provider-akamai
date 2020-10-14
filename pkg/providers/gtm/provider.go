package gtm

import (
	"sync"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
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

		client gtm.GTM
	}

	// Option is a gtm provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider() akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"gtm_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("gtm_section"),
			},
			"gtm": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("gtm"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_gtm_domain":     resourceGTMv1Domain(),
			"akamai_gtm_property":   resourceGTMv1Property(),
			"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
			"akamai_gtm_resource":   resourceGTMv1Resource(),
			"akamai_gtm_asmap":      resourceGTMv1ASmap(),
			"akamai_gtm_geomap":     resourceGTMv1Geomap(),
			"akamai_gtm_cidrmap":    resourceGTMv1Cidrmap(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c gtm.GTM) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *provider) Client(meta akamai.OperationMeta) gtm.GTM {
	if p.client != nil {
		return p.client
	}
	return gtm.Client(meta.Session())
}

func getConfigGTMV1Service(d *schema.ResourceData) (interface{}, error) {

	var section string

	for _, s := range tools.FindStringValues(d, "gtm_section", "config_section") {
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
	return "gtm"
}

// GTMProviderVersion update version string anytime provider adds new features
const GTMProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return GTMProviderVersion
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

	_, err := getConfigGTMV1Service(d)
	if err != nil {
		return nil
	}
	return nil
}
