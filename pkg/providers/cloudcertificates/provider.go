// Package cloudcertificates contains implementation for Akamai Terraform sub-provider responsible for Cloud Certificate Manager.
package cloudcertificates

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers CloudCertificates resources and data sources.
	Subprovider struct{}
)

var (
	_      subprovider.Subprovider = &Subprovider{}
	client ccm.CCM
)

// NewSubprovider returns a new CloudCertificates subprovider.
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// Client returns the CCM interface.
func Client(meta meta.Meta) ccm.CCM {
	if client != nil {
		return client
	}
	return ccm.Client(meta.Session())
}

// SDKResources returns the CloudCertificates resources implemented using terraform-plugin-sdk.
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the CloudCertificates data sources implemented using terraform-plugin-sdk.
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the CloudCertificates resources implemented using terraform-plugin-framework.
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewCertificateResource,
	}
}

// FrameworkDataSources returns the CloudCertificates data sources implemented using terraform-plugin-framework.
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
