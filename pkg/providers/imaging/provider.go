// Package imaging contains implementation for Akamai Terraform sub-provider responsible for managing images and videos
package imaging

import (
	"github.com/akamai/terraform-provider-akamai/v11/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Subprovider gathers imaging resources and data sources
type Subprovider struct{}

var _ subprovider.Subprovider = &Subprovider{}

// NewSubprovider returns a new imaging subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
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
