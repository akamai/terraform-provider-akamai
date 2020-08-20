package registry

import (
	"sync"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
)

var (
	lock sync.Mutex

	allProviders []akamai.Subprovider
)

// RegisterProvider simply adds the provider to the array
func RegisterProvider(p akamai.Subprovider) {
	lock.Lock()
	defer lock.Unlock()

	allProviders = append(allProviders, p)
}

// AllProviders returns all of the registered providers
func AllProviders() []akamai.Subprovider {
	lock.Lock()
	defer lock.Unlock()

	out := make([]akamai.Subprovider, len(allProviders))

	copy(out, allProviders)

	return out
}
