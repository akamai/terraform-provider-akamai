package cps

import "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/registry"

func init() {
	registry.RegisterPluginSubprovider(NewSubprovider())
}
