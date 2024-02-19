package gtm

import "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/registry"

func init() {
	registry.RegisterPluginSubprovider(NewSubprovider())
	registry.RegisterFrameworkSubprovider(NewFrameworkSubprovider())
}
