package imaging

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client imaging.Imaging
	}

	// Option is an imaging provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider(opts ...Option) akamai.Subprovider {
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
		Schema:         map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_imaging_policy_image": resourceImagingPolicyImage(),
			"akamai_imaging_policy_set":   resourceImagingPolicySet(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(i imaging.Imaging) Option {
	return func(p *provider) {
		p.client = i
	}
}

// Client returns the Imaging interface
func (p *provider) Client(meta akamai.OperationMeta) imaging.Imaging {
	if p.client != nil {
		return p.client
	}
	return imaging.Client(meta.Session())
}

func (p *provider) Name() string {
	return "imaging"
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
