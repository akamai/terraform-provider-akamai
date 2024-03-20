package networklists

import "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/registry"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}
