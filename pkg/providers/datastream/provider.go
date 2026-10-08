// Package datastream contains implementation for Akamai Terraform sub-provider responsible for providing scalable, low latency streaming of data
package datastream

import (
	"github.com/akamai/terraform-provider-akamai/v11/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers datastream resources and data sources
	Subprovider struct {
		config subproviderConfig
	}

	// subproviderConfig aggregates the configuration of all datastream resources
	// so that polling intervals and other timing values can be overridden,
	// in particular by tests.
	subproviderConfig struct {
		datastream datastreamResourceConfig
	}
)

var _ subprovider.Subprovider = &Subprovider{}

// defaultSubproviderConfig returns the production defaults for the datastream subprovider.
func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		datastream: defaultDatastreamResourceConfig(),
	}
}

// NewSubprovider returns a new datastream subprovider with the default configuration.
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// newSubproviderWithConfig returns a new datastream subprovider initialized with the given configuration.
func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
}

// SDKResources returns the datastream resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_datastream": resourceDatastream(p.config.datastream),
	}
}

// SDKDataSources returns the datastream data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_datastream_activation_history": dataAkamaiDatastreamActivationHistory(),
		"akamai_datastream_appsec_configs":     dataSourceAppSecConfigs(),
		"akamai_datastream_dataset_fields":     dataSourceDatasetFields(),
		"akamai_datastreams":                   dataAkamaiDatastreamStreams(),
	}
}

// FrameworkResources returns the datastream resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the datastream data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAnswerXServiceIDsDataSource,
	}
}
