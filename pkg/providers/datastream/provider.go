// Package datastream contains implementation for Akamai Terraform sub-provider responsible for providing scalable, low latency streaming of data
package datastream

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
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

var _ subprovider.Plugin = &Subprovider{}

func newSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(c datastream.DS) option {
	return func(p *Subprovider) {
		p.client = c
	}
}

// Client returns the ds interface
func (p *Subprovider) Client(meta meta.Meta) datastream.DS {
	if p.client != nil {
		return p.client
	}
	return datastream.Client(meta.Session())
}

// Resources returns terraform resources for datastream
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_datastream": resourceDatastream(),
	}
}

// DataSources returns terraform data sources for datastream
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_datastream_activation_history": dataAkamaiDatastreamActivationHistory(),
		"akamai_datastream_dataset_fields":     dataSourceDatasetFields(),
		"akamai_datastreams":                   dataAkamaiDatastreamStreams(),
	}
}
