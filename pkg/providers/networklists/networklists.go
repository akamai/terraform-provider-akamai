package networklists

import "github.com/akamai/terraform-provider-akamai/v9/pkg/providers/registry"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}
