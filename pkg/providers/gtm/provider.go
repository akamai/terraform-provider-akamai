// Package gtm contains implementation for Akamai Terraform sub-provider responsible for managing Global Traffic Management (GTM) domain configuration and administration
package gtm

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
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
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_gtm_datacenter":         dataSourceGTMDatacenter(),
			"akamai_gtm_datacenters":        dataSourceGTMDatacenters(),
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

func (p *provider) Configure(log log.Interface, _ *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")
	return nil
}
