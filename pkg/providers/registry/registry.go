// Package registry allows to add and register specific sub-providers in Akamai Terraform
package registry

import (
	"sync"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/subprovider"
)

var (
	lock sync.Mutex

	pluginSubproviders    []subprovider.Plugin
	frameworkSubproviders []subprovider.Framework
)

// RegisterPluginSubprovider registers a terraform-plugin-sdk sub-provider
func RegisterPluginSubprovider(p subprovider.Plugin) {
	lock.Lock()
	defer lock.Unlock()

	pluginSubproviders = append(pluginSubproviders, p)
}

// PluginSubproviders returns all of the registered terraform-plugin-sdk sub-providers
func PluginSubproviders() []subprovider.Plugin {
	lock.Lock()
	defer lock.Unlock()

	out := make([]subprovider.Plugin, len(pluginSubproviders))
	copy(out, pluginSubproviders)

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
