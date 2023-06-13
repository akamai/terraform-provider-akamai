// Package imaging contains implementation for Akamai Terraform sub-provider responsible for managing images and videos
package imaging

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers imaging resources and data sources
	Subprovider struct {
		client imaging.Imaging
	}

	option func(p *Subprovider)
)

var (
	once sync.Once

	inst *Subprovider
)

var _ subprovider.Plugin = &Subprovider{}

// NewSubprovider returns a core sub provider
func NewSubprovider(opts ...option) *Subprovider {
	once.Do(func() {
		inst = &Subprovider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

func withClient(i imaging.Imaging) option {
	return func(p *Subprovider) {
		p.client = i
	}
}

// Client returns the Imaging interface
func (p *Subprovider) Client(meta meta.Meta) imaging.Imaging {
	if p.client != nil {
		return p.client
	}
	return imaging.Client(meta.Session())
}

// Resources returns terraform resources for imaging
func (p *Subprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_imaging_policy_image": resourceImagingPolicyImage(),
		"akamai_imaging_policy_set":   resourceImagingPolicySet(),
		"akamai_imaging_policy_video": resourceImagingPolicyVideo(),
	}
}

// DataSources returns terraform data sources for imaging
func (p *Subprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_imaging_policy_image": dataImagingPolicyImage(),
		"akamai_imaging_policy_video": dataImagingPolicyVideo(),
	}
}
