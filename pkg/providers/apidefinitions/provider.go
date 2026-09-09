// Package apidefinitions contains implementation for Akamai Terraform sub-provider responsible for managing API Definitions
package apidefinitions

import (
	"github.com/akamai/terraform-provider-akamai/v11/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// SubProvider gathers apidefinitions resources and data sources
	SubProvider struct {
		config subproviderConfig
	}

	subproviderConfig struct {
		activation activationResourceConfig
		api        apiResourceConfig
	}
)

func defaultSubproviderConfig() subproviderConfig {
	return subproviderConfig{
		activation: defaultActivationResourceConfig(),
		api:        defaultAPIResourceConfig(),
	}
}

var _ subprovider.Subprovider = &SubProvider{}

func newSubproviderWithConfig(config subproviderConfig) *SubProvider {
	return &SubProvider{config: config}
}

// NewSubprovider returns a new apidefinitions subprovider
func NewSubprovider() *SubProvider {
	return newSubproviderWithConfig(defaultSubproviderConfig())
}

// SDKResources returns the apidefinitions resources implemented using terraform-plugin-sdk
func (p *SubProvider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the apidefinitions resources implemented using terraform-plugin-framework
func (p *SubProvider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewActivationResource(p.config.activation),
		NewAPIResource(p.config.api),
		NewAPIResourceOperationResource,
	}
}

// SDKDataSources returns the apidefinitions data sources implemented using terraform-plugin-sdk
func (p *SubProvider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkDataSources returns the apidefinitions data sources implemented using terraform-plugin-framework
func (p *SubProvider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOpenAPIDataSource,
		NewResourceOperationsDataSource,
		NewAPIDataSource,
	}
}
