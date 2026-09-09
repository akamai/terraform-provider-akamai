// Package subprovider defines contract for a subprovider
package subprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Subprovider is the interface implemented by the akamai sub-providers
type Subprovider interface {
	// SDKResources returns the resources implemented using terraform-plugin-sdk
	SDKResources() map[string]*schema.Resource

	// SDKDataSources returns the data sources implemented using terraform-plugin-sdk
	SDKDataSources() map[string]*schema.Resource

	// FrameworkResources returns the resources implemented using terraform-plugin-framework
	FrameworkResources() []func() resource.Resource

	// FrameworkDataSources returns the data sources implemented using terraform-plugin-framework
	FrameworkDataSources() []func() datasource.DataSource
}

// WithActions is implemented by sub-providers that expose
// actions implemented using terraform-plugin-framework. Implementing this interface
// is optional, as most sub-providers don't define any actions.
type WithActions interface {
	Subprovider

	// FrameworkActions returns the actions implemented using terraform-plugin-framework
	FrameworkActions() []func() action.Action
}
