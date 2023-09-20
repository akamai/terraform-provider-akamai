// Package registry allows to add and register specific sub-providers in Akamai Terraform
package registry

import (
	"sync"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
)

var (
	lock sync.Mutex

	sdkSubproviders       []subprovider.SDK
	frameworkSubproviders []subprovider.Framework
)

// RegisterSDKSubprovider registers a terraform-plugin-sdk sub-provider
func RegisterSDKSubprovider(p subprovider.SDK) {
	lock.Lock()
	defer lock.Unlock()

	sdkSubproviders = append(sdkSubproviders, p)
}

// SDKSubproviders returns all of the registered terraform-plugin-sdk sub-providers
func SDKSubproviders() []subprovider.SDK {
	lock.Lock()
	defer lock.Unlock()

	out := make([]subprovider.SDK, len(sdkSubproviders))
	copy(out, sdkSubproviders)

	return out
}

// RegisterFrameworkSubprovider registers a terraform-plugin-framework sub-provider
func RegisterFrameworkSubprovider(p subprovider.Framework) {
	lock.Lock()
	defer lock.Unlock()

	frameworkSubproviders = append(frameworkSubproviders, p)
}

// FrameworkSubproviders returns all of the registered terraform-plugin-framework sub-providers
func FrameworkSubproviders() []subprovider.Framework {
	lock.Lock()
	defer lock.Unlock()

	out := make([]subprovider.Framework, len(frameworkSubproviders))
	copy(out, frameworkSubproviders)

	return out
}
