// Package imaging contains implementation for Akamai Terraform sub-provider responsible for managing images and videos
package imaging

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new imaging subprovider
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

// SDKResources returns the imaging resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_imaging_policy_image": resourceImagingPolicyImage(),
		"akamai_imaging_policy_set":   resourceImagingPolicySet(),
		"akamai_imaging_policy_video": resourceImagingPolicyVideo(),
	}
}

// SDKDataSources returns the imaging data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_imaging_policy_image": dataImagingPolicyImage(),
		"akamai_imaging_policy_video": dataImagingPolicyVideo(),
	}
}

// FrameworkResources returns the imaging resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the imaging data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
