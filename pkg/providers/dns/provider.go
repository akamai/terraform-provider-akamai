// Package dns contains implementation for Akamai Terraform sub-provider responsible for managing DNS zones configuration
package dns

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers dns resources and data sources
	Subprovider struct {
		client dns.DNS
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.SDK = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider() *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}
	})

	return inst
}

func withClient(c dns.DNS) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *Subprovider) Client(meta meta.Meta) dns.DNS {
	if p.client != nil {
		return p.client
	}
	return dns.Client(meta.Session())
}

// Resources returns terraform resources for imadnsging
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_dns_zone":   resourceDNSv2Zone(),
		"akamai_dns_record": resourceDNSv2Record(),
	}
}

// DataSources returns terraform data sources for dns
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_authorities_set": dataSourceAuthoritiesSet(),
		"akamai_dns_record_set":  dataSourceDNSRecordSet(),
	}
}
