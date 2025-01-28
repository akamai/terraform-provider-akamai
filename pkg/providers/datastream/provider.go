// Package datastream contains implementation for Akamai Terraform sub-provider responsible for providing scalable, low latency streaming of data
package datastream

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers datastream resources and data sources
	Subprovider struct {
		client datastream.DS
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new datastream subprovider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Client returns the DS interface
func (p *Subprovider) Client(meta meta.Meta) datastream.DS {
	if p.client != nil {
		return p.client
	}
	return datastream.Client(meta.Session())
}

// SDKResources returns the datastream resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_datastream": resourceDatastream(),
	}
}

// SDKDataSources returns the datastream data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_datastream_activation_history": dataAkamaiDatastreamActivationHistory(),
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
	return []func() datasource.DataSource{}
}
