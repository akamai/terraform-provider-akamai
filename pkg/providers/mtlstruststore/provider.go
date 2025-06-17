// Package mtlstruststore contains implementation for Akamai Terraform sub-provider responsible for MTLS Truststore
package mtlstruststore

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// Subprovider gathers MTLS Truststore resources and data sources
	Subprovider struct{}
)

var (
	_      subprovider.Subprovider = &Subprovider{}
	client mtlstruststore.MTLSTruststore
)

// NewSubprovider returns a new MTLS Truststore subprovider
func NewSubprovider() *Subprovider {
	return &Subprovider{}
}

// Client returns the MTLS Truststore interface
func Client(meta meta.Meta) mtlstruststore.MTLSTruststore {
	if client != nil {
		return client
	}
	return mtlstruststore.Client(meta.Session())
}

// SDKResources returns the MTLS Truststore resources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SDKDataSources returns the MTLS Truststore data sources implemented using terraform-plugin-sdk
func (p *Subprovider) SDKDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// FrameworkResources returns the MTLS Truststore resources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

// FrameworkDataSources returns the MTLS Truststore data sources implemented using terraform-plugin-framework
func (p *Subprovider) FrameworkDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCASetActivationDataSource,
		NewCASetActivitiesDataSource,
		NewCASetAssociationsDataSource,
		NewCASetDataSource,
		NewCASetsDataSource,
		NewCASetVersionsDataSource,
	}
}
