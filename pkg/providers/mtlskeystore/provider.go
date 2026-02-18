// Package mtlskeystore contains implementation for Akamai Terraform sub-provider responsible for MTLS Keystore
package mtlskeystore

import (
	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers MTLS Keystore resources and data sources
	Subprovider struct{}
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

// NewSubprovider returns a new MTLS Keystore subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// SDKResources returns the MTLS Keystore resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the MTLS Keystore data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the MTLS Keystore resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewClientCertificateThirdPartyResource,
		NewClientCertificateUploadResource,
		NewClientCertificateAkamaiResource,
	}
}

// FrameworkDataSources returns the MTLS Keystore data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountCACertificatesDataSource,
		NewClientCertificatesDataSource,
		NewClientCertificateDataSource,
	}
}
