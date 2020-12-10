package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FactoryFunc = func() (*schema.Provider, error)

// A function that can be used for ProviderFactories member in Terraform resource.TestCase tests
func (p *provider) Factory() (*schema.Provider, error) {
	return p.ProviderSchema(), nil
}

// Convenience method for the sake of unit testing
func (p *provider) ProviderFactories() map[string]FactoryFunc {
	return map[string]FactoryFunc{"akamai": p.Factory}
}
