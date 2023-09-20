// Package subprovider defines contract for a subprovider
package subprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SDK is the interface implemented by the sub-providers using terraform-plugin-sdk
type SDK interface {
	// Resources returns the resources for the sub-provider
	Resources() map[string]*schema.Resource

	// DataSources returns the datasources for the sub-provider
	DataSources() map[string]*schema.Resource
}

// Framework is the interface implemented by the sub-providers using terraform-plugin-framework
type Framework interface {
	// Resources returns the resources for the sub-provider
	Resources() []func() resource.Resource

	// DataSources returns the datasources for the sub-provider
	DataSources() []func() datasource.DataSource
}
