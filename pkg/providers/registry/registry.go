// Package registry allows to add and register specific sub-providers in Akamai Terraform
package registry

import (
	"sync"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
)

var (
	lock sync.Mutex

	allProviders []subprovider.Subprovider
)

// RegisterProvider simply adds the provider to the array
func RegisterProvider(p subprovider.Subprovider) {
	lock.Lock()
	defer lock.Unlock()

	allProviders = append(allProviders, p)
}

// AllProviders returns all of the registered providers
func AllProviders() []subprovider.Subprovider {
	lock.Lock()
	defer lock.Unlock()

	out := make([]subprovider.Subprovider, len(allProviders))

	copy(out, allProviders)

	return out
}
