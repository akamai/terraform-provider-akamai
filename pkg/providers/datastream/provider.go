// Package datastream contains implementation for Akamai Terraform sub-provider responsible for providing scalable, low latency streaming of data
package datastream

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client datastream.DS
	}

	// Option is a ds provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider(opts ...Option) subprovider.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_datastream_activation_history": dataAkamaiDatastreamActivationHistory(),
			"akamai_datastream_dataset_fields":     dataSourceDatasetFields(),
			"akamai_datastreams":                   dataAkamaiDatastreamStreams(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_datastream": resourceDatastream(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c datastream.DS) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the ds interface
func (p *provider) Client(meta meta.Meta) datastream.DS {
	if p.client != nil {
		return p.client
	}
	return datastream.Client(meta.Session())
}

func (p *provider) Name() string {
	return "datastream"
}

// ProviderVersion update version string anytime provider adds new features
const ProviderVersion string = "v0.0.1"

func (p *provider) Version() string {
	return ProviderVersion
}

func (p *provider) Schema() map[string]*schema.Schema {
	return p.Provider.Schema
}

func (p *provider) Resources() map[string]*schema.Resource {
	return p.Provider.ResourcesMap
}

func (p *provider) DataSources() map[string]*schema.Resource {
	return p.Provider.DataSourcesMap
}

func (p *provider) Configure(_ log.Interface, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}
