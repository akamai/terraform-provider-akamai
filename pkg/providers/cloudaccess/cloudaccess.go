package cloudaccess

import "github.com/akamai/terraform-provider-akamai/v12/pkg/providers/registry"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}
