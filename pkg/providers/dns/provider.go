// Package dns contains implementation for Akamai Terraform sub-provider responsible for managing DNS zones configuration
package dns

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new DNS subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}
		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Client returns the DNS interface
func (p *Subprovider) Client(meta meta.Meta) dns.DNS {
	if p.client != nil {
		return p.client
	}
	return dns.Client(meta.Session())
}

// SDKResources returns the DNS resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_dns_zone":   resourceDNSv2Zone(),
		"akamai_dns_record": resourceDNSv2Record(),
	}
}

// SDKDataSources returns the DNS data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_authorities_set": dataSourceAuthoritiesSet(),
		"akamai_dns_record_set":  dataSourceDNSRecordSet(),
	}
}

// FrameworkResources returns the DNS resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the DNS data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewZoneDNSSecStatusDataSource,
	}
}
