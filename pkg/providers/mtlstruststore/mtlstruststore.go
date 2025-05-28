package mtlstruststore

import "github.com/akamai/terraform-provider-akamai/v8/pkg/providers/registry"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}
