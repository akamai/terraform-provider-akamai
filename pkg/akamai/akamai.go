// Package akamai allows to initialize and set up Akamai Provider
package akamai

import (
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v6/version"
)

const (
	// ProviderRegistryPath is the path for the provider in the terraform registry
	ProviderRegistryPath = "registry.terraform.io/akamai/akamai"

	// ProviderName is the legacy name of the provider
	// Deprecated: terrform now uses registry paths, the shortest of which would be akamai/akamai"
	ProviderName = "terraform-provider-akamai"
)

func userAgent(terraformVersion string) string {
	return fmt.Sprintf("Terraform/%s (+https://www.terraform.io) %s/%s", terraformVersion,
		ProviderName, version.ProviderVersion)
}
