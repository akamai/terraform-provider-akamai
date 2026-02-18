// Package cloudcertificates contains implementation for Akamai Terraform sub-provider responsible for Cloud Certificate Manager.
package cloudcertificates

import (
	"github.com/akamai/terraform-provider-akamai/v10/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers CloudCertificates resources and data sources.
	Subprovider struct {
		config subproviderConfig
	}

	subproviderConfig struct {
		certificate       certificateResourceConfig
		uploadCertificate uploadSignedCertificateResourceConfig
	}
)

var (
	_ subprovider.Subprovider = &Subprovider{}
)

func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		certificate:       defaultCertificateResourceConfig(),
		uploadCertificate: defaultUploadSignedCertificateResourceConfig(),
	}
}

func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{
		config: config,
	}
}

// NewSubprovider returns a new CloudCertificates subprovider.
func NewSubprovider() *Subprovider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
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
		NewCertificateResource(p.config.certificate),
		NewUploadSignedCertificateResource(p.config.uploadCertificate),
	}
}

// FrameworkDataSources returns the CloudCertificates data sources implemented using terraform-plugin-framework.
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCertificateDataSource,
		NewCertificatesDataSource,
		NewCloudCertificatesHostnameBindingsDataSource,
	}
}
