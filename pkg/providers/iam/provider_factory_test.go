package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Alias for the function signature expected by the Terraform test driver
type FactoryFunc = func() (*schema.Provider, error)

// A function that can be used for ProviderFactories member in Terraform resource.TestCase tests
func (p *providerOld) Factory() (*schema.Provider, error) {
	return p.ProviderSchema(), nil
}

// Convenience method for the sake of unit testing
func (p *providerOld) ProviderFactories() map[string]FactoryFunc {
	return map[string]FactoryFunc{"akamai": p.Factory}
}
