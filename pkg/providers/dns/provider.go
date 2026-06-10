// Package dns contains implementation for Akamai Terraform sub-provider responsible for managing DNS zones configuration
package dns

import (
	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers dns resources and data sources
	Subprovider struct {
		config subproviderConfig
	}

	// subproviderConfig aggregates the configuration of all DNS resources
	// so that polling intervals and other timing values can be overridden,
	// in particular by tests.
	subproviderConfig struct {
		zone   dnsZoneResourceConfig
		record dnsRecordResourceConfig
	}
)

var _ subprovider.Subprovider = &Subprovider{}

// defaultSubproviderConfig returns the production defaults for the DNS subprovider.
func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		zone:   defaultDNSZoneResourceConfig(),
		record: defaultDNSRecordResourceConfig(),
	}
}

// NewSubprovider returns a new DNS subprovider with the default configuration.
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// newSubproviderWithConfig returns a new DNS subprovider initialized with the
// given configuration.
func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
}

// SDKResources returns the DNS resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_dns_zone":   resourceDNSv2Zone(p.config.zone),
		"akamai_dns_record": resourceDNSv2Record(p.config.record),
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
		func() datasource.DataSource { return newZoneDNSSecStatusDataSource() },
	}
}
