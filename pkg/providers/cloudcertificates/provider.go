// Package cloudcertificates contains implementation for Akamai Terraform sub-provider responsible for Cloud Certificate Manager.
package cloudcertificates

import (
	"github.com/akamai/terraform-provider-akamai/v11/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers CloudCertificates resources and data sources.
	Subprovider struct {
		config subproviderConfig
	}

	// subproviderConfig holds configuration for all CloudCertificates resources and data sources.
	// Keeping this struct allows proper parallelisation of client construction per resource/data source.
	subproviderConfig struct {
		activation activationResourceConfig
		upload     uploadResourceConfig
	}
)

var _ subprovider.Subprovider = &Subprovider{}

func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activation: defaultActivationResourceConfig(),
		upload:     defaultUploadResourceConfig(),
	}
}

func newSubproviderWithConfig(config subproviderConfig) *Subprovider {
	return &Subprovider{config: config}
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
		NewActivationResource(p.config.activation),
		NewLineageResource,
		NewUploadResource(p.config.upload),
	}
}

// FrameworkDataSources returns the CloudCertificates data sources implemented using terraform-plugin-framework.
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewActivationStatusDataSource,
		NewActivationsDataSource,
		NewArchivedGenerationsDataSource,
		NewBindingsDataSource,
		NewCertificatesActivityDataSource,
		NewGenerationDataSource,
		NewLineageDataSource,
		NewLineagesDataSource,
	}
}
