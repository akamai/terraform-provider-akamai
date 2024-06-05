// Package cloudaccess contains implementation for Akamai Terraform sub-provider responsible for cloud access manager
package cloudaccess

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers cloudaccess resources and data sources
	Subprovider struct{}
)

var (
	_      subprovider.Subprovider = &Subprovider{}
	client cloudaccess.CloudAccess
)

// NewSubprovider returns a new cloudaccess subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// Client returns the gtm interface
func Client(meta meta.Meta) cloudaccess.CloudAccess {
	if client != nil {
		return client
	}
	return cloudaccess.Client(meta.Session())
}

// SDKResources returns the cloudaccess resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the cloudaccess data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the cloudaccess resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the cloudaccess data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKeyDataSource,
		NewKeysDataSource,
		NewKeyPropertiesDataSource,
	}
}
