package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-akamai/akamai"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: akamai.Provider,
	})
}
