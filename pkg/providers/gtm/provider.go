// Package gtm contains implementation for Akamai Terraform sub-provider responsible for managing Global Traffic Management (GTM) domain configuration and administration
package gtm

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// PluginSubprovider gathers gtm resources and data sources
	PluginSubprovider struct {
		client gtm.GTM
	}

	// FrameworkSubprovider gathers property resources and data sources written using terraform-plugin-framework
	FrameworkSubprovider struct {
		client gtm.GTM
	}
)

var _ subprovider.Framework = &FrameworkSubprovider{}
var _ subprovider.Plugin = &PluginSubprovider{}

var (
	oncePlugin, onceFramework sync.Once

	inst *PluginSubprovider

	frameworkInst *FrameworkSubprovider
)

// NewSubprovider returns a new GTM subprovider
func NewSubprovider() *PluginSubprovider {
	oncePlugin.Do(func() {
		inst = &PluginSubprovider{}
	})

	return inst
}

// NewFrameworkSubprovider returns a core Framework based sub provider
func NewFrameworkSubprovider() *FrameworkSubprovider {
	onceFramework.Do(func() {
		frameworkInst = &FrameworkSubprovider{}
	})

	return frameworkInst
}

// Client returns the GTM interface
func (p *PluginSubprovider) Client(meta meta.Meta) gtm.GTM {
	if p.client != nil {
		return p.client
	}
	return gtm.Client(meta.Session())
}

// Client returns the GTM interface
func (f *FrameworkSubprovider) Client(meta meta.Meta) gtm.GTM {
	if f.client != nil {
		return f.client
	}
	return gtm.Client(meta.Session())
}

// Resources returns the GTM resources implemented using terraform-plugin-framework
func (f *FrameworkSubprovider) Resources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns the GTM data sources implemented using terraform-plugin-framework
func (f *FrameworkSubprovider) DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGTMDomainDataSource,
		NewGTMDomainsDataSource,
		NewGTMResourceDataSource,
		NewGTMResourcesDataSource,
	}
}

// Resources returns terraform resources for gtm
func (p *PluginSubprovider) Resources() map[string]*schema.Resource {
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
func (p *PluginSubprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_gtm_datacenter":         dataSourceGTMDatacenter(),
		"akamai_gtm_datacenters":        dataSourceGTMDatacenters(),
		"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
	}
}
