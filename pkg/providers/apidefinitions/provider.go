// Package apidefinitions contains implementation for Akamai Terraform sub-provider responsible for managing API Definitions
package apidefinitions

import (
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// SubProvider gathers apidefinitions resources and data sources
	SubProvider struct{}
	option      func(p *SubProvider)
)

var (
	once     sync.Once
	client   apidefinitions.APIDefinitions
	clientV0 v0.APIDefinitions
	inst     *SubProvider
)

var _ subprovider.Subprovider = &SubProvider{}

// NewSubprovider returns a new apidefinitions subprovider
func NewSubprovider(opts ...option) *SubProvider {
	once.Do(func() {
		inst = &SubProvider{}

		for _, opt := range opts {
			opt(inst)
		}
	})

	return inst
}

// SDKResources returns the apidefinitions resources implemented using terraform-plugin-sdk
func (p *SubProvider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the apidefinitions resources implemented using terraform-plugin-framework
func (p *SubProvider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewActivationResource,
		NewAPIResource,
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
	}
}
