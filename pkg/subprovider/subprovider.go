// Package subprovider defines contract for a subprovider
package subprovider

import (
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Subprovider is the interface implemented by the sub providers
type Subprovider interface {
	// Name should return the name of the subprovider
	Name() string

	// Version returns the version of the subprovider
	Version() string

	// Schema returns the schemas for the subprovider
	Schema() map[string]*schema.Schema

	// Resources returns the resources for the subprovider
	Resources() map[string]*schema.Resource

	// DataSources returns the datasources for the subprovider
	DataSources() map[string]*schema.Resource

	// Configure returns the subprovider opaque state object
	Configure(log.Interface, *schema.ResourceData) diag.Diagnostics
}
