// Package gtm contains implementation for Akamai Terraform sub-provider responsible for managing Global Traffic Management (GTM) domain configuration and administration
package gtm

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers gtm resources and data sources
	Subprovider struct {
		client gtm.GTM
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider() *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}
	})

	return inst
}

func withClient(c gtm.GTM) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *Subprovider) Client(meta meta.Meta) gtm.GTM {
	if p.client != nil {
		return p.client
	}
	return gtm.Client(meta.Session())
}

// Resources returns terraform resources for gtm
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_gtm_domain":     resourceGTMv1Domain(),
		"akamai_gtm_property":   resourceGTMv1Property(),
		"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
		"akamai_gtm_resource":   resourceGTMv1Resource(),
		"akamai_gtm_asmap":      resourceGTMv1ASmap(),
		"akamai_gtm_geomap":     resourceGTMv1Geomap(),
		"akamai_gtm_cidrmap":    resourceGTMv1Cidrmap(),
	}
}

// DataSources returns terraform data sources for gtm
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_gtm_datacenter":         dataSourceGTMDatacenter(),
		"akamai_gtm_datacenters":        dataSourceGTMDatacenters(),
		"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
	}
}
