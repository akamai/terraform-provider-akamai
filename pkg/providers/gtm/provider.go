// Package gtm contains implementation for Akamai Terraform sub-provider responsible for managing Global Traffic Management (GTM) domain configuration and administration
package gtm

import (
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers gtm resources and data sources
	Subprovider struct {
		config subproviderConfig
	}
	subproviderConfig struct {
		domain          gtmDomainResourceConfig
		property        gtmPropertyResourceConfig
		defaultInterval time.Duration
	}
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		domain:          defaultGTMDomainResourceConfig(),
		property:        defaultGTMPropertyResourceConfig(),
		defaultInterval: 5 * time.Second,
	}
}

func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
}

// NewSubprovider returns a new gtm subprovider
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// FrameworkResources returns the gtm resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the gtm data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGTMASMapDataSource,
		NewGTMCIDRMapDataSource,
		NewGTMDomainDataSource,
		NewGTMDomainsDataSource,
		NewGTMGeoMapDataSource,
		NewGTMGeoMapsDataSource,
		NewGTMResourceDataSource,
		NewGTMResourcesDataSource,
	}
}

// SDKResources returns the gtm resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_gtm_domain":     resourceGTMv1Domain(p.config.domain),
		"akamai_gtm_property":   resourceGTMv1Property(p.config.property),
		"akamai_gtm_datacenter": resourceGTMv1Datacenter(p.config.defaultInterval),
		"akamai_gtm_resource":   resourceGTMv1Resource(p.config.defaultInterval),
		"akamai_gtm_asmap":      resourceGTMv1ASMap(p.config.defaultInterval),
		"akamai_gtm_geomap":     resourceGTMv1GeoMap(p.config.defaultInterval),
		"akamai_gtm_cidrmap":    resourceGTMv1CIDRMap(p.config.defaultInterval),
	}
}

// SDKDataSources returns the gtm data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_gtm_datacenter":         dataSourceGTMDatacenter(),
		"akamai_gtm_datacenters":        dataSourceGTMDatacenters(),
		"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
	}
}
