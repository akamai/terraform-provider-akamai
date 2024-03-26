// Package gtm contains implementation for Akamai Terraform sub-provider responsible for managing Global Traffic Management (GTM) domain configuration and administration
package gtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers gtm resources and data sources
	Subprovider struct{}
)

var (
	_      subprovider.Subprovider = &Subprovider{}
	client gtm.GTM
)

// NewSubprovider returns a new gtm subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// Client returns the gtm interface
func Client(meta meta.Meta) gtm.GTM {
	if client != nil {
		return client
	}
	return gtm.Client(meta.Session())
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
		NewGTMResourceDataSource,
		NewGTMResourcesDataSource,
	}
}

// SDKResources returns the gtm resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_gtm_domain":     resourceGTMv1Domain(),
		"akamai_gtm_property":   resourceGTMv1Property(),
		"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
		"akamai_gtm_resource":   resourceGTMv1Resource(),
		"akamai_gtm_asmap":      resourceGTMv1ASMap(),
		"akamai_gtm_geomap":     resourceGTMv1GeoMap(),
		"akamai_gtm_cidrmap":    resourceGTMv1CIDRMap(),
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
