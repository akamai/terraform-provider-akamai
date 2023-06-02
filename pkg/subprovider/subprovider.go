// Package subprovider defines contract for a subprovider
package subprovider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Subprovider is the interface implemented by the sub providers
type Subprovider interface {
	// Resources returns the resources for the subprovider
	Resources() map[string]*schema.Resource

	// DataSources returns the datasources for the subprovider
	DataSources() map[string]*schema.Resource
}
