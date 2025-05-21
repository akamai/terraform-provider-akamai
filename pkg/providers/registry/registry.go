// Package registry allows to add and register specific sub-providers in Akamai Terraform
package registry

import (
	"sync"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/subprovider"
)

var (
	lock sync.Mutex

	subproviders []subprovider.Subprovider
)

// RegisterSubprovider registers a terraform-plugin-framework sub-provider
func RegisterSubprovider(s subprovider.Subprovider) {
	lock.Lock()
	defer lock.Unlock()

	subproviders = append(subproviders, s)
}

// Subproviders returns all of the registered terraform-plugin-framework sub-providers
func Subproviders() []subprovider.Subprovider {
	lock.Lock()
	defer lock.Unlock()

	out := make([]subprovider.Subprovider, len(subproviders))
	copy(out, subproviders)

	return out
}
