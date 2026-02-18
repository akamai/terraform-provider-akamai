package mtlskeystore

import "github.com/akamai/terraform-provider-akamai/v10/pkg/providers/registry"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}
